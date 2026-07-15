---
title: "Fireworks AI"
description: "Use Fireworks AI models with docker-agent."
keywords: docker agent, ai agents, model providers, llm, fireworks ai
weight: 100
canonical: https://docs.docker.com/ai/docker-agent/providers/fireworks/
---

_Use Fireworks AI models with docker-agent._

## Overview

[Fireworks AI](https://fireworks.ai/) is a fast inference host for open-weight
models, serving Kimi K2, Llama, Qwen, DeepSeek, GLM and others through an
OpenAI-compatible API. docker-agent includes built-in support for Fireworks AI
as an alias provider.

## Setup

1. Create an API key from the [Fireworks dashboard](https://fireworks.ai/account/api-keys).
2. Set the environment variable:

   ```bash
   export FIREWORKS_API_KEY=your-api-key
   ```

## Usage

### Inline Syntax

The simplest way to use Fireworks AI:

```yaml
agents:
  root:
    model: fireworks/accounts/fireworks/models/kimi-k2-instruct
    description: Assistant using Fireworks AI
    instruction: You are a helpful assistant.
```

### Named Model

For more control over parameters:

```yaml
models:
  fireworks_model:
    provider: fireworks
    model: accounts/fireworks/models/kimi-k2-instruct
    temperature: 0.7
    max_tokens: 8192

agents:
  root:
    model: fireworks_model
    description: Assistant using Fireworks AI
    instruction: You are a helpful assistant.
```

## Available Models

Fireworks serves a broad, changing catalog of open-weight models. Model IDs use
the `accounts/fireworks/models/<name>` form. Check the
[Fireworks model library](https://fireworks.ai/models) for current IDs, context
limits, and pricing.

| Model | Description |
| --- | --- |
| `accounts/fireworks/models/kimi-k2-instruct` | Kimi K2, large open MoE chat and tool-calling model |
| `accounts/fireworks/models/llama-v3p3-70b-instruct` | Llama 3.3 70B instruct |
| `accounts/fireworks/models/qwen3-235b-a22b` | Qwen 3 235B MoE |

> Model IDs are case-sensitive and must be passed exactly as the catalogue lists
> them.

## How It Works

Fireworks AI is implemented as a built-in alias in docker-agent:

- **API Type:** OpenAI-compatible (`openai_chatcompletions`)
- **Base URL:** `https://api.fireworks.ai/inference/v1`
- **Token Variable:** `FIREWORKS_API_KEY`

Because Fireworks fronts open-weight models whose chat templates may reject more
than one leading system message, docker-agent coalesces its per-source system
messages into a single one for this provider.

## Example: Code Assistant

```yaml
agents:
  coder:
    model: fireworks/accounts/fireworks/models/kimi-k2-instruct
    description: Code assistant using Kimi K2 on Fireworks AI
    instruction: |
      You are an expert programmer.
      Write clean, well-documented code and follow language best practices.
    toolsets:
      - type: filesystem
      - type: shell
      - type: think
```
