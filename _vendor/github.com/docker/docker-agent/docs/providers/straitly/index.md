---
title: "Straitly"
description: "Use Straitly models with Docker Agent."
keywords: docker agent, ai agents, model providers, llm, straitly
weight: 247
canonical: https://docs.docker.com/ai/docker-agent/providers/straitly/
---

## Overview

[Straitly](https://straitly.ai) is an AI gateway. One API key and one base URL
reach more than 150 models, including Claude, GPT, Gemini, Grok and DeepSeek,
through an OpenAI-compatible API. Docker Agent includes Straitly as a built-in
provider in its embedded provider catalog.

## Setup

1. Sign up at [straitly.ai](https://straitly.ai/signup), add prepaid credits and
   create an API key on the [API keys page](https://straitly.ai/app/keys).
2. Set the environment variable:

   ```bash
   export STRAITLY_API_KEY=your-api-key
   ```

## Usage

### Inline Syntax

Prefix a Straitly model ID with `straitly/`:

```yaml
agents:
  root:
    model: straitly/deepseek/deepseek-v4-1-flash
    description: Assistant using Straitly
    instruction: You are a helpful assistant.
```

Straitly model IDs usually contain a slash themselves. Docker Agent splits the
inline reference at the first slash, so everything after `straitly/` is passed
to Straitly unchanged.

### Named Model

```yaml
models:
  straitly_model:
    provider: straitly
    model: deepseek/deepseek-v4-1-flash
    max_tokens: 8192

agents:
  root:
    model: straitly_model
    description: Assistant using Straitly
    instruction: You are a helpful assistant.
```

When Straitly is selected automatically from available credentials, the default
model is `deepseek/deepseek-v4-1-flash`. Model IDs are passed through as
supplied.

## Available Models

List the model IDs before choosing another model:

```bash
curl https://api.straitly.ai/v1/models \
  -H "Authorization: Bearer $STRAITLY_API_KEY"
```

Use an ID from the response unchanged as the `model` value. Each entry also
lists its `context_length`, `supported_parameters` and
`architecture.input_modalities`. The same list is on the
[models page](https://straitly.ai/models).

Straitly is absent from Docker Agent's embedded models.dev catalog. Docker Agent
does not import pricing, context limits or attachment capabilities from
Straitly's model list. Attachment detection falls back to text-only, so image
and PDF attachments are dropped, with a warning logged once per model. If the
model you chose lists `image` in `architecture.input_modalities`, set
`models.<name>.capabilities`; see
[Attachment Capability Overrides](../../configuration/models/index.md#attachment-capability-overrides):

```yaml
models:
  straitly_vision:
    provider: straitly
    model: deepseek/deepseek-v4-1-flash
    capabilities:
      image: true
```

## How It Works

Straitly is an entry in Docker Agent's built-in provider catalog
(`pkg/model/provider/catalog/providers.json`) that uses the existing OpenAI
implementation:

- **API type:** OpenAI-compatible Chat Completions (`openai_chatcompletions`)
- **Base URL:** `https://api.straitly.ai/v1`
- **Token variable:** `STRAITLY_API_KEY`

Claude models are callable through this alias like any other model. Use the ID
exactly as `/v1/models` returns it.

## Example: Code Assistant

```yaml
agents:
  coder:
    model: straitly/deepseek/deepseek-v4-1-flash
    description: Code assistant using Straitly
    instruction: |
      You are an expert programmer.
      Write clean code and follow language best practices.
    toolsets:
      - type: filesystem
      - type: shell
      - type: think
```
