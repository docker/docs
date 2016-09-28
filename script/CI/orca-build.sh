#!/bin/bash

# Env Vars - TODO - refine how we assert these are set, and how the CI jobs set them
#
# Required:
# PUSH "true" or "false"
# ORG
# REGISTRY_EMAIL
# REGISTRY_USERNAME
# REGISTRY_PASSWORD
# BUILD_NUMBER
# WORKSPACE
#
# Optional:
# CONTENT_TRUST_KEY


set -eu

VERSION=$(make print-UCP_VERSION)
echo ""
echo "Starting official build for ${VERSION}"
echo ""
cd ${WORKSPACE}
git log -1


# Make sure we've got the latest base images
make freshen

# Set the tag to match the Official Build number, unless overridden
export PRE_RELEASE=${PRE_RELEASE:-ob${BUILD_NUMBER}}

# Start with a clean foundation
make clean

./script/run make -j --output-sync=target ORG=${ORG} DOCKER_BUILD_FLAGS=--no-cache build image

echo ""
echo "Generating bundle"
echo ""
make ORG=${ORG} bundle

if [ "${PUSH}" = "true" ] ; then
    echo ""
    echo "Pushing images to hub in the ${ORG} org"
    echo ""
    export NOTARY_DELEGATION_PASSPHRASE=${REGISTRY_PASSWORD}
    export DOCKER_CONTENT_TRUST_REPOSITORY_PASSPHRASE=${REGISTRY_PASSWORD}
    echo "Logging in"
    docker login -e $REGISTRY_EMAIL -u $REGISTRY_USERNAME -p $REGISTRY_PASSWORD
    if [ -z "${CONTENT_TRUST_KEY:-}" ] ; then
        echo "Not signing the images, set CONTENT_TRUST_KEY to point to a delegation cert"
        _TRUST=0
    else
        echo "Importing notary key"
        cp -f ${CONTENT_TRUST_KEY} $(basename ${CONTENT_TRUST_KEY})
        ./script/run notary -D key import $(basename ${CONTENT_TRUST_KEY}) --role user -d ~/.docker/trust
        rm -f $(basename ${CONTENT_TRUST_KEY})
        _TRUST=1
    fi

    # At this point we don't need to keep rebuilding the dev image
    export SKIP_DEV_IMAGE_BUILD=1

    echo "Pushing..."
    ./script/run make DOCKER_CONTENT_TRUST=${_TRUST} ORG=${ORG} push

    # If we're operating on the "docker" org, never tag latest, since we do that explicitly
    if [ "${ORG}" = "docker" ] ; then
        echo ""
        echo "NOT marking latest - you must do that manually with the orca-bootstrap-update-latest build"
        echo ""
    else
        echo ""
        echo "Marking build latest"
        echo ""
        # Also get a "latest" tag based on the version
        docker tag -f ${ORG}/ucp:$(make print-TAG) ${ORG}/ucp:${VERSION}-latest
        ./script/run docker push ${ORG}/ucp:${VERSION}-latest

        # Make sure "latest" points to this build for lurkers to get bleeding edge
        docker tag -f ${ORG}/ucp:$(make print-TAG) ${ORG}/ucp:latest
        ./script/run docker push ${ORG}/ucp:latest
    fi
else
    echo "Not pushing.  Use the artifacts from this build to test"
fi
