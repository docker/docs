---
title: Choose an environment and interface
linkTitle: Choose how to work
description: Choose how to run and manage Docker Sandboxes based on your workflow and access requirements.
keywords: docker sandboxes, sbx, cloud sandboxes
weight: 10
---

Choose where your sandbox runs, then choose how to interact with it. The CLI
supports local and cloud sandboxes. The Console, API, and SDK use cloud
sandboxes.

| Interface | Environment | Use it for |
| --- | --- | --- |
| `sbx` CLI | Local or cloud | Interactive terminal work and command-line automation |
| Docker Agentic Platform Console | Cloud | Launching agents and using terminals, files, and services in a browser |
| REST API or TypeScript SDK | Cloud | Integrating sandbox operations into an application |

## Choose an environment

Local sandboxes can mount host workspaces and use supported host hardware.
Cloud sandboxes run independently of your computer. Copy files or clone a
repository into a cloud sandbox instead of mounting a local directory.

| Concern | Local | Cloud |
| --- | --- | --- |
| Credentials | Stored on the host | Stored in the cloud; CLI and Console share the store |
| Network controls | Host and sandbox rules, with organization governance where configured | Account and sandbox policies |
| Lifetime | Managed through the local daemon | Includes expiration and supported stop/resume behavior |
| Host configuration | Upstream proxy, registry mirror, and supported GPU passthrough | Host settings do not configure cloud infrastructure |

Read [files and storage](/manuals/ai/sandboxes/concepts/files.md),
[secrets](/manuals/ai/sandboxes/concepts/secrets.md), and [lifecycle](/manuals/ai/sandboxes/concepts/lifecycle.md) before
adapting a workflow that depends on these differences. The
[CLI capability comparison](/manuals/ai/sandboxes/cli/get-started/local-vs-cloud.md) lists command-specific limits.

## Choose an interface

Use [CLI tasks](/manuals/ai/sandboxes/cli/_index.md) for terminal workflows. Local and cloud examples
share a page when the steps fit together. Some tasks, such as configuring an
upstream proxy, apply only to local sandboxes.

Use [the Console](/manuals/ai/sandboxes/console/_index.md) to launch and manage sandboxes in your
browser. Use [the API and SDK](/manuals/ai/sandboxes/api/_index.md) for application workflows that
manage resources, wait for operations, and handle failures programmatically.

A shared cloud service does not make every option interchangeable. Follow the
selected interface's instructions for attaching secrets, configuring MCP, and
launching kits. See [kit compatibility](/manuals/ai/sandboxes/author-kits/compatibility.md) before
reusing a custom kit through another interface.
