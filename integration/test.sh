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
# Test init, add, publish and list
./bin/notary -d $tmpdir init -s $1 $new_repo || echo "FAILED"
./bin/notary -d $tmpdir add $new_repo $new_tag README.md || echo "FAILED"
./bin/notary -d $tmpdir publish -s $1 $new_repo || echo "FAILED"
./bin/notary -d $tmpdir list -s $1 $new_repo | grep $new_tag || echo "FAILED"

# Test remove, publish and list
./bin/notary -d $tmpdir remove $new_repo $new_tag || echo "FAILED"
./bin/notary -d $tmpdir publish -s $1 $new_repo || echo "FAILED"
./bin/notary -d $tmpdir list -s $1 $new_repo | grep $new_tag && echo "FAILED"

# Test key existence/key removal
./bin/notary -d $tmpdir key list | grep -v '^$'| grep $new_repo | wc -l | grep 2 || echo "FAILED"
rootID=`./bin/notary -d $tmpdir key list | grep -v '^$' | grep -v $new_repo | grep -v "#" | awk '{print $0}'`
./bin/notary -d $tmpdir key remove -y -r $rootID || echo "FAILED"
./bin/notary -d $tmpdir key list | grep -v '^$'| grep -v $new_repo | grep -v "#" | wc -l | grep 0 || echo "FAILED"

# Test cert existence/cert removal
./bin/notary -d $tmpdir cert list | grep -v '^$'| grep $new_repo | wc -l | grep 1 || echo "FAILED"
certID=`./bin/notary -d $tmpdir cert list | grep $new_repo | awk '{print $2}'`
./bin/notary -d $tmpdir cert remove -y $certID || echo "FAILED"
./bin/notary -d $tmpdir cert list | grep -v '^$' | grep $new_repo && echo "FAILED"
