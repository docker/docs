#!/bin/bash

set -e

# make sure we have make installed
which make

cd $(dirname $0)
echo If releasing version x.y.z-rca, please enter x y z a. For example, if you were releasing 2.0.0-rc1, you would enter 2 0 0 1.
echo Remember to add '"'config.ssh.forward_agent = true'"' to your Vagrantfile if using vagrant
read x y z a
git clone git@github.com:docker/dhe-deploy.git dhe-deploy-release --branch v$x.$y.$z-rc$a --recursive --depth 1
cd dhe-deploy-release
echo $x.$y.$z-rc$a > version
echo stable > releaseChannel

echo 'Now take the tar using `make tar TAR_UCP_IMAGES=docker/ucp:whatever PULL_UCP_IMAGES=true` and upload it to S3 however you want'
