---
title: Usage
weight: 20
description: Basic sbx commands for creating, managing, and connecting to Docker Sandboxes.
keywords: docker sandboxes, sbx, usage, run, create, stop, remove, ports, workspaces
---

Use this page as a command-oriented guide to day-to-day `sbx` operations. For
scenario-based recommendations, see [Workflow patterns](workflows.md).

## Sign in

Sign in from a terminal:

```console
$ sbx login
```

For scripts or CI runners where a browser isn't available, see
[CI and headless use](workflows.md#ci-and-headless-use).

## Start, stop, and remove

The basic workflow is [`run`](/reference/cli/sbx/run/) to start,
[`ls`](/reference/cli/sbx/ls/) to check status,
[`stop`](/reference/cli/sbx/stop/) to pause, and
[`rm`](/reference/cli/sbx/rm/) to clean up:

```console
$ sbx run claude                    # start an agent
$ sbx ls                            # see what's running
$ sbx stop my-sandbox               # pause it
$ sbx rm my-sandbox                 # delete it entirely
```

If the sandbox has an active session — an open attach, SSH connection, or
in-flight SFTP transfer — `sbx rm` refuses unless you pass `--force`:

```console
$ sbx rm --force my-sandbox
```

If you need a clean slate, remove the sandbox and run it again:

```console
$ sbx stop my-sandbox
$ sbx rm my-sandbox
$ sbx run claude
```

## Reconnect and name sandboxes

Sandboxes persist after the agent exits. Running the same workspace path again
reconnects to the existing sandbox rather than creating another sandbox:

```console
$ sbx run claude ~/my-project  # creates sandbox
$ sbx run claude ~/my-project  # reconnects to same sandbox
```

Use `--name` to give a sandbox an explicit identity:

```console
$ sbx run claude --name my-project
```

Once a named sandbox exists, use `--name` to re-attach to it from any working
directory, with or without the agent positional:

```console
$ sbx run --name my-project        # re-attaches from anywhere
$ sbx run claude --name my-project # same, with agent confirmed
```

To run multiple sandboxes against the same workspace, give each a distinct
name:

```console
$ sbx run claude --name feature ~/my-project
$ sbx run claude --name spike ~/my-project
```

## Create without attaching

[`sbx run`](/reference/cli/sbx/run/) creates the sandbox and attaches you to the
agent. To create a sandbox in the background without attaching:

```console
$ sbx create --name my-project claude .
```

Unlike `run`, `create` requires an explicit workspace path. Attach later with
`sbx run --name`:

```console
$ sbx run --name my-project
```

## Run commands inside a sandbox

To get a shell inside a running sandbox, use [`sbx exec`](/reference/cli/sbx/exec/):

```console
$ sbx exec -it <sandbox-name> bash
```

## Interactive mode

Running `sbx` with no subcommands opens an interactive terminal dashboard:

```console
$ sbx
```

The dashboard shows all your sandboxes as cards with live status, CPU, and
memory usage. From here you can:

- **Create** a sandbox (`c`).
- **Start or stop** a sandbox (`s`).
- **Attach** to an agent session (`Enter`), same as `sbx run`.
- **Open a shell** inside the sandbox (`x`), same as `sbx exec`.
- **Remove** a sandbox (`r`).

The dashboard also includes a network governance panel where you can monitor
outbound connections made by your sandboxes and manage network rules. Use `tab`
to switch between the sandboxes panel and the network panel.

From the network panel you can browse connection logs, allow or block specific
hosts, and add custom network rules. Press `?` to see all keyboard shortcuts.

## Git workflow

When your workspace is a Git repository, you can share it with a sandbox in one
of two ways:

- Direct mode (default): the agent has read-write access to your
  working tree. Changes the agent makes appear on your host immediately.
- [Clone mode](#clone-mode) (`--clone`): the agent works on a private
  Git clone inside the sandbox, with your host repository mounted
  read-only. The sandbox exposes its clone as a Git remote on your host,
  so you fetch the agent's commits the same way you'd fetch from any
  other remote.

For a comparison of approaches and workflow recipes, see
[Workflow patterns](workflows.md#git-workflows). For the security model
behind each mode, see
[Workspace isolation](security/isolation.md#workspace-isolation).

### Direct mode (default)

The agent edits your working tree directly. Stage, commit, and push from the
host as you normally would. If you run multiple agents on the same repository at
the same time, use [clone mode](#clone-mode) to give each agent an isolated
workspace.

### Clone mode

In clone mode, the sandbox becomes a Git remote on your host. Your entire
working directory, including untracked files and files excluded by `.gitignore`,
is mounted read-only inside the sandbox. The agent commits inside the sandbox.
You pull its work back out by fetching from that remote:

```console
$ git fetch sandbox-<name>
```

```console
$ sbx run --clone claude
```

You can also create the sandbox in the background and attach later:

```console
$ sbx create --clone --name my-sandbox claude .
$ sbx run --name my-sandbox
```

The clone follows whichever ref your host repository has checked out at
create time. No branch is created automatically.

Clone mode is fixed at create time. To switch an existing sandbox to clone mode,
remove it and recreate it with `sbx create --clone`.

> [!WARNING]
> Removing a clone-mode sandbox drops the in-sandbox clone along with it.
> Any commits you haven't fetched (`git fetch sandbox-<name>`) or pushed
> to an upstream remote are lost. `sbx rm` prints a warning before
> deleting a clone-mode sandbox — review it before confirming.

Clone mode requires a Git repository as the primary workspace, and is
rejected at create time in two cases:

- `--clone` on a non-Git workspace. Omit `--clone` for non-Git workspaces.
- `--clone` from inside a Git worktree (other than the main one). The
  read-only bind mount can't resolve the worktree's `.git` pointer file.
  Run `sbx create --clone` from the main repository checkout instead.

## Multiple workspaces

You can mount extra directories into a sandbox alongside the main workspace.
The first path is the primary workspace — the agent starts here, and the
sandbox's in-container Git clone is populated from this directory if you
use `--clone`. Extra workspaces are always mounted directly.

All workspaces appear inside the sandbox at their absolute host paths. Append
`:ro` to mount an extra workspace read-only — useful for reference material or
shared libraries the agent shouldn't modify:

```console
$ sbx run claude ~/project-a ~/shared-libs:ro ~/docs:ro
```

Each sandbox is completely isolated, so you can also run separate projects
side-by-side. Remove unused sandboxes when you're done to reclaim disk space:

```console
$ sbx run claude ~/project-a
$ sbx run claude ~/project-b
$ sbx rm <sandbox-name>       # when finished
```

## Copying files between host and sandbox

Use [`sbx cp`](/reference/cli/sbx/cp/) to copy files or directories between
your host and a sandbox. This is useful for one-off files that aren't part of a
mounted workspace, such as generated output, logs, or setup files.

```console
$ sbx cp ./config.json my-sandbox:/home/user/
$ sbx cp my-sandbox:/home/user/output.log ./
$ sbx cp ./src/ my-sandbox:/home/user/src
```

One side of the copy must use `SANDBOX:PATH`. Copying directly between two
sandboxes isn't supported.

## Publish ports

Sandboxes are [network-isolated](security/isolation.md) — your browser or local
tools can't reach a server running inside one by default. Use
[`sbx ports`](/reference/cli/sbx/ports/) to forward traffic from your host into
a running sandbox.

```console
$ sbx ports my-sandbox --publish 8080:3000   # host 8080 → sandbox port 3000
$ open http://localhost:8080
```

To let the OS pick a free host port instead of choosing one yourself:

```console
$ sbx ports my-sandbox --publish 3000        # ephemeral host port
$ sbx ports my-sandbox                       # check which port was assigned
```

`sbx ls` shows active port mappings alongside each sandbox. `sbx ports` lists
them in detail.

```console
$ sbx ls
SANDBOX         AGENT   STATUS   PORTS                    WORKSPACE
my-sandbox      claude  running  127.0.0.1:8080->3000/tcp /home/user/proj
```

To stop forwarding a port:

```console
$ sbx ports my-sandbox --unpublish 8080:3000
```

You can't publish ports at create time — there's no `--publish` flag on
`sbx run` or `sbx create`, so publish them once the sandbox is running. For
dev server and host-service recipes, see
[Local services](workflows.md#local-services).

## What persists

While a sandbox exists, installed packages, Docker images, configuration
changes, and command history all persist across stops and restarts. When you
remove a sandbox, everything inside is deleted — only your workspace files
remain on your host. To preserve a configured environment, create a
[custom template](customize/templates.md) or use a [kit](customize/kits.md).
