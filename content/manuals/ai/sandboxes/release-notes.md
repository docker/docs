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

## 0.28.2

{{< release-date date="2026-04-29" >}}

[GitHub release](https://github.com/docker/sbx-releases/releases/tag/v0.28.2)

### What's New

#### CLI

- Auto-open browser during login flow

#### Templates

- Install `ssh-add` and SSH client tools in the `main` template

#### Bug Fixes

- Prefer Codex OAuth over discovered API-key credentials
- Propagate host TTY size when running `sbx exec -it`
- Reveal trailing characters in masked secrets

## 0.28.1

{{< release-date date="2026-04-28" >}}

[GitHub release](https://github.com/docker/sbx-releases/releases/tag/v0.28.1)

### Highlights

A small release that wires **custom agent kits** through the CLI — discoverable in `--help` and invocable via `--kit` — and brings
**in-process sandbox run/exec** with launch-mode and settings dialogs to the TUI. Two bug fixes round it out: private Docker Hub image pulls work again via `--template`, and the secrets-masking path is tightened.

### What's New

#### CLI

- Make custom agent kits invocable and surface `--kit` in help
- TUI: in-process sandbox run/exec with launch mode dialog, settings dialog + misc fixes

#### Bug Fixes

- Enable private Docker Hub image pulls via `--template`
- Tighten secrets masking and emphasize `set-custom` warning

<!-- END GENERATED RELEASES -->

## Earlier releases

For older versions, see the
[Docker Sandboxes releases on GitHub](https://github.com/docker/sbx-releases/releases).
