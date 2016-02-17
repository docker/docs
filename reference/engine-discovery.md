+++
title = "engine-discovery"
description = "description"
[menu.main]
parent = "ucp_ref"
+++

# engine-discovery

Manage the Engine discovery configuration on a node.

## Usage

```
docker run --rm -it \
     --name ucp \
     -v /var/run/docker.sock:/var/run/docker.sock \
     docker/ucp \
     engine-discovery [options]
```

## Description

Use this command to display and update Engine discovery configuration on a node.
The discovery configuration is used by Engine for cluster membership and
multi-host networking.

Use one or more '--controller' arguments to specify *all* of the
UCP controllers in this cluster.

The '--host-address' argument specifies the public advertise address for the
particular node  you are running the command on. This host-address is how other
nodes in UCP talk to this node.  You may specify an IP or hostname, and the
command automatically detects and fills in the port number.  If you omit the
address, the tool attempts to discover the node's address.

## Options

| Option                    | Description                                                                      |
|---------------------------|----------------------------------------------------------------------------------|
| `--debug`, `-D`           | Enable debug.                                                                    |
| `--jsonlog`               | Produce json formatted output for easier parsing.                                |
| `--interactive`, `-i`     | Enable interactive mode. You are prompted to enter all required information. |
| `--list`                  | Display the Engine discovery configuration                              |
| `--controller [--controller option --controller option]`                  | Update discovery with one or more controller's external IP address or hostname.                            |
| `--host-address`                  | Update the external IP address or hostname this node advertises itself as [`$UCP_HOST_ADDRESS`].                        |
