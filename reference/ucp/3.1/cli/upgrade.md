---
title: docker/ucp upgrade
description: Upgrade the UCP components on this node
keywords: ucp, cli, upgrade
---

Upgrade the UCP cluster

## Usage

```
 docker container run --rm -it \
        --name ucp \
        -v /var/run/docker.sock:/var/run/docker.sock \
        docker/ucp \
        upgrade [command options]
```

## Description

This command upgrades the UCP running on this cluster.

Before performing an upgrade, you should perform a backup by using the
[backup](backup.md) command.

After upgrading UCP, go to the UCP web interface and confirm each node is
healthy and that all nodes have been upgraded successfully.


## Options

| Option                        | Description                                                                         |
|:------------------------------|:------------------------------------------------------------------------------------|
| `--debug, D`                  | Enable debug mode                                                                   |
| `--jsonlog`                   | Produce json formatted output for easier parsing                                    |
| `--interactive, i`            | Run in interactive mode and prompt for configuration values                         |
| `--admin-password` *value*    | The UCP administrator password                                                      |
| `--admin-username` *value*    | The UCP administrator username                                                      |
| `--force-minimums`            | Force the install/upgrade even if the system does not meet the minimum requirements |
| `--host-address` *value*      | Override the previously configured host address with this IP or network interface   |
| `--id`                        | The ID of the UCP instance to upgrade                                               |
| `--pull`                      | Pull UCP images: `always`, when `missing`, or `never`                               |
| `--registry-password` *value* | Password to use when pulling images                                                 |
| `--registry-username` *value* | Username to use when pulling images                                                 |
