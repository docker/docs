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

Follow the Claude Desktop instructions for
[SSH sessions](https://code.claude.com/docs/en/desktop#ssh-sessions). When you
add the connection, enter the sandbox hostname, such as `demo.sbx`, manually in
the **SSH Host** field. Leave **SSH Port** and **Identity File** empty because
the managed SSH config supplies them.

## Related

- [Editor and app integrations](_index.md) — how SSH access works and how to
  set it up
- [Claude Code](../agents/claude-code.md) — run the Claude Code CLI inside a
  sandbox
