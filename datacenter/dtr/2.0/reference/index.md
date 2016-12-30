---
description: Learn about the options available on the docker/dtr image.
keywords: docker, dtr, install, uninstall, configure
title: docker/dtr overview
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

| Option        | Description                                                                     |
|:--------------|:--------------------------------------------------------------------------------|
| `install`     | Install Docker Trusted Registry on this Docker Engine                           |
| `join`        | Add a new replica to an existing DTR cluster                                    |
| `reconfigure` | Change DTR configurations                                                       |
| `remove`      | Remove a replica from a DTR cluster                                             |
| `restore`     | Create a new DTR cluster from an existing backup                                |
| `backup`      | Backup a DTR cluster to a tar file and stream it to stdout                      |
| `migrate`     | Migrate configurations, accounts, and repository metadata from DTR 1.4.3 to 2.0 |
| `upgrade`     | Upgrade a v2.0.0 or later cluster to this version of DTR                        |
| `dumpcerts`   | Dump out the TLS certificates used by this DTR instance                         |
| `images`      | Lists all the images necessary to install DTR                                   |