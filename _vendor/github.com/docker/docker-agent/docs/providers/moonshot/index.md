---
title: "Moonshot AI"
description: "Use Moonshot AI (Kimi) models with docker-agent."
keywords: docker agent, ai agents, model providers, llm, moonshot ai
weight: 180
canonical: https://docs.docker.com/ai/docker-agent/providers/moonshot/
---

_Use Moonshot AI (Kimi) models with docker-agent._

## Overview

[Moonshot AI](https://www.moonshot.ai/) serves its Kimi model family through an
OpenAI-compatible API. The Kimi K2 models have strong momentum for coding and
agentic tasks. docker-agent includes built-in support for Moonshot AI as an
alias provider.

## Setup

1. Create an API key from the [Moonshot AI console](https://platform.moonshot.ai/console/api-keys).
2. Set the environment variable:

   ```bash
   export MOONSHOT_API_KEY=your-api-key
   ```

## Usage

### Inline Syntax

The simplest way to use Moonshot AI:

```yaml
agents:
  root:
    model: moonshot/kimi-k2-0905-preview
    description: Assistant using Moonshot AI
    instruction: You are a helpful assistant.
```

### Named Model

For more control over parameters:

```yaml
models:
  moonshot_model:
    provider: moonshot
    model: kimi-k2-0905-preview
    temperature: 0.7
    max_tokens: 8192

agents:
  root:
    model: moonshot_model
    description: Assistant using Moonshot AI
    instruction: You are a helpful assistant.
```

## Available Models

Moonshot AI exposes a vendor-controlled Kimi model lineup. Check the
[Moonshot API documentation](https://platform.moonshot.ai/docs/api) for current
model IDs, context limits, and pricing.

| Model | Description |
| --- | --- |
| `kimi-k2-0905-preview` | Kimi K2, general-purpose chat, coding, and tool calling |
| `kimi-k2-turbo-preview` | Kimi K2 optimized for higher throughput |
| `kimi-k2-thinking` | Kimi K2 extended-reasoning model |

> Model IDs are case-sensitive and must be passed exactly as the catalogue lists
> them.

## How It Works

Moonshot AI is implemented as a built-in alias in docker-agent:

- **API Type:** OpenAI-compatible (`openai_chatcompletions`)
- **Base URL:** `https://api.moonshot.ai/v1`
- **Token Variable:** `MOONSHOT_API_KEY`

## Example: Code Assistant

```yaml
agents:
  coder:
    model: moonshot/kimi-k2-0905-preview
    description: Code assistant using Kimi K2
    instruction: |
      You are an expert programmer.
      Write clean, well-documented code and follow language best practices.
    toolsets:
      - type: filesystem
      - type: shell
      - type: think
```
