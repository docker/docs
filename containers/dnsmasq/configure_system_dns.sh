#!/bin/bash

# Docker Bridge IP - https://docs.docker.com/articles/networking/
# $DOCKER_HOST will be the IP of the boot2docker or docker-machine
# instance *currently sourced in your shell*. This means something
# like $(docker-machine env dev) or $(boot2docker shellinit)
if [[ $DOCKER_HOST =~ ([0-9]{1,3}[\.]){3}[0-9]{1,3} ]]; then
    DAEMON_IPV4=$BASH_REMATCH
    echo $DAEMON_IPV4
else
    echo "unable to parse string $DOCKER_HOST"
fi

set_dev_resolver() {
    echo "Bagels need your permission to configure system DNS."
    sudo mkdir -p /etc/resolver
    echo "nameserver $DAEMON_IPV4" | sudo tee /etc/resolver/bagels.docker.com
}

if [ ! -f /etc/resolver/bagels.docker.com ]; then
    set_dev_resolver
elif [ "$(cat /etc/resolver/bagels.docker.com)" != "nameserver $DAEMON_IPV4" ]; then
    set_dev_resolver
fi
