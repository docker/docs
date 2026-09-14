---
title: Authenticate command-line tools
linkTitle: Authenticate tools
weight: 30
description: Use authenticated command-line tools and dynamic secret sources inside Docker Sandboxes.
keywords: docker sandboxes, sbx, authentication, github cli, registry, 1password, aws secrets manager
---

The sandbox proxy handles API credentials for model providers automatically,
but agents often also need credentials for tools like `gh`, `docker`, or a
secrets manager. Configure the credential source on your host, and the proxy
injects the resolved value into matching requests from the sandbox. Dynamic
secret sources can retrieve a value from an authenticated host CLI without
copying the value into the secret store.

> [!NOTE]
> Service secrets are global by default, so all future sandboxes can use them.
> Sandboxes that already exist when you run `sbx secret set` do not
> receive the updated value. To update a running sandbox, scope the secret to
> it directly: `sbx secret set <service> --sandbox <sandbox-name>`.

## GitHub CLI

Store your GitHub token as a sandbox secret. The proxy injects it into
outbound requests, so `gh` works inside the sandbox without any additional
configuration:

```console
$ sbx secret set github --command 'gh auth token'
```

The daemon runs `gh auth token` on the host and caches its output for 55 minutes
by default. After the cache expires, it runs the command again, so token updates
from `gh` don't need to be copied into `sbx` manually. Use `--refresh on-demand`
to run the command for every credential use.

The agent can then create pull requests, open issues, comment on PRs, and
interact with the GitHub API the same way it would from your host:

```console
# Inside the sandbox
$ gh pr create --title "feat: my feature" --body "..."
$ gh issue list
```

The token is never stored in plaintext inside the sandbox. See
[GitHub token](../configuration/credentials.md#github-token) for details.

## Docker registry

When using Docker Hub, authentication is handled automatically; `sbx` reuses
your existing login session. For other registries, you need to configure
credentials for `sbx` so it can pull private [templates](../customize/templates.md)
and kits when creating a sandbox:

```console
$ gh auth token | sbx secret set --all-sandboxes --registry ghcr.io \
    --username <github-username> --password-stdin
$ echo "$ACR_PASSWORD" | sbx secret set --all-sandboxes \
    --registry myregistry.azurecr.io \
    --username myuser --password-stdin
```

Add `-g` or a sandbox name when the agent needs to run authenticated
`docker pull` or `docker push` commands from inside the sandbox. The host-side
proxy handles the registry login without writing the credential into the
sandbox.

Images and containers built inside the sandbox run on the sandbox's private
Docker daemon, not your host's. They're deleted when the sandbox is removed.

For information on how registry credentials differ from other secrets,
per-registry username requirements, and all-sandbox versus per-sandbox scoping, see
[Registry credentials](../configuration/credentials.md#registry-credentials).

## Source credentials from 1Password

Install the 1Password CLI, sign in on the host, and pass an `op://` reference to
`sbx secret set`. The secret store records the reference, and the daemon uses
`op read` on the host when the proxy needs the credential:

```console
$ sbx secret set github --ref 'op://Work/GitHub/token'
$ sbx secret set anthropic --ref 'op://Work/Anthropic/credential'
```

The real value stays on your host, and the sandbox sees the proxy-managed
placeholder. Service secrets are cached for 55 minutes by default. To retrieve
the value from 1Password for every credential use, set the refresh policy:

```console
$ sbx secret set anthropic \
    --ref 'op://Work/Anthropic/credential' \
    --refresh on-demand
```

## Source credentials from AWS Secrets Manager

Install and authenticate the AWS CLI on the host, then register the secret's
ARN. The daemon calls AWS Secrets Manager when the proxy needs the value:

```console
$ sbx secret set anthropic \
    --ref 'arn:aws:secretsmanager:us-west-2:123456789012:secret:anthropic-api-key'
```

See [Use a dynamic secret source](../configuration/credentials.md#use-a-dynamic-secret-source)
for refresh policies, verification options, custom secrets, and provider
account or profile selection.
