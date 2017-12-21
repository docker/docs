---
title: docker/ucp restart
description: Start or restart UCP components running on this node
keywords: ucp, cli, restart
---

Start or restart UCP components running on this node

## Usage

```
docker container run --rm -it \
        --name ucp \
        -v /var/run/docker.sock:/var/run/docker.sock \
        docker/ucp \
        restart [command options]
```

## Options

| Option                    | Description                |
|:--------------------------|:---------------------------|
|`--debug, D`|Enable debug mode|
|`--jsonlog`|Produce json formatted output for easier parsing|
