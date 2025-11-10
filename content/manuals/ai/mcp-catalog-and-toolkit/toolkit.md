---
title: Docker MCP Toolkit
linkTitle: MCP Toolkit
description: Use the MCP Toolkit to set up MCP servers and MCP clients.
keywords: Docker MCP Toolkit, MCP server, MCP client, AI agents
weight: 30
aliases:
  - /desktop/features/gordon/mcp/gordon-mcp-server/
  - /ai/gordon/mcp/gordon-mcp-server/
---

{{< summary-bar feature_name="Docker MCP Toolkit" >}}

The Docker MCP Toolkit is a management interface integrated into Docker Desktop
that lets you set up, manage, and run containerized MCP servers and connect
them to AI agents. It removes friction from tool usage by offering secure
defaults, easy setup, and support for a growing ecosystem of LLM-based clients.
It is the fastest way from MCP tool discovery to local execution.

## Key features

- Cross-LLM compatibility: Works with Claude, Cursor, and other MCP clients.
- Integrated tool discovery: Browse and launch MCP servers from the Docker MCP Catalog directly in Docker Desktop.
- Zero manual setup: No dependency management, runtime configuration, or setup required.
- Functions as both an MCP server aggregator and a gateway for clients to access installed MCP servers.

> [!TIP]
> The MCP Toolkit includes [Dynamic MCP](/manuals/ai/mcp-catalog-and-toolkit/dynamic-mcp.md),
> which enables AI agents to discover, add, and compose MCP servers on-demand during
> conversations, without manual configuration. Your agent can search the catalog and
> add tools as needed when you connect to the gateway.

## How the MCP Toolkit works

MCP introduces two core concepts: MCP clients and MCP servers.

- MCP clients are typically embedded in LLM-based applications, such as the
  Claude Desktop app. They request resources or actions.
- MCP servers are launched by the client to perform the requested tasks, using
  any necessary tools, languages, or processes.

Docker standardizes the development, packaging, and distribution of
applications, including MCP servers. By packaging MCP servers as containers,
Docker eliminates issues related to isolation and environment differences. You
can run a container directly, without managing dependencies or configuring
runtimes.

Depending on the MCP server, the tools it provides might run within the same
container as the server or in dedicated containers for better isolation.

## Security

The Docker MCP Toolkit combines passive and active measures to reduce attack
surfaces and ensure safe runtime behavior.

### Passive security

Passive security refers to measures implemented at build-time, when the MCP
server code is packaged into a Docker image.

- Image signing and attestation: All MCP server images under `mcp/` in the [MCP
  Catalog](catalog.md) are built by Docker and digitally signed to verify their
  source and integrity. Each image includes a Software Bill of Materials (SBOM)
  for full transparency.

### Active security

Active security refers to security measures at runtime, before and after tools
are invoked, enforced through resource and access limitations.

- CPU allocation: MCP tools are run in their own container. They are
  restricted to 1 CPU, limiting the impact of potential misuse of computing
  resources.

- Memory allocation: Containers for MCP tools are limited to 2 GB.

- Filesystem access: By default, MCP Servers have no access to the host filesystem.
  The user explicitly selects the servers that will be granted file mounts.

- Interception of tool requests: Requests to and from tools that contain sensitive
  information such as secrets are blocked.

### OAuth authentication

Some MCP servers require authentication to access external services like
GitHub, Notion, and Linear. The MCP Toolkit handles OAuth authentication
automatically. You authorize access through your browser, and the Toolkit
manages credentials securely. You don't need to manually create API tokens or
configure authentication for each service.

#### Authorize a server with OAuth

{{< tabs >}}
{{< tab name="Docker Desktop">}}

1. In Docker Desktop, go to **MCP Toolkit** and select the **Catalog** tab.
2. Find and add an MCP server that requires OAuth.
3. In the server's **Configuration** tab, select the **OAuth** authentication
   method. Follow the link to begin the OAuth authorization.
4. Your browser opens the authorization page for the service. Follow the
   on-screen instructions to complete authentication.
5. Return to Docker Desktop when authentication is complete.

View all authorized services in the **OAuth** tab. To revoke access, select
**Revoke** next to the service you want to disconnect.

{{< /tab >}}
{{< tab name="CLI">}}

Enable an MCP server:

```console
$ docker mcp server enable github-official
```

If the server requires OAuth, authorize the connection:

```console
$ docker mcp oauth authorize github
```

Your browser opens the authorization page. Complete the authentication process,
then return to your terminal.

View authorized services:

```console
$ docker mcp oauth ls
```

Revoke access to a service:

```console
$ docker mcp oauth revoke github
```

{{< /tab >}}
{{< /tabs >}}

## Usage examples

### Example: Use the GitHub Official MCP server with Ask Gordon

To illustrate how the MCP Toolkit works, here's how to enable the GitHub
Official MCP server and use [Ask Gordon](/manuals/ai/gordon/_index.md) to
interact with your GitHub account:

1. From the **MCP Toolkit** menu in Docker Desktop, select the **Catalog** tab
   and find the **GitHub Official** server and add it.
2. In the server's **Configuration** tab, authenticate via OAuth.
3. In the **Clients** tab, ensure Gordon is connected.
4. From the **Ask Gordon** menu, you can now send requests related to your
   GitHub account, in accordance to the tools provided by the GitHub Official
   server. To test it, ask Gordon:

   ```text
   What's my GitHub handle?
   ```

   Make sure to allow Gordon to interact with GitHub by selecting **Always
   allow** in Gordon's answer.

> [!TIP]
> The Gordon client is enabled by default,
> which means Gordon can automatically interact with your MCP servers.

### Example: Use Claude Desktop as a client

Imagine you have Claude Desktop installed, and you want to use the GitHub MCP server,
and the Puppeteer MCP server, you do not have to install the servers in Claude Desktop.
You can simply install these 2 MCP servers in the MCP Toolkit,
and add Claude Desktop as a client:

1. From the **MCP Toolkit** menu, select the **Catalog** tab and find the **Puppeteer** server and add it.
1. Repeat for the **GitHub Official** server.
1. From the **Clients** tab, select **Connect** next to **Claude Desktop**. Restart
   Claude Desktop if it's running, and it can now access all the servers in the MCP Toolkit.
1. Within Claude Desktop, run a test by submitting the following prompt using the Sonnet 3.5 model:

   ```text
   Take a screenshot of docs.docker.com and then invert the colors
   ```

### Example: Use Visual Studio Code as a client

You can interact with all your installed MCP servers in Visual Studio Code:

1. To enable the MCP Toolkit:

   {{< tabs group="" >}}
   {{< tab name="Enable globally">}}

   1. Insert the following in your Visual Studio Code's User `mcp.json`:

      ```json
      "mcp": {
       "servers": {
         "MCP_DOCKER": {
           "command": "docker",
           "args": [
             "mcp",
             "gateway",
             "run"
           ],
           "type": "stdio"
         }
       }
      }
      ```

   {{< /tab >}}
   {{< tab name="Enable for a given project">}}

   1. In your terminal, navigate to your project's folder.
   1. Run:

      ```bash
      docker mcp client connect vscode
      ```

      > [!NOTE]
      > This command creates a `.vscode/mcp.json` file in the current
      > directory. As this is a user-specific file, add it to your `.gitignore`
      > file to prevent it from being committed to the repository.
      >
      > ```console
      > echo ".vscode/mcp.json" >> .gitignore
      > ```

  {{< /tab >}}
  {{</tabs >}}

1. In Visual Studio Code, open a new Chat and select the **Agent** mode:

   ![Copilot mode switching](./images/copilot-mode.png)

1. You can also check the available MCP tools:

   ![Displaying tools in VSCode](./images/tools.png)

For more information about the Agent mode, see the
[Visual Studio Code documentation](https://code.visualstudio.com/docs/copilot/chat/mcp-servers#_use-mcp-tools-in-agent-mode).

## Further reading

- [MCP Catalog](/manuals/ai/mcp-catalog-and-toolkit/catalog.md)
- [MCP Gateway](/manuals/ai/mcp-catalog-and-toolkit/mcp-gateway.md)
