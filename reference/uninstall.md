+++
title = "uninstall"
keywords= ["uninstall, ucp"]
description = "Uninstall a UCP controller and nodes"
[menu.main]
parent = "ucp_ref"
identifier = "ucp_ref_uninstall"
+++

# docker/ucp uninstall

The uninstall command is no longer used.

## Usage

```
docker run --rm -it \
  --name ucp \
  -v /var/run/docker.sock:/var/run/docker.sock \
  docker/ucp \
  uninstall [command options]
```

## Description

The uninstall command is no longer used. To remove a node from the UCP cluster,
run the 'docker swarm leave' and 'docker node rm' commands. To remove the UCP
components but preserve the swarm cluster, use the 'uninstall-cluster' command.

## Options

| Option                | Description                                                                             |
|:----------------------|:----------------------------------------------------------------------------------------|
| `--debug`, `-D`       | Enable debug                                                                            |
| `--jsonlog`           | Produce json formatted output for easier parsing                                        |
| `--interactive`, `-i` | Enable interactive mode.,You are prompted to enter all required information             |
| `--pull`              | Specify image pull behavior ('always', when 'missing', or 'never') (default: "missing") |
| `--registry-username` | Specify the username to pull required images with [$REGISTRY_USERNAME]                  |
| `--registry-password` | Specify the password to pull required images with [$REGISTRY_PASSWORD]                  |
