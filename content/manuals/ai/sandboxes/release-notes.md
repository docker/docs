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

## 0.53.0

{{< release-date date="2026-11-11" >}}

[GitHub release](https://github.com/docker/sbx-releases/releases/tag/v0.53.0)

### Highlights

Docker Sandboxes v0.53.0 adds local disk usage reports, Claude in Chrome access, and automatic loopback port publication for kits. It also improves recovery after host network changes, sandbox startup, credential handling, and command output for scripts.

> [!NOTE]
> v0.52.0 was skipped. These release notes cover changes since v0.51.0.

### Breaking changes

#### Gemini agent retirement

The built-in `gemini` agent name is retired. Existing Gemini sandboxes keep running, but cannot be moved with `sbx move`.

#### Agent arguments

`sbx run <agent> -- <args>` preserves the image's entrypoint arguments when you pass a prompt or subcommand, matching cloud behavior. If the entrypoint already includes a subcommand, omit that subcommand from your arguments to avoid running it twice. For example, use `sbx run kiro -- "explain this project"` without repeating `chat`. To run a different subcommand, invoke the agent directly with `sbx exec` in a running sandbox.

#### Credential providers and settings

- Built-in secret reference providers run host tools from a temporary directory and ignore empty or relative `PATH` entries. Put provider tools in absolute `PATH` directories, use absolute configuration paths, and configure asdf or mise versions globally instead of relying on project-only selections.
- `DOCKER_SANDBOXES_USE_DHI` is no longer read. Use `sbx settings set platform.images.useDHI true` to select Docker Hardened Image templates.
- `sbx logout` requires confirmation even without a running daemon. Use `--yes` or `-y` in scripts. `sbx --cloud logout` warns that signing out affects the whole account, then performs the same sign-out.
- Unsupported `--detach` and `--privileged` flags have been removed from `sbx exec` and `sbx env exec`. Local exec commands also reject `--detach-keys` instead of ignoring it.

#### HTTP policy matching

HTTP policy matching treats ASCII letters in URL paths and query parameter names and values case-insensitively. Review rules that rely on case to distinguish allowed and denied requests.

### What's new

#### Security and credentials

- Sandbox restarts check filesystem policy against the original static mount sources and reject inconsistent saved mount metadata. If a mount source resolves elsewhere, restore the original host path or recreate the sandbox. For mismatched saved metadata, recreate the sandbox. Save accessible sandbox files before removing it.
- Kit install hooks expressed as argument lists pass each argument literally, without shell interpretation.
- Claude Code Remote Control keeps environment and session credentials on the host and omits remotely supplied environment variables and MCP configuration. Access requires the `claude.remoteControl` setting.
- Diagnostic bundles redact account names, home-directory installation paths, and usernames in archive filenames. Sanitized errors also redact paths containing spaces or Unicode, and human-readable inspection output removes terminal control characters.
- `sbx rm`, `sbx prune`, and `sbx env rm` avoid deleting a different sandbox that takes over a name after removal is confirmed.

#### Kits, agents, and skills

- Enable `claude.chrome` to let sandboxed Claude Code use Claude in Chrome with your stored Claude subscription OAuth login. Install the host browser extension, sign in with the same account, and start Claude Code with `--chrome`. API-key authentication alone is insufficient. Turning the setting off blocks further connections; end existing browser sessions to revoke active access.
- Kits granted `ports.publish` with `approval: none` can request an automatically assigned loopback host port without prompting. These bindings disappear when the sandbox stops. Fixed host ports still require approval.
- Unused source-kit build images are eligible for cleanup after seven days. Images referenced by existing sandboxes, saved templates, or in-progress operations are retained. Collected builds may need to be restored or rebuilt on later use.
- Kits launched from Docker Console receive configured proxy-managed API-key variables, including the Anthropic and OpenAI variables used by agents such as Hermes.
- Codex uses the local model selected with `--model` instead of falling back to OpenAI. Agent launchers wrapped in `tini` handle positional prompts and subcommands correctly.
- With `--clone`, kit install hooks can access the cloned project in the workspace. V3 kits receive requested privileged container access, so the kit builder can start Docker inside the sandbox.
- Kit inspection and validation warn about unsupported static-file paths. Install progress identifies commands without descriptions, and `sbx kit add` explains when a recorded mount is denied or its host path is missing.
- Skill import previews group skills by source directory, show missing agent directories only with `--verbose`, and avoid incorrect overwrite warnings.

#### Sandbox lifecycle and workspaces

- `sbx run` reattaches to the current directory's sandbox when another directory with the same name has taken the default sandbox name. Creation rejects names that differ from an existing local name only in case.
- Environment files use the last matching port, MCP server name, or additional workspace target, including across merged files. A replacement replaces the whole entry, so repeat every field you want to retain. Invalid host IP bindings are rejected before planning.
- Sandboxes inherit the host's global Git identity when the workspace root isn't a Git repository. SSH sessions use the login shell of the sandbox user.
- `sbx run` preserves agent output when stdin is `/dev/null` or a closed pipe. Image-pull retries continue into the agent session when a cached image is used.
- Sandbox creation and restore retry temporary registry failures while preserving the original mirror and DHI settings. Concurrent commands on macOS are less likely to cause startup failures while reading cache metadata.
- Sandboxes can restart when Docker was manually installed inside them but cannot start, unless their configuration requires Docker inside the sandbox.
- Canceling creation, runs, or kit additions during engine checks stops promptly. Shell environment recovery stops hung subprocesses, and quitting the TUI cancels pending shell startup while preserving terminals already opened.
- `sbx cp` supports copying a directory to a new path ending in `/` and reports a missing destination directory before uploading.
- Clipboard image paste ignores images larger than 64 MiB, returning an empty clipboard to the agent.

#### Networking and policy

- Stalled outbound connections close after about 30 seconds by default, allowing clients to reconnect after host network changes such as connecting a VPN. Tune the timeout with `egress.stalledConnectionTimeoutSeconds`, or set it to `0` to disable it. Restart the daemon to apply a changed value to existing sandboxes.
- Fixed `EOF` errors during sign-in and other requests through TLS-inspecting proxies that reject handshakes offering ML-DSA certificate signatures, including some Zscaler, Palo Alto, and F5 configurations.
- Published loopback ports work with WSL2 mirrored networking, restoring access to the clone-mode Git remote.
- Unpublishing an automatically assigned port keeps it unpublished after a sandbox restart. Published-port endpoints bracket IPv6 host addresses in list and inspect output, including JSON and YAML.
- Kit pulls work with registries such as Artifactory and self-hosted GitLab whose token service uses a separate host with a private IP address. GitHub release-asset HEAD requests no longer fail with 401 when a GitHub secret is configured.
- DNS answers include CNAME records for aliased hostnames, supporting tools such as Kerberos. HTTPS proxy connections preserve trailing dots in destination names.
- With UDP egress disabled, the first 20 distinct refused external destinations per sandbox network get individual log entries. Further destinations share a `<further blocked attempts>` entry and remain blocked. Denied host UDP traffic is labeled `host.docker.internal`.
- Empty policy rule listings return an empty list in structured output. Declining `sbx policy reset` returns an error, and policy initialization rejects an explicitly empty sandbox scope.

#### Authentication, MCP, and models

- Docker sign-in retries device-code polling after rate-limit access-denied responses. On Linux, a daemon that started before the keychain was available suggests `sbx daemon restart` instead of repeatedly requesting sign-in.
- Unreachable OAuth callbacks time out after 60 seconds. Canceling browser, clipboard, MCP metadata, or image-pull requests also cancels associated helper work and credential reads.
- Environment files can register host MCP commands for multiple sandboxes with independent approval and revocation. Damaged scoped registrations no longer block unrelated MCP gateways.
- Use `--sandbox` with `sbx mcp inspect` or `sbx mcp rm --registration-only` to select a sandbox-owned registration when names are shared. Recovery preserves other sandboxes' registrations and saved credentials.
- MCP, model-provider, and credential errors include corrective commands for unsupported headers, registration conflicts, missing providers, removal failures, and model-server startup timeouts.
- Secret import cancels when confirmation input ends. Conflicting secret-source flags, ambiguous quiet secret listings, and MCP command arguments supplied without `--command` are rejected instead of being silently accepted.

#### Cloud sandboxes

- Sandbox moves reject occupied explicit destination names before capturing or transferring the source. `--name` sets the exact destination name. Cloud-to-local moves derive the default destination name from the source sandbox name even when selected by ID.
- Moving an already-cloud sandbox with `--to cloud` exits cleanly with an explanatory message. Resume recognizes running lifecycle states consistently, and sandboxes whose container creation fails can be deleted instead of remaining stuck.
- Cloud sign-in failures produce warnings and suggest signing in again when appropriate. Invalid SSH certificate lifetimes are rejected even when a cached certificate is available.
- Failed volume and template removals include failure records in JSON and YAML. Port publishing includes the reported protocol in structured output.
- Cloud secret and policy commands reject unsupported local-only flags and filters instead of silently ignoring them. In particular, cloud secret writes reject `--ref` and `--command` instead of storing a literal value.

#### CLI, templates, and diagnostics

- `sbx df` reports local image, container, and managed-volume disk usage and estimated reclaimable space. Use `--verbose` for individual resources. The estimate isn't the amount `sbx prune` alone removes.
- `sbx template ls --all --digests --no-trunc` reveals untagged and digest-only local images. Remove them by repository digest or image ID. JSON and YAML listings preserve full image IDs.
- `sbx inspect` is listed in help and shell completion. Local JSON and YAML inspection includes the sandbox ID and creation time.
- Update notices, confirmation previews, and diagnostic guidance go to stderr to preserve redirected results. Kit pull, sign, and add commands and local template save and load commands write success results to stdout.
- Shell completion covers more supported values, including pull policies, timeout actions, capture modes, shared skills modes, and move destinations. Template save completion suggests source sandboxes, and template flags omit anonymous image IDs.
- Policy checks and MCP authorization accept `-v` for `--verbose`. MCP listings accept `list` as an alias for `ls`, and SSH setup removal accepts `rm` as an alias for `remove`.
- Commands reject more invalid arguments and unsupported options before doing work. Local exec preserves sandbox environment values when a bare host environment variable is unset. Commands requiring confirmation explain how `--no-input` prevents them from proceeding.
- `sbx diagnose --upload` captures daemon goroutine stacks even when the daemon API is unresponsive and retains more crash information and rotated logs. On macOS and Linux, send `SIGUSR1` to `sandboxd` to save goroutine stacks without stopping it.
- Linux packages recommend GNOME Keyring instead of requiring it. Sandbox creation explains missing or inaccessible `/dev/kvm`. Windows handles redirected AppData paths and stops the daemon during per-user MSI upgrades and uninstalls without opening an invisible console.

## 0.51.0

{{< release-date date="2026-11-03" >}}

[GitHub release](https://github.com/docker/sbx-releases/releases/tag/v0.51.0)

### Highlights

Docker Sandboxes v0.51.0 runs built-in agents from published kits and adds custom mount targets for additional workspaces, environment launch command defaults, and more controls for automation. It also strengthens workspace-copy and credential protections and improves network, MCP, and sandbox startup reliability.

### Breaking changes

#### Built-in agents and shell images

Built-in agent names now resolve to published kits at `docker.io/sbx/<name>:latest`. This applies to Claude, Codex, Cursor, Devin, Docker Agent, OpenCode, Kiro, Copilot, and Droid. `sbx run sbx/<kit>` also accepts published kits, including kits named after built-in agents. The `kit.allowExtractedAgents` setting is removed.

- Sandboxes created by earlier releases from embedded Claude or Codex specifications keep running, but must be recreated before using `sbx kit add`, `sbx run --model`, or `sbx move`. Save files and unpushed commits before removing a sandbox.
- The shell agent now uses a published Debian 13 v3 kit. First use requires registry access. To provide shell sandboxes with provider credentials, add the `shell-credentials` mixin. Existing shell sandboxes keep running.
- V1 and v2 kits cannot be combined with v3 kits. Upgrade older mixins to v3 before combining them with the new built-in kits.
- The `docker/sandbox-templates` images for `codex`, `cursor-agent`, `devin`, `docker-agent`, and `opencode`, including their `-docker` variants, are no longer published.
- The `docker/sandbox-templates:shell` and `:shell-docker` images are no longer refreshed. For custom workloads, use `docker.io/sbx/shell-image` and `docker.io/sbx/shell-docker-image`.

#### Clipboard access and custom secrets

- Copying text from a local sandbox to the host clipboard now requires opting in with the `clipboard.textCopy` setting. Native clients such as Codex are supported.
- Custom secrets without destination hosts no longer substitute credentials. Repair existing mappings by adding explicit destination hosts while preserving their placeholders and other settings.

### What's new

#### Security

- `sbx cp` blocks symbolic links in mounted directories from escaping onto other host paths. On Windows, it also blocks workspace escapes through substituted drives.
- Unix clients verify that daemon sockets belong to the same account before sending requests.
- `sbx kit sign` no longer runs configured Git filter drivers while checking signed file contents. Its recovery guidance preserves committed bytes even when repository attributes request line-ending conversion.
- Files replaced by Sandboxes on Windows retain their existing access control lists.
- Network policy commands reject empty `--sandbox` values instead of applying changes globally.

#### Kits and agents

- `sbx kit add` installs added image content in all-v3 sandboxes while preserving filesystem changes. `sbx run --model` generates a v3 model kit when the base uses v3.
- `sbx kit inspect --resolved` shows the kit with its `extends` parent merged in.
- Kit installation logs are saved inside the sandbox at `/var/log/sbx-kits/<kit>/install-<n>.log`. Progress shows command descriptions, and debug mode streams successful installation output and the startup log. Failed startup commands produce a warning with the log tail instead of failing `sbx run`.
- Kits can request publication of a host-loopback port when granted the required permission. Each request requires user approval and expires after 45 seconds.
- Kits declaring `agent-context@2` guidance work with local and cloud sandboxes, including when they also offer `agent-context@1`.
- Optional kit credentials no longer prevent creation when absent. Kit changes after a daemon restart preserve stored Bedrock profiles, and recovery after a failed kit change preserves access to shared skills.
- Shell kits support custom templates, and v3 kits load correctly when their descriptor is named `spec.yaml`.
- Docker Hardened Images (DHI) kits are accepted by the default kit source allow list.
- GitHub Copilot authentication works when the GitHub secret is stored with `--command` or as a reference.

#### Sandbox lifecycle and workspaces

- Additional workspaces accept `PATH:TARGET[:ro]` with `sbx create` and `sbx run`. Environment files support the same custom targets through `additionalWorkspaces[].target`.
- `sbx run --workdir` starts the agent in a directory inside the workspace, including when reattaching to an existing sandbox.
- Environment files support `entrypoint` and `cmd` defaults for `sbx env run`. Command-line arguments can replace the default command. `sandboxOptions` also accepts `rootSize`, `dockerSize`, and `clonedWorkspaceSize`.
- Interactive `sbx env exec` offers to create a missing environment sandbox. Interactive `sbx env run` can adopt a matching existing sandbox when its environment records have been lost, after confirmation.
- New local sandboxes support FUSE filesystems, controlled through `sandbox.fuse` and the environment option `sandboxOptions.fuse`.
- Sandbox memory limits can be set up to the host's total memory, with a warning above 75%.
- Workspace, clone, kit-volume, and mount targets that would hide the shared skills store are rejected before creation or modification. Startup also explains how to recover when a file occupies the shared skills directory.
- `sbx inspect` shows the Sandboxes version that created a sandbox, for sandboxes created with this release onward. `sbx ls`, `sbx inspect`, and `sbx prune --dry-run` JSON include configured host bind mounts, including workspace, kit, and skills mounts. List and prune output also indicate missing mount sources. Existing workspace fields remain available.
- `sbx ls` reports unreadable sandbox state instead of silently omitting the sandbox. Clone failures include Git's underlying reason, such as dubious ownership.
- Resizing an interactive terminal during startup or shutdown no longer produces spurious errors. On Windows, Ctrl+Z no longer ends input forwarding, and automatically started daemons are less likely to stop when the launching command exits.

#### Networking and policy

- Linux hosts using split-horizon DNS, including VPN configurations with corporate nameservers, can reach internal HTTPS hosts through the sandbox proxy and retrieve proxy auto-configuration files from them.
- The proxy reuses upstream connections when injecting credentials, improving parallel downloads such as package installations.
- Bursts of outbound connections no longer crash the daemon when the host network stalls.
- Failed upstream connections for inspected HTTPS requests return a 502 response identifying the host and reason instead of an empty reply.
- `sbx policy check http` checks HTTP governance policy behavior.
- Policy listings keep the implicit default-deny fallback visible and distinguish it from explicit deny rules. Policy setup correctly identifies the local preset, and Gordon preserves HTTP rule restrictions.
- The policy TUI shows setup guidance after a reset, warns about stale rules, and resumes updates after setup. Storage failures report a server error instead of incorrectly requesting policy initialization.
- `sbx reset` reports how many user-added and approval-granted network rules it will remove, including with `--force`, and points to `sbx policy init` afterward.
- `sbx ports` accepts `-p` for `--publish` and identifies the destination sandbox port when a host port is already published.

#### Authentication, MCP, and models

- `sbx secret set-custom --value-stdin` reads a secret from standard input, avoiding secret values in command arguments.
- Daemon image pulls use the `sbx login` session and stored registry credentials. Published agent kits can be fetched through authenticated registry mirrors.
- Enable `mcp.useServerNamePrefixes` to prefix MCP tool names with their server name.
- `sbx create` and `sbx run` accept `--static-mcp none` to use static MCP mode without attaching servers.
- Approval-gated MCP calls proceed after approval, including with newer MCP clients. Server tools ending in `-authorize` or `-revoke-auth` are available again.
- MCP servers have more time to initialize, reducing startup timeouts. Code-mode sessions cancel in-flight tool calls on exit and handle parallel requests more reliably.
- `sbx mcp auth` and `sbx mcp add` suggest `--client-id` when a server refuses dynamic client registration.
- Set `model.server.enabled` to `false`, then run `sbx daemon restart`, to disable the bundled model server and its inference-runtime download.

#### Cloud sandboxes

- Kit updates no longer fail with “Template has N running sandboxes” while sandboxes using an earlier version are running.
- Building source kits uses less temporary disk space during upload.
- Volume names tolerate surrounding whitespace and reject whitespace-only values. Cloud launches reject local-only model overflow flags instead of silently ignoring them.
- Port listings include the reported protocol in JSON and YAML output. `sbx move --ttl` accepts case-insensitive duration units.
- Help and shell completion identify supported cloud options and explain file-transfer limitations for symbolic links and empty directories.

#### CLI, settings, and diagnostics

- `--no-input` prevents interactive prompts, including Docker sign-in, and reports the command needed to continue.
- `sbx run --no-progress` suppresses creation progress while retaining command results. `sbx create --no-progress` still prompts for cloud network-policy confirmation when needed.
- `sbx --version` prints the client version without starting the daemon. Canceled sandbox and environment creation exits with status 130 and a clear message.
- Persistently disable local CLI and daemon telemetry with `sbx settings set telemetry.disabled true`. The next CLI invocation uses the preference; run `sbx daemon restart` to apply it to a running daemon.
- Compatibility errors and stale engine tasks identify `sbx daemon restart` as the recovery command. Daemon startup failures show only the current attempt's error output, preserving line breaks.
- Windows SSH configuration works with both native OpenSSH and Git Bash OpenSSH. `sbx diagnose` detects missing Windows shell support for Claude Desktop SSH connections and warns when the state directory uses eCryptfs, which can cause image-pull I/O errors.
- Command help and empty-list hints provide clearer next steps for secrets, MCP servers, skills, cloud sandboxes, and template recovery.

## 0.50.1

{{< release-date date="2026-10-28" >}}

[GitHub release](https://github.com/docker/sbx-releases/releases/tag/v0.50.1)

- fixes an issue that prevented sandboxes from starting

## 0.50.0

{{< release-date date="2026-10-26" >}}

[GitHub release](https://github.com/docker/sbx-releases/releases/tag/v0.50.0)

### Highlights

Docker Sandboxes v0.50.0 adds configurable kit host mounts, preserves persistent kit volumes across updates and sandbox recreation, and adds cloud sandbox inspection. It also strengthens host-execution, credential, and workspace-copy protections.

### What's new

#### Security

- `sbx cp` rejects symbolic links that escape a sandbox workspace when the workspace is reached through an alternate host path, including a symlink or junction, macOS `/private/tmp`, or a Linux bind mount of a workspace subdirectory.
- Host Git commands no longer run programs configured by a sandbox through the workspace repository's `.git` directory.
- Browser, clipboard, and self-update helpers run only from trusted system locations, preventing sandbox-writable directories on `PATH` from supplying those programs.
- Local source kits are reused and composed only when their digests match recorded local builds. Each local source kit rebuilds once after upgrading.
- Git kit references reject ambiguous paths, non-ASCII hostnames, and SSH transport syntax. Use IDNA ASCII labels for internationalized hostnames. HTTPS source restrictions are checked against the parsed URL host.
- MCP OAuth requests block cross-origin credential redirects. Rejected silent refreshes preserve saved credentials without opening a browser.
- Browser-open logs omit URL query strings and fragments, keeping OAuth parameters and tokens in those components out of daemon logs and diagnostic bundles.
- Windows clients verify that the daemon's named pipes belong to the same Windows user. Before using the updated CLI, older elevated daemons may need to be stopped with the previous CLI or Task Manager.
- Environment MCP commands require approval to execute on the host, remain scoped to their sandbox, and are disabled by `--skip-host-commands`. Approval notices explain that executable paths are pinned, but file contents are not. Keep executables and scripts outside writable sandbox mounts.
- Gordon clears **Approve all** and **Always allow** grants when a conversation resumes, requiring approval before running tools again.

#### Kits

- When creating a sandbox, use `--mount HOST:SANDBOX_PATH` with `sbx create` or `sbx run` to select a host directory for a kit-declared host mount. The sandbox target must exactly match the `path` of a `host-mount@1` capability declared by a kit in the resolved kit set. Environment files support the same bindings through `mounts`, mapping sandbox target paths to host directories. These options cannot add undeclared mount targets.
- Persistent kit volumes retain data across kit updates and sandbox removal and recreation with the same name. This applies to sandboxes created with this version; existing sandbox records retain their previous storage behavior.
- Creating a sandbox from a v3 kit set without a workload kit reports a clear explanation instead of "the kit set is not coherent."

#### Sandbox lifecycle and workspaces

- Host Git identity is copied into sandboxes for repositories that use the reftable ref format.
- Sandboxes created with `--clone` preserve kit-provided workspace files and report clone initialization failures.
- Canceling clone remote configuration preserves the cancellation error even if stopping Git also fails.

#### Authentication and agents

- Kiro IDE AWS IAM Identity Center sign-in callbacks reach the sandbox even when the callback port hasn't been published.
- Claude sandboxes using subscription OAuth are provisioned without the API-key helper.

#### Networking and policy

- HTTP network rules support query-parameter restrictions. Repeat `--query name=value` to require multiple parameters.

#### Cloud sandboxes

- `sbx --cloud inspect <id|name>` shows details for a single cloud sandbox, including published port URLs. Use `--format json` or `--format yaml` for scripting.
- Running `sbx --cloud run --name <name>` without an agent suggests the correct form: `sbx --cloud run <agent> --name <name>`.

#### CLI and diagnostics

- `sbx exec` returns exit status 125 when it cannot run the command, such as when the sandbox doesn't exist.
- On macOS and Windows, `sbx diagnose` and daemon logs warn when `SSL_CERT_FILE` or `SSL_CERT_DIR` replaces the system trust store and can cause HTTPS certificate verification failures.

#### Experimental Runtime backend

- Workloads have network access.
- Startup stops waiting and cleans up when a workload exits or fails before becoming ready.

<!-- END GENERATED RELEASES -->

## Earlier releases

For older versions, see the
[Docker Sandboxes releases on GitHub](https://github.com/docker/sbx-releases/releases).
