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

The first connection installs the T3 server in the sandbox, which needs a
build toolchain. T3 depends on `node-pty`, which ships prebuilt binaries only
for macOS and Windows. On a Linux sandbox, `node-pty` compiles from source and
the build fails without `make`, `python3`, and a compiler such as `g++`.

The [`t3code` kit](https://github.com/docker/sbx-kits-contrib/tree/main/t3code)
prepares a sandbox for T3 Code: it installs the build toolchain and the `t3`
npm package when the sandbox is created, so the first connection starts a
pre-installed server instead of building `node-pty` from source. Pair it with
any agent whose base image ships Node.js 18 or later, which all standard
agent templates do:

```console
$ sbx run claude --kit docker.io/sbx/t3code-kit:latest
```

For an existing sandbox, install the toolchain manually:

```console
$ sbx exec <sandbox> -- sudo apt-get update
$ sbx exec <sandbox> -- sudo DEBIAN_FRONTEND=noninteractive apt-get install -y g++ make python3
```

Verify the toolchain is in place:

```console
$ sbx exec <sandbox> -- sh -lc 'command -v g++ && command -v make && command -v python3'
```

A manual install lasts only until the sandbox is recreated, and the first
connection still builds `node-pty` from source. For a setup that persists,
recreate the sandbox with the [kit](../customize/kits.md) or a custom
[template](../customize/templates.md).

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

The most common causes are a missing C++ toolchain and a full disk, and both
produce this identical error. Get npm's actual output to tell them apart:

```console
$ sbx exec <sandbox> -- sh -lc \
  'rm -rf /tmp/t3probe && mkdir -p /tmp/t3probe && cd /tmp/t3probe \
   && npm init -y >/dev/null && npm install t3@latest 2>&1 | tail -40'
```

A missing compiler fails the native `node-pty` build with `Error 127` from
`make`:

```text
npm ERR! make: g++: No such file or directory
npm ERR! make: *** [pty.target.mk:115: Release/obj.target/pty/src/unix/pty.o] Error 127
npm ERR! gyp ERR! build error
npm ERR! gyp ERR! stack Error: `make` failed with exit code: 2
```

Install the build toolchain as described in [Prerequisites](#prerequisites).

A full disk fails with `ENOSPC`, and no gyp output appears at all because npm
fails before the native build starts:

```text
npm ERR! code ENOSPC
npm ERR! nospc ENOSPC: no space left on device
```

Check free disk space:

```console
$ sbx exec <sandbox> -- df -h /
```

A sandbox can have both problems at once. Fixing one still leaves the same
top-level error, so check both the toolchain and disk space before
concluding the sandbox is ready. Free up space or install the toolchain as
needed, then reconnect.

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
