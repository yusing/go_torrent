package main

import "C"
import (
	"unsafe"

	"github.com/anacrolix/torrent"
	log "github.com/sirupsen/logrus"
)

// func FindTorrent(infoHashStr string) (t *torrent.Torrent, infoHash metainfo.Hash) {
// 	if err := infoHash.FromHexString(infoHashStr); err != nil {
// 		return nil, metainfo.Hash{}
// 	}
// 	t, ok := torrentClient.Torrent(infoHash)
// 	if !ok || t == nil {
// 		return nil, infoHash
// 	}
// 	return t, infoHash
// }

//export GetTorrentInfo
func GetTorrentInfo(torrentPtr unsafe.Pointer) *C.char {
	if torrentPtr == nil {
		log.Debugln("[Torrent-Go] GetTorrentInfo: torrentPtr is nil")
		return jsonify(map[string]interface{}{})
	}
	return jsonify(torrentInfoMap((*torrent.Torrent)(torrentPtr)))
}

func TorrentList() []map[string]interface{} {
	if torrentClient == nil {
		return []map[string]interface{}{}
	}
	var torrents []map[string]interface{}
	for _, t := range torrentClient.Torrents() {
		if t == nil {
			continue
		}
		torrents = append(torrents, torrentInfoMap(t))
	}
	return torrents
}

//export GetTorrentList
func GetTorrentList() *C.char {
	return jsonify(TorrentList())
}
