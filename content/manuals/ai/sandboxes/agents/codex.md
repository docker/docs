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

## Choose where to run

Use the local setup on this page for host workspaces and local authentication.
For a cloud CLI launch, use the cloud setup here. Generic files, ports, and
lifecycle tasks are covered in [Use the CLI](/manuals/ai/sandboxes/cli/_index.md).

## Cloud CLI setup

[Activate cloud access](/manuals/ai/sandboxes/cloud-access/_index.md) and sign in with `sbx login`.
Store the provider credential in the cloud secret store, then launch the agent:

```console
$ sbx --cloud secret set openai
$ sbx --cloud run codex --name cloud-project
```

For account-scoped OpenAI OAuth instead of an API key, use
`sbx --cloud secret set openai --oauth` before launching the sandbox.

Cloud sandboxes do not mount your host project directory. Copy files or clone a
repository into the sandbox. See [Cloud operations](/manuals/ai/sandboxes/cli/cloud-usage.md) and
[Cloud credentials](/manuals/ai/sandboxes/cli/credentials-cloud.md) for the supported workflow.

The setup and host configuration instructions that follow apply to local
sandboxes.

## Quick start

Create a sandbox and run Codex for a project directory:

```console
$ sbx run codex ~/my-project
```

`sbx run` defaults the workspace to the current directory:

```console
$ cd ~/my-project
$ sbx run codex
```

To create a [mountless sandbox](/manuals/ai/sandboxes/cli/usage.md#choose-a-workspace), use
`sbx create` without a workspace path, then attach by name.

## Authentication

For the default OpenAI models, `sbx run codex` prompts you to authenticate on
your host if you haven't stored an OpenAI credential. Authentication happens
before launching the sandbox, so credentials are never exposed inside it.

To set up authentication ahead of time, choose one of the following methods.

**OAuth**: Start the OAuth flow on your host with:

```console
$ sbx secret set openai --oauth
```

This opens a browser window for authentication and stores the resulting tokens
in your OS keychain. The OAuth flow runs on the host, not inside the sandbox,
so browser-based authentication works without any extra setup.

**API key**: Store your OpenAI API key using
[stored secrets](/manuals/ai/sandboxes/cli/credentials.md#stored-secrets):

```console
$ sbx secret set openai
```

See [Credentials](/manuals/ai/sandboxes/cli/credentials.md) for more details.

## Model selection

To use Codex with a local model or another inference provider, see
[Use local and hosted models](/manuals/ai/sandboxes/cli/local/models.md).

## Configuration

Sandboxes don't pick up user-level configuration from your host, such as
`~/.codex`. Only project-level configuration in the working directory is
available inside the sandbox. See
[Why doesn't the sandbox use my user-level agent configuration?](/manuals/ai/sandboxes/faq.md#why-doesnt-the-sandbox-use-my-user-level-agent-configuration)
for workarounds.

### Default startup command

Without extra args, the sandbox runs:

```text
codex --dangerously-bypass-approvals-and-sandbox
```

Arguments after `--` are added after the default flags when the first one is
itself a flag (begins with `-`). A bare word — such as a prompt — replaces the
defaults instead, so lead with the flag to keep bypass mode:

```console
$ sbx run --name <sandbox-name> -- --dangerously-bypass-approvals-and-sandbox "fix the build"
```

## Base image

Template: `docker/sandbox-templates:codex`

See [Customize](/manuals/ai/sandboxes/concepts/kits.md) to pre-install tools or customize this
environment.
