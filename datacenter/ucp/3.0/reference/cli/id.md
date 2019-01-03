---
title: docker/ucp id
description: Print the ID of UCP running on this node
keywords: ucp, cli, id
redirect_from:
 - /reference/ucp/3.0/cli/id/
---

Print the ID of UCP running on this node

## Usage

Run the UCP id command using the correct image version tag for UCP (i.e. docker/ucp:$version)

```
docker container run --rm \
    --name ucp \
    -v /var/run/docker.sock:/var/run/docker.sock \
    docker/ucp:$version \
    id
```

## Description

This command prints the ID of the UCP components running on this node. This ID
matches what you see when running the `docker info` command while using
a client bundle.

This ID is used by other commands as confirmation.

