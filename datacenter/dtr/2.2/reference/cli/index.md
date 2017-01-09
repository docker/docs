---
title: docker/dtr overview
keywords: docker, dtr, install, uninstall, configure
description: Learn about the commands available in the docker/dtr image.
---

This tool has commands to install, configure, and backup Docker
Trusted Registry (DTR). It also allows uninstalling DTR.
By default the tool runs in interactive mode. It prompts you for
the values needed.

Additional help is available for each command with the '--help' option.


## Usage

```bash
docker run -it --rm docker/dtr \
    command [command options]
```


## Commands

| Option        | Description                                        |
|:--------------|:---------------------------------------------------|
| `backup`      | Create a backup of DTR                             |
| `destroy`     | Destroy a DTR replica's data                       |
| `images`      | List all the images necessary to install DTR       |
| `install`     | Install Docker Trusted Registry                    |
| `join`        | Add a new replica to an existing DTR cluster       |
| `reconfigure` | Change DTR configurations                          |
| `remove`      | Remove a DTR replica from a cluster                |
| `restore`     | Install and restore DTR from an existing backup    |
| `upgrade`     | Upgrade DTR 2.0.0 or later cluster to this version |
