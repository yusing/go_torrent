package main

import "C"
import (
	"log"
	"os"
	"path"

	"github.com/anacrolix/torrent"
)

var torrentClient *torrent.Client = nil
var savePath string
var dataPath string

func loadLastSession() {
	log.Println("[Torrent-Go] Loading session...")
	// torrent list is data/*.json
	if files, err := os.ReadDir(dataPath); err != nil {
		log.Printf("[Torrent-Go] Error reading data dir: %s", err)
	} else {
		for _, file := range files {
			if file.IsDir() {
				continue
			}
			if path.Ext(file.Name()) == ".json" {
				// filename is the infohash
				ReadMetadataAndAdd(file.Name()[:len(file.Name())-5])
			}
		}
		log.Println("[Torrent-Go] Session loaded")
	}
}

//export InitTorrentClient
func InitTorrentClient(savePathCStr *C.char) {
	log.Println("[Torrent-Go] Initializing...")
	if torrentClient != nil {
		return // Already initialized, maybe flutter hot reload?
	}
	savePath = C.GoString(savePathCStr)
	dataPath = path.Join(savePath, "data")
	if err := os.MkdirAll(dataPath, 0644); err != nil {
		log.Printf("[Torrent-Go] Error creating data dir: %s", err)
		os.Exit(1)
	}
	config := torrent.NewDefaultClientConfig()
	config.NoDHT = false
	config.NoUpload = false
	config.DataDir = savePath
	config.DisableIPv6 = true
	torrentClient, _ = torrent.NewClient(config)
	loadLastSession()
}

func main() {
	// This is a dummy main function to make possible
}
