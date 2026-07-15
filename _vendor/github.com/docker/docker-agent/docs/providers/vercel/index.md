---
title: "Vercel AI Gateway"
description: "Use Vercel AI Gateway models with docker-agent."
keywords: docker agent, ai agents, model providers, llm, vercel ai gateway
weight: 260
canonical: https://docs.docker.com/ai/docker-agent/providers/vercel/
---

_Use Vercel AI Gateway models with docker-agent._

## Overview

[Vercel AI Gateway](https://vercel.com/docs/ai-gateway) is a single, unified
OpenAI-compatible endpoint that routes to models from OpenAI, Anthropic, Google,
xAI and more at list price with no markup, plus provider routing and failover.
It lets you reach many providers with one API key. docker-agent includes
built-in support for Vercel AI Gateway as an alias provider.

## Setup

1. Create an API key from the [Vercel AI Gateway dashboard](https://vercel.com/docs/ai-gateway).
2. Set the environment variable:

   ```bash
   export AI_GATEWAY_API_KEY=your-api-key
   ```

## Usage

Vercel AI Gateway model IDs use a `creator/model` form (for example
`openai/gpt-5` or `anthropic/claude-sonnet-4.5`); the gateway routes each
request to the underlying provider.

### Inline Syntax

The simplest way to use Vercel AI Gateway:

```yaml
agents:
  root:
    model: vercel/openai/gpt-5
    description: Assistant using Vercel AI Gateway
    instruction: You are a helpful assistant.
```

### Named Model

For more control over parameters:

```yaml
models:
  vercel_model:
    provider: vercel
    model: openai/gpt-5
    temperature: 0.7
    max_tokens: 8192

agents:
  root:
    model: vercel_model
    description: Assistant using Vercel AI Gateway
    instruction: You are a helpful assistant.
```

## Available Models

Vercel AI Gateway exposes models from many providers behind one endpoint. Check
the [Vercel AI Gateway documentation](https://vercel.com/docs/ai-gateway) for
the current model list, IDs, and pricing.

| Model | Description |
| --- | --- |
| `openai/gpt-5` | OpenAI GPT-5 routed through the gateway |
| `anthropic/claude-sonnet-4.5` | Anthropic Claude Sonnet routed through the gateway |
| `google/gemini-2.5-flash` | Google Gemini routed through the gateway |

> Model IDs are case-sensitive and must be passed exactly as the gateway lists
> them, including the `creator/` prefix.

## How It Works

Vercel AI Gateway is implemented as a built-in alias in docker-agent:

- **API Type:** OpenAI-compatible (`openai_chatcompletions`)
- **Base URL:** `https://ai-gateway.vercel.sh/v1`
- **Token Variable:** `AI_GATEWAY_API_KEY`

Because the gateway can route to open-weight models with strict chat templates,
docker-agent coalesces consecutive system messages into a single leading one for
this provider.

## Example: Code Assistant

```yaml
agents:
  coder:
    model: vercel/anthropic/claude-sonnet-4.5
    description: Code assistant via Vercel AI Gateway
    instruction: |
      You are an expert programmer.
      Write clean, well-documented code and follow language best practices.
    toolsets:
      - type: filesystem
      - type: shell
      - type: think
```
