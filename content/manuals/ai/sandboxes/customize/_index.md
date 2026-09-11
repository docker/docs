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

## Builds and runtime configuration

Kits v3 include building as part of authoring a kit. Use the kit's Dockerfile
to install packages, compile tools, and copy static content into an image.
Use its descriptor to declare what the runtime must provide when the sandbox
runs, such as a credential, network access, or a startup hook.

A workload kit can build on a sandbox template with Dockerfile `FROM`. The
template supplies the base environment, and the kit adds its content and
runtime declarations. A mixin can also build and ship its own tools as an
overlay, so adding a tool doesn't require rebuilding the workload image.

Lifecycle install hooks remain available for initialization that needs a
sandbox's runtime inputs. Build steps produce content that sandboxes reuse;
install hooks initialize each sandbox separately.

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
