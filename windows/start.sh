#!/bin/bash
set -e

# clear the MSYS MOTD
clear

cd "$(dirname "$BASH_SOURCE")"

ISO="$HOME/.boot2docker/boot2docker.iso"

if [ ! -e "$ISO" ]; then
	echo 'copying initial boot2docker.iso (run "boot2docker.exe download" to update)'
	mkdir -p "$(dirname "$ISO")"
	cp ./boot2docker.iso "$ISO"
fi

echo 'initializing...'
./boot2docker.exe init
echo

echo 'starting...'
./boot2docker.exe start
echo

echo 'IP address of docker VM:'
./boot2docker.exe ip
echo

echo 'setting environment variables ...'
./boot2docker.exe shellinit | sed  's,\\,\\\\,g' # eval swallows single backslashes in windows style path
eval "$(./boot2docker.exe shellinit 2>/dev/null | sed  's,\\,\\\\,g')"
echo

echo 'You can now use `docker` directly, or `boot2docker ssh` to log into the VM.'

cd
exec "$BASH" --login -i
