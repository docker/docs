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

## 0.39.0

{{< release-date date="2026-08-19" >}}

[GitHub release](https://github.com/docker/sbx-releases/releases/tag/v0.39.0)

### Highlights

**Declarative sandbox environments.** Define a complete, reproducible sandbox in a `.sbxenv.yaml` file, including the agent, workspace, kits, environment variables, secrets, registry credentials, ports, and resource limits. Commit the file with your project so contributors can launch the same environment with `sbx env run`. This feature is experimental.

### What's New

#### Sandbox environments

- Use `sbx env run` to provision an environment from `.sbxenv.yaml` and open an interactive session.
- Use `sbx env create`, `sbx env exec`, and `sbx env rm` to manage the environment lifecycle.
- Combine multiple environment files for shared configuration and local overrides.
- Reference host environment variables in environment files for machine-specific paths and credentials.
- See the [sandbox environment files documentation](https://docs.docker.com/ai/sandboxes/sandbox-environments/).

#### CLI

- Add an experimental `--usb` flag to `sbx create` behind the `DOCKER_SANDBOXES_FEATURE_SANDBOX_USB` environment variable to re-attach the specified USB devices. They will be available inside a sandbox via usbfs. Linux x86_64/ARM64 only.
- sbx run --model now selects the Ollama backend via a new `--provider ollama` flag instead of an `ollama/` prefix on the model name.
- Stopped sandboxes can now be cleaned up in bulk with `sbx prune`, which never removes a running sandbox and can filter on how long each has been stopped.
- `sbx run` and `sbx create` now accept `-e`/`--env` and `--env-file` to set environment variables in a sandbox, following `docker run` precedence rules.

#### Secrets

- `sbx secret set` and `sbx secret set-custom` can now configure dynamic secrets that resolve values from a reference or command, with options to control refreshing, verification, and error output.
- On Linux hosts without an available OS keychain, newly stored secrets are now read and written much faster; secrets already on disk keep their previous cost until they are next written.

#### Daemon

- Sandboxes now expose their own identity as `SANDBOX_NAME` and `SANDBOX_ID` environment variables, matching the name and id shown by `sbx ls --json`; the older `SANDBOX_VM_ID` still carries the sandbox name but is deprecated.

#### Networking

- Claude Code's `/remote-control` can now be used inside sandboxes by enabling `claude.remoteControl` setting: `sbx settings set claude.remoteControl true`.

#### Bug Fixes

- sandboxd now removes the sandbox container immediately when container startup fails, so an interrupted `sbx create`/`sbx run` is less likely to leave the sandbox name unusable.
- Agent kits that declare a persistent volume without a size now get a 512 MB volume instead of a 50 GB one, which significantly reduces sandbox disk usage on the host.
- `sbx` now reports a clear error for an unrecognized command, subcommand, or `sbx help` topic instead of printing help and succeeding, and reports a mistyped command without first asking an unauthenticated user to sign in.
- Claude sandboxes now use around 3.9 GB less disk space on the host.
- Claude sandboxes can connect to required Anthropic services when using the locked-down network policy.
- `sbx kit inspect` now describes kits using kit-spec v2 field names and lists any deprecated fields a kit still relies on, and `sbx kit validate` now rejects OAuth credentials missing sentinels, a service, or a credential-file body.
- DNS lookups in a sandbox now succeed for any host that network policy allows on any port, including hosts allowed only on a non-standard port such as `myhost:2222`.
- `sbx template load` now fails with an error when an image import does not complete, instead of reporting success.
- Correct the `sbx create --name` help text and CLI reference, which incorrectly listed plus signs as valid sandbox-name characters and omitted the leading-alphanumeric and two-character-minimum rules.
- `sbx reset` now removes the Docker Sandboxes-managed block from `~/.ssh/config`.
- `sbx` now reports the exit code when a sandbox container dies at startup, and rejects a template image built for a different CPU architecture with a clear message instead of failing after a 30-second wait.
- Signing in to Claude Code with an Anthropic Console API key now succeeds on repeat logins instead of failing with a 401 error.
- Fixed `sbx cp` failing on Windows when the local path has no directory component (e.g. `sbx cp file.txt sandbox:/tmp/`).
- `sbx daemon restart` now starts the daemon again after a stop that reports a failure but leaves no daemon running.

#### Other

- `sbx diagnose` now reports free disk space on the volume holding sandbox data, and diagnostics bundles include host disk totals.
- `sbx diagnose` now detects broken, shadowed, or stale SSH client configuration.
- Add a `platform.images.registryMirror` setting that redirects Docker Hub-resolving sandbox template and kit images to an organization's registry mirror.
- Filesystem policy denials now include the organization's support contact message, matching network denials.
- sbx now reports when the host cannot provide a hypervisor — including a Windows installation running inside a virtual machine without nested virtualization — instead of a generic "failed to run sandbox container" error, and `sbx diagnose` now checks host virtualization support.
- Kits can now be signed and verified with cosign-compatible Sigstore signatures via `sbx kit sign` / `sbx kit verify`, with optional policy enforcement at load time.
- OAuth kits can declare their credential file with the declarative `credentialFile.structure` form, rendered to well-formed JSON, instead of a free-form Go template.
- Ubuntu 25.10 packages are no longer published; Ubuntu 25.10 is end-of-life.

## 0.38.0

{{< release-date date="2026-08-06" >}}

[GitHub release](https://github.com/docker/sbx-releases/releases/tag/v0.38.0)

### Highlights

**Kit spec v2.** A new schema is available for authoring kits, with a clearer structure for setup, permissions, agent instructions, networking, and credentials. Use `schemaVersion: "2"` for new kits; existing v1 kits continue to load through the legacy path. See the [kit spec reference](https://docs.docker.com/ai/sandboxes/customize/kit-reference/#schema-versions) for migration details.

**MCP management is now a first-class feature.** Register remote or local MCP servers once with `sbx mcp`, then reuse them across supported agents and sandboxes through a built-in MCP gateway. OAuth credentials stay on the host, and organizations can govern server registration and tool calls with Cedar policies. See the [MCP gateway documentation](https://docs.docker.com/ai/sandboxes/mcp-gateway/).

### What's new

#### CLI

- `sbx create` and `sbx run` show detailed structured progress during startup, including environment files loaded, resources provisioned, and each kit command's outcome; kit-install progress streams live during `sbx create --kit`.
- Added `sbx daemon restart` to stop and restart the sandboxd daemon in the background.
- `sbx inspect` now displays custom secrets configured for a sandbox.
- `DOCKER_SANDBOXES_CLONED_WORKSPACE_SIZE` configures the size of the cloned workspace volume.
- `sbx setup ssh` warns when `ssh` is missing from PATH, and on Windows when `sh` (required by Claude Desktop's SSH ProxyCommand) is missing.
- Port publishing failures now identify the affected host port and explain when the OS requires extra daemon privileges.

#### MCP

- The `sbx mcp` subcommand is now available for managing MCP servers.
- Includes dynamic MCP tools (`mcp-find`, `mcp-add`, `mcp-config-set`) for attaching registered servers to sandboxes.
- Govern MCP servers and tools for your organization using Cedar policies.

#### Networking & policy

- `--deny-network HOST` on `sbx run` and `sbx create` records per-sandbox network deny rules at creation time, with layer-aware egress messages.
- `sbx policy allow network` reports a clear "managed by your organization" error when org governance overrides the local allow, and failed rule removals now explain what went wrong using a single rule identifier.
- Signing in refreshes organization policies in the running daemon immediately instead of waiting for the next polling interval.
- IP-literal destinations denied by a CIDR rule fail fast with a policy message instead of timing out.
- Blocked HTTPS proxy connections appear in `sbx policy log` even when the client aborts the TLS handshake.

#### Secrets & credentials

- Service and custom secrets are global by default, with `--sandbox` for sandbox scope; legacy positional and `--global` forms are deprecated with warnings.
- Sandbox-scoped GitHub credentials added after creation now work without recreating the sandbox.
- Pressing Ctrl+C while entering a secret cancels the command without saving it.

#### Agents

- Docker Agent and OpenCode sandboxes can authenticate GitHub Copilot requests with proxy-managed GitHub credentials.
- Codex sandboxes created from the TUI prefer stored OpenAI OAuth credentials over API keys; kit environment variables now reach cloud agents, and git no longer hangs Codex startup prompting for credentials.
- Shared agent skills: directory symlinks under the skills folder are resolved and their contents imported.

#### Kits & templates

- Kit specs use the new v2 grammar.
- Kits using `extends` correctly inherit and override the base image or build source of their parent.
- Kit install commands can consume static files from `files/home`, including binary files.

#### Packaging

- Homebrew installs from a stapled `.dmg` artifact rather than a `.tar.gz` archive, improving Gatekeeper compatibility on macOS.
- Windows: the running sandboxd daemon is stopped during a WinGet/MSI upgrade so client and server end up on the same version.

#### Security

- Claude Desktop SSH sessions no longer expose Desktop OAuth access tokens inside sandboxes.
- Fixed a destination-escape flaw in `sbx cp` copy-out (CVE-2026-17106).
- The daemon's loopback egress proxy only serves the daemon's own traffic, preventing other local users on a multi-user host from reaching the configured upstream proxy through it.

#### Bug fixes

- Fixed a hang where sandboxd stopped answering all endpoints and could not be stopped without SIGKILL after a crash; fatal daemon tracebacks are now included in `sbx diagnose --upload` bundles.
- Fixed intermittent sandboxd startup failures when a running daemon was slow to answer its health check.
- Fixed `sbx daemon stop` hanging when an idle SSH session (for example, Claude Desktop) was connected to a sandbox.
- sandboxd automatically repairs a corrupted local image cache by re-pulling the image, and otherwise reports a clear "run `sbx daemon reset`" error.
- Fixed recreate failures ("base image not found") after daemon restarts when swapping a sandbox's container via `sbx kit add`; recreates self-heal by recomposing from the sandbox's template.
- Fixed reverse DNS (PTR) lookups from sandboxes returning NXDOMAIN for container-resolved addresses.
- Fixed a goroutine and network-endpoint leak from hijacked HTTP CONNECT tunnels that could eventually stall sandbox creation after many delete/recreate cycles.
- Fixed an intermittent 500 error when deleting a sandbox while its network endpoints were being torn down.
- The daemon restores saved sandboxes' network proxies in parallel on restart, speeding up startup with several sandboxes and fixing a potential crash during first-run policy application.
- Fixed host `.git/config` corruption when creating a sandbox for repositories using `includeIf` directives in `~/.gitconfig`; the sandbox now writes git identity only to the container's gitconfig.
- `sbx skills` shows a single usage form and clearer help for importing shared agent skills.
- sbx no longer reports that a stored credential was not injected when the daemon injects it.

### Experimental features

#### Enterprise networking

Settings-driven upstream-proxy configuration with separate sandbox and daemon scopes, integrated NTLM/Kerberos proxy authentication on Windows.

- Configure separate proxy settings for sandbox and daemon traffic using `proxy`, `proxy.sandbox`, `proxy.daemon`, and the matching `no_proxy` settings. These default to the host operating system's proxy settings. The daemon's own traffic, including image pulls and telemetry, also uses the configured proxy.
- On Windows, sbx can authenticate to upstream proxies that require integrated NTLM or Kerberos/Negotiate authentication. Enable this behavior with the `proxy.integratedAuth` setting.
- If a TLS-inspecting proxy issues certificates with negative serial numbers, enable compatibility with `sbx settings set tls.allowNegativeSerial true`, then restart the daemon.

#### GPU passthrough

Run a sandbox with NVIDIA VFIO GPU passthrough on Linux using `sbx run --gpu`. Enable this feature with `sbx settings set feature.sandbox-gpu true`.

#### Local models

Run Claude Code against a local GGUF model with `sbx run --model <name> claude`. To use a model from an existing Ollama installation, prefix the model name with `ollama/`. See [Claude Code > Use a local model](https://docs.docker.com/ai/sandboxes/agents/claude-code/#use-a-local-model).

<!-- END GENERATED RELEASES -->

## Earlier releases

For older versions, see the
[Docker Sandboxes releases on GitHub](https://github.com/docker/sbx-releases/releases).
