---
title: "OpenRouter"
description: "Use OpenRouter models with docker-agent."
keywords: docker agent, ai agents, model providers, llm, openrouter
weight: 230
canonical: https://docs.docker.com/ai/docker-agent/providers/openrouter/
---

_Use OpenRouter models with docker-agent._

## Overview

OpenRouter provides access to models from many providers through an OpenAI-compatible API. docker-agent includes built-in support for OpenRouter as an alias provider.

## Setup

1. Get an API key from [OpenRouter](https://openrouter.ai/settings/keys)
2. Set the environment variable:

   ```bash
   export OPENROUTER_API_KEY=your-api-key
   ```

## Usage

### Inline Syntax

The simplest way to use OpenRouter:

```yaml
agents:
  root:
    model: openrouter/meta-llama/llama-3.3-70b-instruct
    description: Assistant using OpenRouter
    instruction: You are a helpful assistant.
```

OpenRouter model IDs usually include the upstream provider name, such as `anthropic/claude-sonnet-4-5` or `meta-llama/llama-3.3-70b-instruct`. docker-agent splits only the first slash, so the full upstream model ID is preserved.

### Named Model

For more control over parameters:

```yaml
models:
  openrouter_llama:
    provider: openrouter
    model: meta-llama/llama-3.3-70b-instruct
    temperature: 0.7
    max_tokens: 8192

agents:
  root:
    model: openrouter_llama
    description: Assistant using OpenRouter
    instruction: You are a helpful assistant.
```

## Pricing and Model Metadata

docker-agent fetches OpenRouter model metadata from [models.dev](https://models.dev/), including pricing per 1M input/output tokens, cache pricing when available, context limits, output limits, and modalities. This powers cost tracking and the model picker in the same way as other first-class providers.

If models.dev is unavailable, docker-agent falls back to its embedded catalog snapshot.

## How It Works

OpenRouter is implemented as a built-in alias in docker-agent:

- **API Type:** OpenAI-compatible (`openai`)
- **Base URL:** `https://openrouter.ai/api/v1`
- **Token Variable:** `OPENROUTER_API_KEY`

## Example: Code Assistant

```yaml
agents:
  coder:
    model: openrouter/meta-llama/llama-3.3-70b-instruct
    description: Code assistant using OpenRouter
    instruction: |
      You are an expert programmer.
      Write clean, maintainable code.
      Explain trade-offs when helpful.
    toolsets:
      - type: filesystem
      - type: shell
      - type: think
```
