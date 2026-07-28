---
title: Connect ChatGPT to a sandbox
linkTitle: ChatGPT
weight: 40
description: Run Codex in the ChatGPT desktop app against a Docker Sandbox over SSH.
keywords: docker sandboxes, chatgpt, codex, openai, remote ssh, sbx
---

{{< summary-bar feature_name="Docker Sandboxes SSH" >}}

Connect the ChatGPT desktop app to a sandbox over SSH so Codex works inside the
isolated environment instead of on your host.

> [!NOTE]
> This page covers running Codex in the ChatGPT desktop app connected to a
> sandbox over SSH. To run the Codex CLI inside a sandbox directly, see
> [Codex](../agents/codex.md).

## Prerequisites

- SSH access set up. See [Editor and app integrations](_index.md#enable-ssh-access).
- The ChatGPT desktop app installed.
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

In the ChatGPT desktop app, open **Settings > Connections** and add an SSH
connection manually. Enter the sandbox hostname, such as `demo.sbx`, as the
host, then use the remote folder picker to
[select the mounted workspace](_index.md#select-the-workspace-folder) as the
remote project.

For more connection options, see the OpenAI instructions to
[connect to an SSH host](https://learn.chatgpt.com/docs/remote-connections#connect-to-an-ssh-host).

## Related

- [Editor and app integrations](_index.md) — how SSH access works and how to
  set it up
- [Codex](../agents/codex.md) — run the Codex CLI inside a sandbox
