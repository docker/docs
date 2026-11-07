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
    model: gpt-5.6-sol

agents:
  root:
    model: my_model
    instruction: You are a helpful assistant.
```

### Anthropic with shared defaults

Use a named Anthropic provider to share settings across models while allowing model-level overrides. See [Default Inheritance](#default-inheritance) for the example and precedence rules.

A custom Anthropic `base_url` reuses the Messages client but selects a different
host. It does not inherit Claude API capability, beta, price or signature-portability
guarantees. Intermediary prefix rewrites can invalidate signed history. The `azure`
alias uses the OpenAI-compatible client, not native Azure AI Foundry Claude.
These custom configurations are not live-tested by this guide; see
[Claude access paths](../anthropic/index.md#claude-access-paths).

### Google with shared temperature

```yaml
providers:
  my_google:
    provider: google
    temperature: 0.3

models:
  gemini:
    provider: my_google
    model: gemini-3.8-flash
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
| `parallel_tool_calls` | boolean    | Whether to enable parallel tool calls by default. When omitted, the provider/API default is used.         | —                        |
| `track_usage`         | boolean    | Whether to track token usage by default.                                              | —                        |
| `thinking_budget`     | string/int | Default reasoning effort/budget.                                                      | —                        |
| `task_budget`         | int/object | Default total token budget for an agentic task. See [Task Budget](../../configuration/models/index.md#task-budget) for syntax and [Anthropic](../anthropic/index.md#task-budget) for model support. | —                        |
| `compaction_model`    | string     | Default model used for session compaction (summary generation) by agents whose model uses this provider. Named model or inline `provider/model` string. Agent-level and model-level `compaction_model` take precedence. | —                        |
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
    model: claude-sonnet-5

  # Overrides temperature and thinking_budget, inherits the rest
  claude_custom:
    provider: my_anthropic
    model: claude-sonnet-5
    temperature: 0.2
    thinking_budget: low
```

`compaction_model` works slightly differently: it is not merged into the model
but resolved per agent, with the agent-level `compaction_model` winning over
the model-level one, which wins over the provider-level default.

## Shorthand Syntax

Once a provider is defined, you can use the shorthand `provider_name/model` syntax:

```yaml
agents:
  root:
    model: my_gateway/gpt-5.6-terra  # uses the provider's defaults
  researcher:
    model: my_anthropic/claude-sonnet-5  # uses anthropic provider defaults
```

## API Types

Only applicable for OpenAI-compatible providers (when `provider` is `openai` or unset):

- **`openai_chatcompletions`** — Standard OpenAI Chat Completions API. Works with most OpenAI-compatible endpoints.
- **`openai_responses`** — OpenAI Responses API. For newer models that require the Responses API format.

> If `api_type` is not set, Docker Agent automatically selects the API type based on the model name. You only need to set `api_type` explicitly to override the detected default.

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

### Replaying reasoning on later turns

Captured reasoning is sent back to the model on subsequent requests as a `reasoning_content` field on the assistant message. Servers that understand the field read it back, so the model sees its own earlier reasoning instead of re-deriving it after every tool call. Reasoning that is only useful to the model is not free: replaying it can multiply the size of a long request, so it is limited to turns the session keeps anyway (those with content or a tool call).

Chat Completions reasoning replay defaults on. Because `reasoning_content` is a
non-standard extension, set `provider_opts.replay_reasoning: false` for a model
or provider whose endpoint rejects or must not receive it. An invalid,
non-boolean value also fails closed and disables replay.

Set the opt-out on a model, or once on a provider and let its models inherit it.
A model that sets the option itself wins over the provider:

```yaml
providers:
  local_llm:
    base_url: http://localhost:8000/v1
    provider_opts:
      replay_reasoning: false # opt out every local_llm model

models:
  qwen_default:
    provider: local_llm
    model: qwen3 # inherits the opt-out

  qwen_replay:
    provider: local_llm
    model: qwen3
    provider_opts:
      replay_reasoning: true # opts this model back in

agents:
  root:
    model: qwen_default
```

This option is specific to the Chat Completions API and its `reasoning_content`
field. It is not the same as `preserve_reasoning`, which replays the encrypted
reasoning items of first-party OpenAI Responses sessions; neither option
applies to the other API. See [Reasoning state preservation](../openai/index.md#reasoning-state-preservation).

The field is worth a look if you switch models mid-session or use fallbacks:
Docker Agent stores reasoning from every provider in one place, so a later Chat
Completions model may otherwise receive reasoning captured from another vendor.
Set `replay_reasoning: false` on that destination when this is undesirable.

### API Router (Requesty, LiteLLM)

```yaml
providers:
  router:
    base_url: https://router.requesty.ai/v1
    token_key: REQUESTY_API_KEY

agents:
  root:
    model: router/anthropic/claude-sonnet-5
```

### Azure OpenAI

```yaml
models:
  azure_model:
    provider: azure
    model: gpt-5.6-sol
    base_url: https://your-llm.openai.azure.com
    provider_opts:
      api_version: 2024-12-01-preview
```

### Anthropic Team Setup

Agents can use different models backed by the same named provider. Starting from the [Default Inheritance](#default-inheritance) example:

```yaml
agents:
  root:
    model: claude_default
    instruction: You coordinate the development team.
    sub_agents: [code_reviewer]
  code_reviewer:
    model: claude_custom
    instruction: You review code.
```

Each model inherits provider defaults independently; its overrides do not affect the other model.

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
    model: smart_anthropic/claude-sonnet-5
    sub_agents: [helper]
  helper:
    model: fast_openai/gpt-4.1-mini
```

## Global Providers (User Configuration)

Providers defined in an agent file only apply to that file. To make a custom
provider available to every command (`docker agent run`, `new`, `models`, ...),
define it once in your user configuration (`~/.config/cagent/config.yaml`)
under the same `providers` key:

```yaml
# ~/.config/cagent/config.yaml
providers:
  myprovider:
    base_url: https://llm.corp.example.com/v1
    api_type: openai_chatcompletions
    token_key: MYPROVIDER_API_KEY
```

The easiest way to register one is the interactive wizard:

```bash
docker agent setup
# pick "3. Custom OpenAI-compatible endpoint", then enter the base URL,
# API format, and the environment variable holding the API key
```

Once registered, the provider works everywhere:

```bash
docker agent models --provider myprovider   # list the endpoint's models
docker agent new --model myprovider/mymodel # build agents with it
docker agent run --model myprovider/mymodel # chat with it
```

Global providers are merged into every loaded agent configuration; when an
agent file defines a provider with the same name, the agent file wins. Note
that automatic model selection (`model: auto`) never picks a custom provider,
so reference its models explicitly with `--model <name>/<model>` or
`default_model`.

## How It Works

When you reference a provider:

1. The provider's `provider` field determines which API client to use (defaults to `openai`)
2. The provider's `base_url` and `token_key` are applied to the model (if not already set on the model)
3. All model-level defaults (temperature, max_tokens, thinking_budget, etc.) are inherited (model settings take precedence)
4. For OpenAI-compatible providers, the `api_type` is stored in `provider_opts.api_type`
5. The model is used with the appropriate API client

A provider with a `base_url` implies `bypass_models_gateway: true` for every
model that references it: user-chosen endpoints are never routed through a
configured models gateway, and such models authenticate with the provider's
own credentials (`token_key`). See
[Gateway Bypass](../../configuration/models/index.md#gateway-bypass).
