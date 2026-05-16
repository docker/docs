---
title: Run Claude Code in a Docker Sandbox with Docker Model Runner
description: Run Claude Code inside an isolated Docker Sandbox and route requests to Docker Model Runner so the agent uses local models on your host.
summary: |
  Combine Docker Sandboxes with Docker Model Runner to run Claude Code in an
  isolated microVM that talks to a local model on your host through the
  Anthropic-compatible API.
keywords: ai, claude code, docker model runner, docker sandboxes, sbx, anthropic, local models, coding assistant
tags: [ai]
params:
  time: 15 minutes
---

This guide shows how to run Claude Code inside a Docker Sandbox with Docker
Model Runner as the backend model provider. You'll keep the agent isolated
from your host in a microVM, point it at a local model on your machine, and
keep all model traffic on-device.

> **Acknowledgment**
>
> Docker would like to thank [Pradumna Saraf](https://twitter.com/pradumna_saraf) for his contribution to this guide.

In this guide, you'll learn how to:

- Pull a coding model and start Docker Model Runner with TCP enabled
- Allow the sandbox to reach Docker Model Runner on your host
- Create a Claude Code sandbox and set the local endpoint persistently
- Launch Claude Code with a local model and verify the connection
- Package `gpt-oss` with a larger context window for longer prompts

## How the pieces fit together

Three components cooperate at runtime:

- **Docker Model Runner** runs on your host and serves an
  Anthropic-compatible API at `http://localhost:12434`.
- **The Docker Sandbox** runs Claude Code inside an isolated microVM. The
  microVM has its own network and can't reach your host's `localhost`
  directly.
- **The sandbox proxy** sits on your host and brokers every outbound
  request from the sandbox. It enforces network policy and translates the
  special hostname `host.docker.internal` to `localhost`.

Claude Code inside the sandbox sends requests to
`http://host.docker.internal:12434`. The proxy rewrites the destination to
`localhost:12434`, which Docker Model Runner answers. No model traffic
leaves your machine.

## Prerequisites

Before you start, make sure you have:

- [Docker Desktop](../get-started/get-docker.md) or Docker Engine installed
- [Docker Model Runner enabled](../manuals/ai/model-runner/get-started.md#enable-docker-model-runner)
- [Docker Sandboxes (`sbx`) installed and signed in](../manuals/ai/sandboxes/get-started.md#install-and-sign-in)

If you use Docker Desktop, turn on TCP access in **Settings** > **AI**, or
run:

```console
$ docker desktop enable model-runner --tcp 12434
```

## Step 1: Pull a coding model

Pull a model on your host before you create the sandbox:

```console
$ docker model pull ai/devstral-small-2
```

You can also use `ai/qwen3-coder` if you want another coding-focused model
with a large context window.

## Step 2: Allow the sandbox to reach Docker Model Runner

Sandboxes are network-isolated by default, so you need a policy rule before
the sandbox can reach Docker Model Runner.

The rule is matched against the destination the proxy forwards to, not the
hostname the sandbox uses. Because the proxy rewrites
`host.docker.internal` to `localhost` before forwarding, the rule allows
`localhost:12434` even though Claude Code will use `host.docker.internal`
in its requests:

```console
$ sbx policy allow network localhost:12434
```

For background on host access from sandboxes, see
[Accessing host services from a sandbox](../manuals/ai/sandboxes/usage.md#accessing-host-services-from-a-sandbox).

## Step 3: Create a Claude Code sandbox

From your project directory, create a sandbox without launching the agent:

```console
$ cd ~/my-project
$ sbx create claude --name claude-dmr .
```

`sbx run` would also work, but it launches Claude Code immediately. Without
`ANTHROPIC_BASE_URL` set, Claude Code points at `api.anthropic.com` and
either prompts for OAuth or errors out before you can fix the endpoint.
Creating the sandbox first lets you write the local endpoint into it before
the agent starts.

You don't need to set an Anthropic API key or run `sbx secret set
anthropic`. Docker Model Runner doesn't authenticate the local endpoint,
and the sandbox proxy only injects credentials for requests bound for
`api.anthropic.com`. See
[Credentials](../manuals/ai/sandboxes/security/credentials.md) for the full
list of services the proxy authenticates.

## Step 4: Set the local endpoint inside the sandbox

Append `ANTHROPIC_BASE_URL` to the sandbox's persistent environment file so
Claude Code reads it on every launch:

```console
$ sbx exec -d claude-dmr bash -c "echo 'export ANTHROPIC_BASE_URL=http://host.docker.internal:12434' >> /etc/sandbox-persistent.sh"
```

The `bash -c` wrapper ensures the `>>` redirect runs inside the sandbox, not
on your host. For details on this approach, see
[How do I set custom environment variables inside a sandbox?](../manuals/ai/sandboxes/faq.md#how-do-i-set-custom-environment-variables-inside-a-sandbox).

To confirm the variable is set, open a shell in the sandbox:

```console
$ sbx exec -it claude-dmr bash
$ echo $ANTHROPIC_BASE_URL
http://host.docker.internal:12434
```

## Step 5: Verify connectivity to Docker Model Runner

Still inside the sandbox shell, send a test request to the host endpoint:

```console
$ curl http://host.docker.internal:12434/v1/messages \
  -H "Content-Type: application/json" \
  -d '{
    "model": "ai/devstral-small-2",
    "max_tokens": 32,
    "messages": [{"role": "user", "content": "Say hello"}]
  }'
```

A successful response confirms the policy rule and base URL are correct.
Type `exit` to leave the shell. For more details about the request format,
see the
[Anthropic-compatible API reference](../manuals/ai/model-runner/api-reference.md#anthropic-compatible-api).

## Step 6: Launch Claude Code with the local model

Run Claude Code in the sandbox and pass the model flag through to the agent:

```console
$ sbx run claude-dmr -- --model ai/devstral-small-2
```

Everything after `--` is forwarded to the Claude Code CLI. Because
`ANTHROPIC_BASE_URL` is set in the sandbox's persistent environment, Claude
Code routes requests to Docker Model Runner on your host instead of
`api.anthropic.com`.

## Step 7: Inspect Claude Code requests

To inspect the requests Claude Code sends, run on your host:

```console
$ docker model requests --model ai/devstral-small-2 | jq .
```

This helps you debug prompts, context usage, and compatibility issues
without attaching to the sandbox.

## Step 8: Package `gpt-oss` with a larger context window

`ai/gpt-oss` defaults to a smaller context window than coding-focused
models. To use it for repository-scale prompts, package a larger variant on
the host:

```console
$ docker model pull ai/gpt-oss
$ docker model package --from ai/gpt-oss --context-size 32000 gpt-oss:32k
```

Then point Claude Code at the packaged model the next time you run the
sandbox:

```console
$ sbx run claude-dmr -- --model gpt-oss:32k
```

## Clean up

Sandboxes persist after Claude Code exits. To stop the sandbox without
deleting it:

```console
$ sbx stop claude-dmr
```

To remove the sandbox and everything inside, including the persistent
environment file:

```console
$ sbx rm claude-dmr
```

Files in your workspace are unaffected.

## Learn more

- [Use Claude Code with Docker Model Runner](claude-code-model-runner.md)
- [Get started with Docker Sandboxes](../manuals/ai/sandboxes/get-started.md)
- [Claude Code in Docker Sandboxes](../manuals/ai/sandboxes/agents/claude-code.md)
- [Docker Model Runner overview](../manuals/ai/model-runner/_index.md)
- [Docker Model Runner API reference](../manuals/ai/model-runner/api-reference.md)
