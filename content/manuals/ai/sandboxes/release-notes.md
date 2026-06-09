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

## 0.32.0

{{< release-date date="2026-06-09" >}}

[GitHub release](https://github.com/docker/sbx-releases/releases/tag/v0.32.0)

### Highlights

**[Audit logging](https://docs.docker.com/ai/sandboxes/governance/audit/)**: Sandboxes now emit structured JSONL audit records for policy decisions. Records are written to a per-OS log directory and can be forwarded to any SIEM platform for enterprise compliance workflows. Requires a Docker AI Governance subscription.

**[Sign-in enforcement](https://docs.docker.com/ai/sandboxes/governance/sign-in-enforcement/)**: Administrators can now require Docker organization membership verification. Enforcement is deployed via standard endpoint management tooling: configuration profiles on macOS, the registry on Windows, and a JSON policy file on Linux. This closes the gap for organizations that need to ensure only authenticated, authorized users run AI coding agents.

### What's New

#### CLI

- Offer an interactive "Sign in with ChatGPT" OAuth flow on the first `sbx create`/`sbx run codex` when no Codex credentials are configured.
- Pre-select `balanced` as the highlighted default in the first-run network policy prompt, so pressing Enter accepts the recommended policy.
- Make global the default scope for `policy network allow|deny` and `policy rm`; add `--sandbox` to target a specific sandbox and drop the `-g/--global` flag.
- Simplify `sbx version` to a single line by default; gate detailed information behind `-D/--debug`.
- Unhide `sbx secret set-custom`, a command for [setting custom secrets](https://docs.docker.com/ai/sandboxes/security/credentials/#custom-secrets), and mark it as experimental.

#### Secrets

- Add OpenRouter as a built-in service provider, so `sbx secret set <sandbox> openrouter` works without `set-custom` and the proxy injects `Authorization: Bearer <token>` automatically.
- Fall back to an encrypted on-disk secrets store on Linux/WSL hosts where no working keychain is available, with a one-time warning on secret-writing paths including `sbx login`.
- Substitute custom-secret sentinels inside HTTP Basic auth payloads, so credentials referenced in `Basic` Authorization headers are resolved like other sentinel shapes.

#### Networking

- Hide inactive governed policy rules by default in `sbx policy ls` and the TUI Network Rules view, with governance/sync status, hidden-rule indicators, and an `--include-inactive` flag (TUI `i` toggle) to reveal them.
- Route OAuth/browser-open requests to the caller's graphical session, fixing `/login` opening on the host's display instead of the SSH terminal that invoked it.

#### Kits

- Support the v2 OCI kit artifact format end-to-end, so kits are standard OCI images that registries and OCI tooling (Hub, `oras`, `crane`, `skopeo`) can introspect without kit-specific knowledge.
- Write `files/workspace/<path>` kit entries correctly when `sbx run --clone` is used; previously the file hook fired before the in-container clone populated the workspace and failed the sandbox start.

#### Performance

- Keep virtiofs caching enabled for sandboxes using `--clone`, avoiding a FUSE round-trip on every `stat()` and speeding up `git status`, `grep -r`, and tree walks inside the sandbox.

#### Packaging

- Require the system keyring dependency in Linux packages so credential storage works out of the box.

#### Documentation

- Replace stale `--branch`/worktree guidance in generated agent guidance (CLAUDE.md/AGENTS.md) with `--clone`, including how to sync host commits via `/run/sandbox/source`.

#### Bug Fixes

- Fix an issue with `sbx secret set <sandbox> <service>` silently dropping credentials while reporting success.
- Migrate stale runtime `SocketPath` references on daemon restart, so sandboxes upgraded from v0.31.0 stay visible to `sbx ls` after `/tmp` is cleaned.
- Keep non-interactive `sbx exec` output intact by not tearing down the attach-exec bridge on stdin EOF (no more spurious empty output with exit code 0).
- Clear stale pending status in the TUI when a network deny rule is deleted, so a host no longer shows as Blocked after its rule is removed.
- Bind MCP gateway state to the daemon-assigned runtime instance so a same-name sandbox recreate cannot leave Claude pointed at a stale gateway port.
- Set the default network policy before launching the TUI to avoid spurious 412 errors from policy-rule requests.
- Stop counting expected `rm`/`stop`/list-ports "not found" 404s as analytics failures, so routine existence checks no longer inflate error dashboards.
- Require a daemon restart (instead of failing with `405 Method Not Allowed`) when downgrading the CLI below a newer running daemon.

## 0.31.3

{{< release-date date="2026-06-03" >}}

[GitHub release](https://github.com/docker/sbx-releases/releases/tag/v0.31.3)

### Bug Fixes

- Fix a failure to start sandboxes that were created with older versions of the CLI.
- Fix a file descriptor leak on Linux. Each credential lookup left a session
  D-Bus socket open, so long-running processes (such as the daemon) could
  gradually accumulate open file descriptors and eventually hit the session
  bus's connection limit, failing with "The maximum number of active
  connections has been reached." Connections are now closed after each
  operation. macOS and Windows were not affected.

## 0.31.2

{{< release-date date="2026-06-01" >}}

[GitHub release](https://github.com/docker/sbx-releases/releases/tag/v0.31.2)

### Highlights

This patch release resolves two reliability issues. It **fixes a Windows issue** where odd default sandbox memory values could lead to startup timeouts. It also includes a **daemon-compatibility fix** that prevents a silent failure (`405 Method Not Allowed`) when the `sbx` CLI is downgraded while a newer `sandboxd` daemon is still running — the CLI now requires a daemon restart instead.

### What's New

#### Bug Fixes

- Fix a Windows issue where odd default sandbox memory values could lead to startup timeouts.
- Require a daemon restart when downgrading the CLI below a running daemon, instead of silently proceeding into a `405 Method Not Allowed` error.

## 0.31.1

{{< release-date date="2026-05-29" >}}

[GitHub release](https://github.com/docker/sbx-releases/releases/tag/v0.31.1)

### Bug fixes

- Fixes a bug introduced in v0.31.0 where sandboxes from earlier versions were not listed by sbx ls and could fail to run. Upgrading to v0.31.1 restores them.

## 0.31.0

{{< release-date date="2026-05-28" >}}

[GitHub release](https://github.com/docker/sbx-releases/releases/tag/v0.31.0)

### Highlights

#### Clone mode: `--clone`

The `--branch` flag has been removed in favor of `--clone` (clone mode). Using `--branch` now fails with:

```console
$ sbx run claude --branch foo
ERROR: --branch is no longer supported; use --clone instead
```

Clone mode does not create a branch or worktree on your behalf — instead of a host-side worktree, the sandbox now runs against an in-container read-only clone.

- Your source repository is mounted into the sandbox read-only, and the shallow clone sets that mount as a Git remote. The agent only ever writes to the in-container clone, never to your working tree or .git/
- The clone lives on the sandbox's filesystem and is exposed back to the host as a `sandbox-<name>` Git remote served by `git-daemon` (no more `.sbx/<name>-worktrees/...` on the host).
- Forge remotes (`origin`, `upstream`, etc.) on the host are propagated into the in-container clone, so the agent can `git push origin` directly, the same way you would. Local-path remotes are skipped.
- Fetched sandbox refs are mirrored into `refs/sandboxes/<name>/*` on the host and persist after the sandbox is removed. Restore a branch from a removed sandbox with `git branch <local-name> refs/sandboxes/<name>/<branch>`. Commits that were never fetched, or uncommitted changes, are still lost on `sbx rm`.
- The `sandbox-<name>` remote is added to your host on `sbx create --clone` / `sbx run --clone` and removed on `sbx rm`, including across stop and restart.

### What's New

#### CLI

- `sbx create` auto-starts the daemon when it isn't already running.
- `sbx logout` now stops the daemon and running sandboxes.
- Unify terminal environment variables across `sbx run` and `sbx exec`.

#### Policies

- Show policy and rule names in CLI list output and TUI details.
- Add filters to the policies listing.

#### Kits

- Mark kits as experimental.
- Verbose error reporting for kit apply failures.

#### Sandboxes

- Opt a sandbox into virtiofs caching at create time via `DOCKER_SANDBOXES_ENABLE_VIRTIOFS_CACHE=1` (off by default; the choice is persisted in the spec and survives daemon restarts).

#### Networking

- Allow public-CA CRL/OCSP/AIA endpoints in the balanced proxy preset. Applies to new installations or after `sbx policy reset` (which removes any user-added rules).

#### Telemetry

- Surface `port_publish_failed` inner error detail.

#### Secrets

- Store container-registry pull credentials with `sbx secret set --registry`, so `sbx run --template` and `sbx run --kit` can pull from private registries (GHCR, ACR, ECR, Quay, …) without a `docker login`. Manage entries with `sbx secret ls` and remove them with `sbx secret rm --registry <host>`.

> [!WARNING]
> By default the credential is stored **host-side only** and is used just for pulling templates/kits. It is never placed inside a sandbox. If you pass `-g` (or scope it to a sandbox name), the credential is **injected into the sandbox in plaintext**, where the agent and any code running there can read it. Only use `-g`/sandbox scope when the sandbox itself needs to pull from the registry; otherwise omit `-g` to keep it host-only.

#### Bug Fixes

- Sort `template ls` output by repository, then tag.
- Retry `ExecResize` to keep the agent TUI in sync.
- Set `TERM=xterm-256color` when exec'ing with `-t`.
- Move the state directory symlink from `/tmp` to `~/.sbx/run/`.
- Stop `storageRootsGone` from locking the storagekit singleton.
- Use `engineError` and add retry debug logging in sandboxd.
- Retry transient shim start closures.
- Make Cursor session bootstrap proxy-local.
- Add bracketed `[::1]` to `NO_PROXY` for IPv6 loopback.
- Backdate proxy CA `NotBefore` to match the goproxy leaf cert window.

## 0.30.0

{{< release-date date="2026-05-19" >}}

[GitHub release](https://github.com/docker/sbx-releases/releases/tag/v0.30.0)

### Highlights

The CLI gets **non-interactive Docker Hub login** for scripted workflows, and sandboxes now have **a configurable grace period before auto-stopping** when the last session exits. Plus a wave of fixes covering Linux packaging, macOS worktree compatibility, Windows installer paths, network isolation, and recoverable sandbox state when host directories vanish.

### What's New

#### Governance & Policy

- Allow `sbx policy` setup before login

#### Kits & Agents

- Re-run `commands.startup` on every container start so init hooks are idempotent across restarts
- Per-kit memory files for progressive disclosure
- Enumerate installed kits in the AI memory file's Kits section

#### CLI & Auth

- Add non-interactive Docker Hub login for scripted workflows
- Migrate `/reset` to `/daemon/reset`; state-dir wipe is now daemon-side
- Print "Git repository detected" once when using `--branch`
- Skip implicit run options when the user provides explicit args

#### Networking & Sandboxd

- Bind both loopback stacks by default when publishing ports
- Allow raw TCP to `host.docker.internal` when localhost is allowed in policy
- Add grace period before auto-stopping a sandbox when the last session exits

#### Bug Fixes

- Build sailor's `ffi` crate instead of `ffi-krun` for packaged Linux release artifacts
- Keep sandboxes recoverable when workspace or worktree is deleted on the host
- Add macOS `/private` path compatibility for worktrees
- Probe canonical socket path for `sun_path` budget — fixes `krun_start_enter failed` on macOS with long usernames
- Namespace gVisor socket dir and auth/secret stores by `--app-name` so concurrent daemons don't collide
- Sanitize runtime ID when looking up gVisor network
- Check database version before starting the daemon; surface an instructive error instead of crashing
- Report Docker daemon startup time instead of the pre-start message in DinD
- Harden `BuildFileCredential` to check more than just file existence
- Open a sentinel connection in `cp` and `kit add` to prevent auto-stop race
- Remove redundant `ContainerKill` before `ContainerRemove` in sandboxlib
- Use a safe Windows `start` invocation for `OpenURL` in the TUI
- Rename WiX install directory id to `INSTALLFOLDER`

#### Documentation

- Warn agents about worktree path traps with `--branch`
- Improve consistency and wording in CLI help strings

<!-- END GENERATED RELEASES -->

## Earlier releases

For older versions, see the
[Docker Sandboxes releases on GitHub](https://github.com/docker/sbx-releases/releases).
