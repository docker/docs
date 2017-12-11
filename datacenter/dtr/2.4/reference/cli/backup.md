---
title: docker/dtr backup
description: Create a backup of DTR
keywords: dtr, cli, backup
---

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

| Option                        | Environment Variable      | Description                                                                          |
|:------------------------------|:--------------------------|:-------------------------------------------------------------------------------------|
| `--debug` | $DEBUG | Enable debug mode for additional logs. |
| `--existing-replica-id` | $DTR_REPLICA_ID | The ID of an existing DTR replica.To add, remove or modify DTR, you must connect to an existing  healthy replica's database.. |
| `--help-extended` | $DTR_EXTENDED_HELP | Display extended help text for a given command. |
| `--ucp-ca` | $UCP_CA | Use a PEM-encoded TLS CA certificate for UCP.Download the UCP TLS CA certificate from https://<ucp-url>/ca, and  use --ucp-ca "$(cat ca.pem)". |
| `--ucp-insecure-tls` | $UCP_INSECURE_TLS | Disable TLS verification for UCP.The installation uses TLS but always trusts  the TLS certificate used by UCP, which can lead to man-in-the-middle attacks.  For production deployments, use --ucp-ca "$(cat ca.pem)" instead. |
| `--ucp-password` | $UCP_PASSWORD | The UCP administrator password. |
| `--ucp-url` | $UCP_URL | The UCP URL including domain and port. |
| `--ucp-username` | $UCP_USERNAME | The UCP administrator username. |

