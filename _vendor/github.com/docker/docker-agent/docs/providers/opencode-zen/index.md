---
title: "OpenCode Zen"
description: "Use OpenCode Zen models with docker-agent."
keywords: docker agent, ai agents, model providers, llm, opencode zen
weight: 220
canonical: https://docs.docker.com/ai/docker-agent/providers/opencode-zen/
---

_Use OpenCode Zen models with docker-agent._

## Overview

[OpenCode Zen](https://opencode.ai/docs/zen) is a curated gateway of tested and verified AI models provided by the OpenCode team. It offers pay-per-use access to a wide range of models — from GPT and Claude to open-source models — all through a single API key. Several free models are also available.

docker-agent includes built-in support for OpenCode Zen as an alias provider for OpenAI-compatible models. Anthropic and Google models are supported via custom provider definitions.

## Setup

1. Sign in to [OpenCode Zen](https://opencode.ai/auth), add billing information, and copy your API key
2. Set the environment variable:

   ```bash
   export OPENCODE_API_KEY=your-api-key
   ```

3. Verify available models:

   ```bash
   curl https://opencode.ai/zen/v1/models
   ```

## Usage

### Inline Syntax

The simplest way to use OpenCode Zen with a free model:

```yaml
agents:
  root:
    model: opencode-zen/deepseek-v4-flash-free
    description: Assistant using OpenCode Zen (free)
    instruction: You are a helpful assistant.
```

### Named Model

For more control over parameters:

```yaml
models:
  zen_model:
    provider: opencode-zen
    model: gpt-5.5
    temperature: 0.7
    max_tokens: 16384

agents:
  root:
    model: zen_model
    description: Assistant using OpenCode Zen
    instruction: You are a helpful assistant.
```

## Available Models

### Free Models

These models are available at no cost:

| Model                     | Description                        |
| ------------------------- | ---------------------------------- |
| `deepseek-v4-flash-free`  | Free DeepSeek V4 Flash             |
| `mimo-v2.5-free`          | Free MiMo V2.5 model               |
| `qwen3.6-plus-free`       | Free Qwen 3.6 Plus model           |
| `minimax-m3-free`         | Free MiniMax M3 model              |
| `nemotron-3-ultra-free`   | Free Nemotron 3 Ultra model        |
| `north-mini-code-free`    | Free North Mini Code model         |
| `big-pickle`              | Free stealth model                 |

### OpenAI-Compatible (Chat Completions)

These models use the `/v1/chat/completions` endpoint and work directly with the `opencode-zen` alias:

| Model                 | Description                        |
| --------------------- | ---------------------------------- |
| `deepseek-v4-pro`     | DeepSeek V4 Pro model              |
| `deepseek-v4-flash`   | Fast DeepSeek V4 model             |
| `glm-5.2`             | GLM 5.2 model                      |
| `glm-5.1`             | GLM 5.1 model                      |
| `glm-5`               | GLM 5 model                        |
| `kimi-k2.6`           | Kimi K2.6 model                    |
| `kimi-k2.5`           | Kimi K2.5 model                    |
| `minimax-m2.7`        | MiniMax M2.7 model                 |
| `minimax-m2.5`        | MiniMax M2.5 model                 |
| `grok-build-0.1`      | Grok Build 0.1 model               |

### OpenAI-Compatible (Responses API)

These models use the `/v1/responses` endpoint and are auto-detected by docker-agent based on the model name:

| Model                  | Description                       |
| ---------------------- | --------------------------------- |
| `gpt-5.5`              | Latest GPT model                  |
| `gpt-5.5-pro`          | GPT 5.5 Pro, highest capability   |
| `gpt-5.4`              | GPT 5.4 model                     |
| `gpt-5.4-pro`          | GPT 5.4 Pro model                 |
| `gpt-5.4-mini`         | GPT 5.4 Mini model                |
| `gpt-5.4-nano`         | GPT 5.4 Nano, fastest             |
| `gpt-5.3-codex`        | GPT 5.3 Codex for coding          |
| `gpt-5.3-codex-spark`  | GPT 5.3 Codex Spark               |
| `gpt-5.2`              | GPT 5.2 model                     |
| `gpt-5.2-codex`        | GPT 5.2 Codex                     |
| `gpt-5.1`              | GPT 5.1 model                     |
| `gpt-5.1-codex`        | GPT 5.1 Codex                     |
| `gpt-5.1-codex-max`    | GPT 5.1 Codex Max                 |
| `gpt-5.1-codex-mini`   | GPT 5.1 Codex Mini                |
| `gpt-5`                | GPT 5 model                       |
| `gpt-5-codex`          | GPT 5 Codex                       |
| `gpt-5-nano`           | GPT 5 Nano                        |

### Anthropic-Compatible (Messages API)

These models use the `/v1/messages` endpoint and require a [custom provider definition](../custom/index.md):

| Model                   | Description                   |
| ----------------------- | ----------------------------- |
| `claude-fable-5`        | Claude Fable 5 model          |
| `claude-opus-4-8`       | Claude Opus 4.8 model         |
| `claude-opus-4-7`       | Claude Opus 4.7 model         |
| `claude-opus-4-6`       | Claude Opus 4.6 model         |
| `claude-opus-4-5`       | Claude Opus 4.5 model         |
| `claude-opus-4-1`       | Claude Opus 4.1 model         |
| `claude-sonnet-4-6`     | Claude Sonnet 4.6 model       |
| `claude-sonnet-4-5`     | Claude Sonnet 4.5 model       |
| `claude-sonnet-4`       | Claude Sonnet 4 model         |
| `claude-haiku-4-5`      | Claude Haiku 4.5 model        |
| `qwen3.7-max`           | Qwen 3.7 Max model            |
| `qwen3.7-plus`          | Qwen 3.7 Plus model           |
| `qwen3.6-plus`          | Qwen 3.6 Plus model           |
| `qwen3.5-plus`          | Qwen 3.5 Plus model           |

To use an Anthropic-compatible model:

```yaml
providers:
  opencode-zen-claude:
    provider: anthropic
    base_url: https://opencode.ai/zen
    token_key: OPENCODE_API_KEY

models:
  claude:
    provider: opencode-zen-claude
    model: claude-sonnet-4-5

agents:
  root:
    model: claude
    description: Assistant using Claude through OpenCode Zen
    instruction: You are a helpful assistant.
```

### Google-Compatible

These models require a [custom provider definition](../custom/index.md) with a Google-compatible client:

| Model               | Description              |
| ------------------- | ------------------------ |
| `gemini-3.5-flash`  | Gemini 3.5 Flash model   |
| `gemini-3.1-pro`    | Gemini 3.1 Pro model     |
| `gemini-3-flash`    | Gemini 3 Flash model     |

To use a Google model:

```yaml
providers:
  opencode-zen-gemini:
    provider: google
    base_url: https://opencode.ai/zen
    token_key: OPENCODE_API_KEY

models:
  gemini:
    provider: opencode-zen-gemini
    model: gemini-3.5-flash

agents:
  root:
    model: gemini
    description: Assistant using Gemini through OpenCode Zen
    instruction: You are a helpful assistant.
```

## How It Works

OpenCode Zen is implemented as a built-in alias in docker-agent:

- **API Type:** OpenAI-compatible (auto-detects Responses API for GPT models, Chat Completions for others)
- **Base URL:** `https://opencode.ai/zen/v1`
- **Token Variable:** `OPENCODE_API_KEY`

The same API key works for both OpenCode Go and OpenCode Zen — they are part of the same platform. Zen uses a pay-per-use billing model, while Go uses a fixed subscription.

For Anthropic-compatible models, docker-agent uses a custom provider pointing to the Anthropic client at `https://opencode.ai/zen` with the same token. For Google models, a custom provider points to the Google client at `https://opencode.ai/zen` (the Google SDK appends its own `/v1beta/models/...` path segment).

### Differences from OpenCode Go

| Aspect | OpenCode Zen | OpenCode Go |
|--------|-------------|-------------|
| Billing | Pay-per-use | $10/month subscription |
| Models | GPT-5.x, Claude, Gemini, open-source | Open-source only |
| Free models | Yes (7 models) | No |
| Base URL | `https://opencode.ai/zen/v1` | `https://opencode.ai/zen/go/v1` |

## Usage Limits and Pricing

OpenCode Zen uses a pay-per-use model. See the [OpenCode Zen documentation](https://opencode.ai/docs/zen) for current pricing. Automatic top-up and monthly usage limits are available from the console.

You can retrieve the full model catalog at any time:

```bash
curl https://opencode.ai/zen/v1/models
```
