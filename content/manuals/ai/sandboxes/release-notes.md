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

## 0.48.0

{{< release-date date="2026-10-13" >}}

[GitHub release](https://github.com/docker/sbx-releases/releases/tag/v0.48.0)

### Highlights

Docker Sandboxes v0.48.0 updates the CLI across all commands, with consistent terminal output and JSON/YAML formats. This changes flags, structured output, and prompts used by existing scripts. Review the migration details below before upgrading automated workflows.

The release also adds managed Codex and Claude Desktop SSH targets on macOS and Windows, additional v3 kit capabilities, and support for adding kits with startup commands to existing sandboxes.

### Breaking changes

#### CLI flags, output, and automation

The CLI refactor changes both interactive use and scripting interfaces. Existing `--json` parsers may need updates even if their command lines still work. Scripts that parse human-readable output must switch to `--format json` or `--format yaml`.

##### Update command flags

| Command | Previous spelling | Replacement or behavior |
| --- | --- | --- |
| `sbx run`, `sbx env run` | `--detached` | Use `--detach` or `-d`. `--detached` is rejected. |
| `sbx logout` | `--yes`, `-y` | Use `--force` or `-f`. |
| `sbx env plan`, `sbx env create`, `sbx env run` | `--auto-approve`, `-y` | Use `--force` or `-f`. |
| `sbx create` | `--quiet`, `-q` | Use `--no-progress`. This suppresses progress, but still prints the creation summary. |
| `sbx template save` | `-d` for the description | Use `--description`. |
| `sbx diagnose` | `--output`, `-o` | Use `--format table`, `--format json`, `--format yaml`, or `--format github-issue`. |
| `sbx secret ls`, `sbx secret rm` | `--env` | Use `--env-name`. |
| `sbx policy log` | `--quiet`, `-q` | Removed; there is no replacement. |

`--format`, `--json`, and `-q` are mutually exclusive on commands that use the shared output flags. Even `--json --format json` is rejected. Unsupported format values fail during flag parsing instead of falling back to human-readable output.

`sbx settings ls` is the primary spelling; `sbx settings list` remains an alias. `sbx kit builder history rm` keeps `--yes` / `-y` because its `--force` flag is forwarded to Buildx. The deprecated `--global` and `mcp add --skip_auth` flags remain accepted but hidden, with a deprecation warning on stderr; use `--skip-auth` instead of `--skip_auth`.

##### Update structured-output parsers

`--json` is an alias of `--format json` and emits the same document. JSON output structures have changed for the commands listed in the following table. Output uses `snake_case` field names; YAML is derived from the same representation.

| Command | Output change |
| --- | --- |
| `sbx ls` | Returns the sandbox array directly instead of `{"sandboxes": [...]}`. Fields use `snake_case`, `created_at` and `last_used_at` use RFC 3339 timestamps, and rows include `detached`. |
| `sbx template ls` | Returns the image array directly instead of `{"images": [...]}`, for local and cloud sandboxes. |
| `sbx mcp ls` locally | Returns the server array directly instead of `{"gateway": {...}, "servers": [...]}`. Each server row carries `gateway`. |
| `sbx policy log` | Local and cloud output use `{updated_at, allowed: [...], blocked: [...]}`. Entries contain `sandbox_id`, `host`, `proxy_type`, `rule`, `reason`, `count`, `first_seen`, and `last_seen`. `--limit` applies. This replaces local `allowed_hosts` / `blocked_hosts` and the cloud camelCase document. |
| `sbx --cloud policy ls` | Field names change from camelCase to `snake_case`. |
| `sbx ports` | Cloud rows use `sandbox_port` instead of `port`, alongside `url`. Local rows use `host_ip`, `host_port`, `sandbox_port`, and `protocol`. |
| `sbx ports --publish` / `--unpublish` locally | Structured formats emit result rows with `action: published` or `action: unpublished` instead of human-readable lines. Warnings go to stderr. |
| Cloud `volume create`, `volume ls`, and `volume inspect` | Fields such as `ID`, `Name`, and `CreatedAt` become `id`, `name`, and `created_at`. `volume inspect` defaults to a human-readable details view; request `--format json` for JSON. |
| `sbx rm`, `sbx stop`, `sbx template rm`, `sbx volume rm` | Structured formats return an array of results with `action`, `resource`, `name`, and optional `id` and `error`. Actions are lowercase: `removed`, `stopped`, `removing`, `stopping`, or `failed`. Missing and failed targets are included. |

For example, a script that extracts sandbox names from `sbx ls --json` must change its filter from `.sandboxes[].name` to `.[].name`.

Local and cloud output follow the same conventions. MCP list output still represents different resources: local registrations versus server usage across cloud sandboxes, so those rows differ beyond their shared `name` and `type` fields.

`sbx prune --dry-run` accepts `--format json` and `--format yaml`. Cloud `ttl --format yaml` emits YAML instead of JSON, and `sbx diagnose` supports YAML. `sbx events --json` remains an NDJSON stream, and `sbx mcp auth --format text|json` retains its command-specific formats.

##### Update non-interactive workflows

- Confirmation prompts require terminal input. Piping an answer, such as `echo y | sbx secret set …`, no longer confirms the operation. Use `--force` / `-f` where the command provides it. Commands that require confirmation fail with exit status 1 when stdin is not a terminal.
- Declining a confirmation returns a non-zero exit status consistently, including for removal, reset, and prune operations.
- Read-only invocations and runs with `--force` do not start an interactive sign-in flow. Sign in before running automation that requires authentication. Read-only examples include `ports SANDBOX`, `ttl SANDBOX`, `prune --dry-run`, and `diagnose`; their modifying or upload variants are not read-only.
- Without a terminal, `sbx diagnose --upload` saves the diagnostic bundle but skips uploading it and exits with status 0. A successful exit does not mean that an upload occurred.
- Without a terminal, `sbx create` and `sbx run` fail if the workspace directory does not exist and print a `mkdir -p` hint. Create the directory before invoking the command.
- Ctrl-D at a selection or free-text prompt cancels instead of accepting the default. Enter still accepts the default.

##### Terminal output and streams

Human-readable output changes across all commands: headings, status marks, indentation, progress, and result messages use a consistent format. `sbx ls` includes state marks and a creation-time column; `sbx ls` and `sbx inspect` identify detached sandboxes.

Results go to stdout. Errors, warnings, corrective hints, daemon-start progress, image pulls, and kit-build progress go to stderr. Update scripts that depend on exact message text or capture only one stream.

### What's new

#### Security

- Fixed sign-in enforcement ignoring `allowedOrgs` on macOS.
- `sbx cp` blocks local paths that escape a sandbox workspace through symbolic links.
- Cloud `sbx cp` downloads and `sbx kit pack` preserve restrictive destination permissions without temporarily exposing the replacement file. They account for permissions at completion and do not widen them. On Unix, cloud `sbx cp` refuses to replace files owned by another user; on Windows, permission preservation covers the read-only attribute. Special permission bits are retained where the kernel permits. This cloud-copy change does not affect local `sbx cp`.
- `sbx kit add` checks the filesystem mount policy before recreating a sandbox and refuses the operation if its workspace or a saved mount is denied, matching `sbx start`.
- Fixed the daemon crashing on a malformed request sent through an intercepted HTTPS connection.

#### SSH and desktop integration

- On macOS and Windows, `sbx setup ssh --auto-create` registers Docker Sandbox connections in Codex and Claude Desktop. Each connection creates a dedicated sandbox on first use and reuses it on subsequent connections. These sandboxes have no host workspace mounted. Clone or copy your project into the sandbox to work on it.
- Desktop SSH setup warns about missing authentication and provides local sign-in guidance for Codex and Claude. Repeating `sbx setup ssh` recognizes existing registrations even while the applications are running.
- Closing a Cursor or VS Code window connected over SSH no longer stops the daemon and other running sandboxes.
- SSH connections no longer fail because the CLI's initial feature lookup timed out before starting the daemon.
- SSH-agent forwarding supports Windows OpenSSH and 1Password agents when `SSH_AUTH_SOCK` is not set.

#### Kits

- V3 kits support capability groups, `kit.env` references, bundled agent skills, and the `ssh-agent@1` and `git-identity@1` capabilities. Selected optional groups retain their network rules, lifecycle files, and context; unsupported optional groups are skipped. Bundled skills can use directories outside the agent's home.
- `sbx kit add` accepts kits with startup commands. Those commands run after the kit is added and on each subsequent sandbox start.
- Sandbox containers expose the provisioned agent's name through `SBX_AGENT` for kit hooks and scripts.
- `sbx kit validate` rejects kits whose credentials route the same host, identifying the conflicting services. The daemon logs the collision during sandbox creation.
- Sandbox creation supports kits whose resolved lock exceeds Docker's container-label size limit.
- `sbx kit pack` excludes its output from the archive, prevents writes into the kit's `files/` directory, and handles Windows 8.3 output paths without renaming the destination or including a previous archive. The output path must end in `.zip` and cannot be a named pipe, socket, or device node.
- An unresolvable kit-file path no longer causes `sbx kit push` to fail.
- Kits that declare an agent-context directory now have their agent instructions file written in that directory instead of beside the workspace, and large workload and mixin kit combinations no longer fail to create with "failed to run sandbox container".

#### Sandbox lifecycle and workspaces

- Creating a sandbox from a local snapshot saved with this version avoids registry requests by default, including when using the daemon API.
- `sbx exec`, `sbx run`, `sbx cp`, and SSH report an unresponsive guest with stop-and-start guidance instead of hanging. `sbx ls` marks the sandbox as not responding.
- Commands against an already-running sandbox reuse guest setup instead of repeating it on every `exec`, `run`, or `cp` invocation. Setup runs once per container start or daemon restart.
- `sbx reset` no longer leaves stray sandbox resources or stalls when another sandbox operation is in progress.
- Sandbox startup identifies files still in use from a previous run instead of reporting a generic failure.
- `sbx exec` accepts empty command arguments.
- `sbx cp` preserves dot-prefixed paths and restrictive directory permissions. On Windows, downloads can replace a read-only destination instead of failing with “Access is denied”.
- `sbx run` and `sbx create` support workspace directories whose names start with a dot.

#### Cloud sandboxes

- `sbx move --to cloud` and `sbx --cloud env create` accept `--ttl` and `--on-timeout restart`, extending always-on support beyond `sbx --cloud create`.
- Fixed `on_timeout` response decoding errors when using `--on-timeout`.

#### Authentication, agents, and MCP

- Sandbox-scoped custom secrets can override global secrets that use the same environment variable name.
- Fixed Docker Hub sign-outs on macOS when the host sleeps or wakes during login-token refresh.
- The built-in Codex kit registers the MCP gateway with supported settings, removing ignored-setting warnings and sending the gateway its Authorization header. Recreate existing Codex sandboxes with `sbx rm` and `sbx create` to receive the change.
- `sbx mcp unload NAME --sandbox SANDBOX` detaches a loaded MCP server without restarting the sandbox's agent.
- MCP connections and discovery probes identify as `docker-sbx`, and backend recovery respects rate limits to avoid reconnect storms.
- `sbx run --model --provider` uses the API key stored with `sbx secret set <provider>` without exposing it to the sandbox. `sbx run --variant` selects a model thinking variant.
- Fixed `sbx run --model` failing with “401 Unauthorized” when pulling from Docker Hub after a browser-based `docker login`.

#### Networking and policy

- The proxy handles trailing-dot hostnames and mixed-case `host.docker.internal` consistently, fixing refused requests and missing credential injection.
- Sandbox image pulls retry temporary registry failures reported as HTTP 5xx errors, including when the original error type is lost across API or progress-stream boundaries.
- Fixed WebSocket connections through the proxy occasionally losing the first bytes sent by the client.
- `sbx policy ls SANDBOX` returns an error when the sandbox does not exist.
- `sbx policy inspect` accepts policy names that match one of their own rule names, such as `default-deny-all`.
- Removing a kit-managed network rule by resource explains that the kit manages it instead of asking for a rule ID.

#### CLI and diagnostics

- Shell completion scripts generate without starting telemetry, reducing generation time.
- `sbx diagnose` checks for `mkfs.ext4`.

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

<!-- END GENERATED RELEASES -->

## Earlier releases

For older versions, see the
[Docker Sandboxes releases on GitHub](https://github.com/docker/sbx-releases/releases).
