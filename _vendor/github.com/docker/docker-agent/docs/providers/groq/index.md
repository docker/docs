---
title: "Groq"
description: "Use Groq fast-inference models with docker-agent."
keywords: docker agent, ai agents, model providers, llm, groq
weight: 130
canonical: https://docs.docker.com/ai/docker-agent/providers/groq/
---

_Use Groq models with docker-agent._

## Overview

[Groq](https://groq.com/) serves open-weight models on its LPU inference engine
through an OpenAI-compatible API, with a focus on very low latency. docker-agent
includes built-in support for Groq as an alias provider.

## Setup

1. Create an API key from the [Groq Console](https://console.groq.com/keys).
2. Set the environment variable:

   ```bash
   export GROQ_API_KEY=your-api-key
   ```

## Usage

### Inline Syntax

The simplest way to use Groq:

```yaml
agents:
  root:
    model: groq/llama-3.3-70b-versatile
    description: Assistant using Groq
    instruction: You are a helpful assistant.
```

### Named Model

For more control over parameters:

```yaml
models:
  groq_model:
    provider: groq
    model: llama-3.3-70b-versatile
    temperature: 0.7
    max_tokens: 8192

agents:
  root:
    model: groq_model
    description: Assistant using Groq
    instruction: You are a helpful assistant.
```

## Available Models

Groq hosts a rotating catalogue of open-weight models. Check the
[Groq models documentation](https://console.groq.com/docs/models) for current
model IDs, context limits, and rate limits.

| Model | Description |
| --- | --- |
| `llama-3.3-70b-versatile` | Llama 3.3 70B, reliable general-purpose chat and tool calling |
| `llama-3.1-8b-instant` | Llama 3.1 8B, fastest and cheapest |
| `openai/gpt-oss-120b` | GPT-OSS 120B, strong reasoning and tool calling |
| `openai/gpt-oss-20b` | GPT-OSS 20B, compact reasoning model |
| `qwen/qwen3-32b` | Qwen3 32B, reasoning and tool calling |
| `meta-llama/llama-4-scout-17b-16e-instruct` | Llama 4 Scout MoE |

> Model IDs are case-sensitive and must be passed exactly as the catalogue lists
> them.

## How It Works

Groq is implemented as a built-in alias in docker-agent:

- **API Type:** OpenAI-compatible (`openai_chatcompletions`)
- **Base URL:** `https://api.groq.com/openai/v1`
- **Token Variable:** `GROQ_API_KEY`

## Example: Code Assistant

```yaml
agents:
  coder:
    model: groq/llama-3.3-70b-versatile
    description: Code assistant using Llama 3.3
    instruction: |
      You are an expert programmer.
      Write clean, well-documented code and follow language best practices.
    toolsets:
      - type: filesystem
      - type: shell
      - type: think
```
