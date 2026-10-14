---
title: Editor and app integrations
linkTitle: Integrations
weight: 70
description: Connect editors and desktop apps to a Docker Sandbox over SSH.
keywords: docker sandboxes, ssh, integrations, vs code, cursor, remote development, sbx
---

{{< summary-bar feature_name="Docker Sandboxes SSH" >}}

These integrations use local sandbox SSH access. Cloud sandboxes use a
different SSH configuration and address: see
[Connect with SSH](../cloud/usage.md#connect-with-ssh).

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

## Set up desktop app connections

On macOS and Windows, Docker Sandboxes can add connections to ChatGPT and Claude
Desktop that create a dedicated sandbox on first use. This workflow uses a
workspace inside the sandbox. To work directly on a mounted host project,
[create a sandbox](#create-or-identify-a-sandbox) and connect to its hostname
instead.

This setup requires Docker Sandboxes v0.48.0 or later. Launch each installed
app at least once, then quit the apps before running:

```console
$ sbx setup ssh --auto-create
```

Setup adds the following connections to the apps and configures their SSH
hostnames. It creates the sandboxes when the connections are used.

| App | Connection name | SSH hostname | Sandbox name |
| --- | --- | --- | --- |
| ChatGPT | Docker Sandbox (Codex) | `codex.sbx` | `ssh-codex` |
| Claude Desktop | Docker Sandbox (Claude) | `claude.sbx` | `ssh-claude` |

Reopen each app and select its added connection. Connecting for the first time
creates the corresponding sandbox. Later connections reuse that sandbox and
its files, starting it if it is stopped.

These sandboxes have no host workspace mounted. Clone your project inside the
sandbox and select its directory in the app's remote folder picker. Changes
stay inside the sandbox until you push them to a remote repository or copy
them out.

Removing the sandbox with `sbx rm` leaves the app connection available. The
next connection creates a replacement sandbox without restoring deleted files.

If setup reports that an app is running, quit it and rerun `sbx setup ssh`.
Setup also provides sign-in guidance if Docker or agent authentication is
missing.

### Remove desktop app connections

Quit ChatGPT and Claude Desktop, then run:

```console
$ sbx setup ssh remove
```

This removes the managed SSH config and app connections. It leaves existing
sandboxes running or stopped as they were. Remove the sandboxes separately
with `sbx rm` if you no longer need them.

## Create or identify a sandbox

To connect using a sandbox name as `<name>.sbx`, create the sandbox first.
ChatGPT and Claude Desktop also support [creating a dedicated sandbox on first connection](#set-up-desktop-app-connections).
To create a named shell sandbox for the current directory:

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
open the primary workspace automatically. The initial folder depends on the
client. A remote folder picker might open at the sandbox user's home directory,
`/home/agent`, while an interactive `ssh` shell might start in
`/home/agent/workspace`. Select the intended folder explicitly instead of
relying on the initial location.

For a sandbox with a primary workspace, each workspace path appears inside the
sandbox at the same absolute path as on the host. For example, if you pass
`/Users/bob/src/my-project`, select that path in the remote folder picker. For
a mountless sandbox that uses a Docker-provided agent template, select
`/home/agent/workspace`.

## Connect a specific tool

- [VS Code](vscode.md)
- [Cursor](cursor.md)
- [Claude Desktop](claude-desktop.md)
- [ChatGPT](chatgpt.md)
- [T3 Code](t3-code.md)

## How SSH connections work

### Managed SSH configuration

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
# <<< docker sandboxes (managed) <<<
```

You don't edit this block by hand. Its key entries work as follows:

- `Host *.sbx` maps sandbox hostnames to the sandbox daemon. Application host
  pickers don't discover individual sandbox names from this wildcard, so enter
  the hostname, such as `demo.sbx`, manually when you configure an integration.
- `User _default_user_` tells the daemon to use the sandbox image's default
  user, so your host username is never sent.

### Connection and authentication

Connections don't use a network port or an SSH key:

- A `ProxyCommand` relays the SSH stream to the daemon over its local socket
  (a Unix domain socket on macOS and Linux, a named pipe on Windows).
- The daemon accepts the connection only while you have an active Docker login.
  Authentication is tied to your login, not to a stored key.
- The host key is verified on every connection, so a rotated daemon key never
  triggers a host-key mismatch.

Because SSH terminates at the daemon, no SSH server runs inside the sandbox.
For ordinary `<name>.sbx` connections, create the sandbox first. The
[desktop app connections](#set-up-desktop-app-connections) create their dedicated
sandboxes on first use. If a sandbox is stopped, connecting starts it
automatically.

### Environment variables

SSH connections don't forward client environment variables into the sandbox.
The daemon acknowledges SSH environment requests for compatibility but ignores
their names and values.

### Port forwarding

SSH clients can use local port forwarding to make a service listening on the
sandbox's loopback interface available on the host. For example, a remote
development client can map `127.0.0.1:4321` in the sandbox to
`127.0.0.1:55565` on the host, choosing an available host port automatically.
Traffic passes through the SSH connection instead of a published Docker port.

The sandbox daemon accepts forwarded connections only to loopback addresses in
the sandbox, including `localhost`, `127.0.0.0/8`, and `::1`. The SSH client
chooses the bind address for the listener on the host. A listener bound to
`127.0.0.1` or `::1` is reachable only from the host. A client configured to
bind to a non-loopback address can make the forwarded service reachable from
other machines, subject to the host's network and firewall configuration.
