+++
title = "uninstall"
keywords= ["uninstall, ucp"]
description = "Uninstall a UCP controller and nodes"
[menu.main]
identifier = "ucp_uninstall"
parent = "ucp_ref"
+++

# uninstall

```
   docker run --rm -it \
        --name ucp \
        -v /var/run/docker.sock:/var/run/docker.sock \
        docker/ucp \
        uninstall [command options]
```

## Description

Removes UCP from a node. The `uninstall` | does not remove any other containers that are running, except those recognized to be part of UCP.

After you uninstall UCP from a node, the node continues to appear in the Dashboard because the node still has the `ucp` image; the image couldn't be removed while a container was running. Remove the `ucp` tool's image to completely delete the node from the UCP application dashboard.

## Options

| Option | Description |
|-----------------------|------------------------------------------------------------------------------|
| `--debug`, `-D` | Enable debug. |
| `--jsonlog` | Produce json formatted output for easier parsing. |
| `--interactive`, `-i` | Enable interactive mode.,You are prompted to enter all required information. |
| `--id` | The ID of the UCP instance to uninstall. |
| `--preserve-certs` | Don't delete the certs on the host. |
| `--preserve-images` | Don't delete images on the host. |

## Example

The following example illustrates an interactive uninstall.

1. Run the uninstall command.

        $ docker run --rm -it -v /var/run/docker.sock:/var/run/docker.sock --name ucp docker/ucp uninstall -i
        INFO[0000] We're about to uninstall the local components for UCP ID: FEY4:M46O:7OUS:QQA4:HLR3:4HRD:IUTH:LC2W:QPRE:BLYH:UWEM:3TYV
        Do you want proceed with the uninstall? (y/n): y
        INFO[0000] Removing UCP Containers                      
        INFO[0000] Removing UCP images                          
        INFO[0005] Removing UCP volumes                         

2. List the images remaining on the node.

    Because the `ucp` tool is running during the installation, it can't remove
    its own image. Until you remove that image, the node continues to appear in
    the UCP dashboard.

        $ docker images
        REPOSITORY          TAG                 IMAGE ID            CREATED             SIZE
        docker/ucp          latest              788bdcfde423        8 days ago          8.239 MB

3. Remove the `ucp` image from the node.

        $ docker rmi 788bdcfde423
        Untagged: docker/ucp:latest
        Deleted: sha256:788bdcfde423b6226b90ac98e6f233b15c0c527779177d7017a4e17db31404c9
        Deleted: sha256:dee84053b25f9b3edffb734c842a70313021063cc78d9158c63de109e1b3cb72
        Deleted: sha256:93743d5df2362466e2fe116a677ec6a4b0091bd09e889abfc9109047fcfcdebf

4. Return to the UCP application, confirm the node is removed.
