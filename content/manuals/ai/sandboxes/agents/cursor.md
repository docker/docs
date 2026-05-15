---
title: Cursor
weight: 33
description: |
  Use Cursor in Docker Sandboxes with API key or proxy-managed OAuth
  authentication.
keywords: docker sandboxes, cursor, cursor agent, ai agent, sbx
---

{{< summary-bar feature_name="Docker Sandboxes sbx" >}}

This guide covers authentication, configuration, and usage of Cursor in a
sandboxed environment.

Official documentation: [Cursor CLI](https://cursor.com/cli)

## Quick start

Create a sandbox and run Cursor for a project directory:

```console
$ sbx run cursor ~/my-project
```

The workspace parameter is optional and defaults to the current directory:

```console
$ cd ~/my-project
$ sbx run cursor
```

## Authentication

Cursor supports two authentication methods: an API key or OAuth.

**API key**: Store your Cursor API key using
[stored secrets](../security/credentials.md#stored-secrets):

```console
$ sbx secret set -g cursor
```

Alternatively, export the `CURSOR_API_KEY` environment variable in your shell
before running the sandbox. See
[Credentials](../security/credentials.md) for details on both methods.

**OAuth**: If no API key is set, Cursor prompts you to sign in interactively
on first run. The proxy intercepts the token exchange with
`api2.cursor.sh/auth/poll`, so credentials are managed by the host and aren't
stored inside the sandbox.

## Configuration

Sandboxes don't pick up user-level configuration from your host, such as
`~/.cursor`. Only project-level configuration in the working directory is
available inside the sandbox. See
[Why doesn't the sandbox use my user-level agent configuration?](../faq.md#why-doesnt-the-sandbox-use-my-user-level-agent-configuration)
for workarounds.

Cursor reads `AGENTS.md` from the workspace for agent-specific instructions.

The sandbox runs Cursor in YOLO mode by default, which executes commands
without approval prompts. Pass additional `cursor-agent` CLI options after
`--`:

```console
$ sbx run cursor --name <sandbox-name> -- <cursor-options>
```

## Base image

Template: `docker/sandbox-templates:cursor-agent-docker`

Preconfigured to run in YOLO mode with HTTP/1.1 and server-sent events for
agent traffic so requests flow through the host proxy. Authentication state
is persisted across sandbox restarts.

See [Customize](../customize/) to pre-install tools or customize this
environment.
