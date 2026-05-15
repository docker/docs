---
title: Workspace trust
weight: 40
description: |
  How sandboxed agents interact with your workspace files and what to review
  after an agent session.
keywords: docker sandboxes, workspace trust, file access, review, sbx
---

{{< summary-bar feature_name="Docker Sandboxes sbx" >}}

Agents running in sandboxes have full access to the workspace directory without
prompting. With the default direct mount, changes the agent makes appear on
your host immediately. Treat sandbox-modified workspace files the same way
you would treat a pull request from an untrusted contributor: review before
you trust them on your host.

## What the agent can modify

The agent can create, modify, and delete any file in the workspace. This
includes:

- Source code files
- Configuration files (`.eslintrc`, `pyproject.toml`, `.env`, etc.)
- Build files (`Makefile`, `package.json`, `Cargo.toml`)
- Git hooks (`.git/hooks/`)
- CI configuration (`.github/workflows/`, `.gitlab-ci.yml`)
- IDE configuration (`.vscode/tasks.json`, `.idea/` run configurations)
- Hidden files and directories
- Shell scripts and executables

> [!CAUTION]
> Files like Git hooks, CI configuration, IDE task configs, and build scripts
> execute code when triggered by normal development actions such as committing,
> building, or opening the project in an IDE. Review these files after any agent
> session before performing those actions.

## Branch mode

The `--branch` flag lets the agent work on a separate branch. This is a
workflow convenience, not a security boundary: the agent still mounts the full
repository. See the [usage guide](../usage.md) for details.

## Reviewing changes

After an agent session, review changes before executing any code the agent
touched.

With the default direct mount, changes are in your working tree:

```console
$ git diff
```

If you used `--branch`, the agent's changes are on a separate branch:

```console
$ git diff main..my-feature
```

Pay particular attention to:

- **Git hooks** (`.git/hooks/`): run on commit, push, and other Git actions.
  These are inside `.git/` and **do not appear in `git diff` output**. Check
  them separately with `ls -la .git/hooks/`.
- **CI configuration** (`.github/workflows/`, `.gitlab-ci.yml`): runs on push
- **Build files** (`Makefile`, `package.json` scripts, `Cargo.toml`): run
  during build or install steps
- **IDE configuration** (`.vscode/tasks.json`, `.idea/`): can run tasks when
  you open the project
- **Executable files and shell scripts**: can run directly

These files execute code without you explicitly running them. Review them before
committing, building, or opening the project in an IDE.
