---
title: docker/ucp uninstall-ucp
description: Uninstall UCP from this swarm
keywords: ucp, cli, uninstall-ucp
---

Uninstall UCP from this swarm

## Usage

```
 docker container run --rm -it \
        --name ucp \
        -v /var/run/docker.sock:/var/run/docker.sock \
        docker/ucp \
        uninstall-ucp [command options]
```

## Description

This command uninstalls UCP from the swarm, but preserves the swarm so that
your applications can continue running.

After UCP is uninstalled you can use the `docker swarm leave` and
`docker node rm` commands to remove nodes from the swarm.

Once UCP is uninstalled, you can't join nodes to the swarm unless
UCP is installed again.


## Options

| Option                | Description                                                 |
| :-------------------- | :---------------------------------------------------------- |
| `--debug, D`          | Enable debug mode                                           |
| `--jsonlog`           | Produce json formatted output for easier parsing            |
| `--interactive, i`    | Run in interactive mode and prompt for configuration values |
| `--pull`              | Pull UCP images: `always`, when `missing`, or `never`       |
| `--registry-username` | Username to use when pulling images                         |
| `--registry-password` | Password to use when pulling images                         |
| `--id`                | The ID of the UCP instance to uninstall                     |
| `--purge-config`      | Remove UCP configs during uninstallation                    |
