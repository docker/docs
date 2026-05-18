---
title: Droid
weight: 35
description: |
  Use Droid in Docker Sandboxes with API key or OAuth authentication.
keywords: docker sandboxes, droid, factory, ai agent, sbx
---

{{< summary-bar feature_name="Docker Sandboxes sbx" >}}

This guide covers authentication, configuration, and usage of Droid, an AI
coding agent by Factory, in a sandboxed environment.

Official documentation: [Droid](https://docs.factory.ai/)

## Quick start

Create a sandbox and run Droid for a project directory:

```console
$ sbx run droid ~/my-project
```

The workspace parameter is optional and defaults to the current directory:

```console
$ cd ~/my-project
$ sbx run droid
```

## Authentication

Droid requires a [Factory account](https://factory.ai). Both authentication
methods authenticate you to Factory's service directly — unlike other agents
where you supply a model provider key, Factory manages model access through
your Factory account.

**API key**: Store your Factory API key using
[stored secrets](../security/credentials.md#stored-secrets):

```console
$ sbx secret set -g droid
```

Alternatively, export the `FACTORY_API_KEY` environment variable in your shell
before running the sandbox. See
[Credentials](../security/credentials.md) for details on both methods.

**OAuth**: If no API key is set, Droid prompts you to authenticate
interactively on first run. The proxy handles the OAuth flow, so credentials
aren't stored inside the sandbox.

## Configuration

Sandboxes don't pick up user-level configuration from your host. Only
project-level configuration in the working directory is available inside the
sandbox. See
[Why doesn't the sandbox use my user-level agent configuration?](../faq.md#why-doesnt-the-sandbox-use-my-user-level-agent-configuration)
for workarounds.

The sandbox runs Droid without approval prompts by default. Pass additional
`droid` CLI options after `--`:

```console
$ sbx run droid --name <sandbox-name> -- <droid-options>
```

## Base image

Template: `docker/sandbox-templates:droid-docker`

Preconfigured to run without approval prompts. Authentication state is
persisted across sandbox restarts.

See [Customize](../customize/) to pre-install tools or customize this
environment.
