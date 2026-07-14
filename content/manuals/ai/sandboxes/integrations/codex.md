---
title: Connect Codex to a sandbox
linkTitle: Codex
weight: 40
description: Run the Codex app against a Docker Sandbox over SSH.
keywords: docker sandboxes, codex, openai, remote ssh, sbx
---

{{< summary-bar feature_name="Docker Sandboxes SSH" >}}

Connect the Codex app to a sandbox over SSH so it works inside the isolated
environment instead of on your host.

> [!NOTE]
> This page covers the Codex app connecting to a sandbox over SSH. To run the
> Codex CLI inside a sandbox directly, see [Codex](../agents/codex.md).

## Prerequisites

- SSH access set up. See [Editor and app integrations](_index.md#enable-ssh-access).
- The Codex app installed.

## Connect

1. In the Codex app, add a new SSH connection.
2. Enter the sandbox host, `<name>.sbx`. Leave port and identity settings
   empty — the managed SSH config that `sbx setup ssh` wrote supplies them.
3. Connect. Codex runs against the sandbox over SSH.

<!-- TODO: confirm the exact Codex app connection flow and add screenshots -->

## Related

- [Editor and app integrations](_index.md) — how SSH access works and how to
  set it up
- [Codex](../agents/codex.md) — run the Codex CLI inside a sandbox
