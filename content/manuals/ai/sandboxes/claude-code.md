---
title: Configure Claude Code
description: Learn how to configure Claude Code authentication, pass CLI options, and customize your sandboxed agent environment with Docker.
weight: 30
---

{{< summary-bar feature_name="Docker Sandboxes" >}}

This guide covers authentication, configuration files, and common options for
running Claude Code in a sandboxed environment.

## Quick start

To create a sandbox and run Claude Code for a project directory:

```console
$ docker sandbox run claude ~/my-project
```

The workspace parameter is optional and defaults to the current directory:

```console
$ cd ~/my-project
$ docker sandbox run claude
```

### Pass a prompt directly

Start Claude with a specific prompt:

```console
$ docker sandbox run <sandbox-name> -- "Add error handling to the login function"
```

Or:

```console
$ docker sandbox run <sandbox-name> -- "$(cat prompt.txt)"
```

This starts Claude and immediately processes the prompt.

## Authentication

Claude Code requires an Anthropic API key. Credentials are scoped per sandbox
and must be provided through environment variables or interactive login.

### Environment variable (recommended)

The recommended approach is to set the `ANTHROPIC_API_KEY` environment variable in your shell configuration file.

Docker Sandboxes use a daemon process that doesn't inherit environment
variables from your current shell session. To make your API key available to
sandboxes, set it globally in your shell configuration file.

Add the API key to your shell configuration file:

```plaintext {title="~/.bashrc or ~/.zshrc"}
export ANTHROPIC_API_KEY=sk-ant-api03-xxxxx
```

Apply the changes:

1. Source your shell configuration: `source ~/.bashrc` (or `~/.zshrc`)
2. Restart Docker Desktop so the daemon picks up the new environment variable
3. Create and run your sandbox:

```console
$ docker sandbox create claude ~/project
$ docker sandbox run <sandbox-name>
```

The sandbox detects the environment variable and uses it automatically.

### Interactive authentication

If the `ANTHROPIC_API_KEY` environment variable is not set, Claude Code prompts
you to authenticate interactively when it starts. You can also trigger the login
flow manually using the `/login` command within Claude Code.

When using interactive authentication:

- You must authenticate each sandbox separately
- If the sandbox is removed or destroyed, you'll need to authenticate again when you recreate it
- Authentication sessions aren't persisted outside the sandbox
- No fallback authentication methods are used

To avoid repeated authentication, use the `ANTHROPIC_API_KEY` environment variable method described above.

## Configuration

Claude Code can be configured through CLI options. Any arguments you pass after
the sandbox name and a `--` separator are passed directly to Claude Code.

Pass options after the sandbox name:

```console
$ docker sandbox run <sandbox-name> -- [claude-options]
```

For example:

```console
$ docker sandbox run <sandbox-name> -- --continue
```

See the [Claude Code CLI reference](https://docs.claude.com/en/docs/claude-code/cli-reference)
for available options.

## Base image

The Claude Code sandbox template is a container image that runs inside the
sandbox VM. It includes:

- Ubuntu-based environment with Claude Code
- Development tools: Docker CLI, GitHub CLI, Node.js, Go, Python 3, Git, ripgrep, jq
- Non-root `agent` user with sudo access
- Private Docker daemon for running additional containers

Claude launches with `--dangerously-skip-permissions` by default in sandboxes.

You can build custom templates based on `docker/sandbox-templates:claude-code`.
See [Custom templates](templates.md) for details.

## Next steps

- [Using sandboxes effectively](workflows.md)
- [Custom templates](templates.md)
- [Network policies](network-policies.md)
- [Troubleshooting](troubleshooting.md)
- [CLI Reference](/reference/cli/docker/sandbox/)
