---
title: MCP Gateway
description: "Docker's MCP Gateway provides secure, centralized, and scalable orchestration of AI tools through containerized MCP servers—empowering developers, operators, and security teams."
keywords: MCP Gateway
weight: 50
params:
  sidebar:
    group: Open source
---

The MCP Gateway is Docker's open-source enterprise-ready solution for
orchestrating and managing [Model Context Protocol
(MCP)](https://spec.modelcontextprotocol.io/) servers securely across
development and production environments. It is designed to help organizations
connect MCP servers from the [Docker MCP Catalog](https://hub.docker.com/mcp) to
MCP Clients without compromising security, visibility, or control.

By unifying multiple MCP servers into a single, secure endpoint, the MCP Gateway offers
the following benefits:

- Secure by default: MCP servers run in isolated Docker containers with restricted
  privileges, network access, and resource usage.
- Unified management: One gateway endpoint centralizes configuration, credentials,
  and access control for all MCP servers.
- Enterprise observability: Built-in monitoring, logging, and filtering tools ensure
  full visibility and governance of AI tool activity.

## Who is the MCP Gateway designed for?

The MCP Gateway solves problems encountered by various groups:

- Developers: Deploy MCP servers locally and in production using Docker Compose,
  with built-in support for protocol handling, credential management, and security policies.
- Security teams: Achieve enterprise-grade isolation and visibility into AI tool
  behavior and access patterns.
- Operators: Scale effortlessly from local development environments to production
  infrastructure with consistent, low-touch operations.

## Key features

- Server management: List, inspect, and call MCP tools, resources and prompts from multiple servers
- Container-based servers: Run MCP servers as Docker containers with proper isolation
- Secrets management: Secure handling of API keys and credentials via Docker Desktop
- Dynamic discovery and reloading: Automatic tool, prompt, and resource discovery from running servers
- Monitoring: Built-in logging and call tracing capabilities

## Install a pre-release version of the MCP Gateway

If you use Docker Desktop, the MCP Gateway is readily available. Use the
following instructions to test pre-release versions.

### Prerequisites

- Docker Desktop with the [MCP Toolkit feature enabled](../mcp-catalog-and-toolkit/toolkit.md#enable-docker-mcp-toolkit).
- Go 1.24+ (for development)

### Install using a pre-built binary

You can download the latest binary from the [GitHub releases page](https://github.com/docker/mcp-gateway/releases/latest).

Rename the relevant binary and copy it to the destination matching your OS:

| OS      | Binary name      | Destination folder                  |
|---------|------------------|-------------------------------------|
| Linux   | `docker-mcp`     | `$HOME/.docker/cli-plugins`         |
| macOS   | `docker-mcp`     | `$HOME/.docker/cli-plugins`         |
| Windows | `docker-mcp.exe` | `%USERPROFILE%\.docker\cli-plugins` |

Or copy it into one of these folders for installing it system-wide:


{{< tabs group="" >}}
{{< tab name="On Unix environments">}}

* `/usr/local/lib/docker/cli-plugins` OR `/usr/local/libexec/docker/cli-plugins`
* `/usr/lib/docker/cli-plugins` OR `/usr/libexec/docker/cli-plugins`

> [!NOTE]
> You may have to make the binaries executable with `chmod +x`:
> ```bash
> $ chmod +x ~/.docker/cli-plugins/docker-mcp
> ```

{{< /tab >}}
{{< tab name="On Windows">}}

* `C:\ProgramData\Docker\cli-plugins`
* `C:\Program Files\Docker\cli-plugins`

{{< /tab >}}
{{</tabs >}}

You can now use the `mcp` command:

```bash
docker mcp --help
```

## Use the MCP Gateway

1. Select a server of your choice from the [MCP Catalog](https://hub.docker.com/mcp)
   and copy the install command from the **Manual installation** section.

1. For example, run this command in your terminal to install the `duckduckgo`
   MCP server:

   ```console
   docker mcp server enable duckduckgo
   ```

1. Connect a client, like Visual Studio Code:

   ```console
   docker mcp client connect vscode
   ```

1. Run the gateway:

   ```console
   docker mcp gateway run
   ```

Now your MCP gateway is running and you can leverage all the servers set up
behind it from Visual Studio Code.

[View the complete docs on GitHub.](https://github.com/docker/mcp-gateway?tab=readme-ov-file#usage)

## Related pages

- [Docker MCP Toolkit and catalog](/manuals/ai/mcp-catalog-and-toolkit/_index.md)
