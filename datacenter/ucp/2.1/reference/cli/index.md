---
title: docker/ucp overview
description: Learn about the commands available in the docker/ucp image.
keywords:
- docker, ucp, cli, docker run --rm -it \
        --name ucp \
        -v /var/run/docker.sock:/var/run/docker.sock \
        docker/ucp \

---

A new cli application

## Commands

| Option          | Description                                               |
|:----------------|:----------------------------------------------------------|
| `install`       | Install UCP on this node                                  |
| `restart`       | Start or restart UCP components running on this node      |
| `stop`          | Stop UCP components running on this node                  |
| `upgrade`       | Upgrade the UCP components on this node                   |
| `images`        | Verify the UCP images on this node                        |
| `uninstall-ucp` | Uninstall UCP from this swarm                             |
| `dump-certs`    | Print the public certificates used by this UCP web server |
| `fingerprint`   | Print the TLS fingerprint for this UCP web server         |
| `support`       | Create a support dump for this UCP node                   |
| `id`            | Print the ID of UCP running on this node                  |
| `backup`        | Create a backup of a UCP manager node                     |
| `restore`       | Restore a UCP cluster from a backup                       |
