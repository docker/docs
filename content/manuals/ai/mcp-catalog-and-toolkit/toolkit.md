---
title: MCP Toolkit
description: Use the MCP Toolkit to set up MCP servers and MCP clients.
keywords: Docker MCP Toolkit, MCP server, MCP client, AI agents
aliases:
  - /desktop/features/gordon/mcp/gordon-mcp-server/
  - /ai/gordon/mcp/gordon-mcp-server/
---

The Docker MCP Toolkit is a gateway that enables seamless setup, management, and execution of containerized MCP servers and their connections to AI agents. It removes the friction from tool usage by offering secure defaults, one-click setup, and support for a growing ecosystem of LLM-based clients. It is the fastest path from MCP tool discovery to local execution.

## Key features

- Cross-LLM compatibility: Instantly works with Claude Desktop, Cursor, Continue.dev, and [Gordon](/manuals/ai/gordon/_index.md).
- Integrated tool discovery: Browse and launch MCP servers from the Docker MCP Catalog directly in Docker Desktop.
- Zero manual setup: No dependency management, runtime configuration, or server setup required.
- Functions as both an MCP server aggregator and a gateway for clients to access installed MCP servers.

## How the MCP toolkit works

MCP introduces two core concepts: MCP clients and MCP servers. 

- MCP clients are typically embedded in LLM-based applications, such as
  the Claude Desktop App. They request resources or actions.
- MCP servers are launched by the client to perform the requested tasks,
  using any necessary tools, languages, or processes.

Docker standardizes the development, packaging, and distribution of
applications, including MCP servers. By packaging MCP servers as containers,
Docker eliminates issues related to isolation and environment differences. Users
can run a container directly, without managing dependencies or configuring
runtimes.

Depending on the MCP server, the tools it provides may run within the same container
as the server or in dedicated containers:


{{< tabs group="" >}}
{{< tab name="Single container">}}

![Visualisation of the MCP toolkit](/assets/images/mcp_servers.png)

{{< /tab >}}
{{< tab name="Separate containers">}}

![Visualisation of the MCP toolkit](/assets/images/mcp_servers_2.png)

{{< /tab >}}
{{</tabs >}}

## Security

The Docker MCP Toolkit combines passive and active measures to reduce attack
surfaces and ensure safe runtime behavior.

### Passive security

- Image signing and attestation: All MCP server images under `mcp/` in the [catalog](catalog.md) 
  are built by Docker and digitally
  signed to verify their source and integrity. Each image includes a Software
  Bill of Materials (SBOM) for full transparency.

### Active security

Security at runtime is enforced through resource and access limitations:

- CPU allocation: MCP tools are run in their own container. They are
  restricted to 1 CPU, limiting the impact of potential misuse of computing
  resources.

- Memory allocation: Containers for MCP tools are limited to 2 GB.

- Filesystem access: By default, MCP Servers have no access to the host filesystem.
  The user explicitly selects the servers that will be granted file mounts.

- Interception of tool requests: Requests to and from tools that contain sensitive
  information such as secrets are blocked.

## Enable Docker MCP Toolkit

1. Open the Docker Desktop settings and select **Beta features**.
2. Select **Enable Docker MCP Toolkit**.
3. Select **Apply & restart**.

>[!NOTE]
>If you have the MCP Toolkit _extension_ installed, you can uninstall it.

## Install an MCP server

To install an MCP server:

1. In Docker Desktop, select **MCP Toolkit** and select the **Catalog** tab.
   When you select a server you can see the following
   information:

   - Tool name and description 
   - Partner/publisher
   - The list of callable tools the server provides.

2. Find the MCP server of your choice and select the **Plus** icon.
3. Optional: Some servers require extra configuration. To configure them, select
   the **Config** tab and follow the instructions available on the repository of the provider of the MCP server.

> [!TIP]
> By default, the Gordon [client](#install-an-mcp-client) is enabled,
> which means Gordon can automatically interact with your MCP servers.

To learn more about the MCP server catalog, see [Catalog](catalog.md).

### Example: Use the GitHub MCP server

Imagine you want to enable Ask Gordon to interact with your GitHub account:

1. From the **MCP Toolkit** menu, select the **Catalog** tab and find
   the **GitHub Official** server and add it.
2. In the server's **Config** tab, insert your token generated from
   your [GitHub account](https://github.com/settings/personal-access-tokens).
3. In the **Clients** tab, ensure Gordon is connected.
4. From the **Ask Gordon** menu, you can now send requests related to your
   GitHub account, in accordance to the tools provided by the GitHub MCP server. To test it, ask Gordon:

   ```text
   What's my GitHub handle?
   ```

   Make sure to allow Gordon to interact with GitHub by selecting **Always allow** in Gordon's answer.

## Install an MCP client

When you have installed MCP servers, you can add clients to the MCP Toolkit. These clients
can interact with the installed MCP servers, turning the MCP Toolkit into a gateway.

To install a client:

1. In Docker Desktop, select **MCP Toolkit** and select the **Clients** tab.
2. Find the client of your choice and select **Connect**.

Your client can now interact with the MCP Toolkit.

### Example: Use Claude Desktop as a client

Imagine you have Claude Desktop installed, and you want to use the GitHub MCP server, 
and the Puppeteer MCP server, you do not have to install the servers in Claude Desktop.
You can simply install these 2 MCP servers in the MCP Toolkit,
and add Claude Desktop as a client:

1. From the **MCP Toolkit** menu, select the **Catalog** tab and find the **Puppeteer** server and add it.
2. Repeat for the **GitHub** server.
3. From the **Clients** tab, select **Connect** next to **Claude Desktop**. Restart 
   Claude Desktop if it's running, and it can now access all the servers in the MCP Toolkit.
4. Within Claude Desktop, run a test by submitting the following prompt using the Sonnet 3.5 model:

   ```text
   Take a screenshot of docs.docker.com and then invert the colors
   ```
5. From the **Clients** tab, select **Connect** next to **Claude Desktop**. Restart 
   Claude Desktop if it's running, and it can now access all the servers in the MCP Toolkit.
6. Within Claude Desktop, run a test by submitting the following prompt using the Sonnet 3.5 model:

   ```text
   Take a screenshot of docs.docker.com and then invert the colors
   ```

## Related pages

- [Open-source MCP Gateway](/manuals/ai/mcp-gateway/_index.md)
