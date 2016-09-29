#!/bin/sh -ex

VERSION=2

docker build -t aws -f Dockerfile.awscli .

if [ ! -e "s3creds.env" ]; then
  echo Need to initialise S3 config from s3creds.env.in
  exit 1
fi

docker run -v `pwd`/data:/mnt --env-file s3creds.env \
  aws s3 sync --size-only s3://docker-pinata-support/incoming /mnt
