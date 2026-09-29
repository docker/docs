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

## 0.42.1

{{< release-date date="2026-09-07" >}}

[GitHub release](https://github.com/docker/sbx-releases/releases/tag/v0.42.1)

### What's New

#### Bug Fixes

- Fix HTTP/2 upstream responses without bodies being incorrectly framed as chunked by the sandbox proxy.

## 0.42.0

{{< release-date date="2026-09-07" >}}

[GitHub release](https://github.com/docker/sbx-releases/releases/tag/v0.42.0)

### Highlights

- BREAKING: `sbx ports --publish` and kit-declared ports now default to `tcp4` instead of dual-stack `tcp`, so a published port no longer listens on `::1` unless you name the protocol explicitly (`--publish 8080:3000/tcp`); this makes `http://localhost:<port>/` reach a sandbox service that listens only on IPv4.
- `sbx run` and `sbx create` now accept sandbox kit references as the agent positional: `sbx run <sandbox-kit-ref>`. The old form `sbx run <sandbox-kit-name> --kit <sandbox-kit-ref>` is deprecated; use the `--kit` flag for mixins.
- Sandboxes can now be created without a workspace bind mount by omitting the path in `sbx create`. Note that this only affects the `create` command; `sbx run` still defaults to mounting the current directory as the primary workspace.

### Security

- Fixed [CVE-2026-77179](https://www.cve.org/cverecord?id=CVE-2026-77179), a symlink vulnerability in the virtio-fs host server on macOS that could let a malicious guest read or modify arbitrary host files outside the shared workspace, potentially leading to code execution on the host.
- Fixed [CVE-2026-79994](https://www.cve.org/cverecord?id=CVE-2026-79994), a symlink race in the guest-to-host Unix domain socket relay that could let a malicious guest connect to arbitrary host Unix sockets outside the shared workspace, exposing data or host-side capabilities.

### What's New

#### CLI

- Read-only `sbx` commands including `secret ls`, `version`, `mcp ls`, `skills ls`, `policy inspect` and the `kit` verification commands now accept `--json` for machine-readable output.
- Clipboard commands inside local sandboxes can now copy text to the host clipboard.
- Add, update, list, and remove sandbox skills directly from Git repositories with `sbx skills`.

#### Environment files

- `sbx env` now reads a non-hidden `sbxenv.yaml` from a project directory and no longer falls back to a hidden `.sbxenv.yaml` there; it merges a `.sbxenv.yaml` from your home directory beneath the project file as defaults shared across projects.
- `sbx env` now shows a plan of everything an environment file changes on the host — host `lifecycle:` commands, credentials, bindings, MCP servers, workspaces, kits, ports and the sandbox itself — asks before applying it and asks again for every run of a command on this machine unless `env.rememberHostCommands` is set, binds the environment file read-only into the sandbox it describes, and reads a directory for `sbxenv.yaml` alone with `~/.sbxenv.yaml` as the user-level base beneath it.
- `sbx env`: an environment file that declares no `workspace:` now creates a sandbox with no workspace bind mount instead of mounting the directory holding the file; write `workspace: .` to mount the project directory.
- Environment files can now declare their own arguments in an `args:` block, referenced as `${{ env.args.NAME }}` and supplied with `sbx env --env-arg`; `${VAR}` interpolation in `.sbxenv.yaml` is no longer expanded.
- Relative kit paths in an environment file now resolve against the file's directory instead of the directory `sbx` was run from.
- `sbx env create` now shares imported skills by default and accepts display, GPU, and USB options in `sbxenv.yaml`.

#### Daemon

- Sandboxes now get a 10 GB Docker volume instead of 50 GB, which significantly reduces host disk usage; set `DOCKER_SANDBOXES_DOCKER_SIZE` to change it.

#### Agents

- Added Devin as a built-in agent.

#### Kits

- Kits can now declare their arguments in an `args:` block and receive values with `--kit-arg name=value`, or `--kit-arg kit.name=value` to target a single kit.
- A `kits:` entry in sbxenv.yaml carries the arguments for that kit under `kits[].args`.

#### Bug fixes

- Docker Sandboxes no longer opens the setup wizard automatically; run `sbx setup` to launch it explicitly.
- Fixed a vulnerability where a sandboxed process could get the daemon to open a host D-Bus transport and execute an arbitrary command on the host.
- On macOS, sbx now accepts a workspace path whose casing differs from the spelling on disk instead of failing to create the sandbox.
- Fixed a rare case where a spotty network right after your computer woke
from sleep could cause an unexpected Docker Hub sign-out.
- Agent crashes now identify the terminating signal and provide scoped recovery guidance.
- Fixed a vulnerability where a malicious sandbox could hijack another sandbox's OAuth login by pre-claiming its callback port.
- Removing or pruning a local sandbox now also deletes its sandbox-scoped secrets.
- `sbx secret ls` no longer prints a stored secret unmasked when its value happens to match one of the status labels the listing displays.
- `sbx` now warns when a stored credential is not sent to a sandbox because no binding authorizes it, instead of starting the sandbox and failing later with an authentication error.
- Fixed several MCP-related bugs.
- Standardized error message formatting for `sbx rm`, `sbx stop`, MCP authorization, and `sbx reset`.
- Sandbox and agent not-found errors now use one sentence pattern and quote style across commands: sandbox '<name>' not found.
- Docker Hub template pulls created through the TUI now use your Docker Sandboxes login credentials.
- `sbx ports --publish` now automatically starts stopped local sandboxes before publishing ports.
- Docker volume sizes below 512 MiB are now rejected before sandbox creation.
- Creating a sandbox from a Docker Hardened Image template no longer results in a delay.
- Fixed OAuth authentication for custom agent kits that declare resource hosts without a fallback API key.
- Kits can now set a sandbox's CPU and memory limits through the `sandbox.resources` block in their spec.
- Fixed SSH connections from editors by keeping non-interactive probes quiet and delivering their exit status before closing the channel.
- Deleting a sandbox now reliably reclaims its disk volumes, and creating a new sandbox that reuses a deleted sandbox's name no longer inherits its files, Docker images, or agent session history.
- Running `sbx setup` explicitly no longer causes the setup screen to appear again on the next interactive command.
- Fixed sandbox connections to a server that sends data first — including passive FTP transfers and a serial console relayed to the host — failing with a timeout instead of receiving the server's output.
- `sbx kit push` no longer uploads an empty payload layer for kits that ship no files.
- `sbx kit push` now authenticates from the sbx credential store, so a single `sbx login` or `docker login` is enough for pushing, signing, and attaching provenance.
- Kits using `extends:` now inherit the parent's setup commands, credentials, network allowlist, volumes, and environment variables instead of replacing them when the child declares its own.
- Fixed Docker Hub credential refresh retrying without backoff after a rate limit, and a non-interactive login discarding the stored OAuth refresh token.
- Nightly Homebrew installs no longer fail with checksum mismatches while a nightly publish is in flight: the sbx@nightly cask now downloads from immutable per-build release URLs.

#### Other

- Docker Sandboxes now provides a machine-wide Windows MSI for administrator-managed installations.
- A declined `@requireApproval` prompt on a local MCP server now gets its own audit record, with policy attribution and a context digest, instead of leaving the original approval-required decision as the only trace of the exchange.
- A registry mirror configured with `platform.images.registryMirror` is now also used by Docker running inside a sandbox, when the value is a bare host (no path prefix) that is not a loopback or wildcard address.
- `sbx version --json` now reports a `server.state` of `running` or `unavailable`, so scripts can check whether the backend was reachable without parsing the error text.
- SSH agent forwarding can be explicitly disabled and use either each client's current agent socket or a fixed socket path.
- `sbx` terminal output now adapts colors for light terminal backgrounds and
uses a distinct pink spinner glyph; CJK and combining-mark column widths in
table output are now measured correctly; piped and JSON output is unchanged
for ASCII-only content.
- `sbx mcp add` now accepts `--skip-auth` (the old `--skip_auth` still works), and a `--url` on a private, loopback, or cloud-metadata address is resolved and registered with a warning instead of being rejected.
- On Linux hosts without an available OS keychain, newly stored secrets are now read and written much faster; secrets already on disk keep their previous cost until they are next written.
- Fixed a gateway defect where a remote MCP server's reconnect could silently wipe its tool routing, causing the server's tools to disappear from agents and be denied by organization MCP policy as unrecognized built-in tools until the daemon was restarted.
- Docker Sandboxes can now upload a diagnostics bundle automatically when the daemon hits an error, after you opt in.
- Governance resolution issues are now shown in sbx policy output and the dashboard.
- `sbx mcp auth` now requests only the scopes you chose (or the set the resource itself requires, suppressible with the new `--no-scope` flag) instead of every scope a server advertises, explains which scopes a server refused along with a narrower retry command, and reports the scope sets of an existing grant in `sbx mcp auth status`: what was granted, what was requested, and what the server supports — with scopes sorted, duplicates collapsed, and differences such as unrequested or no-longer-advertised grants called out.
- Sandbox listing and creation output now share one rendering package; on a terminal, `sbx ls` column headers are now styled bold.
- Sandbox agent instruction files no longer include generic language-specific development guidance.
- The sandboxd runtime state directory left behind by versions before v0.25.0 is now migrated to its current name instead of being used in place.

<!-- END GENERATED RELEASES -->

## Earlier releases

For older versions, see the
[Docker Sandboxes releases on GitHub](https://github.com/docker/sbx-releases/releases).
