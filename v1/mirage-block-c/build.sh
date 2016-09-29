#!/bin/sh
set -ex

eval `opam config env`

mkdir -p _build
cd _build

echo Building data.qcow2
qcow-tool create data.qcow2 --size 64GiB
qcow-tool write data.qcow2 --text "boot2docker, please format-me"
cp `which qcow-tool` .

mkdir -p root/Contents/Resources/moby
cp data.qcow2 root/Contents/Resources/moby/
mkdir -p root/Contents/MacOS
cp qcow-tool root/Contents/MacOS/
