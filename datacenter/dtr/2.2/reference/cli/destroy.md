---
title: docker/dtr destroy
keywords: docker, dtr, cli, destroy
description: Destroy a DTR replica's data
---

Destroy a DTR replica's data

## Usage

```bash
docker run -it --rm docker/dtr \
    destroy [command options]
```

## Description


This command forcefully removes all containers and volumes associated with the given DTR replica without notifying the rest of the cluster. Use it on all replicas when you want to uninstall DTR.


## Options

| Option                    | Description                |
|:--------------------------|:---------------------------|
|`--replica-id`|The ID of the replica to destroy|
|`--ucp-url`|The UCP URL including domain and port|
|`--ucp-username`|The UCP administrator username|
|`--ucp-password`|The UCP administrator password|
|`--debug`|Enable debug mode for additional logging|
|`--hub-username`|Username to use when pulling images|
|`--hub-password`|Password to use when pulling images|
|`--ucp-insecure-tls`|Disable TLS verification for UCP|
|`--ucp-ca`|Use a PEM-encoded TLS CA certificate for UCP|

