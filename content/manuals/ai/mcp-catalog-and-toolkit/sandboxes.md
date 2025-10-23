---
title: Sandboxes
description: "Learn how sandboxes provide secure, isolated execution environments for AI agents in the MCP ecosystem, enabling safe code execution and protecting production systems."
keywords: Sandboxes, E2B, MCP Gateway, isolated environment, AI agent security
params:
  sidebar:
    badge:
      color: green
      text: New
weight: 50
---

Sandboxes are isolated execution environments that provide secure, controlled spaces for running code and applications without affecting the host system. They create strict boundaries around executing processes, preventing access to unauthorized resources while providing consistent, reproducible environments. Think of it as a virtual "playground" with clearly defined boundaries, where code can execute freely within those boundaries but cannot escape to impact other systems or access sensitive data.

In the Model Context Protocol ecosystem, sandboxes address several critical challenges that arise when AI agents need to execute code or interact with external systems. They enable safe code execution for AI-generated scripts, secure tool validation for MCP servers, and multi-tenant isolation when multiple agents share infrastructure. This ensures that sensitive credentials and data remain protected within appropriate security boundaries while maintaining compliance and audit requirements.

## Key features

- Isolation and Security: Complete separation between executing code and the host environment, with strict controls over file access, network connections, and system calls.
- Resource Management: Fine-grained control over CPU, memory, disk space, and network usage to prevent resource exhaustion.
- Reproducible Environments: Consistent, predictable execution environments. Code that runs successfully in one sandbox instance will behave identically in another.
- Ephemeral Environments: Temporary, disposable environments that can be destroyed after task completion, leaving no persistent artifacts.

## E2B sandboxes

Docker has partnered with [E2B](https://e2b.dev/), a provider of secure cloud sandboxes for AI agents. Through this partnership, every E2B sandbox now includes direct access to Docker’s [MCP Catalog](https://hub.docker.com/mcp), a collection of 200+ tools, including ones from known publishers such as GitHub, Notion, and Stripe, all enabled through the Docker MCP Gateway.

When creating a new sandbox, E2B users can specify which MCP tools the sandbox should access. E2B then launches these MCP tools and provides access through the Docker MCP Gateway.

The following example shows how to set up an E2B sandbox with GitHub and Notion MCP servers.

## Example: Using GitHub and Notion MCP server

The following example demonstrates how to analyze data in Notion and create GitHub issues. By the end, you'll understand how to connect multiple MCP servers in an E2B sandbox and orchestrate cross-platform workflows.

### Prerequisites

Before you begin, make sure you have the following:

- [E2B account](https://e2b.dev/docs/quickstart) with API access
- Anthropic API key for Claude

    >[!Note]
    >
    > This example uses Claude CLI which comes pre-installed in E2B sandboxes. However,
    > you can adapt the example to work with other AI assistants of your choice. See
    > [E2B's MCP documentation](https://e2b.dev/docs/mcp/quickstart) for alternative
    > connection methods.

- Node.js 18+ installed on your machine
- Notion account with:
  - A database containing sample data
  - [Integration token](https://www.notion.com/help/add-and-manage-connections-with-the-api)
- GitHub account with:
    - A repository for testing
    - Personal access token with `repo` scope

### Set up your environment

Create a new directory and initialize a Node.js project:

```bash
mkdir mcp-e2b-quickstart
cd mcp-e2b-quickstart
npm init -y
```

Configure your project for ES modules by updating `package.json`:

```json
{
  "name": "mcp-e2b-quickstart",
  "version": "1.0.0",
  "type": "module",
  "scripts": {
    "start": "node index.js"
  }
}
```

Install required dependencies:

```bash
npm install e2b dotenv
```

Create a `.env` file with your credentials:

```bash
E2B_API_KEY=your_e2b_api_key_here
ANTHROPIC_API_KEY=your_anthropic_api_key_here
NOTION_INTEGRATION_TOKEN=ntn_your_notion_integration_token_here
GITHUB_TOKEN=ghp_your_github_pat_here
```

Protect your credentials:

```bash
echo ".env" >> .gitignore
echo "node_modules/" >> .gitignore
```

### Create an E2B sandbox with MCP servers

{{< tabs group="" >}}
{{< tab name="Typescript">}}

Create a file named `index.ts`:

```typescript
import 'dotenv/config';
import { Sandbox } from 'e2b';

async function quickstart(): Promise<void> {
  console.log("Creating E2B sandbox with Notion and GitHub MCP servers...\n");

  const sbx: Sandbox = await Sandbox.betaCreate({
    envs: {
      ANTHROPIC_API_KEY: process.env.ANTHROPIC_API_KEY as string,
    },
    mcp: {
      notion: {
        internalIntegrationToken: process.env.NOTION_INTEGRATION_TOKEN as string,
      },
      githubOfficial: {
        githubPersonalAccessToken: process.env.GITHUB_TOKEN as string,
      },
    },
  });

  const mcpUrl: string = sbx.betaGetMcpUrl();
  const mcpToken: string = await sbx.betaGetMcpToken();

  console.log("Sandbox created successfully!");
  console.log(`MCP Gateway URL: ${mcpUrl}\n`);

  // Wait for MCP initialization
  await new Promise<void>(resolve => setTimeout(resolve, 1000));

  // Connect Claude CLI to MCP gateway
  console.log("Connecting Claude CLI to MCP gateway...");
  await sbx.commands.run(
    `claude mcp add --transport http e2b-mcp-gateway ${mcpUrl} --header "Authorization: Bearer ${mcpToken}"`,
    {
      timeoutMs: 0,
      onStdout: console.log,
      onStderr: console.log
    }
  );

  console.log("\nConnection successful! Cleaning up...");
  await sbx.kill();
}

quickstart().catch(console.error);
```

Run the script:

```typescript
npx tsx index.ts
```

{{< /tab >}}
{{< tab name="Python">}}

Create a file named `index.py`:

```python
import os
import asyncio
from dotenv import load_dotenv
from e2b import Sandbox

load_dotenv()

async def quickstart():
    print("Creating E2B sandbox with Notion and GitHub MCP servers...\n")

    sbx = await Sandbox.beta_create(
        envs={
            "ANTHROPIC_API_KEY": os.getenv("ANTHROPIC_API_KEY"),
        },
        mcp={
            "notion": {
                "internalIntegrationToken": os.getenv("NOTION_INTEGRATION_TOKEN"),
            },
            "githubOfficial": {
                "githubPersonalAccessToken": os.getenv("GITHUB_TOKEN"),
            },
        },
    )

    mcp_url = sbx.beta_get_mcp_url()
    mcp_token = await sbx.beta_get_mcp_token()

    print("Sandbox created successfully!")
    print(f"MCP Gateway URL: {mcp_url}\n")

    # Wait for MCP initialization
    await asyncio.sleep(1)

    # Connect Claude CLI to MCP gateway
    print("Connecting Claude CLI to MCP gateway...")

    def on_stdout(output):
        print(output, end='')

    def on_stderr(output):
        print(output, end='')

    await sbx.commands.run(
        f'claude mcp add --transport http e2b-mcp-gateway {mcp_url} --header "Authorization: Bearer {mcp_token}"',
        timeout_ms=0,
        on_stdout=on_stdout,
        on_stderr=on_stderr
    )

    print("\nConnection successful! Cleaning up...")
    await sbx.kill()

if __name__ == "__main__":
    try:
        asyncio.run(quickstart())
    except Exception as e:
        print(f"Error: {e}")

```

Run the script:

```python
python index.py
```

{{< /tab >}}
{{</tabs >}}

You should see:

```bash
Creating E2B sandbox with Notion and GitHub MCP servers...

Sandbox created successfully!
MCP Gateway URL: https://50005-xxxxx.e2b.app/mcp

Connecting Claude CLI to MCP gateway...
Added HTTP MCP server e2b-mcp-gateway with URL: https://50005-xxxxx.e2b.app/mcp

Connection successful! Cleaning up...
```

### Test with example workflow

Now, test the setup by running a simple workflow that searches Notion and creates a GitHub issue.

{{< tabs group="" >}}
{{< tab name="Typescript">}}

>[!IMPORTANT]
>
> Replace `owner/repo` in the prompt with your actual GitHub username and repository
> name (for example, `yourname/test-repo`).

Update `index.ts` with the following example:

```typescript
import 'dotenv/config';
import { Sandbox } from 'e2b';

async function exampleWorkflow(): Promise<void> {
  console.log("Creating sandbox...\n");

  const sbx: Sandbox = await Sandbox.betaCreate({
    envs: {
      ANTHROPIC_API_KEY: process.env.ANTHROPIC_API_KEY as string,
    },
    mcp: {
      notion: {
        internalIntegrationToken: process.env.NOTION_INTEGRATION_TOKEN as string,
      },
      githubOfficial: {
        githubPersonalAccessToken: process.env.GITHUB_TOKEN as string,
      },
    },
  });

  const mcpUrl: string = sbx.betaGetMcpUrl();
  const mcpToken: string = await sbx.betaGetMcpToken();

  console.log("Sandbox created successfully\n");

  // Wait for MCP servers to initialize
  await new Promise<void>(resolve => setTimeout(resolve, 3000));

  console.log("Connecting Claude to MCP gateway...\n");
  await sbx.commands.run(
    `claude mcp add --transport http e2b-mcp-gateway ${mcpUrl} --header "Authorization: Bearer ${mcpToken}"`,
    {
      timeoutMs: 0,
      onStdout: console.log,
      onStderr: console.log
    }
  );

  console.log("\nRunning example: Search Notion and create GitHub issue...\n");

  const prompt: string = `Using Notion and GitHub MCP tools:
1. Search my Notion workspace for databases
2. Create a test issue in owner/repo titled "MCP Toolkit Test" with description "Testing E2B + Docker MCP integration"
3. Confirm both operations completed successfully`;

  await sbx.commands.run(
    `echo '${prompt.replace(/'/g, "'\\''")}' | claude -p --dangerously-skip-permissions`,
    {
      timeoutMs: 0,
      onStdout: console.log,
      onStderr: console.log
    }
  );

  await sbx.kill();
}

exampleWorkflow().catch(console.error);
```

Run the script:

```typescript
npx tsx index.ts
```

{{< /tab >}}
{{< tab name="Python">}}

Update `index.py` with this example:

>[!IMPORTANT]
>
> Replace `owner/repo` in the prompt with your actual GitHub username and repository
> name (for example, `yourname/test-repo`).

```python
import os
import asyncio
import shlex
from dotenv import load_dotenv
from e2b import Sandbox

load_dotenv()

async def example_workflow():
    print("Creating sandbox...\n")

    sbx = await Sandbox.beta_create(
        envs={
            "ANTHROPIC_API_KEY": os.getenv("ANTHROPIC_API_KEY"),
        },
        mcp={
            "notion": {
                "internalIntegrationToken": os.getenv("NOTION_INTEGRATION_TOKEN"),
            },
            "githubOfficial": {
                "githubPersonalAccessToken": os.getenv("GITHUB_TOKEN"),
            },
        },
    )

    mcp_url = sbx.beta_get_mcp_url()
    mcp_token = await sbx.beta_get_mcp_token()

    print("Sandbox created successfully\n")

    # Wait for MCP servers to initialize
    await asyncio.sleep(3)

    print("Connecting Claude to MCP gateway...\n")

    def on_stdout(output):
        print(output, end='')

    def on_stderr(output):
        print(output, end='')

    await sbx.commands.run(
        f'claude mcp add --transport http e2b-mcp-gateway {mcp_url} --header "Authorization: Bearer {mcp_token}"',
        timeout_ms=0,
        on_stdout=on_stdout,
        on_stderr=on_stderr
    )

    print("\nRunning example: Search Notion and create GitHub issue...\n")

    prompt = """Using Notion and GitHub MCP tools:
1. Search my Notion workspace for databases
2. Create a test issue in owner/repo titled "MCP Toolkit Test" with description "Testing E2B + Docker MCP integration"
3. Confirm both operations completed successfully"""

    # Escape single quotes for shell
    escaped_prompt = prompt.replace("'", "'\\''")

    await sbx.commands.run(
        f"echo '{escaped_prompt}' | claude -p --dangerously-skip-permissions",
        timeout_ms=0,
        on_stdout=on_stdout,
        on_stderr=on_stderr
    )

    await sbx.kill()

if __name__ == "__main__":
    try:
        asyncio.run(example_workflow())
    except Exception as e:
        print(f"Error: {e}")
```

Run the script:

```bash
python workflow.py
```

{{< /tab >}}
{{</tabs >}}

You should see:

```bash
Creating sandbox...

Running example: Search Notion and create GitHub issue...

## Task Completed Successfully

I've completed both operations using the Notion and GitHub MCP tools:

### 1. Notion Workspace Search

Found 3 databases in your Notion workspace:
- **Customer Feedback** - Database with 12 entries tracking feature requests
- **Product Roadmap** - Planning database with 8 active projects
- **Meeting Notes** - Shared workspace with 45 pages

### 2. GitHub Issue Creation

Successfully created test issue:
- **Repository**: your-org/your-repo
- **Issue Number**: #47
- **Title**: "MCP Test"
- **Description**: "Testing E2B + Docker MCP integration"
- **Status**: Open
- **URL**: https://github.com/your-org/your-repo/issues/47

Both operations completed successfully. The MCP servers are properly configured and working.
```

You've successfully created an E2B sandbox with multiple MCP servers and used Claude to orchestrate a workflow across Notion and GitHub.

You can extend this example to combine any of the 200+ MCP servers in the Docker MCP Catalog to build sophisticated automation workflows for your specific needs.

## Related pages

- [How to build an AI-powered code quality workflow with SonarQube and E2B](/guides/github-sonarqube-sandbox.md)
- [Docker + E2B: Building the Future of Trusted AI](https://www.docker.com/blog/docker-e2b-building-the-future-of-trusted-ai/)
- [Docker MCP Toolkit and Catalog](/manuals/ai/mcp-catalog-and-toolkit/_index.md)
- [Docker MCP Gateway](/manuals/ai/mcp-catalog-and-toolkit/mcp-gateway.md)
- [E2B MCP documentation](https://e2b.dev/docs/mcp)
