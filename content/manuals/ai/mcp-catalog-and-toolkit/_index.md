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
  - title: MCP Profiles
    description: Organize servers into profiles for different projects and share configurations
    icon: folder
    link: /ai/mcp-catalog-and-toolkit/profiles/
  - title: MCP Toolkit
    description: Use Docker Desktop's UI to discover, configure, and manage MCP servers
    icon: /icons/toolkit.svg
    link: /ai/mcp-catalog-and-toolkit/toolkit/
  - title: MCP Gateway
    description: Use the CLI and Gateway to run MCP servers with custom configurations
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

The [MCP Toolkit](/ai/mcp-catalog-and-toolkit/toolkit/) and [MCP
Gateway](/ai/mcp-catalog-and-toolkit/mcp-gateway/) solve these challenges
through centralized management. Instead of configuring each server for every AI
application separately, you set things up once and connect all your clients to
it. The workflow centers on three concepts: catalogs, profiles, and clients.

![MCP overview](./images/mcp_toolkit.avif)

[Catalogs](/ai/mcp-catalog-and-toolkit/catalog/) are curated collections of
MCP servers. The Docker MCP Catalog provides 300+ verified servers packaged as
container images with versioning, provenance, and security updates. Organizations
can create [custom
catalogs](/ai/mcp-catalog-and-toolkit/catalog/#custom-catalogs) with approved
servers for their teams.

[Profiles](/ai/mcp-catalog-and-toolkit/profiles/) organize servers into named
collections for different projects. Your "web-dev" profile might use GitHub and
Playwright; your "backend" profile, database tools. Profiles support both
containerized servers from catalogs and remote MCP servers. Configure a profile
once, then share it across clients or with your team.

Clients are the AI applications that connect to your profiles. Claude Code,
Cursor, Zed, and others connect through the MCP Gateway, which routes requests
to the right server and handles authentication and lifecycle management.

## Learn more

{{< grid >}}
