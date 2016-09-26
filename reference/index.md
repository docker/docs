<!--[metadata]>
+++
title = "UCP tool reference"
description = "Installs Docker Universal Control Plane."
keywords= ["docker, ucp, install"]
[menu.main]
parent = "mn_ucp_installation"
identifier = "ucp_ref"
weight=100
+++
<![end-metadata]-->

# UCP tool reference

Installs Docker Universal Control Plane.

## Usage

```bash
$ docker run --rm -it \
    --name ucp \
    -v /var/run/docker.sock:/var/run/docker.sock \
    docker/ucp \
    command [arguments...]
```

## Description

This image has commands to install and manage Docker Universal Control Plane
(UCP) on a Docker Engine.

You can configure the commands using flags or environment variables. When using
environment variables use the `docker run -e VARIABLE_NAME` syntax to pass the
value from your shell, or `docker run -e VARIABLE_NAME=value` to specify the
value explicitly on the command line.

The container running this image needs to be named 'ucp' and bind-mount the
Docker daemon socket. Below you can find an example of how to run this image.

Additional help is available for each command with the '--help' flag.

## commands

```nohighlight
install            Install UCP on this engine
restart            Start or restart UCP components on this engine
stop               Stop UCP components running on this engine
upgrade            Upgrade the UCP components on this engine
images             Verify the UCP images on this engine
uninstall-cluster  Uninstall UCP from this swarm cluster
dump-certs         Dump out the public certs for this UCP controller
fingerprint        Dump out the TLS fingerprint for the UCP controller running on this engine
support            Generate a support dump for this engine
id                 Dump out the ID of the UCP components running on this engine
backup             Stream a tar file to stdout containing all UCP data volumes
restore            Restore a UCP cluster from a backup tar file.
regen-certs
help               Shows a list of commands or help for one command
```

## Global options

```nohighlight
--help, -h     show help
--version, -v  print the version
```
