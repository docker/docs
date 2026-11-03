---
title: "Anthropic"
description: "Use Claude Sonnet 5, Claude Opus 5, and other Anthropic models with Docker Agent."
keywords: docker agent, ai agents, model providers, llm, anthropic
weight: 20
canonical: https://docs.docker.com/ai/docker-agent/providers/anthropic/
---

_Use Claude Sonnet 5, Claude Opus 5, and other Anthropic models with Docker Agent._

## Setup

```bash
# Set your API key
export ANTHROPIC_API_KEY="sk-ant-..."
```

When no `auth` block is configured, a model or [custom provider](../custom/index.md) that sets `token_key` reads its key from that variable instead of `ANTHROPIC_API_KEY`.

### Workload Identity Federation (no API key)

Authenticate with short-lived tokens minted from your own OIDC identity
provider instead of a long-lived API key. See Anthropic's
[Workload Identity Federation guide](https://platform.claude.com/docs/en/build-with-claude/workload-identity-federation)
to provision a Federation Rule, then configure Docker Agent with a typed
`auth:` block:

```yaml
providers:
  anthropic-wif:
    provider: anthropic
    auth:
      type: workload_identity_federation
      workload_identity_federation:
        federation_rule_id: fdrl_REPLACE_ME
        organization_id: 00000000-0000-0000-0000-000000000000
        # Optional: only required for target_type=SERVICE_ACCOUNT rules.
        service_account_id: svac_REPLACE_ME
        identity_token:
          # Pick exactly one of: file, env, command, url
          file: /var/run/secrets/anthropic.com/token

models:
  claude:
    provider: anthropic-wif
    model: claude-sonnet-5
```

`identity_token` accepts four mutually exclusive sources:

| Source    | When to use                                                                                                              |
| --------- | ------------------------------------------------------------------------------------------------------------------------ |
| `file`    | Kubernetes projected service-account tokens, SPIFFE/SPIRE helpers, Vault sidecars — anything that rotates a file on disk |
| `env`     | The token is already exported in an environment variable                                                                 |
| `command` | Shell out to a CLI on every refresh (`gcloud auth print-identity-token`, `az account get-access-token`, ...)              |
| `url`     | Fetch from an HTTP(S) endpoint (cloud metadata servers, GitHub Actions OIDC token URL, ...)                              |

For `url`, both the URL and any header values support `${env.VAR}` expansion
against the runtime environment (the legacy `${VAR}` form is also accepted),
which lets you wire the GitHub Actions OIDC
token endpoint without putting secrets in YAML:

```yaml
identity_token:
  url: ${env.ACTIONS_ID_TOKEN_REQUEST_URL}&audience=https://api.anthropic.com
  headers:
    Authorization: bearer ${env.ACTIONS_ID_TOKEN_REQUEST_TOKEN}
  response_field: value
```

`auth:` is mutually exclusive with `--models-gateway`. Token-refresh failures are
surfaced through the normal error path with a clear `anthropic workload
identity federation: failed to refresh identity token from <kind> source
(federation_rule=fdrl_…): ...` message in the TUI.

A complete walkthrough of all four sources lives in
[`examples/anthropic_wif.yaml`](https://github.com/docker/docker-agent/blob/main/examples/anthropic_wif.yaml).

## Configuration

### Inline

```yaml
agents:
  root:
    model: anthropic/claude-sonnet-5
```

### Named Model

```yaml
models:
  claude:
    provider: anthropic
    model: claude-sonnet-5
    max_tokens: 64000
```

## Available Models

These are Claude API IDs. See [Claude access paths](#claude-access-paths) for
platform-specific IDs, evidence and compatibility limits.

| Model ID            | Description                                         |
| ------------------- | --------------------------------------------------- |
| `claude-opus-5`     | Highest-capability Opus model; full effort ladder (low–max) |
| `claude-opus-4-8`   | Previous Opus flagship; supports task budget        |
| `claude-sonnet-5`   | Most capable Sonnet; supports extended thinking     |
| `claude-sonnet-4-6` | Previous Sonnet generation, still supported         |
| `claude-haiku-4-5`  | Fast and inexpensive, good for tight loops          |

## Claude access paths

Vendor documentation audited **2026-10-13**; integration behavior checked at
repository revision [`49e4ca6`](https://github.com/docker/docker-agent/tree/49e4ca657fd454e9d46ea258e8cf61c517df9aab).
Choose the protocol, host and credentials together. A Claude model name does not
make two endpoints interchangeable.

| Path | Client and wire protocol | Host | Authentication | Model identity |
| --- | --- | --- | --- | --- |
| Claude API | Anthropic Messages / Beta Messages | `api.anthropic.com` | API key (`ANTHROPIC_API_KEY` or `token_key`), or workload identity federation | Claude API ID |
| Custom Messages host | Anthropic client, Messages dialect | Your `base_url` | Direct auth: `token_key`, or `ANTHROPIC_API_KEY` when unset; configured workload identity federation also applies | Host's model ID |
| Models Gateway | Anthropic client through gateway | `--models-gateway` / `DOCKER_AGENT_MODELS_GATEWAY` | Docker JWT for Docker-domain gateways; `auth:` is excluded | Configured provider/model forwarded to gateway |
| Native Vertex Claude | `google` + `publisher: anthropic`; `:rawPredict` / `:streamRawPredict` | Google prediction endpoint for project/location | ADC OAuth; not an Anthropic API key | Model in URL; `anthropic_version` in body |
| Bedrock runtime | `amazon-bedrock`; **`ConverseStream`** | Regional Bedrock runtime endpoint | AWS SigV4 credentials or Bedrock bearer token | Exact AWS model or inference-profile ID |
| OpenAI-compatible aggregators | OpenAI client, Chat Completions for Claude | Alias/custom endpoint (e.g. OpenRouter, Vercel, Azure, OpenCode) | Host-specific key | Host's catalogue ID, not necessarily a Claude API ID |

These are existing adapter paths, not a guarantee of complete Claude 5.x
compatibility. Model-specific request rules, signed history, thinking and fallback
compatibility remain the subject of the planned [core work (#4583)](https://github.com/docker/docker-agent/issues/4583)
and [fallback work (#4584)](https://github.com/docker/docker-agent/issues/4584).
This guide does not mark those changes as implemented.

- The built-in `opencode` alias is OpenAI-compatible; the
  [custom OpenCode Messages provider](../opencode-zen/index.md#anthropic-compatible-messages-api)
  is a separate configuration. The `azure` alias is not native Azure AI Foundry Claude.
- Native **Mantle, InvokeModel, Bedrock Messages, Azure AI Foundry Claude and Claude
  Platform on AWS** are not wired by this integration. Compatible custom endpoints
  may exist; they have not been tested here.
- [Routing](../../configuration/routing/index.md), `first_available`, gateways and
  instrumentation compose existing clients; they are not additional wire protocols.
  See [model configuration](../../configuration/models/index.md) for composition constraints.
- A custom `base_url` implies [gateway bypass](../../configuration/models/index.md#gateway-bypass).
  It does not isolate credentials: the custom host receives the configured API key
  (defaulting to `ANTHROPIC_API_KEY` when `token_key` is unset) or federation credentials.
  Only use endpoints you trust with those credentials.

The [Anthropic construction](https://github.com/docker/docker-agent/blob/49e4ca657fd454e9d46ea258e8cf61c517df9aab/pkg/model/provider/anthropic/client.go),
[Vertex dispatch](https://github.com/docker/docker-agent/blob/49e4ca657fd454e9d46ea258e8cf61c517df9aab/pkg/model/provider/vertexai/modelgarden.go),
[Bedrock client](https://github.com/docker/docker-agent/blob/49e4ca657fd454e9d46ea258e8cf61c517df9aab/pkg/model/provider/bedrock/client.go)
and [alias definitions](https://github.com/docker/docker-agent/blob/49e4ca657fd454e9d46ea258e8cf61c517df9aab/pkg/model/provider/aliases.go)
are the repository evidence for the table.

### Claude 5.x model availability

**Documented** means corroborated by official vendor documentation, not integrated,
account-enabled or live-tested. **Limited announcement** means Anthropic announces
limited Google access, but the audited Google directory has no dedicated model card.
**Unsupported** means explicitly excluded by the vendor for that endpoint.

| Model | Canonical Claude API ID | Claude API | Native Vertex | Bedrock runtime / Converse |
| --- | --- | --- | --- | --- |
| Haiku 5.5 | `claude-haiku-5-5` | Documented | Documented | Documented |
| Opus 5 | `claude-opus-5` | Documented | Documented | Documented |
| Opus 5.5 | `claude-opus-5-5` | Documented | Documented | Documented |
| Sonnet 5 | `claude-sonnet-5` | Documented | Documented | Documented |
| Sonnet 5.5 | `claude-sonnet-5-5` | Documented | Documented | Documented |
| Fable 5 | `claude-fable-5` | Documented | Documented | Documented |
| Fable 5.1 | `claude-fable-5-1` | Documented | Documented | Documented |
| Mythos 5 | `claude-mythos-5` | Documented, gated | Limited announcement | **Unsupported** (Mantle only) |
| Mythos 5.1 | `claude-mythos-5-1` | Documented, gated | Limited announcement | Documented, gated (not Mantle) |

The audit found **24 corroborated native model/platform pairs** (9 Claude API +
7 Vertex + 8 Bedrock runtime), **2 limited-access Mythos/Google announcements** and
**1 explicitly unsupported Mythos 5/Converse pair**. This is not 26 implemented
or live-tested successes. All documented rows still require model access for your
account and chosen endpoint. An absent Google card proves neither unavailability
nor entitlement.

Sources: [Claude model overview](https://platform.claude.com/docs/en/models/overview),
[Claude on Vertex](https://platform.claude.com/docs/en/build-with-claude/claude-on-vertex-ai),
[Google's Claude directory](https://docs.cloud.google.com/gemini-enterprise-agent-platform/models/partner-models/claude),
[Mythos 5](https://platform.claude.com/docs/en/models/mythos-5/overview),
[Mythos 5.1](https://platform.claude.com/docs/en/models/mythos-5-1/overview) and
[AWS endpoint availability](https://docs.aws.amazon.com/bedrock/latest/userguide/models-endpoint-availability.html).

According to [model IDs and versions](https://platform.claude.com/docs/en/about-claude/models/model-ids-and-versions),
dateless 5.x IDs are **pinned snapshots, not rolling aliases**. Serving infrastructure
can change without replacing the model ID. Vertex's documented IDs match the bare
Claude IDs above; the bundled models.dev snapshot uses `@default` catalogue keys.
Catalogue presence is not endpoint support or an instruction to append that suffix.

### Bedrock model and profile IDs

Use the exact ID from the specific AWS model card. The global profiles below were
copied from those cards, not generated by adding a prefix to every Claude name.
For regional profiles, regions and access requirements, follow the linked card.

| Model card | AWS model ID | Documented global inference-profile ID |
| --- | --- | --- |
| [Haiku 5.5](https://docs.aws.amazon.com/bedrock/latest/userguide/model-card-anthropic-claude-haiku-5-5.html) | `anthropic.claude-haiku-5-5` | `global.anthropic.claude-haiku-5-5` |
| [Opus 5](https://docs.aws.amazon.com/bedrock/latest/userguide/model-card-anthropic-claude-opus-5.html) | `anthropic.claude-opus-5` | `global.anthropic.claude-opus-5` |
| [Opus 5.5](https://docs.aws.amazon.com/bedrock/latest/userguide/model-card-anthropic-claude-opus-5-5.html) | `anthropic.claude-opus-5-5` | `global.anthropic.claude-opus-5-5` |
| [Sonnet 5](https://docs.aws.amazon.com/bedrock/latest/userguide/model-card-anthropic-claude-sonnet-5.html) | `anthropic.claude-sonnet-5` | `global.anthropic.claude-sonnet-5` |
| [Sonnet 5.5](https://docs.aws.amazon.com/bedrock/latest/userguide/model-card-anthropic-claude-sonnet-5-5.html) | `anthropic.claude-sonnet-5-5` | `global.anthropic.claude-sonnet-5-5` |
| [Fable 5](https://docs.aws.amazon.com/bedrock/latest/userguide/model-card-anthropic-claude-fable-5.html) | `anthropic.claude-fable-5` | `global.anthropic.claude-fable-5` |
| [Fable 5.1](https://docs.aws.amazon.com/bedrock/latest/userguide/model-card-anthropic-claude-fable-5-1.html) | `anthropic.claude-fable-5-1` | `global.anthropic.claude-fable-5-1` |
| [Mythos 5](https://docs.aws.amazon.com/bedrock/latest/userguide/model-card-anthropic-claude-mythos-5.html) | `anthropic.claude-mythos-5` (Mantle only) | None; no runtime access |
| [Mythos 5.1](https://docs.aws.amazon.com/bedrock/latest/userguide/model-card-anthropic-claude-mythos-5-1.html) | `anthropic.claude-mythos-5-1` | `global.anthropic.claude-mythos-5-1` |

Do not copy one model's regional prefix set onto another. The bundled models.dev
snapshot's prefix coverage differs from these cards in both directions; it is a
catalogue artifact, not an availability claim. It has no Mythos entries at the
checked revision.

### Feature and access limits

- **Native structured output:** the Anthropic client sends JSON schema via Beta
  Messages (`structured-outputs-2025-11-13`, `output_config.format`). Enforcement
  depends on the endpoint and model, not just the presence of a schema in YAML.
- **Google policy:** [structured outputs](https://docs.cloud.google.com/gemini-enterprise-agent-platform/models/partner-models/claude/structured-outputs)
  are disabled by default by `constraints/vertexai.allowedPartnerModelFeatures`;
  your administrator must explicitly allow `structured_outputs`.
  `constraints/vertexai.allowedModels` can also restrict access to Claude models.
- **Bedrock:** this client has no native JSON-schema structured-output wiring.
  The audited Haiku 5.5, Opus 5/5.5 and Sonnet 5/5.5 cards explicitly mark
  structured outputs unsupported; the Fable/Mythos cards do not document them.
  Missing documentation is not the same as an explicit exclusion.
- **Compatible hosts:** OpenAI-compatible endpoints fronting Claude may ignore
  `response_format.json_schema`. Docker Agent [repeats the schema in the prompt](https://github.com/docker/docker-agent/blob/49e4ca657fd454e9d46ea258e8cf61c517df9aab/pkg/model/provider/openai/claude_structured_output.go)
  as a mitigation, not a constrained-decoding guarantee. [Tool mode](../../configuration/structured-output/index.md#tool)
  provides client-side validation where the model supports the required tool call.
- **Model rules override generic examples:** in particular, the
  [Fable 5.1 / Mythos 5.1](https://platform.claude.com/docs/en/models/fable-5-1/migration-guide)
  and [Opus 5.5](https://platform.claude.com/docs/en/models/opus-5-5/migration-guide)
  guides prohibit assistant prefill and forced tool choice and require thinking.
  Do not infer permission from a platform-wide structured-output example.
- **Signed history and proxies:** those migration guides describe errors or
  dropped blocks after signed prefixes are rewritten. A compatible custom host
  or gateway does not, by default, guarantee preserved thinking/compaction
  signatures, beta headers, constrained decoding, model availability or Claude
  API prices. Confirm the intermediary's contract before depending on these features.

Access prerequisites are separate from integration support. AWS Fable/Mythos
cards require data-retention opt-ins (`aws_review` or `provider_data_share`).
[Google Fable 5.1](https://docs.cloud.google.com/gemini-enterprise-agent-platform/models/partner-models/claude/fable-5-1)
requires the Advanced AI Safety Addendum, per-location model enablement and
prompt-response sharing. Anthropic's Vertex documentation describes Cyber
Verification Program gating for Mythos. **Docker Agent changes no consent,
retention or cloud setting.**

## Thinking Budget

Anthropic accepts either an integer token budget or a string effort value. When `thinking_budget` is unset, the model keeps its API default. Explicit budgets enable interleaved thinking. Sampling settings are omitted on models that reject them.

### Claude 5.x models

Docker Agent recognizes nine Claude 5.x models by ID: `claude-haiku-5-5`, `claude-opus-5`, `claude-opus-5-5`, `claude-sonnet-5`, `claude-sonnet-5-5`, `claude-fable-5`, `claude-fable-5-1`, `claude-mythos-5` and `claude-mythos-5-1`. Bedrock (`us.anthropic.…`) and Vertex (`…@date`) spellings are recognized too. For these models:

- Thinking is adaptive by default and accepts every effort level: `low`, `medium`, `high`, `xhigh` and `max`.
- With `thinking_budget` unset, Haiku 5.5, Opus 5.5, Sonnet 5.5 and Fable 5.1 are sent explicit adaptive thinking (summarized display, no effort) so the server can check replayed thinking blocks; the other five keep their API default.
- An integer token budget becomes adaptive thinking with the model's default effort; `max_tokens` is not raised.
- `temperature`, `top_p` and `top_k` are never sent.
- `thinking_budget: none` (or `0`) sends `thinking: disabled` on Haiku 5.5, Opus 5 and Sonnet 5. The other six cannot disable thinking: they get adaptive thinking with `low` effort, and Docker Agent logs a warning at startup.
- Server-side `fallbacks` involving any of these models are rejected at startup until [#4584](https://github.com/docker/docker-agent/issues/4584).
- A request that ends on an assistant message (a response prefill) is rejected before it is sent. Tool-mode structured output therefore fails after the first plain-text answer instead of retrying with a reminder.

Other Claude IDs keep the existing name-based rules.

**Token budget** (1024–32768 on models that accept manual extended thinking):

```yaml
models:
  claude-deep:
    provider: anthropic
    model: claude-sonnet-4-5
    thinking_budget: 16384 # must be < max_tokens
```

**Adaptive / effort-based** (Claude Opus 4.6+, Sonnet 4.6 — every string value is sent as adaptive thinking via `output_config.effort`):

```yaml
models:
  opus-adaptive:
    provider: anthropic
    model: claude-opus-5
    thinking_budget: adaptive # model decides effort (defaults to high)

  opus-effort:
    provider: anthropic
    model: claude-opus-5
    thinking_budget: high # low | medium | high | xhigh | max (same as adaptive/<effort>)
```

On models that reject token-based thinking (Opus 4.6+, Sonnet 5, Fable and Mythos), an integer budget is automatically coerced to `adaptive` with a logged warning. See the [Thinking / Reasoning guide](../../guides/thinking/index.md) for the full cross-provider reference.

## Interleaved Thinking

Auto-enabled whenever a thinking budget is configured on a Claude model. Allows tool calls during model reasoning for more integrated problem-solving:

```yaml
models:
  claude:
    provider: anthropic
    model: claude-sonnet-5
    provider_opts:
      interleaved_thinking: false # disable if needed
```

## Task Budget

`task_budget` caps the **total** number of tokens the model may spend across a
multi-step agentic task — combined thinking, tool calls, and final output. It
is forwarded as
[`output_config.task_budget`](https://platform.claude.com/docs/en/build-with-claude/task-budgets)
and is ideal for letting long-running agents self-regulate effort without
tightening `max_tokens` on every call.

Docker Agent automatically attaches the required `task-budgets-2026-03-13`
beta header whenever this field is set. You can configure `task_budget` on
**any** Claude model — Docker Agent never gates it by model name. **Claude
Opus 4.7, 4.8, and 5** all honor the field identically; other Claude models
(Sonnet, Haiku, etc.) are expected to reject requests that include it. Check
the Anthropic docs linked above for the current list of supported models.

```yaml
models:
  opus:
    provider: anthropic
    model: claude-opus-5
    task_budget: 128000 # integer shorthand → { type: tokens, total: 128000 }
    thinking_budget: adaptive
```

Object form (forward-compatible with future budget types):

```yaml
  opus:
    provider: anthropic
    model: claude-opus-5
    task_budget:
      type: tokens
      total: 128000
```

See the full schema on the [Model Configuration](../../configuration/models/index.md#task-budget) page.

## Server-Side Fallbacks

When the primary model refuses a request (stop reason `refusal`), Anthropic
can retry the request with backup models in a single round trip. Set
`fallbacks` in `provider_opts` to a list of model IDs, in priority order:

```yaml
models:
  opus:
    provider: anthropic
    model: claude-opus-4-8
    provider_opts:
      fallbacks:
        - claude-opus-4-7
        - claude-sonnet-4-6
```

> [!WARNING]
> **Not yet available with Claude 5.x**
>
> A fallback can change the serving model mid-response, which Docker Agent does not handle yet. A `fallbacks` list is rejected at startup when the primary model or any fallback is a Claude 5.x model (Haiku 5.5, Opus 5/5.5, Sonnet 5/5.5, Fable 5/5.1, Mythos 5/5.1). Support is tracked in [#4584](https://github.com/docker/docker-agent/issues/4584).

Docker Agent automatically attaches the required
`server-side-fallback-2026-06-01` beta header and forwards the option as
`fallbacks: [{"model": "..."}]`. The response's `model` field reports which
model actually served the request.

Fallback models receive the exact same request as the primary model
(thinking configuration, task budget, beta features, ...), so list only
models that accept the same request shape. Not available on Bedrock, Vertex
AI, or the Message Batches API.

## Native Compaction

Set `native_compaction` in `provider_opts` to let the Claude API compact the
conversation server-side instead of running the LLM summarization strategy:

```yaml
models:
  sonnet:
    provider: anthropic
    model: claude-sonnet-4-6
    provider_opts:
      native_compaction: true
```

Compaction sends the assembled conversation, current tools and stored instructions,
and replaces the whole history with a signed compaction block. This first
implementation is synchronous and compacts the whole history: it does not run
in the background or preserve a recent tail. Legacy per-call hook extras are
not included in the compaction snapshot; they are applied again on the next
normal turn. The block is persisted alongside a readable summary; the
Anthropic provider replays the block on later requests (including after a model
switch), while other providers receive the readable text.

- The compaction request goes through the same `before_llm_call` hooks and
  message transforms as a regular model call, so builtins such as
  `redact_secrets` apply to it.
- Editing the summary afterwards — through a hook rewrite or a manual edit —
  disables replay of the block; the edited text is sent instead.
- The option is explicit: a model that does not support native compaction
  fails compaction with an error rather than silently falling back to the LLM
  strategy, and a `compaction_model` pointing at another model is rejected.

Uses `compact-2026-09-04` on the Claude API. Anthropic documents
[Google Cloud beta support for on-demand compaction](https://platform.claude.com/docs/en/build-with-claude/compaction-on-demand),
but Docker Agent's Vertex client rejects `native_compaction` at construction:
**not integrated**, not unsupported by Google. AWS documents
[both compaction mechanisms as requiring InvokeModel, not Converse](https://docs.aws.amazon.com/bedrock/latest/userguide/claude-messages-compaction.html);
Docker Agent uses `ConverseStream`, so it has no Bedrock compaction path.
Summarization uses the conversation model and charges its normal rates. Compaction usage is read from `usage.iterations`. An incomplete, unsigned
or failed response never replaces the stored history. Compact before exceeding
the model's context window.

## Cache-Preserving Updates

Anthropic hashes the request prefix in order: `tools`, then `system`, then
`messages`. Any change to the system prompt, the tool list or the top-level
`output_config.effort` therefore invalidates the prompt cache for the whole
conversation. Set `cache_preserving_updates` in `provider_opts` to express
those changes as
[mid-conversation system messages](https://platform.claude.com/docs/en/build-with-claude/mid-conversation-system-messages)
instead:

```yaml
models:
  opus:
    provider: anthropic
    model: claude-opus-5
    provider_opts:
      cache_preserving_updates: true
```

The first request of a conversation fixes a baseline: its `system` blocks,
`tools` and effort. Later requests resend that baseline byte for byte and
append a `role: system` message after the last user turn carrying only what
changed since:

- system blocks appended to the prompt become `text` blocks;
- tools that disappear or reappear become `tool_removal` / `tool_addition`
  blocks (the full definitions stay in `tools`); a tool never seen before is
  declared at the end of `tools` with `defer_loading: true` and surfaced by a
  `tool_addition`;
- an effort change becomes a per-message `output_config.effort` on Claude
  Opus 5, Fable 5.1 and Mythos 5.1;
- per-call system extras (environment information, prompt files and hook
  reminders) become text-only, turn-scoped system messages. They remain in
  history verbatim but stop applying at the next user turn. If the prompt ends
  with an assistant turn, reminders stay top-level instead of using an invalid
  system-message position.

Each assistant reply records the update that preceded it, so the next request
replays every update at its original position and the cached prefix keeps
matching. The record travels with the session (`provider_state.request_context`)
and survives reloads, message edits and native compaction. Baselines persist
system instructions and tool schemas in the session store; treat session
exports with the same care as your configuration and conversation history.

Changes that have no cache-preserving form restart the baseline from the
current request (one cache miss, logged as a warning) rather than keeping stale
instructions in effect: a system block rewritten or removed, a tool definition
edited, an effort reset to the model default, or an effort change on Opus 4.8 / Fable 5. Only append-only system
prompt changes are expressed as updates.

Docker Agent attaches the `mid-conversation-tool-changes-2026-07-01`,
`mid-conversation-output-config-2026-07-01` and
`mid-conversation-system-clear-at-2026-08-21` beta headers only on requests
that actually carry a tool change, a per-message effort or a turn-scoped
reminder.

Supported on Claude Opus 4.8, Opus 5, Fable 5, Fable 5.1, Mythos 5 and
Mythos 5.1 (and their server-side fallbacks); Sonnet 5 and older models are
rejected at configuration time. Fallback models must also support the request's per-message effort features. The option requires the Beta Messages API.

## Thinking Display

Controls whether thinking blocks are returned in responses when thinking is enabled. Newer Claude models (Opus 4.7+, Fable 5) hide thinking content by default (`omitted`); Docker Agent counters this by requesting `summarized` thinking whenever an adaptive/effort-based budget is used without an explicit `thinking_display`, so reasoning stays visible in the UI. Set `thinking_display` in `provider_opts` to override:

```yaml
models:
  claude-opus-5:
    provider: anthropic
    model: claude-opus-5
    thinking_budget: adaptive
    provider_opts:
      thinking_display: omitted # "summarized" or "omitted" ("display" on pre-4.6 models only)
```

Valid values:

- `summarized`: thinking blocks are returned with summarized thinking text (Docker Agent's default for adaptive/effort-based budgets).
- `display`: thinking blocks are returned for display. Only accepted by pre-4.6 token-thinking models (e.g. Sonnet 4.5, Haiku 4.5); models from the adaptive-thinking generation onward (Opus/Sonnet 4.6+, Sonnet 5, Fable 5) reject it, and Docker Agent fails fast with a configuration error instead of sending a request the API would refuse.
- `updates`: short progress updates between tool calls, with reasoning text omitted. Available on Fable 5/5.1 and Mythos 5.1; automatically enables `thinking-display-updates-2026-08-18`. No explicit budget is required, and an unset effort keeps the model default. Other models (including configured fallbacks) are rejected locally.
- `omitted`: thinking blocks are returned with an empty thinking field; the signature is still returned for multi-turn continuity. Useful to reduce time-to-first-text-token when streaming.

Note: `thinking_display` applies to both `thinking_budget` with token counts and adaptive/effort-based budgets. For token-count budgets no default is applied (the API already defaults to `summarized`). Full thinking tokens are billed regardless of the `thinking_display` value.

> [!NOTE]
> Anthropic thinking budget values below 1024 or greater than or equal to `max_tokens` are ignored (a warning is logged).

## Preserved Thinking Compatibility

Docker Agent preserves the original ordered assistant content blocks, including
multiple thinking blocks, omitted thinking, redacted thinking and their signatures.
A hook or user edit invalidates raw replay so stale content cannot override it.
When only the rest of the turn changed (for example a sanitized tool name or a
redacted tool argument), its recorded thinking is still replayed in its original
place, with the current text and tool calls substituted. When the thinking itself changed, the turn loses it,
with a warning in the logs. If that turn belongs to the tool loop the next
request continues, and thinking is on, the API needs its original signed
thinking: Docker Agent refuses the request with an error instead, and a new user
turn starts a fresh loop.

On Claude models whose server checks thinking prefixes (Haiku 5.5, Opus 5.5,
Sonnet 5.5 and Fable 5.1), Docker Agent defaults to
`thinking.block_binding.prefix_mismatch_behavior: drop_block`: when earlier
history changes, the API drops invalidated thinking instead of rejecting the
request. Dropped blocks are reported in logs. Docker Agent does not default to
`error` because it changes earlier history itself during a session: automatic
compaction (on by default) and agent handoffs rewrite the history, or the system
prompt and tools, often in the middle of a tool loop, and under `error` every
later request in the session would then be rejected. To enforce append-only
history strictly anyway, set `provider_opts.thinking_prefix_mismatch: error`, and
avoid compaction, handoffs, `num_history_items` and toolsets that start late. You can explicitly
select `drop_block` on other supported Claude models too. Both modes attach
`thinking-binding-controls-2026-08-01`; neither repairs corrupted signatures.

## Strict Tool Arguments

Set `provider_opts.strict_tools: true` to enable strict argument generation for
compatible tool schemas, or use a list such as `strict_tools: [read_file]` to
require strict mode for named tools present in the request. Names absent from
an agent's toolset are ignored because model configurations can be shared.

- Boolean mode leaves incompatible tools non-strict and logs the reason at debug level.
- Named mode fails before sending a request if a selected tool is incompatible.
- Schemas are not rewritten: optional properties remain optional; object schemas
  must already set `additionalProperties: false`.
- The eligibility check is conservative (for example, regex patterns and recursive
  references are left non-strict). Limits cover 20 strict tools, 24 optional
  parameters and 16 union parameters, including a configured JSON output schema.

Strict arguments are separate from structured JSON responses. Normal requests
without `strict_tools` retain their existing tool schema representation.

## Cache Diagnostics

Set `provider_opts.cache_diagnostics: true` to enable Claude API cache diagnostics
(`cache-diagnosis-2026-04-07`). Docker Agent sends the preceding Anthropic response
ID from this conversation and records the diagnosis alongside the reply. Debug
logs report the reason; misses with a token count also produce a runtime warning.
Reasons include changed instructions, tools, messages or model, and an unavailable
previous request. A null/pending diagnosis is not reported as a miss.

IDs and request state belong to the conversation, not the provider client, so
parallel sessions cannot borrow one another's cache history. Docker Agent's
Vertex client rejects `cache_diagnostics` at construction; the Bedrock
`ConverseStream` integration does not wire this Claude API option.

See [`examples/anthropic-context.yaml`](https://github.com/docker/docker-agent/blob/main/examples/anthropic-context.yaml)
for a combined configuration.
