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

Use a sandbox created with the Claude agent type. The Claude sandbox template
configures Anthropic credentials and network access for the remote Claude Code
session.

## Connect

> [!WARNING]
> Connecting Claude Desktop to a sandbox over SSH transmits Anthropic
> credentials into the Claude Code process within the sandbox, reducing
> isolation guarantees.

Create a named Claude sandbox for the current directory if you don't already
have one:

```console
$ sbx create --name demo claude .
```

Confirm that you can connect to the sandbox from a terminal:

```console
$ ssh demo.sbx
```

In Claude Desktop, open the environment drop-down before starting a session and
select **+ Add SSH connection**. Enter a name for the connection and enter the
sandbox hostname, such as `demo.sbx`, in **SSH Host**. Leave **SSH Port** and
**Identity File** empty because the managed SSH config supplies them.

Select the connection from the environment drop-down, then use the remote
folder picker to
[select the mounted workspace](_index.md#select-the-workspace-folder). The
picker might initially open at `/home/agent`.

For more connection options, see the Claude Desktop instructions for
[SSH sessions](https://code.claude.com/docs/en/desktop#ssh-sessions).

## Troubleshoot a broken SSH connection after token refresh

The SSH connection drops when the Anthropic token expires and needs to be
refreshed. To work around this, run the sandbox manually from your host:

```console
$ sbx run --name <sandbox-name>
```

## Troubleshoot SSH connection timeouts on Windows

Claude Desktop requires Git on Windows. If an SSH connection times out and the
Claude Desktop logs include `ProxyCommand error: spawn sh ENOENT`, install
[Git for Windows](https://git-scm.com/download/win).

If Git is already installed, verify that `sh.exe` is available on your `PATH`:

```powershell
PS> where.exe sh
```

If the command doesn't find `sh.exe`, add the Git `bin` directory to your user
`Path`. The default directory is `C:\Program Files\Git\bin`. Quit and restart
Claude Desktop after updating `Path`.

## Related

- [Editor and app integrations](_index.md) — how SSH access works and how to
  set it up
- [Claude Code](../agents/claude-code.md) — run the Claude Code CLI inside a
  sandbox
