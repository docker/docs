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
        --unmanaged-cni <true|false>
        upgrade [command options]
```

## Description

This command upgrades the UCP running on this cluster.

Before performing an upgrade, you should perform a backup by using the
[backup](backup.md) command.

After upgrading UCP, go to the UCP web UI and confirm each node is
healthy and that all nodes have been upgraded successfully.


## Options

| Option                | Description                                                                                           |
|:----------------------|:------------------------------------------------------------------------------------------------------|
| `--debug, D`          | Enable debug mode                                                                                     |
| `--jsonlog`           | Produce json formatted output for easier parsing                                                      |
| `--interactive, i`    | Run in interactive mode and prompt for configuration values                                           |
| `--admin-username`    | The UCP administrator username                                                                        |
| `--admin-password`    | The UCP administrator password                                                                        |
| `--pull`              | Pull UCP images: `always`, when `missing`, or `never`                                                 |
| `--registry-username` | Username to use when pulling images                                                                   |
| `--registry-password` | Password to use when pulling images                                                                   |
| `--unmanaged-cni`        | This determines who manages the CNI plugin, using `true` or `false`. The default is false. The `true` value installs UCP without a managed CNI plugin. UCP and the Kubernetes components will be running but pod to pod networking will not function until a CNI plugin is manually installed. This will impact some functionality of UCP until a CNI plugin is running.    
|
| `--id`                | The ID of the UCP instance to upgrade                                                                 |
| `--host-address`      | Override the previously configured host address with this IP or network interface                     |
| `--force-minimums`    | Force the install/upgrade even if the system does not meet the minimum requirements                   |
| `--pod-cidr`          | Kubernetes cluster IP pool for the pods to allocated IP from (Default: 192.168.0.0/16                 |
| `--nodeport-range`    | Allowed port range for Kubernetes services of type NodePort (Default: 32768-35535)                    |
| `--cloud-provider`    | The cloud provider for the cluster                                                                    |
