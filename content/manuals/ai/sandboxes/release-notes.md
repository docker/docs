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

## 0.43.0

{{< release-date date="2026-09-15" >}}

[GitHub release](https://github.com/docker/sbx-releases/releases/tag/v0.43.0)

### What's New

#### Breaking changes

- `shareSkills` in `sbxenv.yaml` ([experimental feature](https://docs.docker.com/ai/sandboxes/configuration/environment-files/)) has been replaced with `skills`, `skills` may be set to `off|readonly|readwrite`.

- MCP OAuth client secrets are renamed to `mcp:<server>:client_secret` (was `mcp:<server>.client_secret`), matching the header-secret naming; a secret stored under the old name is no longer read and must be re-set with `sbx secret set mcp:<server>:client_secret`.

#### Environment files

- Environment files can reference `${{ env.projectDir }}` and `${{ env.fileDir }}`, the user-level `~/.sbxenv.yaml` can mount each project's own directory by declaring `workspace: ${{ env.projectDir }}`, relative workspace paths now resolve against the file that declares them, and every `sbx env` subcommand accepts `--name` to override the sandbox name.
- `sbx env run`, `sbx env create`, and `sbx env rm` detect name conflicts with sandboxes created outside `sbx env` and provide guidance instead of treating them as environment-managed sandboxes.

#### Agents and models

- Formerly built-in agents that moved to public kits (kiro, copilot, droid) can be launched by name again — `sbx run kiro` resolves the pinned replacement kit and its stored credentials work without extra approval steps.
- `sbx run --provider` now accepts hosted models.dev providers, served through llmman.
- `sbx run --model` gains `--overflow-provider`/`--overflow-model` to pair a local model with a hosted one for oversized requests; `--provider` now works with codex for providers lacking the Responses API; codex sandboxes no longer spend seconds retrying WebSocket connections to the local model server.
- Claude Code sandboxes started with `sbx run --model` now use the model you selected instead of the harness's own default model.
- A slow first launch of the bundled llmman no longer fails `sbx run --model`.

#### Kits and skills

- `sbx create`/`sbx run` now share skills read-only by default via a new tri-state `--skills=off|readonly|readwrite` flag; the retired `--no-share-skills` flag still works as a deprecated alias for `--skills=off`. There is also a new `skills.defaultMode` to set the desired default behaviour.
- Commit-pinned git kits now resolve offline from a local content-addressed cache, and every cache hit verifies the checkout against a per-file manifest, so a tampered cache entry is quarantined and refetched instead of being served.
- Signed git kits now verify on every host: a kit checkout is materialized from the commit's blobs alone, so smudge filters, line-ending conversion, LFS, hooks, and other host git configuration can no longer alter the checked-out bytes.
- Fixed kit-argument (`${{ kit.args.* }}`) substitution silently not applying when a kit reference is a symlinked directory.
- Kits can now be installed through registry mirrors configured with an explicit port.
- Hardened git kit cloning against command-line config injection (`GIT_CONFIG_PARAMETERS` and its numbered counterparts) carried in the inherited environment.

#### Sandbox lifecycle and workspaces

- Add a last-used timestamp to `sbx ls --json` and Docker-style `until` filtering to `sbx prune`.
- Sandboxes created with `sbx create` now stop automatically after becoming idle.
- The minimum memory for a sandbox has been decreased to 512 MiB. Note: This is only suitable for shell use cases.
- Sandbox names are now validated to reject names longer than 63 characters or ending in a hyphen or period.
- `daemon inspect` and `inspect` will now show mount information.
- Clone-mode sandboxes now support shallow Git repositories.
- Fixed an issue where the `sbx` CLI could select the wrong repository during Git-related setup tasks, such as loading kits or configuring workspaces, when Git environment variables were set on the host.
- Fix UNC path resolution on Windows so that the same folder is identified correctly.
- SSH connections now remain bound to the original sandbox identity while preserving existing sandboxes during UUID migration.
- Windows clients can now connect to `sandboxd` through filesystem `AF_UNIX` sockets.
- The local daemon now verifies connecting operating-system users on Unix sockets and Windows named pipes.

#### Authentication and credentials

- Docker sign-in now explains how to recover when macOS Keychain denies access to stored credentials.
- Fixed a bug where a single Docker Hub sign-in timeout could permanently lock the daemon out of Docker Hub, requiring a manual sign-in to recover.
- Concurrent sandbox creates now reuse one Docker Hub authentication request.
- Private registries can use an explicitly trusted cross-host authentication endpoint for sandbox pulls.
- `sbx secret rm --sandbox` now immediately revokes the removed credential from the sandbox proxy.
- Prevent `sbx exec` from synchronizing credentials that were not configured for the sandbox.
- Credential-binding consent now defaults to decline and clearly identifies when API-key secrets will be sent to new domains.
- Fixed the OAuth credential gate so a third-party kit re-declaring a built-in agent's OAuth service can no longer inherit that agent's trust and receive a real token without an explicit binding; sandboxes created before this fix now self-heal on the next daemon restart or kit add instead of requiring a manual recreate.
- Unrelated credentials no longer switch Claude sandboxes into Anthropic API-key mode.
- `sbx secret import` and `sbx secret ls` now list copilot's GitHub credential correctly.

#### Networking and policy

- The CLI honors configured proxy settings for host-side HTTP requests, including login, diagnostic uploads, and update checks.
- Fix HTTP/2 upstream responses without bodies being incorrectly framed as chunked by the sandbox proxy.
- Hardened credential handling in the sandbox egress proxy so a client-supplied credential the proxy did not issue is not forwarded to managed provider hosts.
- Network allow rules for IP-literal targets (for example `sbx policy allow network [::1]:8080` or CIDR rules such as `10.0.0.0/8`) are enforced correctly again; the proxy no longer blocks them with a default-deny after the governance approval-callback migration.
- Fixed: agents no longer suggest `sbx policy allow` for a host blocked by an org-governed default-deny policy — it now reads as `Blocked by org policy`, same as an explicit org deny rule.
- The `balanced` policy preset now allows access to the NodeSource APT repository.
- Claude sandboxes can now access the Claude Code documentation.

#### MCP

- `sbx mcp add` can now send custom request headers to remote MCP servers via `--header`, with header values substituted from the local secret vault.
- MCP authorization supports private OAuth discovery with `--skip-ssrf-check`, falls back to advertised common OIDC scopes, and honors `--no-scope` for local OAuth registrations.

#### CLI, diagnostics, and updates

- Local `sbx` commands no longer wait on slow or unreachable update services before exiting.
- Plain `sbx version` invocations now return embedded version information without full CLI startup.
- `sbx diagnose` checks whether `mkfs.erofs` is usable and warns if its default block size exceeds the sandbox kernel's page size.
- `sbx diagnose` no longer reports a missing SSH `ProxyCommand` in Git Bash when the Windows OpenSSH configuration is healthy.
- The `tls.allowNegativeSerial` setting no longer prints an informational log line on every `sbx` command while remaining visible in daemon diagnostics.
- Terminal output now uses default text colors when the background theme cannot be detected.
- Fix terminal cursor flickering issue on Windows.
- Nightly and development builds now report a version based on the latest stable release instead of a release-candidate tag.
- `sbx@rc` brew users on macOS will be updated to the latest stable build when it is released.

<!-- END GENERATED RELEASES -->

## Earlier releases

For older versions, see the
[Docker Sandboxes releases on GitHub](https://github.com/docker/sbx-releases/releases).
