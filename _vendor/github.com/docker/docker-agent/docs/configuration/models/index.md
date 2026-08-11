---
title: "Model Configuration"
description: "Complete reference for defining models with providers, parameters, and reasoning settings."
keywords: docker agent, ai agents, configuration, yaml, model configuration
linkTitle: "Model Config"
weight: 40
canonical: https://docs.docker.com/ai/docker-agent/configuration/models/
---

_Complete reference for defining models with providers, parameters, and reasoning settings._

## Full Schema

<!-- yaml-lint:skip -->
```yaml
models:
  model_name:
    first_available: [list] # Optional: candidate model refs, tried in order by available credentials.
                            # Mutually exclusive with other model settings.
    provider: string # Required unless using first_available. One of: openai, anthropic, google, amazon-bedrock,
                     # dmr, mistral, xai, nebius, nvidia, minimax, baseten, ovhcloud, groq, fireworks, deepseek, cerebras, together, huggingface, moonshot, vercel, cloudflare-workers-ai, cloudflare-ai-gateway, requesty, openrouter,
                     # azure, ollama, github-copilot, or a named provider defined
                     # under the top-level `providers:` section.
    model: string # Required: model identifier
    description: string # Optional: human-readable summary of the model's purpose or strengths
    temperature: float # Optional: 0.0–2.0 (provider-dependent; e.g. Anthropic caps at 1.0)
    max_tokens: integer # Optional: response length limit
    top_p: float # Optional: 0.0–1.0
    frequency_penalty: float # Optional: -2.0–2.0
    presence_penalty: float # Optional: -2.0–2.0
    base_url: string # Optional: custom API endpoint
    token_key: string # Optional: env var for API token
    thinking_budget: string|int # Optional: reasoning effort
    task_budget: int|object # Optional: total task token budget (Anthropic)
    parallel_tool_calls: boolean # Optional: allow parallel tool calls
    track_usage: boolean # Optional: track token usage
    routing: [list] # Optional: rule-based model routing
    capabilities: # Optional: override attachment capabilities
      image: boolean # Optional: whether the model accepts image attachments
      pdf: boolean # Optional: whether the model accepts PDF attachments
    cost: # Optional: explicit token pricing (USD per 1M tokens)
      input: float # Optional: price per 1M input tokens
      output: float # Optional: price per 1M output tokens
      cache_read: float # Optional: price per 1M cached input tokens
      cache_write: float # Optional: price per 1M cache-write tokens
    provider_opts: # Optional: provider-specific options
      key: value
    title_model: string # Optional: model used for session-title generation
    compaction_model: string # Optional: model used for session-compaction (summary generation)
    compaction_threshold: float # Optional: context-window fraction that triggers auto-compaction (0–1, default: 0.9)
    bypass_models_gateway: boolean # Optional: skip the models gateway for this model (implied by a custom base_url)
```

## Properties Reference

| Property              | Type       | Required | Description                                                                           |
| --------------------- | ---------- | -------- | ------------------------------------------------------------------------------------- |
| `first_available`     | array      | ✗        | Candidate model references tried in order; selects the first whose credentials are configured. Mutually exclusive with other model settings. |
| `provider`            | string     | ✓/✗      | Required for regular model definitions; omitted for `first_available` selectors. Provider: `openai`, `anthropic`, `google`, `amazon-bedrock`, `dmr`, `mistral`, `xai`, `nebius`, `nvidia`, `minimax`, `baseten`, `ovhcloud`, `groq`, `fireworks`, `deepseek`, `cerebras`, `together`, `huggingface`, `moonshot`, `vercel`, `cloudflare-workers-ai`, `cloudflare-ai-gateway`, `requesty`, `openrouter`, `azure`, `ollama`, `github-copilot`, `chatgpt`, or any [named provider](../../providers/custom/index.md). |
| `model`               | string     | ✓/✗      | Required for regular model definitions; omitted for `first_available` selectors. Model name (e.g., `gpt-4o`, `claude-sonnet-4-5`, `gemini-3.5-flash`) |
| `description`         | string     | ✗        | Informational, human-readable summary of the model's purpose or strengths (e.g., "fast and cheap, good for summaries"). Not sent to the model. Can be combined with `first_available` (a selector's description is kept when it resolves). |
| `temperature`         | float      | ✗        | Sampling randomness. Range is provider-dependent — typically `0.0–2.0` (Anthropic caps at `1.0`). `0.0` is deterministic. |
| `max_tokens`          | int        | ✗        | Maximum response length in tokens                                                     |
| `top_p`               | float      | ✗        | Nucleus sampling threshold (`0.0–1.0`)                                                |
| `frequency_penalty`   | float      | ✗        | Penalize repeated tokens (`-2.0–2.0`)                                                 |
| `presence_penalty`    | float      | ✗        | Encourage topic diversity (`-2.0–2.0`)                                                |
| `base_url`            | string     | ✗        | Custom API endpoint URL (for self-hosted or proxied endpoints)                        |
| `token_key`           | string     | ✗        | Environment variable name containing the API token (overrides provider default)       |
| `thinking_budget`     | string/int | ✗        | Reasoning effort control                                                              |
| `task_budget`         | int/object | ✗        | Total token budget for an agentic task (forwarded to Anthropic; see [Task Budget](#task-budget)). |
| `parallel_tool_calls` | boolean    | ✗        | Allow model to call multiple tools at once                                            |
| `track_usage`         | boolean    | ✗        | Track and report token usage for this model                                           |
| `routing`             | array      | ✗        | Rule-based routing to different models. See [Model Routing](../routing/index.md). |
| `capabilities`        | object     | ✗        | Override attachment capabilities for this model. See [Attachment Capability Overrides](#attachment-capability-overrides). |
| `cost`                | object     | ✗        | Explicit token pricing in USD per 1M tokens, overriding the built-in catalogue. See [Custom Token Pricing](#custom-token-pricing). |
| `provider_opts`       | object     | ✗        | Provider-specific options (see provider pages)                                        |
| `title_model`         | string     | ✗        | Model used for session-title generation. Can be a named model from the `models:` section or an inline `provider/model` string. When omitted, the agent's primary model generates titles. Cannot be combined with `first_available`. |
| `compaction_model`    | string     | ✗        | Model used for session compaction (summary generation). Can be a named model or an inline `provider/model` string. The agent-level `compaction_model` takes precedence over this value, which in turn takes precedence over a provider-level default. When none is set, the primary model compacts. Cannot be combined with `first_available`. See the [Context & Compaction guide](../../guides/compaction/index.md). |
| `compaction_threshold` | float     | ✗        | Fraction of the context window at which proactive auto-compaction triggers for agents running this model. Must be greater than `0` and at most `1`. Takes precedence over the agent-level `compaction_threshold`. Cannot be combined with `first_available`. Default: `0.9`. See the [Context & Compaction guide](../../guides/compaction/index.md). |
| `bypass_models_gateway` | boolean  | ✗        | When `true`, this model connects directly to its provider even when a models gateway (`--models-gateway` / `CAGENT_MODELS_GATEWAY`) is configured. Implied by a custom `base_url`. See [Gateway Bypass](#gateway-bypass). |

## Attachment Capability Overrides

For custom OpenAI-compatible providers, local models (Ollama, DMR), and any
model the built-in catalogue does not describe, Docker Agent cannot
auto-detect whether the endpoint accepts image or PDF attachments. When the
model is absent from the catalogue, Docker Agent logs a diagnostic and falls
back to text-only, silently dropping attachments.

Declare `capabilities` to make the model's attachment support authoritative
and skip the catalogue lookup entirely:

```yaml
models:
  llava-local:
    provider: ollama
    model: llava
    capabilities:
      image: true   # accepts image attachments
      pdf: false    # does not accept PDFs

  proxy-vision:
    provider: vision-proxy
    model: gpt-4o
    capabilities:
      image: true
      pdf: true
```

| Field                  | Type    | Description                                       |
| ---------------------- | ------- | ------------------------------------------------- |
| `capabilities.image`   | boolean | Whether the model accepts image attachments       |
| `capabilities.pdf`     | boolean | Whether the model accepts PDF attachments         |

The flags must match what the endpoint actually accepts. Claiming a modality
that the endpoint does not support leads to a provider-side API error. When
`capabilities` is omitted the behaviour is unchanged (catalogue lookup then
conservative text-only fallback).

See [`examples/capability-overrides.yaml`](https://github.com/docker/docker-agent/blob/main/examples/capability-overrides.yaml) for a complete example.

## Custom Token Pricing

Docker Agent prices each model call from the [models.dev](https://models.dev/)
catalogue. Models the catalogue does not know — custom OpenAI-compatible
providers, local models, private deployments — are "unpriced": every call is
recorded at $0 despite consuming tokens, with only a log warning.

Declare `cost` to price a model explicitly, in **USD per one million tokens**.
When set, it takes precedence over the catalogue and makes an uncatalogued
model priced:

```yaml
models:
  internal-gpt:
    provider: internal-llm
    model: gpt-4o
    cost:
      input: 1.25 # USD per 1M input tokens
      output: 5.00 # USD per 1M output tokens
      cache_read: 0.125 # USD per 1M cached input tokens
      cache_write: 1.5625 # USD per 1M cache-write tokens

  # Also works for catalogued models, e.g. a negotiated enterprise discount:
  discounted-sonnet:
    provider: anthropic
    model: claude-sonnet-4-5
    cost:
      input: 2.4
      output: 12.0
```

| Field              | Type  | Description                             |
| ------------------ | ----- | --------------------------------------- |
| `cost.input`       | float | USD price per 1M input tokens           |
| `cost.output`      | float | USD price per 1M output tokens          |
| `cost.cache_read`  | float | USD price per 1M cached input tokens    |
| `cost.cache_write` | float | USD price per 1M cache-write tokens     |

The declared prices feed per-turn cost computation, session cost tracking, the
`/model` picker, and the [`after_llm_call` hook](../hooks/index.md)'s `cost`
field. Prices must not be negative; omitted fields default to `0`. An all-zero
table means "priced, free" — distinct from omitting `cost` entirely
(unpriced). Cannot be combined with `first_available` (set it on the candidate
models instead).

See [`examples/custom-pricing.yaml`](https://github.com/docker/docker-agent/blob/main/examples/custom-pricing.yaml) for a complete example.

## Delegating Session-Title Generation

The `title_model` field lets a heavyweight primary model hand off the cheap
title-generation call to a smaller, faster model:

```yaml
model: anthropic/claude-opus-4-7
title_model: anthropic/claude-haiku-4-5
```

The value can be a named entry from the `models` stanza or an inline
`provider/model` string. When omitted, the agent's primary model generates
titles.

> [!WARNING]
> **Constraint**
>
> `title_model` cannot be combined with `first_available` model selection — the combination is rejected at validation time.

## Delegating Session Compaction

> [!TIP]
> **Full guide**
>
> For a task-oriented walkthrough of automatic vs. on-demand compaction, trimming tool results, and reading the context gauge, see [Managing Context & Compaction](../../guides/compaction/index.md). This section covers the `compaction_model` and `compaction_threshold` fields themselves.

The `compaction_model` field lets a heavyweight primary model hand off the expensive
compaction (summary generation) call to a smaller, faster model:

```yaml
models:
  primary:
    provider: anthropic
    model: claude-sonnet-4-5
    compaction_model: fast
  fast:
    provider: anthropic
    model: claude-haiku-4-5
```

The value can be a named entry from the `models` stanza or an inline
`provider/model` string. Resolution priority: an agent-level `compaction_model`
wins, then the model-level value, then a provider-level default set in the
`providers` section; when none is set, the primary model compacts. For an
agent listing several models (`model: a,b`), the first listed model that sets
a value (or whose provider sets a default) wins at that level.

```yaml
providers:
  my_anthropic:
    provider: anthropic
    # Default for every agent whose model uses this provider.
    compaction_model: anthropic/claude-haiku-4-5
```

If the compaction model has a **smaller context window** than the primary,
Docker Agent triggers compaction against the smaller window so the summary
call can always ingest the full conversation. Pair the primary with a
compaction model whose window is at least as large to keep the proactive
trigger aligned with the primary's window.

By default the proactive trigger fires when the estimated token usage crosses
**90%** of the context window. The `compaction_threshold` field tunes that
fraction (greater than `0`, at most `1`): lower values compact earlier and
keep requests smaller, higher values compact later and keep more verbatim
history. It can be set on the model (as above, taking precedence) or on the
agent, and automatic compaction can be disabled entirely per agent with
`session_compaction: false` — see [Agent Config](../agents/index.md#properties-reference).

```yaml
models:
  primary:
    provider: anthropic
    model: claude-sonnet-4-5
    compaction_model: fast
    # Compact at 80% of the window instead of the default 90%.
    compaction_threshold: 0.8
```

> [!WARNING]
> **Constraint**
>
> `compaction_model` cannot be combined with `first_available` model selection — the combination is rejected at validation time.

See [`examples/compaction_model.yaml`](https://github.com/docker/docker-agent/blob/main/examples/compaction_model.yaml)
and [`examples/compaction_threshold.yaml`](https://github.com/docker/docker-agent/blob/main/examples/compaction_threshold.yaml)
for complete examples.

## Gateway Bypass

When a models gateway (`--models-gateway` / `CAGENT_MODELS_GATEWAY`) is configured,
models without a custom `base_url` route through it by default. Set
`bypass_models_gateway: true` on a specific model to make it connect directly
to its provider instead:

```yaml
models:
  gateway-model:
    provider: openai
    model: gpt-5

  direct-model:
    provider: anthropic
    model: claude-sonnet-4-5
    bypass_models_gateway: true  # uses ANTHROPIC_API_KEY directly
```

The bypassed model authenticates with the provider's own credentials
(`OPENAI_API_KEY`, `ANTHROPIC_API_KEY`, `token_key`, etc.) rather than the
gateway's short-lived token. The rest of the agent's models continue routing
through the gateway as before.

Bypass is propagated transparently through router models: a bypass-flagged routing
model passes the flag to all of its routed targets automatically.

> [!WARNING]
> **Security note**
>
> On an untrusted config, a malicious `base_url` combined with `bypass_models_gateway: true` could route provider credentials to an attacker-controlled endpoint. Only enable this on configs you control.

> [!WARNING]
> **Constraint**
>
> `bypass_models_gateway: true` cannot be combined with `first_available` — the combination is rejected at validation time.

See [`examples/bypass_models_gateway.yaml`](https://github.com/docker/docker-agent/blob/main/examples/bypass_models_gateway.yaml) for a complete example.

## First Available Models

Use `first_available` when the same agent should work with whichever provider credentials are available in the current environment. Docker Agent checks the candidates in order at load time and replaces the selector with the first candidate whose required environment variables are configured.

```yaml
models:
  smart:
    first_available:
      - anthropic/claude-sonnet-4-6
      - openai/gpt-5
      - google/gemini-3.5-flash
      - dmr/ai/qwen3 # local fallback; no API key required

agents:
  root:
    model: smart
    instruction: You are a helpful assistant.
```

Candidates can be inline `provider/model` references or names from the same `models:` section. Local providers such as `dmr` and `ollama` do not require credentials, so they are useful as final fallbacks.

If none of the candidates has credentials configured, Docker Agent reports the missing environment variables grouped by candidate. You only need to configure one group of credentials, not every provider in the list.

A `first_available` model is only a selector. Except for the informational `description`, it cannot be combined with `provider`, `model`, `routing`, `token_key`, budgets, sampling options, or other model settings. Put those settings on named candidate models instead:

```yaml
models:
  claude:
    provider: anthropic
    model: claude-sonnet-4-6
    max_tokens: 64000

  gpt:
    provider: openai
    model: gpt-5
    thinking_budget: low

  smart:
    first_available:
      - claude
      - gpt
      - dmr/ai/qwen3
```

See [`examples/first_available.yaml`](https://github.com/docker/docker-agent/blob/main/examples/first_available.yaml) for a complete example.

## Thinking Budget

Control how much reasoning the model does before responding. This varies by provider:

### OpenAI

Uses effort levels as strings:

```yaml
models:
  gpt:
    provider: openai
    model: gpt-5.6
    thinking_budget: low # none | minimal | low | medium | high | xhigh | max (xhigh needs gpt-5.2+; none/max need gpt-5.6+; minimal dropped on gpt-5.6+)
```

### Anthropic

Uses an integer token budget (1024–32768), or — on adaptive-capable models (Opus 4.6+) — `adaptive`, `adaptive/<effort>`, or a bare effort level:

```yaml
models:
  claude:
    provider: anthropic
    model: claude-sonnet-4-5
    thinking_budget: 16384 # must be < max_tokens

  opus:
    provider: anthropic
    model: claude-opus-4-6
    thinking_budget: adaptive # or adaptive/<effort>, or low | medium | high | xhigh | max
```

### Google Gemini 2.5

Uses an integer token budget. `0` disables, `-1` lets the model decide:

```yaml
models:
  gemini:
    provider: google
    model: gemini-2.5-flash
    thinking_budget: -1 # dynamic (default)
```

### Google Gemini 3

Uses effort levels like OpenAI:

```yaml
models:
  gemini3:
    provider: google
    model: gemini-3-flash
    thinking_budget: medium # minimal | low | medium | high
```

### Disabling Thinking

```yaml
thinking_budget: none # or 0
```

`none` and `0` both clear Docker Agent's local thinking configuration (omitting `thinking_budget` has the same effect); neither is guaranteed to reach the API as a real "off" switch:

- **OpenAI gpt-5.6+** (Sol/Terra/Luna) is the only case with a genuine API-level `none` reasoning effort: Docker Agent sends it as-is and the model does not reason.
- **Older OpenAI reasoning models** (o-series, gpt-5 through gpt-5.5) have no such switch: `none`/`0` just clear the local config, and the model falls back to the API's own default effort and still reasons internally. Same for other always-reasoning models (Gemini 3).
- Providers with a true optional-thinking switch (Gemini 2.5, Claude, local models) are fully disabled by `none`/`0`.

```yaml
models:
  fast-responder:
    provider: openai
    model: gpt-5.6
    thinking_budget: none # real API-level disable on gpt-5.6+
```

See the [Thinking / Reasoning guide](../../guides/thinking/index.md) for per-provider details, including AWS Bedrock and Docker Model Runner.

## Task Budget

**Anthropic-only.**

`task_budget` caps the **total** number of tokens the model may spend across a
multi-step agentic task — combining thinking, tool calls, and final output
tokens. It lets long-running agents self-regulate effort without having to
choose a tight per-call `max_tokens`.

It is forwarded to Anthropic's
[`output_config.task_budget`](https://platform.claude.com/docs/en/about-claude/models/whats-new-claude-4-7)
request field. Docker Agent automatically attaches the required
`task-budgets-2026-03-13` beta header whenever this field is set.

You can configure `task_budget` on **any** Claude model — Docker Agent never
gates it by model name. At the time of writing only **Claude Opus 4.7**
actually honors the field; other Claude models will reject requests that
include it. Check the Anthropic release notes linked above for the current
list of supported models.

### Integer shorthand

```yaml
models:
  opus:
    provider: anthropic
    model: claude-opus-4-7
    task_budget: 128000 # total tokens for the whole task
    thinking_budget: adaptive # works nicely together
```

### Object form

Equivalent, and forward-compatible with future budget types:

```yaml
models:
  opus:
    provider: anthropic
    model: claude-opus-4-7
    task_budget:
      type: tokens # only "tokens" is supported today
      total: 128000
```

Setting `task_budget: 0` (or omitting the field) disables the feature — the
model falls back to the provider's default behavior.

Like other inheritable model settings, `task_budget` can also be declared on a
[provider definition](../../providers/custom/index.md) and is
inherited by every model that references that provider.

See [`examples/task_budget.yaml`](https://github.com/docker/docker-agent/blob/main/examples/task_budget.yaml) for a complete example.

## Interleaved Thinking

For Anthropic and Bedrock Claude models, interleaved thinking allows tool calls during model reasoning. It is auto-enabled whenever a thinking budget is configured:

```yaml
models:
  claude:
    provider: anthropic
    model: claude-sonnet-4-5
    thinking_budget: 8192
    # interleaved_thinking is auto-enabled when thinking_budget is set
    provider_opts:
      interleaved_thinking: false # disable if needed
```

## Thinking Display (Anthropic)

For Anthropic Claude models, `thinking_display` controls whether thinking blocks are returned in responses when thinking is enabled. Newer Claude models (Opus 4.7+, Fable 5) hide thinking content by default (`omitted`); Docker Agent requests `summarized` thinking by default for adaptive/effort-based budgets so reasoning stays visible. Set this provider option to override:

```yaml
models:
  opus-4-7:
    provider: anthropic
    model: claude-opus-4-7
    thinking_budget: adaptive
    provider_opts:
      thinking_display: omitted # "summarized" or "omitted" ("display" on pre-4.6 models only)
```

`display` (full thinking blocks) is only accepted by pre-4.6 token-thinking models (e.g. Sonnet 4.5, Haiku 4.5); newer models (Opus/Sonnet 4.6+, Sonnet 5, Fable 5) only accept `summarized` and `omitted`, and Docker Agent rejects the configuration at startup.

See the [Anthropic provider page](../../providers/anthropic/index.md#thinking-display) for details.

## Custom HTTP Headers

For OpenAI-compatible providers (`openai`, `github-copilot`, `mistral`, `xai`,
`nebius`, `nvidia`, `minimax`, `baseten`, `ovhcloud`, `groq`, `fireworks`, `deepseek`, `cerebras`, `together`, `huggingface`, `moonshot`, `vercel`, `cloudflare-workers-ai`, `cloudflare-ai-gateway`, `requesty`, `openrouter`, `ollama`, and any custom provider using the OpenAI API),
`provider_opts.http_headers` adds arbitrary HTTP headers to every outgoing
request:

```yaml
models:
  my_model:
    provider: openai
    model: gpt-4o
    provider_opts:
      http_headers:
        X-Request-Source: docker-agent
        X-Tenant-Id: my-team
```

Header names are matched case-insensitively. The `github-copilot` provider
automatically sets `Copilot-Integration-Id: copilot-developer-cli` — see the
[GitHub Copilot provider page](../../providers/github-copilot/index.md)
for details.

## Examples by Provider

```yaml
models:
  # OpenAI
  gpt:
    provider: openai
    model: gpt-5

  # Anthropic
  claude:
    provider: anthropic
    model: claude-sonnet-4-5
    max_tokens: 64000

  # Google Gemini
  gemini:
    provider: google
    model: gemini-3.5-flash
    temperature: 0.5

  # AWS Bedrock
  bedrock:
    provider: amazon-bedrock
    model: global.anthropic.claude-sonnet-4-5-20250929-v1:0
    provider_opts:
      region: us-east-1

  # OpenRouter
  openrouter:
    provider: openrouter
    model: meta-llama/llama-3.3-70b-instruct

  # Docker Model Runner (local)
  local:
    provider: dmr
    model: ai/qwen3
    max_tokens: 8192
```

For detailed provider setup, see the [Model Providers](../../providers/overview/index.md) section.

## Custom Endpoints

Use `base_url` to point to custom or self-hosted endpoints:

```yaml
models:
  # Azure OpenAI
  azure_gpt:
    provider: openai
    model: gpt-4o
    base_url: https://my-resource.openai.azure.com/openai/deployments/gpt-4o
    token_key: AZURE_OPENAI_API_KEY

  # Self-hosted vLLM
  local_llama:
    provider: openai # vLLM is OpenAI-compatible
    model: meta-llama/Llama-3.2-3B-Instruct
    base_url: http://localhost:8000/v1

  # Proxy or gateway
  proxied:
    provider: openai
    model: gpt-4o
    base_url: https://proxy.internal.company.com/openai/v1
    token_key: INTERNAL_API_KEY
```

The `model` and `base_url` fields accept `${env.VAR}` (or `${VAR}`) references, which are substituted from the environment when the model is loaded. This keeps the model id or endpoint out of the config when it is supplied by the environment, e.g. a Docker Compose / DMR setup:

```yaml
models:
  nemotron3:
    provider: dmr
    model: "${env.NEMOTRON3_MODEL}"
    base_url: "${env.DMR_BASE_URL}"
```

See [Variable Expansion in Config Fields](../overview/index.md#variable-expansion-in-config-fields) for the full set of fields and supported syntaxes.

See [Local Models](../../providers/local/index.md) for more examples of custom endpoints.

## Inheriting from Provider Definitions

Models can reference a named provider to inherit shared defaults. Model-level settings always take precedence:

```yaml
providers:
  my_anthropic:
    provider: anthropic
    token_key: MY_ANTHROPIC_KEY
    max_tokens: 16384
    thinking_budget: 8192
    temperature: 0.5

models:
  claude:
    provider: my_anthropic
    model: claude-sonnet-4-5
    # Inherits max_tokens, thinking_budget, temperature from provider

  claude_fast:
    provider: my_anthropic
    model: claude-haiku-4-5
    thinking_budget: 1024  # Overrides provider default
```

See [Provider Definitions](../../providers/custom/index.md) for the full list of inheritable properties.
