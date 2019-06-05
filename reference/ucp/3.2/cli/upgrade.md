---
title: docker/ucp upgrade
description: Upgrade the UCP components on this node
keywords: ucp, cli, upgrade
---

Upgrade the UCP cluster.

## Usage

```
 docker container run --rm -it \
        --name ucp \
        -v /var/run/docker.sock:/var/run/docker.sock \
        docker/ucp \
        upgrade [command options]
```

## Description

This command upgrades the UCP running on this cluster. To upgrade UCP:

- (Optional) Upgrade the Docker Engine in all nodes.
- Run the upgrade command on one manager node.

Before performing an upgrade, you should perform a backup by using the
[backup](backup.md) command.

After upgrading UCP, go to the UCP web interface and confirm each node is
healthy and that all nodes have been upgraded successfully.


## Options

| Option                | Description                                                                                           |
|:----------------------|:------------------------------------------------------------------------------------------------------|
| `--admin-username`    | The UCP administrator username                                                                        |
| `--admin-password`    | The UCP administrator password                                                                        |
| `--cloud-provider`    | The cloud provider for the cluster                                                                    |
| `--debug, D`          | Enable debug mode                                                                                     |
| `--force-minimums`    | Force the install/upgrade even if the system does not meet the minimum requirements                   |
| `--host-address`      | Override the previously configured host address with this IP or network interface                     |
| `--id`                | The ID of the UCP instance to upgrade                                                                 |
| `--jsonlog`           | Produce json formatted output for easier parsing                                                      |
| `--interactive, i`    | Run in interactive mode and prompt for configuration values                                           |
| `--manual-worker-upgrade`    | Specifies whether to manually upgrade worker nodes. Defaults to `false`.                                           |
| `--nodeport-range`    | Allowed port range for Kubernetes services of type `NodePort`. The default port range is `32768-35535`.                   |
| `--pod-cidr`          | Kubernetes cluster IP pool for the pods to allocated IP. The default IP pool is `192.168.0.0/16`.                 |
| `--pull`              | Pull UCP images: `always`, when `missing`, or `never`                                                 |
| `--registry-username` | Username to use when pulling images                                                                   |
| `--registry-password` | Password to use when pulling images                                                                   |
