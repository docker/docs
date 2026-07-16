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

Confirm that you can connect to the sandbox from a terminal:

```console
$ ssh demo.sbx
```

In Claude Desktop, open the environment drop-down before starting a session and
select **+ Add SSH connection**. Enter a name for the connection and enter the
sandbox hostname, such as `demo.sbx`, in **SSH Host**. Leave **SSH Port** and
**Identity File** empty because the managed SSH config supplies them.

Select the connection from the environment drop-down and choose the sandbox
workspace folder.

For more connection options, see the Claude Desktop instructions for
[SSH sessions](https://code.claude.com/docs/en/desktop#ssh-sessions).

## Related

- [Editor and app integrations](_index.md) — how SSH access works and how to
  set it up
- [Claude Code](../agents/claude-code.md) — run the Claude Code CLI inside a
  sandbox
