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

<!-- END GENERATED RELEASES -->

## Earlier releases

For older versions, see the
[Docker Sandboxes releases on GitHub](https://github.com/docker/sbx-releases/releases).
