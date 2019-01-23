---
title: docker/dtr backup
description: Create a backup of DTR
keywords: dtr, cli, backup
redirect_from:
 - /reference/dtr/2.5/cli/backup/
---

Create a backup of DTR

## Usage

```bash
docker run -i --rm docker/dtr \
    backup [command options] > backup.tar
```

### Example Usage
```bash 
docker run -i --rm docker/dtr \
    backup --ucp-ca "$(cat ca.pem)" --existing-replica-id 5eb9459a7832 > backup.tar
```

## Description


This command creates a `tar` file with the contents of the volumes used by
DTR, and prints it. You can then use `docker/dtr restore` to restore the data
from an existing backup. 

Note:

  * This command only creates backups of configurations, and image metadata.
    It does not back up users and organizations. Users and organizations can be
    backed up during a UCP backup.

    It also doesn't back up Docker images stored in your registry.
    You should implement a separate backup policy for the Docker images stored
    in your registry, taking into consideration whether your DTR installation is
    configured to store images on the filesystem or is using a cloud provider.

  * This backup contains sensitive information and should be
    stored securely.

  * Using the `--offline-backup` flag temporarily shuts down the RethinkDB container. 
    Take the replica out of your load balancer to avoid downtime.


## Options

| Option                        | Environment Variable      | Description                                                                          |
|:------------------------------|:--------------------------|:-------------------------------------------------------------------------------------|
| `--debug` | $DEBUG | Enable debug mode for additional logs. |
| `--existing-replica-id` | $DTR_REPLICA_ID | The ID of an existing DTR replica. To add, remove or modify a DTR replica, you must connect to an existing healthy replica's database. |
| `--help-extended` | $DTR_EXTENDED_HELP | Display extended help text for a given command. |
| `--offline-backup` | $DTR_OFFLINE_BACKUP | This flag takes RethinkDB down during backup and takes a more reliable backup. If you back up DTR with this flag, RethinkDB will go down during backup. However, offline backups are guaranteed to be more consistent than online backups. |
| `--ucp-ca` | $UCP_CA | Use a PEM-encoded TLS CA certificate for UCP. Download the UCP TLS CA certificate from `https://<ucp-url>/ca`, and  use `--ucp-ca "$(cat ca.pem)"`. |
| `--ucp-insecure-tls` | $UCP_INSECURE_TLS | Disable TLS verification for UCP. The installation uses TLS but always trusts the TLS certificate used by UCP, which can lead to MITM (man-in-the-middle) attacks.  For production deployments, use `--ucp-ca "$(cat ca.pem)"` instead. |
| `--ucp-password` | $UCP_PASSWORD | The UCP administrator password. |
| `--ucp-url` | $UCP_URL | The UCP URL including domain and port. |
| `--ucp-username` | $UCP_USERNAME | The UCP administrator username. |

