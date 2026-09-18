---
title: Kits
description: Use sandbox kits and share environments as kit sets composed from workloads and mixins.
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

This page covers v3 kits. See [Version compatibility](#version-compatibility)
for preview-build requirements and the earlier kit formats.

## What is a kit?

A kit packages software and configuration for a sandbox. Its YAML descriptor
specifies the runtime behavior it needs, such as network access, credentials,
setup commands, and agent instructions. When published, the kit is a container
image that carries both its files and its descriptor.

You give Docker Sandboxes a kit reference to use that package. Sandboxes
prepares its files and applies its runtime configuration. You can use one kit
for a complete environment or combine kits that contribute different parts.

## Workloads and mixins

Kits have two roles in a sandbox: one workload supplies the base environment
and launch command, and mixins add tools or behavior to it.

| Kind | What it supplies | How you use it |
| --- | --- | --- |
| `workload` | The environment and command to run, such as an agent or a shell | Pass it to `sbx run` or `sbx create` |
| `mixin` | Additional tools, configuration, or runtime behavior | Add it with `--kit` |

For example, a shell workload can provide the base environment, a Claude Code
mixin can add the agent, and another mixin can add an Agent Client Protocol
(ACP) adapter. Together, they give you a shell with Claude Code and its adapter
available. Each component can be maintained and reused separately.

## Kit sets

A set brings those components together in a descriptor of its own. It lists
the kits to include and can add capabilities, lifecycle hooks, instructions,
and arguments for the combined environment. Publish the set to give consumers
one reference with the component versions already selected.

Docker publishes the shell, Claude Code, and ACP adapter combination described
above as a set. Run it with:

```console
$ sbx run docker.io/docker/sbx-kit-claude-acp-set:2.1.274
```

This starts the shell with Claude Code and its ACP adapter available. See
[Use kits](/manuals/ai/sandboxes/customize/use-kits.md) for prerequisites and usage.

Workload and mixin describe how a published kit is used. A set describes how
an author builds it from other kits. Publishing a set that contains a workload
produces a workload kit. A set containing only mixins produces a mixin kit,
which you can add to a workload with `--kit`.

To customize and publish your own combination, see
[Compose a kit set](/manuals/ai/sandboxes/customize/author/kit-sets.md).

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

## Version compatibility

> [!NOTE]
> V3 kits are experimental. The format and CLI commands are subject to change.
> Share spec feedback or report issues in
> [docker/sandbox-kit-spec](https://github.com/docker/sandbox-kit-spec/issues).

The [sandbox-kit-spec repository](https://github.com/docker/sandbox-kit-spec)
contains the authoritative v3 specification, build frontend, and examples.

V3 requires an `sbx` nightly with v3 support; it isn't available in the stable
release. See [Install a v3-capable build](/manuals/ai/sandboxes/customize/use-kits.md#install-a-v3-capable-build).

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

## Choose your next step

- [Use kits](/manuals/ai/sandboxes/customize/use-kits.md) to run a published environment or combine
  a workload with mixins.
- [Author kits](/manuals/ai/sandboxes/customize/author/_index.md) to compose a kit set or build workload and
  mixin components with their own runtime requirements.

For the earlier format used by built-in agents, see [Kits v2](/manuals/ai/sandboxes/customize/kits-v2/_index.md).
To capture an interactively configured sandbox's container filesystem for reuse,
see [Save a sandbox as a template](/manuals/ai/sandboxes/usage.md#saving-a-sandbox-as-a-template).
