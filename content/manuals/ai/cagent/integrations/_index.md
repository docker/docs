---
title: Integrations
description: Connect cagent agents to editors, MCP clients, and other agents
keywords: [cagent, integration, acp, mcp, a2a, editor, protocol]
weight: 60
---

cagent agents can integrate with different environments depending on how you
want to use them. Each integration type serves a specific purpose.

## Integration types

### ACP - Editor integration

Run cagent agents directly in your editor (Neovim, Zed). The agent sees your
editor's file context and can read and modify files through the editor's
interface. Use ACP when you want an AI coding assistant embedded in your
editor.

See [ACP integration](./acp.md) for setup instructions.

### MCP - Tool integration

Expose cagent agents as tools in MCP clients like Claude Desktop or Claude
Code. Your agents appear in the client's tool list, and the client can call
them when needed. Use MCP when you want Claude Desktop (or another MCP client)
to have access to your specialized agents.

See [MCP integration](./mcp.md) for setup instructions.

### A2A - Agent-to-agent communication

Run cagent agents as HTTP servers that other agents or systems can call using
the Agent-to-Agent protocol. Your agent becomes a service that other systems
can discover and invoke over the network. Use A2A when you want to build
multi-agent systems or expose your agent as an HTTP service.

See [A2A integration](./a2a.md) for setup instructions.

## Choosing the right integration

| Feature       | ACP                | MCP                | A2A                  |
| ------------- | ------------------ | ------------------ | -------------------- |
| Use case      | Editor integration | Agents as tools    | Agent-to-agent calls |
| Transport     | stdio              | stdio/SSE          | HTTP                 |
| Discovery     | Editor plugin      | Server manifest    | Agent card           |
| Best for      | Code editing       | Tool integration   | Multi-agent systems  |
| Communication | Editor calls agent | Client calls tools | Between agents       |

Choose ACP if you want your agent embedded in your editor while you code.
Choose MCP if you want Claude Desktop (or another MCP client) to be able to
call your specialized agents as tools. Choose A2A if you're building
multi-agent systems where agents need to call each other over HTTP.
