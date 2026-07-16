---
title: Connect Claude Desktop to a sandbox
linkTitle: Claude Desktop
weight: 30
description: Run Claude Code from Claude Desktop against a Docker Sandbox over SSH.
keywords: docker sandboxes, claude desktop, claude code, remote ssh, sbx
---

{{< summary-bar feature_name="Docker Sandboxes SSH" >}}

Claude Desktop can run Claude Code on a remote machine over SSH. Point it at a
sandbox so the agent works inside the isolated environment instead of on your
host.

> [!NOTE]
> This page covers Claude Desktop connecting to a sandbox over SSH. To run the
> Claude Code CLI inside a sandbox directly, see
> [Claude Code](../agents/claude-code.md).

## Prerequisites

- SSH access set up. See [Editor and app integrations](_index.md#enable-ssh-access).
- Claude Desktop installed.

## Connect

1. Open the environment picker and select **Add SSH host**.
2. In the **Add SSH connection** dialog, enter the connection details:
   - **Name**: a friendly label for the connection, such as `My sandbox`.
   - **SSH Host**: enter the sandbox host manually as `<name>.sbx`.
   - Leave **SSH Port** and **Identity File** empty. The managed SSH config
     that `sbx setup ssh` wrote supplies everything else.
3. Select **Add SSH connection**. Claude Desktop connects to the sandbox and
   runs Claude Code there.

<!-- TODO: add screenshots of the environment picker and the Add SSH connection dialog -->

## Related

- [Editor and app integrations](_index.md) — how SSH access works and how to
  set it up
- [Claude Code](../agents/claude-code.md) — run the Claude Code CLI inside a
  sandbox
