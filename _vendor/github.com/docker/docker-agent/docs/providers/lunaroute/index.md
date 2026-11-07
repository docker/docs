---
title: "Lunaroute"
description: "Use Lunaroute models with Docker Agent."
keywords: docker agent, ai agents, model providers, llm, lunaroute
weight: 232
canonical: https://docs.docker.com/ai/docker-agent/providers/lunaroute/
---

_Use Lunaroute models with Docker Agent._

## Setup

Obtain an API key from [Lunaroute](https://lunaroute.com/). Lunaroute keys are
prefixed with `lr_`. Set the environment variable:

```bash
export LUNAROUTE_API_KEY=lr_your-api-key
```

## Usage

Use the inline model reference:

```yaml
agents:
  root:
    model: lunaroute/glm-5.3
    description: Assistant using Lunaroute
    instruction: You are a helpful assistant.
```

Or configure a named model:

```yaml
models:
  lunaroute_glm:
    provider: lunaroute
    model: glm-5.3

agents:
  root:
    model: lunaroute_glm
    instruction: You are a helpful assistant.
```

When Lunaroute is selected automatically from available credentials, the default
model is `glm-5.3`. Model IDs are passed through as supplied.

## Model discovery and capabilities

Run `docker agent models --provider lunaroute` with `LUNAROUTE_API_KEY` set to
query the authenticated `https://gw.lunaroute.com/v1/models` endpoint. The
listing also includes the configured default. Discovery is best-effort: if the
endpoint fails, the default remains available. The gateway's listing may contain
models that are not usable for chat.

Lunaroute is absent from the embedded models.dev catalog. Discovery imports model
IDs, not pricing, context limits, or capabilities; cost tracking and automatic
capability detection therefore lack catalog metadata. Unknown attachment
capabilities fall back to text-only, so image and PDF attachments are dropped
with a warning. If you have verified a model supports them, set a
[capability override](../../configuration/models/index.md#attachment-capability-overrides).

## How it works

Direct requests use OpenAI-compatible Chat Completions at
`https://gw.lunaroute.com/v1` with `LUNAROUTE_API_KEY`. When a Docker Agent
models gateway is configured, the built-in alias routes through that gateway
instead; the gateway must support the Lunaroute provider and model. To use
Lunaroute directly, set `bypass_models_gateway: true` on a named model. See
[Gateway Bypass](../../configuration/models/index.md#gateway-bypass).
