---
title: Connect T3 Code to a sandbox
linkTitle: T3 Code
weight: 50
description: Run T3 Code against a Docker Sandbox over SSH.
keywords: docker sandboxes, t3 code, remote ssh, remote development, sbx
---

{{< summary-bar feature_name="Docker Sandboxes SSH" >}}

These connection instructions use a local sandbox. For cloud SSH setup, see
[Connect with SSH](../cloud/usage.md#connect-with-ssh).

T3 Code's SSH integration lets the desktop app drive coding agents inside a
sandbox. T3 Code has no dedicated Docker Sandboxes integration — it treats the
sandbox as an ordinary SSH host, connects to it, and starts a T3 server inside
that tunnels back to the app.

## Prerequisites

- SSH access set up. See [Editor and app integrations](_index.md#enable-ssh-access).
- T3 Code installed.

The first connection installs the T3 server in the sandbox. `t3` ships
prebuilt native modules for Linux, so nothing compiles, but those binaries
link against `libatomic1`. Docker-provided templates include it.

The [`t3code` kit](https://github.com/docker/sbx-kits-contrib/tree/main/t3code)
goes further: it installs the `t3` npm package when the sandbox is created,
along with `libatomic1` if the base image lacks it, so the first connection
starts a server that is already there. Pair it with any agent whose base image ships
Node.js 18 or later, which all standard agent templates do:

```console
$ sbx run claude --kit docker.io/sbx/t3code-kit:latest
```

A sandbox built from a custom
[base image](/manuals/ai/sandboxes/customize/author/base-images.md) without
`libatomic1` needs it installed before the first connection, unless you apply
the kit above, which installs it for you:

```console
$ sbx exec <sandbox> -- sudo apt-get update
$ sbx exec <sandbox> -- sudo DEBIAN_FRONTEND=noninteractive apt-get install -y libatomic1
```

Verify it's in place:

```console
$ sbx exec <sandbox> -- sh -lc 'ldconfig -p | grep libatomic.so.1'
```

## Connect

Confirm that you can connect to the sandbox from a terminal:

```console
$ ssh demo.sbx
```

In T3 Code, add an SSH environment and enter the sandbox hostname, such as
`demo.sbx`, as the host. The first connection installs the T3 server inside
the sandbox unless the `t3code` kit pre-installed it, so it can take a
moment. Later connections are faster.

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

The `npm WARN EBADENGINE` lines warn about the transitive `ini` dependency
and are separate from the failure: npm enforces engine requirements only
when `engine-strict` is set, which is off by default, so this warning alone
still lets the install proceed.

The most common causes are a missing shared library and a full disk. Install
`t3` by hand and run it to tell them apart:

```console
$ sbx exec <sandbox> -- sh -lc \
  'rm -rf /tmp/t3probe && mkdir -p /tmp/t3probe && cd /tmp/t3probe \
   && npm init -y >/dev/null && npm install t3@latest 2>&1 | tail -20 \
   && ./node_modules/.bin/t3 --version'
```

A missing `libatomic1` lets the install finish and leaves a binary that can't
start:

```text
node_modules/@t3code/t3-linux-x64/t3: error while loading shared libraries:
libatomic.so.1: cannot open shared object file: No such file or directory
```

Install it as described in [Prerequisites](#prerequisites).

A full disk fails the install itself with `ENOSPC`:

```text
npm ERR! code ENOSPC
npm ERR! nospc ENOSPC: no space left on device
```

Check free disk space:

```console
$ sbx exec <sandbox> -- df -h /
```

A sandbox can have both problems at once. Fixing one still leaves the same
top-level error, so check the library and disk space before concluding the
sandbox is ready. Free up space or install `libatomic1` as needed, then
reconnect.

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
