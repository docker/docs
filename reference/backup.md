<!--[metadata]>
+++
title ="backup"
description="Create a backup of DTR"
keywords= ["docker, dtr, cli, backup"]
[menu.main]
parent="dtr_menu_reference"
identifier="dtr_reference_backup"
+++
<![end-metadata]-->

# docker/dtr backup

Create a backup of DTR

## Usage

```bash
docker run -i --rm docker/dtr \
    backup [command options] > backup.tar
```

## Description


This command creates a tar file with the contents of the volumes used by
DTR, and prints it. You can then use the 'restore' command to restore the data
from an existing backup.

Note:

  * This command only creates backups of configurations, and image metadata.
    It doesn't backup users and organizations. Users and organizations can be
    backed up when performing a UCP backup.

    It also doesn't backup the Docker images stored in your registry.
    You should implement a separate backup policy for the Docker images stored
    in your registry, taking in consideration whether your DTR installation is
    configured to store images on the filesystem or using a cloud provider.

  * This backup contains sensitive information and should be
    stored securely.


## Options

| Option                    | Description                |
|:--------------------------|:---------------------------|
|`--ucp-url`|The UCP URL including domain and port|
|`--ucp-username`|The UCP administrator username|
|`--ucp-password`|The UCP administrator password|
|`--debug`|Enable debug mode for additional logging|
|`--hub-username`|Username to use when pulling images|
|`--hub-password`|Password to use when pulling images|
|`--ucp-insecure-tls`|Disable TLS verification for UCP|
|`--ucp-ca`|Use a PEM-encoded TLS CA certificate for UCP|
|`--existing-replica-id`|The ID of an existing DTR replica|
|`--config-only`|Backup/restore only the configurations of DTR and not the database|

