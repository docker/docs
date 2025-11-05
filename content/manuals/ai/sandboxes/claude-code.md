---
title: Configure Claude Code
description: Learn how to configure Claude Code authentication, pass CLI options, and customize your sandboxed agent environment with Docker.
weight: 30
---

{{< summary-bar feature_name="Docker Sandboxes" >}}

This guide covers authentication, configuration files, and common options for
running Claude Code in a sandboxed environment.

## Quick start

The simplest way to start Claude in a sandbox:

```console
$ docker sandbox run claude
```

This starts a sandboxed Claude Code agent with the current working directory as
its workspace.

Or specify a different workspace:

```console
$ docker sandbox run -w ~/my-project claude
```

## Passing CLI options to Claude

Claude Code supports various command-line options that you can pass through
`docker sandbox run`. Any arguments after the agent name (`claude`) are passed
directly to Claude Code inside the sandbox.

### Continue previous conversation

Resume your most recent conversation:

```console
$ docker sandbox run claude -c
```

Or use the long form:

```console
$ docker sandbox run claude --continue
```

### Pass a prompt directly

Start Claude with a specific prompt:

```console
$ docker sandbox run claude "Add error handling to the login function"
```

This starts Claude and immediately processes the prompt.

### Combine options

You can combine sandbox options with Claude options:

```console
$ docker sandbox run -e DEBUG=1 claude -c
```

This creates a sandbox with `DEBUG` set to `1`, enabling debug output for
troubleshooting, and continues the previous conversation.

### Available Claude options

All Claude Code CLI options work through `docker sandbox run`:

- `-c, --continue` - Continue the most recent conversation
- `-p, --prompt` - Read prompt from stdin (useful for piping)
- `--dangerously-skip-permissions` - Skip permission prompts (enabled by default in sandboxes)
- And more - see the [Claude Code documentation](https://docs.claude.com/en/docs/claude-code) for a complete list

## Authentication

Claude sandboxes support the following credential management strategies.

### Strategy 1: `sandbox` (Default)

```console
$ docker sandbox run claude
```

On first run, Claude prompts you to enter your Anthropic API key. The
credentials are stored in a persistent Docker volume named
`docker-claude-sandbox-data`. All future Claude sandboxes automatically use
these stored credentials, and they persist across sandbox restarts and deletion.

Sandboxes mount this volume at `/mnt/claude-data` and create symbolic links in
the sandbox user's home directory.

> [!NOTE]
> If your workspace contains a `.claude.json` file with a `primaryApiKey`
> field, you'll receive a warning about potential conflicts. You can choose to
> remove the `primaryApiKey` field from your `.claude.json` or proceed and
> ignore the warning.

### Strategy 2: `none`

No automatic credential management.

```console
$ docker sandbox run --credentials=none claude
```

Docker does not discover, inject, or store any credentials. You must
authenticate manually inside the container. Credentials are not shared with
other sandboxes but persist for the lifetime of the container.

## Configuration

Claude Code can be configured through CLI options. Any arguments you pass after
the agent name are passed directly to Claude Code inside the container.

Pass options after the agent name:

```console
$ docker sandbox run claude [claude-options]
```

For example:

```console
$ docker sandbox run claude --continue
```

See the [Claude Code CLI reference](https://docs.claude.com/en/docs/claude-code/cli-reference)
for a complete list of available options.

## Advanced usage

For more advanced configurations including environment variables, volume mounts,
Docker socket access, and custom templates, see
[Advanced configurations](advanced-config.md).

## Base image

The `docker/sandbox-templates:claude-code` image includes Claude Code with
automatic credential management, plus development tools (Docker CLI, GitHub
CLI, Node.js, Go, Python 3, Git, ripgrep, jq). It runs as a non-root `agent`
user with `sudo` access and launches Claude with
`--dangerously-skip-permissions` by default.

## Next Steps

- [Advanced configurations](advanced-config.md)
- [Troubleshooting](troubleshooting.md)
- [CLI Reference](/reference/cli/docker/sandbox/)
