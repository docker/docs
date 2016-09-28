#!/bin/bash

# WARNING!!!
#
# The docker build is not very good about reproducibility, and incremental
# builds across different commits Things WILL blow up sporadically and you'll
# have to stop the bisect, go scorched earth, and resume to get it back in good
# shape.  If you see a lot of skips happening, stop it and investigate.

# Run this from within the docker/docker tree
MACHINE=dh-manual-test1
IMAGE=docker/ucp:1.1.0-dev

# Usage:
# git bisect start HEAD 9b630197111fabedd3e382eb94812ecba0435be7
# git bisect run <this script>
#
# Then in another window, periodically
# git bisect log

# Refined ranges
# GOOD: 39b799a, 0b882cc, cdd8d39, f7fa83c, 6eaf443, 719a43c, d274456, 264b5b6, 0300361, c48439a
# BAD: dd51e85, 68674f7, 731c59f, upstream/master

set -e

echo ""
echo "Applying patches"
echo ""
cp Dockerfile Dockerfile.orig
patch Dockerfile < Dockerfile.patch || (rm Dockerfile.orig; /bin/true)
patch -p1 < build.patch || /bin/true
patch -p1 < containerd.patch || /bin/true

echo ""
echo "Building docker"
echo ""
make binary || exit 125

# Revert the patches
mv Dockerfile.orig Dockerfile
git checkout hack/make/binary
git checkout hack/make/gccgo
git checkout libcontainerd/remote_linux.go


echo ""
echo "Updating target machine ${MACHINE}"
echo ""
for bin in docker  docker-containerd  docker-containerd-ctr  docker-containerd-shim  docker-runc runc containerd containerd-shim ctr; do
    # Ignore failures...
    echo "Moving ${MACHINE} /usr/bin/${bin} -> /usr/bin/${bin}.prior"
    docker-machine ssh ${MACHINE} sudo mv /usr/bin/${bin} /usr/bin/${bin}.prior || /bin/true
    if [ -f bundles/latest/binary/$bin ] ; then
        echo "Uploading $bin..."
        cat bundles/latest/binary/$bin | docker-machine ssh ${MACHINE} sudo sh -c "\" cat - > /usr/bin/${bin}; chmod a+x /usr/bin/${bin}\""
    fi
done

echo "Generating empty daemon config"
docker-machine ssh ${MACHINE} sudo sh -c "\" echo '{}' > /etc/docker/daemon.json \"" || exit 125

echo "Bouncing daemon"
docker-machine ssh ${MACHINE} sudo systemctl restart docker.service || exit 125
sleep 3
echo "Checking version"
docker-machine ssh ${MACHINE} sudo docker version || (docker-machine ssh ${MACHINE} sudo journalctl -n 20 -u docker.service; exit 125)

# Now get down to the test case
eval $(docker-machine env ${MACHINE})

# Clean out any prior cruft, but preserve images
ID=$(docker run --rm --name ucp -v /var/run/docker.sock:/var/run/docker.sock docker/ucp:1.1.0-dev id||/bin/true)
if [ -n "${ID}" ] ; then
    echo ""
    echo "Cleaning prior run"
    echo ""
    docker run --rm -it \
        --name ucp \
        -v /var/run/docker.sock:/var/run/docker.sock \
        docker/ucp:1.1.0-dev \
        uninstall --id ${ID} --preserve-images || /bin/true
fi

echo ""
echo "Initial install"
echo ""
docker run --rm -it \
        --name ucp \
        -v /var/run/docker.sock:/var/run/docker.sock \
        docker/ucp:1.1.0-dev \
        install --swarm-port 3376 || exit 125


echo ""
echo "Re-install"
echo ""
docker run --rm -it \
        --name ucp \
        -v /var/run/docker.sock:/var/run/docker.sock \
        docker/ucp:1.1.0-dev \
        install -D --swarm-port 3376 --fresh-install

docker version
echo ""
echo "PASS!!! The daemon is OK"
