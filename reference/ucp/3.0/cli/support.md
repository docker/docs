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
|`--servicedriller`|Run the swarm service driller (ssd) tool. For more information on this tool, see [Docker Swarm Service Driller(ssd)](https://github.com/sanimej/ssd) Not run by default.|
|`--nodes`|Select specific nodes on which to produce a support dump. Comma-separated node IDs are allowed. The default selects all nodes.|
