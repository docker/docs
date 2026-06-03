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

To get a shell inside a running sandbox — useful for inspecting the environment,
checking Docker containers, or manually installing something:

```console
$ sbx exec -it <sandbox-name> bash
```

If you need a clean slate, remove the sandbox and re-run:

```console
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
  Best when you're collaborating turn-by-turn with the agent on a single
  repository.
- **[Clone mode](#clone-mode) (`--clone`)** — the agent works on a private
  Git clone inside the sandbox, with your host repository mounted
  read-only. The sandbox exposes its clone as a Git remote on your host,
  so you fetch the agent's commits the same way you'd fetch from any
  other remote. Best when you want the agent isolated from your host
  repository — for running multiple agents in parallel, working with
  untrusted code, or keeping your working tree clean while the agent
  works.

See [Workspace isolation](security/isolation.md#workspace-isolation) for the
security model behind each mode.

### Direct mode (default)

The agent edits your working tree directly. Stage, commit, and push as you
normally would. If you run multiple agents on the same repository at the
same time, they may step on each other's changes — use
[clone mode](#clone-mode) to give each agent its own isolated workspace.

### Clone mode

In clone mode, the sandbox becomes a Git remote on your host. The agent
commits inside the sandbox; you pull its work back out by fetching from
that remote.

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
$ sbx run my-sandbox
```

The clone follows whichever ref your host repository has checked out at
create time. No new branch is created automatically. If you want the agent
to work on a dedicated branch, instruct it to run `git checkout -b
my-feature` inside the sandbox before it starts editing. Alternatively,
open a shell with `sbx exec` and create the branch yourself.

> [!NOTE]
> Clone mode is fixed at create time. To switch an existing sandbox to
> clone mode, remove it and recreate it with `sbx create --clone`.

#### Reviewing and merging the agent's commits

The CLI wires the in-sandbox clone as a `sandbox-<sandbox-name>` Git remote
on your host. Pull the agent's commits the same way you'd fetch any other
remote — no `cd` into a separate directory, no extra tooling:

```console
$ git fetch sandbox-my-sandbox
$ git log sandbox-my-sandbox/<branch>
$ git diff main..sandbox-my-sandbox/<branch>
$ git checkout -b my-feature sandbox-my-sandbox/<branch>
$ git push -u origin my-feature
$ gh pr create
```

If you asked the agent to work on a dedicated branch, `<branch>` is that
branch name. Otherwise it's whatever ref your host repository was on at
create time.

Some agents don't commit automatically. If `git log sandbox-<name>/<branch>`
shows nothing new, open a shell in the sandbox and commit from there:

```console
$ sbx exec -it my-sandbox bash
$ git commit -am "save work"
```

#### Pushing to your fork from inside the sandbox

When the sandbox starts, the CLI copies the Git remotes from your host
repository (`origin`, `upstream`, and so on) into the in-sandbox clone
with their existing URLs. The agent can push to your fork on GitHub
directly — for example, by prompting:

> Commit these changes and push them to a new branch on `origin`.

The push uses the same `git push origin ...` invocation the agent would
run on the host. This is interchangeable with fetching the commits to
your host first and pushing from there.

Local-path remotes (`file://` URLs, filesystem paths) aren't copied, since
they aren't reachable from inside the sandbox.

#### Running multiple branches in parallel

A single sandbox can hold several branches at once. Each branch the
agent commits to appears as a separate ref on the `sandbox-<name>`
remote, so you can fetch them independently from the host:

```console
$ git fetch sandbox-my-sandbox
$ git log sandbox-my-sandbox/feature-a
$ git log sandbox-my-sandbox/feature-b
```

A few common ways to have the agent start each task on its own branch:

- A subagent orchestrator such as Claude Code's
  [agents view](agents/claude-code.md#agents-view) dispatches each task
  to a subagent that creates its own worktree inside the clone.
- Agent-level instructions in `CLAUDE.md`, an orchestration skill, or a
  system prompt include a rule to start each task on a new branch.
- For one-off tasks, ask the agent to switch to a new branch before it
  starts.

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
Clone mode is the supported alternative for working on a separate branch.

### Signed commits

Sandboxes can sign Git commits with SSH keys from your host agent. The private
key stays on your host.

On the host, load the key into your SSH agent:

```console
$ ssh-add ~/.ssh/id_ed25519
```

Inside the sandbox, check that the forwarded agent exposes the key:

```console
$ ssh-add -L
```

Configure Git globally inside the sandbox to use SSH commit signing. This
writes to the sandbox user's Git config, not your repository's `.git/config`.
Use an inline public key instead of a key file path, because host paths such as
`~/.ssh/id_ed25519.pub` might not exist in the sandbox:

```console
$ git config --global gpg.format ssh
$ git config --global user.signingkey "key::$(ssh-add -L | head -n 1)"
```

Then commit as usual:

```console
$ git commit -S
```

For common signing failures, see
[Sandbox commits aren't signed](troubleshooting.md#sandbox-commits-arent-signed).

## Reconnecting and naming

Sandboxes persist after the agent exits. Running the same workspace path again
reconnects to the existing sandbox rather than creating a new one:

```console
$ sbx run claude ~/my-project  # creates sandbox
$ sbx run claude ~/my-project  # reconnects to same sandbox
```

Use `--name` to make this explicit and avoid ambiguity:

```console
$ sbx run claude --name my-project
```

## Creating without attaching

[`sbx run`](/reference/cli/sbx/run/) creates the sandbox and attaches you to
the agent. To create a sandbox in the background without attaching:

```console
$ sbx create claude .
```

Unlike `run`, `create` requires an explicit workspace path. It uses direct
mode by default, or pass `--clone` for [clone mode](#clone-mode). Attach
later with `sbx run`:

```console
$ sbx run claude-my-project
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
[Compose](https://docs.docker.com/compose/). Everything runs inside the sandbox's private Docker
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

A few things to keep in mind:

- **Services must listen on all interfaces** — a service listening only on
  `127.0.0.1` inside the sandbox won't be reachable through a published port.
  Bind to `0.0.0.0` for IPv4, or `[::]` to accept both IPv4 and IPv6. Most dev
  servers default to `127.0.0.1`, so you'll usually need to pass a flag like
  `--host 0.0.0.0` or `--host '[::]'` when starting them.
- **`localhost` on the host can resolve to IPv6** — by default, `--publish`
  listens on both `127.0.0.1` and `::1`. Your browser or client may pick IPv6
  when resolving `localhost`. If the sandboxed service only listens on IPv4,
  the IPv6 connection fails with "connection reset by peer" — even though
  `http://127.0.0.1:<port>/` works. To fix it, bind the sandboxed service to
  `[::]` so it accepts both families, or restrict the published port to one
  family with `--publish 8080:3000/tcp4` (IPv4) or `/tcp6` (IPv6).
- **Not persistent** — published ports are lost when the sandbox stops or the
  daemon restarts. Re-publish after restarting.
- **No create-time flag** — unlike `docker run -p`, there's no `--publish`
  option on `sbx run` or `sbx create`. Ports can only be published after the
  sandbox is running.
- **Unpublish requires the host port** — `--unpublish 3000` is rejected; you
  must use `--unpublish 8080:3000`. Run `sbx ports my-sandbox` first if you
  used an ephemeral port and need to find the assigned host port.

## Accessing host services from a sandbox

Services running on your host are reachable from inside a sandbox using the
hostname `host.docker.internal`.
Use this instead of `127.0.0.1` or your machine's local network IP address,
which are not reachable from inside the sandbox.

The sandbox proxy translates `host.docker.internal` to `localhost` before
forwarding the request, so you must add the `localhost` address with the
specific port to your network policy allowlist:

```console
$ sbx policy allow network -g localhost:11434
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
