---
title: "DeepSeek"
description: "Use DeepSeek models with docker-agent."
keywords: docker agent, ai agents, model providers, llm, deepseek
weight: 80
canonical: https://docs.docker.com/ai/docker-agent/providers/deepseek/
---

_Use DeepSeek models with docker-agent._

## Overview

[DeepSeek](https://www.deepseek.com/) serves its frontier chat and reasoning
models through an OpenAI-compatible API, with strong price/performance on coding
and reasoning tasks. docker-agent includes built-in support for DeepSeek as an
alias provider.

## Setup

1. Create an API key from the [DeepSeek Platform](https://platform.deepseek.com/api_keys).
2. Set the environment variable:

   ```bash
   export DEEPSEEK_API_KEY=your-api-key
   ```

## Usage

### Inline Syntax

The simplest way to use DeepSeek:

```yaml
agents:
  root:
    model: deepseek/deepseek-chat
    description: Assistant using DeepSeek
    instruction: You are a helpful assistant.
```

### Named Model

For more control over parameters:

```yaml
models:
  deepseek_model:
    provider: deepseek
    model: deepseek-chat
    temperature: 0.7
    max_tokens: 8192

agents:
  root:
    model: deepseek_model
    description: Assistant using DeepSeek
    instruction: You are a helpful assistant.
```

## Available Models

DeepSeek exposes a small, vendor-controlled model lineup. Check the
[DeepSeek models documentation](https://api-docs.deepseek.com/quick_start/pricing)
for current model IDs, context limits, and pricing.

| Model | Description |
| --- | --- |
| `deepseek-chat` | DeepSeek-V3, general-purpose chat and tool calling |
| `deepseek-reasoner` | DeepSeek-R1, extended-reasoning model |

> Model IDs are case-sensitive and must be passed exactly as the catalogue lists
> them.

## How It Works

DeepSeek is implemented as a built-in alias in docker-agent:

- **API Type:** OpenAI-compatible (`openai_chatcompletions`)
- **Base URL:** `https://api.deepseek.com/v1`
- **Token Variable:** `DEEPSEEK_API_KEY`

## Example: Code Assistant

```yaml
agents:
  coder:
    model: deepseek/deepseek-chat
    description: Code assistant using DeepSeek-V3
    instruction: |
      You are an expert programmer.
      Write clean, well-documented code and follow language best practices.
    toolsets:
      - type: filesystem
      - type: shell
      - type: think
```
