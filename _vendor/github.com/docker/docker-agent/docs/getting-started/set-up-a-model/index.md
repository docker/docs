---
title: "Set Up a Model"
description: "Make a model available to Docker Agent: bring a cloud provider API key or run a local model with Docker Model Runner."
keywords: docker agent, ai agents, getting started, set up a model, api key, local model, docker model runner
weight: 25
canonical: https://docs.docker.com/ai/docker-agent/getting-started/set-up-a-model/
---

_Most agents need a model to think with: bring an API key for a cloud provider, or run a model locally with Docker Model Runner. This page walks through both paths end to end — plus the exception: agents that delegate to the Claude Code CLI on a Claude subscription, which need no model at all._

## Pick a Path

|              | Cloud provider (Path A)                 | Local model (Path B)                        |
| ------------ | ---------------------------------------- | ------------------------------------------- |
| You need     | An account and an API key                | Docker Desktop with Model Runner enabled    |
| Cost         | Pay per token                            | Free once the model is downloaded           |
| Your data    | Sent to the provider                     | Never leaves your machine                   |
| Model quality| Frontier models (Claude, GPT-5, Gemini)  | Open models sized to your hardware          |

You can set up both. When you don't name a model, Docker Agent's `auto` selection picks the first cloud provider with a configured key and falls back to a locally pulled Docker Model Runner model.

There is also a third path that is neither: if you have a **Claude subscription**, the [Claude Code harness](#path-c-claude-code-harness-claude-subscription) runs the official `claude` CLI as the agent — no API key and no local model required.

> [!TIP]
> **Prefer a wizard?**
>
> `docker agent setup` walks through the same choices interactively: pick a provider and store its key, check Docker Model Runner and pull a local model, or set up the Claude Code harness. This page is the manual version. See the [CLI reference](../../features/cli/index.md#docker-agent-setup).

## Path A: Cloud Provider (API Key)

### 1. Get an API key

Create a key in your provider's console:

| Provider      | Environment variable | Get a key at                                                        |
| ------------- | -------------------- | ------------------------------------------------------------------- |
| Anthropic     | `ANTHROPIC_API_KEY`  | [console.anthropic.com](https://console.anthropic.com/settings/keys) |
| OpenAI        | `OPENAI_API_KEY`     | [platform.openai.com](https://platform.openai.com/api-keys)          |
| Google Gemini | `GOOGLE_API_KEY`     | [aistudio.google.com](https://aistudio.google.com/apikey)            |

Every other provider works the same way. See [Model Providers](../../providers/overview/index.md) for the full list and each provider's environment variable.

### 2. Store the key

The fastest option is an environment variable in your shell:

```bash
$ export ANTHROPIC_API_KEY=sk-ant-...
```

That lasts for the current shell session. To set a key up once, use any other built-in secret source:

```bash
# Env file, passed at run time
$ echo 'ANTHROPIC_API_KEY=sk-ant-...' > .env
$ docker agent run --env-from-file .env

# Docker Agent env file, read automatically on every run
# (`docker agent setup` writes it for you with owner-only permissions)
$ echo 'ANTHROPIC_API_KEY=sk-ant-...' >> ~/.config/cagent/.env
$ chmod 600 ~/.config/cagent/.env
```

The entry name must match the environment variable the provider expects. [Managing Secrets](../../guides/secrets/index.md) covers every source (Docker Compose secrets, credential helpers, 1Password references) and the order they are checked in.

> [!IMPORTANT]
> Keys never go in `agent.yaml`. If you use an env file, add it to `.gitignore`.

### 3. Verify

`docker agent doctor` shows whether the key is visible and where it comes from:

```bash
$ docker agent doctor
```

```text
User configuration
  ~/.config/cagent/config.yaml: ok

Model provider credentials
  PROVIDER    STATUS    CREDENTIAL          SOURCE
  anthropic   found     ANTHROPIC_API_KEY   environment
  openai      not set   OPENAI_API_KEY      -
  ...

Docker Model Runner
  Status: not installed (https://docs.docker.com/ai/model-runner/get-started/)

Model auto-selection
  auto -> anthropic/claude-sonnet-4-6

No issues found.
```

### 4. Run

```bash
$ docker agent run
```

With no config file, the default agent picks the provider you configured. To name a model explicitly, use `--model` or the `model` field in your config:

```bash
$ docker agent run --model anthropic/claude-sonnet-4-5
```

```yaml
agents:
  root:
    model: anthropic/claude-sonnet-4-5
    description: A helpful coding assistant
    instruction: You are an expert software developer.
```

## Path B: Local Model (Docker Model Runner)

Docker Model Runner (DMR) runs open models on your own machine: no API key, no per-token cost, and prompts never leave your computer.

### 1. Install Docker Model Runner

Model Runner ships with [Docker Desktop](https://www.docker.com/products/docker-desktop/) (enable it under **Settings > AI**) and is also available for Docker Engine. Check that it responds:

```bash
$ docker model status
```

If the command is missing or fails, follow the [Model Runner get-started guide](https://docs.docker.com/ai/model-runner/get-started/).

### 2. Pull a model

```bash
$ docker model pull ai/qwen3
```

`ai/qwen3` is the model Docker Agent reaches for by default, but any model from the [Docker Hub `ai` catalog](https://hub.docker.com/u/ai) works. Pick one sized for your machine's memory. List what you have locally:

```bash
$ docker model ls
```

### 3. Verify

```bash
$ docker agent doctor
```

```text
User configuration
  ~/.config/cagent/config.yaml: ok

Model provider credentials
  PROVIDER    STATUS    CREDENTIAL          SOURCE
  anthropic   not set   ANTHROPIC_API_KEY   -
  ...

Docker Model Runner
  Status: reachable, 1 model(s) pulled:
    - ai/qwen3:latest

Model auto-selection
  auto -> dmr/ai/qwen3:latest

No issues found.
```

### 4. Run

```bash
$ docker agent run --model dmr/ai/qwen3
```

Or in your config:

```yaml
agents:
  root:
    model: dmr/ai/qwen3
    description: A local assistant
    instruction: You are a helpful assistant.
```

When no cloud key is configured, bare `docker agent run` auto-selects a pulled local model, so after `docker model pull` you can run with no flags at all. The [Docker Model Runner provider page](../../providers/dmr/index.md) covers context size, runtime flags, and other tuning options.

## Path C: Claude Code Harness (Claude Subscription)

If you already pay for a Claude subscription, an agent can delegate its work to
the official Claude Code CLI instead of calling a model API. This is an
**external CLI, not provider API access**: Docker Agent launches `claude`,
which authenticates with its own subscription login — no `ANTHROPIC_API_KEY`,
no Docker Model Runner, and no token ever passes through Docker Agent.

### 1. Install and log in

Install [Claude Code](https://docs.anthropic.com/en/docs/claude-code), then
log in **as the same OS user and environment that run `docker agent`**:

```bash
$ claude auth login --claudeai   # interactive, opens a browser
$ claude auth status --text      # verify
```

### 2. Create a harness agent

`docker agent setup` (pick "Claude Code harness") generates this file for you,
or write it yourself:

```yaml
# claude-code-agent.yaml
agents:
  root:
    description: Claude Code running on your Claude subscription
    harness:
      type: claude-code
      effort: medium # low | medium | high | xhigh | max; omit for the Claude Code default
```

### 3. Verify and run

```bash
$ docker agent doctor claude-code-agent.yaml   # checks the CLI is installed and logged in
$ docker agent run claude-code-agent.yaml
```

The harness runs the CLI non-interactively and bypasses Claude Code's
permission prompts, so use it in a repository you trust — see the security
notes and full field reference in [Coding Harnesses](../../features/harnesses/index.md).

## Check Your Setup Anytime

`docker agent doctor` reports which providers have credentials (and from which source), whether Docker Model Runner is reachable and which models are pulled, and which model `auto` would pick. Secret values are never printed.

```bash
$ docker agent doctor                     # credential, DMR, and auto-selection state
$ docker agent doctor ./agent.yaml        # also check that file's requirements
```

It exits non-zero when something would block a run, which makes it usable as a CI preflight. See the [CLI reference](../../features/cli/index.md#docker-agent-doctor).

## What's Next?

- [**Quick Start**](../quickstart/index.md) — run your first agent now that a model is available.
- [**Models**](../../concepts/models/index.md) — inline vs. named models, fallbacks, and `auto` selection.
- [**Managing Secrets**](../../guides/secrets/index.md) — every way to store credentials, compared.
- [**Troubleshooting**](../../community/troubleshooting/index.md#missing-credentials-or-model-errors) — decode "no model available" and credential errors.
