---
title: Claude Code
weight: 10
description: |
  Use Claude Code in Docker Sandboxes with authentication, local models,
  configuration, and YOLO mode for AI-assisted development.
keywords: docker sandboxes, claude code, anthropic, ai agent, sbx, local models, llmman, ollama
---

Official documentation: [Claude Code](https://code.claude.com/docs)

## Quick start

Launch Claude Code in a sandbox by pointing it at a project directory:

```console
$ sbx run claude ~/my-project
```

To start Claude with a specific prompt in the current directory:

```console
$ sbx run --name my-sandbox claude -- "Add error handling to the login function"
```

Everything after `--` is passed directly to Claude Code. You can also pipe in a
prompt from a file with `-- "$(cat prompt.txt)"`.

To create a [mountless sandbox](../usage.md#choose-a-workspace), use
`sbx create` without a workspace path, then attach by name.

## Authentication

For the default Anthropic models, Claude Code requires either an Anthropic
API key or a Claude subscription. For other models, see
[Use a local model](#use-a-local-model).

**API key**: Store your key using
[stored secrets](../configuration/credentials.md#stored-secrets):

```console
$ sbx secret set anthropic
```

**Claude subscription**: If no API key is set, use the `/login` command inside
Claude Code to authenticate via OAuth.

## Configuration

Sandboxes don't pick up user-level configuration from your host, such as
`~/.claude`. Only project-level configuration in the working directory is
available inside the sandbox. See
[Why doesn't the sandbox use my user-level agent configuration?](../faq.md#why-doesnt-the-sandbox-use-my-user-level-agent-configuration)
for workarounds.

### Remote control

To use Claude Code's `/remote-control` command inside a sandbox, turn on
[`claude.remoteControl`](../configuration/settings.md#clauderemotecontrol):

```console
$ sbx settings set claude.remoteControl true
```

### Default startup command

Without extra args, the sandbox runs:

```text
claude --dangerously-skip-permissions
```

Arguments after `--` are added after the default flags when the first one is
itself a flag (begins with `-`), so `--dangerously-skip-permissions` is
preserved:

```console
$ sbx run --name <sandbox-name> -- -c   # runs claude --dangerously-skip-permissions -c
```

When the first argument is a bare word, such as the `agents` subcommand, it
replaces the defaults instead.

See the [Claude Code CLI reference](https://code.claude.com/docs/en/cli-reference)
for available options.

## Agents view

Claude Code's [agents view](https://code.claude.com/docs/en/agent-view)
starts background sessions that run tasks in parallel. Pair it with
[clone mode](../workflows/git.md#clone-mode) to keep their changes inside the
sandbox:

```console
$ sbx run --clone claude . -- agents
```

This invocation replaces the
[default startup command](#default-startup-command), so it doesn't
include `--dangerously-skip-permissions` and you can't switch to
bypass-permissions mode inside the sandbox. To work around this, either
use Claude Code's auto mode or pass the flag explicitly:

```console
$ sbx run --clone claude . -- --dangerously-skip-permissions agents
```

Claude Code may use branches or worktrees to keep changes from its background
sessions separate. This depends on the task, Claude Code configuration, and
project instructions. The `--clone` flag doesn't control this behavior. Claude
Code creates any branches and worktrees inside the sandbox, not in your host
checkout.

To review a branch created by a session, fetch the
`sandbox-<sandbox-name>` remote from the host:

```console
$ git fetch sandbox-<sandbox-name>
$ git diff main..sandbox-<sandbox-name>/<branch>
```

See [Git workflows](../workflows/git.md) for clone-mode details.

## Base image

The sandbox uses `docker/sandbox-templates:claude-code`. See
[Templates](../customize/templates.md) to build your own image on top of
this base.

## Use a local model

For local models, hosted providers, and custom inference endpoints, see
[Use local and hosted models](../configuration/models.md).
