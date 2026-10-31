---
title: "OrcaRouter"
description: "Use OrcaRouter models with Docker Agent."
keywords: docker agent, ai agents, model providers, llm, orcarouter
weight: 231
canonical: https://docs.docker.com/ai/docker-agent/providers/orcarouter/
---

_Use OrcaRouter models with Docker Agent._

## Overview

[OrcaRouter](https://www.orcarouter.ai) provides access to models from many providers through an OpenAI-compatible API. Docker Agent includes built-in support for OrcaRouter as an alias provider.

## Setup

1. Get an API key from [OrcaRouter](https://www.orcarouter.ai)
2. Set the environment variable:

   ```bash
   $ export ORCAROUTER_API_KEY=your-api-key
   ```

## Usage

### Inline Syntax

Set the provider and model in the agent definition:

```yaml
agents:
  root:
    model: orcarouter/orcarouter/auto
    description: Assistant using OrcaRouter
    instruction: You are a helpful assistant.
```

`orcarouter/auto` is OrcaRouter's adaptive router. You can also reference a model served by the gateway with its upstream provider prefix, such as `orcarouter/anthropic/claude-sonnet-4.5` or `orcarouter/deepseek/deepseek-v4-pro`. Docker Agent splits only the first slash, so the full upstream model ID is preserved.

### Named Model

For more control over parameters:

```yaml
models:
  orca_router:
    provider: orcarouter
    model: orcarouter/auto
    temperature: 0.7
    max_tokens: 8192

agents:
  root:
    model: orca_router
    description: Assistant using OrcaRouter
    instruction: You are a helpful assistant.
```

## Pricing and Model Metadata

Docker Agent fetches OrcaRouter model metadata from [models.dev](https://models.dev/), including pricing per 1M input/output tokens, cache pricing when available, context limits, output limits, and modalities. This powers cost tracking and the model picker in the same way as other first-class providers.

If models.dev is unavailable, Docker Agent falls back to its embedded catalog snapshot.

## How It Works

OrcaRouter is implemented as a built-in alias in Docker Agent:

- **API Type:** OpenAI-compatible (`openai`)
- **Base URL:** `https://api.orcarouter.ai/v1`
- **Token Variable:** `ORCAROUTER_API_KEY`

## Example: Code Assistant

```yaml
agents:
  coder:
    model: orcarouter/orcarouter/auto
    description: Code assistant using OrcaRouter
    instruction: |
      You are an expert programmer.
      Write clean, maintainable code.
      Explain trade-offs when helpful.
    toolsets:
      - type: filesystem
      - type: shell
      - type: think
```
