+++
title = "uninstall"
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

Removes UCP from the controller and the nodes. The `uninstall` | does not remove any other containers that are running, except those recognized to be part of UCP.

## Options

| Option | Description |
|-----------------------|------------------------------------------------------------------------------|
| `--debug`, `-D` | Enable debug. |
| `--jsonlog` | Produce json formatted output for easier parsing. |
| `--interactive`, `-i` | Enable interactive mode.,You are prompted to enter all required information. |
| `--id` | The ID of the UCP instance to uninstall. |
| `--preserve-certs` | Don't delete the certs on the host. |
| `--preserve-images` | Don't delete images on the host. |
