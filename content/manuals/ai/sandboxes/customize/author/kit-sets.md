---
title: Compose a kit set
description: Author a sandbox kit set with published components, model settings, arguments, and team instructions, then publish it as one image.
keywords: sandboxes, sbx, kits, v3, sets, composition, workloads, mixins
weight: 5
---

{{< summary-bar feature_name="Docker Sandboxes sbx" >}}

A kit set combines published kits with your own runtime declarations into one
image you can version and share. Alongside its component list, a set can declare
capabilities, lifecycle hooks, instructions, and arguments. Consumers use one
reference for the resulting environment or bundle of mixins.

V3 kits require a preview build of `sbx`. See
[Install a v3-capable build](/manuals/ai/sandboxes/customize/use-kits.md#install-a-v3-capable-build).
To publish the set, you also need Docker Buildx and a registry namespace you
can push to. To run the Claude example, have an Anthropic API key available.

## Choose the components

A set can contain one workload and any number of mixins, or only mixins.
A workload set supplies a complete sandbox environment. A set of mixins
packages tools and settings to add to a workload with `--kit`.

This example combines two published kits, then adds model settings and team
instructions directly to the set:

| Kit | Role |
| --- | --- |
| `docker.io/docker/sbx-kit-shell:1.0.0` | Shell workload and base tools |
| `docker.io/docker/sbx-kit-claude-mixin:2.1.274` | Claude Code and requests for credentials and session storage |

You can try the components together before publishing a set:

```console
$ sbx run docker.io/docker/sbx-kit-shell:1.0.0 \
    --kit docker.io/docker/sbx-kit-claude-mixin:2.1.274
```

This starts a shell with Claude Code available. Component requirements
determine composition order; the order of `--kit` flags or entries
in the set doesn't control it.

## Write the set descriptor

The set adds a model choice and team instructions to the components. Its
`model` argument selects a value when you create a sandbox. A lifecycle
capability writes that value into a settings file for Claude to read.

Create a `team-claude` directory containing `team-claude.yaml`:

```yaml {title="team-claude/team-claude.yaml"}
# syntax=docker/sandbox-kit:3
schemaVersion: "3"
kind: set
displayName: Team Claude environment
version: "1.0.0"

args:
  model:
    default: sonnet
    enum: [sonnet, opus, haiku]

kits:
  - ref: docker.io/docker/sbx-kit-shell:1.0.0
  - ref: docker.io/docker/sbx-kit-claude-mixin:2.1.274

capabilities:
  - type: com.docker.sandbox/lifecycle@1
    config:
      files:
        - path: /home/agent/.config/team-claude/settings.json
          content: '{"model": "${{ kit.args.model }}"}'
          mode: "0644"
  - type: com.docker.sandbox/agent-context@1
    config:
      content: |
        Run the project's tests before reporting a task complete.
```

The `kits:` list is the set's content recipe. Don't add a companion
Dockerfile, a `dockerfile:` field, or a `build:` block. To include additional
software or static files in the image, publish a workload or mixin that
contains them and add it to the list. The set can generate runtime files with
lifecycle declarations, as this example does.

Every `ref` must be a published registry reference. Local paths and Git URLs
aren't accepted as components, even when building the set from local source.

The `model` argument belongs to the set. At sandbox creation, its value is
validated and written to the settings file. The set also contributes team
instructions. Neither customization needs a separate mixin or a change to the
Claude component. The component supplies Claude's network and credential
requests, so the set doesn't repeat them.

## Build and publish the set

Build with the set descriptor and push the result to your namespace:

```console
$ docker login
$ docker buildx build ./team-claude -f ./team-claude/team-claude.yaml \
    -t docker.io/<NAMESPACE>/team-claude:1.0.0 --push
```

Replace `<NAMESPACE>` with a namespace you can push to. The frontend pulls the
listed kits, checks their declared compatibility, and merges their layers and
declarations. A missing dependency, conflicting feature provider, or second
workload makes the build fail.

The published descriptor has `kind: workload`, because the set contains the
shell workload. It keeps a `kits:` record with the resolved manifest digest
of every component. A set containing only mixins publishes as `kind: mixin`.
`kind: set` is an authoring form; consumers receive an ordinary kit.

## Run with the set's settings

Store your Anthropic API key on the host, then run the published set with a
model choice:

```console
$ sbx secret set anthropic
$ sbx run docker.io/<NAMESPACE>/team-claude:1.0.0 \
    --name team-claude-opus --kit-arg model=opus
```

Approve the credential request when prompted. See
[Credential bindings](/manuals/ai/sandboxes/configuration/credentials.md#credential-bindings)
for authentication setup.

The result retains the shell workload's launch command. From that shell,
start Claude with the generated settings and the set's packaged instructions:

```console
$ claude --settings /home/agent/.config/team-claude/settings.json \
    --append-system-prompt-file /usr/share/sandbox/kit/team-claude/context.md
```

The shell workload selects `AGENTS.md` as its instruction profile. The
`--append-system-prompt-file` flag explicitly gives Claude the merged guidance
file, including the set's team instructions. The frontend stages that file
under the set descriptor's filename stem, `team-claude`.

Claude uses the selected model. Omitting `--kit-arg model=opus` uses the set's
`sonnet` default. Changing this create-time argument doesn't require publishing
another image; create another sandbox to use another value.

This set configures an existing agent through its settings interface. To
control the image's entrypoint or build the agent on your own Linux base,
see [Build an agent workload](/manuals/ai/sandboxes/customize/author/build-an-agent.md).

## Configure component arguments

The previous example declares a model argument for the set's own settings
file. Components can also declare arguments. Use a component's `args` map to
choose its create-time values during publication:

```yaml
kits:
  - ref: docker.io/my-org/linter-kit:1.0.0
    args:
      mode: check
```

This illustrative component must declare a create-time `mode` argument.
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

## Add an ACP adapter

For an Agent Client Protocol (ACP) client, add this component to `kits:` and
rebuild the set:

```yaml
  - ref: docker.io/docker/sbx-kit-claude-acp:0.79.0
```

The adapter requires `claude`, supplied by the Claude mixin. Use that mixin
with this adapter: it installs the binary at the path the adapter expects.
The Claude workload is a separate kit and isn't interchangeable with the
mixin in this combination.

The adapter runs as `claude-agent-acp`. The settings file in this guide is
passed to the interactive `claude` command; it doesn't configure the adapter.
Use your ACP client's configuration for its agent session.

Docker publishes the shell, Claude mixin, and adapter combination as
`docker.io/docker/sbx-kit-claude-acp-set:2.1.274`, with guidance for the adapter.
Use that reference if you don't need your own set declarations.

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
