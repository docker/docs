---
title: OpenCode
weight: 60
description: |
  Use OpenCode in Docker Sandboxes with multi-provider authentication and TUI
  interface for AI development.
keywords: docker sandboxes, opencode, ai agent, authentication, sbx
---

{{< summary-bar feature_name="Docker Sandboxes sbx" >}}

This guide covers authentication, configuration, and usage of OpenCode in a
sandboxed environment.

Official documentation: [OpenCode](https://opencode.ai/docs)

## Quick start

Create a sandbox and run OpenCode for a project directory:

```console
$ sbx run opencode ~/my-project
```

The workspace parameter is optional and defaults to the current directory:

```console
$ cd ~/my-project
$ sbx run opencode
```

OpenCode launches a TUI (text user interface) where you can select your
preferred LLM provider and interact with the agent.

## Authentication

OpenCode supports multiple providers. Store keys for the providers you want to
use with [stored secrets](../security/credentials.md#stored-secrets):

```console
$ sbx secret set -g openai
$ sbx secret set -g anthropic
$ sbx secret set -g google
$ sbx secret set -g xai
$ sbx secret set -g groq
$ sbx secret set -g aws
```

You only need to configure the providers you want to use. OpenCode detects
available credentials and offers those providers in the TUI.

You can also use environment variables (`OPENAI_API_KEY`, `ANTHROPIC_API_KEY`,
`GOOGLE_API_KEY`, `XAI_API_KEY`, `GROQ_API_KEY`, `AWS_ACCESS_KEY_ID`). See
[Credentials](../security/credentials.md) for details on both methods.

## Configuration

Sandboxes don't pick up user-level configuration from your host. Only
project-level configuration in the working directory is available inside the
sandbox. See
[Why doesn't the sandbox use my user-level agent configuration?](../faq.md#why-doesnt-the-sandbox-use-my-user-level-agent-configuration)
for workarounds.

OpenCode uses a TUI interface and doesn't require extensive configuration
files. The agent prompts you to select a provider when it starts, and you can
switch providers during a session.

### Pass options at runtime

Pass OpenCode CLI options after `--`:

```console
$ sbx run opencode --name <sandbox-name> -- <opencode-options>
```

For example, to resume an existing session in a named sandbox:

```console
$ sbx run <sandbox-name> -- -s <session-id>
```

### TUI mode

OpenCode launches in TUI mode by default. The interface shows:

- Available LLM providers (based on configured credentials)
- Current conversation history
- File operations and tool usage
- Real-time agent responses

Use keyboard shortcuts to navigate the interface and interact with the agent.

## Base image

Template: `docker/sandbox-templates:opencode`

OpenCode supports multiple LLM providers with automatic credential injection
through the sandbox proxy.

See [Customize](../customize/) to pre-install tools or customize this
environment.
