---
title: MCP connections and gateways
linkTitle: MCP connections
description: Understand MCP registration, authorization, and attachment for local and cloud sandboxes.
keywords: docker sandboxes, sbx, cloud sandboxes
weight: 60
---

An MCP gateway gives a sandboxed agent access to tools from Model Context
Protocol servers through one endpoint. The gateway manages connections to
those servers on the agent's behalf.

## Registration, authorization, and attachment

Connecting a server involves three separate choices:

1. Register or configure the server so the gateway knows how to reach it.
2. Authorize the gateway to use it, if the server requires credentials.
3. Attach the server to the sandbox so its tools are available to the agent.

Registering a server does not necessarily attach it. Selecting a server does
not replace the authorization required by its provider.

## Local and cloud gateways

For local sandboxes, `sbx mcp add` stores server registrations on the host.
A server can be a remote endpoint or a stdio process launched on the host.
Host-run servers are outside the sandbox isolation boundary.

For cloud CLI workflows, connect a server through Docker Agentic Platform
using the same Docker account, then load it into the sandbox. Cloud gateways
do not read registrations from your local `sbx mcp` store.

The local gateway's static mode preloads selected servers. Dynamic mode lets
supported agents discover and attach registered servers during a session.
Do not assume those local discovery controls apply to cloud gateways.

See [Connect MCP servers with the CLI](/manuals/ai/sandboxes/cli/access/mcp.md),
[Console MCP setup](/manuals/ai/sandboxes/console/mcp.md), or the
[SDK gateway recipe](/manuals/ai/sandboxes/api/cookbook/give-a-sandbox-an-mcp-gateway.md).

## Credentials and policy

The gateway uses server credentials outside the sandbox. These are separate
from the credentials an agent uses to call its model provider.

MCP governance can control local server registration and tool use. Cloud
network policies control network access and do not establish equivalent
MCP-specific governance. See [Network access and policies](/manuals/ai/sandboxes/concepts/network.md).

The Docker Sandboxes gateway is separate from Docker Desktop's MCP Toolkit.
Their server configurations are not shared.
