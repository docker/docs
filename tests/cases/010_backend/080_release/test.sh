#!/bin/sh
# SUMMARY: Basic docker-release test
# LABELS: checkout, !win, !circleci
# AUTHOR: French Ben <frenchben@docker.com>

set -e # Exit on error
. "${RT_PROJECT_ROOT}/_lib/lib.sh"

# Import the release check
. "${OSX_SCRIPTS}/Docker.Utils"

updateChannel=""

# Set the channel per release check
export CIRCLE_TAG="mac-v1.12.0"
SetUploadChannel
# Check if they are equal
assert_equals "stable" "$updateChannel"

# Set the channel per release check
export CIRCLE_TAG="mac-v1.12.0-beta16.2"
SetUploadChannel
# Check if they are equal
assert_equals "beta" "$updateChannel"

# Set the channel per release check
export CIRCLE_TAG="mac-rc-v1.12.0-beta16.2"
SetUploadChannel
# Check if they are equal
assert_equals "rc" "$updateChannel"

# Set the channel per release check
export CIRCLE_TAG="mac-test-v1.12.0-beta16.2"
SetUploadChannel
# Check if they are equal
assert_equals "test" "$updateChannel"

# Set the channel per release check
export CIRCLE_TAG="mac-CS-1.12.v1.12.0-beta16.2"
SetUploadChannel
# Check if they are equal
assert_equals "CS-1.12" "$updateChannel"

# cleanup
unset CIRCLE_TAG

# Set the channel per release check
export CI_PULL_REQUEST="pull/1234"
SetUploadChannel
# Check if they are equal
assert_equals "pr" "$updateChannel"

# cleanup
unset CI_PULL_REQUEST

# Set the channel per release check
export CIRCLE_BRANCH="master"
SetUploadChannel
# Check if they are equal
assert_equals "master" "$updateChannel"

# Set the channel per release check
export CIRCLE_BRANCH="b14-release-fixes"
SetUploadChannel
# Check if they are equal
assert_equals "default" "$updateChannel"

# cleanup
unset CIRCLE_BRANCH

exit 0
