#!/bin/sh

chmod +x clangwrap.sh
echo "Building go_torrent for iOS"
mkdir -p build/ios
go mod init yusing/go_torrent/v2
go mod tidy

export MIN_VERSION=11
export CC=$(PWD)/clangwrap.sh
export CXX=$(PWD)/clangwrap.sh

export SDK=iphoneos
export CGO_CFLAGS='-fembed-bitcode'
export CGO_ENABLED=1
export GOOS=ios
export GOARCH=arm64
. ./target.sh
go build -tags ios -o build/$SDK/go_torrent.a -buildmode=c-archive
$CC -fpic -shared -Wl,-all_load build/$SDK/go_torrent.a -framework Corefoundation -framework Security -lresolv -lstdc++ -o build/$SDK/go_torrent.dylib
lipo -info build/$SDK/go_torrent.dylib

export SDK=iphonesimulator
export GOOS=ios
export GOARCH=amd64
. ./target.sh
go build -tags ios -o build/$SDK/go_torrent_amd64.a -buildmode=c-archive
$CC -fpic -shared -Wl,-all_load build/$SDK/go_torrent_amd64.a -framework Corefoundation -framework Security -lSystem -lresolv -lstdc++ -o build/$SDK/go_torrent_amd64.dylib

export CARCH='arm64'
export GOOS=ios
export GOARCH=arm64
. ./target.sh
go build -tags ios -o build/$SDK/go_torrent_arm64.a -buildmode=c-archive
$CC -fpic -shared -Wl,-all_load build/$SDK/go_torrent_arm64.a -framework Corefoundation -framework Security -lSystem -lresolv -lstdc++ -o build/$SDK/go_torrent_arm64.dylib

lipo build/$SDK/go_torrent_*.dylib -output build/$SDK/go_torrent.dylib -create
lipo -info build/$SDK/go_torrent.dylib