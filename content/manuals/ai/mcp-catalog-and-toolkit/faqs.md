---
title: MCP Toolkit FAQs
linkTitle: FAQs
description: Frequently asked questions related to MCP Catalog and Toolkit security
keywords: MCP, Toolkit, MCP server, MCP client, security, faq
tags: [FAQ]
weight: 70
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
> When using the images with [Docker MCP gateway](/manuals/ai/mcp-catalog-and-toolkit/mcp-gateway.md),
> you can verify attestations at runtime using the `docker mcp gateway run
--verify-signatures` CLI command.


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

### Why don't I see remote MCP servers in the catalog?

If remote MCP servers aren't visible in the Docker Desktop catalog, your local
catalog may be out of date. Remote servers are indicated by a cloud icon and
include services like GitHub, Notion, and Linear.

Update your catalog by running:

```console
$ docker mcp catalog update
```

After the update completes, refresh the **Catalog** tab in Docker Desktop.

### What's the difference between profiles and the catalog?

The [catalog](/manuals/ai/mcp-catalog-and-toolkit/catalog.md) is the source of
available MCP servers - a library of tools you can choose from.
[Profiles](/manuals/ai/mcp-catalog-and-toolkit/profiles.md) are collections of
servers you've added to organize your work. Think of the catalog as a library,
and profiles as your personal bookshelves containing the books you've selected
for different purposes.

### Can I share profiles with my team?

Yes. Profiles can be pushed to OCI-compliant registries using
`docker mcp profile push my-profile registry.example.com/profiles/my-profile:v1`.
Team members can pull your profile with
`docker mcp profile pull registry.example.com/profiles/my-profile:v1`. Note
that credentials aren't included in shared profiles for security reasons - team
members need to configure OAuth and other credentials separately.

### Do I need to create a profile to use MCP Toolkit?

Yes, MCP Toolkit requires a profile to run servers. If you're upgrading from a
version before profiles were introduced, a default profile is automatically
created for you with your existing server configurations. You can create
additional named profiles to organize servers for different projects or
environments.

### What happens to servers when I switch profiles?

Each profile contains its own set of servers and configurations. When you run
the gateway with `--profile profile-name`, only servers in that profile are
available to clients. The default profile is used when no profile is specified.
Switching between profiles changes which servers your AI applications can
access.

### Can I use the same server in multiple profiles?

Yes. You can add the same MCP server to multiple profiles, each with different
configurations if needed. This is useful when you need the same server with
different settings for different projects or environments.

## Related pages

- [Get started with MCP Toolkit](/manuals/ai/mcp-catalog-and-toolkit/get-started.md)
- [Open-source MCP Gateway](/manuals/ai/mcp-catalog-and-toolkit/mcp-gateway.md)
