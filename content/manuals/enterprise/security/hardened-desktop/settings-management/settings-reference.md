---
description: Reference for all settings and features that are configured with Settings Management
keywords: admin, controls, settings management, reference
title: Settings reference
linkTitle: Settings reference
aliases: 
 - /security/for-admins/hardened-desktop/settings-management/settings-reference/
---

This reference lists all Docker Desktop settings, including where they are configured, which operating systems they apply to, and whether they're available in the Docker Desktop GUI, the Docker Admin Console, or the `admin-settings.json` file. Settings are grouped to match the structure of the Docker Desktop interface.

Each setting includes:

- The display name used in Docker Desktop
- A table of values, default values, and required format
- A description and use cases
- OS compatibility
- Configuration methods: via [Docker Desktop](/manuals/desktop/settings-and-maintenance/settings.md), the Admin Console, or the `admin-settings.json` file

Use this reference to compare how settings behave across different configuration
methods and platforms.

## General

### Start Docker Desktop when you sign in to your computer

| Default value | Accepted values | Format |
|---------------|-----------------|--------|
| `false`       | `true`, `false` | Boolean |

- **Description:** Start Docker Desktop automatically when booting machine.
- **OS:** {{< badge color=blue text="All" >}}
- **Use case:** Ensure Docker Desktop is always running after boot.
- **Configure this setting with:**
    - **General** settings in [Docker Desktop GUI](/manuals/desktop/settings-and-maintenance/settings.md)

### Open Docker Dashboard when Docker Desktop starts

| Default value | Accepted values            | Format |
|---------------|----------------------------|--------|
| `false`      | `true`, `false`  | Boolean   |

- **Description:** Open the Docker Dashboard automatically when Docker Desktop starts.
- **OS:** {{< badge color=blue text="All" >}}
- **Use case:** Quickly access containers, images, and volumes in the Docker Dashboard after starting Docker Desktop.
- **Configure this setting with:**
    - **General** settings in [Docker Desktop GUI](/manuals/desktop/settings-and-maintenance/settings.md)

### Choose theme for Docker Desktop

| Default value | Accepted values            | Format |
|---------------|----------------------------|--------|
| `system`      | `light`, `dark`, `system`  | Enum   |

- **Description:** Choose the Docker Desktop GUI theme.
- **OS:** {{< badge color=blue text="All" >}}
- **Use case:** Personalize Docker Desktop appearance.
- **Configure this setting with:**
    - **General** settings in [Docker Desktop GUI](/manuals/desktop/settings-and-maintenance/settings.md)

### Configure shell completions

| Default value | Accepted values         | Format |
|---------------|-------------------------|--------|
| `integrated`  | `integrated`, `system`  | String |

- **Description:** If installed, automatically edits your shell configuration.
- **OS:** {{< badge color=blue text="All" >}}
- **Use case:** Customize developer experience with shell completions.
- **Configure this setting with:**
    - **General** settings in [Docker Desktop GUI](/manuals/desktop/settings-and-maintenance/settings.md)

### Choose container terminal

| Default value | Accepted values         | Format |
|---------------|-------------------------|--------|
| `integrated`  | `integrated`, `system`  | String |

- **Description:** Select default terminal for launching Docker CLI from Docker
Desktop.
- **OS:** {{< badge color=blue text="All" >}}
- **Use case:** Customize developer experience with preferred terminal.
- **Configure this setting with:**
    - **General** settings in [Docker Desktop GUI](/manuals/desktop/settings-and-maintenance/settings.md)

### Enable Docker terminal

| Default value | Accepted values | Format |
|---------------|-----------------|--------|
| `false`       | `true`, `false` | Boolean |

- **Description:** Enable access to the Docker Desktop integrated terminal. If
the value is set to `false`, users can't use the Docker terminal to interact
with the host machine and execute commands directly from Docker Desktop.
- **OS:** {{< badge color=blue text="All" >}}
- **Use case:** Allow or restrict developer access to the built-in terminal.
- **Configure this setting with:**
    - **General** setting in [Docker Desktop GUI](/manuals/desktop/settings-and-maintenance/settings.md)
    - Settings Management: `desktopTerminalEnabled` setting in the [`admin-settings.json` file](/manuals/enterprise/security/hardened-desktop/settings-management/configure-json-file.md)

### Enable Docker Debug by default

| Default value | Accepted values | Format   |
|---------------|-----------------|----------|
| `false`       | `true`, `false` | Boolean  |

- **Description:** Enable debug logging by default for Docker CLI commands.
- **OS:** {{< badge color=blue text="All" >}}
- **Use case:** Assist with debugging support issues.
- **Configure this setting with:**
    - **General** settings in [Docker Desktop GUI](/manuals/desktop/settings-and-maintenance/settings.md)

### Include VM in Time Machine backup

| Default value | Accepted values | Format |
|---------------|-----------------|--------|
| `false`       | `true`, `false` | Boolean |

- **Description:** Back up the Docker Desktop virtual machine.
- **OS:** {{< badge color=blue text="Mac only" >}}
- **Use case:** Manage persistence of application data.
- **Configure this setting with:**
    - **General** settings in [Docker Desktop GUI](/manuals/desktop/settings-and-maintenance/settings.md)

### Use containerd for pulling and storing images

| Default value | Accepted values | Format   |
|---------------|-----------------|----------|
| `false`       | `true`, `false` | Boolean  |

- **Description:** Use containerd native snapshotter instead of legacy
snapshotters.
- **OS:** {{< badge color=blue text="All" >}}
- **Use case:** Improve image handling performance and compatibility.
- **Configure this setting with:**
    - **General** settings in [Docker Desktop GUI](/manuals/desktop/settings-and-maintenance/settings.md)

### Choose Virtual Machine Manager

#### Docker VMM

| Default value | Accepted values | Format   |
|---------------|-----------------|----------|
| `true`        | `true`, `false` | Boolean  |

#### Apple Virtualization framework

| Default value | Accepted values | Format   |
|---------------|-----------------|----------|
| `true`        | `true`, `false` | Boolean  |

- **Description:** Use Apple Virtualization Framework to run Docker containers.
- **OS:** {{< badge color=blue text="Mac only" >}}
- **Use case:** Improve VM performance on Apple Silicon.
- **Configure this setting with:**
    - **General** settings in [Docker Desktop GUI](/manuals/desktop/settings-and-maintenance/settings.md)

#### Rosetta

| Default value | Accepted values | Format   |
|---------------|-----------------|----------|
| `true`        | `true`, `false` | Boolean  |

- **Description:** Use Rosetta to emulate `amd64` on Apple Silicon. If value
is set to `true`, Docker Desktop turns on Rosetta to accelerate
x86_64/amd64 binary emulation on Apple Silicon.
- **OS:** {{< badge color=blue text="Mac only" >}} 13+
- **Use case:** Run Intel-based containers on Apple Silicon hosts.

> [!NOTE]
>
> In hardened environments, disable and lock this setting so only ARM-native
images are permitted.

- **Configure this setting with:**
    - **General** settings in [Docker Desktop GUI](/manuals/desktop/settings-and-maintenance/settings.md)
    - Settings Management:`useVirtualizationFrameworkRosetta` setting in the [`admin-settings.json` file](/manuals/enterprise/security/hardened-desktop/settings-management/configure-json-file.md)
    - Settings Management: **Use Rosetta for x86_64/amd64 emulation on Apple Silicon** setting in the [Admin Console](/manuals/enterprise/security/hardened-desktop/settings-management/configure-admin-console.md)

> [!NOTE]
>
> Rosetta requires enabling Apple Virtualization framework.

#### QEMU

> [!WARNING]
>
> QEMU has been deprecated in Docker Desktop versions 4.44 and later. For more information, see the [blog announcement](https://www.docker.com/blog/docker-desktop-for-mac-qemu-virtualization-option-to-be-deprecated-in-90-days/) 

| Default value | Accepted values | Format   |
|---------------|-----------------|----------|
| `true`        | `true`, `false` | Boolean  |

### Choose file sharing implementation

#### VirtioFS

| Default value | Accepted values | Format   |
|---------------|-----------------|----------|
| `true`        | `true`, `false` | Boolean  |

- **Description:** Use VirtioFS for fast, native file sharing between host and
containers. If value is set to `true`, VirtioFS is set as the file sharing
mechanism. If both VirtioFS and gRPC are set to `true`, VirtioFS takes
precedence.
- **OS:** {{< badge color=blue text="Mac only" >}} 12.5+
- **Use case:** Improve volume mount performance and compatibility.

> [!NOTE]
>
> In hardened environments, enable and lock this setting for macOS 12.5 and
later.

- **Configure this setting with:**
    - **General settings** in [Docker Desktop GUI](/manuals/desktop/settings-and-maintenance/settings.md)
    - Settings Management: `useVirtualizationFrameworkVirtioFS` setting in the [`admin-settings.json` file](/manuals/enterprise/security/hardened-desktop/settings-management/configure-json-file.md)
    - Settings Management: **Use VirtioFS for file sharing** setting in the [Admin Console](/manuals/enterprise/security/hardened-desktop/settings-management/configure-admin-console.md)

#### gRPC FUSE

| Default value | Accepted values | Format   |
|---------------|-----------------|----------|
| `true`        | `true`, `false` | Boolean  |

- **Description:** Enable gRPC FUSE for macOS file sharing. If value is set to
`true`, gRPC Fuse is set as the file sharing mechanism.
- **OS:** {{< badge color=blue text="Mac only" >}}
- **Use case:** Improve performance and compatibility of file mounts.

> [!NOTE]
>
> In hardened environments, disable and lock this setting.

- **Configure this setting with:**
    - **General** settings in [Docker Desktop GUI](/manuals/desktop/settings-and-maintenance/settings.md)
    - Settings Management: `useGrpcfuse` setting in the [`admin-settings.json` file](/manuals/enterprise/security/hardened-desktop/settings-management/configure-json-file.md)
    - Settings Management: **Use gRPC FUSE for file sharing** setting in the [Admin Console](/manuals/enterprise/security/hardened-desktop/settings-management/configure-admin-console.md)

#### osxfs

| Default value | Accepted values | Format  |
| ------------- | --------------- | ------- |
| `false`       | `true`, `false` | Boolean |

- **Description:** Enable the legacy osxfs file sharing driver for macOS. When
set to true, Docker Desktop uses osxfs instead of VirtioFS or gRPC FUSE to mount
host directories into containers.
- **OS:** {{< badge color=blue text="Mac only" >}}
- **Use case:** Use the original file sharing implementation when compatibility
with older tooling or specific workflows is required.
- **Configure this setting with:**
    - **General** settings in [Docker Desktop GUI](/manuals/desktop/settings-and-maintenance/settings.md)

### Send usage statistics

| Default value | Accepted values | Format |
|---------------|-----------------|--------|
| `true`        | `true`, `false` | Boolean |

- **Description:** Controls whether Docker Desktop collects and sends local
usage statistics and crash reports to Docker. This setting affects telemetry
gathered from the Docker Desktop application itself. It does not affect
server-side telemetry collected via Docker Hub or other backend services, such
as login timestamps, pulls, or builds.
- **OS:** {{< badge color=blue text="All" >}}
- **Use case:** Enable analytics to help Docker improve the product based on
usage data.

> [!NOTE]
>
> In hardened environments, disable and lock this setting. This allows you
to control all your data flows and collect support logs via secure channels
if needed.

- **Configure this setting with:**
    - **General** settings in [Docker Desktop GUI](/manuals/desktop/settings-and-maintenance/settings.md)
    - Settings Management: `analyticsEnabled` setting in the [`admin-settings.json` file](/manuals/enterprise/security/hardened-desktop/settings-management/configure-json-file.md)
    - Settings Management: **Send usage statistics** setting in the [Admin Console](/manuals/enterprise/security/hardened-desktop/settings-management/configure-admin-console.md)

> [!NOTE]
>
> Organizations using the Insights Dashboard may need this setting enabled to
ensure that developer activity is fully visible. If users opt out and the
setting is not locked, their activity may be excluded from analytics
views.

### Use Enhanced Container Isolation

| Default value | Accepted values | Format   |
|---------------|-----------------|----------|
| `false`       | `true`, `false` | Boolean  |

- **Description:** Enable Enhanced Container Isolation for secure container
execution.
- **OS:** {{< badge color=blue text="All" >}}
- **Use case:** Prevent containers from modifying configuration or sensitive
host areas.

> [!NOTE]
>
> In hardened environments, disable and lock this setting.

- **Configure this setting with:**
    - **General settings** in [Docker Desktop GUI](/manuals/desktop/settings-and-maintenance/settings.md)
    - Settings Management: `enhancedContainerIsolation` setting in the [`admin-settings.json` file](/manuals/enterprise/security/hardened-desktop/settings-management/configure-json-file.md)
    - Settings Management: **Enable enhanced container isolation** setting in the [Admin Console](/manuals/enterprise/security/hardened-desktop/settings-management/configure-admin-console.md)

### Show CLI hints

| Default value | Accepted values | Format   |
|---------------|-----------------|----------|
| `true`       | `true`, `false` | Boolean  |

- **Description:** Display helpful CLI tips in the terminal when using Docker commands.
- **OS:** {{< badge color=blue text="All" >}}
- **Use case:** Help users discover and learn Docker CLI features through inline suggestions.
- **Configure this setting with:**
    - **General** settings in [Docker Desktop GUI](/manuals/desktop/settings-and-maintenance/settings.md)

### Enable Scout image analysis

| Default value | Accepted values | Format   |
|---------------|-----------------|----------|
| `true`        | `true`, `false` | Boolean  |

- **Description:** Enable Docker Scout to generate and display SBOM data for container images.
- **OS:** {{< badge color=blue text="All" >}}
- **Use case:** Turn on Docker Scout analysis features to view vulnerabilities, packages, and metadata associated with images.

> [!NOTE]
>
> In hardened environments, enable and lock this setting to ensure SBOMs are
always built to satisfy compliance scans.

- **Configure this setting with:**
    - **General settings** in [Docker Desktop GUI](/manuals/desktop/settings-and-maintenance/settings.md)
    - Settings Management: `sbomIndexing` setting in the [`admin-settings.json` file](/manuals/enterprise/security/hardened-desktop/settings-management/configure-json-file.md)
    - Settings Management: **SBOM indexing** settings in the [Admin Console](/manuals/enterprise/security/hardened-desktop/settings-management/configure-admin-console.md)

### Enable background Scout SBOM indexing

| Default value | Accepted values | Format   |
|---------------|-----------------|----------|
| `false`        | `true`, `false` | Boolean  |

- **Description:** Automatically index SBOM data for images in the background without requiring user interaction.
- **OS:** {{< badge color=blue text="All" >}}
- **Use case:** Keep image metadata up to date by allowing Docker to perform SBOM indexing during idle time or after image pull operations.

> [!NOTE]
>
> In hardened environments, enable and lock this setting.

- **Configure this setting with:**
    - **General settings** in [Docker Desktop GUI](/manuals/desktop/settings-and-maintenance/settings.md)

### Automatically check configuration

| Default value         | Accepted values | Format  |
|-----------------------|-----------------|---------|
| `CurrentSettingsVersions` | Integer         | Integer |

- **Description:** Regularly checks your configuration to ensure no unexpected changes have been made by another application
- **OS:** {{< badge color=blue text="All" >}}
- **Use case:** Track versions for compatibility
- **Configure this setting with:**
    - **General** settings in [Docker Desktop GUI](/manuals/desktop/settings-and-maintenance/settings.md)
    - Settings Management: `configurationFileVersion` setting in the [`admin-settings.json` file](/manuals/enterprise/security/hardened-desktop/settings-management/configure-json-file.md)

## Resources

### CPU limit

| Default value                                 | Accepted values | Format  |
|-----------------------------------------------|-----------------|---------|
| Number of logical CPU cores available on host | Integer         | Integer |

- **Description:** Number of CPUs assigned to the Docker Desktop virtual machine.
- **OS:** {{< badge color=blue text="All" >}}
- **Use case:** Resource allocation control.
- **Configure this setting with:**
    - **Advanced** Resources settings in [Docker Desktop GUI](/manuals/desktop/settings-and-maintenance/settings.md)

### Memory limit

| Default value              | Accepted values | Format  |
|---------------------------|-----------------|---------|
| Based on system resources | Integer         | Integer |

- **Description:** Amount of RAM (in MiB) assigned to the Docker virtual machine.
- **OS:** {{< badge color=blue text="All" >}}
- **Use case:** Control how much memory Docker can use on the host.
- **Configure this setting with:**
    - **Advanced** Resources settings in [Docker Desktop GUI](/manuals/desktop/settings-and-maintenance/settings.md)

### Swap

| Default value | Accepted values | Format  |
|---------------|-----------------|---------|
| `1024`        | Integer         | Integer |

- **Description:** Amount of swap space (in MiB) assigned to the Docker virtual machine
- **OS:** {{< badge color=blue text="All" >}}
- **Use case:** Extend memory availability via swap
- **Configure this setting with:**
    - **Advanced** Resources settings in [Docker Desktop GUI](/manuals/desktop/settings-and-maintenance/settings.md)

### Disk usage limit

| Default value                  | Accepted values | Format  |
|-------------------------------|-----------------|---------|
| Default disk size of machine. | Integer         | Integer |

- **Description:** Maximum disk size (in MiB) allocated for Docker Desktop.
- **OS:** {{< badge color=blue text="All" >}}
- **Use case:** Constrain Docker's virtual disk size for storage management.
- **Configure this setting with:**
    - **Advanced** Resources settings in [Docker Desktop GUI](/manuals/desktop/settings-and-maintenance/settings.md)

### Disk image location

| Default value                                                                 | Accepted values | Format |
|--------------------------------------------------|-----------------|--------|
| macOS: `~/Library/Containers/com.docker.docker/Data/vms/0`  <br> Windows: `%USERPROFILE%\AppData\Local\Docker\wsl\data` | File path       | String |

- **Description:** Path where Docker Desktop stores virtual machine data.
- **OS:** {{< badge color=blue text="All" >}}
- **Use case:** Redirect Docker data to a custom location.
- **Configure this setting with:**
    - **Advanced** Resources settings in [Docker Desktop GUI](/manuals/desktop/settings-and-maintenance/settings.md)

### Enable Resource Saver

| Default value | Accepted values | Format   |
|---------------|-----------------|----------|
| `true`        | `true`, `false` | Boolean  |

- **Description:** Enable Docker Desktop to pause when idle.
- **OS:** {{< badge color=blue text="All" >}}
- **Use case:** Save system resources during periods of inactivity.
- **Configure this setting with:**
    - **Advanced** Resources settings in [Docker Desktop GUI](/manuals/desktop/settings-and-maintenance/settings.md)

### File sharing directories

| Default value                           | Accepted values                 | Format                  |
|----------------------------------------|---------------------------------|--------------------------|
| Varies by OS                           | List of file paths as strings   | Array list of strings   |

- **Description:** List of allowed directories shared between the host and
containers. When a path is added, its subdirectories are allowed.
- **OS:** {{< badge color=blue text="All" >}}
- **Use case:** Restrict or define what file paths are available to containers.

> [!NOTE]
>
> In hardened environments, lock to an explicit whitelist and disable end-user
edits.

- **Configure this setting with:**
    - **File sharing** Resources settings in [Docker Desktop GUI](/manuals/desktop/settings-and-maintenance/settings.md)
    - Settings Management: `filesharingAllowedDirectories` setting in the [`admin-settings.json` file](/manuals/enterprise/security/hardened-desktop/settings-management/configure-json-file.md)
    - Settings Management: **Allowed file sharing directories** setting in the [Admin Console](/manuals/enterprise/security/hardened-desktop/settings-management/configure-admin-console.md)

### Proxy exclude

| Default value | Accepted values    | Format |
|---------------|--------------------|--------|
| `""`          | List of addresses  | String |

- **Description:** Configure addresses that containers should bypass from proxy
settings.
- **OS:** {{< badge color=blue text="All" >}}
- **Use case:** Fine-tune proxy exceptions for container networking.

> [!NOTE]
>
> In hardened environments, disable and lock this setting.

- **Configure this setting with:**
    - **Proxies** Resources settings in [Docker Desktop GUI](/manuals/desktop/settings-and-maintenance/settings.md)
    - Settings Management: `proxy` setting with `manual` and `exclude` modes in the [`admin-settings.json` file](/manuals/enterprise/security/hardened-desktop/settings-management/configure-json-file.md)

### Docker subnet

| Default value     | Accepted values | Format |
|-------------------|-----------------|--------|
| `192.168.65.0/24` | IP address      | String |

- **Description:** Overrides the network range used for vpnkit DHCP/DNS for
`*.docker.internal`.
- **OS:** {{< badge color=blue text="Mac only" >}}
- **Use case:** Customize the subnet used for Docker container networking.
- **Configure this setting with:**
    - Settings Management: `vpnkitCIDR` setting in the [`admin-settings.json` file](/manuals/enterprise/security/hardened-desktop/settings-management/configure-json-file.md)
    - Settings Management: **VPN Kit CIDR** setting in the [Admin Console](/manuals/enterprise/security/hardened-desktop/settings-management/configure-admin-console.md)

### Use kernel networking for UDP

| Default value | Accepted values | Format   |
|---------------|-----------------|----------|
| `false`       | `true`, `false` | Boolean  |

- **Description:** Use the host’s kernel network stack for UDP traffic instead of Docker’s virtual network driver. This enables faster and more direct UDP communication, but may bypass some container isolation features.
- **OS:** {{< badge color=blue text="All" >}}
- **Use case:** Improve performance or compatibility for workloads that rely heavily on UDP traffic, such as real-time media, DNS, or game servers.
- **Configure this setting with:**
    - **Network** Resources settings in [Docker Desktop GUI](/manuals/desktop/settings-and-maintenance/settings.md)

### Enable host networking

| Default value | Accepted values | Format   |
|---------------|-----------------|----------|
| `false`       | `true`, `false` | Boolean  |

- **Description:** Enable experimental host networking support.
- **OS:** {{< badge color=blue text="All" >}}
- **Use case:** Allow containers to use the host network stack.
- **Configure this setting with:**
    - **Network** Resources settings in [Docker Desktop GUI](/manuals/desktop/settings-and-maintenance/settings.md)

### Networking mode

| Default value | Accepted values | Format   |
|---------------|-----------------|----------|
| `dual-stack` | `ipv4only`, `ipv6only` | String  |

- **Description:** Set the networking mode.
- **OS:** {{< badge color=blue text="Windows and Mac" >}}
- **Use case:** Choose the default IP protocol used when Docker creates new networks.
- **Configure this setting with:**
    - **Network** Resources settings in [Docker Desktop GUI](/manuals/desktop/settings-and-maintenance/settings.md)
    - Settings Management: `defaultNetworkingMode` setting in the [`admin-settings.json` file](/manuals/enterprise/security/hardened-desktop/settings-management/configure-json-file.md)

For more information, see [Networking](/manuals/desktop/features/networking.md#networking-mode-and-dns-behaviour-for-mac-and-windows).

#### Inhibit DNS resolution for IPv4/IPv6

| Default value | Accepted values | Format   |
|---------------|-----------------|----------|
| `auto` | `ipv4`, `ipv6`, `none` | String  |

- **Description:** Filters unsupported DNS record types. Requires Docker Desktop
version 4.43 and up.
- **OS:** {{< badge color=blue text="Windows and Mac" >}}
- **Use case:** Control how Docker filters DNS records returned to containers, improving reliability in environments where only IPv4 or IPv6 is supported.
- **Configure this setting with:**
    - **Network** Resources settings in [Docker Desktop GUI](/manuals/desktop/settings-and-maintenance/settings.md)
    - Settings Management: `dnsInhibition` setting in the [`admin-settings.json` file](/manuals/enterprise/security/hardened-desktop/settings-management/configure-json-file.md)

For more information, see [Networking](/manuals/desktop/features/networking.md#networking-mode-and-dns-behaviour-for-mac-and-windows).

### Enable WSL engine

| Default value | Accepted values | Format   |
|---------------|-----------------|----------|
| `true`        | `true`, `false` | Boolean  |

- **Description:** If the value is set to `true`, Docker Desktop uses the WSL2
based engine. This overrides anything that may have been set at installation
using the `--backend=<backend name>` flag.
- **OS:** {{< badge color=blue text="Windows only" >}} + WSL
- **Use case:** Enable Linux containers via WSL 2 backend.

> [!NOTE]
>
> In hardened environments, enable and lock this setting.

- **Configure this setting with:**
    - **WSL Integration** Resources settings in [Docker Desktop GUI](/manuals/desktop/settings-and-maintenance/settings.md)
    - Settings Management: `wslEngineEnabled` setting in the [`admin-settings.json` file](/manuals/enterprise/security/hardened-desktop/settings-management/configure-json-file.md)
    - Settings Management: **Windows Subsystem for Linux (WSL) Engine** setting in the [Admin Console](/manuals/enterprise/security/hardened-desktop/settings-management/configure-admin-console.md)

## Docker Engine

The Docker Engine settings let you configure low-level daemon settings through a raw JSON object. These settings are passed directly to the dockerd process that powers container management in Docker Desktop.

| Key                   | Example                     | Description                                        | Accepted values / Format       | Default |
| --------------------- | --------------------------- | -------------------------------------------------- | ------------------------------ | ------- |
| `debug`               | `true`                      | Enable verbose logging in the Docker daemon        | Boolean                        | `false` |
| `experimental`        | `true`                      | Enable experimental Docker CLI and daemon features | Boolean                        | `false` |
| `insecure-registries` | `["myregistry.local:5000"]` | Allow pulling from HTTP registries without TLS     | Array of strings (`host:port`) | `[]`    |
| `registry-mirrors`    | `["https://mirror.gcr.io"]` | Define alternative registry endpoints              | Array of URLs                  | `[]`    |

- **Description:** Customize the behavior of the Docker daemon using a structured JSON config passed directly to dockerd.
- **OS:** {{< badge color=blue text="All" >}}
- **Use case:** Fine-tune registry access, enable debug mode, or opt into experimental features.
- **Configure this setting with:**
    - **Docker Engine** settings in [Docker Desktop GUI](/manuals/desktop/settings-and-maintenance/settings.md)

> [!NOTE]
>
> Values for this setting are passed as-is to the Docker daemon. Invalid or unsupported fields may prevent Docker Desktop from starting.

## Builders

Builders settings lets you manage Buildx builder instances for advanced image-building scenarios, including multi-platform builds and custom backends.

| Key         | Example                          | Description                                                                | Accepted values / Format  | Default   |
| ----------- | -------------------------------- | -------------------------------------------------------------------------- | ------------------------- | --------- |
| `name`      | `"my-builder"`                   | Name of the builder instance                                               | String                    | —         |
| `driver`    | `"docker-container"`             | Backend used by the builder (`docker`, `docker-container`, `remote`, etc.) | String                    | `docker`  |
| `platforms` | `["linux/amd64", "linux/arm64"]` | Target platforms supported by the builder                                  | Array of platform strings | Host arch |

- **Description:** Configure custom Buildx builders for Docker Desktop, including driver type and supported platforms.
- **OS:** {{< badge color=blue text="All" >}}
- **Use case:** Set up advanced build configurations like cross-platform images or remote builders.
- **Configure this setting with:**
    - **Builders** settings in [Docker Desktop GUI](/manuals/desktop/settings-and-maintenance/settings.md)

> [!NOTE]
>
> Builder definitions are structured as an array of objects, each describing a builder instance. Conflicting or unsupported configurations may cause build errors.

## Kubernetes

### Enable Kubernetes

| Default value | Accepted values | Format   |
|---------------|-----------------|----------|
| `false`       | `true`, `false` | Boolean  |

- **Description:** Enable the integrated Kubernetes cluster in Docker Desktop.
- **OS:** {{< badge color=blue text="All" >}}
- **Use case:** Enable or disable Kubernetes support for developers.

> [!NOTE]
>
> In hardened environments, disable and lock this setting.

> [!IMPORTANT]
>
> When Kubernetes is enabled through Settings Management policies, only the
`kubeadm` cluster provisioning method is supported. The `kind` provisioning
method is not yet supported by Settings Management.

- **Configure this setting with:**
    - **Kubernetes** settings in [Docker Desktop GUI](/manuals/desktop/settings-and-maintenance/settings.md)
    - Settings Management: `kubernetes` setting in the [`admin-settings.json` file](/manuals/enterprise/security/hardened-desktop/settings-management/configure-json-file.md)
    - Settings Management: **Allow Kubernetes** setting in the [Admin Console](/manuals/enterprise/security/hardened-desktop/settings-management/configure-admin-console.md)

### Choose cluster provisioning method

| Default value | Accepted values | Format |
|---------------|-----------------|--------|
| `kubeadm`     | `kubeadm`, `kind`  | String |

- **Description:** Set the Kubernetes node mode (single-node or multi-node).
- **OS:** {{< badge color=blue text="All" >}}
- **Use case:** Control the topology of the integrated Kubernetes cluster.
- **Configure this setting with:**
    - **Kubernetes** settings in [Docker Desktop GUI](/manuals/desktop/settings-and-maintenance/settings.md)

### Kubernetes node count (kind provisioning)

| Default value | Accepted values | Format  |
|---------------|-----------------|---------|
| `1`           | Integer         | Integer |

- **Description:** Number of nodes to create in a multi-node Kubernetes cluster.
- **OS:** {{< badge color=blue text="All" >}}
- **Use case:** Scale the number of Kubernetes nodes for development or testing.
- **Configure this setting with:**
    - **Kubernetes** settings in [Docker Desktop GUI](/manuals/desktop/settings-and-maintenance/settings.md)

### Kubernetes node version (kind provisioning)

| Default value | Accepted values               | Format |
|---------------|-------------------------------|--------|
| `1.31.1`      | Semantic version (e.g., 1.29.1) | String |

- **Description:** Version of Kubernetes used for cluster node creation.
- **OS:** {{< badge color=blue text="All" >}}
- **Use case:** Pin a specific Kubernetes version for consistency or
compatibility.
- **Configure this setting with:**
    - **Kubernetes** settings in [Docker Desktop GUI](/manuals/desktop/settings-and-maintenance/settings.md)

### Show system containers

| Default value | Accepted values | Format   |
|---------------|-----------------|----------|
| `false`       | `true`, `false` | Boolean  |

- **Description:** Show Kubernetes system containers in the Docker Dashboard container list
- **OS:** {{< badge color=blue text="All" >}}
- **Use case:** Allow developers to view kube-system containers for debugging

> [!NOTE]
>
> In hardened environments, disable and lock this setting.

- **Configure this setting with:**
    - **Kubernetes** settings in [Docker Desktop GUI](/manuals/desktop/settings-and-maintenance/settings.md)

### Custom Kubernetes image repository

| Default value | Accepted values | Format   |
|---------------|-----------------|----------|
| `""`          | Registry URL    | String   |

- **Description**: Configure a custom image repository for Kubernetes control
plane images. This allows Docker Desktop to pull Kubernetes system
images from a private registry or mirror instead of Docker Hub. This setting
overrides the `[registry[:port]/][namespace]` portion of image names.
- **OS**: {{< badge color=blue text="All" >}}
- **Use case**: Use private registries in air-gapped environments or
when Docker Hub access is restricted.

> [!NOTE]
>
> The images must be cloned/mirrored from Docker Hub with matching tags. The
specific images required depend on the cluster provisioning method (`kubeadm`
or `kind`). See the Kubernetes documentation for the complete list
of required images and detailed setup instructions.

- **Configure this setting with**:
    - Settings Management: `KubernetesImagesRepository` settings in the
    [`admin-settings.json` file](/manuals/enterprise/security/hardened-desktop/settings-management/configure-json-file.md)
    - Settings Management: **Kubernetes Images Repository** setting in the
    [Admin Console](/manuals/enterprise/security/hardened-desktop/settings-management/configure-admin-console.md)

> [!IMPORTANT]
>
> When using `KubernetesImagesRepository` with Enhanced Container Isolation (ECI)
enabled, you must add the following images to the ECI Docker socket mount image
list: `[imagesRepository]/desktop-cloud-provider-kind:*` and
`[imagesRepository]/desktop-containerd-registry-mirror:*`.

## Software updates

### Automatically check for updates

| Default value | Accepted values | Format |
|---------------|-----------------|--------|
| `false`       | `true`, `false` | Boolean |

- **Description:** Disable automatic update polling for Docker Desktop. If the
value is set to `true`, checking for updates and notifications about Docker
Desktop updates are disabled.
- **OS:** {{< badge color=blue text="All" >}}
- **Use case:** Freeze the current version in enterprise environments.

> [!NOTE]
>
> In hardened environments, enable this setting and lock. This guarantees that
only internally vetted versions are installed.

- **Configure this setting with:**
    - Settings Management: `disableUpdate` setting in the [`admin-settings.json` file](/manuals/enterprise/security/hardened-desktop/settings-management/configure-json-file.md)
    - Settings Management: **Disable update** setting in the [Admin Console](/manuals/enterprise/security/hardened-desktop/settings-management/configure-admin-console.md)

### Always download updates

| Default value | Accepted values | Format |
|---------------|-----------------|--------|
| `false`       | `true`, `false` | Boolean |

- **Description:** Automatically download Docker Desktop updates when available.
- **OS:** {{< badge color=blue text="All" >}}
- **Use case:** Manage auto update behavior.
- **Configure this setting with:**
    - **Software updates** settings in [Docker Desktop GUI](/manuals/desktop/settings-and-maintenance/settings.md)
    - Settings Management: **Disable updates** setting in the [Admin Console](/manuals/enterprise/security/hardened-desktop/settings-management/configure-admin-console.md)

## Extensions

### Enable Docker extensions

| Default value | Accepted values | Format   |
|---------------|-----------------|----------|
| `true`        | `true`, `false` | Boolean  |

- **Description:** Enable or disable Docker Extensions.
- **OS:** {{< badge color=blue text="All" >}}
- **Use case:** Control access to the Extensions Marketplace and installed
extensions.

> [!NOTE]
>
> In hardened environments, disable and lock this setting. This prevents
third-party or unvetted plugins from being installed.

- **Configure this setting with:**
    - **Extensions** settings in [Docker Desktop GUI](/manuals/desktop/settings-and-maintenance/settings.md)
    - Settings Management: `extensionsEnabled` setting in the [`admin-settings.json` file](/manuals/enterprise/security/hardened-desktop/settings-management/configure-json-file.md)
    - Settings Management: **Allow Extensions** setting in the [Admin Console](/manuals/enterprise/security/hardened-desktop/settings-management/configure-admin-console.md)

### Allow only extensions distributed through the Docker Marketplace

| Default value | Accepted values | Format   |
|---------------|-----------------|----------|
| `true`        | `true`, `false` | Boolean  |

- **Description:** Restrict Docker Desktop to only run Marketplace extensions.
- **OS:** {{< badge color=blue text="All" >}}
- **Use case:** Prevent running third-party or local extensions.
- **Configure this setting with:**
    - **Extensions** settings in [Docker Desktop GUI](/manuals/desktop/settings-and-maintenance/settings.md)

### Show Docker Extensions system containers

| Default value | Accepted values | Format   |
|---------------|-----------------|----------|
| `false`       | `true`, `false` | Boolean  |

- **Description:** Show system containers used by Docker Extensions in the container list
- **OS:** {{< badge color=blue text="All" >}}
- **Use case:** Help developers troubleshoot or view extension system containers
- **Configure this setting with:**
    - **Extensions** settings in [Docker Desktop GUI](/manuals/desktop/settings-and-maintenance/settings.md)

## Beta features

> [!IMPORTANT]
>
> For Docker Desktop versions 4.41 and earlier, these settings lived under the **Experimental features** tab on the **Features in development** page.

### Enable Docker AI

| Default value | Accepted values | Format   |
|---------------|-----------------|----------|
| `false`       | `true`, `false` | Boolean  |

- **Description:** Enable Docker AI features in the Docker Desktop experience.
- **OS:** {{< badge color=blue text="All" >}}
- **Use case:** Enable or disable AI features like "Ask Gordon".
- **Configure this setting with:**
    - **Beta** settings in [Docker Desktop GUI](/manuals/desktop/settings-and-maintenance/settings.md)
    - Settings Management: `enableDockerAI` setting in the [`admin-settings.json` file](/manuals/enterprise/security/hardened-desktop/settings-management/configure-json-file.md)

### Enable Docker Model Runner

| Default value | Accepted values | Format   |
|---------------|-----------------|----------|
| `true`       | `true`, `false` | Boolean  |

- **Description:** Enable Docker Model Runner features in Docker Desktop.
- **OS:** {{< badge color=blue text="All" >}}
- **Use case:** Enable or disable Docker Model Runner features.
- **Configure this setting with:**
    - **Beta** settings in [Docker Desktop GUI](/manuals/desktop/settings-and-maintenance/settings.md)
    - Settings Management: `enableDockerAI` setting in the [`admin-settings.json` file](/manuals/enterprise/security/hardened-desktop/settings-management/configure-json-file.md)

#### Enable host-side TCP support

| Default value | Accepted values | Format   |
|---------------|-----------------|----------|
| `false`       | `true`, `false` | Boolean  |

- **Description:** Enable Docker Model Runner features in Docker Desktop.
- **OS:** {{< badge color=blue text="All" >}}
- **Use case:** Enable or disable Docker Model Runner features.
- **Configure this setting with:**
    - **Beta** settings in [Docker Desktop GUI](/manuals/desktop/settings-and-maintenance/settings.md)
    - Settings Management: `enableDockerAI` setting in the [`admin-settings.json` file](/manuals/enterprise/security/hardened-desktop/settings-management/configure-json-file.md)
    
> [!NOTE]
>
> This setting requires Docker Model Runner setting to be enabled first.

##### Port

| Default value | Accepted values | Format  |
|---------------|-----------------|---------|
| 12434         | Integer         | Integer |

- **Description:** Specifies the exposed TCP port.
- **OS:** {{< badge color=blue text="All" >}}
- **Use case:** Connect to the Model Runner via TCP.
- **Configure this setting with:**
    - **Beta features** settings in [Docker Desktop GUI](/manuals/desktop/settings-and-maintenance/settings.md)
    - Settings Management: `enableInferenceTCP` setting in the [`admin-settings.json` file](/manuals/enterprise/security/hardened-desktop/settings-management/configure-json-file.md)

##### CORS Allowed Origins

| Default value | Accepted values                                                                 | Format |
|---------------|---------------------------------------------------------------------------------|--------|
| Empty string  | Empty string to deny all,`*` to accept all, or a list of comma-separated values | String |

- **Description:** Specifies the allowed CORS origins.
- **OS:** {{< badge color=blue text="All" >}}
- **Use case:** Integration with a web app.
- **Configure this setting with:**
    - **Beta features** settings in [Docker Desktop GUI](/manuals/desktop/settings-and-maintenance/settings.md)
    - Settings Management: `enableInferenceCORS` setting in the [`admin-settings.json` file](/manuals/enterprise/security/hardened-desktop/settings-management/configure-json-file.md)

### Enable Docker MCP Toolkit

| Default value | Accepted values | Format   |
|---------------|-----------------|----------|
| `true`       | `true`, `false` | Boolean  |

- **Description:** Enable [Docker MCP Toolkit](/manuals/ai/mcp-catalog-and-toolkit/_index.md) in Docker Desktop.
- **OS:** {{< badge color=blue text="All" >}}
- **Configure this setting with:**
    - **Beta** settings in [Docker Desktop GUI](/manuals/desktop/settings-and-maintenance/settings.md)
    - Settings Management: `enableDockerMCPToolkit` setting in the [`admin-settings.json` file](/manuals/enterprise/security/hardened-desktop/settings-management/configure-json-file.md)
    

### Enable Wasm

| Default value | Accepted values | Format   |
|---------------|-----------------|----------|
| `true`       | `true`, `false` | Boolean  |

- **Description:** Enable [Wasm](/manuals/desktop/features/wasm.md) to run Wasm workloads.
- **OS:** {{< badge color=blue text="All" >}}
- **Configure this setting with:**
    - **Beta** settings in [Docker Desktop GUI](/manuals/desktop/settings-and-maintenance/settings.md)    

### Enable Compose Bridge

| Default value | Accepted values | Format   |
|---------------|-----------------|----------|
| `true`       | `true`, `false` | Boolean  |

- **Description:** Enable [Compose Bridge](/manuals/compose/bridge/_index.md).
- **OS:** {{< badge color=blue text="All" >}}
- **Configure this setting with:**
    - **Beta** settings in [Docker Desktop GUI](/manuals/desktop/settings-and-maintenance/settings.md)

## Notifications

### Status updates on tasks and processes

| Default value | Accepted values | Format   |
|---------------|-----------------|----------|
| `true`        | `true`, `false` | Boolean  |

- **Description:** Display general informational messages inside Docker Desktop
- **OS:** {{< badge color=blue text="All" >}}
- **Use case:** Customize in-app communication visibility
- **Configure this setting with:**
    - **Notifications** settings in [Docker Desktop GUI](/manuals/desktop/settings-and-maintenance/settings.md)

### Recommendations from Docker

| Default value | Accepted values | Format   |
|---------------|-----------------|----------|
| `true`        | `true`, `false` | Boolean  |

- **Description:** Display promotional announcements and banners inside Docker Desktop
- **OS:** {{< badge color=blue text="All" >}}
- **Use case:** Control exposure to Docker news and feature promotion
- **Configure this setting with:**
    - **Notifications** settings in [Docker Desktop GUI](/manuals/desktop/settings-and-maintenance/settings.md)

### Docker announcements

| Default value | Accepted values | Format   |
|---------------|-----------------|----------|
| `true`        | `true`, `false` | Boolean  |

- **Description:** Display general announcements inside Docker Desktop.
- **OS:** {{< badge color=blue text="All" >}}
- **Use case:** Enable or suppress Docker-wide announcements in the GUI.
- **Configure this setting with:**
    - **Notifications** settings in [Docker Desktop GUI](/manuals/desktop/settings-and-maintenance/settings.md)

### Docker surveys

| Default value | Accepted values | Format   |
|---------------|-----------------|----------|
| `true`        | `true`, `false` | Boolean  |

- **Description:** Display notifications inviting users to participate in surveys
- **OS:** {{< badge color=blue text="All" >}}
- **Use case:** Enable or disable in-product survey prompts
- **Configure this setting with:**
    - **Notifications** settings in [Docker Desktop GUI](/manuals/desktop/settings-and-maintenance/settings.md)

### Docker Scout Notification pop-ups

| Default value | Accepted values | Format   |
|---------------|-----------------|----------|
| `true`        | `true`, `false` | Boolean  |

- **Description:** Enable Docker Scout popups inside Docker Desktop.
- **OS:** {{< badge color=blue text="All" >}}
- **Use case:** Show or hide vulnerability scan notifications
- **Configure this setting with:**
    - **Notifications** settings in [Docker Desktop GUI](/manuals/desktop/settings-and-maintenance/settings.md)

### Docker Scout OS notifications

| Default value | Accepted values | Format   |
|---------------|-----------------|----------|
| `false`       | `true`, `false` | Boolean  |

- **Description:** Enable Docker Scout notifications through the operating system.
- **OS:** {{< badge color=blue text="All" >}}
- **Use case:** Push Scout updates via system notification center
- **Configure this setting with:**
    - **Notifications** settings in [Docker Desktop GUI](/manuals/desktop/settings-and-maintenance/settings.md)

## Advanced

### Configure installation of Docker CLI

| Default value | Accepted values | Format   |
|---------------|-----------------|----------|
| `system`      | File path       | String   |

- **Description:** Install location for Docker CLI binaries.
- **OS:** {{< badge color=blue text="All" >}}
- **Use case:** Customize CLI install location for compliance or tooling.
- **Configure this setting with:**
    - **Advanced** settings in [Docker Desktop GUI](/manuals/desktop/settings-and-maintenance/settings.md)

### Allow the default Docker socket to be used

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
    - **Advanced** settings in [Docker Desktop GUI](/manuals/desktop/settings-and-maintenance/settings.md)
    - Settings Management: `dockerSocketMount` setting in the [`admin-settings.json` file](/manuals/enterprise/security/hardened-desktop/settings-management/configure-json-file.md)

### Allow privileged port mapping

| Default value | Accepted values | Format   |
|---------------|-----------------|----------|
| `true`        | `true`, `false` | Boolean  |

- **Description:** Starts the privileged helper process which binds privileged ports that are between 1 and 1024
- **OS:** {{< badge color=blue text="Mac only" >}}
- **Use case:** Enforce elevated privileges for networking support
- **Configure this setting with:**
    - **Advanced** settings in [Docker Desktop GUI](/manuals/desktop/settings-and-maintenance/settings.md)

## Settings not available in the Docker Desktop GUI

The following settings aren’t shown in the Docker Desktop GUI. You can only configure them using Settings Management with the Admin Console or the `admin-settings.json` file.

### Block `docker load`

| Default value | Accepted values | Format |
|---------------|-----------------|--------|
| `false`       | `true`, `false` | Boolean |

- **Description:** Prevent users from loading local Docker images using the `docker load` command.
- **OS:** {{< badge color=blue text="All" >}}
- **Use case:** Enforce image provenance by restricting local image imports.

> [!NOTE]
>
> In hardened environments, enable and lock this setting. This forces all images
to come from your secure, scanned registry.

- **Configure this setting with:**
    - Settings Management: `blockDockerLoad` setting in the [`admin-settings.json` file](/manuals/enterprise/security/hardened-desktop/settings-management/configure-json-file.md)

### Expose Docker API on TCP 2375

| Default value | Accepted values | Format |
|---------------|-----------------|--------|
| `false`       | `true`, `false` | Boolean |

- **Description:** Exposes the Docker API over an unauthenticated TCP socket on port 2375. Only recommended for isolated and protected environments.
- **OS:** {{< badge color=blue text="Windows only" >}}
- **Use case:** Required for legacy integrations or environments without named pipe support.

> [!NOTE]
>
> In hardened environments, disable and lock this setting. This ensures the
Docker API is only reachable via the secure internal socket.

- **Configure this setting with:**
    - Settings Management: `exposeDockerAPIOnTCP2375` in the [`admin-settings.json` file](/manuals/enterprise/security/hardened-desktop/settings-management/configure-json-file.md)

### Air-gapped container proxy

| Default value | Accepted values | Format      |
| ------------- | --------------- | ----------- |
| See example   | Object          | JSON object |

- **Description:** Configure a manual HTTP/HTTPS proxy for containers. Useful in air-gapped environments where containers need restricted access.
- **OS:** {{< badge color=blue text="All" >}}
- **Use case:** Redirect or block container networking to comply with offline or secured network environments.
- **Configure this setting with:**
    - Settings Management: `containersProxy` setting in the [`admin-settings.json` file](/manuals/enterprise/security/hardened-desktop/settings-management/configure-json-file.md)

#### Example

```json
"containersProxy": {
  "locked": true,
  "mode": "manual",
  "http": "",
  "https": "",
  "exclude": [],
  "pac": "",
  "transparentPorts": ""
}
```

Docker socket access control (ECI exceptions)

| Default value | Accepted values | Format      |
| ------------- | --------------- | ----------- |
| -           | Object          | JSON object |

- **Description:** Allow specific images or commands to use the Docker socket when Enhanced Container Isolation is enabled.
- **OS:** {{< badge color=blue text="All" >}}
- **Use case:** Support tools like Testcontainers or LocalStack that need Docker socket access while maintaining secure defaults.
- Configure this setting with:
    - Settings Management: `enhancedContainerIsolation` > `dockerSocketMount` in the [`admin-settings.json` file](/manuals/enterprise/security/hardened-desktop/settings-management/configure-json-file.md)

#### Example

```json
"enhancedContainerIsolation": {
  "locked": true,
  "value": true,
  "dockerSocketMount": {
    "imageList": {
      "images": [
        "docker.io/localstack/localstack:*",
        "docker.io/testcontainers/ryuk:*"
      ]
    },
    "commandList": {
      "type": "deny",
      "commands": ["push"]
    }
  }
}
```

### Allow beta features

| Default value | Accepted values | Format   |
|---------------|-----------------|----------|
| `false`       | `true`, `false` | Boolean  |

- **Description:** Enable access to beta features in Docker Desktop.
- **OS:** {{< badge color=blue text="All" >}}
- **Use case:** Give developers early access to features that are in public beta.

> [!NOTE]
>
> In hardened environments, disable and lock this setting.

- **Configure this setting with:**
    - Settings Management: `allowBetaFeatures` setting in the [`admin-settings.json` file](/manuals/enterprise/security/hardened-desktop/settings-management/configure-json-file.md)

### Docker daemon options (Linux or Windows)

| Default value | Accepted values | Format   |
|---------------|-----------------|----------|
| `{}`          | JSON object     | Stringified JSON |

- **Description:** Override the Docker daemon configuration used in Linux or Windows containers.
- **OS:** {{< badge color=blue text="All" >}}
- **Use case:** Configure low-level Docker daemon options (e.g., logging, storage drivers) without editing the local config files.

> [!NOTE]
>
> In hardened environments, provide a vetted JSON config and lock it so no
overrides are possible.

- **Configure this setting with:**
    - Settings Management: `linuxVM.dockerDaemonOptions` or `windowsContainers.dockerDaemonOptions` in the [`admin-settings.json` file](/manuals/enterprise/security/hardened-desktop/settings-management/configure-json-file.md)

### VPNKit CIDR

| Default value     | Accepted values | Format |
|-------------------|-----------------|--------|
| `192.168.65.0/24` | CIDR notation   | String |

- **Description:** Set the subnet used for internal VPNKit DHCP/DNS services.
- **OS:** {{< badge color=blue text="Mac only" >}}
- **Use case:** Prevent IP conflicts in environments with overlapping subnets.

> [!NOTE]
>
> In hardened environments, lock to an approved, non-conflicting CIDR.

- **Configure this setting with:**
    - Settings Management: `vpnkitCIDR` setting in the [`admin-settings.json` file](/manuals/enterprise/security/hardened-desktop/settings-management/configure-json-file.md)
    - Settings Management: **VPN Kit CIDR** setting in the [Admin Console](/manuals/enterprise/security/hardened-desktop/settings-management/configure-admin-console.md)

### Enable Kerberos and NTLM authentication

| Default value | Accepted values | Format |
|---------------|-----------------|--------|
| `false`       | `true`, `false` | Boolean |

- **Description:** Enables Kerberos and NTLM proxy authentication for enterprise environments.
- **OS:** {{< badge color=blue text="All" >}}
- **Use case:** Allow users to authenticate with enterprise proxy servers that require Kerberos or NTLM.
- **Configure this setting with:**
    - Settings Management: `proxy.enableKerberosNtlm` in the [`admin-settings.json` file](/manuals/enterprise/security/hardened-desktop/settings-management/configure-json-file.md)
