#!/bin/sh

if [ -z "$(docker ps)" ]; then
  echo Failed to docker ps
  exit 1
fi
