---
description: Explore Docker Hub's generative AI catalogs and sandbox kits.
keywords: Docker Hub, Hub, generative AI, catalog
title: Docker Hub Generative AI
linkTitle: Generative AI
weight: 60
---

Docker Hub groups its generative AI content, including curated catalogs and
sandbox kits, under Generative AI. Catalogs are your go-to collections of
trusted, ready-to-use container images and resources, tailored to meet
specific development needs. They make it easier to find high-quality,
pre-verified content so you can quickly build, deploy, and manage your
applications with confidence. Catalogs in Docker Hub:

- Simplify content discovery: Organized and curated content makes it easy to
  discover tools and resources tailored to your specific domain or technology.
- Reduce complexity: Trusted resources, vetted by Docker and its partners,
  ensure security, reliability, and adherence to best practices.
- Accelerate development: Quickly integrate advanced capabilities into your
  applications without the hassle of extensive research or setup.

The following sections provide an overview of the generative AI content
available in Docker Hub.

## MCP Catalog

The [MCP Catalog](https://hub.docker.com/mcp/) is a centralized, trusted
registry for discovering, sharing, and running Model Context Protocol
(MCP)-compatible tools. Seamlessly integrated into Docker Hub, the catalog
includes:

- Over 100 verified MCP servers packaged as Docker images
- Tools from partners such as New Relic, Stripe, and Grafana
- Versioned releases with publisher verification
- Simplified pull-and-run support through Docker Desktop and Docker CLI

Each server runs in an isolated container to ensure consistent behavior and
minimize configuration headaches. For developers working with Claude Desktop or
other MCP clients, the catalog provides an easy way to extend functionality with
drop-in tools.

To learn more about MCP servers, see [MCP Catalog and Toolkit](../../ai/mcp-catalog-and-toolkit/_index.md).

## AI Models Catalog

The [AI Models Catalog](https://hub.docker.com/catalogs/models/) provides
curated, trusted models that work with [Docker Model
Runner](../../ai/model-runner/_index.md). This catalog is designed to make AI
development more accessible by offering pre-packaged, ready-to-use models that
you can pull, run, and interact with using familiar Docker tools.

With the AI Models Catalog and Docker Model Runner, you can:

- Pull and serve models from Docker Hub or any OCI-compliant registry
- Interact with models via OpenAI-compatible APIs
- Run and test models locally using Docker Desktop or CLI
- Package and publish models using the `docker model` CLI

Whether you're building generative AI applications, integrating LLMs into your
workflows, or experimenting with machine learning tools, the AI Models Catalog
simplifies the model management experience.

## Sandbox kits

[Sandbox kits](https://hub.docker.com/search?type=sbx_kit) package what an AI
coding agent needs to run in a [Docker Sandbox](../../ai/sandboxes/_index.md):
its image, tools, credentials, network rules, and startup commands.

Each kit's Hub page shows the command to run it, in the form
`sbx run <agent> --kit docker.io/<namespace>/<kit-name>`. To learn more about
building and running kits, see [Kits](../../ai/sandboxes/customize/kits.md).
