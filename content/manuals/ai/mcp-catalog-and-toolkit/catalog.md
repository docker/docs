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

- Over 100 verified MCP servers in one place.
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
