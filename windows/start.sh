#!/bin/bash

set -e

ISO=$HOME/.docker/machine/cache/boot2docker.iso
VM=dev
DOCKER_MACHINE=./docker-machine.exe

BLUE='\033[1;34m'
GREEN='\033[0;32m'
NC='\033[0m'

mkdir -p ~/.docker/machine/cache
if [ ! -f $ISO ]; then
	mkdir -p "$(dirname "$ISO")"
	cp ./boot2docker.iso "$ISO"
fi

machine=$($DOCKER_MACHINE ls -q | grep "^$VM$") || :
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

eval "$($DOCKER_MACHINE env $VM 2>/dev/null | sed  's,\\,\\\\,g')"

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
echo -e "${BLUE}docker${NC} is configured to use the ${GREEN}dev${NC} machine with IP ${GREEN}$($DOCKER_MACHINE ip $VM)${NC}"
echo
cd
