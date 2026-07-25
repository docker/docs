---
title: Workflow patterns
linkTitle: Workflows
weight: 30
description: Workflow patterns for Docker Sandboxes, covering shared agent skills, git strategies, local services, authenticated tools, and CI integration.
keywords: docker sandboxes, sbx, workflows, agent skills, shared skills, clone mode, git, branches, commit signing, github cli, local services, ci, headless
---

Use this page when you need to choose an approach for a specific way of working
with sandboxes. For command syntax and lifecycle basics, see
[Usage](usage.md).

## Share agent skills

Shared agent skills make skills from supported agents on your host available
inside your sandboxes. Importing copies the skills into a persistent store that
survives sandbox deletion and is shared by default with new sandboxes that run
a supported agent.

> [!NOTE]
> Shared agent skills are experimental.

Preview the skills that `sbx` finds without copying them:

```console
$ sbx skills import --dry-run
```

The command scans the following directories in order and copies each skill
subdirectory into the shared store. When the sandbox starts, `sbx` mounts the
store at the path the agent reads inside the sandbox.

| Agent       | Host source         | Sandbox mount target          |
| ----------- | ------------------- | ----------------------------- |
| Claude Code | `~/.claude/skills`  | `/home/agent/.claude/skills`  |
| Codex       | `~/.agents/skills`  | `/home/agent/.agents/skills`  |
| Copilot     | `~/.copilot/skills` | `/home/agent/.copilot/skills` |
| Cursor      | `~/.cursor/skills`  | `/home/agent/.cursor/skills`  |
| Droid       | `~/.factory/skills` | `/home/agent/.factory/skills` |

All imported skills go into the same store, regardless of their source. If
more than one source contains a skill with the same directory name, the skill
from the first source in the table wins and `sbx` warns about the others.

Import the skills:

```console
$ sbx skills import
```

The final output reports the shared store path. The default locations are:

| Platform | Shared store path                                                           |
| -------- | --------------------------------------------------------------------------- |
| macOS    | `~/Library/Application Support/com.docker.sandboxes/sandboxes/agent-skills` |
| Linux    | `~/.local/state/sandboxes/sandboxes/agent-skills`                           |
| Windows  | `%LOCALAPPDATA%\DockerSandboxes\sandboxes\state\agent-skills`               |

On Linux, `sbx` uses `$XDG_STATE_HOME/sandboxes/sandboxes/agent-skills` when
`XDG_STATE_HOME` is set.

When a skill already exists in the store, `sbx` prompts before replacing it.
Use `--force` to replace existing skills without prompts. Importing replaces
the complete skill directory rather than merging files. Run the import command
again when you want to copy updates from the host. Running `sbx reset` clears
the shared store.

Sandboxes created with `sbx` version 0.37.0 or later for a supported agent are
configured to mount the store read-write by default. These sandboxes mount the
current contents of the store each time they start, so you can import skills
before or after creating them. To create a sandbox without the shared store,
use `--no-share-skills`:

```console
$ sbx run --no-share-skills claude
```

Upgrading `sbx` does not enable shared skills for sandboxes created with an
earlier version. Remove and recreate those sandboxes after upgrading. The
`--no-share-skills` option also only applies when the sandbox is created. To
turn off shared skills for an existing sandbox, remove it and recreate it with
the option.

> [!WARNING]
> The shared skills store is mounted read-write. A sandbox can modify any skill
> in the store, and another sandbox can later load the modified instructions or
> run the modified scripts. The store is dedicated sandbox state, so this does
> not by itself execute the modified skill on your host. It does put every
> sandbox that shares the store in the same trust boundary. Use
> `--no-share-skills` to keep a sandbox outside that boundary.

Some agents scan for skills when a session starts. If imported skills don't
appear in an existing session, start another agent session.

## Git workflows

Sandboxes support three approaches for working with Git repositories. The
right choice depends on whether you want branch isolation and whether you
plan to run tasks in parallel:

|                           | Direct mode      | Clone mode (`--clone`)       | Host worktree                 |
| ------------------------- | ---------------- | ---------------------------- | ----------------------------- |
| Branch management         | You, on the host | Agent, inside the clone      | You, on the host              |
| Changes visible on host   | Immediately      | After fetch or agent push    | Immediately                   |
| Agent can use Git         | Yes              | Yes                          | No                            |
| Parallelism               | No               | Multiple agents, one sandbox | One sandbox per parallel task |
| Mode fixed at create time | No               | Yes                          | —                             |

### Direct mode

The simplest approach. The sandbox mounts your host working tree directly —
the agent edits files in place and changes appear immediately. You manage
branches yourself.

1. Check out the branch you want to work on:

   ```console
   $ git checkout -b feat/my-feature
   ```

2. Start the sandbox. No special flags needed:

   ```console
   $ sbx run claude
   ```

3. The agent edits files in your working tree. Review diffs, stage, and
   commit as you normally would:

   ```console
   $ git diff
   $ git add -p
   $ git commit
   $ git push -u origin feat/my-feature
   ```

Because the sandbox mounts your working tree, switching branches on the host
also changes what the agent sees. This makes direct mode well-suited for
focused, single-branch work where you're collaborating with the agent
turn-by-turn.

### Clone mode

In clone mode, `sbx` creates a separate Git clone inside the sandbox. The agent
edits this clone instead of your host working tree. Its changes stay inside the
sandbox until you fetch a branch or the agent pushes one to a remote. Your host
repository is also available at `/run/sandbox/source`, but only with read
access. The sandbox clone is not a Git worktree linked to your host checkout.

A single clone-mode sandbox can hold multiple branches and worktrees for
parallel tasks. The `--clone` flag creates the clone, but it doesn't separate
one task from another. To keep parallel tasks isolated, instruct your agent tool
to create a separate branch or worktree for each task.

> [!NOTE]
> `--clone` is a create-time flag and cannot be changed on an existing
> sandbox. To change a sandbox from clone mode to direct mode, remove and
> recreate it. To run both modes against the same repository, create separate
> sandboxes with distinct names.

#### Sandbox remote behavior

The CLI copies Git remotes from your host repository, such as `origin` and
`upstream`, into the in-sandbox clone. Local-path remotes, such as `file://`
URLs and filesystem paths, aren't copied because they aren't reachable from
inside the sandbox.

The Git daemon that exposes the in-sandbox clone runs as part of the sandbox.
It's only reachable while the sandbox is running:

- `sbx stop` shuts down the daemon. `git fetch sandbox-<name>` fails until the
  sandbox starts again.
- Restarting the sandbox assigns another ephemeral port to the daemon. The CLI
  updates the `sandbox-<name>` remote URL in your host repository's Git config,
  so fetching continues without manual reconfiguration.
- `sbx rm` removes the sandbox, the daemon, the published port, and the
  `sandbox-<name>` remote entry from your host repository.

#### Single task

1. Start a clone-mode sandbox:

   ```console
   $ sbx run --clone claude
   ```

2. Ask the agent to create a branch before it starts editing:

   > Create a branch `feat/my-feature` and make the changes.

3. Fetch the agent's branch when it's done:

   ```console
   $ git fetch sandbox-<name>
   $ git log sandbox-<name>/feat/my-feature
   $ git diff main..sandbox-<name>/feat/my-feature
   ```

4. Pull the branch to the host and push, or ask the agent to push directly:

   ```console
   # Pull to host, then push
   $ git checkout -b feat/my-feature sandbox-<name>/feat/my-feature
   $ git push -u origin feat/my-feature
   $ gh pr create

   # Or ask the agent
   # "Push feat/my-feature to origin and open a PR."
   ```

#### Parallel tasks

1. Start a clone-mode sandbox and open the
   [agents view](agents/claude-code.md#agents-view):

   ```console
   $ sbx run --clone claude
   ```

2. Dispatch each independent task to a separate background session. Your agent
   tool may use branches or worktrees to keep their changes separate. If it
   doesn't, add a project instruction such as:

   ```markdown
   Always start each task on its own git branch before making changes.
   ```

3. Fetch all branches when the agents are done:

   ```console
   $ git fetch sandbox-<name>
   $ git log sandbox-<name>/feat/task-a
   $ git log sandbox-<name>/feat/task-b
   ```

4. Check out the branches you want to keep and open PRs as normal.

### Host worktree

You can create a Git worktree on your host and point the sandbox at it. The
agent edits files directly in the worktree — but because the sandbox mounts
only the worktree directory (not the parent repository), it can't resolve the
`.git` pointer file and has no Git access. The agent can read and write files,
but can't commit, branch, or check status.

This is useful when you want branch isolation without the create-time
commitment of clone mode, and you're comfortable committing from the host
yourself after reviewing the changes.

1. Create the worktree on the host:

   ```console
   $ git worktree add -b feat/my-feature ../my-feature-work
   ```

2. Start the sandbox with the worktree as the workspace:

   ```console
   $ sbx run claude ../my-feature-work
   ```

3. The agent edits files. When it's done, commit and push from the host:

   ```console
   $ cd ../my-feature-work
   $ git diff
   $ git add -p && git commit
   $ git push -u origin feat/my-feature
   $ gh pr create
   ```

## Build and test inside a sandbox

Agents have sudo access inside the sandbox, so they can install packages,
start databases, run test dependencies, and prepare the environment they need.
Installed packages persist for the sandbox's lifetime. For repeated setup, use
[Customize](customize/) to package the environment as a template or kit.

Agents can also build Docker images, run containers, and use
[Compose](/manuals/compose/_index.md). Everything runs inside the sandbox's
private Docker daemon, so containers started by the agent never appear in your
host's `docker ps`. When you remove the sandbox, all images, containers, and
volumes inside it are deleted with it.

This pattern works well for tasks where the agent needs to run the project's
test suite or inspect a service it started. If you need to reach that service
from your host, publish a port after the sandbox is running.

## Local services

Use this workflow when a sandboxed agent starts a dev server, or when the agent
needs to call a service running on your host.

### Accessing services in the sandbox

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
IPv4 and IPv6. Most dev servers need a flag like `--host 0.0.0.0` to do this.
On the host, `--publish` listens on both `127.0.0.1` and `::1`, so a client
resolving `localhost` might pick IPv6 and fail with "connection reset by peer"
if the sandboxed service only listens on IPv4, even when
`http://127.0.0.1:<port>/` works. To fix that, bind the service to `[::]`, or
pin the published port to one family with `--publish 8080:3000/tcp4` or
`/tcp6`.

Published ports survive restarts: `sbx` re-publishes them when the sandbox or
the daemon restarts. Explicit host ports are reused, while a port published with
an OS-assigned host port, such as `--publish 3000`, gets a different host port
on each start. Check `sbx ports my-sandbox` to find it. If an explicit host port
is already in use at restart, the CLI or the dashboard prompts you to choose
another. Removing the sandbox releases its ports.

You can't publish ports at create time — there's no `--publish` flag on
`sbx run` or `sbx create`, so publish them once the sandbox is running. To stop
forwarding, `--unpublish 8080:3000` removes a single mapping, and
`--unpublish 3000` removes every host port mapped to sandbox port 3000.

### Accessing host services from a sandbox

Services running on your host are reachable from inside a sandbox using the
hostname `host.docker.internal`. Use this instead of `127.0.0.1` or your
machine's local network IP address, which are not reachable from inside the
sandbox.

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

## Commit signing

Sandboxes forward your host SSH agent into the sandbox, so the agent can
sign commits with your SSH key without the private key ever leaving your
host.

1. On your host, make sure the signing key is loaded in your SSH agent:

   ```console
   $ ssh-add ~/.ssh/id_ed25519
   $ ssh-add -L  # confirm the key appears
   ```

2. Inside the sandbox, configure Git to sign with SSH. Use the forwarded key
   directly rather than a file path, since host paths don't exist inside the
   sandbox:

   ```console
   $ git config --global gpg.format ssh
   $ git config --global user.signingkey "key::$(ssh-add -L | head -n 1)"
   ```

3. Sign commits as usual:

   ```console
   $ git commit -S -m "feat: my change"
   ```

To apply this configuration automatically to every sandbox, use the
[`git-ssh-sign`](https://github.com/docker/sbx-kits-contrib/tree/main/git-ssh-sign)
community kit, which handles all of the above setup. See [Kits](customize/kits.md)
if you want to package it alongside other sandbox customizations.

For troubleshooting, see
[Sandbox commits aren't signed](troubleshooting.md#sandbox-commits-arent-signed).

## Authenticated CLI tools

The sandbox proxy handles API credentials for model providers automatically,
but agents often also need credentials for tools like `gh`, `docker`, or a
secrets manager. The pattern is the same in each case: configure the
credential on your host once, and the sandbox either forwards it via the
proxy or via SSH agent forwarding.

> [!NOTE]
> The `-g` flag stores a secret globally so all future sandboxes can use it.
> Sandboxes that already exist when you run `sbx secret set -g` do not
> receive the updated value. To update a running sandbox, scope the secret to
> it directly: `sbx secret set <sandbox-name> <service>`.

### GitHub CLI

Store your GitHub token as a sandbox secret. The proxy injects it into
outbound requests, so `gh` works inside the sandbox without any additional
configuration:

```console
$ echo "$(gh auth token)" | sbx secret set -g github
```

The agent can then create pull requests, open issues, comment on PRs, and
interact with the GitHub API the same way it would from your host:

```console
# Inside the sandbox
$ gh pr create --title "feat: my feature" --body "..."
$ gh issue list
```

The token is never stored in plaintext inside the sandbox. See
[GitHub token](security/credentials.md#github-token) for details.

### Docker registry

When using Docker Hub, authentication is handled automatically; `sbx` reuses
your existing login session. For other registries, you need to configure
credentials for `sbx` so you can push or pull private images, including
[templates](customize/templates.md). The agent can then run `docker build` and
`docker push`, and `sbx` can resolve private image references, without any
extra authentication:

```console
$ gh auth token | sbx secret set --registry ghcr.io \
    --username <github-username> --password-stdin
$ echo "$ACR_PASSWORD" | sbx secret set --registry myregistry.azurecr.io \
    --username myuser --password-stdin
```

Images and containers built inside the sandbox run on the sandbox's private
Docker daemon, not your host's. They're deleted when the sandbox is removed.

For information on how registry credentials differ from other secrets,
per-registry username requirements, and global versus per-sandbox scoping, see
[Registry credentials](security/credentials.md#registry-credentials).

### Sourcing credentials from 1Password

#### Populating stored secrets with `op read`

Use `op read` to populate stored secrets without pasting values manually. Store
the value once and it's available to all future sandboxes:

```console
$ op read "op://Work/GitHub/token" | sbx secret set -g github
$ op read "op://Work/Anthropic/credential" | sbx secret set -g anthropic
```

The real value stays on your host; the sandbox sees the proxy-managed
placeholder as usual.

#### Per-launch injection with `op run`

To resolve credentials fresh from your vault on each launch without storing
them via `sbx secret set`, use `op run`:

```console
$ ANTHROPIC_API_KEY="op://Work/Anthropic/credential" op run -- sbx run claude
$ OPENAI_API_KEY="op://Work/OpenAI/key" op run -- sbx run codex
$ GEMINI_API_KEY="op://Work/Google/key" op run -- sbx run gemini
```

`op run` resolves each `op://` reference in the environment before executing
`sbx`. The sandbox reads the
[built-in service environment variables](security/credentials.md#built-in-services)
at launch and routes them through its proxy — the credential is never stored in
sbx's state and never appears inside the sandbox container.

This only applies to those specific credential variables. The sandbox does not
forward arbitrary environment variables from the host into the sandbox.

For multiple credentials at once, use `--env-file` with a file of `op://`
references:

```console
$ cat .sbx-secrets.env
ANTHROPIC_API_KEY=op://Work/Anthropic/credential
GITHUB_TOKEN=op://Work/GitHub/token

$ op run --env-file=.sbx-secrets.env -- sbx run claude
```

## CI and headless use

For CI environments and scripts where a browser isn't available, authenticate
with a Docker Personal Access Token (PAT):

```console
$ echo "$DOCKER_PAT" | sbx login --username <your-docker-id> --password-stdin
```

Generate a PAT from your
[Docker account settings](https://app.docker.com/settings/personal-access-tokens)
with at least **Read** scope.

From there, the rest of the `sbx` workflow is the same as interactive use.
Create the sandbox in the background with `sbx create`, run agent tasks with
`sbx exec`, and clean up with `sbx rm`:

```console
$ sbx create --name ci-task --clone claude
$ sbx run --name ci-task  # attach and give instructions, or use sbx exec for one-off commands
$ git fetch sandbox-ci-task
$ sbx rm ci-task
```

Agent credentials (API keys, GitHub token) can be preconfigured as global
secrets so they're available to any sandbox the CI runner creates. If the
relevant environment variables are already set in the CI environment (see the
[built-in services table](security/credentials.md#built-in-services) for which
variables each service reads), import them all at once:

```console
$ sbx secret import --all
```

To overwrite an existing stored entry, add `--force`. To pass a value from your
CI provider's secret store, use `-t`. For example, in a GitHub Actions step:

```yaml
- run: sbx secret set -g anthropic -t "${{ secrets.ANTHROPIC_API_KEY }}"
```

## Share setup across a team

When several people use sandboxes on the same project, separate repeatable
environment setup from policy enforcement.

Use [custom templates and kits](customize/) for project-level setup: agent
configurations, MCP servers, base images, setup scripts, and per-project
defaults. Version kit specs and template definitions with the project, and
publish reusable template images to your registry. This gives each developer the
same starting environment.

Use [organization governance](governance/org.md) for rules that admins need to
apply across developers, such as network and filesystem policies. Organization
rules are managed in the Docker Admin Console, take precedence over local
policy, and require a separate paid subscription.

You can use both. Templates and kits describe the development environment;
governance defines the boundaries it runs within.
