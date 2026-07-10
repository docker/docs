---
title: Usage
weight: 20
description: Common patterns for working with sandboxes.
keywords: docker sandboxes, sbx, usage, run, policy, secrets, branches, git, workspaces, ssh
---

## Working with sandboxes

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

If the sandbox has an active session (an open attach, SSH connection, or
in-flight SFTP transfer), `sbx rm` refuses unless you pass `--force`:

```console
$ sbx rm --force my-sandbox
```

To get a shell inside a running sandbox — useful for inspecting the environment,
checking Docker containers, or manually installing something:

```console
$ sbx exec -it <sandbox-name> bash
```

If you need a clean slate, remove the sandbox and re-run:

```console
$ sbx stop my-sandbox
$ sbx rm my-sandbox
$ sbx run claude
```

## Non-interactive login

For CI environments and scripts where a browser is not available, use a
Docker Personal Access Token (PAT) with `--username` and `--password-stdin`:

```console
$ echo "$DOCKER_PAT" | sbx login --username <your-docker-id> --password-stdin
```

`--password-stdin` reads the token from standard input to keep it out of
your shell history. Generate a PAT from your
[Docker account settings](https://app.docker.com/settings/personal-access-tokens)
with at least **Read** scope.

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

When your workspace is a Git repository, you can choose one of two ways
to share it with a sandbox. You make the choice when you create the
sandbox:

- **Direct mode (default)** — the agent has read-write access to your
  working tree. Changes the agent makes appear on your host immediately.
- **[Clone mode](#clone-mode) (`--clone`)** — the agent works on a private
  Git clone inside the sandbox, with your host repository mounted
  read-only. The sandbox exposes its clone as a Git remote on your host,
  so you fetch the agent's commits the same way you'd fetch from any
  other remote.

For a comparison of approaches and step-by-step recipes, see
[Workflow patterns](workflows.md#git-workflows). For the security model
behind each mode, see
[Workspace isolation](security/isolation.md#workspace-isolation).

### Direct mode (default)

The agent edits your working tree directly. Stage, commit, and push as you
normally would. If you run multiple agents on the same repository at the
same time, they may step on each other's changes — use
[clone mode](#clone-mode) to give each agent its own isolated workspace.

### Clone mode

In clone mode, the sandbox becomes a Git remote on your host. Your entire
working directory, including untracked files and files excluded by `.gitignore`, is mounted
read-only inside the sandbox. The agent commits inside the sandbox; you pull its work back
out by fetching from that remote.

> [!NOTE]
> Clone mode was introduced in `sbx` v0.31.0 and replaces the `--branch`
> flag used in earlier versions. If your CLI doesn't recognize `--clone`,
> update to the latest version.

```console
$ sbx run --clone claude
```

You can also create the sandbox in the background and attach later:

```console
$ sbx create --clone --name my-sandbox claude .
$ sbx run --name my-sandbox
```

The clone follows whichever ref your host repository has checked out at
create time. No new branch is created automatically.

> [!NOTE]
> Clone mode is fixed at create time. To switch an existing sandbox to
> clone mode, remove it and recreate it with `sbx create --clone`.

The CLI copies the Git remotes from your host repository (`origin`,
`upstream`, and so on) into the in-sandbox clone. The agent can push to
your fork directly using the same remote names. Local-path remotes
(`file://` URLs, filesystem paths) aren't copied, since they aren't
reachable from inside the sandbox.

#### Sandbox lifecycle and the Git remote

The Git daemon that exposes the in-sandbox clone runs as part of the
sandbox itself. It's only reachable while the sandbox is running:

- `sbx stop` shuts down the daemon. `git fetch sandbox-<name>` fails until
  the sandbox starts again.
- Restarting the sandbox assigns a new ephemeral port to the daemon. The
  CLI updates the `sandbox-<name>` remote URL in your host repository's
  Git config automatically, so fetching continues to work without manual
  reconfiguration.
- `sbx rm` removes the sandbox, the daemon, the published port, and the
  `sandbox-<name>` remote entry from your host repository.

> [!WARNING]
> Removing a clone-mode sandbox drops the in-sandbox clone along with it.
> Any commits you haven't fetched (`git fetch sandbox-<name>`) or pushed
> to an upstream remote are lost. `sbx rm` prints a warning before
> deleting a clone-mode sandbox — review it before confirming.

#### Restrictions

Clone mode requires a Git repository as the primary workspace, and is
rejected at create time in two cases:

- `--clone` on a non-Git workspace. Omit `--clone` for non-Git workspaces.
- `--clone` from inside a Git worktree (other than the main one). The
  read-only bind mount can't resolve the worktree's `.git` pointer file.
  Run `sbx create --clone` from the main repository checkout instead.

You can also create a Git worktree yourself and run an agent inside it
without `--clone`, but the sandbox won't have access to the `.git`
directory in the parent repository, so the agent can't use Git at all.
See [Host worktree](workflows.md#host-worktree) in Workflow patterns.

## Reconnecting and naming

Sandboxes persist after the agent exits. Running the same workspace path again
reconnects to the existing sandbox rather than creating a new one:

```console
$ sbx run claude ~/my-project  # creates sandbox
$ sbx run claude ~/my-project  # reconnects to same sandbox
```

Use `--name` to give a sandbox an explicit identity:

```console
$ sbx run claude --name my-project
```

Once a named sandbox exists, `--name` is how you re-attach to it — from any
working directory, with or without the agent positional:

```console
$ sbx run --name my-project        # re-attaches from anywhere
$ sbx run claude --name my-project # same, with agent confirmed
```

Re-running a command that previously created a sandbox reconnects to it rather
than returning an error, so you can up-arrow and re-enter a session without
first looking up the sandbox name.

To run multiple sandboxes against the same workspace — for example, one for a
feature branch and one for exploratory changes — give each a distinct name:

```console
$ sbx run claude --name feature ~/my-project
$ sbx run claude --name spike ~/my-project
```

Both sandboxes share the same workspace but are otherwise independent.

## Creating without attaching

[`sbx run`](/reference/cli/sbx/run/) creates the sandbox and attaches you to
the agent. To create a sandbox in the background without attaching:

```console
$ sbx create --name my-project claude .
```

Unlike `run`, `create` requires an explicit workspace path. It uses direct
mode by default, or pass `--clone` for [clone mode](#clone-mode). Attach
later with `sbx run --name`:

```console
$ sbx run --name my-project
```

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

## Installing dependencies and using Docker

Ask the agent to install what's needed — it has sudo access, and installed
packages persist for the sandbox's lifetime. For teams or repeated setups,
see [Customize](customize/) for reusable templates and declarative kits.

Agents can also build Docker images, run containers, and use
[Compose](/manuals/compose/_index.md). Everything runs inside the sandbox's private Docker
daemon, so containers started by the agent never appear in your host's
`docker ps`. When you remove the sandbox, all images, containers, and volumes
inside it are deleted with it.

## Accessing services in the sandbox

Sandboxes are [network-isolated](security/isolation.md) — your browser or local
tools can't reach a server running inside one by default. Use
[`sbx ports`](/reference/cli/sbx/ports/) to forward traffic from your host into
a running sandbox.

The common case: an agent has started a dev server or API, and you want to open
it in your browser or run tests against it.

```console
$ sbx ports my-sandbox --publish 8080:3000   # host 8080 → sandbox port 3000
$ open http://localhost:8080
```

To let the OS pick a free host port instead of choosing one yourself:

```console
$ sbx ports my-sandbox --publish 3000        # ephemeral host port
$ sbx ports my-sandbox                       # check which port was assigned
```

`sbx ls` shows active port mappings alongside each sandbox, and `sbx ports`
lists them in detail:

```console
$ sbx ls
SANDBOX         AGENT   STATUS   PORTS                    WORKSPACE
my-sandbox      claude  running  127.0.0.1:8080->3000/tcp /home/user/proj
```

To stop forwarding a port:

```console
$ sbx ports my-sandbox --unpublish 8080:3000
```

For a service to be reachable, it must listen on all interfaces inside the
sandbox, not only `127.0.0.1`. Bind it to `0.0.0.0` for IPv4 or `[::]` for both
IPv4 and IPv6; most dev servers need a flag like `--host 0.0.0.0` to do this. On
the host, `--publish` listens on both `127.0.0.1` and `::1`, so a client
resolving `localhost` might pick IPv6 and fail with "connection reset by peer"
if the sandboxed service only listens on IPv4 — even when
`http://127.0.0.1:<port>/` works. To fix that, bind the service to `[::]`, or
pin the published port to one family with `--publish 8080:3000/tcp4` or `/tcp6`.

Published ports survive restarts: `sbx` re-publishes them when the sandbox or
the daemon restarts. Explicit host ports are reused, while a port published with
an OS-assigned host port (such as `--publish 3000`) gets a new host port on each
start, so check `sbx ports my-sandbox` to find it. If an explicit host port is
already in use at restart, the CLI or the dashboard prompts you to choose
another. Removing the sandbox releases its ports.

You can't publish ports at create time — there's no `--publish` flag on
`sbx run` or `sbx create`, so publish them once the sandbox is running. To stop
forwarding, `--unpublish 8080:3000` removes a single mapping, and
`--unpublish 3000` removes every host port mapped to sandbox port 3000.

## Accessing host services from a sandbox

Services running on your host are reachable from inside a sandbox using the
hostname `host.docker.internal`.
Use this instead of `127.0.0.1` or your machine's local network IP address,
which are not reachable from inside the sandbox.

The sandbox proxy translates `host.docker.internal` to `localhost` before
forwarding the request, so you must add the `localhost` address with the
specific port to your network policy allowlist:

```console
$ sbx policy allow network localhost:11434
```

Then use `host.docker.internal` in any configuration or request that points at
the host service. For example, to verify connectivity from a sandbox shell:

```console
$ curl http://host.docker.internal:11434
```

## Rolling out to a team

When rolling sandboxes out across a team, two features handle different
needs:

- [Custom templates and kits](customize/) let you package reusable agent
  configurations, MCP servers, base images, and per-project policies. Every
  developer pulls them down with their workspace.
- [Organization governance](governance/org.md) lets admins define
  network and filesystem rules in the Docker Admin Console. The rules apply
  across every developer's sandboxes and take precedence over local policy.
  Available on a separate paid subscription.

Customization gives developers shared starting points. Governance gives
admins centralized enforcement.

## What persists

While a sandbox exists, installed packages, Docker images, configuration
changes, and command history all persist across stops and restarts. When you
remove a sandbox, everything inside is deleted — only your workspace files
remain on your host. To preserve a configured environment, create a
[custom template](customize/templates.md).
