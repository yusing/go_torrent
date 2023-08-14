#!/bin/sh
CLANG=$(xcrun --sdk "$SDK" --find clang)
exec $CLANG -target "$TARGET" -isysroot $SDK_PATH -mios-version-min=11.0 "$@"