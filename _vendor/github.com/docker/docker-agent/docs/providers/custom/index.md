---
title: "Provider Definitions"
description: "Define reusable provider configurations with shared defaults for any provider type — OpenAI, Anthropic, Google, Bedrock, and more."
keywords: docker agent, ai agents, model providers, llm, provider definitions
weight: 280
canonical: https://docs.docker.com/ai/docker-agent/providers/custom/
---

_Define reusable provider configurations with shared defaults for any provider type — OpenAI, Anthropic, Google, Bedrock, and more._

## Overview

The `providers` section in your agent YAML lets you define named provider configurations that models can reference. This is useful for:

- **Grouping shared defaults** — Set temperature, max_tokens, thinking_budget once and share across models
- **Custom endpoints** — Connect to self-hosted models, API proxies, or gateways
- **Centralizing credentials** — Define token_key once for all models using a provider
- **Any provider type** — Works with OpenAI, Anthropic, Google, Bedrock, and any OpenAI-compatible API

> [!NOTE]
> **Works with any provider**
>
> The `providers` section supports all provider types: `openai`, `anthropic`, `google`, `amazon-bedrock`, `dmr`, and any built-in alias. When the `provider` field is not set, it defaults to `openai` for backward compatibility.

## Configuration

### OpenAI-compatible endpoint

```yaml
providers:
  my_gateway:
    base_url: https://api.example.com/v1
    token_key: MY_API_KEY

models:
  my_model:
    provider: my_gateway
    model: gpt-4o

agents:
  root:
    model: my_model
    instruction: You are a helpful assistant.
```

### Anthropic with shared defaults

```yaml
providers:
  my_anthropic:
    provider: anthropic
    token_key: MY_ANTHROPIC_KEY
    max_tokens: 16384
    thinking_budget: 8192

models:
  claude_smart:
    provider: my_anthropic
    model: claude-sonnet-4-5
    # Inherits max_tokens: 16384, thinking_budget: 8192

  claude_fast:
    provider: my_anthropic
    model: claude-haiku-4-5
    thinking_budget: 1024  # Overrides provider default

agents:
  root:
    model: claude_smart
    instruction: You are a helpful assistant.
```

### Google with shared temperature

```yaml
providers:
  my_google:
    provider: google
    temperature: 0.3

models:
  gemini:
    provider: my_google
    model: gemini-2.5-flash
    # Inherits temperature: 0.3

agents:
  root:
    model: gemini
    instruction: You are a helpful assistant.
```

## Provider Properties

| Property              | Type       | Description                                                                           | Default                  |
| --------------------- | ---------- | ------------------------------------------------------------------------------------- | ------------------------ |
| `provider`            | string     | Underlying provider type: `openai`, `anthropic`, `google`, `amazon-bedrock`, `dmr`, etc. | `openai`                 |
| `api_type`            | string     | API schema: `openai_chatcompletions` or `openai_responses`. Only for OpenAI-compatible providers. When omitted, the API type is selected automatically based on the model name: newer models (gpt-4.1, o-series, gpt-5, Codex) default to `openai_responses`; all others default to `openai_chatcompletions`. | `auto (model-dependent)` |
| `base_url`            | string     | Base URL for the API endpoint. Required for OpenAI-compatible providers, optional for native providers. | —                        |
| `token_key`           | string     | Environment variable name containing the API token.                                   | —                        |
| `unload_api`          | string     | Optional path (or absolute URL) to the provider's model-unload endpoint. Used by the [`unload`](../../configuration/hooks/index.md#available-built-ins) built-in hook to release model resources between agent switches. Relative paths resolve against `base_url`'s scheme + host; absolute URLs are used verbatim. Today only Docker Model Runner ships a provider that calls this endpoint; cloud providers don't implement the underlying interface and the hook silently skips them. | —                        |
| `temperature`         | float      | Default sampling temperature (0.0–2.0).                                               | —                        |
| `max_tokens`          | int        | Default maximum response tokens.                                                      | —                        |
| `top_p`               | float      | Default nucleus sampling threshold (0.0–1.0).                                         | —                        |
| `frequency_penalty`   | float      | Default frequency penalty (-2.0–2.0).                                                 | —                        |
| `presence_penalty`    | float      | Default presence penalty (-2.0–2.0).                                                  | —                        |
| `parallel_tool_calls` | boolean    | Whether to enable parallel tool calls by default.                                     | —                        |
| `track_usage`         | boolean    | Whether to track token usage by default.                                              | —                        |
| `thinking_budget`     | string/int | Default reasoning effort/budget.                                                      | —                        |
| `task_budget`         | int/object | Default total token budget for an agentic task (forwarded to Anthropic; honored by Claude Opus 4.7 today). Integer shorthand or `{type: tokens, total: N}`. | —                        |
| `provider_opts`       | object     | Provider-specific options passed through to the client.                               | —                        |

## Default Inheritance

Models referencing a provider inherit all its defaults. Model-level settings always take precedence:

```yaml
providers:
  my_anthropic:
    provider: anthropic
    token_key: MY_ANTHROPIC_KEY
    max_tokens: 16384
    temperature: 0.7
    thinking_budget: high

models:
  # Inherits everything from provider
  claude_default:
    provider: my_anthropic
    model: claude-sonnet-4-5

  # Overrides temperature and thinking_budget, inherits the rest
  claude_custom:
    provider: my_anthropic
    model: claude-sonnet-4-5
    temperature: 0.2
    thinking_budget: low
```

## Shorthand Syntax

Once a provider is defined, you can use the shorthand `provider_name/model` syntax:

```yaml
agents:
  root:
    model: my_gateway/gpt-4o-mini  # uses the provider's defaults
  researcher:
    model: my_anthropic/claude-sonnet-4-5  # uses anthropic provider defaults
```

## API Types

Only applicable for OpenAI-compatible providers (when `provider` is `openai` or unset):

- **`openai_chatcompletions`** — Standard OpenAI Chat Completions API. Works with most OpenAI-compatible endpoints.
- **`openai_responses`** — OpenAI Responses API. For newer models that require the Responses API format.

> If `api_type` is not set, docker-agent automatically selects the API type based on the model name. You only need to set `api_type` explicitly to override the detected default.

## Examples

### vLLM / Ollama

```yaml
providers:
  local_llm:
    base_url: http://localhost:8000/v1

agents:
  root:
    model: local_llm/llama-3.1-8b
```

> [!NOTE]
> **Reasoning tokens from OpenAI-compatible providers**
>
> Models that stream reasoning under `delta.reasoning` (e.g. Qwen3 served via OVHcloud AI Endpoints, OpenRouter, or a self-hosted vLLM / SGLang deployment) are fully supported. Docker Agent reads both the `delta.reasoning_content` and `delta.reasoning` fields from the stream, so thinking blocks are captured and shown in the TUI regardless of which field the server uses.

### API Router (Requesty, LiteLLM)

```yaml
providers:
  router:
    base_url: https://router.requesty.ai/v1
    token_key: REQUESTY_API_KEY

agents:
  root:
    model: router/anthropic/claude-sonnet-4-5
```

### Azure OpenAI

```yaml
models:
  azure_model:
    provider: azure
    model: gpt-4o
    base_url: https://your-llm.openai.azure.com
    provider_opts:
      api_version: 2024-12-01-preview
```

### Anthropic Team Setup

```yaml
providers:
  team_anthropic:
    provider: anthropic
    token_key: TEAM_ANTHROPIC_KEY
    max_tokens: 32768
    thinking_budget: high
    temperature: 0.5

models:
  architect:
    provider: team_anthropic
    model: claude-sonnet-4-5

  reviewer:
    provider: team_anthropic
    model: claude-haiku-4-5
    thinking_budget: low  # faster reviews

agents:
  root:
    model: architect
    sub_agents: [code_reviewer]
  code_reviewer:
    model: reviewer
```

### Multi-Provider with Shared Defaults

```yaml
providers:
  fast_openai:
    base_url: https://api.openai.com/v1
    token_key: OPENAI_API_KEY
    temperature: 0.3
    max_tokens: 8192

  smart_anthropic:
    provider: anthropic
    token_key: ANTHROPIC_API_KEY
    max_tokens: 64000
    thinking_budget: high

agents:
  root:
    model: smart_anthropic/claude-sonnet-4-5
    sub_agents: [helper]
  helper:
    model: fast_openai/gpt-4o-mini
```

## How It Works

When you reference a provider:

1. The provider's `provider` field determines which API client to use (defaults to `openai`)
2. The provider's `base_url` and `token_key` are applied to the model (if not already set on the model)
3. All model-level defaults (temperature, max_tokens, thinking_budget, etc.) are inherited (model settings take precedence)
4. For OpenAI-compatible providers, the `api_type` is stored in `provider_opts.api_type`
5. The model is used with the appropriate API client
