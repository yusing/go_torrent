package main

import (
	"bytes"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type ErrorReporter struct {
}

// Write implements io.Writer.
func (ErrorReporter) Write(p []byte) (n int, err error) {
	req, err := http.NewRequest("POST", "https://err.yumar.org", bytes.NewBuffer(p))
	if err != nil {
		log.Errorf("[Torrent-Go] Error reporting error: %s", err)
		n = -1
		return n, err
	}
	req.Header.Set("Content-Type", "text/plain")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Errorf("[Torrent-Go] Error reporting error: %s (Code: %d)", err, resp.StatusCode)
		n = -1
		return n, err
	}
	n = 0
	return n, nil
}
