---
title: Customizing sandboxes
linkTitle: Customize
description: Build and share sandbox environments with v3 workload and mixin kits, using Docker-provided base images or your own Linux image.
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

Use kits to package the tools, configuration, and runtime behavior your
sandboxes need. A workload kit defines the environment and command to run.
Mixin kits extend it with tools, configuration, or instructions.

A workload's build recipe can start from a Docker-provided agent image or
another Linux base image. Its descriptor declares network access, credentials,
and lifecycle hooks alongside that build.

Kits are experimental. The format and CLI commands are subject to change.
Share feedback in the
[docker/sbx-releases](https://github.com/docker/sbx-releases) repository.

## Choose a customization

| Goal | Option |
| --- | --- |
| Run a v3 workload and add mixins | [Use kits](kits.md#use-kits) |
| Define an agent or another sandbox workload | [Build an agent](build-an-agent.md) |
| Add tools, configuration, or instructions to a v3 workload | [Mixin examples](kit-examples.md) |
| Choose an existing agent image or your own Linux base | [Base images](base-images.md) |

## Start with v3 kits

V3 is the recommended format for kit development. To use it, select an
explicit v3 workload reference and combine it with v3 mixins. Built-in agent
shortcuts such as `claude` and `codex` select v2 kits, so they can't be used
with v3 mixins. See [Run a kit](kits.md#run-a-kit) for the opt-in workflow.

V2 is deprecated but remains supported. For existing customizations, see
[Kits v2](kits-v2/_index.md), including guidance for moving a complete
environment to v3.

To capture an interactively configured sandbox's container filesystem for reuse,
see [Save a sandbox as a template](/manuals/ai/sandboxes/usage.md#saving-a-sandbox-as-a-template).
