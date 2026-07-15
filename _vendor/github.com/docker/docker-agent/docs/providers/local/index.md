---
title: "Local Models (Ollama, vLLM, LocalAI)"
description: "Run docker-agent with locally hosted models for privacy, offline use, or cost savings."
keywords: docker agent, ai agents, model providers, llm, local models (ollama, vllm, localai)
linkTitle: "Local Models"
weight: 150
canonical: https://docs.docker.com/ai/docker-agent/providers/local/
aliases:
  - /ai/docker-agent/local-models/
---

_Run docker-agent with locally hosted models for privacy, offline use, or cost savings._

## Overview

docker-agent can connect to any OpenAI-compatible local model server. This guide covers the most popular options:

- **Ollama** — Easy-to-use local model runner
- **vLLM** — High-performance inference server
- **LocalAI** — OpenAI-compatible API for various backends

> [!TIP]
> **Docker Model Runner**
>
> For the easiest local model experience, consider [Docker Model Runner](../dmr/index.md) which is built into Docker Desktop and requires no additional setup.

## Ollama

Ollama is a popular tool for running LLMs locally. docker-agent includes a built-in `ollama` alias for easy configuration.

### Setup

1. Install Ollama from [ollama.ai](https://ollama.ai/)
2. Pull a model:

   ```bash
   ollama pull llama3.2
   ollama pull qwen2.5-coder
   ```

3. Start the Ollama server (usually runs automatically):

   ```bash
   ollama serve
   ```

### Configuration

Use the built-in `ollama` alias:

```yaml
agents:
  root:
    model: ollama/llama3.2
    description: Local assistant
    instruction: You are a helpful assistant.
```

The `ollama` alias automatically uses:

- **Base URL:** `http://localhost:11434/v1`
- **API Type:** OpenAI-compatible
- **No API key required**

### Custom Port or Host

If Ollama runs on a different host or port:

```yaml
models:
  my_ollama:
    provider: ollama
    model: llama3.2
    base_url: http://192.168.1.100:11434/v1

agents:
  root:
    model: my_ollama
    description: Remote Ollama assistant
    instruction: You are a helpful assistant.
```

### Popular Ollama Models

| Model            | Size | Best For              |
| ---------------- | ---- | --------------------- |
| `llama3.2`       | 3B   | General purpose, fast |
| `llama3.1`       | 8B   | Better reasoning      |
| `qwen2.5-coder`  | 7B   | Code generation       |
| `mistral`        | 7B   | General purpose       |
| `codellama`      | 7B   | Code tasks            |
| `deepseek-coder` | 6.7B | Code generation       |

## vLLM

vLLM is a high-performance inference server optimized for throughput.

### Setup

```bash
# Install vLLM
pip install vllm

# Start the server
python -m vllm.entrypoints.openai.api_server \
  --model meta-llama/Llama-3.2-3B-Instruct \
  --port 8000
```

### Configuration

```yaml
providers:
  vllm:
    api_type: openai_chatcompletions
    base_url: http://localhost:8000/v1

agents:
  root:
    model: vllm/meta-llama/Llama-3.2-3B-Instruct
    description: vLLM-powered assistant
    instruction: You are a helpful assistant.
```

## LocalAI

LocalAI provides an OpenAI-compatible API that works with various backends.

### Setup

```bash
# Run with Docker
docker run -p 8080:8080 --name local-ai \
  -v ./models:/models \
  localai/localai:latest-cpu
```

### Configuration

```yaml
providers:
  localai:
    api_type: openai_chatcompletions
    base_url: http://localhost:8080/v1

agents:
  root:
    model: localai/gpt4all-j
    description: LocalAI assistant
    instruction: You are a helpful assistant.
```

## Generic Custom Provider

For any OpenAI-compatible server:

```yaml
providers:
  my_server:
    api_type: openai_chatcompletions
    base_url: http://localhost:8000/v1
    # token_key: MY_API_KEY  # if auth required

agents:
  root:
    model: my_server/model-name
    description: Custom server assistant
    instruction: You are a helpful assistant.
```

## Performance Tips

> [!NOTE]
> **Local Model Considerations**
>
> - **Memory:** Larger models need more RAM/VRAM. A 7B model typically needs 8-16GB RAM.
> - **GPU:** GPU acceleration dramatically improves speed. Check your server's GPU support.
> - **Context length:** Local models often have smaller context windows than cloud models.
> - **Tool calling:** Not all local models support function/tool calling. Test your model's capabilities.

## Example: Offline Development Agent

```yaml
agents:
  developer:
    model: ollama/qwen2.5-coder
    description: Offline code assistant
    instruction: |
      You are a software developer working offline.
      Focus on code quality and clear explanations.
    max_iterations: 20
    toolsets:
      - type: filesystem
      - type: shell
      - type: think
      - type: todo
```

## Troubleshooting

### Connection Refused

Ensure your model server is running and accessible:

```bash
curl http://localhost:11434/v1/models  # Ollama
curl http://localhost:8000/v1/models   # vLLM
```

### Model Not Found

Verify the model is downloaded/available:

```bash
ollama list  # List available Ollama models
```

### Slow Responses

- Check if GPU acceleration is enabled
- Try a smaller model
- Reduce `max_tokens` in your config
