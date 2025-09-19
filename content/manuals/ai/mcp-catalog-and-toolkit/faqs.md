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

When using the images with [Docker MCP gateway](/manuals/ai/mcp-gateway/_index.md), these attestations are verified at runtime.

In addition to Docker-built servers, the catalog includes select servers from trusted registries such as GitHub and HashiCorp. Each third-party server undergoes a strict verification process that includes:

- Pulling and building the code in an ephemeral build environment.
- Ensuring only trusted materials are used during the build.
- Testing initialization and functionality.
- Verifying that tools, resources, and prompts can be successfully listed.

### Does Docker take accountability for malicious MCP servers in the Toolkit?

Docker’s security measures currently represent a best-effort approach. While Docker implements automated testing, scanning, and metadata extraction for each server in the catalog, these security measures are not yet exhaustive. Docker is actively working to enhance its security processes and expand testing coverage.

### How are credentials  managed for MCP servers?

Starting with Docker Desktop version 4.43.0, credentials are stored securely in the Docker Desktop VM. The storage implementation depends on the platform (for example, macOS, WSL2, etc.). You can manage the credentials using `docker mcp secret ls` and `docker mcp secret rm`  CLI commands. If OAuth is used, you can run `docker mcp oauth revoke` command to remove the credentials.

In the upcoming versions of Docker Desktop, we aim to support pluggable  storage for these secrets, and a few out-of-the-box storage providers for enhanced flexibility for credential management.

### Are credentials removed when an MCP server is uninstalled?

MCP servers are not technically uninstalled since they exist as Docker containers pulled to your local Docker Desktop. when you disable an MCP server, the server stops running, but the container image remains on your system.


## Related pages

- [Get started with MCP Toolkit](/manuals/ai/mcp-catalog-and-toolkit/get-started.md)
- [Open-source MCP Gateway](/manuals/ai/mcp-gateway/_index.md)