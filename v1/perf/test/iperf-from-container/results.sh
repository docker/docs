#!/bin/sh

. ../../lib/functions

if [ "${network}" = "slirp" ]; then
  echo '# Skipping test because we are using user-space networking'
  exit 0
fi

echo '# time speed-mbit/sec'
GBITS=$(cat logs/stdout | grep '0.0-30.0 sec' | grep 'Mbits/sec' | awk '{print $7}')
echo "$(date +%s) ${GBITS}"
