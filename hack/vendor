#!/usr/bin/env bash

set -eu -o pipefail

output=$(mktemp -d -t hugo-vendor-output.XXXXXXXXXX)

function clean {
  rm -rf "$output"
}

trap clean EXIT

docker buildx bake vendor \
	"--set=*.output=type=local,dest=${output}"
rm -r _vendor
cp -R $output/* .
