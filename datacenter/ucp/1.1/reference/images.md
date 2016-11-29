---
description: Verify the UCP images on this Docker Engine.
keywords: images, ucp, images
title: docker/ucp images
---

Verify the UCP images on this Docker Engine.

## Usage

```
docker run --rm -it \
  --name ucp \
  -v /var/run/docker.sock:/var/run/docker.sock \
  docker/ucp \
  images [command options]
```

## Description

This command will verify all the required images used by UCP on the current
engine.
By default, this will pull any missing images. Use the '--pull' argument
to change behavior.

## Options

| Option                | Description                                                                             |
|:----------------------|:----------------------------------------------------------------------------------------|
| `--debug, -D`         | Enable debug                                                                            |
| `--jsonlog`           | Produce json formatted output for easier parsing                                        |
| `--registry-username` | Specify the username to pull required images with [$REGISTRY_USERNAME]                  |
| `--registry-password` | Specify the password to pull required images with [$REGISTRY_PASSWORD]                  |
| `--pull "missing"`    | Specify image pull behavior ('always', when 'missing', or 'never') (default: "missing") |
| `--list`              | Don`t do anything, just list the images used by UCP                                     |