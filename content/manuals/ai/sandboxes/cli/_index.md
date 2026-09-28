---
title: Use the CLI
description: Run and manage local and cloud sandboxes with the sbx CLI.
keywords: docker sandboxes, sbx, cloud sandboxes
weight: 40
---

Use `sbx` to work with local sandboxes or cloud sandboxes from your terminal.
For supported cloud operations, add the global `--cloud` flag. Host-only
settings remain local.

## Get started

[Install and sign in](/manuals/ai/sandboxes/cli/install.md), then follow the
[local quickstart](/manuals/ai/sandboxes/cli/get-started-local.md) or the
[cloud quickstart](/manuals/ai/sandboxes/cli/get-started-cloud.md).

## Manage your work

| Task | Guide |
| --- | --- |
| Create, stop, and remove sandboxes | [Manage sandboxes](/manuals/ai/sandboxes/cli/manage.md) |
| Run commands, use workspaces, and save templates locally | [Local operations](/manuals/ai/sandboxes/cli/usage.md) |
| Transfer files, expose services, and manage cloud expiration | [Cloud operations](/manuals/ai/sandboxes/cli/cloud-usage.md) |
| Configure credentials | [Local credentials](/manuals/ai/sandboxes/cli/credentials.md) and [cloud credentials](/manuals/ai/sandboxes/cli/credentials-cloud.md) |
| Configure network access | [Local policies](/manuals/ai/sandboxes/cli/network-local.md) and [cloud policies](/manuals/ai/sandboxes/cli/network-cloud.md) |
| Connect tools | [MCP servers](/manuals/ai/sandboxes/cli/mcp.md) |
| Apply reusable environments | [Kits and mixins](/manuals/ai/sandboxes/cli/kits.md) and [environment files](/manuals/ai/sandboxes/cli/environment-files.md) |
| Transfer sandbox state | [Move between local and cloud](/manuals/ai/sandboxes/cli/move.md) |

## Configure your workflow

Use [Git workflows](/manuals/ai/sandboxes/cli/git.md), [local development](/manuals/ai/sandboxes/cli/development.md),
[authenticated command-line tools](/manuals/ai/sandboxes/cli/authentication.md), or
[automation](/manuals/ai/sandboxes/cli/automation.md) for complete workflows. Connect an editor through
[integrations](/manuals/ai/sandboxes/cli/integrations/_index.md).

[Local configuration](/manuals/ai/sandboxes/cli/local/_index.md) covers model endpoints, agent skills,
GPU passthrough, upstream proxies, and registry mirrors.

For unsupported options and host dependencies, see
[Local and cloud differences](/manuals/ai/sandboxes/cli/local-vs-cloud.md). To diagnose a problem, use
[Troubleshooting](/manuals/ai/sandboxes/cli/troubleshooting.md).
