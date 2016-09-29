#!/bin/bash

set -e

repo="git@github.com:docker/hyperkit"

usage() {
    echo "usage: update.sh [branch]"
    echo ""
    echo "Default branch: $repo master"
    echo ""
    echo "Alternative branches:"
    echo " - None"
}

# git subtree requires that it is run at the top of the working tree.
if [ -f update.sh -a -d hyperkit -a -d libvmnetd -a hyperkit/src/xhyve.c ] ; then
    cd ../../../
fi

case ${1:-master} in
    master)
	git subtree pull -P v1/cmd/com.docker.hyperkit/hyperkit git@github.com:docker/hyperkit master
	;;
    --help|help)
	usage
	exit 0
	;;
    *)
	echo "Unknown branch: %1"
	usage
	exit 1
	;;
esac
