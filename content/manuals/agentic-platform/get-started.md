---
title: Get started with Docker Agentic Platform
linkTitle: Get started
description: Review a sample project with a coding agent or build and preview it from a shell in a cloud sandbox, without setting up a local development environment.
keywords: docker agentic platform, get started, cloud sandboxes, shell, agents, ssh, vscode, cursor
weight: 10
aliases:
  - /agentic-platform/concepts/platform-model/
---

Use a cloud sandbox to review a project with a coding agent or build and run
an application from a shell, without setting up the project on your computer.

Both paths introduce the same cloud sandbox lifecycle: launch an environment,
run a task, save work, and pause or delete it. Complete either path in your
browser, then use the [guides](guides/_index.md) for sharing, remote development,
or moving local work to the cloud.

## Before you begin

You need a personal Docker account with an active
[Docker Agentic Platform subscription](signup.md#activate-cloud-access).

Both workflows use a public sample project, so you don't need a GitHub token.
Review the compute price and timer before launching. Choose **Stop** as the
timeout action to keep files. Compute is billed while the sandbox runs;
closing the browser doesn't stop it.

## Choose a workflow

- [Use a coding agent](#use-a-coding-agent) to have Claude Code review a sample
  project and save its findings. You need an Anthropic API key.
- [Use a shell](#use-a-shell) to build and run the sample application and preview
  it in your browser. You don't need a model provider API key.

## Use a coding agent

Add an Anthropic API key under **Secrets** or enter it when you launch the
sandbox. Your model provider bills inference separately from Docker's sandbox
compute charges. For other agents and their credentials, see [Kits](kits.md)
and [Secrets](secrets.md).

### Launch Claude Code

1. Open [Docker Agentic Platform](https://agentic-platform.docker.com/) and
   select **+ New sandbox**.
2. Choose the Claude Code kit.
3. Add your Anthropic API key if it isn't already saved. Don't enter it in
   the agent's prompt or in a file inside the sandbox.
4. Keep the defaults and select **Run**.

On your first launch, the default **Open** network policy permits outbound
access to any host. If you're using a custom policy, make sure it permits
GitHub and your model provider. See [Network policies](policies.md).

When the sandbox is ready, its terminal opens with Claude Code.

### Give the agent a task

Your local files aren't mounted in a cloud sandbox. Ask Claude Code to clone
and review a public sample project:

```text
Clone https://github.com/docker/welcome-to-docker.git into
/home/agent/workspace/project, then review the application.
Read the source, package.json, and Dockerfile. Don't modify the application
or install dependencies.
Write /home/agent/workspace/review.md with:
- A short explanation of what the application does.
- The commands defined in package.json for running and building it.
- An explanation of how the Dockerfile builds and serves the application.
- One small improvement you'd recommend, with the relevant file path.
Base the review on the files you read. Don't claim to have run the application.
Summarize your findings in your response.
```

Claude Code summarizes the project in the conversation and saves the review
at `/home/agent/workspace/review.md`. You can ask follow-up questions or give
it another task in the same session.

Continue to [keep your work](#keep-your-work).

## Use a shell

### Launch a shell

The Shell kit starts a Bash shell without an agent.

1. Open [Docker Agentic Platform](https://agentic-platform.docker.com/) and
   select **+ New sandbox**.
2. Open the kit picker and choose **shell**.
3. If saved secrets or a GitHub token are selected, deselect them for this
   walkthrough. You don't need MCP tools either.
4. In the network policy picker, select **Balanced** to permit access to
   GitHub, Docker Hub, and npm. Deselect **Open** if selected: combining it
   with Balanced doesn't restrict Open's access. Review the compute price and
   timer, choose **Stop** as the timeout action, and select **Run**.

If you use a custom network policy, make sure it permits GitHub, Docker Hub,
and npm. See [Network policies](policies.md).

When the sandbox is ready, its terminal opens with a Bash shell.

### Build and run an application

Your local files aren't mounted in a cloud sandbox. In the browser terminal,
clone a public sample project:

```console
$ git clone https://github.com/docker/welcome-to-docker.git project
$ cd project
```

The Shell kit includes Docker Engine. Build the application's container image:

```console
$ docker build -t welcome-to-docker .
```

The build installs the application's dependencies and compiles it inside the
sandbox. When it finishes, start a container and publish its port:

```console
$ docker run -d --name welcome-to-docker -p 3000:3000 welcome-to-docker
```

### Preview the application

1. Select **Connect** on the sandbox's detail page.
2. Under **Connect via a Port**, enter `3000` and select **Open Port**.
3. Open the generated URL in your browser.

You should see the Welcome to Docker application. The application runs in the
cloud sandbox, so you can try it without installing Docker or Node.js on your
computer.

The generated URL is public, and Docker account authentication doesn't protect
it. Use only non-sensitive data. Close the port in **Connect** when you've
finished previewing, then reload the URL to confirm access is withdrawn.
For a complete sharing workflow, including checking access from a separate
browser session, follow [Share an application](guides/share-app.md).

Continue to [keep your work](#keep-your-work).

## Keep your work

The cloned project and any files you create are stored in the sandbox. Files
don't sync automatically with your computer or a remote repository. They remain available
when you pause and resume the sandbox, but are lost when you delete it.

To keep files outside the sandbox, use the CLI to
[transfer files](/manuals/ai/sandboxes/cloud/usage.md#transfer-files), or commit
and push changes to a repository. Include a
[GitHub credential](secrets.md#github-credential) when launching a sandbox that
needs to push changes.

## Pause, keep running, or delete the sandbox

Closing the browser leaves the sandbox running.
Commands and background processes can continue until you pause or delete the
sandbox, or its timer expires. Compute is billed while the sandbox runs, even
when no commands are running.

Return to **Sandboxes** in the Console and choose what to do next:

- Pause the sandbox to stop compute charges and keep its files. Confirm that
  it is paused in the list. To continue later, reopen the sandbox and resume
  it. Resuming starts a fresh timer with the original duration.
- Leave it running if the agent still has work to do or to keep the application available. Check the expiration
  time: the default **Stop** action stops the sandbox after one hour, including
  background processes, and keeps its files.
- Delete it when you no longer need it. Push any source changes
  you want to keep first. Deletion removes files stored only in the sandbox.

For more detail, see [Lifecycle actions](sandboxes.md#manage-the-lifecycle).

## Sandbox options

For launch settings, see [Kits](kits.md), [Secrets](secrets.md),
[Network policies](policies.md), and [MCP tools](mcp.md). For compute platforms
and timer actions, see [Sandboxes](sandboxes.md). Choose these settings before
launching: credentials, tools, network policies, and compute size can't be
changed afterward.

## Next steps

To run a task from your own terminal, follow the
[cloud CLI walkthrough](/manuals/ai/sandboxes/cloud/_index.md#get-started).
For programmatic workflows, see the
[Sandboxes API quickstart](/manuals/ai/sandboxes-api/get-started.md).

Continue with a targeted workflow:

- [Share an application](guides/share-app.md) through a public URL.
- [Develop remotely](guides/remote-development.md) with VS Code or Cursor.
- [Continue local work in the cloud](guides/local-handoff.md) and retrieve the results.
