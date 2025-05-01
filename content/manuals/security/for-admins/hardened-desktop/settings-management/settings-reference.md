---
description: Reference for all settings and features that are configured with Settings Management
keywords: admin, controls, settings management, reference
title: Settings reference
linkTitle: Settings reference
---

This reference lists all Docker Desktop settings, including where they live,
which operating systems they apply to, and whether they're configurable via the
Docker Admin Console or the `admin-settings.json` file.

Each setting includes:

- Name
- Description
- OS compatibility
- The default value when a user first downloads Docker Desktop
- Accepted values
- Format of accepted values
- Use cases
- Details of how to configure the setting, either with [Docker Desktop](/manuals/desktop/settings-and-maintenance/settings.md) or
Settings Management (either Admin Console or `admin-settings.json` file)

For details on the format and usage of the `admin-settings.json` file, see
[Configure Settings Management with a JSON file](/manuals/security/for-admins/hardened-desktop/settings-management/configure-json-file.md).

## `AcceptCanaryUpdates`

- **Description:** Opt in to early access of Docker Desktop updates.
- **OS compatibility:** All
- **Default value:** `false`
- **Accepted values:** `true`, `false`
- **Format:** Boolean
- **Use case:** Enable early access to test new releases before general
availability.
- **Configure this setting with:**
    - `settings-store.json` or `settings.json` files

## `ActiveOrganizationName`

- **Description:** Stores the active organization name for Docker Business
accounts
- **OS compatibility:** All
- **Default value:** `""`
- **Accepted values:** String
- **Format:** String
- **Use case:** Manage organization-specific Docker settings.
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## `AllowBetaFeatures`
- **Description:** Allow access to Beta features in Docker Desktop.
- **OS compatibility:** All
- **Default value:** `true`
- **Accepted values:** `true`, `false`
- **Format:** Boolean
- **Use case:** Enable early features for testing upcoming functionality.
- **Configure this setting with:**
    - Features in development settings in [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files
    - Settings Management: `allowBetaFeatures` setting in the [`admin-settings.json` file](/manuals/security/for-admins/hardened-desktop/settings-management/configure-json-file.md)
    - Settings Management: **Access beta features** setting in the [Admin Console](/manuals/security/for-admins/hardened-desktop/settings-management/configure-admin-console.md)

## `AllowExperimentalFeatures`

- **Description:** Allow access to Experimental features in Docker Desktop.
- **OS compatibility:** All
- **Default value:** `true`
- **Accepted values:** `true`, `false`
- **Format:** Boolean
- **Use case:** Enable experimental features.
- **Configure this setting with:**
    - Features in development settings in [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files
    - Settings Management: `allowExperimentalFeatures` setting in the [`admin-settings.json` file](/manuals/security/for-admins/hardened-desktop/settings-management/configure-json-file.md)
    - Settings Management: **Access experimental features** setting in the [Admin Console](/manuals/security/for-admins/hardened-desktop/settings-management/configure-admin-console.md)

## `AnalyticsEnabled`

- **Description:** Send usage statistics and crash reports to Docker. If set to
`false`, Docker Desktop doesn't send usage statistics to Docker.
- **OS compatibility:** All
- **Default value:** `true`
- **Accepted values:** `true`, `false`
- **Format:** Boolean
- **Use case:** Enable analytics to help Docker improve the product based on
usage data.
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files
    - Settings Management: `analyticsEnabled` setting in the [`admin-settings.json` file](/manuals/security/for-admins/hardened-desktop/settings-management/configure-json-file.md)
    - Settings Management: **Send usage statistics** setting in the [Admin Console](/manuals/security/for-admins/hardened-desktop/settings-management/configure-admin-console.md)

## `AutoDownloadUpdates`

- **Description:** Automatically download Docker Desktop updates when available.
- **OS compatibility:** All
- **Default value:** `false`
- **Accepted values:** `true`, `false`
- **Format:** Boolean
- **Use case:** Manage auto update behavior.
- **Configure this setting with:**
    - Software updates settings in [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files
    - Settings Management: [Admin Console](/manuals/security/for-admins/hardened-desktop/settings-management/configure-admin-console.md)

## `AutoPauseTimedActivitySeconds`

- **Description:** Number of seconds before Docker Desktop auto-pauses due to
inactivity.
- **OS compatibility:** All
- **Default value:** `30`
- **Accepted values:** Integer (seconds)
- **Format:** Integer
- **Use case:** Save system resources during periods of inactivity.
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## `AutoPauseTimeoutSeconds`

- **Description:** Maximum idle time allowed before Docker Desktop pauses.
- **OS compatibility:** All
- **Default value:** `300`
- **Accepted values:** Integer (seconds)
- **Format:** Integer
- **Use case:** Manage Docker Desktop pause behavior during long idle periods.
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## `AutoStart`

- **Description:** Start Docker Desktop automatically when booting machine.
- **OS compatibility:** All
- **Default value:** `false`
- **Accepted values:** `true`, `false`
- **Format:** Boolean
- **Use case:** Ensure Docker Desktop is always running after boot.
- **Configure this setting with:**
    - General settings in [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## `BackupData`

- **Description:** Enable or disable backup of Docker Desktop application data.
- **OS compatibility:** All
- **Default value:** `false`
- **Accepted values:** `true`, `false`
- **Format:** Boolean
- **Use case:** Manage persistence of application data.
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## `BlockDockerLoad`

- **Description:** Block the `docker load` command to prevent loading local images.
If the value is set to `true`, users are no longer able to run `docker load`
and receive an error if they try to.
- **OS compatibility:** All
- **Default value:** `false`
- **Accepted values:** `true`, `false`
- **Format:** Boolean
- **Use case:** Harden security by restricting local image loading.
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files
    - Settings Management: `blockDockerLoad` setting in the [`admin-settings.json` file](/manuals/security/for-admins/hardened-desktop/settings-management/configure-json-file.md)
    - Settings Management: **Block Docker load** setting in the [Admin Console](/manuals/security/for-admins/hardened-desktop/settings-management/configure-admin-console.md)

## `ContainerTerminal`

- **Description:** Select default terminal for launching Docker CLI from Docker
Desktop.
- **OS compatibility:** All
- **Default value:** `integrated`
- **Accepted values:** `integrated`, `system`
- **Format:** String
- **Use case:** Customize developer experience with preferred terminal.
- **Configure this setting with:**
    - General settings in [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## `ContainersOverrideProxyExclude`

- **Description:** Configure addresses that containers should bypass from proxy
settings.
- **OS compatibility:** All
- **Default value:** `""`
- **Accepted values:** List of addresses
- **Format:** String
- **Use case:** Fine-tune proxy exceptions for container networking.
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## `ContainersOverrideProxyHTTP`

- **Description:** HTTP proxy setting for container networking.
- **OS compatibility:** All
- **Default value:** `""`
- **Accepted values:** URL string
- **Format:** String
- **Use case:** Set up container traffic to use a custom HTTP proxy.
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## `ContainersOverrideProxyHTTPS`

- **Description:** HTTPS proxy setting for container networking.
- **OS compatibility:** All
- **Default value:** `""`
- **Accepted values:** URL string
- **Format:** String
- **Use case:** Set up container traffic to use a custom HTTPS proxy.
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## `ContainersOverrideProxyPAC`

- **Description:** PAC (Proxy Auto-config) URL for container networking.
- **OS compatibility:** Windows
- **Default value:** `""`
- **Accepted values:** URL string
- **Format:** String
- **Use case:** Automatically configure container proxy routing via PAC file.
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files
    - Settings Management: **PAC** setting in the [Admin Console](/manuals/security/for-admins/hardened-desktop/settings-management/configure-admin-console.md)

## `ContainersOverrideProxyTCP`

- **Description:** TCP proxy setting for container networking.
- **OS compatibility:** All
- **Default value:** `""`
- **Accepted values:** String
- **Format:** String
- **Use case:** Configure advanced TCP proxy for containers.
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## `ContainersOverrideProxyTransparentPorts`

- **Description:** List of ports to bypass transparent proxying in containers.
- **OS compatibility:** All
- **Default value:** `80,443`
- **Accepted values:** List of ports
- **Format:** String
- **Use case:** Exclude specific ports from transparent proxy behavior.
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files
    - Settings Management: **Transparent ports** setting in the [Admin Console](/manuals/security/for-admins/hardened-desktop/settings-management/configure-admin-console.md)

## `ContainersProxyHTTPMode`

- **Description:** Creates air-gapped containers. For more information, see
[Air-Gapped Containers](/manuals/security/for-admins/hardened-desktop/air-gapped-containers.md).
- **OS compatibility:** All
- **Default value:** `system`
- **Accepted values:** `manual`, `system`
- **Format:** String
- **Use case:** Fine-tune container proxy behavior.
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files
    - Settings Management: `containersProxy` setting in the [`admin-settings.json` file](/manuals/security/for-admins/hardened-desktop/settings-management/configure-json-file.md)
    - Settings Management: **Proxy mode** setting in the [Admin Console](/manuals/security/for-admins/hardened-desktop/settings-management/configure-admin-console.md)

## `Cpus`

- **Description:** Number of CPUs assigned to the Docker Desktop virtual machine.
- **OS compatibility:** All
- **Default value:** The number of logical CPU cores available on the host system.
- **Accepted values:** Integer
- **Format:** Integer
- **Use case:** Resource allocation control.
- **Configure this setting with:**
    - Resources settings in [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## `CredentialHelper`

- **Description:** Credential storage helper to use for `docker login`.
- **OS compatibility:** macOS
- **Default value:** `desktop`
- **Accepted values:** String
- **Format:** String
- **Use case:** Manage secure storage of Docker credentials.
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## `CustomWslDistroDir`

- **Description:** Custom path for WSL2 distributions managed by Docker.
- **OS compatibility:** Windows + WSL
- **Default value:** `%USERPROFILE%\AppData\Local\Docker\wsl\distro`
- **Accepted values:** File path
- **Format:** String
- **Use case:** Control where Docker stores WSL2 distributions.
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## `DataFolder`

- **Description:** Path where Docker Desktop stores virtual machine data.
- **OS compatibility:** All
- **Default value:**
    - macOS: `~/Library/Containers/com.docker.docker/Data/vms/0`
    - Windows: `%USERPROFILE%\AppData\Local\Docker\wsl\data`
- **Accepted values:** File path
- **Format:** String
- **Use case:** Redirect Docker data to a custom location.
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## `DefaultSnapshotter`

- **Description:** Set the default container snapshotter.
- **OS compatibility:** All
- **Default value:** `overlayfs`
- **Accepted values:** String
- **Format:** String
- **Use case:** Control storage backend for container layers.
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## `DeprecatedCgroupv1`

- **Description:** Enable cgroup v1 support if needed for compatibility.
- **OS compatibility:** All
- **Default value:** `false`
- **Accepted values:** `true`, `false`
- **Format:** Boolean
- **Use case:** Maintain compatibility with legacy software.
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## `DesktopTerminalEnabled`

- **Description:** Enable access to the Docker Desktop integrated terminal. If
the value is set to `false`, users can't use the Docker terminal to interact
with the host machine and execute commands directly from Docker Desktop.
- **OS compatibility:** All
- **Default value:** `false`
- **Accepted values:** `true`, `false`
- **Format:** Boolean
- **Use case:** Allow or restrict developer access to the built-in terminal.
- **Configure this setting with:**
    - General settings in [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files
    - Settings Management: `desktopTerminalEnabled` setting in the [`admin-settings.json` file](/manuals/security/for-admins/hardened-desktop/settings-management/configure-json-file.md)

## `DefaultNetworkingMode`

- **Description:** Set the default networking mode for containers.
- **OS compatibility:** All
- **Default value:** `ipv4only`
- **Accepted values:** `ipv4only`, `ipv6only`, `dual-stack`
- **Format:** Enum
- **Use case:** Specify a custom container network mode.
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## `DevEnvironmentsEnabled`

- **Description:** Enable the Docker Dev Environments feature.
- **OS compatibility:** All
- **Default value:** `false`
- **Accepted values:** `true`, `false`
- **Format:** Boolean
- **Use case:** Control access to experimental development workflows.
- **Configure this setting with:**
    - Features in development settings in [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## `DisableHardwareAcceleration`

- **Description:** Disable hardware (GPU) acceleration support.
- **OS compatibility:** Windows
- **Default value:** `false`
- **Accepted values:** `true`, `false`
- **Format:** Boolean
- **Use case:** Work around graphics driver issues or run in VMs.
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## `DisableUpdate`

- **Description:** Disable automatic update polling for Docker Desktop. If the
value is set to `true`, checking for updates and notifications about Docker
Desktop updates are disabled.
- **OS compatibility:** All
- **Default value:** `false`
- **Accepted values:** `true`, `false`
- **Format:** Boolean
- **Use case:** Freeze the current version in enterprise environments.
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files
    - Settings Management: `disableUpdate` setting in the [`admin-settings.json` file](/manuals/security/for-admins/hardened-desktop/settings-management/configure-json-file.md)
    - Settings Management: **Disable update** setting in the [Admin Console](/manuals/security/for-admins/hardened-desktop/settings-management/configure-admin-console.md)

## `DiskFlush`

- **Description:** Control when data flushing occurs for the VM disk.
- **OS compatibility:** All
- **Default value:** `os`
- **Accepted values:** String
- **Format:** String
- **Use case:** Tune disk performance versus safety.
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## `DiskSizeMiB`

- **Description:** Maximum disk size (in MiB) allocated for Docker Desktop.
- **OS compatibility:** All
- **Default value:** Default disk size of machine.
- **Accepted values:** Integer
- **Format:** Integer
- **Use case:** Constrain Docker's virtual disk size for storage management.
- **Configure this setting with:**
    - Resources settings in [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## `DiskStats`

- **Description:** Disk usage statistics.
- **OS compatibility:** macOS
- **Default value:** `""`
- **Accepted values:** String
- **Format:** String
- **Use case:** Monitor or debug disk usage performance on Unix-based systems.
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## `DiskTRIM`

- **Description:** Enable TRIM operation support to reclaim unused disk space.
- **OS compatibility:** All
- **Default value:** `true`
- **Accepted values:** `true`, `false`
- **Format:** Boolean
- **Use case:** Optimize disk usage over time.
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## `DisplayRestartDialog`

- **Description:** Show a restart notification when settings changes require a
restart.
- **OS compatibility:** All
- **Default value:** `true`
- **Accepted values:** `true`, `false`
- **Format:** Boolean
- **Use case:** Provide user feedback about restart requirements.
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## `DisplaySwitchWinLinContainers`

- **Description:** Allow users to switch between Linux and Windows containers.
- **OS compatibility:** Windows
- **Default value:** `false`
- **Accepted values:** `true`, `false`
- **Format:** Boolean
- **Use case:** Flexibility in development environments.
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## `Displayed18362Deprecation`

- **Description:** Show the deprecation warning for Windows build 18362.
- **OS compatibility:** Windows
- **Default value:** `false`
- **Accepted values:** `true`, `false`
- **Format:** Boolean
- **Use case:** Prevent showing the same Windows version deprecation warning
multiple times.
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## `DisplayedElectronPopup`

- **Description:** Show Electron (tips, alerts, announcements) pop-ups for users.
- **OS compatibility:** All
- **Default value:** `[]`
- **Accepted values:** List of strings
- **Format:** Array with list of strings
- **Use case:** Prevents Docker Desktop from repeatedly showing the same popup
messages.
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## `DisplayedOnboarding`

- **Description:** Display the onboarding survey for Docker Desktop. If the
value is set to `true`, the onboarding survey will not be displayed to new
users. Settings the value to `false` has no effecct.
- **OS compatibility:** All
- **Default value:** `false`
- **Accepted values:** `true`, `false`
- **Format:** Boolean
- **Use case:** Prevents Docker Desktop from repeatedly showing onboarding.
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files
    - Settings Management: `displayedOnboarding` setting in the [`admin-settings.json` file](/manuals/security/for-admins/hardened-desktop/settings-management/configure-json-file.md)
    - Settings Management: **Hide onboarding survey** setting in the [Admin Console](/manuals/security/for-admins/hardened-desktop/settings-management/configure-admin-console.md)

## `DockerAppLaunchPath`

- **Description:** Path to the Docker Desktop application executable on macOS.
- **OS compatibility:** macOS
- **Default value:** `/Applications/Docker.app`
- **Accepted values:** File path
- **Format:** String
- **Use case:** Custom install management or scripting.
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## `DockerBinInstallPath`

- **Description:** Install location for Docker CLI binaries.
- **OS compatibility:** All
- **Default value:** `system`
- **Accepted values:** File path
- **Format:** String
- **Use case:** Customize CLI install location for compliance or tooling.
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## `DockerDebugDefaultEnabled`

- **Description:** Enable debug logging by default for Docker CLI commands.
- **OS compatibility:** All
- **Default value:** `false`
- **Accepted values:** `true`, `false`
- **Format:** Boolean
- **Use case:** Assist with debugging support issues.
- **Configure this setting with:**
    - General settings in [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## `ECIDockerSocketAllowDerivedImages`

- **Description:** Allow Enhanced Container Isolation (ECI) to use derived images.
- **OS compatibility:** All
- **Default value:** `true`
- **Accepted values:** `true`, `false`
- **Format:** Boolean
- **Use case:** Permit use of base images with layered builds in ECI mode.
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files
    - Settings Management: **Allow derived images** setting in the [Admin Console](/manuals/security/for-admins/hardened-desktop/settings-management/configure-admin-console.md)

## `ECIDockerSocketCmdList`

- **Description:** Restricts the commands that containers can issue via the
bind-mounted Docker Engine socket.
- **OS compatibility:** All
- **Default value:** `[]`
- **Accepted values:** List of strings
- **Format:** Array with list of strings
- **Use case:** Fine-tune developer CLI access in hardened environments.
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files
    - Settings Management: `commandList` setting in the [`admin-settings.json` file](/manuals/security/for-admins/hardened-desktop/settings-management/configure-json-file.md)
    - Settings Management: **Command list** setting in the [Admin Console](/manuals/security/for-admins/hardened-desktop/settings-management/configure-admin-console.md)

## `ECIDockerSocketCmdListType`

- **Description:** Whether the ECI command list is an allow-list or deny-list.
- **OS compatibility:** All
- **Default value:** `deny`
- **Accepted values:** `allow`, `deny`
- **Format:** String
- **Use case:** Determine behavior of `ECIDockerSocketCmdList`.
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## `ECIDockerSocketImgList`

- **Description:** 	Indicates which container images are allowed to bind-mount
the Docker Engine socket.
- **OS compatibility:** All
- **Default value:** `[]`
- **Accepted values:** List of strings
- **Format:** Array list of strings
- **Use case:** Restrict containers to a known set of images.
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files
    - Settings Management: `imageList` setting in the [`admin-settings.json` file](/manuals/security/for-admins/hardened-desktop/settings-management/configure-json-file.md)
    - Settings Manageament: **Image list** setting in the [Admin Console](/manuals/security/for-admins/hardened-desktop/settings-management/configure-admin-console.md)

## `EnableDefaultDockerSocket`

- **Description:** By default, enhanced container isolation blocks bind-mounting
the Docker Engine socket into containers
(e.g., `docker run -v /var/run/docker.sock:/var/run/docker.sock ...`). This lets
you relax this in a controlled way. See ECI Configuration for more info.
- **OS compatibility:** All
- **Default value:** `true`
- **Accepted values:** `true`, `false`
- **Format:** Boolean
- **Use case:** Allow containers to access the Docker socket for scenarios like
Docker-in-Docker or containerized CI agents.
- **Configure this setting with:**
    - Advanced settings in [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files
    - Settings Management: `dockerSocketMount` setting in the [`admin-settings.json` file](/manuals/security/for-admins/hardened-desktop/settings-management/configure-json-file.md)

## `EnableDockerAI`

- **Description:** Enable Docker AI features in the Docker Desktop experience.
- **OS compatibility:** All
- **Default value:** `false`
- **Accepted values:** `true`, `false`
- **Format:** Boolean
- **Use case:** Enable or disable AI features like "Ask Gordon".
- **Configure this setting with:**
    - Features in development settings in [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files
    - Settings Management: `enableDockerAI` setting in the [`admin-settings.json` file](/manuals/security/for-admins/hardened-desktop/settings-management/configure-json-file.md)

## `EnableIntegrationWithDefaultWslDistro`

- **Description:** Automatically integrate Docker with the default WSL
distribution.
- **OS compatibility:** Windows + WSL
- **Default value:** `false`
- **Accepted values:** `true`, `false`
- **Format:** Boolean
- **Use case:** Ensure Docker integrates with default WSL distro automatically.
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## `EnableIntegrityCheck`

- **Description:** Perform integrity checks on Docker Desktop binaries.
- **OS compatibility:** All
- **Default value:** `true`
- **Accepted values:** `true`, `false`
- **Format:** Boolean
- **Use case:** Enforce binary verification for security.
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## `EnableSegmentDebug`

- **Description:** Enable debug logging for Docker Desktop’s Segment analytics
events.
- **OS compatibility:** All
- **Default value:** `false`
- **Accepted values:** `true`, `false`
- **Format:** Boolean
- **Use case:** Troubleshoot or inspect analytics event delivery during
development or support sessions.
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## `EnableWasmShims`

- **Description:** Enable WebAssembly (Wasm) shims to run Wasm containers.
- **OS compatibility:** All
- **Default value:** `false`
- **Accepted values:** `true`, `false`
- **Format:** Boolean
- **Use case:** Run Wasm workloads in Docker Desktop.
- **Configure this setting with:**
    - Features in development settings in [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## `EnhancedContainerIsolation`

- **Description:** Enable Enhanced Container Isolation for secure container
execution.
- **OS compatibility:** All
- **Default value:** `false`
- **Accepted values:** `true`, `false`
- **Format:** Boolean
- **Use case:** Prevent containers from modifying configuration or sensitive
host areas.
- **Configure this setting with:**
    - General settings in [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files
    - Settings Management: `enhancedContainerIsolation` setting in the [`admin-settings.json` file](/manuals/security/for-admins/hardened-desktop/settings-management/configure-json-file.md)
    - Settings Management: **Enable enhanced container isolation** setting in the [Admin Console](/manuals/security/for-admins/hardened-desktop/settings-management/configure-admin-console.md)

## `ExposeDockerAPIOnTCP2375`

- **Description:** Expose the Docker API over TCP on a specified port. If value
is set to `true`, the Docker API is exposed on port 2375. This port is
unauthenticated and should only be enabled if protected by suitable firewall
rules.
- **OS compatibility:** Windows
- **Default value:** `false`
- **Accepted values:** `true`, `false`
- **Format:** Boolean
- **Use case:** Allow non-TLS API access for development/testing.
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files
    - Settings Management: `exposeDockerAPIOnTCP2375` setting in the [`admin-settings.json` file](/manuals/security/for-admins/hardened-desktop/settings-management/configure-json-file.md)
    - Settings Management: **Expose Docker API** setting in the [Admin Console](/manuals/security/for-admins/hardened-desktop/settings-management/configure-admin-console.md)

## `ExtensionsEnabled`

- **Description:** Enable or disable Docker Extensions.
- **OS compatibility:** All
- **Default value:** `true`
- **Accepted values:** `true`, `false`
- **Format:** Boolean
- **Use case:** Control access to the Extensions Marketplace and installed
extensions.
- **Configure this setting with:**
    - Extensions settings in [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files
    - Settings Management: `extensionsEnabled` setting in the [`admin-settings.json` file](/manuals/security/for-admins/hardened-desktop/settings-management/configure-json-file.md)
    - Settings Management: **Allow Extensions** setting in the [Admin Console](/manuals/security/for-admins/hardened-desktop/settings-management/configure-admin-console.md)

## `ExtensionsPrivateMarketplace`

- **Description:** Enable a private marketplace for Docker Extensions.
- **OS compatibility:** All
- **Default value:** `false`
- **Accepted values:** `true`, `false`
- **Format:** Boolean
- **Use case:** Restrict extension installation to curated extensions.
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## `ExtensionsPrivateMarketplaceAdminContactURL`

- **Description:** Set a contact URL for admins on the private extensions
marketplace page.
- **OS compatibility:** All
- **Default value:** `""`
- **Accepted values:** URL string
- **Format:** String
- **Use case:** Help users contact support if they can’t find an extension.
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## `FilesharingDirectories`

- **Description:** List of allowed directories shared between the host and
containers. When a path is added, its subdirectories are allowed.
- **OS compatibility:** All
- **Default value:** Varies by OS (typically includes user and temp directories)
- **Accepted values:** List of file paths as strings. This setting also accepts
`$HOME`, `$TMP`, or `$TEMP` as path variables.
- **Format:** Array list of strings
- **Use case:** Restrict or define what file paths are available to containers.
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files
    - Settings Management: `filesharingAllowedDirectories` setting in the [`admin-settings.json` file](/manuals/security/for-admins/hardened-desktop/settings-management/configure-json-file.md)
    - Settings Management: **Allowed file sharing directories** setting in the [Admin Console](/manuals/security/for-admins/hardened-desktop/settings-management/configure-admin-console.md)

## `HostNetworkingEnabled`

- **Description:** Enable experimental host networking support.
- **OS compatibility:** All
- **Default value:** `false`
- **Accepted values:** `true`, `false`
- **Format:** Boolean
- **Use case:** Allow containers to use the host network stack.
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## `KubernetesEnabled`

- **Description:** Enable the integrated Kubernetes cluster in Docker Desktop.
- **OS compatibility:** All
- **Default value:** `false`
- **Accepted values:** `true`, `false`
- **Format:** Boolean
- **Use case:** Enable or disable Kubernetes support for developers.
- **Configure this setting with:**
    - Kubernetes settings in [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files
    - Settings Management: `kubernetes` setting in the [`admin-settings.json` file](/manuals/security/for-admins/hardened-desktop/settings-management/configure-json-file.md)
    - Settings Management: **Allow Kubernetes** setting in the [Admin Console](/manuals/security/for-admins/hardened-desktop/settings-management/configure-admin-console.md)

## `KubernetesImagesRepository`

- **Description:** Set a custom repository for Kubernetes images.
- **OS compatibility:** All
- **Default value:** `""`
- **Accepted values:** URL string
- **Format:** String
- **Use case:** Support Kubernetes use in restricted or offline environments.
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files
    - Settings Management: `imagesRepository` setting in the [`admin-settings.json` file](/manuals/security/for-admins/hardened-desktop/settings-management/configure-json-file.md)
    - Settings Management: **Kubernetes images repository** setting in the [Admin Console](/manuals/security/for-admins/hardened-desktop/settings-management/configure-admin-console.md)

## `KubernetesMode`

- **Description:** Set the Kubernetes node mode (single-node or multi-node).
- **OS compatibility:** All
- **Default value:** `kubeadm`
- **Accepted values:** String
- **Format:** String
- **Use case:** Control the topology of the integrated Kubernetes cluster.
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## `KubernetesNodesCount`

- **Description:** Number of nodes to create in a multi-node Kubernetes cluster.
- **OS compatibility:** All
- **Default value:** `1`
- **Accepted values:** Integer
- **Format:** Integer
- **Use case:** Scale the number of Kubernetes nodes for development or testing.
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## `KubernetesNodesVersion`

- **Description:** Version of Kubernetes used for cluster node creation.
- **OS compatibility:** All
- **Default value:** `1.31.1`
- **Accepted values:** Semantic version (e.g., `1.29.1`)
- **Format:** String
- **Use case:** Pin a specific Kubernetes version for consistency or
compatibility.
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## `LastLoginDate`

- **Description:** Timestamp of last successful Docker Desktop login.
- **OS compatibility:** All
- **Default value:** `0`
- **Accepted values:** `int64` values
- **Format:** Integer in `int64` format
- **Use case:** Display usage activity.
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## `LatestBannerKey`

- **Description:** Tracks the most recently shown in-app banner.
- **OS compatibility:** All
- **Default value:** `""`
- **Accepted values:** String
- **Format:** String
- **Use case:** Prevent repeated display of the same banner across sessions.
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## `LicenseTermsVersion`

- **Description:** Version of Docker Desktop license terms accepted by the user.
- **OS compatibility:** All
- **Default value:** `0`
- **Accepted values:** Integer
- **Format:** Integer
- **Use case:** Audit license terms agreement.
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## `LifecycleTimeoutSeconds`

- **Description:** Number of seconds Docker Desktop waits for the Docker Engine
to start before timing out.
- **OS compatibility**: All
- **Default value:** `600`
- **Accepted values:** Integer (seconds)
- **Format:** Integer
- **Use case:** Extend or reduce the timeout window for environments where the
engine may start slowly.
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## `MemoryMiB`

- **Description:** Amount of RAM (in MiB) assigned to the Docker virtual machine.
- **OS compatibility:** All
- **Default value:** Based on system resources
- **Accepted values:** Integer
- **Format:** Integer
- **Use case:** Control how much memory Docker can use on the host.
- **Configure this setting with:**
    - Resources settings in [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## `OnlyMarketplaceExtensions`

- **Description:** Restrict Docker Desktop to only run Marketplace extensions.
- **OS compatibility:** All
- **Default value:** `true`
- **Accepted values:** `true`, `false`
- **Format:** Boolean
- **Use case:** Prevent running third-party or local extensions.
- **Configure this setting with:**
    - Extensions settings in [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## `OpenUIOnStartupDisabled`

- **Description:** Prevent the Docker Desktop UI from opening automatically at
startup.
- **OS compatibility:** All
- **Default value:** `false`
- **Accepted values:** `true`, `false`
- **Format:** Boolean
- **Use case:** Streamline startup experience.
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## `OverrideProxyExclude`

- **Description:** Comma-separated list of domain patterns that should bypass
the proxy.
- **OS compatibility:** All
- **Default value:** `""`
- **Accepted values:** String
- **Format:** String
- **Use case:** Exclude internal services or domains from being routed through
the proxy.
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## `OverrideProxyHTTP`

- **Description:** Override the default HTTP proxy used by Docker Desktop and
its containers.
- **OS compatibility**: All
- **Default value:** `""`
- **Accepted values:** URL string
- **Format:** String
- **Use case:** Route container HTTP traffic through a specific proxy.
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## `OverrideProxyHTTPS`

- **Description:** Override the default HTTPS proxy used by Docker Desktop and
its containers.
- **OS compatibility:** All
- **Default value:** `""`
- **Accepted values:** URL string
- **Format:** String
- **Use case:** Route container HTTPS traffic through a specific proxy.
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## `OverrideProxyPAC`

- **Description:** URL to a Proxy Auto-Config (PAC) file to dynamically
configure proxy rules.
- **OS compatibility:** All
- **Default value:** `""`
- **Accepted values:** URL string
- **Format:** String
- **Use case:** Load dynamic proxy rules from a PAC file.
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## `OverrideProxyTCP`

- **Description:** Override the TCP proxy settings used by Docker Desktop.
- **OS compatibility:** All
- **Default value:** `""`
- **Accepted values:** URL string
- **Format:** String
- **Use case:** Define a custom proxy for TCP traffic not covered by
HTTP/HTTPS proxies.
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## `OverrideWindowsDockerdPort`

- **Description:** Exposes Docker Desktop's internal proxy locally on this port
for the Windows Docker daemon to connect to. If set to `0`, a random free port
is chosen. If the value is greater than 0, it uses that exact value for the port.
-1 disables the option.
- **OS:** Windows
- **Description:** Override the port used by the Windows Docker deamon.
- **Default value:** `-1`
- **Use case:** Allow precise control of how Docker Desktop exposes its
internal proxy for `dockerd.exe`.
- **Configure this settings with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files
    - Settings Management: `windowsDockerdPort` setting in the [`admin-settings.json` file](/manuals/security/for-admins/hardened-desktop/settings-management/configure-json-file.md)
    - Settings Management: **Override Windows “dockerd” port** setting in the [Admin Console](/manuals/security/for-admins/hardened-desktop/settings-management/configure-admin-console.md)

## `ProxyEnableKerberosNTLM`

- **Description:** Enable Kerberos and NTLM authentication for the proxy.
- **OS:** Windows
- **Default value:** `false`
- **Accepted values:** `true`, `false`
- **Format:** Boolean
- **Use case:** Allow Docker Desktop to authenticate with enterprise proxies
that require Kerberos or NTLM credentials.
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files
    - Settings Management: `enableKerberosNtlm` setting in the [`admin-settings.json` file](/manuals/security/for-admins/hardened-desktop/settings-management/configure-json-file.md)
    - Settings Management: **Kerberos NTLM** setting in the [Admin Console](/manuals/security/for-admins/hardened-desktop/settings-management/configure-admin-console.md)

## `ProxyHTTPMode`

- **Description:** Proxy mode setting. If mode is set to `system` instead of
`manual`, Docker Desktop gets the proxy values from the system and ignores
values set for `http`, `https`, and `exclude`. To manually configure proxy
servers, use `manual`.
- **OS compatibility:** All
- **Default value:** `system`
- **Accepted values:** `system`, `manual`
- **Format:** String
- **Use case:** Control how Docker Desktop uses or ignores system proxy settings.
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files
    - Settings Management: `proxy` setting in the [`admin-settings.json` file](/manuals/security/for-admins/hardened-desktop/settings-management/configure-json-file.md)
    - Settings Management: [Admin Console](/manuals/security/for-admins/hardened-desktop/settings-management/configure-admin-console.md)

## `ProxyLocalhostPort`

- **Description:** Specifies the local port used by Docker Desktop’s internal
proxy to route container traffic through the host network.
- **OS compatibility:** All
- **Default value:** `0`
- **Accepted values:** Integer (port number)
- **Format:** Integer
- **Use case:** Bind the internal proxy to a fixed localhost port for debugging
or compatibility with network security tools.
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## `RequireVmnetd`

- **Description:** Require the privileged helper (`vmnetd`) for networking on
macOS.
- **OS compatibility:** macOS
- **Default value:** `true`
- **Accepted values:** `true`, `false`
- **Format:** Boolean
- **Use case:** Enforce elevated privileges for networking support
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## `RunWinServiceInWslMode`

- **Description:** Allow the Windows service that supports Docker Desktop to
run in WSL mode for enhanced integration.
- **OS compatibility:** Windows + WSL
- **Default value:** `false`
- **Accepted values:** `true`, `false`
- **Format:** Boolean
- **Use case:** Enable deeper integration between the Windows service layer and
the WSL-based Docker backend.
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## `SbomIndexing`

- **Description:** Enable SBOM indexing for container images
- **OS compatibility:** All
- **Default value:** `true`
- **Accepted values:** `true`, `false`
- **Format:** Boolean
- **Use case:** Control whether Docker indexes SBOMs for images
- **Configure this setting with:**
    - General settings in [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files
    - Settings Management: `sbomIndexing` setting in the [`admin-settings.json` file](/manuals/security/for-admins/hardened-desktop/settings-management/configure-json-file.md)
    - Settings Management: **SBOM indexing** settings in the [Admin Console](/manuals/security/for-admins/hardened-desktop/settings-management/configure-admin-console.md)

## `ScoutNotificationPopupsEnabled`

- **Description:** Enable Docker Scout popups inside Docker Desktop.
- **OS compatibility:** All
- **Default value:** `true`
- **Accepted values:** `true`, `false`
- **Format:** Boolean
- **Use case:** Show or hide vulnerability scan notifications
- **Configure this setting with:**
    - Notifications settings in [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## `ScoutOsNotificationsEnabled`

- **Description:** Enable Docker Scout notifications through the operating system.
- **OS compatibility:** All
- **Default value:** `false`
- **Accepted values:** `true`, `false`
- **Format:** Boolean
- **Use case:** Push Scout updates via system notification center
- **Configure this setting with:**
    - Notifications settings in [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## `SettingsVersion`

- **Description:** Specifies the version of the settings configuration file format
- **OS compatibility:** All
- **Default value:** `CurrentSettingsVersions`
- **Accepted values:** Integer
- **Format:** Integer
- **Use case:** Track schema versions for compatibility
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files
    - Settings Management: `configurationFileVersion` setting in the [`admin-settings.json` file](/manuals/security/for-admins/hardened-desktop/settings-management/configure-json-file.md)

## `ShowAnnouncementNotifications`

- **Description:** Display general announcements inside Docker Desktop.
- **OS compatibility:** All
- **Default value:** `true`
- **Accepted values:** `true`, `false`
- **Format:** Boolean
- **Use case:** Enable or suppress Docker-wide announcements in the UI.
- **Configure this setting with:**
    - Notifications settings in [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## `ShowExtensionsSystemContainers`

- **Description:** Show system containers used by Docker Extensions in the container list
- **OS compatibility:** All
- **Default value:** `false`
- **Accepted values:** `true`, `false`
- **Format:** Boolean
- **Use case:** Help developers troubleshoot or view extension system containers
- **Configure this setting with:**
    - Extensions settings in [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## `ShowGeneralNotifications`

- **Description:** Display general informational messages inside Docker Desktop
- **OS compatibility:** All
- **Default value:** `true`
- **Accepted values:** `true`, `false`
- **Format:** Boolean
- **Use case:** Customize in-app communication visibility
- **Configure this setting with:**
    - Notifications settings in [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## `ShowInstallScreen`

- **Description:** Show the installation onboarding screen in Docker Desktop
- **OS compatibility:** All
- **Default value:** `true`
- **Accepted values:** `true`, `false`
- **Format:** Boolean
- **Use case:** Control whether onboarding screens are shown after installation
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## `ShowKubernetesSystemContainers`

- **Description:** Show Kubernetes system containers in the Docker Dashboard container list
- **OS compatibility:** All
- **Default value:** `false`
- **Accepted values:** `true`, `false`
- **Format:** Boolean
- **Use case:** Allow developers to view kube-system containers for debugging
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## `ShowPromotionalNotifications`

- **Description:** Display promotional announcements and banners inside Docker Desktop
- **OS compatibility:** All
- **Default value:** `true`
- **Accepted values:** `true`, `false`
- **Format:** Boolean
- **Use case:** Control exposure to Docker news and feature promotion
- **Configure this setting with:**
    - Notifications settings in [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## `ShowSurveyNotifications`

- **Description:** Display notifications inviting users to participate in surveys
- **OS compatibility:** All
- **Default value:** `true`
- **Accepted values:** `true`, `false`
- **Format:** Boolean
- **Use case:** Enable or disable in-product survey prompts
- **Configure this setting with:**
    - Notifications settings in [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## `SkipUpdateToWSLPrompt`

- **Description:** Skip prompting users to upgrade to the WSL 2 backend
- **OS compatibility:** Windows + WSL
- **Default value:** `false`
- **Accepted values:** `true`, `false`
- **Format:** Boolean
- **Use case:** Silence UI nudges to switch WSL versions
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## `SkipWSLMountPerfWarning`

- **Description:** Skip the performance warning about WSL mount speed.
- **OS compatibility:** Windows + WSL
- **Default value:** `false`
- **Accepted values:** `true`, `false`
- **Format:** Boolean
- **Use case:** Suppress warnings for known limitations or user preference
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## `SocksProxyPort`

- **Description:** Local SOCKS proxy port for Docker Desktop.
- **OS compatibility:** All
- **Default value:** `0`
- **Accepted values:** Integer (port)
- **Format:** Integer
- **Use case:** Route Docker traffic through a SOCKS proxy
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## `SwapMiB`

- **Description:** Amount of swap space (in MiB) assigned to the Docker virtual machine
- **OS compatibility:** All
- **Default value:** `1024`
- **Accepted values:** Integer
- **Format:** Integer
- **Use case:** Extend memory availability via swap
- **Configure this setting with:**
    - Resources settings in [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## `SynchronizedDirectories`

- **Description:** Directories that should be synchronized between host and
container filesystems.
- **OS compatibility:** All
- **Default value:** Varies by system and user configurations
- **Accepted values:** Array of file paths as strings
- **Format:** Array
- **Use case:** Improve performance for bind mounts and volume sharing.
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## `ThemeSource`

- **Description:** Choose the Docker Desktop UI theme
- **OS compatibility:** All
- **Default value:** `system`
- **Accepted values:** `light`, `dark`, `system`
- **Format:** Enum
- **Use case:** Personalize Docker Desktop appearance
- **Configure this setting with:**
    - General settings in [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## `UpdateAvailableTime`

- **Description:** Timestamp of last update availability check
- **OS compatibility:** All
- **Default value:** `0`
- **Accepted values:** ISO 8601 timestamp
- **Format:** String
- **Use case:** Telemetry and internal logic
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## `UpdateHostsFile`

- **Description:** Allow Docker Desktop to update the system `hosts` file
- **OS compatibility:** All
- **Default value:** `false`
- **Accepted values:** `true`, `false`
- **Format:** Boolean
- **Use case:** Support DNS resolution for internal services
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## `UpdateInstallTime`

- **Description:** Timestamp of last Docker Desktop update installation.
- **OS compatibility:** All
- **Default value:** `0`
- **Accepted values:** ISO 8601 timestamp
- **Format:** String
- **Use case:** Track install history.
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## `UseBackgroundIndexing`

- **Description:** Enable background indexing of local Docker images for Docker
Scout.
- **OS compatibility:** All
- **Default value:** `false`
- **Accepted values:** `true`, `false`
- **Format:** Boolean
- **Use case:** Improve performance of features like image search.
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files
    - Settings Management: `useBackgroundIndexing` setting in the [`admin-settings.json` file](/manuals/security/for-admins/hardened-desktop/settings-management/configure-json-file.md)
    - Settings Management: **Background indexing** setting in the [Admin Console](/manuals/security/for-admins/hardened-desktop/settings-management/configure-admin-console.md)

## `UseContainerdSnapshotter`

- **Description:** Use containerd native snapshotter instead of legacy
snapshotters.
- **OS compatibility:** All
- **Default value:** `false`
- **Accepted values:** `true`, `false`
- **Format:** Boolean
- **Use case:** Improve image handling performance and compatibility.
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## `UseCredentialHelper`

- **Description:** Use the configured credential helper to securely store and
retrieve Docker registry credentials.
- **OS compatibility:** All
- **Default value:** `true`
- **Accepted values:** `true`, `false`
- **Format:** Boolean
- **Use case:** Enable secure, system-integrated storage of Docker login
credentials instead of plain-text config files.
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## `UseGrpcfuse`

- **Description:** Enable gRPC FUSE for macOS file sharing. If value is set to
`true`, gRPC Fuse is set as the file sharing mechanism.
- **OS compatibility:** macOS
- **Default value:** `true`
- **Accepted values:** `true`, `false`
- **Format:** Boolean
- **Use case:** Improve performance and compatibility of file mounts.
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files
    - Settings Management: `useGrpcfuse` setting in the [`admin-settings.json` file](/manuals/security/for-admins/hardened-desktop/settings-management/configure-json-file.md)
    - Settings Management: **Use gRPC FUSE for file sharing** setting in the [Admin Console](/manuals/security/for-admins/hardened-desktop/settings-management/configure-admin-console.md)

## `UseLibkrun`

- **Description:** Enable lightweight VM virtualization via libkrun.
- **OS compatibility:** macOS
- **Default value:** `false`
- **Accepted values:** `true`, `false`
- **Format:** Boolean
- **Use case:** Run containers in microVMs using libkrun.
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## `UseNightlyBuildUpdates`

- **Description:** Enable updates from the Docker Desktop nightly build channel
instead of the stable release channel.
- **OS compatibility:** All
- **Default value:** `false`
- **Accepted values:** `true`, `false`
- **Format:** Boolean
- **Use case:** Receive early access to experimental features and fixes by
subscribing to nightly builds.
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## `UseResourceSaver`

- **Description:** Enable Docker Desktop to pause when idle.
- **OS compatibility:** All
- **Default value:** `true`
- **Accepted values:** `true`, `false`
- **Format:** Boolean
- **Use case:** Save system resources during periods of inactivity.
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## `UseVirtualizationFramework`

- **Description:** Use Apple Virtualization Framework to run Docker containers.
- **OS compatibility:** macOS
- **Default value:** `true`
- **Accepted values:** `true`, `false`
- **Format:** Boolean
- **Use case:** Improve VM performance on Apple Silicon.
- **Configure this setting with:**
    - General settings in [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## `UseVirtualizationFrameworkRosetta`

- **Description:** Use Rosetta to emulate `amd64` on Apple Silicon. If value
is set to `true`, Docker Desktop turns on Rosetta to accelerate
x86_64/amd64 binary emulation on Apple Silicon.
- **OS compatibility:** macOS 13+
- **Default value:** `true`
- **Accepted values:** `true`, `false`
- **Format:** Boolean
- **Use case:** Run Intel-based containers on Apple Silicon hosts.
- **Configure this setting with:**
    - General settings in [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files
    - Settings Management:`useVirtualizationFrameworkRosetta` setting in the [`admin-settings.json` file](/manuals/security/for-admins/hardened-desktop/settings-management/configure-json-file.md)
    - Settings Management: **Use Rosetta for x86_64/amd64 emulation on Apple Silicon** setting in the [Admin Console](/manuals/security/for-admins/hardened-desktop/settings-management/configure-admin-console.md)

## `UseVirtualizationFrameworkVirtioFS`

- **Description:** Use VirtioFS for fast, native file sharing between host and
containers. If value is set to `true`, VirtioFS is set as the file sharing
mechanism. If both VirtioFS and gRPC are set to `true`, VirtioFS takes
precedence.
- **OS compatibility:** macOS 12.5+
- **Default value:** `true`
- **Accepted values:** `true`, `false`
- **Format:** Boolean
- **Use case:** Improve volume mount performance and compatibility.
- **Configure this setting with:**
    - General settings in [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files
    - Settings Management: `useVirtualizationFrameworkVirtioFS` setting in the [`admin-settings.json` file](/manuals/security/for-admins/hardened-desktop/settings-management/configure-json-file.md)
    - Settings Management: **Use VirtioFS for file sharing** setting in the [Admin Console](/manuals/security/for-admins/hardened-desktop/settings-management/configure-admin-console.md)

## `UseVpnkit`

- **Description:** Use vpnkit for Docker Desktop networking on macOS.
- **OS compatibility:** macOS
- **Default value:** `true`
- **Accepted values:** `true`, `false`
- **Format:** Boolean
- **Use case:** Enable or disable vpnkit as the networking backend.
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## `UseWindowsContainers`

- **Description:** Enable Windows container mode in Docker Desktop.
- **OS compatibility:** Windows
- **Default value:** `false`
- **Accepted values:** `true`, `false`
- **Format:** Boolean
- **Use case:** Switch between Linux and Windows container runtimes.
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files
    - Settings Management: `windowContainters` setting in the [`admin-settings.json` file](/manuals/security/for-admins/hardened-desktop/settings-management/configure-json-file.md)

## `VpnKitAllowedBindAddresses`

- **Description:** Specify which local IP addresses vpnkit is allowed to bind
to for handling network traffic.
- **OS compatibility:** All
- **Default value:** `0.0.0.0`
- **Accepted values:** IP address
- **Format:** String
- **Use case:** Restrict or allow vpnkit to bind to specific interfaces for
security or debugging purposes.
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## `VpnKitMTU`

- **Description:** Set the Maximum Transmission Unit (MTU) for vpnkit’s virtual
network interface.
- **OS compatibility:** All
- **Default value:** `1500`
- **Accepted values:** Integer
- **Format:** Integer
- **Use case:** Tune network performance or resolve issues with packet
fragmentation when using vpnkit.
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## `VpnKitMaxConnections`

- **Description:** Set the maximum number of simultaneous network connections
vpnkit can handle.
- **OS compatibility:** All
- **Default value:** `2000`
- **Accepted values:** Integer
- **Format:** Integer
- **Use case:** Control resource usage or support high-connection workloads
inside containers.
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## `VpnKitMaxPortIdleTime`

- **Description:** Maximum idle time in seconds before vpnkit closes an
unused port.
- **OS compatibility:** All
- **Default value:** `300`
- **Accepted values:** Integer (seconds)
- **Format:** Integer
- **Use case:** Improve performance and free up unused ports by closing
idle connections.
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## `VpnKitTransparentProxy`

- **Description:** Enable transparent proxying in vpnkit.
- **OS compatibility:** macOS
- **Default value:** `true`
- **Accepted values:** `true`, `false`
- **Format:** Boolean
- **Use case:** Seamlessly forward traffic through proxies using vpnkit.
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## `VpnkitCIDR`

- **Description:** Overrides the network range used for vpnkit DHCP/DNS for
`*.docker.internal`.
- **OS compatibility:** macOS
- **Default value:** `192.168.65.0/24`
- **Accepted values:** IP address
- **Format:** String
- **Use case:** Customize the subnet used for Docker container networking.
- **Configure this setting with:**
    - Resources settings in [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files
    - Settings Management: `vpnkitCIDR` setting in the [`admin-settings.json` file](/manuals/security/for-admins/hardened-desktop/settings-management/configure-json-file.md)
    - Settings Management: **VPN Kit CIDR** setting in the [Admin Console](/manuals/security/for-admins/hardened-desktop/settings-management/configure-admin-console.md)


## `WslDiskCompactionThresholdGb`

- **Description:** Minimum free disk space required to trigger WSL disk
compaction.
- **OS compatibility:** Windows + WSL
- **Default value:** `0`
- **Accepted values:** Integer (GB)
- **Format:** Integer
- **Use case:** Automatically reclaim unused space from WSL disks.
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## `WslEnableGrpcfuse`

- **Description:** Enable gRPC FUSE file sharing in WSL2 mode.
- **OS compatibility:** Windows + WSL
- **Default value:** `false`
- **Accepted values:** `true`, `false`
- **Format:** Boolean
- **Use case:** Improve performance and compatibility for file mounts in WSL.
- **Configure this setting with:**
    - General settings in [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## `WslEngineEnabled`

- **Description:** If the value is set to `true`, Docker Desktop uses the WSL2
based engine. This overrides anything that may have been set at installation
using the `--backend=<backend name>` flag.
- **OS compatibility:** Windows + WSL
- **Default value:** `true`
- **Accepted values:** `true`, `false`
- **Format:** Boolean
- **Use case:** Enable Linux containers via WSL 2 backend.
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files
    - Settings Management: `wslEngineEnabled` setting in the [`admin-settings.json` file](/manuals/security/for-admins/hardened-desktop/settings-management/configure-json-file.md)
    - Settings Management: **Windows Subsystem for Linux (WSL) Engine** setting in the [Admin Console](/manuals/security/for-admins/hardened-desktop/settings-management/configure-admin-console.md)

## `WslInstallMode`

- **Description:** Select how Docker Desktop installs and manages WSL on
Windows systems.
- **OS compatibility:** Windows + WSL
- **Default value:** `installLatestWsl`
- **Accepted values:** `installLatestWsl`, `manualInstall`
- **Format:** String
- **Use case:** Control whether Docker Desktop installs WSL automatically or
relies on a pre-installed version.
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## `WslUpdateRequired`

- **Description:** Indicates whether a WSL update is required for Docker Desktop
to function.
- **OS compatibility:** Windows + WSL
- **Default value:** `false`
- **Accepted values:** `true`, `false`
- **Format:** Boolean
- **Use case:** Internal check for platform support.
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files
