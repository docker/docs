---
title: Copilot sandbox
description: |
  Use GitHub Copilot in Docker Sandboxes with GitHub token authentication and
  trusted folder configuration.
keywords: docker, sandboxes, copilot, github, ai agent, authentication, trusted folders
weight: 30
---

{{< summary-bar feature_name="Docker Sandboxes" >}}

This guide covers authentication, configuration, and usage of GitHub Copilot
in a sandboxed environment.

Official documentation: [GitHub Copilot CLI](https://docs.github.com/en/copilot/how-tos/copilot-cli)

## Quick start

Create a sandbox and run Copilot for a project directory:

```console
$ docker sandbox run copilot ~/my-project
```

The workspace parameter is optional and defaults to the current directory:

```console
$ cd ~/my-project
$ docker sandbox run copilot
```

## Authentication

Copilot requires a GitHub token with Copilot access. Credentials are scoped
per sandbox and must be provided through environment variables on the host.

### Environment variable (recommended)

Set the `GH_TOKEN` or `GITHUB_TOKEN` environment variable in your shell
configuration file.

Docker Sandboxes use a daemon process that doesn't inherit environment
variables from your current shell session. To make your token available to
sandboxes, set it globally in your shell configuration file.

Add the token to your shell configuration file:

```plaintext {title="~/.bashrc or ~/.zshrc"}
export GH_TOKEN=ghp_xxxxx
```

Or use `GITHUB_TOKEN`:

```plaintext {title="~/.bashrc or ~/.zshrc"}
export GITHUB_TOKEN=ghp_xxxxx
```

Apply the changes:

1. Source your shell configuration: `source ~/.bashrc` (or `~/.zshrc`)
2. Restart Docker Desktop so the daemon picks up the new environment variable
3. Create and run your sandbox:

```console
$ docker sandbox create copilot ~/project
$ docker sandbox run <sandbox-name>
```

The sandbox detects the environment variable and uses it automatically.

## Configuration

Copilot can be configured to trust specific folders, disabling safety prompts
for those locations. Configure trusted folders in `~/.copilot/config.json`:

```json
{
  "trusted_folders": ["/workspace", "/home/agent/projects"]
}
```

Workspaces are mounted at `/workspace` by default, so trusting this path
allows Copilot to operate without repeated confirmations.

### Pass options at runtime

Pass Copilot CLI options after the sandbox name and a `--` separator:

```console
$ docker sandbox run <sandbox-name> -- --yolo
```

The `--yolo` flag disables approval prompts for a single session without
modifying the configuration file.

## Base image

Template: `docker/sandbox-templates:copilot`

Copilot launches with `--yolo` by default when trusted folders are configured.

See [Custom templates](../templates.md) to build your own agent images.
