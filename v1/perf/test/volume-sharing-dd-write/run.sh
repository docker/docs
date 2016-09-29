#!/bin/bash

. common.sh

# Run the test
mkdir -p volumes
if [ ! -e volumes/128MiB ]; then
  echo creating 128MiB volume
  dd if=/dev/zero of=./volumes/128MiB bs=1048576 count=128
fi

mkdir -p logs
for BS in ${BLOCK_SIZES}; do
  count=$(echo "134217728 / ${BS}" | bc)
  echo Writing 128MiB in ${count} blocks of size ${BS}
  docker run -v `pwd`/volumes:/volumes volume-sharing-dd-write /dd.sh ${BS} ${count}
  mv volumes/${BS}.time logs/
done
