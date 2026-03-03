---
title: Codex sandbox
description: |
  Use OpenAI Codex in Docker Sandboxes with API key authentication and YOLO
  mode configuration.
keywords: docker, sandboxes, codex, openai, ai agent, authentication, yolo mode
weight: 20
---

{{< summary-bar feature_name="Docker Sandboxes" >}}

This guide covers authentication, configuration, and usage of Codex in a
sandboxed environment.

Official documentation: [Codex CLI](https://developers.openai.com/codex/cli)

## Quick start

Create a sandbox and run Codex for a project directory:

```console
$ docker sandbox run codex ~/my-project
```

The workspace parameter is optional and defaults to the current directory:

```console
$ cd ~/my-project
$ docker sandbox run codex
```

## Authentication

Codex requires an OpenAI API key. Credentials are scoped per sandbox.

Set the `OPENAI_API_KEY` environment variable in your shell configuration file.

Docker Sandboxes use a daemon process that doesn't inherit environment
variables from your current shell session. To make your API key available to
sandboxes, set it globally in your shell configuration file.

Add the API key to your shell configuration file:

```plaintext {title="~/.bashrc or ~/.zshrc"}
export OPENAI_API_KEY=sk-xxxxx
```

Apply the changes:

1. Source your shell configuration: `source ~/.bashrc` (or `~/.zshrc`)
2. Restart Docker Desktop so the daemon picks up the new environment variable
3. Create and run your sandbox:

```console
$ docker sandbox create codex ~/project
$ docker sandbox run <sandbox-name>
```

The sandbox detects the environment variable and uses it automatically.

## Configuration

Codex supports a YOLO mode that disables safety checks and approval prompts.
This mode grants the agent full access to your sandbox environment without
interactive confirmation.

Configure YOLO mode in `~/.codex/config.toml`:

```toml
approval_policy = "never"
sandbox_mode = "danger-full-access"
```

With these settings, Codex runs without approval prompts.

### Pass options at runtime

Pass Codex CLI options after the sandbox name and a `--` separator:

```console
$ docker sandbox run <sandbox-name> -- --dangerously-bypass-approvals-and-sandbox
```

This flag enables YOLO mode for a single session without modifying the
configuration file.

## Base image

Template: `docker/sandbox-templates:codex`

Codex launches with `--dangerously-bypass-approvals-and-sandbox` by default when YOLO mode is configured.

See [Custom templates](../templates.md) to build your own agent images.
