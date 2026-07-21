---
title: "GitHub Copilot"
description: "Use GitHub Copilot's hosted models (GPT-4o, Claude, Gemini, and more) with docker-agent through your GitHub subscription."
keywords: docker agent, ai agents, model providers, llm, github copilot
weight: 110
canonical: https://docs.docker.com/ai/docker-agent/providers/github-copilot/
---

_Use GitHub Copilot's hosted models with docker-agent through your existing GitHub subscription._

## Overview

GitHub Copilot exposes an OpenAI-compatible Chat Completions API at
`https://api.githubcopilot.com`. docker-agent ships with built-in support for
it as the `github-copilot` provider, so any user with a paid GitHub Copilot
subscription can reuse their entitlement from docker-agent.

## Prerequisites

- An active **GitHub Copilot** subscription (Individual, Business, or Enterprise).
- A **personal access token** with the `copilot` scope, exported as `GITHUB_TOKEN`.

```bash
export GITHUB_TOKEN="ghp_..."
```

## Configuration

### Inline

```yaml
agents:
  root:
    model: github-copilot/gpt-4o
    instruction: You are a helpful assistant.
```

### Named model

```yaml
models:
  copilot:
    provider: github-copilot
    model: gpt-4o
    temperature: 0.7
    max_tokens: 4000

agents:
  root:
    model: copilot
```

## Available Models

The exact set of models you can call depends on your Copilot plan. The most
commonly available ones today are:

| Model                    | Best For                            |
| ------------------------ | ----------------------------------- |
| `gpt-4o`                 | Multimodal, balanced performance    |
| `gpt-4o-mini`            | Fast and cheap                      |
| `claude-sonnet-4`        | Strong coding and analysis          |
| `gemini-2.5-pro`         | Google's flagship, large context    |
| `o3-mini`                | Reasoning-focused                   |

Check the
[GitHub Copilot documentation](https://docs.github.com/en/copilot)
for the current model list.

## `Copilot-Integration-Id` Header

GitHub's Copilot API rejects requests that don't carry a
`Copilot-Integration-Id` header with a `Bad Request` error. docker-agent
automatically sends `copilot-developer-cli` for the `github-copilot`
provider, so PAT-based usage works out of the box.

We specifically chose `copilot-developer-cli` (instead of, say,
`vscode-chat`) because it is the integration id accepted by the Copilot
API for **both** OAuth tokens and Personal Access Tokens. Most
docker-agent users authenticate with a PAT exported as `GITHUB_TOKEN`,
and `vscode-chat` is rejected for those tokens.

If you need to send a different integration id — for example if your
organization allows-lists a specific value — you can override it via
`provider_opts.http_headers`:

```yaml
models:
  copilot:
    provider: github-copilot
    model: gpt-4o
    provider_opts:
      http_headers:
        Copilot-Integration-Id: my-custom-integration
```

Header names are matched case-insensitively, so `copilot-integration-id`
works too.

## Chat Completions vs. Responses API

GitHub Copilot proxies OpenAI models behind two endpoints: the legacy
`/chat/completions` and the newer `/responses`. Newer models (the `gpt-5`
family, Codex variants, etc.) are only served via `/responses` and reject
`/chat/completions` with a `400 Bad Request`. docker-agent auto-selects the
right endpoint per model, so no configuration is needed in the common case.

If you ever need to force one or the other, set `api_type` explicitly:

```yaml
models:
  copilot:
    provider: github-copilot
    model: gpt-5
    provider_opts:
      api_type: openai_responses # or openai_chatcompletions
```

## Custom HTTP Headers

`provider_opts.http_headers` is a generic escape hatch that works for any
OpenAI-compatible provider, not just GitHub Copilot. Every key/value pair
is added to every outgoing request:

```yaml
models:
  my_model:
    provider: openai
    model: gpt-4o
    provider_opts:
      http_headers:
        X-Request-Source: docker-agent
        X-Tenant-Id: my-team
```

## How It Works

GitHub Copilot is implemented as a built-in alias in docker-agent:

- **API type:** OpenAI-compatible (Chat Completions)
- **Base URL:** `https://api.githubcopilot.com`
- **Token variable:** `GITHUB_TOKEN`
- **Default headers:** `Copilot-Integration-Id: copilot-developer-cli`

This means the same client as OpenAI is used, so every OpenAI feature
supported by docker-agent (tool calling, structured output, multimodal
inputs, etc.) is available when the underlying model supports it.
