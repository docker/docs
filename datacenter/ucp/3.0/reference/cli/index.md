---
title: docker/ucp overview
description: Learn about the commands available in the docker/ucp image.
keywords: ucp, cli, ucp
---

This image has commands to install and manage
Docker Universal Control Plane (UCP) on a Docker Engine.

You can configure the commands using flags or environment variables. When using
environment variables, use the `docker container run -e VARIABLE_NAME` syntax to pass the
value from your shell, or `docker container run -e VARIABLE_NAME=value` to specify the
value explicitly on the command line.

The container running this image needs to be named `ucp` and bind-mount the
Docker daemon socket. Below you can find an example of how to run this image.

Additional help is available for each command with the `--help` flag.

## Usage

```bash
docker container run -it --rm \
    --name ucp \
    -v /var/run/docker.sock:/var/run/docker.sock \
    docker/ucp \
    command [command arguments]
```

## Commands

| Option           | Description                                               |
|:-----------------|:----------------------------------------------------------|
| `install`        | Install UCP on this node                                  |
| `restart`        | Start or restart UCP components running on this node      |
| `stop`           | Stop UCP components running on this node                  |
| `upgrade`        | Upgrade the UCP cluster                                   |
| `images`         | Verify the UCP images on this node                        |
| `uninstall-ucp`  | Uninstall UCP from this swarm                             |
| `dump-certs`     | Print the public certificates used by this UCP web server |
| `support`        | Create a support dump for this UCP node                   |
| `id`             | Print the ID of UCP running on this node                  |
| `backup`         | Create a backup of a UCP manager node                     |
| `restore`        | Restore a UCP cluster from a backup                       |
| `example-config` | Display an example configuration file for UCP             |
