#!/bin/sh
SDK_PATH=`xcrun --sdk $SDK --show-sdk-path`
CLANG=`xcrun --sdk $SDK --find clang`
exec $CLANG -arch $CARCH -isysroot $SDK_PATH -mios-version-min=10.0 "$@"