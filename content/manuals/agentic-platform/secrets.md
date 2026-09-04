---
title: Secrets
description: Manage model provider and service credentials for Docker Agentic Platform sandboxes.
keywords: docker agentic platform, secrets, api keys, service credentials, proxy
weight: 40
aliases:
  - /agentic-platform/concepts/secrets/
  - /agentic-platform/guides/manage-secrets/
---

Secrets provide credentials to agents without placing their values in a
sandbox. Docker Agentic Platform stores each value outside sandboxes and uses
the sandbox proxy to apply it to matching requests.

Docker Agentic Platform supports the following credentials:

| Service       | Secret ID           | Used by                                           |
| ------------- | ------------------- | ------------------------------------------------- |
| Anthropic     | `ANTHROPIC_API_KEY` | Claude Code and OpenCode with Anthropic models    |
| OpenAI        | `OPENAI_API_KEY`    | Codex and OpenCode with OpenAI models             |
| Google Gemini | `GEMINI_API_KEY`    | Gemini CLI and OpenCode with Google models        |
| Groq          | `GROQ_API_KEY`      | OpenCode with Groq models                         |
| xAI           | `XAI_API_KEY`       | OpenCode with xAI models                          |
| GitHub        | `GITHUB_TOKEN`      | Copilot and any sandbox type that accesses GitHub |

Provider-specific credentials are applied only to matching requests from
compatible sandbox types. For example, an Anthropic credential is not applied
to requests from a Codex sandbox. OpenCode can use any configured provider that
it supports.

## GitHub credential

`GITHUB_TOKEN` authenticates both Copilot and GitHub repository operations.
Entering it in the launcher stores it as the same GitHub secret shown under
**Secrets**. Any sandbox type can use it for matching GitHub requests, such as
cloning a private repository or pushing changes.

## Configure a secret

To configure a credential before launching a sandbox:

1. Open **Secrets** and find the service.
2. Select the edit icon for that service.
3. Enter its API key or token and save the value.

You can also provide a required credential from **New** when you launch a
sandbox.

## Manage secrets

The **Secrets** page shows each secret ID and the sandboxes that use it. Select
the copy icon to copy the secret ID, or select the edit icon to change the
stored value.

Agents see the secret ID, not the stored value. When an outbound request matches
the secret's service binding, the sandbox proxy applies the value to the
request. A secret for one service does not become a general-purpose credential
inside the sandbox.

Do not put API keys or tokens in prompts or files inside the sandbox. Rotate or
revoke a credential at its provider when it is no longer needed.
