---
description: Print the ID of UCP running on this node
keywords: docker, dtr, cli, id
title: docker/ucp id
---

Print the ID of UCP running on this node

## Usage

```bash

docker run -it --rm \
    --name ucp \
    -v /var/run/docker.sock:/var/run/docker.sock \
    docker/ucp \
    id [command options]

```

## Description

This command prints the ID of the UCP components running on this node. This ID
matches what you see when running the 'docker info' command while using
a client bundle.

This ID is used by other commands as confirmation.
