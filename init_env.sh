#!/bin/sh
export SDK_PATH=`xcrun --sdk $SDK --show-sdk-path`
export CLANG=`xcrun --sdk $SDK --find clang`
export CGO_ENABLED=1
export GOOS=ios
export GOARCH=arm64
export SDK=iphoneos
export SDK_PATH=`xcrun --sdk $SDK --show-sdk-path`
export CLANG=`xcrun --sdk $SDK --find clang`
export CC=$(PWD)/clangwrap.sh
export CARCH='arm64'
export CGO_CFLAGS='-fembed-bitcode'