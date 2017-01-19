#!/bin/bash

ARCHIVE_BRANCHES=("v1.4" "v1.5" "v1.6" "v1.7" "v1.8" "v1.9" "v1.10" "v1.11" "master")

for VERSION in ${ARCHIVE_BRANCHES[@]}; do
  echo "STARTING to build image for $VERSION"
  docker build --no-cache -t docs:${VERSION} --build-arg VER=$VERSION .
  echo "Finished building image for $VERSION"
done
