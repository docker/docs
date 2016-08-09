<!--[metadata]>
+++
title = "Uninstall UCP"
description = "Learn how to uninstall a Docker Universal Control Plane cluster."
keywords = ["docker, ucp, uninstall"]
[menu.main]
parent="mn_ucp_installation"
identifier="ucp_uninstall"
weight=70
+++
<![end-metadata]-->

# Uninstall UCP

Use the docker/ucp `uninstall-cluster` command to uninstall Docker Universal Control
Plane from your swarm cluster. This command only removes the UCP services,
containers and doesn’t affect any other services or containers. Your swarm
cluster will be left intact.

To remove an individual node from UCP, and keep UCP intact, see the
documentation on joining nodes.

## Example

In this example we’ll be running the `uninstall-cluster` command interactively, so that
the command prompts for the necessary configuration values.
You can also use flags to pass values to the `uninstall-cluster` command.

1. Run the `uninstall-cluster` command.

    ```bash
    $ docker run --rm -it \
      -v /var/run/docker.sock:/var/run/docker.sock
      --name ucp \
      docker/ucp `uninstall-cluster` -i

    time="2016-08-09T20:36:36Z" level=info msg="Your engine version 1.12.0, build 8eab29e (4.4.16-boot2docker) is compatible" 
    time="2016-08-09T20:36:36Z" level=info msg="We're about to uninstall the local components for UCP ID: ZB6V:R3ZR:VMMJ:WM7B:M3US:VHMS:HZZ6:SHEL:RGXF:BHAE:2FPV:K7WH" 
    Do you want proceed with the uninstall? (y/n): y
    time="2016-08-09T20:36:38Z" level=info msg="Uninstalling UCP on each node..." 
    time="2016-08-09T20:36:56Z" level=info msg="UCP has been removed from this cluster successfully." 
    time="2016-08-09T20:36:58Z" level=info msg="Removing UCP Services" 
    ```

2. List the images remaining on the node.

    ```
    $ docker images

    REPOSITORY          TAG                 IMAGE ID            CREATED             SIZE
    docker/ucp          latest              788bdcfde423        8 days ago          8.239 MB
    ```

    The uninstall command removes all UCP-related images except the
    `docker/ucp` image.

3. Remove the docker/ucp image.

    ```
    $ docker rmi docker/ucp

    Untagged: docker/ucp:latest
    Deleted: sha256:788bdcfde423b6226b90ac98e6f233b15c0c527779177d7017a4e17db31404c9
    Deleted: sha256:dee84053b25f9b3edffb734c842a70313021063cc78d9158c63de109e1b3cb72
    Deleted: sha256:93743d5df2362466e2fe116a677ec6a4b0091bd09e889abfc9109047fcfcdebf
    ```
