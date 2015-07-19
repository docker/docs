#!/bin/sh

set -x
set -e

rm -rf ~/.docker/trust
make binaries
./bin/notary init $1
./bin/notary add $1 v1 README.md
./bin/notary publish $1
./bin/notary list $1
