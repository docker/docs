---
description: How to configure Settings Management for Docker Desktop
keywords: admin, controls, rootless, enhanced container isolation
title: Configure Settings Management with a JSON file
linkTitle: Use a JSON file
weight: 10
aliases: 
 - /desktop/hardened-desktop/settings-management/configure/
 - /security/for-admins/hardened-desktop/settings-management/configure/
---

{{< summary-bar feature_name="Hardened Docker Desktop" >}}

This page contains information on how to configure Settings Management with an `admin-settings.json` file. You can specify and lock configuration parameters to create a standardized Docker Desktop environment across your company or organization.

Settings Management is designed specifically for organizations who don’t give developers root access to their machines.

## Prerequisites

You first need to [enforce sign-in](/manuals/security/for-admins/enforce-sign-in/_index.md) to ensure that all Docker Desktop developers authenticate with your organization. Since Settings Management requires a Docker Business subscription, enforced sign-in guarantees that only authenticated users have access and that the feature consistently takes effect across all users, even though it may still work without enforced sign-in.

## Step one: Create the `admin-settings.json` file and save it in the correct location

You can either use the `--admin-settings` installer flag on [macOS](/manuals/desktop/setup/install/mac-install.md#install-from-the-command-line) or [Windows](/manuals/desktop/setup/install/windows-install.md#install-from-the-command-line) to automatically create the `admin-settings.json` and save it in the correct location, or set it up manually.

To set it up manually:
1. Create a new, empty JSON file and name it `admin-settings.json`.
2. Save the `admin-settings.json` file on your developers' machines in the following locations:
   - Mac: `/Library/Application\ Support/com.docker.docker/admin-settings.json`
   - Windows: `C:\ProgramData\DockerDesktop\admin-settings.json`
   - Linux: `/usr/share/docker-desktop/admin-settings.json`

   By placing this file in a protected directory, developers are unable to modify it.

   > [!IMPORTANT]
   >
   > It is assumed that you have the ability to push the `admin-settings.json` settings file to the locations specified through a device management software such as [Jamf](https://www.jamf.com/lp/en-gb/apple-mobile-device-management-mdm-jamf-shared/?attr=google_ads-brand-search-shared&gclid=CjwKCAjw1ICZBhAzEiwAFfvFhEXjayUAi8FHHv1JJitFPb47C_q_RCySTmF86twF1qJc_6GST-YDmhoCuJsQAvD_BwE).

## Step two: Configure the settings you want to lock in

> [!NOTE]
>
> Some of the configuration parameters only apply to certain platforms or to specific Docker Desktop versions. This is highlighted in the following table.

The `admin-settings.json` file requires a nested list of configuration parameters, each of which must contain the `locked` parameter. You can add or remove configuration parameters as per your requirements.

If `locked: true`, users aren't able to edit this setting from Docker Desktop or the CLI.

If `locked: false`, it's similar to setting a factory default in that:
   - For new installs, `locked: false` pre-populates the relevant settings in the Docker Desktop Dashboard, but users are able to modify it.

   - If Docker Desktop is already installed and being used, `locked: false` is ignored. This is because existing users of Docker Desktop may have already updated a setting, which in turn will have been written to the relevant config file, for example the `settings-store.json` (or `settings.json` for Docker Desktop versions 4.34 and earlier) or `daemon.json`. In these instances, the user's preferences are respected and the values aren't altered. These can be controlled by setting `locked: true`.

The following `admin-settings.json` code and table provides an example of the required syntax and descriptions for parameters and values:

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
  "allowExperimentalFeatures": {
    "locked": false,
    "value": false
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
|`kubernetes`|  | If `enabled` is set to true, a Kubernetes single-node cluster is started when Docker Desktop starts. If `showSystemContainers` is set to true, Kubernetes containers are displayed in the Docker Desktop Dashboard and when you run `docker ps`.  `imagesRepository` lets you specify which repository Docker Desktop pulls the Kubernetes images from. For example, `"imagesRepository": "registry-1.docker.io/docker"`.  |  |

### Features in development 

|Parameter|OS|Description|Version|
|:-------------------------------|---|:-------------------------------|---|
| `allowExperimentalFeatures`| | If `value` is set to `false`, experimental features are disabled.|  |
| `allowBetaFeatures`| | If `value` is set to `false`, beta features are disabled.|  |
| `enableDockerAI` | | If `value` is set to `false`, Docker AI (Ask Gordon) features are disabled. |  |

### Enhanced Container Isolation 

|Parameter|OS|Description|Version|
|:-------------------------------|---|:-------------------------------|---|
|`enhancedContainerIsolation`|  | If `value` is set to true, Docker Desktop runs all containers as unprivileged, via the Linux user-namespace, prevents them from modifying sensitive configurations inside the Docker Desktop VM, and uses other advanced techniques to isolate them. For more information, see [Enhanced Container Isolation](../enhanced-container-isolation/_index.md).|  |
| &nbsp; &nbsp; &nbsp; &nbsp;`dockerSocketMount` |  | By default, enhanced container isolation blocks bind-mounting the Docker Engine socket into containers (e.g., `docker run -v /var/run/docker.sock:/var/run/docker.sock ...`). This lets you relax this in a controlled way. See [ECI Configuration](../enhanced-container-isolation/config.md) for more info. |  |
| &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; `imageList` |  | Indicates which container images are allowed to bind-mount the Docker Engine socket. |  |
| &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; `commandList` |  | Restricts the commands that containers can issue via the bind-mounted Docker Engine socket. |  |

## Step three: Re-launch Docker Desktop

> [!NOTE]
>
> Test the changes made through the `admin-settings.json` file locally to see if the settings work as expected.

For settings to take effect:
- On a new install, developers need to launch Docker Desktop and authenticate to their organization.
- On an existing install, developers need to quit Docker Desktop through the Docker menu, and then re-launch Docker Desktop. If they are already signed in, they don't need to sign in again for the changes to take effect.
  > [!IMPORTANT]
  >
  > Selecting **Restart** from the Docker menu isn't enough as it only restarts some components of Docker Desktop.

So as not to disrupt your developers' workflow, Docker doesn't automatically mandate that developers re-launch and re-authenticate once a change has been made.

In Docker Desktop, developers see the relevant settings grayed out and the message **Locked by your administrator**.

![Proxy settings grayed out with Settings Management](/assets/images/grayed-setting.png)
