#!/bin/bash

SCRIPT="build"

if [ "$DEV" == "1" ]; then
	SCRIPT="build-dev"
fi;

# Copy the linked build files into the cached NPM install directory with package.json
cp -nrT /build /opt
# Compile JS using either the prod "build" script or dev "build-dev" script
# need to set cache path due to permission issues
BABEL_CACHE_PATH=/opt/.babel.json npm run-script $SCRIPT
# Copy the built files back to the linked directory
cp /opt/src/bundle.js /build/src/
