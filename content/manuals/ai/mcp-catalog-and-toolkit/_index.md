---
title: Docker MCP Catalog and Toolkit
linkTitle: MCP Catalog and Toolkit
params:
  sidebar:
    group: AI
    badge:
      color: blue
      text: Beta
weight: 10
description: Learn about Docker's MCP catalog on Docker Hub
keywords: Docker, ai, mcp servers, ai agents, extension, docker desktop, llm, docker hub
grid:
  - title: Get started with MCP Toolkit
    description: Learn how to quickly install and use the MCP Toolkit to set up servers and clients.
    icon: explore
    link: /ai/mcp-catalog-and-toolkit/get-started/
  - title: MCP Catalog
    description: Browse Docker's curated collection of verified MCP servers
    icon: hub
    link: /ai/mcp-catalog-and-toolkit/catalog/
  - title: MCP Toolkit
    description: Learn about the MCP Toolkit to manage MCP servers and clients
    icon: /icons/toolkit.svg
    link: /ai/mcp-catalog-and-toolkit/toolkit/
  - title: MCP Gateway
    description: Learn about the underlying technology that powers the MCP Toolkit
    icon: developer_board
    link: /ai/mcp-catalog-and-toolkit/mcp-gateway/
  - title: Dynamic MCP
    description: Discover and add MCP servers on-demand using natural language
    icon: search
    link: /ai/mcp-catalog-and-toolkit/dynamic-mcp/
  - title: Docker Hub MCP server
    description: Use the Docker Hub MCP server to search images and manage repositories
    icon: device_hub
    link: /ai/mcp-catalog-and-toolkit/hub-mcp/
  - title: Security FAQs
    description: Common questions about MCP security, credentials, and server verification
    icon: security
    link: /ai/mcp-catalog-and-toolkit/faqs/
  - title: E2B sandboxes
    description: Cloud sandboxes for AI agents with built-in MCP Catalog access
    icon: cloud
    link: /ai/mcp-catalog-and-toolkit/e2b-sandboxes/
---

{{< summary-bar feature_name="Docker MCP Catalog and Toolkit" >}}

[Model Context Protocol](https://modelcontextprotocol.io/introduction) (MCP) is
an open protocol that standardizes how AI applications access external tools
and data sources. By connecting LLMs to local development tools, databases,
APIs, and other resources, MCP extends their capabilities beyond their base
training.

The challenge is that running MCP servers locally creates operational friction.
Each server requires separate installation and configuration for every
application you use. You run untrusted code directly on your machine, manage
updates manually, and troubleshoot dependency conflicts yourself. Configure a
GitHub server for Claude, then configure it again for Cursor, and so on. Each
time you manage credentials, permissions, and environment setup.

## Docker MCP features

Docker solves these challenges by packaging MCP servers as containers and
providing tools to manage them centrally. Docker provides three integrated
components: the [MCP Catalog](/ai/mcp-catalog-and-toolkit/catalog/) for
discovering servers, the [MCP Gateway](/ai/mcp-catalog-and-toolkit/mcp-gateway/)
for running them, and the [MCP Toolkit](/ai/mcp-catalog-and-toolkit/toolkit/)
for managing everything through Docker Desktop.

The [MCP Catalog](/ai/mcp-catalog-and-toolkit/catalog/) is where you find
servers. Docker maintains 300+ verified servers, packaged as container images
with versioning, provenance, and security updates. Servers run isolated in
containers rather than directly on your machine. Organizations can create
[custom catalogs](/ai/mcp-catalog-and-toolkit/catalog/#custom-catalogs) with
approved servers for their teams.

The [MCP Gateway](/ai/mcp-catalog-and-toolkit/mcp-gateway/) runs your servers
and routes requests from AI applications to the right server. It handles
containerized servers, remote servers, authentication, and lifecycle
management. Every AI application connects to the Gateway, which means you
configure credentials and permissions once instead of per-application.

The [MCP Toolkit](/ai/mcp-catalog-and-toolkit/toolkit/) provides a graphical
interface in Docker Desktop for browsing catalogs, enabling servers, and
connecting clients. You can also use the `docker mcp` CLI to manage everything
from the terminal.

![MCP overview](./images/mcp-overview.svg)

## Learn more

{{< grid >}}
