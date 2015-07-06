#!/bin/bash

# Uninstall Script

if [ "$(which boot2docker)" == "" ]; then
	echo "boot2docker does not exist on your machine!"
	exit 1
fi

if [ "${USER}" != "root" ]; then
	echo "$0 must be run as root!"
	exit 2
fi

echo "Stopping boot2docker processes..."
boot2docker stop && boot2docker delete

echo "Removing boot2docker executable..."
rm -f /usr/local/bin/boot2docker

echo "Removing boot2docker ISO and socket files..."
rm -rf ~/.boot2docker
rm -rf /usr/local/share/boot2docker

echo "Removing boot2docker SSH keys..."
rm -f ~/.ssh/id_boot2docker*

echo "Removing boot2docker OSX files..."
rm -f /private/var/db/receipts/io.boot2docker.*
rm -f /private/var/db/receipts/io.boot2dockeriso.*

echo "Removing Docker executable..."
rm -f /usr/local/bin/docker

echo "All Done!"
