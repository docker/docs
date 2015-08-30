#!/bin/bash

set -o pipefail
./docker-machine -D create -d virtualbox --virtualbox-memory 2048 --virtualbox-import-boot2docker-vm boot2docker-vm default | sed -e '/BEGIN/,/END/d'
