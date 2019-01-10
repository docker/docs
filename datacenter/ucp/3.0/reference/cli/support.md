---
title: docker/ucp support
description: Create a support dump for UCP nodes
keywords: ucp, cli, support, support dump, troubleshooting
redirect_from:
 - /reference/ucp/3.0/cli/support/
---

Create a support dump for specified UCP nodes. You create a support dump to help [Docker Support](http://success.docker.com/support) understand your environment and more effectively troubleshoot issues in resolving your support case.

## Usage

```
docker container run --rm \
        --name ucp \
        -v /var/run/docker.sock:/var/run/docker.sock \
        docker/ucp \
        support [command options] > docker-support.tgz
```

## Description

This command creates a support dump file for this node, and prints it to stdout.

## Options

| Option                    | Description                |
|:--------------------------|:---------------------------|
|`--debug, D`|Enable debug mode|
|`--jsonlog`|Produce json formatted output for easier parsing|
