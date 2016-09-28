#!/bin/sh
set -x

# we need to not have a tty in order to print newlines normally instead of with
# \r\n at the end

docker run -i --rm \
  -v $(pwd):$(pwd) \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -v ${HOME}/.docker:/root/.docker \
  -w $(pwd) \
  dockerhubenterprise/doit:withscript sh -c "$@"
