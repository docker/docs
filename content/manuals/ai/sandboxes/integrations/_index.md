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

## Prerequisites

- The `sbx` CLI installed and signed in. See [Get started](../get-started.md).
- An SSH client. macOS and most Linux distributions include OpenSSH. On
  Windows, install the OpenSSH client.
- The editor or app you want to connect, with its remote-over-SSH support
  installed.

## Enable SSH access

Run the SSH setup command once:

```console
$ sbx setup ssh
```

The command starts the Docker Sandboxes daemon if needed and configures your
SSH client. You can re-run it at any time.

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

## Select the workspace folder

Connecting an app to a sandbox selects the remote environment, but it might not
open the mounted workspace automatically. Use the app's remote folder picker to
select the workspace when you configure the connection or start a session.

The folder picker might open at the sandbox user's home directory, typically
`/home/agent`. Workspaces retain their absolute host paths inside the sandbox.
For example, if you mount `/Users/bob/src/my-project`, select
`/Users/bob/src/my-project` in the remote folder picker.

## Connect a specific tool

- [VS Code](vscode.md)
- [Cursor](cursor.md)
- [Claude Desktop](claude-desktop.md)
- [ChatGPT](chatgpt.md)

## How SSH connections work

`sbx setup ssh` writes a managed block to your SSH config: `~/.ssh/config` on
macOS and Linux, or `%USERPROFILE%\.ssh\config` on Windows. The block is similar
to the following:

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

`SendEnv *` offers host environment variables to the daemon, but the daemon
accepts only variables in its `ssh.acceptEnv` allowlist. Execution-sensitive
variables such as `PATH`, `LD_*`, and `NODE_OPTIONS` are always blocked, even
if added to the allowlist. Accepted values apply only to the SSH session and
aren't stored in the sandbox configuration.

The `*.sbx` wildcard maps sandbox hostnames to the sandbox daemon, but it
doesn't add individual sandbox names to application host pickers. Enter the
sandbox hostname, such as `demo.sbx`, manually when you configure an
integration.

Connections don't use a network port or an SSH key:

- A `ProxyCommand` relays the SSH stream to the daemon over its local socket
  (a Unix domain socket on macOS and Linux, a named pipe on Windows).
- The daemon accepts the connection only while you have an active Docker login.
  Authentication is tied to your login, not to a stored key.
- The host key is verified on every connection, so a rotated daemon key never
  triggers a host-key mismatch.

Because SSH terminates at the daemon, no SSH server runs inside the sandbox.
The sandbox must already be created. If it is stopped, connecting to
`<name>.sbx` starts it automatically.
