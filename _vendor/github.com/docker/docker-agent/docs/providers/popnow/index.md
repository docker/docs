---
title: "PopNow"
description: "Use PopNow models with Docker Agent."
keywords: docker agent, ai agents, model providers, llm, popnow
weight: 25
canonical: https://docs.docker.com/ai/docker-agent/providers/popnow/
---

_Use PopNow models with Docker Agent._

## Overview

[PopNow](https://popnow.ai/api) provides chat models from several vendors through one OpenAI-compatible API, billed per token in USD. Docker Agent includes built-in support for PopNow as a provider alias.

## Setup

1. [Create a PopNow API key](https://popnow.ai/api/console/keys)
2. Set the environment variable:

   ```bash
   export POPNOW_API_KEY=your-api-key
   ```

## Usage

### Inline Syntax

Use a PopNow model by prefixing its model ID with `popnow/`:

```yaml
agents:
  root:
    model: popnow/deepseek/v4-flash/chat
    description: Assistant using PopNow
    instruction: You are a helpful assistant.
```

Docker Agent splits only the first slash, so the vendor prefix and task suffix in PopNow model IDs are preserved.

### Named Model

For more control over parameters:

```yaml
models:
  popnow_deepseek:
    provider: popnow
    model: deepseek/v4-flash/chat
    temperature: 0.7
    max_tokens: 8192

agents:
  root:
    model: popnow_deepseek
    description: Assistant using PopNow
    instruction: You are a helpful assistant.
```

## Available Models

PopNow model IDs end in a task suffix, and chat models use `/chat`:

| Model ID | Model |
| --- | --- |
| `deepseek/v4-flash/chat` | DeepSeek V4 Flash |
| `moonshot/kimi-k3/chat` | Kimi K3 |

Model IDs are case-sensitive. Check the [PopNow model catalog](https://popnow.ai/api/models) for current IDs and availability. Only models with the `chat` task work with Docker Agent.

PopNow is not in the models.dev catalog, so Docker Agent has no catalog-derived pricing or context limits for these models. Attachment capabilities conservatively fall back to text-only; image and PDF attachments are dropped, with a warning logged once per model.

If you have verified that your chosen model accepts a particular attachment type, set `models.<name>.capabilities` to override detection. See [Attachment Capability Overrides](../../configuration/models/index.md#attachment-capability-overrides) for all supported flags.

## How It Works

PopNow is implemented as a built-in alias in Docker Agent:

- **API Type:** OpenAI-compatible Chat Completions (`openai_chatcompletions`)
- **Base URL:** `https://api.popnow.ai/v1`
- **Token Variable:** `POPNOW_API_KEY`

Because PopNow serves models from many vendors, some of which reject multiple leading system messages, Docker Agent coalesces its per-source system messages into one.

## Example: Code Assistant

```yaml
agents:
  coder:
    model: popnow/deepseek/v4-flash/chat
    description: Code assistant using PopNow
    instruction: |
      You are an expert programmer.
      Write clean code and follow language best practices.
    toolsets:
      - type: filesystem
      - type: shell
      - type: think
```
