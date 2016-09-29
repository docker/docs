#!/bin/sh
set -e

# These are binaries from a mirror of
#  http://tinycorelinux.net
# with the following patch applied:
# Upstream source is available http://www.tinycorelinux.net/6.x/x86/release/src/
#BASE_URL="http://www.tinycorelinux.net/"

BASE_URL="http://distro.ibiblio.org/tinycorelinux/"

TMP_DIR=$(mktemp -d -t hyperkit)
INITRD_DIR="${TMP_DIR}"/initrd

echo Downloading tinycore linux
curl -s -o vmlinuz "${BASE_URL}/6.x/x86/release/distribution_files/vmlinuz64"
curl -s -o "${TMP_DIR}"/initrd.gz "${BASE_URL}/6.x/x86/release/distribution_files/core.gz"

mkdir "${INITRD_DIR}"
( cd "${INITRD_DIR}"; gzip -dc "${TMP_DIR}"/initrd.gz | sudo cpio -idm )
sudo sed -i -e '/^# ttyS0$/s#^..##' "${INITRD_DIR}"/etc/securetty 
sudo sed -i -e '/^tty1:/s#tty1#ttyS0#g' "${INITRD_DIR}"/etc/inittab
( cd "${INITRD_DIR}" ; find . | sudo cpio -o -H newc ) | gzip -c > initrd.gz && sudo rm -rf "${TMP_DIR}"
