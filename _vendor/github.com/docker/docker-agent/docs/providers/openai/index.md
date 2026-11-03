---
title: "OpenAI"
description: "Use GPT-5.6, GPT-4o, GPT-4.1, and other OpenAI models with Docker Agent."
keywords: docker agent, ai agents, model providers, llm, openai
weight: 200
canonical: https://docs.docker.com/ai/docker-agent/providers/openai/
---

_Use GPT-5.6, GPT-4o, GPT-4.1, and other OpenAI models with Docker Agent._

## Setup

```bash
# Set your API key
export OPENAI_API_KEY="sk-..."
```

> [!TIP]
> No API key? A ChatGPT Plus/Pro/Business subscription can be used instead
> through the [`chatgpt` provider](../chatgpt/index.md): sign in once with
> `docker agent setup` (pick chatgpt).

## Configuration

### Inline

```yaml
agents:
  root:
    model: openai/gpt-5.6
```

### Named Model

```yaml
models:
  gpt:
    provider: openai
    model: gpt-5.6
    max_tokens: 4000
```

## Available Models

| Model            | Best For                                             |
| ---------------- | ----------------------------------------------------- |
| `gpt-5.6`         | Alias for `gpt-5.6-sol`; tracks the flagship model    |
| `gpt-5.6-sol`     | Frontier model, most capable, complex reasoning       |
| `gpt-5.6-terra`   | Everyday workhorse; successor to the `-mini` tier     |
| `gpt-5.6-luna`    | High-volume, cost-efficient; successor to `-nano` tier |
| `gpt-4.1`         | Previous-generation flagship                          |
| `gpt-4.1-mini`    | Previous-generation fast, cost-effective model        |
| `gpt-4o`          | Multimodal, balanced performance                      |
| `gpt-4o-mini`     | Cheapest, fast for simple tasks                       |

Starting with GPT-5.6, OpenAI renamed the `-mini`/`-nano` size tiers to `-terra`/`-luna` (with `-sol` denoting the frontier tier previously left unsuffixed).

Find more model names at [modelnames.ai](https://modelnames.ai/) or in the [official OpenAI docs](https://platform.openai.com/docs/models).

## Service Tier (Fast and Ultrafast Modes)

Set `provider_opts.service_tier` to request OpenAI's [Fast mode](https://developers.openai.com/api/docs/guides/fast-mode):

```yaml
models:
  fast-gpt:
    provider: openai
    model: gpt-5.6
    provider_opts:
      service_tier: fast
```

OpenAI also accepts `priority` for Fast mode. It provides faster processing at premium pricing on supported models, without reducing reasoning effort. This is independent of `thinking_budget` and applies to all requests using the configured model, including internal calls such as title generation and compaction.

The value is forwarded unchanged to Chat Completions (including reranking) and Responses requests, over either SSE or WebSocket. OpenAI-compatible providers using these APIs also receive the option when set; the endpoint must support it. Other tiers, such as `auto`, `default`, and `flex`, can also be requested; availability and valid values depend on the API and model. When omitted or empty, no `service_tier` is sent, leaving the API's default behavior unchanged. Non-string values are ignored.

For models that support it, such as GPT-6 Astra and GPT-6.1 Sol, request [Ultrafast mode](https://developers.openai.com/api/docs/guides/ultrafast-mode) with `service_tier: ultrafast`.

Cost estimates use the **actual response tier**, not the requested tier. For models served by the official OpenAI endpoint (including an explicitly configured `base_url: https://api.openai.com/v1`), an actual Fast tier (`fast` or `priority`) uses a catalogue-derived premium when matching Vercel Standard/Fast entries have the same base rates and a uniform premium across token types. Otherwise it is estimated at 2× catalogue rates. `ultrafast` is estimated at 6× rates. The response selects the tier; Docker Agent does not maintain a model allowlist. These adjustments cover input, cached input, cache writes, and output, including the applicable long-context band. A response reporting `default` retains standard pricing, even if Fast was requested.

> [!WARNING]
> Missing or unrecognized response tiers, custom endpoints (including `OPENAI_BASE_URL`), rule-based routers, and other providers retain catalogue pricing. Gateway model IDs already priced as `-fast` are not multiplied again. The response reports the tier, not a price. Catalogue-derived premiums cover differing Fast rates such as GPT-4o and GPT-5.5; missing or inconsistent tier metadata falls back to the estimate. Estimates do not include endpoint surcharges or negotiated discounts. Use a model's [`cost` override](../../configuration/models/index.md#custom-token-pricing) when a model's published tier rates differ or automatic pricing does not apply. Overrides replace the entire price table, including context bands, and are never adjusted for the service tier—even after a downgrade.

See [`examples/openai-service-tier.yaml`](https://github.com/docker/docker-agent/blob/main/examples/openai-service-tier.yaml) for a complete example.

## Hosted Tool Search

Set `provider_opts.native_tool_search: true` to let OpenAI's hosted [tool search](https://developers.openai.com/api/docs/guides/tools-tool-search) discover the agent's [deferred tools](../../configuration/tools/index.md#provider-native-tool-search-openai) instead of the built-in `search_tool`/`add_tool` pair:

```yaml
models:
  gpt:
    provider: openai
    model: gpt-5.6-sol
    provider_opts:
      native_tool_search: true
```

Deferred tools are declared with `defer_loading: true` alongside a server-executed `tool_search` tool on Responses API requests. Strictly opt-in: only `openai` models that support deferred tools (gpt-5.4 and later) qualify; the option is ignored elsewhere (`chatgpt`, OpenAI-compatible endpoints including a custom `base_url`, Chat Completions requests) and the legacy behaviour is kept.

See [`examples/deferred_native_tool_search.yaml`](https://github.com/docker/docker-agent/blob/main/examples/deferred_native_tool_search.yaml) for a complete example.

## Thinking Budget

OpenAI reasoning models (o-series, gpt-5, gpt-5-mini, gpt-5.6 family, gpt-6 family) support extended thinking through the `reasoning_effort` API parameter. Set `thinking_budget` to control the effort level:

```yaml
models:
  gpt-thinker:
    provider: openai
    model: gpt-5.6
    thinking_budget: high   # none | minimal | low | medium | high | xhigh | max
```

**Effort levels:**

| Level     | Description                                              |
| --------- | -------------------------------------------------------- |
| `none`    | No reasoning. On GPT-5.x from 5.6 onward and GPT-6 Sol/Luna this is a real API value that is sent as-is; on Astra and older models it just disables the local `thinking_budget` (the API's own default still applies). |
| `minimal` | Fastest; lightest reasoning pass. Not accepted on `gpt-5.6`+ (dropped from the API). |
| `low`     | Quick reasoning for straightforward tasks.               |
| `medium`  | Balanced default.                                        |
| `high`    | More thorough; recommended for complex tasks.            |
| `xhigh`   | Near-maximum effort; slower but most accurate. Requires `gpt-5.2`+. |
| `max`     | Maximum effort. Requires `gpt-5.6`+ (Sol/Terra/Luna).    |

Token counts, `adaptive`, and `adaptive/<effort>` are rejected with a configuration error at request time. Older models (o1, o3-mini) only accept `low`/`medium`/`high`; `xhigh` requires `gpt-5.2`+; `max` requires `gpt-5.6`+; `none` is supported on GPT-5.x from 5.6 onward and GPT-6 Sol/Luna, but not Astra; `minimal` is not accepted on `gpt-5.6`+.

GPT-6 Sol/Luna also support function tools with `thinking_budget: none` on Chat Completions. Tool calling with reasoning requires the Responses API, which Docker Agent selects automatically for these models on the `openai` provider. Astra requires Responses for tool calling.

> [!WARNING]
> **Hidden reasoning tokens**
>
> OpenAI reasoning models always produce hidden reasoning tokens that count against `max_tokens` — even with `thinking_budget: none` on older models. Docker Agent automatically raises the output-token floor for its internal low-effort calls so reasoning cannot starve visible text output.

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
    model: gpt-5.6-terra
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

## Prompt Cache Diagnostics

Enable `provider_opts.cache_diagnostics` to compare a response's cache reuse with
its preceding assistant turn. Available on the OpenAI Responses API with
GPT-5.6 and later supported models:

```yaml
models:
  gpt:
    provider: openai
    model: gpt-5.6-sol
    provider_opts:
      cache_diagnostics: true
```

Run with `--debug` to log the diagnostic result, cache-miss reason, reusable and
missed token estimates, and actual cached tokens. Reasons can include changed
tools, input, reasoning effort, or service tier. Diagnostics do not change caching
behavior or load conversation history.

The comparison comes from the current session's preceding assistant response,
not shared client state. The first turn, a model switch, or a session without a
comparison ID skips the comparison. Expired diagnostic records and unavailable
results are logged without failing the conversation. This option is ignored for
Chat Completions, the `chatgpt` provider, and custom `base_url` endpoints.

## Extra Request Body

Use [`provider_opts.extra_body`](../../configuration/models/index.md#extra-request-body)
to pass backend-specific fields to Chat Completions requests, including when you
use a custom `base_url`. For example, a local Qwen server may support
`chat_template_kwargs: {enable_thinking: false}`. Explicit fields override
generated values. This option does not apply to the Responses API.

## Reasoning State Preservation

Reasoning replay is opt-in, so existing configurations keep their request and
session-storage behavior. Enable it for first-party OpenAI Responses sessions,
including the `chatgpt` provider:

```yaml
provider_opts:
  preserve_reasoning: true
```

This preserves encrypted reasoning items and assistant message phases in the
session and replays the ordered output alongside tool results. Enabling
`native_tool_search` also enables output replay because hosted search needs to
retain its loaded tools. `cache_diagnostics` alone saves only response metadata;
it does not enable reasoning replay.

Opaque state is only replayed to the same provider and resolved model. When hooks
edit assistant text, reasoning summaries, or tool calls, Docker Agent falls back
to the edited conversation instead of replaying stale output. Hosted tool-search
history is discarded when it references tools no longer in the allowed catalog.
Other providers and custom endpoints retain their existing behavior.

When enabled, this adds opaque data to the local session store and may increase
input-token usage and trigger earlier compaction; it does not expose hidden
reasoning as readable text. Reasoning tokens remain part of the context estimate,
and the existing session compaction strategy still applies. With replay enabled,
WebSocket requests replay the full history without also chaining it through `previous_response_id`.

See [OpenAI's reasoning guide](https://developers.openai.com/api/docs/guides/reasoning)
and [cache diagnostics guide](https://developers.openai.com/api/docs/guides/prompt-caching/diagnostics).
