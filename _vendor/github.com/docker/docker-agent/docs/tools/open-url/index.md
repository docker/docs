---
title: "Open URL Tool"
description: "Open a fixed URL in the user's default browser."
keywords: docker agent, ai agents, tools, toolsets, open url tool
linkTitle: "Open URL"
weight: 40
canonical: https://docs.docker.com/ai/docker-agent/tools/open-url/
---

_Open a fixed URL in the user's default browser._

## Overview

The `open_url` toolset exposes a single, argument-less tool that opens a URL
baked into the toolset definition in the user's default browser. The model
never supplies the URL — it just calls the tool by name. Launching the browser
is cross-platform: docker-agent uses `open` on macOS, `xdg-open` on Linux, and
`rundll32` on Windows.

> [!NOTE]
> **When to Use**
>
> - Letting an agent open a dashboard, documentation page, or deep link on demand
> - Deep-linking into a desktop app via a custom URI scheme (e.g. `docker-desktop://`)
> - Any "take me there" action where the destination is fixed and known up front

## Configuration

```yaml
agents:
  assistant:
    model: openai/gpt-4o
    description: Assistant that can open the dashboard
    instruction: When the user asks to see the dashboard, call open_dashboard.
    toolsets:
      - type: open_url
        name: open_dashboard
        url: https://example.com/dashboard
```

## Properties

| Property | Type   | Required | Description                                                                                          |
| -------- | ------ | -------- | ---------------------------------------------------------------------------------------------------- |
| `url`    | string | ✓        | URL to open. Supports `${env.VAR}` interpolation. Any scheme the OS can dispatch is allowed.         |
| `name`   | string | ✗        | Tool name the agent references. Defaults to `open_url`. Use a descriptive name when configuring several. |

## Multiple URLs

Add one toolset entry per destination, each with its own `name`:

```yaml
toolsets:
  - type: open_url
    name: open_dashboard
    url: https://example.com/dashboard
  - type: open_url
    name: open_docs
    url: https://docs.example.com/${env.DOCS_VERSION}
```

## URL Interpolation

The `url` field supports `${env.VAR}` placeholders, expanded at call time
against the runtime environment:

```yaml
toolsets:
  - type: open_url
    name: open_docs
    url: https://docs.example.com/${env.DOCS_VERSION}
```

## Custom URI Schemes

Any scheme the operating system knows how to dispatch works, including deep
links into desktop applications:

```yaml
toolsets:
  - type: open_url
    name: open_in_docker_desktop
    url: docker-desktop://dashboard/apps
```

## Limitations

- The URL must include a scheme (e.g. `https://`); bare paths are rejected.
- URLs that look like a command-line flag (starting with `-`) are refused to
  prevent argument injection into the platform `open` helper.
- The tool opens the URL on the **host** running docker-agent; in headless or
  remote environments where no browser/launcher is available, the call fails
  gracefully and reports the error to the agent.

See [`examples/open_url.yaml`](https://github.com/docker/docker-agent/blob/main/examples/open_url.yaml) for a complete configuration.
