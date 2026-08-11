---
title: "Together AI"
description: "Use Together AI models with docker-agent."
keywords: docker agent, ai agents, model providers, llm, together ai
weight: 250
canonical: https://docs.docker.com/ai/docker-agent/providers/together/
---

_Use Together AI models with docker-agent._

## Overview

[Together AI](https://www.together.ai/) is one of the largest hosts of open
models, serving Llama, Qwen, DeepSeek, Kimi, GLM and others through an
OpenAI-compatible API. docker-agent includes built-in support for Together AI as
an alias provider.

## Setup

1. Create an API key from the [Together AI settings](https://api.together.ai/settings/api-keys).
2. Set the environment variable:

   ```bash
   export TOGETHER_API_KEY=your-api-key
   ```

## Usage

### Inline Syntax

The simplest way to use Together AI:

```yaml
agents:
  root:
    model: together/meta-llama/Llama-3.3-70B-Instruct-Turbo
    description: Assistant using Together AI
    instruction: You are a helpful assistant.
```

### Named Model

For more control over parameters:

```yaml
models:
  together_model:
    provider: together
    model: meta-llama/Llama-3.3-70B-Instruct-Turbo
    temperature: 0.7
    max_tokens: 8192

agents:
  root:
    model: together_model
    description: Assistant using Together AI
    instruction: You are a helpful assistant.
```

## Available Models

Together AI serves a broad, changing catalog of open-weight models. Check the
[Together AI model library](https://docs.together.ai/docs/serverless-models) for
current model IDs, context limits, and pricing.

| Model | Description |
| --- | --- |
| `meta-llama/Llama-3.3-70B-Instruct-Turbo` | Llama 3.3 70B, general-purpose chat and tool calling |
| `Qwen/Qwen3-235B-A22B-Instruct-2507-tput` | Qwen3 235B mixture-of-experts instruct model |
| `deepseek-ai/DeepSeek-V3` | DeepSeek-V3, strong coding and reasoning |

> Model IDs are case-sensitive and must be passed exactly as the catalogue lists
> them.

## How It Works

Together AI is implemented as a built-in alias in docker-agent:

- **API Type:** OpenAI-compatible (`openai_chatcompletions`)
- **Base URL:** `https://api.together.xyz/v1`
- **Token Variable:** `TOGETHER_API_KEY`

Because Together AI fronts open-weight models whose chat templates may reject
more than one leading system message, docker-agent coalesces its per-source
system messages into a single one for this provider.

## Example: Code Assistant

```yaml
agents:
  coder:
    model: together/Qwen/Qwen3-235B-A22B-Instruct-2507-tput
    description: Code assistant using Qwen3 on Together AI
    instruction: |
      You are an expert programmer.
      Write clean, well-documented code and follow language best practices.
    toolsets:
      - type: filesystem
      - type: shell
      - type: think
```
