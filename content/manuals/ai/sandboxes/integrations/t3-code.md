---
title: Connect T3 Code to a sandbox
linkTitle: T3 Code
weight: 50
description: Run T3 Code against a Docker Sandbox over SSH.
keywords: docker sandboxes, t3 code, remote ssh, remote development, sbx
---

{{< summary-bar feature_name="Docker Sandboxes SSH" >}}

T3 Code's SSH integration lets the desktop app drive coding agents inside a
sandbox. T3 Code has no dedicated Docker Sandboxes integration — it treats the
sandbox as an ordinary SSH host, connects to it, and starts a T3 server inside
that tunnels back to the app.

## Prerequisites

- SSH access set up. See [Editor and app integrations](_index.md#enable-ssh-access).
- T3 Code installed.

## Connect

Confirm that you can connect to the sandbox from a terminal:

```console
$ ssh demo.sbx
```

In T3 Code, add an SSH environment and enter the sandbox hostname, such as
`demo.sbx`, as the host. The first connection installs the T3 server inside
the sandbox, so it can take a moment. Later connections are faster.

Then add a new project, select the SSH environment from the list, and
[choose the mounted workspace](_index.md#select-the-workspace-folder) as the
project directory inside the sandbox.

## Troubleshoot a server that never becomes ready

T3 Code can fail to connect with an error like the following, wrapped here
for readability. It concatenates the connection failure with npm's install
output from inside the sandbox into a single error dialog:

```text
Could not prepare the SSH environment: ... SshCommandError: Connecting to
sandbox "sandboxes"… Remote T3 server did not become ready on
127.0.0.1:3773. npm WARN EBADENGINE Unsupported engine { package:
'ini@7.0.0', required: { node: '^22.22.2 || ^24.15.0 || >=26.0.0' },
current: { node: 'v22.22.1', npm: '9.2.0' } }
```

The `npm WARN EBADENGINE` lines warn about the transitive `ini` dependency.
They aren't the cause of the failure — npm doesn't enforce engine
requirements by default, so this warning alone doesn't stop the install.
Check free disk space in the sandbox. A full disk can fail the T3 server
install partway through, and T3 Code reports this the same way as a
connection timeout:

```console
$ sbx exec <sandbox> -- df -h /
```

Free up space and reconnect if the sandbox is close to full.

## Troubleshoot `turn/setPermissionMode failed`

If your organization manages Claude Code with a policy file, a local T3 Code
thread can fail to start with `turn/setPermissionMode failed`. T3 Code's
default runtime mode is Full access, which maps to the Claude Agent SDK's
`bypassPermissions` mode. A managed policy that disables that mode rejects
the request.

On macOS, check whether this applies to you:

```console
$ cat "/Library/Application Support/ClaudeCode/managed-settings.json"
```

If `permissions.disableBypassPermissionsMode` is set to `disable`, switch T3
Code to a different runtime mode, such as Supervised, Auto-accept edits, or
Auto, then start a new thread. The permission mode is captured once when a
thread starts, so switching modes in an already-failing thread doesn't
recover it.

This restriction applies to the host running Claude Code, not to a sandbox.
A thread connected to a sandbox isn't subject to the host's managed policy,
so Full access works normally there.

## Related

- [Editor and app integrations](_index.md) — how SSH access works and how to
  set it up
