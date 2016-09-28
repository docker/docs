#!/bin/bash

set -e

# make sure we have make installed
which make

cd $(dirname $0)
echo If releasing version x.y.z, please enter x y z. For example, if you were releasing 2.0, you would enter 2 0 0.
echo Remember to add '"'config.ssh.forward_agent = true'"' to your Vagrantfile if using vagrant
read x y z
git clone git@github.com:docker/dhe-deploy.git dhe-deploy-release --branch v$x.$y.$z --recursive --depth 1
cd dhe-deploy-release
echo $x.$y.$z > version
echo stable > releaseChannel

echo "Type PUSH (all caps) to upload all non-latest images for DTR $x.$y.$z"
read confirm
if [ "$confirm" != 'PUSH' ]; then
    echo $confirm != PUSH\; aborting.
    exit 1
fi

make push PUSH_NO_LATEST=true

echo "Type PUSH (all caps) to retag latest for DTR $x.$y.$z"
read confirm
if [ "$confirm" != 'PUSH' ]; then
    echo $confirm != PUSH\; aborting.
    exit 1
fi

make push

echo 'Now take the tar using `make tar TAR_UCP_IMAGES=docker/ucp:whatever PULL_UCP_IMAGES=true` and upload it to S3 however you want'
