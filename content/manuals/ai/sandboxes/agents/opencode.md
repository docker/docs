---
title: OpenCode sandbox
description: |
  Use OpenCode in Docker Sandboxes with multi-provider authentication and TUI
  interface for AI development.
keywords: docker, sandboxes, opencode, ai agent, multi-provider, authentication, tui
weight: 50
---

{{< summary-bar feature_name="Docker Sandboxes v0.12" >}}

This guide covers authentication, configuration, and usage of OpenCode in a
sandboxed environment.

Official documentation: [OpenCode](https://opencode.ai/docs)

## Quick start

Create a sandbox and run OpenCode for a project directory:

```console
$ docker sandbox run opencode ~/my-project
```

The workspace parameter is optional and defaults to the current directory:

```console
$ cd ~/my-project
$ docker sandbox run opencode
```

OpenCode launches a TUI (text user interface) where you can select your
preferred LLM provider and interact with the agent.

## Authentication

OpenCode uses proxy-managed authentication for all supported providers. Docker
Sandboxes intercepts API requests and injects credentials transparently. You
provide your API keys through environment variables on the host, and the
sandbox handles credential management.

### Supported providers

Configure one or more providers by setting environment variables:

```plaintext {title="~/.bashrc or ~/.zshrc"}
export OPENAI_API_KEY=sk-xxxxx
export ANTHROPIC_API_KEY=sk-ant-xxxxx
export GOOGLE_API_KEY=AIzaSyxxxxx
export XAI_API_KEY=xai-xxxxx
export GROQ_API_KEY=gsk_xxxxx
export AWS_ACCESS_KEY_ID=AKIA_xxxxx
export AWS_SECRET_ACCESS_KEY=xxxxx
export AWS_REGION=us-west-2
```

You only need to configure the providers you want to use. OpenCode detects
available credentials and offers those providers in the TUI.

### Environment variable setup

Docker Sandboxes use a daemon process that doesn't inherit environment
variables from your current shell session. To make your API keys available to
sandboxes, set them globally in your shell configuration file.

Apply the changes:

1. Source your shell configuration: `source ~/.bashrc` (or `~/.zshrc`)
2. Restart Docker Desktop so the daemon picks up the new environment variables
3. Create and run your sandbox:

```console
$ docker sandbox create opencode ~/project
$ docker sandbox run <sandbox-name>
```

The sandbox detects the environment variables and uses them automatically.

## Configuration

OpenCode uses a TUI interface and doesn't require extensive configuration
files. The agent prompts you to select a provider when it starts, and you can
switch providers during a session.

### TUI mode

OpenCode launches in TUI mode by default. The interface shows:

- Available LLM providers (based on configured credentials)
- Current conversation history
- File operations and tool usage
- Real-time agent responses

Use keyboard shortcuts to navigate the interface and interact with the agent.

## Base image

Template: `docker/sandbox-templates:opencode`

OpenCode supports multiple LLM providers with automatic credential injection
through the sandbox proxy.

See [Custom templates](../templates.md) to build your own agent images.
