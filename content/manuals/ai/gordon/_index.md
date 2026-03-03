---
title: Gordon
description: AI agent for Docker workflows - execute tasks, debug issues, and
  manage containers with intelligent assistance
weight: 40
params:
  sidebar:
    badge:
      color: blue
      text: Beta
    group: AI
aliases:
  - /desktop/features/gordon/
---

{{< summary-bar feature_name="Gordon" >}}

Gordon is an AI agent that takes action on your Docker workflows. It analyzes
your environment, proposes solutions, and executes commands with your
permission. Available in Docker Desktop and via the `docker ai` CLI command.

## What Gordon does

Gordon takes action to help you with Docker tasks:

- Explains Docker concepts and commands
- Searches Docker documentation and web resources for solutions
- Writes and modifies Dockerfiles following best practices
- Debugs container failures by reading logs and proposing fixes
- Manages containers, images, volumes, and networks

Gordon proposes every action before executing. You approve what it does.

## Get started

### Prerequisites

Before you begin:

- Docker Desktop 4.61.0 or later
- Sign in to your Docker account

> [!NOTE]
> Gordon is enabled by default for Personal, Pro, and Team subscriptions. For
> Business subscriptions, an administrator must enable Gordon for the
> organization before users can access it.

### Quick start

{{< tabs >}}
{{< tab name="Docker Desktop" >}}

1. Open Docker Desktop.
2. Select **Ask Gordon** in the sidebar.
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

   ![Gordon running in the terminal](./images/gordon_tui.avif?border=true)

3. Review Gordon's proposed actions and approve by typing `y`.

{{< /tab >}}
{{< /tabs >}}

### Permissions

By default, Gordon asks for approval before executing actions. You can approve
individual actions or allow all actions for the current session.

![Gordon permission request](./images/permissions.avif)

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

## Usage and availability

Gordon is available with all Docker subscriptions. Usage limits vary by tier:

- Personal: Baseline usage
- Pro and Team: 3x more usage than Personal
- Business: 6x more usage than Personal

For details, see [Usage and limits](./usage-and-limits/).
