---
title: Restore Docker Swarm
description: Learn how to restore Docker Swarm from an existing backup
keywords: enterprise, restore, swarm
---

## Prerequisites

-   You must use the same IP as the node from which you made the backup. The command to force the new cluster does not reset the IP in the Swarm data.  
-   You must restore the backup on the same Docker Engine version.
-   You can find the list of manager IP addresses in `state.json` in the zip file.
-   If `auto-lock` was enabled on the old Swarm, the unlock key is required to perform the restore.

## Perform Swarm restore
Use the following procedure on each manager node to restore data to a new swarm.

1. Shut down the Docker Engine on the node you select for the restore:

    ```
    systemctl stop docker
    ``` 
2. Remove the contents of the `/var/lib/docker/swarm` directory on the new Swarm if it exists.
3. Restore the `/var/lib/docker/swarm` directory with the contents of the backup.  
    
    > **Note**: The new node uses the same encryption key for on-disk
    > storage as the old one. It is not possible to change the on-disk storage
    > encryption keys at this time. In the case of a swarm with auto-lock enabled, 
    > the unlock key is also the same as on the old swarm, and the unlock key is 
    > needed to restore the swarm.

4. Start Docker on the new node.  Unlock the swarm if necessary. 

    ```
    systemctl start docker
    ```
5. Re-initialize the swarm so that the node does not attempt to connect to nodes that were part of the old swarm, and presumably no longer exist:

    ```
    $ docker swarm init --force-new-cluster
    ```  
    
6.  Verify that the state of the swarm is as expected. This may include
    application-specific tests or simply checking the output of
    `docker service ls` to be sure that all expected services are present.

7.  If you use auto-lock,
    [rotate the unlock key](/engine/swarm/swarm_manager_locking.md#rotate-the-unlock-key).
8.  Add the manager and worker nodes to the new swarm.
9.  Reinstate your previous backup regimen on the new swarm.

### Where to go next

- [Restore UCP](restore-ucp.md)
