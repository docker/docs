---
description: Stop UCP components running on this node
keywords: docker, dtr, cli, stop
title: docker/ucp stop
---

Stop UCP components running on this node

## Usage

```bash
docker run -it --rm \
    --name ucp \
    -v /var/run/docker.sock:/var/run/docker.sock \
    docker/ucp \
    stop [command options]

```

## Options

| Option                    | Description                |
|:--------------------------|:---------------------------|
|`--debug, D`|Enable debug mode|
|`--jsonlog`|Produce json formatted output for easier parsing|
