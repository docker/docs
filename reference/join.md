<!--[metadata]>
+++
title = "join"
description = "Docker Trusted Registry join command reference."
keywords = ["docker, registry, reference, join"]
[menu.main]
parent="dtr_menu_reference"
identifier="dtr_reference_join"
+++
<![end-metadata]-->

# docker/dtr join

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


| Option                     | Description                                                                                  |
|:---------------------------|:---------------------------------------------------------------------------------------------|
| `--ucp-url`                | Specify the UCP controller URL [$UCP_URL]                                                    |
| `--ucp-username`           | Specify the UCP admin username [$UCP_USERNAME]                                               |
| `--ucp-password`           | Specify the UCP admin password [$UCP_PASSWORD]                                               |
| `--debug`                  | Enable debug mode, provides additional logging [$DEBUG]                                      |
| `--hub-username`           | Specify the Docker Hub username for pulling images [$HUB_USERNAME]                           |
| `--hub-password`           | Specify the Docker Hub password for pulling images [$HUB_PASSWORD]                           |
| `--ucp-ca`                 | Use a PEM-encoded TLS CA certificate for UCP [$UCP_CA]                                       |
| `--ucp-node`               | Specify the host to install Docker Trusted Registry [$UCP_NODE]                              |
| `--replica-id`             | Specify the replica Id. Must be unique per replica, leave blank for random [$DTR_REPLICA_ID] |
| `--existing-replica-id`    | ID of an existing replica in a cluster [$DTR_EXISTING_REPLICA_ID]                            |
| `--replica-http-port "0"`  | Specify the public HTTP port for the DTR replica [$REPLICA_HTTP_PORT]                        |
| `--replica-https-port "0"` | Specify the public HTTPS port for the DTR replica [$REPLICA_HTTPS_PORT]                      |
| `--skip-network-test`      | Whether to skip the overlay networking test or not [$DTR_SKIP_NETWORK_TEST]                  |
