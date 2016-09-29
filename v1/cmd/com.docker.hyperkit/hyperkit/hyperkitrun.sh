#!/bin/sh

HYPERKIT="build/com.docker.hyperkit"

# Linux
KERNEL="test/vmlinuz"
INITRD="test/initrd.gz"
CMDLINE="earlyprintk=serial console=ttyS0"

# FreeBSD
#USERBOOT="test/userboot.so"
#BOOTVOLUME="/somepath/somefile.{img | iso}"
#KERNELENV=""

MEM="-m 1G"
#SMP="-c 2"
#NET="-s 2:0,virtio-net"
#IMG_CD="-s 3,ahci-cd,/somepath/somefile.iso"
#IMG_HDD="-s 4,virtio-blk,/somepath/somefile.img"
PCI_DEV="-s 0:0,hostbridge -s 31,lpc"
LPC_DEV="-l com1,stdio"
ACPI="-A"
#UUID="-U deadbeef-dead-dead-dead-deaddeafbeef"

if [ ! -x  "$HYPERKIT" ]; then
  make
  if [ ! $? -eq 0 ]; then
    echo "Error whilst building $HYPERKIT"
    exit 1
  fi
fi

# Linux
if [ ! -f "$KERNEL" ]; then
 pushd test
 ./tinycore.sh
 popd
fi
build/com.docker.hyperkit $ACPI $MEM $SMP $PCI_DEV $LPC_DEV $NET $IMG_CD $IMG_HDD $UUID -f kexec,$KERNEL,$INITRD,"$CMDLINE"

# FreeBSD
#build/com.docker.hyperkit $ACPI $MEM $SMP $PCI_DEV $LPC_DEV $NET $IMG_CD $IMG_HDD $UUID -f fbsd,$USERBOOT,$BOOTVOLUME,"$KERNELENV"
