---
title: Connect Cursor to a sandbox
linkTitle: Cursor
weight: 20
description: Use Cursor's Remote - SSH support to develop inside a Docker Sandbox.
keywords: docker sandboxes, cursor, remote ssh, remote development, sbx
---

{{< summary-bar feature_name="Docker Sandboxes SSH" >}}

Cursor is built on VS Code, so it connects to a sandbox the same way, using
Remote - SSH. Your editor stays on your host while files, terminals, and
extensions run in the isolated sandbox.

> [!NOTE]
> This page covers the Cursor editor connecting to a sandbox over SSH. To run
> the Cursor agent CLI inside a sandbox instead, see
> [Cursor agent](../agents/cursor.md).

## Prerequisites

- SSH access set up. See [Editor and app integrations](_index.md#enable-ssh-access).
- Cursor's Remote - SSH support installed.

## Connect

Confirm that you can connect to the sandbox from a terminal:

```console
$ ssh demo.sbx
```

1. Open the Command Palette and run **Remote-SSH: Connect to Host**.
2. Enter the sandbox host manually as `<name>.sbx`.
3. Cursor opens a new window connected to the sandbox. Use the remote folder
   picker to [select the mounted workspace](_index.md#select-the-workspace-folder).

## Notes

- The first connection installs the editor server inside the sandbox, so it
  can take a moment. Later connections are faster.

## Related

- [Editor and app integrations](_index.md) — how SSH access works and how to
  set it up
- [Cursor agent](../agents/cursor.md) — run the Cursor CLI inside a sandbox
