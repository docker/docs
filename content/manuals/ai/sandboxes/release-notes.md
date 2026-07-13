---
title: Docker Sandboxes release notes
linkTitle: Release notes
description: New features, bug fixes, and changes in Docker Sandboxes
keywords: docker sandboxes, sbx, release notes, changelog
toc_min: 1
toc_max: 2
tags:
  - Release notes
---

This page lists changes in recent stable releases of Docker Sandboxes. For
the full release history, including pre-releases and downloads, see the
[Docker Sandboxes releases on GitHub](https://github.com/docker/sbx-releases/releases).

<!-- BEGIN GENERATED RELEASES -->

## 0.35.0

{{< release-date date="2026-07-10" >}}

[GitHub release](https://github.com/docker/sbx-releases/releases/tag/v0.35.0)

### Highlights

This release revamps **policy tooling** with a concise `sbx policy ls`, a new `sbx policy inspect`, and a `sbx policy check network` command for testing whether the current policy would allow an access request before you run. 
Networking gains a **SOCKS5 upstream-proxy transport**.  
Secrets get a new **`sbx secret import`** with clearer env-source visibility.

### What's New

#### Networking & Proxy

- The sandbox proxy can chain upstream egress through a SOCKS5 proxy (`socks5://` / `socks5h://`, with optional auth) via `DOCKER_SANDBOXES_PROXY`, `HTTP_PROXY`, or `HTTPS_PROXY`.
- Add `DOCKER_SANDBOXES_NO_PROXY` to exclude destinations from `DOCKER_SANDBOXES_PROXY`, using standard `NO_PROXY` matching semantics.
- Droid OAuth credentials are now proxy-managed: real tokens stay on the host and never land in the sandbox.
- Faster sandbox startup: the TLS-proxy CA is installed by merging into the trust bundle instead of running `update-ca-certificates`, saving several hundred milliseconds.

#### Policy

- Simplify `sbx policy ls` and add `--wide`, `--source`, and `--decision` filters
- Add `sbx policy check` to test whether the current policy would allow an access request
- Balanced network preset now allows VS Code domains, Azure Blob Storage (`*.blob.core.windows.net`), and `dhi.io` over HTTP.

#### Kits

- `sbx kit add` now recreates the sandbox container with the augmented kit set instead of injecting at runtime. State is preserved with the re-creation.
- `sbx kit add` applies the added kit's network allow/deny rules and composed policy on the running sandbox.
- Re-attaching to a sandbox created from a custom `--kit` agent now works with `sbx run --name <name>` without re-passing `--kit`.
- Kits can inject the user's Docker login token into requests to docker.com hosts via a credential with service `sbx-login`.

#### CLI

- `sbx rm` now won't delete an active session unless `--force` is passed.
- `sbx inspect` now lists the sandbox's kits, injected secrets, and sandbox information.
- Added `sbx daemon` command (`start`, `stop`, `status`, `log-level`)

#### Secrets

- `sbx secret import` imports credential env vars into the keychain; `sbx secret ls` flags env-only and OAuth-shadowed entries. Host env vars no longer auto-inject at runtime — use `sbx secret import` to migrate.

#### Runtime & images

- Enable virtiofs caching by default on all operating systems by default for faster filesystem performance (`DOCKER_SANDBOXES_ENABLE_VIRTIOFS_CACHE=0` to opt out).

### Bug Fixes

- Fix "container not found" errors when copying files with `sbx cp` on a sandbox that has had a kit added.
- Enforce the one-credential-per-service rule on credential capture paths so a stale API key no longer shadows a newly captured credential.
- Fix `sbx login` failing with "The specified item already exists in the keychain" when signing back into a previously used account; logout now clears all stored Docker credentials.
- Restarted sandboxes keep GitHub access by rehydrating the stored `github` credential on daemon restart.
- Fix a custom kit clearing the proxy's built-in GitHub auth header mapping for the whole daemon until a restart.
- Tunnel plain-HTTP forward traffic (e.g. `apt`, port 80) via CONNECT when the upstream proxy only supports CONNECT.
- Sandbox egress through an upstream proxy identifies as `sbx-proxy` on the CONNECT handshake.
- Fix IPv6 policy allow rules using bracket notation (e.g. `[fdcb::1]:22`) not matching.
- Fix `sbx` connecting to the wrong Docker daemon when `DOCKER_HOST` is set in the environment.
- Serialize Docker Hub token refresh across the CLI and daemon so sign-in sessions aren't unexpectedly lost.

### Platform support

- Block installation on Windows versions older than Windows 11 (the only currently supported version).

## 0.34.0

{{< release-date date="2026-06-26" >}}

[GitHub release](https://github.com/docker/sbx-releases/releases/tag/v0.34.0)

### Highlights

Kit installs are now restricted to an allowlist of sources, defaulting to Docker Hub only — a **breaking change** if you install kits from a Git URL or another registry.

This release also renames `sbx policy set-default` to `sbx policy init`, restores published ports when a sandbox restarts, fixes a number of bugs, and adds two experimental previews: a native SSH endpoint and an `sbx setup` command for smoother first-time onboarding.

### What's New

#### SSH

- Add an experimental native SSH endpoint in sandboxd: connect with `ssh <sandbox-name>@127.0.0.1 -p 2222` (publickey auth, connect-to-create, interactive shell and exec; no SFTP yet). Enable with `sbx settings set feature.ssh true`.

#### Setup & Onboarding

- Add an experimental `sbx setup` command that imports agent credentials from environment variables.

#### Agents

- Cursor sandboxes no longer show the workspace trust prompt on launch.

#### Kits

- Add OCI v2 kit artifact streaming that decompresses the layer once to a cache directory and uses seek-based random access, so file content is not held in memory between reads.
- Restrict kit installs to an allowlist of sources, defaulting to Docker Hub (`docker.io/`) only.

  **Breaking:** installing a kit from another registry or a Git URL fails until you add its prefix with `sbx settings set kit.allowedSources`. See [Docs: Restrict kit sources](https://docs.docker.com/ai/sandboxes/customize/kits#restrict-kit-sources) for details.

#### CLI & Behavior Changes

- Rename `sbx policy set-default` to `sbx policy init`; the old name keeps working as a hidden, deprecated alias.
- Published sandbox ports are restored on restart, and the CLI/TUI can recover explicit host-port conflicts by choosing a new host port.

#### Bug Fixes

- Fix a daemon hang where a slow or stuck sandbox creation/deletion blocked `sbx ls`, the TUI, and new sessions until the daemon was restarted.
- Fix a kit mixin regression where adding `network.serviceDomains` for a service already provided by the base agent failed with a "credential … defined in both" error.
- Reject `+` in sandbox names with a clear validation error instead of panicking.
- Fix the interactive host-port conflict recovery prompt not appearing on Windows when restarting a sandbox whose published port is already in use.

## 0.33.0

{{< release-date date="2026-06-17" >}}

[GitHub release](https://github.com/docker/sbx-releases/releases/tag/v0.33.0)

### Highlights

`sbx run --name <sandbox>` now re-attaches to an existing sandbox by name. You can now also create multiple sandboxes of the same agent type and workspace by specifying unique sandbox names with `--name`.

Consequently, re-attaching to existing sandboxes with `sbx run <name>` is deprecated; the preferred form is `sbx run --name <name>`. The positional argument for `sbx run` should be an agent (e.g. `claude`, or `codex`). Sandbox name as the positional argument for run is still supported but will be removed in a future release.

This release also improves network isolation and policy enforcement. Sandbox DNS is now gated on network policy (closing a DNS-based exfiltration channel), ICMP egress is blocked across daemon restarts, and the MITM proxy publishes a CRL so revocation-strict clients keep working.

### What's New

#### Sandbox Identity & CLI

- `sbx run --name` now identifies a sandbox independent of the working directory: run multiple independently-named sandboxes in the same workspace, re-attach from any directory (agent may be omitted), and re-run a create command to re-enter. It no longer auto-creates numbered sibling sandboxes, prompts before entering a same-named sandbox from a different workspace, and errors when the requested agent doesn't match the named sandbox. The TUI follows the same rules.
- `sbx run <sandbox>` now prints a deprecation warning when re-attaching to an existing sandbox; use `sbx run --name <sandbox>` instead.
- `sbx ls --json` now reports a stable per-sandbox `id`.
- `sbx create` now fails with a clear missing-agent error when run without arguments.
- `sbx exec` now uses the same working directory as `sbx run`.
- `sbx cp -L` now follows symlinks in the source path for sandbox-to-host copies.
- Daemon inspect output is now included in the diagnostics bundle.

#### Networking & Proxy

- Sandbox DNS lookups are now gated on the network policy: a sandboxed process can no longer resolve domains that policy denies, closing a DNS-based data-exfiltration channel. Loopback names (e.g. `localhost`) are exempt to avoid breaking local OAuth callback flows. [CVE-2026-12039](https://www.cve.org/CVERecord?id=CVE-2026-12039)
- Outgoing ICMP from sandboxes is now blocked across daemon restarts. [CVE-2026-12539](https://www.cve.org/CVERecord?id=CVE-2026-12539)
- CIDR subnet allow rules (e.g. `sbx policy allow network 10.10.14.0/24`) now correctly permit traffic to IP addresses within the subnet.
- The MITM proxy now publishes a CRL and embeds a CRL distribution point in generated certificates, fixing clients that require certificate revocation checking (e.g. .NET `CheckCertificateRevocationList=true`).
- Removed the bracketed `[::1]` entry from the sandbox `NO_PROXY` default, fixing credential injection for HTTP clients that mis-parsed it.
- Claude connectors (Slack, Gmail, Notion, Atlassian, etc.) now work inside sandboxed Claude Code without manual policy overrides.

#### Secrets & Credentials

- `sbx secret set-custom --host`, and `serviceDomains` in kits, now accept wildcard host patterns (`*` matches one label, `**` matches any number) and is repeatable, so one custom secret can cover multiple subdomains/domains.

#### Agents

- Cursor OAuth is now supported

#### Platform & Performance

- The virtiofs cache is now enabled by default on macOS and Linux.
- Build packages for `linux/arm64` are now produced.
- On Linux, the keychain backend now falls back to the encrypted on-disk store when `dbus-launch` is unavailable, fixing headless/server hosts.

#### Bug Fixes

- Suppress a misleading warning when saving OAuth credentials while the daemon is not running.
- Fixed a TTY sizing issue on Windows.
- Keep agent entrypoint flags when arguments after `--` are themselves flags.
- Inject git identity from subdirectories and `[include]`d Git config when cloning.
- Proxy service detection now supports middle-position wildcards.
- Sandboxes blocked by mount policies are no longer filtered out on daemon startup.

<!-- END GENERATED RELEASES -->

## Earlier releases

For older versions, see the
[Docker Sandboxes releases on GitHub](https://github.com/docker/sbx-releases/releases).
