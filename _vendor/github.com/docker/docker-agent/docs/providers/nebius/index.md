---
title: "Nebius Token Factory"
description: "Use Nebius Token Factory models with Docker Agent."
keywords: docker agent, ai agents, model providers, llm, nebius, token factory
weight: 190
canonical: https://docs.docker.com/ai/docker-agent/providers/nebius/
---

_Use Nebius Token Factory models with Docker Agent._

## Overview

[Nebius Token Factory](https://tokenfactory.nebius.com/) (formerly Nebius AI Studio) provides hosted AI models through an OpenAI-compatible API. Docker Agent includes built-in support for Token Factory as the `nebius` alias provider.

## Setup

1. Create an API key in the [Nebius Token Factory console](https://tokenfactory.nebius.com/project/api-keys)
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
    model: nebius/MiniMaxAI/MiniMax-M3.1-Flash
    description: Assistant using Nebius Token Factory
    instruction: You are a helpful assistant.
```

### Named Model

For more control over parameters:

```yaml
models:
  nebius_model:
    provider: nebius
    model: openai/gpt-oss-120b
    temperature: 0.7
    max_tokens: 8192

agents:
  root:
    model: nebius_model
    description: Assistant using Nebius Token Factory
    instruction: You are a helpful assistant.
```

## Available Models

Token Factory's catalog changes as models launch and retire, so this page does not keep a static model table. With `NEBIUS_API_KEY` set, list the Nebius models Docker Agent knows about:

```bash
docker agent models --provider nebius
```

Add `--all` to list them without a key. Recent additions include `MiniMaxAI/MiniMax-M3.1-Flash`, `deepseek-ai/DeepSeek-V4.1-Flash-1017`, `zai-org/GLM-5.3` and `moonshotai/Kimi-K3`.

The listing comes from the [models.dev](https://models.dev/providers/nebius/) catalog, which can lag behind Token Factory. The [Token Factory console](https://tokenfactory.nebius.com/) is the authoritative list of models available to your project. Model IDs are case-sensitive; copy the full ID, including the vendor prefix.

## How It Works

Nebius is implemented as a built-in alias in Docker Agent:

- **API Type:** OpenAI-compatible (`openai_chatcompletions`)
- **Base URL:** `https://api.tokenfactory.nebius.com/v1`
- **Token Variable:** `NEBIUS_API_KEY`

## Example: Code Assistant

```yaml
agents:
  coder:
    model: nebius/moonshotai/Kimi-K3
    description: Code assistant using Nebius Token Factory
    instruction: |
      You are an expert programmer.
      Write clean, well-documented code.
      Follow best practices for the language being used.
    toolsets:
      - type: filesystem
      - type: shell
      - type: think
```
