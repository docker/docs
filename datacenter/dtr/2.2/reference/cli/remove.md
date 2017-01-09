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


This command removes a replica from a DTR deployment. All DTR containers and
volumes are removed from the node.


## Options

| Option                    | Description                |
|:--------------------------|:---------------------------|
|`--debug`|Enable debug mode for additional logging|
|`--existing-replica-id`|The ID of an existing DTR replica|
|`--hub-password`|Password to use when pulling images|
|`--hub-username`|Username to use when pulling images|
|`--replica-id`|Assign an ID to the DTR replica. By default the ID is random|
|`--ucp-ca`|Use a PEM-encoded TLS CA certificate for UCP|
|`--ucp-insecure-tls`|Disable TLS verification for UCP|
|`--ucp-password`|The UCP administrator password|
|`--ucp-url`|The UCP URL including domain and port|
|`--ucp-username`|The UCP administrator username|

