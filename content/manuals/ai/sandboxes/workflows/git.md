---
title: Use Git with sandboxes
linkTitle: Git workflows
weight: 10
description: Choose where your Git checkout lives and manage branches, parallel tasks, and signed commits with Docker Sandboxes.
keywords: docker sandboxes, sbx, git, clone mode, direct mode, worktrees, branches, commit signing
---

These workspace modes apply to local sandboxes. In cloud sandboxes,
[transfer files or clone a remote repository](../cloud/usage.md#transfer-files).
To copy a sandbox filesystem between environments, see
[Move a sandbox](../cloud/move.md). Host mounts and clone-mode volumes are not
included in that snapshot.

Choose where your checkout lives based on how you want to review and run the
agent's changes:

- Share a host checkout with [direct mode](#direct-mode) to see edits in your
  local editor as they happen. A [host worktree](#host-worktree) gives a task
  its own branch and directory while keeping edits visible on the host.
- Use [clone mode](#clone-mode) to give the agent a separate checkout based on
  your local repository. Fetch or push committed work to bring it back.
- Use a [mountless sandbox](#mountless-sandbox) to clone repositories inside
  the sandbox without attaching a host workspace. Push branches to a Git
  remote to share the results.

| Checkout | Changes visible in host checkout | Agent can use Git |
| --- | --- | --- |
| Direct mount | Immediately | Yes |
| Mounted host worktree | Immediately, in the worktree | No |
| Clone mode (`--clone`) | After checking out or merging a fetched branch | Yes |
| Mountless sandbox | After pulling committed work from a Git remote | Yes |

For simultaneous tasks, use a separate worktree or clone for each task.
Branches alone don't isolate edits in a shared working directory. Worktrees
inside one sandbox share its tools and services. Use separate sandboxes when
tasks also need independent environments.

## Direct mode

The sandbox mounts your host working tree directly. The agent edits files in
place, and changes appear immediately. Both you and the agent can use Git,
so agree on who manages branches and commits. Switching branches from either
side changes the files both sides see.

To manage Git operations from the host:

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

Direct mode suits focused, single-branch work where you collaborate with the
agent turn-by-turn.

## Clone mode

In clone mode, `sbx` creates a separate Git clone inside the sandbox. The agent
edits this clone instead of your host working tree. Its changes stay inside the
sandbox until you fetch a branch or the agent pushes one to a remote. Your host
repository is also available at `/run/sandbox/source`, but only with read
access. The sandbox clone is not a Git worktree linked to your host checkout.

A single clone-mode sandbox can hold multiple branches and worktrees for
parallel tasks. The `--clone` flag creates the clone, but it doesn't separate
one task from another. To keep parallel tasks isolated, instruct your agent tool
to create a separate worktree or clone for each task.

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
   $ sbx run --clone claude .
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
   $ sbx run --clone claude .
   ```

2. Dispatch each independent task to a separate background session. Your agent
   tool controls how it separates their working directories. To require a
   separate worktree for each task, add a project instruction such as:

   ```markdown
   For each simultaneous task, create a separate Git worktree with its own branch
   and make all changes for that task in that worktree.
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

## Mountless sandbox

A mountless sandbox has no host workspace attached. Clone repositories inside
it to keep the checkout, dependencies, and development tools in the sandbox.
Files persist across stops and restarts, but removing the sandbox deletes them.

1. Create a named sandbox without a workspace path, then attach:

   ```console
   $ sbx create --name my-project claude
   $ sbx run --name my-project
   ```

   Omitting the path from `sbx create` creates a mountless sandbox. Omitting
   it from `sbx run claude` mounts your current host directory instead.
   See [Choose a workspace](../usage.md#choose-a-workspace).

2. Ask the agent to clone your repository inside the sandbox and create a
   branch before editing. For simultaneous tasks, ask it to create a separate
   worktree for each task.

3. Review and commit the changes inside the sandbox, then push the branch to
   your Git remote. You can fetch that branch in a host checkout or review it
   in a pull request. Push any work you want to keep before removing the sandbox.

You can also create a mountless sandbox with an `sbxenv.yaml` environment file.
Omit `workspace` and any `additionalWorkspaces` mounts, then run `sbx env run`.
See [Sandbox environment files](../configuration/environment-files.md#workspace)
for configuration options.

For running the application and opening a preview from your host browser, see
[Develop and test locally](development.md#accessing-services-in-the-sandbox).
To work with several host repositories in one sandbox instead, see
[Multiple workspaces](../usage.md#multiple-workspaces).

## Git ignore rules

Git inside a local sandbox preserves host-specific ignore decisions for paths
that exist when you create the sandbox. This applies when the primary workspace
is at a Git repository root.

The sandbox receives a snapshot of ignored paths, not the host-specific patterns.
Those patterns don't apply to files created later, but ignored directories still
cover files added inside them. Additional repositories and primary workspaces
below a repository root don't receive this snapshot. Repository `.gitignore`
files continue to work as usual.

Git ignore rules don't restrict the agent's access to files in the workspace.

## Commit signing

SSH agent forwarding is enabled by default. When `SSH_AUTH_SOCK` is set,
sandboxes forward your host SSH agent into the sandbox, so the agent can sign
commits with your SSH key without the private key ever leaving your host. If
you turned off forwarding or use a fixed SSH agent socket, see
[SSH agent configuration](../configuration/credentials.md#ssh-agent).

1. Make sure the signing key is loaded in your host SSH agent:

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
community kit, which handles all of the above setup. For using it with the
built-in agents, see [Kits v2](../customize/kits-v2.md).

For troubleshooting, see
[Sandbox commits aren't signed](../troubleshooting.md#sandbox-commits-arent-signed).
