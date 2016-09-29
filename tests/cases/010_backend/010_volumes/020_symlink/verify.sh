#!/bin/sh
set -e

cd /tmp
ls -l foo bar
if [ "$(readlink bar)" != "foo" ]; then
  echo Symlink /tmp/bar is not pointing to foo
  exit 1
fi
