---
title: MCP
description: Connect predefined and custom MCP servers for Docker Agentic Platform sandboxes.
keywords: docker agentic platform, mcp servers, mcp tools, authorization, remote url
weight: 30
aliases:
  - /agentic-platform/concepts/mcp-tools/
  - /agentic-platform/guides/configure-mcp-tools/
---

Connect Model Context Protocol (MCP) servers to give your agent tools for
working with external services. From **MCP**, choose a predefined server or
add a custom server by URL. Authorize access if prompted.

You can also select or connect servers from the tools control in the sandbox
launcher. Select a server to use it in the sandbox, and authorize access if
prompted. To connect a custom server there, enter its URL and
select **Connect**.

Connecting an MCP server doesn't restrict or inspect the sandbox's other
network traffic. To control which hosts and services the sandbox can reach,
use [network policies](/manuals/agentic-platform/policies.md). The initial
release doesn't support MCP-specific policies.

## Connect a predefined server

1. Open **MCP** and choose a predefined server.
2. Connect the server.
3. Authorize access if prompted.

## Add a server by URL

To connect a server that is not predefined:

1. Open **MCP**.
2. Choose the option to add a server and enter its URL.
3. Connect the server and authorize access if prompted.

## Use tools from another client

On the **MCP** page, use **MCP Gateway endpoint URL** to connect your external
MCP client to the gateway. Select **Add to your client**, choose your
client, and follow the connection and authorization instructions. Supported
options include VS Code and Codex CLI. The VS Code instructions include an
install link and a manual configuration example.
