+++
title = "support"
description = "Generate a support dump for this engine."
keywords= ["docker, ucp, support, logs"]
[menu.main]
parent = "ucp_ref"
identifier = "ucp_ref_support"
+++

# docker/ucp support

Generate a support dump for this engine.

## Usage

```
docker run --rm \
  --name ucp \
  -v /var/run/docker.sock:/var/run/docker.sock \
  docker/ucp \
  support > docker-support.tgz
```

## Description

This utility will produce a support dump file on stdout for this local node.
