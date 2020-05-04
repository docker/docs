---
title: Back up Docker Swarm
description: Learn how to create a backup of Docker Swarm
keywords: enterprise, backup, swarm
---

Docker manager nodes store the swarm state and manager logs in the `/var/lib/docker/swarm/` directory. Swarm raft logs contain crucial information for re-creating Swarm specific resources, including services, secrets, configurations and node cryptographic identity. In 1.13 and higher, this data includes the keys used to encrypt the raft logs. Without these keys, you cannot restore the swarm. 

You must perform a manual backup on each manager node, because logs contain node IP address information and are not transferable to other nodes. If you do not backup the raft logs, you cannot verify workloads or Swarm resource provisioning after restoring the cluster.

> You can avoid performing Swarm backup by storing stacks, services definitions, secrets, and networks definitions in a *Source Code Management* or *Config Management* tool.

## Swarm backup contents

| Data               | Description                                                                          | Backed up |
| :------------------|:-------------------------------------------------------------------------------------|:----------|
| Raft keys          | Used to encrypt communication among Swarm nodes and to encrypt and decrypt Raft logs | yes
| Membership         | List of the nodes in the cluster                                                     | yes
| Services           | Stacks and services stored in Swarm-mode                                             | yes
| Networks (overlay) | The overlay networks created on the cluster                                          | yes
| Configs            | The configs created in the cluster                                                   | yes
| Secrets            | Secrets saved in the cluster                                                         | yes
| Swarm unlock key   | **Must be saved on a password manager !**                                            | no  

## Procedure

1.  Retrieve your Swarm unlock key if `auto-lock` is enabled to be able
    to restore the swarm from backup. Retrieve the unlock key if necessary and
    store it in a safe location. If you are unsure, read
    [Lock your swarm to protect its encryption key](../../../engine/swarm/swarm_manager_locking.md).

2.  Because you must stop the engine of the manager node before performing the backup, having three manager 
    nodes is recommended for high availability (HA). For a cluster to be operational, a majority of managers 
    must be online. If less than 3 managers exists, the cluster is unavailable during the backup.

    > **Note**: During the time that a manager is shut down, your swarm is more vulnerable to
    > losing the quorum if further nodes are lost. A loss of quorum means that the swarm is unavailabile 
    > until quorum is recovered. Quorum is only recovered when more than 50% of the nodes are again available. 
    > If you regularly take down managers to do backups, consider running a 5-manager swarm, so that you 
    > can lose an additional manager while the backup is running without disrupting services. 

3.  Select a manager node. Try not to select the leader in order to avoid a new election inside the cluster:

    ```
    docker node ls -f "role=manager" | tail -n+2 | grep -vi leader
    ```
>  Optional: Store the Docker version in a variable for easy addition to your backup name.
   
{% raw %}
    ```
    ENGINE=$(docker version -f '{{.Server.Version}}')
    ```
{% endraw %}

4.  Stop the Docker Engine on the manager before backing up the data, so that no data is changed during the backup:

    ```
    systemctl stop docker
    ```

5. Back up the entire `/var/lib/docker/swarm` folder:

    ```
    tar cvzf "/tmp/swarm-${ENGINE}-$(hostname -s)-$(date +%s%z).tgz" /var/lib/docker/swarm/
    ```

   Note: _You can decode the Unix epoch in the filename by typing `date -d @timestamp`._  For example:
   
   ```
   date -d @1531166143
   Mon Jul  9 19:55:43 UTC 2018
   ```

6.  Restart the manager Docker Engine:

    ```
    systemctl start docker
    ```
7. Except for step 1, repeat the previous steps for each manager node.

### Where to go next

- [Back up UCP](back-up-ucp.md)
 
