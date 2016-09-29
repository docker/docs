#!/bin/sh
# SUMMARY: Test that proxy settings appear inside the container
# LABELS: skip
# AUTHOR: Dave Tucker <dt@docker.com>, Magnus Skjegstad <magnus.skjegstad@docker.com>

# Source libraries. Uncomment if needed/defined
#. ${RT_ROOT}/lib/lib.sh
. "${RT_PROJECT_ROOT}/_lib/lib.sh"

set -e

adapter=""
IMAGE_NAME=test_backend_network_proxy
CONTAINER_NAME=${IMAGE_NAME}_container

cleanup () {
    echo "Cleaning up: removing images..."
    docker ps -a | grep $IMAGE_NAME | cut -f1 -d" " | xargs docker kill || true
    docker ps -a | grep $IMAGE_NAME | cut -f1 -d" " | xargs docker rm || true
    docker rmi $IMAGE_NAME || true
    docker rmi hello-world
    rm -f out.tmp
}

get_first_active_if() {
    networksetup -listnetworkserviceorder | grep 'Hardware Port' | \
    while read -r line; do
        sname=$(echo "$line" | awk -F  "(, )|(: )|[)]" '{print $2}')
        sdev=$(echo "$line" | awk -F  "(, )|(: )|[)]" '{print $4}')
        if [ -n "$sdev" ]; then
            ifconfig "$sdev" 2>/dev/null | grep 'status: active' > /dev/null 2>&1 || continue
            echo "$sname"
        fi
    done
}

get_adapter() {
    adapter=$(get_first_active_if)
    if [ -z "${adapter}" ]; then
        echo "Unable to find active network interface"
        exit 1
    fi
}

setup() {
    docker build -t ${IMAGE_NAME} .
    docker run --restart always --name ${CONTAINER_NAME} -d -p 3128:3128 ${IMAGE_NAME}
    sleep 5
    docker ps | grep ${CONTAINER_NAME}
    get_adapter
    sudo networksetup -setwebproxy "${adapter}" "192.168.65.2" 3128
    sudo networksetup -setwebproxystate "${adapter}" on
    sleep 20
}

run() {
    # Test that Docker pull is working
    docker pull hello-world || (echo "unable to pull hello-world image" && exit 1)
    docker rmi hello-world
    
    # Check the env variables are present
    docker run -it --rm alpine env > out.tmp || (echo "unable to get alpine env" && exit 1)
    grep -i 'http' out.tmp
    grep 'http_proxy' out.tmp || (echo "http_proxy not set in container" && exit 1)

    # Check internet access from a container
    docker run -it --rm alpine wget -O- http://www.docker.com > out.tmp || ( echo "unable to run wget in alpine shell, got:" && cat out.tmp && exit 1)
    grep '</html>' out.tmp || ( echo "---- unexpected response: ----" && cat out.tmp && exit 1)

    # Check that we did actually use the proxy for the last command
    docker exec -it ${CONTAINER_NAME}  grep 'www.docker.com' /var/log/squid/access.log || (echo "proxy access request not found in squid logs" && exit 1)

    # Check that when we set and proxy env variable, it overrides
    docker run -it --rm -e HTTP_PROXY=1.1.1.1 alpine env > out.tmp || (echo "unable to pull alpine env" && exit 1)
    grep 'HTTP_PROXY=1.1.1.1' out.tmp || (echo "alpine env did not contain the correct HTTP_PROXY variable" && exit 1)
    grep -v "http_proxy=192.168.65.2" out.tmp
    
    # Unset proxy settings
    if [ ! -z "${adapter}" ]; then
        echo "Disabling proxy on ${adapter}..."
        sudo networksetup -setwebproxy "${adapter}" "" ""
        sudo networksetup -setwebproxystate "${adapter}" off
    fi
    sleep 20
    
    # Test that Docker pull is working
    docker pull hello-world || (echo "unable to pull hello-world image" && exit 1)
    
    # Check the env variables are not present
    docker run -it --rm alpine env > out.tmp || (echo "unable to get alpine env" && exit 1)
    grep -i -v 'http' out.tmp
    grep -v 'http_proxy' out.tmp || (echo "http_proxy is set in container" && exit 1)

    # Check internet access from a container
    docker run -it --rm alpine wget -O- http://www.docker.com > out.tmp || ( echo "unable to run wget in alpine shell, got:" && cat out.tmp && exit 1)
    grep '</html>' out.tmp || ( echo "---- unexpected response: ----" && cat out.tmp && exit 1)
    
}

trap cleanup EXIT

setup
run
