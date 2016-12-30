---
description: Create a support dump for this UCP node
keywords: docker, dtr, cli, support
title: docker/ucp support
---

Create a support dump for this UCP node

## Usage

```bash

docker run --rm \
    --name ucp \
    -v /var/run/docker.sock:/var/run/docker.sock \
    docker/ucp \
    support [command options] > docker-support.tgz

```

## Description

This command creates a support dump file for this node, and prints it to stdout.
