#!/bin/bash

set -e

if [ $# -ne 1 ] ; then
    echo "usage: update-vendor.sh <TAG>" 2>&1
    exit 1
fi

tag=$1

[ -d docker.git ] || git clone https://github.com/docker/docker.git docker.git

pushd docker.git >/dev/null
git checkout "$tag"
popd >/dev/null

# First the docker/docker vendoring
rm -rf vendor
cp -r docker.git/vendor/src/ vendor/

# Then the sub-directories of docker/docker itself. Add to the list as
# new dependencies arise.
mkdir -p vendor/github.com/docker/docker
mkdir -p vendor/github.com/docker/docker/pkg
for p in runconfig api pkg/ioutils pkg/broadcaster pkg/longpath pkg/random pkg/stringid pkg/system volume ; do
    cp -r ./docker.git/$p vendor/github.com/docker/docker/$p
    find vendor/github.com/docker/docker/$p -type f \( -name \*_test.go -o \! -name \*.go \) -exec rm {} \;
done

git add vendor
