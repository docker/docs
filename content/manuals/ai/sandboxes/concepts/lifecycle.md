---
title: Lifecycle and persistence
linkTitle: Lifecycle and persistence
description: Understand sandbox execution, retained state, connections, expiration, and compute billing.
keywords: docker sandboxes, sbx, cloud sandboxes
weight: 30
---

A sandbox's lifetime is separate from the terminal or browser session you use
to reach it. Closing a connection does not necessarily stop the sandbox.
Stopping a sandbox and deleting it have different effects on its state.

## Execution and memory

A stopped sandbox does not execute work. This does not, by itself, tell you
whether its process memory is retained. A service that preserves memory can
hibernate the sandbox and resume that state later.

The cloud API makes memory preservation dependent on the backend's supported
stop guarantee. Do not assume that every sandbox can resume process memory
because its interface offers a stop or pause action. See the
[API reference](/reference/api/sandboxes/latest/) for the stop and start contract.

For local sandboxes, stopping retains filesystem changes. Follow the
[local lifecycle instructions](/manuals/ai/sandboxes/cli/manage/local.md#start-stop-and-remove) to start
the agent again.

## Disk and other stored data

Files in a sandbox's own filesystem persist across supported stops and starts.
Deleting the sandbox removes that filesystem. Files in a mounted local
workspace live on the host and remain after sandbox deletion.

Snapshots, template images, and volumes have their own lifetimes and capture
rules. A local template does not capture mounted workspaces or the mounted
Docker store. Cloud volume data is saved at sandbox exit, rather than
continuously. See [Files, workspaces, and storage](/manuals/ai/sandboxes/concepts/files.md).

## Connections

A terminal, process stream, SSH connection, or published service connection
can end independently of the sandbox. After stopping and starting a sandbox,
reconnect through your interface. SDK clients must not assume that an earlier
process connection remains usable.

Detaching from an agent session leaves it running when the interface supports
detachment. It does not pause the cloud expiration timer.

## Cloud expiration

Cloud sandboxes have a time-to-live and an action to take at expiration. That
action can delete state or stop the sandbox, depending on the configuration
and supported capabilities. Set it deliberately when retaining work matters.

A volume-backed cloud sandbox cannot be stopped through the documented CLI
workflow. It requires deletion at expiration, which saves the volume snapshot.

Inspect the expiration after resuming a cloud sandbox. Follow your interface's
instructions for configuring or extending its timer:

- [CLI expiration](/manuals/ai/sandboxes/cli/manage/cloud.md#configure-expiration)
- [Console lifecycle](/manuals/ai/sandboxes/console/sandboxes.md#manage-the-lifecycle)
- [SDK stop and restart](/manuals/ai/sandboxes/api/cookbook/stop-and-restart-a-sandbox.md)

## Billing and quotas

Cloud compute is billed while the sandbox runs. Stopped sandboxes still count
toward stored sandbox quotas. Stopping an ordinary sandbox releases its
concurrency slot; an always-on sandbox retains its reservation.

Your model provider bills inference separately. See
[cloud access and billing](/manuals/ai/sandboxes/cloud-access/_index.md) and
[compute sizes and quotas](/manuals/ai/sandboxes/cloud-access/resources.md).
