---
title: Custom shell
description: |
  Use the custom shell sandbox for manual agent installation and custom
  development environments in Docker Sandboxes.
keywords: docker, sandboxes, shell, custom, manual setup, development environment
weight: 80
---

{{< summary-bar feature_name="Docker Sandboxes" >}}

This guide covers the Shell sandbox, a minimal environment for custom agent
installation and development. Unlike other agent sandboxes, Shell doesn't
include a pre-installed agent binary. Instead, it provides a clean environment
where you can install and configure any agent or tool.

## Quick start

Create a sandbox and launch a shell environment:

```console
$ docker sandbox run shell ~/my-project
```

The workspace parameter is optional and defaults to the current directory:

```console
$ cd ~/my-project
$ docker sandbox run shell
```

This launches a bash login shell inside the sandbox.

## Use cases

The Shell sandbox serves several purposes:

- Custom agent installation

  Install agents not officially supported by Docker Sandboxes. The environment
  includes package managers and development tools for installing arbitrary
  software.

- Agent development

  Test custom agent implementations or modifications in an isolated environment
  with a private Docker daemon.

- Manual configuration

  Configure agents with complex setup requirements or custom authentication
  flows that aren't supported by the standard templates.

- Troubleshooting

  Debug agent issues by manually running commands and inspecting the sandbox
  environment.

## Authentication

The Shell sandbox uses proxy credential injection. The proxy automatically
injects credentials into API requests for supported providers (OpenAI,
Anthropic, Google, GitHub, etc.).

Set your API keys in your shell configuration file:

```plaintext {title="~/.bashrc or ~/.zshrc"}
export OPENAI_API_KEY=sk-xxxxx
export ANTHROPIC_API_KEY=sk-ant-xxxxx
export GOOGLE_API_KEY=AIzaSyxxxxx
export GH_TOKEN=ghp_xxxxx
```

Apply the changes:

1. Source your shell configuration: `source ~/.bashrc` (or `~/.zshrc`)
2. Restart Docker Desktop so the daemon picks up the environment variables
3. Create and run your sandbox:

```console
$ docker sandbox create shell ~/project
$ docker sandbox run <sandbox-name>
```

The proxy reads credentials from your host environment and injects them into
API requests automatically. Credentials are never stored inside the sandbox.

## Installing agents

Once inside the shell sandbox, install agents using their standard installation
methods.

### Example: Installing Continue

[Continue](https://continue.dev) is an AI code assistant. Since Node.js is
pre-installed, you can install it directly:

```console
$ npm install -g @continuedev/cli
$ cn --version
1.5.43
```

For containerized agents or complex setups, consider creating a [custom
template](../templates.md) based on the shell template instead of installing
interactively.

## Running commands

Pass shell options after the `--` separator to execute commands:

```console
$ docker sandbox run <sandbox-name> -- -c "echo 'Hello from sandbox'"
```

## Base image

Template: `docker/sandbox-templates:shell`

The shell template provides the base environment without a pre-installed agent,
making it suitable for manual agent installation.

See [Custom templates](../templates.md) to build your own agent images.

The minimal nature of this template makes it suitable as a base for any custom
agent installation.
