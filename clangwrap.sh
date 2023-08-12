#!/bin/sh
exec $CLANG -arch $CARCH -isysroot $SDK_PATH -mios-version-min=10.0 "$@"