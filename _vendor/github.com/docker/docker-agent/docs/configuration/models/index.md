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
                     # dmr, mistral, xai, nebius, nvidia, minimax, baseten, daoxe, ovhcloud, groq, fireworks-ai, deepseek, cerebras, togetherai, huggingface, moonshotai, vercel, cloudflare-workers-ai, cloudflare-ai-gateway, requesty, openrouter, api-route, atlascloud,
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
    parallel_tool_calls: boolean # Optional: allow parallel tool calls. Omit to use the provider/API default.
    track_usage: boolean # Optional: track token usage
    routing: [list] # Optional: rule-based model routing
    capabilities: # Optional: override attachment (input) capabilities
      image: boolean # Optional: whether the model accepts image attachments
      pdf: boolean # Optional: whether the model accepts PDF attachments
      audio: boolean # Optional: whether the model accepts audio attachments
      video: boolean # Optional: whether the model accepts video attachments
    output_capabilities: # Optional: override generative output capabilities (otherwise detected from models.dev)
      image: boolean # Optional: whether the model can generate image output
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
| `provider`            | string     | ✓/✗      | Required for regular model definitions; omitted for `first_available` selectors. Provider: `openai`, `anthropic`, `google`, `amazon-bedrock`, `dmr`, `mistral`, `xai`, `nebius`, `nvidia`, `minimax`, `baseten`, `daoxe`, `ovhcloud`, `groq`, `fireworks-ai`, `deepseek`, `cerebras`, `togetherai`, `huggingface`, `moonshotai`, `vercel`, `cloudflare-workers-ai`, `cloudflare-ai-gateway`, `requesty`, `openrouter`, `api-route`, `atlascloud`, `azure`, `ollama`, `github-copilot`, `chatgpt`, `opencode`, `opencode-go`, or any [named provider](../../providers/custom/index.md). |
| `model`               | string     | ✓/✗      | Required for regular model definitions; omitted for `first_available` selectors. Model name (e.g., `gpt-4o`, `claude-sonnet-4-5`, `gemini-3.5-flash`) |
| `description`         | string     | ✗        | Informational, human-readable summary of the model's purpose or strengths (e.g., "fast and cheap, good for summaries"). Not sent to the model. Can be combined with `first_available` (a selector's description is kept when it resolves). |
| `temperature`         | float      | ✗        | Sampling randomness. Range is provider-dependent — typically `0.0–2.0` (Anthropic caps at `1.0`). `0.0` is deterministic. |
| `max_tokens`          | int        | ✗        | Maximum response length in tokens                                                     |
| `top_p`               | float      | ✗        | Nucleus sampling threshold (`0.0–1.0`)                                                |
| `frequency_penalty`   | float      | ✗        | Penalize repeated tokens (`-2.0–2.0`)                                                 |
| `presence_penalty`    | float      | ✗        | Encourage topic diversity (`-2.0–2.0`)                                                |
| `base_url`            | string     | ✗        | Custom API endpoint URL (for self-hosted or proxied endpoints)                        |
| `token_key`           | string     | ✗        | Environment variable name containing the API token (overrides provider default)       |
| `thinking_budget`     | string/int | ✗        | Reasoning effort control. See [Thinking Budget](#thinking-budget).                    |
| `task_budget`         | int/object | ✗        | Total token budget for an agentic task (Anthropic only). See [Task Budget](#task-budget). |
| `parallel_tool_calls` | boolean    | ✗        | Allow model to call multiple tools at once. When omitted, Docker Agent leaves the setting unset so the selected provider or API can apply its own default. |
| `track_usage`         | boolean    | ✗        | Track and report token usage for this model                                           |
| `routing`             | array      | ✗        | Rule-based routing to different models. See [Model Routing](../routing/index.md). |
| `capabilities`        | object     | ✗        | Override attachment (input) capabilities for this model. See [Attachment Capability Overrides](#attachment-capability-overrides). |
| `output_capabilities` | object     | ✗        | Override generative output capabilities for this model, e.g. image generation. Omitted flags are detected from models.dev; explicit values take precedence. Cannot be combined with `first_available`. See [Output Capabilities](#output-capabilities). |
| `cost`                | object     | ✗        | Explicit token pricing in USD per 1M tokens, overriding the built-in catalogue. See [Custom Token Pricing](#custom-token-pricing). |
| `provider_opts`       | object     | ✗        | Provider-specific options (see provider pages)                                        |
| `title_model`         | string     | ✗        | Model used for session-title generation. Can be a named model from the `models:` section or an inline `provider/model` string. When omitted, the agent's primary model generates titles. Cannot be combined with `first_available`. |
| `compaction_model`    | string     | ✗        | Model used for session compaction (summary generation). Can be a named model or an inline `provider/model` string. The agent-level `compaction_model` takes precedence over this value, which in turn takes precedence over a provider-level default. When none is set, the primary model compacts. Cannot be combined with `first_available`. See the [Context & Compaction guide](../../guides/compaction/index.md). |
| `compaction_threshold` | float     | ✗        | Fraction of the context window at which proactive auto-compaction triggers for agents running this model. Must be greater than `0` and at most `1`. Takes precedence over the agent-level `compaction_threshold`. Cannot be combined with `first_available`. Default: `0.9`. See the [Context & Compaction guide](../../guides/compaction/index.md). |
| `bypass_models_gateway` | boolean  | ✗        | When `true`, this model connects directly to its provider even when a models gateway (`--models-gateway` / `DOCKER_AGENT_MODELS_GATEWAY`) is configured. Implied by a custom `base_url`. See [Gateway Bypass](#gateway-bypass). |

Built-in provider IDs use models.dev names: `fireworks-ai`, `togetherai`,
`moonshotai`, and `opencode`. Existing configurations using `fireworks`,
`together`, `moonshot`, or `opencode-zen` continue to work without warnings.
New configurations using canonical IDs require a version that supports them.
Explicitly named custom providers still take precedence.

## Attachment Capability Overrides

For custom OpenAI-compatible providers, local models (Ollama, DMR), and any
model the built-in catalogue does not describe, Docker Agent cannot
auto-detect whether the endpoint accepts image, PDF, audio, or video
attachments. When the model is absent from the catalogue, Docker Agent logs a
diagnostic and falls back to text-only, silently dropping attachments. For the
`chatgpt` provider, a missing entry first falls back to the matching `openai`
entry for **image input only**; PDF, audio, and video are not inferred.

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

  proxy-multimodal:
    provider: vision-proxy
    model: gemini-2.5-pro
    capabilities:
      image: true
      pdf: true
      audio: true
      video: true
```

| Field                  | Type    | Description                                       |
| ---------------------- | ------- | -------------------------------------------------- |
| `capabilities.image`   | boolean | Whether the model accepts image attachments       |
| `capabilities.pdf`     | boolean | Whether the model accepts PDF attachments         |
| `capabilities.audio`   | boolean | Whether the model accepts audio attachments       |
| `capabilities.video`   | boolean | Whether the model accepts video attachments       |

The flags must match what the endpoint actually accepts. Claiming a modality
that the endpoint does not support leads to a provider-side API error. When
`capabilities` is omitted, input support comes from the catalogue, the scoped
ChatGPT image fallback above, or conservative text-only fallback.

### Unsupported media is stripped before the call

Before each model call, Docker Agent removes image, audio, and video message
parts that the resolved capabilities of the active model do not cover, instead
of letting the provider fail the whole request. Adjacent text (and PDF) parts
are preserved in their original order, and each stripped part is reported in
the debug log (`--debug`) with its media kind and reason.

The stripping decision uses the same capability resolution as attachment
routing: an explicit `capabilities` declaration is authoritative, so a model
declared with `audio: true` keeps its audio parts even when the catalogue says
otherwise. Models absent from the catalogue (without an override) resolve to
the conservative text-only default and have their media parts stripped.

See [`examples/capability-overrides.yaml`](https://github.com/docker/docker-agent/blob/main/examples/capability-overrides.yaml) for a complete example, and
[`examples/strip-unsupported-media.yaml`](https://github.com/docker/docker-agent/blob/main/examples/strip-unsupported-media.yaml) for a fixture demonstrating the
stripping behaviour with and without an override.

### Output capabilities

`output_capabilities` overrides what a model can generate, as opposed to
`capabilities`, which overrides what it accepts as input. Resolution follows
one precedence chain: explicit `false`, explicit `true`, then an exact
models.dev record whose `Modalities.Output` contains `image`. An omitted image
flag (including `output_capabilities: {}`) therefore uses catalogue metadata;
an unknown model or unavailable catalogue leaves image output disabled. Docker
Agent never infers this capability from the model name.

```yaml
models:
  gemini-image:
    provider: google
    model: gemini-2.5-flash-image
    output_capabilities:
      image: true # this model is declared able to generate image output
```

| Field                       | Type    | Description                                                  |
| --------------------------- | ------- | -------------------------------------------------------------|
| `output_capabilities.image` | boolean | Whether the model is declared able to generate image output  |

Omitting `output_capabilities`, using an empty block, or omitting `image` uses
models.dev metadata for that exact model when available. Setting `image`
explicitly overrides the catalogue; an explicit `false` has highest precedence
and disables image response modalities even when the catalogue lists image
output. Enabling image output only opts the model into behavior that keys off
that capability (for example, a provider-specific image-output request
contract); it does not guarantee that a provider will return an image.

Which requests ask for image output, and which request shapes are rejected
when it is enabled, is provider-specific — see
[Google Gemini: Generated Images](../../providers/google/index.md#generated-images).
Where the returned images land and how they are rendered is covered under
[Generated Media Files](../../features/sessions/index.md#generated-media-files)
and the TUI's [Generated Media](../../features/tui/index.md#generated-media).

> [!WARNING]
> **Constraint**
>
> `output_capabilities` cannot be combined with `first_available` model selection — the combination is rejected at validation time. Declare it on the concrete candidate models instead.

See [`examples/gemini_image_output.yaml`](https://github.com/docker/docker-agent/blob/main/examples/gemini_image_output.yaml) for a complete example.

## Custom Token Pricing

Docker Agent prices each model call from the [models.dev](https://models.dev/)
catalogue, including long-context tiers. When the total prompt (fresh, cached,
and cache-written input) exceeds a tier's threshold, its rates apply to the
whole call. Thresholds are model-specific: for example, GPT-5.4 uses 272k tokens
and Gemini 2.5 Pro uses 200k. Models the catalogue does not know — custom
OpenAI-compatible providers, local models, private deployments — are "unpriced": every call is
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

When a models gateway (`--models-gateway` / `DOCKER_AGENT_MODELS_GATEWAY`) is configured,
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

Control how much reasoning the model does before responding with `thinking_budget` (string or integer). The accepted values and defaults depend on the provider and model. See the [Thinking / Reasoning guide](../../guides/thinking/index.md#quick-reference) for the comparison and how to choose an effort level.

```yaml
models:
  claude:
    provider: anthropic
    model: claude-sonnet-4-5
    max_tokens: 32768
    thinking_budget: 16384
```

### OpenAI

Use a string effort level. See [OpenAI thinking budgets](../../providers/openai/index.md#thinking-budget) for supported levels and model restrictions.

### Anthropic

Use an integer token budget or an adaptive effort setting, depending on the model. See [Anthropic thinking budgets](../../providers/anthropic/index.md#thinking-budget) for accepted values, the `max_tokens` constraint, and model restrictions.

### Google Gemini 2.5

Use an integer token budget. See [Gemini thinking budgets](../../providers/google/index.md#thinking-budget) for defaults, limits, and dynamic thinking.

### Google Gemini 3

Use a string effort level. See [Gemini thinking budgets](../../providers/google/index.md#thinking-budget) for supported values and examples.

### Disabling Thinking

```yaml
thinking_budget: none # or 0
```

`none` and `0` both clear Docker Agent's local thinking configuration (omitting `thinking_budget` has the same effect). Whether that reaches the API as a real "off" switch depends on the model — see [Disabling Thinking](../../guides/thinking/index.md#disabling-thinking) in the guide.

## Task Budget

**Anthropic-only.**

`task_budget` caps the **total** number of tokens the model may spend across a
multi-step agentic task — thinking, tool calls, and final output combined.
Docker Agent never gates it by model name; which Claude models honor or
reject the field, and how it is forwarded to the API, is covered on the
[Anthropic provider page](../../providers/anthropic/index.md#task-budget).

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

`provider_opts.interleaved_thinking` controls reasoning between tool calls on Claude models. See [Anthropic](../../providers/anthropic/index.md#interleaved-thinking) or [Bedrock](../../providers/bedrock/index.md#interleaved-thinking-claude-on-bedrock) for automatic enablement, opt-out syntax, and beta-header handling.

## Thinking Display (Anthropic)

`provider_opts.thinking_display` controls the thinking content returned in responses. See [Anthropic: Thinking Display](../../providers/anthropic/index.md#thinking-display) for accepted values, defaults, an override example, and startup validation.

## Extra Request Body

For OpenAI-compatible providers and Docker Model Runner, `provider_opts.extra_body`
adds arbitrary top-level fields to **Chat Completions** request bodies. Use it for
backend-specific options that Docker Agent does not expose directly:

```yaml
models:
  local:
    provider: openai
    model: mlx-community/Qwen3.6-35B-A3B-8bit
    base_url: http://localhost:8080/v1
    provider_opts:
      extra_body:
        chat_template_kwargs:
          enable_thinking: false
```

Choose the switch your server supports. Other backends may accept
`extra_body: {reasoning_effort: none}` or
`extra_body: {thinking: {type: disabled}}` instead. Docker Agent does not translate
these options or infer them from `thinking_budget: none`.

`extra_body` must be an object; other values cause a request error. Its contents
are merged at the top level, without an `extra_body` wrapper. Nested objects are
replaced, not deep-merged. Explicit values override generated fields, including
sampling parameters, token limits, and internal no-thinking settings. These
overrides also apply to title generation, compaction, and OpenAI-client reranking
when they use the configured model's Chat Completions endpoint. Streaming calls
always force `stream: true`, regardless of `extra_body`. Only set fields your
server accepts, and avoid overriding request structure such as `messages`.

This option does not apply to the Responses API, embeddings, or DMR's model-wide
`_configure` endpoint. Other `provider_opts` keys remain client-side unless
explicitly supported.

See [`examples/extra_body.yaml`](https://github.com/docker/docker-agent/blob/main/examples/extra_body.yaml)
for a complete local-server configuration.

## Custom HTTP Headers

For OpenAI-compatible providers (`openai`, `github-copilot`, `mistral`, `xai`,
`nebius`, `nvidia`, `minimax`, `baseten`, `daoxe`, `ovhcloud`, `groq`, `fireworks-ai`, `deepseek`, `cerebras`, `togetherai`, `huggingface`, `moonshotai`, `vercel`, `cloudflare-workers-ai`, `cloudflare-ai-gateway`, `requesty`, `openrouter`, `atlascloud`, `ollama`, and any custom provider using the OpenAI API),
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

Models can reference a named provider to inherit shared defaults. Model-level settings take precedence. See [Default Inheritance](../../providers/custom/index.md#default-inheritance) for a complete example and the special precedence rules for `compaction_model`, and [Provider Properties](../../providers/custom/index.md#provider-properties) for the inheritable fields.
