---
description: Print the TLS fingerprint for this UCP web server
keywords: docker, dtr, cli, fingerprint
title: docker/ucp fingerprint
---

Print the TLS fingerprint for this UCP web server

## Usage

```bash

docker run --rm \
    --name ucp \
    -v /var/run/docker.sock:/var/run/docker.sock \
    docker/ucp \
    fingerprint [command options]

```

## Description

This command displays the fingerprint of the certificate used in the UCP web
server running on this node.
