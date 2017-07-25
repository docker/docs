---
title: docker/ucp example-config
description: Display an example configuration file for UCP
keywords: ucp, cli, config, configuration
---

Display an example configuration file for UCP

## Usage

```
docker container run --rm -i \
    --name ucp \
    -v /var/run/docker.sock:/var/run/docker.sock \
    docker/ucp \
    example-config > ucp.config
```

## Description

This command emits an example configuration file for setting up UCP.
[Learn about UCP configuration files](../../guides/admin/configure/ucp-configuration-file.md). 