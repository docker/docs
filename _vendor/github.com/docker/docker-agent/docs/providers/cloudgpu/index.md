---
title: "CloudGPU"
description: "Use CloudGPU models with Docker Agent."
keywords: docker agent, ai agents, model providers, llm, cloudgpu
weight: 26
canonical: https://docs.docker.com/ai/docker-agent/providers/cloudgpu/
---

_Use CloudGPU models with Docker Agent._

## Overview

[CloudGPU](https://cloudgpu.app/api/docs) serves DeepSeek, Kimi, GLM and MiniMax
models through an OpenAI-compatible API. Docker Agent includes built-in support
for CloudGPU as an alias provider. Current model IDs are on the live
[model list](https://cloudgpu.app/v1/models).

## Setup

1. Create an API key in the [CloudGPU console](https://cloudgpu.app/)
2. Set the environment variable:

   ```bash
   export CLOUDGPU_API_KEY=your-api-key
   ```

## Usage

### Inline Syntax

Use a CloudGPU model by prefixing its model ID with `cloudgpu/`:

```yaml
agents:
  root:
    model: cloudgpu/deepseek-v4.1-flash
    description: Assistant using CloudGPU
    instruction: You are a helpful assistant.
```

### Named Model

For more control over parameters:

```yaml
models:
  cloudgpu_kimi:
    provider: cloudgpu
    model: kimi-k3
    temperature: 0.7
    max_tokens: 8192

agents:
  root:
    model: cloudgpu_kimi
    description: Assistant using CloudGPU
    instruction: You are a helpful assistant.
```

When CloudGPU is selected automatically from available credentials, the default
model is `deepseek-v4.1-flash`.

## Available Models

Some chat model IDs CloudGPU serves:

| Model ID | Model |
| --- | --- |
| `deepseek-v4.1-flash` | DeepSeek V4.1 Flash (default) |
| `deepseek-v4-pro` | DeepSeek V4 Pro |
| `kimi-k3` | Kimi K3 |
| `glm-5.2` | GLM 5.2 |
| `minimax-m3` | MiniMax M3 |

Model IDs are case-sensitive. The live [`/v1/models` endpoint](https://cloudgpu.app/v1/models)
lists current IDs, and the [API docs](https://cloudgpu.app/api/docs) describe
the endpoints. A listing does not guarantee that inference is currently available.

CloudGPU is absent from the embedded models.dev catalog. Docker Agent does not
import context limits or attachment capabilities from CloudGPU's model list.
Attachment capabilities fall back to text-only, so image and PDF attachments
are dropped, with a warning logged once per model.

If you have verified that your chosen model accepts a particular attachment
type, set `models.<name>.capabilities` to override detection. See
[Attachment Capability Overrides](../../configuration/models/index.md#attachment-capability-overrides)
for all supported flags.

## How It Works

CloudGPU is a built-in alias in Docker Agent:

- **API type:** OpenAI-compatible Chat Completions (`openai_chatcompletions`)
- **Base URL:** `https://cloudgpu.app/v1`
- **Token variable:** `CLOUDGPU_API_KEY`

Some models CloudGPU serves require a single leading system message, so Docker
Agent coalesces consecutive system messages into one.

## Example: Code Assistant

```yaml
agents:
  coder:
    model: cloudgpu/deepseek-v4.1-flash
    description: Code assistant using CloudGPU
    instruction: |
      You are an expert programmer.
      Write clean code and follow language best practices.
    toolsets:
      - type: filesystem
      - type: shell
      - type: think
```
