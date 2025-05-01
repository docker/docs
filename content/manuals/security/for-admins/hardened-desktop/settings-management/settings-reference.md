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

- Desktop setting name
- A values table that includes: the default value when a user first downloads
Docker Desktop, accepted values, and the format of accepted values
- Description
- OS compatibility
- Use cases
- How to configure the setting: Wwith [Docker Desktop](/manuals/desktop/settings-and-maintenance/settings.md) or
Settings Management (either Admin Console or `admin-settings.json` file)

For details on the format and usage of the `admin-settings.json` file, see
[Configure Settings Management with a JSON file](/manuals/security/for-admins/hardened-desktop/settings-management/configure-json-file.md).

## Accept canary updates

| Default value | Accepted values | Format |
|---------------|-----------------|--------|
| `false`       | `true`, `false` | Boolean |

- **Description:** Opt in to early access of Docker Desktop updates.
- **OS:** {{< badge color=blue text="All" >}}
- **Use case:** Enable early access to test new releases before general
availability.
- **Configure this setting with:**
    - `AcceptCanaryUpdates` in `settings-store.json` or `settings.json` files

## Active organization name

| Default value | Accepted values | Format |
|---------------|-----------------|--------|
| `""`          | String          | String |

- **Description:** Stores the active organization name for Docker Business
accounts
- **OS:** {{< badge color=blue text="All" >}}
- **Use case:** Manage organization-specific Docker settings.
- **Configure this setting with:**
    - `ActiveOrganizationName` in `settings-store.json` or `settings.json` files

## Allow beta features

| Default value | Accepted values | Format |
|---------------|-----------------|--------|
| `true`        | `true`, `false` | Boolean |

- **Description:** Allow access to Beta features in Docker Desktop.
- **OS:** {{< badge color=blue text="All" >}}
- **Use case:** Enable early features for testing upcoming functionality.
- **Configure this setting with:**
    - `AllowBetaFeatures` in `settings-store.json` or `settings.json` files
    - Settings Management: `allowBetaFeatures` setting in the [`admin-settings.json` file](/manuals/security/for-admins/hardened-desktop/settings-management/configure-json-file.md)
    - Settings Management: **Access beta features** setting in the [Admin Console](/manuals/security/for-admins/hardened-desktop/settings-management/configure-admin-console.md)

## Access experimental features

| Default value | Accepted values | Format |
|---------------|-----------------|--------|
| `true`        | `true`, `false` | Boolean |

- **Description:** Allow access to Experimental features in Docker Desktop.
- **OS:** {{< badge color=blue text="All" >}}
- **Use case:** Enable experimental features.
- **Configure this setting with:**
    - **Features in development** settings in [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `AllowExperimentalFeatures` in `settings-store.json` or `settings.json` files
    - Settings Management: `allowExperimentalFeatures` setting in the [`admin-settings.json` file](/manuals/security/for-admins/hardened-desktop/settings-management/configure-json-file.md)
    - Settings Management: **Access experimental features** setting in the [Admin Console](/manuals/security/for-admins/hardened-desktop/settings-management/configure-admin-console.md)

## Always download updates

| Default value | Accepted values | Format |
|---------------|-----------------|--------|
| `false`       | `true`, `false` | Boolean |

- **Description:** Automatically download Docker Desktop updates when available.
- **OS:** {{< badge color=blue text="All" >}}
- **Use case:** Manage auto update behavior.
- **Configure this setting with:**
    - **Software updates** settings in [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `AutoDownloadUpdates` in `settings-store.json` or `settings.json` files
    - Settings Management: **Disable updates** setting in the [Admin Console](/manuals/security/for-admins/hardened-desktop/settings-management/configure-admin-console.md)

## Auto pause activity

| Default value | Accepted values     | Format  |
|---------------|---------------------|---------|
| `30`          | Integer (seconds)   | Integer |

- **Description:** Number of seconds before Docker Desktop auto-pauses due to
inactivity.
- **OS:** {{< badge color=blue text="All" >}}
- **Use case:** Save system resources during periods of inactivity.
- **Configure this setting with:**
    - `AutoPauseTimedActivitySeconds` in `settings-store.json` or `settings.json` files

## Auto pause timeout

| Default value | Accepted values     | Format  |
|---------------|---------------------|---------|
| `300`         | Integer (seconds)   | Integer |

- **Description:** Maximum idle time allowed before Docker Desktop pauses.
- **OS:** {{< badge color=blue text="All" >}}
- **Use case:** Manage Docker Desktop pause behavior during long idle periods.
- **Configure this setting with:**
    - `AutoPauseTimeoutSeconds` in `settings-store.json` or `settings.json` files

## Block `docker load`

| Default value | Accepted values | Format |
|---------------|-----------------|--------|
| `false`       | `true`, `false` | Boolean |

- **Description:** Block the `docker load` command to prevent loading local images.
If the value is set to `true`, users are no longer able to run `docker load`
and receive an error if they try to.
- **OS:** {{< badge color=blue text="All" >}}
- **Use case:** Harden security by restricting local image loading.
- **Configure this setting with:**
    - `BlockDockerLoad` in `settings-store.json` or `settings.json` files
    - Settings Management: `blockDockerLoad` setting in the [`admin-settings.json` file](/manuals/security/for-admins/hardened-desktop/settings-management/configure-json-file.md)
    - Settings Management: **Block Docker load** setting in the [Admin Console](/manuals/security/for-admins/hardened-desktop/settings-management/configure-admin-console.md)

## Choose container terminal

| Default value | Accepted values         | Format |
|---------------|-------------------------|--------|
| `integrated`  | `integrated`, `system`  | String |

- **Description:** Select default terminal for launching Docker CLI from Docker
Desktop.
- **OS:** {{< badge color=blue text="All" >}}
- **Use case:** Customize developer experience with preferred terminal.
- **Configure this setting with:**
    - General settings in [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `ContainerTerminal` in `settings-store.json` or `settings.json` files

## Include VM in Time Machine backup

| Default value | Accepted values | Format |
|---------------|-----------------|--------|
| `false`       | `true`, `false` | Boolean |

- **Description:** Back up the Docker Desktop virtual machine.
- **OS:** {{< badge color=blue text="Mac only" >}}
- **Use case:** Manage persistence of application data.
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `BackupData` in `settings-store.json` or `settings.json` files

## Send usage statistics

| Default value | Accepted values | Format |
|---------------|-----------------|--------|
| `true`        | `true`, `false` | Boolean |

- **Description:** Send usage statistics and crash reports to Docker. If set to
`false`, Docker Desktop doesn't send usage statistics to Docker.
- **OS:** {{< badge color=blue text="All" >}}
- **Use case:** Enable analytics to help Docker improve the product based on
usage data.
- **Configure this setting with:**
    - Send usage statistics [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `AnalyticsEnabled` in `settings-store.json` or `settings.json` files
    - Settings Management: `analyticsEnabled` setting in the [`admin-settings.json` file](/manuals/security/for-admins/hardened-desktop/settings-management/configure-json-file.md)
    - Settings Management: **Send usage statistics** setting in the [Admin Console](/manuals/security/for-admins/hardened-desktop/settings-management/configure-admin-console.md)

## Start Docker Desktop when you sign in to your computer

| Default value | Accepted values | Format |
|---------------|-----------------|--------|
| `false`       | `true`, `false` | Boolean |

- **Description:** Start Docker Desktop automatically when booting machine.
- **OS:** {{< badge color=blue text="All" >}}
- **Use case:** Ensure Docker Desktop is always running after boot.
- **Configure this setting with:**
    - General settings in [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `AutoStart` in `settings-store.json` or `settings.json` files

## `ContainersOverrideProxyExclude`

| Default value | Accepted values    | Format |
|---------------|--------------------|--------|
| `""`          | List of addresses  | String |

- **Description:** Configure addresses that containers should bypass from proxy
settings.
- **OS:** {{< badge color=blue text="All" >}}
- **Use case:** Fine-tune proxy exceptions for container networking.
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## `ContainersOverrideProxyHTTP`

| Default value | Accepted values | Format |
|---------------|-----------------|--------|
| `""`          | URL string      | String |

- **Description:** HTTP proxy setting for container networking.
- **OS:** {{< badge color=blue text="All" >}}
- **Use case:** Set up container traffic to use a custom HTTP proxy.
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## `ContainersOverrideProxyHTTPS`

| Default value | Accepted values | Format |
|---------------|-----------------|--------|
| `""`          | URL string      | String |

- **Description:** HTTPS proxy setting for container networking.
- **OS:** {{< badge color=blue text="All" >}}
- **Use case:** Set up container traffic to use a custom HTTPS proxy.
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## `ContainersOverrideProxyPAC`

| Default value | Accepted values | Format |
|---------------|-----------------|--------|
| `""`          | URL string      | String |

- **Description:** PAC (Proxy Auto-config) URL for container networking.
- **OS:** {{< badge color=blue text="Windows only" >}}
- **Use case:** Automatically configure container proxy routing via PAC file.
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files
    - Settings Management: **PAC** setting in the [Admin Console](/manuals/security/for-admins/hardened-desktop/settings-management/configure-admin-console.md)

## `ContainersOverrideProxyTCP`

| Default value | Accepted values | Format |
|---------------|-----------------|--------|
| `""`          | String          | String |

- **Description:** TCP proxy setting for container networking.
- **OS:** {{< badge color=blue text="All" >}}
- **Use case:** Configure advanced TCP proxy for containers.
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## `ContainersOverrideProxyTransparentPorts`

| Default value | Accepted values | Format |
|---------------|-----------------|--------|
| `80,443`      | List of ports   | String |

- **Description:** List of ports to bypass transparent proxying in containers.
- **OS:** {{< badge color=blue text="All" >}}
- **Use case:** Exclude specific ports from transparent proxy behavior.
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files
    - Settings Management: **Transparent ports** setting in the [Admin Console](/manuals/security/for-admins/hardened-desktop/settings-management/configure-admin-console.md)

## `ContainersProxyHTTPMode`

| Default value | Accepted values     | Format |
|---------------|---------------------|--------|
| `system`      | `manual`, `system`  | String |

- **Description:** Creates air-gapped containers. For more information, see
[Air-Gapped Containers](/manuals/security/for-admins/hardened-desktop/air-gapped-containers.md).
- **OS:** {{< badge color=blue text="All" >}}
- **Use case:** Fine-tune container proxy behavior.
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files
    - Settings Management: `containersProxy` setting in the [`admin-settings.json` file](/manuals/security/for-admins/hardened-desktop/settings-management/configure-json-file.md)
    - Settings Management: **Proxy mode** setting in the [Admin Console](/manuals/security/for-admins/hardened-desktop/settings-management/configure-admin-console.md)

## `Cpus`

| Default value                                 | Accepted values | Format  |
|-----------------------------------------------|-----------------|---------|
| Number of logical CPU cores available on host | Integer         | Integer |

- **Description:** Number of CPUs assigned to the Docker Desktop virtual machine.
- **OS:** {{< badge color=blue text="All" >}}
- **Use case:** Resource allocation control.
- **Configure this setting with:**
    - Resources settings in [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## `CredentialHelper`

| Default value | Accepted values | Format |
|---------------|-----------------|--------|
| `desktop`     | String          | String |

- **Description:** Credential storage helper to use for `docker login`.
- **OS:** {{< badge color=blue text="Mac only" >}}
- **Use case:** Manage secure storage of Docker credentials.
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## `CustomWslDistroDir`

| Default value                                 | Accepted values | Format |
|----------------------------------------------|-----------------|--------|
| `%USERPROFILE%\AppData\Local\Docker\wsl\distro` | File path       | String |

- **Description:** Custom path for WSL2 distributions managed by Docker.
- **OS:** {{< badge color=blue text="Windows only" >}} + WSL
- **Use case:** Control where Docker stores WSL2 distributions.
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## `DataFolder`

| Default value                                                                 | Accepted values | Format |
|-------------------------------------------------------------------------------|-----------------|--------|
| macOS: `~/Library/Containers/com.docker.docker/Data/vms/0`  <br> Windows: `%USERPROFILE%\AppData\Local\Docker\wsl\data` | File path       | String |

- **Description:** Path where Docker Desktop stores virtual machine data.
- **OS:** {{< badge color=blue text="All" >}}
- **Use case:** Redirect Docker data to a custom location.
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## `DefaultSnapshotter`

| Default value | Accepted values | Format |
|---------------|-----------------|--------|
| `overlayfs`   | String          | String |

- **Description:** Set the default container snapshotter.
- **OS:** {{< badge color=blue text="All" >}}
- **Use case:** Control storage backend for container layers.
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## `DeprecatedCgroupv1`

| Default value | Accepted values | Format |
|---------------|-----------------|--------|
| `false`       | `true`, `false` | Boolean |

- **Description:** Enable cgroup v1 support if needed for compatibility.
- **OS:** {{< badge color=blue text="All" >}}
- **Use case:** Maintain compatibility with legacy software.
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## Enable Desktop terminal

| Default value | Accepted values | Format |
|---------------|-----------------|--------|
| `false`       | `true`, `false` | Boolean |

- **Description:** Enable access to the Docker Desktop integrated terminal. If
the value is set to `false`, users can't use the Docker terminal to interact
with the host machine and execute commands directly from Docker Desktop.
- **OS:** {{< badge color=blue text="All" >}}
- **Use case:** Allow or restrict developer access to the built-in terminal.
- **Configure this setting with:**
    - **Enable Docker terminal** setting in [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `DesktopTerminalEnabled` in `settings-store.json` or `settings.json` files
    - Settings Management: `desktopTerminalEnabled` setting in the [`admin-settings.json` file](/manuals/security/for-admins/hardened-desktop/settings-management/configure-json-file.md)

## Default networking mode

| Default value | Accepted values                    | Format |
|---------------|------------------------------------|--------|
| `ipv4only`    | `ipv4only`, `ipv6only`, `dual-stack` | Enum   |

- **Description:** Set the default networking mode for containers.
- **OS:** {{< badge color=blue text="All" >}}
- **Use case:** Specify a custom container network mode.
- **Configure this setting with:**
    - `DefaultNetworkingMode` in `settings-store.json` or `settings.json` files

## `DevEnvironmentsEnabled`

| Default value | Accepted values | Format |
|---------------|-----------------|--------|
| `false`       | `true`, `false` | Boolean |

- **Description:** Enable the Docker Dev Environments feature.
- **OS:** {{< badge color=blue text="All" >}}
- **Use case:** Control access to experimental development workflows.
- **Configure this setting with:**
    - Features in development settings in [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## `DisableHardwareAcceleration`

| Default value | Accepted values | Format |
|---------------|-----------------|--------|
| `false`       | `true`, `false` | Boolean |

- **Description:** Disable hardware (GPU) acceleration support.
- **OS:** {{< badge color=blue text="Windows only" >}}
- **Use case:** Work around graphics driver issues or run in VMs.
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## Disable automatic updates

| Default value | Accepted values | Format |
|---------------|-----------------|--------|
| `false`       | `true`, `false` | Boolean |

- **Description:** Disable automatic update polling for Docker Desktop. If the
value is set to `true`, checking for updates and notifications about Docker
Desktop updates are disabled.
- **OS:** {{< badge color=blue text="All" >}}
- **Use case:** Freeze the current version in enterprise environments.
- **Configure this setting with:**
    - `DisableUpdate` in `settings-store.json` or `settings.json` files
    - Settings Management: `disableUpdate` setting in the [`admin-settings.json` file](/manuals/security/for-admins/hardened-desktop/settings-management/configure-json-file.md)
    - Settings Management: **Disable update** setting in the [Admin Console](/manuals/security/for-admins/hardened-desktop/settings-management/configure-admin-console.md)

## `DiskFlush`

| Default value | Accepted values | Format |
|---------------|-----------------|--------|
| `os`          | String          | String |

- **Description:** Control when data flushing occurs for the VM disk.
- **OS:** {{< badge color=blue text="All" >}}
- **Use case:** Tune disk performance versus safety.
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## `DiskSizeMiB`

| Default value                  | Accepted values | Format  |
|-------------------------------|-----------------|---------|
| Default disk size of machine. | Integer         | Integer |

- **Description:** Maximum disk size (in MiB) allocated for Docker Desktop.
- **OS:** {{< badge color=blue text="All" >}}
- **Use case:** Constrain Docker's virtual disk size for storage management.
- **Configure this setting with:**
    - Resources settings in [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## `DiskStats`

| Default value | Accepted values | Format |
|---------------|-----------------|--------|
| `""`          | String          | String |

- **Description:** Disk usage statistics.
- **OS:** {{< badge color=blue text="Mac only" >}}
- **Use case:** Monitor or debug disk usage performance on Unix-based systems.
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## `DiskTRIM`

| Default value | Accepted values | Format |
|---------------|-----------------|--------|
| `true`        | `true`, `false` | Boolean |

- **Description:** Enable TRIM operation support to reclaim unused disk space.
- **OS:** {{< badge color=blue text="All" >}}
- **Use case:** Optimize disk usage over time.
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## `DisplayRestartDialog`

| Default value | Accepted values | Format |
|---------------|-----------------|--------|
| `true`        | `true`, `false` | Boolean |

- **Description:** Show a restart notification when settings changes require a
restart.
- **OS:** {{< badge color=blue text="All" >}}
- **Use case:** Provide user feedback about restart requirements.
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## `DisplaySwitchWinLinContainers`

| Default value | Accepted values | Format |
|---------------|-----------------|--------|
| `false`       | `true`, `false` | Boolean |

- **Description:** Allow users to switch between Linux and Windows containers.
- **OS:** {{< badge color=blue text="Windows only" >}}
- **Use case:** Flexibility in development environments.
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## `Displayed18362Deprecation`

| Default value | Accepted values | Format |
|---------------|-----------------|--------|
| `false`       | `true`, `false` | Boolean |

- **Description:** Show the deprecation warning for Windows build 18362.
- **OS:** {{< badge color=blue text="Windows only" >}}
- **Use case:** Prevent showing the same Windows version deprecation warning
multiple times.
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## `DisplayedElectronPopup`

| Default value | Accepted values     | Format                      |
|---------------|---------------------|-----------------------------|
| `[]`          | List of strings     | Array with list of strings |

- **Description:** Show Electron (tips, alerts, announcements) pop-ups for users.
- **OS:** {{< badge color=blue text="All" >}}
- **Use case:** Prevents Docker Desktop from repeatedly showing the same popup
messages.
- **Configure this setting with:**
    - `settings-store.json` or `settings.json` files

## Display onboarding survey

| Default value | Accepted values | Format |
|---------------|-----------------|--------|
| `false`       | `true`, `false` | Boolean |

- **Description:** Display the onboarding survey for Docker Desktop. If the
value is set to `true`, the onboarding survey will not be displayed to new
users. Settings the value to `false` has no effecct.
- **OS:** {{< badge color=blue text="All" >}}
- **Use case:** Prevents Docker Desktop from repeatedly showing onboarding.
- **Configure this setting with:**
    - `DisplayedOnboarding` in `settings-store.json` or `settings.json` files
    - Settings Management: `displayedOnboarding` setting in the [`admin-settings.json` file](/manuals/security/for-admins/hardened-desktop/settings-management/configure-json-file.md)
    - Settings Management: **Hide onboarding survey** setting in the [Admin Console](/manuals/security/for-admins/hardened-desktop/settings-management/configure-admin-console.md)

## `DockerAppLaunchPath`

| Default value             | Accepted values | Format |
|--------------------------|-----------------|--------|
| `/Applications/Docker.app` | File path       | String |

- **Description:** Path to the Docker Desktop application executable on macOS.
- **OS:** {{< badge color=blue text="Mac only" >}}
- **Use case:** Custom install management or scripting.
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## `DockerBinInstallPath`

| Default value | Accepted values | Format   |
|---------------|-----------------|----------|
| `system`      | File path       | String   |

- **Description:** Install location for Docker CLI binaries.
- **OS:** {{< badge color=blue text="All" >}}
- **Use case:** Customize CLI install location for compliance or tooling.
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## `DockerDebugDefaultEnabled`

| Default value | Accepted values | Format   |
|---------------|-----------------|----------|
| `false`       | `true`, `false` | Boolean  |

- **Description:** Enable debug logging by default for Docker CLI commands.
- **OS:** {{< badge color=blue text="All" >}}
- **Use case:** Assist with debugging support issues.
- **Configure this setting with:**
    - General settings in [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## Allow ECI to use derived images

| Default value | Accepted values | Format   |
|---------------|-----------------|----------|
| `true`        | `true`, `false` | Boolean  |

- **Description:** Allow Enhanced Container Isolation (ECI) to use derived images.
- **OS:** {{< badge color=blue text="All" >}}
- **Use case:** Permit use of base images with layered builds in ECI mode.
- **Configure this setting with:**
    - `ECIDockerSocketAllowDerivedImages` in `settings-store.json` or `settings.json` files
    - Settings Management: **Allow derived images** setting in the [Admin Console](/manuals/security/for-admins/hardened-desktop/settings-management/configure-admin-console.md)

## ECI command list

| Default value | Accepted values | Format                      |
|---------------|-----------------|-----------------------------|
| `[]`          | List of strings | Array with list of strings |

- **Description:** Restricts the commands that containers can issue via the
bind-mounted Docker Engine socket.
- **OS:** {{< badge color=blue text="All" >}}
- **Use case:** Fine-tune developer CLI access in hardened environments.
- **Configure this setting with:**
    - `ECIDockerSocketCmdList` in `settings-store.json` or `settings.json` files
    - Settings Management: `commandList` setting in the [`admin-settings.json` file](/manuals/security/for-admins/hardened-desktop/settings-management/configure-json-file.md)
    - Settings Management: **Command list** setting in the [Admin Console](/manuals/security/for-admins/hardened-desktop/settings-management/configure-admin-console.md)

## ECI command list type

| Default value | Accepted values | Format |
|---------------|-----------------|--------|
| `deny`        | `allow`, `deny` | String |

- **Description:** Whether the ECI command list is an allow-list or deny-list.
- **OS:** {{< badge color=blue text="All" >}}
- **Use case:** Determine behavior of `ECIDockerSocketCmdList`.
- **Configure this setting with:**
    - `ECIDockerSocketCmdListType` in `settings-store.json` or `settings.json` files
    - Settings Management: `????` setting in the [`admin-settings.json` file](/manuals/security/for-admins/hardened-desktop/settings-management/configure-json-file.md)
    - Settings Management: **????** setting in the [Admin Console](/manuals/security/for-admins/hardened-desktop/settings-management/configure-admin-console.md)

## ECI image list

| Default value | Accepted values | Format                      |
|---------------|-----------------|-----------------------------|
| `[]`          | List of strings | Array list of strings       |

- **Description:** 	Indicates which container images are allowed to bind-mount
the Docker Engine socket.
- **OS:** {{< badge color=blue text="All" >}}
- **Use case:** Restrict containers to a known set of images.
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files
    - Settings Management: `imageList` setting in the [`admin-settings.json` file](/manuals/security/for-admins/hardened-desktop/settings-management/configure-json-file.md)
    - Settings Manageament: **Image list** setting in the [Admin Console](/manuals/security/for-admins/hardened-desktop/settings-management/configure-admin-console.md)

## `EnableDefaultDockerSocket`

| Default value | Accepted values | Format   |
|---------------|-----------------|----------|
| `true`        | `true`, `false` | Boolean  |

- **Description:** By default, enhanced container isolation blocks bind-mounting
the Docker Engine socket into containers
(e.g., `docker run -v /var/run/docker.sock:/var/run/docker.sock ...`). This lets
you relax this in a controlled way. See ECI Configuration for more info.
- **OS:** {{< badge color=blue text="All" >}}
- **Use case:** Allow containers to access the Docker socket for scenarios like
Docker-in-Docker or containerized CI agents.
- **Configure this setting with:**
    - ???? in [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files
    - Settings Management: `dockerSocketMount` setting in the [`admin-settings.json` file](/manuals/security/for-admins/hardened-desktop/settings-management/configure-json-file.md)

## `EnableDockerAI`

| Default value | Accepted values | Format   |
|---------------|-----------------|----------|
| `false`       | `true`, `false` | Boolean  |

- **Description:** Enable Docker AI features in the Docker Desktop experience.
- **OS:** {{< badge color=blue text="All" >}}
- **Use case:** Enable or disable AI features like "Ask Gordon".
- **Configure this setting with:**
    - **Features in development** settings in [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `EnableDockerAI` in `settings-store.json` or `settings.json` files
    - Settings Management: `enableDockerAI` setting in the [`admin-settings.json` file](/manuals/security/for-admins/hardened-desktop/settings-management/configure-json-file.md)

## `EnableIntegrationWithDefaultWslDistro`

| Default value | Accepted values | Format   |
|---------------|-----------------|----------|
| `false`       | `true`, `false` | Boolean  |

- **Description:** Automatically integrate Docker with the default WSL
distribution.
- **OS:** {{< badge color=blue text="Windows only" >}} + WSL
- **Use case:** Ensure Docker integrates with default WSL distro automatically.
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `EnableIntegrationWithDefaultWslDistro` in `settings-store.json` or `settings.json` files

## `EnableIntegrityCheck`

| Default value | Accepted values | Format   |
|---------------|-----------------|----------|
| `true`        | `true`, `false` | Boolean  |

- **Description:** Perform integrity checks on Docker Desktop binaries.
- **OS:** {{< badge color=blue text="All" >}}
- **Use case:** Enforce binary verification for security.
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## `EnableSegmentDebug`

| Default value | Accepted values | Format   |
|---------------|-----------------|----------|
| `false`       | `true`, `false` | Boolean  |

- **Description:** Enable debug logging for Docker Desktop’s Segment analytics
events.
- **OS:** {{< badge color=blue text="All" >}}
- **Use case:** Troubleshoot or inspect analytics event delivery during
development or support sessions.
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## `EnableWasmShims`

| Default value | Accepted values | Format   |
|---------------|-----------------|----------|
| `false`       | `true`, `false` | Boolean  |

- **Description:** Enable WebAssembly (Wasm) shims to run Wasm containers.
- **OS:** {{< badge color=blue text="All" >}}
- **Use case:** Run Wasm workloads in Docker Desktop.
- **Configure this setting with:**
    - Features in development settings in [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## Enable Enhanced Container Isolation (ECI)

| Default value | Accepted values | Format   |
|---------------|-----------------|----------|
| `false`       | `true`, `false` | Boolean  |

- **Description:** Enable Enhanced Container Isolation for secure container
execution.
- **OS:** {{< badge color=blue text="All" >}}
- **Use case:** Prevent containers from modifying configuration or sensitive
host areas.
- **Configure this setting with:**
    - **General settings** in [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `EnhancedContainerIsolation` in `settings-store.json` or `settings.json` files
    - Settings Management: `enhancedContainerIsolation` setting in the [`admin-settings.json` file](/manuals/security/for-admins/hardened-desktop/settings-management/configure-json-file.md)
    - Settings Management: **Enable enhanced container isolation** setting in the [Admin Console](/manuals/security/for-admins/hardened-desktop/settings-management/configure-admin-console.md)

## Expose Docker API

| Default value | Accepted values | Format   |
|---------------|-----------------|----------|
| `false`       | `true`, `false` | Boolean  |

- **Description:** Expose the Docker API over TCP on a specified port. If value
is set to `true`, the Docker API is exposed on port 2375. This port is
unauthenticated and should only be enabled if protected by suitable firewall
rules.
- **OS:** {{< badge color=blue text="Windows only" >}}
- **Use case:** Allow non-TLS API access for development/testing.
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `ExposeDockerAPIOnTCP2375` in `settings-store.json` or `settings.json` files
    - Settings Management: `exposeDockerAPIOnTCP2375` setting in the [`admin-settings.json` file](/manuals/security/for-admins/hardened-desktop/settings-management/configure-json-file.md)
    - Settings Management: **Expose Docker API** setting in the [Admin Console](/manuals/security/for-admins/hardened-desktop/settings-management/configure-admin-console.md)

## Enable extensions

| Default value | Accepted values | Format   |
|---------------|-----------------|----------|
| `true`        | `true`, `false` | Boolean  |

- **Description:** Enable or disable Docker Extensions.
- **OS:** {{< badge color=blue text="All" >}}
- **Use case:** Control access to the Extensions Marketplace and installed
extensions.
- **Configure this setting with:**
    - **Extensions** settings in [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `ExtensionsEnabled` in `settings-store.json` or `settings.json` files
    - Settings Management: `extensionsEnabled` setting in the [`admin-settings.json` file](/manuals/security/for-admins/hardened-desktop/settings-management/configure-json-file.md)
    - Settings Management: **Allow Extensions** setting in the [Admin Console](/manuals/security/for-admins/hardened-desktop/settings-management/configure-admin-console.md)

## Enable private extensions marketplace

| Default value | Accepted values | Format   |
|---------------|-----------------|----------|
| `false`       | `true`, `false` | Boolean  |

- **Description:** Enable a private marketplace for Docker Extensions.
- **OS:** {{< badge color=blue text="All" >}}
- **Use case:** Restrict extension installation to curated extensions.
- **Configure this setting with:**
    - `ExtensionsPrivateMarketplace` in `settings-store.json` or `settings.json` files

## Set private extension contact URL

| Default value | Accepted values | Format   |
|---------------|-----------------|----------|
| `""`          | URL string      | String   |

- **Description:** Set a contact URL for admins on the private extensions
marketplace page.
- **OS:** {{< badge color=blue text="All" >}}
- **Use case:** Help users contact support if they can’t find an extension.
- **Configure this setting with:**
    - `ExtensionsPrivateMarketplaceAdminContactURL` in `settings-store.json` or `settings.json` files

## Filesharing directories

| Default value                           | Accepted values                 | Format                  |
|----------------------------------------|---------------------------------|--------------------------|
| Varies by OS                           | List of file paths as strings   | Array list of strings   |

- **Description:** List of allowed directories shared between the host and
containers. When a path is added, its subdirectories are allowed.
- **OS:** {{< badge color=blue text="All" >}}
- **Use case:** Restrict or define what file paths are available to containers.
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `FilesharingDirectories` in `settings-store.json` or `settings.json` files
    - Settings Management: `filesharingAllowedDirectories` setting in the [`admin-settings.json` file](/manuals/security/for-admins/hardened-desktop/settings-management/configure-json-file.md)
    - Settings Management: **Allowed file sharing directories** setting in the [Admin Console](/manuals/security/for-admins/hardened-desktop/settings-management/configure-admin-console.md)

## `HostNetworkingEnabled`

| Default value | Accepted values | Format   |
|---------------|-----------------|----------|
| `false`       | `true`, `false` | Boolean  |

- **Description:** Enable experimental host networking support.
- **OS:** {{< badge color=blue text="All" >}}
- **Use case:** Allow containers to use the host network stack.
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## Enable Kubernetes

| Default value | Accepted values | Format   |
|---------------|-----------------|----------|
| `false`       | `true`, `false` | Boolean  |

- **Description:** Enable the integrated Kubernetes cluster in Docker Desktop.
- **OS:** {{< badge color=blue text="All" >}}
- **Use case:** Enable or disable Kubernetes support for developers.
- **Configure this setting with:**
    - **Kubernetes** settings in [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `KubernetesEnabled` in `settings-store.json` or `settings.json` files
    - Settings Management: `kubernetes` setting in the [`admin-settings.json` file](/manuals/security/for-admins/hardened-desktop/settings-management/configure-json-file.md)
    - Settings Management: **Allow Kubernetes** setting in the [Admin Console](/manuals/security/for-admins/hardened-desktop/settings-management/configure-admin-console.md)

## `KubernetesImagesRepository`

| Default value | Accepted values | Format   |
|---------------|-----------------|----------|
| `""`          | URL string      | String   |

- **Description:** Set a custom repository for Kubernetes images.
- **OS:** {{< badge color=blue text="All" >}}
- **Use case:** Support Kubernetes use in restricted or offline environments.
- **Configure this setting with:**
    - `settings-store.json` or `settings.json` files
    - Settings Management: `imagesRepository` setting in the [`admin-settings.json` file](/manuals/security/for-admins/hardened-desktop/settings-management/configure-json-file.md)
    - Settings Management: **Kubernetes images repository** setting in the [Admin Console](/manuals/security/for-admins/hardened-desktop/settings-management/configure-admin-console.md)

## Set Kubernetes mode

| Default value | Accepted values | Format |
|---------------|-----------------|--------|
| `kubeadm`     | String          | String |

- **Description:** Set the Kubernetes node mode (single-node or multi-node).
- **OS:** {{< badge color=blue text="All" >}}
- **Use case:** Control the topology of the integrated Kubernetes cluster.
- **Configure this setting with:**
    - **Kubernetes** settings in [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `KubernetesMode` in `settings-store.json` or `settings.json` files

## Kubernetes node count

| Default value | Accepted values | Format  |
|---------------|-----------------|---------|
| `1`           | Integer         | Integer |

- **Description:** Number of nodes to create in a multi-node Kubernetes cluster.
- **OS:** {{< badge color=blue text="All" >}}
- **Use case:** Scale the number of Kubernetes nodes for development or testing.
- **Configure this setting with:**
    - **Kubernetes** settings in [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `KubernetesNodesCount` in `settings-store.json` or `settings.json` files

## Kubernetes node version

| Default value | Accepted values               | Format |
|---------------|-------------------------------|--------|
| `1.31.1`      | Semantic version (e.g., 1.29.1) | String |

- **Description:** Version of Kubernetes used for cluster node creation.
- **OS:** {{< badge color=blue text="All" >}}
- **Use case:** Pin a specific Kubernetes version for consistency or
compatibility.
- **Configure this setting with:**
    - **Kubernetes** settings in [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `KubernetesNodesVersion` in `settings-store.json` or `settings.json` files

## `LastLoginDate`

| Default value | Accepted values | Format                |
|---------------|-----------------|-----------------------|
| `0`           | `int64` values  | Integer in `int64` format |

- **Description:** Timestamp of last successful Docker Desktop login.
- **OS:** {{< badge color=blue text="All" >}}
- **Use case:** Display usage activity.
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## `LatestBannerKey`

| Default value | Accepted values | Format |
|---------------|-----------------|--------|
| `""`          | String          | String |

- **Description:** Tracks the most recently shown in-app banner.
- **OS:** {{< badge color=blue text="All" >}}
- **Use case:** Prevent repeated display of the same banner across sessions.
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## `LicenseTermsVersion`

| Default value | Accepted values | Format  |
|---------------|-----------------|---------|
| `0`           | Integer         | Integer |

- **Description:** Version of Docker Desktop license terms accepted by the user.
- **OS:** {{< badge color=blue text="All" >}}
- **Use case:** Audit license terms agreement.
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## `LifecycleTimeoutSeconds`

| Default value | Accepted values     | Format  |
|---------------|---------------------|---------|
| `600`         | Integer (seconds)   | Integer |

- **Description:** Number of seconds Docker Desktop waits for the Docker Engine
to start before timing out.
- **OS compatibility**: All
- **Use case:** Extend or reduce the timeout window for environments where the
engine may start slowly.
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## `MemoryMiB`

| Default value              | Accepted values | Format  |
|---------------------------|-----------------|---------|
| Based on system resources | Integer         | Integer |

- **Description:** Amount of RAM (in MiB) assigned to the Docker virtual machine.
- **OS:** {{< badge color=blue text="All" >}}
- **Use case:** Control how much memory Docker can use on the host.
- **Configure this setting with:**
    - Resources settings in [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## Allow only Marketplace extensions

| Default value | Accepted values | Format   |
|---------------|-----------------|----------|
| `true`        | `true`, `false` | Boolean  |

- **Description:** Restrict Docker Desktop to only run Marketplace extensions.
- **OS:** {{< badge color=blue text="All" >}}
- **Use case:** Prevent running third-party or local extensions.
- **Configure this setting with:**
    - **Extensions** settings in [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `OnlyMarketplaceExtensions` in `settings-store.json` or `settings.json` files

## `OpenUIOnStartupDisabled`

| Default value | Accepted values | Format   |
|---------------|-----------------|----------|
| `false`       | `true`, `false` | Boolean  |

- **Description:** Prevent the Docker Desktop UI from opening automatically at
startup.
- **OS:** {{< badge color=blue text="All" >}}
- **Use case:** Streamline startup experience.
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## Override proxy exclude

| Default value | Accepted values | Format |
|---------------|-----------------|--------|
| `""`          | String          | String |

- **Description:** Comma-separated list of domain patterns that should bypass
the proxy.
- **OS:** {{< badge color=blue text="All" >}}
- **Use case:** Exclude internal services or domains from being routed through
the proxy.
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `OverrideProxyExclude` in `settings-store.json` or `settings.json` files

## Override proxy HTTP

| Default value | Accepted values | Format   |
|---------------|-----------------|----------|
| `""`          | URL string      | String   |

- **Description:** Override the default HTTP proxy used by Docker Desktop and
its containers.
- **OS compatibility**: All
- **Use case:** Route container HTTP traffic through a specific proxy.
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `OverrideProxyHTTP` in `settings-store.json` or `settings.json` files

## Override proxy HTTPS

| Default value | Accepted values | Format   |
|---------------|-----------------|----------|
| `""`          | URL string      | String   |

- **Description:** Override the default HTTPS proxy used by Docker Desktop and
its containers.
- **OS:** {{< badge color=blue text="All" >}}
- **Use case:** Route container HTTPS traffic through a specific proxy.
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `OverrideProxyHTTPS` in `settings-store.json` or `settings.json` files

## `OverrideProxyPAC`

| Default value | Accepted values | Format   |
|---------------|-----------------|----------|
| `""`          | URL string      | String   |

- **Description:** URL to a Proxy Auto-Config (PAC) file to dynamically
configure proxy rules.
- **OS:** {{< badge color=blue text="All" >}}
- **Use case:** Load dynamic proxy rules from a PAC file.
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## `OverrideProxyTCP`

| Default value | Accepted values | Format   |
|---------------|-----------------|----------|
| `""`          | URL string      | String   |

- **Description:** Override the TCP proxy settings used by Docker Desktop.
- **OS:** {{< badge color=blue text="All" >}}
- **Use case:** Define a custom proxy for TCP traffic not covered by
HTTP/HTTPS proxies.
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## `OverrideWindowsDockerdPort`

| Default value | Accepted values | Format |
|---------------|-----------------|--------|
| `-1`          | Integer         | Integer |

- **Description:** Exposes Docker Desktop's internal proxy locally on this port
for the Windows Docker daemon to connect to. If set to `0`, a random free port
is chosen. If the value is greater than 0, it uses that exact value for the port.
-1 disables the option.
- **OS:** Windows
- **Use case:** Allow precise control of how Docker Desktop exposes its
internal proxy for `dockerd.exe`.
- **Configure this settings with:**
    - `OverrideWindowsDockerdPort` in `settings-store.json` or `settings.json` files
    - Settings Management: `windowsDockerdPort` setting in the [`admin-settings.json` file](/manuals/security/for-admins/hardened-desktop/settings-management/configure-json-file.md)
    - Settings Management: **Override Windows “dockerd” port** setting in the [Admin Console](/manuals/security/for-admins/hardened-desktop/settings-management/configure-admin-console.md)

## Proxy enable Kerberos NTLM

| Default value | Accepted values | Format   |
|---------------|-----------------|----------|
| `false`       | `true`, `false` | Boolean  |

- **Description:** Enable Kerberos and NTLM authentication for the proxy.
- **OS:** Windows
- **Use case:** Allow Docker Desktop to authenticate with enterprise proxies
that require Kerberos or NTLM credentials.
- **Configure this setting with:**
    - `ProxyEnableKerberosNTLM` in `settings-store.json` or `settings.json` files
    - Settings Management: `enableKerberosNtlm` setting in the [`admin-settings.json` file](/manuals/security/for-admins/hardened-desktop/settings-management/configure-json-file.md)
    - Settings Management: **Kerberos NTLM** setting in the [Admin Console](/manuals/security/for-admins/hardened-desktop/settings-management/configure-admin-console.md)

## Proxy HTTP mode

| Default value | Accepted values     | Format |
|---------------|---------------------|--------|
| `system`      | `system`, `manual`  | String |

- **Description:** Proxy mode setting. If mode is set to `system` instead of
`manual`, Docker Desktop gets the proxy values from the system and ignores
values set for `http`, `https`, and `exclude`. To manually configure proxy
servers, use `manual`.
- **OS:** {{< badge color=blue text="All" >}}
- **Use case:** Control how Docker Desktop uses or ignores system proxy settings.
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `ProxyHTTPMode` in `settings-store.json` or `settings.json` files
    - Settings Management: `proxy` setting in the [`admin-settings.json` file](/manuals/security/for-admins/hardened-desktop/settings-management/configure-json-file.md)
    - Settings Management: [Admin Console](/manuals/security/for-admins/hardened-desktop/settings-management/configure-admin-console.md)

## `ProxyLocalhostPort`

| Default value | Accepted values     | Format  |
|---------------|---------------------|---------|
| `0`           | Integer (port number) | Integer |

- **Description:** Specifies the local port used by Docker Desktop’s internal
proxy to route container traffic through the host network.
- **OS:** {{< badge color=blue text="All" >}}
- **Use case:** Bind the internal proxy to a fixed localhost port for debugging
or compatibility with network security tools.
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## `RequireVmnetd`

| Default value | Accepted values | Format   |
|---------------|-----------------|----------|
| `true`        | `true`, `false` | Boolean  |

- **Description:** Require the privileged helper (`vmnetd`) for networking on
macOS.
- **OS:** {{< badge color=blue text="Mac only" >}}
- **Use case:** Enforce elevated privileges for networking support
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## `RunWinServiceInWslMode`

| Default value | Accepted values | Format   |
|---------------|-----------------|----------|
| `false`       | `true`, `false` | Boolean  |

- **Description:** Allow the Windows service that supports Docker Desktop to
run in WSL mode for enhanced integration.
- **OS:** {{< badge color=blue text="Windows only" >}} + WSL
- **Use case:** Enable deeper integration between the Windows service layer and
the WSL-based Docker backend.
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## SBOM indexing

| Default value | Accepted values | Format   |
|---------------|-----------------|----------|
| `true`        | `true`, `false` | Boolean  |

- **Description:** Enable SBOM indexing for container images
- **OS:** {{< badge color=blue text="All" >}}
- **Use case:** Control whether Docker indexes SBOMs for images
- **Configure this setting with:**
    - **General settings** in [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `SbomIndexing` in `settings-store.json` or `settings.json` files
    - Settings Management: `sbomIndexing` setting in the [`admin-settings.json` file](/manuals/security/for-admins/hardened-desktop/settings-management/configure-json-file.md)
    - Settings Management: **SBOM indexing** settings in the [Admin Console](/manuals/security/for-admins/hardened-desktop/settings-management/configure-admin-console.md)

## `ScoutNotificationPopupsEnabled`

| Default value | Accepted values | Format   |
|---------------|-----------------|----------|
| `true`        | `true`, `false` | Boolean  |

- **Description:** Enable Docker Scout popups inside Docker Desktop.
- **OS:** {{< badge color=blue text="All" >}}
- **Use case:** Show or hide vulnerability scan notifications
- **Configure this setting with:**
    - Notifications settings in [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## `ScoutOsNotificationsEnabled`

| Default value | Accepted values | Format   |
|---------------|-----------------|----------|
| `false`       | `true`, `false` | Boolean  |

- **Description:** Enable Docker Scout notifications through the operating system.
- **OS:** {{< badge color=blue text="All" >}}
- **Use case:** Push Scout updates via system notification center
- **Configure this setting with:**
    - Notifications settings in [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## `SettingsVersion`

| Default value         | Accepted values | Format  |
|-----------------------|-----------------|---------|
| `CurrentSettingsVersions` | Integer         | Integer |

- **Description:** Specifies the version of the settings configuration file format
- **OS:** {{< badge color=blue text="All" >}}
- **Use case:** Track schema versions for compatibility
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files
    - Settings Management: `configurationFileVersion` setting in the [`admin-settings.json` file](/manuals/security/for-admins/hardened-desktop/settings-management/configure-json-file.md)

## `ShowAnnouncementNotifications`

| Default value | Accepted values | Format   |
|---------------|-----------------|----------|
| `true`        | `true`, `false` | Boolean  |

- **Description:** Display general announcements inside Docker Desktop.
- **OS:** {{< badge color=blue text="All" >}}
- **Use case:** Enable or suppress Docker-wide announcements in the UI.
- **Configure this setting with:**
    - Notifications settings in [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## `ShowExtensionsSystemContainers`

| Default value | Accepted values | Format   |
|---------------|-----------------|----------|
| `false`       | `true`, `false` | Boolean  |

- **Description:** Show system containers used by Docker Extensions in the container list
- **OS:** {{< badge color=blue text="All" >}}
- **Use case:** Help developers troubleshoot or view extension system containers
- **Configure this setting with:**
    - Extensions settings in [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## `ShowGeneralNotifications`

| Default value | Accepted values | Format   |
|---------------|-----------------|----------|
| `true`        | `true`, `false` | Boolean  |

- **Description:** Display general informational messages inside Docker Desktop
- **OS:** {{< badge color=blue text="All" >}}
- **Use case:** Customize in-app communication visibility
- **Configure this setting with:**
    - Notifications settings in [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## `ShowInstallScreen`

| Default value | Accepted values | Format   |
|---------------|-----------------|----------|
| `true`        | `true`, `false` | Boolean  |

- **Description:** Show the installation onboarding screen in Docker Desktop
- **OS:** {{< badge color=blue text="All" >}}
- **Use case:** Control whether onboarding screens are shown after installation
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## `ShowKubernetesSystemContainers`

| Default value | Accepted values | Format   |
|---------------|-----------------|----------|
| `false`       | `true`, `false` | Boolean  |

- **Description:** Show Kubernetes system containers in the Docker Dashboard container list
- **OS:** {{< badge color=blue text="All" >}}
- **Use case:** Allow developers to view kube-system containers for debugging
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## `ShowPromotionalNotifications`

| Default value | Accepted values | Format   |
|---------------|-----------------|----------|
| `true`        | `true`, `false` | Boolean  |

- **Description:** Display promotional announcements and banners inside Docker Desktop
- **OS:** {{< badge color=blue text="All" >}}
- **Use case:** Control exposure to Docker news and feature promotion
- **Configure this setting with:**
    - Notifications settings in [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## `ShowSurveyNotifications`

| Default value | Accepted values | Format   |
|---------------|-----------------|----------|
| `true`        | `true`, `false` | Boolean  |

- **Description:** Display notifications inviting users to participate in surveys
- **OS:** {{< badge color=blue text="All" >}}
- **Use case:** Enable or disable in-product survey prompts
- **Configure this setting with:**
    - Notifications settings in [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## `SkipUpdateToWSLPrompt`

| Default value | Accepted values | Format   |
|---------------|-----------------|----------|
| `false`       | `true`, `false` | Boolean  |

- **Description:** Skip prompting users to upgrade to the WSL 2 backend
- **OS:** {{< badge color=blue text="Windows only" >}} + WSL
- **Use case:** Silence UI nudges to switch WSL versions
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## `SkipWSLMountPerfWarning`

| Default value | Accepted values | Format   |
|---------------|-----------------|----------|
| `false`       | `true`, `false` | Boolean  |

- **Description:** Skip the performance warning about WSL mount speed.
- **OS:** {{< badge color=blue text="Windows only" >}} + WSL
- **Use case:** Suppress warnings for known limitations or user preference
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## `SocksProxyPort`

| Default value | Accepted values | Format  |
|---------------|-----------------|---------|
| `0`           | Integer (port)  | Integer |

- **Description:** Local SOCKS proxy port for Docker Desktop.
- **OS:** {{< badge color=blue text="All" >}}
- **Use case:** Route Docker traffic through a SOCKS proxy
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## `SwapMiB`

| Default value | Accepted values | Format  |
|---------------|-----------------|---------|
| `1024`        | Integer         | Integer |

- **Description:** Amount of swap space (in MiB) assigned to the Docker virtual machine
- **OS:** {{< badge color=blue text="All" >}}
- **Use case:** Extend memory availability via swap
- **Configure this setting with:**
    - Resources settings in [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## `SynchronizedDirectories`

| Default value                    | Accepted values             | Format |
|----------------------------------|-----------------------------|--------|
| Varies by system/user configs    | Array of file paths as strings | Array |

- **Description:** Directories that should be synchronized between host and
container filesystems.
- **OS:** {{< badge color=blue text="All" >}}
- **Use case:** Improve performance for bind mounts and volume sharing.
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## `ThemeSource`

| Default value | Accepted values            | Format |
|---------------|----------------------------|--------|
| `system`      | `light`, `dark`, `system`  | Enum   |

- **Description:** Choose the Docker Desktop UI theme.
- **OS:** {{< badge color=blue text="All" >}}
- **Use case:** Personalize Docker Desktop appearance.
- **Configure this setting with:**
    - General settings in [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## `UpdateAvailableTime`

| Default value | Accepted values     | Format |
|---------------|---------------------|--------|
| `0`           | ISO 8601 timestamp  | String |

- **Description:** Timestamp of last update availability check.
- **OS:** {{< badge color=blue text="All" >}}
- **Use case:** Telemetry and internal logic.
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## `UpdateHostsFile`

| Default value | Accepted values | Format   |
|---------------|-----------------|----------|
| `false`       | `true`, `false` | Boolean  |

- **Description:** Allow Docker Desktop to update the system `hosts` file.
- **OS:** {{< badge color=blue text="All" >}}
- **Use case:** Support DNS resolution for internal services.
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## `UpdateInstallTime`

| Default value | Accepted values     | Format |
|---------------|---------------------|--------|
| `0`           | ISO 8601 timestamp  | String |

- **Description:** Timestamp of last Docker Desktop update installation.
- **OS:** {{< badge color=blue text="All" >}}
- **Use case:** Track install history.
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## Use background indexing

| Default value | Accepted values | Format   |
|---------------|-----------------|----------|
| `false`       | `true`, `false` | Boolean  |

- **Description:** Enable background indexing of local Docker images for Docker
Scout.
- **OS:** {{< badge color=blue text="All" >}}
- **Use case:** Improve performance of features like image search.
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `UseBackgroundIndexing` in `settings-store.json` or `settings.json` files
    - Settings Management: `useBackgroundIndexing` setting in the [`admin-settings.json` file](/manuals/security/for-admins/hardened-desktop/settings-management/configure-json-file.md)
    - Settings Management: **Background indexing** setting in the [Admin Console](/manuals/security/for-admins/hardened-desktop/settings-management/configure-admin-console.md)

## `UseContainerdSnapshotter`

| Default value | Accepted values | Format   |
|---------------|-----------------|----------|
| `false`       | `true`, `false` | Boolean  |

- **Description:** Use containerd native snapshotter instead of legacy
snapshotters.
- **OS:** {{< badge color=blue text="All" >}}
- **Use case:** Improve image handling performance and compatibility.
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## `UseCredentialHelper`

| Default value | Accepted values | Format   |
|---------------|-----------------|----------|
| `true`        | `true`, `false` | Boolean  |

- **Description:** Use the configured credential helper to securely store and
retrieve Docker registry credentials.
- **OS:** {{< badge color=blue text="All" >}}
- **Use case:** Enable secure, system-integrated storage of Docker login
credentials instead of plain-text config files.
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## Use gRPC Fuse

| Default value | Accepted values | Format   |
|---------------|-----------------|----------|
| `true`        | `true`, `false` | Boolean  |

- **Description:** Enable gRPC FUSE for macOS file sharing. If value is set to
`true`, gRPC Fuse is set as the file sharing mechanism.
- **OS:** {{< badge color=blue text="Mac only" >}}
- **Use case:** Improve performance and compatibility of file mounts.
- **Configure this setting with:**
    - **Choose file sharing implementation for your containers** setting in [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `UseGrpcfuse` in `settings-store.json` or `settings.json` files
    - Settings Management: `useGrpcfuse` setting in the [`admin-settings.json` file](/manuals/security/for-admins/hardened-desktop/settings-management/configure-json-file.md)
    - Settings Management: **Use gRPC FUSE for file sharing** setting in the [Admin Console](/manuals/security/for-admins/hardened-desktop/settings-management/configure-admin-console.md)

## `UseLibkrun`

| Default value | Accepted values | Format   |
|---------------|-----------------|----------|
| `false`       | `true`, `false` | Boolean  |

- **Description:** Enable lightweight VM virtualization via libkrun.
- **OS:** {{< badge color=blue text="Mac only" >}}
- **Use case:** Run containers in microVMs using libkrun.
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## Use nightly build updates

| Default value | Accepted values | Format   |
|---------------|-----------------|----------|
| `false`       | `true`, `false` | Boolean  |

- **Description:** Enable updates from the Docker Desktop nightly build channel
instead of the stable release channel.
- **OS:** {{< badge color=blue text="All" >}}
- **Use case:** Receive early access to experimental features and fixes by
subscribing to nightly builds.
- **Configure this setting with:**
    - `UseNightlyBuildUpdates` in `settings-store.json` or `settings.json` files

## `UseResourceSaver`

| Default value | Accepted values | Format   |
|---------------|-----------------|----------|
| `true`        | `true`, `false` | Boolean  |

- **Description:** Enable Docker Desktop to pause when idle.
- **OS:** {{< badge color=blue text="All" >}}
- **Use case:** Save system resources during periods of inactivity.
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## `UseVirtualizationFramework`

| Default value | Accepted values | Format   |
|---------------|-----------------|----------|
| `true`        | `true`, `false` | Boolean  |

- **Description:** Use Apple Virtualization Framework to run Docker containers.
- **OS:** {{< badge color=blue text="Mac only" >}}
- **Use case:** Improve VM performance on Apple Silicon.
- **Configure this setting with:**
    - General settings in [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## Use virtualization framework: Rosetta

| Default value | Accepted values | Format   |
|---------------|-----------------|----------|
| `true`        | `true`, `false` | Boolean  |

- **Description:** Use Rosetta to emulate `amd64` on Apple Silicon. If value
is set to `true`, Docker Desktop turns on Rosetta to accelerate
x86_64/amd64 binary emulation on Apple Silicon.
- **OS:** {{< badge color=blue text="Mac only" >}} 13+
- **Use case:** Run Intel-based containers on Apple Silicon hosts.
- **Configure this setting with:**
    - **General settings** in [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `UseVirtualizationFrameworkRosetta` in `settings-store.json` or `settings.json` files
    - Settings Management:`useVirtualizationFrameworkRosetta` setting in the [`admin-settings.json` file](/manuals/security/for-admins/hardened-desktop/settings-management/configure-json-file.md)
    - Settings Management: **Use Rosetta for x86_64/amd64 emulation on Apple Silicon** setting in the [Admin Console](/manuals/security/for-admins/hardened-desktop/settings-management/configure-admin-console.md)

## Use virtualization framework: VirtioFS

| Default value | Accepted values | Format   |
|---------------|-----------------|----------|
| `true`        | `true`, `false` | Boolean  |

- **Description:** Use VirtioFS for fast, native file sharing between host and
containers. If value is set to `true`, VirtioFS is set as the file sharing
mechanism. If both VirtioFS and gRPC are set to `true`, VirtioFS takes
precedence.
- **OS:** {{< badge color=blue text="Mac only" >}} 12.5+
- **Use case:** Improve volume mount performance and compatibility.
- **Configure this setting with:**
    - **General settings** in [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `UseVirtualizationFrameworkVirtioFS` in `settings-store.json` or `settings.json` files
    - Settings Management: `useVirtualizationFrameworkVirtioFS` setting in the [`admin-settings.json` file](/manuals/security/for-admins/hardened-desktop/settings-management/configure-json-file.md)
    - Settings Management: **Use VirtioFS for file sharing** setting in the [Admin Console](/manuals/security/for-admins/hardened-desktop/settings-management/configure-admin-console.md)

## `UseVpnkit`

| Default value | Accepted values | Format   |
|---------------|-----------------|----------|
| `true`        | `true`, `false` | Boolean  |

- **Description:** Use vpnkit for Docker Desktop networking on macOS.
- **OS:** {{< badge color=blue text="Mac only" >}}
- **Use case:** Enable or disable vpnkit as the networking backend.
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## `UseWindowsContainers`

| Default value | Accepted values | Format   |
|---------------|-----------------|----------|
| `false`       | `true`, `false` | Boolean  |

- **Description:** Enable Windows container mode in Docker Desktop.
- **OS:** {{< badge color=blue text="Windows only" >}}
- **Use case:** Switch between Linux and Windows container runtimes.
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files
    - Settings Management: `windowContainters` setting in the [`admin-settings.json` file](/manuals/security/for-admins/hardened-desktop/settings-management/configure-json-file.md)

## `VpnKitAllowedBindAddresses`

| Default value | Accepted values | Format   |
|---------------|-----------------|----------|
| `0.0.0.0`     | IP address      | String   |

- **Description:** Specify which local IP addresses vpnkit is allowed to bind
to for handling network traffic.
- **OS:** {{< badge color=blue text="All" >}}
- **Use case:** Restrict or allow vpnkit to bind to specific interfaces for
security or debugging purposes.
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## `VpnKitMTU`

| Default value | Accepted values | Format  |
|---------------|-----------------|---------|
| `1500`        | Integer         | Integer |

- **Description:** Set the Maximum Transmission Unit (MTU) for vpnkit’s virtual
network interface.
- **OS:** {{< badge color=blue text="All" >}}
- **Use case:** Tune network performance or resolve issues with packet
fragmentation when using vpnkit.
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## `VpnKitMaxConnections`

| Default value | Accepted values | Format  |
|---------------|-----------------|---------|
| `2000`        | Integer         | Integer |

- **Description:** Set the maximum number of simultaneous network connections
vpnkit can handle.
- **OS:** {{< badge color=blue text="All" >}}
- **Use case:** Control resource usage or support high-connection workloads
inside containers.
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## `VpnKitMaxPortIdleTime`

| Default value | Accepted values     | Format  |
|---------------|---------------------|---------|
| `300`         | Integer (seconds)   | Integer |

- **Description:** Maximum idle time in seconds before vpnkit closes an
unused port.
- **OS:** {{< badge color=blue text="All" >}}
- **Use case:** Improve performance and free up unused ports by closing
idle connections.
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## `VpnKitTransparentProxy`

| Default value | Accepted values | Format   |
|---------------|-----------------|----------|
| `true`        | `true`, `false` | Boolean  |

- **Description:** Enable transparent proxying in vpnkit.
- **OS:** {{< badge color=blue text="Mac only" >}}
- **Use case:** Seamlessly forward traffic through proxies using vpnkit.
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## vpnkit CIDR

| Default value     | Accepted values | Format |
|-------------------|-----------------|--------|
| `192.168.65.0/24` | IP address      | String |

- **Description:** Overrides the network range used for vpnkit DHCP/DNS for
`*.docker.internal`.
- **OS:** {{< badge color=blue text="Mac only" >}}
- **Use case:** Customize the subnet used for Docker container networking.
- **Configure this setting with:**
    - `VpnkitCIDR` in `settings-store.json` or `settings.json` files
    - Settings Management: `vpnkitCIDR` setting in the [`admin-settings.json` file](/manuals/security/for-admins/hardened-desktop/settings-management/configure-json-file.md)
    - Settings Management: **VPN Kit CIDR** setting in the [Admin Console](/manuals/security/for-admins/hardened-desktop/settings-management/configure-admin-console.md)


## `WslDiskCompactionThresholdGb`

| Default value | Accepted values | Format  |
|---------------|-----------------|---------|
| `0`           | Integer (GB)    | Integer |

- **Description:** Minimum free disk space required to trigger WSL disk
compaction.
- **OS:** {{< badge color=blue text="Windows only" >}} + WSL
- **Use case:** Automatically reclaim unused space from WSL disks.
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## `WslEnableGrpcfuse`

| Default value | Accepted values | Format   |
|---------------|-----------------|----------|
| `false`       | `true`, `false` | Boolean  |

- **Description:** Enable gRPC FUSE file sharing in WSL2 mode.
- **OS:** {{< badge color=blue text="Windows only" >}} + WSL
- **Use case:** Improve performance and compatibility for file mounts in WSL.
- **Configure this setting with:**
    - General settings in [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## Enable WSL engine

| Default value | Accepted values | Format   |
|---------------|-----------------|----------|
| `true`        | `true`, `false` | Boolean  |

- **Description:** If the value is set to `true`, Docker Desktop uses the WSL2
based engine. This overrides anything that may have been set at installation
using the `--backend=<backend name>` flag.
- **OS:** {{< badge color=blue text="Windows only" >}} + WSL
- **Use case:** Enable Linux containers via WSL 2 backend.
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `WslEngineEnabled` in `settings-store.json` or `settings.json` files
    - Settings Management: `wslEngineEnabled` setting in the [`admin-settings.json` file](/manuals/security/for-admins/hardened-desktop/settings-management/configure-json-file.md)
    - Settings Management: **Windows Subsystem for Linux (WSL) Engine** setting in the [Admin Console](/manuals/security/for-admins/hardened-desktop/settings-management/configure-admin-console.md)

## `WslInstallMode`

| Default value       | Accepted values                | Format |
|---------------------|--------------------------------|--------|
| `installLatestWsl`  | `installLatestWsl`, `manualInstall` | String |

- **Description:** Select how Docker Desktop installs and manages WSL on
Windows systems.
- **OS:** {{< badge color=blue text="Windows only" >}} + WSL
- **Use case:** Control whether Docker Desktop installs WSL automatically or
relies on a pre-installed version.
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files

## `WslUpdateRequired`

| Default value | Accepted values | Format   |
|---------------|-----------------|----------|
| `false`       | `true`, `false` | Boolean  |

- **Description:** Indicates whether a WSL update is required for Docker Desktop
to function.
- **OS:** {{< badge color=blue text="Windows only" >}} + WSL
- **Use case:** Internal check for platform support.
- **Configure this setting with:**
    - [Docker Desktop UI](/manuals/desktop/settings-and-maintenance/settings.md)
    - `settings-store.json` or `settings.json` files
