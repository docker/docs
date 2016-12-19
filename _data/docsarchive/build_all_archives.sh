#!/bin/bash

ARCHIVE_BRANCHES=("v1.4" "v1.5" "v1.6" "v1.7" "v1.8" "v1.9" "v1.10" "v1.11" "master")
CURRENT_VERSION="v1.12"

# Build the base image
docker build -t docs:base -f Dockerfile.base .

for BRANCH in ${ARCHIVE_BRANCHES[@]}; do
	if [ "$BRANCH" == "master" ]
	then
		VERSION="$CURRENT_VERSION"
	else
		VERSION="$BRANCH"
	fi

	BASEURL="$VERSION/"

	echo "STARTING to build image for $VERSION using branch $BRANCH"
	docker build --no-cache -t docs:${VERSION} --build-arg VERSION=$VERSION --build-arg BRANCH=$BRANCH --build-arg BASEURL=$BASEURL .
	echo "Finished building image for $VERSION"
	if [ "$BRANCH" == "master" ]
	then
		# Build the "current" version without permalink (/version/)
		docker build --no-cache -t docs:latest --build-arg VERSION=$VERSION  --build-arg BRANCH=$BRANCH .
	fi
done
