---
title: "Cerebras"
description: "Use Cerebras models with docker-agent."
keywords: docker agent, ai agents, model providers, llm, cerebras
weight: 50
canonical: https://docs.docker.com/ai/docker-agent/providers/cerebras/
---

_Use Cerebras models with docker-agent._

## Overview

[Cerebras](https://www.cerebras.ai/) serves open-weight models such as GPT-OSS
and GLM through an OpenAI-compatible API on its wafer-scale hardware, delivering
some of the highest tokens/sec available. That speed makes it a strong fit for
latency-sensitive coding workflows. docker-agent includes built-in support for
Cerebras as an alias provider.

## Setup

1. Create an API key from the [Cerebras Cloud console](https://cloud.cerebras.ai/).
2. Set the environment variable:

   ```bash
   export CEREBRAS_API_KEY=your-api-key
   ```

## Usage

### Inline Syntax

The simplest way to use Cerebras:

```yaml
agents:
  root:
    model: cerebras/gpt-oss-120b
    description: Assistant using Cerebras
    instruction: You are a helpful assistant.
```

### Named Model

For more control over parameters:

```yaml
models:
  cerebras_model:
    provider: cerebras
    model: gpt-oss-120b
    temperature: 0.7
    max_tokens: 8192

agents:
  root:
    model: cerebras_model
    description: Assistant using Cerebras
    instruction: You are a helpful assistant.
```

## Available Models

Cerebras hosts a curated set of open-weight models. Check the
[Cerebras models documentation](https://inference-docs.cerebras.ai/models/overview)
for current model IDs, context limits, and pricing.

| Model | Description |
| --- | --- |
| `gpt-oss-120b` | Open-weight GPT-OSS reasoning model with tool calling |
| `zai-glm-4.7` | Z.AI GLM-4.7 reasoning model with tool calling |

> Model IDs are case-sensitive and must be passed exactly as the catalogue lists
> them. Cerebras may serve additional models not in the built-in catalog; those
> still work but resolve to default capability metadata locally.

## How It Works

Cerebras is implemented as a built-in alias in docker-agent:

- **API Type:** OpenAI-compatible (`openai_chatcompletions`)
- **Base URL:** `https://api.cerebras.ai/v1`
- **Token Variable:** `CEREBRAS_API_KEY`

Because Cerebras fronts open-weight models whose chat templates may only accept
a single leading system message, docker-agent coalesces its per-source system
messages (agent instruction plus each toolset's instructions) into one before
sending the request.

## Example: Code Assistant

```yaml
agents:
  coder:
    model: cerebras/gpt-oss-120b
    description: Fast code assistant using Cerebras
    instruction: |
      You are an expert programmer.
      Write clean, well-documented code and follow language best practices.
    toolsets:
      - type: filesystem
      - type: shell
      - type: think
```
