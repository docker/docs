---
title: Build a code quality check workflow
linkTitle: Build workflow
summary: Learn to use GitHub and SonarQube MCP servers in E2B sandboxes through progressive examples.
description: Create E2B sandboxes, discover MCP tools, test individual operations, and build complete quality-gated PR workflows.
weight: 10
---

In this section, you'll build a complete code quality automation workflow
step-by-step. You'll start by creating an E2B sandbox with GitHub and
SonarQube MCP servers, then progressively add functionality until you have a
production-ready workflow that analyzes code quality and creates pull requests.

By working through each step sequentially, you'll learn how MCP servers work,
how to interact with them through Claude, and how to chain operations together
to build powerful automation workflows.

## Prerequisites

Before you begin, make sure you have:

- E2B account with [API access](https://e2b.dev/docs/api-key)
- [Anthropic API key](https://docs.claude.com/en/api/admin-api/apikeys/get-api-key)

  > [!NOTE]
  >
  > This example uses Claude CLI which comes pre-installed in E2B sandboxes, but you can adapt the example to work with other AI assistants of your choice. See [E2B's MCP documentation](https://e2b.dev/docs/mcp/quickstart) for alternative connection methods.

- GitHub account with:
  - A repository containing code to analyze
  - [Personal access token](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/managing-your-personal-access-tokens) with `repo` scope
- SonarCloud account with:
  - [Organization](https://docs.sonarsource.com/sonarqube-cloud/administering-sonarcloud/resources-structure/organization) created
  - [Project configured](https://docs.sonarsource.com/sonarqube-community-build/project-administration/creating-and-importing-projects) for your repository
  - [User token](https://docs.sonarsource.com/sonarqube-server/instance-administration/security/administering-tokens) generated
- Language runtime installed:
  - TypeScript: [Node.js 18+](https://nodejs.org/en/download)
  - Python: [Python 3.8+](https://www.python.org/downloads/)

> [!NOTE]
>
> This guide uses Claude's `--dangerously-skip-permissions` flag to enable
> automated command execution in E2B sandboxes. This flag bypasses permission
> prompts, which is appropriate for isolated container environments like E2B
> where sandboxes are disposable and separate from your local machine.
>
> However, be aware that Claude can execute any commands within the sandbox,
> including accessing files and credentials available in that environment. Only
> use this approach with trusted code and workflows. For more information,
> see [Anthropic's guidance on container security](https://docs.anthropic.com/en/docs/claude-code/devcontainer).

## Step 1: Set up your project

{{< tabs group="language" >}}
{{< tab name="TypeScript" >}}

1. Create a new directory for your workflow and initialize Node.js:

```bash
mkdir github-sonarqube-workflow
cd github-sonarqube-workflow
npm init -y
```

2. Open `package.json` and configure it for ES modules:

```json
{
  "name": "github-sonarqube-workflow",
  "version": "1.0.0",
  "description": "Automated code quality workflow using E2B, GitHub, and SonarQube",
  "type": "module",
  "main": "quality-workflow.ts",
  "scripts": {
    "start": "tsx quality-workflow.ts"
  },
  "keywords": ["e2b", "github", "sonarqube", "mcp", "code-quality"],
  "author": "",
  "license": "MIT"
}
```

3. Install required dependencies:

```bash
npm install e2b dotenv
npm install -D typescript tsx @types/node
```

4. Create a `.env` file in your project root:

```bash
touch .env
```

5. Add your API keys and configuration, replacing the placeholders with your actual credentials:

```plaintext
E2B_API_KEY=your_e2b_api_key_here
ANTHROPIC_API_KEY=your_anthropic_api_key_here
GITHUB_TOKEN=ghp_your_personal_access_token_here
GITHUB_OWNER=your_github_username
GITHUB_REPO=your_repository_name
SONARQUBE_ORG=your_sonarcloud_org_key
SONARQUBE_TOKEN=your_sonarqube_user_token
SONARQUBE_URL=https://sonarcloud.io
```

6. Protect your credentials by adding `.env` to `.gitignore`:

```bash
echo ".env" >> .gitignore
echo "node_modules/" >> .gitignore
```

{{< /tab >}}
{{< tab name="Python" >}}

1. Create a new directory for your workflow:

```bash
mkdir github-sonarqube-workflow
cd github-sonarqube-workflow
```

2. Create a virtual environment and activate it:

```bash
python3 -m venv venv
source venv/bin/activate  # On Windows: venv\Scripts\activate
```

3. Install required dependencies:

```bash
pip install e2b python-dotenv
```

4. Create a `.env` file in your project root:

```bash
touch .env
```

5. Add your API keys and configuration, replacing the placeholders with your actual credentials:

```plaintext
E2B_API_KEY=your_e2b_api_key_here
ANTHROPIC_API_KEY=your_anthropic_api_key_here
GITHUB_TOKEN=ghp_your_personal_access_token_here
GITHUB_OWNER=your_github_username
GITHUB_REPO=your_repository_name
SONARQUBE_ORG=your_sonarcloud_org_key
SONARQUBE_TOKEN=your_sonarqube_user_token
SONARQUBE_URL=https://sonarcloud.io
```

6. Protect your credentials by adding `.env` to `.gitignore`:

```bash
echo ".env" >> .gitignore
echo "venv/" >> .gitignore
echo "__pycache__/" >> .gitignore
```

{{< /tab >}}
{{< /tabs >}}

## Step 2: Create your first sandbox

Let's start by creating a sandbox and verifying the MCP servers are configured correctly.

{{< tabs group="language" >}}
{{< tab name="TypeScript" >}}

Create a file named `01-test-connection.ts` in your project root:

```typescript
import "dotenv/config";
import { Sandbox } from "e2b";

async function testConnection() {
  console.log(
    "Creating E2B sandbox with GitHub and SonarQube MCP servers...\n",
  );

  const sbx = await Sandbox.betaCreate({
    envs: {
      ANTHROPIC_API_KEY: process.env.ANTHROPIC_API_KEY!,
      GITHUB_TOKEN: process.env.GITHUB_TOKEN!,
      SONARQUBE_TOKEN: process.env.SONARQUBE_TOKEN!,
    },
    mcp: {
      githubOfficial: {
        githubPersonalAccessToken: process.env.GITHUB_TOKEN!,
      },
      sonarqube: {
        org: process.env.SONARQUBE_ORG!,
        token: process.env.SONARQUBE_TOKEN!,
        url: "https://sonarcloud.io",
      },
    },
  });

  const mcpUrl = sbx.betaGetMcpUrl();
  const mcpToken = await sbx.betaGetMcpToken();

  console.log(" Sandbox created successfully!");
  console.log(`MCP Gateway URL: ${mcpUrl}\n`);

  // Wait for MCP initialization
  await new Promise((resolve) => setTimeout(resolve, 1000));

  // Configure Claude to use the MCP gateway
  console.log("Connecting Claude CLI to MCP gateway...");
  await sbx.commands.run(
    `claude mcp add --transport http e2b-mcp-gateway ${mcpUrl} --header "Authorization: Bearer ${mcpToken}"`,
    {
      timeoutMs: 0,
      onStdout: console.log,
      onStderr: console.log,
    },
  );

  console.log("\nConnection successful! Cleaning up...");
  await sbx.kill();
}

testConnection().catch(console.error);
```

Run this script to verify your setup:

```bash
npx tsx 01-test-connection.ts
```

{{< /tab >}}
{{< tab name="Python" >}}

Create a file named `01_test_connection.py` in your project root:

```python
import os
import asyncio
from dotenv import load_dotenv
from e2b import AsyncSandbox

load_dotenv()

async def test_connection():
    print("Creating E2B sandbox with GitHub and SonarQube MCP servers...\n")

    sbx = await AsyncSandbox.beta_create(
        envs={
            "ANTHROPIC_API_KEY": os.getenv("ANTHROPIC_API_KEY"),
            "GITHUB_TOKEN": os.getenv("GITHUB_TOKEN"),
            "SONARQUBE_TOKEN": os.getenv("SONARQUBE_TOKEN"),
        },
        mcp={
            "githubOfficial": {
                "githubPersonalAccessToken": os.getenv("GITHUB_TOKEN"),
            },
            "sonarqube": {
                "org": os.getenv("SONARQUBE_ORG"),
                "token": os.getenv("SONARQUBE_TOKEN"),
                "url": "https://sonarcloud.io",
            },
        },
    )

    mcp_url = sbx.beta_get_mcp_url()
    mcp_token = await sbx.beta_get_mcp_token()

    print(" Sandbox created successfully!")
    print(f"MCP Gateway URL: {mcp_url}\n")

    # Wait for MCP initialization
    await asyncio.sleep(1)

    # Configure Claude to use the MCP gateway
    print("Connecting Claude CLI to MCP gateway...")
    await sbx.commands.run(
        f'claude mcp add --transport http e2b-mcp-gateway {mcp_url} --header "Authorization: Bearer {mcp_token}"',
        timeout=0,
        on_stdout=print,
        on_stderr=print,
    )

    print("\n Connection successful! Cleaning up...")
    await sbx.kill()

if __name__ == "__main__":
    asyncio.run(test_connection())
```

Run this script to verify your setup:

```bash
python 01_test_connection.py
```

{{< /tab >}}
{{< /tabs >}}

Your output should look similar to the following example:

```console {collapse=true}
Creating E2B sandbox with GitHub and SonarQube MCP servers...

✓ Sandbox created successfully!
MCP Gateway URL: https://50005-xxxxx.e2b.app/mcp

Connecting Claude CLI to MCP gateway...
Added HTTP MCP server e2b-mcp-gateway with URL: https://50005-xxxxx.e2b.app/mcp to local config
Headers: {
  "Authorization": "Bearer xxxxx-xxxx-xxxx"
}
File modified: /home/user/.claude.json [project: /home/user]

✓ Connection successful! Cleaning up...
```

You've just learned how to create an E2B sandbox with multiple MCP servers
configured. The `betaCreate` method initializes a cloud environment
with Claude CLI and your specified MCP servers.

## Step 3: Discover available MCP tools

MCP servers expose tools that Claude can call. The GitHub MCP server provides
repository management tools, while SonarQube provides code analysis tools.
By listing their tools, you know what operations are possible.

To try listing MCP tools:

{{< tabs group="language" >}}
{{< tab name="TypeScript" >}}

Create `02-list-tools.ts`:

```typescript
import "dotenv/config";
import { Sandbox } from "e2b";

async function listTools() {
  console.log("Creating sandbox...\n");

  const sbx = await Sandbox.betaCreate({
    envs: {
      ANTHROPIC_API_KEY: process.env.ANTHROPIC_API_KEY!,
      GITHUB_TOKEN: process.env.GITHUB_TOKEN!,
      SONARQUBE_TOKEN: process.env.SONARQUBE_TOKEN!,
    },
    mcp: {
      githubOfficial: {
        githubPersonalAccessToken: process.env.GITHUB_TOKEN!,
      },
      sonarqube: {
        org: process.env.SONARQUBE_ORG!,
        token: process.env.SONARQUBE_TOKEN!,
        url: "https://sonarcloud.io",
      },
    },
  });

  const mcpUrl = sbx.betaGetMcpUrl();
  const mcpToken = await sbx.betaGetMcpToken();

  // Wait for MCP initialization
  await new Promise((resolve) => setTimeout(resolve, 1000));

  await sbx.commands.run(
    `claude mcp add --transport http e2b-mcp-gateway ${mcpUrl} --header "Authorization: Bearer ${mcpToken}"`,
    { timeoutMs: 0, onStdout: console.log, onStderr: console.log },
  );

  console.log("\nDiscovering available MCP tools...\n");

  const prompt =
    "List all MCP tools you have access to. For each tool, show its exact name and a brief description.";

  await sbx.commands.run(
    `echo '${prompt}' | claude -p --dangerously-skip-permissions`,
    { timeoutMs: 0, onStdout: console.log, onStderr: console.log },
  );

  await sbx.kill();
}

listTools().catch(console.error);
```

Run the script:

```bash
npx tsx 02-list-tools.ts
```

{{< /tab >}}
{{< tab name="Python" >}}

Create `02_list_tools.py`:

```python
import os
import asyncio
from dotenv import load_dotenv
from e2b import AsyncSandbox

load_dotenv()

async def list_tools():
    print("Creating sandbox...\n")

    sbx = await AsyncSandbox.beta_create(
        envs={
            "ANTHROPIC_API_KEY": os.getenv("ANTHROPIC_API_KEY"),
            "GITHUB_TOKEN": os.getenv("GITHUB_TOKEN"),
            "SONARQUBE_TOKEN": os.getenv("SONARQUBE_TOKEN"),
        },
        mcp={
            "githubOfficial": {
                "githubPersonalAccessToken": os.getenv("GITHUB_TOKEN"),
            },
            "sonarqube": {
                "org": os.getenv("SONARQUBE_ORG"),
                "token": os.getenv("SONARQUBE_TOKEN"),
                "url": "https://sonarcloud.io",
            },
        },
    )

    mcp_url = sbx.beta_get_mcp_url()
    mcp_token = await sbx.beta_get_mcp_token()

    # Wait for MCP initialization
    await asyncio.sleep(1)

    await sbx.commands.run(
        f'claude mcp add --transport http e2b-mcp-gateway {mcp_url} --header "Authorization: Bearer {mcp_token}"',
        timeout=0,
        on_stdout=print,
        on_stderr=print,
    )

    print("\nDiscovering available MCP tools...\n")

    prompt = "List all MCP tools you have access to. For each tool, show its exact name and a brief description."

    await sbx.commands.run(
        f"echo '{prompt}' | claude -p --dangerously-skip-permissions",
        timeout=0,
        on_stdout=print,
        on_stderr=print,
    )

    await sbx.kill()

if __name__ == "__main__":
    asyncio.run(list_tools())
```

Run the script:

```bash
python 02_list_tools.py
```

{{< /tab >}}
{{< /tabs >}}

In the console, you should see a list of MCP tools:

```console {collapse=true}
Creating sandbox...

Sandbox created
Connecting to MCP gateway...

Discovering available MCP tools...

I have access to the following MCP tools:

**GitHub Tools:**
1. mcp__create_repository - Create a new GitHub repository
2. mcp__list_issues - List issues in a repository
3. mcp__create_issue - Create a new issue
4. mcp__get_file_contents - Get file contents from a repository
5. mcp__create_or_update_file - Create or update files in a repository
6. mcp__create_pull_request - Create a pull request
7. mcp__create_branch - Create a new branch
8. mcp__push_files - Push multiple files in a single commit
... (30+ more GitHub tools)

**SonarQube Tools:**
1. mcp__get_projects - List projects in organization
2. mcp__get_quality_gate_status - Get quality gate status for a project
3. mcp__list_project_issues - List quality issues in a project
4. mcp__search_issues - Search for specific quality issues
... (SonarQube analysis tools)
```

## Step 4: Test GitHub MCP tools

Let's try testing GitHub using MCP tools. We'll start simple by listing
repository issues.

{{< tabs group="language" >}}
{{< tab name="TypeScript" >}}

Create `03-test-github.ts`:

```typescript
import "dotenv/config";
import { Sandbox } from "e2b";

async function testGitHub() {
  console.log("Creating sandbox...\n");

  const sbx = await Sandbox.betaCreate({
    envs: {
      ANTHROPIC_API_KEY: process.env.ANTHROPIC_API_KEY!,
      GITHUB_TOKEN: process.env.GITHUB_TOKEN!,
    },
    mcp: {
      githubOfficial: {
        githubPersonalAccessToken: process.env.GITHUB_TOKEN!,
      },
    },
  });

  const mcpUrl = sbx.betaGetMcpUrl();
  const mcpToken = await sbx.betaGetMcpToken();

  await new Promise((resolve) => setTimeout(resolve, 1000));

  await sbx.commands.run(
    `claude mcp add --transport http e2b-mcp-gateway ${mcpUrl} --header "Authorization: Bearer ${mcpToken}"`,
    { timeoutMs: 0, onStdout: console.log, onStderr: console.log },
  );

  const repoPath = `${process.env.GITHUB_OWNER}/${process.env.GITHUB_REPO}`;

  console.log(`\nListing issues in ${repoPath}...\n`);

  const prompt = `Using the GitHub MCP tools, list all open issues in the repository "${repoPath}". Show the issue number, title, and author for each.`;

  await sbx.commands.run(
    `echo '${prompt.replace(/'/g, "'\\''")}' | claude -p --dangerously-skip-permissions`,
    {
      timeoutMs: 0,
      onStdout: console.log,
      onStderr: console.log,
    },
  );

  await sbx.kill();
}

testGitHub().catch(console.error);
```

Run the script:

```bash
npx tsx 03-test-github.ts
```

{{< /tab >}}
{{< tab name="Python" >}}

Create `03_test_github.py`:

```python
import os
import asyncio
from dotenv import load_dotenv
from e2b import AsyncSandbox

load_dotenv()

async def test_github():
    print("Creating sandbox...\n")

    sbx = await AsyncSandbox.beta_create(
        envs={
            "ANTHROPIC_API_KEY": os.getenv("ANTHROPIC_API_KEY"),
            "GITHUB_TOKEN": os.getenv("GITHUB_TOKEN"),
        },
        mcp={
            "githubOfficial": {
                "githubPersonalAccessToken": os.getenv("GITHUB_TOKEN"),
            },
        },
    )

    mcp_url = sbx.beta_get_mcp_url()
    mcp_token = await sbx.beta_get_mcp_token()

    await asyncio.sleep(1)

    await sbx.commands.run(
        f'claude mcp add --transport http e2b-mcp-gateway {mcp_url} --header "Authorization: Bearer {mcp_token}"',
        timeout=0,
        on_stdout=print,
        on_stderr=print,
    )

    repo_path = f"{os.getenv('GITHUB_OWNER')}/{os.getenv('GITHUB_REPO')}"

    print(f"\nListing issues in {repo_path}...\n")

    prompt = f'Using the GitHub MCP tools, list all open issues in the repository "{repo_path}". Show the issue number, title, and author for each.'

    await sbx.commands.run(
        f"echo '{prompt}' | claude -p --dangerously-skip-permissions",
        timeout=0,
        on_stdout=print,
        on_stderr=print,
    )

    await sbx.kill()

if __name__ == "__main__":
    asyncio.run(test_github())
```

Run the script:

```bash
python 03_test_github.py
```

{{< /tab >}}
{{< /tabs >}}

You should see Claude use the GitHub MCP tools to list your repository's issues:

```console {collapse=true}
Creating sandbox...
Connecting to MCP gateway...

Listing issues in <your-repo>...

Here are the first 10 open issues in the <your-repo> repository:

1. **Issue #23577**: Update README (author: user1)
2. **Issue #23575**: release-notes for Compose v2.40.1 version (author: user2)
3. **Issue #23570**: engine-cli: fix `docker volume prune` output (author: user3)
4. **Issue #23568**: Engdocs update (author: user4)
5. **Issue #23565**: add new section (author: user5)
... (continues with more issues)
```

You can now send prompts to Claude and interact with GitHub through
natural language. Claude decides what tool to call based on your prompt.

## Step 5: Test SonarQube MCP tools

Let's analyze code quality using SonarQube MCP tools.

{{< tabs group="language" >}}
{{< tab name="TypeScript" >}}

Create `04-test-sonarqube.ts`:

```typescript
import "dotenv/config";
import { Sandbox } from "e2b";

async function testSonarQube() {
  console.log("Creating sandbox...\n");

  const sbx = await Sandbox.betaCreate({
    envs: {
      ANTHROPIC_API_KEY: process.env.ANTHROPIC_API_KEY!,
      GITHUB_TOKEN: process.env.GITHUB_TOKEN!,
      SONARQUBE_TOKEN: process.env.SONARQUBE_TOKEN!,
    },
    mcp: {
      githubOfficial: {
        githubPersonalAccessToken: process.env.GITHUB_TOKEN!,
      },
      sonarqube: {
        org: process.env.SONARQUBE_ORG!,
        token: process.env.SONARQUBE_TOKEN!,
        url: "https://sonarcloud.io",
      },
    },
  });

  const mcpUrl = sbx.betaGetMcpUrl();
  const mcpToken = await sbx.betaGetMcpToken();

  await new Promise((resolve) => setTimeout(resolve, 1000));

  await sbx.commands.run(
    `claude mcp add --transport http e2b-mcp-gateway ${mcpUrl} --header "Authorization: Bearer ${mcpToken}"`,
    { timeoutMs: 0, onStdout: console.log, onStderr: console.log },
  );

  console.log("\nAnalyzing code quality with SonarQube...\n");

  const prompt = `Using the SonarQube MCP tools:
    1. List all projects in my organization
    2. For the first project, show:
    - Quality gate status (pass/fail)
    - Number of bugs
    - Number of code smells
    - Number of security vulnerabilities
    3. List the top 5 most critical issues found`;

  await sbx.commands.run(
    `echo '${prompt.replace(/'/g, "'\\''")}' | claude -p --dangerously-skip-permissions`,
    {
      timeoutMs: 0,
      onStdout: console.log,
      onStderr: console.log,
    },
  );

  await sbx.kill();
}

testSonarQube().catch(console.error);
```

Run the script:

```bash
npx tsx 04-test-sonarqube.ts
```

{{< /tab >}}
{{< tab name="Python" >}}

Create `04_test_sonarqube.py`:

```python
import os
import asyncio
from dotenv import load_dotenv
from e2b import AsyncSandbox

load_dotenv()

async def test_sonarqube():
    print("Creating sandbox...\n")

    sbx = await AsyncSandbox.beta_create(
        envs={
            "ANTHROPIC_API_KEY": os.getenv("ANTHROPIC_API_KEY"),
            "GITHUB_TOKEN": os.getenv("GITHUB_TOKEN"),
            "SONARQUBE_TOKEN": os.getenv("SONARQUBE_TOKEN"),
        },
        mcp={
            "githubOfficial": {
                "githubPersonalAccessToken": os.getenv("GITHUB_TOKEN"),
            },
            "sonarqube": {
                "org": os.getenv("SONARQUBE_ORG"),
                "token": os.getenv("SONARQUBE_TOKEN"),
                "url": "https://sonarcloud.io",
            },
        },
    )

    mcp_url = sbx.beta_get_mcp_url()
    mcp_token = await sbx.beta_get_mcp_token()

    await asyncio.sleep(1)

    await sbx.commands.run(
        f'claude mcp add --transport http e2b-mcp-gateway {mcp_url} --header "Authorization: Bearer {mcp_token}"',
        timeout=0,
        on_stdout=print,
        on_stderr=print,
    )

    print("\nAnalyzing code quality with SonarQube...\n")

    prompt = """Using the SonarQube MCP tools:
    1. List all projects in my organization
    2. For the first project, show:
    - Quality gate status (pass/fail)
    - Number of bugs
    - Number of code smells
    - Number of security vulnerabilities
    3. List the top 5 most critical issues found"""

    await sbx.commands.run(
        f"echo '{prompt}' | claude -p --dangerously-skip-permissions",
        timeout=0,
        on_stdout=print,
        on_stderr=print,
    )

    await sbx.kill()

if __name__ == "__main__":
    asyncio.run(test_sonarqube())
```

Run the script:

```bash
python 04_test_sonarqube.py
```

{{< /tab >}}
{{< /tabs >}}

> [!NOTE]
>
> This script may take a few minutes to run.

You should see Claude output SonarQube analysis results:

```console {collapse=true}
Creating sandbox...

Analyzing code quality with SonarQube...

## SonarQube Analysis Results

### 1. Projects in Your Organization

Found **1 project**:
- **Project Name**: project-1
- **Project Key**: project-testing

### 2. Project Analysis

...

### 3. Top 5 Most Critical Issues

Found 1 total issues (all are code smells with no critical/blocker severity):

1. **MAJOR Severity** - test.js:2
   - **Rule**: javascript:S1854
   - **Message**: Remove this useless assignment to variable "unusedVariable"
   - **Status**: OPEN

**Summary**: The project is in good health with no bugs or vulnerabilities detected.
```

You can now use SonarQube MCP tools to analyze code quality through
natural language. You can retrieve quality metrics, identify issues,
and understand what code needs fixing.

## Step 6: Create a branch and make code changes

Now, let's teach Claude to fix code based on quality issues discovered
by SonarQube.

{{< tabs group="language" >}}
{{< tab name="TypeScript" >}}

Create `05-fix-code-issue.ts`:

```typescript
import "dotenv/config";
import { Sandbox } from "e2b";

async function fixCodeIssue() {
  console.log("Creating sandbox...\n");

  const sbx = await Sandbox.betaCreate({
    envs: {
      ANTHROPIC_API_KEY: process.env.ANTHROPIC_API_KEY!,
      GITHUB_TOKEN: process.env.GITHUB_TOKEN!,
      SONARQUBE_TOKEN: process.env.SONARQUBE_TOKEN!,
    },
    mcp: {
      githubOfficial: {
        githubPersonalAccessToken: process.env.GITHUB_TOKEN!,
      },
      sonarqube: {
        org: process.env.SONARQUBE_ORG!,
        token: process.env.SONARQUBE_TOKEN!,
        url: "https://sonarcloud.io",
      },
    },
  });

  const mcpUrl = sbx.betaGetMcpUrl();
  const mcpToken = await sbx.betaGetMcpToken();

  await new Promise((resolve) => setTimeout(resolve, 1000));

  await sbx.commands.run(
    `claude mcp add --transport http e2b-mcp-gateway ${mcpUrl} --header "Authorization: Bearer ${mcpToken}"`,
    { timeoutMs: 0, onStdout: console.log, onStderr: console.log },
  );

  const repoPath = `${process.env.GITHUB_OWNER}/${process.env.GITHUB_REPO}`;
  const branchName = `quality-fix-${Date.now()}`;

  console.log("\nFixing a code quality issue...\n");

  const prompt = `Using GitHub and SonarQube MCP tools:

    1. Analyze code quality in repository "${repoPath}" with SonarQube
    2. Find ONE simple issue that can be confidently fixed (like an unused variable or code smell)
    3. Create a new branch called "${branchName}"
    4. Read the file containing the issue using GitHub tools
    5. Fix the issue in the code
    6. Commit the fix to the new branch with a clear commit message

    Important: Only fix issues you're 100% confident about. Explain what you're fixing and why.`;

  await sbx.commands.run(
    `echo '${prompt.replace(/'/g, "'\\''")}' | claude -p --dangerously-skip-permissions`,
    {
      timeoutMs: 0,
      onStdout: console.log,
      onStderr: console.log,
    },
  );

  console.log(`\nCheck your repository for branch: ${branchName}`);

  await sbx.kill();
}

fixCodeIssue().catch(console.error);
```

Run the script:

```bash
npx tsx 05-fix-code-issue.ts
```

{{< /tab >}}
{{< tab name="Python" >}}

Create `05_fix_code_issue.py`:

```python
import os
import asyncio
import time
from dotenv import load_dotenv
from e2b import AsyncSandbox

load_dotenv()

async def fix_code_issue():
    print("Creating sandbox...\n")

    sbx = await AsyncSandbox.beta_create(
        envs={
            "ANTHROPIC_API_KEY": os.getenv("ANTHROPIC_API_KEY"),
            "GITHUB_TOKEN": os.getenv("GITHUB_TOKEN"),
            "SONARQUBE_TOKEN": os.getenv("SONARQUBE_TOKEN"),
        },
        mcp={
            "githubOfficial": {
                "githubPersonalAccessToken": os.getenv("GITHUB_TOKEN"),
            },
            "sonarqube": {
                "org": os.getenv("SONARQUBE_ORG"),
                "token": os.getenv("SONARQUBE_TOKEN"),
                "url": "https://sonarcloud.io",
            },
        },
    )

    mcp_url = sbx.beta_get_mcp_url()
    mcp_token = await sbx.beta_get_mcp_token()

    await asyncio.sleep(1)

    await sbx.commands.run(
        f'claude mcp add --transport http e2b-mcp-gateway {mcp_url} --header "Authorization: Bearer {mcp_token}"',
        timeout=0,
        on_stdout=print,
        on_stderr=print,
    )

    repo_path = f"{os.getenv('GITHUB_OWNER')}/{os.getenv('GITHUB_REPO')}"
    branch_name = f"quality-fix-{int(time.time() * 1000)}"

    print("\nFixing a code quality issue...\n")

    prompt = f"""Using GitHub and SonarQube MCP tools:

    1. Analyze code quality in repository "{repo_path}" with SonarQube
    2. Find ONE simple issue that can be confidently fixed (like an unused variable or code smell)
    3. Create a new branch called "{branch_name}"
    4. Read the file containing the issue using GitHub tools
    5. Fix the issue in the code
    6. Commit the fix to the new branch with a clear commit message

    Important: Only fix issues you're 100% confident about. Explain what you're fixing and why."""

    await sbx.commands.run(
        f"echo '{prompt}' | claude -p --dangerously-skip-permissions",
        timeout=0,
        on_stdout=print,
        on_stderr=print,
    )

    print(f"\n Check your repository for branch: {branch_name}")

    await sbx.kill()

if __name__ == "__main__":
    asyncio.run(fix_code_issue())
```

Run the script:

```bash
python 05_fix_code_issue.py
```

{{< /tab >}}
{{< /tabs >}}

> [!NOTE]
>
> This script may take a few minutes to run.

Claude will analyze your repository and fix a code quality issue:

```console {collapse=true}
Creating sandbox...

Fixing a code quality issue...

## Task Completed Successfully!

I've successfully analyzed the repository and fixed a code quality issue. Here's what I did:

### 1. Analysis
Found 1 SonarQube issues in the `<your-repo>` repository:
- **Unused variable** (MINOR severity, line 2)
- Useless assignment (MAJOR severity, line 2)
- For-of loop suggestion (MINOR severity, lines 5-7)

### 2. Issue Fixed
**Location:** test.js:2
**Issue:** Unused variable declaration
**SonarQube Rule:** javascript:S1481
**What was wrong:** The line `const unusedVariable = "this is not used";` declared a variable that was never used anywhere in the function.

### 3. Why This Fix is Safe
- The variable was only declared and assigned, never referenced
- Removing it has zero impact on functionality
- It eliminates dead code and improves code quality
- This is a textbook example of a code smell that should be removed

### 4. Changes Made
- Created branch: `quality-fix-xxxxx`
- Removed the unused variable declaration
- Committed with a clear message explaining the fix and referencing the SonarQube rule

**Commit SHA:** `xxxxxxxxxxxxxxxxxxxxxxxxx`
**Branch:** https://github.com/<github-org>/<your-repo>/tree/quality-fix-xxxxx

The fix is ready for review and can be merged to improve the code quality of the repository!
```

You can now use GitHub and SonarQube MCP tools in the same workflow to read
files, make code changes, and commit them.

## Step 7: Create quality-gated pull requests

Finally, let's build the complete workflow: analyze quality, fix issues,
and create a PR only if improvements are made.

{{< tabs group="language" >}}
{{< tab name="TypeScript" >}}

Create `06-quality-gated-pr.ts`:

```typescript
import "dotenv/config";
import { Sandbox } from "e2b";

async function qualityGatedPR() {
  console.log("Creating sandbox for quality-gated PR workflow...\n");

  const sbx = await Sandbox.betaCreate({
    envs: {
      ANTHROPIC_API_KEY: process.env.ANTHROPIC_API_KEY!,
      GITHUB_TOKEN: process.env.GITHUB_TOKEN!,
      SONARQUBE_TOKEN: process.env.SONARQUBE_TOKEN!,
    },
    mcp: {
      githubOfficial: {
        githubPersonalAccessToken: process.env.GITHUB_TOKEN!,
      },
      sonarqube: {
        org: process.env.SONARQUBE_ORG!,
        token: process.env.SONARQUBE_TOKEN!,
        url: "https://sonarcloud.io",
      },
    },
  });

  const mcpUrl = sbx.betaGetMcpUrl();
  const mcpToken = await sbx.betaGetMcpToken();

  await new Promise((resolve) => setTimeout(resolve, 1000));

  await sbx.commands.run(
    `claude mcp add --transport http e2b-mcp-gateway ${mcpUrl} --header "Authorization: Bearer ${mcpToken}"`,
    { timeoutMs: 0, onStdout: console.log, onStderr: console.log },
  );

  const repoPath = `${process.env.GITHUB_OWNER}/${process.env.GITHUB_REPO}`;
  const branchName = `quality-improvements-${Date.now()}`;

  console.log("\nRunning quality-gated PR workflow...\n");

  const prompt = `You are a code quality engineer. Using GitHub and SonarQube MCP tools:

    STEP 1: ANALYSIS
    - Get current code quality status from SonarQube for "${repoPath}"
    - Record the current number of bugs, code smells, and vulnerabilities
    - Identify 1-3 issues that you can confidently fix

    STEP 2: FIX ISSUES
    - Create branch "${branchName}"
    - For each issue you're fixing:
        * Read the file with the issue
        * Make the fix
        * Commit with a descriptive message
    - Only fix issues where you're 100% confident the fix is correct

    STEP 3: VERIFICATION
        - After your fixes, check if quality metrics would improve
        - Calculate: Would this reduce bugs/smells/vulnerabilities?

    STEP 4: QUALITY GATE
        - Only proceed if your changes improve quality
        - If quality would not improve, explain why and stop

    STEP 5: CREATE PR (only if quality gate passes)
        - Create a pull request from "${branchName}" to main
        - Title: "Quality improvements: [describe what you fixed]"
        - Description should include:
            * What issues you fixed
            * Before/after quality metrics
            * Why these fixes improve code quality
        - Add a comment with detailed SonarQube analysis

    Be thorough and explain your decisions at each step.`;

  await sbx.commands.run(
    `echo '${prompt.replace(/'/g, "'\\''")}' | claude -p --dangerously-skip-permissions`,
    {
      timeoutMs: 0,
      onStdout: console.log,
      onStderr: console.log,
    },
  );

  console.log(`\n Workflow complete! Check ${repoPath} for new pull request.`);

  await sbx.kill();
}

qualityGatedPR().catch(console.error);
```

Run the script:

```bash
npx tsx 06-quality-gated-pr.ts
```

{{< /tab >}}
{{< tab name="Python" >}}

Create `06_quality_gated_pr.py`:

```python
import os
import asyncio
import time
from dotenv import load_dotenv
from e2b import AsyncSandbox

load_dotenv()

async def quality_gated_pr():
    print("Creating sandbox for quality-gated PR workflow...\n")

    sbx = await AsyncSandbox.beta_create(
        envs={
            "ANTHROPIC_API_KEY": os.getenv("ANTHROPIC_API_KEY"),
            "GITHUB_TOKEN": os.getenv("GITHUB_TOKEN"),
            "SONARQUBE_TOKEN": os.getenv("SONARQUBE_TOKEN"),
        },
        mcp={
            "githubOfficial": {
                "githubPersonalAccessToken": os.getenv("GITHUB_TOKEN"),
            },
            "sonarqube": {
                "org": os.getenv("SONARQUBE_ORG"),
                "token": os.getenv("SONARQUBE_TOKEN"),
                "url": "https://sonarcloud.io",
            },
        },
    )

    mcp_url = sbx.beta_get_mcp_url()
    mcp_token = await sbx.beta_get_mcp_token()

    await asyncio.sleep(1)

    await sbx.commands.run(
        f'claude mcp add --transport http e2b-mcp-gateway {mcp_url} --header "Authorization: Bearer {mcp_token}"',
        timeout=0,
        on_stdout=print,
        on_stderr=print,
    )

    repo_path = f"{os.getenv('GITHUB_OWNER')}/{os.getenv('GITHUB_REPO')}"
    branch_name = f"quality-improvements-{int(time.time() * 1000)}"

    print("\nRunning quality-gated PR workflow...\n")

    prompt = f"""You are a code quality engineer. Using GitHub and SonarQube MCP tools:

    STEP 1: ANALYSIS
    - Get current code quality status from SonarQube for "{repo_path}"
    - Record the current number of bugs, code smells, and vulnerabilities
    - Identify 1-3 issues that you can confidently fix

    STEP 2: FIX ISSUES
    - Create branch "{branch_name}"
    - For each issue you are fixing:
        Read the file with the issue
        Make the fix
        Commit with a descriptive message
    - Only fix issues where you are 100 percent confident the fix is correct

    STEP 3: VERIFICATION
        - After your fixes, check if quality metrics would improve
        - Calculate: Would this reduce bugs/smells/vulnerabilities?

    STEP 4: QUALITY GATE
        - Only proceed if your changes improve quality
        - If quality would not improve, explain why and stop

    STEP 5: CREATE PR (only if quality gate passes)
        - Create a pull request from "{branch_name}" to main
        - Title: "Quality improvements: [describe what you fixed]"
        - Description should include:
            What issues you fixed
            Before/after quality metrics
            Why these fixes improve code quality
        - Add a comment with detailed SonarQube analysis

    Be thorough and explain your decisions at each step."""

    await sbx.commands.run(
        f"echo '{prompt.replace(chr(39), chr(39) + chr(92) + chr(39) + chr(39))}' | claude -p --dangerously-skip-permissions",
        timeout=0,
        on_stdout=print,
        on_stderr=print,
    )

    print(f"\n Workflow complete! Check {repo_path} for new pull request.")

    await sbx.kill()

if __name__ == "__main__":
    asyncio.run(quality_gated_pr())
```

Run the script:

```bash
python 06_quality_gated_pr.py
```

{{< /tab >}}
{{< /tabs >}}

> [!NOTE]
>
> This script may take a few minutes to run.

Claude will run the entire workflow, creating a quality improvement
and opening a PR in GitHub:

```console {collapse=true}
Creating sandbox for quality-gated PR workflow...

Running quality-gated PR workflow...

## Complete - Code Quality Improvement Successfully Executed

### Summary of Work Completed

**Pull Request Created:** https://github.com/<github-org>/<your-repo>/pull/1

---

### STEP 1: ANALYSIS

**Current State:**
- **Bugs:** 0
- **Vulnerabilities:** 0
- **Code Smells:** 1
- **Lines of Code:** 12
- **All Ratings:** A (Reliability, Security, Maintainability)

**Issues Identified (all in test.js):**
1. xxxxxxxxx - Unused variable declaration (MINOR)

---

### STEP 2: FIX ISSUES

**Branch Created:** `quality-improvements-xxxxxxxx`

**Changes Made:**
1. **Removed unused variable** (line 2) - Eliminated dead code that served no purpose
2. **Modernized loop pattern** (lines 5-7) - Converted `for (let i = 0; i < items.length; i++)` to `for (const item of items)`

**Commit:** xxxxxxxxxx

---

### STEP 3: VERIFICATION

**Expected Impact:**
- Code Smells: 1 → 0 (100% reduction)
- Bugs: 0 → 0 (maintained)
- Vulnerabilities: 0 → 0 (maintained)
- All quality ratings maintained at A

---

### STEP 4: QUALITY GATE PASSED

**Decision Criteria Met:**
- ✅ Reduces code smells by 100%
- ✅ No new bugs or vulnerabilities introduced
- ✅ Code is more readable and maintainable
- ✅ Follows modern JavaScript best practices
- ✅ All fixes are low-risk refactorings with no behavioral changes

---

### STEP 5: CREATE PR

**Pull Request Details:**
- **Number:** #1
- **Title:** Quality improvements: Remove unused variable and modernize for loop
- **Branch:** quality-improvements-xxxxxxxx → main
- **URL:** https://github.com/<github-org)/<your-repo>/pull/1

**PR Includes:**
- Comprehensive description with before/after metrics
- Detailed SonarQube analysis comment with issue breakdown
- Code comparison showing improvements
- Quality metrics table

The pull request is now ready for review and merge!
```

You've now built a complete, multi-step workflow with conditional logic.
Claude analyzes quality with SonarQube, makes fixes using GitHub tools,
verifies improvements, and only creates a PR if quality actually improves.

## Step 8: Add error handling

Production workflows need error handling. Let's make our workflow more robust.

{{< tabs group="language" >}}
{{< tab name="TypeScript" >}}

Create `07-robust-workflow.ts`:

```typescript
import "dotenv/config";
import { Sandbox } from "e2b";

async function robustWorkflow() {
  let sbx: Sandbox | undefined;

  try {
    console.log("Creating sandbox...\n");

    sbx = await Sandbox.betaCreate({
      envs: {
        ANTHROPIC_API_KEY: process.env.ANTHROPIC_API_KEY!,
        GITHUB_TOKEN: process.env.GITHUB_TOKEN!,
        SONARQUBE_TOKEN: process.env.SONARQUBE_TOKEN!,
      },
      mcp: {
        githubOfficial: {
          githubPersonalAccessToken: process.env.GITHUB_TOKEN!,
        },
        sonarqube: {
          org: process.env.SONARQUBE_ORG!,
          token: process.env.SONARQUBE_TOKEN!,
          url: "https://sonarcloud.io",
        },
      },
    });

    const mcpUrl = sbx.betaGetMcpUrl();
    const mcpToken = await sbx.betaGetMcpToken();

    await new Promise((resolve) => setTimeout(resolve, 1000));

    await sbx.commands.run(
      `claude mcp add --transport http e2b-mcp-gateway ${mcpUrl} --header "Authorization: Bearer ${mcpToken}"`,
      { timeoutMs: 0, onStdout: console.log, onStderr: console.log },
    );

    const repoPath = `${process.env.GITHUB_OWNER}/${process.env.GITHUB_REPO}`;

    console.log("\nRunning workflow with error handling...\n");

    const prompt = `Run a quality improvement workflow for "${repoPath}".

    ERROR HANDLING RULES:
    1. If SonarQube is unreachable, explain the error and stop gracefully
    2. If GitHub API fails, retry once, then explain and stop
    3. If no fixable issues are found, explain why and exit (this is not an error)
    4. If file modifications fail, explain which file and why
    5. At each step, check for errors before proceeding

    Run the workflow and handle any errors you encounter professionally.`;

    await sbx.commands.run(
      `echo '${prompt.replace(/'/g, "'\\''")}' | claude -p --dangerously-skip-permissions`,
      {
        timeoutMs: 0,
        onStdout: console.log,
        onStderr: console.log,
      },
    );

    console.log("\n Workflow completed");
  } catch (error) {
    const err = error as Error;
    console.error("\n Workflow failed:", err.message);

    if (err.message.includes("403")) {
      console.error("\n Check your E2B account has MCP gateway access");
    } else if (err.message.includes("401")) {
      console.error("\n Check your API tokens are valid");
    } else if (err.message.includes("Credit balance")) {
      console.error("\n Check your Anthropic API credit balance");
    }

    process.exit(1);
  } finally {
    if (sbx) {
      console.log("\n Cleaning up sandbox...");
      await sbx.kill();
    }
  }
}

robustWorkflow().catch(console.error);
```

Run the script:

```bash
npx tsx 07-robust-workflow.ts
```

{{< /tab >}}
{{< tab name="Python" >}}

Create `07_robust_workflow.py`:

```python
import os
import asyncio
import sys
from dotenv import load_dotenv
from e2b import AsyncSandbox

load_dotenv()

async def robust_workflow():
    sbx = None

    try:
        print("Creating sandbox...\n")

        sbx = await AsyncSandbox.beta_create(
            envs={
                "ANTHROPIC_API_KEY": os.getenv("ANTHROPIC_API_KEY"),
                "GITHUB_TOKEN": os.getenv("GITHUB_TOKEN"),
                "SONARQUBE_TOKEN": os.getenv("SONARQUBE_TOKEN"),
            },
            mcp={
                "githubOfficial": {
                    "githubPersonalAccessToken": os.getenv("GITHUB_TOKEN"),
                },
                "sonarqube": {
                    "org": os.getenv("SONARQUBE_ORG"),
                    "token": os.getenv("SONARQUBE_TOKEN"),
                    "url": "https://sonarcloud.io",
                },
            },
        )

        mcp_url = sbx.beta_get_mcp_url()
        mcp_token = await sbx.beta_get_mcp_token()

        await asyncio.sleep(1)

        await sbx.commands.run(
            f'claude mcp add --transport http e2b-mcp-gateway {mcp_url} --header "Authorization: Bearer {mcp_token}"',
            timeout=0,  # Fixed: was timeout_ms
            on_stdout=print,
            on_stderr=print,
        )

        repo_path = f"{os.getenv('GITHUB_OWNER')}/{os.getenv('GITHUB_REPO')}"

        print("\nRunning workflow with error handling...\n")

        prompt = f"""Run a quality improvement workflow for "{repo_path}".

        ERROR HANDLING RULES:
        1. If SonarQube is unreachable, explain the error and stop gracefully
        2. If GitHub API fails, retry once, then explain and stop
        3. If no fixable issues are found, explain why and exit (this is not an error)
        4. If file modifications fail, explain which file and why
        5. At each step, check for errors before proceeding

        Run the workflow and handle any errors you encounter professionally."""

        await sbx.commands.run(
            f"echo '{prompt}' | claude -p --dangerously-skip-permissions",
            timeout=0,
            on_stdout=print,
            on_stderr=print,
        )

        print("\n Workflow completed")

    except Exception as error:
        print(f"\n✗ Workflow failed: {str(error)}")

        error_msg = str(error)
        if "403" in error_msg:
            print("\n Check your E2B account has MCP gateway access")
        elif "401" in error_msg:
            print("\n Check your API tokens are valid")
        elif "Credit balance" in error_msg:
            print("\n Check your Anthropic API credit balance")

        sys.exit(1)

    finally:
        if sbx:
            print("\n Cleaning up sandbox...")
            await sbx.kill()

if __name__ == "__main__":
    asyncio.run(robust_workflow())
```

Run the script:

```bash
python 07_robust_workflow.py
```

{{< /tab >}}
{{< /tabs >}}

Claude will run the entire workflow, and if it encounters an error, respond
with robust error messaging.

## Next steps

In the next section, you'll customize your workflow for your needs.
