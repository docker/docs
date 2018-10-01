---
title: docker/ucp support
description: Create a support dump for this UCP node
keywords: ucp, cli, support
---

Create a support dump for this UCP node

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
|`--debug, D`|Enable debug mode.|
|`--jsonlog`|Produce json formatted output for easier parsing.|
|`--loglines`|Specify number of lines to grab from `journalctl`. The default is 10,000 lines.|
|`--servicedriller`|Run the swarm service driller tool.  The default is off.|
|`--nodes`|Select specific nodes on which to produce a support dump. Comma-separated node IDs are allowed. The default is all nodes.|
|`--verbose`|Collect additional data , including a list of all non-system tasks or services running within the cluster. The default is off. Customer-specific info is not included.|
