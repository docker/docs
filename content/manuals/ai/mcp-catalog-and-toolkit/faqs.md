---
title: Security FAQs
linkTitle: Security FAQs
description: Frequently asked questions related to MCP Catalog and Toolkit security
keywords: MCP, Toolkit, MCP server, MCP client, security, faq
tags: [FAQ]
weight: 50
---

Docker MCP Catalog and Toolkit is a solution for securely building, sharing, and
running MCP tools. This page answers common questions about MCP Catalog and Toolkit security.

### What process does Docker follow to add a new MCP server to the catalog?

Developers can submit a pull request to the [Docker MCP Registry](https://github.com/docker/mcp-registry) to propose new servers. Docker provides detailed [contribution guidelines](https://github.com/docker/mcp-registry/blob/main/CONTRIBUTING.md) to help developers meet the required standards.

Currently, a majority of the servers in the catalog are built directly by Docker. Each server includes attestations such as:

- Build attestation: Servers are built on Docker Build Cloud.
- Source provenance: Verifiable source code origins.
- Signed SBOMs: Software Bill of Materials with cryptographic signatures.

> [!NOTE]
> When using the images with [Docker MCP gateway](/manuals/ai/mcp-gateway/_index.md),
> you can verify attestations at runtime using the `docker mcp gateway run
> --verify-signatures` CLI command.


In addition to Docker-built servers, the catalog includes select servers from trusted registries such as GitHub and HashiCorp. Each third-party server undergoes a verification process that includes:

- Pulling and building the code in an ephemeral build environment.
- Testing initialization and functionality.
- Verifying that tools can be successfully listed.

### Under what conditions does Docker reject MCP server submissions?

Docker rejects MCP server submissions that fail automated testing and validation processes during pull request review. Additionally, Docker reviewers evaluate submissions against specific requirements and reject MCP servers that don't meet these criteria.

### Does Docker take accountability for malicious MCP servers in the Toolkit?

Docker’s security measures currently represent a best-effort approach. While Docker implements automated testing, scanning, and metadata extraction for each server in the catalog, these security measures are not yet exhaustive. Docker is actively working to enhance its security processes and expand testing coverage. Enterprise customers can contact their Docker account manager for specific security requirements and implementation details.

### How are credentials managed for MCP servers?

Starting with Docker Desktop version 4.43.0, credentials are stored securely in the Docker Desktop VM. The storage implementation depends on the platform (for example, macOS, WSL2). You can manage the credentials using the following CLI commands:

- `docker mcp secret ls` - List stored credentials
- `docker mcp secret rm` - Remove specific credentials
- `docker mcp oauth revoke` - Revoke OAuth-based credentials

In the upcoming versions of Docker Desktop, Docker plans to support pluggable storage for these secrets and additional out-of-the-box storage providers to give users more flexibility in managing credentials.

### Are credentials removed when an MCP server is uninstalled?

No. MCP servers are not technically uninstalled since they exist as Docker containers pulled to your local Docker Desktop. Removing an MCP server stops the container but leaves the image on your system. Even if the container is deleted, credentials remain stored until you remove them manually.

## Related pages

- [Get started with MCP Toolkit](/manuals/ai/mcp-catalog-and-toolkit/get-started.md)
- [Open-source MCP Gateway](/manuals/ai/mcp-gateway/_index.md)
