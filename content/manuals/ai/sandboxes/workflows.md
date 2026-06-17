---
title: Workflow patterns
linkTitle: Workflows
weight: 30
description: Common workflow patterns for Docker Sandboxes, covering git strategies, authenticated tools, commit signing, and CI integration.
keywords: docker sandboxes, sbx, workflows, clone mode, git, branches, commit signing, github cli, ci, headless
---

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

In clone mode, the sandbox gets a private Git clone. The agent manages its
own branches and commits inside that clone; your host working tree is never
touched. When the agent is done, you either fetch its branches to the host or
ask the agent to push directly to your fork.

Clone mode is designed for parallelism: a single clone-mode sandbox can hold
many branches at once, and subagent orchestrators (such as Claude Code's
[agents view](agents/claude-code.md#agents-view)) can dispatch independent
tasks to separate agents, each working on its own branch or worktree inside
the clone.

> [!NOTE]
> `--clone` is a create-time flag and cannot be changed on an existing
> sandbox. If you need to run additional non-clone sandboxes for the same
> repository, you would have to remove the clone-mode sandbox first.
> Keep a clone-mode sandbox running across multiple tasks rather than
> recreating it per task.

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

2. Dispatch each independent task to a separate subagent. Claude Code handles
   branch isolation for subagents automatically in agents view. For other
   agents (such as Codex), add an instruction to `AGENTS.md` to get the same
   behavior:

   ```markdown
   Always start each task on a new git branch before making changes.
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
> receive the new value. To update a running sandbox, scope the secret to
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

Agent credentials (API keys, GitHub token) can be pre-configured as global
secrets so they're available to any sandbox the CI runner creates:

```console
$ echo "$ANTHROPIC_API_KEY" | sbx secret set -g anthropic
$ echo "$GITHUB_TOKEN" | sbx secret set -g github
```
