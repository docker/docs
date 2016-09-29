#!/bin/sh -ex

eval `opam config env`
make
# bundler dylib dependencies
mkdir -p _build/root/Contents/MacOS
cp com.docker.osxfs _build/root/Contents/MacOS/com.docker.osxfs
mkdir -p _build/root/Contents/Resources/lib
cp libffi/libffi.6.dylib _build/root/Contents/Resources/lib/libffi.6.dylib
install_name_tool -change /usr/local/opt/libffi/lib/libffi.6.dylib @executable_path/../Resources/lib/libffi.6.dylib _build/root/Contents/MacOS/com.docker.osxfs
