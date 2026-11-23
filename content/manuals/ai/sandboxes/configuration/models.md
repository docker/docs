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
`--model` isn't supported with cloud sandboxes.

> [!NOTE]
> Model selection is experimental. Enable it before following these examples.

## Bundled model service

Docker Sandboxes includes `llmman`, a model management tool installed in the
same directory as `sbx`. It serves local models and forwards requests to hosted
providers or custom endpoints. You don't need to install it separately. To run
it yourself, see [Install Docker Sandboxes](../install.md#model-service-binary).

By default, the Docker Sandboxes daemon starts `llmman` on your host and stops
it when the daemon stops. The service starts even if you haven't used
`--model`. It keeps running after a sandbox or CLI exits, and sandboxes share
its model store and loaded models.

The daemon runs `llmman serve` bound to `127.0.0.1:17434`, with its shell route
turned off. If an `llmman` of a different version already holds that port, such
as one left behind by an older install, the daemon stops it first. Starting the
service doesn't pull model files, load a model, or run inference. Only a
`--model` or `--provider` request from `sbx run` or `sbx create` asks the
service for a model.

`llmman` inherits the daemon's environment, so provider API-key variables set
there are visible to it. See [Provider authentication](#provider-authentication).
`sbx` has no setting to limit the service's CPU, memory, or disk use. A model
runs on the host, outside the sandbox's resource limits.

To prevent all of this, [disable the bundled service](#disable-the-bundled-service).

### Where llama.cpp comes from

`llmman` gets the llama.cpp inference server itself. `sbx` downloads and
installs nothing. On Linux, `llmman` first tries a container runtime on your
host, then a download, then `llama-server` on your `PATH`. On Windows and
macOS, it tries the download, then `llama-server` on your `PATH`.

For the download, `llmman` detects your hardware and picks the llama.cpp build
for it. To see its choice, run `llmman gpu-discover`.

| Platform | Builds |
| --- | --- |
| Windows | CPU, Vulkan, CUDA 12 or 13, ROCm, and OpenCL for Adreno on Arm |
| Linux | CPU, Vulkan, ROCm, and CUDA 12 or 13, including CUDA 13 on Arm (NVIDIA DGX Spark) |
| macOS | Metal |

`llmman` downloads the build from Docker's mirror in the public
[`docker/sbx-releases`](https://github.com/docker/sbx-releases) repository, from
the release named `llama.cpp-<BUILD>`. The Docker Sandboxes release pipeline
publishes the mirror after it code signs the Windows and macOS binaries. Linux
builds are copied unchanged. Upstream Windows builds aren't signed, and Windows
Smart App Control, Windows Defender Application Control (WDAC), and AppLocker
can block them.

`llmman` gets llama.cpp when it starts, before it listens for requests. The
first `sbx run --model` waits for it, up to 15 minutes. If the download fails,
`llmman` still starts and logs a warning. The first model load that needs
llama.cpp tries again and reports the error if it fails again.

The download works as follows:

- Progress: `llmman` prints it to `llmman/serve.log`. See
  [Find logs and status](#find-logs-and-status).
- Verification: `llmman` checks each file against the digest that GitHub lists
  for it. It never uses a partial install.
- Location: `~/.local/share/llmman/llama.cpp/` on Linux and macOS, and
  `%LOCALAPPDATA%\llmman\llama.cpp\` on Windows. An install is reused without
  network access. When `llmman` installs another build, it removes the older ones.
- Network: allow `api.github.com`, `github.com`, and `*.githubusercontent.com`.
  The download is anonymous and uses no Docker or provider credentials. It
  honors the standard proxy variables, such as `HTTPS_PROXY` and `NO_PROXY`, in
  the daemon's environment.
- Size: about 15 to 40 MB for CPU, Vulkan, Metal, and OpenCL builds, about
  280 MB for ROCm, and 170 to 320 MB for CUDA, plus a CUDA runtime bundle of
  about 280 to 600 MB.

### Use another llama.cpp source

To make `llmman` download from another GitHub repository, set
`LLMMAN_LLAMA_CPP_REPO` to `<OWNER>/<REPO>` in the daemon's environment, then
restart the daemon:

```console
$ LLMMAN_LLAMA_CPP_REPO=ggml-org/llama.cpp sbx daemon restart
```

With `ggml-org/llama.cpp`, `llmman` downloads the same build straight from
upstream, without Docker's signing. On Windows, a policy that blocks unsigned
code blocks it. A repository that isn't upstream must publish each build as a
release named `llama.cpp-<BUILD>` with the upstream file names. For details,
see the
[llmman documentation](https://github.com/llmmanorg/llmman/blob/main/docs/backends.md#mirroring-the-llamacpp-release).

### Disable the bundled service

If you don't use the bundled model service, you can turn it off to skip its
llama.cpp download:

```console
$ sbx settings set model.server.enabled false
$ sbx daemon restart
```

Restart the daemon to apply the change. If the daemon wasn't running, the
settings command starts it first, so the service might start and begin its
download before you restart.

With the service turned off, the daemon doesn't run `llmman` or download
llama.cpp. It only watches port `17434` and stops a leftover `llmman` that has
its shell route enabled. A later `sbx run --model` reports that the local model
server is disabled.

You can still [use Ollama](#use-ollama) with the bundled service turned off.
For local models managed by `llmman`, hosted providers, and custom endpoints,
turn the service back on:

```console
$ sbx settings set model.server.enabled true
$ sbx daemon restart
```

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

`sbx run` and `sbx create` reject a model flag with a blank value, such as
`--model "$MODEL"` when `MODEL` is unset.

## Run a local model

Pass a GGUF model reference or short name to `--model`:

```console
$ sbx run --model gemma4 claude
```

Docker Sandboxes downloads the model if needed. The model runs on the host,
so its memory and compute requirements are separate from the sandbox's
resource limits. Replace `claude` with `codex` or `opencode` to use another
agent with the same model.

The download shows the same progress bar as `llmman pull`, because `sbx` runs
the `llmman` installed next to it:

```text
gemma4               [======>               ]   1.02 GiB/3.10 GiB    48.31 MiB/s
```

If the model files are already in the store under another name, `llmman`
prints `Already have <MODEL>` instead of a bar. If `llmman` can't run on the
host, `sbx` prints plain status lines instead.

The first time the service starts, `sbx run --model` also waits for llama.cpp.
See [Troubleshoot the model service](#troubleshoot-the-model-service).

### Use Ollama

Install and start Ollama on your host, then select it with `--provider`:

```console
$ sbx run --model gemma4 --provider ollama claude
```

Docker Sandboxes connects to Ollama at `localhost:11434`. It doesn't install,
start, or manage the Ollama process. If the model isn't in Ollama yet, `sbx`
pulls it through Ollama's API and shows the same `llmman` progress bar.

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
Hosted providers need `--model`.

If the provider has no credential, `sbx` prints a warning and continues.
Requests to the provider fail until you store a key. See
[Hosted provider has no credentials](#hosted-provider-has-no-credentials).

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

## Expose several models

To let the agent choose between models, expose more than one:

- Pass `--provider` without `--model`. For `ollama` and for a provider defined
  in `model.providers`, the agent gets every chat model the provider lists.
  The first listed model is the default.
- Repeat `--model`. The first value is the default. Duplicates are ignored.

```console
$ sbx run --provider company opencode
$ sbx run --model gemma4 --model qwen3.5:0.8b opencode
```

OpenCode lists every exposed model in its model picker. Claude Code and Codex
start with the default model.

`sbx` leaves out models that can't chat. For Ollama, it uses the capabilities
that Ollama reports, or the model's family and name if it reports none. For a
custom endpoint, it leaves out IDs that contain `embed`, `rerank`, `whisper`,
or `tts`. At most 64 models are exposed. If a provider lists more, `sbx`
exposes the first 64 and prints a warning. Repeating `--model` more than 64
times is an error.

Some combinations aren't accepted:

- `--provider` without `--model` for a hosted provider, or for `llmman`. These
  catalogs are too large to expose, so pass `--model`.
- `--variant`, `--overflow-provider`, and `--overflow-model` with more than one
  model. They apply to one model.

## Create a sandbox with a model

`sbx create` accepts `--model`, `--provider`, `--overflow-provider`,
`--overflow-model`, and `--variant`, with the same meaning as for `sbx run`:

```console
$ sbx create --name my-sandbox --model gemma4 claude .
$ sbx run --name my-sandbox
```

`sbx create` downloads the model and builds the model configuration into the
sandbox, as `sbx run --model` does. You can adjust the sandbox's network rules
before you start it.

## Change an existing sandbox's model

Pass the sandbox name and your model selection:

```console
$ sbx run --name <SANDBOX_NAME> --model <MODEL_NAME> --provider <PROVIDER_ID>
```

Changing the model recreates the sandbox container. The workspace and
kit-owned volumes persist. Omit `--provider` to select a local model managed
by `llmman`.

Before `sbx` downloads the model or touches the sandbox, it checks which
additional permissions the change needs. It grants the model server's port
(`localhost:17434`) automatically when the model stays on your host, such as a
local model or a `model.providers` endpoint on a loopback address, and prints
the permission. Any other permission, such as Ollama's `localhost:11434` or a
hosted provider, asks for approval. The prompt needs a terminal and fails with
`--no-input`. If you decline, nothing is downloaded and the sandbox isn't
changed. To avoid the prompt in scripts, select the model with
[`sbx create`](#create-a-sandbox-with-a-model) instead.

## Choose a reasoning level

Use `--variant` with `--model` to select a reasoning level supported by your
model. For example, replace `<MODEL_NAME>` with an OpenAI model that supports
`high` reasoning effort to start Codex at that level:

```console
$ sbx run --provider openai --model <MODEL_NAME> --variant high codex
```

Available values depend on the model and agent. Reasoning levels are `minimal`,
`low`, `medium`, `high`, `xhigh`, and `max`. Some models support `none` to disable
thinking or `thinking` to enable it. Omit `--variant` to use the default.

`--variant` works with Claude Code, Codex, and OpenCode, on both `sbx run` and
`sbx create`. The sandbox keeps the variant, so later `sbx run` commands start
at it. To change it, pass `--model` with another `--variant` to `sbx run`.
Claude Code doesn't support `none` or `minimal`. For Claude Code and Codex,
`--variant` requires an attached session and can't be combined with `--detach`.

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

## Troubleshoot the model service

### Find logs and status

The model service writes to the sandbox state directory:

| Platform | State directory |
| --- | --- |
| macOS | `~/Library/Application Support/com.docker.sandboxes/sandboxes/sandboxd` |
| Linux | `~/.local/state/sandboxes/sandboxes/sandboxd` |
| Windows | `%LOCALAPPDATA%\DockerSandboxes\sandboxes\state\sandboxd` |

On Linux, `sbx` uses `$XDG_STATE_HOME/sandboxes/sandboxes/sandboxd` when
`XDG_STATE_HOME` is set and you run it as an ordinary user.

| Path in the state directory | Contents |
| --- | --- |
| `llmman/serve.log` | Output of `llmman serve`. Earlier runs rotate to `serve-1.log` through `serve-5.log`. |
| `llmman-status.json` | What the daemon last did with the service. |

The `state` field of `llmman-status.json` is one of:

| State | Meaning |
| --- | --- |
| `failed` | The daemon couldn't run the service. `lastError` says why. |
| `started`, `running` | `llmman serve` was started, or answers on port `17434`. |
| `exited` | `llmman serve` exited. The daemon restarts it. `exits` counts consecutive early exits. |
| `standby` | Another `llmman` holds port `17434`. The daemon waits for it to exit. |
| `disabled` | `model.server.enabled` is `false`. |

For the daemon's own log and state, use:

```console
$ sbx daemon status
$ sbx daemon logs --tail 100
$ sbx daemon restart
```

### See which llmman is used

`sbx` runs only the `llmman` in the same directory as `sbx`, after following
symlinks. An `llmman` elsewhere on your `PATH` isn't used by the daemon. To
see what the running service is, query it:

```console
$ curl -s http://127.0.0.1:17434/api/version
```

The response has the service's `version`, the path of its executable in
`exe`, and its process ID in `pid`. On Windows, use `curl.exe`.

### Model service doesn't start

If `sbx run --model` keeps waiting, or fails with one of these errors, check
the status file first:

- `the local model server is disabled`: `model.server.enabled` is `false`. Turn
  it on and restart the daemon. See [Disable the bundled service](#disable-the-bundled-service).
- `the local model server keeps exiting`: `llmman serve` exits right after it
  starts. Read `llmman/serve.log` for the cause. If `llmman` can't run on the
  host at all, see [llmman can't run on the host](#llmman-cant-run-on-the-host).
- `locate llmman binary`: there is no `llmman` next to `sbx`. Reinstall `sbx`,
  or use `--provider ollama`. See [Install Docker Sandboxes](../install.md#model-service-binary).
- `did not become ready within 15m0s`: the service didn't answer in 15
  minutes. See [First start is slow or blocked](#first-start-is-slow-or-blocked).

Then run `sbx daemon restart` and try again.

### First start is slow or blocked

`llmman` gets llama.cpp before it listens, so the first `sbx run --model` waits
for the download. While it waits, `sbx` prints the path of `llmman/serve.log`,
where `llmman` prints the progress. A firewall or proxy that blocks the
download is the usual cause of a failure. To fix it:

- Allow the hosts listed in [Where llama.cpp comes from](#where-llamacpp-comes-from).
  They are in addition to the
  [domains `sbx` needs](../faq.md#which-domains-do-i-need-to-allow-for-docker-sandboxes-to-work).
- Set `HTTPS_PROXY` in the shell that starts the daemon, then run
  `sbx daemon restart`.

After the download works once, later starts need no network access for
llama.cpp.

### Smart App Control or WDAC blocks a file

Smart App Control, WDAC, and AppLocker block unsigned code. `sbx`, `llmman`, and
the llama.cpp build from Docker's mirror are signed. If a policy still blocks
the service, check the Windows event logs for the blocked file. Report the
problem with `sbx diagnose --upload`. See
[Report an issue](../troubleshooting.md#report-an-issue).

If `LLMMAN_LLAMA_CPP_REPO` is set to `ggml-org/llama.cpp` in the daemon's
environment, `llmman` downloads unsigned builds that these policies block. Unset
it and run `sbx daemon restart`.

### Ollama isn't running or the model is missing

With `--provider ollama`, `sbx` doesn't start Ollama. If Ollama doesn't answer
at `http://127.0.0.1:11434`, `sbx run` fails with
`Ollama does not appear to be running`. Install and start Ollama, then run the
command again. To use the bundled service instead, drop `--provider ollama`.

If the model isn't in Ollama, `sbx` pulls it. A failed pull reports
`pull <MODEL> via Ollama`, followed by Ollama's error. Check the model name
with `ollama list` or `ollama pull <MODEL>`.

For `--provider ollama` without `--model`, the error `lists no chat models to
expose` means Ollama has no model that can chat. Pull one first.

### Hosted provider has no credentials

If `sbx run` warns `no API key for <PROVIDER>`, no stored secret or
environment variable provides a key. The run continues, but requests to the
provider fail. Store the key and run again:

```console
$ sbx secret set <PROVIDER_ID>
```

If `sbx` reports `invalid --provider`, the provider isn't one `llmman` can
route to. Run `llmman providers` for the list, or define the endpoint in
[`model.providers`](#connect-a-custom-endpoint).

If you set the key as an environment variable instead, restart the daemon from
the shell that has it. See [Provider authentication](#provider-authentication).

### Sandbox can't reach the service

The model configuration allows the sandbox to reach the model service at port
`17434`, or `11434` for Ollama, on the host. The agent connects through
`host.docker.internal`. If requests fail inside the sandbox:

- Check that the service answers on the host with
  `curl -s http://127.0.0.1:17434/api/version`.
- Check `sbx policy log` for blocked requests to `localhost:17434`.
  Organization governance can remove the permission.
- Check that another program doesn't hold port `17434`. The daemon waits when
  another `llmman` holds the port, and the status file shows `standby`.

For other reasons a sandbox can't reach a host service, see
[Can't reach a service running on the host](../troubleshooting.md#cant-reach-a-service-running-on-the-host).

### llmman can't run on the host

If the system lacks a library `llmman` needs, `sbx run --model` fails with
`the bundled llmman cannot run on this host`, or the service keeps exiting.
`llmman/serve.log` names the cause. On Linux, `llmman` needs glibc 2.28 or
newer.

Use Ollama instead. It doesn't need `llmman`:

```console
$ sbx run --provider ollama --model <MODEL_NAME> claude
```
