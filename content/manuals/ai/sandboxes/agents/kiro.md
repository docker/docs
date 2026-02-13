---
title: Kiro sandbox
description: |
  Use Kiro in Docker Sandboxes with device flow authentication for interactive
  AI-assisted development.
keywords: docker, sandboxes, kiro, ai agent, authentication, device flow
weight: 70
---

{{< summary-bar feature_name="Docker Sandboxes" >}}

This guide covers authentication, configuration, and usage of Kiro in a
sandboxed environment.

Official documentation: [Kiro CLI](https://kiro.dev/docs/cli/)

## Quick start

Create a sandbox and run Kiro for a project directory:

```console
$ docker sandbox run kiro ~/my-project
```

The workspace parameter is optional and defaults to the current directory:

```console
$ cd ~/my-project
$ docker sandbox run kiro
```

On first run, Kiro prompts you to authenticate using device flow.

## Authentication

Kiro uses device flow authentication, which requires interactive login through
a web browser. This method provides secure authentication without storing API
keys directly.

### Device flow login

When you first run Kiro, it prompts you to authenticate:

1. Kiro displays a URL and a verification code
2. Open the URL in your web browser
3. Enter the verification code
4. Complete the authentication flow in your browser
5. Return to the terminal - Kiro proceeds automatically

The authentication session is persisted in the sandbox and doesn't require
repeated login unless you destroy and recreate the sandbox.

### Manual login

You can trigger the login flow manually:

```console
$ docker sandbox run <sandbox-name> -- login --use-device-flow
```

This command initiates device flow authentication without starting a coding
session.

### Authentication persistence

Kiro stores authentication state in `~/.local/share/kiro-cli/data.sqlite3`
inside the sandbox. This database persists as long as the sandbox exists. If
you destroy the sandbox, you'll need to authenticate again when you recreate
it.

## Configuration

Kiro requires minimal configuration. The agent runs with trust-all-tools mode
by default, which allows it to execute commands without repeated approval
prompts.

### Pass options at runtime

Pass Kiro CLI options after the sandbox name and a `--` separator:

```console
$ docker sandbox run <sandbox-name> -- chat --trust-all-tools
```

The `chat --trust-all-tools` command starts Kiro with approval prompts
disabled.

## Base image

Template: `docker/sandbox-templates:kiro`

Kiro manages authentication through an interactive device flow. The authentication database is persisted across sandbox restarts. Launches with `chat --trust-all-tools` by default.

See [Custom templates](../templates.md) to build your own agent images.
