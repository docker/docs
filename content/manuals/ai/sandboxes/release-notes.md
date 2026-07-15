---
title: Docker Sandboxes release notes
linkTitle: Release notes
description: New features, bug fixes, and changes in Docker Sandboxes
keywords: docker sandboxes, sbx, release notes, changelog
weight: 120
toc_min: 1
toc_max: 2
tags:
  - Release notes
---

This page lists changes in recent stable releases of Docker Sandboxes. For
the full release history, including pre-releases and downloads, see the
[Docker Sandboxes releases on GitHub](https://github.com/docker/sbx-releases/releases).

<!-- BEGIN GENERATED RELEASES -->

## 0.37.1

{{< release-date date="2026-07-29" >}}

[GitHub release](https://github.com/docker/sbx-releases/releases/tag/v0.37.1)

### Highlights

This patch release stops SSH sessions from **forwarding credential environment variables into sandboxes by default**. Variables such as `ANTHROPIC_API_KEY`, `OPENAI_API_KEY`, and `GH_TOKEN` are no longer sent from the client to the sandbox unless explicitly allowed via the `ssh.acceptEnv` setting.

### What's New

#### Bug Fixes

- SSH sessions no longer forward credential environment variables (`ANTHROPIC_API_KEY`, `OPENAI_API_KEY`, `GH_TOKEN`, ...) from the client into the sandbox by default; use the `ssh.acceptEnv` setting to opt back in for specific variables

## 0.37.0

{{< release-date date="2026-07-24" >}}

[GitHub release](https://github.com/docker/sbx-releases/releases/tag/v0.37.0)

### Highlights

**SSH access to sandboxes (experimental).** Docker Sandboxes can now be used as SSH targets. After enabling SSH access, run `sbx setup ssh` once, then connect to an existing sandbox by name with `ssh my-sandbox.sbx`. Use the connection for interactive shells, one-shot commands, and SSH-based remote development.

**Shared agent skills.** Docker Sandboxes can now import skills from supported host agents into a persistent store shared across sandboxes. Run sbx skills import to import them. New sandboxes mount the store read-write by default; use --no-share-skills to opt out.

### What's New

#### SSH

- `sbx setup ssh` adds a managed `*.sbx` entry to your SSH config, making existing sandboxes available at `<name>.sbx`.
- SSH connections start the local Docker Sandboxes daemon and the target sandbox automatically when needed.
- Connect using OpenSSH-compatible clients and remote-development tools such as VS Code, Cursor, Claude Desktop, and ChatGPT.

#### Shared skills

- `sbx skills import` discovers and imports skills from the host and makes them available to sandboxed agents. Use `--dry-run` to preview imports and `--force` to replace existing skills.
- Imported skills persist after sandbox deletion and are mounted into new sandboxes for supported agents.
- Pass `--no-share-skills` to `sbx run` or `sbx create` when creating a sandbox to opt out.

#### CLI

- `sbx create` and `sbx run` accept `-p/--publish` to publish sandbox ports at creation time.

#### Networking & Policy

- `DOCKER_SANDBOXES_PROXY=system` routes sandbox egress through the host operating system's proxy configuration (macOS/Windows), including any PAC auto-config URL.
- Governance-policy denials can now display an organization-configured support message (for example, who to contact).

#### Security & Audit

- Audit now emits execution-outcome records for network egress (per allowed connection) and filesystem mounts (per allowed path) — success, latency, and error class — alongside the policy-attributed decision records.
- sandboxd excludes itself from Windows Error Reporting so daemon crash dumps cannot capture in-memory credentials.

#### Performance

- `sbx secret ls` and sandbox startup are faster on Linux hosts without an OS keychain — stored secrets are no longer all decrypted just to list or resolve credentials.

### Bug Fixes

- Fixed sandboxd failing to start on Linux hosts without an OS keychain, where the on-disk secret store's key derivation could peg a CPU during startup and the CLI would kill the still-starting daemon.
- Fixed an intermittent "failed to fully delete sandbox" error when removing a running sandbox, caused by a network-teardown race with the engine's endpoint cleanup.

## 0.35.0

{{< release-date date="2026-07-10" >}}

[GitHub release](https://github.com/docker/sbx-releases/releases/tag/v0.35.0)

### Notice

There are no Linux/ARM64 builds for v0.35.x due to stability issues that were encountered during this release. We plan to bring them back for the next release.

### Highlights

- **Host environment variables are no longer used for authentication.** Previous versions automatically detected API keys in predefined environment variables (such as `ANTHROPIC_API_KEY`) and injected them into model provider requests. Starting with this release, sandboxes only authenticate using credentials you've explicitly stored, or OAuth for agents that support it. If you relied on environment variables, run the new `sbx secret import` command once to move your keys into the keychain. See the [credentials documentation](https://docs.docker.com/ai/sandboxes/security/credentials/) for details.
- Policy commands are revamped with a more concise `sbx policy ls`, a new `sbx policy inspect`, and a `sbx policy check network` command for testing whether the current policy would allow an access request before you run. 
- Networking gains a **SOCKS5 upstream-proxy transport**.  

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

<!-- END GENERATED RELEASES -->

## Earlier releases

For older versions, see the
[Docker Sandboxes releases on GitHub](https://github.com/docker/sbx-releases/releases).
