---
title: Kit spec reference
linkTitle: Spec reference
description: Reference for the v3 kit descriptor, including capabilities, composition, arguments, build recipes, lifecycle hooks, and the published image format.
keywords: sandboxes, sbx, kits, v3, schema, capabilities, workload, mixin
weight: 50
---

{{< summary-bar feature_name="Docker Sandboxes sbx" >}}

> [!NOTE]
> Kits are experimental. The kit file format, CLI commands, and experience
> for creating, loading, and managing kits are subject to change. Share
> feedback in [docker/sbx-releases](https://github.com/docker/sbx-releases).

This page describes the v3 kit descriptor and its capability configs. Use it
when authoring a workload or mixin. For concepts and usage, see [Kits](kits.md).
For complete examples, see [Kit examples](kit-examples.md).

Some capability types describe functionality beyond the `sbx` integration.
The [capability table](#runtime-capabilities) identifies these types.

## Schema versions

Use `schemaVersion: "3"` for kit authoring. Docker Sandboxes also supports v1
and v2 kits. For `schemaVersion: "2"`, see the
[v2 spec reference](kits-v2/kit-reference.md), including the
[v1-to-v2 field mapping](kits-v2/kit-reference.md#schema-versions).
Each descriptor uses one grammar. V2 fields aren't accepted in a v3
descriptor.

## Descriptor fields

A descriptor is a YAML document with `schemaVersion: "3"` and a `kind`.
The `# syntax` line selects the kit BuildKit frontend when you build it.

```yaml
# syntax=docker/runtime-kit:3
schemaVersion: "3"
kind: mixin
displayName: Team guidelines
description: Development conventions for the team
version: "1.0.0"
provides: [team-guidelines]
capabilities:
  - type: com.docker.runtime/agent-context@1
    config:
      content: |
        Run the project's tests before committing changes.
```

Only `schemaVersion` and `kind` are required at the top level. A workload also
needs a [build recipe](#authoring-forms) to supply its filesystem and launch
command. Unknown fields are errors, including unknown keys in the configs of
recognized capability types.

| Field | Type | Description |
| --- | --- | --- |
| `schemaVersion` | String | Required. Exactly `"3"`. |
| `kind` | String | Required. `workload` or `mixin`. |
| `displayName` | String | Human-readable label. |
| `description` | String | Short summary. |
| `author` | String | Publisher display text, such as `Name <email>`. Metadata, not verified publisher identity. |
| `sourceUrl` | String | Source repository or documentation URL. |
| `iconUrl` | String | Absolute HTTPS URL to an image for catalogs and pickers. |
| `version` | String | Fallback version for `provides` entries with no version. See [Versions](#versions). |
| `licenses` | List of strings | SPDX license identifiers. |
| `provides` | List of strings | Features this kit supplies. |
| `requires` | List of strings | Features the composition must supply. |
| `integrates` | List of strings | Features this kit integrates with when present. |
| `conflicts` | List of strings | Features that must be absent from the composition. |
| `capabilities` | List of objects | Typed requests for runtime resources and behavior. |
| `args` | Map | Named build-time or create-time arguments. |
| `build` | String | Inline Dockerfile. Mutually exclusive with `dockerfile`. |
| `dockerfile` | String | Path to a companion Dockerfile, relative to the descriptor's directory. Mutually exclusive with `build`. |

Field names are case-sensitive. For example, use `sourceUrl`, not the v2
spelling `sourceURL`. A v3 descriptor has no `name`, `extends`, `sandbox`,
`environment`, or `setup` field. Its consumed reference identifies the kit,
the image config defines how it runs, and capabilities declare runtime
behavior.

## Composition fields

Exactly one `workload` supplies the sandbox's root filesystem and launch
command. Zero or more `mixin` kits contribute filesystem overlays and
declarations. A mixin can contain declarations alone.

The following fields describe relationships between kits. They are distinct
from `capabilities`, which requests something from the runtime.

```yaml
provides: ["com.example/tooling@1.4.0"]
requires: ["node >= 20.0.0"]
integrates: ["docker-engine >= 25.0.0"]
conflicts: [podman]
```

| Field | Entry syntax | Resolution rule |
| --- | --- | --- |
| `provides` | `name` or `name@version` | Advertises a feature. No feature is provided implicitly. |
| `requires` | `name` or `name >= version` | A provider in the selected kit set must satisfy the requirement. |
| `integrates` | `name` or `name >= version` | An absent provider is accepted. A present provider must satisfy the minimum. |
| `conflicts` | `name` | Resolution fails if the selected set provides this name. |

Requirements validate the kits you select. They don't search a registry or
download dependencies. Providers compose before their dependents, including
matched `integrates` entries. The dependency graph determines this order,
rather than the order of `--kit` flags.

Bare names such as `node` normalize to `com.docker.kit/node`. A qualified
name such as `com.example/node` is a separate feature. Names use lowercase
letters, digits, and hyphens, with an optional namespace before `/`.

### Versions

Versions contain dot-separated segments, starting with digits, with no `v`
prefix. Examples include `20`, `20.0.0`, and `1.2.3-rc1`. Segments compare
numerically when numeric, and lexically otherwise. Requirements support
minimum versions only, with no exact pins or ranges.

An explicit version in `provides`, such as `node@20.0.0`, takes precedence.
For a provide with no version, a version-shaped consumption tag takes precedence
over the descriptor's `version` fallback. At publication, each `provides`
entry must have an explicit version or a descriptor `version` fallback.

Build-time argument references are accepted in `provides` and `version` and
expanded before publication. Keep `requires`, `integrates`, and `conflicts`
literal.

## Arguments

Declare arguments under `args` and reference them as
`${{ kit.args.<name> }}`. Arguments resolve in one of two phases:

- An argument with `buildArg` resolves while building the kit. Its value is
  passed to the Dockerfile and expanded into the published descriptor.
- An argument without `buildArg` resolves when creating the sandbox. Its
  value is substituted into the descriptor for that installation.

```yaml
args:
  version:
    default: "2.98.0"
    pattern: '^[0-9]+\.[0-9]+\.[0-9]+$'
    buildArg: GH_VERSION
  timeout:
    default: "30000"
    pattern: '^[0-9]+$'
    env: BROWSER_TIMEOUT
  team:
    required: true
    enum: [alpha, beta]
```

| Field | Description |
| --- | --- |
| Argument name | Must match `[A-Za-z_][A-Za-z0-9_]*`. Hyphens aren't accepted. |
| `default` | String used when no value is supplied. `""` is a valid default. Mutually exclusive with `required: true`. |
| `required` | If `true`, the caller must supply a value. Defaults to `false`. |
| `description` | Help text for the argument. |
| `enum` | List of accepted string values. Mutually exclusive with `pattern`. |
| `pattern` | Go RE2 regular expression matched against the whole value. Mutually exclusive with `enum`. |
| `env` | Exports a resolved create-time value under this environment variable name. Mutually exclusive with `buildArg`. |
| `buildArg` | Dockerfile build argument name. Mutually exclusive with `env`. |

`env` and `buildArg` names follow the same identifier rules as argument names.
An optional argument without a default can be omitted only if the descriptor
doesn't reference it. Every argument reference must have a declaration and a
resolved value. Use `default: ""` when an empty value is valid. Missing required
values, unknown supplied arguments, and values outside their constraints are
errors.

Arguments aren't exported as environment variables unless you set `env`.
Use credentials for secrets. Build arguments and published descriptor values
are part of the kit's build and distribution process.

Substitution happens before decoding. Quote placeholders in string fields.
Create-time values produce an effective descriptor that is validated again,
including capability configs. The published descriptor remains unchanged.
Shell expressions such as `$HOME` and `${HOME}` aren't kit argument references.

See [Pass arguments to kits](kits.md#pass-arguments-to-kits) for CLI syntax.

## Authoring forms

The content recipe is a Dockerfile. It builds the files that ship with the
kit and, for a workload, sets `ENTRYPOINT`, `CMD`, `ENV`, `USER`, and
`WORKDIR`. Choose one authoring form:

| Form | Descriptor and recipe |
| --- | --- |
| Companion files | `<stem>.yaml` and `<stem>.dockerfile` in the same directory. Use `dockerfile: <path>` to name the recipe explicitly. |
| Inline recipe | A YAML descriptor with Dockerfile text in `build: \|`. |
| Comment descriptor | A Dockerfile with a `# kit:` comment block containing the descriptor. The remaining file is its recipe. |

An explicitly named Dockerfile must exist and stay within the descriptor's
directory. A missing conventional companion is accepted for a mixin with no
content recipe. A workload must have a recipe.

For the inline form, `build` contains Dockerfile text, not a map of build
options:

```yaml
# syntax=docker/runtime-kit:3
schemaVersion: "3"
kind: mixin
build: |
  FROM scratch
  COPY review-checklist.md /usr/local/share/team/review-checklist.md
```

The comment form starts with a descriptor comment block:

```dockerfile
# syntax=docker/runtime-kit:3
# kit:
#   schemaVersion: "3"
#   kind: mixin
FROM scratch
COPY review-checklist.md /usr/local/share/team/review-checklist.md
```

A comment descriptor can't also declare `build` or `dockerfile`.
Dockerfile semantics apply to every recipe, including multi-stage builds and
build mounts. See [Directory and build layout](kits.md#directory-and-build-layout)
for organizing source files, and
[Packaging and distribution](kits.md#packaging-and-distribution) for build commands.

### Base image requirements

A workload must provide `bash`, `sh`, `curl`, `git`, a populated CA
certificate store, and a non-root `agent` user with UID 1000 and home
directory `/home/agent`. Its image must define an `ENTRYPOINT` or `CMD`.
Install or ship any additional tools the workload needs.

The workload's image config owns the launch command and working directory.
Mixin recipes add files and additive image settings such as environment
variables. A mixin's `ENTRYPOINT`, `CMD`, `USER`, and `WORKDIR` don't replace
the workload's launch contract during composition.

## Runtime capabilities

Each `capabilities` entry names a runtime contract and its config-schema
version. The type version is independent of `schemaVersion`.

```yaml
capabilities:
  - type: com.docker.runtime/port@1
    description: Development server
    optional: false
    config:
      container: 3000
```

| Field | Description |
| --- | --- |
| `type` | Required. `<namespace>/<name>@<version>`, such as `com.docker.runtime/port@1`. Names use lowercase letters, digits, and hyphens. Namespaces can also contain dots. The version is a positive integer. |
| `optional` | Defaults to `false`. The spec requires a runtime to refuse unavailable required capabilities and to skip unavailable optional ones. |
| `description` | Human-readable explanation of the request. |
| `config` | Type-specific fields. Omit for capability types with no config. |

The following table lists the standard types. Every type has the prefix
`com.docker.runtime/`.

| Type | Entries per descriptor | `sbx` support |
| --- | --- | --- |
| [`network-policy@1`](#network-policy) | At most one | Install and runtime network policy. |
| [`credential@1`](#credentials) | One per service and phase | Credential bindings and proxy injection. See field limitations. |
| [`lifecycle@1`](#lifecycle) | At most one | Install, startup, generated files, and interactive arguments. See field limitations. |
| [`agent-context@1`](#agent-context) | At most one | Agent profile and kit instructions. |
| [`port@1`](#ports) | One per container port and transport | Publishes ports to the host. |
| [`resources@1`](#resources) | At most one | Workload CPU and memory settings. |
| [`kit-registry@1`](#kit-registry) | At most one | Restricted to approved builder kits. |
| [`volume@1`](#volumes) | One per path | Schema accepted. Runtime integration pending. |
| [`usb-device@1`](#usb-devices) | Multiple distinct requests | Schema accepted. Runtime integration pending. |
| [`privileged@1`](#privileged-mode) | At most one | Schema accepted. Runtime integration pending. |
| [`agent-sessions@1`](#agent-sessions) | At most one | Schema accepted. Runtime integration pending. |

Exact duplicate entries are errors. Unknown capability types can carry a
custom config, but acceptance by the parser doesn't establish runtime support.
`sbx` doesn't implement general rejection of unsupported required capability
types. Don't rely on an unsupported request being enforced or blocking
sandbox creation.

The following sections describe each type's `config` fields.

## Network policy

`com.docker.runtime/network-policy@1` declares outbound access separately for
install hooks and the running workload. These phases concern sandbox
execution. They don't control Dockerfile build networking.

```yaml
capabilities:
  - type: com.docker.runtime/network-policy@1
    config:
      install:
        allow: [registry.npmjs.org]
      runtime:
        allow: [api.github.com, "*.example.com"]
        deny: [telemetry.example.com]
```

| Field | Description |
| --- | --- |
| `install.allow` | Domains permitted while install hooks run. |
| `install.deny` | Domains blocked while install hooks run. |
| `runtime.allow` | Domains permitted during workload execution. |
| `runtime.deny` | Domains blocked during workload execution. |

`sbx` makes runtime rules available during installation too, then removes the
install-only rules before launching the workload. An omitted phase adds no
grants for that phase. Rules from composed kits combine, and deny takes
precedence over allow.

Use exact hosts, hosts with ports such as `api.example.com:443`, or
single-label wildcards such as `*.example.com`. A bare `*` or `**` permits
all destinations. Quote wildcard strings in YAML.

Every credential injection domain must also appear in the matching phase's
allow list. Validation ignores port suffixes for this membership check. A
bare `*` or `**` covers all injection domains. A narrower wildcard doesn't
replace an explicit injection-domain entry for validation.

## Credentials

`com.docker.runtime/credential@1` declares a service and how its credential
is presented. The user stores the value in the secret store and approves its
use through [credential bindings](../configuration/credentials.md#credential-bindings).
The descriptor doesn't name a host file or environment variable to read.

```yaml
capabilities:
  - type: com.docker.runtime/network-policy@1
    config:
      runtime:
        allow: [api.github.com]
  - type: com.docker.runtime/credential@1
    optional: true
    description: GitHub API access
    config:
      service: github
      phase: runtime
      apiKey:
        name: GH_TOKEN
        proxyManaged: true
        inject:
          - domain: api.github.com
            header: Authorization
            format: "Bearer %s"
```

| Field | Description |
| --- | --- |
| `service` | Required. Lowercase service identifier in the host credential store. |
| `phase` | Required. `install` or `runtime`. |
| `apiKey` | API key presentation. Declare `apiKey`, `oauth`, or both. |
| `oauth` | OAuth token interception and credential-file config. |

Set `optional` on the capability entry, outside `config`. In `sbx`, a missing
binding withholds the credential instead of blocking sandbox creation, even
for a required request. Install-only credential injection is removed after
install hooks. When the same service is declared for both phases, `sbx` uses
the runtime entry's presentation config.

### API keys

| Field | Description |
| --- | --- |
| `apiKey.name` | Required. Environment variable name in the sandbox, such as `GH_TOKEN`. |
| `apiKey.proxyManaged` | If `true`, keeps the real credential on the host and puts a sentinel value in the sandbox. Defaults to `false`. |
| `apiKey.inject[].domain` | Required for each rule. Domain where the proxy injects the credential. Must be allowed in the same phase's network policy. |
| `apiKey.inject[].header` | HTTP header to set, such as `Authorization`. |
| `apiKey.inject[].format` | Header value format, such as `"Bearer %s"`. |
| `apiKey.inject[].scheme` | Schema field for an auth scheme, such as `basic`. Runtime mapping in `sbx` is pending. Use explicit `header` and `format` for header injection. |
| `apiKey.inject[].username` | Username for an authentication scheme that requires one. |

### OAuth

| Field | Description |
| --- | --- |
| `oauth.tokenEndpoint.host` | Required when `tokenEndpoint` is set. Host whose OAuth token response the proxy intercepts. |
| `oauth.tokenEndpoint.path` | Token endpoint path. |
| `oauth.resourceHosts` | API hosts where the proxy uses the token. |
| `oauth.sentinels.accessToken` | Placeholder access token presented inside the sandbox. |
| `oauth.sentinels.refreshToken` | Placeholder refresh token presented inside the sandbox. |
| `oauth.credentialFile.path` | Credential-file destination. `~` expands to the agent's home. |
| `oauth.credentialFile.structure` | Nested map rendered as JSON after placeholder substitution. |
| `oauth.responseFields.accessToken` | Provider field to read instead of `access_token`. |
| `oauth.responseFields.expiresIn` | Provider field to read instead of `expires_in`. |
| `oauth.passthrough` | If `true`, returns real tokens to the sandbox instead of sentinels. Defaults to `false`. |

Credential-file leaf values support `{{.AccessToken}}`,
`{{.RefreshToken}}`, `{{.ExpiresAt}}`, `{{.Scopes}}`, and
`{{.PrimaryApiKey}}`. The last placeholder's containing key is omitted when
no primary API key is captured. Unknown placeholders are errors.

The v2 fields `provider`, `skipIfEnv`, and `credentialFile.template` aren't
part of the v3 schema.

## Lifecycle

`com.docker.runtime/lifecycle@1` declares work performed inside the sandbox
after the kit image has been built. Use a Dockerfile to install software that
can ship in the image. Use lifecycle hooks for initialization or work that
needs the sandbox's runtime state.

```yaml
capabilities:
  - type: com.docker.runtime/lifecycle@1
    config:
      install:
        - command: [sh, -c, 'printf "%s\n" "$WORKSPACE_DIR" > /home/agent/workspace-path']
          user: "1000"
          env: [WORKSPACE_DIR]
      startup:
        - command: [my-service, --foreground]
          background: true
      files:
        - path: /home/agent/.config/tool/settings.json
          content: '{"telemetry": false}'
          mode: "0644"
          overwrite: false
      interactive: [--interactive]
```

### Install and startup hooks

| Field | Description |
| --- | --- |
| `install[].command` | Required. Shell string passed to `sh -c`, or a list of command arguments. Runs synchronously when the sandbox is created. |
| `install[].user` | Execution user. Defaults to root. Use a username or quoted UID. |
| `install[].env` | Additional environment variable names the hook can read. |
| `install[].description` | Human-readable explanation. |
| `startup[].command` | Required. Shell string or list of command arguments. Runs on each sandbox start. |
| `startup[].user` | Execution user. Defaults to the agent user. Use `"0"`, `"1000"`, or a login name. |
| `startup[].background` | If `true`, lets the startup dispatcher continue without waiting for this command. Defaults to `false`. |
| `startup[].env` | Schema field for the hook's environment inputs. Filtering in `sbx` is pending. |
| `startup[].description` | Human-readable explanation. |

Install hooks retain basic process variables, proxy settings, and certificate
paths. Declare other inputs in `env`, including `WORKSPACE_DIR` and any
credential sentinel variables the hook reads. Hooks from dependency providers
run before hooks from their dependents.

Startup hooks must tolerate running again. They run without an interactive
terminal. In `sbx`, the startup dispatcher runs alongside the workload:
`background: false` orders startup commands but doesn't delay the workload's
entrypoint. Use install hooks or generated files for initialization the
workload must see before it launches.

### Generated files

| Field | Description |
| --- | --- |
| `files[].path` | Required. Absolute path inside the sandbox. |
| `files[].content` | File body. Kit argument references expand at create time. |
| `files[].mode` | Octal permissions as a string. Set explicitly, such as `"0644"`, for consistent file permissions. |
| `files[].overwrite` | Defaults to `true`. Set to `false` to retain an existing file. |
| `files[].description` | Human-readable explanation. |

Files are written after install hooks, as the agent user, before the workload
runs. An install hook can read files from the image but can't depend on
generated files. Choose a path the agent can write. For root-owned locations,
create the file in the image or use an install hook. For static content, use
Dockerfile `COPY` instead.

### Interactive arguments

`interactive` is a list of arguments for TTY sessions. For a workload with
image `ENTRYPOINT` and `CMD`, the launch forms are:

| Mode | Command |
| --- | --- |
| Default | Image `ENTRYPOINT` followed by image `CMD`. |
| Interactive | Image `ENTRYPOINT` followed by lifecycle `interactive`. |

An omitted or empty `interactive` list in `sbx` falls back to the default
command. A lifecycle capability must declare at least one hook, file, or
nonempty interactive argument list.

## Agent context

`com.docker.runtime/agent-context@1` supplies instructions the agent reads.
The workload chooses the profile filename. Mixins contribute instructions to
that profile's kit index.

```yaml
capabilities:
  - type: com.docker.runtime/agent-context@1
    config:
      filename: AGENTS.md
      contentFile: ./agent-context.md
```

| Field | Description |
| --- | --- |
| `filename` | Agent profile filename, such as `AGENTS.md` or `CLAUDE.md`. Accepted only on workload kits. |
| `contentFile` | Path to an instruction file relative to the build context. Mutually exclusive with `content`. |
| `content` | Inline instruction text. Mutually exclusive with `contentFile`. |

The frontend copies a `contentFile` into the kit image and rewrites the
published descriptor to point to that location. The agent profile indexes
each kit's instructions so the agent can read them when needed. A mixin omits
`filename` and uses `contentFile` or `content`.

## Ports

`com.docker.runtime/port@1` publishes an inbound port. Declare a separate
entry for each port and transport.

| Field | Description |
| --- | --- |
| `container` | Required. Container port from 1 to 65535. |
| `transport` | `tcp` or `udp`. TCP is the default. |
| `name` | Optional label for port listings. |

The runtime allocates the host port. To choose a host port, use
`sbx ports <SANDBOX> --publish <HOST_PORT>:<CONTAINER_PORT>` instead of a descriptor
field. Port publication doesn't grant outbound network access.

## Resources

`com.docker.runtime/resources@1` declares workload resource settings.

| Field | Description |
| --- | --- |
| `cpu` | Non-negative number of CPU cores in the schema. `sbx` requires a whole number. |
| `memory` | Byte-size string, such as `4096m` or `8g`. |
| `gpu` | GPU selector interpreted by the runtime. Enforcement in `sbx` is pending. |

Unset fields add no resource constraint. `sbx` applies CPU and memory from the
workload kit unless overridden at creation. Resource declarations on mixins
aren't applied.

## Kit registry

`com.docker.runtime/kit-registry@1` requests access to the runtime's registry
for kit builds. It takes no `config`.

In `sbx`, this capability is restricted to approved OCI builder kits, with
`docker/sbx-kit-builder` approved by default. Local and Git kit sources don't
receive this grant. Declaring it doesn't grant general network access.

## Volumes

`com.docker.runtime/volume@1` describes a storage mount. `sbx` accepts this
schema, but doesn't apply v3 volume requests to the sandbox.

| Field | Description |
| --- | --- |
| `path` | Required. Absolute mount path inside the sandbox. |
| `size` | Optional byte-size string, such as `2g`. |
| `tmpfs` | If `true`, requests RAM-backed storage. Defaults to `false`. |
| `mode` | Optional octal permissions, such as `"0755"`. |

## USB devices

`com.docker.runtime/usb-device@1` describes a USB device request. `sbx`
accepts this schema, but doesn't apply v3 USB requests.

| Field | Description |
| --- | --- |
| `vendorId` | Vendor identifier. Must be paired with `productId`. |
| `productId` | Product identifier. Must be paired with `vendorId`. |
| `class` | Device class. Mutually exclusive with the vendor and product pair. |

Declare exactly one match form: a vendor/product pair, or a class.

## Privileged mode

`com.docker.runtime/privileged@1` requests elevated runtime privileges. It
takes no `config`. `sbx` accepts this schema, but doesn't apply v3 privileged
requests.

## Agent sessions

`com.docker.runtime/agent-sessions@1` describes a workload's session commands.
`sbx` accepts this schema, but doesn't use it to drive sessions.

| Field | Description |
| --- | --- |
| `prompt` | Argument list appended to the workload's launch command. Must contain `{{.Prompt}}`. |
| `resume` | Argument list appended to the launch command. Must contain `{{.SessionID}}`. |
| `continue` | Argument list to reopen the most recent session. |
| `list` | Complete command, as a shell string or argument list, that outputs session IDs, one per line, most recent first. |

At least one command is required. `prompt`, `resume`, and `continue` are
argument tails. `list` is a complete command.

## Published image format

A published v3 kit is an OCI image. Its manifest annotations carry the
descriptor, its image config carries the launch settings, and its layers
carry content.

| Annotation | Value |
| --- | --- |
| `vnd.docker.runtime.kit.descriptor` | Published descriptor as compact JSON. |
| `vnd.docker.runtime.kit.schema-version` | `"3"`. |
| `vnd.docker.runtime.kit.capabilities` | Sorted, comma-separated capability types with duplicates removed. Omitted when none are requested. |

Each image stages its descriptor at
`/usr/share/runtime/kit/<stem>/kit.yaml` and its recipe, when present, at
`/usr/share/runtime/kit/<stem>/kit.dockerfile`. Instruction files referenced
by `contentFile` are staged under the same kit directory. Even a mixin with
no recipe has a layer containing its descriptor.

The published descriptor has a 512 KiB limit, with a build warning above
64 KiB. Keep substantial instruction text in `contentFile` and other content
in image layers. See [Compose kits](kits.md#compose-kits) for combining images,
and [Packaging and distribution](kits.md#packaging-and-distribution) for publishing them.

## Move from v2 to v3

A v3 kit combines image content with a descriptor. Changing `schemaVersion`
alone doesn't convert a v2 kit. Separate reusable build work from runtime
initialization, then declare the runtime capabilities the kit needs.

| V2 surface | V3 equivalent |
| --- | --- |
| `kind: sandbox` | `kind: workload` with a Dockerfile recipe |
| `sandbox.image` | Dockerfile `FROM` |
| `sandbox.entrypoint`, `sandbox.command`, `environment.variables` | Dockerfile `ENTRYPOINT`, `CMD`, and `ENV` |
| `extends` | A mixin for composition, or a derived workload image with its own descriptor |
| `setup.install` | Dockerfile `RUN` for reusable content; lifecycle `install` for sandbox initialization |
| `setup.startup` and `setup.files` | Lifecycle capability `startup` and `files` |
| `setup.files[].onlyIfMissing: true` | Lifecycle `files[].overwrite: false` |
| Automatic `files/home/` and `files/workspace/` injection | Dockerfile `COPY`, with lifecycle hooks for destinations provided by runtime mounts |
| `permissions.network` and `credentials` | Network-policy and credential capabilities |
| `agentInstructions` | Agent-context capability |

Review the [runtime support table](#runtime-capabilities) before migrating
features such as volumes. Publish the converted kit as an image, and select
v3 workload and mixin kits together.
