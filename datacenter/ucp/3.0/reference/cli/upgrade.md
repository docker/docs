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

This command upgrades the UCP running on this node.
To upgrade UCP:

  * Upgrade the Docker Engine in all nodes (optional)
  * Run the upgrade command in all manager nodes
  * Run the upgrade command in all worker nodes

Before performing an upgrade, you should perform a backup by using the 
[backup](backup.md) command.

After upgrading UCP in a node, go to the UCP web UI and confirm the node is
healthy, before upgrading other nodes.


## Options

| Option                    | Description                |
|:--------------------------|:---------------------------|
|`--debug, D`|Enable debug mode|
|`--jsonlog`|Produce json formatted output for easier parsing|
|`--interactive, i`|Run in interactive mode and prompt for configuration values|
|`--admin-username`|The UCP administrator username|
|`--admin-password`|The UCP administrator password|
|`--pull`|Pull UCP images: `always`, when `missing`, or `never`|
|`--registry-username`|Username to use when pulling images|
|`--registry-password`|Password to use when pulling images|
|`--id`|The ID of the UCP instance to upgrade|
|`--host-address`|Override the previously configured host address with this IP or network interface|
