---
title: "DaoXE"
description: "Use DaoXE models with Docker Agent."
keywords: docker agent, ai agents, model providers, llm, daoxe
weight: 245
canonical: https://docs.docker.com/ai/docker-agent/providers/daoxe/
---

_Use DaoXE models with Docker Agent._

## Overview

[DaoXE](https://daoxe.com) is a multi-model LLM gateway that serves models from
multiple vendors behind one OpenAI-compatible endpoint. Docker Agent includes
built-in support for DaoXE as an alias provider.

## Setup

1. Create an API key from your DaoXE account dashboard.
2. Set the environment variable:

   ```bash
   export DAOXE_API_KEY=your-api-key
   ```

## Usage

### Inline Syntax

```yaml
agents:
  root:
    model: daoxe/<model-id>
    description: Assistant using DaoXE
    instruction: You are a helpful assistant.
```

### Named Model

```yaml
models:
  daoxe_model:
    provider: daoxe
    model: <model-id>
    temperature: 0.7
    max_tokens: 8192

agents:
  root:
    model: daoxe_model
    description: Assistant using DaoXE
    instruction: You are a helpful assistant.
```

## Available Models

DaoXE's catalog is account-scoped and changes over time: new models appear as
upstream vendors release them. List the model IDs your account can call with:

```bash
curl https://api.daoxe.com/v1/models \
  -H "Authorization: Bearer $DAOXE_API_KEY"
```

Use any ID from that response as the `model` value. Docker Agent's default for
DaoXE is `claude-sonnet-4-6`.
