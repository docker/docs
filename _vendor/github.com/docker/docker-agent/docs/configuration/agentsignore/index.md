---
title: "Ignoring files"
description: "Hide files from the agent with a .agentsignore file, using the same syntax as .gitignore."
keywords: docker agent, ai agents, configuration, agentsignore, gitignore, ignore, secrets
weight: 75
canonical: https://docs.docker.com/ai/docker-agent/configuration/agentsignore/
---

_Hide files from the agent with a `.agentsignore` file, using the same syntax as `.gitignore`._

## Overview

Put a `.agentsignore` file in your project and the agent stops seeing the paths it lists. Matched files are absent from directory listings and searches, and reads, writes and edits targeting them are refused.

```text
# .agentsignore
secrets.env
*.key
.env.*
build/
docs/**/*.draft.md
!public.key
```

There is nothing to configure. The file's presence is the opt-in, and it is picked up automatically by the [filesystem toolset](../../tools/filesystem/index.md).

## Syntax

The syntax is `.gitignore` syntax — parsed with the same library git uses, so patterns behave identically to `.gitignore` and `.dockerignore`:

| Pattern | Matches |
| --- | --- |
| `secrets.env` | that name at any depth (`secrets.env`, `config/secrets.env`) |
| `/secrets.env` | that name at the project root only |
| `*.key` | any file with the extension |
| `build/` | the directory and everything under it |
| `docs/**/*.draft.md` | nested matches via `**` |
| `!public.key` | re-includes a path an earlier pattern excluded |
| `# comment` | ignored, as are blank lines |

## Where the file is found

The nearest `.agentsignore` at or above the working directory is used, so starting a run in a subdirectory still honours the project's file. Patterns are anchored to the directory containing the file, exactly as git anchors to the directory containing `.gitignore`.

Unlike `.gitignore`, **a git repository is not required** — `.agentsignore` works in any directory.

## What it affects

| Behaviour | Effect |
| --- | --- |
| `list_directory`, `directory_tree` | matched entries are omitted |
| `search_files_content` | matched files are skipped |
| `read_file`, `read_multiple_files` | refused |
| `write_file`, `edit_file` | refused, including for files that do not exist yet |
| `create_directory`, `remove_directory` | refused |
| `permissions` | matching deny rules are derived automatically, so `/permissions` shows them |

Paths are resolved before matching — symlinks, `./` prefixes and `..` segments are all normalised — so an ignored file cannot be reached by spelling it differently.

The `.agentsignore` file is itself always hidden: it names the very things being kept back, so handing it to the agent would be a map of what to look for.

> [!NOTE]
> `.agentsignore` is independent of the filesystem toolset's [`ignore_vcs`](../tools/index.md) option. Setting `ignore_vcs: false` turns off `.gitignore` filtering but does **not** un-hide `.agentsignore` entries.

## Relationship to `.gitignore`

They are separate, and they do different amounts of work.

`.gitignore` is respected by default (`ignore_vcs`), but only as a **display filter**: gitignored files are hidden from listings and searches while `read_file` still opens them. That is reasonable for build output, which is noise rather than secret.

`.agentsignore` is for content the agent should not have at all, so it blocks reads and writes as well. Use `.gitignore` for noise, `.agentsignore` for secrets.

## Limits

> [!WARNING]
> `.agentsignore` governs the filesystem toolset. It is **not** a sandbox.
>
> An agent with the [`shell`](../../tools/shell/index.md) toolset can still run `cat secrets.env`, because the shell runs commands the toolset never inspects. The same applies to any toolset that reaches the filesystem on its own, such as [`lsp`](../../tools/lsp/index.md).
>
> Treat `.agentsignore` as a strong default that keeps sensitive files out of the agent's view and context — not as a boundary against an agent actively trying to reach them. When you need a real boundary, combine it with [permissions](../permissions/index.md) that restrict `shell`, or run in [sandbox mode](../sandbox/index.md).

The derived permission rules are best-effort for the same reason: permission patterns match the argument string as the model wrote it, without resolving it first, so `./secrets.env` can slip past a rule written for `secrets.env`. The filesystem toolset's own check resolves paths first and is the part that actually enforces.

## Example

```text
# .agentsignore

# Secrets
.env
.env.*
secrets.env
*.pem
*.key
!public.key          # this one is safe to read

# Credentials directories
.aws/
.ssh/

# Large build output the agent doesn't need
build/
dist/
node_modules/
```
