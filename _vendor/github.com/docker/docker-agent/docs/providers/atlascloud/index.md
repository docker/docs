---
title: "Atlas Cloud"
description: "Use Atlas Cloud models with Docker Agent."
keywords: docker agent, ai agents, model providers, llm, atlas cloud
weight: 25
canonical: https://docs.docker.com/ai/docker-agent/providers/atlascloud/
---

_Use Atlas Cloud models with Docker Agent._

## Overview

Atlas Cloud provides access to multiple model families through an OpenAI-compatible API. Docker Agent includes built-in support for Atlas Cloud as an alias provider.

## Setup

1. [Create an Atlas Cloud API key](https://www.atlascloud.ai/docs/api-keys)
2. Set the environment variable:

   ```bash
   export ATLASCLOUD_API_KEY=your-api-key
   ```

## Usage

### Inline Syntax

Use an Atlas Cloud model by prefixing its model ID with `atlascloud/`:

```yaml
agents:
  root:
    model: atlascloud/qwen/qwen3.8-max
    description: Assistant using Atlas Cloud
    instruction: You are a helpful assistant.
```

Docker Agent splits only the first slash, so model IDs that include a model family prefix are preserved.

### Named Model

For more control over parameters:

```yaml
models:
  atlas_qwen:
    provider: atlascloud
    model: qwen/qwen3.8-max
    temperature: 0.7
    max_tokens: 8192

agents:
  root:
    model: atlas_qwen
    description: Assistant using Atlas Cloud
    instruction: You are a helpful assistant.
```

## Available Models

Atlas Cloud's catalog includes the following exact model ID used by the examples:

| Model ID | Model |
| --- | --- |
| `qwen/qwen3.8-max` | Qwen3.8 Max |

Model IDs are case-sensitive. Check the [Atlas Cloud model catalog](https://www.atlascloud.ai/models)
or the live [`/v1/models` endpoint](https://api.atlascloud.ai/v1/models) for current IDs
and availability. A catalog listing does not guarantee that inference is currently available.

Atlas Cloud is not in the models.dev catalog. Docker Agent does not import metadata
from Atlas Cloud's model list, so it has no catalog-derived pricing or context limits
for these models. Attachment capabilities conservatively fall back to text-only;
image and PDF attachments are dropped, with a warning logged once per model.

If you have verified that your chosen model and endpoint accept a particular
attachment type, set `models.<name>.capabilities` to override detection. For example,
`capabilities: {image: true, pdf: false}` enables image attachments only. This is not
a claim that the example model supports them: enable only types the endpoint actually
accepts. See [Attachment Capability Overrides](../../configuration/models/index.md#attachment-capability-overrides)
for all supported flags.

## How It Works

Atlas Cloud is implemented as a built-in alias in Docker Agent:

- **API Type:** OpenAI-compatible Chat Completions (`openai_chatcompletions`)
- **Base URL:** `https://api.atlascloud.ai/v1`
- **Token Variable:** `ATLASCLOUD_API_KEY`

Because Atlas Cloud can serve models whose chat templates reject multiple leading
system messages, Docker Agent coalesces its per-source system messages into one.

## Example: Code Assistant

```yaml
agents:
  coder:
    model: atlascloud/qwen/qwen3.8-max
    description: Code assistant using Atlas Cloud
    instruction: |
      You are an expert programmer.
      Write clean code and follow language best practices.
    toolsets:
      - type: filesystem
      - type: shell
      - type: think
```
