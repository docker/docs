#!/bin/bash

./docker-machine -D create -d virtualbox --virtualbox-import-boot2docker-vm boot2docker-vm default | sed -e '/BEGIN/,/END/d'
