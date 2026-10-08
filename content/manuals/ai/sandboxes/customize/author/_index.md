---
title: Author kits
description: Build and share sandbox environments with kits. Choose a tutorial, organize your source files, and decide what happens during build and setup.
keywords: sandboxes, sbx, kits, v3, workloads, mixins
weight: 20
---

{{< summary-bar feature_name="Docker Sandboxes sbx" >}}

Build a kit to give your team a repeatable sandbox environment. You can
package an agent, add a tool to use with different agents, or combine existing
kits into one kit your team can run. The guides in this section walk through
each approach.

> [!NOTE]
> Select a v3 workload and v3 mixins together.
> Built-in shortcuts such as `claude` and `codex` use v2 and can't be combined
> with v3 mixins. See [Version compatibility](/manuals/ai/sandboxes/customize/_index.md#version-compatibility)
> or the [v2 reference](/manuals/ai/sandboxes/customize/kits-v2.md).

## Choose what to author

- [Compose a kit set](/manuals/ai/sandboxes/customize/author/kit-sets.md) to combine
  kits you've published or chosen from a registry. A set can also add network
  access, setup commands, and agent instructions, and let users choose settings
  such as which model to use.
- [Build a tool mixin](/manuals/ai/sandboxes/customize/author/tool-mixins.md) to
  package a reusable tool with its network access and credentials.
- [Build an agent workload](/manuals/ai/sandboxes/customize/author/build-an-agent.md)
  to control the base environment, agent installation, and launch command.
  You can also [use an existing agent image](/manuals/ai/sandboxes/customize/author/base-images.md#use-an-image-in-a-v3-workload).
- [Debug a kit](/manuals/ai/sandboxes/customize/author/debug.md) when one of its
  install or startup commands fails. Read the progress output and the logs
  inside the sandbox.

For complete kits you can study and adapt, see
[Kit examples](https://github.com/docker/sandbox-kit-spec/tree/main/examples)
in the Docker Sandbox Kit Specification repository.

## Directory and build layout

Keep each kit's source in its own directory. A kit usually starts with two
files:

- A YAML descriptor identifies the kit as a workload, mixin, or set and
  declares its settings and requirements.
- A Dockerfile installs software and copies files into the image.

Use the kit's name for the directory and descriptor. If the kit has a
Dockerfile, use the same filename stem as the descriptor: `my-kit.yaml`
pairs with `my-kit.dockerfile`.

For example:

```text
my-kit/
├── my-kit.yaml
├── my-kit.dockerfile
├── context.md
└── files/
    └── settings.json
```

Organize supporting files in whatever way suits your kit. This example keeps
agent instructions in `context.md` and configuration files under `files/`.

### How the files work together

The Dockerfile defines what goes into the image, such as installed tools
and configuration files. For a workload, it also sets the command to launch.
The YAML descriptor defines the kit's settings and requirements, such as
network access, credentials, and agent instructions.

The descriptor starts with a syntax declaration:

```yaml
# syntax=docker/sandbox-kit:3
```

This selects the kit build frontend, which reads the descriptor and its
matching Dockerfile. The build produces a container image containing the
software, supporting files, and validated descriptor.

A kit set uses the descriptor's `kits:` list to combine published kits.

### Build and use the kit

When building with Docker Buildx, pass the YAML descriptor to `-f` and
the source directory as the build context:

```console
$ docker buildx build -f my-kit/my-kit.yaml -t my-kit:dev my-kit/
```

You can also give `sbx` a reference to a local kit directory. It builds the kit
when creating the sandbox. After publishing the image to a registry, you can
use its image reference instead.

For build options and publishing instructions, see
[Build and distribute kits](/manuals/ai/sandboxes/customize/author/distribute.md).

### Other source layouts

The separate descriptor and Dockerfile are one way to organize a kit.
You can also write a Dockerfile inline under `build: |`, select one with
`dockerfile:`, or embed the descriptor in a Dockerfile comment.

See [Authoring forms](https://github.com/docker/sandbox-kit-spec/blob/main/docs/spec/SPEC-v3.md#3-authoring-forms)
for the syntax.

## Capabilities

Installing a tool is often only part of the job. The tool might also need to
reach an API, authenticate with a credential, or run a setup command before
the agent starts. Describe those needs in the descriptor's `capabilities`
list. Each entry asks Docker Sandboxes to provide one of these features.

A kit's Dockerfile defines how its image is built. Its capabilities describe
what Docker Sandboxes needs to do when preparing and running the sandbox.
For example, a tool mixin can install an API client through its Dockerfile
and use a network capability to request access to that API. Building or
running the image with Docker alone doesn't apply these capability settings.

You can add capabilities to a workload, mixin, or set. Keep each request with
the kit that needs it, so a tool's access rules follow it when you use it with
another agent.

Each capability entry identifies the feature in `type`. Capabilities that
need settings take them in `config`. For example, `com.docker.sandbox/network-policy@1`
requests network access, and its `config` lists the domains to allow.

The `@1` identifies the capability's version. It is independent of the kit
format version and the `sbx` release. The
[upstream capability definitions](https://github.com/docker/sandbox-kit-spec/blob/main/docs/spec/SPEC-v3.md#72-well-known-types)
describe the available capabilities and their settings.

> [!NOTE]
> `sbx` doesn't apply `usb-device@1`, `privileged@1`, or `agent-sessions@1`
> requests. Its capability enforcement can let sandbox creation succeed even
> when a required capability is unsupported. Don't rely on this behavior:
> choose capabilities supported by the runtime where your kit will run.
> The `kit-registry@1` capability is restricted to approved OCI builder kits;
> local and Git kit sources don't receive it.

The guides here show how to use capabilities with Docker Sandboxes. For all
descriptor fields and the rules for combining kits, see the
[upstream v3 specification](https://github.com/docker/sandbox-kit-spec/blob/main/docs/spec/SPEC-v3.md).

## Choose when setup runs

Install tools and copy static files during the image build so you can reuse
them in every sandbox. Some setup needs to wait until the sandbox exists. For
example, a mixin can bring a CA certificate, but it needs the workload's tools
to register that certificate. Use a lifecycle capability to run commands or
generate files at the right point:

| Where to put the work | When it runs in `sbx` | Example |
| --- | --- | --- |
| Dockerfile `RUN` and `COPY` | Image build | Install a tool and copy its default configuration |
| Lifecycle `install` hooks | Once during sandbox creation, before the agent launches | Register a CA or populate a mounted directory |
| Lifecycle `files` | During creation, after install hooks and before the agent launches | Generate settings from kit arguments |
| Lifecycle `startup` hooks | Every sandbox start, alongside the agent | Start a background service or refresh state after a restart |

Lifecycle hooks are commands that run during sandbox creation or startup.
Make startup hooks safe to run more than once. In `sbx`, they run alongside
the agent, so the agent might start before they finish. If a command must
finish before every agent launch, put it in the workload's entrypoint instead.

### Pass environment variables to hooks

In `sbx`, install hooks receive a limited set of environment variables:
basic process variables, proxy settings, and certificate paths. If your command needs another
variable, name it in the hook's `env` list. For example, `env: [WORKSPACE_DIR]`
gives the command the mounted workspace's path. In `sbx`, startup hooks receive
the sandbox's environment without this filtering.

See [Kit authoring patterns](/manuals/ai/sandboxes/customize/author/patterns.md#run-setup-after-combining-kits)
for setup examples and the upstream
[lifecycle definition](https://github.com/docker/sandbox-kit-spec/blob/main/docs/spec/capabilities/com.docker.sandbox/lifecycle@1.md)
for the fields.

### Generated files

In lifecycle `files` content, write `${{ kit.env.NAME }}` to insert a
container environment variable. For example, write
`${{ kit.env.WORKSPACE_DIR }}` for the workspace path. Plain `$VAR` and
`${VAR}` are written unchanged, without substituting their values. In hook
commands, use shell syntax such as `$WORKSPACE_DIR` instead.

Docker Sandboxes substitutes `${{ kit.env.NAME }}` once, when it creates the
sandbox. `kit.env` reads the final container environment, independently of a
hook's `env` list.

## Set workload compute requirements

Set the workload's default CPU and memory allocation with the
[`resources@1` capability](https://github.com/docker/sandbox-kit-spec/blob/main/docs/spec/capabilities/com.docker.sandbox/resources@1.md)
in its descriptor. Use a whole number of CPU cores. Anyone running the kit
can override these defaults with `--cpus` and `--memory` when creating a sandbox.

Put these settings on the workload. `sbx` ignores resource settings on
mixins added separately, and it doesn't use the capability's `gpu` field
to select GPUs.

## Keep the sandbox running after sessions end

By default, a local sandbox stops automatically shortly after its last session
disconnects. A kit that runs a service, such as a development server on a
published port, can ask the sandbox to keep running instead. Declare the
[`long-running@1` capability](https://github.com/docker/sandbox-kit-spec/blob/main/docs/spec/capabilities/com.docker.sandbox/long-running@1.md),
which takes no `config`:

```yaml
capabilities:
  - type: com.docker.sandbox/long-running@1
```

A sandbox created with this kit behaves like one started with
[`sbx run --detached`](/manuals/ai/sandboxes/usage.md#keep-a-sandbox-running-in-the-background):
it keeps running until someone stops or removes it, including after you attach
to it and disconnect. When a workload, mixin, or set declares the capability,
it applies to the whole sandbox.

A kit that declares the capability runs only where the runtime can provide it.
If the kit also works when the sandbox stops after its sessions end, set
`optional: true` on the entry. A runtime that can't provide the capability
then skips it instead of refusing the kit. In `sbx`, two cases can't provide
it:

- Cloud sandboxes can't provide this capability. `sbx --cloud` refuses a kit
  that requires it and skips an optional entry.
- `sbx kit add` can't apply the capability to a running sandbox. A required
  entry fails, and an optional entry is skipped with a warning. Recreate the
  sandbox with the kit to apply it.
