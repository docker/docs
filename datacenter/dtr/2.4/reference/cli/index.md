---
title: docker/dtr overview
description: Learn about the commands available in the docker/dtr image.
keywords: dtr, install, uninstall, configure
---

This tool has commands to install, configure, and back up Docker
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

| Option                                    | Description                |
|:------------------------------------------|:---------------------------|
|[install](install)| Install Docker Trusted Registry                 |
|[join](join)| Add a new replica to an existing DTR cluster                 |
|[reconfigure](reconfigure)| Change DTR configurations                 |
|[remove](remove)| Remove a DTR replica from a cluster                 |
|[destroy](destroy)| Destroy a DTR replica's data                 |
|[restore](restore)| Install and restore DTR from an existing backup                 |
|[backup](backup)| Create a backup of DTR                 |
|[upgrade](upgrade)| Upgrade DTR 2.0.0 or later cluster to this version                 |
|[dumpcerts](dumpcerts)| Print the TLS certificates used by DTR                 |
|[images](images)| List all the images necessary to install DTR                 |

