---
description: Upgrade Docker Universal Control Plane.
keywords: docker, ucp, upgrade
title: docker/ucp upgrade
---

Upgrade the UCP components on this Docker Engine.

## Usage

```
docker run --rm -it \
  --name ucp \
  -v /var/run/docker.sock:/var/run/docker.sock \
  docker/ucp \
  upgrade [OPTIONS]
```

## Description

When upgrading UCP, you must run the 'upgrade' command against every
engine in your cluster.  You should upgrade your controller and replica
nodes first, followed by your compute nodes.  If you plan to upgrade your
engine as well, upgrade the engine first, before upgrading UCP on the node.

After upgrading each node, confirm the node is present in the UCP console
before proceeding to the next node.


## Options

| Option                | Description                                                                             |
|:----------------------|:----------------------------------------------------------------------------------------|
| `--debug, -D`         | Enable debug                                                                            |
| `--jsonlog`           | Produce json formatted output for easier parsing                                        |
| `--interactive, -i`   | Enable interactive mode.,You will be prompted to enter all required information         |
| `--admin-username`    | Specify the UCP admin username [$UCP_ADMIN_USER]                                        |
| `--admin-password`    | Specify the UCP admin password [$UCP_ADMIN_PASSWORD]                                    |
| `--registry-username` | Specify the username to pull required images with [$REGISTRY_USERNAME]                  |
| `--registry-password` | Specify the password to pull required images with [$REGISTRY_PASSWORD]                  |
| `--id`                | The ID of the UCP instance to upgrade                                                   |
| `--pull "always"`     | Specify image pull behavior ('always', when 'missing', or 'never') (default: "missing") |