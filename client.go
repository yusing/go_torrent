package main

import "C"
import (
	"os"
	"path"
	"runtime"

	"github.com/anacrolix/torrent"

	log "github.com/sirupsen/logrus"
)

var torrentClient *torrent.Client = nil
var savePath string
var dataPath string

const IS_MOBILE = runtime.GOOS == "android" || runtime.GOOS == "ios"

func loadLastSession() {
	log.Debugln("[Torrent-Go] Loading session...")
	// torrent list is data/*.json
	if files, err := os.ReadDir(dataPath); err != nil {
		log.Debugf("[Torrent-Go] Error reading data dir: %s", err)
	} else {
		log.Debugf("[Torrent-Go] Found %d files", len(files))
		for _, file := range files {
			if file.IsDir() {
				continue
			}
			log.Debugf("[Torrent-Go] Found file %s", file.Name())
			if path.Ext(file.Name()) == ".json" {
				// filename is the infohash
				ReadMetadataAndAdd(file.Name()[:len(file.Name())-5])
			}
		}
		log.Debugln("[Torrent-Go] Session loaded")
	}
}

//export InitTorrentClient
func InitTorrentClient(savePathCStr *C.char) {
	// fix logcat android
	if runtime.GOOS == "android" {
		addAndroidLogHook()
	}
	log.Debugln("[Torrent-Go] Initializing...")
	if torrentClient != nil {
		return // Already initialized, maybe flutter hot reload?
	}
	savePath = C.GoString(savePathCStr)
	if IS_MOBILE {
		dataPath = savePath
	} else {
		dataPath = path.Join(savePath, "data")
	}
	config := torrent.NewDefaultClientConfig()
	config.NoDHT = false
	config.NoUpload = IS_MOBILE
	config.DataDir = savePath
	config.DisableIPv6 = true
	torrentClient, _ = torrent.NewClient(config)
	loadLastSession()
}

func main() {
	// This is a dummy main function to make possible
}
