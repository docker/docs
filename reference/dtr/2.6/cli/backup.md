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

### Example Commands

#### Basic

```bash 
docker run -i --rm --log-driver none docker/dtr:2.6.5 \
    backup --ucp-ca "$(cat ca.pem)" --existing-replica-id 5eb9459a7832 > backup.tar
```

#### Advanced (with chained commands)

The following command has been tested on Linux:

{% raw %}
```none
DTR_VERSION=$(docker container inspect $(docker container ps -f \
  name=dtr-registry -q) | grep -m1 -Po '(?<=DTR_VERSION=)\d.\d.\d'); \
REPLICA_ID=$(docker inspect -f '{{.Name}}' $(docker ps -q -f name=dtr-rethink) | cut -f 3 -d '-')); \
read -p 'ucp-url (The UCP URL including domain and port): ' UCP_URL; \
read -p 'ucp-username (The UCP administrator username): ' UCP_ADMIN; \
read -sp 'ucp password: ' UCP_PASSWORD; \
docker run --log-driver none -i --rm \
  --env UCP_PASSWORD=$UCP_PASSWORD \
  docker/dtr:$DTR_VERSION backup \
  --ucp-username $UCP_ADMIN \
  --ucp-url $UCP_URL \
  --ucp-ca "$(curl https://${UCP_URL}/ca)" \
  --existing-replica-id $REPLICA_ID > \
  dtr-metadata-${DTR_VERSION}-backup-$(date +%Y%m%d-%H_%M_%S).tar
```
{% endraw %}

For a detailed explanation on the advanced example, see 
[Back up your DTR metadata](/ee/dtr/admin/disaster-recovery/create-a-backup/#back-up-dtr-metadata).
To learn more about the `--log-driver` option for `docker run`, see [docker run reference](/engine/reference/run/#logging-drivers---log-driver). 

## Description

This command creates a `tar` file with the contents of the volumes used by
DTR, and prints it. You can then use `docker/dtr restore` to restore the data
from an existing backup.


Note:

  * This command only creates backups of configurations, and image metadata.
    It does not back up users and organizations. Users and organizations can be
    backed up during a UCP backup.

    It also does not back up Docker images stored in your registry.
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

