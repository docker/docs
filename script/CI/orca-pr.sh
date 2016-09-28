#!/bin/bash

# Required Env Vars - TODO - refine how we assert these are set, and how the CI jobs set them
#
# MACHINE_DRIVER
# AWS_VPC_ID -- assuming AWS driver
# AWS_DEFAULT_REGION  -- assuming AWS driver
# AWS_ACCESS_KEY_ID
# AWS_SECRET_ACCESS_KEY
# MACHINE_CREATE_FLAGS
# REGISTRY_EMAIL
# REGISTRY_USERNAME
# REGISTRY_PASSWORD


set -e


echo ""
echo "Starting PR build"
echo ""
cd ${WORKSPACE}
git log -1

# Make sure we're logged in
echo "Logging in"
docker login -e $REGISTRY_EMAIL -u $REGISTRY_USERNAME -p $REGISTRY_PASSWORD


# Get our machine named something unique to this build so we can run PRs in parallel
export MACHINE_PREFIX=orca-pr-${BUILD_NUMBER}

# Set the tag to match the PR number
export PRE_RELEASE=pr${ghprbPullId}

# Force the ORG to dockerorcadev for PR builds
export ORG=dockerorcadev

# Make sure we always clean up our machine
function cleanup {
    echo "Cleaning up (regardless of pass/fail)"
    ./script/run make clean-test-machine
    echo "All done"
}
trap cleanup EXIT

echo ""
echo "Updating base images"
echo ""

make freshen

echo ""
echo "Setting environment for AWS machines"
echo ""
eval $(./script/matrix_helper engine=oss-test platform=aws-ubuntu15.10 swarm=default)

echo ""
echo "Beginning build and unit tests"
echo ""

# Build and unit test
./script/run make ORG=${ORG} -j --output-sync=target image test create-test-machine

echo ""
echo "Loading images to the test VM"
echo ""

# Load the image up to the test machine
./script/run make ORG=${ORG} load-test-machine

echo ""
echo "Pushing, and running integration tests"
echo ""

# Now run the tests and push in parallel
TEST_FLAGS="-v --short" ./script/run make -j --output-sync=target ORG=${ORG} push integration USE_TEST_MACHINE=1 TEST_TIMEOUT=30m INTEGRATION_TEST_SCOPE=./integration/acceptance/c1_w0/...

echo ""
echo "Tests passed!  Updating tags on hub for easier consumption"
echo ""


# Finally clean up the tags to make it a little easier to consume
TAG=$(make print-TAG)
# Wire up an extra tag for the slack messages that isn't version specific
docker tag -f ${ORG}/ucp:${TAG} ${ORG}/ucp:pr${ghprbPullId}
docker push ${ORG}/ucp:pr${ghprbPullId}
