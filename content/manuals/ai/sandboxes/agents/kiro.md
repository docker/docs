---
title: Kiro
weight: 80
description: |
  Use Kiro in Docker Sandboxes with device flow authentication for interactive
  AI-assisted development.
keywords: docker sandboxes, kiro, ai agent, authentication, sbx
---

{{% include "sandboxes-local-scope.md" %}}

This guide covers authentication, configuration, and usage of Kiro in a
sandboxed environment.

Official documentation: [Kiro CLI](https://kiro.dev/docs/cli/)

## Quick start

Create a sandbox and run Kiro for a project directory:

```console
$ sbx run kiro ~/my-project
```

The workspace parameter is optional and defaults to the current directory:

```console
$ cd ~/my-project
$ sbx run kiro
```

On first run, Kiro prompts you to authenticate using device flow.

## Authentication

Kiro uses device flow authentication, which requires interactive login through
a web browser. This method provides secure authentication without storing API
keys directly.

### Device flow login

When you first run Kiro, it prompts you to authenticate:

1. Kiro displays a URL and a verification code
2. Open the URL in your web browser
3. Enter the verification code
4. Complete the authentication flow in your browser
5. Return to the terminal - Kiro proceeds automatically

The authentication session is persisted in the sandbox and doesn't require
repeated login unless you destroy and recreate the sandbox.

### Manual login

To trigger the login flow manually in a running sandbox, use `sbx exec` to
run Kiro directly:

```console
$ sbx exec -it <sandbox-name> kiro login --use-device-flow
```

This command initiates device flow authentication without starting a coding
session.

### Authentication persistence

Kiro stores authentication state in `~/.local/share/kiro-cli/data.sqlite3`
inside the sandbox. This database persists as long as the sandbox exists. If
you destroy the sandbox, you'll need to authenticate again when you recreate
it.

## Configuration

Sandboxes don't pick up user-level configuration from your host. Only
project-level configuration in the working directory is available inside the
sandbox. See
[Why doesn't the sandbox use my user-level agent configuration?](../faq.md#why-doesnt-the-sandbox-use-my-user-level-agent-configuration)
for workarounds.

Kiro requires minimal configuration. The agent runs with trust-all-tools mode
by default, which lets it execute commands without repeated approval prompts.

### Default startup command

Without extra args, the sandbox runs:

```text
kiro chat --trust-all-tools
```

The Kiro image stores `chat --trust-all-tools` in its `ENTRYPOINT`, so
forwarded arguments keep that prefix. Pass a prompt directly without repeating
`chat`:

```console
$ sbx run kiro -- "explain this project"
```

This runs `kiro chat --trust-all-tools "explain this project"`. You can also
pass flags, such as `sbx run kiro -- --resume`.

To run a different Kiro subcommand, use `sbx exec` in a running sandbox, as
shown in [Manual login](#manual-login).

## Base image

Template: `docker/sandbox-templates:kiro`

Authentication state is persisted across sandbox restarts.

See [Customize](../customize/) to pre-install tools or customize this
environment.
