---
title: "Git Tool"
description: "Read-only inspection of the working git repository."
keywords: docker agent, ai agents, tools, toolsets, git tool
linkTitle: "Git"
weight: 125
canonical: https://docs.docker.com/ai/docker-agent/tools/git/
---

_Read-only inspection of the working git repository._

## Overview

The git toolset gives an agent structured, **read-only** access to the working repository — status, history, branches, a commit's changes, and line-level authorship. It is implemented with go-git, so it needs **no `git` binary**.

Compared with running `git` through the `shell` tool, the git toolset returns clean, structured output the model can read reliably, is **safe by construction** (no command can modify the repository), and works even when `shell` is disabled or no `git` binary is installed.

> [!NOTE]
> The git toolset is read-only. To stage, commit, or check out, use the [`shell`](../shell/index.md) tool.

## Configuration

```yaml
toolsets:
  - type: git
```

No configuration options. The repository is opened at the agent's working directory; a subdirectory still resolves to the repository root.

> [!WARNING]
> **The repository is discovered by walking up parent directories.** If the working
> directory is not itself a repository but an ancestor is (for example a
> home directory tracked as dotfiles), the toolset resolves to that ancestor and
> `git_show` / `git_blame` can expose its full history and file contents. The
> filesystem toolset's allow/deny lists do **not** apply here. Only enable this
> toolset where the surrounding repository is safe to read.

> [!NOTE]
> **Performance.** go-git is pure Go, which costs speed on large repositories:
> `git_status` rehashes the whole worktree, and `git_blame` scales with history
> depth times file size — its 400-line output cap is applied *after* the full
> computation, so it does not make blaming a large file cheaper.

## Tools

| Tool | Description |
| --- | --- |
| `git_status` | Current branch and changed files (staged / unstaged / untracked). |
| `git_log` | Recent commits (hash, date, author, subject). |
| `git_branches` | Local branches, current one marked with `*`. |
| `git_show` | A commit's metadata, message, and changed files with +/- counts. |
| `git_blame` | Line-by-line authorship for a file. |

### `git_log`

| Parameter | Required | Description |
| --- | --- | --- |
| `limit` | No | Maximum number of commits to return (default 20). |
| `path` | No | Only show commits that touch this path. |

### `git_show`

| Parameter | Required | Description |
| --- | --- | --- |
| `ref` | No | Commit hash or revision to show (default HEAD). |

### `git_blame`

| Parameter | Required | Description |
| --- | --- | --- |
| `path` | Yes | File path to blame, relative to the repository root. |
| `rev` | No | Commit or revision to blame at (default HEAD). |

## Example

```yaml
agents:
  root:
    model: openai/gpt-5-mini
    description: A code review assistant
    instruction: |
      Review the working changes: check git_status, then git_show the latest
      commit, and summarize what changed.
    toolsets:
      - type: git
      - type: filesystem
```

Example `git_status` output:

```text
On branch master
1 changed file(s) [XY = staged/worktree; M=modified A=added D=deleted R=renamed ?=untracked]:
   M main.go
```

> [!TIP]
> **When to use**
>
> Use the git toolset whenever the agent needs repository context — before editing, to review recent history, or to find who last touched a line — without exposing the writable `shell` surface.
