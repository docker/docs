---
title: Droid
weight: 60
description: |
  Use Droid in Docker Sandboxes with API key or OAuth authentication.
keywords: docker sandboxes, droid, factory, ai agent, sbx
---

This guide covers authentication, configuration, and usage of Droid, an AI
coding agent by Factory, in a sandboxed environment.

Official documentation: [Droid](https://docs.factory.ai/)

## Quick start

Create a sandbox and run Droid for a project directory:

```console
$ sbx run droid ~/my-project
```

Use `.` to mount the current directory:

```console
$ cd ~/my-project
$ sbx run droid .
```

Omit the workspace path to create a
[mountless sandbox](../usage.md#choose-a-workspace) instead.

## Authentication

Droid requires a [Factory account](https://factory.ai). Both authentication
methods authenticate you to Factory's service directly — unlike other agents
where you supply a model provider key, Factory manages model access through
your Factory account.

**API key**: Store your Factory API key using
[stored secrets](../configuration/credentials.md#stored-secrets):

```console
$ sbx secret set droid
```

**OAuth**: If no API key is set, Droid prompts you to authenticate
interactively on first run. The proxy handles the OAuth flow, so credentials
aren't stored inside the sandbox.

## Configuration

Sandboxes don't pick up user-level configuration from your host. Only
project-level configuration in the working directory is available inside the
sandbox. See
[Why doesn't the sandbox use my user-level agent configuration?](../faq.md#why-doesnt-the-sandbox-use-my-user-level-agent-configuration)
for workarounds.

### Default startup command

The sandbox runs `droid` with no implicit flags. Args after `--` are passed
straight through:

```console
$ sbx run --name <sandbox-name> -- exec "fix the build"
```

## Base image

Template: `docker/sandbox-templates:droid-docker`

Preconfigured to run without approval prompts. Authentication state is
persisted across sandbox restarts.

See [Customize](../customize/) to pre-install tools or customize this
environment.
