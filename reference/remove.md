<!--[metadata]>
+++
title = "remove"
description = "Docker Trusted Registry remove command reference."
keywords = ["docker, registry, reference, remove"]
[menu.main]
parent="dtr_menu_reference"
identifier="dtr_reference_remove"
+++
<![end-metadata]-->

# docker/dtr remove

Remove a replica from a DTR cluster

## Usage

```bash
$ docker run -it --rm docker/dtr \
  remove [command options]
```

## Description

This command removes a replica from the cluster, stops and removes all
DTR containers, and deletes all DTR volumes.

| Option                  | Description                                                                                  |
|:------------------------|:---------------------------------------------------------------------------------------------|
| `--ucp-url`             | Specify the UCP controller URL [$UCP_URL]                                                    |
| `--ucp-username`        | Specify the UCP admin username [$UCP_USERNAME]                                               |
| `--ucp-password`        | Specify the UCP admin password [$UCP_PASSWORD]                                               |
| `--debug`               | Enable debug mode, provides additional logging [$DEBUG]                                      |
| `--hub-username`        | Specify the Docker Hub username for pulling images [$HUB_USERNAME]                           |
| `--hub-password`        | Specify the Docker Hub password for pulling images [$HUB_PASSWORD]                           |
| `--ucp-insecure-tls`    | Disable TLS verification for UCP [$UCP_INSECURE_TLS]                                         |
| `--ucp-ca`              | Use a PEM-encoded TLS CA certificate for UCP [$UCP_CA]                                       |
| `--replica-id`          | Specify the replica Id. Must be unique per replica, leave blank for random [$DTR_REPLICA_ID] |
| `--existing-replica-id` | ID of an existing replica in a cluster [$DTR_EXISTING_REPLICA_ID]                            |
