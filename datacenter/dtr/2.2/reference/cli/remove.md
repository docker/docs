---
title: docker/dtr remove
keywords: docker, dtr, cli, remove
description: Remove a DTR replica from a cluster
---

Remove a DTR replica from a cluster

## Usage

```bash
docker run -it --rm docker/dtr \
    remove [command options]
```

## Description


This command gracefully scales down your DTR cluster by removing exactly one replica. All other replicas must be healthy and will remain healthy after this operation.


## Options

| Option                    | Description                |
|:--------------------------|:---------------------------|
|`--debug`|Enable debug mode for additional logging|
|`--existing-replica-id`|The ID of an existing DTR replica|
|`--hub-password`|Password to use when pulling images|
|`--hub-username`|Username to use when pulling images|
|`--replica-id`|The ID of the replica you want to remove from the cluster|
|`--ucp-ca`|Use a PEM-encoded TLS CA certificate for UCP|
|`--ucp-insecure-tls`|Disable TLS verification for UCP|
|`--ucp-password`|The UCP administrator password|
|`--ucp-url`|The UCP URL including domain and port|
|`--ucp-username`|The UCP administrator username|

