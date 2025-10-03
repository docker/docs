---
title: Docker MCP Catalog
linkTitle: MCP Catalog
description: Learn about the benefits of the MCP Catalog, how you can use it, and how you can contribute
keywords: docker hub, mcp, mcp servers, ai agents, catalog, docker
weight: 20
---

{{< summary-bar feature_name="Docker MCP Catalog" >}}

The [Docker MCP Catalog](https://hub.docker.com/mcp) is a centralized, trusted
registry for discovering, sharing, and running MCP-compatible tools. Integrated
with Docker Hub, it offers verified, versioned, and curated MCP servers
packaged as Docker images. The catalog is also available in Docker Desktop.

The catalog solves common MCP server challenges:

- Environment conflicts. Tools often need specific runtimes that might clash
  with existing setups.
- Lack of isolation. Traditional setups risk exposing the host system.
- Setup complexity. Manual installation and configuration slow adoption.
- Inconsistency across platforms. Tools might behave unpredictably on different
  operating systems.

With Docker, each MCP server runs as a self-contained container. This makes it
portable, isolated, and consistent. You can launch tools instantly using the
Docker CLI or Docker Desktop, without worrying about dependencies or
compatibility.

## Key features

- Extensive collection of verified MCP servers in one place.
- Publisher verification and versioned releases.
- Pull-based distribution using Docker infrastructure.
- Tools provided by partners such as New Relic, Stripe, Grafana, and more.

## How it works

Each tool in the MCP Catalog is packaged as a Docker image with metadata.

- Discover tools on Docker Hub under the `mcp/` namespace.
- Connect tools to your preferred agents with simple configuration through the
  [MCP Toolkit](toolkit.md).
- Pull and run tools using Docker Desktop or the CLI.

Each catalog entry displays:

- Tool description and metadata.
- Version history.
- List of tools provided by the MCP server.
- Example configuration for agent integration.

## Server deployment types

The Docker MCP Catalog supports both local and remote server deployments, each optimized for different use cases and requirements.

### Local MCP servers

Local MCP servers are containerized applications that run directly on your machine. All local servers are built and digitally signed by Docker, providing enhanced security through verified provenance and integrity. These servers run as containers on your local environment and function without internet connectivity once downloaded. Local servers display a Docker icon {{< inline-image src="../../desktop/images/whale-x.svg" alt="docker whale icon" >}} to indicate they are built by Docker.

Local servers offer predictable performance, complete data privacy, and independence from external service availability. They work well for development workflows, sensitive data processing, and scenarios requiring offline functionality.

### Remote MCP servers

Remote MCP servers are hosted services that you connect to through the internet. Service providers maintain and update these servers, ensuring access to current features and live data without requiring local updates or maintenance. Remote servers display a cloud icon {{< inline-image src="../../offload/images/cloud-mode.png" alt="cloud icon" >}} to indicate their hosted nature and external connectivity requirements.

Remote servers excel when you need always-current data, want to minimize local resource usage, or require capabilities that benefit from provider-managed infrastructure and scaling.

## Use an MCP server from the catalog

To use an MCP server from the catalog, see [MCP Toolkit](toolkit.md).

## Contribute an MCP server to the catalog

The MCP server registry is available at
https://github.com/docker/mcp-registry. To submit an MCP server, follow the
[contributing guidelines](https://github.com/docker/mcp-registry/blob/main/CONTRIBUTING.md).

When your pull request is reviewed and approved, your MCP server is available
within 24 hours on:

- Docker Desktop's [MCP Toolkit feature](toolkit.md).
- The [Docker MCP Catalog](https://hub.docker.com/mcp).
- The [Docker Hub](https://hub.docker.com/u/mcp) `mcp` namespace (for MCP
  servers built by Docker).
