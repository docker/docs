---
title: cagent sandbox
description: |
  Use Docker cagent in Docker Sandboxes with multi-provider authentication
  supporting OpenAI, Anthropic, and more.
keywords: docker, sandboxes, cagent, ai agent, multi-provider, authentication
weight: 60
---

{{< summary-bar feature_name="Docker Sandboxes" >}}

This guide covers authentication, configuration, and usage of Docker cagent in
a sandboxed environment. [cagent](/ai/cagent/) is Docker's open source coding
agent that supports multiple providers.

## Quick start

Create a sandbox and run cagent for a project directory:

```console
$ docker sandbox run cagent ~/my-project
```

The workspace parameter is optional and defaults to the current directory:

```console
$ cd ~/my-project
$ docker sandbox run cagent
```

## Authentication

cagent uses proxy-managed authentication for all supported providers. Docker
Sandboxes intercepts API requests and injects credentials transparently. You
provide your API keys through environment variables, and the sandbox handles
credential management.

### Supported providers

Configure one or more providers by setting environment variables:

```plaintext {title="~/.bashrc or ~/.zshrc"}
export OPENAI_API_KEY=sk-xxxxx
export ANTHROPIC_API_KEY=sk-ant-xxxxx
export GOOGLE_API_KEY=AIzaSyxxxxx
export XAI_API_KEY=xai-xxxxx
export NEBIUS_API_KEY=xxxxx
export MISTRAL_API_KEY=xxxxx
```

You only need to configure the providers you want to use. cagent detects
available credentials and routes requests to the appropriate provider.

### Environment variable setup

Docker Sandboxes use a daemon process that doesn't inherit environment
variables from your current shell session. To make your API keys available to
sandboxes, set them globally in your shell configuration file.

Apply the changes:

1. Source your shell configuration: `source ~/.bashrc` (or `~/.zshrc`)
2. Restart Docker Desktop so the daemon picks up the new environment variables
3. Create and run your sandbox:

```console
$ docker sandbox create cagent ~/project
$ docker sandbox run <sandbox-name>
```

The sandbox detects the environment variables and uses them automatically.

## Configuration

cagent supports YOLO mode that disables safety checks and approval prompts.
This mode grants the agent full access to your sandbox environment without
interactive confirmation.

### Pass options at runtime

Pass cagent CLI options after the sandbox name and a `--` separator:

```console
$ docker sandbox run <sandbox-name> -- run --yolo
```

The `run --yolo` command starts cagent with approval prompts disabled.

## Base image

Template: `docker/sandbox-templates:cagent`

cagent supports multiple LLM providers with automatic credential injection
through the sandbox proxy. Launches with `run --yolo` by default.

See [Custom templates](../templates.md) to build your own agent images.
