---
title: "Snapshots"
description: "Shadow-git snapshots capture your workspace at turn boundaries so you can review what an agent changed and undo it without touching your real git history."
keywords: docker agent, ai agents, features, snapshots
weight: 20
canonical: https://docs.docker.com/ai/docker-agent/features/snapshots/
---

_Shadow-git snapshots capture your workspace at turn boundaries so you can review what an agent changed and undo it without touching your real git history._

## Overview

When snapshots are enabled, docker-agent records the state of your working
directory as the agent runs. Each checkpoint is stored in a **shadow git
repository** kept under the docker-agent data directory — completely separate
from your project's own `.git`. This lets you:

- **Review** exactly which files an agent touched on a given turn.
- **Undo** the most recent change with `/undo`.
- **Reset** the workspace back to any earlier checkpoint — or all the way to
  its pre-agent state.

Snapshots only ever touch files on disk. They never write commits to your
repository, never move your `HEAD`, and never appear in `git log` or
`git status`.

> [!NOTE]
> **Git repositories only**
>
> The snapshot machinery is a no-op outside a git worktree. It needs a git
> repository to scope which files belong to your project and to mirror the
> repository's ignore rules.

## Enabling Snapshots

The simplest way is to enable snapshots globally for every agent in your user
config:

```yaml
# ~/.config/cagent/config.yaml
settings:
  snapshot: true
```

Omit `snapshot` or set it to `false` to leave automatic snapshots off.

You can also wire the [`snapshot` built-in hook](../../configuration/hooks/index.md#available-built-ins)
into a single agent instead of enabling it globally:

```yaml
hooks:
  turn_start:
    - type: builtin
      command: snapshot
  turn_end:
    - type: builtin
      command: snapshot
  session_end:
    - type: builtin
      command: snapshot
```

Manually configured snapshot hooks always run, even when `settings.snapshot`
is unset. See [`examples/snapshot_hooks.yaml`](https://github.com/docker/docker-agent/blob/main/examples/snapshot_hooks.yaml)
for a complete configuration.

## Using Snapshots in the TUI

When snapshots are enabled, the [TUI](../tui/index.md)
exposes two slash commands:

- **`/undo`** restores files from the most recent snapshot (one step back).
- **`/snapshots`** opens a dialog listing every captured snapshot and the
  number of files in each. Use <kbd>↑</kbd>/<kbd>↓</kbd> (or
  <kbd>j</kbd>/<kbd>k</kbd>) to highlight an entry, then press <kbd>r</kbd> to
  reset the workspace to that point. Pick `<original>` to revert every snapshot
  and bring the workspace back to its pre-agent state. <kbd>Esc</kbd> closes the
  dialog without changing anything.

Neither command removes messages from the session transcript — they only touch
files on disk. Both commands (and their command-palette entries) are hidden
when snapshots are turned off.

## How It Works

- **Shadow repository.** The first time a snapshot is taken for a worktree,
  docker-agent initializes a separate shadow git directory under the data
  directory (`~/.cagent/snapshot/...` by default), keyed by a hash of the
  worktree path. The shadow repo stores tree objects only — it never writes
  commits and never touches your source repository's `.git`.
- **Ignore rules are mirrored.** Before each capture, the source repository's
  `.gitignore` and `info/exclude` rules are mirrored into the shadow repo so
  ignored files never appear in snapshots.
- **Large files are skipped.** Newly-added files larger than 2 MiB are excluded
  from snapshots to keep the shadow repo small.
- **Checkpoints only on change.** A checkpoint is recorded only when files
  actually changed, so a final no-op model response does not hide the last
  meaningful snapshot.
- **Scoped to the working directory.** Snapshot operations are scoped to the
  agent's working directory within the worktree, so a sub-directory agent only
  captures and restores files under that directory.
- **Garbage collection.** Wiring `snapshot` into `session_end` (or enabling it
  globally) runs `git gc` against the shadow repo so old, unreferenced objects
  are pruned over time.

## See Also

- [Hooks](../../configuration/hooks/index.md) — the `snapshot`
  built-in and the events it can run on.
- [Terminal UI](../tui/index.md) — the `/undo` and
  `/snapshots` commands.
- [`examples/snapshot_hooks.yaml`](https://github.com/docker/docker-agent/blob/main/examples/snapshot_hooks.yaml) — a complete snapshot hook configuration.
