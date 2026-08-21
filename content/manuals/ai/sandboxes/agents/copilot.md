---
title: Copilot
weight: 30
description: |
  Use GitHub Copilot in Docker Sandboxes with GitHub token authentication and
  trusted folder configuration.
keywords: docker sandboxes, github copilot, ai agent, github token, sbx
---

This guide covers authentication, configuration, and usage of GitHub Copilot
in a sandboxed environment.

Official documentation: [GitHub Copilot CLI](https://docs.github.com/en/copilot/how-tos/copilot-cli)

## Quick start

Create a sandbox and run Copilot for a project directory:

```console
$ sbx run copilot ~/my-project
```

`sbx run` defaults the workspace to the current directory:

```console
$ cd ~/my-project
$ sbx run copilot
```

To create the sandbox without mounting a host workspace, omit the path from
[`sbx create`](../usage.md#choose-a-workspace), then attach by name.

## Authentication

Copilot requires a GitHub token with Copilot access. Store your token using
[stored secrets](../configuration/credentials.md#stored-secrets):

```console
$ sbx secret set github --command 'gh auth token'
```

## Configuration

Sandboxes don't pick up user-level configuration from your host. Only
project-level configuration in the working directory is available inside the
sandbox. See
[Why doesn't the sandbox use my user-level agent configuration?](../faq.md#why-doesnt-the-sandbox-use-my-user-level-agent-configuration)
for workarounds.

Copilot is configured to trust the workspace directory by default, so it
operates without repeated confirmations for workspace files.

### Default startup command

Without extra args, the sandbox runs:

```text
copilot --yolo
```

Arguments after `--` are added after the default flags when the first one is
itself a flag (begins with `-`), so `--yolo` is preserved:

```console
$ sbx run --name <sandbox-name> -- -p "review this PR"   # runs copilot --yolo -p "review this PR"
```

When the first argument is a bare word — a subcommand or prompt — it replaces
the defaults instead.

## Base image

Template: `docker/sandbox-templates:copilot`

Preconfigured to trust the workspace directory.

See [Customize](../customize/) to pre-install tools or customize this
environment.
