---
title: Secrets and credential injection
linkTitle: Secrets
description: Understand where sandbox credentials live and how agents use them through a proxy.
keywords: docker sandboxes, sbx, cloud sandboxes
weight: 40
---

A managed secret lets a sandbox use a credential without storing its value in
the sandbox filesystem. The agent sends a request through a proxy, which
applies the credential for the configured service or destination.

Docker account authentication and agent authentication serve different
purposes. Signing in to Docker grants access to sandbox resources. An agent
still needs credentials for its model provider or other services.

## Where credentials live

Local sandboxes use a credential store on your host. Cloud sandboxes use a
cloud secret store. Local credentials are not copied to the cloud when you
choose `--cloud`.

The cloud CLI and Docker Agentic Platform Console share the same cloud secret
store. They provide different controls for selecting and attaching credentials;
they do not create independent stores for each interface.

| Environment | Store | Setup instructions |
| --- | --- | --- |
| Local | Host OS credential store, with a file fallback on some Linux hosts | [Local CLI credentials](/manuals/ai/sandboxes/cli/access/credentials-local.md) |
| Cloud | Cloud secret store | [Cloud CLI credentials](/manuals/ai/sandboxes/cli/access/credentials-cloud.md), [Console secrets](/manuals/ai/sandboxes/console/secrets.md), or [SDK secrets](/manuals/ai/sandboxes/api/cookbook/manage-cloud-secrets.md) |

## How injection works

With proxy-managed authentication, the sandbox can receive a placeholder or
sentinel instead of the real API key or token. The proxy applies the real
credential to requests matching the configured service or destination.
Configuring a credential does not grant network access: the request must also
pass [network policy](/manuals/ai/sandboxes/concepts/network.md).

Custom secrets describe which hosts may receive a credential and how to format
it in a request. Follow the relevant setup guide for the header, value format,
and placeholder your client needs. A placeholder is not itself the secret.

MCP gateway credentials authenticate the gateway's connection to a server.
They are distinct from credentials used by an agent to call a model provider.
See [MCP connections](/manuals/ai/sandboxes/concepts/mcp.md).

## Scope and attachment

Storage, scope, and attachment are separate choices:

- Storage determines where the credential value is held.
- Scope determines which sandboxes can use it.
- Attachment or kit configuration determines which credentials a sandbox uses.

The CLI supports a general scope and a sandbox-specific scope for service
secrets. In local mode the general scope is on the host; in cloud mode it is
in the Docker account. A sandbox-scoped credential takes precedence for the
same service. Cloud OAuth credentials use account scope.

The Console launcher selects credentials supported by the chosen kit. SDK
workflows attach stored secrets by their resource names. Follow the interface's
instructions rather than assuming an identifier shown in one interface can be
passed unchanged to another.

## Keep credentials outside sandbox state

Proxy-managed secrets stay outside the sandbox filesystem. Credentials you
write to a file, paste into a prompt, or obtain through an agent's interactive
sign-in can enter sandbox state. Templates and snapshots can capture such
files. A kit's OAuth passthrough option can also expose tokens inside the
sandbox.

Use managed secrets where supported, remove credential files before sharing
captured state, and rotate or revoke credentials at their provider when they
are no longer needed.
