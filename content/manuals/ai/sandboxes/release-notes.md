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

## 0.31.1

{{< release-date date="2026-05-29" >}}

[GitHub release](https://github.com/docker/sbx-releases/releases/tag/v0.31.1)

### Bug fixes

- Fixes a bug introduced in v0.31.0 where sandboxes from earlier versions were not listed by sbx ls and could fail to run. Upgrading to v0.31.1 restores them.

## 0.31.0

{{< release-date date="2026-05-28" >}}

[GitHub release](https://github.com/docker/sbx-releases/releases/tag/v0.31.0)

### Highlights

> [!IMPORTANT]
> This release has a known issue where sandboxes from earlier versions may not be listed by sbx ls and can fail to run. This is fixed in [v0.31.1](https://github.com/docker/sbx-releases/releases/tag/v0.31.1) — please upgrade.

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

## 0.29.0

{{< release-date date="2026-05-13" >}}

[GitHub release](https://github.com/docker/sbx-releases/releases/tag/v0.29.0)

### Highlights

This release brings **per-sandbox network policies**, giving callers fine-grained control over which domains each sandbox can reach, including an explicit `deniedDomains` list and allowance for binary TCP protocols like SSH. Sandboxes now carry **daemon-assigned UUIDs**, enabling reliable identification across restarts and telemetry. Several **agent improvements** land in this release: Gemini gets SSO browser relay, Codex auth is more robust, and the OpenAI OAuth flow now auto-opens the browser. A round of **bug fixes** improves daemon robustness on macOS (long-username `sun_path` overflow), gVisor isolation under `--app-name`, and database-version handling.

### What's New

#### Networking & Policy

- Support per-sandbox scoped network policies
- Add `deniedDomains` to network kit policy
- Allow binary TCP protocols (e.g. SSH) through domain allow rules
- Pipe in policykit error handler for better diagnostics

#### Sandboxes

- Add daemon-assigned UUID to sandbox runtimes

#### Agents

- Enable SSO browser relay for Gemini
- Auto-open browser during OpenAI OAuth flow
- Skip auth.json placeholder for Codex when no host credentials
- Expose Claude guidance to Codex sandboxes

#### CLI

- Require confirmation for `sbx rm <name>` to prevent accidental deletion
- Unhide `kit` command in help output

#### Bug Fixes

- Namespace gVisor socket dir by `--app-name` so concurrent daemons don't share state
- Probe canonical socket path for `sun_path` budget — fixes `krun_start_enter failed` for macOS users with long usernames
- Check database version before starting the daemon and surface an instructive error instead of crashing
- Route gVisor sockets to a persistent, sandboxd-owned location
- Delete stranded tracker after failed auto-stop with no active sessions
- Clean up DinD volume even when container inspect fails
- Apply `SANDBOXES_STORAGE_ROOT` override to storage config
- Report running binary (not first `sbx` on PATH) in `diagnose`
- Explain how to configure OpenAI credentials in no-creds warning
- Allow MCR layer-blob CDN in default-code-and-containers policy
- Improve empty state of `sbx ls` with actionable guidance

<!-- END GENERATED RELEASES -->

## Earlier releases

For older versions, see the
[Docker Sandboxes releases on GitHub](https://github.com/docker/sbx-releases/releases).
