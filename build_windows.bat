@echo off
set CGO_ENABLED=1
set GOARCH=amd64

go build -o build\windows\torrent_go.dll -buildmode=c-shared -ldflags="-s -w"

REM set GOARCH=386
REM go build -o build/windows/x86/libtorrent_go.dll -buildmode=c-shared

echo Build finished