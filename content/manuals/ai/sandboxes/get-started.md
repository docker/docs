---
title: Get started with Docker Sandboxes
linkTitle: Get started
description: Run Claude Code in an isolated sandbox. Quick setup guide with prerequisites and essential commands.
weight: 20
---

{{< summary-bar feature_name="Docker Sandboxes" >}}

This guide shows how to run Claude Code in an isolated sandbox for the first time.

> [!NOTE]
> Upgrading from an earlier version of Docker Desktop? See the
> [migration guide](migration.md) for information about the new microVM
> architecture.

## Prerequisites

Before you begin, ensure you have:

- Docker Desktop 4.58 or later
- macOS, or Windows {{< badge color=violet text=Experimental >}}
- A Claude API key (can be provided via environment variable or interactively)

## Run your first sandbox

Follow these steps to run Claude Code:

1. (Optional but recommended) Set your Anthropic API key as an environment variable.

   Add the API key to your shell configuration file:

   ```plaintext {title="~/.bashrc or ~/.zshrc"}
   export ANTHROPIC_API_KEY=sk-ant-api03-xxxxx
   ```

   Docker Sandboxes use a daemon process that runs independently of your
   current shell session. This means setting the environment variable inline or
   in your current session will not work. You must set it globally in your
   shell configuration file to ensure the daemon can access it.

   Apply the changes:
   1. Source your shell configuration.
   2. Restart Docker Desktop so the daemon picks up the new environment variable.

   Alternatively, you can skip this step and authenticate interactively when
   Claude Code starts. If no credentials are found, you'll be prompted to log
   in. Note that interactive authentication requires you to authenticate for
   each workspace separately.

2. Create and run a sandbox for Claude Code for your workspace:

   ```console
   $ docker sandbox run claude ~/my-project
   ```

   This creates a microVM sandbox. Docker assigns it a name automatically.

3. Claude Code starts and you can begin working. The first run takes longer
   while Docker initializes the microVM and pulls the template image.

## What just happened?

When you ran `docker sandbox run`:

- Docker created a lightweight microVM with a private Docker daemon
- The sandbox was assigned a name based on the workspace path
- Your workspace synced into the VM
- Docker started the Claude Code agent as a container inside the sandbox VM

The sandbox persists until you remove it. Installed packages and configuration
remain available. Run `docker sandbox run <sandbox-name>` again to reconnect.

> [!NOTE]
> Agents can modify files in your workspace. Review changes before executing
> code or performing actions that auto-run scripts. See
> [Security considerations](workflows.md#security-considerations) for details.

## Basic commands

Here are essential commands to manage your sandboxes:

### List sandboxes

```console
$ docker sandbox ls
```

Shows all your sandboxes with their IDs, names, status, and creation time.

> [!NOTE]
> Sandboxes don't appear in `docker ps` because they're microVMs, not
> containers. Use `docker sandbox ls` to see them.

### Access a running sandbox

```console
$ docker sandbox exec -it <sandbox-name> bash
```

Executes a command inside the container in the sandbox. Use `-it` to open an
interactive shell for debugging or installing additional tools.

### Remove a sandbox

```console
$ docker sandbox rm <sandbox-name>
```

Deletes the sandbox VM and all installed packages inside it. You can remove
multiple sandboxes at once by specifying multiple names:

```console
$ docker sandbox rm <sandbox-1> <sandbox-2>
```

### Recreate a sandbox

To start fresh with a clean environment, remove and recreate the sandbox:

```console
$ docker sandbox rm <sandbox-name>
$ docker sandbox run claude ~/project
```

Configuration like custom templates and workspace paths are set when you create
the sandbox. To change these settings, remove and recreate.

For a complete list of commands and options, see the
[CLI reference](/reference/cli/docker/sandbox/).

## Next steps

Now that you have Claude running in a sandbox, learn more about:

- [Claude Code configuration](claude-code.md)
- [Supported agents](agents.md)
- [Using sandboxes effectively](workflows.md)
- [Custom templates](templates.md)
- [Network policies](network-policies.md)
- [Troubleshooting](troubleshooting.md)
