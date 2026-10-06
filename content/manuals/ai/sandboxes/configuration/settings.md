---
title: Docker Sandboxes settings
linkTitle: Settings
description: Configure Docker Sandboxes with sbx settings. View effective values, manage persistent overrides, and use environment variables for supported settings.
keywords: docker sandboxes, sbx, settings, configuration, environment variables, overrides
weight: 5
---

Use `sbx settings` to configure Docker Sandboxes on your host, including
clipboard access, kit sources, and defaults for sandbox creation. Settings
apply across your local sandboxes. For project-specific configuration, use
[environment files](environment-files.md).

The commands read and write settings through the local daemon, starting it if
necessary. Overrides persist across CLI invocations and daemon restarts.

## View settings

List settings with their effective value, type, source, and description:

```console
$ sbx settings list
```

Long values and descriptions are truncated in the table. Use
`sbx settings list --no-trunc` for complete text, or `sbx settings list --json`
for JSON output. The JSON records also include defaults and environment
variable names where available.

To inspect one setting:

```console
$ sbx settings get clipboard.imagePaste
false
$ sbx settings get clipboard.imagePaste --json
```

Without `--json`, `get` prints only the effective value. With `--json`, it
prints the complete setting record, including the source and default.

## Change a setting

Set an override by passing the setting key and a value of the required type.
For example, allow sandboxed agents to read images from your host clipboard:

```console
$ sbx settings set clipboard.imagePaste true
```

For JSON values, quote the argument so your shell passes it intact. This
example permits kits from Docker Hub and your organization's GitHub repositories:

```console
$ sbx settings set kit.allowedSources '["docker.io/","github.com/myorg/"]'
```

JSON arrays and objects replace the whole value. Include any entries you want
to keep.

Remove an override with `unset`:

```console
$ sbx settings unset clipboard.imagePaste
```

The setting falls back to its environment variable, if set, or its default.
Setting a value equal to its built-in default also removes the stored override.

### Value precedence

For each setting, the first available value wins:

1. The setting's environment variable, if it has one
2. A user override written with `sbx settings set`
3. The built-in default

The `SOURCE` column in `sbx settings list` identifies the selected source as
`envvar`, `override`, or `default`. If an environment variable takes precedence,
`set` still updates the stored override and reports why the effective value
hasn't changed. `unset` removes only the stored override, not the environment
variable.

### When changes take effect

Most changes take effect within about five seconds. Settings marked `yes` in
the `RESTART` column require a daemon restart for existing daemon-side consumers:

```console
$ sbx daemon restart
```

The `set` and `unset` commands print a restart reminder when needed. Some
settings apply only when creating a sandbox: changing a template default or disk
size doesn't update existing sandboxes. Proxy settings have separate timing for
CLI requests, daemon traffic, and sandbox traffic. See
[When proxy changes take effect](upstream-proxy.md#when-changes-take-effect).

## Environment variables

Some settings have an environment variable equivalent, listed in the reference
entries below. These variables configure Docker Sandboxes on the host. They
don't set environment variables inside a sandbox.

The daemon inherits environment variables when it starts. Export a variable
before the first `sbx` command, or restart an existing daemon from the shell
where you set it. For example:

```console
$ export DOCKER_SANDBOXES_CLIPBOARD_IMAGE_PASTE=true
$ sbx daemon restart
$ sbx settings get clipboard.imagePaste
true
```

Changing your shell environment doesn't change the environment of a running
daemon, even for settings that normally take effect without a restart. To
return to a stored override or default, remove the variable and restart:

```console
$ unset DOCKER_SANDBOXES_CLIPBOARD_IMAGE_PASTE
$ sbx daemon restart
```

CLI operations that read settings locally use their own environment on each
invocation. Keep the CLI and daemon environments consistent when using these
variables. A persistent override with `sbx settings set` avoids needing to
export a variable in each shell.

## Settings reference

Each entry lists its built-in default, before overrides. An environment
variable appears only when the setting has a direct equivalent. Use
`sbx settings list` to inspect the values supported by your installed version.

### Agents and host access

#### clipboard.imagePaste {.wrap-anywhere}

{{< setting-metadata type="boolean" default="false"
  env="DOCKER_SANDBOXES_CLIPBOARD_IMAGE_PASTE" >}}

Turn this on to paste screenshots and other host clipboard images into agents
such as Claude Code and Codex with `Ctrl+V`. Text paste doesn't need this setting.

```console
$ sbx settings set clipboard.imagePaste true
```

This grants sandboxed processes access to host clipboard images. The change
applies to running sandboxes without recreating them. Set it to `false` to
withdraw that access. See [image paste](../faq.md#can-i-paste-images-into-an-agent)
for supported behavior.

#### claude.remoteControl {.wrap-anywhere}

{{< setting-metadata type="boolean" default="false"
  env="DOCKER_SANDBOXES_CLAUDE_REMOTE_CONTROL" >}}

Turn this on before using Claude Code's `/remote-control` command inside a
sandbox:

```console
$ sbx settings set claude.remoteControl true
```

The remote-control connection must authenticate with its own session token.
This setting lets it do so instead of having the sandbox proxy replace that
token with the host credential. Set it to `false` to restore credential
replacement on that connection. Changes apply to requests from running
sandboxes. See [Claude Code remote control](../agents/claude-code.md#remote-control).

#### env.rememberHostCommands {.wrap-anywhere}

{{< setting-metadata type="boolean" default="false" >}}

Use this when you repeatedly run a trusted environment file and want to
approve its host commands once, until those commands change:

```console
$ sbx settings set env.rememberHostCommands true
```

The first approval is still required. Commands run on your host with your
permissions, outside the sandbox. Leave the setting at `false` to require
approval on every invocation, or use `--auto-approve` to approve only one
invocation. See [environment lifecycle commands](environment-files.md#lifecycle).

#### ssh.agentForwardingEnabled {.wrap-anywhere}

{{< setting-metadata type="boolean" default="true" >}}

Set this to `false` when sandboxes should not be able to request signatures
from your host SSH agent, including for Git authentication and commit signing:

```console
$ sbx settings set ssh.agentForwardingEnabled false
$ sbx daemon restart
```

Restarting the daemon applies the change to existing forwarders. When forwarding
is enabled, the private keys remain on the host, but sandboxed processes can
ask the agent to use them. See [SSH agent credentials](credentials.md#ssh-agent).

#### ssh.agentSocketPath {.wrap-anywhere}

{{< setting-metadata type="string" default="Empty string" >}}

Use a fixed path when every sandbox should use the same host SSH agent, such
as a password manager's agent, regardless of which shell starts the sandbox:

```console
$ sbx settings set ssh.agentSocketPath "$SSH_AUTH_SOCK"
$ sbx daemon restart
```

Run this from a shell whose `SSH_AUTH_SOCK` points to the intended agent. The
command stores that path, not a reference to the variable. With an empty value,
Docker Sandboxes uses the socket supplied by each client instead. Remove the
fixed path with `sbx settings unset ssh.agentSocketPath` and restart the daemon
to update existing forwarders.

### Images and storage

#### platform.images.registryMirror {.wrap-anywhere}

{{< setting-metadata type="string" default="Empty string" >}}

Use this when your organization routes Docker Hub pulls through a registry
mirror. Specify a host, optionally with a port and path prefix, without a URL
scheme:

```console
$ sbx settings set platform.images.registryMirror registry.example.com/docker-remote
```

The mirror applies to template and kit references that resolve to Docker Hub.
References to other registries stay unchanged. An empty value disables
mirroring.

A bare, non-loopback host can also configure Docker pulls inside sandboxes
created after the change. A mirror with a path prefix doesn't configure those
pulls. See [registry mirrors](registry-mirror.md) for authentication, certificate
requirements, and how to apply changes to existing sandboxes.

#### platform.images.useDHI {.wrap-anywhere}

{{< setting-metadata type="boolean" default="false"
  env="DOCKER_SANDBOXES_USE_DHI" >}}

Turn this on to use Docker Hardened Image variants when creating sandboxes
with default agent templates:

```console
$ sbx settings set platform.images.useDHI true
```

For example, the default template `docker/sandbox-templates:claude-code-docker`
becomes `dhi/sbx-templates:claude-code-docker`. The image tag stays the same.
An explicit `--template` or a custom kit image takes precedence.

Existing sandboxes keep their template. Create a sandbox after changing the
setting to use the selected image variant.

#### sandbox.disk.dockerVolume {.wrap-anywhere}

{{< setting-metadata type="string" default="10g" >}}

Increase this default if sandboxes need more room for Docker images,
containers, and volumes under `/var/lib/docker`. The minimum is 512 MiB:

```console
$ sbx settings set sandbox.disk.dockerVolume 20g
```

The size applies when creating a sandbox. It doesn't resize existing volumes
or increase the size of the sandbox workspace.

To choose a size for one sandbox without changing the default:

```console
$ DOCKER_SANDBOXES_DOCKER_SIZE=30g sbx create claude ~/my-project
```

This variable applies to that creation command. It doesn't change the value
reported by `sbx settings get sandbox.disk.dockerVolume`.

### Kits

#### kit.allowedSources {.wrap-anywhere}

{{< setting-metadata type="JSON" default="[\"docker.io/\"]"
  env="DOCKER_SANDBOXES_KIT_ALLOWED_SOURCES" >}}

Add a publisher here before installing its kits from a remote registry or Git
repository. The value replaces the entire allowlist, so retain any sources you
still need:

```console
$ sbx settings set kit.allowedSources '["docker.io/","github.com/myorg/"]'
```

Prefixes match on a path-segment boundary: `github.com/myorg/` permits that
organization's repositories, but not `github.com/myorg-other/`. The value
`["*"]` permits any remote source. Local directories and ZIP files are controlled
by [kit.allowLocalKits](#kitallowlocalkits), and pinned agent kits have the
[kit.allowExtractedAgents](#kitallowextractedagents) exception.

See [restrict kit sources](../customize/use-kits.md#restrict-kit-sources) for source
formats and examples.

#### kit.allowLocalKits {.wrap-anywhere}

{{< setting-metadata type="boolean" default="true"
  env="DOCKER_SANDBOXES_KIT_ALLOW_LOCAL" >}}

Set this to `false` to require kits to come from a remote source instead of a
local directory or ZIP file:

```console
$ sbx settings set kit.allowLocalKits false
```

Remote sources must still satisfy [kit.allowedSources](#kitallowedsources).
Keep this enabled while developing kits locally. Allowing local kits doesn't
exempt them from signature requirements when
[kit.requireSignature](#kitrequiresignature) is enabled.

#### kit.allowExtractedAgents {.wrap-anywhere}

{{< setting-metadata type="boolean" default="true"
  env="DOCKER_SANDBOXES_KIT_ALLOW_EXTRACTED_AGENTS" >}}

Some agents that previously shipped inside Docker Sandboxes are distributed
as kits. By default, Docker Sandboxes permits the exact pinned references it
uses for those agents even if their sources aren't in `kit.allowedSources`,
and exempts them from `kit.requireSignature`.

To apply your source and signature requirements to those kits too:

```console
$ sbx settings set kit.allowExtractedAgents false
```

After this change, launching one of those agents can fail unless you allow its
source and, when signatures are required, it has a signature from a trusted
signer. The default exception applies only to the pinned references, not to
other kits from the same publisher.

#### kit.requireSignature {.wrap-anywhere}

{{< setting-metadata type="boolean" default="false"
  env="DOCKER_SANDBOXES_KIT_REQUIRE_SIGNATURE" >}}

Turn this on to reject unsigned kits and kits whose signatures don't match
your trusted signers. Configure [kit.trustedSigners](#kittrustedsigners) first,
then require signatures:

```console
$ sbx settings set kit.requireSignature true
```

The check applies when installing kits from local directories, Git, or OCI
registries. ZIP kits can't carry verifiable signatures and are rejected.
Pinned agent kits remain exempt while
[kit.allowExtractedAgents](#kitallowextractedagents) is `true`.

A signature covers the kit's `spec.yaml` and `files/` content. It doesn't pin
image tags or verify downloads performed by the kit's commands. See
[sign and verify kits](../customize/kits-v2.md#sign-and-verify-kits).

#### kit.trustedSigners {.wrap-anywhere}

{{< setting-metadata type="JSON" default="Docker employee identities"
  env="DOCKER_SANDBOXES_KIT_TRUSTED_SIGNERS" >}}

Set the identities or public keys whose kit signatures you trust. The value is
a JSON array. Each entry describes either a keyless signer or a public key, and
a signature can match any entry.

For a keyless signer, specify both the identity and its OpenID Connect issuer:

```console
$ sbx settings set kit.trustedSigners \
    '[{"identity":"release-bot@example.com","issuer":"https://accounts.google.com"}]'
```

For a public key:

```console
$ sbx settings set kit.trustedSigners '[{"key":"/path/to/cosign.pub"}]'
```

Each command replaces the whole list. To trust both, include both objects in
one array. The default policy trusts Docker employee identities ending in
`@docker.com`, attested by Google's issuer. Setting trusted signers alone
doesn't require signatures: also enable
[kit.requireSignature](#kitrequiresignature).

#### kit.ignoreTransparencyLog {.wrap-anywhere}

{{< setting-metadata type="boolean" default="false"
  env="DOCKER_SANDBOXES_KIT_IGNORE_TLOG" >}}

Use this for private kits whose keyless signatures were created with
`--tlog-upload=false`, so verification doesn't require a public Rekor
transparency log entry:

```console
$ sbx settings set kit.ignoreTransparencyLog true
```

Signature and signer verification still apply. These signatures must provide a
timestamp from a timestamp authority instead of a transparency log timestamp.
Leave this at `false` when your signing workflow uses the public transparency
log. It has no effect on signatures verified with a public key.

### MCP gateway

#### mcp.forceLocalGateway {.wrap-anywhere}

{{< setting-metadata type="boolean" default="false" >}}

Set this to `true` to use the local MCP gateway when your account would
otherwise use the SaaS gateway:

```console
$ sbx settings set mcp.forceLocalGateway true
$ sbx daemon restart
```

The daemon caches its gateway selection, so a restart is required after changing
this setting. It doesn't override a gateway selected by organization governance.
Set it back to `false` and restart to return to automatic selection. See
[MCP gateway](../mcp-gateway.md) for server registration and agent setup.

### Diagnostics

#### diagnostics.autoUpload {.wrap-anywhere}

{{< setting-metadata type="string" default="Empty string" >}}

Choose whether Docker Sandboxes may automatically upload diagnostic bundles
after eligible daemon errors:

- `yes`: Consent to automatic uploads.
- `no`: Decline automatic uploads and suppress further consent prompts.
- Empty string: No recorded decision. Docker Sandboxes may prompt for consent,
  but doesn't automatically upload without it.

For example, decline automatic uploads:

```console
$ sbx settings set diagnostics.autoUpload no
```

Use `sbx settings unset diagnostics.autoUpload` to clear the decision. Read
[automatic diagnostics uploads](../troubleshooting.md#enable-automatic-diagnostics-uploads)
for what bundles contain before opting in. This controls automatic uploads,
not an explicit `sbx diagnose --upload` request.

#### diagnostics.autoUploadErrorCooldownInDays {.wrap-anywhere}

{{< setting-metadata type="integer" default="1" >}}

Increase this value to upload automatic diagnostic bundles less frequently.
For example, allow at most one automatic upload per seven days:

```console
$ sbx settings set diagnostics.autoUploadErrorCooldownInDays 7
```

The cooldown applies across all eligible errors, not separately to each error
type. Values below one are treated as one day. This setting doesn't grant
upload consent: [diagnostics.autoUpload](#diagnosticsautoupload) must be `yes`.

### Upstream proxies and TLS

Upstream proxy support is experimental. Use these settings to control how
outbound traffic reaches the network. See [upstream proxies](upstream-proxy.md)
for setup and [when changes take effect](upstream-proxy.md#when-changes-take-effect).

The standard `HTTP_PROXY`, `HTTPS_PROXY`, and `NO_PROXY` variables and their
lowercase forms are fallbacks for proxy selection, not direct overrides of the
`proxy` and `no_proxy` settings.

#### proxy {.wrap-anywhere}

{{< setting-metadata type="string" default="Empty string" >}}

Set a shared upstream proxy when both sandbox traffic and host-side Docker
Sandboxes traffic should pass through it:

```console
$ sbx settings set proxy http://proxy.corp:3128
$ sbx daemon restart
```

The value can be an HTTP, HTTPS, or SOCKS5 proxy URL, a PAC source, `system` to
use the operating system's proxy, or `direct` to bypass upstream proxies. An
empty value falls back to the standard proxy environment variables and then the
operating system's proxy settings.

The scope-specific settings below take precedence over this shared value.
See [proxy value formats](upstream-proxy.md#set-a-proxy-manually) and
[proxy precedence](upstream-proxy.md#precedence).

#### proxy.sandbox {.wrap-anywhere}

{{< setting-metadata type="string" default="Empty string"
  env="DOCKER_SANDBOXES_PROXY" >}}

Use this when traffic from inside sandboxes needs a different route from the
daemon and CLI. For example, send sandbox traffic through a SOCKS5 proxy:

```console
$ sbx settings set proxy.sandbox socks5h://proxy.corp:1080
$ sbx daemon restart
```

With `socks5h://`, the proxy resolves destination names. Use `direct` to bypass
the shared proxy for sandbox traffic, or unset this setting to inherit `proxy`.
It doesn't affect daemon traffic. Sandboxes created after the change use the
updated value; existing sandbox proxies require a daemon restart.

#### proxy.daemon {.wrap-anywhere}

{{< setting-metadata type="string" default="Empty string" >}}

Use this to route the daemon's requests, such as image pulls, separately from
sandbox traffic. Supported host CLI requests, including `sbx login` and
`sbx diagnose --upload`, also use this scope.

For example, keep a shared proxy for sandbox traffic while letting the daemon
connect directly:

```console
$ sbx settings set proxy.daemon direct
$ sbx daemon restart
```

Unset this setting to inherit `proxy`. Supported CLI requests use changes on
the next invocation; the daemon's own requests require a restart.

#### no_proxy {.wrap-anywhere}

{{< setting-metadata type="string" default="Empty string" >}}

List destinations that should bypass the selected upstream proxy. Use a
comma-separated string of hosts, domain suffixes, IP addresses, or CIDR ranges:

```console
$ sbx settings set no_proxy "registry.internal,10.0.0.0/8"
$ sbx daemon restart
```

This shared list applies to sandbox and daemon traffic unless a scope-specific
list replaces it. A value of `*` bypasses the upstream proxy for all destinations.
Proxy exclusions don't grant network access: sandbox requests must still pass
[network policy](../governance/access-controls/network.md).

#### no_proxy.sandbox {.wrap-anywhere}

{{< setting-metadata type="string" default="Empty string"
  env="DOCKER_SANDBOXES_NO_PROXY" >}}

Use this when only sandbox traffic should bypass the upstream proxy for a set
of destinations. For example, connect directly to cluster services:

```console
$ sbx settings set no_proxy.sandbox "*.svc.cluster.local"
$ sbx daemon restart
```

A non-empty value replaces the shared `no_proxy` list for sandbox traffic; it
doesn't append to it. Include any shared exclusions that sandboxes still need.
Unset it to inherit `no_proxy` again. Existing sandbox proxies require a daemon
restart.

#### no_proxy.daemon {.wrap-anywhere}

{{< setting-metadata type="string" default="Empty string" >}}

Use this when the daemon and supported CLI requests need different proxy
exclusions from sandbox traffic. For example, pull images directly from an
internal registry:

```console
$ sbx settings set no_proxy.daemon "registry.internal"
$ sbx daemon restart
```

A non-empty value replaces the shared `no_proxy` list for this scope. It doesn't
change sandbox exclusions. Unset it to inherit `no_proxy` again. Supported CLI
requests use changes on the next invocation; daemon requests require a restart.

#### proxy.integratedAuth {.wrap-anywhere}

{{< setting-metadata type="boolean" default="false" >}}

Turn this on when a corporate proxy requires NTLM or Kerberos authentication
using your Windows sign-in identity:

```console
$ sbx settings set proxy.integratedAuth true
$ sbx daemon restart
```

It applies to both sandbox and daemon proxy traffic. Authentication happens on
the host, so the Windows credentials don't enter the sandbox. It has no effect
on macOS or Linux. For cross-platform authentication with credentials in a
proxy URL, see [proxy authentication](upstream-proxy.md#authentication).

#### tls.allowNegativeSerial {.wrap-anywhere}

{{< setting-metadata type="boolean" default="false"
  env="DOCKER_SANDBOXES_TLS_ALLOW_NEGATIVE_SERIAL" >}}

Use this compatibility setting when HTTPS requests fail with
`x509: negative serial number` because a TLS-inspecting proxy issues
certificates with negative serial numbers:

```console
$ sbx settings set tls.allowNegativeSerial true
$ sbx daemon restart
```

It relaxes certificate validation to accept those serial numbers. It doesn't
make an untrusted certificate authority trusted or resolve other certificate
errors. For an untrusted internal CA, see
[certificate troubleshooting](../troubleshooting.md#api-calls-fail-with-a-certificate-error).

### Shared agent skills

Shared agent skills are experimental.

#### skills.defaultMode {.wrap-anywhere}

{{< setting-metadata type="string" default="readonly" >}}

Choose how future sandboxes use the shared agent skills store when you omit
`--skills`:

- `readonly`: Agents can read shared skills but can't change the store.
- `readwrite`: Agents can read and modify skills used by other sandboxes.
- `off`: Don't mount the shared store.

For example, omit the shared store by default:

```console
$ sbx settings set skills.defaultMode off
```

An explicit `--skills` value or environment-file `skills` value overrides this
default. Existing sandboxes retain their mounts: recreate them to change their
access mode. See [share agent skills](../workflows/agent-skills.md) for setup
and the consequences of sharing a writable store.

## Command reference

### sbx settings list

List settings. Alias: `sbx settings ls`.

Use `--json` for complete records as JSON, or `--no-trunc` for full text. These
options are mutually exclusive.

### sbx settings get \<KEY\>

Print one effective value.

Use `--json` for the complete setting record.

### sbx settings set \<KEY\> \<VALUE\>

Write a user override.

Values are parsed as `bool`, `int`, `float`, `string`, or `json`, according to
the setting's type.

### sbx settings unset \<KEY\>

Remove a user override.

Use `sbx settings <COMMAND> --help` for command help.
