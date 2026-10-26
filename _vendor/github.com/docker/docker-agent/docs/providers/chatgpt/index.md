---
title: "ChatGPT (OpenAI account)"
description: "Use your ChatGPT Plus/Pro/Business subscription with Docker Agent by signing in with your OpenAI account, no API key needed."
keywords: docker agent, ai agents, model providers, llm, chatgpt, openai, codex, subscription
weight: 55
canonical: https://docs.docker.com/ai/docker-agent/providers/chatgpt/
---

_Use your ChatGPT subscription with Docker Agent by signing in with your OpenAI account. No API key needed._

## Overview

The `chatgpt` provider authenticates with a ChatGPT account (the same
"Sign in with ChatGPT" flow used by OpenAI's Codex CLI) instead of an
`OPENAI_API_KEY`. Usage is billed against your ChatGPT Plus, Pro, or
Business plan rather than pay-per-token API credits.

Under the hood, Docker Agent talks to the ChatGPT Codex backend
(`https://chatgpt.com/backend-api/codex`), which serves the `gpt-5` model
family over the OpenAI Responses API. GPT-5.6 (Sol/Terra/Luna) is served
there too; GPT-5.2 and GPT-5.3-Codex are deprecated for ChatGPT sign-in.

## Prerequisites

- A paid **ChatGPT** subscription (Plus, Pro, or Business).
- A browser on the machine running the sign-in (the OAuth flow uses a
  fixed `localhost:1455` callback).

## Sign In

```bash
docker agent setup
```

Pick **chatgpt** in the provider list: instead of asking for an API key, the
wizard opens your browser on the ChatGPT sign-in page and stores the
resulting OAuth credential in the Docker Agent config directory
(`~/.config/cagent/chatgpt-auth.json`, owner-only permissions). The access
token is refreshed automatically; you only need to sign in again if the
refresh token is revoked.

Related commands:

```bash
docker agent doctor                        # the chatgpt row shows the credential state
rm ~/.config/cagent/chatgpt-auth.json      # sign out (remove the stored sign-in)
```

## Configuration

### Inline

```yaml
agents:
  root:
    model: chatgpt/gpt-5.6
    instruction: You are a helpful assistant.
```

### Named model

```yaml
models:
  gpt:
    provider: chatgpt
    model: gpt-5.6
    thinking_budget: medium

agents:
  root:
    model: gpt
```

## Available Models

The Codex backend serves the models available to your ChatGPT plan,
typically:

| Model                | Best For                              |
| -------------------- | -------------------------------------- |
| `gpt-5.6`            | Alias for `gpt-5.6-sol`; general purpose, strong reasoning |
| `gpt-5.6-sol`        | Frontier model, most capable          |
| `gpt-5.6-terra`      | Everyday workhorse                    |
| `gpt-5.6-luna`       | High-volume, cost-efficient           |
| `gpt-5.2`            | Deprecated for ChatGPT sign-in        |
| `gpt-5.2-codex`      | Deprecated for ChatGPT sign-in        |

The effort picker exposes Low/Medium/High/XHigh/Max on the GPT-5.6 family
(no Minimal).

## How It Works

- **Auth:** the `docker agent setup` sign-in runs an OAuth 2.0
  authorization-code + PKCE flow against `auth.openai.com`. The stored login
  is exposed to credential checks (doctor, `first_available`, auto model
  selection) as the virtual `CHATGPT_OAUTH_TOKEN` variable.
- **API:** requests go to the Responses API only; the backend has no Chat
  Completions endpoint, so `api_type` is pinned automatically.
- **Request shape:** the backend requires stateless requests (`store: false`)
  and a top-level `instructions` field, so Docker Agent moves system messages
  there. Client-side sampling parameters (`temperature`, `top_p`,
  `max_tokens`) are not supported by the backend and are dropped.

## Image Input

For image-capable models such as `gpt-6.1-sol`, you can attach an image:

```bash
docker agent run examples/chatgpt-images.yaml --exec \
  --attach examples/chatgpt-images.png "Describe the shapes and colors."
```

The sample includes a red square centered on a blue background. Running it
requires ChatGPT sign-in and access to the selected model. Image input for
`gpt-6.1-sol` is inferred from the matching OpenAI catalogue entry; no
`capabilities` override is needed.

When models.dev has no direct `chatgpt/<model>` entry, Docker Agent uses the
matching `openai/<model>` entry for **image input only**. This applies both to
`--attach` images and image content returned by MCP tools (for example,
screenshots). Tool images are sent in a follow-up user message alongside the
text function-call output. PDF, audio, and video support, prices, context limits,
tool support, output capabilities, and `max_tokens` are not inherited.
A direct ChatGPT catalogue entry takes precedence; unknown models remain
text-only.

To intentionally override catalogue capabilities or declare support for an
unknown model, see [Attachment Capability Overrides](../../configuration/models/index.md#attachment-capability-overrides).
Image input support does not enable image generation or editing.

## Setting the Token Explicitly

`CHATGPT_OAUTH_TOKEN` can also be set like any other credential (shell
environment, `--env-from-file`, ...). An explicitly set value takes
precedence over the stored sign-in. This is useful for short-lived CI runs
with a pre-minted access token, but note that such a token expires and is not
refreshed.

## ChatGPT Subscription vs. OpenAI API Key

| | `chatgpt` | `openai` |
| --- | --- | --- |
| Credential | ChatGPT account sign-in | `OPENAI_API_KEY` |
| Billing | Included in the ChatGPT plan (rate-limited) | Pay per token |
| Models | `gpt-5` family served by the Codex backend | Full OpenAI API catalog |
| Sampling controls (`temperature`, ...) | Not supported | Supported |
| Embeddings / reranking | Not supported | Supported |

When both credentials are configured, automatic model selection prefers
`openai`; pin `--model chatgpt/gpt-5.6` (or use a named model) to use the
subscription.

> [!NOTE]
> Use of the Codex backend is governed by OpenAI's terms for ChatGPT and
> Codex. Sign-in is per user; do not share the stored credential.
