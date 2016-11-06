+++
title = "engine-discovery"
description = "Manage the engine discovery configuration on a node."
keywords= ["docker, ucp, discovery"]
[menu.main]
parent = "ucp_ref"
identifier="ucp_ref_engine_discovery"
+++

# docker/ucp engine-discovery

The engine-discovery command is no longer used.

## Usage

```
docker run --rm -it \
  --name ucp \
  -v /var/run/docker.sock:/var/run/docker.sock \
  docker/ucp \
  engine-discovery [options]
```

## Description

The engine-discovery command is no longer used.  Overlay networking is enabled
automatically via swarm-mode.

## Options

```nohighlight
--debug, -D  Enable debug mode
--jsonlog    Produce json formatted output for easier parsing
```
