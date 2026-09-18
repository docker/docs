---
title: Use kits
description: Run published v3 kit sets and workloads, add compatible mixins, and configure kit arguments and runtime access.
keywords: sandboxes, sbx, kits, v3, workloads, mixins
weight: 10
---

{{< summary-bar feature_name="Docker Sandboxes sbx" >}}

If you've run `sbx run claude` or `sbx run codex`, you've already used a kit.
The built-in agents are kits that package an environment, tools, and runtime
settings. Kits you build yourself or get from another publisher use the same
model: you give `sbx` a reference, and Docker Sandboxes prepares the environment
and applies its settings.

This page shows how to run published kits, combine them with mixins, and
configure their runtime behavior. See
[Kits](/manuals/ai/sandboxes/customize/_index.md) for an introduction to workloads,
mixins, and sets.

## Run a kit

The built-in agent names are shortcuts for kit references. To run another kit,
replace the agent name with that kit's reference. For example, run Docker's
published v3 Claude Code workload:

```console
$ sbx run docker.io/docker/sbx-kit-claude:2.1.274 --name claude-project
```

The sandbox uses your current directory as its workspace. To use another
project directory, append its path. See
[Choose a workspace](/manuals/ai/sandboxes/usage.md#choose-a-workspace).

This starts Claude Code. Follow the prompts to authenticate, or prepare
credentials on the host as described in
[Credential configuration](/manuals/ai/sandboxes/configuration/credentials.md).
A published set containing a workload runs the same way: pass the set's
reference in place of the workload's reference.

Browse [Docker's published kits](https://hub.docker.com/orgs/docker/repositories?search=sbx-kit)
for other environments. For local and Git references, see
[Choose a kit source](#choose-a-kit-source).

To create a sandbox without launching the workload, use `sbx create`.
Include `.` as the workspace argument to mount your current directory;
omitting the workspace creates a mountless sandbox.

## Name and reuse a sandbox

Use `--name` to give a sandbox a name when creating it. Reconnect to the
Claude Code sandbox from the previous example with:

```console
$ sbx run --name claude-project
```

Running an existing sandbox reuses its recorded kit configuration. To try a
different workload, mixin combination, or argument value, create a sandbox
with a different name or recreate the existing one. The following examples
use separate names for separate combinations. `sbx kit add` doesn't change
a v3 composition in place.

## Runtime access and instructions

Kits declare their runtime needs through capabilities. The Claude workload
requests network access and credentials so Claude can reach
Anthropic. When using an API key, storing it on the host provides its value.
Approving the credential request authorizes the kit to use it. See
[Credential bindings](/manuals/ai/sandboxes/configuration/credentials.md#credential-bindings)
for preparing unattended runs.

When choosing a kit, check the publisher's documentation for its services,
credentials, and startup behavior. Proxy-managed credentials stay on your
host, where the sandbox proxy authenticates matching requests on the kit's
behalf.

Network requests must also meet your sandbox's network policy. Kit allow rules
can't grant access beyond your organization's policy. If a connection fails,
check the [policy log](#debug-kits) to see which rule blocked it.

Kits can also supply instructions for using their tools. The workload chooses
the instruction profile, and mixins contribute guidance to it. The generated
profile sits outside your workspace and doesn't overwrite project instructions.

## Add mixins

The examples on this page use v3 kits. Select a v3 workload and v3 mixins
together. Built-in shortcuts such as `claude` and `codex` select v2 kits, so
use [v2 mixins](/manuals/ai/sandboxes/customize/kits-v2/_index.md) with those
shortcuts. See [Version compatibility](/manuals/ai/sandboxes/customize/_index.md#version-compatibility).

Add a tool to the agent environment with `--kit`. For example, an internal
CLI mixin can package your company's executable along with network rules
and a credential request for its API:

```console
$ sbx run docker.io/docker/sbx-kit-claude:2.1.274 --name claude-tools \
    --kit docker.io/<NAMESPACE>/company-cli:1.0.0
```

Replace the mixin reference with one your organization has published. To
build this example, see [Package an internal CLI](/manuals/ai/sandboxes/customize/author/kit-examples.md#package-an-internal-cli).
Store the credential requested by the mixin on the host and approve access
when prompted.

The workload and mixin form a composition: their files and runtime settings
combine into one environment. Claude Code still starts as the agent, with
`company-cli` available for it to use.

Repeat `--kit` to add more mixins. You can also add compatible mixins to a
published workload set, or add a set containing only mixins with `--kit`.
Avoid adding components already included in a set: duplicate feature providers
can cause composition errors.

Once you have a combination to share, [author a kit set](/manuals/ai/sandboxes/customize/author/kit-sets.md)
to publish it as one reference. For reusable component examples, see
[Mixin examples](/manuals/ai/sandboxes/customize/author/kit-examples.md).

## Compose kits

Some mixins need a feature another kit supplies. For example, an Agent Client
Protocol (ACP) adapter lets a compatible editor or other client drive an
agent session. Docker's Claude ACP adapter needs the Claude Code mixin.
Include both with a shell workload:

```console
$ sbx run docker.io/docker/sbx-kit-shell:1.0.0 --name claude-acp-components \
    --kit docker.io/docker/sbx-kit-claude-mixin:2.1.274 \
    --kit docker.io/docker/sbx-kit-claude-acp:0.79.0
```

This opens a shell with `claude` and `claude-agent-acp` installed. An ACP
client communicates with `claude-agent-acp` over standard input and output.
Use the Claude mixin with this adapter: it supplies the binary at the path
the adapter expects.

The adapter declares its dependency in its descriptor:

```yaml
requires: [claude]
```

The Claude Code mixin declares the feature it provides:

```yaml
provides: ["claude@2.1.274"]
```

Docker Sandboxes checks these declarations and applies the provider before
the adapter. Reordering `--kit` flags doesn't change that dependency order.
If you omit the provider, sandbox creation fails. A requirement names a
feature, not an image to download: you select the kit that supplies it.
Requirements can also specify a minimum feature version, and kits can declare
incompatible combinations.

Docker publishes this combination as
`docker.io/docker/sbx-kit-claude-acp-set:2.1.274`. With `--kit`, Docker Sandboxes
checks the combination when creating the sandbox. With a set, the publisher
resolves and merges the components during the build, and consumers receive
one kit with a digest-pinned record of its components.

To declare these relationships in your own kits, see
[Authoring compositions](/manuals/ai/sandboxes/customize/author/_index.md#compose-kits).

## Pass arguments to kits

A kit can expose arguments for choices such as a linter's mode. Check the
publisher's documentation for argument names, defaults, and accepted values.
For example, a set that exposes the `lint_mode` argument shown in
[Configure component arguments](/manuals/ai/sandboxes/customize/author/kit-sets.md#configure-component-arguments)
can select a mode with `--kit-arg`:

```console
$ sbx run docker.io/<NAMESPACE>/claude-tools:1.0.0 \
    --name claude-tools-fix --kit-arg lint_mode=fix
```

This assumes you have published a workload set with that linter component
and argument. Replace the reference with your set's published reference.

A bare argument name applies to every selected kit that declares it. To target
one kit, prefix the name with its handle and a period:

```console
$ sbx run docker.io/<NAMESPACE>/claude-tools:1.0.0 \
    --name claude-tools-fix --kit-arg claude-tools.lint_mode=fix
```

The handle is the local directory name, the Git subdirectory or repository
name, or the last repository segment in an OCI reference. Scoped values take
precedence over shared values. Use `--kit-args-file <FILE>` for reusable
`name=value` entries; `--kit-arg` values take precedence over file values.

A published set exposes only arguments declared on the set. Component
arguments are fixed during publication unless the author exposes them through
set arguments. Build-time choices, such as an installed tool's version,
require rebuilding the image.

Argument values are plain text and can be recorded in shell history and
sandbox state. Use [stored credentials](/manuals/ai/sandboxes/configuration/credentials.md)
for secrets.

## Choose a kit source

Workload and mixin references can point to a published image, a local
directory, or a kit in a Git repository:

| Source | Example reference |
| --- | --- |
| Published image | `docker.io/my-org/agent-kit:1.0.0` |
| Local directory | `./my-agent` |
| Git repository | `git+https://github.com/<ORG>/<REPOSITORY>.git#ref=<COMMIT>&dir=my-agent` |

For Git sources, `ref` selects a revision and `dir` selects the kit's
subdirectory. Quote Git URLs in shell commands because they can contain `&`:

```console
$ sbx run "git+https://github.com/<ORG>/<REPOSITORY>.git#ref=<COMMIT>&dir=my-agent"
```

`sbx` pulls published images and builds local or Git sources when creating
the sandbox. Builds with unchanged source content and supplied kit arguments
reuse cached results. See
[Restrict kit sources](#restrict-kit-sources) for permitted sources and
[Registry credentials](/manuals/ai/sandboxes/configuration/credentials.md#registry-credentials)
for private image access.

## Restrict kit sources

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

## Verify kit signatures

To verify a published kit against its author's signing identity:

```console
$ sbx kit verify docker.io/<NAMESPACE>/my-kit:1.0.0 \
    --certificate-identity <SIGNER_IDENTITY> \
    --certificate-oidc-issuer <ISSUER_URL>
```

For key-based signatures, use `--key cosign.pub` instead of the certificate
identity and issuer options. Obtain the expected identity or public key from
the publisher.

To require trusted signatures when loading kits, configure the trusted signer
policy, then turn on the requirement:

```console
$ sbx settings set kit.trustedSigners \
    '[{"identity":"release-bot@example.com","issuer":"https://accounts.google.com"}]'
$ sbx settings set kit.requireSignature true
```

V3 source directories and Git sources don't support the source-signing
workflow. Use a signed OCI image when signatures are required.

## Debug kits

When a tool is missing or a request fails, inspect the running sandbox:

```console
$ sbx exec <SANDBOX> -- which <TOOL>
$ sbx exec <SANDBOX> -- cat /home/agent/.config/<TOOL>/settings.json
$ sbx policy log
```

The policy log shows outbound requests and the rules they matched. Use it to
find blocked package registries or API hosts. If downloads fail after adding a kit that requests
credentials, ask its author to check that credential injection targets only
the service hosts that need it.

To try an updated v3 kit, create a sandbox with a different name. Reusing
an existing sandbox keeps its recorded kit configuration.
