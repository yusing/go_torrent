#!/bin/sh

echo "Building go_torrent for Linux"
mkdir -p build/linux
go mod init yusing/go_torrent/v2
go mod tidy
export CGO_ENABLED=1

export GOOS=linux
export GOARCH=arm64
export CC=aarch64-linux-gnu-gcc
export CXX=aarch64-linux-gnu-g++
go build -o build/linux/go_torrent_arm64.so -buildmode=c-shared -ldflags="-s -w"

export GOOS=linux
export GOARCH=amd64
export CC=x86_64-linux-gnu-gcc
export CXX=x86_64-linux-gnu-g++
go build -o build/linux/go_torrent.so -buildmode=c-shared -ldflags="-s -w"

file build/linux/go_torrent_arm64.so
file build/linux/go_torrent.so
