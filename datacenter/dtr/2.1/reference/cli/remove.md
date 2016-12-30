---
description: Remove a replica from a DTR cluster
keywords: docker, dtr, cli, remove
title: docker/dtr remove
---

Remove a replica from a DTR cluster

## Usage

```bash
docker run -it --rm docker/dtr \
    remove [command options]
```

## Description


This command removes a replica from the cluster, stops and removes all
DTR containers, and deletes all DTR volumes.


## Options

| Option                    | Description                |
|:--------------------------|:---------------------------|
|`--ucp-url`|Specify the UCP controller URL including domain and port|
|`--ucp-username`|Specify the UCP admin username|
|`--ucp-password`|Specify the UCP admin password|
|`--debug`|Enable debug mode, provides additional logging|
|`--hub-username`|Specify the Docker Hub username for pulling images|
|`--hub-password`|Specify the Docker Hub password for pulling images|
|`--ucp-insecure-tls`|Disable TLS verification for UCP|
|`--ucp-ca`|Use a PEM-encoded TLS CA certificate for UCP|
|`--force-remove`|Force removal of replica even if it can break your cluster's state. Necessary only when --existing-replica-id == --replica-id.|
|`--replica-id`|Specify the replica ID. Must be unique per replica, leave blank for random|
|`--existing-replica-id`|ID of an existing replica in a cluster|
