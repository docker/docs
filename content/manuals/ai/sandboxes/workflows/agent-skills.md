---
title: Share agent skills
linkTitle: Agent skills
weight: 40
description: Add skills from Git repositories or import them from supported host agents into a store shared with Docker Sandboxes.
keywords: docker sandboxes, sbx, agent skills, shared skills, git repository, claude code, codex, copilot, cursor, devin, droid
---

Shared agent skills are reusable instruction packs for coding agents (for
example Claude Code, Codex, Copilot, Cursor, or Devin).

You can:

- Install skills from a Git repository into Docker Sandboxes
- Import skills that are already installed for supported agents on your host

Imported and installed skills are copied into a persistent store on the host.
That store survives sandbox deletion, and new sandboxes that run a supported
agent can use those skills by default.

> [!NOTE]
> Shared agent skills are experimental.

## Add skills from a repository

Add every skill from a Git repository:

```console
$ sbx skills add anthropics/skills
```

The repository must contain one or more valid `SKILL.md` files. You can use
GitHub `owner/repository` shorthand or a Git URL, including HTTPS and SSH URLs.

To add specific skills, use `--skill` with each name or pass a comma-separated
list:

```console
$ sbx skills add https://github.com/anthropics/skills --skill frontend-design --skill pdf
```

When a skill with the same name is already installed, `sbx` prompts before
replacing it. Use `--force` to replace existing skills without prompts.

## Manage installed skills

List the skills in the shared store:

```console
$ sbx skills ls
```

Update every skill installed from a repository:

```console
$ sbx skills update
```

To update specific skills, pass one or more names:

```console
$ sbx skills update frontend-design pdf
```

`sbx skills update` only refreshes skills installed with `sbx skills add`. To
refresh a skill imported from the host, run `sbx skills import` again. To manage
it with `sbx skills update`, install it from a repository instead.

Remove one or more installed skills:

```console
$ sbx skills rm frontend-design pdf
```

`sbx` asks for confirmation before removing skills that running agents may be
using. Use `--force` to skip confirmation in scripts.

## Import skills from the host

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
| Codex and Devin | `~/.agents/skills` | `/home/agent/.agents/skills` |
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
again when you want to copy updates from the host. If an import replaces a
repository-installed skill, `sbx` no longer associates that skill with its
repository, so `sbx skills update` won't refresh it.

## Shared store behavior

Running `sbx reset` clears the shared store.

Sandboxes created for a supported agent mount the shared store read-only by
default. These sandboxes mount the contents of the store each time they start,
so you can install skills before or after creating them.

Use `--skills` with `sbx run` or `sbx create` to choose the access mode when
creating a sandbox:

- `readonly`: Mount the store so the agent can read skills but cannot modify them.
- `readwrite`: Mount the store so the agent can read and modify shared skills.
- `off`: Omit the shared store mount.

For example, create a sandbox without the shared store:

```console
$ sbx run --skills=off claude
```

To change the default for future sandboxes, set `skills.defaultMode` to `off`,
`readonly`, or `readwrite`:

```console
$ sbx settings set skills.defaultMode readonly
```

When no mode is specified, the daemon uses `skills.defaultMode`, whose built-in
value is `readonly`. An explicit `--skills` value overrides that default.

The mode is applied only when a sandbox is created. Upgrading `sbx` or changing
`skills.defaultMode` leaves existing sandbox mounts unchanged. Remove and
recreate a sandbox to change its mode. Sandboxes created without shared skills
also need to be recreated to mount the store.

> [!WARNING]
> A sandbox with `readwrite` access can modify skills that other sandboxes load,
> including sandboxes with `readonly` access. Read-only access prevents writes
> from that sandbox but does not isolate it from changes to the store. The store
> is dedicated sandbox state, so this does not by itself execute modified skills
> on your host. Use `--skills=off` when creating a sandbox to keep it outside
> this shared trust boundary.

Some agents scan for skills when a session starts. If installed skills don't
appear in an existing session, start another agent session.
