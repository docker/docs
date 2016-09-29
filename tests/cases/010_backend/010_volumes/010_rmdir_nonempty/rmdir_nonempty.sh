#!/bin/sh

cd /tmp || exit 1
rm -rf foo
mkdir foo
touch foo/bar
(rmdir foo 2> output || true)
cat output
grep "Directory not empty" output
