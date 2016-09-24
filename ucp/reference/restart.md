---
description: Restart Docker Universal Control Plane containers.
keywords:
- install, ucp, restart
menu:
  main:
    identifier: ucp_ref_restart
    parent: ucp_ref
title: restart
---

# docker/ucp restart

Start or restart UCP components on this engine

## Usage

```bash
$ docker run --rm -it \
     --name ucp \
     -v /var/run/docker.sock:/var/run/docker.sock \
     docker/ucp \
     restart [command options]
```

## Options

| Option        | Description                                      |
|:--------------|:-------------------------------------------------|
| `--debug, -D` | Enable debug mode                                |
| `--jsonlog`   | Produce json formatted output for easier parsing |
