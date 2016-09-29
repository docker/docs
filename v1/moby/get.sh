#!/bin/sh

set -e 

rm -rf tmp
mkdir -p tmp

regextract mobylinux/media:$(cat COMMIT) | tar xf - -C tmp

cp tmp/initrd.img .
cp tmp/mobylinux-efi.iso ../../win/src/mobylinux.iso
cp tmp/vmlinuz64 .

rm -rf tmp
