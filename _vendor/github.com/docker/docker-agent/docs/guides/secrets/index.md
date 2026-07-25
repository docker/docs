---
title: "Managing Secrets"
description: "How to securely provide API keys and credentials to docker-agent using environment variables, env files, Docker Compose secrets, macOS Keychain, pass, and 1Password references."
keywords: docker agent, ai agents, guides, managing secrets
weight: 30
canonical: https://docs.docker.com/ai/docker-agent/guides/secrets/
---

_How to securely provide API keys and credentials to docker-agent._

## Overview

docker-agent needs API keys to talk to model providers (OpenAI, Anthropic, etc.) and MCP tool servers (GitHub, Slack, etc.). These keys are **never stored in config files**. Instead, docker-agent resolves them at runtime through a chain of secret providers, checked in order (see `pkg/environment/default.go`):

| Priority | Provider | Description |
| --- | --- | --- |
| 1 | [Environment variables](#environment-variables) | `export OPENAI_API_KEY=sk-...` |
| 2 | [Docker Compose secrets](#docker-compose-secrets) | Files in `/run/secrets/` |
| 3 | [docker agent env file](#docker-agent-env-file) | `~/.config/cagent/.env`, written by `docker agent setup` |
| 4 | [Credential helper](#credential-helper) | Custom command declared in `~/.config/cagent/config.yaml` under `credential_helper:` |
| 5 | [Docker Desktop](#docker-desktop) | Secrets stored by the Docker Desktop backend (no setup on a Desktop install) |
| 6 | [`pass` password manager](#pass-password-manager) | `pass insert OPENAI_API_KEY` |
| 7 | [macOS Keychain](#macos-keychain) | `security add-generic-password` |

The first provider that has a value wins. You can mix and match — for example, use environment variables for one key and Keychain for another.

Whatever provider returns the value, if that value looks like a [1Password secret reference](#1password-references) (it starts with `op://`), docker-agent resolves it through the `op` CLI before handing it to a model provider or tool.

When docker-agent runs inside a Docker sandbox (detected via `SANDBOX_VM_ID`), a sandbox token provider is prepended to the chain so that `DOCKER_TOKEN` is read from a continuously-refreshed file instead of a stale environment variable.

## Environment Variables

The simplest approach. Set variables in your shell before running docker-agent:

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

For convenience, you can store secrets in a `.env` file and pass it to docker-agent with `--env-from-file`:

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

## docker agent env file

A `.env` file (same format as above) at `~/.config/cagent/.env` is read automatically on every run — no `--env-from-file` flag needed. It is where [`docker agent setup`](../../features/cli/index.md#docker-agent-setup) stores API keys when you choose the env-file location, and you can edit it by hand:

```bash
# ~/.config/cagent/.env
OPENAI_API_KEY=sk-...
```

The file is created with owner-only permissions (`0600`), but the values are stored in plain text: prefer the OS keychain or `pass` when available.

## Docker Compose Secrets

When running docker-agent in a container with Docker Compose, you can use [Compose secrets](https://docs.docker.com/compose/how-tos/use-secrets/) to inject credentials securely. Compose mounts secrets as files under `/run/secrets/`, and docker-agent reads from this location automatically.

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

Docker Compose mounts the file as `/run/secrets/ANTHROPIC_API_KEY`. docker-agent picks it up with no extra configuration.

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

docker-agent can shell out to an external credential helper you define in your user config. This is useful when your organisation already has a secrets daemon you want to reuse (HashiCorp Vault, 1Password CLI, `bitwarden-cli`, etc.).

Declare the helper in `~/.config/cagent/config.yaml`:

```yaml
# ~/.config/cagent/config.yaml
credential_helper:
  command: op
  args: ["read", "op://Personal/docker-agent"]
```

The command is invoked with the variable name appended as the final argument, and must print the secret value to stdout.

## Docker Desktop

On machines where Docker Desktop is installed, docker-agent queries Docker Desktop's backend for secrets stored against your signed-in Docker account. This is transparent — no extra configuration — and it is how signed-in Docker users get provider API keys without setting any environment variables.

## `pass` Password Manager

docker-agent integrates with [`pass`](https://www.passwordstore.org/), the standard Unix password manager. Secrets are stored as GPG-encrypted files in `~/.password-store/`.

### Store a secret

```bash
pass insert ANTHROPIC_API_KEY
```

The entry name must match the environment variable name that docker-agent expects.

### Verify it works

```bash
pass show ANTHROPIC_API_KEY
```

Once `pass` is set up, docker-agent resolves secrets from it automatically.

## macOS Keychain

On macOS, docker-agent can read secrets from the system Keychain. This is useful for local development — you store the key once and it's available across all your projects.

### Store a secret

```bash
security add-generic-password -a "$USER" -s ANTHROPIC_API_KEY -w "sk-ant-your-key-here"
```

The `-s` (service name) must match the environment variable name that docker-agent expects.

### Verify it works

```bash
security find-generic-password -s ANTHROPIC_API_KEY -w
```

### Delete a secret

```bash
security delete-generic-password -s ANTHROPIC_API_KEY
```

Once stored, docker-agent finds the secret automatically — no flags or config needed.

## 1Password References

Any secret value resolved through the chain above can be a **1Password secret reference** instead of the literal secret. If the value starts with `op://`, docker-agent resolves it by invoking the [1Password CLI](https://developer.1password.com/docs/cli/) (`op read <reference>`) and uses the result.

This works with every provider — most commonly an environment variable or env file:

```bash
export OPENAI_API_KEY="op://Personal/OpenAI/api-key"
docker agent run agent.yaml
```

References follow the `op://<vault>/<item>/<field>` format. Make sure the `op` CLI is installed and you are signed in (`op signin`) so that non-interactive reads succeed.

> [!WARNING]
> **Behaviour when resolution fails**
>
> If the value starts with `op://` but the `op` CLI is not installed, or the reference cannot be read (not signed in, wrong path, locked vault), docker-agent logs a warning and uses an **empty value** — it never forwards the raw `op://` reference to a model provider or tool. Resolved references (and deterministic failures) are cached for the lifetime of the run; transient failures such as a cancelled lookup are not cached, so a later attempt can retry.

## Choosing a Method

| Method | Best for | Setup effort |
| --- | --- | --- |
| Environment variables | Quick local development, scripts | Low |
| Env files | Team projects, multiple keys | Low |
| docker agent env file | Keys used across all projects, written by `docker agent setup` | Low |
| Docker Compose secrets | Containerized deployments, CI/CD | Medium |
| `pass` | Linux/macOS, GPG-based workflows | Medium |
| macOS Keychain | macOS local development | Low |
| 1Password references (`op://`) | Teams already using 1Password | Low |

You can combine methods. For example, store long-lived provider keys in macOS Keychain and pass project-specific MCP tokens via env files.

## Preventing Secret Leaks

Provider keys live in the secret store and are passed to docker-agent through the chain above — the agent itself never receives them as input. But the **content of a conversation** can still leak credentials: a user pasting a token, a tool returning a config file with embedded keys, a transcript dumped into a prompt.

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
