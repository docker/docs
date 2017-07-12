---
title: docker/ucp fingerprint
description: Print the TLS fingerprint for this UCP web server
keywords: docker, ucp, cli, fingerprint
---

Print the TLS fingerprint for this UCP web server

## Usage
```
docker run --rm \
    --name ucp \
    -v /var/run/docker.sock:/var/run/docker.sock \
    docker/ucp \
    fingerprint
```


## Description

This command displays the fingerprint of the certificate used in the UCP web
server running on this node.

