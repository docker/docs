---
title: docker/ucp support
description: Create a support dump for UCP nodes
keywords: ucp, cli, support, support dump, troubleshooting
---

>{% include enterprise_label_shortform.md %}

Create a support dump for specified UCP nodes.

## Usage

```
docker container run --rm \
        --name ucp \
        --log-driver none \
        --volume /var/run/docker.sock:/var/run/docker.sock \
        docker/ucp \
        support [command options] > docker-support.tgz
```

## Description

This command creates a support dump file for the specified node(s), and prints
it to stdout. This includes the ID of the UCP components running on the node.
The ID matches what you see when running the `docker info` command while using
a client bundle, and is used by other commands as confirmation.

## Options

| Option       | Description                                      |
|:-------------|:-------------------------------------------------|
| `--debug, D` | Enable debug mode                                |
| `--jsonlog`  | Produce json formatted output for easier parsing |
