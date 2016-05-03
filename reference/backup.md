<!--[metadata]>
+++
title = "backup"
description = "Docker Trusted Registry backup command reference."
keywords = ["docker, registry, reference, backup"]
[menu.main]
parent="dtr_menu_reference"
identifier="dtr_reference_backup"
+++
<![end-metadata]-->

# docker/dtr backup

Backup a DTR cluster to a tar file and stream it to stdout

## Usage

```bash
$ docker run -it --rm docker/dtr \
 backup [command options] > backupt.tar
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


| Option                  | Description                                                                           |
|:------------------------|:--------------------------------------------------------------------------------------|
| `--ucp-url`             | Specify the UCP controller URL [$UCP_URL]                                             |
| `--ucp-username`        | Specify the UCP admin username [$UCP_USERNAME]                                        |
| `--ucp-password`        | Specify the UCP admin password [$UCP_PASSWORD]                                        |
| `--debug`               | Enable debug mode, provides additional logging [$DEBUG]                               |
| `--hub-username`        | Specify the Docker Hub username for pulling images [$HUB_USERNAME]                    |
| `--hub-password`        | Specify the Docker Hub password for pulling images [$HUB_PASSWORD]                    |
| `--ucp-insecure-tls`    | Disable TLS verification for UCP [$UCP_INSECURE_TLS]                                  |
| `--ucp-ca`              | Use a PEM-encoded TLS CA certificate for UCP [$UCP_CA]                                |
| `--existing-replica-id` | ID of an existing replica in a cluster [$DTR_EXISTING_REPLICA_ID]                     |
| `--config-only`         | Backup/restore only the configurations of DTR and not the database [$DTR_CONFIG_ONLY] |
