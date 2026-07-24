---
title: "Managing Secrets"
description: "How to securely provide API keys and credentials to Docker Agent using environment variables, env files, Docker Compose secrets, and 1Password references."
keywords: docker agent, ai agents, guides, managing secrets
weight: 30
canonical: https://docs.docker.com/ai/docker-agent/guides/secrets/
---

_How to securely provide API keys and credentials to Docker Agent._

## Overview

Docker Agent needs API keys to talk to model providers (OpenAI, Anthropic, etc.) and MCP tool servers (GitHub, Slack, etc.). These keys are **never stored in config files**. Instead, Docker Agent resolves them at runtime through a chain of secret providers, checked in order (see `pkg/environment/default.go`):

| Priority | Provider | Description |
| --- | --- | --- |
| 1 | [Environment variables](#environment-variables) | `export OPENAI_API_KEY=sk-...` |
| 2 | [Docker Compose secrets](#docker-compose-secrets) | Files in `/run/secrets/` |
| 3 | [Docker Agent env file](#docker-agent-env-file) | `~/.config/cagent/.env`, written by `docker agent setup` |
| 4 | [Credential helper](#credential-helper) | Custom command declared in `~/.config/cagent/config.yaml` under `credential_helper:` |
| 5 | [Docker Desktop](#docker-desktop) | Secrets stored by the Docker Desktop backend (no setup on a Desktop install) |

The first provider that has a value wins. You can mix and match — for example, use environment variables for one key and the Docker Agent env file for another.

> [!NOTE]
> Older Docker Agent versions could also read secrets from the macOS Keychain and the `pass` password manager. These sources are no longer consulted: migrate any keys stored there to one of the sources above, e.g. by re-running `docker agent setup`.

Whatever provider returns the value, if that value looks like a [1Password secret reference](#1password-references) (it starts with `op://`), Docker Agent resolves it through the `op` CLI before handing it to a model provider or tool.

When Docker Agent runs inside a Docker sandbox (detected via `SANDBOX_VM_ID`), a sandbox token provider is prepended to the chain so that `DOCKER_TOKEN` is read from a continuously-refreshed file instead of a stale environment variable.

## Environment Variables

The simplest approach. Set variables in your shell before running Docker Agent:

```bash
export OPENAI_API_KEY=sk-...
export ANTHROPIC_API_KEY=sk-ant-...
docker agent run agent.yaml
```

Common variables:

| Variable | Provider |
| --- | --- |
| `OPENAI_API_KEY` | OpenAI |
| `ANTHROPIC_API_KEY` | Anthropic |
| `GOOGLE_API_KEY` | Google Gemini |
| `MISTRAL_API_KEY` | Mistral |
| `OPENROUTER_API_KEY` | OpenRouter |
| `XAI_API_KEY` | xAI |
| `NEBIUS_API_KEY` | Nebius |

MCP tools may require additional variables. For example, the GitHub MCP server needs `GITHUB_PERSONAL_ACCESS_TOKEN`. These are passed to tools via the `env` field in your config:

```yaml
toolsets:
  - type: mcp
    ref: docker:github-official
    env:
      GITHUB_PERSONAL_ACCESS_TOKEN: $GITHUB_PERSONAL_ACCESS_TOKEN
```

## Env Files

For convenience, you can store secrets in a `.env` file and pass it to Docker Agent with `--env-from-file`:

```bash
# .env
OPENAI_API_KEY=sk-...
ANTHROPIC_API_KEY=sk-ant-...
GITHUB_PERSONAL_ACCESS_TOKEN=ghp_...
```

```bash
docker agent run agent.yaml --env-from-file .env
```

The file format supports:

- `KEY=VALUE` pairs, one per line
- Comments starting with `#`
- Quoted values: `KEY="value with spaces"`
- Blank lines are ignored

> [!IMPORTANT]
> Add `.env` to your `.gitignore` to avoid committing secrets to version control.

## Docker Agent env file

A `.env` file (same format as above) at `~/.config/cagent/.env` is read automatically on every run — no `--env-from-file` flag needed. It is where [`docker agent setup`](../../features/cli/index.md#docker-agent-setup) stores API keys when you choose the env-file location, and you can edit it by hand:

```bash
# ~/.config/cagent/.env
OPENAI_API_KEY=sk-...
```

The file is created with owner-only permissions (`0600`), but the values are stored in plain text.

## Docker Compose Secrets

When running Docker Agent in a container with Docker Compose, you can use [Compose secrets](https://docs.docker.com/compose/how-tos/use-secrets/) to inject credentials securely. Compose mounts secrets as files under `/run/secrets/`, and Docker Agent reads from this location automatically.

### From a file

Store each secret in its own file, then reference it in `compose.yaml`:

```bash
echo -n "sk-ant-your-key-here" > .anthropic_api_key
```

```yaml
# compose.yaml
services:
  agent:
    image: docker/docker-agent
    command: run --exec /app/agent.yaml "Hello!"
    secrets:
      - ANTHROPIC_API_KEY
    volumes:
      - ./agent.yaml:/app/agent.yaml:ro

secrets:
  ANTHROPIC_API_KEY:
    file: ./.anthropic_api_key
```

Docker Compose mounts the file as `/run/secrets/ANTHROPIC_API_KEY`. Docker Agent picks it up with no extra configuration.

### From a host environment variable

In CI/CD pipelines, secrets are often injected as environment variables. Compose can forward these to `/run/secrets/`:

```yaml
secrets:
  ANTHROPIC_API_KEY:
    environment: "ANTHROPIC_API_KEY"
```

### Multiple secrets

```yaml
services:
  agent:
    image: docker/docker-agent
    command: run --exec /app/agent.yaml "Summarize my GitHub issues"
    secrets:
      - ANTHROPIC_API_KEY
      - GITHUB_PERSONAL_ACCESS_TOKEN
    volumes:
      - ./agent.yaml:/app/agent.yaml:ro

secrets:
  ANTHROPIC_API_KEY:
    file: ./.anthropic_api_key
  GITHUB_PERSONAL_ACCESS_TOKEN:
    file: ./.github_token
```

### Why use Compose secrets over environment variables?

| Aspect | Environment Variables | Compose Secrets |
| --- | --- | --- |
| Storage | In memory, visible via `docker inspect` | Mounted as tmpfs files under `/run/secrets/` |
| Visibility | Shown in process list and inspect output | Not exposed in `docker inspect` |
| Best for | Development | Production and CI/CD |

## Credential Helper

Docker Agent can shell out to an external credential helper you define in your user config. This is useful when your organisation already has a secrets daemon you want to reuse (HashiCorp Vault, 1Password CLI, `bitwarden-cli`, etc.).

Declare the helper in `~/.config/cagent/config.yaml`:

```yaml
# ~/.config/cagent/config.yaml
credential_helper:
  command: op
  args: ["read", "op://Personal/docker-agent"]
```

The command is invoked with the variable name appended as the final argument, and must print the secret value to stdout.

## Docker Desktop

On machines where Docker Desktop is installed, Docker Agent queries Docker Desktop's backend for secrets stored against your signed-in Docker account. This is transparent — no extra configuration — and it is how signed-in Docker users get provider API keys without setting any environment variables.

## 1Password References

Any secret value resolved through the chain above can be a **1Password secret reference** instead of the literal secret. If the value starts with `op://`, Docker Agent resolves it by invoking the [1Password CLI](https://developer.1password.com/docs/cli/) (`op read <reference>`) and uses the result.

This works with every provider — most commonly an environment variable or env file:

```bash
export OPENAI_API_KEY="op://Personal/OpenAI/api-key"
docker agent run agent.yaml
```

References follow the `op://<vault>/<item>/<field>` format. Make sure the `op` CLI is installed and you are signed in (`op signin`) so that non-interactive reads succeed.

> [!WARNING]
> **Behaviour when resolution fails**
>
> If the value starts with `op://` but the `op` CLI is not installed, or the reference cannot be read (not signed in, wrong path, locked vault), Docker Agent logs a warning and uses an **empty value** — it never forwards the raw `op://` reference to a model provider or tool. Resolved references (and deterministic failures) are cached for the lifetime of the run; transient failures such as a cancelled lookup are not cached, so a later attempt can retry.

## Choosing a Method

| Method | Best for | Setup effort |
| --- | --- | --- |
| Environment variables | Quick local development, scripts | Low |
| Env files | Team projects, multiple keys | Low |
| Docker Agent env file | Keys used across all projects, written by `docker agent setup` | Low |
| Docker Compose secrets | Containerized deployments, CI/CD | Medium |
| Credential helper | Reusing an existing secrets daemon (Vault, 1Password CLI, ...) | Medium |
| 1Password references (`op://`) | Teams already using 1Password | Low |

You can combine methods. For example, store long-lived provider keys in the Docker Agent env file and pass project-specific MCP tokens via env files.

## Preventing Secret Leaks

Provider keys live in the secret store and are passed to Docker Agent through the chain above — the agent itself never receives them as input. But the **content of a conversation** can still leak credentials: a user pasting a token, a tool returning a config file with embedded keys, a transcript dumped into a prompt.

For that defense-in-depth case, set `redact_secrets: true` on an agent. It scrubs detected secrets out of:

- the arguments of every outgoing tool call (before the tool sees them),
- every outgoing chat message (before the model provider sees them), and
- every tool's output (before it reaches event consumers, the persisted session file, the `post_tool_use` hook input, or the next LLM call).

```yaml
agents:
  root:
    model: openai/gpt-5
    description: A helpful assistant
    instruction: You are a helpful assistant.
    redact_secrets: true
    toolsets:
      - type: shell
```

The ruleset covers GitHub PATs, AWS / GCP / Azure credentials, Stripe / Slack / GitLab / Hugging Face tokens, JWTs, PEM-encoded private keys, Docker Hub PATs, and many others. Each detected span is replaced with the literal `[REDACTED]`. See the [Redacting Secrets](../../configuration/agents/index.md#redacting-secrets) section in the agent configuration reference for the full picture and important caveats about false negatives.
