---
description: Backup a DTR cluster to a tar file and stream it to stdout
keywords: docker, dtr, cli, backup
title: docker/dtr backup
---

Backup a DTR cluster to a tar file and stream it to stdout

## Usage

```bash
docker run -i --rm docker/dtr \
    backup [command options] > backup.tar
```

## Description


This command creates a backup of configurations and image metadata, in a tar
file. It can be used to make periodic backups of a DTR installation.

This command only creates backups of configurations, and image metadata.
It does not backup the Docker images stored in your registry.

You should implement a separate backup policy for the Docker images stored
in your registry, taking in consideration whether your DTR installation is
configured to store images on the filesystem or using a cloud provider.

WARNING: This backup contains sensitive information and should be
stored securely.


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
|`--existing-replica-id`|ID of an existing replica in a cluster|
|`--config-only`|Backup/restore only the configurations of DTR and not the database|
