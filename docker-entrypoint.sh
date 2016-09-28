#!/bin/bash -ex

wrapdocker &

# Give the service a moment to warm up...
sleep 5

if [[ -z "$@" ]]
then
  # No args provided, just run make
  make
else
  exec "$@"
fi
