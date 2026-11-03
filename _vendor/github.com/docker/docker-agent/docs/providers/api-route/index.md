---
title: "API Route"
description: "Use API Route models with Docker Agent."
keywords: docker agent, ai agents, model providers, llm, api route
weight: 23
canonical: https://docs.docker.com/ai/docker-agent/providers/api-route/
---

_Use API Route models with Docker Agent._

## Overview

[API Route](https://api-route.com) provides access to multiple model families
through an OpenAI-compatible API. Docker Agent includes API Route as a built-in
alias provider using the global endpoint.

## Setup

1. Create an API key in your [API Route account](https://api-route.com).
2. Set the environment variable:

   ```bash
   export API_ROUTE_API_KEY=your-api-key
   ```

## Usage

### Inline Syntax

```yaml
agents:
  root:
    model: api-route/gpt-6.1-sol
    description: Assistant using API Route
    instruction: You are a helpful assistant.
```

### Named Model

```yaml
models:
  api_route_model:
    provider: api-route
    model: gpt-6.1-sol

agents:
  root:
    model: api_route_model
    description: Assistant using API Route
    instruction: You are a helpful assistant.
```

## Available Models

The examples use the exact model ID `gpt-6.1-sol`. Retrieve your account's current
model IDs before choosing another model:

```bash
curl https://global.api-route.com/v1/models \
  -H "Authorization: Bearer $API_ROUTE_API_KEY"
```

Use an ID from the response unchanged as the `model` value. A catalog listing
does not guarantee that inference is currently available.

API Route is absent from Docker Agent's embedded models.dev catalog. Docker
Agent does not import pricing, context limits or attachment capabilities from
API Route's model list. Attachment detection falls back to text-only, so image
and PDF attachments are dropped with a warning. If you have verified support
for an attachment type on your chosen model and endpoint, set
`models.<name>.capabilities`; see
[Attachment Capability Overrides](../../configuration/models/index.md#attachment-capability-overrides).

## How It Works

- **API type:** OpenAI-compatible Chat Completions (`openai_chatcompletions`)
- **Base URL:** `https://global.api-route.com/v1`
- **Token variable:** `API_ROUTE_API_KEY`

As with other multi-model hosts, Docker Agent combines consecutive system
messages for compatibility with models that require a single leading system
message.

## Example: Code Assistant

```yaml
agents:
  coder:
    model: api-route/gpt-6.1-sol
    description: Code assistant using API Route
    instruction: |
      You are an expert programmer.
      Write clean code and follow language best practices.
    toolsets:
      - type: filesystem
      - type: shell
      - type: think
```
