---
description: Dump out the ID of the UCP components running on this engine.
keywords: docker, ucp, id
title: docker/ucp id
---

Dump out the ID of the UCP components running on this engine.

## Usage

```
docker run --rm \
  --name ucp \
  -v /var/run/docker.sock:/var/run/docker.sock \
  docker/ucp \
  id
```

## Description

This utility will display the ID of the local UCP components running
on this node. This ID matches what you see when you run 'docker info'
pointed to the UCP controller(s) and is required by various commands
in this tool as confirmation.