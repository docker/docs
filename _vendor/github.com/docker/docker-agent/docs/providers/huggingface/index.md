---
title: "Hugging Face"
description: "Use Hugging Face Inference Providers with docker-agent."
keywords: docker agent, ai agents, model providers, llm, hugging face
weight: 140
canonical: https://docs.docker.com/ai/docker-agent/providers/huggingface/
---

_Use Hugging Face Inference Providers with docker-agent._

## Overview

[Hugging Face Inference Providers](https://huggingface.co/docs/inference-providers/index)
routes requests to open models (Llama, Qwen, DeepSeek, Kimi, GLM and others)
across many backends through a single OpenAI-compatible endpoint. docker-agent
includes built-in support for Hugging Face as an alias provider.

## Setup

1. Create a token from the [Hugging Face token settings](https://huggingface.co/settings/tokens).
2. Set the environment variable:

   ```bash
   export HF_TOKEN=your-token
   ```

## Usage

### Inline Syntax

The simplest way to use Hugging Face:

```yaml
agents:
  root:
    model: huggingface/meta-llama/Llama-3.3-70B-Instruct
    description: Assistant using Hugging Face Inference Providers
    instruction: You are a helpful assistant.
```

### Named Model

For more control over parameters:

```yaml
models:
  huggingface_model:
    provider: huggingface
    model: meta-llama/Llama-3.3-70B-Instruct
    temperature: 0.7
    max_tokens: 8192

agents:
  root:
    model: huggingface_model
    description: Assistant using Hugging Face Inference Providers
    instruction: You are a helpful assistant.
```

## Available Models

Hugging Face routes to a broad, changing catalog of open-weight models. Check the
[Hugging Face models page](https://huggingface.co/models?inference_provider=all)
for current model IDs, context limits, and pricing.

| Model | Description |
| --- | --- |
| `meta-llama/Llama-3.3-70B-Instruct` | Llama 3.3 70B, general-purpose chat and tool calling |
| `Qwen/Qwen3-235B-A22B` | Qwen3 235B mixture-of-experts instruct model |
| `deepseek-ai/DeepSeek-V3.2` | DeepSeek, strong coding and reasoning |

> Model IDs are case-sensitive and must be passed exactly as the catalogue lists
> them.

## How It Works

Hugging Face is implemented as a built-in alias in docker-agent:

- **API Type:** OpenAI-compatible (`openai_chatcompletions`)
- **Base URL:** `https://router.huggingface.co/v1`
- **Token Variable:** `HF_TOKEN`

Because Hugging Face fronts open-weight models whose chat templates may reject
more than one leading system message, docker-agent coalesces its per-source
system messages into a single one for this provider.

## Example: Code Assistant

```yaml
agents:
  coder:
    model: huggingface/Qwen/Qwen3-Coder-480B-A35B-Instruct
    description: Code assistant using Qwen3 Coder on Hugging Face
    instruction: |
      You are an expert programmer.
      Write clean, well-documented code and follow language best practices.
    toolsets:
      - type: filesystem
      - type: shell
      - type: think
```
