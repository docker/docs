---
title: Back up UCP
description: Learn how to create a backup of UCP
keywords: enterprise, backup, ucp
redirect_from:
 - /ee/ucp/admin/backups-and-disaster-recovery/
---

>{% include enterprise_label_shortform.md %}

UCP backups no longer require pausing the reconciler and deleting UCP containers, and backing up a UCP manager does not disrupt the manager’s activities. 

Because UCP stores the same data on all manager nodes, you only need to back up a single UCP manager node.

User resources, such as services, containers, and stacks are not affected by this
operation and continue operating as expected.

## Limitations
  
- Backups should not be utilized for restoring clusters on a cluster with a newer version of Docker Enterprise. For example, if backups occur on version N, then a restore on version N+1 is not supported. 
- More than one backup at the same time is not supported. If a backup is attempted while another backup is in progress, or if two backups are scheduled at the same time, a message is displayed to indicate that the second backup failed because another backup is running.
- For crashed clusters, backup capability is not guaranteed. Perform regular backups to avoid this situation. 
- UCP backup does not include swarm workloads.

## UCP backup contents
Backup contents are stored in a `.tar` file. Backups contain UCP configuration metadata to re-create configurations such as **Administration Settings** values such as LDAP and SAML, and RBAC configurations (Collections, Grants, Roles, User, and more): 

| Data                  | Description                                                                        | Backed up |
| :---------------------|:-----------------------------------------------------------------------------------|:----------|
| Configurations        | UCP configurations, including Docker Engine - Enterprise license. Swarm, and client CAs          | yes
| Access control        | Permissions for teams to swarm resources, including collections, grants, and roles | yes
| Certificates and keys | Certificates and public and private keys used for authentication and mutual TLS communication         | yes
| Metrics data          | Monitoring data gathered by UCP                                                    | yes
| Organizations         | Users, teams, and organizations                                                        | yes
| Volumes               | All [UCP named volumes](/ee/ucp/ucp-architecture/#volumes-used-by-ucp/), including all UCP component certificates and data. [Learn more about UCP named volumes](/ee/ucp/ucp-architecture/).              | yes
| Overlay Networks      | Swarm-mode overlay network definitions, including port information                                             | no
| Configs, Secrets      | Create a Swarm backup to backup these data                                         | no
| Services              | Stacks and services are stored in Swarm-mode or SCM/Config Management              | no  

> Note
> 
> Because Kubernetes stores the state of resources on `etcd`, a backup of `etcd` is sufficient for stateless backups and is described [here](https://kubernetes.io/docs/tasks/administer-cluster/configure-upgrade-etcd/#backing-up-an-etcd-cluster).

## Data not included in the backup
* `ucp-metrics-data`: holds the metrics server's data.
* `ucp-node-certs` : holds certs used to lock down UCP system components
* Routing mesh settings. Interlock L7 ingress configuration information is not captured in UCP backups. A manual backup and restore process is possible and should be performed.

## Kubernetes settings, data, and state

UCP backups include all Kubernetes declarative objects (pods, deployments, replicasets, configurations, and so on), including secrets. These objects are stored in the `ucp-kv etcd` database that is backed up (and restored) as part of UCP backup/restore.


> Note
>
> You cannot back up Kubernetes volumes and node labels. Instead, upon restore, Kubernetes declarative objects are re-created. Containers are re-created and IP addresses are resolved. 

For more information, see [Backing up an etcd cluster](https://kubernetes.io/docs/tasks/administer-cluster/configure-upgrade-etcd/#backing-up-an-etcd-cluster).

## Specify a backup file
To avoid directly managing backup files, you can specify a file name and host directory on a secure and configured storage backend, such as NFS or another networked file system. The file system location is the backup folder on the manager node file system. This location must be writable by the `nobody` user, which is specified by changing the folder ownership to `nobody`. This operation requires administrator permissions to the manager node, and must only be run once for a given file system location. 

```
sudo chown nobody:nogroup /path/to/folder
```
> Important
> 
> Specify a different name for each backup file. Otherwise, the existing backup file with the same name is overwritten. Specify a location that is mounted on a fault-tolerant file system (such as NFS) rather than the node's local disk. Otherwise, it is important to regularly move backups from the manager node's local disk to ensure adequate space for ongoing backups.
 
## UCP backup steps
There are several options for creating a UCP backup:

- [CLI](#create-a-ucp-backup-using-the-cli)
- [UI](#create-a-ucp-backup-using-the-ui)
- [API](#create-list-and-retrieve-ucp-backups-using-the-api) 

The backup process runs on one manager node.

### Create a UCP backup using the CLI
The following example shows how to create a UCP manager node backup, encrypt it
by using a passphrase, decrypt it, verify its contents, and store it locally on
the node at `/tmp/mybackup.tar`:

Run the `{{ page.ucp_org }}/{{ page.ucp_repo }}:{{ page.ucp_version }} backup` command on a single UCP manager and include the `--file` and `--include-logs`options. This creates a tar archive with the contents of all [volumes used by UCP](/ee/ucp-architecture/) and streams it to `stdout`.
Replace `{{ page.ucp_version }}` with the version you are currently running.

```bash
$ docker container run \
    --rm \
    --log-driver none \
    --name ucp \
    --volume /var/run/docker.sock:/var/run/docker.sock \
    --volume /tmp:/backup \
    {{ page.ucp_org }}/{{ page.ucp_repo }}:{{ page.ucp_version }} backup \
    --file mybackup.tar \
    --passphrase "secret12chars" \
    --include-logs=false
```

> Note
>
> If you are running with Security-Enhanced Linux (SELinux) enabled,
> which is typical for RHEL hosts, you must include `--security-opt
> label=disable` in the `docker` command (replace `version` with the version
> you are currently running):

```bash
$ docker container run \
    --rm \
    --log-driver none \
    --security-opt label=disable \
    --name ucp \
    --volume /var/run/docker.sock:/var/run/docker.sock \
    {{ page.ucp_org }}/{{ page.ucp_repo }}:{{ page.ucp_version }} backup \
    --passphrase "secret12chars" > /tmp/mybackup.tar
 ```   

> Note
>
> To determine whether SELinux is enabled in the engine, view the host’s `/etc/docker/daemon.json` file, and search for the string `"selinux-enabled":"true"`.

#### View log and progress information
To view backup progress and error reporting, view the contents of the stderr streams of the running backup container during the backup. Progress is updated for each backup step, for example, after validation, after volumes are backed up, after `etcd` is backed up, and after `rethinkDB`. Progress is not preserved after the backup has completed. 

#### Verify a UCP backup
In a valid backup file, 27 or more files are displayed in the list and the `./ucp-controller-server-certs/key.pem` file is present. Ensure the backup is a valid tar file by listing its contents, as shown in the following example: 

```
$ gpg --decrypt /directory1/directory2/backup.tar | tar --list
```

If decryption is not needed, you can list the contents by removing the `--decrypt flag`, as shown in the following example:

```
$ tar --list -f /directory1/directory2/backup.tar
```

### Create a UCP backup using the UI

To create a UCP backup using the UI:

1. In the UCP UI, navigate to **Admin Settings**.
2. Select **Backup Admin**. 
3. Select **Backup Now** to trigger an immediate backup.

The UI also provides the following options:
 - Display the status of a running backup
 - Display backup history
 - View backup contents

### Create, list, and retrieve UCP backups using the API

The UCP API provides three endpoints for managing UCP backups. You must be a UCP administrator to access these API endpoints. 

#### Create a UCP backup using the API
You can create a backup with the `POST: /api/ucp/backup` endpoint. This is a JSON endpoint with the following arguments: 
 
| field name 	|   JSON data type*  	|                description               	|
|:----------:	|:-------:	|:----------------------------------------:	|
|  passphrase 	| string 	|      Encryption passphrase     	|
|  noPassphrase 	| bool 	|      Set to `true` if not using a passphrase     	|
|  fileName  	|  string 	|          Backup file name          	|
| includeLogs 	|   bool  	|       Specifies whether to include a log file      	|
| hostPath 	| string  	| [File system location](#specify-a-backup-file) 	|

The request returns one of the following HTTP status codes, and, if successful, a backup ID.

 - 200: Success
 - 500: Internal server error
 - 400: Malformed request (payload fails validation)

##### Example 

```
$ curl -sk -H 'Authorization: Bearer $AUTHTOKEN' https://$UCP_HOSTNAME/api/ucp/backup \
   -X POST \
   -H "Content-Type: application/json" \
   --data  '{"encrypted": true, "includeLogs": true, "fileName": "backup1.tar", "logFileName": "backup1.log", "hostPath": "/secure-location"}'
200 OK
```

where:

 - `$AUTHTOKEN` is your authentication bearer token if using auth token identification.
 - `$UCP_HOSTNAME` is your UCP hostname.

#### List all backups using the API

You can view all existing backups with the `GET: /api/ucp/backups` endpoint. This request does not expect a payload and returns a list of backups, each as a JSON object following the schema found in the [Backup schema](#backup-schema) section.

The request returns one of the following HTTP status codes and, if successful, a list of existing backups:

 - 200: Success
 - 500: Internal server error

##### Example

```
curl -sk -H 'Authorization: Bearer $AUTHTOKEN' https://$UCP_HOSTNAME/api/ucp/backups
[
  {
    "id": "0d0525dd-948a-41b4-9f25-c6b4cd6d9fe4",
    "encrypted": true,
    "fileName": "backup2.tar",
    "logFileName": "backup2.log",
    "backupPath": "/secure-location",
    "backupState": "SUCCESS",
    "nodeLocation": "ucp-node-ubuntu-0",
    "shortError": "",
    "created_at": "2019-04-10T21:55:53.775Z",
    "completed_at": "2019-04-10T21:56:01.184Z"
  },
  {
    "id": "2cf210df-d641-44ca-bc21-bda757c08d18",
    "encrypted": true,
    "fileName": "backup1.tar",
    "logFileName": "backup1.log",
    "backupPath": "/secure-location",
    "backupState": "IN_PROGRESS",
    "nodeLocation": "ucp-node-ubuntu-0",
    "shortError": "",
    "created_at": "2019-04-10T01:23:59.404Z",
    "completed_at": "0001-01-01T00:00:00Z"
  }
]
```

#### Retrieve backup details using the API

You can retrieve details for a specific backup using the `GET: /api/ucp/backup/{backup_id}` endpoint, where `{backup_id}` is the ID of an existing backup. This request returns the backup, if it exists, for the specified ID, as a JSON object following the schema found in the [Backup schema](#backup-schema) section.

The request returns one of the following HTTP status codes, and if successful, the backup for the specified ID:

 - 200: Success
 - 404: Backup not found for the given `{backup_id}`
 - 500: Internal server error

#### Backup schema

The following table describes the backup schema returned by the `GET` and `LIST` APIs:

|  field name  	| JSON data type* 	|                             description                             	|
|:------------:	|:---------------:	|:-------------------------------------------------------------------:	|
|      id      	|      string     	|                             Unique ID                             	|
|   encrypted  	|     boolean     	|                 Set to `true` if encrypted with a passphrase                 	|
|   fileName   	|      string     	|  Backup file name if backing up to a file, empty otherwise  	|
|  logFileName 	|      string     	| Backup log file name if saving  backup logs, empty otherwise 	|
|  backupPath  	|      string     	|                    Host path where backup resides                   	|
|  backupState 	|      string     	|     Current state of the backup  (`IN_PROGRESS`, `SUCCESS`, `FAILED`)     	|
| nodeLocation 	|      string     	|                  Node on which the backup was taken                 	|
|  shortError  	|      string     	|           Short error. Empty unless `backupState` is set to `FAILED`           	|
|  created_at  	|      string     	|                     Time of backup creation                     	|
| completed_at 	|      string     	|                    Time of backup completion                    	|

> Note
> 
> *= JSON data type as defined per [JSON RFC 7159](https://tools.ietf.org/html/rfc7159).


### Where to go next

- [Back up the Docker Trusted Registry](back-up-dtr.md)
 
