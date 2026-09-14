---
title: Customizing sandboxes
linkTitle: Customize
description: Build and share sandbox environments with v3 workload and mixin kits, or customize and save sandbox template images.
keywords: sandboxes, sbx, customize, templates, kits, mixins, workloads, custom agents
weight: 90
aliases:
  - /ai/sandboxes/agents/custom-environments/
params:
  sidebar:
    badge:
      color: blue
      text: Early Access
---

{{< summary-bar feature_name="Docker Sandboxes sbx" >}}

Use kits and templates to package the tools, configuration, and runtime
behavior your sandboxes need.

- [Kits](kits.md) combine reusable image content with declarations for network
  access, credentials, lifecycle hooks, and agent instructions. A workload kit
  defines what the sandbox runs. Mixin kits extend it.
- [Templates](templates.md) are reusable sandbox images. Extend a base image
  with a Dockerfile, or save a configured running sandbox as a template.

Kits are experimental. The format and CLI commands are subject to change.
Share feedback in the
[docker/sbx-releases](https://github.com/docker/sbx-releases) repository.

## Choose a customization

| Goal | Option |
| --- | --- |
| Run a v3 workload and add mixins | [Use kits](kits.md#use-kits) |
| Define an agent or another sandbox workload | [Build an agent](build-an-agent.md) |
| Add tools, configuration, or instructions to a v3 workload | [Mixin examples](kit-examples.md) |
| Customize the image used by a built-in agent | [Template](templates.md#build-a-custom-template) |
| Capture a configured running sandbox for reuse | [Saved template](templates.md#saving-a-sandbox-as-a-template) |

## Start with v3 kits

V3 is the recommended format for kit development. To use it, select an
explicit v3 workload reference and combine it with v3 mixins. Built-in agent
shortcuts such as `claude` and `codex` select v2 kits, so they can't be used
with v3 mixins. See [Run a kit](kits.md#run-a-kit) for the opt-in workflow.

V2 is deprecated but remains supported. For existing customizations, see
[Kits v2](kits-v2/_index.md), including guidance for moving a complete
environment to v3.
