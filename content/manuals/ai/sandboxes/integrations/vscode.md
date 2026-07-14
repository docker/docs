---
title: Connect VS Code to a sandbox
linkTitle: VS Code
weight: 10
description: Use VS Code Remote - SSH to develop inside a Docker Sandbox.
keywords: docker sandboxes, vs code, remote ssh, remote development, sbx
---

{{< summary-bar feature_name="Docker Sandboxes SSH" >}}

Use the Remote - SSH extension to open a VS Code window that runs inside a
sandbox. Your editor stays on your host while files, terminals, and extensions
run in the isolated sandbox.

## Prerequisites

- SSH access set up. See [Editor and app integrations](_index.md#enable-ssh-access).
- The [Remote - SSH](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.remote-ssh)
  extension (`ms-vscode-remote.remote-ssh`) installed in VS Code.

## Connect

1. Open the Command Palette and run **Remote-SSH: Connect to Host**.
2. Select the sandbox host, `<name>.sbx`. Because `sbx setup ssh` added it to
   your SSH config, it appears in the host list automatically. You can also
   type it in directly.
3. VS Code opens a new window connected to the sandbox. Open a folder from the
   sandbox workspace to start working.

<!-- TODO: add screenshot of the Remote-SSH host picker showing <name>.sbx -->

## Notes

- The first connection installs the VS Code server inside the sandbox, so it
  can take a moment. Later connections are faster.
- Port forwarding is limited to the sandbox's loopback address. To reach a
  sandbox port from your host, use
  [`sbx ports`](../usage.md#accessing-services-in-the-sandbox).

## Related

- [Editor and app integrations](_index.md) — how SSH access works and how to
  set it up
