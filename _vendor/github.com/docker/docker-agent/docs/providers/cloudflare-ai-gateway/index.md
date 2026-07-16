---
title: "Cloudflare AI Gateway"
description: "Use Cloudflare AI Gateway models with docker-agent."
keywords: docker agent, ai agents, model providers, llm, cloudflare ai gateway
weight: 60
canonical: https://docs.docker.com/ai/docker-agent/providers/cloudflare-ai-gateway/
---

_Use Cloudflare AI Gateway models with docker-agent._

## Overview

[Cloudflare AI Gateway](https://developers.cloudflare.com/ai-gateway/) is a
single OpenAI-compatible endpoint that routes to models from OpenAI, Anthropic,
Workers AI and more, with caching, rate limiting and observability. docker-agent
includes built-in support for AI Gateway as an alias provider.

The alias sends your token in the standard `Authorization: Bearer` header, so it
works out of the box with a gateway that has **authentication disabled** (the
default), typically to route to your own Workers AI models through a gateway you
own. See [Authentication](#authentication) below for the unified-billing /
authenticated-gateway caveat.

## Setup

The gateway endpoint is account- and gateway-scoped, so its base URL is resolved
from your account ID and gateway ID. Three environment variables are required:

```bash
export CLOUDFLARE_ACCOUNT_ID=your-account-id
export CLOUDFLARE_GATEWAY_ID=your-gateway-id
export CLOUDFLARE_API_TOKEN=your-api-token
```

Create a gateway from the
[AI Gateway dashboard](https://dash.cloudflare.com/?to=/:account/ai/ai-gateway)
and an API token with the appropriate permissions.

## Usage

AI Gateway model IDs use the gateway's `provider/model` form (for example
`workers-ai/@cf/meta/llama-3.1-8b-instruct` or `openai/gpt-4o`); the gateway
routes each request to the underlying provider.

### Inline Syntax

```yaml
agents:
  root:
    model: cloudflare-ai-gateway/workers-ai/@cf/meta/llama-3.1-8b-instruct
    description: Assistant using Cloudflare AI Gateway
    instruction: You are a helpful assistant.
```

### Named Model

For more control over parameters:

```yaml
models:
  cloudflare_model:
    provider: cloudflare-ai-gateway
    model: "workers-ai/@cf/meta/llama-3.1-8b-instruct"
    temperature: 0.7
    max_tokens: 8192

agents:
  root:
    model: cloudflare_model
    description: Assistant using Cloudflare AI Gateway
    instruction: You are a helpful assistant.
```

## Available Models

AI Gateway exposes models from many providers behind one endpoint. Check the
[AI Gateway documentation](https://developers.cloudflare.com/ai-gateway/) for
the current provider list, model IDs, and how billing works.

> Model IDs are case-sensitive and must be passed exactly as the gateway lists
> them, including the `provider/` prefix.

## How It Works

Cloudflare AI Gateway is implemented as a built-in alias in docker-agent:

- **API Type:** OpenAI-compatible (`openai_chatcompletions`)
- **Base URL:** `https://gateway.ai.cloudflare.com/v1/${CLOUDFLARE_ACCOUNT_ID}/${CLOUDFLARE_GATEWAY_ID}/compat`
- **Token Variable:** `CLOUDFLARE_API_TOKEN`

The base URL is templated: `${CLOUDFLARE_ACCOUNT_ID}` and
`${CLOUDFLARE_GATEWAY_ID}` are substituted from the environment when the provider
is built, so both must be set in addition to `CLOUDFLARE_API_TOKEN`. Because the
gateway can route to open-weight models with strict chat templates, docker-agent
coalesces consecutive system messages into a single leading one for this
provider.

## Authentication

docker-agent authenticates by sending `CLOUDFLARE_API_TOKEN` in the standard
`Authorization: Bearer` header. On the `.../compat` endpoint that header is
treated as the **provider** key, so this alias works when:

- the gateway has **authentication disabled** (the default), and
- the routed models accept that token as their provider key, which is the case
  for **Workers AI** models (`workers-ai/@cf/...`).

A gateway with **authentication enabled** (required for
[unified billing](https://developers.cloudflare.com/ai-gateway/features/unified-billing/))
instead expects the token in Cloudflare's `cf-aig-authorization` header. The
alias does not send that header, and custom `provider_opts.http_headers` values
are not environment-expanded, so an authenticated gateway is **not supported out
of the box** today. For that setup, use an unauthenticated gateway, or configure
a [custom provider](../custom/index.md) against the
Cloudflare AI Gateway REST API.
