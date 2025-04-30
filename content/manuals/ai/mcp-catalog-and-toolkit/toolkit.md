---
title: MCP Toolkit
description:
keywords: 
---

The Docker MCP Toolkit is a Docker Desktop extension local that enables seamless setup, management, and execution of containerized MCP servers and their connections to AI agents. It removes the friction from tool usage by offering secure defaults, one-click setup, and support for a growing ecosystem of LLM-based clients. It is the fastest path from MCP tool discovery to local execution.

## Key features

- Cross-LLM compatibility: Works out of the box with Claude Desktop, Cursor, Continue.dev, and [Gordon](/manuals/ai/gordon/_index.md).
- Integrated tool discovery: Browse and launch MCP servers that are available in the Docker MCP Catalog, directly from Docker Desktop.
- No manual setup: Skip dependency management, runtime setup, and manual server configuration.

## How it works

The **MCP Servers** tab lists all available servers from the Docker MCP Catalog. Each entry includes:

- Tool name and description
- Partner/publisher
- Number of callable tools and what they are

To enable an MCP server, simply use the toggle switch to toggle it on.

> [!NOTE]
>
> Some MCP servers requires secrets or tokens to be configured before it can be enabled. Instructions on how to do this can be found on each MCP servers' repository.

The **MCP Clients** tab lets you connect your enabled MCP servers to supported agents. Connection is as simple as selecting **Connect**, so you can switch between LLM providers without altering your MCP server integrations or security configurations.

## Installation

To install the Docker MCP Toolkit extension:

1. In the Docker Desktop Dashboard, select the **Extensions** view, and then select **Manage**.
2. Select the **Browse** tab and search for **Docker MCP Toolkit**.
3. On the **Docker MCP Toolkit** result, select install.

The extension then appears under the **My extensions** tab.

### Example

The following example assumes you have already installed and set up Claude Desktop.

1. In the Docker MCP Toolkit extension, search for the Puppeteer MCP server in the **MCP Servers** tab, and toggle it on to enable.
2. From the **MCP Clients** tab, select the **Connect** button for Claude Desktop. 
3. Within Claude Desktop, submit the following prompt using the Sonnet 3.5 model:

   ```text
   Take a screenshot of docs.docker.com and then invert the colors
   ```

Once you've given your consent to use the new tools, Claude spins up the Puppeteer MCP server inside a container, navigates to the target URL, captures and modify the page, and returns the screenshot.