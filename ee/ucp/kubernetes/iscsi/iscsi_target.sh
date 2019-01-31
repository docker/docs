#!/bin/bash

# Main customization needed for this script is
# setting the iscsi target_name. Default is
# "iqn.2019-01.org.iscsi.docker:targetd". Feel free to change it,
# as long as it follows IQN naming rules.

# Run this script as root
if [[ $EUID -ne 0 ]]; then
   echo "This script must be run as root"
   exit 1
fi

# Setup volume group on Loopback device.
mkdir /var/lib/loopback
cd /var/lib/loopback
dd if=/dev/zero of=disk.img bs=1G count=2

export LOOP=`sudo losetup -f`
losetup $LOOP disk.img
vgcreate vg-targetd $LOOP


# Install targetd and targetcli
yum install -y targetcli targetd

# Enable targetcli
systemctl enable target
systemctl start target

# Configure targetd

echo "password: ciao

# defaults below; uncomment and edit
# if using a thin pool, use <volume group name>/<thin pool name>
# e.g vg-targetd/pool
pool_name: vg-targetd
user: admin
ssl: false
target_name: iqn.2019-01.org.iscsi.docker:targetd" > /etc/target/targetd.yaml

# Enable targetd
systemctl enable targetd
systemctl start targetd
