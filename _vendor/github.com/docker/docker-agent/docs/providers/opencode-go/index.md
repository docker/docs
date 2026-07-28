---
title: "OpenCode Go"
description: "Use OpenCode Go models with docker-agent."
keywords: docker agent, ai agents, model providers, llm, opencode go
weight: 210
canonical: https://docs.docker.com/ai/docker-agent/providers/opencode-go/
---

_Use OpenCode Go models with docker-agent._

## Overview

[OpenCode Go](https://opencode.ai/docs/go) is a low-cost subscription service ($5 first month, then $10/month) that provides reliable access to popular open-source coding models. It serves models through both OpenAI-compatible and Anthropic-compatible APIs from globally distributed endpoints.

docker-agent includes built-in support for OpenCode Go as an alias provider.

## Setup

1. Subscribe to OpenCode Go at [opencode.ai/auth](https://opencode.ai/auth)
2. Copy your API key from the console
3. Set the environment variable:

   ```bash
   export OPENCODE_API_KEY=your-api-key
   ```

## Usage

### Inline Syntax

The simplest way to use OpenCode Go:

```yaml
agents:
  root:
    model: opencode-go/deepseek-v4-flash
    description: Assistant using OpenCode Go
    instruction: You are a helpful assistant.
```

### Named Model

For more control over parameters:

```yaml
models:
  my_model:
    provider: opencode-go
    model: deepseek-v4-pro
    temperature: 0.7
    max_tokens: 8192

agents:
  root:
    model: my_model
    description: Assistant using OpenCode Go
    instruction: You are a helpful assistant.
```

## Available Models

You can retrieve the full, up-to-date model list at any time:

```bash
curl https://opencode.ai/zen/go/v1/models
```

### OpenAI-Compatible

These models use the `/v1/chat/completions` endpoint and work directly with the `opencode-go` alias:

| Model               | Description                           |
| ------------------- | ------------------------------------- |
| `deepseek-v4-flash` | Fast and cost-effective DeepSeek model |
| `deepseek-v4-pro`   | Most capable DeepSeek model           |
| `kimi-k2.7-code`    | Kimi K2.7 optimized for code          |
| `kimi-k2.6`         | Kimi K2.6 model                       |
| `kimi-k2.5`         | Kimi K2.5 model                       |
| `glm-5.2`           | GLM 5.2 flagship model                |
| `glm-5.1`           | GLM 5.1 model                         |
| `glm-5`             | GLM 5 model                           |
| `mimo-v2.5`         | MiMo V2.5 efficient model             |
| `mimo-v2.5-pro`     | MiMo V2.5 Pro model                   |
| `mimo-v2-pro`       | MiMo V2 Pro model                     |
| `mimo-v2-omni`      | MiMo V2 Omni model                    |
| `hy3-preview`       | HY3 preview model                     |

### Anthropic-Compatible

These models use the `/v1/messages` endpoint and require a [custom provider definition](../custom/index.md):

| Model             | Description              |
| ----------------- | ------------------------ |
| `minimax-m3`      | MiniMax M3 model         |
| `minimax-m2.7`    | MiniMax M2.7 model       |
| `minimax-m2.5`    | MiniMax M2.5 model       |
| `qwen3.7-max`     | Qwen 3.7 Max model       |
| `qwen3.7-plus`    | Qwen 3.7 Plus model      |
| `qwen3.6-plus`    | Qwen 3.6 Plus model      |
| `qwen3.5-plus`    | Qwen 3.5 Plus model      |

To use an Anthropic-compatible model, define a custom provider:

```yaml
providers:
  opengo-ant:
    provider: anthropic
    base_url: https://opencode.ai/zen/go
    token_key: OPENCODE_API_KEY

models:
  qwen:
    provider: opengo-ant
    model: qwen3.7-max

agents:
  root:
    model: qwen
    description: Assistant using Qwen through OpenCode Go
    instruction: You are a helpful assistant.
```

## How It Works

OpenCode Go is implemented as a built-in alias in docker-agent:

- **API Type:** OpenAI-compatible (`openai_chatcompletions`)
- **Base URL:** `https://opencode.ai/zen/go/v1`
- **Token Variable:** `OPENCODE_API_KEY`

This means OpenCode Go uses the same client as OpenAI, making it fully compatible with all OpenAI features supported by docker-agent.

For Anthropic-compatible models (MiniMax, Qwen), docker-agent uses a custom provider pointing to the Anthropic client at `https://opencode.ai/zen/go` with the same token.

## Example: Code Assistant

```yaml
agents:
  coder:
    model: opencode-go/deepseek-v4-flash
    description: Expert code assistant
    instruction: |
      You are an expert programmer using DeepSeek V4 Flash.
      Write clean, efficient, well-documented code.
      Explain your reasoning when helpful.
    toolsets:
      - type: filesystem
      - type: shell
      - type: think
```

## Usage Limits

OpenCode Go subscriptions include the following limits:

- **5-hour rolling limit** — $12 of usage
- **Weekly limit** — $30 of usage
- **Monthly limit** — $60 of usage

Limits are defined as dollar values. More expensive models allow fewer requests per limit period. You can also [add Zen balance](https://opencode.ai/auth) to continue usage beyond the limits.
