#!/bin/bash
ISO=$HOME/.docker/machine/cache/boot2docker.iso
VM=dev
DOCKER_MACHINE=./docker-machine.exe

mkdir -p ~/.docker/machine/cache
if [ ! -f $ISO ]; then
	mkdir -p "$(dirname "$ISO")"
	cp ./boot2docker.iso "$ISO"
fi

machine=$($DOCKER_MACHINE ls -q | grep "^$VM$")
if [ -z $machine ]; then
  echo "Creating Machine $VM..."
   $DOCKER_MACHINE create -d virtualbox --virtualbox-memory 2048 $VM
else
  echo "Machine $VM already exists."
fi

echo "Starting machine $VM..."
$DOCKER_MACHINE start dev

echo "Setting environment variables for machine $VM..."
./docker-machine.exe env $VM | sed  's,\\,\\\\,g' # eval swallows single backslashes in windows style path

eval "$(./docker-machine.exe env $VM 2>/dev/null | sed  's,\\,\\\\,g')"

cd
clear

cat << EOF


                        ##         .
                  ## ## ##        ==
               ## ## ## ## ##    ===
           /"""""""""""""""""\___/ ===
      ~~~ {~~ ~~~~ ~~~ ~~~~ ~~~ ~ /  ===- ~~~
           \______ o           __/
             \    \         __/
              \____\_______/

EOF
echo "The Quick Start CLI is configured to use Docker with the $VM VM"
echo

exec "$BASH" --login -i
