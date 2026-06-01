---
title: Claude Code
weight: 10
description: |
  Use Claude Code in Docker Sandboxes with authentication, configuration, and
  YOLO mode for AI-assisted development.
keywords: docker sandboxes, claude code, anthropic, ai agent, sbx
---

Official documentation: [Claude Code](https://code.claude.com/docs)

## Quick start

Launch Claude Code in a sandbox by pointing it at a project directory:

```console
$ sbx run claude ~/my-project
```

The workspace parameter defaults to the current directory, so `sbx run claude`
from inside your project works too. To start Claude with a specific prompt:

```console
$ sbx run claude --name my-sandbox -- "Add error handling to the login function"
```

Everything after `--` is passed directly to Claude Code. You can also pipe in a
prompt from a file with `-- "$(cat prompt.txt)"`.

## Authentication

Claude Code requires either an Anthropic API key or a Claude subscription.

**API key**: Store your key using
[stored secrets](../security/credentials.md#stored-secrets):

```console
$ sbx secret set -g anthropic
```

Alternatively, export the `ANTHROPIC_API_KEY` environment variable in your
shell before running the sandbox. See
[Credentials](../security/credentials.md) for details on both methods.

**Claude subscription**: If no API key is set, Claude Code prompts you to
authenticate interactively using OAuth. The proxy handles the OAuth flow, so
credentials aren't stored inside the sandbox.

## Configuration

Sandboxes don't pick up user-level configuration from your host, such as
`~/.claude`. Only project-level configuration in the working directory is
available inside the sandbox. See
[Why doesn't the sandbox use my user-level agent configuration?](../faq.md#why-doesnt-the-sandbox-use-my-user-level-agent-configuration)
for workarounds.

### Default startup command

Without extra args, the sandbox runs:

```text
claude --dangerously-skip-permissions
```

Args after `--` replace these defaults rather than being appended. To keep
`--dangerously-skip-permissions`, include it yourself:

```console
$ sbx run claude -- --dangerously-skip-permissions -c
```

See the [Claude Code CLI reference](https://code.claude.com/docs/en/cli-reference)
for available options.

## Agents view

Claude Code's [agents view](https://code.claude.com/docs/en/agent-view)
dispatches tasks to subagents that work in parallel, each in its own
Git worktree. Pair it with [clone mode](../usage.md#clone-mode) for an
isolated multi-agent workflow:

```console
$ sbx run --clone claude -- agents
```

This invocation replaces the
[default startup command](#default-startup-command), so it doesn't
include `--dangerously-skip-permissions` and you can't switch to
bypass-permissions mode inside the sandbox. To work around this, either
use Claude Code's auto mode or pass the flag explicitly:

```console
$ sbx run --clone claude -- --dangerously-skip-permissions agents
```

The subagents' worktrees live inside the sandbox's private clone — none
of them touches your host repository. Each subagent commits to its own
branch, and you review the work from the host by fetching the
`sandbox-<sandbox-name>` remote:

```console
$ git fetch sandbox-<sandbox-name>
$ git diff main..sandbox-<sandbox-name>/<branch>
```

See [Git workflow](../usage.md#git-workflow) for clone-mode details.

## Base image

The sandbox uses `docker/sandbox-templates:claude-code`. See
[Templates](../customize/templates.md) to build your own image on top of
this base.

## Use a local model

To run Claude Code in a sandbox against a local model on your host through
Docker Model Runner, see
[Run Claude Code in a Docker Sandbox with Docker Model Runner](/guides/claude-code-sandbox-model-runner/).
For the host-only version without a sandbox, see
[Use Claude Code with Docker Model Runner](/guides/claude-code-model-runner/).
