---
description: settings management for desktop
keywords: admin, controls, rootless, enhanced container isolation
title: Configure Settings Management
--- 

>**Note**
>
>Settings Management is available to Docker Business customers only. 

This page contains information for admins on how to configure Settings Management to specify and lock configuration parameters to create a standardized Docker Desktop environment across the organization.

Settings Management is designed specifically for organizations who don’t give developers root access to their machines.

### Prerequisites

- [Download and install Docker Desktop 4.13.0 or later](../../release-notes.md).
- As an admin, you need to [configure a registry.json to enforce sign-in](../../../docker-hub/configure-sign-in.md). This is because this feature requires a Docker Business subscription and therefore your Docker Desktop users must authenticate to your organization for this configuration to take effect.

### Step one: Create the `admin-settings.json` file and save it in the correct location

You can either use the `--admin-settings` installer flag on [macOS](../../install/mac-install.md#install-from-the-command-line) or [Windows](../../install/windows-install.md#install-from-the-command-line) to automatically create the `admin-settings.json` and save it in the correct location, or set it up manually. 

To set it up manually:
1. Create a new, empty JSON file and name it `admin-settings`.
2. Save the `admin-settings.json` file on your developers' machines in the following locations:

    - Mac: `/Library/Application\ Support/com.docker.docker/admin-settings.json`
    - Windows: `C:\ProgramData\DockerDesktop\admin-settings.json`
    - Linux: `/usr/share/docker-desktop/admin-settings.json`

    By placing this file in the above protected directories, end users are unable to modify it.

    >**Note**
    >
    > It is assumed that you have the ability to push the `admin-settings.json` settings file to the locations specified above through a device management software such as [Jamf](https://www.jamf.com/lp/en-gb/apple-mobile-device-management-mdm-jamf-shared/?attr=google_ads-brand-search-shared&gclid=CjwKCAjw1ICZBhAzEiwAFfvFhEXjayUAi8FHHv1JJitFPb47C_q_RCySTmF86twF1qJc_6GST-YDmhoCuJsQAvD_BwE).

### Step two: Configure the settings you want to lock in

>**Note**
>
>Some of the configuration parameters only apply to Windows. This is highlighted in the table below.

The `admin-settings.json` file requires a nested list of configuration parameters, each of which must contain the  `locked` parameter. You can add or remove configuration parameters as per your requirements.

If `locked: true`, users aren't able to edit this setting from Docker Desktop or the CLI.

If `locked: false`, it's similar to setting a factory default in that:
- For new installs, `locked: false` pre-populates the relevant settings in the Docker Desktop UI, but users are able to modify it.

- If Docker Desktop is already installed and being used, `locked: false` is ignored. This is because existing users of Docker Desktop may have already updated a setting, which in turn will have been written to the relevant config file, for example the `settings.json` or `daemon.json`. In these instances, the user's preferences are respected and we don't alter these values. These can be controlled by the admin by setting `locked: true`.

The following `admin-settings.json` code and table provides an example of the required syntax and descriptions for parameters and values:

```json
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
    "exclude": []
  },
  "enhancedContainerIsolation": {
    "locked": true,
    "value": true
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
  }
}
```

| Parameter                        |   | Description                      |  
| :------------------------------- |---| :------------------------------- |
| `configurationFileVersion`        |   |Specifies the version of the configuration file format.    |
| `exposeDockerAPIOnTCP2375` | <span class="badge badge-info">Windows only</span>| Exposes the Docker API on a specified port. If `value` is set to true, the Docker API is exposed on port 2375. Note: This is unauthenticated and should only be enabled if protected by suitable firewall rules.| 
| `proxy` |   |If `mode` is set to `system` instead of `manual`, Docker Desktop gets the proxy values from the system and ignores and values set for `http`, `https` and `exclude`. Change `mode` to `manual` to manually configure proxy servers. If the proxy port is custom, specify it in the `http` or `https` property, for example `"https": "http://myotherproxy.com:4321"`. The `exclude` property specifies a comma-separated list of hosts and domains to bypass the proxy. |
| `enhancedContainerIsolation`  |  | If `value` is set to true, Docker Desktop runs all containers as unprivileged, via the Linux user-namespace, prevents them from modifying sensitive configurations inside the Docker Desktop VM, and uses other advanced techniques to isolate them. For more information, see [Enhanced Container Isolation](../enhanced-container-isolation/index.md). Note: Enhanced Container Isolation is currently [incompatible with WSL](../enhanced-container-isolation/faq.md#incompatibility-with-windows-subsystem-for-linux-wsl). | 
| `linuxVM` |   |Parameters and settings related to Linux VM options - grouped together here for convenience. |
| &nbsp; &nbsp; &nbsp; &nbsp;`wslEngineEnabled`  | <span class="badge badge-info">Windows only</span> | If `value` is set to true, Docker Desktop uses the WSL 2 based engine. This overrides anything that may have been set at installation using the `--backend=<backend name>` flag. It is also incompatible with Enhanced Container Isolation. See [Known issues](../enhanced-container-isolation/faq.md) for more information.| 
| &nbsp;&nbsp; &nbsp; &nbsp;`dockerDaemonOptions`|  |If `value` is set to true, it overrides the options in the Docker Engine config file. See the [Docker Engine reference](/engine/reference/commandline/dockerd/#daemon-configuration-file). Note that for added security, a few of the config attributes may be overridden when Enhanced Container Isolation is enabled. |
| &nbsp;&nbsp; &nbsp; &nbsp;`vpnkitCIDR` |  |Overrides the network range used for vpnkit DHCP/DNS for `*.docker.internal`  |
| `windowsContainers` |  | Parameters and settings related to `windowsContainers` options - grouped together here for convenience.                  |
| &nbsp; &nbsp; &nbsp; &nbsp;`dockerDaemonOptions` |  | Overrides the options in the linux daemon config file. See the [Docker Engine reference](/engine/reference/commandline/dockerd/#daemon-configuration-file).|                                |
|`disableUpdate`|  |If `value` is set to true, checking for and notifications about Docker Desktop updates is disabled.|
|`analyticsEnabled`|  |If `value` is set to false, Docker Desktop doesn't send usage statistics to Docker. |

### Step three: Re-launch Docker Desktop
>**Note**
>
>Administrators should test the changes made through the `admin-settings.json` file locally to see if the settings work as expected.

For settings to take effect:
- On a new install, developers need to launch Docker Desktop and authenticate to their organization.
- On an existing install, developers need to quit Docker Desktop through the Docker menu, and then relaunch Docker Desktop. If they are already signed in, they don't need to sign in again for the changes to take effect.
  >**Important**
  >
  >Selecting **Restart** from the Docker menu isn't enough as it only restarts some components of Docker Desktop.
  {: .important}

Docker doesn't automatically mandate that developers re-launch and sign in once a change has been made so as not to disrupt your developers' workflow. 


In Docker Desktop, developers see the relevant settings grayed out and the message **Locked by your administrator**.

![Proxy settings grayed out](/assets/images/grayed-setting.png){:width="750px"}
