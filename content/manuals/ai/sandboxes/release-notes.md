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

## 0.54.1

{{< release-date date="2026-11-24" >}}

[GitHub release](https://github.com/docker/sbx-releases/releases/tag/v0.54.1)

### What's New

- On macOS, secrets saved before an upgrade can be removed without manually deleting their Keychain entries.

## 0.54.0

{{< release-date date="2026-11-16" >}}

[GitHub release](https://github.com/docker/sbx-releases/releases/tag/v0.54.0)

### Highlights

Docker Sandboxes v0.54.0 adds dedicated port management commands and a TLS compatibility setting for networks that block post-quantum key exchange. It also strengthens generated SSH configuration and clone-mode Git access, improves cloud sandbox creation, and makes CLI errors and output more consistent.

### Breaking changes and upgrade actions

#### Refresh generated SSH configuration

Generated SSH configuration prevents crafted host or user names from executing shell commands with older OpenSSH clients. After upgrading, rerun `sbx setup ssh` for local connections and `sbx --cloud setup ssh` for cloud connections to update existing configuration.

#### Storage paths when running as root

On Unix, running `sbx` as root uses root's registered home directory and the default storage, credential, and daemon endpoint locations. Inherited overrides, including `HOME`, `XDG_CONFIG_HOME`, `XDG_STATE_HOME`, `XDG_CACHE_HOME`, `SANDBOXES_STORAGE_ROOT`, and `DOCKER_SANDBOXES_API`, are ignored. Ordinary users retain their existing paths and environment overrides.

Before upgrading, stop any older root daemon that uses custom storage. The custom storage data is left untouched.

#### Clone-mode Git remote transport

Fetching from a clone-mode sandbox's `sandbox-<name>` Git remote now uses the daemon socket, restricted to your host user, instead of an unauthenticated localhost TCP port. Starting the sandbox refreshes its remote configuration. Fetching requires the sandbox to be running.

For Git clients that ignore the remote's `uploadpack` setting, enable TCP compatibility with `sbx settings set clone.legacyTcpRemote true`. This takes effect the next time the sandbox starts and exposes the repository and its history to other local accounts on the host. Leave the setting disabled to use the daemon socket. Clone-mode sandboxes using the default transport reject publishing sandbox TCP port 9418.

### What's new

#### Port management

Use dedicated subcommands to manage ports on local and cloud sandboxes:

- `sbx port ls <SANDBOX>` lists published ports.
- `sbx port publish <SANDBOX> <HOST_PORT>:<SANDBOX_PORT>` publishes a port. Omit the host port to assign one automatically.
- `sbx port unpublish <SANDBOX> <SANDBOX_PORT>` removes all host mappings for that sandbox port. Specify `<HOST_PORT>:<SANDBOX_PORT>` to remove a single mapping.

The legacy `sbx ports` command still works but is deprecated.

#### Networking and policy

- Set `tls.postQuantumKeyExchange` to `false` if a firewall or TLS-inspecting proxy drops or stalls Sandboxes HTTPS connections because it cannot handle post-quantum key exchange. Run `sbx settings set tls.postQuantumKeyExchange false`, then `sbx daemon restart`. The setting defaults to `true`; restore it and restart the daemon to re-enable post-quantum key exchange.
- `sbx policy check network` names the rule that allowed a target. Policy checks explain when further evaluation is needed and how to check a specific HTTP request.
- `sbx policy --help` explains hostname-blocking limitations when an allowed and denied host share an IP address.
- SSH reverse port forwarding no longer consumes a CPU core after its input socket is half-closed.

#### Sandbox lifecycle, kits, and agents

- `sbx run` no longer fails with “network already exists” after an earlier start failed partway.
- When a saved mount's source resolves to a different directory, startup refuses the mount. Remove it with `sbx umount`, then mount the intended directory again.
- Kit installation failures preserve command diagnostics and report network denials observed during installation.
- Image cleanup preserves assembled images referenced by saved kit locks or containers.
- Kit validation checks inherited configuration and reports inherited warnings.
- Existing kit permissions may require confirmation again after upgrading because permission checks use stronger service digests.
- Reattaching to a Kiro or Cursor sandbox relaunches the agent with the command used by its published kit.

#### Authentication and MCP

- `sbx login` reports failure with diagnostic steps when it cannot reach Docker's sign-in service after browser sign-in. Settings errors suggest signing in when administrator settings are unavailable.
- OAuth secret sign-in respects `--no-input` instead of opening browser authorization.
- Declining secret replacement or environment removal cancels the command with a nonzero exit status.
- `sbx inspect` shows a local sandbox's static or dynamic MCP mode and connected servers. These describe the sandbox gateway; use `sbx mcp ls` to list host registrations.
- Loading an MCP server into a cloud sandbox prints a success acknowledgment.

#### Cloud sandboxes

- Creating a sandbox honors the requested architecture when combining published kits with locally built mixins.
- A kit's network allow list no longer replaces the account's network policy when neither `--allow-network` nor `--deny-network` is supplied.
- Cloud creation rejects TTL values outside 15 seconds to 24 hours, explicitly empty timeout actions, and explicitly empty sandbox names before provisioning.
- Declining a sandbox move cancels it with a nonzero exit status. An explicitly empty destination name is rejected before work starts.
- Cloud API connection failures behind HTTP proxies report that the API is unreachable and suggest `sbx --cloud diagnose`, instead of returning a generic 500 error.
- File copies distinguish failed path checks from missing source paths.
- Cloud exec rejects `--detach-keys`, matching local exec. Omit this flag for both `sbx exec` and `sbx env exec`.
- Scheduled runs accept `--workdir` and `-w`. Schedule durations accept case-insensitive units.
- Help clarifies cloud template selection, inherited sizing, incompatible launch options, and the supported cloud SSH setup workflow.

#### CLI and diagnostics

- `sbx diagnose --format json` and `--format yaml` include a stable `id` for each check. Match checks by `id` in scripts; `name` is display text that can change.
- Diagnostic bundles use unique output files so existing workspace links cannot cause other host files to be overwritten.
- MCP progress, warnings, client-secret guidance, and secret-list diagnostics go to stderr, preserving stdout for command results.
- Sandbox and kit inspection report output write failures instead of succeeding silently.
- Cloud-only commands report the required mode without starting a local daemon.
- Installed skills can be removed with the `sbx skills delete` alias.
- Argument validation, shell completion, and secret-removal help more consistently reflect supported options. Kit signing rejects conflicting credential sources, and kit push rejects signing options unless signing is enabled.

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
- Unpublishing an automatically assigned port keeps it unpublished after a sandbox restart. (CVE-2026-126330)
- Published-port endpoints bracket IPv6 host addresses in list and inspect output, including JSON and YAML.
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

- Create and manage recurring cloud sandbox tasks with `sbx --cloud schedule`. Test tasks with manual runs, inspect run history, and pause or resume recurring runs.
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

<!-- END GENERATED RELEASES -->

## Earlier releases

For older versions, see the
[Docker Sandboxes releases on GitHub](https://github.com/docker/sbx-releases/releases).
