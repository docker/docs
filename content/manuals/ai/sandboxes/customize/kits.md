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

Every v3 sandbox composition contains exactly one workload kit and zero or more
mixin kits:

| Kind | What it supplies | How you use it |
| --- | --- | --- |
| `workload` | A complete root filesystem and the command to run, such as an agent or a shell | Pass it as the first argument to `sbx run` or `sbx create` |
| `mixin` | Additional files, tools, or runtime declarations | Add it with `--kit` |

Both kinds use the same descriptor schema. A mixin that only declares network
access or agent instructions needs no Dockerfile. A mixin that ships a tool
includes a build recipe. A workload always includes a build recipe because it
supplies the sandbox's filesystem and launch configuration.

For example, with a local v3 workload and two local v3 mixins:

```console
$ sbx run ./my-agent --name my-project --kit ./my-tool --kit ./team-config .
```

`sbx` builds the source directories, checks that the kits can work together,
and assembles their content into the image the sandbox runs. Unchanged source
builds and compositions reuse cached results. For complete examples, see
[Kit examples](kit-examples.md). To build an agent workload, follow
[Build an agent](build-an-agent.md).

`--kit` applies when creating a sandbox. An existing sandbox with the same
name is reused, so use a different `--name` to try a different kit set or
source revision. Recreate the sandbox when you want to replace its composition.

## Capabilities

Kits declare two kinds of contract. The distinction determines who supplies
what a kit needs.

### Requests to the runtime

The `capabilities` list describes what the kit needs the runtime to do. Each
entry has a namespaced, versioned `type` and a `config` specific to that type.
For example, this mixin requests outbound access to the GitHub API:

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

Other capability types request credentials, volumes, ports, lifecycle hooks,
or agent instructions. The type's `@1` identifies its config schema. A
capability can evolve independently of the kit descriptor's `schemaVersion`.

The capability contract distinguishes required requests from requests marked
`optional: true`. Mark a capability optional only when the kit can work without
it. Optional requests still participate in permission review when granted.

Docker Sandboxes implements a subset of the capability schemas and doesn't
reject every unsupported required request. Check the support notes in the
reference before relying on a capability.

For each type's fields and Docker Sandboxes support, see
[Runtime capabilities](kit-reference.md#runtime-capabilities).

### Contracts between kits

Use `provides`, `requires`, `integrates`, and `conflicts` to describe how kits
fit together:

```yaml
provides: ["my-tool@1.0.0"]
requires: ["node >= 22.0.0"]
integrates: ["docker-engine >= 25.0.0"]
conflicts: ["incompatible-tool"]
```

This kit supplies `my-tool` version `1.0.0`, needs another selected kit to
provide Node.js version `22.0.0` or later, and can work with Docker Engine if
it is present at a compatible version. It rejects a composition that supplies
`incompatible-tool`.

These names are declarations by kit authors. They don't install packages or
search a registry. You select the complete set of kits, and the runtime
validates it. Installing Node.js in a Dockerfile doesn't automatically declare
`provides: ["node@22.0.0"]`.

See [Composition fields](kit-reference.md#composition-fields) for naming and
version rules.

## Build content and runtime setup

Choose where to put a customization based on when its inputs are available
and whether its result belongs in the reusable image or an individual sandbox.

| Mechanism | When it runs | Use it for |
| --- | --- | --- |
| Dockerfile `RUN` and `COPY` | When the kit is built | Install tools, compile binaries, and package static content |
| Lifecycle `install` hooks | Once per sandbox, during creation | Initialize sandbox state using credentials, mounts, or other runtime inputs |
| Lifecycle `startup` hooks | Each sandbox start | Start a service or refresh state that must be restored after a restart |
| Lifecycle `files` | During sandbox setup | Write configuration with kit argument values or preserve an existing file with `overwrite: false` |

Building a kit is separate from running its lifecycle hooks. A package
installed by a Dockerfile becomes reusable image content. An install hook
runs again for each sandbox you create from the kit. Put tool installations
and compilation in the build when they don't require sandbox-specific inputs.

### Build a workload

The descriptor declares `kind: workload`. Its Dockerfile sets the launch
configuration with `ENTRYPOINT`, `CMD`, `ENV`, `USER`, and `WORKDIR`. For
example, a shell workload can use this companion pair:

```yaml {title="my-shell/my-shell.yaml"}
# syntax=docker/runtime-kit:3
schemaVersion: "3"
kind: workload
provides: ["team-shell@1.0.0"]
```

```dockerfile {title="my-shell/my-shell.dockerfile"}
FROM docker/sandbox-templates:shell
USER root
RUN apt-get update && apt-get install -y jq \
    && rm -rf /var/lib/apt/lists/*
USER agent
ENTRYPOINT ["bash"]
CMD []
```

The sandbox templates provide the runtime's base requirements, including
Bash, Git, a CA certificate store, and the `agent` user with UID 1000. See
[Base image requirements](kit-reference.md#base-image-requirements) before
choosing another base.

### Build a mixin

A mixin contributes an overlay: files added or changed by its recipe. The
recipe's base image is a build environment, and its unchanged filesystem
doesn't become part of the overlay. For a tool built in a separate stage,
use a final `FROM scratch` stage and copy the tool and everything it needs
into that stage.

A copied binary must be compatible with the workload's architecture and
libraries. Declare relevant compatibility requirements with `requires` and
ship dependencies that the workload doesn't supply. See
[Build a tool overlay](kit-examples.md#build-a-tool-overlay).

A mixin's environment additions become part of the composed image, but its
`ENTRYPOINT`, `CMD`, `USER`, and `WORKDIR` don't replace the workload's launch
configuration.

### Static files

Static content belongs in the image. Use Dockerfile `COPY` to put a tool config,
helper script, or reference document at its destination, or stage it at a
kit-specific path for a lifecycle hook to copy later.

Files destined for a mounted workspace or persistent volume need that second
step: a mount can hide files baked into the image at its mount path. Copy from
the staged image path after the mount is available. V3 doesn't automatically
inject a source directory named `files/home/` or `files/workspace/`. See
[Copy shared configuration](kit-examples.md#copy-shared-configuration).

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
String commands run through
`sh -c`; an argument list invokes the command directly. See
[Lifecycle](kit-reference.md#lifecycle) for command fields, file permissions,
and runtime behavior.

## Control network access

Use `com.docker.runtime/network-policy@1` to declare network rules. The schema
separates `install` rules, intended for runtime install hooks, from `runtime`
rules for the workload. Neither block configures the Dockerfile build's network.

```yaml
capabilities:
  - type: com.docker.runtime/network-policy@1
    config:
      runtime:
        allow: [api.example.com]
        deny: [telemetry.example.com]
```

Network declarations compose across the selected kits. A deny rule takes
precedence over a kit allow rule. Kit rules also participate in the sandbox's
[policy precedence](../governance/concepts.md#precedence).

See [Network policy](kit-reference.md#network-policy) for phase support and
pattern syntax. Use `sbx policy log` to investigate refused connections.

## Authenticate to external services

A credential capability declares the service and how the runtime presents its
credential. The user supplies the value on the host with
[`sbx secret set`](../configuration/credentials.md#stored-secrets).

```yaml
capabilities:
  - type: com.docker.runtime/network-policy@1
    config:
      runtime:
        allow: [api.example.com]
  - type: com.docker.runtime/credential@1
    config:
      service: my-service
      phase: runtime
      apiKey:
        name: MY_SERVICE_TOKEN
        proxyManaged: true
        inject:
          - domain: api.example.com
            header: Authorization
            format: "Bearer %s"
```

With `proxyManaged: true`, the sandbox receives a sentinel value in
`MY_SERVICE_TOKEN`. The host proxy inserts the real secret in outbound requests
matching the injection rule. Every injection domain must also appear in the
same phase's network allow list.

Store the credential using the kit's service identifier:

```console
$ sbx secret set my-service
```

On the first interactive run, `sbx` also asks you to approve the service's
credential mechanism and domains. Storing a secret doesn't grant that approval.
Without a binding, the sandbox starts with the credential withheld. For
unattended runs, prepare the binding in advance.

See [Credentials](kit-reference.md#credentials) for API key and OAuth fields,
and [Credential configuration](../configuration/credentials.md) for host-side
storage and approval.

## Compose kits

Composition validates the set you selected and orders providers before the
kits that require or integrate with them. Independent mixins are ordered by
reference. Reordering `--kit` flags isn't an override mechanism.

The workload supplies the root filesystem and launch configuration. Mixins
contribute overlays and runtime declarations. Files contributed by more than
one kit cause a composition error. Conflicting image environment values also
cause an error, while `PATH` additions are combined. Give each kit its own
paths for staged content and avoid having multiple kits manage the same config.

V3 has no `extends` field. To derive a workload, build its Dockerfile from the
base image you want and declare its runtime capabilities in its descriptor.
Dockerfile `FROM` inherits image content and config, but it doesn't merge the
parent kit's descriptor into the child. To extend a workload without replacing
it, use a mixin.

Select all kits when creating the sandbox. To change a v3 kit set, recreate the
sandbox with the desired workload and mixins. `sbx kit add` doesn't apply v3
changes to an existing sandbox.

## Pass arguments to kits

Kit arguments use `${{ kit.args.<name> }}` in the descriptor. They resolve in
one of two phases:

| Declaration | Phase | Supply a value with |
| --- | --- | --- |
| `buildArg: VERSION` | Kit build | `docker buildx build --build-arg version=...` using the kit argument's name |
| `env: TOOL_MODE`, or neither mapping | Sandbox creation | `sbx run --kit-arg mode=...` |

For example:

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

`env` exports the resolved value to the sandbox. Without `env`, a create-time
argument is available only through descriptor substitution. Build arguments
are validated and expanded before the kit is published. Changing a build
argument requires rebuilding the kit.

A bare argument name applies to every kit that declares it. To target one kit,
prefix the name with its handle and a period:

```console
$ sbx run ./my-agent --kit ./my-tool --kit-arg my-tool.mode=fix .
```

The handle is the local directory name, the Git subdirectory or repository
name, or the last repository segment in an OCI reference. Scoped values take
precedence over shared values. Use `--kit-args-file <FILE>` for reusable
`name=value` entries; `--kit-arg` values take precedence over file values.

Argument values are plain text and can be recorded in shell history and
sandbox state. Use credential capabilities for secrets. See
[Arguments](kit-reference.md#arguments) for validation and mapping fields.

## Directory and build layout

A companion pair keeps the YAML descriptor separate from its Dockerfile recipe:

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

For a single-file kit, use `build: |` with literal Dockerfile text in the
YAML. You can also select a differently named recipe with `dockerfile:` or
embed a descriptor in a Dockerfile comment block. See
[Authoring forms](kit-reference.md#authoring-forms) for syntax and discovery
rules. The `files/` name in this example is an authoring convention, not a
special runtime directory.

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

### Sign and verify a published kit

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

The image layers carry the workload filesystem or mixin overlay. The image
config carries environment and launch settings. The manifest annotation
`vnd.docker.runtime.kit.descriptor` carries the published descriptor as compact
JSON. Build-time argument values are resolved in this descriptor; create-time
values are resolved for each sandbox.

Every kit also includes its published descriptor at
`/usr/share/runtime/kit/<stem>/kit.yaml` and, when it has a recipe, that recipe
at `kit.dockerfile` in the same directory. This makes the source declarations
available for inspection inside the sandbox. Agent guidance supplied through
`contentFile` is staged in the image and referenced by its published path.

Ordinary Docker image tools can inspect and distribute the image. Running a
workload with `docker run` uses its image configuration, but doesn't apply the
kit's capability declarations or lifecycle hooks. Use `sbx` for those runtime
behaviors. A mixin is intended to be composed with a workload.

See the [Kit spec reference](kit-reference.md) for the descriptor schema.
