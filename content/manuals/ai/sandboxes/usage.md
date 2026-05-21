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

When your workspace is a Git repository, the agent edits your working tree
directly by default. Changes appear in your working tree immediately, the same
as working in a normal terminal.

If you want the agent to work in isolation from your host repository — for
example to run multiple agents in parallel, or to prevent any chance of an
agent rewriting your local Git state — use [clone mode](#clone-mode). The
agent runs against a private Git clone inside the sandbox; your host repository
sees the agent's commits only after you explicitly fetch them.

### Direct mode (default)

The agent edits your working tree directly. Stage, commit, and push as you
normally would. If you run multiple agents on the same repository at the same
time, they may step on each other's changes. See
[clone mode](#clone-mode) for an alternative.

### Clone mode

Pass `--clone` to run the agent on a private Git clone living entirely
inside the sandbox, instead of bind-mounting your working tree. Your host
repository is mounted read-only as the clone's reference, so the agent — even
with full root inside the VM — cannot modify any byte of your `.git`
directory or working tree. You can set `--clone` on `create` or, equivalently,
on `run` at create time.

When `--clone` is active:

- The agent works on a private clone populated by
  `git clone --reference` from your repository, on the sandbox's overlay
  filesystem. The clone has its own index, refs, and working tree.
  Object storage is shared via `.git/objects/info/alternates`, so the
  clone is space-efficient and full history is walkable, but writes to
  the clone never reach your host's object database.
- Your repository's Git root is bind-mounted at `/run/sandbox/source` as a
  read-only mount. The agent's `git clone --reference` reads from this
  mount; nothing on the host is writable from inside the sandbox.
- The clone follows whatever HEAD your host repository is on at create time.
  No branch is created automatically — if you want the agent to work on a
  dedicated branch, instruct the agent to `git checkout -b my-feature`
  inside the sandbox before it starts editing.
- The sandbox runs a `git-daemon` over a `127.0.0.1`-bound ephemeral port
  that exports the in-container clone. The CLI registers it as a Git remote
  named `sandbox-<sandbox-name>` on your host repository, so you can pull
  the agent's commits with `git fetch`.
- Forge remotes you have on the host (`origin`, `upstream`, …) are
  propagated into the in-container clone with their existing URLs, so the
  agent can `git push origin …` to your GitHub fork as you would on the
  host. Local-path remotes (`file://`, paths) are skipped because they
  aren't reachable from inside the sandbox.

See [Source-repository isolation](security/isolation.md#source-repository-isolation)
for the security boundary.

#### Starting a sandbox in clone mode

```console
$ sbx run --clone claude   # private clone of the current repository
```

You can also create the sandbox first and attach later:

```console
$ sbx create --clone --name my-sandbox claude .
$ sbx run my-sandbox                       # resumes in the in-container clone
```

> [!NOTE]
> Clone mode is fixed at create time. Recreate the sandbox with
> `sbx create --clone ...` to switch an existing sandbox into clone mode.

#### Reviewing and pushing changes

The CLI wires the agent's in-container clone as a `sandbox-<sandbox-name>`
Git remote on your host repository. Review the agent's work with the same
commands you'd use for any other remote — no extra tooling, no `cd` into
a separate directory:

```console
$ git fetch sandbox-my-sandbox                            # pull the agent's commits
$ git log sandbox-my-sandbox/<branch-the-agent-used>      # see what the agent did
$ git diff main..sandbox-my-sandbox/<branch-the-agent-used>
$ git checkout -b my-feature sandbox-my-sandbox/<branch-the-agent-used>
$ git push -u origin my-feature
$ gh pr create
```

If the agent committed on a dedicated branch (because you asked it to
`git checkout -b ...`), that branch name appears on the `sandbox-<name>`
remote. If it stayed on the HEAD it inherited at create time, its commits
extend that branch instead — you'll see them by fetching and diffing.

Some agents don't commit automatically. If `git log sandbox-<name>/...`
shows nothing new, open a shell in the sandbox and commit from there
before fetching. `sbx exec` drops you into the in-container clone:

```console
$ sbx exec -it my-sandbox bash
$ git commit -am "save work"
```

See [Workspace trust](security/workspace.md) for security considerations when
reviewing agent changes.

#### Cleanup

`sbx rm` deletes the sandbox, its in-container clone, the published Git
port, and the `sandbox-<sandbox-name>` remote on your host repository.

> [!WARNING]
> Any commits the agent made inside the sandbox that you have not yet
> fetched (via `git fetch sandbox-<name>`) or pushed to an upstream
> remote will be lost — the in-container clone lives on the sandbox's
> overlay filesystem and is dropped with it. `sbx rm` prints a warning
> for clone-mode sandboxes; review it before confirming the removal.

#### Restrictions

A few configurations are incompatible with clone mode and are rejected at
create time:

- `--clone` together with `--workspace-volume`: the source-repository
  isolation relies on bind-mounting your Git root, which is incompatible
  with a volume-backed workspace.
- `--clone` from inside a host Git worktree: the bind mount can't resolve
  the worktree's `.git` pointer file. Run `sbx create --clone ...` from
  the main repository instead.
- `--clone` on a non-Git workspace: clone mode requires a Git repository.
  Run `sbx create` without `--clone` for non-Git workspaces.

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
