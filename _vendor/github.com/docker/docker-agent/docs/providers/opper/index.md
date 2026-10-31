---
title: "Opper"
description: "Use Opper models with Docker Agent."
keywords: docker agent, ai agents, model providers, llm, opper
weight: 225
canonical: https://docs.docker.com/ai/docker-agent/providers/opper/
---

_Use Opper models with Docker Agent._

## Overview

[Opper](https://opper.ai) is an EU-hosted AI gateway that serves models from
many providers behind one OpenAI-compatible endpoint, with one API key.
Docker Agent includes built-in support for Opper as an alias provider.

## Setup

1. Create an API key at [platform.opper.ai](https://platform.opper.ai).
2. Set the environment variable:

   ```bash
   export OPPER_API_KEY=your-api-key
   ```

## Usage

### Inline Syntax

```yaml
agents:
  root:
    model: opper/claude-sonnet-4-6
    description: Assistant using Opper
    instruction: You are a helpful assistant.
```

### Named Model

```yaml
models:
  opper_model:
    provider: opper
    model: claude-sonnet-4-6
    temperature: 0.7
    max_tokens: 8192

agents:
  root:
    model: opper_model
    description: Assistant using Opper
    instruction: You are a helpful assistant.
```

## Available Models

Opper model IDs come in two forms. A bare pool name such as
`claude-sonnet-4-6` or `gpt-5.4-mini` lets Opper pick a provider for each
request. A `provider/model` ID such as `anthropic/claude-sonnet-4-6` pins one
route. List the model IDs your key can call with:

```bash
curl https://api.opper.ai/v3/compat/models \
  -H "Authorization: Bearer $OPPER_API_KEY"
```

The public catalogue is also available without a key at
`https://api.opper.ai/v3/models`. Docker Agent's default for Opper is
`claude-sonnet-4-6`.
