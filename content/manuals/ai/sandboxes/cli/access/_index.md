---
title: Configure access and connections
linkTitle: Access and connections
description: Configure sandbox credentials, network permissions, and MCP connections with sbx.
keywords: docker sandboxes, sbx, local sandboxes, cloud sandboxes
weight: 30
---

Give your sandbox credentials for the services it needs, configure network
access, and connect MCP tools. These are separate tasks: storing a credential
does not grant network access, and registering an MCP server does not attach it
to a sandbox.

## Credentials

Configure [local credentials](/manuals/ai/sandboxes/cli/access/credentials-local.md)
or [cloud credentials](/manuals/ai/sandboxes/cli/access/credentials-cloud.md)
before launching an agent. For GitHub CLI, registries, and external secret
providers, see
[Authenticate command-line tools](/manuals/ai/sandboxes/cli/access/authentication.md).

[Secrets and credential injection](/manuals/ai/sandboxes/concepts/secrets.md)
explains storage, scope, and how the proxy uses credentials.

## Network access

Set [local network rules](/manuals/ai/sandboxes/cli/access/network-local.md)
or [cloud network rules](/manuals/ai/sandboxes/cli/access/network-cloud.md)
and inspect connection decisions when a request is blocked.

For centrally enforced organization policies, see
[AI Governance](/manuals/ai/governance/_index.md). To route permitted local
traffic through a corporate proxy, see
[Upstream proxy configuration](/manuals/ai/sandboxes/cli/local/upstream-proxy.md).

## MCP tools

[Connect MCP servers](/manuals/ai/sandboxes/cli/access/mcp.md) covers both local
registration and cloud account setup, authorization, and attachment.
