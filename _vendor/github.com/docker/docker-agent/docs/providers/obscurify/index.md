---
title: "Obscurify.ai"
description: "Use Obscurify.ai models with Docker Agent."
keywords: docker agent, ai agents, model providers, llm, obscurify
weight: 25
canonical: https://docs.docker.com/ai/docker-agent/providers/obscurify/
---

_Use Obscurify.ai models with Docker Agent._

## Overview

[Obscurify.ai](https://obscurify.ai/?utm_source=docker-agent) is a
privacy-focused provider with an OpenAI-compatible API and 200+ models.
Accounts are a random account number with no email, and prompts and responses
are not stored on its servers. Obscurify.ai is a built-in provider in Docker
Agent's embedded provider catalog.

The API works on every plan with a positive balance. Frontier models such as
`openai/gpt-5.5` or `anthropic/claude-sonnet-5.5` need a Casual or Pro plan;
see the [model list](https://obscurify.ai/model_cards) for prices.

## Setup

1. Create an account at [obscurify.ai](https://obscurify.ai/?utm_source=docker-agent)
   and copy your key from the [account page](https://obscurify.ai/account)
   ("Show API Key").
2. Set the environment variable:

   ```bash
   export OBSCURIFY_API_KEY=your-api-key
   ```

## Usage

### Inline Syntax

Prefix an Obscurify.ai model ID with `obscurify/`:

```yaml
agents:
  root:
    model: obscurify/deepseek/deepseek-v4-flash
    description: Assistant using Obscurify.ai
    instruction: You are a helpful assistant.
```

Only the first slash separates the provider, so the rest of the model ID
(`deepseek/deepseek-v4-flash`) is sent as is.

### Named Model

For more control over parameters:

```yaml
models:
  obscurify_kimi:
    provider: obscurify
    model: moonshotai/kimi-k2-0905
    temperature: 0.7
    max_tokens: 8192

agents:
  root:
    model: obscurify_kimi
    description: Assistant using Obscurify.ai
    instruction: You are a helpful assistant.
```

## Available Models

A few of the chat model IDs:

| Model ID | Context |
| --- | --- |
| `deepseek/deepseek-v4-flash` | 1,048,576 (default) |
| `z-ai/glm-5.3` | 1,048,576 |
| `moonshotai/kimi-k2-0905` | 262,144 |
| `qwen/qwen3-coder` | 262,000 |
| `meta-llama/llama-3.3-70b-instruct` | 131,072 |

The public [`/v1/models` endpoint](https://obscurify.ai/v1/models) lists every
current ID with its context length and pricing; no key is needed to read it.

Obscurify.ai is not in the models.dev catalog, so Docker Agent has no
catalog-derived pricing or context limits for these models. Attachment
capabilities fall back to text-only; image and PDF attachments are dropped,
with a warning logged once per model. If you know your model accepts a given
attachment type, set `models.<name>.capabilities` to override detection. See
[Attachment Capability Overrides](../../configuration/models/index.md#attachment-capability-overrides).

## How It Works

Obscurify.ai is an entry in Docker Agent's built-in provider catalog
(`pkg/model/provider/catalog/providers.json`) that uses the existing OpenAI
implementation:

- **API Type:** OpenAI-compatible Chat Completions (`openai_chatcompletions`)
- **Base URL:** `https://obscurify.ai/v1`
- **Token Variable:** `OBSCURIFY_API_KEY`

Many of its models are open-weight models whose chat templates can reject
several leading system messages, so Docker Agent merges its per-source system
messages into one.

## Example: Code Assistant

```yaml
agents:
  coder:
    model: obscurify/qwen/qwen3-coder
    description: Code assistant using Obscurify.ai
    instruction: |
      You are an expert programmer.
      Write clean code and follow language best practices.
    toolsets:
      - type: filesystem
      - type: shell
      - type: think
```
