---
title: Editor and app integrations
linkTitle: Integrations
weight: 37
description: Connect editors and desktop apps to a Docker Sandbox over SSH.
keywords: docker sandboxes, ssh, integrations, vs code, cursor, remote development, sbx
---

{{< summary-bar feature_name="Docker Sandboxes SSH" >}}

You can connect an external editor or desktop app to a running sandbox over
SSH. This lets you use the tools you already know — VS Code, Cursor, Claude
Desktop, and others — while your code runs, builds, and executes inside the
isolated sandbox instead of on your host.

Each sandbox is reachable at `<name>.sbx`, where `<name>` is the sandbox name.
Once SSH is set up, `<name>.sbx` behaves like any other SSH host, so any tool
that supports remote development over SSH can connect to it.

> [!NOTE]
> SSH access is experimental and off by default. The command surface and
> behavior may change.

## How it works

`sbx setup ssh` writes a managed block to your `~/.ssh/config` that maps the
`*.sbx` host pattern to the sandbox daemon. Connections don't use a network
port or an SSH key:

- A `ProxyCommand` relays the SSH stream to the daemon over its local socket
  (a Unix domain socket on macOS and Linux, a named pipe on Windows).
- The daemon accepts the connection only while you have an active Docker login.
  Authentication is tied to your login, not to a stored key.
- The host key is verified on every connection, so a rotated daemon key never
  triggers a host-key mismatch.

Because SSH terminates at the daemon, no SSH server runs inside the sandbox.
Connecting to `<name>.sbx` starts the sandbox if it isn't running. The sandbox
must already exist.

## Prerequisites

- The `sbx` CLI installed and signed in. See [Get started](../get-started.md).
- An SSH client. macOS and most Linux distributions include OpenSSH. On
  Windows, install the OpenSSH client.
- The editor or app you want to connect, with its remote-over-SSH support
  installed.

## Enable SSH access

SSH access is an experimental feature. Turn it on, then stop the daemon so it
reloads with the new setting. The daemon reads `feature.ssh` only at startup,
so the change takes effect the next time it starts:

```console
$ sbx settings set platform.allowExperimentalFeatures true
$ sbx settings set feature.ssh true
$ sbx daemon stop
$ sbx setup ssh
```

The `sbx` CLI starts the daemon automatically in the background when a command
needs it, so `sbx setup ssh` brings it back up with SSH enabled — you don't
start it by hand. To start it yourself instead, use `sbx daemon start -d`; the
`-d` flag runs it in the background rather than holding your terminal.

`sbx setup ssh` is idempotent — you can re-run it at any time. It adds a
managed block to `~/.ssh/config` similar to the following:

```text
# >>> docker sandboxes (managed) >>>
Host *.sbx
    User _default_user_
    ProxyCommand "sbx" ssh proxy %n
    IdentityAgent none
    IdentityFile /dev/null
    IdentitiesOnly yes
    ControlMaster no
    ControlPath none
    UserKnownHostsFile "~/.ssh/sbx_known_hosts"
    KnownHostsCommand "sbx" ssh known-hosts %H
    StrictHostKeyChecking yes
    SendEnv *
# <<< docker sandboxes (managed) <<<
```

You don't edit this block by hand. The `User _default_user_` sentinel tells the
daemon to log you in as the sandbox image's default user, so your host username
is never sent.

The wildcard entry configures how SSH clients connect, but it doesn't add
individual sandbox names to application host pickers. Enter the sandbox
hostname, such as `demo.sbx`, manually when you configure an integration.

## Create or identify a sandbox

SSH connections require an existing sandbox. To create a named shell sandbox
for the current directory:

```console
$ sbx create --name demo shell .
```

To identify an existing sandbox, list your sandboxes:

```console
$ sbx ls
```

## Connect to a sandbox over SSH

Use the sandbox name with the `.sbx` suffix. For example, to connect to a
sandbox named `demo`:

```console
$ ssh demo.sbx
```

## Connect a specific tool

- [VS Code](vscode.md)
- [Cursor](cursor.md)
- [Claude Desktop](claude-desktop.md)
- [ChatGPT](chatgpt.md)
