---
title: docker/ucp support
description: Create a support dump for UCP nodes
keywords: ucp, cli, support, support dump, troubleshooting
---

Create a support dump for specified UCP nodes

## Usage

```
docker container run --rm \
        --name ucp \
        -v /var/run/docker.sock:/var/run/docker.sock \
        docker/ucp \
        support [command options] > docker-support.tgz
```

## Description

This command creates a support dump file for the specified node(s), and prints it to stdout.

## Options

| Option       | Description                                      |
|:-------------|:-------------------------------------------------|
| `--debug, D` | Enable debug mode                                |
| `--jsonlog`  | Produce json formatted output for easier parsing |
