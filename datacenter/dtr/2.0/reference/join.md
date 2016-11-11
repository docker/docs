---
description: Docker Trusted Registry join command reference.
keywords: docker, registry, reference, join
title: docker/dtr join
---

Add a new replica to an existing DTR cluster

## Usage

```bash
$ docker run -it --rm docker/dtr \
   join [command options]
```


## Description

This command installs DTR on the Docker Engine that runs the command,
and joins the new installation to an existing cluster.

To set up a cluster with high-availability, add 3, 5, or 7 nodes to
the cluster.

## Options

| Option                  | Description                                                                                                                                                                                                                                                                                             |
|:------------------------|:--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `--ucp-url`             | Specify the UCP controller URL including domain and port                                                                                                                                                                                                                                                |
| `--ucp-username`        | Specify the UCP admin username                                                                                                                                                                                                                                                                          |
| `--ucp-password`        | Specify the UCP admin password                                                                                                                                                                                                                                                                          |
| `--debug`               | Enable debug mode, provides additional logging                                                                                                                                                                                                                                                          |
| `--hub-username`        | Specify the Docker Hub username for pulling images                                                                                                                                                                                                                                                      |
| `--hub-password`        | Specify the Docker Hub password for pulling images                                                                                                                                                                                                                                                      |
| `--ucp-insecure-tls`    | Disable TLS verification for UCP                                                                                                                                                                                                                                                                        |
| `--ucp-ca`              | Use a PEM-encoded TLS CA certificate for UCP                                                                                                                                                                                                                                                            |
| `--ucp-node`            | Specify the host to install Docker Trusted Registry                                                                                                                                                                                                                                                     |
| `--replica-id`          | Specify the replica ID. Must be unique per replica, leave blank for random                                                                                                                                                                                                                              |
| `--unsafe`              | Enable this flag to skip safety checks when installing or joining                                                                                                                                                                                                                                       |
| `--existing-replica-id` | ID of an existing replica in a cluster                                                                                                                                                                                                                                                                  |
| `--replica-http-port`   | Specify the public HTTP port for the DTR replica; 0 means unchanged/default                                                                                                                                                                                                                             |
| `--replica-https-port`  | Specify the public HTTPS port for the DTR replica; 0 means unchanged/default                                                                                                                                                                                                                            |
| `--skip-network-test`   | Enable this flag to skip the overlay networking test                                                                                                                                                                                                                                                    |
| `--extra-envs`          | List of extra environment variables to use for deploying the DTR containers for the replica. This can be used to specify swarm constraints. Separate the environment variables with ampersands (&). You can escape actual ampersands with backslashes (\). Can't be used in combination with --ucp-node |