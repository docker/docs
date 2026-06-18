---
title: Docker Agent
weight: 70
description: |
  Use Docker Agent in Docker Sandboxes with multi-provider authentication
  supporting OpenAI, Anthropic, and more.
keywords: docker sandboxes, docker agent, openai, anthropic, sbx
---

Official documentation: [Docker Agent](/manuals/ai/docker-agent/_index.md)

## Quick start

Create a sandbox and run Docker Agent for a project directory:

```console
$ sbx run docker-agent ~/my-project
```

The workspace parameter defaults to the current directory, so
`sbx run docker-agent` from inside your project works too.

## Authentication

Docker Agent supports multiple providers. Store keys for the providers you want
to use with [stored secrets](../security/credentials.md#stored-secrets):

```console
$ sbx secret set -g openai
$ sbx secret set -g anthropic
$ sbx secret set -g google
$ sbx secret set -g xai
$ sbx secret set -g nebius
$ sbx secret set -g mistral
$ sbx secret set -g openrouter
```

You only need to configure the providers you want to use. Docker Agent detects
available credentials and routes requests to the appropriate provider.

You can also source these from environment variables (`OPENAI_API_KEY`,
`ANTHROPIC_API_KEY`, `GOOGLE_API_KEY`, `XAI_API_KEY`, `NEBIUS_API_KEY`,
`MISTRAL_API_KEY`, `OPENROUTER_API_KEY`) through
[credential bindings](../security/credentials.md#credential-bindings); the
sandbox prompts you to approve one per provider on first run. See
[Credentials](../security/credentials.md) for details.

## Configuration

Sandboxes don't pick up user-level configuration from your host. Only
project-level configuration in the working directory is available inside the
sandbox. See
[Why doesn't the sandbox use my user-level agent configuration?](../faq.md#why-doesnt-the-sandbox-use-my-user-level-agent-configuration)
for workarounds.

### Default startup command

Without extra args, the sandbox runs:

```text
docker-agent run --yolo
```

Arguments after `--` are added after the default flags when the first one is
itself a flag (begins with `-`). When the first argument is a bare word — such
as the `run` subcommand or a config file — it replaces the defaults, so include
`run --yolo` yourself:

```console
$ sbx run docker-agent -- run --yolo agent.yml
```

## Base image

The sandbox uses `docker/sandbox-templates:docker-agent`. See
[Templates](../customize/templates.md) to build your own image on top of
this base.
