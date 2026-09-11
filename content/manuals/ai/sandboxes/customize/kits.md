---
title: Kits
description: Build, compose, and distribute sandbox workloads and extensions with kits v3, using image content and declared runtime capabilities.
keywords: sandboxes, sbx, kits, v3, workloads, mixins, capabilities, builds, composition
weight: 20
---

{{< summary-bar feature_name="Docker Sandboxes sbx" >}}

A kit packages a sandbox workload or an extension to one. It combines image
content, such as an agent or a toolchain, with declarations for the runtime:
network access, credentials, storage, lifecycle hooks, and agent instructions.
Build the kit once and reuse its content across sandboxes, while the runtime
applies its declarations when you create and run each sandbox.

This page covers kits v3, the preferred format for authoring kits. Docker
Sandboxes also supports v1 and v2. See [Kits v2](kits-v2/_index.md) for the
`spec.yaml` format used by existing v2 kits. A single composition cannot combine
v3 kits with v1 or v2 kits.

> [!NOTE]
> Kits are experimental. The format and CLI commands are subject to change.
> Share feedback in the
> [docker/sbx-releases](https://github.com/docker/sbx-releases) repository.

## Workloads and mixins

A sandbox runs one workload kit. You can add mixin kits to customize that
workload:

| Kind | What it supplies | How you use it |
| --- | --- | --- |
| `workload` | The environment and command to run, such as an agent or a shell | Pass it to `sbx run` or `sbx create` |
| `mixin` | Additional tools, configuration, or runtime behavior | Add it with `--kit` |

For example, a team might use a Claude Code workload with a mixin that adds a
linter and another that supplies the team's review instructions. Each kit can
be maintained and shared separately.

## Kit files and images

When you author a kit, you work in a directory of source files. A typical kit
has a YAML file and a Dockerfile:

```text
my-shell/
├── my-shell.yaml
└── my-shell.dockerfile
```

The YAML file is the kit's descriptor. It identifies the kit as a workload or
mixin and declares what Docker Sandboxes should do when running it. The
Dockerfile defines the software and files to include and, for a workload,
the command to launch.

Building this directory produces a container image containing the kit's files and
its descriptor. You can share that image through a container registry. During
development, `sbx` can build directly from the directory when you create a
sandbox.

### Build a workload

For example, these two files define a shell workload with the `jq` tool
installed. The descriptor identifies it as a v3 workload:

```yaml {title="my-shell/my-shell.yaml"}
# syntax=docker/runtime-kit:3
schemaVersion: "3"
kind: workload
```

The Dockerfile starts from a sandbox template, installs `jq`, and sets Bash
as the command to run:

```dockerfile {title="my-shell/my-shell.dockerfile"}
FROM docker/sandbox-templates:shell
USER root
RUN apt-get update && apt-get install -y jq \
    && rm -rf /var/lib/apt/lists/*
USER agent
ENTRYPOINT ["bash"]
CMD []
```

The template supplies the sandbox's base environment, including the `agent`
user. The installation runs as root, then the Dockerfile switches back to
`agent` for the shell. For an agent workload, the Dockerfile would install
and launch the agent instead. See [Build an agent](build-an-agent.md) for a
complete walkthrough.

## Run a kit

Save the two files in `my-shell` and run this command from its parent
directory:

```console
$ sbx run ./my-shell
```

`sbx` builds the kit and opens a sandbox shell with `jq` installed. Your current
directory is mounted as the workspace. Unchanged builds reuse cached results.

The workload reference is a positional argument. It can also be a published
image or a Git URL:

```console
$ sbx run docker.io/<NAMESPACE>/my-shell:1.0.0
$ sbx run "git+https://github.com/<ORG>/<REPOSITORY>.git#ref=<COMMIT>&dir=my-shell"
```

The same positional syntax applies to `sbx create` and to v1 and v2 sandbox
kits. Use `--kit` to add mixins when creating a sandbox:

```console
$ sbx run ./my-shell --name shell-with-tools --kit ./my-tool --kit ./team-config .
```

`sbx` combines the workload and mixins into the sandbox's environment. This
combination is called a composition. An existing sandbox with the same name
is reused, so choose a different `--name` when trying a different kit set.
For complete mixins to try, see [Kit examples](kit-examples.md).

## Capabilities

A kit can include a tool in its image, but using that tool might also require
network access or a credential. The descriptor tells Docker Sandboxes what to
provide through its `capabilities` list. Capabilities also cover behavior such
as running setup commands and supplying instructions to an agent.

Workloads and mixins use the same capability format. For example, a workload
can declare the network access its agent needs, and a mixin can request access
to an additional service.

### Control network access

This mixin permits requests to the GitHub API. Save it as
`github-access/github-access.yaml`:

```yaml {title="github-access/github-access.yaml"}
# syntax=docker/runtime-kit:3
schemaVersion: "3"
kind: mixin

capabilities:
  - type: com.docker.runtime/network-policy@1
    config:
      runtime:
        allow: [api.github.com]
```

The entry's `type` names the capability, and `config` contains its settings.
Here, `runtime.allow` lists a domain the running sandbox can reach.

This kit needs only its YAML file: it configures network access without adding
software or files to the image. Run it with the shell workload:

```console
$ sbx run ./my-shell --name shell-github --kit ./github-access
```

Network rules from the selected kits combine. A kit's `deny` entries take
precedence over kit allow entries, and the resulting rules participate in the
sandbox's [policy precedence](../governance/concepts.md#precedence). Use
`sbx policy log` to investigate refused connections.

The same pattern of `type` and `config` applies to other capabilities. Each
type has its own settings; `@1` identifies the version of those settings.
See [Runtime capabilities](kit-reference.md#runtime-capabilities) for the
available types and their Docker Sandboxes support, including limits on
required and optional requests.

### Authenticate to external services

A credential capability names the service a kit needs and declares how to
authenticate to it. Users supply the secret on the host. For example, for a
service named `my-service`:

```console
$ sbx secret set my-service
```

On the first interactive run, `sbx` also asks the user to approve how the kit
uses that credential. A kit can request proxy-managed authentication, where
the host proxy inserts the secret into outbound requests and the real value
stays outside the sandbox.

See [Credentials](kit-reference.md#credentials) for a descriptor example and
API key and OAuth fields. [Credential configuration](../configuration/credentials.md)
covers host-side storage and approval, including preparation for unattended runs.

## Build content and runtime setup

The shell example installs `jq` while building the kit. Every sandbox using
that image starts with the tool available. Other setup needs information that
exists only when a sandbox runs, such as the mounted workspace path or a
host-provided credential.

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

Hooks are declared through `com.docker.runtime/lifecycle@1`. For example, write
a default config and start a service that the kit's image already contains:

```yaml
capabilities:
  - type: com.docker.runtime/lifecycle@1
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
directly. See [Lifecycle](kit-reference.md#lifecycle) for command fields,
file permissions, and runtime behavior.

### Static files

Static content belongs in the image. Use Dockerfile `COPY` to put a tool config,
helper script, or reference document at its destination, or stage it at a
kit-specific path for a lifecycle hook to copy later.

Files destined for a mounted workspace or persistent volume need that second
step: a mount can hide files baked into the image at its mount path. Copy from
the staged image path after the mount is available. See
[Copy shared configuration](kit-examples.md#copy-shared-configuration).

## Compose kits

Selecting a workload and mixins with `sbx run` is enough to combine them. When
a kit depends on another kit's tools, its descriptor can also declare that
relationship.

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
[Composition fields](kit-reference.md#composition-fields) for these rules.

### Avoid conflicting customizations

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

## Pass arguments to kits

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

```console
$ sbx run ./my-agent --kit ./my-tool --kit-arg mode=fix .
```

`args.mode` declares the input, its default, and its accepted values. The `env`
field exports the chosen value as `TOOL_MODE`. You can also reference the value
in the descriptor as `${{ kit.args.mode }}`, for example in the content of a
configuration file. Without `env`, the value is available only through that
substitution.

A bare argument name applies to every kit that declares it. To target one kit,
prefix the name with its handle and a period:

```console
$ sbx run ./my-agent --kit ./my-tool --kit-arg my-tool.mode=fix .
```

The handle is the local directory name, the Git subdirectory or repository
name, or the last repository segment in an OCI reference. Scoped values take
precedence over shared values. Use `--kit-args-file <FILE>` for reusable
`name=value` entries; `--kit-arg` values take precedence over file values.

Arguments can also select build-time inputs, such as a tool version. Declare
those with `buildArg` and supply them to `docker buildx build --build-arg`
using the kit argument's name. These values are resolved into the published
kit, so changing them requires rebuilding. For an example, see
[Build a tool overlay](kit-examples.md#build-a-tool-overlay).

Argument values are plain text and can be recorded in shell history and
sandbox state. Use credential capabilities for secrets. See
[Arguments](kit-reference.md#arguments) for validation and mapping fields.

## Directory and build layout

The shell example uses two files with matching names. Kits can also include
configuration, scripts, and instructions that the Dockerfile copies into the
image. A larger directory might look like this:

```text
my-kit/
├── my-kit.yaml
├── my-kit.dockerfile
├── context.md
└── files/
    └── settings.json
```

The descriptor's first line, `# syntax=docker/runtime-kit:3`, selects the kit
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
runtime behavior, like the GitHub network example, can omit it.

For a single-file kit, use `build: |` with literal Dockerfile text in the
YAML. You can also select a differently named recipe with `dockerfile:` or
embed a descriptor in a Dockerfile comment block. See
[Authoring forms](kit-reference.md#authoring-forms) for syntax and discovery
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
See [Build a tool overlay](kit-examples.md#build-a-tool-overlay).

A mixin's environment additions become part of the composed image, but its
`ENTRYPOINT`, `CMD`, `USER`, and `WORKDIR` don't replace the workload's launch
configuration.

## Packaging and distribution

A published v3 kit is an OCI image. Use Docker Buildx to build and push it:

```console
$ docker login
$ docker buildx build ./my-kit -f ./my-kit/my-kit.yaml \
    -t docker.io/<NAMESPACE>/my-kit:1.0.0 --push
```

For a workload, launch the published reference. For a mixin, pass it with
`--kit`:

```console
$ sbx run docker.io/<NAMESPACE>/my-agent:1.0.0 \
    --kit docker.io/<NAMESPACE>/my-kit:1.0.0 .
```

Build for both supported Linux architectures when distributing across machines:

```console
$ docker buildx build ./my-kit -f ./my-kit/my-kit.yaml \
    --platform linux/amd64,linux/arm64 \
    -t docker.io/<NAMESPACE>/my-kit:1.0.0 --push
```

An image present only in the host Docker image store isn't available to the
sandbox runtime by registry reference. Push it to a registry, or pass a local
source directory to `sbx` for the build-and-run development loop. The
`sbx kit pack`, `push`, and `pull` packaging commands belong to
[Kits v2](kits-v2/_index.md#packaging-and-distribution).

You can also share source through Git. Select the kit directory with `dir` and
pin the source with `ref`:

```console
$ sbx run "git+https://github.com/<ORG>/<REPOSITORY>.git#ref=<COMMIT>&dir=my-agent" .
```

For private registries, configure
[Registry credentials](../configuration/credentials.md#registry-credentials).
Include `docker.io/` explicitly for Docker Hub references.

### Restrict kit sources

`kit.allowedSources` controls permitted remote kit sources. Its default permits
Docker Hub. To include a Git publisher, set the complete list of prefixes you
want to permit:

```console
$ sbx settings set kit.allowedSources '["docker.io/","github.com/docker/"]'
```

Prefixes match at path-segment boundaries. Local source directories are
controlled separately by `kit.allowLocalKits`, which defaults to `true`:

```console
$ sbx settings set kit.allowLocalKits false
```

For non-interactive configuration, use `DOCKER_SANDBOXES_KIT_ALLOWED_SOURCES`
and `DOCKER_SANDBOXES_KIT_ALLOW_LOCAL`.

### Sign and verify kits

Sign the OCI image after pushing it:

```console
$ sbx kit sign docker.io/<NAMESPACE>/my-kit:1.0.0
$ sbx kit verify docker.io/<NAMESPACE>/my-kit:1.0.0 \
    --certificate-identity <SIGNER_IDENTITY> \
    --certificate-oidc-issuer <ISSUER_URL>
```

These commands use Cosign-compatible Sigstore signatures. For keyless signing,
verification must specify the signer's certificate identity and OpenID Connect
issuer. For key-based signing, pass `--key cosign.key` to `sign` and
`--key cosign.pub` to `verify`.

To require trusted signatures when loading kits, configure the trusted signer
policy, then turn on the requirement:

```console
$ sbx settings set kit.trustedSigners \
    '[{"identity":"release-bot@example.com","issuer":"https://accounts.google.com"}]'
$ sbx settings set kit.requireSignature true
```

V3 source directories and Git sources don't support the source-signing
workflow. Publish and sign an OCI image when signatures are required.

## Published format

The image you publish contains both the kit's content and its descriptor.
Docker image tools can inspect and distribute it, and Docker Sandboxes reads
the descriptor when creating the sandbox.

The built kit also includes its descriptor under
`/usr/share/runtime/kit/<stem>/`, so you can inspect it inside the sandbox.
For the image annotations and file layout, see
[Published image format](kit-reference.md#published-image-format).

Running a workload with `docker run` uses its image configuration, but doesn't
apply the kit's capability declarations or lifecycle hooks. Use `sbx` to run
it with those behaviors.
