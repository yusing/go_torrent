@echo off
set CGO_ENABLED=1
set GOOS=android
set CGO_CFLAGS="--sysroot=%NDK%/toolchains/llvm/prebuilt/windows-x86_64/sysroot"


set GOARCH=arm64
set CC="%NDK%/toolchains/llvm/prebuilt/windows-x86_64/bin/aarch64-linux-android23-clang"
set CXX="%NDK%/toolchains/llvm/prebuilt/windows-x86_64/bin/aarch64-linux-android23-clang++"
go build -o build/android/arm64-v8a/libtorrent_go.so -buildmode=c-shared
xcopy %NDK%\sources\cxx-stl\llvm-libc++\libs\arm64-v8a\libc++_shared.so build\android\arm64-v8a\ /Q /Y

set GOARCH=arm
set GOARM=7
set CC="%NDK%/toolchains/llvm/prebuilt/windows-x86_64/bin/armv7a-linux-androideabi23-clang"
set CXX="%NDK%/toolchains/llvm/prebuilt/windows-x86_64/bin/armv7a-linux-androideabi23-clang++"
go build -o build/android/armeabi-v7a/libtorrent_go.so -buildmode=c-shared
xcopy %NDK%\sources\cxx-stl\llvm-libc++\libs\armeabi-v7a\libc++_shared.so build\android\armeabi-v7a\ /Q /Y

set GOARCH=386
set GOARM=
set CC="%NDK%/toolchains/llvm/prebuilt/windows-x86_64/bin/i686-linux-android23-clang"
set CXX="%NDK%/toolchains/llvm/prebuilt/windows-x86_64/bin/i686-linux-android23-clang++"
go build -o build/android/x86/libtorrent_go.so -buildmode=c-shared
xcopy %NDK%\sources\cxx-stl\llvm-libc++\libs\x86\libc++_shared.so build\android\x86\ /Q /Y

set GOARCH=amd64
set CC="%NDK%/toolchains/llvm/prebuilt/windows-x86_64/bin/x86_64-linux-android23-clang"
set CXX="%NDK%/toolchains/llvm/prebuilt/windows-x86_64/bin/x86_64-linux-android23-clang++"
go build -o build/android/x86_64/libtorrent_go. -buildmode=c-shared
xcopy %NDK%\sources\cxx-stl\llvm-libc++\libs\x86_64\libc++_shared.so build\android\x86_64\ /Q /Y

echo Build finished