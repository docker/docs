---
title: Use local and hosted models
linkTitle: Models
weight: 15
description: Choose a local model, hosted provider, or custom inference endpoint for Claude Code, Codex, or OpenCode in Docker Sandboxes.
keywords: docker sandboxes, sbx, models, providers, llmman, ollama, inference, local models
---

Use `sbx run --model` to choose the model and service that answer your agent's
requests. The agent runs inside a local sandbox. The model can run on your
host, at a hosted provider, or at an inference endpoint you configure.

This page covers the built-in `claude`, `codex`, and `opencode` agents. Choose
a model that supports the tool calls and context length your agent needs.
`--model` isn't supported with cloud sandboxes or v3 kits.

> [!NOTE]
> Model selection is experimental. Enable it before following these examples.

## Bundled model service

Docker Sandboxes includes `llmman`, a model management tool installed alongside
`sbx`. It serves local models and forwards requests to hosted providers or
custom endpoints. You don't need to install it separately.

The Docker Sandboxes daemon starts `llmman` on your host and stops it when
the daemon stops. At startup, `llmman` downloads its llama.cpp runtime if
needed, so the first startup requires network access for that download. The service starts even if you haven't used `--model`.

The service keeps running after a sandbox or CLI exits. Sandboxes share its
model store and loaded models.

On Linux, `llmman` first tries Docker or Podman to run the inference server
in a container. If neither is usable, it tries a downloaded binary, then
`llama-server` on the host's `PATH`.

## Enable model selection

Run these commands on your host:

```console
$ sbx settings set platform.allowExperimentalFeatures true
$ sbx settings set feature.model true
```

The `--provider` flag selects where the model runs:

| Provider | Model destination |
| --- | --- |
| Omitted, or `llmman` | A local model managed by `llmman` |
| `ollama` | An existing Ollama installation on your host |
| A hosted provider ID | A provider supported by `llmman`, such as `openai` or `anthropic` |
| An ID from `model.providers` | An endpoint you configure |

## Run a local model

Pass a GGUF model reference or short name to `--model`:

```console
$ sbx run --model gemma4 claude
```

Docker Sandboxes downloads the model if needed. The model runs on the host,
so its memory and compute requirements are separate from the sandbox's
resource limits. Replace `claude` with `codex` or `opencode` to use another
agent with the same model.

### Use Ollama

Install and start Ollama on your host, then select it with `--provider`:

```console
$ sbx run --model gemma4 --provider ollama claude
```

Docker Sandboxes connects to Ollama at `localhost:11434`. It doesn't install,
start, or manage the Ollama process.

For Docker Model Runner, see
[Run Claude Code in a Docker Sandbox with Docker Model Runner](/guides/claude-code-sandbox-model-runner/).

## Use a hosted provider

Select a provider supported by `llmman` and a model available from that
provider. For example, store your OpenAI API key on the host, then run Codex
with an OpenAI model:

```console
$ sbx secret set openai
$ sbx run --provider openai --model gpt-5-nano codex
```

The host's `llmman` service forwards requests to the provider. It doesn't
download or run the hosted model. For available providers and their API-key
variable names, see the [llmman provider documentation](https://github.com/llmmanorg/llmman/blob/main/docs/providers.md).

### Provider authentication

Store the API key with `sbx secret set <PROVIDER_ID>`, using the same ID as
`--provider` or `--overflow-provider`. This also works for a custom provider
ID defined in `model.providers`. Sandbox-scoped secrets take precedence over
global secrets.

The sandbox sends a placeholder in place of the key. The host proxy replaces
it with the stored key only in authentication headers on requests to the host's
`llmman` service. That service uses the key to authenticate requests to the
selected provider. The real key stays outside the sandbox.

Docker Sandboxes configures this injection automatically when you select the
provider. You don't need `sbx secret set-custom` or domain injection rules for
this workflow.

> [!WARNING]
> Model-provider credentials are protected differently from other sandbox
> credentials. For other credentials, the proxy restricts injection to their
> configured domains. For model providers, it supplies the key to the host's
> `llmman` service instead. That service can forward the key to other endpoints
> it supports; selecting `--provider` does not restrict the key to that
> provider's endpoint.

If no stored secret is available, make the key available in the host shell
that starts the daemon. For a custom endpoint, choose the variable name with
[`apiKeyEnv`](#connect-a-custom-endpoint).

`llmman` inherits the daemon's environment. After setting or changing a key
in that environment, run `sbx daemon restart` from the shell containing the
updated variable.
This restarts the model service and interrupts model requests from sandboxes
using it. Setting a variable only in the shell where you run `sbx run --model`
doesn't update an already-running service.

Using a stored secret avoids restarting the daemon to load a key from its
environment. For the agents' default authentication flows, see
[Manage credentials](credentials.md).

## Connect a custom endpoint

Use `model.providers` to connect to an OpenAI- or Anthropic-compatible
inference endpoint, such as an internal GPU server. The endpoint must be
reachable from your host.

The setting is a JSON object keyed by provider ID. Check its existing value
before changing it:

```console
$ sbx settings get model.providers
```

For an OpenAI-compatible endpoint, define a provider named `company`:

```console
$ sbx settings set model.providers '{"company":{"url":"https://inference.example.com/v1","wire":"openai"}}'
```

Replace the URL with your endpoint's base URL. Setting `model.providers`
replaces the whole object, so include any existing providers you want to keep.

| Field | Description |
| --- | --- |
| `url` | Required HTTP or HTTPS base URL, usually ending in `/v1`. Use the base URL, without `/chat/completions` or `/messages`. |
| `wire` | The endpoint's API format: `openai` (default) or `anthropic`. This describes the endpoint, regardless of which agent you run. |
| `apiKeyEnv` | Name of an optional host environment variable containing the API key. Used as a fallback when no stored provider secret is available. |
| `name` | Optional display name. Defaults to the provider ID. |

If the endpoint requires an API key, store it using the provider ID
(`company`). Docker Sandboxes configures authentication when you select that
provider with `--provider`:

```console
$ sbx secret set company
$ sbx run --provider company --model <MODEL_NAME> claude
```

Docker Sandboxes applies the provider configuration when you run with
`--model`. You can use the same provider with `codex` or `opencode`.

## Change an existing sandbox's model

Pass the sandbox name and your model selection:

```console
$ sbx run --name <SANDBOX_NAME> --model <MODEL_NAME> --provider <PROVIDER_ID>
```

Changing the model recreates the sandbox container. The workspace and
kit-owned volumes persist. Omit `--provider` to select a local model managed
by `llmman`.

## Use another provider for larger requests

Pair a local model with another provider to handle requests that exceed the
local model's context capacity:

```console
$ sbx run --model gemma4 \
    --overflow-provider openai --overflow-model gpt-5-nano claude
```

Configure [provider authentication](#provider-authentication) before starting
the model service. You can also use a provider defined in `model.providers`.
Both overflow flags are required, and the local model must use the default
`llmman` provider. This option can't be combined with `--provider ollama` or
a hosted provider selected with `--provider`.

Requests that fit the local model stay local. Requests routed to the overflow
provider send their contents to that endpoint and can incur provider charges.
