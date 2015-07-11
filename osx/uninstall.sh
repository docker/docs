#!/bin/bash

# Uninstall Script

if [ "${USER}" != "root" ]; then
	echo "$0 must be run as root!"
	exit 2
fi

echo "Removing dev VirtualBox VM..."
docker-machine rm -f $(docker-machine ls -q)

echo "Removing docker binaries..."
rm -f /usr/local/bin/docker
rm -f /usr/local/bin/docker-machine
rm -f /usr/local/bin/docker-compose

echo "Removing boot2docker.iso and socket files..."
rm -rf ~/.docker
rm -rf /usr/local/share/boot2docker

echo "Removing boot2docker OSX files..."
rm -f /private/var/db/receipts/io.boot2docker.*
rm -f /private/var/db/receipts/io.boot2dockeriso.*

echo "All Done!"
