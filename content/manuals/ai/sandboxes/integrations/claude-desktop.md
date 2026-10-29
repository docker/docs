---
title: Connect Claude Desktop to a sandbox
linkTitle: Claude Desktop
weight: 30
description: Run Claude Code from Claude Desktop against a Docker Sandbox over SSH.
keywords: docker sandboxes, claude desktop, claude code, remote ssh, sbx
---

{{< summary-bar feature_name="Docker Sandboxes SSH" >}}

These connection instructions use a local sandbox. For cloud SSH setup, see
[Connect with SSH](../cloud/usage.md#connect-with-ssh).

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
configures credential injection and network access for the remote Claude Code
session. You must supply the credential before connecting.

## Configure Claude credentials

Before connecting Claude Desktop, configure an Anthropic credential in Docker
Sandboxes using one of these methods:

- API key: run `sbx secret set anthropic` on your host to store the key.
- Claude subscription: run `sbx run claude` from a terminal, then use `/login`
  inside Claude Code to complete OAuth sign-in.

These credentials are available globally, including to sandboxes created on
first connection. To restrict an API key to an existing sandbox, use
`sbx secret set anthropic --sandbox <sandbox-name>`. See
[Manage credentials](../configuration/credentials.md#stored-secrets).

Signing in to Claude Desktop doesn't authenticate Claude Code inside the
sandbox. The SSH bridge replaces the credential from Claude Desktop with a
placeholder to keep the real token out of the sandbox. The sandbox proxy uses
the credential you configured in Docker Sandboxes to authenticate requests.

## Connect

For a dedicated sandbox created on first connection, see
[Set up desktop app connections](_index.md#set-up-desktop-app-connections).
The following steps connect to a sandbox with a mounted host workspace.

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
