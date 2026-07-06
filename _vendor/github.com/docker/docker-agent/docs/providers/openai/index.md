---
title: "OpenAI"
description: "Use GPT-4o, GPT-5, GPT-5-mini, and other OpenAI models with docker-agent."
keywords: docker agent, ai agents, model providers, llm, openai
weight: 200
---

_Use GPT-4o, GPT-5, GPT-5-mini, and other OpenAI models with docker-agent._

## Setup

```bash
# Set your API key
export OPENAI_API_KEY="sk-..."
```

## Configuration

### Inline

```yaml
agents:
  root:
    model: openai/gpt-5
```

### Named Model

```yaml
models:
  gpt:
    provider: openai
    model: gpt-5
    temperature: 0.7
    max_tokens: 4000
```

## Available Models

| Model         | Best For                             |
| ------------- | ------------------------------------ |
| `gpt-5`       | Most capable, complex reasoning      |
| `gpt-5-mini`  | Fast, cost-effective, good reasoning |
| `gpt-4o`      | Multimodal, balanced performance     |
| `gpt-4o-mini` | Cheapest, fast for simple tasks      |

Find more model names at [modelnames.ai](https://modelnames.ai/) or in the [official OpenAI docs](https://platform.openai.com/docs/models).

## Thinking Budget

OpenAI reasoning models (o-series, gpt-5, gpt-5-mini) support extended thinking through the `reasoning_effort` API parameter. Set `thinking_budget` to control the effort level:

```yaml
models:
  gpt-thinker:
    provider: openai
    model: gpt-5-mini
    thinking_budget: high   # minimal | low | medium | high | xhigh
```

**Effort levels:**

| Level     | Description                                              |
| --------- | -------------------------------------------------------- |
| `none`    | Don't request extra reasoning (alias for `0`); the API's own default still applies. |
| `minimal` | Fastest; lightest reasoning pass.                        |
| `low`     | Quick reasoning for straightforward tasks.               |
| `medium`  | Balanced default.                                        |
| `high`    | More thorough; recommended for complex tasks.            |
| `xhigh`   | Near-maximum effort; slower but most accurate.           |

These are the **only** values OpenAI accepts — token counts, `max`, `adaptive`, and `adaptive/<effort>` are rejected with a configuration error at request time. Older models (o1, o3-mini) only accept `low`/`medium`/`high`.

> [!WARNING]
> **Hidden reasoning tokens**
>
> OpenAI reasoning models always produce hidden reasoning tokens that count against `max_tokens` — even with `thinking_budget: none`. docker-agent automatically raises the output-token floor for its internal low-effort calls so reasoning cannot starve visible text output.

See the [Thinking / Reasoning guide](../../guides/thinking/index.md) for a cross-provider overview.

> [!TIP]
> **Custom endpoints**
>
> Use `base_url` for proxies and OpenAI-compatible services. See [Custom Providers](../custom/index.md) for full setup.

## Custom Endpoint

Use `base_url` to connect to OpenAI-compatible APIs:

```yaml
models:
  custom:
    provider: openai
    model: gpt-5-mini
    base_url: https://your-proxy.example.com/v1
```

## WebSocket Transport

For OpenAI Responses API models (gpt-4.1+, o-series, gpt-5), you can use WebSocket streaming instead of the default SSE (Server-Sent Events):

```yaml
models:
  fast-gpt:
    provider: openai
    model: gpt-4.1
    provider_opts:
      transport: websocket  # Use WebSocket instead of SSE
```

### Benefits

- **~40% faster** for workflows with 20+ tool calls
- **Persistent connection** reduces per-turn overhead
- **Server-side caching** of connection state
- **Automatic fallback** to SSE if WebSocket fails

### Requirements

- Only works with Responses API models: `gpt-4.1+`, `o1`, `o3`, `o4`, `gpt-5`
- NOT compatible with the `--models-gateway` flag (automatically falls back to SSE when a gateway is configured)
- Requires `OPENAI_API_KEY` environment variable

### Example

See [`examples/websocket_transport.yaml`](https://github.com/docker/docker-agent/blob/main/examples/websocket_transport.yaml) for a complete example.
