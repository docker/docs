---
title: "Model Providers"
description: "docker-agent supports multiple AI model providers. Choose the right one for your use case, or use multiple providers in the same configuration."
keywords: docker agent, ai agents, model providers, llm
linkTitle: "Overview"
weight: 10
aliases:
  - /ai/docker-agent/model-providers/
---

_docker-agent supports multiple AI model providers. Choose the right one for your use case, or use multiple providers in the same configuration._

## Supported Providers

- [**OpenAI**](../openai/index.md) — GPT-5, GPT-5-mini, GPT-4o. The most widely used AI models.
- [**Anthropic**](../anthropic/index.md) — Claude Sonnet 4.5, Claude Opus 4.7. Excellent for coding and analysis.
- [**Google Gemini**](../google/index.md) — Gemini 2.5 Flash, Gemini 3 Pro. Fast and cost-effective.
- [**AWS Bedrock**](../bedrock/index.md) — access Claude, Nova, Llama, and more through AWS infrastructure.
- [**Docker Model Runner**](../dmr/index.md) — run models locally with Docker. No API keys, no costs.
- [**Provider Definitions**](../custom/index.md) — define reusable provider configurations with shared defaults for any provider type.

## Quick Comparison

| Provider            | Key              | Local? | Strengths                                             |
| ------------------- | ---------------- | ------ | ----------------------------------------------------- |
| OpenAI              | `openai`         | No     | Broad model selection, tool calling, multimodal       |
| Anthropic           | `anthropic`      | No     | Strong coding, extended thinking, large context       |
| Google              | `google`         | No     | Fast inference, competitive pricing, multimodal       |
| AWS Bedrock         | `amazon-bedrock` | No     | Enterprise features, multiple models, AWS integration |
| Docker Model Runner | `dmr`            | Yes    | No API costs, data privacy, offline capable           |

## Additional Built-in Providers

docker-agent also includes built-in aliases for these providers:

| Provider       | Alias            | API Key / Env Variable              |
| -------------- | ---------------- | ----------------------------------- |
| OpenCode Zen   | `opencode-zen`   | `OPENCODE_API_KEY`                  |
| OpenCode Go    | `opencode-go`    | `OPENCODE_API_KEY`                  |
| Mistral        | `mistral`        | `MISTRAL_API_KEY`                   |
| xAI (Grok)     | `xai`            | `XAI_API_KEY`                       |
| Nebius         | `nebius`         | `NEBIUS_API_KEY`                    |
| MiniMax        | `minimax`        | `MINIMAX_API_KEY`                   |
| Baseten        | `baseten`        | `BASETEN_API_KEY`                   |
| OVHcloud       | `ovhcloud`       | `OVH_AI_ENDPOINTS_ACCESS_TOKEN`     |
| Groq           | `groq`           | `GROQ_API_KEY`                      |
| Fireworks AI   | `fireworks`      | `FIREWORKS_API_KEY`                 |
| DeepSeek       | `deepseek`       | `DEEPSEEK_API_KEY`                  |
| Cerebras       | `cerebras`       | `CEREBRAS_API_KEY`                  |
| Together AI    | `together`       | `TOGETHER_API_KEY`                  |
| Hugging Face   | `huggingface`    | `HF_TOKEN`                          |
| Cloudflare Workers AI | `cloudflare-workers-ai` | `CLOUDFLARE_API_TOKEN` + `CLOUDFLARE_ACCOUNT_ID` |
| Moonshot AI    | `moonshot`       | `MOONSHOT_API_KEY`                  |
| Vercel AI Gateway | `vercel`      | `AI_GATEWAY_API_KEY`                |
| Cloudflare AI Gateway | `cloudflare-ai-gateway` | `CLOUDFLARE_API_TOKEN` + `CLOUDFLARE_ACCOUNT_ID` + `CLOUDFLARE_GATEWAY_ID` |
| Requesty       | `requesty`       | `REQUESTY_API_KEY`                  |
| OpenRouter     | `openrouter`     | `OPENROUTER_API_KEY`                |
| Azure OpenAI   | `azure`          | `AZURE_API_KEY` + `base_url`        |
| Ollama         | `ollama`         | None (local; optional `base_url`)   |
| GitHub Copilot | `github-copilot` | `GITHUB_TOKEN` (PAT with `copilot` scope) |

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
    model: claude-sonnet-4-5
    max_tokens: 64000
  gpt:
    provider: openai
    model: gpt-5
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
