---
title: "Cloudflare Workers AI"
description: "Use Cloudflare Workers AI models with docker-agent."
keywords: docker agent, ai agents, model providers, llm, cloudflare workers ai
weight: 70
canonical: https://docs.docker.com/ai/docker-agent/providers/cloudflare-workers-ai/
---

_Use Cloudflare Workers AI models with docker-agent._

## Overview

[Cloudflare Workers AI](https://developers.cloudflare.com/workers-ai/) runs
open-weight models (Llama, Mistral, Qwen, Gemma, and more) on Cloudflare's
global edge network through an OpenAI-compatible endpoint. No separate provider
accounts are needed for the supported models. docker-agent includes built-in
support for Workers AI as an alias provider.

## Setup

Workers AI is account-scoped, so its base URL is resolved from your account ID.
Two environment variables are required:

```bash
export CLOUDFLARE_ACCOUNT_ID=your-account-id
export CLOUDFLARE_API_TOKEN=your-api-token
```

Create an API token with the `Workers AI` permission from the
[Cloudflare dashboard](https://dash.cloudflare.com/profile/api-tokens). Your
account ID is shown on the Workers AI page.

## Usage

Workers AI model IDs use the `@cf/...` form (for example
`@cf/meta/llama-3.1-8b-instruct`).

### Inline Syntax

```yaml
agents:
  root:
    model: cloudflare-workers-ai/@cf/meta/llama-3.1-8b-instruct
    description: Assistant using Cloudflare Workers AI
    instruction: You are a helpful assistant.
```

### Named Model

For more control over parameters:

```yaml
models:
  cloudflare_model:
    provider: cloudflare-workers-ai
    model: "@cf/meta/llama-3.1-8b-instruct"
    temperature: 0.7
    max_tokens: 8192

agents:
  root:
    model: cloudflare_model
    description: Assistant using Cloudflare Workers AI
    instruction: You are a helpful assistant.
```

## Available Models

Check the
[Workers AI models catalog](https://developers.cloudflare.com/workers-ai/models/)
for the current list, IDs, and pricing.

| Model | Description |
| --- | --- |
| `@cf/meta/llama-3.1-8b-instruct` | Meta Llama 3.1 8B Instruct |
| `@cf/mistralai/mistral-small-3.1-24b-instruct` | Mistral Small 3.1 24B Instruct |
| `@cf/qwen/qwen2.5-coder-32b-instruct` | Qwen 2.5 Coder 32B Instruct |

## How It Works

Cloudflare Workers AI is implemented as a built-in alias in docker-agent:

- **API Type:** OpenAI-compatible (`openai_chatcompletions`)
- **Base URL:** `https://api.cloudflare.com/client/v4/accounts/${CLOUDFLARE_ACCOUNT_ID}/ai/v1`
- **Token Variable:** `CLOUDFLARE_API_TOKEN`

The base URL is templated: `${CLOUDFLARE_ACCOUNT_ID}` is substituted from the
environment when the provider is built, so `CLOUDFLARE_ACCOUNT_ID` must be set in
addition to `CLOUDFLARE_API_TOKEN`. Because Workers AI serves open-weight models
with strict chat templates, docker-agent coalesces consecutive system messages
into a single leading one for this provider.
