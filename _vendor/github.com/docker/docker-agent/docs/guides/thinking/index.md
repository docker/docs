---
title: "Thinking / Reasoning"
description: "Control how much a model reasons before responding. Works across OpenAI, Anthropic, Google Gemini, AWS Bedrock, and Docker Model Runner."
keywords: docker agent, ai agents, guides, thinking / reasoning
weight: 20
canonical: https://docs.docker.com/ai/docker-agent/guides/thinking/
---

_Control how much a model reasons before responding. Works across OpenAI, Anthropic, Google Gemini, AWS Bedrock, and Docker Model Runner._

## What Is Thinking?

Several modern models support an extended reasoning phase that happens before they produce visible output. During this phase the model plans, evaluates options, and works through the problem — internally, not shown in the response by default. This typically improves accuracy on complex tasks like coding, math, and multi-step planning, at the cost of higher token usage and latency.

docker-agent exposes this through a single `thinking_budget` field on any named model. The value format differs slightly by provider, but the semantics are the same: higher effort means more thorough reasoning.

> [!NOTE]
> **Think tool vs. thinking budget**
>
> The [think tool](../../tools/think/index.md) is a scratchpad for models that lack native reasoning. If your model supports `thinking_budget`, you do not need the think tool.

## Quick Reference

| Provider            | Format     | Values                                                                                  | Default            |
| ------------------- | ---------- | --------------------------------------------------------------------------------------- | ------------------ |
| OpenAI              | string     | `minimal`, `low`, `medium`, `high`, `none`; `xhigh` on gpt-5.2+ only                    | `medium` (API default) |
| Anthropic           | int or str | 1024–32768 tokens, or `minimal`–`max`, `adaptive`, `adaptive/<effort>`, `none`          | off                |
| Gemini 2.5          | int        | `0` (off), `-1` (dynamic), or token count (max 24576 / 32768)                           | `-1` (dynamic)     |
| Gemini 3            | string     | `minimal`, `low`, `medium`, `high`                                                      | API default (model-dependent) |
| AWS Bedrock         | int or str | 1024–32768 tokens (`minimal`–`max` mapped to tokens); `adaptive`, `adaptive/<effort>` for Opus 4.6+ (rejected by older Claude models) | off                |
| Docker Model Runner | int or str | token count, `minimal`–`max` (mapped to tokens), `adaptive` (unlimited), `none`         | engine default     |

String values are case-insensitive. The full set of accepted strings is `none`, `minimal`, `low`, `medium`, `high`, `xhigh`, `max`, `adaptive`, and `adaptive/<effort>` — but each provider only honors the subset listed above. Unsupported values either fail at request time (OpenAI) or are mapped/ignored as described per provider below.

> `thinking_budget` is only applied by the providers listed above. Other OpenAI-compatible providers (xAI, Mistral, Ollama, …) currently ignore it — see [xAI and Mistral](#xai-grok-and-mistral).

## OpenAI

OpenAI reasoning models (o-series, gpt-5, gpt-5-mini) use a string effort level that maps to their `reasoning_effort` API parameter. The `xhigh` level is only accepted by gpt-5.2 and later minor versions; o-series and earlier gpt-5 releases top out at `high`.

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

These effort levels (`minimal`–`xhigh`) are the **only** values accepted for OpenAI. Token counts, `max`, `adaptive`, and `adaptive/<effort>` are rejected with a configuration error at request time. The `xhigh` level is only supported by gpt-5.2 and later minor versions (e.g. gpt-5.2, gpt-5.4-mini); o-series and earlier gpt-5 releases top out at `high`. Older models (o1, o3-mini) only accept `low`/`medium`/`high` — sending an unsupported level returns an API error.

> [!WARNING]
> **Tokens and max_tokens**
>
> OpenAI reasoning models always reason internally — even with `thinking_budget: none` there are hidden reasoning tokens that count against `max_tokens`. docker-agent automatically raises the output-token floor for its internal low-effort calls (e.g. title generation) so hidden reasoning cannot starve visible text output.

## Anthropic

Anthropic Claude supports two thinking modes: a **token budget** (older models) and **adaptive / effort-based** thinking (newer models).

### Token budget (Claude 4 and earlier)

Set an explicit number of thinking tokens (1024–32768). This must be less than `max_tokens`:

```yaml
models:
  claude-thinker:
    provider: anthropic
    model: claude-sonnet-4-5
    thinking_budget: 16384   # tokens reserved for internal reasoning
```

docker-agent auto-adjusts `max_tokens` when you set a thinking budget but leave `max_tokens` at its default. If you set `max_tokens` explicitly, it must be greater than `thinking_budget`.

### Adaptive thinking (Opus 4.6+ and Sonnet 4.6)

Newer Claude models support adaptive thinking, where the model decides how much to think. **Claude Opus 4.6, 4.7, 4.8, and Sonnet 4.6 only support adaptive thinking** — they reject token-based budgets. Use `adaptive`, `adaptive/<effort>`, or a bare effort level — on Anthropic, a bare effort level like `high` is shorthand for adaptive thinking at that effort:

```yaml
models:
  claude-adaptive:
    provider: anthropic
    model: claude-opus-4-6
    thinking_budget: adaptive          # model decides effort (defaults to high)

  claude-adaptive-low:
    provider: anthropic
    model: claude-opus-4-6
    thinking_budget: low               # same as adaptive/low

  claude-adaptive-max:
    provider: anthropic
    model: claude-opus-4-6
    thinking_budget: adaptive/max      # adaptive/low | adaptive/medium | adaptive/high | adaptive/xhigh | adaptive/max
```

**Adaptive effort levels and per-model support:**

| Level     | Opus 4.5 | Sonnet 4.5 / Haiku | Sonnet 4.6 | Opus 4.6 | Opus 4.7 / 4.8 | Fable 5 | Mythos 5 | Mythos preview |
| --------- | :------: | :----------------: | :--------: | :------: | :------------: | :-----: | :------: | :------------: |
| `low`     | ✓        | ✓                  | ✓          | ✓        | ✓              | ✓       | ✓        | ✓              |
| `medium`  | ✓        | ✓                  | ✓          | ✓        | ✓              | ✓       | ✓        | ✓              |
| `high`    | ✓        | ✓                  | ✓          | ✓        | ✓              | ✓       | ✓        | ✓              |
| `xhigh`   | —        | —                  | —          | —        | ✓              | ✓       | ✓        | —              |
| `max`     | —        | —                  | ✓          | ✓        | ✓              | ✓       | ✓        | ✓              |

`minimal` is treated as `low` (bare form only). `high` is the default when `adaptive` is used without an effort level.

> [!WARNING]
> **Effort strings require adaptive-capable models**
>
> Every string effort value on Anthropic is sent as adaptive thinking (`output_config.effort`), which only newer Claude models (Opus 4.6+, Sonnet 4.6) accept. For older models like Sonnet 4.5, use an integer token budget instead. Conversely, models that _only_ support adaptive thinking (Opus 4.6, 4.7, 4.8, Sonnet 4.6) automatically have token budgets coerced to `adaptive` (a warning is logged).

### Disabling thinking

```yaml
thinking_budget: none   # or 0
```

### Interleaved thinking

Interleaved thinking lets the model reason between tool calls — useful for complex agentic tasks. docker-agent auto-enables it whenever a thinking budget is configured on a Claude model, so you only need to set it explicitly to turn it off:

```yaml
models:
  claude-interleaved:
    provider: anthropic
    model: claude-sonnet-4-5
    thinking_budget: 16384
    # interleaved_thinking is auto-enabled; disable it explicitly if needed:
    provider_opts:
      interleaved_thinking: false
```

> [!NOTE]
> **Temperature and top_p**
>
> When extended thinking is enabled, Anthropic requires `temperature=1.0`. docker-agent automatically suppresses any `temperature` or `top_p` settings you have configured — they are silently ignored while thinking is active.

### Thinking display

Newer Claude models (Opus 4.7+, Fable 5) hide thinking content by default at the API level. To keep reasoning visible, docker-agent requests `summarized` thinking whenever adaptive/effort-based thinking is used without an explicit `thinking_display`. Use `thinking_display` in `provider_opts` to override:

```yaml
models:
  opus-47:
    provider: anthropic
    model: claude-opus-4-7
    thinking_budget: adaptive
    provider_opts:
      thinking_display: omitted   # summarized | display | omitted
```

| Value        | Behavior                                                                              |
| ------------ | ------------------------------------------------------------------------------------- |
| `summarized` | Thinking blocks returned with a text summary (docker-agent default for adaptive thinking). |
| `display`    | Full thinking blocks returned for display.                                            |
| `omitted`    | Thinking blocks hidden — only the signature is returned.                               |

Full thinking tokens are billed regardless of `thinking_display`.

### Task budget (Anthropic)

`task_budget` caps total tokens across an entire multi-step agentic task (thinking + tool calls + output combined):

```yaml
models:
  opus-bounded:
    provider: anthropic
    model: claude-opus-4-7
    thinking_budget: adaptive
    task_budget: 128000   # total token ceiling for the whole task
```

See the [Anthropic provider page](../../providers/anthropic/index.md#task-budget) for details.

## Google Gemini

Gemini 2.5 and Gemini 3 use different formats.

### Gemini 2.5 (token budget)

```yaml
models:
  gemini-off:
    provider: google
    model: gemini-2.5-flash
    thinking_budget: 0      # disable thinking

  gemini-dynamic:
    provider: google
    model: gemini-2.5-flash
    thinking_budget: -1     # let the model decide (default)

  gemini-fixed:
    provider: google
    model: gemini-2.5-flash
    thinking_budget: 8192   # fixed token budget (max 24576 for Flash, 32768 for Pro)
```

### Gemini 3 (level-based)

```yaml
models:
  gemini3-flash:
    provider: google
    model: gemini-3-flash
    thinking_budget: medium   # minimal | low | medium | high

  gemini3-pro:
    provider: google
    model: gemini-3-pro
    thinking_budget: high     # low | high (Pro supports fewer levels)
```

## AWS Bedrock (Claude)

Bedrock Claude uses a token budget like Anthropic. String effort levels (`minimal`–`max`) are mapped automatically:

| Effort level | Token budget |
| ------------ | ------------ |
| `minimal`    | 1,024        |
| `low`        | 2,048        |
| `medium`     | 8,192        |
| `high`       | 16,384       |
| `xhigh`/`max`| 32,768       |

```yaml
models:
  bedrock-claude-thinker:
    provider: amazon-bedrock
    model: global.anthropic.claude-sonnet-4-5-20250929-v1:0
    thinking_budget: 8192   # or use an effort level: medium
    provider_opts:
      region: us-east-1

  bedrock-claude-interleaved:
    provider: amazon-bedrock
    model: global.anthropic.claude-sonnet-4-5-20250929-v1:0
    thinking_budget: high
    provider_opts:
      region: us-east-1
      # interleaved_thinking is auto-enabled when thinking_budget is set
```

**Claude Opus 4.6+ on Bedrock requires adaptive thinking** — these models reject `thinking.type=enabled` (token budgets). Configure them with `adaptive` or `adaptive/<effort>`; docker-agent auto-coerces token budgets and effort levels on these models with a warning:

```yaml
models:
  bedrock-opus-adaptive:
    provider: amazon-bedrock
    model: global.anthropic.claude-opus-4-8
    thinking_budget: adaptive/high
    provider_opts:
      region: us-east-1
```

> [!WARNING]
> **Bedrock thinking requirements**
>
> Bedrock Claude requires token-based `thinking_budget` values to be ≥ 1024 and less than `max_tokens`. docker-agent logs a warning and ignores the budget if either condition is violated. Interleaved thinking requires the `interleaved-thinking-2025-05-14` beta header, which docker-agent adds automatically; it is auto-enabled whenever a token thinking budget is set on a Bedrock-hosted Claude model (adaptive thinking interleaves on its own).

## Docker Model Runner (local models)

For local models, `thinking_budget` is forwarded to the inference engine. Both token counts and effort strings work; effort strings map to the same token scale as Bedrock (`minimal`=1024 … `xhigh`/`max`=32768), and `adaptive` means unlimited:

```yaml
models:
  local:
    provider: dmr
    model: ai/qwen3
    thinking_budget: medium   # llama.cpp: reasoning-budget=8192; vLLM: thinking_token_budget=8192
```

- **llama.cpp**: sent as `reasoning-budget` at model-configure time.
- **vLLM**: sent as `thinking_token_budget` on each request.
- **MLX / SGLang**: no reasoning-budget knob; the value is silently ignored.

See the [Docker Model Runner provider page](../../providers/dmr/index.md) for details.

## xAI (Grok) and Mistral

xAI and Mistral run through docker-agent's OpenAI-compatible client, but the `reasoning_effort` parameter is only sent for OpenAI reasoning model names (o-series, gpt-5). **Setting `thinking_budget` on Grok or Mistral models currently has no effect** — the value is accepted by config validation but never sent to the API.

Grok and Mistral reasoning models (e.g. `grok-3-mini`, `magistral`) manage reasoning on their own; for non-reasoning models, consider the [think tool](../../tools/think/index.md) instead.

## Disabling Thinking

Use `none` or `0` to disable thinking on any provider:

```yaml
models:
  fast-model:
    provider: openai
    model: gpt-5-mini
    thinking_budget: none

  gemini-no-think:
    provider: google
    model: gemini-2.5-flash
    thinking_budget: 0
```

`none` and `0` clear docker-agent's thinking configuration — no thinking parameter is sent. Models that always reason (OpenAI o-series, gpt-5, Gemini 3) then fall back to the API's default behavior and still reason internally; only models with optional thinking (Gemini 2.5, Claude, local models) are fully disabled.

## Choosing an Effort Level

| Task complexity                  | Recommended level       |
| -------------------------------- | ----------------------- |
| Simple factual Q&A               | `none` / `minimal`      |
| General-purpose chat             | `low` / `medium`        |
| Coding, debugging, analysis      | `medium` / `high`       |
| Complex reasoning, planning      | `high` / `xhigh`        |
| Research, difficult math/logic   | `xhigh` / `max`         |
| Long agentic tasks (Anthropic)   | `adaptive`              |

## Changing Thinking Level at Runtime

While running in the TUI, press **Shift+Tab** to cycle the thinking effort level for the current model without editing your YAML config, or type `/effort <level>` to jump straight to a specific level (e.g. `/effort high`). Running `/effort` without an argument opens a picker listing the levels the current model supports:

- The level steps through the model's supported range (model-specific), wrapping around — for example `none → minimal → low → medium → high → none` on OpenAI gpt-5/o-series, `none → minimal → low → medium → high → xhigh → none` on gpt-5.2+, `none → low → medium → high → max → none` on Anthropic Opus 4.6 and Sonnet 4.6, and `none → low → medium → high → xhigh → max → none` on Anthropic Opus 4.7+, Fable 5, and Mythos 5. For older Anthropic models (e.g. Sonnet 4.5) that only accept token budgets, effort-string cycling has no effect — use an integer `thinking_budget` in your YAML config instead.
- The current level is shown in the sidebar next to the model name (e.g. `openai/gpt-5 • high`).
- This applies as a session override — it is **not** saved to the config file. The next session starts from the level defined in your YAML.
- For models that don't support reasoning, and for remote runtimes, Shift+Tab is a no-op and an informational message is displayed.
- `/effort` only accepts levels the current model supports; requesting an unsupported level shows the model's supported list. Like Shift+Tab, it is unavailable for non-reasoning models and remote runtimes.

## Sharing Thinking Config Across Models

Define a provider with a default `thinking_budget` and all models that reference it inherit it:

```yaml
providers:
  deep-anthropic:
    provider: anthropic
    thinking_budget: adaptive/high
    max_tokens: 32768

models:
  claude-smart:
    provider: deep-anthropic
    model: claude-opus-4-6     # inherits thinking_budget: adaptive/high

  claude-faster:
    provider: deep-anthropic
    model: claude-opus-4-6
    thinking_budget: low       # overrides to adaptive/low
```

## Full Example

See [`examples/thinking_budget.yaml`](https://github.com/docker/docker-agent/blob/main/examples/thinking_budget.yaml) for a runnable config covering all providers.
