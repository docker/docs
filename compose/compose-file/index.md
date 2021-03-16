---
description: Compose file reference
keywords: fig, composition, compose, docker
redirect_from:
- /compose/yml
- /compose/compose-file/compose-file-v1/
title: Compose file
toc_max: 4
toc_min: 1
---

## Reference and guidelines

These topics describe the Docker Compose implementation of the Compose format.
Docker Compose **1.27.0+** implements the format defined by the [Compose Specification](https://github.com/compose-spec/compose-spec/blob/master/spec.md). Previous Docker Compose versions have support for several Compose file formats – 2, 2.x, and 3.x. The Compose specification is an unified 2.x and 3.x file format, aggregating properties across these formats.

## Compose and Docker compatibility matrix

There are several versions of the Compose file format – 1, 2, 2.x, and 3.x. The
table below provides a snapshot of various versions. For full details on what each version includes and
how to upgrade, see **[About versions and upgrading](compose-versioning.md)**.

{% include content/compose-matrix.md %}

## Compose documentation

- [User guide](../index.md)
- [Installing Compose](../install.md)
- [Compose file versions and upgrading](compose-versioning.md)
- [Sample apps with Compose](../samples-for-compose.md)
- [Enabling GPU access with Compose](../gpu-support.md)
- [Command line reference](../reference/index.md)
