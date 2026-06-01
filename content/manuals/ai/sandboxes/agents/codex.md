---
title: Codex
weight: 20
description: |
  Use OpenAI Codex in Docker Sandboxes with API key authentication and YOLO
  mode configuration.
keywords: docker sandboxes, codex, openai, ai agent, sbx
---

This guide covers authentication, configuration, and usage of Codex in a
sandboxed environment.

Official documentation: [Codex CLI](https://developers.openai.com/codex/cli)

## Quick start

Create a sandbox and run Codex for a project directory:

```console
$ sbx run codex ~/my-project
```

The workspace parameter is optional and defaults to the current directory:

```console
$ cd ~/my-project
$ sbx run codex
```

## Authentication

Codex supports two authentication methods: an API key or OAuth.

**API key**: Store your OpenAI API key using
[stored secrets](../security/credentials.md#stored-secrets):

```console
$ sbx secret set -g openai
```

Alternatively, export the `OPENAI_API_KEY` environment variable in your shell
before running the sandbox.

**OAuth**: If you prefer not to use an API key, start the OAuth flow on your
host with:

```console
$ sbx secret set -g openai --oauth
```

This opens a browser window for authentication and stores the resulting tokens
in your OS keychain. The OAuth flow runs on the host, not inside the sandbox,
so browser-based authentication works without any extra setup.

See [Credentials](../security/credentials.md) for more details.

## Configuration

Sandboxes don't pick up user-level configuration from your host, such as
`~/.codex`. Only project-level configuration in the working directory is
available inside the sandbox. See
[Why doesn't the sandbox use my user-level agent configuration?](../faq.md#why-doesnt-the-sandbox-use-my-user-level-agent-configuration)
for workarounds.

### Default startup command

Without extra args, the sandbox runs:

```text
codex --dangerously-bypass-approvals-and-sandbox
```

Args after `--` replace these defaults rather than being appended. To keep
the flag, include it yourself:

```console
$ sbx run codex -- --dangerously-bypass-approvals-and-sandbox "fix the build"
```

## Base image

Template: `docker/sandbox-templates:codex`

See [Customize](../customize/) to pre-install tools or customize this
environment.
