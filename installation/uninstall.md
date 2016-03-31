<!--[metadata]>
+++
title = "Uninstall UCP"
description = "Learn how to uninstall a Docker Universal Control Plane cluster."
keywords = ["docker, ucp, uninstall"]
[menu.main]
parent="mn_ucp_installation"
weight=50
+++
<![end-metadata]-->


# Uninstall UCP

Use the docker/ucp uninstall command, to uninstall Docker Universal Control
Plane from a node. This command only removes the UCP containers, and doesn’t
affect any other containers.

To see what options are available in the uninstall command, check the
[uninstall command reference](../reference/uninstall.md), or run:

```bash
$ docker run --rm -it docker/ucp uninstall --help
```

To uninstall Docker UCP from a cluster, you should:

1. Uninstall UCP from every node joined in the cluster,
2. Uninstall UCP from every controller node, one at a time.

## Example

The following example illustrates an interactive uninstall.

1. Run the uninstall command.

    ```
    $ docker run --rm -it -v /var/run/docker.sock:/var/run/docker.sock --name ucp docker/ucp uninstall -i

    INFO[0000] We're about to uninstall the local components for UCP ID: FEY4:M46O:7OUS:QQA4:HLR3:4HRD:IUTH:LC2W:QPRE:BLYH:UWEM:3TYV
    Do you want proceed with the uninstall? (y/n): y

    INFO[0000] Removing UCP Containers
    INFO[0000] Removing UCP images
    INFO[0005] Removing UCP volumes
    ```

2. List the images remaining on the node.

    ```
    $ docker images

    REPOSITORY          TAG                 IMAGE ID            CREATED             SIZE
    docker/ucp          latest              788bdcfde423        8 days ago          8.239 MB
    ```

3. Remove the `ucp` image from the node.

    Because the `ucp` tool is running during the installation, it can't remove
    its own image. If you want to free up that space, delete the ucp image.

    ```
    $ docker rmi docker/ucp:latest

    Untagged: docker/ucp:latest
    Deleted: sha256:788bdcfde423b6226b90ac98e6f233b15c0c527779177d7017a4e17db31404c9
    Deleted: sha256:dee84053b25f9b3edffb734c842a70313021063cc78d9158c63de109e1b3cb72
    Deleted: sha256:93743d5df2362466e2fe116a677ec6a4b0091bd09e889abfc9109047fcfcdebf
    ```

4. Return to the UCP web application to confirm the node is removed.
