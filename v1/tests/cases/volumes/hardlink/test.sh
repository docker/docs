#!/bin/sh

cd /tmp
echo hello > foo
ln foo bar
if [ "$(cat bar)" != "hello" ]; then
  echo Hard link /tmp/bar does not have the same contents as /tmp/foo
  exit 1
fi
