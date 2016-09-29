#!/bin/sh

if [ "$(docker run busybox echo hello)" != "hello" ]; then
  echo Failed to echo hello
  exit 1
fi
