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
- An existing sandbox created from the Codex template. The template includes
  the `codex` command required by the app's remote server.

## Connect

Create a named Codex sandbox for the current directory if you don't already
have one:

```console
$ sbx create --name demo codex .
```

Confirm that you can connect to the sandbox from a terminal:

```console
$ ssh demo.sbx
```

In the Codex app, open **Settings > Connections** and add an SSH connection
manually. Enter the sandbox hostname, such as `demo.sbx`, as the host, then
choose the sandbox workspace folder as the remote project.

For more connection options, see the Codex instructions to
[connect to an SSH host](https://learn.chatgpt.com/docs/remote-connections#connect-to-an-ssh-host).

## Related

- [Editor and app integrations](_index.md) — how SSH access works and how to
  set it up
- [Codex](../agents/codex.md) — run the Codex CLI inside a sandbox
