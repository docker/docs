---
title: Customizing sandboxes
linkTitle: Customize
description: Build reusable sandbox images, extend agents with tools and credentials, and define custom agents using templates and kits.
keywords: sandboxes, sbx, customize, templates, kits, mixins, custom agents
weight: 35
aliases:
  - /ai/sandboxes/agents/custom-environments/
---

{{< summary-bar feature_name="Docker Sandboxes sbx" >}}

Docker Sandboxes offers two ways to customize a sandbox beyond the built-in
defaults:

- [Templates](templates.md) — reusable sandbox images with tools, packages,
  and configuration baked in. Extend a base image with a Dockerfile, or
  save a running sandbox as a template.
- [Kits](kits.md) — declarative YAML artifacts that extend an agent with
  tools, credentials, network rules, and files at runtime, or define a new
  agent from scratch.

Kits are experimental. The kit file format, CLI commands, and experience for
creating, loading, and managing kits are subject to change as the feature
evolves. Share feedback and bug reports in the
[docker/sbx-releases](https://github.com/docker/sbx-releases) repository.

## Templates and kits, side by side

A template is a Docker image that the sandbox runs. It's built ahead
of time with a Dockerfile (or saved from a running sandbox), pushed to a
registry, and pulled when a sandbox is created. Use templates for things
that belong in an image: system packages, language toolchains, large
dependencies — anything you'd rather not reinstall on every sandbox start.

A kit is a YAML artifact applied at sandbox creation. The kit can run
install commands, drop files into the sandbox, declare network and
credential rules, and (for agent kits) define which template image the
agent runs in. Use kits for things that vary per agent or per team:
shared linter config, project-specific install steps, credential
injection for a service the agent talks to.

Templates and kits work together. An agent kit's `agent.image` field
points at a template: the template provides the base environment, the
kit layers config, secrets, and runtime behavior on top. A team can ship
one heavy template and several thin kits without rebuilding the image
each time something changes.

## When to use which

| Goal                                                      | Option                                                        |
| --------------------------------------------------------- | ------------------------------------------------------------- |
| Pre-install tools and packages into a reusable base image | [Template](templates.md)                                      |
| Capture a configured running sandbox for reuse            | [Saved template](templates.md#saving-a-sandbox-as-a-template) |
| Add a tool, credential, or config to agent runs via YAML  | [Kit (mixin)](kits.md)                                        |
| Define a new agent from scratch                           | [Kit (agent)](kits.md#defining-an-agent)                      |

Templates and kits can be used together. A template bakes heavy tools into
the image for fast sandbox startup; a kit layered on top adds per-run
credentials, config, or extra capabilities.

## Tutorials

- [Build your own agent kit](build-an-agent.md) — step-by-step walkthrough
  for packaging [Amp](https://ampcode.com/) as an agent kit.
