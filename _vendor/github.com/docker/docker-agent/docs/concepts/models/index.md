---
title: "Models"
description: "Models are the AI brains behind your agents. Docker Agent supports multiple providers and flexible configuration."
keywords: docker agent, ai agents, concepts, models
weight: 20
canonical: https://docs.docker.com/ai/docker-agent/concepts/models/
---

_Models are the AI brains behind your agents. Docker Agent supports multiple providers and flexible configuration._

## Inline vs. Named Models

There are two ways to assign a model to an agent:

### Inline (Quick)

Use the `provider/model` shorthand directly in the agent definition:

```yaml
agents:
  root:
    model: openai/gpt-5
    instruction: You are a helpful assistant.
```

### Named (Full Control)

Define models in a `models` section and reference them by name:

```yaml
models:
  claude:
    provider: anthropic
    model: claude-sonnet-4-5
    max_tokens: 64000
    temperature: 0.7

agents:
  root:
    model: claude
    instruction: You are a helpful assistant.
```

Named models let you configure temperature, token limits, thinking budgets, and other parameters. They're also reusable across multiple agents.

## First Available Models

A named model can also select the first usable model from a priority list. This
is useful for shared configs that should prefer paid cloud models when their API
keys are present, but still work with a local fallback:

```yaml
models:
  smart:
    first_available:
      - anthropic/claude-sonnet-4-5
      - openai/gpt-5
      - dmr/ai/qwen3

agents:
  root:
    model: smart
    instruction: You are a helpful assistant.
```

At load time, Docker Agent selects the first candidate whose credentials are
configured. You only need credentials for one candidate. See
[Model Configuration](../../configuration/models/index.md#first-available-models)
for details.

## Supported Providers

| Provider            | Key              | Example Models                       | API Key Env Var                     |
| ------------------- | ---------------- | ------------------------------------ | ----------------------------------- |
| OpenAI              | `openai`         | gpt-5, gpt-5-mini, gpt-4o            | `OPENAI_API_KEY`                    |
| Anthropic           | `anthropic`      | claude-sonnet-4-5, claude-opus-4-7   | `ANTHROPIC_API_KEY`                 |
| Google              | `google`         | gemini-3.5-flash, gemini-3-pro       | `GOOGLE_API_KEY` / `GEMINI_API_KEY` |
| AWS Bedrock         | `amazon-bedrock` | Claude, Nova, Llama models           | AWS credentials                     |
| Docker Model Runner | `dmr`            | ai/qwen3, ai/llama3.2                | None (local)                        |
| Mistral             | `mistral`        | Mistral models                       | `MISTRAL_API_KEY`                   |
| xAI                 | `xai`            | Grok models                          | `XAI_API_KEY`                       |
| Nebius              | `nebius`         | Open-source and specialised models   | `NEBIUS_API_KEY`                    |
| NVIDIA NIM          | `nvidia`         | Nemotron, Llama, Qwen, DeepSeek (open models) | `NVIDIA_API_KEY`               |
| MiniMax             | `minimax`        | MiniMax models                       | `MINIMAX_API_KEY`                   |
| Baseten             | `baseten`        | DeepSeek, Kimi, GLM, Llama models    | `BASETEN_API_KEY`                   |
| OVHcloud            | `ovhcloud`       | Qwen, Llama, Mistral, DeepSeek (EU-hosted) | `OVH_AI_ENDPOINTS_ACCESS_TOKEN` |
| Groq                | `groq`           | Llama, Qwen, GPT-OSS (fast inference) | `GROQ_API_KEY`                     |
| Fireworks AI        | `fireworks`      | Kimi, Llama, Qwen, DeepSeek, GLM (open models) | `FIREWORKS_API_KEY`         |
| DeepSeek            | `deepseek`       | DeepSeek-V3 chat and R1 reasoner     | `DEEPSEEK_API_KEY`                  |
| Cerebras            | `cerebras`       | GPT-OSS, GLM (fast inference)         | `CEREBRAS_API_KEY`                  |
| Together AI         | `together`       | Llama, Qwen, DeepSeek, Kimi (open models) | `TOGETHER_API_KEY`             |
| Hugging Face        | `huggingface`    | Llama, Qwen, DeepSeek, GLM (open models) | `HF_TOKEN`                      |
| Cloudflare Workers AI | `cloudflare-workers-ai` | Llama, Mistral, Qwen, Gemma (edge-hosted open models) | `CLOUDFLARE_API_TOKEN` + `CLOUDFLARE_ACCOUNT_ID` |
| Moonshot AI         | `moonshot`       | Kimi K2 chat, reasoning, and coding models | `MOONSHOT_API_KEY`             |
| Vercel AI Gateway   | `vercel`         | Multi-provider gateway               | `AI_GATEWAY_API_KEY`                |
| Cloudflare AI Gateway | `cloudflare-ai-gateway` | Multi-provider gateway         | `CLOUDFLARE_API_TOKEN` + `CLOUDFLARE_ACCOUNT_ID` + `CLOUDFLARE_GATEWAY_ID` |
| Requesty            | `requesty`       | Multi-provider gateway               | `REQUESTY_API_KEY`                  |
| OpenRouter          | `openrouter`     | Multi-provider gateway               | `OPENROUTER_API_KEY`                |
| Azure OpenAI        | `azure`          | gpt-4o, gpt-5 on Azure               | `AZURE_API_KEY` + `base_url`        |
| Ollama              | `ollama`         | Any local Ollama model               | None (local; optional `base_url`)   |
| GitHub Copilot      | `github-copilot` | Copilot-hosted OpenAI/Anthropic      | `GITHUB_TOKEN` (PAT with `copilot`) |
| ChatGPT (OpenAI account) | `chatgpt`   | gpt-5 family via ChatGPT subscription | None (sign in via `docker agent setup`) |

See the [Model Providers](../../providers/overview/index.md) section for detailed configuration guides.

## Model Properties

| Property            | Type       | Description                                       |
| ------------------- | ---------- | ------------------------------------------------- |
| `provider`          | string     | Provider identifier (required)                    |
| `model`             | string     | Model name (required)                             |
| `description`       | string     | Human-readable summary of the model's purpose     |
| `temperature`       | float      | Randomness: 0.0 (deterministic) to 1.0 (creative) |
| `max_tokens`        | int        | Maximum response length                           |
| `top_p`             | float      | Nucleus sampling: 0.0 to 1.0                      |
| `frequency_penalty` | float      | Reduce repetition: 0.0 to 2.0                     |
| `presence_penalty`  | float      | Encourage topic diversity: 0.0 to 2.0             |
| `base_url`          | string     | Custom API endpoint                               |
| `thinking_budget`   | string/int | Reasoning effort configuration                    |
| `task_budget`       | int/object | Total token budget for an agentic task (Anthropic; honored by Opus 4.7 today) |
| `provider_opts`     | object     | Provider-specific options                         |

## Reasoning / Thinking Budget

Control how much the model "thinks" before responding:

| Provider   | Format     | Values                                                              | Default                          |
| ---------- | ---------- | ------------------------------------------------------------------- | -------------------------------- |
| OpenAI     | string     | `minimal`, `low`, `medium`, `high`, `xhigh`, `max`                  | `medium` (always-reasoning models only) |
| Anthropic  | int or str | 1024–32768 tokens, or `adaptive`, `adaptive/<effort>`, effort level | off                              |
| Gemini 2.5 | int        | 0 (off), -1 (dynamic), or token count                               | -1 (dynamic)                     |
| Gemini 3   | string     | `minimal`, `low`, `medium`, `high`                                  | varies                           |
| All        | string/int | `none` or `0` clears Docker Agent's local config                    | —                                |

`none` and `0` are not universal API-level disable switches. On genuine OpenAI
gpt-5.6+ endpoints (Sol/Terra/Luna), `none` is a real `reasoning_effort` value
that Docker Agent sends as-is and the model does not reason. On older OpenAI
models, `none`/`0` only clear the local `thinking_budget` — omitting the field
has the same effect — and the model falls back to the API's own default effort
(still reasoning internally for always-reasoning models like the o-series).
Providers with a true optional-thinking switch (Gemini 2.5, Claude, local
models) are fully disabled by `none`/`0`. See the
[Thinking / Reasoning guide](../../guides/thinking/index.md#disabling-thinking)
for the full per-provider breakdown.

```yaml
models:
  deep-thinker:
    provider: anthropic
    model: claude-sonnet-4-5
    thinking_budget: 16384

  fast-responder:
    provider: openai
    model: gpt-5.6
    thinking_budget: none # real API-level disable on gpt-5.6+
```

> [!NOTE]
> **Multi-provider teams**
>
> Different agents can use different providers in the same config. See [Multi-Agent](../multi-agent/index.md) for patterns.

## Alloy Models

"Alloy models" let you use more than one model in the same conversation — Docker Agent alternates between them to leverage the strengths of each:

```yaml
agents:
  root:
    model: anthropic/claude-sonnet-4-5,openai/gpt-5
    instruction: You are a helpful assistant.
```

Read more about the alloy model concept at [xbow.com/blog/alloy-agents](https://xbow.com/blog/alloy-agents).
