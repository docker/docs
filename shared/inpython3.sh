#!/bin/sh
set -x

# we have to trick doit into  thinking that it's running over a tty for it to
# not buffer output

docker run -i --rm \
  -v $(pwd):$(pwd) \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -v ${HOME}/.docker:/root/.docker \
  -w $(pwd) \
  dockerhubenterprise/doit:withscript script --return -qfc "$@" /dev/null
