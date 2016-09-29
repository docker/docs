#!/bin/sh
# SUMMARY:
# LABELS:
# REPEAT:
# AUTHOR:

set -e # Exit on error
# Source libraries. Uncomment if needed/defined
#. "${RT_ROOT}/lib/lib.sh"
#. "${RT_PROJECT_ROOT}/_lib/lib.sh"
# IMAGE_NAME=  # Use a env variable to name images/containers

clean_up() {
    # remove any files, containers, images etc
}
trap clean_up EXIT

# Test code goes here

exit 0
