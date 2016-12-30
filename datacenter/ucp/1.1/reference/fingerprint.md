---
description: Dump out TLS certificates.
keywords: fingerprint, ucp, tool, reference, ucp
title: docker/ucp fingerprint
---

Dump out the TLS fingerprint for the UCP controller running on this
Docker Engine.

## Usage

```
docker run --rm \
  --name ucp \
  -v /var/run/docker.sock:/var/run/docker.sock \
  docker/ucp \
  fingerprint
```

## Description

This utility will display the certificate fingerprint of the UCP controller
running on the local engine.  This can be used when scripting 'join'
operations for the '--fingerprint' flag.