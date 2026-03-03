---
title: Gemini sandbox
description: |
  Use Google Gemini in Docker Sandboxes with proxy-managed authentication and
  API key configuration.
keywords: docker, sandboxes, gemini, google, ai agent, authentication, proxy
weight: 40
---

{{< summary-bar feature_name="Docker Sandboxes" >}}

This guide covers authentication, configuration, and usage of Google Gemini in
a sandboxed environment.

Official documentation: [Gemini CLI](https://geminicli.com/docs/)

## Quick start

Create a sandbox and run Gemini for a project directory:

```console
$ docker sandbox run gemini ~/my-project
```

The workspace parameter is optional and defaults to the current directory:

```console
$ cd ~/my-project
$ docker sandbox run gemini
```

## Authentication

Gemini uses proxy-managed authentication. Docker Sandboxes intercepts API
requests and injects credentials transparently. You provide your API key
through environment variables on the host, and the sandbox handles credential
management.

### Environment variable (recommended)

Set the `GEMINI_API_KEY` or `GOOGLE_API_KEY` environment variable in your
shell configuration file.

Docker Sandboxes use a daemon process that doesn't inherit environment
variables from your current shell session. To make your API key available to
sandboxes, set it globally in your shell configuration file.

Add the API key to your shell configuration file:

```plaintext {title="~/.bashrc or ~/.zshrc"}
export GEMINI_API_KEY=AIzaSyxxxxx
```

Or use `GOOGLE_API_KEY`:

```plaintext {title="~/.bashrc or ~/.zshrc"}
export GOOGLE_API_KEY=AIzaSyxxxxx
```

Apply the changes:

1. Source your shell configuration: `source ~/.bashrc` (or `~/.zshrc`)
2. Restart Docker Desktop so the daemon picks up the new environment variable
3. Create and run your sandbox:

```console
$ docker sandbox create gemini ~/project
$ docker sandbox run <sandbox-name>
```

The sandbox detects the environment variable and uses it automatically.

### Interactive authentication

If neither `GEMINI_API_KEY` nor `GOOGLE_API_KEY` is set, Gemini prompts you to
sign in when it starts.

When using interactive authentication:

- You must authenticate each sandbox separately
- If the sandbox is removed or destroyed, you'll need to authenticate again when you recreate it
- Authentication sessions aren't persisted outside the sandbox
- No fallback authentication methods are used

To avoid repeated authentication, set the `GEMINI_API_KEY` or `GOOGLE_API_KEY` environment variable.

## Configuration

Configure Gemini behavior in `~/.gemini/settings.json`:

```json
{
  "disable_sandbox_tool": true,
  "trusted_folders": ["/workspace"]
}
```

These settings disable safety checks and allow Gemini to operate without
repeated confirmations for workspace files.

### Pass options at runtime

Pass Gemini CLI options after the sandbox name and a `--` separator:

```console
$ docker sandbox run <sandbox-name> -- --yolo
```

The `--yolo` flag disables approval prompts for a single session without
modifying the configuration file.

## Base image

Template: `docker/sandbox-templates:gemini`

Gemini is configured to disable its built-in OAuth flow. Authentication is managed through the Docker Sandbox proxy with API keys.

See [Custom templates](../templates.md) to build your own agent images.
