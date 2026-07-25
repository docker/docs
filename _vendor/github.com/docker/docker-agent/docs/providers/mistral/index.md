---
title: "Mistral"
description: "Use Mistral AI models with docker-agent."
keywords: docker agent, ai agents, model providers, llm, mistral
weight: 170
canonical: https://docs.docker.com/ai/docker-agent/providers/mistral/
---

_Use Mistral AI models with docker-agent._

## Overview

Mistral AI provides powerful language models through an OpenAI-compatible API. docker-agent includes built-in support for Mistral as an alias provider.

## Setup

1. Get an API key from [Mistral Console](https://console.mistral.ai/)
2. Set the environment variable:

   ```bash
   export MISTRAL_API_KEY=your-api-key
   ```

## Usage

### Inline Syntax

The simplest way to use Mistral:

```yaml
agents:
  root:
    model: mistral/mistral-large-latest
    description: Assistant using Mistral
    instruction: You are a helpful assistant.
```

### Named Model

For more control over parameters:

```yaml
models:
  mistral:
    provider: mistral
    model: mistral-large-latest
    temperature: 0.7
    max_tokens: 8192

agents:
  root:
    model: mistral
    description: Assistant using Mistral
    instruction: You are a helpful assistant.
```

## Available Models

| Model                   | Description                       | Context |
| ----------------------- | --------------------------------- | ------- |
| `mistral-large-latest`  | Most capable Mistral model        | 128K    |
| `mistral-medium-latest` | Balanced performance and cost     | 128K    |
| `mistral-small-latest`  | Fast and cost-effective (default) | 128K    |
| `codestral-latest`      | Optimized for code generation     | 32K     |
| `open-mistral-nemo`     | Open-weight model                 | 128K    |
| `ministral-8b-latest`   | Compact 8B parameter model        | 128K    |
| `ministral-3b-latest`   | Smallest Mistral model            | 128K    |

Check the [Mistral Models documentation](https://docs.mistral.ai/getting-started/models/) for the latest available models.

## Auto-Detection

When you run `docker agent run` without specifying a config, docker-agent automatically detects available providers. If `MISTRAL_API_KEY` is set and higher-priority providers (OpenAI, Anthropic, Google) are not available, Mistral will be used with `mistral-small-latest` as the default model.

## Extended Thinking

docker-agent's `thinking_budget` field is **not applied** to Mistral models: the underlying OpenAI-compatible client only sends `reasoning_effort` for OpenAI reasoning model names (o-series, gpt-5). Setting `thinking_budget` on a Mistral model passes config validation but has no effect on the request.

Mistral reasoning models (e.g. `magistral`) reason on their own without configuration. For non-reasoning models, use the [think tool](../../tools/think/index.md) instead.

## How It Works

Mistral is implemented as a built-in alias in docker-agent:

- **API Type:** OpenAI-compatible (`openai_chatcompletions`)
- **Base URL:** `https://api.mistral.ai/v1`
- **Token Variable:** `MISTRAL_API_KEY`

This means Mistral uses the same client as OpenAI, making it fully compatible with all OpenAI features supported by docker-agent.

## Example: Code Assistant

```yaml
agents:
  coder:
    model: mistral/codestral-latest
    description: Expert code assistant
    instruction: |
      You are an expert programmer using Codestral.
      Write clean, efficient, well-documented code.
      Explain your reasoning when helpful.
    toolsets:
      - type: filesystem
      - type: shell
      - type: think
```
