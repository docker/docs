---
title: Customize sandbox environments
linkTitle: Customize environments
description: Apply kits and environment files to create repeatable sandbox environments with sbx.
keywords: docker sandboxes, sbx, local sandboxes, cloud sandboxes
weight: 40
---

Use a kit to supply an agent or add tools. Use an environment file to record
how a project creates its sandbox.

| Task | Guide |
| --- | --- |
| Run a published kit or add mixins | [Use kits](/manuals/ai/sandboxes/cli/customize/kits.md) |
| Save reusable project configuration in `sbxenv.yaml` | [Environment files](/manuals/ai/sandboxes/cli/customize/environment-files.md) |

Before reusing a kit through another interface, check
[Kit compatibility](/manuals/ai/sandboxes/author-kits/compatibility.md).
To build and publish a kit, use
[Author kits](/manuals/ai/sandboxes/author-kits/_index.md).

Host settings, including GPU passthrough and registry mirrors, belong under
[Local configuration](/manuals/ai/sandboxes/cli/local/_index.md).
To capture an environment configured interactively, see
[Save a local template](/manuals/ai/sandboxes/cli/manage/local.md#saving-a-sandbox-as-a-template)
or [save a cloud template](/manuals/ai/sandboxes/cli/manage/cloud.md#customize-a-cloud-sandbox).
