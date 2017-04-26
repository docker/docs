---
title: docker/dtr backup
keywords: docker, dtr, cli, backup
description: Create a backup of DTR
---

Create a backup of DTR

## Usage

```bash
docker run -i --rm docker/dtr \
    backup [command options] > backup.tar
```

## Description


This command creates a tar file with the contents of the volumes used by
DTR, and prints it. You can then use the `restore` command to restore the data
from an existing backup.

Note:

  * This command only creates backups of configurations, and image metadata.
    It doesn't back up users and organizations. Users and organizations can be
    backed up when performing a UCP backup.

    It also doesn't back up the Docker images stored in your registry.
    You should implement a separate backup policy for the Docker images stored
    in your registry, taking in consideration whether your DTR installation is
    configured to store images on the filesystem or using a cloud provider.

  * This backup contains sensitive information and should be
    stored securely.


## Options

| Option                    | Description                |
|:--------------------------|:---------------------------|
|`--debug`|Enable debug mode for additional logging|
|`--existing-replica-id`|The ID of an existing DTR replica|
|`--hub-password`|Password to use when pulling images|
|`--hub-username`|Username to use when pulling images|
|`--ucp-ca`|Use a PEM-encoded TLS CA certificate for UCP|
|`--ucp-insecure-tls`|Disable TLS verification for UCP|
|`--ucp-password`|The UCP administrator password|
|`--ucp-url`|The UCP URL including domain and port|
|`--ucp-username`|The UCP administrator username|

