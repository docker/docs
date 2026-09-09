---
title: Authenticate cloud agents
linkTitle: Credentials
description: Configure cloud-specific API keys and OAuth credentials for agents without storing credentials in the cloud sandbox filesystem.
keywords: docker sandboxes, cloud credentials, cloud secrets, oauth, anthropic, openai
weight: 30
---

Cloud agents authenticate with credentials from the cloud secret store. Set up
these credentials before launching an agent so authentication doesn't depend
on files stored inside the sandbox.

## Use the cloud secret store

The following secret interfaces are separate:

- `sbx secret` manages secrets for local sandboxes
- `sbx --cloud secret` manages secrets for cloud sandboxes created by the CLI
- Docker Agentic Platform manages its own secret names through its web
  interface

A credential created through one interface isn't available through the
others. If you already configured an Anthropic or OpenAI credential for local
sandboxes, configure it again with `sbx --cloud secret` before starting a
cloud sandbox.

## Choose an authentication method

Claude Code in cloud sandboxes requires an Anthropic API key. Anthropic OAuth
and Claude subscription sign-in are not supported for cloud sandboxes.

Configure your agent with the corresponding command:

| Agent or provider | Recommended command | Authentication |
| --- | --- | --- |
| Claude Code | `sbx --cloud secret set anthropic` | Prompts for an Anthropic API key |
| Codex with OpenAI OAuth | `sbx --cloud secret set openai --oauth` | Opens the OpenAI OAuth flow and stores the resulting credential at account scope |
| Service API key | `sbx --cloud secret set <service>` | Prompts for an API key or token |

For Codex, you can use OpenAI OAuth at account scope or store an OpenAI API
key with `sbx --cloud secret set openai`.

After storing the Anthropic API key, launch Claude Code:

```console
$ sbx --cloud run claude --name cloud-project
```

## Keep credentials out of the sandbox filesystem

Credentials configured with `sbx --cloud secret` stay in the cloud secret store,
outside the sandbox filesystem.

An agent's interactive sign-in can write credentials inside the sandbox.
Those files can be included in templates and `sbx move` snapshots. If you
signed in inside an agent, follow the provider's sign-out guidance and remove
those credentials before capturing or moving the sandbox.

## Set credential scope

Credentials use account scope by default. An account-scoped credential is
available to cloud sandboxes in the Docker account:

```console
$ sbx --cloud secret set github
```

Scope an API key or token to one sandbox by name:

```console
$ sbx --cloud secret set openai --sandbox cloud-project
```

A sandbox-scoped secret takes precedence over an account-scoped secret for the
same service. Set it before creating the sandbox so the CLI can include it when
the sandbox starts. OAuth credentials can't use sandbox scope.

## Service identifiers

The following table shows cloud secret support for the
[built-in services documented for local sandboxes](../configuration/credentials.md#built-in-services):

| Service | Cloud secret authentication |
| --- | --- |
| `anthropic` | API key |
| `cursor` | Not supported |
| `droid` | API key; cloud-managed OAuth is not supported |
| `github` | Token |
| `google` | API key |
| `groq` | API key |
| `mistral` | API key |
| `nebius` | API key |
| `openai` | API key or OAuth |
| `openrouter` | Not supported |
| `xai` | API key |

To configure a supported service, run `sbx --cloud secret set <service>`.
For OpenAI OAuth, add `--oauth` and use account scope. An agent's interactive
sign-in does not use the cloud secret store.

Services declared by local kits aren't automatically supported in cloud
sandboxes. A kit cannot add a service to the cloud secret store.

Registry credentials, dynamic `--ref` values, and host-run `--command`
resolvers aren't cloud secret workflows.

## List and remove credentials

List cloud secret metadata without revealing values:

```console
$ sbx --cloud secret ls
```

Remove a credential from the matching scope:

```console
$ sbx --cloud secret rm github
$ sbx --cloud secret rm openai --sandbox cloud-project
```

Cloud mode doesn't support removing every secret in one operation.
