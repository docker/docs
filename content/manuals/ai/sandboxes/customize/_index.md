---
title: Kits
description: Use sandbox kits and share environments as kit sets composed from workloads and mixins.
keywords: sandboxes, sbx, kits, v3, workloads, mixins
weight: 90
linkTitle: Kits
aliases:
  - /ai/sandboxes/agents/custom-environments/
---

{{< summary-bar feature_name="Docker Sandboxes sbx" >}}

Kits let you shape a sandbox around the way you work. Use a custom base image,
add the tools your project needs, and give your agent instructions for using
them. You also control which services the sandbox can access, how it
authenticates, and what runs when the sandbox starts.

A kit can define the whole environment or add something to an existing one,
such as a toolchain or your team's shared configuration. Package those choices
once, then reuse them across projects and share them with your team.

This page covers v3 kits, which require `sbx` v0.45 or later. See
[Version compatibility](#version-compatibility) for compatibility with earlier
kit formats.

## What is a kit?

A kit packages software and configuration for a sandbox. A YAML file, called
its descriptor, tells Docker Sandboxes what the kit needs: network access,
credentials, setup commands, or instructions for the agent. You publish the
kit as a container image containing its files and descriptor.

You identify a published kit by its container image reference, such as
`me/my-kit:latest`. When you use the kit, `sbx` prepares its files and applies
its settings. You can use one kit for a complete environment or combine kits
that contribute different parts.

## What kits can do

Use kits to:

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

## Workloads and mixins

Kits have two roles in a sandbox: a workload supplies the base environment
and launch command, and mixins add tools or behavior to it.

| Kind | What it supplies | How you use it |
| --- | --- | --- |
| `workload` | The environment and command to run, such as an agent or a shell | Pass it to `sbx run` or `sbx create` |
| `mixin` | Additional tools, configuration, or runtime behavior | Add it with `--kit` |

For example, this command runs a workload with a mixin:

```console
$ sbx run me/my-agent-kit:latest --kit me/my-mixin:latest
```

The workload supplies the environment and launch command. The mixin adds its
tools and configuration to that environment.

## Kit sets

A kit set combines kits into a single kit that you can publish and reuse.
Use it to package a workload with the mixins you need, so you and your team
can run the environment from one reference.

The set specifies which kits and versions to combine. It can also add setup
commands, agent instructions, network access, and other settings for the
combined environment.

If the set includes a workload, you run the published kit with `sbx run`.
If it contains only mixins, you add it to a workload with `--kit`.

See [Use kits](/manuals/ai/sandboxes/customize/use-kits.md) for how to run kits
and add mixins. To customize and publish your own combination, see
[Compose a kit set](/manuals/ai/sandboxes/customize/author/kit-sets.md).

## Version compatibility

V3 kits cannot be combined with v1 or v2 kits in the same sandbox. To use
v3, select a v3 workload and use v3 for every mixin you add.

The built-in agent names, such as `claude` and `codex`, select v2 kits.
You can't add a v3 mixin to these built-ins. For example,
`sbx run claude --kit ./some-v3-mixin` fails because it mixes kit versions.
Instead, select a v3 workload by its published image, local path, or Git
reference, as shown in [Run a kit](/manuals/ai/sandboxes/customize/use-kits.md#run-a-kit).

V2 remains supported, including the built-in agents and
existing v2 customizations. See [Kits v2](/manuals/ai/sandboxes/customize/kits-v2.md) for maintenance
and migration guidance.

## Choose your next step

- [Use kits](/manuals/ai/sandboxes/customize/use-kits.md) to run a published environment or combine
  a workload with mixins.
- [Author kits](/manuals/ai/sandboxes/customize/author/_index.md) to build an agent environment,
  package a tool, or combine kits into a set to share with your team.
- Explore the [sandbox-kit-spec repository](https://github.com/docker/sandbox-kit-spec)
  for the authoritative v3 specification, build frontend, and examples.

For the earlier format used by built-in agents, see [Kits v2](/manuals/ai/sandboxes/customize/kits-v2.md).
To save an environment you've configured inside a sandbox,
see [Save a sandbox as a template](/manuals/ai/sandboxes/usage.md#saving-a-sandbox-as-a-template).
