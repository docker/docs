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
- Cursor's Remote - SSH support installed. Cursor bundles a Remote - SSH
  extension compatible with the VS Code one.

## Connect

1. Open the Command Palette and run **Remote-SSH: Connect to Host**.
2. Enter the sandbox host manually as `<name>.sbx`.
3. Cursor opens a new window connected to the sandbox. Open a folder from the
   sandbox workspace to start working.

<!-- TODO: add screenshot of the Cursor Remote-SSH host picker showing <name>.sbx -->

## Notes

- The first connection installs the editor server inside the sandbox, so it
  can take a moment. Later connections are faster.
- SSH port forwarding can reach services on the sandbox's loopback address.
  Other destination addresses aren't supported.

## Related

- [Editor and app integrations](_index.md) — how SSH access works and how to
  set it up
- [Cursor agent](../agents/cursor.md) — run the Cursor CLI inside a sandbox
