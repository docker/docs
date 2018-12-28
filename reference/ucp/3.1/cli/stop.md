---
title: docker/ucp stop
description: Stop UCP components running on this node
keywords: ucp, cli, stop
---

Stop UCP components running on this node

## Usage

```
docker container run --rm -it \
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
