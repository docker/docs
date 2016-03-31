<!--[metadata]>
+++
title = "UCP tool reference"
keywords= ["tool, reference, ucp"]
description = "Run UCP commands"
[menu.main]
identifier = "ucp_ref"
parent = "mn_ucp_installation"
weight=100
+++
<![end-metadata]-->

# ucp tool Reference


```
docker run --rm -it \
    --name ucp \
    -v /var/run/docker.sock:/var/run/docker.sock \
    docker/ucp \
     command [arguments...]
```

The UCP installation consists of using the Docker Engine CLI to run the `ucp`
tool. The `ucp` tool is an image with subcommands to `install` a controller or
`join` a node to a UCP controller. The general format of these commands are:

| Docker client | `run` command with options | `ucp` image  | Subcommand with options |
|:--------------|:---------------------------|:-------------|:------------------------|
| `docker`      | `run --rm -it`             | `docker/ucp` | `install --help`        |
| `docker`      | `run --rm -it`             | `docker/ucp` | `join --help`           |
| `docker`      | `run --rm -it`             | `docker/ucp` | `uninstall --help`      |

You can these two subcommands interactively by passing them the `-i`
option or by passing command-line options. This installation guide's steps
assume both are run interactively.

To list all the possible subcommands, use:

```
$ docker run --rm -it docker/ucp  --help
```


## Description

This tool has commands to 'install' the UCP initial controller and
'join' nodes to that controller.  The tool can also 'uninstall' the product.
This tool must run as a container with a well-known name and with the
docker.sock volume mounted, which you can cut-and-paste from the usage
example below.

This tool generates TLS certificates and attempts to discover the local
systems's hostname and primary IP addresses.  The tool may be unable to discover
your externally visible fully qualified hostname.  You can use  the
'--host-address' option to specify a hostname or primary IP address
specifically.

For proper certificate verification, you should pass one or more subject
alternative names (SANs) with '--san' during 'install' and 'join' that matches
the fully qualified hostname you intend to use to access the given system.

Additional help is available for each command with the '--help' option.

## Options
`--help`, `-h` Show help
`--version`, `-v`	Print the version

## Subcommands

| Command                         | Description                                                                 |
|:--------------------------------|:----------------------------------------------------------------------------|
| [`install`](install.md)         | Install UCP on this engine.                                                 |
| [`join`](join.md)               | Join this engine to an existing UCP.                                        |
| [`upgrade`](upgrade.md)         | Upgrade the UCP components on this Engine.                                  |
| [`images`](images.md)           | Verify the UCP images on this Engine.                                       |
| [`uninstall`](uninstall.md)     | Uninstall UCP components from this Engine.                                  |
| [`dump-certs`](dump-certs.md)   | Dump out the public certs for this UCP controller.                          |
| [`fingerprint`](fingerprint.md) | Dump out the TLS fingerprint for the UCP controller running on this Engine. |
| [`help`](help.md)               | Shows a list of commands or help for one command.                           |
