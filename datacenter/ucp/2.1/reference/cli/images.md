---
title: docker/ucp images
description: Verify the UCP images on this node
keywords: docker, ucp, cli, images
---

Verify the UCP images on this node

## Description

This command checks the UCP images that are available in this node, and pulls
the ones that are missing.


## Options

| Option                    | Description                |
|:--------------------------|:---------------------------|
|`--debug, D`|Enable debug mode|
|`--jsonlog`|Produce json formatted output for easier parsing|
|`--pull`|Pull UCP images: `always`, when `missing`, or `never`|
|`--registry-username`|Username to use when pulling images|
|`--registry-password`|Password to use when pulling images|
|`--list`|List all images used by UCP but don't pull them|
