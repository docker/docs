---
description: Uninstall a UCP controller and nodes
keywords: uninstall, ucp
redirect_from:
- /ucp/reference/uninstall/
title: docker/ucp uninstall
---

Uninstall UCP components from this Docker Engine.

## Usage

```
docker run --rm -it \
  --name ucp \
  -v /var/run/docker.sock:/var/run/docker.sock \
  docker/ucp \
  uninstall [command options]
```

## Description

When uninstalling UCP, you must run the 'uninstall' command against every
engine in your cluster.

## Options

| Option                | Description                                                                             |
|:----------------------|:----------------------------------------------------------------------------------------|
| `--debug`, `-D`       | Enable debug                                                                            |
| `--jsonlog`           | Produce json formatted output for easier parsing                                        |
| `--interactive`, `-i` | Enable interactive mode.,You are prompted to enter all required information             |
| `--pull`              | Specify image pull behavior ('always', when 'missing', or 'never') (default: "missing") |
| `--registry-username` | Specify the username to pull required images with [$REGISTRY_USERNAME]                  |
| `--registry-password` | Specify the password to pull required images with [$REGISTRY_PASSWORD]                  |
| `--id`                | The ID of the UCP instance to uninstall                                                 |
| `--preserve-certs`    | Don't delete the certs on the host                                                      |
| `--preserve-images`   | Don't delete images on the host                                                         |
