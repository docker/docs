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
 - title: Get started with MCP Toolkit
   description: Learn how to quickly install and use the MCP Toolkit to set up servers and clients.
   icon: explore
   link: /ai/mcp-catalog-and-toolkit/get-started/
 - title: MCP Catalog
   description: Learn about the benefits of the MCP Catalog, how you can use it, and how you can contribute
   icon: hub
   link: /ai/mcp-catalog-and-toolkit/catalog/
 - title: MCP Toolkit
   description: Learn about the MCP Toolkit to manage MCP servers and clients
   icon: /icons/toolkit.svg
   link: /ai/mcp-catalog-and-toolkit/toolkit/
 - title: Dynamic MCP
   description: Discover and add MCP servers on-demand using natural language
   icon: search
   link: /ai/mcp-catalog-and-toolkit/dynamic-mcp/
 - title: MCP Gateway
   description: Learn about the underlying technology that powers the MCP Toolkit
   icon: developer_board
   link: /ai/mcp-catalog-and-toolkit/mcp-gateway/
 - title: Docker Hub MCP server
   description: Explore about the Docker Hub server for searching images, managing repositories, and more
   icon: device_hub
   link: /ai/mcp-catalog-and-toolkit/hub-mcp/
---

{{< summary-bar feature_name="Docker MCP Catalog and Toolkit" >}}

[Model Context Protocol](https://modelcontextprotocol.io/introduction) (MCP) is
an open protocol that standardizes how AI applications access external tools
and data sources. By connecting LLMs to local development tools, databases,
APIs, and other resources, MCP extends their capabilities beyond their base
training.

Through a client-server architecture, applications such as Claude, ChatGPT, and
[Gordon](/manuals/ai/gordon/_index.md) act as clients that send requests to MCP
servers, which then process these requests and deliver the necessary context to
AI models.

MCP servers extend the utility of AI applications, but running servers locally
also presents several operational challenges. Typically, servers must be
installed directly on your machine and configured individually for each
application. Running untrusted code locally requires careful vetting, and the
responsibility of keeping servers up-to-date and resolving environment
conflicts falls on the user.

## Docker MCP features

Docker provides three integrated components that address the challenges of
running local MCP servers:

MCP Catalog
: A curated collection of verified MCP servers, packaged and distributed as
container images via Docker Hub. All servers are versioned, come with full
provenance and SBOM metadata, and are continuously maintained and updated with
security patches.

MCP Toolkit
: A graphical interface in Docker Desktop for discovering, configuring, and
managing MCP servers. The Toolkit provides a unified way to search for servers,
handle authentication, and connect them to AI applications.

MCP Gateway
: The core open source component that powers the MCP Toolkit. The MCP Gateway
manages MCP containers provides a unified endpoint that exposes your enabled
servers to all AI applications you use.

This integrated approach ensures:

- Simplified discovery and setup of trusted MCP servers from a curated catalog
  of tools
- Centralized configuration and authentication from within Docker Desktop
- A secure, consistent execution environment by default
- Improved performance since applications can share a single server runtime,
  compared to having to spin up duplicate servers for each application.

![MCP overview](./images/mcp-overview.svg)

## Learn more

{{< grid >}}
