---
title: Use Claude Code with Docker Model Runner
description: Configure Claude Code to use Docker Model Runner so you can code with local models.
summary: |
  Connect Claude Code to Docker Model Runner with the Anthropic-compatible API,
  package `gpt-oss` with a larger context window, and inspect requests.
keywords: ai, claude code, docker model runner, anthropic, local models, coding assistant
tags: [ai]
params:
  time: 10 minutes
---

This guide shows how to run Claude Code with Docker Model Runner as the backend
model provider. You'll point Claude Code at the local Anthropic-compatible API,
run a coding model, and package `gpt-oss` with a larger context window for
longer repository prompts.

> **Acknowledgment**
>
> Docker would like to thank [Pradumna Saraf](https://twitter.com/pradumna_saraf) for his contribution to this guide.

In this guide, you'll learn how to:

- Pull a coding model and start Claude Code with Docker Model Runner
- Make the endpoint configuration persistent
- Verify the local API endpoint and inspect requests
- Package `gpt-oss` with a larger context window for longer prompts

## Prerequisites

Before you start, make sure you have:

- [Docker Desktop](../get-started/get-docker.md) or Docker Engine installed
- [Docker Model Runner enabled](../manuals/ai/model-runner/get-started.md#enable-docker-model-runner)
- [Claude Code installed](https://code.claude.com/docs/en/quickstart)

If you use Docker Desktop, turn on TCP access in **Settings** > **AI**, or run:

```console
$ docker desktop enable model-runner --tcp 12434
```

## Step 1: Pull a coding model

Pull a model before you start Claude Code:

```console
$ docker model pull ai/devstral-small-2
```

You can also use `ai/qwen3-coder` if you want another coding-focused model with
a large context window.

## Step 2: Start Claude Code with Docker Model Runner

Set `ANTHROPIC_BASE_URL` to your local Docker Model Runner endpoint when you run
Claude Code.

On macOS or Linux:

```console
$ ANTHROPIC_BASE_URL=http://localhost:12434 claude --model ai/devstral-small-2
```

On Windows PowerShell:

```powershell
$env:ANTHROPIC_BASE_URL="http://localhost:12434"
claude --model ai/devstral-small-2
```

Claude Code now sends requests to Docker Model Runner instead of Anthropic's
hosted API.

## Step 3: Troubleshoot your first launch

If Claude Code can't connect, check Docker Model Runner status:

```console
$ docker model status
```

If Claude Code can't find the model, list local models:

```console
$ docker model ls
```

If the model is missing, pull it first. If needed, use the fully qualified
model name, such as `ai/devstral-small-2`.

## Step 4: Make the endpoint persistent

To avoid setting the environment variable each time, add it to your shell
profile:

```bash {title="~/.bashrc or ~/.zshrc"}
export ANTHROPIC_BASE_URL=http://localhost:12434
```

On Windows PowerShell, add it to your PowerShell profile:

```powershell {title="$PROFILE"}
$env:ANTHROPIC_BASE_URL = "http://localhost:12434"
```

After you reload your shell, you can run Claude Code with only the model flag:

```console
$ claude --model ai/devstral-small-2
```

## Step 5: Verify the API endpoint

Send a test request to confirm the Anthropic-compatible API is reachable:

```console
$ curl http://localhost:12434/v1/messages \
  -H "Content-Type: application/json" \
  -d '{
    "model": "ai/devstral-small-2",
    "max_tokens": 32,
    "messages": [{"role": "user", "content": "Say hello"}]
  }'
```

For more details about the request format, see the
[Anthropic-compatible API reference](../manuals/ai/model-runner/api-reference.md#anthropic-compatible-api).

## Step 6: Inspect Claude Code requests

To inspect the requests Claude Code sends to Docker Model Runner, run:

```console
$ docker model requests --model ai/devstral-small-2 | jq .
```

This helps you debug prompts, context usage, and compatibility issues.

## Step 7: Package `gpt-oss` with a larger context window

`ai/gpt-oss` defaults to a smaller context window than coding-focused models. If
you want to use it for repository-scale prompts, package a larger variant:

```console
$ docker model pull ai/gpt-oss
$ docker model package --from ai/gpt-oss --context-size 32000 gpt-oss:32k
```

Then run Claude Code with the packaged model:

```console
$ ANTHROPIC_BASE_URL=http://localhost:12434 claude --model gpt-oss:32k
```

## Learn more

- [Docker Model Runner overview](../manuals/ai/model-runner/_index.md)
- [Docker Model Runner API reference](../manuals/ai/model-runner/api-reference.md)
- [IDE and tool integrations](../manuals/ai/model-runner/ide-integrations.md)
