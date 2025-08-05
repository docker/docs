---
description: How to configure Settings Management for Docker Desktop
keywords: admin, controls, rootless, enhanced container isolation
title: Configure Settings Management with a JSON file
linkTitle: Use a JSON file
weight: 10
aliases:
 - /desktop/hardened-desktop/settings-management/configure/
 - /security/for-admins/hardened-desktop/settings-management/configure/
 - /security/for-admins/hardened-desktop/settings-management/configure-json-file/
---

{{< summary-bar feature_name="Hardened Docker Desktop" >}}

This page explains how to use an `admin-settings.json` file to configure and
enforce Docker Desktop settings. Use this method to standardize Docker
Desktop environments in your organization.

## Prerequisites

- [Enforce sign-in](/manuals/enterprise/security/enforce-sign-in/_index.md) to
ensure all users authenticate with your organization.
- A Docker Business subscription is required.

Docker Desktop only applies settings from the `admin-settings.json` file if both
authentication and Docker Business license checks succeed.

> [!IMPORTANT]
>
> If a user isn't signed in or isn't part of a Docker Business organization,
the settings file is ignored.

## Limitation

- The `admin-settings.json` file doesn't work in air-gapped or offline
environments.
- The file is not compatible with environments that restrict authentication
with Docker Hub.

## Step one: Create the settings file

You can:

- Use the `--admin-settings` installer flag to auto-generate the file. See:
    - [macOS](/manuals/desktop/setup/install/mac-install.md#install-from-the-command-line) install guide
    - [Windows](/manuals/desktop/setup/install/windows-install.md#install-from-the-command-line) install guide
- Or create it manually and place it in the following locations:
    - Mac: `/Library/Application\ Support/com.docker.docker/admin-settings.json`
    - Windows: `C:\ProgramData\DockerDesktop\admin-settings.json`
    - Linux: `/usr/share/docker-desktop/admin-settings.json`

> [!IMPORTANT]
>
> Place the file in a protected directory to prevent modification. Use MDM tools
like [Jamf](https://www.jamf.com/lp/en-gb/apple-mobile-device-management-mdm-jamf-shared/?attr=google_ads-brand-search-shared&gclid=CjwKCAjw1ICZBhAzEiwAFfvFhEXjayUAi8FHHv1JJitFPb47C_q_RCySTmF86twF1qJc_6GST-YDmhoCuJsQAvD_BwE) to distribute it at scale.

## Step two: Define settings

> [!TIP]
>
> For a complete list of available settings, their supported platforms, and which configuration methods they work with, see the [Settings reference](settings-reference.md).

The `admin-settings.json` file uses structured keys to define what can
be configured and whether the values are enforced.

Each setting supports the `locked` field. When `locked` is set to `true`, users
can't change that value in Docker Desktop, the CLI, or config files. When
`locked` is set to `false`, the value acts like a default suggestion and users
can still update it.

Settings where `locked` is set to `false` are ignored on existing installs if
a user has already customized that value in `settings-store.json`,
`settings.json`, or `daemon.json`.

> [!NOTE]
>
> Some settings are platform-specific or require a minimum Docker Desktop
version. See the [Settings reference](/manuals/enterprise/security/hardened-desktop/settings-management/settings-reference.md) for details.

### Example settings file

The following file is an example `admin-settings.json` file. For a full list
of configurable settings for the `admin-settings.json` file, see [`admin-settings.json` configurations](#admin-settingsjson-configurations).

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
    "enableKerberosNtlm": false
  },
  "containersProxy": {
    "locked": true,
    "mode": "manual",
    "http": "",
    "https": "",
    "exclude": [],
    "pac":"",
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
  }
}
```

## Step three: Restart and apply settings

Settings apply after Docker Desktop is restarted and the user is signed in.

- New installs: Launch Docker Desktop and sign in.
- Existing installs: Quit Docker Desktop fully and relaunch it.

> [!IMPORTANT]
>
> Restarting Docker Desktop from the menu isn't enough. It must be fully
quit and reopened.

## `admin-settings.json` configurations

### General

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

### Proxy

|Parameter|OS|Description|Version|
|:-------------------------------|---|:-------------------------------|---|
|`proxy`|   |If `mode` is set to `system` instead of `manual`, Docker Desktop gets the proxy values from the system and ignores and values set for `http`, `https` and `exclude`. Change `mode` to `manual` to manually configure proxy servers. If the proxy port is custom, specify it in the `http` or `https` property, for example `"https": "http://myotherproxy.com:4321"`. The `exclude` property specifies a comma-separated list of hosts and domains to bypass the proxy. |  |
|&nbsp; &nbsp; &nbsp; &nbsp;`windowsDockerdPort`| Windows only | Exposes Docker Desktop's internal proxy locally on this port for the Windows Docker daemon to connect to. If it is set to 0, a random free port is chosen. If the value is greater than 0, use that exact value for the port. The default value is -1 which disables the option. |  |
|&nbsp; &nbsp; &nbsp; &nbsp;`enableKerberosNtlm`|  |When set to `true`, Kerberos and NTLM authentication is enabled. Default is `false`. For more information, see the settings documentation. | Docker Desktop version 4.32 and later. |

### Container proxy

|Parameter|OS|Description|Version|
|:-------------------------------|---|:-------------------------------|---|
|`containersProxy` | | Creates air-gapped containers. For more information see [Air-Gapped Containers](../air-gapped-containers.md).| Docker Desktop version 4.29 and later. |

### Linux VM

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

### Kubernetes

|Parameter|OS|Description|Version|
|:-------------------------------|---|:-------------------------------|---|
|`kubernetes`|  | If `enabled` is set to true, a Kubernetes single-node cluster is started when Docker Desktop starts. If `showSystemContainers` is set to true, Kubernetes containers are displayed in the Docker Desktop Dashboard and when you run `docker ps`. The [imagesRepository](../../../../desktop/features/kubernetes.md#configuring-a-custom-image-registry-for-kubernetes-control-plane-images) setting lets you specify which repository Docker Desktop pulls control-plane Kubernetes images from. |  |

> [!NOTE]
>
> When using the `imagesRepository` setting and Enhanced Container Isolation (ECI), add the following images to the [ECI Docker socket mount image list](#enhanced-container-isolation):
>
> * [imagesRepository]/desktop-cloud-provider-kind:*
> * [imagesRepository]/desktop-containerd-registry-mirror:*
>
> These containers mount the Docker socket, so you must add the images to the ECI images list. If not, ECI will block the mount and Kubernetes won't start.

### Networking

|Parameter|OS|Description|Version|
|:-------------------------------|---|:-------------------------------|---|
| `defaultNetworkingMode` | Windows and Mac only | Defines the default IP protocol for new Docker networks: `dual-stack` (IPv4 + IPv6, default), `ipv4only`, or `ipv6only`. | Docker Desktop version 4.43 and later. |
| `dnsInhibition` | Windows and Mac only | Controls DNS record filtering returned to containers. Options: `auto` (recommended), `ipv4`, `ipv6`, `none`| Docker Desktop version 4.43 and later. |

For more information, see [Networking](/manuals/desktop/features/networking.md#networking-mode-and-dns-behaviour-for-mac-and-windows).

### Beta features

> [!IMPORTANT]
>
> For Docker Desktop versions 4.41 and earlier, some of these settings lived under the **Experimental features** tab on the **Features in development** page.

| Parameter                                            | OS | Description                                                                                                                                                                                                                                               | Version                                 |
|:-----------------------------------------------------|----|:----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|-----------------------------------------|
| `allowBetaFeatures`                                  |    | If `value` is set to `true`, beta features are enabled.                                                                                                                                                                                                   |                                         |
| `enableDockerAI`                                     |    | If `allowBetaFeatures` is true, setting `enableDockerAI` to `true` enables [Docker AI (Ask Gordon)](/manuals/ai/gordon/_index.md) by default. You can independently control this setting from the `allowBetaFeatures` setting.                            |                                         |
| `enableInference`                                    |    | If `allowBetaFeatures` is true, setting `enableInference` to `true` enables [Docker Model Runner](/manuals/ai/model-runner/_index.md) by default. You can independently control this setting from the `allowBetaFeatures` setting.                        |                                         |
| &nbsp; &nbsp; &nbsp; &nbsp; `enableInferenceTCP`     |    | Enable host-side TCP support. This setting requires Docker Model Runner setting to be enabled first.                                                                                                                                                      |                                         |
| &nbsp; &nbsp; &nbsp; &nbsp; `enableInferenceTCPPort` |    | Specifies the exposed TCP port. This setting requires Docker Model Runner setting to be enabled first.                                                                                                                                                    |                                         |
| &nbsp; &nbsp; &nbsp; &nbsp; `enableInferenceCORS`    |    | Specifies the allowed CORS origins. Empty string to deny all,`*` to accept all, or a list of comma-separated values. This setting requires Docker Model Runner setting to be enabled first.                                                                                                                                                    |                                         |
| `enableDockerMCPToolkit`                             |    | If `allowBetaFeatures` is true, setting `enableDockerMCPToolkit` to `true` enables the [MCP Toolkit feature](/manuals/ai/mcp-catalog-and-toolkit/toolkit.md) by default. You can independently control this setting from the `allowBetaFeatures` setting. |                                         |
| `allowExperimentalFeatures`                          |    | If `value` is set to `true`, experimental features are enabled.                                                                                                                                                                                           | Docker Desktop version 4.41 and earlier |

### Enhanced Container Isolation

|Parameter|OS|Description|Version|
|:-------------------------------|---|:-------------------------------|---|
|`enhancedContainerIsolation`|  | If `value` is set to true, Docker Desktop runs all containers as unprivileged, via the Linux user-namespace, prevents them from modifying sensitive configurations inside the Docker Desktop VM, and uses other advanced techniques to isolate them. For more information, see [Enhanced Container Isolation](../enhanced-container-isolation/_index.md).|  |
| &nbsp; &nbsp; &nbsp; &nbsp;`dockerSocketMount` |  | By default, enhanced container isolation blocks bind-mounting the Docker Engine socket into containers (e.g., `docker run -v /var/run/docker.sock:/var/run/docker.sock ...`). This lets you relax this in a controlled way. See [ECI Configuration](../enhanced-container-isolation/config.md) for more info. |  |
| &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; `imageList` |  | Indicates which container images are allowed to bind-mount the Docker Engine socket. |  |
| &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; `commandList` |  | Restricts the commands that containers can issue via the bind-mounted Docker Engine socket. |  |
