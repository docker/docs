#!/bin/bash

set -eu

# Required Env Vars:
# INTEGRATION_TEST_SCOPE
# JOB_NAME
# REGISTRY_PASSWORD
# REGISTRY_USERNAME
# WORKSPACE

# Optional Env Vars:
PULL_IMAGES="${PULL_IMAGES:-1}"
ORCA_ORG="${ORCA_ORG:-dockerorcadev}"
TEST_PARALLEL="${TEST_PARALLEL:-2}"
SWARM_IMAGE="${SWARM_IMAGE:-}"
TAG="${TAG:-$(make print-TAG PRE_RELEASE=latest)}"

# Make a prefix suitable for use in machine names.
export MACHINE_PREFIX=matrix$(echo $JOB_NAME | md5sum | grep -o '.\{8\}' | head -1)
echo prefix: $JOB_NAME=$MACHINE_PREFIX

export TEST_FLAGS="-v --short"

# Use run_inc to grab integration.log
./script/run make \
    SELENIUM_URL="" \
    TEST_TIMEOUT=60m \
    TAG=$TAG \
    TEST_PARALLEL=${TEST_PARALLEL} \
    INTEGRATION_TEST_SCOPE=${INTEGRATION_TEST_SCOPE} \
    ORCA_ORG=${ORCA_ORG} \
    PULL_IMAGES=${PULL_IMAGES} \
    PURGE_MACHINES=1 \
    SWARM_IMAGE=${SWARM_IMAGE} \
    integration
