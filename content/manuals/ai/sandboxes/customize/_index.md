---
title: Kits
description: Understand sandbox workload and mixin kits, then choose how to use existing kits or author your own.
keywords: sandboxes, sbx, kits, v3, workloads, mixins
weight: 90
linkTitle: Kits
aliases:
  - /ai/sandboxes/agents/custom-environments/
  - /ai/sandboxes/customize/kits/
---

{{< summary-bar feature_name="Docker Sandboxes sbx" >}}

Kits let you shape a sandbox around the way you work. Choose its base image,
add the tools your project needs, and give your agent instructions for using
them. You also control which services the sandbox can access, how it
authenticates, and what runs when the sandbox starts.

A kit can define the whole environment or add something to an existing one,
such as a toolchain or your team's shared configuration. Package those choices
once, then reuse them across projects and share them with your team.

This page covers v3, the recommended format for new kit development. V3 brings
image builds and runtime capabilities together in the kit format.

> [!NOTE]
> V3 kits are experimental. The format and CLI commands are subject to change.
> Share spec feedback or report issues in
> [docker/sandbox-kit-spec](https://github.com/docker/sandbox-kit-spec/issues).

The [sandbox-kit-spec repository](https://github.com/docker/sandbox-kit-spec)
contains the authoritative v3 specification, build frontend, and examples.

## Workloads and mixins

A sandbox runs one workload kit. You can add mixin kits to customize that
workload:

| Kind | What it supplies | How you use it |
| --- | --- | --- |
| `workload` | The environment and command to run, such as an agent or a shell | Pass it to `sbx run` or `sbx create` |
| `mixin` | Additional tools, configuration, or runtime behavior | Add it with `--kit` |

For example, a team might use an OpenCode workload with a mixin that adds a
linter and another that supplies the team's review instructions. Each kit can
be maintained and shared separately.

## Version compatibility

V3 kits cannot be combined with v1 or v2 kits in the same sandbox. To use
v3, select a v3 workload and use v3 for every mixin you add.

The built-in agent names, such as `claude` and `codex`, select v2 kits.
You can't add a v3 mixin to these built-ins. For example,
`sbx run claude --kit ./some-v3-mixin` fails because it mixes kit versions.
Instead, select a v3 workload by its published image, local path, or Git
reference, as shown in [Run a kit](/manuals/ai/sandboxes/customize/use-kits.md#run-a-kit).

V2 remains supported, including the built-in agents and
existing v2 customizations. See [Kits v2](/manuals/ai/sandboxes/customize/kits-v2/_index.md) for maintenance
and migration guidance.

## What kits can do

Use kits to give agents a repeatable working environment and share it with
your team:

- Package a custom agent, or configure an existing agent for your team's
  projects.
- Include the tools the agent needs, such as linters, language runtimes,
  test runners, and compilers.
- Share linter rules, editor settings, helper scripts, and reference material.
  Give the agent instructions and skills for using them.
- Connect the agent to services through network rules and credentials,
  including internal APIs and private package registries.
- Initialize each sandbox and run supporting services when it starts, such
  as a development server for previewing the agent's work.

## Choose your next step

- [Use kits](/manuals/ai/sandboxes/customize/use-kits.md) to run a workload, add mixins, and configure a
  sandbox from kits someone else has published.
- [Author kits](/manuals/ai/sandboxes/customize/author/_index.md) to define a workload or mixin, build its
  content, and declare its runtime requirements.

For the earlier format used by built-in agents, see [Kits v2](/manuals/ai/sandboxes/customize/kits-v2/_index.md).
To capture an interactively configured sandbox's container filesystem for reuse,
see [Save a sandbox as a template](/manuals/ai/sandboxes/usage.md#saving-a-sandbox-as-a-template).
