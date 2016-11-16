---
description: Learn how to uninstall a Docker Universal Control Plane cluster.
keywords: docker, ucp, uninstall
redirect_from:
- /ucp/installation/uninstall/
title: Uninstall UCP
---

Use the docker/ucp uninstall command, to uninstall Docker Universal Control
Plane from a node. This command only removes the UCP containers, and doesn’t
affect any other containers.

To see what options are available in the uninstall command, check the
[uninstall command reference](../reference/uninstall.md).

To uninstall Docker UCP from a cluster, you need to:

1. Uninstall UCP from every node joined in the cluster,
2. Uninstall UCP from every controller node, one at a time,
3. Restart the Docker engine on all the nodes.


## Example

In this example we’ll be running the uninstall command interactively, so that
the command prompts for the necessary configuration values.
You can also use flags to pass values to the uninstall command.

1.  Run the uninstall command.

    ```none
    $ docker run --rm -it \
      -v /var/run/docker.sock:/var/run/docker.sock \
      --name ucp \
      docker/ucp uninstall -i

    INFO[0000] Were about to uninstall the local components for UCP ID: FEY4:M46O:7OUS:QQA4:HLR3:4HRD:IUTH:LC2W:QPRE:BLYH:UWEM:3TYV
    Do you want proceed with the uninstall? (y/n): y

    WARN[0000] We detected a daemon advertisement configuration. Proceed with caution, as the daemon will require a restart. Press ctrl-c to cancel uninstall within 4 seconds.
    INFO[0004] Removing UCP Containers
    INFO[0005] Removing UCP images
    WARN[0006] Configuration updated. You will have to manually restart the docker daemon for the changes to take effect.
    WARN[0006] Engine discovery configuration removed. You will need to restart the daemon.
    INFO[0010] Removing UCP volumes
    ```

2.  List the images remaining on the node.

    ```none
    $ docker images

    REPOSITORY          TAG                 IMAGE ID            CREATED             SIZE
    docker/ucp          latest              788bdcfde423        8 days ago          8.239 MB
    ```

    The uninstall command removes all UCP-related images except the
    `docker/ucp` image.

3.  Remove the docker/ucp image.

    ```none
    $ docker rmi docker/ucp

    Untagged: docker/ucp:latest
    Deleted: sha256:788bdcfde423b6226b90ac98e6f233b15c0c527779177d7017a4e17db31404c9
    Deleted: sha256:dee84053b25f9b3edffb734c842a70313021063cc78d9158c63de109e1b3cb72
    Deleted: sha256:93743d5df2362466e2fe116a677ec6a4b0091bd09e889abfc9109047fcfcdebf
    ```

4.  Restart the Docker daemon.

    When you install or join a node, UCP configures the Docker engine on that
    node for multi-host networking. When uninstalling, the configuration is
    reverted to its original state, but you need to restart the Docker engine
    for the configurations to take effect.

    As an example, to restart the Docker engine on a Ubuntu distribution:

    ```bash
    $ sudo service docker restart
    ```

5. Confirm the node was removed from the cluster.

    In the UCP web application, confirm the node is no longer listed. It
    might take a few minutes for UCP to stop listing that node.