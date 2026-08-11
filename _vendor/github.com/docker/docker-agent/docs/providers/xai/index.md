---
title: "xAI (Grok)"
description: "Use xAI's Grok models with docker-agent."
keywords: docker agent, ai agents, model providers, llm, xai (grok)
weight: 270
canonical: https://docs.docker.com/ai/docker-agent/providers/xai/
---

_Use xAI's Grok models with docker-agent._

## Overview

xAI provides the Grok family of models through an OpenAI-compatible API. docker-agent includes built-in support for xAI as an alias provider.

## Setup

1. Get an API key from [xAI Console](https://console.x.ai/)
2. Set the environment variable:

   ```bash
   export XAI_API_KEY=your-api-key
   ```

## Usage

### Inline Syntax

The simplest way to use xAI:

```yaml
agents:
  root:
    model: xai/grok-3
    description: Assistant using Grok
    instruction: You are a helpful assistant.
```

### Named Model

For more control over parameters:

```yaml
models:
  grok:
    provider: xai
    model: grok-3
    temperature: 0.7
    max_tokens: 8192

agents:
  root:
    model: grok
    description: Assistant using Grok
    instruction: You are a helpful assistant.
```

## Available Models

| Model              | Description                        | Context |
| ------------------ | ---------------------------------- | ------- |
| `grok-3`           | Latest and most capable Grok model | 131K    |
| `grok-3-fast`      | Faster variant with lower latency  | 131K    |
| `grok-3-mini`      | Compact model for simpler tasks    | 131K    |
| `grok-3-mini-fast` | Fast variant of the mini model     | 131K    |
| `grok-2`           | Previous generation model          | 128K    |
| `grok-vision`      | Vision-capable model               | 32K     |

Check the [xAI documentation](https://docs.x.ai/docs) for the latest available models.

## Extended Thinking

docker-agent's `thinking_budget` field is **not applied** to xAI models: the underlying OpenAI-compatible client only sends `reasoning_effort` for OpenAI reasoning model names (o-series, gpt-5). Setting `thinking_budget` on a Grok model passes config validation but has no effect on the request.

Grok reasoning models (e.g. `grok-3-mini`) reason on their own without configuration. For non-reasoning models, use the [think tool](../../tools/think/index.md) instead.

## How It Works

xAI is implemented as a built-in alias in docker-agent:

- **API Type:** OpenAI-compatible (`openai_chatcompletions`)
- **Base URL:** `https://api.x.ai/v1`
- **Token Variable:** `XAI_API_KEY`

## Example: Research Assistant

```yaml
agents:
  researcher:
    model: xai/grok-3
    description: Research assistant with real-time knowledge
    instruction: |
      You are a research assistant using Grok.
      Provide well-researched, factual responses.
      Cite sources when available.
    toolsets:
      - type: mcp
        ref: docker:duckduckgo
      - type: think
```
