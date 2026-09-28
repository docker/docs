---
title: Files, workspaces, and storage
linkTitle: Files and storage
description: Choose how files reach a sandbox and which data survives its removal.
keywords: docker sandboxes, sbx, cloud sandboxes
weight: 20
---

A sandbox has its own filesystem. Decide which files belong there and which
must remain available independently of that sandbox.

## Local workspaces

A local sandbox can work without a host mount, mount a host directory directly,
or use clone mode for a Git repository. Direct mounts expose the host working
tree, so edits affect your host files. Clone mode gives the agent a private
clone and leaves the source repository mounted read-only.

Deleting the sandbox removes files stored only inside it. Host files remain.
Fetch or push work from a private clone before deleting its sandbox. See
[CLI workspace modes](/manuals/ai/sandboxes/cli/manage/local.md#git-workspace-modes) and
[Git workflows](/manuals/ai/sandboxes/cli/workflows/git.md).

## Cloud files

Cloud sandboxes do not mount your local directories. Copy files into a sandbox
or clone a repository there. A copy is a point-in-time transfer, not ongoing
synchronization. Push commits or copy results out before deleting the sandbox.

Use [CLI file transfers](/manuals/ai/sandboxes/cli/manage/cloud.md#transfer-files),
[Console files](/manuals/ai/sandboxes/console/sandboxes.md#source-code-and-files), or the
[SDK cookbook](/manuals/ai/sandboxes/api/cookbook/_index.md).

## Captured state and volumes

A kit defines a reproducible environment. A template or snapshot captures state
from an environment. These serve different purposes: a capture can contain
changes and credential files that were never part of a kit's source.

Local templates capture the container filesystem, excluding mounted
workspaces and the Docker store. Cloud images and snapshots follow their
service's capture rules. Check whether memory is included before depending on
process resumption.

Experimental cloud volumes store data independently of one sandbox. The CLI
workflow saves their data at sandbox exit. Concurrent writers can overwrite
each other's snapshots. See [Cloud volumes](/manuals/ai/sandboxes/cli/manage/cloud.md#use-persistent-volumes)
before using a volume for shared data.

See [Lifecycle and persistence](/manuals/ai/sandboxes/concepts/lifecycle.md) for how stopping, resuming, and
deleting affect these resources.
