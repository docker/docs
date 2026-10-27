---
title: "Fireworks AI"
description: "Use Fireworks AI models with Docker Agent."
keywords: docker agent, ai agents, model providers, llm, fireworks ai
weight: 100
canonical: https://docs.docker.com/ai/docker-agent/providers/fireworks/
---

_Use Fireworks AI models with Docker Agent._

## Overview

[Fireworks AI](https://fireworks.ai/) is a fast inference host for open-weight
models, serving Kimi, Qwen, DeepSeek, GLM and others through an
OpenAI-compatible API. Docker Agent includes built-in support for Fireworks AI
as an alias provider.

Use the provider ID `fireworks-ai`. The legacy ID `fireworks` is also accepted.

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
    model: fireworks-ai/accounts/fireworks/models/kimi-k3
    description: Assistant using Fireworks AI
    instruction: You are a helpful assistant.
```

### Named Model

For more control over parameters:

```yaml
models:
  fireworks_model:
    provider: fireworks-ai
    model: accounts/fireworks/models/kimi-k3
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
| `accounts/fireworks/models/kimi-k3` | Kimi K3, large open MoE chat and tool-calling model |
| `accounts/fireworks/models/kimi-k2p7-code` | Kimi K2.7 Code, coding-focused variant |
| `accounts/fireworks/models/glm-5p3` | GLM 5.3 |
| `accounts/fireworks/models/qwen3p8-max` | Qwen 3.8 Max |
| `accounts/fireworks/models/deepseek-v4-pro-0813` | DeepSeek V4 Pro |
| `accounts/fireworks/models/gpt-oss-120b` | GPT OSS 120B |

> Model IDs are case-sensitive and must be passed exactly as the catalogue lists
> them.

## How It Works

Fireworks AI is implemented as a built-in alias in Docker Agent:

- **API Type:** OpenAI-compatible (`openai_chatcompletions`)
- **Base URL:** `https://api.fireworks.ai/inference/v1`
- **Token Variable:** `FIREWORKS_API_KEY`

Because Fireworks fronts open-weight models whose chat templates may reject more
than one leading system message, Docker Agent coalesces its per-source system
messages into a single one for this provider.

## Example: Code Assistant

```yaml
agents:
  coder:
    model: fireworks-ai/accounts/fireworks/models/kimi-k2p7-code
    description: Code assistant using Kimi K2.7 Code on Fireworks AI
    instruction: |
      You are an expert programmer.
      Write clean, well-documented code and follow language best practices.
    toolsets:
      - type: filesystem
      - type: shell
      - type: think
```
