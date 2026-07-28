---
title: "NVIDIA NIM"
description: "Use NVIDIA NIM models with docker-agent."
keywords: docker agent, ai agents, model providers, llm, nvidia, nim, nemotron
weight: 290
canonical: https://docs.docker.com/ai/docker-agent/providers/nvidia/
---

_Use NVIDIA NIM models with docker-agent._

## Overview

NVIDIA provides access to Nemotron and many other open-weight models through
[build.nvidia.com](https://build.nvidia.com/) (with a free tier) via an
OpenAI-compatible API. docker-agent includes built-in support for NVIDIA as an
alias provider. The same alias also works against a self-hosted
[NVIDIA NIM](https://docs.nvidia.com/nim/) deployment by overriding `base_url`.

## Setup

1. Get an API key from [build.nvidia.com](https://build.nvidia.com/)
2. Set the environment variable:

   ```bash
   export NVIDIA_API_KEY=your-api-key
   ```

## Usage

### Inline Syntax

The simplest way to use NVIDIA NIM:

```yaml
agents:
  root:
    model: nvidia/nvidia/nemotron-3-super-120b-a12b
    description: Assistant using NVIDIA NIM
    instruction: You are a helpful assistant.
```

### Named Model

For more control over parameters:

```yaml
models:
  nemotron:
    provider: nvidia
    model: nvidia/nemotron-3-super-120b-a12b
    temperature: 0.7
    max_tokens: 8192

agents:
  root:
    model: nemotron
    description: Assistant using NVIDIA NIM
    instruction: You are a helpful assistant.
```

## Available Models

NVIDIA NIM hosts Nemotron alongside many other open models (Llama, Qwen,
DeepSeek, Mistral, ...). Check the [NVIDIA API catalog](https://build.nvidia.com/)
for the current model list.

| Model                                       | Description                       |
| -------------------------------------------- | ---------------------------------- |
| `nvidia/nemotron-3-super-120b-a12b`         | Nemotron 3 Super, reasoning + tool calling |
| `nvidia/nemotron-3-nano-30b-a3b`            | Nemotron 3 Nano, smaller/faster    |
| `meta/llama-3.3-70b-instruct`               | Llama 3.3 70B instruction-tuned    |
| `qwen/qwen3-coder-480b-a35b-instruct`       | Qwen3 Coder, code-focused          |

## On-Prem / Self-Hosted NIM

For self-hosted NIM deployments, point `base_url` at your own endpoint instead
of the hosted `integrate.api.nvidia.com` API:

```yaml
models:
  local_nim:
    provider: nvidia
    model: meta/llama-3.3-70b-instruct
    base_url: http://localhost:8000/v1
```

## How It Works

NVIDIA is implemented as a built-in alias in docker-agent:

- **API Type:** OpenAI-compatible (`openai_chatcompletions`)
- **Base URL:** `https://integrate.api.nvidia.com/v1`
- **Token Variable:** `NVIDIA_API_KEY`

Because NIM fronts open-weight models whose chat templates often only accept
a single system message, docker-agent coalesces the agent instruction and any
toolset instructions into one leading system message before sending the
request.

## Example: Code Assistant

```yaml
agents:
  coder:
    model: nvidia/nvidia/nemotron-3-super-120b-a12b
    description: Code assistant using Nemotron
    instruction: |
      You are an expert programmer using NVIDIA Nemotron.
      Write clean, well-documented code.
      Follow best practices for the language being used.
    toolsets:
      - type: filesystem
      - type: shell
      - type: think
```
