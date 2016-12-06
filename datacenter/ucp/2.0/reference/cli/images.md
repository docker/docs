---
description: Verify the UCP images on this node
keywords: docker, dtr, cli, images
title: docker/ucp images
---

Verify the UCP images on this node

## Usage

```bash

docker run -it --rm \
    --name ucp \
    -v /var/run/docker.sock:/var/run/docker.sock \
    docker/ucp \
    images [command options]

```

## Description

This command checks the UCP images that are available in this node, and pulls
the ones that are missing.
  

## Options

| Option                    | Description                |
|:--------------------------|:---------------------------|
|`--debug, D`|Enable debug mode|
|`--jsonlog`|Produce json formatted output for easier parsing|
|`--pull`|Pull UCP images: 'always', when 'missing', or 'never'|
|`--registry-username`|Username to use when pulling images|
|`--registry-password`|Password to use when pulling images|
|`--list`|List all images used by UCP but don't pull them|
