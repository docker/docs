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
    model: global.anthropic.claude-sonnet-4-5-20250929-v1:0
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
    model: global.anthropic.claude-sonnet-4-5-20250929-v1:0
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
    model: global.anthropic.claude-sonnet-4-5-20250929-v1:0
    provider_opts:
      profile: my-aws-profile
      region: us-east-1
```

### With IAM Role Assumption

```yaml
models:
  bedrock:
    provider: amazon-bedrock
    model: global.anthropic.claude-sonnet-4-5-20250929-v1:0
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

Use inference profile prefixes for optimal routing:

| Prefix    | Routes To                                |
| --------- | ---------------------------------------- |
| `global.` | All commercial AWS regions (recommended) |
| `us.`     | US regions only                          |
| `eu.`     | EU regions only (GDPR compliance)        |
| `apac.`   | Asia Pacific regions only                |

> [!TIP]
> **Inference profiles**
>
> Use `global.` prefix on model IDs for automatic cross-region routing. Use `eu.` prefix for GDPR compliance.

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
    model: global.anthropic.claude-sonnet-4-5-20250929-v1:0
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

docker-agent recognizes these models (including Bedrock-style IDs) and transparently coerces a configured token budget or effort level to adaptive thinking, logging a warning — so `thinking_budget: 32768` on Opus 4.8 won't fail, but `adaptive` or `adaptive/<effort>` is the recommended configuration. On older models that still use token-based thinking (e.g. Sonnet 4.5), `adaptive` is forwarded as-is and rejected by Bedrock — use a token count or effort level there instead.

> [!NOTE]
> **Temperature and top_p**
>
> Bedrock Claude suppresses `temperature` and `top_p` while extended thinking is active — Anthropic requires `temperature=1.0` internally.

## Interleaved Thinking (Claude on Bedrock)

Interleaved thinking lets the model reason between tool calls, not just at the start. This is useful for complex agentic tasks. Enable it alongside a thinking budget:

```yaml
models:
  bedrock-claude-interleaved:
    provider: amazon-bedrock
    model: global.anthropic.claude-sonnet-4-5-20250929-v1:0
    thinking_budget: high
    provider_opts:
      region: us-east-1
      # interleaved_thinking is auto-enabled when thinking_budget is set
```

docker-agent auto-enables `interleaved_thinking` whenever a thinking budget is configured on a Bedrock-hosted Claude model and automatically adds the `interleaved-thinking-2025-05-14` beta header. If you set `interleaved_thinking: false` while a thinking budget is active, a warning is logged because the budget may be ignored by Bedrock without the beta header.

See the [Thinking / Reasoning guide](../../guides/thinking/index.md) for a full cross-provider overview.

## Prompt Caching

Automatically enabled for supported models to reduce latency and costs. System prompts, tool definitions, and recent messages are cached with a 5-minute TTL.

```bash
# Disable if needed
provider_opts:
  disable_prompt_caching: true
```
