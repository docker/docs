---
title: Customizing sandboxes
linkTitle: Customize
description: Build reusable sandbox workloads and extensions with kits, or customize and save sandbox template images.
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

- [Kits](kits.md) combine reusable image content with runtime declarations for
  credentials, network access, storage, lifecycle hooks, and agent instructions.
  A workload kit defines what the sandbox runs. Mixin kits extend it.
- [Templates](templates.md) are reusable sandbox images. Extend a base image
  with a Dockerfile, or save a configured running sandbox as a template.

Kits are experimental. The format and CLI commands are subject to change.
Share feedback in the
[docker/sbx-releases](https://github.com/docker/sbx-releases) repository.

## Choose a customization

| Goal | Option |
| --- | --- |
| Define an agent or another sandbox workload | [Workload kit](build-an-agent.md) |
| Add a tool, configuration, or runtime capability to a v3 workload | [Mixin kit](kit-examples.md) |
| Customize the image used by a built-in agent | [Template](templates.md#build-a-custom-template) |
| Capture a configured running sandbox for reuse | [Saved template](templates.md#saving-a-sandbox-as-a-template) |
| Maintain an existing v2 sandbox or mixin kit | [Kits v2](kits-v2/_index.md) |

Use v3 for authoring kits. Docker Sandboxes also supports v1 and v2, whose
`spec.yaml` workflow is documented in [Kits v2](kits-v2/_index.md). A single
composition cannot combine v3 kits with v1 or v2 kits.

## Builds and runtime configuration

A kit starts as source files: a YAML descriptor for sandbox behavior and a
Dockerfile for software and files to package. Building them produces an image
you can share. Docker Sandboxes reads the descriptor when it runs the kit.

For example, an agent kit can install the agent during its build, then request
access to its API and credentials when the sandbox runs. The installed agent
is reused across sandboxes, while each sandbox receives its own configuration.

A workload kit can use a sandbox template as its base image. Mixins add tools
or behavior to that workload. See [Kits](kits.md) for the file layout, a
complete example, and how the parts work together.
