#!/bin/bash

# DMG Creation Script
# Usage: makedmg <imagename> <imagetitle> <contentdir>
#
# Based on makedmg by Jon Cowie
#
# imagename: The output file name of the image, ie foo.dmg
# imagetitle: The title of the DMG File as displayed in OS X
# contentdir: The directory containing the content you want the DMG file to contain

if [ ! $# == 3 ]; then
    echo "Usage: $0 <imagename> <imagetitle> <contentdir>"
else
    OUTPUT=$1
    TITLE=$2
    CONTENTDIR=$3
    FILESIZE=$(du -sm "${CONTENTDIR}" | cut -f1)
    FILESIZE=$((${FILESIZE} + 5))
    USER=$(whoami)
    TMPDIR="/tmp/dmgdir"

    if [ "${USER}" != "root" ]; then
        echo "$0 must be run as root!"
    else
        echo "Creating DMG File..."
        dd if=/dev/zero of="${OUTPUT}" bs=1M count=$FILESIZE
        mkfs.hfsplus -v "${TITLE}" "${OUTPUT}"

        echo "Mounting DMG File..."
        mkdir -p ${TMPDIR}
        mount -t hfsplus -o loop "${OUTPUT}" "${TMPDIR}"

        echo "Copying content to DMG File..."
        cp -R "${CONTENTDIR}"/* "${TMPDIR}"

        echo "Unmounting DMG File..."
        umount "${TMPDIR}"
        rm -rf "${TMPDIR}"

        echo "All Done!"
    fi
fi
