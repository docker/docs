---
title: Gordon
description:
  AI assistant for Docker workflows - execute tasks, debug issues, and
  manage containers with intelligent assistance
weight: 40
params:
  sidebar:
    group: AI and agents
aliases:
  - /desktop/features/gordon/
---

{{< summary-bar feature_name="Gordon" >}}

Gordon is an AI-powered assistant that takes action on your Docker workflows.
It analyzes your environment, proposes solutions, and executes commands with
your permission.

## What Gordon does

Gordon takes action to help you with Docker tasks:

- Explains Docker concepts and commands
- Searches Docker documentation and web resources for solutions
- Writes and modifies Dockerfiles following best practices
- Debugs container failures by reading logs and proposing fixes
- Manages containers, images, volumes, and networks

Gordon proposes every action before executing. You approve what it does.

## Where to use Gordon

Gordon is available on four surfaces:

- Open the Gordon view from the Docker Desktop sidebar to run Docker commands
  with your approval. See [Using Gordon in Docker
  Desktop](./how-to/docker-desktop.md).
- Run `docker ai` in the terminal to use the full assistant from the command
  line. See [Using Gordon via CLI](./how-to/cli.md).
- Select the Gordon icon on any repository page at
  [hub.docker.com](https://hub.docker.com) to ask about a repository's
  images, tags, and metadata. Hand off to Docker Desktop to take action.
- Select the Gordon icon on any page at
  [docs.docker.com](https://docs.docker.com) to ask Docker questions.

Docker Desktop and the CLI count against your Gordon plan's [usage
limits](./usage-limits.md). Gordon on Docker Hub and docs.docker.com is free
and does not require a Docker account or a Docker Desktop install. It has
its own shared public usage limit and does not access your Docker
environment.

## Get started

### Prerequisites

Before you begin:

- Docker Desktop 4.74 or later
- Sign in to your Docker account

> [!NOTE]
> Gordon is enabled by default for signed-in Docker users. If your account
> belongs to an organization with a Business subscription, access requires two
> additional steps:
>
> 1. Contact Docker Support to activate Gordon for your organization. Docker
>    will confirm when activation is complete.
> 2. Once confirmed, an organization administrator must turn on Gordon via
>    [Settings Management](/manuals/enterprise/security/hardened-desktop/settings-management/_index.md).
>    Set **Enable Gordon** to **Enabled** or **Always enabled**. Ensure all
>    Settings Management prerequisites are met for the setting to take effect
>    on Docker Desktop clients.

### Quick start

{{< tabs >}}
{{< tab name="Docker Desktop" >}}

1. Open Docker Desktop.
2. Select **Gordon** in the sidebar.
3. Select your project directory.
4. Type a question: "What containers are running?"

   ![Gordon running in Docker Desktop](./images/gordon_gui.avif)

5. Review Gordon's proposed actions and approve.

{{< /tab >}}
{{< tab name="CLI" >}}

1. Open your terminal and run:

   ```console
   $ docker ai
   ```

   This opens the Terminal User Interface (TUI) for Gordon.

2. Type a question: "what containers are running?" and press <kbd>Enter</kbd>.

   ![Gordon running in the terminal](./images/gordon_tui.avif)

3. Review Gordon's proposed actions and approve by typing `y`.

{{< /tab >}}
{{< /tabs >}}

### Permissions

By default, Gordon asks for approval before executing actions. You can approve
individual actions or allow all actions for the current session.

![Gordon permission request](./images/gordon_permissions_prompt.avif)

Permissions reset for each session. To configure default permissions or enable
auto-approve mode, see [Permissions](./how-to/permissions.md).

### Try these examples

Container inspection:

```console
$ docker ai "show me logs from my nginx container"
```

Dockerfile review:

```console
$ docker ai "review my Dockerfile for best practices"
```

Image management:

```console
$ docker ai "list my local images and their sizes"
```
