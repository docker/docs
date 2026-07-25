---
title: "OVHcloud"
description: "Use OVHcloud AI Endpoints models with docker-agent."
keywords: docker agent, ai agents, model providers, llm, ovhcloud
weight: 240
canonical: https://docs.docker.com/ai/docker-agent/providers/ovhcloud/
---

_Use OVHcloud AI Endpoints models with docker-agent._

## Overview

[OVHcloud AI Endpoints](https://endpoints.ai.cloud.ovh.net/) serves open-weight
models through an OpenAI-compatible API, hosted in the EU. docker-agent includes
built-in support for OVHcloud as an alias provider.

## Setup

1. Create an access token from the
   [OVHcloud AI Endpoints portal](https://endpoints.ai.cloud.ovh.net/).
2. Set the environment variable:

   ```bash
   export OVH_AI_ENDPOINTS_ACCESS_TOKEN=your-access-token
   ```

## Usage

### Inline Syntax

```yaml
agents:
  root:
    model: ovhcloud/Qwen3.5-397B-A17B
    description: Assistant using OVHcloud
    instruction: You are a helpful assistant.
```

### Named Model

```yaml
models:
  ovhcloud_model:
    provider: ovhcloud
    model: Qwen3.5-397B-A17B
    temperature: 0.7
    max_tokens: 8192

agents:
  root:
    model: ovhcloud_model
    description: Assistant using OVHcloud
    instruction: You are a helpful assistant.
```

## Available Models

OVHcloud hosts a rotating catalogue of open-weight models. Check the
[AI Endpoints catalogue](https://endpoints.ai.cloud.ovh.net/) for current model
IDs, context limits, and free-tier availability.

| Model | Description |
| --- | --- |
| `Qwen3.5-397B-A17B` | Large Qwen3.5 MoE — strong general, coding, and reasoning |
| `Qwen3-32B` | Mid-size Qwen3 — fast, tool-calling, reasoning |
| `Qwen3.6-27B` | Compact Qwen3.6 — fast and efficient |
| `Qwen3.5-9B` | Small Qwen3.5 — lightweight, free-tier friendly |
| `Qwen3-Coder-30B-A3B-Instruct` | Qwen3 Coder MoE — optimised for code generation |
| `Meta-Llama-3_3-70B-Instruct` | Llama 3.3 70B — reliable general-purpose chat |
| `Mistral-Small-3.2-24B-Instruct-2506` | Compact, fast, tool-calling |

> Model IDs are case-sensitive and must be passed exactly as the catalogue lists
> them.

## How It Works

OVHcloud is implemented as a built-in alias in docker-agent:

- **API Type:** OpenAI-compatible (`openai_chatcompletions`)
- **Base URL:** `https://oai.endpoints.kepler.ai.cloud.ovh.net/v1`
- **Token Variable:** `OVH_AI_ENDPOINTS_ACCESS_TOKEN`

docker-agent automatically coalesces consecutive system messages into one for
OVHcloud, because some OVHcloud models return an empty stream when a request
carries more than one system message.

## Free tier

OVHcloud offers rate-limited free access to several models. Under heavy
rate-limiting the endpoint may return an empty response; docker-agent surfaces
this as a warning rather than failing. For sustained use, an access token with a
paid plan avoids the free-tier request-rate cap.

## Example: Code Assistant

```yaml
agents:
  coder:
    model: ovhcloud/Qwen3.5-397B-A17B
    description: Code assistant using Qwen3.5
    instruction: |
      You are an expert programmer.
      Write clean, well-documented code and follow language best practices.
    toolsets:
      - type: filesystem
      - type: shell
      - type: think
```
