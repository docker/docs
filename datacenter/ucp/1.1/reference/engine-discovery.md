---
description: Manage the engine discovery configuration on a node.
keywords: docker, ucp, discovery
title: docker/ucp engine-discovery
---

Manage the Engine discovery configuration.

## Usage

```
docker run --rm -it \
  --name ucp \
  -v /var/run/docker.sock:/var/run/docker.sock \
  docker/ucp \
  engine-discovery [options]
```

## Description

This command will display and update the local Docker engine discovery
configuration. This is used by the engine for cluster membership and multi-host
networking.

By default, this command will check if the configuration is up to date. Use
'--update' to update the configuration. Use '--debug' to show more information
including current and proposed configuration.

UCP controllers in this cluster are auto-detected if possible. Alternately you
may use one or more '--controller' arguments to manually specify ALL of the UCP
controllers.

The '--host-address' argument specifies the public advertise address for THIS
node (how other nodes in the system talk to this node.)  You may specify an IP
address, and the port number will be automatically filled in. If omitted, the
tool will attempt to discover the address of this node.

This command uses the exit status of 0 for success. An exit status of 1 is used
when run without the '--update' flag and when the configuration needs updating,
and 2 is used for any failures.

## Options

| Option                                                     | Description                                                                                   |
|:-----------------------------------------------------------|:----------------------------------------------------------------------------------------------|
| `--debug, -D`                                              | Enable debug                                                                                  |
| `--jsonlog`                                                | Produce json formatted output for easier parsing                                              |
| `--update`                                                 | Apply engine discovery configuration changes                                                  |
| `--controller` `[--controller option --controller option]` | Update discovery with the external IP address or hostname of the controller(s)                |
| `--host-address`                                           | Update the external IP address or hostname this node advertises itself as [$UCP_HOST_ADDRESS] |