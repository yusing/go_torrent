package main

import "C"
import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
	"time"
	"unsafe"

	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/metainfo"
	"github.com/anacrolix/torrent/types/infohash"

	log "github.com/sirupsen/logrus"
)

func GetTrackers() []string {
	URL := "https://github.com/XIU2/TrackersListCollection/raw/master/best.txt"
	trackersTxtPath := path.Join(dataPath, "trackers.txt")
	// if trackers.txt is older than 1 day, download it again
	if fileInfo, err := os.Stat(trackersTxtPath); err != nil || time.Since(fileInfo.ModTime()).Hours() > 24 {
		// download trackers.txt
		resp, err := http.Get(URL)
		if err != nil {
			log.Debugf("[Torrent-Go] Error downloading trackers.txt: %s", err)
			return nil
		}
		defer resp.Body.Close()
		trackersTxtFile, err := os.OpenFile(trackersTxtPath, os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			log.Debugf("[Torrent-Go] Error creating trackers.txt: %s", err)
			return nil
		}
		defer trackersTxtFile.Close()
		_, err = io.Copy(trackersTxtFile, resp.Body)
		if err != nil {
			log.Debugf("[Torrent-Go] Error writing trackers.txt: %s", err)
			return nil
		}
	}

	if file, err := os.OpenFile(trackersTxtPath, os.O_RDONLY, 0644); err != nil {
		log.Debugf("[Torrent-Go] Error reading trackers.txt: %s", err)
		return nil
	} else {
		defer file.Close()
		content := make([]byte, 8192)
		length, err := file.Read(content)
		if err != nil {
			log.Debugf("[Torrent-Go] Error reading trackers.txt: %s", err)
			return nil
		}
		lines := strings.Split(string(content[:length-1]), "\n")
		trackers := []string{}
		for _, line := range lines {
			// if line is empty, skip
			if line == "" {
				continue
			}
			trackers = append(trackers, line)
		}
		return trackers
	}
}

func AddTorrentFromInfoHash(infoHashStr string) *torrent.Torrent {
	infoHash_ := infohash.T{}
	if parseErr := infoHash_.FromHexString(infoHashStr); parseErr != nil {
		log.Debugf("[Torrent-Go] Error parsing infoHash: %s %s", infoHashStr, parseErr)
		return nil
	} else {
		t, ok := torrentClient.AddTorrentInfoHash(infoHash_)
		if t == nil || !ok {
			log.Debugf("[Torrent-Go] Error adding torrent from infoHash %s", infoHashStr)
		}
		<-t.GotInfo()
		t.DownloadAll()
		SaveMetadata(t)
		return t
	}
}

func SaveMetadata(t *torrent.Torrent) {
	metaInfoJson, err := json.Marshal(t.Metainfo())
	if err != nil {
		log.Debugf("[Torrent-Go] Error saving metadata: %s", err)
		return
	}
	file, err := os.OpenFile(path.Join(dataPath, t.InfoHash().HexString()+".json"), os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Debugf("[Torrent-Go] Error saving metadata: %s", err)
		return
	}
	defer file.Close()
	file.Write(metaInfoJson)
}

func ReadMetadataAndAdd(infoHashStr string) *torrent.Torrent {
	// fallback to AddTorrentFromInfoHash if metadata not found
	metaInfoJsonBytes, err := os.ReadFile(path.Join(dataPath, infoHashStr+".json"))
	if err != nil {
		log.Debugf("[Torrent-Go] Error reading metadata: %s", err)
		AddTorrentFromInfoHash(infoHashStr)
		return nil
	}
	metaInfo := metainfo.MetaInfo{}
	if metaInfoParseErr := json.Unmarshal(metaInfoJsonBytes, &metaInfo); metaInfoParseErr != nil {
		log.Debugf("[Torrent-Go] Error parsing metaInfo: %s", metaInfoParseErr)
		return AddTorrentFromInfoHash(infoHashStr)
	} else {
		t, _ := torrentClient.AddTorrent(&metaInfo)
		<-t.GotInfo()
		t.DownloadAll()
		return t
	}
}

//export AddMagnet
func AddMagnet(magnetCString *C.char) *C.char {
	magnet := C.GoString(magnetCString)
	log.Debugf("Adding torrent: %s", magnet)
	t, err := torrentClient.AddMagnet(magnet)
	if t == nil || err != nil {
		return jsonify([]map[string]interface{}{})
	}
	<-t.GotInfo()
	// log.Debugf("Added %p", t)
	if trackers := GetTrackers(); trackers != nil {
		t.AddTrackers([][]string{trackers})
	}
	SaveMetadata(t)
	t.DownloadAll()
	torrentInfoMap := torrentInfoMap(t)
	return jsonify(torrentInfoMap)
}

//export AddTorrent
func AddTorrent(torrentUrlCStr *C.char) *C.char {
	torrentUrl := C.GoString(torrentUrlCStr)
	torrentPath := path.Join(dataPath, path.Base(torrentUrl))
	torrentFile, err := os.Create(torrentPath)
	if err != nil {
		log.Debugf("[Torrent-Go] Error creating torrent file: %s", err)
		return jsonify([]map[string]interface{}{})
	}
	defer torrentFile.Close()
	resp, err := http.Get(torrentUrl)
	if err != nil {
		log.Debugf("[Torrent-Go] Error downloading torrent file: %s", err)
		return jsonify([]map[string]interface{}{})
	}
	defer resp.Body.Close()
	_, err = io.Copy(torrentFile, resp.Body)
	if err != nil {
		log.Debugf("[Torrent-Go] Error reading torrent file: %s", err)
		return jsonify([]map[string]interface{}{})
	}
	t, err := torrentClient.AddTorrentFromFile(torrentPath)
	os.Remove(torrentPath)
	if t == nil || err != nil {
		log.Debugf("[Torrent-Go] Error adding torrent: %s", err)
		return jsonify([]map[string]interface{}{})
	}
	<-t.GotInfo()
	if trackers := GetTrackers(); trackers != nil {
		t.AddTrackers([][]string{trackers})
	}
	torrentInfoMap := torrentInfoMap(t)
	t.DownloadAll()
	SaveMetadata(t)
	return jsonify(torrentInfoMap)
}

//export PauseTorrent
func PauseTorrent(torrentPtr unsafe.Pointer) {
	if torrentPtr == nil {
		return
	}
	t := (*torrent.Torrent)(torrentPtr)
	SaveMetadata(t)
	t.Drop()
}

//export ResumeTorrent
func ResumeTorrent(infoHashCStr *C.char) uintptr {
	t := ReadMetadataAndAdd(C.GoString(infoHashCStr))
	return uintptr(unsafe.Pointer(t))
}

//export DeleteTorrent
func DeleteTorrent(torrentPtr unsafe.Pointer) {
	if torrentPtr == nil {
		return
	}
	t := (*torrent.Torrent)(torrentPtr)
	// remove files/directory
	infoHashStr := t.InfoHash().HexString()
	t.Drop()
	var err error
	if t.Info().IsDir() {
		err = os.RemoveAll(path.Join(savePath, t.Name()))
	} else {
		err = os.Remove(path.Join(savePath, t.Name()))
	}
	if err != nil {
		log.Debugf("[Torrent-Go] Warning: Error deleting files: %s", err)
	}
	if (os.Remove(path.Join(dataPath, infoHashStr+".json"))) != nil {
		log.Debugf("[Torrent-Go] Warning: Error deleting metadata: %s", err)
	}
}
