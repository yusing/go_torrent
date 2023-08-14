#!/bin/sh

chmod +x clangwrap.sh
echo "Building go_torrent for iOS"
mkdir -p build/ios
go mod init yusing/go_torrent/v2
go mod tidy
SDK=iphoneos \
SDK_PATH=`xcrun --sdk $SDK --show-sdk-path` \
CLANG=`xcrun --sdk $SDK --find clang` \
CC=$(PWD)/clangwrap.sh \
CXX=$(PWD)/clangwrap.sh \
CARCH='arm64' \
CGO_CFLAGS='-fembed-bitcode' \
CGO_ENABLED=1 \
GOOS=ios \
GOARCH=arm64 \
go build -tags ios -o build/ios/libtorrent_go_ios.a -buildmode=c-archive
echo "Building go_torrent for iOS completed"
file build/ios/libtorrent_go_ios.a