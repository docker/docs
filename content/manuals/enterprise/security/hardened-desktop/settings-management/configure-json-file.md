---
title: Configure Settings Management with a JSON file
linkTitle: Use a JSON file
description: Configure and enforce Docker Desktop settings using an admin-settings.json file
keywords: admin controls, settings management, configuration, enterprise, docker desktop, json file
weight: 10
aliases:
 - /desktop/hardened-desktop/settings-management/configure/
 - /security/for-admins/hardened-desktop/settings-management/configure/
 - /security/for-admins/hardened-desktop/settings-management/configure-json-file/
---

{{< summary-bar feature_name="Hardened Docker Desktop" >}}

Settings Management lets you configure and enforce Docker Desktop settings across your organization using an `admin-settings.json` file. This standardizes Docker Desktop environments and ensures consistent configurations for all users.

## Prerequisites

Before you begin, make sure you have:

- [Enforce sign-in](/manuals/enterprise/security/enforce-sign-in/_index.md) for
your organization
- A Docker Business subscription

Docker Desktop only applies settings from the `admin-settings.json` file when both authentication and Docker Business license checks succeed.

> [!IMPORTANT]
>
> Users must be signed in and part of a Docker Business organization. If either condition isn't met, the settings file is ignored.

## Step one: Create the settings file

You can create the `admin-settings.json` file in two ways:

- Use the `--admin-settings` installer flag to auto-generate the file:
    - [macOS](/manuals/desktop/setup/install/mac-install.md#install-from-the-command-line) installation guide
    - [Windows](/manuals/desktop/setup/install/windows-install.md#install-from-the-command-line) installation guide
- Create it manually and place it in the following locations:
    - Mac: `/Library/Application\ Support/com.docker.docker/admin-settings.json`
    - Windows: `C:\ProgramData\DockerDesktop\admin-settings.json`
    - Linux: `/usr/share/docker-desktop/admin-settings.json`

> [!IMPORTANT]
>
> Place the file in a protected directory to prevent unauthorized changes. Use Mobile Device Management (MDM) tools like Jamf to distribute the file at scale across your organization.

## Step two: Configure settings

> [!TIP]
>
> For a complete list of available settings, their supported platforms, and which configuration methods they work with, see the [Settings reference](settings-reference.md).

The `admin-settings.json` file uses structured keys to define configurable settings and whether values are enforced.

Each setting supports a `locked` field that controls user permissions:

- When `locked` is set to `true`, users can't change that value in Docker Desktop, the CLI, or config files.
- When `locked` is set to `false`, the value acts like a default suggestion and users
can still update it.

Settings where `locked` is set to `false` are ignored on existing installs if
a user has already customized that value in `settings-store.json`,
`settings.json`, or `daemon.json`.

### Grouped settings

Docker Desktop groups some settings together with a single toggle that controls
the entire section. These include:

- Enhanced Container Isolation (ECI): Uses a main toggle (`enhancedContainerIsolation`) that enables/disables the entire feature, with sub-settings for specific configurations
- Kubernetes: Uses a main toggle (`kubernetes.enabled`) with sub-settings for cluster configuration
- Docker Scout: Groups settings under the `scout` object

When configuring grouped settings:

1. Set the main toggle to enable the feature
1. Configure sub-settings within that group
1. When you lock the main toggle, users cannot modify any settings in that group

Example for `enhancedContainerIsolation`:

```json
"enhancedContainerIsolation": {
  "locked": true,  // This locks the entire ECI section
  "value": true,   // This enables ECI
  "dockerSocketMount": {  // These are sub-settings
    "imageList": {
      "images": ["docker.io/testcontainers/ryuk:*"]
    }
  }
}
```

### Example `admin-settings.json` file

The following sample is an `admin-settings.json` file with common enterprise settings configured. You can use this example as a template with the [`admin-settings.json` configurations](#admin-settingsjson-configurations):

```json {collapse=true}
{
  "configurationFileVersion": 2,
  "exposeDockerAPIOnTCP2375": {
    "locked": true,
    "value": false
  },
  "proxy": {
    "locked": true,
    "mode": "system",
    "http": "",
    "https": "",
    "exclude": [],
    "windowsDockerdPort": 65000,
    "enableKerberosNtlm": false,
    "pac": "",
    "embeddedPac": ""
  },
  "containersProxy": {
    "locked": true,
    "mode": "manual",
    "http": "",
    "https": "",
    "exclude": [],
    "pac":"",
    "embeddedPac": "",
    "transparentPorts": ""
  },
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
  },
  "linuxVM": {
    "wslEngineEnabled": {
      "locked": false,
      "value": false
    },
    "dockerDaemonOptions": {
      "locked": false,
      "value":"{\"debug\": false}"
    },
    "vpnkitCIDR": {
      "locked": false,
      "value":"192.168.65.0/24"
    }
  },
  "kubernetes": {
     "locked": false,
     "enabled": false,
     "showSystemContainers": false,
     "imagesRepository": ""
  },
  "windowsContainers": {
    "dockerDaemonOptions": {
      "locked": false,
      "value":"{\"debug\": false}"
    }
  },
  "disableUpdate": {
    "locked": false,
    "value": false
  },
  "analyticsEnabled": {
    "locked": false,
    "value": true
  },
  "extensionsEnabled": {
    "locked": true,
    "value": false
  },
  "scout": {
    "locked": false,
    "sbomIndexing": true,
    "useBackgroundIndexing": true
  },
  "allowBetaFeatures": {
    "locked": false,
    "value": false
  },
  "blockDockerLoad": {
    "locked": false,
    "value": true
  },
  "filesharingAllowedDirectories": [
    {
      "path": "$HOME",
      "sharedByDefault": true
    },
    {
      "path":"$TMP",
      "sharedByDefault": false
    }
  ],
  "useVirtualizationFrameworkVirtioFS": {
    "locked": true,
    "value": true
  },
  "useVirtualizationFrameworkRosetta": {
    "locked": true,
    "value": true
  },
  "useGrpcfuse": {
    "locked": true,
    "value": true
  },
  "displayedOnboarding": {
    "locked": true,
    "value": true
  },
  "desktopTerminalEnabled": {
    "locked": false,
    "value": false
  },
  "enableInference": {
    "locked": false,
    "value": true
  },
  "enableInferenceTCP": {
    "locked": false,
    "value": true
  },
  "enableInferenceTCPPort": {
    "locked": true,
    "value": 12434
  },
  "enableInferenceCORS": {
    "locked": true,
    "value": ""
  },
  "enableInferenceGPUVariant": {
    "locked": true,
    "value": true
  }
}
```

## Step three: Apply the settings

Settings take effect after Docker Desktop restarts and the user signs in.

For new installations:

1. Launch Docker Desktop.
1. Sign in with your Docker account.

For existing installations:

1. Quit Docker Desktop completely.
1. Relaunch Docker Desktop.

> [!IMPORTANT]
>
> You must fully quit and reopen Docker Desktop. Restarting from the menu isn't sufficient.

## `admin-settings.json` configurations

The following tables describe all available settings in the `admin-settings.json` file.

> [!NOTE]
>
> Some settings are platform-specific or require minimum Docker Desktop versions. Check the Version column for requirements.

### General settings

|Parameter|OS|Description|Version|
|:-------------------------------|---|:-------------------------------|---|
|`configurationFileVersion`|   |Specifies the version of the configuration file format.|   |
|`analyticsEnabled`|  |If `value` is set to false, Docker Desktop doesn't send usage statistics to Docker. |  |
|`disableUpdate`|  |If `value` is set to true, checking for and notifications about Docker Desktop updates is disabled.|  |
|`extensionsEnabled`|  |If `value` is set to false, Docker extensions are disabled. |  |
| `blockDockerLoad` | | If `value` is set to `true`, users are no longer able to run [`docker load`](/reference/cli/docker/image/load/) and receive an error if they try to.|  |
| `displayedOnboarding` |  | If `value` is set to `true`, the onboarding survey will not be displayed to new users. Setting `value` to `false` has no effect. |  Docker Desktop version 4.30 and later |
| `desktopTerminalEnabled` |  | If `value` is set to `false`, developers cannot use the Docker terminal to interact with the host machine and execute commands directly from Docker Desktop. |  |
|`exposeDockerAPIOnTCP2375`| Windows only| Exposes the Docker API on a specified port. If `value` is set to true, the Docker API is exposed on port 2375. Note: This is unauthenticated and should only be enabled if protected by suitable firewall rules.|  |
| `silentModulesUpdate` | | If `value` is set to `true`, Docker Desktop automatically updates components that don't require a restart. For example, the Docker CLI or Docker Scout components. | Docker Desktop version 4.46 and later. |

### File sharing and emulation

|Parameter|OS|Description|Version|
|:-------------------------------|---|:-------------------------------|---|
| `filesharingAllowedDirectories` |  | Specify which paths your developers can add file shares to. Also accepts `$HOME`, `$TMP`, or `$TEMP` as `path` variables. When a path is added, its subdirectories are allowed. If `sharedByDefault` is set to `true`, that path will be added upon factory reset or when Docker Desktop first starts. |  |
| `useVirtualizationFrameworkVirtioFS`|  macOS only | If `value` is set to `true`, VirtioFS is set as the file sharing mechanism. Note: If both `useVirtualizationFrameworkVirtioFS` and `useGrpcfuse` have `value` set to `true`, VirtioFS takes precedence. Likewise, if both `useVirtualizationFrameworkVirtioFS` and `useGrpcfuse` have `value` set to `false`, osxfs is set as the file sharing mechanism. |  |
| `useGrpcfuse` | macOS only | If `value` is set to `true`, gRPC Fuse is set as the file sharing mechanism. |  |
| `useVirtualizationFrameworkRosetta`|  macOS only | If `value` is set to `true`, Docker Desktop turns on Rosetta to accelerate x86_64/amd64 binary emulation on Apple Silicon. Note: This also automatically enables `Use Virtualization framework`. | Docker Desktop version 4.29 and later. |

### Docker Scout

|Parameter|OS|Description|Version|
|:-------------------------------|---|:-------------------------------|---|
|`scout`| | Setting `useBackgroundIndexing` to `false` disables automatic indexing of images loaded to the image store. Setting `sbomIndexing` to `false` prevents users from being able to index image by inspecting them in Docker Desktop or using `docker scout` CLI commands. |  |

### Proxy settings

|Parameter|OS|Description|Version|
|:-------------------------------|---|:-------------------------------|---|
|`proxy`|   |If `mode` is set to `system` instead of `manual`, Docker Desktop gets the proxy values from the system and ignores and values set for `http`, `https` and `exclude`. Change `mode` to `manual` to manually configure proxy servers. If the proxy port is custom, specify it in the `http` or `https` property, for example `"https": "http://myotherproxy.com:4321"`. The `exclude` property specifies a comma-separated list of hosts and domains to bypass the proxy. |  |
| `windowsDockerdPort`| Windows only | Exposes Docker Desktop's internal proxy locally on this port for the Windows Docker daemon to connect to. If it is set to 0, a random free port is chosen. If the value is greater than 0, use that exact value for the port. The default value is -1 which disables the option. |  |
|`enableKerberosNtlm`|  |When set to `true`, Kerberos and NTLM authentication is enabled. Default is `false`. For more information, see the settings documentation. | Docker Desktop version 4.32 and later. |
| `pac` | | Specifies a PAC file URL. For example, `"pac": "http://proxy/proxy.pac"`. | |
| `embeddedPac`  | | Specifies an embedded PAC (Proxy Auto-Config) script. For example, `"embeddedPac": "function FindProxyForURL(url, host) { return \"DIRECT\"; }"`. This setting takes precedence over HTTP, HTTPS, Proxy bypass and PAC server URL. |  Docker Desktop version 4.46 and later. |

### Container proxy

|Parameter|OS|Description|Version|
|:-------------------------------|---|:-------------------------------|---|
|`containersProxy` | | Creates air-gapped containers. For more information see [Air-Gapped Containers](../air-gapped-containers.md).| Docker Desktop version 4.29 and later. |
| `pac` | | Specifies a PAC file URL. For example, `"pac": "http://containerproxy/proxy.pac"`. | |
| `embeddedPac`  | | Specifies an embedded PAC (Proxy Auto-Config) script. For example, `"embeddedPac": "function FindProxyForURL(url, host) { return \"PROXY 192.168.92.1:2003\"; }"`. This setting takes precedence over HTTP, HTTPS, Proxy bypass and PAC server URL. |  Docker Desktop version 4.46 and later. |

### Linux VM settings

|Parameter|OS|Description|Version|
|:-------------------------------|---|:-------------------------------|---|
| `linuxVM` |   |Parameters and settings related to Linux VM options - grouped together here for convenience. |  |
| &nbsp; &nbsp; &nbsp; &nbsp;`wslEngineEnabled`  | Windows only | If `value` is set to true, Docker Desktop uses the WSL 2 based engine. This overrides anything that may have been set at installation using the `--backend=<backend name>` flag. |  |
| &nbsp; &nbsp; &nbsp; &nbsp;`dockerDaemonOptions` |  |If `value` is set to true, it overrides the options in the Docker Engine config file. See the [Docker Engine reference](/reference/cli/dockerd/#daemon-configuration-file). Note that for added security, a few of the config attributes may be overridden when Enhanced Container Isolation is enabled. |  |
| &nbsp; &nbsp; &nbsp; &nbsp;`vpnkitCIDR` |  |Overrides the network range used for vpnkit DHCP/DNS for `*.docker.internal`  |  |

### Windows containers

|Parameter|OS|Description|Version|
|:-------------------------------|---|:-------------------------------|---|
| `windowsContainers` |  | Parameters and settings related to `windowsContainers` options - grouped together here for convenience.  |  |
| &nbsp; &nbsp; &nbsp; &nbsp;`dockerDaemonOptions` |  | Overrides the options in the Linux daemon config file. See the [Docker Engine reference](/reference/cli/dockerd/#daemon-configuration-file).|  |

> [!NOTE]
>
> This setting is not available to configure via the Docker Admin Console.

### Kubernetes settings

|Parameter|OS|Description|Version|
|:-------------------------------|---|:-------------------------------|---|
|`kubernetes`|  | If `enabled` is set to true, a Kubernetes single-node cluster is started when Docker Desktop starts. If `showSystemContainers` is set to true, Kubernetes containers are displayed in the Docker Desktop Dashboard and when you run `docker ps`. The [imagesRepository](../../../../desktop/features/kubernetes.md#configuring-a-custom-image-registry-for-kubernetes-control-plane-images) setting lets you specify which repository Docker Desktop pulls control-plane Kubernetes images from. |  |

> [!NOTE]
>
> When using `imagesRepository` with Enhanced Container Isolation (ECI), add these images to the [ECI Docker socket mount image list](#enhanced-container-isolation):
>
> `[imagesRepository]/desktop-cloud-provider-kind:`
> `[imagesRepository]/desktop-containerd-registry-mirror:`
>
> These containers mount the Docker socket, so you must add them to the ECI images list. Otherwise, ECI blocks the mount and Kubernetes won't start.

### Networking settings

|Parameter|OS|Description|Version|
|:-------------------------------|---|:-------------------------------|---|
| `defaultNetworkingMode` | Windows and Mac only | Defines the default IP protocol for new Docker networks: `dual-stack` (IPv4 + IPv6, default), `ipv4only`, or `ipv6only`. | Docker Desktop version 4.43 and later. |
| `dnsInhibition` | Windows and Mac only | Controls DNS record filtering returned to containers. Options: `auto` (recommended), `ipv4`, `ipv6`, `none`| Docker Desktop version 4.43 and later. |

For more information, see [Networking](/manuals/desktop/features/networking.md#networking-mode-and-dns-behaviour-for-mac-and-windows).

### AI settings

| Parameter                   | OS            | Description                                                                                                                                                                                                                         | Version |
|:----------------------------|---------------|:------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|---------|
| `enableInference`           |               | Setting `enableInference` to `true` enables [Docker Model Runner](/manuals/ai/model-runner/_index.md).                                                                                                                              |         |
| `enableInferenceTCP`        |               | Enable host-side TCP support. This setting requires the Docker Model Runner setting to be enabled first.                                                                                                                                |         |
| `enableInferenceTCPPort`    |               | Specifies the exposed TCP port. This setting requires the Docker Model Runner and Enable host-side TCP support settings to be enabled first.                                                                                            |         |
| `enableInferenceCORS`       |               | Specifies the allowed CORS origins. Empty string to deny all,`*` to accept all, or a list of comma-separated values. This setting requires the Docker Model Runner and Enable host-side TCP support settings to be enabled first.       |         |
| `enableInferenceGPUVariant` | Windows only  | Setting `enableInferenceGPUVariant` to `true` enables GPU-backed inference. The additional components required for this don't come by default with Docker Desktop, therefore they will be downloaded to `~/.docker/bin/inference`.  |         |

### Beta features

> [!IMPORTANT]
>
> For Docker Desktop versions 4.41 and earlier, some of these settings lived under the **Experimental features** tab on the **Features in development** page.

| Parameter                                            | OS | Description                                                                                                                                                                                                                                               | Version                                 |
|:-----------------------------------------------------|----|:----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|-----------------------------------------|
| `allowBetaFeatures`                                  |    | If `value` is set to `true`, beta features are enabled.                                                                                                                                                                                                   |                                         |
| `enableDockerAI`                                     |    | If `allowBetaFeatures` is true, setting `enableDockerAI` to `true` enables [Docker AI (Ask Gordon)](/manuals/ai/gordon/_index.md) by default. You can independently control this setting from the `allowBetaFeatures` setting.                            |                                         |
| `enableDockerMCPToolkit`                             |    | If `allowBetaFeatures` is true, setting `enableDockerMCPToolkit` to `true` enables the [MCP Toolkit feature](/manuals/ai/mcp-catalog-and-toolkit/toolkit.md) by default. You can independently control this setting from the `allowBetaFeatures` setting. |                                         |
| `allowExperimentalFeatures`                          |    | If `value` is set to `true`, experimental features are enabled.                                                                                                                                                                                           | Docker Desktop version 4.41 and earlier |

### Enhanced Container Isolation

|Parameter|OS|Description|Version|
|:-------------------------------|---|:-------------------------------|---|
|`enhancedContainerIsolation`|  | If `value` is set to true, Docker Desktop runs all containers as unprivileged, via the Linux user-namespace, prevents them from modifying sensitive configurations inside the Docker Desktop VM, and uses other advanced techniques to isolate them. For more information, see [Enhanced Container Isolation](../enhanced-container-isolation/_index.md).|  |
| &nbsp; &nbsp; &nbsp; &nbsp;`dockerSocketMount` |  | By default, enhanced container isolation blocks bind-mounting the Docker Engine socket into containers (e.g., `docker run -v /var/run/docker.sock:/var/run/docker.sock ...`). This lets you relax this in a controlled way. See [ECI Configuration](../enhanced-container-isolation/config.md) for more info. |  |
| &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; `imageList` |  | Indicates which container images are allowed to bind-mount the Docker Engine socket. |  |
| &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; `commandList` |  | Restricts the commands that containers can issue via the bind-mounted Docker Engine socket. |  |
