---
title: Use kits
description: Run published v3 kit sets and workloads, add compatible mixins, and configure kit arguments and runtime access.
keywords: sandboxes, sbx, kits, v3, workloads, mixins
weight: 10
---

{{< summary-bar feature_name="Docker Sandboxes sbx" >}}

If you've run `sbx run claude` or `sbx run codex`, you've already used a kit.
The built-in agents are kits that package an environment, tools, and runtime
settings. You run kits you build yourself or get from another publisher in the
same way: give `sbx` a reference, and Docker Sandboxes prepares the environment
and applies the kit's settings.

This page shows how to run published kits, combine them with mixins, and
customize their settings.

## Run a kit

The built-in agent names are shortcuts for kit references. To run another kit,
replace the agent name with that kit's reference. For example, run Docker's
published v3 Codex workload:

```console
$ sbx run docker.io/docker/sbx-kit-codex:0.155.1 --name codex-v3-kit
```

This starts Codex using a v3 kit from Docker Hub. The built-in `codex`
shortcut uses a v2 kit to preserve compatibility with existing customizations.

### Choose a kit source

The previous example uses a published image from Docker Hub. You can also
run a kit from a local directory or a Git repository. The same source types
work for mixins added with `--kit`:

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
reuse cached results.

By default, remote kit sources are limited to Docker Hub. To use a Git
source or another registry, see [Restrict kit sources](#restrict-kit-sources).
For private images, see
[Registry credentials](/manuals/ai/sandboxes/configuration/credentials.md#registry-credentials).

## Reuse a sandbox

A sandbox keeps the kit configuration it was created with. To return to the
sandbox from the previous example, specify its name:

```console
$ sbx run --name codex-v3-kit
```

You don't need to specify the kit reference again. To try a different workload,
mixin combination, or argument value, create a sandbox with a different name
or recreate the existing one. To add a mixin to an existing sandbox, see
[Add mixins](#add-mixins).

## Runtime access and instructions

A kit can request access to the services it needs. For example, the Codex
workload needs to reach OpenAI and authenticate. If you use an API key,
store it on your host and approve the kit's request to use it. These are
separate steps: storing a secret doesn't give a third-party kit permission to
use it. See
[Credential bindings](/manuals/ai/sandboxes/configuration/credentials.md#credential-bindings)
for preparing unattended runs.

Check the kit's documentation for the services it contacts, the credentials
it needs, and the commands it runs at startup. For proxy-managed credentials,
the kit's credential binding tells the proxy which credential to add to
requests to the service. Once you approve the binding, the proxy adds the
credential to outgoing requests. The credential value stays on your host
and isn't exposed inside the sandbox.

Network requests must also meet your sandbox's network policy. Kit allow rules
can't grant access beyond your organization's policy. If a connection fails,
check the [policy log](#debug-kits) to see which rule blocked it.

Kits can also give the agent instructions for using their tools. The workload
chooses the instruction file, and mixins add their guidance to it. Docker
Sandboxes writes that file outside your workspace, leaving your project's
instructions intact.

## Add mixins

Mixins add tools and configuration to a workload. Use `--kit` to add a mixin
when creating a sandbox. For v3 kits, both the workload and the mixin must
use v3.

For example, an internal CLI mixin can package your company's executable
along with network rules and a credential request for its API:

```console
$ sbx run docker.io/docker/sbx-kit-codex:0.155.1 --name codex-tools \
    --kit docker.io/<NAMESPACE>/company-cli:1.0.0
```

Replace the mixin reference with one your organization has published.
To build your own, see
[Build a tool mixin](/manuals/ai/sandboxes/customize/author/tool-mixins.md).
Store the credential requested by the mixin on the host and approve access
when prompted.

Codex starts with `company-cli` available to use. Docker Sandboxes
applies the files and settings from both the workload and the mixin.

Built-in shortcuts such as `claude` and `codex` use v2 kits and require
[v2 mixins](/manuals/ai/sandboxes/customize/kits-v2.md).
See [Version compatibility](/manuals/ai/sandboxes/customize/_index.md#version-compatibility)
for details.

### Combine multiple mixins

Repeat `--kit` to add more mixins. You can also add compatible mixins to a
published workload set, or pass a set containing only mixins with `--kit`.

Avoid adding a kit that's already in the set. If two kits provide the same
feature, Docker Sandboxes can reject the combination.

To publish your combination as one reference, see
[Compose a kit set](/manuals/ai/sandboxes/customize/author/kit-sets.md).

### Change a sandbox's mixins

For v3 kits, you choose the mixins when you create the sandbox. To use a
different combination, create another sandbox with a different name.
Specify the workload and all the mixins you want to include.

The `sbx kit add` command can't add mixins to an existing v3 sandbox.

The new sandbox doesn't inherit changes or kit volume data from the previous sandbox.

## Compose kits

Some mixins need a tool another kit supplies. For example, a mixin whose
scripts call `gojq` can declare that dependency:

```yaml
requires: [gojq]
```

The [gojq mixin example](/manuals/ai/sandboxes/customize/author/patterns.md#leave-build-tools-out-of-the-mixin)
declares the feature it provides:

```yaml
provides: ["gojq@0.12.17"]
```

Include both mixins when you create the sandbox. Docker Sandboxes applies
the gojq mixin before the mixin whose scripts need it, regardless of the order
of your `--kit` flags. If you leave out gojq, sandbox creation fails.

A `requires` entry tells Docker Sandboxes what a kit needs. You still choose
and include the kit that provides it. Docker Sandboxes doesn't download a
kit for you based on this entry. A kit can also require a minimum feature
version or rule out incompatible kits.

With `--kit`, Docker Sandboxes checks and combines the kits when it creates
the sandbox. For a set, this happens when the publisher builds it. The
published kit includes a record of the exact components used, identified by
their image digests.

To declare these relationships in your own kits, see
[Composition fields](https://github.com/docker/sandbox-kit-spec/blob/main/docs/spec/SPEC-v3.md#5-provides-requires-integrates-conflicts).

## Pass arguments to kits

Kits can expose arguments for settings such as a linter's mode. Check the
kit's documentation for argument names, defaults, and accepted values.
Use `--kit-arg name=value` to set an argument:

```console
$ sbx run docker.io/<NAMESPACE>/codex-tools:1.0.0 \
    --name codex-tools-fix --kit-arg lint_mode=fix
```

This example uses a workload set that exposes a `lint_mode` argument.
Replace the image reference with your set's published reference. To expose
arguments in your own set, see
[Configure component arguments](/manuals/ai/sandboxes/customize/author/kit-sets.md#configure-component-arguments).

Argument values are plain text and can be recorded in shell history and
sandbox state. Use
[stored credentials](/manuals/ai/sandboxes/configuration/credentials.md)
for secrets.

### Target a specific kit

An argument without a kit prefix applies to every selected kit that declares
it. To target one kit, use `--kit-arg <HANDLE>.<ARGUMENT>=<VALUE>`:

```console
$ sbx run docker.io/<NAMESPACE>/codex-tools:1.0.0 \
    --name codex-tools-fix --kit-arg codex-tools.lint_mode=fix
```

The handle identifies the kit and comes from its reference:

| Kit source | Handle |
| --- | --- |
| Published image | Final part of the repository name, such as `codex-tools` |
| Local directory | Directory name |
| Git repository | Selected subdirectory name, or repository name if no subdirectory is selected |

A value targeted at one kit overrides a value supplied to all kits.

### Load arguments from a file

Use `--kit-args-file <FILE>` to load arguments from a reusable file.
Write one `name=value` entry per line. You can prefix names with a kit
handle, as with `--kit-arg`.

Values passed with `--kit-arg` override those in the file.

### Which settings can you change?

For a kit set, you can change only the arguments the author has exposed.
Other component arguments are fixed when the set is published.

Settings chosen during the image build, such as a tool's version, require
rebuilding the image.

## Restrict kit sources

[`kit.allowedSources`](../configuration/settings.md#kitallowedsources) controls
permitted remote kit sources. Its default permits
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

For defaults and environment variable equivalents, see the
[kit settings reference](../configuration/settings.md#kits).

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
