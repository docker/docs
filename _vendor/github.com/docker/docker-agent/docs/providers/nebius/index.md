---
title: "Nebius"
description: "Use Nebius AI models with docker-agent."
keywords: docker agent, ai agents, model providers, llm, nebius
weight: 190
canonical: https://docs.docker.com/ai/docker-agent/providers/nebius/
---

_Use Nebius AI models with docker-agent._

## Overview

Nebius provides AI models through an OpenAI-compatible API. docker-agent includes built-in support for Nebius as an alias provider.

## Setup

1. Get an API key from [Nebius AI](https://nebius.ai/)
2. Set the environment variable:

   ```bash
   export NEBIUS_API_KEY=your-api-key
   ```

## Usage

### Inline Syntax

The simplest way to use Nebius:

```yaml
agents:
  root:
    model: nebius/deepseek-ai/DeepSeek-V3
    description: Assistant using Nebius
    instruction: You are a helpful assistant.
```

### Named Model

For more control over parameters:

```yaml
models:
  nebius_model:
    provider: nebius
    model: deepseek-ai/DeepSeek-V3
    temperature: 0.7
    max_tokens: 8192

agents:
  root:
    model: nebius_model
    description: Assistant using Nebius
    instruction: You are a helpful assistant.
```

## Available Models

Nebius hosts various open models. Check the [Nebius documentation](https://nebius.ai/docs) for the current model catalog.

| Model                               | Description                     |
| ----------------------------------- | ------------------------------- |
| `deepseek-ai/DeepSeek-V3`           | DeepSeek V3 model               |
| `Qwen/Qwen2.5-72B-Instruct`         | Qwen 2.5 72B instruction-tuned  |
| `meta-llama/Llama-3.3-70B-Instruct` | Llama 3.3 70B instruction-tuned |

## How It Works

Nebius is implemented as a built-in alias in docker-agent:

- **API Type:** OpenAI-compatible (`openai_chatcompletions`)
- **Base URL:** `https://api.studio.nebius.com/v1`
- **Token Variable:** `NEBIUS_API_KEY`

## Example: Code Assistant

```yaml
agents:
  coder:
    model: nebius/deepseek-ai/DeepSeek-V3
    description: Code assistant using DeepSeek
    instruction: |
      You are an expert programmer using DeepSeek V3.
      Write clean, well-documented code.
      Follow best practices for the language being used.
    toolsets:
      - type: filesystem
      - type: shell
      - type: think
```
