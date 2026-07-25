---
title: "Baseten"
description: "Use Baseten AI models with docker-agent."
keywords: docker agent, ai agents, model providers, llm, baseten
weight: 40
canonical: https://docs.docker.com/ai/docker-agent/providers/baseten/
---

_Use Baseten AI models with docker-agent._

## Overview

Baseten provides AI models through an OpenAI-compatible API. docker-agent includes built-in support for Baseten as an alias provider.

## Setup

1. Get an API key from [Baseten](https://www.baseten.co/)
2. Set the environment variable:

   ```bash
   export BASETEN_API_KEY=your-api-key
   ```

## Usage

### Inline Syntax

The simplest way to use Baseten:

```yaml
agents:
  root:
    model: baseten/deepseek-ai/DeepSeek-V3.1
    description: Assistant using Baseten
    instruction: You are a helpful assistant.
```

### Named Model

For more control over parameters:

```yaml
models:
  baseten_model:
    provider: baseten
    model: deepseek-ai/DeepSeek-V3.1
    temperature: 0.7
    max_tokens: 8192

agents:
  root:
    model: baseten_model
    description: Assistant using Baseten
    instruction: You are a helpful assistant.
```

## Available Models

Baseten hosts various open models through its Model APIs. Check the [Baseten documentation](https://docs.baseten.co/) for the current model catalog.

| Model                          | Description                    |
| ------------------------------ | ------------------------------ |
| `deepseek-ai/DeepSeek-V3.1`    | DeepSeek V3.1 model            |
| `moonshotai/Kimi-K2.5`         | Moonshot Kimi K2.5 model       |
| `openai/gpt-oss-120b`          | GPT-OSS 120B model             |
| `zai-org/GLM-5`                | GLM-5 model                    |

## How It Works

Baseten is implemented as a built-in alias in docker-agent:

- **API Type:** OpenAI-compatible (`openai_chatcompletions`)
- **Base URL:** `https://inference.baseten.co/v1`
- **Token Variable:** `BASETEN_API_KEY`

## Example: Code Assistant

```yaml
agents:
  coder:
    model: baseten/deepseek-ai/DeepSeek-V3.1
    description: Code assistant using DeepSeek
    instruction: |
      You are an expert programmer using DeepSeek V3.1.
      Write clean, well-documented code.
      Follow best practices for the language being used.
    toolsets:
      - type: filesystem
      - type: shell
      - type: think
```
