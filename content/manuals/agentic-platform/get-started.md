---
title: Get started with Docker Agentic Platform
linkTitle: Get started
description: Start a shell in your first cloud sandbox, explore a sample project, save a file, and pause or delete the sandbox when you're finished.
keywords: docker agentic platform, get started, cloud sandboxes, shell, agents, ssh, vscode, cursor
weight: 10
aliases:
  - /agentic-platform/concepts/platform-model/
---

Use a cloud sandbox to work on a project without setting it up on your
computer. In this walkthrough, you'll start a shell, clone a public sample
application, and save a record of its latest commit. Then you'll pause the
sandbox to return later or delete it when you're finished.

You can complete the walkthrough in the browser. An optional step shows how
to open the same project in VS Code or Cursor over SSH.

## Before you begin

You need a personal Docker account with an active
[Docker Agentic Platform subscription](signup.md#activate-cloud-access).

The Shell kit doesn't require a model provider API key. This example uses a
public repository, so you don't need a GitHub token either. For agents and
their credentials, see [Kits](kits.md) and [Secrets](secrets.md).

## Start a sandbox

1. Open [Docker Agentic Platform](https://agentic-platform.docker.com/) and
   select **+ New sandbox**.
2. Open the kit picker and choose **shell**.
3. If saved secrets or a GitHub token are selected, deselect them for this
   walkthrough. You don't need MCP tools either.
4. In the network policy picker, select **Balanced** to permit access to
   GitHub. Keep the compute and timer defaults and select **Run**.

If you use a custom network policy, make sure it permits GitHub. See
[Network policies](policies.md).

When the sandbox is ready, its terminal opens with a Bash shell. After launch,
you can't change the selected credentials, tools, network policies, or compute
size. See [Sandbox options](#sandbox-options) for other launch settings.

## Explore a project

Your local files aren't mounted in a cloud sandbox. In the browser terminal,
clone a public sample project:

```console
$ cd /home/agent/workspace
$ git clone https://github.com/docker/welcome-to-docker.git project
$ cd project
$ ls
```

Read the project's package scripts and Dockerfile:

```console
$ cat package.json
$ cat Dockerfile
```

Save a record of the latest commit to a file in the workspace:

```console
$ git log -1 --oneline > /home/agent/workspace/latest-commit.txt
$ cat /home/agent/workspace/latest-commit.txt
```

The output shows the commit's short ID and message. The file stays in the
sandbox so you can open it later or copy it to your computer.

## Keep your work

The cloned project and commit record are stored in the sandbox. Files don't sync
automatically with your computer or a remote repository. They remain available
when you pause and resume the sandbox, but are lost when you delete it.

To keep files outside the sandbox, use the CLI to
[transfer files](/manuals/ai/sandboxes/cloud/usage.md#transfer-files), or commit
and push changes to a repository. Include a
[GitHub credential](secrets.md#github-credential) when launching a sandbox that
needs to push changes.

## Open the project in your editor (optional)

Open the project and commit record in VS Code or Cursor over SSH. Your editor
runs on your computer, while the files and terminal commands stay in the
cloud sandbox.

1. Install the [Docker Sandboxes CLI](/manuals/ai/sandboxes/install.md),
   version 0.45.0 or later.
2. In a terminal on your computer, sign in with the same Docker account you
   used in the Console and configure cloud SSH access:

   ```console
   $ sbx login
   $ sbx --cloud setup ssh
   $ sbx --cloud ls
   ```

3. Find your sandbox in the list and copy its `sbx_`-prefixed ID. Test the
   connection, replacing `<SANDBOX_ID>` with that ID:

   ```console
   $ ssh <SANDBOX_ID>@sbx_cloud
   ```

   A successful connection opens a shell inside the cloud sandbox. Run `exit`
   to return to your computer's terminal.
4. In VS Code, install the
   [Remote - SSH extension](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.remote-ssh).
   In Cursor, install its Remote - SSH support.
5. Open your editor's Command Palette and run **Remote-SSH: Connect to Host**.
   Enter `<SANDBOX_ID>@sbx_cloud`, using the same ID as in the SSH command.
6. Open `/home/agent/workspace` in the remote window. You'll find the cloned
   application in `project` and the commit record in `latest-commit.txt`.

The first connection installs the editor's server inside the sandbox.

To download the commit record, run this command in a terminal on your computer:

```console
$ sbx --cloud cp <SANDBOX_ID>:/home/agent/workspace/latest-commit.txt ./latest-commit.txt
```

## Pause, keep running, or delete the sandbox

Closing the browser or disconnecting your editor leaves the sandbox running.
Commands and background processes can continue until you pause or delete the
sandbox, or its timer expires. Compute is billed while the sandbox runs, even
when no commands are running.

Return to **Sandboxes** in the Console and choose what to do next:

- Pause the sandbox to stop compute charges and keep its files. Confirm that
  it is paused in the list. To continue later, reopen the sandbox and resume
  it. Resuming starts a fresh timer with the original duration.
- Leave it running if a command or background process still has work to do. Check the expiration
  time: the default **Stop** action stops the sandbox after one hour, including
  background processes, and keeps its files.
- Delete it when you no longer need it. Save the commit record or push any changes
  you want to keep first. Deletion removes files stored only in the sandbox.

For more detail, see [Lifecycle actions](sandboxes.md#manage-the-lifecycle).
To run a task from your own terminal, follow the
[cloud CLI walkthrough](/manuals/ai/sandboxes/cloud/_index.md#get-started).
For programmatic workflows, see the
[Sandboxes API quickstart](/manuals/ai/sandboxes-api/get-started.md).

## Sandbox options

Configure these options before selecting **Run**:

| Option | What to choose |
| --- | --- |
| Agent credentials | Add the API keys or tokens your kit needs. Saved credentials are reused by default; deselect optional credentials you don't want to include. See [Secrets](secrets.md). |
| GitHub token | Provide a token for Copilot or for cloning private repositories and pushing changes. Copilot uses the same token for both. Repository access is optional for other agents. See [GitHub credential](secrets.md#github-credential). |
| Network policies | Review the hosts and services the sandbox can reach. See [Network policies](policies.md). |
| MCP tools | Choose the servers your agent can use and authorize them if prompted. See [MCP](mcp.md). |
| Compute and platform | Choose CPU and memory resources. Under **Platform**, use **auto** or a platform supported by the kit's image. See [Sandbox platforms](sandboxes.md#choose-a-platform). |
| Timer | Under **Run for**, choose 1 to 24 hours. Under **When the time is up**, choose **Stop**, **Delete**, or **Restart** if offered. See [Lifecycle actions](sandboxes.md#manage-the-lifecycle). |

The launcher remembers your policy selections from the previous launch in the
same browser. Policies marked **Always applied** are selected automatically.
Review the selected policies before launching. Selecting **Balanced** alongside
**Open** doesn't restrict Open's access; see [How policies combine](policies.md#how-policies-combine).
