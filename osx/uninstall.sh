#!/bin/bash

# Uninstall Script

if [ "${USER}" != "root" ]; then
	echo "$0 must be run as root!"
	exit 2
fi

while true; do
  read -p "Remove all Docker Machine VMs? (Y/N): " yn
  case $yn in
    [Yy]* ) docker-machine rm -f $(docker-machine ls -q); break;;
    [Nn]* ) break;;
    * ) echo "Please answer yes or no.";;
  esac
done

echo "Removing Applications..."
rm -rf /Applications/Docker

echo "Removing docker binaries..."
rm -f /usr/local/bin/docker
rm -f /usr/local/bin/docker-machine
rm -f /usr/local/bin/docker-compose

echo "Removing boot2docker.iso"
rm -rf /usr/local/share/boot2docker

echo "All Done!"
