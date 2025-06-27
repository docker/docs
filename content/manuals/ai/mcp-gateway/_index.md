---
title: MCP Gateway
description: "Docker's MCP Gateway provides secure, centralized, and scalable orchestration of AI tools through containerized MCP servers—empowering developers, operators, and security teams."
keywords: MCP Gateway
params:
  sidebar:
    group: Open source
---

The MCP Gateway is Docker's open-source enterprise-ready solution for orchestrating and
managing [Model Context Protocol (MCP)](https://spec.modelcontextprotocol.io/) servers
securely across development and production environments.
It is designed to help organizations connect MCP servers from the [Docker MCP Catalog](https://hub.docker.com/mcp) to MCP Clients without compromising security, visibility, or control.

By unifying multiple MCP servers into a single, secure endpoint, the MCP Gateway offers
the following benefits:

- Secure by Default: MCP servers run in isolated Docker containers with restricted
  privileges, network access, and resource usage.
- Unified Management: One gateway endpoint centralizes configuration, credentials,
  and access control for all MCP servers.
- Enterprise Observability: Built-in monitoring, logging, and filtering tools ensure
  full visibility and governance of AI tool activity.

## Who is the MCP Gateway designed for?

The MCP Gateway solves problems encountered by various groups:

- Developers: Deploy MCP servers locally and in production using Docker Compose,
  with built-in support for protocol handling, credential management, and security policies.
- Security Teams: Achieve enterprise-grade isolation and visibility into AI tool
  behavior and access patterns.
- Operators: Scale effortlessly from local development environments to production
  infrastructure with consistent, low-touch operations.

## Key features

- Server Management: List, inspect, and call MCP tools, resoures and prompts from multiple servers
- Container-based Servers: Run MCP servers as Docker containers with proper isolation
- Secrets Management: Secure handling of API keys and credentials via Docker Desktop
- Server Catalog: Manage and configure multiple MCP catalogs
- Dynamic Discovery and Reloading: Automatic tool, prompt, and resource discovery from running servers
- Monitoring: Built-in logging and call tracing capabilities

## Install the MCP Gateway

### Prerequisites

- [Docker Engine](/manuals/engine/_index.md)
- Go 1.24+ (for development)

### Install as Docker CLI Plugin

The MCP CLI is already installed on recent versions of Docker Desktop.
To update to the latest version:

```bash
# Clone the repository
git clone https://github.com/docker/docker-mcp.git
cd docker-mcp

# Build and install the plugin
make docker-mcp
```

You can now use the `mcp` command:

```bash
docker mcp --help
```

## Use the MCP Gateway

To view all the commands and configuration options, go to the [docker-mcp repository](https://github.com/docker/docker-mcp).

## Related pages

- [Docker MCP toolkit and catalog](/manuals/ai/mcp-catalog-and-toolkit/_index.md)
