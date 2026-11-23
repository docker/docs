---
title: API v2 recipes
description: Automate sandbox management, compose environments, and connect tools through API v2.
keywords: docker sandboxes, API, SDK, local sandboxes, cloud sandboxes
weight: 40
linkTitle: Recipes
---

Use these recipes to automate tasks you perform with Docker Sandboxes. The
recipes are outlines while the API interface is being finalized. Examples and
exact request fields are deferred.

- [Customize an environment](customize-an-environment.md): combine a workload
  with mixins using v3 kit references.
- [Manage existing sandboxes](manage-sandboxes.md): find environments and manage
  their lifecycle.
- [Work with project files](project-files.md): distinguish local workspace
  mounts from cloud file transfers.
- [Connect agents to services](connect-agents-to-services.md): provide
  credentials declared by kits.
- [Connect MCP tools](mcp.md): configure a managed cloud MCP gateway.

Each recipe states which execution environment it covers. See
[Local and cloud capabilities](../capabilities.md) before adapting a workflow
between environments.
