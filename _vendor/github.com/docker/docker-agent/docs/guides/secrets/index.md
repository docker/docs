---
title: "Managing Secrets"
description: "How to securely provide API keys and credentials to Docker Agent using environment variables, env files, Docker Compose secrets, and 1Password references."
keywords: docker agent, ai agents, guides, managing secrets
weight: 30
canonical: https://docs.docker.com/ai/docker-agent/guides/secrets/
---

_How to securely provide API keys and credentials to Docker Agent._

## Overview

Docker Agent needs API keys to talk to model providers (OpenAI, Anthropic, etc.) and MCP tool servers (GitHub, Slack, etc.). Avoid putting literal secrets in agent configuration fields such as `env`, headers, or URLs. Docker Agent can resolve them at runtime through a chain of secret providers, checked in order (see `pkg/environment/default.go`):

| Priority | Provider | Description |
| --- | --- | --- |
| 1 | [Environment variables](#environment-variables) | `export OPENAI_API_KEY=sk-...` |
| 2 | [Docker Compose secrets](#docker-compose-secrets) | Files in `/run/secrets/` |
| 3 | [Docker Agent env file](#docker-agent-env-file) | `~/.config/cagent/.env`, written by `docker agent setup` |
| 4 | [ChatGPT login](../../providers/chatgpt/index.md) | `CHATGPT_OAUTH_TOKEN` from the account login created by `docker agent setup` |
| 5 | [Credential helper](#credential-helper) | Custom command declared in `~/.config/cagent/config.yaml` under `credential_helper:`; supplies `DOCKER_TOKEN` only |
| 6 | [Docker Desktop](#docker-desktop) | Secrets stored by the Docker Desktop backend (no setup on a Desktop install) |

The first provider that has a value wins. You can mix and match — for example, use environment variables for one key and the Docker Agent env file for another.

> [!NOTE]
> Older Docker Agent versions could also read secrets from the macOS Keychain and the `pass` password manager. These sources are no longer consulted: migrate any keys stored there to one of the sources above, e.g. by re-running `docker agent setup`.

Whatever provider returns the value, if that value looks like a [1Password secret reference](#1password-references) (it starts with `op://`), Docker Agent resolves it through the `op` CLI before handing it to a model provider or tool. This applies to every variable, API keys included, and needs no credential helper: `ANTHROPIC_API_KEY=op://Personal/Anthropic/api-key` in any source above is enough.

When Docker Agent runs an agent with `--sandbox` and a Docker AI gateway is configured, the sandbox proxy authenticates gateway requests itself: it injects your Docker login token into HTTPS requests to the gateway host (docker.com domains only). The token never enters the sandbox — inside it, `DOCKER_TOKEN` is a proxy-managed placeholder.

## Environment Variables

The simplest approach. Set variables in your shell before running Docker Agent:

```bash
export OPENAI_API_KEY=sk-...
export ANTHROPIC_API_KEY=sk-ant-...
docker agent run agent.yaml
```

See [Provider Credentials](../../providers/overview/index.md#provider-credentials) for the full list of provider-specific environment variables and authentication methods.

MCP tools may require additional variables. For example, the GitHub MCP server needs `GITHUB_PERSONAL_ACCESS_TOKEN`. These are passed to tools via the `env` field in your config:

```yaml
toolsets:
  - type: mcp
    ref: docker:github-official
    env:
      GITHUB_PERSONAL_ACCESS_TOKEN: $GITHUB_PERSONAL_ACCESS_TOKEN
```

### Docker Agent-only variables

Other tools read the same variables: Claude Code reads `ANTHROPIC_API_KEY`, Codex reads `OPENAI_API_KEY`. To give Docker Agent its own value, prefix the name with `DOCKER_AGENT_ENV_`. Docker Agent looks up `DOCKER_AGENT_ENV_<NAME>` before `<NAME>`, and other tools ignore it:

```bash
export ANTHROPIC_API_KEY=sk-ant-...                                        # used by other tools
export DOCKER_AGENT_ENV_ANTHROPIC_API_KEY="op://Personal/Anthropic/api-key" # used by Docker Agent
docker agent run agent.yaml
```

- The prefix works for every variable Docker Agent resolves through its secret chain: provider API keys, `${VAR}` references in MCP `env`, `DOCKER_TOKEN`, and the OpenAI key used by `/speak`.
- `DOCKER_AGENT_ENV_<NAME>` wins over `<NAME>` across all [sources](#overview) and `--env-from-file`: an override in the Docker Agent env file beats a plain shell export.
- An empty override is ignored: `export DOCKER_AGENT_ENV_ANTHROPIC_API_KEY=` falls back to an override in a later source, then to `ANTHROPIC_API_KEY`.
- Docker Agent's own settings (`DOCKER_AGENT_*`) are never overridden this way.
- The value can be a [1Password reference](#1password-references).
- `docker agent doctor` shows which name supplied each credential.

> [!NOTE]
> With `--sandbox`, the provider API keys that the sandbox proxy injects itself (such as `OPENAI_API_KEY` and `ANTHROPIC_API_KEY`) are not read through Docker Agent, so neither the prefix nor `op://` references apply to them. Other variables are resolved on the host, prefix included, and forwarded under their plain name.

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

Docker Agent can shell out to an external credential helper you define in your user config to obtain the **Docker token** (`DOCKER_TOKEN`) used by the [Docker models gateway](../../configuration/models/index.md). This is useful when your organisation already keeps that token in a secrets daemon (HashiCorp Vault, 1Password CLI, `bitwarden-cli`, etc.).

> [!IMPORTANT]
> **The credential helper does not serve API keys.** It is consulted for `DOCKER_TOKEN` only. To load provider API keys or MCP tokens from 1Password, store an `op://` reference in any other source instead, for example in the [Docker Agent env file](#docker-agent-env-file):
>
> ```bash
> # ~/.config/cagent/.env
> ANTHROPIC_API_KEY=op://Personal/Anthropic/api-key
> ```
>
> See [1Password References](#1password-references).

Declare the helper in `~/.config/cagent/config.yaml`:

```yaml
# ~/.config/cagent/config.yaml
credential_helper:
  command: op
  args: ["read", "op://Personal/Docker/token"]
```

The command is run exactly as configured, with no extra arguments, and must print the token to stdout.

## Docker Desktop

On machines where Docker Desktop is installed, Docker Agent queries Docker Desktop's backend for secrets stored against your signed-in Docker account. This is transparent — no extra configuration — and it is how signed-in Docker users get provider API keys without setting any environment variables.

## Docker Authentication

Routing model traffic through the [Docker models gateway](../../configuration/models/index.md) needs a Docker token. A `DOCKER_TOKEN` you provide yourself — an environment variable, a Compose secret, an env file or a credential helper, as is common in CI — always takes precedence, and none of the sources below are consulted.

Otherwise, Docker Agent tries these sources in order and uses the first token that is not expired or about to expire:

1. **The Docker secrets engine**, which Docker Desktop serves on a local socket and which holds the session of the account signed in to Desktop. It is skipped when it is unavailable (Docker Desktop is not installed or not running, or predates the engine), holds no session, or holds a session that is incomplete: an access token that is missing or not a JWT, or no account or expiry. A request waits for the engine for at most five seconds, and never more than half of its own time limit, so the other sources still get a turn. An engine that fails or does not answer is skipped for the next 30 seconds. Production sessions are read unless `DOCKER_AGENT_HUB_LOGIN_URL` points at a Docker Hub staging host, such as `hub-stage.docker.com`; staging sessions are read instead. Docker Agent cannot tell which environment Docker Desktop is signed in to, so set that variable when Desktop is signed in to staging.
2. **Docker Desktop's backend.**
3. **The access token `docker login` stored.** Docker Agent exchanges the long-lived access token in your credential store for a fresh Docker token, the same exchange `docker login` itself performs. Signing in with `docker login` is therefore enough; Docker Desktop is not required.
4. **A refresh of Docker Desktop's session**, as a last resort.

The tokens Docker Desktop hands out are valid for 15 minutes, and only Desktop can renew them. That is why the exchange exists: it keeps Docker Agent working when Desktop has nothing usable to offer because it is signed out, not running, or its own refresh is stuck. Tokens are looked up again as they near expiry, so a long-running session keeps working.

Only Docker access tokens are exchanged — the `dckr_…` secrets `docker login` stores — never an account password, and the exchange goes to Docker Hub over HTTPS. The resulting bearer token is cached in a private file under Docker Agent's cache directory so sibling processes reuse it instead of minting their own, and it stops being used within seconds of a `docker logout` or an account switch. Run `docker agent debug auth` to see which token is in use and where it came from.

Set `DOCKER_AGENT_NO_TOKEN_EXCHANGE=1` to opt out: Docker Agent then relies on Docker Desktop alone.

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

> [!TIP]
> Other tools that read the same variables (Claude Code, Codex, provider SDKs) do not understand `op://` values. To keep references out of your shell, put them in the [Docker Agent env file](#docker-agent-env-file), which only Docker Agent reads, or point a model at a dedicated variable with [`token_key`](../../configuration/models/index.md).

## Choosing a Method

| Method | Best for | Setup effort |
| --- | --- | --- |
| Environment variables | Quick local development, scripts | Low |
| Env files | Team projects, multiple keys | Low |
| Docker Agent env file | Keys used across all projects, written by `docker agent setup` | Low |
| Docker Compose secrets | Containerized deployments, CI/CD | Medium |
| Credential helper | Fetching `DOCKER_TOKEN` from an existing secrets daemon (Vault, 1Password CLI, ...) | Medium |
| 1Password references (`op://`) | Teams already using 1Password | Low |

You can combine methods. For example, store long-lived provider keys in the Docker Agent env file and pass project-specific MCP tokens via env files.

## Preventing Secret Leaks

Provider keys live in the secret store and are passed to Docker Agent through the chain above — the agent itself never receives them as input. But the **content of a conversation** can still leak credentials: a user pasting a token, a tool returning a config file with embedded keys, a transcript dumped into a prompt.

Secret redaction is enabled by default as a defense in depth. Set `redact_secrets: false` on an agent to opt out; the explicit `true` below is equivalent to omitting the field. It scrubs detected secrets out of:

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
