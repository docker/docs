---
title: Usage
weight: 30
description: Basic sbx commands for creating, managing, and connecting to Docker Sandboxes.
keywords: docker sandboxes, sbx, usage, run, create, stop, remove, ports, workspaces
---

Use this page as a command-oriented guide to day-to-day `sbx` operations. For
scenario-based recommendations, see [Workflow patterns](workflows/).

## Sign in

Sign in from a terminal:

```console
$ sbx login
```

For scripts or CI runners where a browser isn't available, see
[CI and headless use](workflows/automation.md).

## Start, stop, and remove

The basic workflow is [`run`](/reference/cli/sbx/run/) to start,
[`ls`](/reference/cli/sbx/ls/) to check status,
[`stop`](/reference/cli/sbx/stop/) to pause, and
[`rm`](/reference/cli/sbx/rm/) to clean up:

```console
$ sbx run claude .                  # start an agent in the current directory
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
$ sbx run claude .
```

To remove all stopped local sandboxes, use `sbx prune`. Running sandboxes are
never removed. Preview the sandboxes that would be removed, or filter out
sandboxes stopped within the last week:

```console
$ sbx prune --dry-run
$ sbx prune --filter since=168h
```

Run `sbx prune` without flags to confirm and remove all stopped sandboxes.

## Choose a workspace

`sbx run` mounts the current directory when you don't pass a workspace path.
Pass a path to mount another directory instead:

```console
$ sbx run claude
$ sbx run claude ~/my-project
```

The first workspace path is the primary workspace. The agent starts there, and
`sbx exec` uses it as the default working directory. The host directory is
mounted at the same absolute path inside the sandbox. When you don't pass a
path to `sbx run`, the current directory is the primary workspace.

Starting with `sbx` version 0.40.0, workspace paths are optional for
`sbx create`. Omit them to create a mountless sandbox without a host workspace
bind mount, then attach to the sandbox by name:

```console
$ sbx create --name scratch claude
$ sbx run --name scratch
```

Mountless mode is the default for `sbx create` when you don't pass a path. The
agent works in the container image's working directory. Built-in images use
`/home/agent/workspace`; a custom template can define another absolute working
directory. Files in a mountless workspace persist when you stop and restart
the sandbox, but they are deleted with the sandbox. Use `--name` to give a
mountless sandbox a stable identity for reconnecting. Use
[`sbx cp`](#copy-files-between-host-and-sandbox) to move files between the
mountless workspace and your host.

## Reconnect and name sandboxes

Sandboxes persist after the agent exits. Running the same workspace path again
reconnects to the existing sandbox rather than creating another sandbox:

```console
$ sbx run claude ~/my-project  # creates sandbox
$ sbx run claude ~/my-project  # reconnects to same sandbox
```

Use `--name` to give a sandbox an explicit identity:

```console
$ sbx run claude --name my-project .
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
agent. To create a sandbox with the current directory mounted in the background
without attaching:

```console
$ sbx create --name my-project claude .
```

Omit the path to create a mountless sandbox instead. Attach later with
`sbx run --name`:

```console
$ sbx create --name scratch claude
$ sbx run --name scratch
```

## Set environment variables

> [!NOTE]
> The `-e`/`--env` and `--env-file` flags require `sbx` version 0.39.0 or
> later.

Pass `-e` or `--env` to `sbx run` or `sbx create` to set an environment
variable in the sandbox:

```console
$ sbx run -e LOG_LEVEL=debug claude
```

Specify a variable name without a value to copy its value from the host
environment:

```console
$ export API_URL=https://api.example.com
$ sbx run -e API_URL claude
```

To load multiple variables, pass one or more environment files:

```console
$ sbx create --name my-project --env-file .env.sandbox claude .
```

The flags follow `docker run` precedence rules. Values passed with `-e`
override values from environment files. When you pass multiple environment
files, a value in a later file overrides the same variable in an earlier file.

When either command creates a sandbox, the variables are stored with the
sandbox. They are also available to the agent session started by `sbx run`.
When `sbx run` re-attaches to an existing sandbox, the variables apply to that
agent session without changing the sandbox's stored environment. To set
variables for one command instead, use `sbx exec -e` or
`sbx exec --env-file`.

To persist a variable across future sessions of an existing sandbox, append an
export to `/etc/sandbox-persistent.sh`:

```console
$ sbx exec -d <sandbox-name> bash -c "echo 'export INTERNAL_API_URL=https://api.example.com' >> /etc/sandbox-persistent.sh"
```

The `bash -c` wrapper ensures the `>>` redirect runs inside the sandbox instead
of on your host. The file is sourced when Bash starts inside the sandbox,
including for interactive sessions and agents started with `sbx run`. A command
passed directly to `sbx exec` doesn't start a shell. Wrap that command in
`bash -c` if it needs variables from the persistent environment file.

A variable added to the file only takes effect for sessions and agents started
afterward. Restart a running agent, or stop and start the sandbox, to pick up
the new value.

Environment variables are readable by processes inside the sandbox. For API
keys and other credentials, use [`sbx secret set`](configuration/credentials.md#store-a-secret)
for a supported service or the experimental
[`sbx secret set-custom`](configuration/credentials.md#custom-secrets) for a
credential sent to known hosts. The host-side proxy can then inject the real
value without exposing it to the agent.

## Run commands inside a sandbox

To get a shell inside a running sandbox, use [`sbx exec`](/reference/cli/sbx/exec/):

```console
$ sbx exec -it <sandbox-name> bash
```

Without `--workdir`, the command starts in the sandbox's primary workspace. In
a mountless sandbox, it starts in the container image's working directory.

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

## Git workspace modes

When your primary workspace is a Git repository, choose how the sandbox receives
it when you create the sandbox:

- Direct mode is the default for `sbx run`. It also applies when you pass a
  workspace path to `sbx create`. The agent has read-write access to your
  working tree, and changes appear on your host immediately.
- [Clone mode](#clone-mode) uses `--clone`. The agent edits a separate Git clone
  inside the sandbox. Its changes stay there until you fetch them or the agent
  pushes them. Your host repository is also available at
  `/run/sandbox/source`, but only with read access.

For guidance on branch strategy, fetching work from a sandbox, and parallel
agent workflows, see [Git workflows](workflows/git.md). For the
security model behind each mode, see
[Workspace isolation](security/isolation.md#workspace-isolation).

### Clone mode

To create a clone-mode sandbox, pass `--clone` when you run or create it:

```console
$ sbx run --clone claude .
```

You can also create the sandbox in the background and attach later:

```console
$ sbx create --clone --name my-sandbox claude .
$ sbx run --name my-sandbox
```

Clone mode has a few create-time constraints:

- Clone mode is fixed at create time. To switch an existing sandbox to clone
  mode, remove it and recreate it with `sbx create --clone`.
- The clone follows whichever ref your host repository has checked out at create
  time. No branch is created automatically.
- The primary workspace must be a Git repository. Omit `--clone` for non-Git
  workspaces.
- Clone mode is rejected from inside a Git worktree other than the main one. The
  read-only bind mount can't resolve the worktree's `.git` pointer file. Run
  `sbx create --clone <agent> .` from the main repository checkout instead.
- Removing a clone-mode sandbox drops the in-sandbox clone. Fetch or push any
  commits you want to keep before you remove it.

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

You can also run separate projects side-by-side. Remove unused sandboxes when
you're done to reclaim disk space:

```console
$ sbx run claude ~/project-a
$ sbx run claude ~/project-b
$ sbx rm <sandbox-name>       # when finished
```

## Copy files between host and sandbox

Use [`sbx cp`](/reference/cli/sbx/cp/) to copy files or directories between
your host and a sandbox. This is useful for one-off files that aren't part of a
mounted workspace, such as generated output, logs, or setup files. For example,
copy files to or from the workspace used by a built-in image:

```console
$ sbx cp ./config.json my-sandbox:/home/agent/workspace/
$ sbx cp my-sandbox:/home/agent/workspace/output.log ./
$ sbx cp ./src/ my-sandbox:/home/agent/workspace/src
```

One side of the copy must use `SANDBOX:PATH`. Copying directly between two
sandboxes isn't supported.

## Publish ports

Sandboxes are [network-isolated](security/isolation.md) — your browser or local
tools can't reach a server running inside one by default. A port mapping of
`8080:3000` publishes sandbox port 3000 on host port 8080.

If you know which ports you need, publish them when you create the sandbox:

```console
$ sbx run --publish 8080:3000 --name my-sandbox claude
```

For an existing sandbox, use [`sbx ports`](/reference/cli/sbx/ports/) to
forward traffic from your host:

```console
$ sbx ports my-sandbox --publish 8080:3000
$ open http://localhost:8080
```

To let the OS pick a free host port instead of choosing one yourself, specify
only the sandbox port. Then use `sbx ports` to check which host port was
assigned:

```console
$ sbx ports my-sandbox --publish 3000
$ sbx ports my-sandbox
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

When `sbx run` re-attaches to an existing sandbox, it ignores `--publish`. Use
`sbx ports` to publish ports on that sandbox. For dev server and host-service
recipes, see
[Local services](workflows/development.md#local-services).

## What persists

While a sandbox exists, installed packages, Docker images, configuration
changes, command history, and mountless workspace files all persist across
stops and restarts. When you remove a sandbox, everything inside is deleted.
Host workspace files, including repositories used as clone sources, and the
[shared agent skills store](workflows/agent-skills.md) remain on your host. To
preserve a configured environment, create a [custom
template](customize/templates.md) or use a [kit](customize/kits.md).
