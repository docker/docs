---
title: Get started with Docker MCP Toolkit
linkTitle: Get started
description: Learn how to quickly install and use the MCP Toolkit to set up servers and clients.
keywords: Docker MCP Toolkit, MCP server, MCP client, AI agents
weight: 10
params:
  test_prompt: Use the GitHub MCP server to show me my open pull requests
---

{{< summary-bar feature_name="Docker MCP Toolkit" >}}

The Docker MCP Toolkit makes it easy to set up, manage, and run containerized
Model Context Protocol (MCP) servers, and connect them to AI agents. It
provides secure defaults and support for a growing ecosystem of LLM-based
clients. This page shows you how to get started quickly with the Docker MCP
Toolkit.

## Setup

Before you begin, make sure you meet the following requirements to get started with Docker MCP Toolkit.

1. Download and install the latest version of [Docker Desktop](/get-started/get-docker/).
2. Open the Docker Desktop settings and select **Beta features**.
3. Select **Enable Docker MCP Toolkit**.
4. Select **Apply**.

The **Learning center** in Docker Desktop provides walkthroughs and resources
to help you get started with Docker products and features. On the **MCP
Toolkit** page, the **Get started** walkthrough that guides you through
installing an MCP server, connecting a client, and testing your setup.

Alternatively, follow the step-by-step instructions on this page to:

- [Install MCP servers](#install-mcp-servers)
- [Connect clients](#connect-clients)
- [Verify connections](#verify-connections)

## Install MCP servers

{{< tabs >}}
{{< tab name="Docker Desktop">}}

1. In Docker Desktop, select **MCP Toolkit** and select the **Catalog** tab.
2. Search for the **GitHub Official** server from the catalog and then select the plus icon to add it.
3. In the **GitHub Official** server page, select the **Configuration** tab and select **OAuth**.

   > [!NOTE]
   >
   > The type of configuration required depends on the server you select. For the GitHub Official server, you must authenticate using OAuth.

   Your browser opens the GitHub authorization page. Follow the on-screen instructions to [authenticate via OAuth](/manuals/ai/mcp-catalog-and-toolkit/toolkit.md#authenticate-via-oauth).

4. Return to Docker Desktop when the authentication process is complete.
5. Search for the **Playwright** server from the catalog and add it.

{{< /tab >}}
{{< tab name="CLI">}}

1. Add the GitHub Official MCP server. Run:

   ```console
   $ docker mcp server enable github-official
   ```

2. Authenticate the server by running the following command:

   ```console
   $ docker mcp oauth authorize github
   ```

   > [!NOTE]
   >
   > The type of configuration required depends on the server you select. For the GitHub Official server, you must authenticate using OAuth.

   Your browser opens the GitHub authorization page. Follow the on-screen instructions to [authenticate via OAuth](/manuals/ai/mcp-catalog-and-toolkit/toolkit.md#authenticate-via-oauth).

3. Add the **Playwright** server. Run:

   ```console
   $ docker mcp server enable playwright
   ```

   {{< /tab >}}
   {{< /tabs >}}

You’ve now successfully added an MCP server. Next, connect an MCP client to use
the MCP Toolkit in an AI application.

## Connect clients

To connect a client to MCP Toolkit:

1. In Docker Desktop, select **MCP Toolkit** and select the **Clients** tab.
2. Find your application in the list.
3. Select **Connect** to configure the client.

If your client isn't listed, you can connect the MCP Toolkit manually over
`stdio` by configuring your client to run the following command:

```plaintext
docker mcp gateway run
```

For example, if your client uses a JSON file to configure MCP servers, you may
add an entry like:

```json {title="Example configuration"
{
  "servers": {
    "MCP_DOCKER": {
      "command": "docker",
      "args": ["mcp", "gateway", "run"],
      "type": "stdio"
    }
  }
}
```

Consult the documentation of the application you're using for instructions on
how to set up MCP servers manually.

## Verify connections

Refer to the relevant section for instructions on how to verify that your setup
is working:

- [Claude Code](#claude-code)
- [Claude Desktop](#claude-desktop)
- [OpenAI Codex](#codex)
- [Continue](#continue)
- [Cursor](#cursor)
- [Gemini](#gemini)
- [Goose](#goose)
- [Gordon](#gordon)
- [LM Studio](#lm-studio)
- [OpenCode](#opencode)
- [Sema4.ai](#sema4)
- [Visual Studio Code](#vscode)
- [Zed](#zed)

### Claude Code

If you configured the MCP Toolkit for a specific project, navigate to the
relevant project directory. Then run `claude mcp list`. The output should show
`MCP_DOCKER` with a "connected" status:

```console
$ claude mcp list
Checking MCP server health...

MCP_DOCKER: docker mcp gateway run - ✓ Connected
```

Test the connection by submitting a prompt that invokes one of your installed
MCP servers:

```console
$ claude "{{% param test_prompt %}}"
```

### Claude Desktop

Restart Claude Desktop and check the **Search and tools** menu in the chat
input. You should see the `MCP_DOCKER` server listed and enabled:

![Claude Desktop](images/claude-desktop.avif)

Test the connection by submitting a prompt that invokes one of your installed
MCP servers:

```plaintext
{{% param test_prompt %}}
```

### Codex

Run `codex mcp list` to view active MCP servers and their statuses. The
`MCP_DOCKER` server should appear in the list with an "enabled" status:

```console
$ codex mcp list
Name        Command  Args             Env  Cwd  Status   Auth
MCP_DOCKER  docker   mcp gateway run  -    -    enabled  Unsupported
```

Test the connection by submitting a prompt that invokes one of your installed
MCP servers:

```console
$ codex "{{% param test_prompt %}}"
```

### Continue

Launch the Continue terminal UI by running `cn`. Use the `/mcp` command to view
active MCP servers and their statuses. The `MCP_DOCKER` server should appear in
the list with a "connected" status:

```plaintext
   MCP Servers

   ➤ 🟢 MCP_DOCKER (🔧75 📝3)
     🔄 Restart all servers
     ⏹️ Stop all servers
     🔍 Explore MCP Servers
     Back

   ↑/↓ to navigate, Enter to select, Esc to go back
```

Test the connection by submitting a prompt that invokes one of your installed
MCP servers:

```console
$ cn "{{% param test_prompt %}}"
```

### Cursor

Open Cursor. If you configured the MCP Toolkit for a specific project, open the
relevant project directory. Then navigate to **Cursor Settings > Tools & MCP**.
You should see `MCP_DOCKER` under **Installed MCP Servers**:

![Cursor](images/cursor.avif)

Test the connection by submitting a prompt that invokes one of your installed
MCP servers:

```plaintext
{{% param test_prompt %}}
```

### Gemini

Run `gemini mcp list` to view active MCP servers and their statuses. The
`MCP_DOCKER` should appear in the list with a "connected" status.

```console
$ gemini mcp list
Configured MCP servers:

✓ MCP_DOCKER: docker mcp gateway run (stdio) - Connected
```

Test the connection by submitting a prompt that invokes one of your installed
MCP servers:

```console
$ gemini "{{% param test_prompt %}}"
```

### Goose

{{< tabs >}}
{{< tab name="Desktop app" >}}

Open the Goose desktop application and select **Extensions** in the sidebar.
Under **Enabled Extensions**, you should see an extension named `Mcpdocker`:

![Goose desktop app](images/goose.avif)

{{< /tab >}}
{{< tab name="CLI" >}}

Run `goose info -v` and look for an entry named `mcpdocker` under extensions.
The status should show `enabled: true`:

```console
$ goose info -v
…
    mcpdocker:
      args:
      - mcp
      - gateway
      - run
      available_tools: []
      bundled: null
      cmd: docker
      description: The Docker MCP Toolkit allows for easy configuration and consumption of MCP servers from the Docker MCP Catalog
      enabled: true
      env_keys: []
      envs: {}
      name: mcpdocker
      timeout: 300
      type: stdio
```

{{< /tab >}}
{{< /tabs >}}

Test the connection by submitting a prompt that invokes one of your installed
MCP servers:

```plaintext
{{% param "test_prompt" %}}
```

### Gordon

Open the **Ask Gordon** view in Docker Desktop and select the toolbox icon in
the chat input area. The **MCP Toolkit** tab shows whether MCP Toolkit is
enabled and displays all the provided tools:

![MCP Toolkit in the Ask Gordon UI](images/ask-gordon.avif)

Test the connection by submitting a prompt that invokes one of your installed
MCP servers, either directly in Docker Desktop or using the CLI:

```console
$ docker ai "{{% param "test_prompt" %}}"
```

### LM Studio

Restart LM Studio and start a new chat. Open the integrations menu and look for
an entry named `mcp/mcp-docker`. Use the toggle to enable the server:

![LM Studio](images/lm-studio.avif)

Test the connection by submitting a prompt that invokes one of your installed
MCP servers:

```plaintext
{{% param "test_prompt" %}}
```

### OpenCode

The OpenCode configuration file (at `~/.config/opencode/opencode.json` by
default) contains the setup for MCP Toolkit:

```json
{
  "mcp": {
    "MCP_DOCKER": {
      "type": "local",
      "command": ["docker", "mcp", "gateway", "run"],
      "enabled": true
    }
  },
  "$schema": "https://opencode.ai/config.json"
}
```

Test the connection by submitting a prompt that invokes one of your installed
MCP servers:

```console
$ opencode "{{% param "test_prompt" %}}"
```

### Sema4.ai Studio {#sema4}

In Sema4.ai Studio, select **Actions** in the sidebar, then select the **MCP
Servers** tab. You should see Docker MCP Toolkit in the list:

![Docker MCP Toolkit in Sema4.ai Studio](./images/sema4-mcp-list.avif)

To use MCP Toolkit with Sema4.ai, add it as an agent action. Find the agent you
want to connect to the MCP Toolkit and open the agent editor. Select **Add
Action**, enable Docker MCP Toolkit in the list, then save your agent:

![Editing an agent in Sema4.ai Studio](images/sema4-edit-agent.avif)

Test the connection by submitting a prompt that invokes one of your installed
MCP servers:

```plaintext
{{% param test_prompt %}}
```

### Visual Studio Code {#vscode}

Open Visual Studio Code. If you configured the MCP Toolkit for a specific
project, open the relevant project directory. Then open the **Extensions**
pane. You should see the `MCP_DOCKER` server listed under installed MCP
servers.

![MCP_DOCKER installed in Visual Studio Code](images/vscode-extensions.avif)

Test the connection by submitting a prompt that invokes one of your installed
MCP servers:

```plaintext
{{% param test_prompt %}}
```

### Zed

Launch Zed and open agent settings:

![Opening Zed agent settings from command palette](images/zed-cmd-palette.avif)

Ensure that `MCP_DOCKER` is listed and enabled in the MCP Servers section:

![MCP_DOCKER in Zed's agent settings](images/zed-agent-settings.avif)

Test the connection by submitting a prompt that invokes one of your installed
MCP servers:

```plaintext
{{% param test_prompt %}}
```

## Further reading

- [MCP Toolkit](/manuals/ai/mcp-catalog-and-toolkit/toolkit.md)
- [MCP Catalog](/manuals/ai/mcp-catalog-and-toolkit/catalog.md)
- [MCP Gateway](/manuals/ai/mcp-catalog-and-toolkit/mcp-gateway.md)
