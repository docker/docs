---
title: Use kits
description: Run v3 workload kits, add compatible mixins, supply arguments, and manage kit sources and runtime access.
keywords: sandboxes, sbx, kits, v3, workloads, mixins
weight: 10
---

{{< summary-bar feature_name="Docker Sandboxes sbx" >}}

You can use a kit someone else has published, or build your own. Kit authors
describe the environment in a YAML file and package any software and files it
needs into an image. To use it, give `sbx` the kit's reference. Docker
Sandboxes prepares the environment and applies its settings.

> [!NOTE]
> V3 kits are experimental. Select a v3 workload and v3 mixins together.
> Built-in shortcuts such as `claude` and `codex` use v2 and can't be combined
> with v3 mixins. See [Version compatibility](/manuals/ai/sandboxes/customize/_index.md#version-compatibility)
> or the [v2 reference](/manuals/ai/sandboxes/customize/kits-v2/_index.md).

## Run a kit

Run a workload by passing its kit reference to `sbx run`:

```console
$ sbx run docker.io/my-org/agent-kit:1.0.0
```

This example uses an image reference. Kits can also come from a local
directory or Git repository; see [Choose a kit source](#choose-a-kit-source).
Replace the example references with the v3 kits you want to use.

The sandbox uses your current directory as its workspace. To use another
project directory, append its path to the command. See
[Choose a workspace](/manuals/ai/sandboxes/usage.md#choose-a-workspace).

To create a sandbox without launching the workload, use `sbx create`.
Include `.` as the workspace argument to mount your current directory;
omitting the workspace creates a mountless sandbox.

If you don't have a v3 workload to run, the
[OpenCode workload example](/manuals/ai/sandboxes/customize/author/_index.md#build-a-workload) provides a complete source kit.
For a step-by-step authoring walkthrough, see [Build an agent](/manuals/ai/sandboxes/customize/author/build-an-agent.md).

## Add mixins

A sandbox runs one workload kit. Add v3 mixins with `--kit`, repeating the
flag for each one. For example, add a Neovim mixin and shared team configuration:

```console
$ sbx run docker.io/my-org/agent-kit:1.0.0 \
    --kit docker.io/my-org/neovim-kit:1.0.0 \
    --kit docker.io/my-org/team-config-kit:1.0.0
```

The workload and its mixins form a composition: their tools, files, and
runtime settings combine to define the sandbox's environment. See
[Compose kits](#compose-kits) for dependency and compatibility rules, and
[Kit examples](/manuals/ai/sandboxes/customize/author/kit-examples.md) for complete mixins.

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

## Name and reuse a sandbox

Use `--name` to give the sandbox a name:

```console
$ sbx run docker.io/my-org/agent-kit:1.0.0 --name my-project
```

Running an existing sandbox reuses its recorded kit configuration. Kit
selection applies when creating a sandbox. To use a different v3 workload or
mixin set, choose another name or recreate the sandbox with the desired kits.

## Compose kits

Choose a workload and mixins that are designed to work together. A kit can
require functionality provided by another selected kit or reject an incompatible
combination. Docker Sandboxes checks these relationships when creating the
sandbox; it doesn't find or download missing dependencies for you.

For example, if a mixin requires a tool provided by another kit, include both
kits with `--kit`. Installing a similarly named tool in the workload doesn't
satisfy that requirement unless its kit declares that it provides it.

Providers are applied before kits that require them. Independent mixins are
ordered by reference, so changing the order of `--kit` flags doesn't resolve
conflicts. Two kits contributing the same image file or conflicting environment
values cause a composition error; `PATH` additions are combined. Choose compatible
kits or ask their authors to resolve overlapping customizations.

The workload controls the launch command and working directory. Mixins add
content and runtime settings without replacing that launch configuration.
Select the whole composition at creation time. `sbx kit add` doesn't apply v3
changes to an existing sandbox.

## Pass arguments to kits

Kits can expose arguments for choices such as a tool's operating mode. Use
`--kit-arg` to supply a value declared by the kit:

```console
$ sbx run docker.io/my-org/agent-kit:1.0.0 \
    --kit docker.io/my-org/linter-kit:1.0.0 --kit-arg mode=fix
```

A bare argument name applies to every kit that declares it. To target one kit,
prefix the name with its handle and a period:

```console
$ sbx run docker.io/my-org/agent-kit:1.0.0 \
    --kit docker.io/my-org/linter-kit:1.0.0 --kit-arg linter-kit.mode=fix
```

The handle is the local directory name, the Git subdirectory or repository
name, or the last repository segment in an OCI reference. Scoped values take
precedence over shared values. Use `--kit-args-file <FILE>` for reusable
`name=value` entries; `--kit-arg` values take precedence over file values.

Use the kit's documentation to find its argument names, defaults, and accepted
values. Arguments resolved when an image is built, such as a bundled tool's
version, require the author to rebuild that image. Creating a sandbox doesn't
change those values in a published kit.

Argument values are plain text and can be recorded in shell history and
sandbox state. Use [stored credentials](/manuals/ai/sandboxes/configuration/credentials.md)
for secrets.

## Runtime access and instructions

Kits declare capabilities: requests for network access, credentials, storage,
and runtime behavior. Docker Sandboxes applies these declarations when creating
and running the sandbox.

Network rules from the selected kits combine. Kit deny rules take precedence
over kit allow rules. When organization governance is active, only organization
allow rules grant access; kit allow rules don't expand it. See
[Policy precedence](/manuals/ai/sandboxes/governance/concepts.md#precedence).

A kit can request credentials for a named service. Store the credential on your
host and approve the kit's use when prompted. With proxy-managed credentials,
the agent receives a placeholder and the host proxy inserts the real value into
requests to the declared service. See
[Credential configuration](/manuals/ai/sandboxes/configuration/credentials.md) for storing keys,
approval, and unattended use.

Kits can also contribute agent instructions. The workload chooses a profile
filename, such as `AGENTS.md` or `CLAUDE.md`. Docker Sandboxes generates it above
the mounted workspace, outside the mount, and includes or references guidance
from the selected kits. It doesn't replace instructions in your project.

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
