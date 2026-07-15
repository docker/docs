---
title: "MiniMax"
description: "Use MiniMax AI models with docker-agent."
keywords: docker agent, ai agents, model providers, llm, minimax
weight: 160
canonical: https://docs.docker.com/ai/docker-agent/providers/minimax/
---

_Use MiniMax AI models with docker-agent._

## Overview

MiniMax provides AI models through an OpenAI-compatible API. docker-agent includes built-in support for MiniMax as an alias provider.

## Setup

1. Get an API key from [MiniMax](https://www.minimaxi.com/)
2. Set the environment variable:

   ```bash
   export MINIMAX_API_KEY=your-api-key
   ```

## Usage

### Inline Syntax

The simplest way to use MiniMax:

```yaml
agents:
  root:
    model: minimax/MiniMax-M2.5
    description: Assistant using MiniMax
    instruction: You are a helpful assistant.
```

### Named Model

For more control over parameters:

```yaml
models:
  minimax_model:
    provider: minimax
    model: MiniMax-M2.5
    temperature: 0.7
    max_tokens: 8192

agents:
  root:
    model: minimax_model
    description: Assistant using MiniMax
    instruction: You are a helpful assistant.
```

## Available Models

Check the [MiniMax documentation](https://www.minimaxi.com/document/introduction) for the current model catalog.

| Model                    | Description                                     |
| ------------------------ | ----------------------------------------------- |
| `MiniMax-M2.5`           | Peak performance, 204K context                  |
| `MiniMax-M2.5-highspeed` | Same as M2.5 but faster (~100 tps)              |
| `MiniMax-M2.1`           | Multi-language programming capabilities         |
| `MiniMax-M2.1-highspeed` | Faster variant of M2.1 (~100 tps)               |
| `MiniMax-M2`             | Agentic capabilities, advanced reasoning        |

## How It Works

MiniMax is implemented as a built-in alias in docker-agent:

- **API Type:** OpenAI-compatible (`openai`)
- **Base URL:** `https://api.minimax.io/v1`
- **Token Variable:** `MINIMAX_API_KEY`

## Example: Code Assistant

```yaml
agents:
  coder:
    model: minimax/MiniMax-M2.5
    description: Code assistant using MiniMax
    instruction: |
      You are an expert programmer using MiniMax M2.5.
      Write clean, well-documented code.
      Follow best practices for the language being used.
    toolsets:
      - type: filesystem
      - type: shell
      - type: think
```
