---
title: Docker Sandboxes release notes
linkTitle: Release notes
description: New features, bug fixes, and changes in Docker Sandboxes
keywords: docker sandboxes, sbx, release notes, changelog
weight: 150
toc_min: 1
toc_max: 2
tags:
  - Release notes
---

This page lists changes in recent stable releases of Docker Sandboxes. For
the full release history, including pre-releases and downloads, see the
[Docker Sandboxes releases on GitHub](https://github.com/docker/sbx-releases/releases).

<!-- BEGIN GENERATED RELEASES -->

## 0.47.0

{{< release-date date="2026-10-05" >}}

[GitHub release](https://github.com/docker/sbx-releases/releases/tag/v0.47.0)

### Highlights

Docker Sandboxes v0.47.0 adds automatic cleanup after agent sessions with `sbx run --rm` and a kit capability for keeping local sandboxes running after sessions disconnect.

### What's new

#### Security

- Fixed OAuth token interception when a provider hostname uses different capitalization or a trailing dot, preventing real tokens from reaching the sandbox instead of the proxy's placeholders.
- The proxy rejects unrecognized OAuth token grants and prevents their responses from replacing host-managed credentials. Supported in-sandbox sign-in flows remain available.
- The proxy returns an error when it cannot safely mask a successful Anthropic API-key creation response, instead of forwarding the unmasked response to the sandbox.
- Included since v0.46.0: fixed raw TCP connections to denied hostnames being permitted by an allow rule for the hostname's resolved IP address. This fix applies to TCP; the related UDP case with multiple tracked hostnames remains outside its scope.
- Included since v0.45.0: CLI-created cloud hostname allowlists no longer gain implicit `0.0.0.0/0` and `::/0` rules. Existing stored policies are unchanged; remove those rules or recreate the policy to apply the restriction. Explicit IP and CIDR allowances remain supported.
- The proxy rejects TLS handshakes that fill its inspection buffer before a complete ClientHello can be checked on a non-MITM CONNECT tunnel or during the transparent proxy's late handshake check.
- The proxy closes incomplete TLS handshakes on those paths after a two-minute timeout instead of retaining the connections for the sandbox's lifetime.
- Sandboxes reject UDP to multicast, link-local, and unspecified destinations, limited broadcast, and broadcast addresses derived from the host's interface prefixes, regardless of network policy. UDP to `host.docker.internal` is unaffected. Broadcast addresses configured outside that derivation and networks reachable only through routes are not covered by this check.
- Kit pulls enforce limits on registry-declared blob sizes, decompressed content, and archive entry counts.

#### Sandbox lifecycle and workspaces

- `sbx run --rm` removes a local or cloud sandbox after its agent session ends. It cannot be combined with `--detached`. If a cloud session is interrupted, such as by a dropped connection, the sandbox is kept and the CLI prints the command to remove it.
- Running `sbx run -d` against an existing local sandbox keeps it running after sessions disconnect, until you stop or remove it.
- Dynamic mounts and permissions created through symlink paths can be removed without reappearing after a restart. Incompatible saved records include recovery guidance.

#### Kits

- Local kit inspection, validation, and pulling require the daemon to be running.
- Kits can declare `com.docker.sandbox/long-running@1` to keep local sandboxes running after all sessions disconnect. Cloud sandboxes and `sbx kit add` cannot provide this capability: required entries are rejected, and optional entries are skipped.
- `sbx kit sign` and `sbx kit push --sign` succeed on registries that refuse manifest deletion, including GitHub Container Registry and Docker Hub, instead of reporting failure after attaching the signature.
- Adding a kit to a sandbox whose guest has stopped responding fails without leaving the sandbox unusable until the daemon restarts.

#### Cloud sandboxes

- `sbx --cloud ttl` reports stopped sandboxes as stopped instead of expired. The time-to-live restarts when the sandbox resumes. JSON output includes `stopped` and `ttl_paused`; while the sandbox is resuming, only `ttl_paused` is true.
- Creating a cloud sandbox or moving a local sandbox to the cloud applies cloud policy and the kit's network rules without copying locally added network rules. Moving a sandbox with local HTTP method or path restrictions warns that those restrictions will not apply in the cloud.

#### Settings and proxies

- Upstream proxy recovery no longer blocks unrelated sandboxes or deletes isolated container data when credentials are unavailable. Local host services remain reachable.
- `sbx settings set` rejects invalid upstream proxy values before saving them.

#### MCP

- MCP authorization requests `offline_access` when the server advertises it and authorization uses saved defaults or resource-required scopes, so the server can issue refresh tokens. Explicit `--scope` values are unchanged. Run `sbx mcp auth` again to obtain a grant with a refresh token for existing credentials.
- Fixed MCP gateway availability in sandboxes created from the TUI and when connecting over SSH after a daemon restart or automatic sandbox creation.

#### CLI and updates

- `sbx env` reports the correct file, line, and column for unrecognized keys even when another entry in the file fails custom validation. Multiple validation errors appear on separate lines.
- Help remains available when the settings directory is unwritable or local settings cannot be opened, with a warning instead of a panic.
- Host-port collisions from `sbx ports --publish` identify the occupied binding and offer a retry with an automatically allocated port when safe.
- Windows update notices appear only after WinGet confirms that the version is available in its catalog.
<!-- release-notes:end -->

## 0.46.0

{{< release-date date="2026-09-28" >}}

[GitHub release](https://github.com/docker/sbx-releases/releases/tag/v0.46.0)

### What's new

#### Breaking changes

- Secret commands configured with `sbx secret set --command`, `sbx secret set-custom --command`, or `secrets.<name>.command` in an environment file execute from a fresh temporary directory on the host. Relative paths such as `./credential-helper` no longer resolve from the project directory or the directory where you ran `sbx`. Store helpers and their dependencies outside writable sandbox mounts. Run helpers by name from an absolute directory on the host's `PATH`, use absolute paths, or explicitly change to their private directory in the command. For existing environments that declare secret commands, the next `sbx env run` asks you to approve a one-time plan change for the working directory. The execution change takes effect after upgrading and restarting the daemon, even before you approve that plan.

#### Cloud sandboxes

- `sbx --cloud create --on-timeout restart` accepts `restart` as the timeout action. When the sandbox reaches its time limit, it stops and immediately restarts instead of remaining stopped.
- `sbx --cloud create` passes `--kit-arg` and `--kit-args-file` values to kits supplied with `--kit`.

#### Kits and skills

- Kits can install files in the agent's skills directory when shared skills are read-only. The shared skills store remains read-only, while kit installation and startup commands can write their own skills without a read-only filesystem error.
- Kits added through the runtime API retain their network rules and applicable agent instructions when another kit addition recreates the sandbox container.
- `sbx kit add` warns if it cannot save the updated sandbox record. The warning explains which kit settings could be lost and whether a daemon restart or another container replacement would cause the loss.

#### Agents and models

- The local model server starts and stops with the Docker Sandboxes daemon and downloads its llama.cpp runtime when the daemon starts. The macOS and Windows bundles include llmman v0.1.418, which manages the runtime download instead of relying on a separately bundled `llama-server`.
- Image paste in WSL2 falls back to the Windows clipboard when Linux clipboard tools return no image. Requires `clipboard.imagePaste` to be enabled.

#### Sandbox lifecycle and workspaces

- Fixed a shutdown bug affecting templates that use dash as `/bin/sh`, including the built-in Ubuntu-based templates. The shutdown handler forwards `SIGTERM` correctly, giving sandbox processes a chance to exit gracefully instead of waiting five seconds for a forced shutdown.
- On Linux arm64 hosts, the default CPU allocation is capped at 16 CPUs per sandbox. This fixes startup failures with `VM did not connect within 15s` when several sandboxes start together on hosts with many CPU cores. Use `--cpus` to request a larger allocation.
- `sbx umount` can remove a saved mount from a stopped sandbox using the original host path even after that directory has been deleted.
- On macOS, mounting, unmounting, and restoring saved mounts consistently recognize host paths whose capitalization differs.
- If the runtime fails to mount the workspace, sandbox startup reports the mount failure and points to the daemon log for the cause instead of reporting a generic container startup error.
- Starting a second daemon against a state directory already in use fails with an error instead of disrupting the running daemon.
- `sbx reset` stops background feature-flag updates and log writes before deleting local state, preventing leftover files and recreated directories. On Windows, it also closes daemon log files before deleting them to avoid cleanup retries caused by open file handles.

#### Authentication and credentials

- Removing secrets in bulk revokes credentials from running sandboxes. Failed revocations can be retried even after the stored secrets have been deleted.

#### Networking and policy

- Sandboxes created with the `balanced` network policy preset can download Playwright browser binaries from `cdn.playwright.dev` over HTTPS. Existing sandboxes keep their saved policies. To use the updated preset, create a new sandbox with the `balanced` network policy.

#### CLI and diagnostics

- `sbx env` reports unrecognized environment-file keys with the file, line, and column where they were declared, including when multiple files are merged.
- Canceling a batch `sbx rm` or `sbx stop` stops processing the remaining sandboxes instead of printing a cancellation error for each one.
- `sbx diagnose --upload` returns a non-zero exit status if the requested diagnostics upload fails, so scripts can detect the failure.
- `sbx diagnose` reports the socket-path length limit used by the container runtime as `runtime_socket_limit_bytes` and clarifies the meaning of the reported socket-path values.

#### Packaging and installation

- Windows MSI installations include `llmman` and the guest kernel, fixing local model serving with `sbx run --model` and sandbox launches that failed with `No kernel specified`.
- Uninstalling Docker Sandboxes through the Windows MSI stops the daemon.
- Fixed the Linux static tarball failing to start on distributions with older glibc versions. The tarball is built against glibc 2.34.

## 0.45.1

{{< release-date date="2026-09-22" >}}

[GitHub release](https://github.com/docker/sbx-releases/releases/tag/v0.45.1)

### Fixes and improvements

- Improved sandbox moves and support for private kit images in cloud sandboxes.

## 0.45.0

{{< release-date date="2026-09-21" >}}

[GitHub release](https://github.com/docker/sbx-releases/releases/tag/v0.45.0)

### Highlights

#### Compose reusable environments with v3 kits

Docker Sandboxes now supports v3 kits: OCI-based packages that combine an agent workload with reusable mixins for tools, configuration, credentials, network access, and agent instructions. Compose compatible kits directly when creating a sandbox, or publish the combination as a kit set that your team can run from a single reference.

V2 kits remain supported for built-in agents and existing customizations. V3 workloads and mixins must be used together; they can't be combined with v1 or v2 kits. [Learn more about kits](https://docs.docker.com/ai/sandboxes/customize/).

#### Run agents in cloud sandboxes

Run AI agents on Docker-managed cloud infrastructure with `sbx --cloud`. Cloud support is experimental and requires an active Docker Agentic Platform subscription. See [Get started with cloud sandboxes](https://docs.docker.com/ai/sandboxes/cloud/).

### What's new

#### Breaking changes

- `sbx mcp catalog` has been removed. To authorize a remote MCP server, register it with `sbx mcp add` before running `sbx mcp auth`.
- `sbx secret rm` now returns an error on stderr when the requested secret doesn't exist.
- `sbx mcp rm` now returns an error when the requested MCP server isn't registered.

#### Security

- Fixed an issue where revoking a sandbox's OAuth or API-key credential could leave its running proxy authorized until the sandbox was recreated.

#### Kits

- V3 kits introduce separate workload and mixin roles. A workload supplies the base environment and command to run; mixins add tools, configuration, and runtime behavior. Dependencies and compatibility declarations determine composition order.
- Kit sets let authors combine a workload and mixins, pin their component versions, and publish the result as a single OCI reference. Sets can also add capabilities, lifecycle hooks, instructions, and arguments of their own.
- V3 kits can scope network access by HTTP method and path, declare install-phase network access and credential use, and specify where an agent reads shared skills.
- Multiple OAuth-backed agents can be composed in the same sandbox with credentials scoped to the kits that request them.
- HTTP Basic credentials declared by a kit now produce the expected `Authorization: Basic` header. Composition fails when kits declare conflicting ownership of a Basic-auth service instead of silently dropping the username.
- `sbx kit validate` now rejects malformed API-key declarations, including invalid names, missing injection domains, invalid format placeholders, and Basic-auth usernames containing a colon. It also warns about declarations that have no effect or target domains outside the kit's network allowlist.
- Reusing an unchanged local kit no longer rebuilds its composed image.
- Adding a mixin to an existing sandbox through the daemon API now writes the mixin's agent instructions as expected.

#### Agents and models

- `sbx run --model` can use any OpenAI- or Anthropic-compatible endpoint configured in the new `model.providers` setting.
- `sbx run opencode --model` now exposes the model's supported thinking levels as OpenCode variants, selectable with <kbd>Ctrl</kbd>+<kbd>T</kbd>.
- The OpenCode kit now configures GitHub Copilot from the account's stored GitHub credential, so Copilot models work without a separate device login.
- Codex sandboxes now install Codex with its native installer instead of npm.

#### CLI and output

- MCP server, secret, skill, template, volume, and policy-profile list commands now support `--quiet` (`-q`) for name-only output.
- List commands now use consistent table formatting, and errors use a consistent format with clearer recovery guidance.
- Commands that remove resources now ask for confirmation. Use `--force`, or `--yes`/`-y` for `sbx kit builder history rm`, in non-interactive workflows. Declining a destructive-action or required-restart prompt now returns a non-zero exit code.
- Running `sbx secret rm` without a service opens a picker showing existing local secrets and their scope, type, and name.
- Unsupported detached execution with `sbx exec -d` or `--detach` now fails immediately instead of running in the foreground.
- `sbx settings` now appears in `sbx --help` and the CLI reference.
- `sbx ls --json` now includes `created_at`. `sbx ls --json` and `sbx inspect --json` also report recorded CPU and memory limits for local sandboxes.
- The updater no longer asks to switch channels when the requested version is already installed.
- `sbx env rm` now warns about data loss for a cloned workspace before asking for confirmation.
- `sbx logout` no longer warns about stopped sandboxes when the daemon isn't running.

#### Sandbox lifecycle and workspaces

- Sandboxes now recover when the guest kernel crashes instead of becoming permanently unusable. If a guest stops responding, affected operations fail with an explanation, held proxy connections are released, and `sbx ls` and `sbx inspect` report the unresponsive state.
- Dynamic mounts are restored after a sandbox restart. Startup fails clearly if a saved mount can't be restored, and `sbx umount` can remove a saved mount while the sandbox is stopped. A missing unmount target no longer disrupts existing mounts.
- Clone-mode sandboxes restore their host Git remotes on every restart, preserve complete remote configuration during concurrent lifecycle operations, and provide recovery instructions if configuration fails.
- Newly created or recreated sandboxes have a writable `/etc/hosts` file.
- Image pulls retry transient registry network failures before sandbox creation fails.
- Cached-image recovery is reported as successful without also showing a registry error, and mount-policy evaluation failures are distinguished from access denials.
- Container swaps remove obsolete registry-mirror allowances even if saving the previous swap state fails.
- Updated containerd to fix image layers being dropped.

#### Authentication and credentials

- Adding, updating, or removing global service secrets now updates existing local sandboxes without a restart while preserving sandbox-specific credentials. Sandbox-scoped command and reference secrets also take effect immediately.
- Registry and service-secret revocation failures are now reported and can be retried, including after a stored OAuth token has been deleted.
- OAuth refreshes are coordinated across sandboxes that share credentials, preventing simultaneous refreshes from forcing another sign-in.

#### Networking and policy

- Network policy now treats hostnames with a trailing dot the same as their canonical form for routing, interception, credential injection, and `host.docker.internal` handling. The policy log also records cleartext HTTP requests whose `Host` header differs from the connection destination.
- Experimental outbound UDP now follows sandbox network policy. New local allow rules cover TCP by default; select UDP explicitly with `--protocol` in the CLI or the TCP+UDP option in the TUI. UDP is refused when the destination requires an HTTP, SOCKS5, system, or PAC-selected proxy, because those proxies can't carry it.
- Reverse-DNS lookups are now allowed only for destination IPs already authorized by policy, including IP, CIDR, and allow-all rules. This closes the previous policy bypass without blocking PTR lookups for permitted addresses.
- DNS resolution is no longer allowed when no network rule permits it.
- Connections allowed only by a CIDR rule no longer wait for hostname detection before connecting, improving protocols such as SSH where the server speaks first.
- Network and filesystem access now fail closed with accurate errors when policy evaluation fails, governance can't be resolved, a policy snapshot is stale, or a request is malformed.
- `sbx policy allow network`, `sbx policy deny network`, `--allow-network`, and `--deny-network` now reject malformed patterns before saving them.
- `sbx policy ls` now shows how each rule was created and supports filtering with `--created-via`.
- Fixed excessive daemon CPU use caused by reading the settings file for every blocked UDP packet.
- The governance-rules table no longer reserves space for a hidden profile column, keeping host values readable in narrow terminals.

#### MCP

- Fixed gateway creation failures caused by parentheses or other sandbox-ID punctuation in generated gateway names.
- MCP gateways now recover correctly after a daemon restart when using `sbx exec` or `sbx env run`.
- Internal MCP discovery and OAuth informational logs no longer appear in normal command output.
- OAuth authorization errors now suggest explicit scopes when the authorization server rejects a request without scopes.
- OAuth metadata discovery for private addresses now warns and continues by default. `--skip-ssrf-check` remains available to suppress the warning for trusted providers and their discovery destinations.
- `sbx mcp add --disable-http2` disables HTTP/2 for a remote MCP transport, providing a workaround for servers whose HTTP/2 handling stalls long-lived streams.

#### Packaging and installation

- macOS distributions now contain a single signed `Sbx.app` bundle. Homebrew and tarball PATH installs continue to work through a symlink into the bundle.

<!-- END GENERATED RELEASES -->

## Earlier releases

For older versions, see the
[Docker Sandboxes releases on GitHub](https://github.com/docker/sbx-releases/releases).
