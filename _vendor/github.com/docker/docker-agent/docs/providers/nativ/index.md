---
title: "Nativ"
description: "Run Docker Agent with local MLX models served by Nativ on Apple Silicon."
keywords: docker agent, ai agents, model providers, local models, nativ, mlx, apple silicon
weight: 185
canonical: https://docs.docker.com/ai/docker-agent/providers/nativ/
---

_Run Docker Agent with local MLX models served by Nativ on Apple Silicon._

## Overview

[Nativ](https://github.com/Blaizzy/nativ) is a macOS app for downloading and
serving MLX models locally. Docker Agent connects to its OpenAI-compatible API
through a [provider definition](../custom/index.md). No built-in `nativ` alias
or additional provider plugin is required.

Nativ requires Apple Silicon and macOS 26 or newer. Model downloads require
network access; inference runs locally after the model is downloaded.

## Setup

1. Install Nativ from its [releases](https://github.com/Blaizzy/nativ/releases/latest)
   or with Homebrew:

   ```bash
   $ brew install --cask nativ
   ```

2. Launch Nativ and complete its initial setup. In **Models**, download and
   select a compatible language model. For a small first download, use
   [`mlx-community/Qwen3-0.6B-4bit`](https://huggingface.co/mlx-community/Qwen3-0.6B-4bit),
   approximately **351 MB** including tokenizer files.
3. Start the server from Nativ's **Developer** page or menu-bar controls. Keep
   Nativ running while you use Docker Agent. The default OpenAI base URL is
   `http://127.0.0.1:8080/v1`; the Developer page shows the configured host,
   port, endpoints, and logs.
4. Check that the server responds:

   ```bash
   $ curl http://127.0.0.1:8080/health
   $ curl http://127.0.0.1:8080/v1/models
   ```

If you enabled server authentication, include the configured Bearer token in
requests. See [Authentication](#authentication) below.

## Configuration

Save this as `nativ.yaml`:

```yaml
providers:
  nativ:
    api_type: openai_chatcompletions
    base_url: http://127.0.0.1:8080/v1
    provider_opts:
      extra_body:
        enable_thinking: false

models:
  tiny:
    provider: nativ
    model: mlx-community/Qwen3-0.6B-4bit
    temperature: 0
    max_tokens: 256

agents:
  root:
    model: tiny
    description: Local chat assistant using Nativ on Apple Silicon
    instruction: You are a concise helpful assistant. Answer directly.
```

Run a short chat request:

```bash
$ docker agent run nativ.yaml --exec "What is the capital of France? Answer with only the city name."
```

A copy is available in
[`examples/nativ.yaml`](https://github.com/docker/docker-agent/blob/main/examples/nativ.yaml).
Use the repository ID of your downloaded model in `models.tiny.model` if you
choose a different model. Update `base_url` if you change Nativ's port.

`provider_opts.extra_body.enable_thinking: false` sends Nativ's request-level
reasoning control to the server. It avoids spending the small output budget
on reasoning. See [Extra Request Body](../../configuration/models/index.md#extra-request-body)
for other pass-through settings.

Nativ also exposes the OpenAI Responses API. To use it, change
`providers.nativ.api_type` to `openai_responses`; keep the same base URL.
Docker Agent does not forward `provider_opts.extra_body` with Responses
requests, so the thinking override above applies only to Chat Completions.
Use Chat Completions when you need that request-level control; with Responses,
check Nativ's server defaults and allow a larger output budget if thinking is
enabled.

## Authentication

If you enable a server API key in Nativ, add `token_key: NATIV_API_KEY` under
`providers.nativ` and set that environment variable to the same key:

```bash
$ export NATIV_API_KEY=your-server-api-key
$ curl http://127.0.0.1:8080/v1/models -H "Authorization: Bearer $NATIV_API_KEY"
$ docker agent run nativ.yaml --exec "Hello"
```

`token_key` names an environment variable, not a literal token. Omit it when
server authentication is disabled. Nativ's server key is separate from a
Hugging Face token used to download gated models.

## Tool Calling

The example is chat-only. Nativ's API supports tool calls, but reliability
depends on the model, instructions, and tool schemas. A 0.6B model is useful
for testing the connection, not a dependable coding agent: it can ignore a
tool or give an incorrect answer even when instructed to use one.

Before adding filesystem or shell tools, choose a tool-capable model and test
that it calls the tool with correct arguments and uses the returned result.
Keep tool descriptions and parameter schemas clear, and retain tool approvals.
See [Tools](../../concepts/tools/index.md) for configuration.

## Troubleshooting

- **Connection refused:** Start Nativ's server and check the host and port in
  its Developer page. Installing the app alone does not make the API available.
- **Model not found:** Confirm that the model is downloaded and use its exact
  repository ID. Check `/v1/models` and Nativ's server logs.
- **Unauthorized:** Match `NATIV_API_KEY` to Nativ's configured server key, or
  remove `token_key` when authentication is disabled.
- **No answer before the token limit:** Keep thinking disabled for the tiny
  model, or increase `max_tokens` for a model that needs longer responses.
- **Docker Agent runs in a container:** Container-local `127.0.0.1` does not
  refer to your Mac. On Docker Desktop, use
  `http://host.docker.internal:8080/v1` and ensure Nativ listens on an interface
  reachable from the container. Enable authentication before exposing the
  server beyond loopback; changing its bind address can expose it to your LAN.

For server controls and endpoint details, see Nativ's
[Developer documentation](https://github.com/Blaizzy/nativ/blob/main/Docs/features/developer.md).
