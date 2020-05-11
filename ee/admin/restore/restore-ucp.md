---
title: Restore UCP
description: Learn how to restore UCP from a backup
keywords: enterprise, restore, swarm
---

>{% include enterprise_label_shortform.md %}

To restore UCP, select one of the following options:

* Run the restore on the machines from which the backup originated or on new machines. You can use the same swarm from which the backup originated or a new swarm.  
* On a manager node of an existing swarm that does not have UCP installed.
  In this case, UCP restore uses the existing swarm and runs instead of any install.
* Run the restore on a docker engine that isn't participating in a swarm, in which case it performs `docker swarm init` in the same way as the install operation would. A new swarm is created and UCP is restored on top.  

## Limitations

- To restore an existing UCP installation from a backup, you need to
uninstall UCP from the swarm by using the `uninstall-ucp` command.
[Learn to uninstall UCP](/ee/ucp/admin/install/uninstall/).
- Restore operations must run using the same major/minor UCP version (and `docker/ucp` image version) as the backed up cluster. Restoring to a later patch release version is allowed. 
- If you restore UCP using a different Docker swarm than the one where UCP was
previously deployed on, UCP will start using new TLS certificates. Existing
client bundles won't work anymore, so you must download new ones.

## Kubernetes settings, data, and state
During the UCP restore, Kubernetes declarative objects are re-created, containers are re-created, and IPs are resolved.

For more information, see [Restoring an etcd cluster](https://kubernetes.io/docs/tasks/administer-cluster/configure-upgrade-etcd/#restoring-an-etcd-cluster).

## Perform UCP restore

When the restore operations starts, it looks for the UCP version used in the backup and performs one of the following actions: 

    - Fails if the restore operation is running using an image that does not match the UCP version from the backup (a `--force` flag is  available to override this if necessary)
    - Provides instructions how to run the restore process using the matching UCP version from the backup

Volumes are placed onto the host on which the UCP restore command occurs. 

The following example shows how to restore UCP from an existing backup file, presumed to be located at `/tmp/backup.tar` (replace `<UCP_VERSION>` with the version of your backup):

```
$ docker container run \
  --rm \
  --interactive \
  --name ucp \
  --volume /var/run/docker.sock:/var/run/docker.sock  \
  {{ page.ucp_org }}/{{ page.ucp_repo }}:{{ page.ucp_version }} restore < /tmp/backup.tar
```

If the backup file is encrypted with a passphrase, provide the passphrase to the restore operation(replace `<UCP_VERSION>` with the version of your backup):

```
$ docker container run \
  --rm \
  --interactive \
  --name ucp \
  --volume /var/run/docker.sock:/var/run/docker.sock  \
  {{ page.ucp_org }}/{{ page.ucp_repo }}:{{ page.ucp_version }} restore --passphrase "secret" < /tmp/backup.tar  
```

The restore command can also be invoked in interactive mode, in which case the
backup file should be mounted to the container rather than streamed through
`stdin`:

```none
$ docker container run \
  --rm \
  --interactive \
  --name ucp \
  --volume /var/run/docker.sock:/var/run/docker.sock \
  -v /tmp/backup.tar:/config/backup.tar \
  {{ page.ucp_org }}/{{ page.ucp_repo }}:{{ page.ucp_version }} restore -i
```

## Regenerate Certs
The current certs volume containing cluster specific information (such as SANs) is invalid on new clusters with different IPs. For volumes that are not backed up (`ucp-node-certs`, for example), the restore regenerates certs. For certs that are backed up, (ucp-controller-server-certs), the restore does not perform a regeneration and you must correct those certs when the restore completes.

After you successfully restore UCP, you can add new managers and workers the same way you would after a fresh installation. 

## Restore operation status
For restore operations, view the output of the restore command.

## Verify the UCP restore
A successful UCP restore involves verifying the following items:

- All swarm managers are healthy after running the following command:

```
"curl -s -k https://localhost/_ping". 
```

Alternatively, check the UCP UI **Nodes** page for node status, and monitor the UI for warning banners about unhealthy managers.

**Note**: 
- Monitor all swarm managers for at least 15 minutes to ensure no degradation.
- Ensure no containers on swarm managers are marked as "unhealthy".
- No swarm managers or nodes are running containers with the old version, except for Kubernetes Pods that use the "ucp-pause" image.

### Where to go next

- [Restore DTR](restore-dtr.md)
