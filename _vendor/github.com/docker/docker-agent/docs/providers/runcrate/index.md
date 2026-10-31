---
title: "Runcrate"
description: "Use Runcrate models with Docker Agent."
keywords: docker agent, ai agents, model providers, llm, runcrate
weight: 25
canonical: https://docs.docker.com/ai/docker-agent/providers/runcrate/
---

_Use Runcrate models with Docker Agent._

## Overview

[Runcrate](https://www.runcrate.ai/?utm_source=docker-agent&utm_medium=directory&utm_campaign=gpu_pricing)
is a GPU cloud that also serves open models (DeepSeek, GLM, Qwen, gpt-oss and
others) through an OpenAI-compatible API. Docker Agent includes built-in
support for Runcrate as an alias provider.

Some open models come with a limited free allowance for GitHub accounts at
least 30 days old. Once the allowance is used up, or when free capacity is
busy, requests from a workspace with credit are billed to that credit. See
[Runcrate's free inference terms](https://www.runcrate.ai/inference/free?utm_source=docker-agent&utm_medium=directory&utm_campaign=gpu_pricing)
for the current models, limits and prices.

## Setup

1. [Create a Runcrate API key](https://www.runcrate.ai/dashboard/api-keys)
2. Set the environment variable:

   ```bash
   export RUNCRATE_API_KEY=your-api-key
   ```

## Usage

### Inline Syntax

Use a Runcrate model by prefixing its model ID with `runcrate/`:

```yaml
agents:
  root:
    model: runcrate/deepseek-ai/DeepSeek-V4.1-Flash
    description: Assistant using Runcrate
    instruction: You are a helpful assistant.
```

Docker Agent splits only the first slash, so model IDs that include an
organization prefix are preserved.

### Named Model

For more control over parameters:

```yaml
models:
  runcrate_glm:
    provider: runcrate
    model: zai-org/GLM-5.3
    temperature: 0.7
    max_tokens: 8192

agents:
  root:
    model: runcrate_glm
    description: Assistant using Runcrate
    instruction: You are a helpful assistant.
```

## Available Models

Some chat model IDs Runcrate serves:

| Model ID | Model |
| --- | --- |
| `deepseek-ai/DeepSeek-V4.1-Flash` | DeepSeek V4.1 Flash (default) |
| `zai-org/GLM-5.3` | GLM-5.3 |
| `Qwen/Qwen3.8-Flash` | Qwen3.8 Flash |
| `openai/gpt-oss-120b` | gpt-oss 120B |

Model IDs are case-sensitive. The live [`/v1/models` endpoint](https://api.runcrate.ai/v1/models)
lists current IDs; it also lists image, video, speech and embedding models,
which this chat provider does not use.

Runcrate is not in the models.dev catalog. Docker Agent does not import metadata
from Runcrate's model list, so it has no catalog-derived pricing or context limits
for these models. Attachment capabilities conservatively fall back to text-only;
image and PDF attachments are dropped, with a warning logged once per model.

If you have verified that your chosen model accepts a particular attachment
type, set `models.<name>.capabilities` to override detection. See
[Attachment Capability Overrides](../../configuration/models/index.md#attachment-capability-overrides)
for all supported flags.

## How It Works

Runcrate is implemented as a built-in alias in Docker Agent:

- **API Type:** OpenAI-compatible Chat Completions (`openai_chatcompletions`)
- **Base URL:** `https://api.runcrate.ai/v1`
- **Token Variable:** `RUNCRATE_API_KEY`

Because Runcrate serves open-weight models whose chat templates can reject
multiple leading system messages, Docker Agent coalesces its per-source system
messages into one.

## Example: Code Assistant

```yaml
agents:
  coder:
    model: runcrate/deepseek-ai/DeepSeek-V4.1-Flash
    description: Code assistant using Runcrate
    instruction: |
      You are an expert programmer.
      Write clean code and follow language best practices.
    toolsets:
      - type: filesystem
      - type: shell
      - type: think
```
