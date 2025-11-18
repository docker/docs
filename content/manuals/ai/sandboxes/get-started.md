---
title: Get started with Docker Sandboxes
linkTitle: Get started
description: Run Claude Code in an isolated sandbox. Quick setup guide with prerequisites and essential commands.
weight: 20
---

{{< summary-bar feature_name="Docker Sandboxes" >}}

This guide will help you run Claude Code in a sandboxed environment for the first time.

## Prerequisites

Before you begin, ensure you have:

- Docker Desktop 4.50 or later
- A Claude Code subscription

## Run a sandboxed agent

Follow these steps to run Claude Code in a sandboxed environment:

1. Navigate to Your Project

   ```console
   $ cd ~/my-project
   ```

2. Start Claude in a sandbox

   ```console
   $ docker sandbox run claude
   ```

3. Authenticate: on first run, Claude will prompt you to authenticate.

   Once you've authenticated, the credentials are stored in a persistent Docker
   volume and reused for future sessions.

4. Claude Code launches inside the container.

## What just happened?

When you ran `docker sandbox run claude`:

- Docker created a container from a template image
- Your current directory was mounted at the same path inside the container
- Your Git name and email were injected into the container
- Your API key was stored in a Docker volume (`docker-claude-sandbox-data`)
- Claude Code started with bypass permissions enabled

The container continues running in the background. Running `docker sandbox run
claude` again in the same directory reuses the existing container, allowing the
agent to maintain state (installed packages, temporary files) across sessions.

## Basic commands

Here are a few essential commands to manage your sandboxes:

### List your sandboxes

```console
$ docker sandbox ls
```

Shows all your sandboxes with their IDs, names, status, and creation time.

### Remove a sandbox

```console
$ docker sandbox rm <sandbox-id>
```

Deletes a sandbox when you're done with it. Get the sandbox ID from `docker sandbox ls`.


### View sandbox details

```console
$ docker sandbox inspect <sandbox-id>
```

Shows detailed information about a specific sandbox in JSON format.

For a complete list of all commands and options, see the [CLI reference](/reference/cli/docker/sandbox/).

## Next Steps

Now that you have Claude running in a sandboxed environment, learn more about:

- [Authentication strategies](claude-code.md#authentication)
- [Configuration options](claude-code.md#configuration)
- [Advanced configurations](advanced-config.md)
- [Troubleshooting guide](troubleshooting.md)
