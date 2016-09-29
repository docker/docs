#!/bin/bash
BS=$1
count=$2
if [ "$(uname)" = "Darwin" ]; then
  FLAGS=
else
  FLAGS=iflag=fullblock,dsync
fi
./timethis.py volumes/${BS}.time  /bin/dd if=volumes/128MiB of=/dev/null bs=${BS} count=${count} ${FLAGS}
