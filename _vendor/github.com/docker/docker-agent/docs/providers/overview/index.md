---
title: "Model Providers"
description: "Docker Agent supports multiple AI model providers. Choose the right one for your use case, or use multiple providers in the same configuration."
keywords: docker agent, ai agents, model providers, llm
linkTitle: "Overview"
weight: 10
canonical: https://docs.docker.com/ai/docker-agent/providers/overview/
aliases:
  - /ai/docker-agent/model-providers/
---

_Docker Agent supports multiple AI model providers. Choose the right one for your use case, or use multiple providers in the same configuration._

## Supported Providers

- [**OpenAI**](../openai/index.md) — GPT-5, GPT-5-mini, GPT-4o. The most widely used AI models.
- [**Anthropic**](../anthropic/index.md) — Claude models for coding and analysis. See [Claude access paths](../anthropic/index.md#claude-access-paths) for protocol, host and model-specific limits.
- [**Google Gemini**](../google/index.md) — Gemini 2.5 Flash, Gemini 3 Pro. Fast and cost-effective.
- [**AWS Bedrock**](../bedrock/index.md) — access Claude, Nova, Llama, and more through AWS infrastructure.
- [**Docker Model Runner**](../dmr/index.md) — run models locally with Docker. No API keys, no costs.
- [**Local Models**](../local/index.md) — run Ollama, vLLM, or LocalAI locally. No API key required.
- [**Nativ**](../nativ/index.md) — run MLX models locally on Apple Silicon through a custom provider.
- [**Provider Definitions**](../custom/index.md) — define reusable provider configurations with shared defaults for any provider type.

## Quick Comparison

| Provider            | Key              | Local? | Strengths                                             |
| ------------------- | ---------------- | ------ | ----------------------------------------------------- |
| OpenAI              | `openai`         | No     | Broad model selection, tool calling, multimodal       |
| Anthropic           | `anthropic`      | No     | Strong coding, extended thinking, large context       |
| Google              | `google`         | No     | Fast inference, competitive pricing, multimodal       |
| AWS Bedrock         | `amazon-bedrock` | No     | Enterprise features, multiple models, AWS integration |
| Docker Model Runner | `dmr`            | Yes    | No API costs, data privacy, offline capable           |
| Local Models (Ollama / vLLM) | `ollama` / custom | Yes | No API costs, full data privacy, any OpenAI-compatible server |

## Provider Credentials

Use this table to find a built-in provider's config key and authentication method. For ways to supply credentials without embedding secrets in YAML, see [Managing Secrets](../../guides/secrets/index.md). Individual provider pages cover setup and provider-specific overrides.

| Provider       | Alias            | API Key / Env Variable              |
| -------------- | ---------------- | ----------------------------------- |
| [OpenAI](../openai/index.md) | `openai` | `OPENAI_API_KEY` |
| [Anthropic](../anthropic/index.md) | `anthropic` | `ANTHROPIC_API_KEY` (override with `token_key`) |
| [Google Gemini](../google/index.md) | `google` | `GOOGLE_API_KEY` or `GEMINI_API_KEY` (override with `token_key`) |
| [AWS Bedrock](../bedrock/index.md) | `amazon-bedrock` | `AWS_BEARER_TOKEN_BEDROCK` or the standard AWS credentials chain |
| [Docker Model Runner](../dmr/index.md) | `dmr` | None (local) |
| ChatGPT (OpenAI account) | [`chatgpt`](../chatgpt/index.md) | None (sign in via `docker agent setup`) |
| OpenCode Zen   | `opencode`   | `OPENCODE_API_KEY`                  |
| OpenCode Go    | `opencode-go`    | `OPENCODE_API_KEY`                  |
| Mistral        | `mistral`        | `MISTRAL_API_KEY`                   |
| xAI (Grok)     | `xai`            | `XAI_API_KEY`                       |
| Nebius         | `nebius`         | `NEBIUS_API_KEY`                    |
| NVIDIA NIM     | `nvidia`         | `NVIDIA_API_KEY`                    |
| MiniMax        | `minimax`        | `MINIMAX_API_KEY`                   |
| Baseten        | `baseten`        | `BASETEN_API_KEY`                   |
| DaoXE          | `daoxe`          | `DAOXE_API_KEY`                     |
| OVHcloud       | `ovhcloud`       | `OVH_AI_ENDPOINTS_ACCESS_TOKEN`     |
| Groq           | `groq`           | `GROQ_API_KEY`                      |
| Fireworks AI   | `fireworks-ai`      | `FIREWORKS_API_KEY`                 |
| DeepSeek       | `deepseek`       | `DEEPSEEK_API_KEY`                  |
| Cerebras       | `cerebras`       | `CEREBRAS_API_KEY`                  |
| Together AI    | `togetherai`       | `TOGETHER_API_KEY`                  |
| Hugging Face   | `huggingface`    | `HF_TOKEN`                          |
| Cloudflare Workers AI | `cloudflare-workers-ai` | `CLOUDFLARE_API_TOKEN` + `CLOUDFLARE_ACCOUNT_ID` |
| Moonshot AI    | `moonshotai`       | `MOONSHOT_API_KEY`                  |
| Vercel AI Gateway | `vercel`      | `AI_GATEWAY_API_KEY`                |
| Cloudflare AI Gateway | `cloudflare-ai-gateway` | `CLOUDFLARE_API_TOKEN` + `CLOUDFLARE_ACCOUNT_ID` + `CLOUDFLARE_GATEWAY_ID` |
| Requesty       | `requesty`       | `REQUESTY_API_KEY`                  |
| OpenRouter     | `openrouter`     | `OPENROUTER_API_KEY`                |
| API Route      | `api-route`      | `API_ROUTE_API_KEY`                 |
| Atlas Cloud    | [`atlascloud`](../atlascloud/index.md) | `ATLASCLOUD_API_KEY`                |
| Azure OpenAI   | `azure`          | `AZURE_API_KEY` + `base_url` (override the key with `token_key`) |
| [Ollama](../local/index.md) | `ollama` | None (local; optional `base_url`) |
| GitHub Copilot | `github-copilot` | `GITHUB_TOKEN` or `GH_TOKEN` (PAT with `copilot` scope) |

## Additional Built-in Providers

The [credential table](#provider-credentials) includes built-in aliases as well as the core providers. Select one using the inline form below:

```bash
# Use built-in providers inline
agents:
  root:
    model: mistral/mistral-large-latest
```

> [!TIP]
> **Multi-provider teams**
>
> Use expensive models for complex reasoning and cheaper/local models for routine tasks. See the example below.

## Using Multiple Providers

Different agents can use different providers in the same configuration:

```yaml
models:
  claude:
    provider: anthropic
    model: claude-sonnet-5
    max_tokens: 64000
  gpt:
    provider: openai
    model: gpt-5.6-sol
  local:
    provider: dmr
    model: ai/qwen3

agents:
  root:
    model: claude # coordinator uses Claude
    sub_agents: [coder, helper]
  coder:
    model: gpt # coder uses GPT-5
  helper:
    model: local # helper runs locally for free
```

Built-in provider IDs follow the models.dev catalogue. The legacy IDs
`fireworks`, `together`, `moonshot`, and `opencode-zen` remain accepted as
`fireworks-ai`, `togetherai`, `moonshotai`, and `opencode`, respectively. New
configurations use the canonical IDs; older Docker Agent binaries may not
recognize them. API model names, URLs, and environment variables are unchanged.
Custom provider definitions take precedence over built-in IDs.
