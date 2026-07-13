---
title: Model providers
description: Get API keys and configure cloud model providers for Docker Agent
keywords: [docker agent, model providers, api keys, anthropic, openai, google, gemini, groq, deepseek, cerebras, fireworks, together, moonshot, openrouter, ovhcloud, baseten, vercel, cloudflare]
weight: 10
---

To run Docker Agent, you need a model provider. You can either use a cloud provider
with an API key or run models locally with [Docker Model
Runner](local-models.md).

This guide covers cloud providers. For the local alternative, see [Local
models with Docker Model Runner](local-models.md).

## Supported providers

Docker Agent supports these cloud model providers:

- Anthropic — Claude models
- Baseten — open-weight models via Baseten
- Cerebras — fast inference for open-weight models
- Cloudflare AI Gateway — multi-provider gateway with caching and observability
- Cloudflare Workers AI — open-weight models on the Cloudflare edge
- DeepSeek — DeepSeek chat and reasoning models
- Fireworks AI — fast inference for open-weight models
- Google — Gemini models
- Groq — ultra-low-latency open-weight models
- Moonshot AI — Kimi models
- OpenAI — GPT models
- OpenRouter — unified gateway to hundreds of models
- OVHcloud — EU-hosted open-weight models
- Together AI — large catalog of open-weight models
- Vercel AI Gateway — unified gateway to OpenAI, Anthropic, Google, and more

## Anthropic

Anthropic provides the Claude family of models, including Claude Sonnet and
Claude Opus.

To get an API key:

1. Go to [console.anthropic.com](https://console.anthropic.com/).
2. Sign up or sign in to your account.
3. Navigate to the API Keys section.
4. Create a new API key.
5. Copy the key.

Set your API key as an environment variable:

```console
$ export ANTHROPIC_API_KEY=your_key_here
```

Use Anthropic models in your agent configuration:

```yaml
agents:
  root:
    model: anthropic/claude-sonnet-4-5
    instruction: You are a helpful coding assistant
```

Available models include:

- `anthropic/claude-sonnet-4-5`
- `anthropic/claude-opus-4-5`
- `anthropic/claude-haiku-4-5`

## Baseten

[Baseten](https://www.baseten.co/) provides AI models through an
OpenAI-compatible API. It is a good choice for deploying your own models or
accessing hosted open-weight models.

1. Get an API key from [Baseten](https://www.baseten.co/).
2. Set the environment variable:

```console
$ export BASETEN_API_KEY=your_key_here
```

Use Baseten models in your agent configuration:

```yaml
agents:
  root:
    model: baseten/deepseek-ai/DeepSeek-V3.1
    instruction: You are a helpful assistant
```

Or with a named model for more control:

```yaml
models:
  baseten_model:
    provider: baseten
    model: deepseek-ai/DeepSeek-V3.1
    max_tokens: 8192

agents:
  root:
    model: baseten_model
    instruction: You are a helpful assistant
```

## Cerebras

[Cerebras](https://www.cerebras.ai/) serves open-weight models (including
GPT-OSS and GLM) on its wafer-scale hardware through an OpenAI-compatible API,
delivering some of the highest inference speeds available.

1. Create an API key from the [Cerebras Cloud console](https://cloud.cerebras.ai/).
2. Set the environment variable:

```console
$ export CEREBRAS_API_KEY=your_key_here
```

Use Cerebras models in your agent configuration:

```yaml
agents:
  root:
    model: cerebras/gpt-oss-120b
    instruction: You are a helpful assistant
```

Available models include:

- `cerebras/gpt-oss-120b`
- `cerebras/zai-glm-4.7`

## Cloudflare AI Gateway

[Cloudflare AI Gateway](https://developers.cloudflare.com/ai-gateway/) is a
single OpenAI-compatible endpoint that routes to models from OpenAI, Anthropic,
Workers AI, and more, with caching, rate limiting, and observability.

The gateway endpoint is account- and gateway-scoped. Three environment variables
are required:

```console
$ export CLOUDFLARE_ACCOUNT_ID=your_account_id
$ export CLOUDFLARE_GATEWAY_ID=your_gateway_id
$ export CLOUDFLARE_API_TOKEN=your_api_token
```

Use Cloudflare AI Gateway in your agent configuration:

```yaml
agents:
  root:
    model: cloudflare-ai-gateway/workers-ai/@cf/meta/llama-3.1-8b-instruct
    instruction: You are a helpful assistant
```

Or with a named model:

```yaml
models:
  cf_gateway_model:
    provider: cloudflare-ai-gateway
    model: "workers-ai/@cf/meta/llama-3.1-8b-instruct"

agents:
  root:
    model: cf_gateway_model
    instruction: You are a helpful assistant
```

> [!NOTE]
> The alias sends your token in the standard `Authorization: Bearer` header,
> which works for unauthenticated gateways (the default) routing to Workers AI
> models. Gateways with authentication enabled require the `cf-aig-authorization`
> header, which is not supported by this alias. For that setup, use a [custom
> provider](reference/config.md#models) instead.

## Cloudflare Workers AI

[Cloudflare Workers AI](https://developers.cloudflare.com/workers-ai/) runs
open-weight models (Llama, Mistral, Qwen, Gemma, and more) on Cloudflare's
global edge network.

Workers AI is account-scoped, so two environment variables are required:

```console
$ export CLOUDFLARE_ACCOUNT_ID=your_account_id
$ export CLOUDFLARE_API_TOKEN=your_api_token
```

Use Cloudflare Workers AI in your agent configuration:

```yaml
agents:
  root:
    model: cloudflare-workers-ai/@cf/meta/llama-3.1-8b-instruct
    instruction: You are a helpful assistant
```

Available models include `@cf/meta/llama-3.1-8b-instruct`,
`@cf/mistralai/mistral-small-3.1-24b-instruct`, and more. See the
[Workers AI models catalog](https://developers.cloudflare.com/workers-ai/models/)
for the full list.

## DeepSeek

[DeepSeek](https://www.deepseek.com/) serves its frontier chat and reasoning
models through an OpenAI-compatible API, with strong price/performance on coding
and reasoning tasks.

1. Create an API key from the [DeepSeek Platform](https://platform.deepseek.com/api_keys).
2. Set the environment variable:

```console
$ export DEEPSEEK_API_KEY=your_key_here
```

Use DeepSeek models in your agent configuration:

```yaml
agents:
  root:
    model: deepseek/deepseek-chat
    instruction: You are a helpful coding assistant
```

Available models include:

- `deepseek/deepseek-chat` — DeepSeek-V3, general-purpose chat and tool calling
- `deepseek/deepseek-reasoner` — DeepSeek-R1, extended-reasoning model

## Fireworks AI

[Fireworks AI](https://fireworks.ai/) is a fast inference host for open-weight
models, serving Kimi K2, Llama, Qwen, DeepSeek, GLM, and others through an
OpenAI-compatible API.

1. Create an API key from the [Fireworks dashboard](https://fireworks.ai/account/api-keys).
2. Set the environment variable:

```console
$ export FIREWORKS_API_KEY=your_key_here
```

Use Fireworks AI models in your agent configuration:

```yaml
agents:
  root:
    model: fireworks/accounts/fireworks/models/kimi-k2-instruct
    instruction: You are a helpful assistant
```

Fireworks model IDs use the `accounts/fireworks/models/<name>` form. See the
[Fireworks model library](https://fireworks.ai/models) for current IDs.

## Google Gemini

Google provides the Gemini family of models.

To get an API key:

1. Go to [aistudio.google.com/apikey](https://aistudio.google.com/apikey).
2. Sign in with your Google account.
3. Create an API key.
4. Copy the key.

Set your API key as an environment variable:

```console
$ export GOOGLE_API_KEY=your_key_here
```

Use Gemini models in your agent configuration:

```yaml
agents:
  root:
    model: google/gemini-2.5-flash
    instruction: You are a helpful coding assistant
```

Available models include:

- `google/gemini-2.5-flash`
- `google/gemini-2.5-pro`

## Groq

[Groq](https://groq.com/) serves open-weight models on its LPU inference engine
through an OpenAI-compatible API, with a focus on very low latency.

1. Create an API key from the [Groq Console](https://console.groq.com/keys).
2. Set the environment variable:

```console
$ export GROQ_API_KEY=your_key_here
```

Use Groq models in your agent configuration:

```yaml
agents:
  root:
    model: groq/llama-3.3-70b-versatile
    instruction: You are a helpful assistant
```

Available models include `llama-3.3-70b-versatile`, `llama-3.1-8b-instant`, and
more. See the [Groq models documentation](https://console.groq.com/docs/models)
for current model IDs.

## Moonshot AI

[Moonshot AI](https://www.moonshot.ai/) serves its Kimi model family through an
OpenAI-compatible API. Kimi K2 models are well-suited for coding and agentic
tasks.

1. Create an API key from the [Moonshot AI console](https://platform.moonshot.ai/console/api-keys).
2. Set the environment variable:

```console
$ export MOONSHOT_API_KEY=your_key_here
```

Use Moonshot AI models in your agent configuration:

```yaml
agents:
  root:
    model: moonshot/kimi-k2-0905-preview
    instruction: You are a helpful assistant
```

Available models include:

- `moonshot/kimi-k2-0905-preview`
- `moonshot/kimi-k2-turbo-preview`
- `moonshot/kimi-k2-thinking`

## OpenAI

OpenAI provides the GPT family of models, including GPT-5 and GPT-5 mini.

To get an API key:

1. Go to [platform.openai.com/api-keys](https://platform.openai.com/api-keys).
2. Sign up or sign in to your account.
3. Navigate to the API Keys section.
4. Create a new API key.
5. Copy the key.

Set your API key as an environment variable:

```console
$ export OPENAI_API_KEY=your_key_here
```

Use OpenAI models in your agent configuration:

```yaml
agents:
  root:
    model: openai/gpt-5
    instruction: You are a helpful coding assistant
```

Available models include:

- `openai/gpt-5`
- `openai/gpt-5-mini`

## OpenRouter

[OpenRouter](https://openrouter.ai/) provides access to hundreds of models from
many providers through a single OpenAI-compatible API, with automatic failover
and unified billing.

1. Get an API key from [OpenRouter](https://openrouter.ai/settings/keys).
2. Set the environment variable:

```console
$ export OPENROUTER_API_KEY=your_key_here
```

Use OpenRouter in your agent configuration:

```yaml
agents:
  root:
    model: openrouter/meta-llama/llama-3.3-70b-instruct
    instruction: You are a helpful assistant
```

OpenRouter model IDs include the upstream provider name (for example
`anthropic/claude-sonnet-4-5` or `meta-llama/llama-3.3-70b-instruct`). Docker
Agent preserves the full upstream model ID after the first slash. See the
[OpenRouter models list](https://openrouter.ai/models) for available models.

## OVHcloud

[OVHcloud AI Endpoints](https://endpoints.ai.cloud.ovh.net/) serves open-weight
models through an OpenAI-compatible API, hosted in the EU. Several models are
available on a rate-limited free tier with no billing setup required.

1. Create an access token from the
   [OVHcloud AI Endpoints portal](https://endpoints.ai.cloud.ovh.net/).
2. Set the environment variable:

```console
$ export OVH_AI_ENDPOINTS_ACCESS_TOKEN=your_token_here
```

Use OVHcloud models in your agent configuration:

```yaml
agents:
  root:
    model: ovhcloud/Qwen3.5-397B-A17B
    instruction: You are a helpful assistant
```

Available models include `Qwen3.5-397B-A17B`, `Qwen3-32B`,
`Meta-Llama-3_3-70B-Instruct`, and more. See the
[AI Endpoints catalogue](https://endpoints.ai.cloud.ovh.net/) for current model
IDs.

## Together AI

[Together AI](https://www.together.ai/) is one of the largest hosts of
open-weight models, serving Llama, Qwen, DeepSeek, Kimi, GLM, and others
through an OpenAI-compatible API.

1. Create an API key from the [Together AI settings](https://api.together.ai/settings/api-keys).
2. Set the environment variable:

```console
$ export TOGETHER_API_KEY=your_key_here
```

Use Together AI models in your agent configuration:

```yaml
agents:
  root:
    model: together/meta-llama/Llama-3.3-70B-Instruct-Turbo
    instruction: You are a helpful assistant
```

See the [Together AI model library](https://docs.together.ai/docs/serverless-models)
for current model IDs.

## Vercel AI Gateway

[Vercel AI Gateway](https://vercel.com/docs/ai-gateway) is a unified
OpenAI-compatible endpoint that routes to models from OpenAI, Anthropic, Google,
xAI, and more at list price with no markup.

1. Create an API key from the [Vercel AI Gateway dashboard](https://vercel.com/docs/ai-gateway).
2. Set the environment variable:

```console
$ export AI_GATEWAY_API_KEY=your_key_here
```

Use Vercel AI Gateway in your agent configuration:

```yaml
agents:
  root:
    model: vercel/openai/gpt-5
    instruction: You are a helpful assistant
```

Vercel AI Gateway model IDs use the `creator/model` form (for example
`openai/gpt-5` or `anthropic/claude-sonnet-4.5`). See the
[Vercel AI Gateway documentation](https://vercel.com/docs/ai-gateway) for the
current model list.

## OpenAI-compatible providers

You can use the `openai` provider type to connect to any model or provider that
implements the OpenAI API specification. This includes services like Azure
OpenAI, local inference servers, and other compatible endpoints.

Define a named model in the `models` section with the `openai` provider and a
`base_url`, then reference it from your agent:

```yaml
models:
  my-model:
    provider: openai
    model: your-model-name
    base_url: https://your-provider.example.com/v1

agents:
  root:
    model: my-model
    instruction: You are a helpful coding assistant
```

By default, Docker Agent uses the `OPENAI_API_KEY` environment variable for
authentication. If your provider uses a different variable, specify it with
`token_key`:

```yaml
models:
  my-model:
    provider: openai
    model: your-model-name
    base_url: https://your-provider.example.com/v1
    token_key: YOUR_PROVIDER_API_KEY

agents:
  root:
    model: my-model
    instruction: You are a helpful coding assistant
```

## What's next

- Follow the [tutorial](tutorial.md) to build your first agent
- Learn about [local models with Docker Model Runner](local-models.md) as an
  alternative to cloud providers
- Review the [configuration reference](reference/config.md) for advanced model
  settings
