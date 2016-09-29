#!/bin/sh

cd /tmp || exit 1
touch foo
ln -s foo bar
if [ "$(readlink bar)" != "foo" ]; then
  echo Symlink /tmp/bar is not pointing to foo
  exit 1
fi
