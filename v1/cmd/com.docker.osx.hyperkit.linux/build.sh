#!/bin/sh -ex

make
# bundler dylib dependencies
mkdir -p _build/root/Contents/MacOS
cp com.docker.osx.hyperkit.linux _build/root/Contents/MacOS/com.docker.osx.hyperkit.linux
dylibbundler -od -b \
  -x _build/root/Contents/MacOS/com.docker.osx.hyperkit.linux \
  -d _build/root/Contents/Resources/lib \
  -p @executable_path/../Resources/lib
