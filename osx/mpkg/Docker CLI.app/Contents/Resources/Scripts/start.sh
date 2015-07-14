#!/bin/bash

set -e

ISO=$HOME/.docker/machine/cache/boot2docker.iso
VM=dev
DOCKER_MACHINE=/usr/local/bin/docker-machine

BLUE='\033[0;34m'
GREEN='\033[0;32m'
NC='\033[0m'

unset DYLD_LIBRARY_PATH
unset LD_LIBRARY_PATH

clear

mkdir -p ~/.docker/machine/cache
if [ ! -f $ISO ]; then
  cp /usr/local/share/boot2docker/boot2docker.iso $ISO
fi


machine=$($DOCKER_MACHINE ls -q | grep "^$VM$") || :
if [ -z $machine ]; then
  echo "Creating Machine $VM..."
   $DOCKER_MACHINE create -d virtualbox --virtualbox-memory 2048 $VM
else
  echo "Machine $VM already exists."
fi

echo "Starting machine $VM..."
$DOCKER_MACHINE start $VM

echo "Setting environment variables for machine $VM..."
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
echo -e "${BLUE}docker${NC} is configured to use the ${GREEN}dev${NC} machine with IP ${GREEN}$($DOCKER_MACHINE ip dev)${NC}"
echo

eval $($DOCKER_MACHINE env $VM --shell=bash)
