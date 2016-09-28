#!/bin/sh

set -ex

ORCA_PKG="github.com/docker/orca"
ENZI_PKG="$ORCA_PKG/enzi"

# Create a container which will build the enzi binary.
C_ID=$(docker create \
	-w /go/src/github.com/docker/orca/enzi \
	golang:1.6.2-alpine \
		go build -o enzi $ENZI_PKG
)

# Copy our vendored source code into the orca vendor directory.
docker cp ../vendor "$C_ID:/go/src/$ORCA_PKG/vendor"
# Copy the enzi source code into the target directory in the container.
docker cp . "$C_ID:/go/src/$ENZI_PKG"

# Start and attach to the container.
docker start -a "$C_ID"

# Copy out the built binary.
docker cp "$C_ID:/go/src/$ENZI_PKG/enzi" ./enzi

# Cleanup by removing the container.
docker rm "$C_ID"
