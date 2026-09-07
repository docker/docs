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

Use the authentication method supported by the provider:

| Agent or provider | Recommended command | Authentication |
| --- | --- | --- |
| Claude Code with Anthropic OAuth | `sbx --cloud secret set anthropic --oauth` | Opens the Anthropic OAuth flow and stores the resulting credential at account scope |
| Codex with OpenAI OAuth | `sbx --cloud secret set openai --oauth` | Opens the OpenAI OAuth flow and stores the resulting credential at account scope |
| Service API key | `sbx --cloud secret set <service>` | Prompts for an API key or token |

OAuth setup is available for Anthropic and OpenAI at account scope. To use an
API key instead, omit `--oauth`:

```console
$ sbx --cloud secret set anthropic
```

After configuring the credential, launch the agent:

```console
$ sbx --cloud run claude --name cloud-project
```

## Keep credentials out of the sandbox filesystem

> [!IMPORTANT]
>
> An agent's interactive sign-in doesn't use the cloud secret store. It can
> write credentials to the sandbox filesystem. For example, Claude Code can
> store its credential in plain text at `~/.claude/.credentials.json`.

Files created by interactive sign-in remain part of the sandbox filesystem.
They can persist until the sandbox is deleted and can be included in a
filesystem snapshot created by `sbx move`.

When a cloud-managed OAuth flow is available, use it instead of the agent's
`/login` command. For Claude Code, run:

```console
$ sbx --cloud secret set anthropic --oauth
```

The actual credential remains in the cloud secret store. Agent credential files
can contain non-secret sentinel values that direct requests through the
credential proxy.

If you already signed in from inside an agent, follow the provider's sign-out
guidance and remove the in-sandbox credential before configuring the
cloud-managed credential.

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
| `anthropic` | API key or OAuth |
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
Cloud-managed OAuth is available only for `anthropic` and `openai`, at account
scope. Other providers' interactive sign-in flows do not use the cloud secret
store.

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
