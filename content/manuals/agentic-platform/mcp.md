---
title: MCP
description: Connect predefined and custom MCP servers for Docker Agentic Platform sandboxes.
keywords: docker agentic platform, mcp servers, mcp tools, authorization, remote url
weight: 30
aliases:
  - /agentic-platform/concepts/mcp-tools/
  - /agentic-platform/guides/configure-mcp-tools/
---

Model Context Protocol (MCP) servers connect agents to external services and
expose operations from those services as tools. From the **MCP** page, connect a
predefined server or add a custom server by URL. Complete authorization when a
server requires it.

MCP configuration grants agents tools they can invoke. It does not restrict or
inspect ordinary network access from a sandbox. Use
[network policies](/manuals/agentic-platform/policies.md) to control outbound
destinations. MCP-specific policies are not part of the initial release.

## Connect a predefined server

1. Open **MCP** and choose a predefined server.
2. Connect the server.
3. Complete authorization if prompted.

## Add a server by URL

To connect a server that is not predefined:

1. Open **MCP**.
2. Choose the option to add a server and enter its URL.
3. Connect the server and complete authorization if prompted.
