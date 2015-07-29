#!/bin/sh

set -x
set -e

export NOTARY_ROOT_PASSPHRASE="ponies"
export NOTARY_SNAPSHOT_PASSPHRASE="ponies"
export NOTARY_TARGET_PASSPHRASE="ponies"

tmpdir=`mktemp -d -t notary-integration-XXXXXXXXXXXXXXX`
new_repo=$(cat /dev/urandom | tr -dc 'a-zA-Z0-9' | fold -w 32 | head -n 1)
new_tag=$(cat /dev/urandom | tr -dc 'a-zA-Z0-9' | fold -w 10 | head -n 1)

make binaries
./bin/notary -d $tmpdir init -s $1 $new_repo || echo "FAILED"
./bin/notary -d $tmpdir add $new_repo $new_tag README.md || echo "FAILED"
./bin/notary -d $tmpdir publish -s $1 $new_repo || echo "FAILED"
./bin/notary -d $tmpdir list -s $1 $new_repo | grep $new_tag || echo "FAILED"
./bin/notary -d $tmpdir remove $new_repo $new_tag || echo "FAILED"
./bin/notary -d $tmpdir publish -s $1 $new_repo || echo "FAILED"
./bin/notary -d $tmpdir list -s $1 $new_repo | grep $new_tag && echo "FAILED"
./bin/notary -d $tmpdir key list | grep $new_repo | wc -l | grep 3 || echo "FAILED"