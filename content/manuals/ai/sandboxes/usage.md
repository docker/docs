---
title: Usage
weight: 20
description: Common patterns for working with sandboxes.
keywords: docker sandboxes, sbx, usage, run, policy, secrets, branches, git, workspaces, ssh
---

{{< summary-bar feature_name="Docker Sandboxes sbx" >}}

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

When your workspace is a Git repository, the agent edits your working tree
directly by default. Changes appear in your working tree immediately, the same
as working in a normal terminal.

If you run multiple agents on the same repository at once, use [branch
mode](#branch-mode) to give each agent its own branch and working directory.

### Direct mode (default)

The agent edits your working tree directly. Stage, commit, and push as you
normally would. If you run multiple agents on the same repository at the same
time, they may step on each other's changes. See
[branch mode](#branch-mode) for an alternative.

### Branch mode

Pass `--branch <name>` to give the agent its own branch and an isolated clone
of your repository inside the sandbox. This prevents conflicts when multiple
agents, or you and an agent, write to the same files at the same time. You
can set `--branch` on `create` or, equivalently, on `run` at create time.

When `--branch` is active:

- The agent works on a private clone living entirely inside the sandbox.
  Your repository's `.git` directory is bind-mounted as read-only into the
  sandbox, serving as a reference for the clone's object database, so the
  agent reuses your local history without consuming extra disk space.
- The CLI creates the new branch on your host repository and checks it out
  if your working tree is clean. If it's dirty, the branch ref is still
  created but the checkout is skipped with a warning, so your uncommitted
  changes are preserved.
- The sandbox runs a `git-daemon` over a `127.0.0.1`-bound ephemeral port
  that exports the in-container clone. The CLI registers it as a Git remote
  named `sandbox-<sandbox-name>` on your host repository, so you can pull
  the agent's commits with `git fetch`.
- The agent's clone has its own index, refs, and working tree. Concurrent
  Git operations on the host and inside the sandbox can't corrupt each
  other. See [Source-repository isolation](security/isolation.md#source-repository-isolation)
  for the security implications.

#### Starting a branch

```console
$ sbx run claude --branch my-feature   # agent works on the my-feature branch
```

Use `--branch auto` to let the CLI generate a branch name for you:

```console
$ sbx run claude --branch auto
```

You can also create the sandbox first and attach later:

```console
$ sbx create --name my-sandbox --branch my-feature claude .
$ sbx run my-sandbox                       # resumes in the my-feature clone
```

> [!NOTE]
> A sandbox is bound to the branch chosen at create time. To work on a
> different branch, create a new sandbox with `sbx create --branch
> other-feature ...`. Running `sbx run --branch ...` on an existing sandbox
> with a different branch is rejected.

#### Reviewing and pushing changes

The CLI wires the agent's in-container clone as a `sandbox-<sandbox-name>`
Git remote on your host repository. Review the agent's work with the same
commands you'd use for any other remote — no `cd` into a worktree, no extra
tooling:

```console
$ git fetch sandbox-my-sandbox                            # pull the agent's commits
$ git log sandbox-my-sandbox/my-feature                   # see what the agent did
$ git diff main..sandbox-my-sandbox/my-feature            # full diff
$ git checkout my-feature && git merge --ff-only \
    sandbox-my-sandbox/my-feature                         # fast-forward your local branch
$ git push -u origin my-feature
$ gh pr create
```

Some agents don't commit automatically. If `git log sandbox-<name>/<branch>`
shows nothing new, open a shell in the sandbox and commit from there before
fetching. `sbx exec` drops you into the in-container clone, so plain `git`
commands work without changing directory:

```console
$ sbx exec -it my-sandbox bash
$ git commit -am "save work"
```

See [Workspace trust](security/workspace.md) for security considerations when
reviewing agent changes.

#### Cleanup

`sbx rm` deletes the sandbox, its in-container clone, the published Git
port, and the `sandbox-<sandbox-name>` remote on your host repository. The
local branch the agent worked on stays on your host so you don't lose any
commits you've already fetched.

#### Restrictions

A few configurations are incompatible with branch mode and are rejected at
create time:

- `--branch` together with `--workspace-volume`: the source-repository
  isolation relies on bind-mounting your Git root, which is incompatible
  with a volume-backed workspace.
- `--branch` from inside a host Git worktree: the bind mount can't resolve
  the worktree's `.git` pointer file. Run `sbx create --branch ...` from
  the main repository instead.

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
mode by default, or pass `--branch` for [branch mode](#branch-mode). Attach
later with `sbx run`:

```console
$ sbx run claude-my-project
```

## Multiple workspaces

You can mount extra directories into a sandbox alongside the main workspace.
The first path is the primary workspace — the agent starts here, and the
sandbox's branch-mode clone is created from this directory if you use
`--branch`. Extra workspaces are always mounted directly.

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

- **Services must bind to `0.0.0.0`** — a service listening on `127.0.0.1`
  inside the sandbox won't be reachable through a published port. Most dev
  servers default to `127.0.0.1`, so you'll usually need to pass a flag like
  `--host 0.0.0.0` when starting them.
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
- [Organization governance](security/governance.md) lets admins define
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
