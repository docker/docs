---
title: Gemini
weight: 40
description: |
  Use Google Gemini in Docker Sandboxes with proxy-managed authentication and
  API key configuration.
keywords: docker sandboxes, gemini, google, ai agent, sbx
---

{{< summary-bar feature_name="Docker Sandboxes sbx" >}}

This guide covers authentication, configuration, and usage of Google Gemini in
a sandboxed environment.

Official documentation: [Gemini CLI](https://geminicli.com/docs/)

## Quick start

Create a sandbox and run Gemini for a project directory:

```console
$ sbx run gemini ~/my-project
```

The workspace parameter is optional and defaults to the current directory:

```console
$ cd ~/my-project
$ sbx run gemini
```

## Authentication

Gemini requires either a Google API key or a Google account with Gemini access.

**API key**: Store your key using
[stored secrets](../security/credentials.md#stored-secrets):

```console
$ sbx secret set -g google
```

Alternatively, export the `GEMINI_API_KEY` or `GOOGLE_API_KEY` environment
variable in your shell before running the sandbox. See
[Credentials](../security/credentials.md) for details on both methods.

**Google account**: If no API key is set, Gemini prompts you to sign in
interactively when it starts. Interactive authentication is scoped to the
sandbox and doesn't persist if you remove and recreate it.

## Configuration

Sandboxes don't pick up user-level configuration from your host, such as
`~/.gemini`. Only project-level configuration in the working directory is
available inside the sandbox. See
[Why doesn't the sandbox use my user-level agent configuration?](../faq.md#why-doesnt-the-sandbox-use-my-user-level-agent-configuration)
for workarounds.

The sandbox runs Gemini without approval prompts by default and disables
Gemini's built-in sandbox tool (since the sandbox itself provides isolation).
Pass additional Gemini CLI options after `--`:

```console
$ sbx run gemini --name <sandbox-name> -- <gemini-options>
```

## Base image

Template: `docker/sandbox-templates:gemini`

Gemini is configured to disable its built-in OAuth flow. Authentication is
managed through the proxy with API keys. Preconfigured to run without
approval prompts.

See [Customize](../customize/) to pre-install tools or customize this
environment.
