---
title: Author kits
description: Define v3 workloads and mixins with build recipes, capabilities, lifecycle hooks, and composition requirements.
keywords: sandboxes, sbx, kits, v3, workloads, mixins
weight: 20
---

{{< summary-bar feature_name="Docker Sandboxes sbx" >}}

> [!NOTE]
> V3 kits are experimental. Select a v3 workload and v3 mixins together.
> Built-in shortcuts such as `claude` and `codex` use v2 and can't be combined
> with v3 mixins. See [Version compatibility](/manuals/ai/sandboxes/customize/_index.md#version-compatibility)
> or the [v2 reference](/manuals/ai/sandboxes/customize/kits-v2/_index.md).

When you author a kit, you work in a directory of source files. A typical kit
has a YAML file and a Dockerfile:

```text
opencode-python/
├── opencode-python.yaml
└── opencode-python.dockerfile
```

The YAML file is the kit's descriptor. It identifies the kit as a workload or
mixin and declares what Docker Sandboxes should do when running it. The
Dockerfile defines the software and files to include and, for a workload,
the command to launch.

Building this directory produces a container image containing the kit's files and
its descriptor. You can share that image through a container registry. During
development, `sbx` can build directly from the directory when you create a
sandbox.

## Build a workload

Suppose your team uses OpenCode to work on Python projects. Package it with
Ruff and instructions to check Python changes before handing work back to you.
Everyone using the kit gets the same linter version and review workflow.

The Dockerfile starts from Docker's OpenCode [base image](/manuals/ai/sandboxes/customize/author/base-images.md) and
installs Ruff:

```dockerfile {title="opencode-python/opencode-python.dockerfile"}
FROM docker/sandbox-templates:opencode
USER agent
RUN uv tool install ruff==0.12.12
ENTRYPOINT ["opencode"]
CMD []
```

The template supplies OpenCode, Python, uv, and the `agent` user. Ruff is
installed during the build, so it is ready when the agent starts.

The descriptor declares this as a workload, connects OpenCode to the
Anthropic API, and gives it the team's review instructions:

```yaml {title="opencode-python/opencode-python.yaml"}
# syntax=docker/sandbox-kit:3
schemaVersion: "3"
kind: workload

capabilities:
  - type: com.docker.sandbox/network-policy@1
    config:
      runtime:
        allow:
          - api.anthropic.com
          - opencode.ai
          - models.dev
          - registry.npmjs.org
          - pypi.org
          - files.pythonhosted.org
  - type: com.docker.sandbox/credential@1
    config:
      service: anthropic
      phase: runtime
      apiKey:
        name: ANTHROPIC_API_KEY
        proxyManaged: true
        inject:
          - domain: api.anthropic.com
            header: x-api-key
            format: "%s"
  - type: com.docker.sandbox/agent-context@1
    config:
      filename: AGENTS.md
      content: |
        Ruff is installed. Run `ruff check` on Python files you change,
        and fix any lint errors before reporting completion.
        Follow the project's existing configuration and test commands.
```

The `capabilities` list describes what the sandbox provides at runtime:
network access, authentication, and instructions for the agent. The credential
entry names the service; you store the actual API key on your host.

This example extends an existing agent environment. To prepare your own Linux
base image, install an agent, and configure its runtime needs step by step,
see [Build an agent](/manuals/ai/sandboxes/customize/author/build-an-agent.md).

## Capabilities

A capability declares a resource or behavior that a kit needs from Docker
Sandboxes at runtime, such as network access, credentials, lifecycle hooks,
or agent instructions. Declare these requests in the descriptor's
`capabilities` list.

Workloads and mixins use the same capability format. For example, a workload
can declare the network access its agent needs, and a mixin can request access
to an additional service.

Each entry's `type` identifies the capability and its settings version, and
`config` contains those settings. For a complete descriptor using network,
credential, and instruction capabilities, see the
[OpenCode workload example](#build-a-workload). The
[capability reference](/manuals/ai/sandboxes/customize/author/kit-reference.md#runtime-capabilities) lists the
available types, their fields, and Docker Sandboxes support.

### Control network access

Use the network-policy capability to declare which domains a sandbox can
reach. For example, this mixin permits requests to the GitHub API:

```yaml {title="github-access/github-access.yaml"}
# syntax=docker/sandbox-kit:3
schemaVersion: "3"
kind: mixin

capabilities:
  - type: com.docker.sandbox/network-policy@1
    config:
      runtime:
        allow: [api.github.com]
```

The type `com.docker.sandbox/network-policy@1` selects version 1 of the
network-policy settings. `runtime.allow` lists domains the running sandbox
can reach.

This kit needs only its YAML file: it configures network access without adding
software or files to the image. Save the descriptor in `github-access` and
run it with a v3 workload, such as the
[OpenCode workload example](#build-a-workload):

```console
$ sbx run ./opencode-python --name python-github --kit ./github-access
```

Network rules from the selected kits combine. A kit's `deny` entries take
precedence over kit allow entries, and the resulting rules participate in the
sandbox's [policy precedence](/manuals/ai/sandboxes/governance/concepts.md#precedence). Use
`sbx policy log` to investigate refused connections.

When organization governance is active, only organization allow rules grant
access. Kit allow rules don't grant additional access, but kit deny rules
still restrict it.

### Authenticate to external services

A credential capability names the service a kit needs and declares how to
authenticate to it. The [OpenCode workload example](#build-a-workload)
declares `service: anthropic`. Store its API key on the host using that
service name:

```console
$ sbx secret set anthropic
```

The example sets `proxyManaged: true`. OpenCode receives a placeholder in
`ANTHROPIC_API_KEY`, and the host proxy inserts the real key into requests to
`api.anthropic.com`. The key stays on the host. First-run approval authorizes
the kit to use it through that mechanism and on that domain.

See [Credentials](/manuals/ai/sandboxes/customize/author/kit-reference.md#credentials) for a descriptor example and
API key and OAuth fields. [Credential configuration](/manuals/ai/sandboxes/configuration/credentials.md)
covers host-side storage and approval, including preparation for unattended runs.

### Give the agent instructions

Installing a tool makes it available. Agent instructions tell the agent when
and how to use it, where to find shared configuration, and which checks to run
before reporting a task complete.

A mixin can contribute instructions through the agent-context capability:

```yaml
capabilities:
  - type: com.docker.sandbox/agent-context@1
    config:
      content: |
        Follow the project's contribution guide when changing code.
        Run the project's lint and test commands before reporting completion.
        If a check fails, include the failure in your response.
```

The workload chooses the profile filename, such as `AGENTS.md` or `CLAUDE.md`.
Docker Sandboxes generates this file in the parent directory of the mounted
workspace inside the sandbox. For example, a workspace at
`/home/agent/workspace` has its generated profile at `/home/agent/AGENTS.md`.
The profile sits outside the mount and doesn't replace instructions in your
project.

Inline workload instructions, such as those in the
[OpenCode workload example](#build-a-workload), go directly into the profile.
Inline mixin instructions go into separate files under `kits-agent-context/`
beside the profile. The profile indexes those files for the agent to read
when needed. See
[Contribute agent instructions](/manuals/ai/sandboxes/customize/author/kit-examples.md#contribute-agent-instructions).

Use `contentFile` to keep longer guidance in a
Markdown file beside the descriptor. The build packages that file in the
image, and the profile points to it instead of copying its text. The frontend
also stages these files for mixins without a build recipe. See
[Add agent instructions](/manuals/ai/sandboxes/customize/author/build-an-agent.md#add-agent-instructions) for a
workload that uses `contentFile`.

You can also package an agent skill: a reusable set of instructions for a
particular task, such as reviewing a Dockerfile. Skills are files installed
where that agent looks for them. See
[Ship a Claude Code skill](/manuals/ai/sandboxes/customize/author/kit-examples.md#ship-a-claude-code-skill).

## Build content and runtime setup

Build tools and static content into the kit's image so sandboxes can reuse
them. Reserve runtime setup for work that needs an individual sandbox's
state, such as its mounted workspace path or a host-provided credential.

Use lifecycle hooks for that work. A hook is a command Docker Sandboxes runs
at a particular point in the sandbox's life. Install hooks initialize each
sandbox during creation; startup hooks run each time it starts. The lifecycle
capability can also write configuration files during creation.

| Mechanism | When it runs | Use it for |
| --- | --- | --- |
| Dockerfile `RUN` and `COPY` | Kit build | Install tools, compile binaries, and package static content |
| Lifecycle `install` hooks | Once per sandbox, during creation | Initialize state that needs runtime credentials or mounted paths |
| Lifecycle `startup` hooks | Each sandbox start | Start a service or refresh state after a restart |
| Lifecycle `files` | Sandbox creation | Write configuration for an individual sandbox |

Build results are reusable across sandboxes. Install hooks run for each
sandbox you create. Put software installation and compilation in the build
when they don't require sandbox-specific inputs.

### Lifecycle hooks

Hooks are declared through `com.docker.sandbox/lifecycle@1`. For example, write
a default config and start a service that the kit's image already contains:

```yaml
capabilities:
  - type: com.docker.sandbox/lifecycle@1
    config:
      files:
        - path: /home/agent/.config/my-service/config.json
          content: '{"port": 8080}'
          overwrite: false
      startup:
        - command: [my-service, --config, /home/agent/.config/my-service/config.json]
          user: "1000"
          background: true
```

Startup hooks must tolerate repeated execution. In Docker Sandboxes they run
through a background dispatcher and don't block the agent's launch. Use an
install hook or a workload entrypoint script for setup the agent must wait for.
String commands run through `sh -c`; an argument list invokes the command
directly. See [Lifecycle](/manuals/ai/sandboxes/customize/author/kit-reference.md#lifecycle) for command fields,
file permissions, and runtime behavior.

### Static files

Static content belongs in the image. Use Dockerfile `COPY` to put a tool config,
helper script, or reference document at its destination, or stage it at a
kit-specific path for a lifecycle hook to copy later.

Files destined for a mounted workspace or persistent volume need that second
step: a mount can hide files baked into the image at its mount path. Copy from
the staged image path after the mount is available. See
[Copy shared configuration](/manuals/ai/sandboxes/customize/author/kit-examples.md#copy-shared-configuration).

### Customize agent settings

When an agent supports additional settings files, keep team defaults in a
separate file and point the agent to it. For example, Claude Code has a
`--settings` launch option, and OpenCode reads the path in `OPENCODE_CONFIG`.

Package fixed settings with Dockerfile `COPY`. Use lifecycle `files` for
settings that contain kit argument values chosen when creating the sandbox.
The workload controls its launch options; a mixin's `ENTRYPOINT` doesn't
change how the agent starts. See
[Write the model settings](/manuals/ai/sandboxes/customize/author/build-an-agent.md#write-the-model-settings) for a
complete example.

### Set environment variables

Use Dockerfile `ENV` for tool settings that apply to every sandbox using the
kit. For example, make Python write output without buffering:

```dockerfile
ENV PYTHONUNBUFFERED=1
```

For values the user chooses when creating a sandbox, declare a kit argument
with an `env` mapping. See [Declare arguments](#declare-arguments).
Use credential capabilities for secrets.

Docker Sandboxes sets `HTTP_PROXY`, `HTTPS_PROXY`, `NO_PROXY`, and their
lowercase equivalents to route traffic through its policy and credential
proxy. Leave those variables to the sandbox. Configure a corporate proxy on
the host using [Upstream proxy](/manuals/ai/sandboxes/architecture.md#upstream-proxy).

Some tools need a shell initialization script, such as a version manager's
`init.sh`. With Docker sandbox templates, append the initialization commands
to `/etc/sandbox-persistent.sh` in a lifecycle install hook. The templates
source this file for interactive and non-interactive Bash commands. Append
to it so other kits' settings remain, and keep shell completion scripts out
of it because they can fail in non-interactive shells.

## Compose kits

Composition combines one workload kit with its mixins. Select the kits
when [creating a sandbox](/manuals/ai/sandboxes/customize/use-kits.md#add-mixins). A kit's descriptor can also declare
relationships with other kits, such as a dependency on a tool they supply.

Every selected kit must use v3. A v3 mixin can't extend a v2 built-in agent,
and a v2 mixin can't extend a v3 workload. Choose the workload and compatible
mixins together when creating the sandbox.

For example, a tool kit can advertise what it supplies:

```yaml
provides: ["my-tool@1.0.0"]
```

A mixin that needs that tool can declare a requirement:

```yaml
requires: ["my-tool >= 1.0.0"]
```

Include both kits when creating the sandbox. Docker Sandboxes checks that
the selected set satisfies the requirement; it doesn't search a registry or
install a package to satisfy it. Kit authors declare these names explicitly.
Installing a tool in a Dockerfile doesn't automatically add a `provides` entry.

Providers are applied before the kits that require them. Independent mixins
are ordered by reference, so reordering `--kit` flags isn't an override
mechanism. Use `integrates` for a relationship that applies only when another
kit is present, and `conflicts` to reject an incompatible combination. See
[Composition fields](/manuals/ai/sandboxes/customize/author/kit-reference.md#composition-fields) for these rules.

### Avoid conflicting customizations

Give each declared feature one provider. If your workload provides `node`,
don't add a mixin that also provides `node`, even at the same version. The
spec rejects multiple owners of a normalized feature name. Credential
requests also have one owner per service and phase across the selected kits.

The workload supplies the environment and launch command. Mixins add files
and runtime declarations. Two kits contributing the same image file cause a
composition error. Conflicting image environment values also cause an error,
while `PATH` additions are combined. Give each kit its own paths for staged
content and avoid having multiple kits manage the same config.

To derive a workload, build its Dockerfile from the base image you want and
declare its runtime capabilities in its descriptor. Dockerfile `FROM` inherits
image content and config, but doesn't merge a parent kit's descriptor. V3 has
no `extends` field. To add to a workload without deriving another workload,
use a mixin.

Select all kits when creating the sandbox. To change a v3 kit set, recreate the
sandbox with the desired workload and mixins. `sbx kit add` doesn't apply v3
changes to an existing sandbox.

## Declare arguments

Kit arguments let users choose values without editing the kit's source. For
example, a tool kit can offer a mode that becomes an environment variable in
the sandbox:

```yaml
args:
  mode:
    default: check
    enum: [check, fix]
    env: TOOL_MODE
```

`args.mode` declares the input, its default, and its accepted values. The `env`
field exports the chosen value as `TOOL_MODE`. You can also reference the value
in the descriptor as `${{ kit.args.mode }}`, for example in the content of a
configuration file. Without `env`, the value is available only through that
substitution.

For shared and kit-scoped CLI values, see [Pass arguments to kits](/manuals/ai/sandboxes/customize/use-kits.md#pass-arguments-to-kits).

Arguments can also select build-time inputs, such as a tool version. Declare
those with `buildArg` and supply them to `docker buildx build --build-arg`
using the kit argument's name. These values are resolved into the published
kit, so changing them requires rebuilding. For an example, see
[Build a tool overlay](/manuals/ai/sandboxes/customize/author/kit-examples.md#build-a-tool-overlay).

Argument values are plain text and can be recorded in shell history and
sandbox state. Use credential capabilities for secrets. See
[Arguments](/manuals/ai/sandboxes/customize/author/kit-reference.md#arguments) for validation and mapping fields.

## Directory and build layout

Keep a kit's descriptor, Dockerfile, and supporting files in one source
directory. Use matching filename stems for the YAML descriptor and its
companion Dockerfile so the build can discover the recipe. For a kit with
instructions and a configuration file, the layout could be:

```text
my-kit/
├── my-kit.yaml
├── my-kit.dockerfile
├── context.md
└── files/
    └── settings.json
```

The descriptor's first line, `# syntax=docker/sandbox-kit:3`, selects the kit
BuildKit frontend. Pass the YAML file to `docker buildx build -f`, with its
directory as the build context. The frontend finds the same-stem
`my-kit.dockerfile`, builds the content, validates the declarations, and
publishes them together.

For `sbx` to discover a local source kit, keep exactly one v3 descriptor at the
directory root, named with the `.yaml` extension. A `.dockerfile` with a
comment descriptor also works. Avoid `spec.yaml` and `spec.yml`: these names
select the v1/v2 loader. Matching the descriptor stem to the directory name
also keeps build and create argument scopes consistent.

A workload needs a Dockerfile to supply its environment and launch command.
A mixin needs one when it adds image content. A mixin that only declares
runtime behavior, such as the [GitHub network mixin](#control-network-access),
can omit it.

For a single-file kit, use `build: |` with literal Dockerfile text in the
YAML. You can also select a differently named recipe with `dockerfile:` or
embed a descriptor in a Dockerfile comment block. See
[Authoring forms](/manuals/ai/sandboxes/customize/author/kit-reference.md#authoring-forms) for syntax and discovery
rules. The `files/` name in this example is an authoring convention, not a
special runtime directory.

### Build a mixin

A mixin contributes an overlay: files added or changed by its recipe. The
recipe's base image is a build environment, and its unchanged filesystem
doesn't become part of the overlay. For a tool built in a separate stage,
use a final `FROM scratch` stage and copy the tool and everything it needs
into that stage.

A copied binary must be compatible with the workload's architecture and
libraries. Ship dependencies that the workload doesn't supply, and declare
any requirements on other kits as described in [Compose kits](#compose-kits).
See [Build a tool overlay](/manuals/ai/sandboxes/customize/author/kit-examples.md#build-a-tool-overlay).

A mixin's environment additions become part of the composed image, but its
`ENTRYPOINT`, `CMD`, `USER`, and `WORKDIR` don't replace the workload's launch
configuration.

## Debug lifecycle hooks

For a background service, redirect its startup command's output to a file
inside the sandbox, then read that file with `sbx exec`. Set
`background: true` on the hook rather than adding `&` to the shell command.
See [Run a hook on every start](/manuals/ai/sandboxes/customize/author/kit-examples.md#run-a-hook-on-every-start).

## Continue authoring

- [Base images](/manuals/ai/sandboxes/customize/author/base-images.md) covers Docker-provided images and choosing
  your own Linux base.
- [Build an agent](/manuals/ai/sandboxes/customize/author/build-an-agent.md) walks through a complete workload from
  a custom base image.
- [Mixin examples](/manuals/ai/sandboxes/customize/author/kit-examples.md) includes tools,
  shared configuration, instructions, and lifecycle hooks.
- [Build and distribute kits](/manuals/ai/sandboxes/customize/author/distribute.md) covers publishing and signing.
- [Spec reference](/manuals/ai/sandboxes/customize/author/kit-reference.md) documents descriptor fields and runtime support.
