---
title: Docker MCP Catalog and Toolkit
linkTitle: MCP Catalog and Toolkit
params:
  sidebar:
    group: AI
    badge:
      color: blue
      text: Beta
weight: 30
description: Learn about Docker's MCP catalog on Docker Hub
keywords: Docker, ai, mcp servers, ai agents, extension, docker desktop, llm, docker hub
grid:
 - title: MCP Catalog
   description: Learn about the benefits of the MCP Catalog, how you can use it, and how you can contribute
   icon: hub
   link: /ai/mcp-catalog-and-toolkit/catalog/
 - title: MCP Toolkit
   description: Learn about the MCP Toolkit to manage MCP servers and clients
   icon: /icons/toolkit.svg
   link: /ai/mcp-catalog-and-toolkit/toolkit/
---

## MCP servers

The Model Context Protocol (MCP) is a modern standard that transforms AI agents
from passive responders into action-oriented systems. By standardizing how tools
are described, discovered, and invoked, MCP enables agents to securely query
APIs, access data, and run services across different environments.

As agents move into production, MCP solves common integration challenges —
interoperability, reliability, and security — by providing a consistent,
decoupled, and scalable interface between agents and tools. Just as containers
redefined software deployment, MCP is reshaping how AI systems interact with the
world.

> **MCP servers in simple terms**
>
> An MCP server is a way for an LLM to interact with an external system.
>
> For example:
> If you ask a model to create a meeting, it needs to communicate with your calendar app to do that.
> An MCP server for your calendar app provides _tools_ that perform atomic actions, such as:
> "getting the details of a meeting" or "creating a new meeting".

## Docker MCP Catalog and Toolkit

Docker MCP Catalog and Toolkit is a solution for securely building, sharing, and
running MCP tools. It simplifies the developer experience across these areas:

- Discovery: A central catalog with verified, versioned tools.
- Credential management: OAuth-based and secure by default.
- Execution: Tools run in isolated, containerized environments.
- Portability: Use MCP tools across Claude, Cursor, VS Code, and more—no code
  changes needed.

With Docker Hub and the MCP Toolkit, you can:

- Launch MCP servers in seconds.
- Add tools using the CLI or GUI.
- Rely on Docker's pull-based infrastructure for trusted delivery.

{{< grid >}}
