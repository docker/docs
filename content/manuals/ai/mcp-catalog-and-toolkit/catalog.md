---
title: Docker MCP Catalog
linkTitle: Catalog
description: Learn about the benefits of the MCP Catalog, how you can use it, and how you can contribute
keywords: docker hub, mcp, mcp servers, ai agents, catalog, docker
weight: 20
---

{{< summary-bar feature_name="Docker MCP Catalog" >}}

The [Docker MCP Catalog](https://hub.docker.com/mcp) is a curated collection of
verified MCP servers, packaged as Docker images and distributed through Docker
Hub. It solves common challenges with running MCP servers locally: environment
conflicts, setup complexity, and security concerns.

The catalog serves as the source of available MCP servers. Each server runs as
an isolated container, making it portable and consistent across different
environments.

> [!NOTE]
> E2B sandboxes now include direct access to the Docker MCP Catalog, giving developers
> access to over 200 tools and services to seamlessly build and run AI agents. For
> more information, see [E2B Sandboxes](sandboxes.md).

## What's in the catalog

The Docker MCP Catalog includes:

- Verified servers: All local servers are versioned with full provenance and SBOM
  metadata
- Partner tools: Servers from New Relic, Stripe, Grafana, and other trusted
  partners
- Docker-built servers: Locally-running servers built and digitally signed by
  Docker for enhanced security
- Remote services: Cloud-hosted servers that connect to external services like
  GitHub, Notion, and Linear

You can browse the catalog at [hub.docker.com/mcp](https://hub.docker.com/mcp)
or through the **Catalog** tab in Docker Desktop's MCP Toolkit.

### Local versus remote servers

The catalog contains two types of servers based on where they run:

Local servers run as containers on your machine. They work offline once
downloaded and offer predictable performance and complete data privacy. Docker
builds and signs all local servers in the catalog.

Remote servers run on the provider's infrastructure and connect to external
services. Many remote servers use OAuth authentication, which the MCP Toolkit
handles automatically through your browser.

## Using servers from the catalog

To start using MCP servers from the catalog:

1. Browse servers in the [MCP Catalog](https://hub.docker.com/mcp) or in Docker
   Desktop
2. Enable servers through the MCP Toolkit
3. Configure any required authentication (OAuth is handled automatically)
4. Connect your AI applications to use the servers

For detailed step-by-step instructions, see:

- [Get started with MCP Toolkit](/ai/mcp-catalog-and-toolkit/get-started/) -
  Quick start guide
- [MCP Toolkit](/ai/mcp-catalog-and-toolkit/toolkit/) - Detailed usage
  instructions

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

## Custom catalogs

Custom catalogs let you curate focused collections of recommended servers. You
can package custom server implementations alongside public servers, distribute
curated lists to your team, and define what agents can discover when using
Dynamic MCP.

Common use cases:

- Curate a subset of servers from the Docker MCP Catalog that your organization
  approves
- Include community registry servers that aren't in the Docker catalog
- Add your organization's private MCP servers
- Control which versions of servers your team uses

### Custom catalogs with Dynamic MCP

Custom catalogs work particularly well with
[Dynamic MCP](/ai/mcp-catalog-and-toolkit/dynamic-mcp/), where agents
discover and add MCP servers on-demand during conversations. When you specify a
custom catalog with the gateway, the `mcp-find` tool searches only within your
curated catalog. If your catalog contains 20 servers instead of 300+, agents
work within that focused set and can dynamically add servers as needed without
manual configuration each time.

This gives agents the autonomy to discover and use tools while keeping their
options within boundaries your team defines.

### Create and curate a catalog

The most practical way to create a custom catalog is to fork the Docker catalog
and then curate which servers to keep:

```console
$ docker mcp catalog fork docker-mcp my-catalog
```

This creates a copy of the Docker catalog with all available servers. Export it
to a file where you can edit which servers to include:

```console
$ docker mcp catalog export my-catalog ./my-catalog.yaml
```

Edit `my-catalog.yaml` to remove servers you don't want, keeping only the ones
your team needs. Each server is listed in the `registry` section. Import the
edited catalog back:

```console
$ docker mcp catalog import ./my-catalog.yaml
```

View your curated catalog:

```console
$ docker mcp catalog show my-catalog
```

#### Alternative: Build incrementally

You can also build a catalog from scratch. Start with an empty catalog or a
template:

```console
$ docker mcp catalog create my-catalog
```

Or create a starter template with example servers:

```console
$ docker mcp catalog bootstrap ./starter-catalog.yaml
```

Add servers from other catalog files:

```console
$ docker mcp catalog add my-catalog notion ./other-catalog.yaml
```

### Use a custom catalog

Use your custom catalog when running the MCP gateway. For static server
configuration, specify which servers to enable:

```console
$ docker mcp gateway run --catalog my-catalog.yaml --servers notion,brave
```

For Dynamic MCP, where agents discover and add servers during conversations,
specify just the catalog:

```console
$ docker mcp gateway run --catalog my-catalog.yaml
```

Agents can then use `mcp-find` to search for servers within your catalog and
`mcp-add` to enable them dynamically.

The `--catalog` flag points to a catalog file in `~/.docker/mcp/catalogs/`.

### Share your catalog

Share your catalog with your team by distributing the YAML file or hosting it
at a URL:

```console
$ docker mcp catalog export my-catalog ./team-catalog.yaml
```

Team members can import it:

```console
$ docker mcp catalog import ./team-catalog.yaml
$ docker mcp catalog import https://example.com/team-catalog.yaml
```
