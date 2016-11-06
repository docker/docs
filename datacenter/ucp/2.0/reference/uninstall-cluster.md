+++
title = "uninstall-cluster"
keywords= ["uninstall, ucp"]
description = "Uninstall UCP from this swarm clusters"
[menu.main]
parent = "ucp_ref"
identifier = "ucp_ref_uninstall-cluster"
+++

# docker/ucp uninstall-cluster

Uninstall UCP from this swarm cluster.

## Usage

```
docker run --rm -it \
           --name ucp \
           -v /var/run/docker.sock:/var/run/docker.sock \
           docker/ucp \
           uninstall-cluster [command options]
```

## Description

Uninstall UCP from this swarm cluster, while preserving the cluster.

## Options
```nohighlight
--debug, -D                Enable debug mode
--jsonlog                  Produce json formatted output for easier parsing
--interactive, -i          Enable interactive mode.  You will be prompted to enter all required information
--pull value               Specify image pull behavior ('always', when 'missing', or 'never') (default: "missing")
--registry-username value  Specify the username to pull required images with [$REGISTRY_USERNAME]
--registry-password value  Specify the password to pull required images with [$REGISTRY_PASSWORD]
--id value                 The ID of the UCP instance to uninstall
```
