---
title: Compose a kit set
description: Combine a Claude Code workload with an internal CLI mixin, customize runtime access, and publish the environment as one kit.
keywords: sandboxes, sbx, kits, v3, sets, composition, workloads, mixins
weight: 5
---

{{< summary-bar feature_name="Docker Sandboxes sbx" >}}

A kit set combines published kits with your own runtime declarations into one
image you can version and share. Alongside its component list, a set can declare
capabilities, lifecycle hooks, instructions, and arguments. Consumers use one
reference for the resulting environment or bundle of mixins.

This guide combines Claude Code with an internal CLI so the agent can use
your company's service. First, build and publish the mixin from
[Package an internal CLI](/manuals/ai/sandboxes/customize/author/kit-examples.md#package-an-internal-cli),
adapting it to your tool and API. You also need Docker Buildx and a registry
namespace you can push to. To run the result, you need `sbx` and credentials
for Claude Code and your internal API.

## Choose the components

A set can contain one workload and any number of mixins, or only mixins.
A workload set supplies a complete sandbox environment. A set of mixins
packages tools and settings to add to a workload with `--kit`.

This example uses two published components:

| Kit | Role |
| --- | --- |
| `docker.io/docker/sbx-kit-claude:2.1.274` | Claude Code agent and its environment |
| `docker.io/<NAMESPACE>/company-cli:1.0.0` | Your internal CLI, API network rule, and credential request |

Replace `<NAMESPACE>` with the namespace where you published the mixin.
You can try the components together before publishing a set:

```console
$ sbx run docker.io/docker/sbx-kit-claude:2.1.274 --name claude-components \
    --kit docker.io/<NAMESPACE>/company-cli:1.0.0
```

This starts Claude Code with the CLI available. Component requirements
determine composition order. The order of `--kit` flags or entries in the
set doesn't control it.

## Write the set descriptor

Create a `claude-tools` directory containing `claude-tools.yaml`:

```yaml {title="claude-tools/claude-tools.yaml"}
# syntax=docker/sandbox-kit:3
schemaVersion: "3"
kind: set
displayName: Claude Code with company tools
version: "1.0.0"

kits:
  - ref: docker.io/docker/sbx-kit-claude:2.1.274
  - ref: docker.io/<NAMESPACE>/company-cli:1.0.0
```

Replace `<NAMESPACE>` in the descriptor too. The components supply their own
network rules and credential requests, so the set doesn't repeat them.

The `kits:` list is the set's content recipe. Don't add a companion
Dockerfile, a `dockerfile:` field, or a `build:` block. To include additional
software or static files in the image, publish a workload or mixin that
contains them and add it to the list.

Every `ref` must be a published registry reference. Local paths and Git URLs
aren't accepted as components, even when building the set from local source.

## Add runtime access to the set

A set can also declare behavior of its own. For example, if this environment
needs an internal package registry, add a network capability to the set's
descriptor:

```yaml
capabilities:
  - type: com.docker.sandbox/network-policy@1
    config:
      runtime:
        allow:
          - packages.company.example:443
```

Replace the example domain with your registry's host. This permits network
access within the sandbox's policy. If the registry requires authentication,
add an appropriate credential declaration too. See
[Credentials](/manuals/ai/sandboxes/customize/author/kit-reference.md#credentials).

The internal CLI's API access stays with its mixin so it follows the tool to
other workloads. Access needed by this combined environment can live on the
set. You can also declare lifecycle hooks, generated files, instructions, and
arguments on the set without creating another mixin.

## Build and publish the set

Build with the set descriptor and push the result to your namespace:

```console
$ docker login
$ docker buildx build ./claude-tools -f ./claude-tools/claude-tools.yaml \
    -t docker.io/<NAMESPACE>/claude-tools:1.0.0 --push
```

The frontend pulls the listed kits, checks their declared compatibility, and
merges their layers and declarations. A missing dependency, conflicting
feature provider, or second workload makes the build fail.

The published descriptor has `kind: workload`, because the set contains the
Claude Code workload. It keeps a `kits:` record with the resolved manifest
digest of every component. A set containing only mixins publishes as
`kind: mixin`. `kind: set` is an authoring form. Consumers receive an ordinary
kit.

## Run the set

Store the internal API credential on the host, then run the published set:

```console
$ sbx secret set company-api
$ sbx run docker.io/<NAMESPACE>/claude-tools:1.0.0 --name claude-tools
```

Use the credential service name you declared in the mixin. Authenticate
Claude Code and approve the credential requests when prompted. See
[Credential configuration](/manuals/ai/sandboxes/configuration/credentials.md)
for authentication options and preparing unattended runs.

The result starts Claude Code with your CLI installed and its runtime access
declared. Consumers select one kit reference. They don't need to supply each
component with `--kit`.

To control the agent's base image or launch command, see
[Build an agent workload](/manuals/ai/sandboxes/customize/author/build-an-agent.md).

## Configure component arguments

Components can declare arguments for choices such as a linter's mode. To
extend the set with a linter, add its published reference to `kits:` and use
the component's `args` map to choose its create-time values during publication:

```yaml
kits:
  - ref: docker.io/my-org/linter-kit:1.0.0
    args:
      mode: check
```

These snippets show only the linter entry. Keep the workload and internal
CLI entries in your set. The illustrative linter must declare a create-time
`mode` argument.
Its build-time arguments, such as an installed tool version, are already
resolved in the component image. To change those, rebuild the component.

To let consumers choose a value, declare an argument on the set and pass it
to the component:

```yaml
args:
  lint_mode:
    default: check
    enum: [check, fix]

kits:
  - ref: docker.io/my-org/linter-kit:1.0.0
    args:
      mode: ${{ kit.args.lint_mode }}
```

This example assumes the component declares `mode` with the same default and
enum. The set must preserve the component's argument constraints. Consumers
supply `--kit-arg lint_mode=fix` to the published set; they don't address its
internal components. Only arguments declared on the set remain configurable.
See [Kit sets in the reference](/manuals/ai/sandboxes/customize/author/kit-reference.md#kit-sets)
for the argument rules.

## Review and update the composition

Before sharing a set, check its files, hooks, and capabilities together.
A set merges image layers in dependency order, with later layers taking
precedence for overlapping paths. This differs from composition at sandbox
creation, which rejects file collisions. Reordering `kits:` isn't a way to
choose which file wins. Give components separate paths and inspect the built
image for unintended replacements.

Component declarations also combine. Network rules accumulate, lifecycle hooks
run in composition order, and agent instructions combine into one staged file.
The workload selects the instruction filename. Avoid adding another filename
in the set's own agent-context declaration.

To update a component, change its reference, rebuild the set, and publish a
set version. An already published set keeps the component content it was
built with, even if a component tag moves. Add a `digest` beside each `ref`
in the source to pin rebuild inputs. Create another sandbox to try the updated
set; existing sandboxes retain their recorded kit configuration.

For signing and multi-platform builds, see
[Build and distribute kits](/manuals/ai/sandboxes/customize/author/distribute.md).
