---
title: Copilot
weight: 30
description: |
  Use GitHub Copilot in Docker Sandboxes with GitHub token authentication and
  trusted folder configuration.
keywords: docker sandboxes, github copilot, ai agent, github token, sbx
---

{{< summary-bar feature_name="Docker Sandboxes sbx" >}}

This guide covers authentication, configuration, and usage of GitHub Copilot
in a sandboxed environment.

Official documentation: [GitHub Copilot CLI](https://docs.github.com/en/copilot/how-tos/copilot-cli)

## Quick start

Create a sandbox and run Copilot for a project directory:

```console
$ sbx run copilot ~/my-project
```

The workspace parameter is optional and defaults to the current directory:

```console
$ cd ~/my-project
$ sbx run copilot
```

## Authentication

Copilot requires a GitHub token with Copilot access. Store your token using
[stored secrets](../security/credentials.md#stored-secrets):

```console
$ echo "$(gh auth token)" | sbx secret set -g github
```

Alternatively, export the `GH_TOKEN` or `GITHUB_TOKEN` environment variable in
your shell before running the sandbox. See
[Credentials](../security/credentials.md) for details on both methods.

## Configuration

Sandboxes don't pick up user-level configuration from your host. Only
project-level configuration in the working directory is available inside the
sandbox. See
[Why doesn't the sandbox use my user-level agent configuration?](../faq.md#why-doesnt-the-sandbox-use-my-user-level-agent-configuration)
for workarounds.

Copilot is configured to trust the workspace directory by default, so it
operates without repeated confirmations for workspace files.

### Pass options at runtime

Pass Copilot CLI options after `--`:

```console
$ sbx run copilot --name <sandbox-name> -- <copilot-options>
```

## Base image

Template: `docker/sandbox-templates:copilot`

Preconfigured to trust the workspace directory and run without approval prompts.

See [Customize](../customize/) to pre-install tools or customize this
environment.
