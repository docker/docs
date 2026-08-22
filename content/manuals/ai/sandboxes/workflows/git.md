---
title: Use Git with sandboxes
linkTitle: Git workflows
weight: 10
description: Choose a Git workspace mode and manage branches, parallel tasks, and signed commits with Docker Sandboxes.
keywords: docker sandboxes, sbx, git, clone mode, direct mode, worktrees, branches, commit signing
---

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

## Direct mode

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

## Clone mode

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

### Sandbox remote behavior

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

### Single task

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

### Parallel tasks

1. Start a clone-mode sandbox and open the
   [agents view](../agents/claude-code.md#agents-view):

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

## Host worktree

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

With SSH agent forwarding enabled, sandboxes forward your host SSH agent into
the sandbox, so the agent can sign commits with your SSH key without the
private key ever leaving your host.

1. Enable SSH agent forwarding, which is disabled by default, and restart the
   daemon:

   ```console
   $ sbx settings set ssh.agentForwardingEnabled true
   $ sbx daemon restart
   ```

   For how the forwarded socket is selected, including using a fixed socket
   path, see [SSH agent](../configuration/credentials.md#ssh-agent).

2. On your host, make sure the signing key is loaded in your SSH agent:

   ```console
   $ ssh-add ~/.ssh/id_ed25519
   $ ssh-add -L  # confirm the key appears
   ```

3. Inside the sandbox, configure Git to sign with SSH. Use the forwarded key
   directly rather than a file path, since host paths don't exist inside the
   sandbox:

   ```console
   $ git config --global gpg.format ssh
   $ git config --global user.signingkey "key::$(ssh-add -L | head -n 1)"
   ```

4. Sign commits as usual:

   ```console
   $ git commit -S -m "feat: my change"
   ```

To apply this configuration automatically to every sandbox, use the
[`git-ssh-sign`](https://github.com/docker/sbx-kits-contrib/tree/main/git-ssh-sign)
community kit, which handles all of the above setup. See [Kits](../customize/kits.md)
if you want to package it alongside other sandbox customizations.

For troubleshooting, see
[Sandbox commits aren't signed](../troubleshooting.md#sandbox-commits-arent-signed).
