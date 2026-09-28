---
title: Sandboxes
description: Create and manage cloud sandboxes in Docker Agentic Platform.
keywords: docker agentic platform, sandboxes, agents, cloud runtime, terminal, compute
weight: 20
aliases:
  - /agentic-platform/concepts/sandboxes/
  - /agentic-platform/guides/manage-sandboxes/
---

A sandbox is an isolated environment for running agents and tools in the cloud.
You can work with your agent through a terminal in the Console.

The [kit](kits.md) you select determines the sandbox's base image, agent, and
installed tools. Each sandbox has its own filesystem, network settings, and
terminal. Choose a compute size to set how much CPU and
memory your sandbox has. Check your kit for installed tools such as Docker
Engine. For example, the Hermes kit doesn't include it.

The Shell kit opens a Bash shell without a pre-installed agent. It uses the
same environment as [`sbx run shell`](/manuals/ai/sandboxes/agents/shell.md)
and is useful for working manually or installing your own agent.

Your sandbox keeps running when you leave the Console. It runs until you pause
or delete it. When its timer expires, it performs the action you selected.

## Choose a platform

In the compute picker, under **Platform**, choose **auto**, **linux/amd64**
(Intel/AMD), or **linux/arm64** (Arm). The default, **auto**, uses the platform
provided by the kit's image. An explicit platform choice requires the image
to support that platform; otherwise, creation fails.

## Source code and files

Each sandbox starts with a fresh filesystem. Your local repositories,
directories, and workspaces aren't mounted in it by default.

Files don't sync automatically with your computer or a remote repository.
Commit and push any work you want to keep before deleting the sandbox. Files
left only in a deleted sandbox won't be available in a later sandbox.

Use GitHub to clone source code into your sandbox and save changes remotely.
To clone a private repository or push changes, use the **GitHub token**
control in the launcher. Copilot uses this token for both the agent and GitHub
repository access. A custom kit must declare a GitHub credential to offer this
option.
Make sure your token has the required repository permissions. You can clone
public repositories without a token, but pushing to them still requires
authentication.

## Open a sandbox

After you select **Run**, the Console shows provisioning progress. When the
sandbox is ready, its detail page opens. Use the terminal to work with your
agent or shell. Reloading the launch page resumes tracking the launch.

Open **Sandboxes** to review each sandbox's name, type, status, hourly rate,
expiration, and age. Select a sandbox to reopen its detail page and terminal.

The **New** page also lists **Recent sandboxes**. Select a sandbox to reopen it,
or select **See all sandboxes** to open the full list.

In a Claude Code sandbox, paste an image from your clipboard or drag an image
file onto the terminal to attach it to your prompt. Wait for the upload to
finish, type your question, and press Enter to send it.

## Open additional shells

Use the **+** button in the terminal tab bar to open another shell in the same
sandbox. Switch tabs to work with multiple shells. The agent remains in its
own tab while you run commands in another.

## Connect from your computer

If the sandbox detail page shows **Connect**, open it to find SSH connection
options. The panel includes a Docker Sandboxes CLI option and a connect script.
For the script, manage your public SSH key under **Settings**. Follow the
instructions in the panel for your chosen connection method.

### Access a service by port

To access a web application or development server running in your sandbox:

1. Start the service in the sandbox and note its port.
2. Open **Connect** and find **Connect via a Port**.
3. Enter the port and select **Open Port**.
4. Copy the public URL to access the service.

The sandbox must be running. Port 2222 is reserved for SSH and can't be
published through this control. To stop exposing a service, use the close
action next to its port.

## Pause, resume, or delete a sandbox {#manage-the-lifecycle}

A sandbox can be running or paused:

- Pause a running sandbox to stop it without deleting its files.
- Resume a paused sandbox to continue working with it. Resuming starts a fresh
  timer with the original duration.
- Delete a sandbox when you no longer need it.

When you create a sandbox, set a timer from 1 to 24 hours and choose
what happens when it expires. **Stop** stops the sandbox, while **Delete**
deletes the sandbox and its files. If **Restart** is offered, select it to
restart the sandbox automatically when the timer expires. When a sandbox stops,
all processes inside it stop too, including background processes.

After launch, you can't change the selected credentials, tools, network
policies, or compute size.

You pay for compute by the second while your sandbox runs. Your model provider
bills inference separately. For account, usage, and payment information, see
[Docker Billing](/subscription-billing/).

## Account quotas

The following default quotas apply across your cloud sandbox account, whether
resources are created through the Console, CLI, or API:

| Resource | Default limit |
| --- | ---: |
| Concurrent sandboxes | 10 |
| Stored sandboxes | 50 |
| Volumes | 100 |
| Secrets | 100 |
| Images being prepared at the same time | 3 |

Your account can have different quotas. Confirm your account's limits with
Docker before planning a workload that depends on a particular allowance.

Stopping an ordinary sandbox releases its concurrency slot, but the sandbox
still counts toward stored usage. Resuming it needs a concurrency slot. An
always-on sandbox retains its concurrency reservation while stopped. Delete
sandboxes you no longer need to release stored usage.

If the Console reports a running sandbox limit, stop or delete a sandbox
before trying again. If it reports a sandbox storage limit, delete a sandbox;
stopping it doesn't free a stored-sandbox slot.

## Check sandbox configuration

If your agent cannot reach a service or use a tool:

- Check that the network policies allow access to the service.
- If the sandbox needs an MCP tool, confirm that its server is connected and
  authorized.
- If the service requires authentication, check that its credential was
  included when you created the sandbox.
