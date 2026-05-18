---
title: Docker Agent
weight: 70
description: |
  Use Docker Agent in Docker Sandboxes with multi-provider authentication
  supporting OpenAI, Anthropic, and more.
keywords: docker sandboxes, docker agent, openai, anthropic, sbx
---

{{< summary-bar feature_name="Docker Sandboxes sbx" >}}

Official documentation: [Docker Agent](https://docs.docker.com/ai/docker-agent/)

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
```

You only need to configure the providers you want to use. Docker Agent detects
available credentials and routes requests to the appropriate provider.

Alternatively, export the environment variables (`OPENAI_API_KEY`,
`ANTHROPIC_API_KEY`, `GOOGLE_API_KEY`, `XAI_API_KEY`, `NEBIUS_API_KEY`,
`MISTRAL_API_KEY`) in your shell before running the sandbox. See
[Credentials](../security/credentials.md) for details on both methods.

## Configuration

Sandboxes don't pick up user-level configuration from your host. Only
project-level configuration in the working directory is available inside the
sandbox. See
[Why doesn't the sandbox use my user-level agent configuration?](../faq.md#why-doesnt-the-sandbox-use-my-user-level-agent-configuration)
for workarounds.

The sandbox runs Docker Agent without approval prompts by default. Pass
additional CLI options after `--`:

```console
$ sbx run docker-agent --name my-sandbox -- <options>
```

For example, to specify a custom `agent.yml` configuration file:

```console
$ sbx run docker-agent -- agent.yml
```

## Base image

The sandbox uses `docker/sandbox-templates:docker-agent` and launches Docker
Agent without approval prompts by default. See
[Templates](../customize/templates.md) to build your own image on top of
this base.
