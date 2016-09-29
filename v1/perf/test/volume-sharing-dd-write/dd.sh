#!/bin/bash
BS=$1
count=$2
if [ "$(uname)" = "Darwin" ]; then
  FLAGS=
else
  FLAGS=oflag=dsync
fi

./timethis.py volumes/${BS}.time  /bin/dd if=/dev/zero of=volumes/128MiB bs=${BS} count=${count} ${FLAGS}

