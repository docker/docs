---
title: "AWS Bedrock"
description: "Access Claude, Nova, Llama, and more through AWS infrastructure with enterprise-grade security and compliance."
keywords: docker agent, ai agents, model providers, llm, aws bedrock
weight: 30
canonical: https://docs.docker.com/ai/docker-agent/providers/bedrock/
---

_Access Claude, Nova, Llama, and more through AWS infrastructure with enterprise-grade security and compliance._

## Prerequisites

- AWS account with Bedrock enabled in your region
- Model access granted in the [Bedrock Console](https://console.aws.amazon.com/bedrock/) (some models require approval)
- AWS credentials configured (see authentication below)

## Configuration

```yaml
models:
  bedrock-claude:
    provider: amazon-bedrock
    model: global.anthropic.claude-sonnet-5
    max_tokens: 64000
    provider_opts:
      region: us-east-1
```

## Authentication

### Option 1: Bedrock API Key (Simplest)

```bash
export AWS_BEARER_TOKEN_BEDROCK="your-key"
```

```yaml
models:
  bedrock:
    provider: amazon-bedrock
    model: global.anthropic.claude-sonnet-5
    token_key: AWS_BEARER_TOKEN_BEDROCK # env var name
    provider_opts:
      region: us-east-1
```

### Option 2: AWS Credentials (Default)

Uses the standard AWS SDK credential chain: env vars → shared credentials → config → IAM roles.

```yaml
models:
  bedrock:
    provider: amazon-bedrock
    model: global.anthropic.claude-sonnet-5
    provider_opts:
      profile: my-aws-profile
      region: us-east-1
```

### With IAM Role Assumption

```yaml
models:
  bedrock:
    provider: amazon-bedrock
    model: global.anthropic.claude-sonnet-5
    provider_opts:
      role_arn: "arn:aws:iam::123456789012:role/BedrockAccessRole"
      external_id: "my-external-id"
```

## Provider Options

| Option                   | Type   | Default                | Description                          |
| ------------------------ | ------ | ---------------------- | ------------------------------------ |
| `region`                 | string | us-east-1              | AWS region                           |
| `profile`                | string | —                      | AWS profile name                     |
| `role_arn`               | string | —                      | IAM role ARN for assume role         |
| `role_session_name`      | string | docker-agent-bedrock-session | Session name for assumed role        |
| `external_id`            | string | —                      | External ID for role assumption      |
| `endpoint_url`           | string | —                      | Custom endpoint (VPC/testing)        |
| `interleaved_thinking`   | bool   | auto                   | Allow reasoning between tool calls (Claude); auto-enabled when a thinking budget is set on a Claude model; adds the required beta header automatically |
| `disable_prompt_caching` | bool   | false                  | Disable automatic prompt caching     |

## Inference Profiles

Inference profiles route requests across the regions supported by that profile:

| Prefix    | Routes To                                |
| --------- | ---------------------------------------- |
| `global.` | Supported commercial AWS regions         |
| `us.`     | US regions only                          |
| `eu.`     | Supported European regions               |
| `apac.`   | Asia Pacific regions only                |

> [!TIP]
> **Inference profiles**
>
> Choose an exact profile from the set on that model's AWS card; prefix availability differs per model. A routing prefix alone does not establish GDPR compliance.

## Claude on Bedrock: transport and limits

Docker Agent calls **`ConverseStream`** on Bedrock runtime. Native Mantle,
InvokeModel, Bedrock Messages and Claude Platform on AWS are not wired by this
integration. A custom endpoint does not switch the runtime client's API.

[Mythos 5](https://docs.aws.amazon.com/bedrock/latest/userguide/model-card-anthropic-claude-mythos-5.html)
is Mantle-only and explicitly unsupported on runtime/Converse.
[Mythos 5.1](https://docs.aws.amazon.com/bedrock/latest/userguide/model-card-anthropic-claude-mythos-5-1.html)
is documented on runtime, with gated access, not Mantle. Vendor support does not
mean this client's complete 5.x compatibility work has shipped.

[AWS compaction](https://docs.aws.amazon.com/bedrock/latest/userguide/claude-messages-compaction.html)
requires InvokeModel for both `compact-2026-01-12` and `compact-2026-09-04`;
neither works through Converse. Haiku 5.5, Opus 5/5.5 and Sonnet 5/5.5 cards
document the latter; Fable 5/5.1 and Mythos 5.1 cards do not document compaction.
The current client also does not wire Claude API cache diagnostics, server-side
fallbacks or native JSON-schema structured output. General Bedrock feature
statements do not override a specific model card.

See [Claude access paths](../anthropic/index.md#claude-access-paths) for the dated
nine-model evidence matrix, exact card-sourced IDs, structured-output limits and
account prerequisites. `apac.` in the general profile table is not a documented
Claude 5.x profile in this audit; do not generate IDs by copying prefixes.

## Thinking Budget (Claude on Bedrock)

Bedrock Claude models support extended thinking — an internal reasoning phase before the model produces its response. Set `thinking_budget` to a token count (1024–32768) or an effort level string that maps automatically:

| Effort level | Token budget |
| ------------ | ------------ |
| `minimal`    | 1,024        |
| `low`        | 2,048        |
| `medium`     | 8,192        |
| `high`       | 16,384       |
| `xhigh`/`max`| 32,768       |

```yaml
models:
  bedrock-claude-thinking:
    provider: amazon-bedrock
    model: global.anthropic.claude-sonnet-5
    thinking_budget: 8192   # tokens, or use an effort string like "medium"
    max_tokens: 16384       # must be > thinking_budget
    provider_opts:
      region: us-east-1
```

`thinking_budget` must be ≥ 1024 and less than `max_tokens`. Values outside this range are logged as a warning and ignored.

### Adaptive thinking (Opus 4.6+)

Newer Claude Opus models (4.6, 4.7, 4.8) **reject token-based thinking** — Bedrock returns a `ValidationException` asking for `thinking.type=adaptive`. For these models, use adaptive thinking:

```yaml
models:
  bedrock-opus-adaptive:
    provider: amazon-bedrock
    model: global.anthropic.claude-opus-4-8
    thinking_budget: adaptive/high   # adaptive | adaptive/low | adaptive/medium | adaptive/high | adaptive/xhigh | adaptive/max
    provider_opts:
      region: us-east-1
```

Docker Agent recognizes these models (including Bedrock-style IDs) and transparently coerces a configured token budget or effort level to adaptive thinking, logging a warning — so `thinking_budget: 32768` on Opus 4.8 won't fail, but `adaptive` or `adaptive/<effort>` is the recommended configuration. On older models that still use token-based thinking (e.g. Sonnet 4.5), `adaptive` is forwarded as-is and rejected by Bedrock — use a token count or effort level there instead.

> [!NOTE]
> **Temperature and top_p**
>
> Bedrock Claude suppresses `temperature` and `top_p` while extended thinking is active — Anthropic requires `temperature=1.0` internally.

### Claude 5.x models

The nine Claude 5.x models listed in the [Anthropic provider guide](../anthropic/index.md#claude-5x-models) follow the same rules on Bedrock, including Bedrock-style IDs such as `us.anthropic.claude-sonnet-5-5-v1:0`:

- Effort levels and `adaptive` are sent as adaptive thinking. An integer token budget becomes adaptive thinking with the model's default effort.
- `temperature`, `top_p` and `top_k` are never sent.
- `thinking_budget: none` sends `thinking: {"type": "disabled"}` on Haiku 5.5, Opus 5 and Sonnet 5. The other six get adaptive thinking with `low` effort and a startup warning.
- Reasoning blocks are stored and replayed exactly as Bedrock returned them, in order: signed thinking (including signature-only thinking, the default display on these models), redacted thinking and the text and tool calls around them. If a message changed after it was received but its reasoning did not (for example a sanitized tool name or a redacted tool argument), its recorded reasoning is still replayed in its original place, with the current text and tool calls substituted. If the reasoning itself changed, it is sent without it, with a warning in the logs; inside the tool loop the next request continues, with thinking on, the request is refused instead, and a new user turn starts a fresh loop. Sessions saved by older versions keep that flattened replay; older versions ignore the stored blocks.
- Mythos 5 is refused at startup: AWS serves it only through Bedrock Mantle, which Docker Agent does not support. Use Mythos 5.1, which is available on Converse, or the Claude API.

> [!NOTE]
> **Not yet verified live**
>
> The `disabled` value follows the Claude Messages contract; it has not yet been verified against live Bedrock Converse.

## Interleaved Thinking (Claude on Bedrock)

Interleaved thinking lets the model reason between tool calls, not just at the start. This is useful for complex agentic tasks. Enable it alongside a thinking budget:

```yaml
models:
  bedrock-claude-interleaved:
    provider: amazon-bedrock
    model: global.anthropic.claude-sonnet-5
    thinking_budget: high
    provider_opts:
      region: us-east-1
      # interleaved_thinking is auto-enabled when thinking_budget is set
```

Docker Agent auto-enables `interleaved_thinking` whenever a thinking budget is configured on a Bedrock-hosted Claude model and automatically adds the `interleaved-thinking-2025-05-14` beta header. If you set `interleaved_thinking: false` while a thinking budget is active, a warning is logged because the budget may be ignored by Bedrock without the beta header.

See the [Thinking / Reasoning guide](../../guides/thinking/index.md) for a full cross-provider overview.

## Prompt Caching

Automatically enabled for supported models to reduce latency and costs. System prompts, tool definitions, and recent messages are cached with a 5-minute TTL.

```bash
# Disable if needed
provider_opts:
  disable_prompt_caching: true
```
