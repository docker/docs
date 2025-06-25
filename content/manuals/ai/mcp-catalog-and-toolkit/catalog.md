---
title: Docker MCP Catalog
description: Learn about the benefits of the MCP Catalog, how you can use it, and how you can contribute
keywords: docker hub, mcp, mcp servers, ai agents, calatog, docker
---

The [Docker MCP Catalog](https://hub.docker.com/catalogs/mcp) is a centralized, trusted registry for discovering, sharing, and running MCP-compatible tools. Seamlessly integrated into Docker Hub, it offers verified, versioned, and curated MCP servers packaged as Docker images. The catalog is also available in Docker Desktop.

The catalog solves common MCP server challenges:

- Environment conflicts: Tools often need specific runtimes that may clash with existing setups.
- Lack of isolation: Traditional setups risk exposing the host system.
- Setup complexity: Manual installation and configuration result in slow adoption.
- Inconsistency across platforms: Tools may behave unpredictably on different OSes.

With Docker, each MCP server runs as a self-contained container so it is
portable, isolated, and consistent. You can launch tools instantly using Docker
CLI or Docker Desktop, without worrying about dependencies or compatibility.

## Key features

- Over 100 verified MCP servers in one place
- Publisher verification and versioned releases
- Pull-based distribution using Docker's infrastructure
- Tools provided by partners such as New Relic, Stripe, Grafana, and more

## How it works

Each tool in the MCP Catalog is packaged as a Docker image with metadata:

- Discover tools via Docker Hub under the `mcp/` namespace.
- Connect tools to their preferred agents with simple configuration through the [MCP Toolkit](toolkit.md).
- Pull and run tools using Docker Desktop or the CLI.

Each catalog entry provides:

- Tool description and metadata
- Version history
- Example configuration for agent integration

## Use an MCP server from the catalog

To use an MCP server from the catalog, see [MCP toolkit](toolkit.md).

## Contribute an MCP server to the catalog

The MCP server registry is available at https://github.com/docker/mcp-registry. To submit an MCP server:
follow the [contributing guidelines](https://github.com/docker/mcp-registry/blob/main/CONTRIBUTING.md).

When your pull request is reviewed and approved, your MCP server is available in 24 hours on:

- Docker Desktop's [MCP Toolkit feature](toolkit.md)
- The [Docker MCP catalog](https://hub.docker.com/mcp)
- The [Docker Hub](https://hub.docker.com/u/mcp) mcp namespace (for MCP servers built by Docker)
