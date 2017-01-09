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


This command cleans up any data associated with a replica. It is useful for cleaning up garbage in case of failures or to uninstall DTR. For scaling down a cluster see the remove command.


## Options

| Option                    | Description                |
|:--------------------------|:---------------------------|
|`--replica-id`|Assign an ID to the DTR replica. By default the ID is random|
|`--ucp-url`|The UCP URL including domain and port|
|`--ucp-username`|The UCP administrator username|
|`--ucp-password`|The UCP administrator password|
|`--debug`|Enable debug mode for additional logging|
|`--hub-username`|Username to use when pulling images|
|`--hub-password`|Password to use when pulling images|
|`--ucp-insecure-tls`|Disable TLS verification for UCP|
|`--ucp-ca`|Use a PEM-encoded TLS CA certificate for UCP|

