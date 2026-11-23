---
title: Work with project files through API v2
description: Choose local workspace mounts or cloud file transfers for sandbox projects.
keywords: docker sandboxes, API, SDK, local sandboxes, cloud sandboxes
weight: 30
linkTitle: Work with project files
---

How you make project files available depends on where the sandbox runs. Local
sandbox creation can mount directories from the machine running the daemon.
Cloud sandboxes need files transferred to their remote environment.

A host path belongs to the daemon's machine. It is not a path on the machine
running an arbitrary API client, and it does not identify files to upload to
cloud.

The cloud API supports file-content transfers and directory archives. These
operations do not imply that the API exposes every filesystem operation
available through the CLI or inside a sandbox.

<!-- Add separate local mount and cloud transfer examples after validation.
Verify mount roles, read-only access, paths, archive semantics, and transfer
limits. Do not imply that file-transfer routes are served by the local host. -->
