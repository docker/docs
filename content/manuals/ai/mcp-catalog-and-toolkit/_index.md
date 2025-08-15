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

{{< summary-bar feature_name="Docker MCP Catalog and Toolkit" >}}

Docker MCP Catalog and Toolkit is a solution for securely building, sharing, and
running MCP tools.

It simplifies the developer experience across these areas:

- Discovery: A central catalog with verified, versioned tools.
- Credential management: OAuth-based and secure by default.
- Execution: Tools run in isolated, containerized environments.
- Portability: Use MCP tools across Claude, Cursor, VS Code, and more—no code
  changes needed.

With Docker Hub and the MCP Toolkit, you can:

- Launch MCP servers in seconds.
- Add tools using the CLI or GUI.
- Rely on Docker's pull-based infrastructure for trusted delivery.

## MCP servers

MCP servers are systems that use the [Model Context Protocol](https://www.anthropic.com/news/model-context-protocol) (MCP) to help manage
and run AI or machine learning models more efficiently. MCP allows different
parts of a system, like the model, data, and runtime environment, to
communicate in a standardized way. You can see them as
add-ons that provide specific tools to an LLM.

> [!TIP]
> **Example**:
> If you ask a model to create a meeting, it needs to communicate with your calendar app to do that.
>
> An MCP server provided by your calendar app provides _tools_ to your model to perform atomic
> actions, like:
>
> - `get the details of a meeting`
> - `create a new meeting`
> - ...


{{< grid >}}
