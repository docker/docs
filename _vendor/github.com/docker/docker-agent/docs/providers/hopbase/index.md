---
title: "HopBase"
description: "Use HopBase models with Docker Agent."
keywords: docker agent, ai agents, model providers, llm, hopbase
weight: 135
canonical: https://docs.docker.com/ai/docker-agent/providers/hopbase/
---

_Use HopBase models with Docker Agent._

## Overview

[HopBase](https://hop-base.com) is a hosted AI API gateway. It serves GPT,
Claude, Gemini and Grok models next to DeepSeek, Zhipu GLM, Moonshot Kimi,
Alibaba Qwen and Doubao behind one base URL with an OpenAI-compatible API.
Docker Agent includes HopBase as a built-in provider in its embedded provider
catalog.

## Setup

1. Create an API key in the [HopBase console](https://api.hop-base.com/login).
2. Set the environment variable:

   ```bash
   export HOPBASE_API_KEY=your-api-key
   ```

## Usage

### Inline Syntax

Prefix a HopBase model ID with `hopbase/`:

```yaml
agents:
  root:
    model: hopbase/glm-5.3
    description: Assistant using HopBase
    instruction: You are a helpful assistant.
```

### Named Model

For more control over parameters:

```yaml
models:
  hopbase_gpt:
    provider: hopbase
    model: gpt-6.1-sol
    max_tokens: 8192

agents:
  root:
    model: hopbase_gpt
    description: Assistant using HopBase
    instruction: You are a helpful assistant.
```

When HopBase is selected automatically from available credentials, the default
model is `glm-5.3`. Model IDs are passed through as supplied.

## Available Models

The examples on this page use `glm-5.3`, `gpt-6.1-sol` and `claude-sonnet-5-5`.
The models endpoint returns exactly the model IDs your key can call, so
retrieve the list before choosing another model:

```bash
curl https://api.hop-base.com/v1/models \
  -H "Authorization: Bearer $HOPBASE_API_KEY"
```

Use an ID from the response unchanged as the `model` value. The list can also
contain image, video and speech models, which are not usable for chat. The
[model catalog](https://hop-base.com/models) and the
[HopBase docs](https://hop-base.com/en/docs) describe each model.

HopBase is absent from Docker Agent's embedded models.dev catalog. Docker Agent
does not import pricing, context limits or attachment capabilities from
HopBase's model list. Attachment detection falls back to text-only, so image
and PDF attachments are dropped, with a warning logged once per model. If you
have verified that your chosen model accepts an attachment type, set
`models.<name>.capabilities`; see
[Attachment Capability Overrides](../../configuration/models/index.md#attachment-capability-overrides).
For example, `gpt-6.1-sol` accepts image input and `glm-5.3` does not:

```yaml
models:
  hopbase_vision:
    provider: hopbase
    model: gpt-6.1-sol
    capabilities:
      image: true
```

## How It Works

HopBase is an entry in Docker Agent's built-in provider catalog
(`pkg/model/provider/catalog/providers.json`) that uses the existing OpenAI
implementation:

- **API type:** OpenAI-compatible Chat Completions (`openai_chatcompletions`)
- **Base URL:** `https://api.hop-base.com/v1`
- **Token variable:** `HOPBASE_API_KEY`

As with other multi-model hosts, Docker Agent combines consecutive system
messages for compatibility with models that require a single leading system
message.

Claude models are callable through this alias like any other model, for
example `hopbase/claude-sonnet-5-5`.

## Anthropic Messages Endpoint

The same host also exposes an Anthropic Messages-compatible endpoint at
`https://api.hop-base.com` (no `/v1`). To send Claude requests in the Messages
format instead of Chat Completions, define a
[named provider](../custom/index.md) that reuses the same key:

```yaml
providers:
  hopbase_messages:
    provider: anthropic
    base_url: https://api.hop-base.com
    token_key: HOPBASE_API_KEY

models:
  hopbase_claude:
    provider: hopbase_messages
    model: claude-sonnet-5-5

agents:
  root:
    model: hopbase_claude
    description: Assistant using Claude through HopBase
    instruction: You are a helpful assistant.
```

A custom Anthropic `base_url` does not inherit Claude API capability, beta or
price metadata; see
[Claude access paths](../anthropic/index.md#claude-access-paths).

## Example: Code Assistant

```yaml
agents:
  coder:
    model: hopbase/glm-5.3
    description: Code assistant using HopBase
    instruction: |
      You are an expert programmer.
      Write clean code and follow language best practices.
    toolsets:
      - type: filesystem
      - type: shell
      - type: think
```
