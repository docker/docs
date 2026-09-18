---
title: Use kits
description: Run published v3 kit sets and workloads, add compatible mixins, and configure kit arguments and runtime access.
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

## Install a v3-capable build

V3 kits require an `sbx` nightly with v3 support. The stable release doesn't
support v3. Get a build from
[sbx releases](https://github.com/docker/sbx-releases/releases).
On macOS, install the preview channel with:

```console
$ brew install docker/tap/sbx@rc
```

## Run a kit

A published kit set gives you a complete environment through one reference.
For example, Docker's Claude ACP set includes a shell workload, Claude Code,
and the Agent Client Protocol (ACP) adapter:

```console
$ sbx run docker.io/docker/sbx-kit-claude-acp-set:2.1.274
```

This opens a shell with `claude` and `claude-agent-acp` installed. Run `claude`
in the sandbox for an interactive session; `claude-agent-acp` serves an ACP
client over standard input and output. The set includes Claude's credential
requests; see [Runtime access and instructions](#runtime-access-and-instructions)
for storing and approving credentials.

For a shell without an agent, run the shell workload on its own:

```console
$ sbx run docker.io/docker/sbx-kit-shell:1.0.0
```

Browse [Docker's published kits](https://hub.docker.com/orgs/docker/repositories?search=sbx-kit)
for other environments and components. Kits can also come from a local
directory or Git repository; see [Choose a kit source](#choose-a-kit-source).

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

A published set containing a workload runs as one workload kit. You can add
compatible mixins to it, or assemble a workload and mixins yourself. Add v3
mixins with `--kit`, repeating the flag for each one. For example, add Neovim for editing and Ruff for Python linting:

```console
$ sbx run docker.io/my-org/agent-kit:1.0.0 \
    --kit docker.io/my-org/neovim-kit:1.0.0 \
    --kit docker.io/my-org/ruff-kit:1.0.0
```

The workload and its mixins form a composition: their tools, files, and
runtime settings combine to define the sandbox's environment. See
[Compose kits](#compose-kits) for dependency and compatibility rules, and
[Kit examples](/manuals/ai/sandboxes/customize/author/kit-examples.md) for complete mixins.

Once you have a combination to share, [publish a kit set](/manuals/ai/sandboxes/customize/author/kit-sets.md)
so consumers can use one reference. A set containing only mixins is added with
`--kit`, like an individual mixin. Avoid adding components already included in
a set: duplicate feature providers can cause composition errors.

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

Composition lets you assemble an environment from kits with different roles.
A set's publisher resolves its components at build time. The consumer receives
one merged kit, with a digest-pinned record of those components. When you add
kits with `--kit`, Docker Sandboxes resolves that combination at sandbox creation.
In the [Neovim and Ruff example](#add-mixins), the workload supplies the agent,
Neovim adds an editor, and Ruff adds a Python linter. The agent runs in one
sandbox with all three available. Adding the tools doesn't change which agent
starts or how it launches.

Some kits are designed to work independently; others need something another
kit supplies. Suppose a browser-testing kit declares that it requires Chromium
from a browser kit. Include both alongside your agent workload:

```console
$ sbx run docker.io/my-org/agent-kit:1.0.0 \
    --kit docker.io/my-org/browser-testing-kit:1.0.0 \
    --kit docker.io/my-org/chromium-kit:1.0.0
```

You can recognize this relationship in the kits' descriptors. The testing kit
might declare:

```yaml
requires: ["chromium >= 120.0.0"]
```

The browser kit declares what it supplies:

```yaml
provides: ["chromium@120.0.0"]
```

The testing kit requires Chromium version 120 or later. The browser kit declares
that it provides version 120, so it satisfies that requirement. These names
describe what kits provide; they aren't image references that `sbx` downloads.

Docker Sandboxes checks that the selected kits satisfy the declared requirement
and applies the browser kit before the testing kit. If you leave out the browser
kit, sandbox creation fails with an unmet requirement. It doesn't download an
extra kit automatically: you choose which provider to include, using the testing
kit's documentation to find a compatible one.

Kits can also declare that they can't work together. For example, two kits
might configure a tool in incompatible ways. Choose a compatible combination.

To declare these relationships in your own kits, see
[Authoring compositions](/manuals/ai/sandboxes/customize/author/_index.md#compose-kits).

A sandbox keeps the composition it was created with. To try another combination,
create a sandbox with a different name or recreate the existing one. `sbx kit add`
doesn't change a v3 composition in place.

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
values. A published set exposes only the arguments its author declares on the
set. Component arguments are fixed during publication unless the author
exposes them through set arguments.

Arguments resolved when an image is built, such as a bundled tool's version, require the author to rebuild that image. Creating a sandbox doesn't
change those values in a published kit.

Argument values are plain text and can be recorded in shell history and
sandbox state. Use [stored credentials](/manuals/ai/sandboxes/configuration/credentials.md)
for secrets.

## Runtime access and instructions

A kit can bring the setup needed to use its tools: access to external services,
authentication, and instructions for the agent. These requests are called
capabilities. When choosing a kit, look at its documentation for the services
it connects to, credentials it needs, and behavior it adds to your sandbox.

For example, suppose you add a GitHub CLI kit so your agent can read issues
and open pull requests. The kit can request access to GitHub, declare how to
authenticate, and give the agent guidance for using `gh`. You don't need to
configure each of those pieces yourself, but you do need to supply a credential
and authorize the kit to use it.

Store the credential on your host using the service name specified by the kit.
For this example:

```console
$ sbx secret set github
```

Approve the kit's credential request when prompted if you want it to use that
credential. For proxy-managed credentials, the real value stays on your host;
the sandbox proxy authenticates requests to the service on the kit's behalf.
See [Credential configuration](/manuals/ai/sandboxes/configuration/credentials.md)
for storing credentials, approving their use, and preparing unattended runs.

A kit's network requests still have to meet your sandbox's network policy.
For example, a GitHub kit can't grant access to GitHub if your organization's
policy doesn't allow it. If a connection fails, check the
[policy log](#debug-kits) to see which rule blocked it.

Agent instructions take effect without you copying files into your project.
For example, a Ruff kit can tell the agent to run the linter after changing
Python files. The agent receives that guidance alongside your project's own
instructions; adding the kit doesn't overwrite them.

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
