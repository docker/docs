---
description: admin controls for desktop
keywords: admin, controls, rootless, enhanced container isolation
title: Configure Admin Controls
--- 

>Note
>
>Admin Controls is available to Docker Business customers only. 

This page contains information about how admins can configure Admin Controls to specify and lock configuration parameters to create a standardized Docker Desktop environment across the organization.

Admin Controls is designed specifically for organizations who don’t give developers root access to their machines.

## Prerequisite

You need to [configure a registry.json to enforce sign-in](../../../docker-hub/configure-sign-in.md). For this configuration to take effect, Docker Desktop users must authenticate to your organization. 

## Step one: Place the `admin-settings.json` file in the correct location

Place the `admin-settings.json` file on your developers' machines in the following locations:

- Mac: `/Library/Application Support/com.docker.docker/admin-settings.json`
- Windows: `/ProgramData/DockerDesktop/admin-settings.json`
- Linux: `/usr/share/docker-desktop/admin-settings.json`

By placing this file in the above protected directories, end users are unable to modify it.

>Note
>
> It is assumed that you have the ability to push the `admin-settings.json` settings file to the locations specified above through a device management software such as [Jamf](https://www.jamf.com/lp/en-gb/apple-mobile-device-management-mdm-jamf-shared/?attr=google_ads-brand-search-shared&gclid=CjwKCAjw1ICZBhAzEiwAFfvFhEXjayUAi8FHHv1JJitFPb47C_q_RCySTmF86twF1qJc_6GST-YDmhoCuJsQAvD_BwE).

## Step two: Configure the admin controls you want to lock in

>Note
>
>Some of the configuration parameters only apply to Windows. This is highlighted in the table below.

The `admin-settings.json` file requires a nested list of configuration parameters, each of which must contain the  `locked` parameter. You can add or remove configuration parameters as per your requirements.
If set to `true`, users are not able to edit this setting from Docker Desktop or the CLI. 
If set to `false`, the configuration value acts as a default value, but users can change this setting from Docker Desktop or the CLI by directly editing the `settings.json` file.

<div class="panel panel-default">
<div class="panel-heading collapsed" data-toggle="collapse" data-target="#collapseSample2"  style="cursor: pointer"> Additional <code>Locked: false</code> information<i class="chevron fa fa-fw"></i></div>
<div class="collapse block" id="collapseSample2">
<p>
<code>Locked: false</code> is similar to having a setting be the factory default.
<br>
<li>For new installs, <code>Locked: false</code> pre-populates the relevant settings in the Desktop UI.</li>
<br>
<li> If Docker Desktop is already installed and being used, `locked: false` is ignored. This is because existing users of Docker Desktop may have already updated a setting which may have been written to the relevant config file, for example the <code>settings.json</code> or <code>daemon.json</code> file. In these instances, Docker respects the user's preference. This can be overridden by the admin by setting <code>locked: true</code>.</li>
</p>
 </div>
</div>

The following `admin-settings.json` code and table provides an example of the required syntax and descriptions for parameters and values:

```json
{
  "configurationFileVersion": 2,
  "exposeDockerAPIOnTCP2375": {
    "locked": true,
    "value": false
  },
  "proxy": {
    "locked": false,
    "mode": "system",
    "server": "myproxy.com",
    "port":3129,
    "exclude": ["foo.com", "bar.com"]
  },
  "enhancedContainerIsolation": {
    "locked": false,
    "value": false
  },
  "useWindowsContainers": {
      "locked": false,
      "value": false
  },
  "linuxVM": {
    "wslEngineEnabled": {
      "locked": false,
      "value": false
    },
    "dockerDaemonOptions": {
      "locked": false,
      "value":"<json string>"
    },
    "vpnkitCIDR": {
      "locked": false,
      "value":"192.168.65.0/24"
    }
  },
  "windowsContainers": {
    "dockerDaemonOptions": {
      "locked": false,
      "value":"<json string>"
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
}
```

| Parameter                        | Description                      |
| :------------------------------- | :------------------------------- |
| `configurationFileVersion`        | Specifies the version of the configuration file format.    |
| `exposeDockerAPIOnTCP2375` |<span class="badge badge-info">Windows only</span> Exposes the Docker API on a specified port. If `value` is set to true, the Docker API is exposed on port 2375. Note: This is unauthenticated and should only be enabled if protected by suitable firewall rules.|
| `proxy` | It is used for `http` and `https`.  If the port is custom, specify it in the property. |
| `enhancedContainerIsolation`  | If `value` is set to true, Docker Desktop runs all containers as unprivileged, via the Linux user-namespace, prevents them from modifying sensitive configurations inside the Docker Desktop VM, and uses other advanced techniques to isolate them. For more information, see [Enhanced Container Isolation](../enhanced-container-isolation/index.md). |
|`useWindowsContainers` | <span class="badge badge-info">Windows only</span> If `value` is set to true, it switches Docker Desktop to toggle the Docker CLI to talk to the Windows daemon, enabling Windows containers. If false, switches Docker Desktop to toggle the Docker CLI to talk to the Linux daemon, enabling Linux containers. This overrides anything that may have been set at installation using the `--no-windows-containers` flag.|
| `linuxVM` | Parameters and settings related to Linux VM options - grouped together here for convenience. |
| &nbsp; &nbsp; &nbsp; &nbsp;`wslEngineEnabled`  |<span class="badge badge-info">Windows only</span> If `value` is set to true, Docker Desktop uses the WSL 2 based engine. This overrides anything that may have been set at installation using the `--backend=<backend name>` flag. It is also incompatible with Enhanced Container Isolation. See [Known issues](faq.md) for more information.|
| &nbsp;&nbsp; &nbsp; &nbsp;`dockerDaemonOptions`|If `value` is set to true, it overrides the options in the Linux daemon config file. See the [Docker Engine reference](../../../engine/reference/commandline/dockerd/#daemon-configuration-file). |
| &nbsp;&nbsp; &nbsp; &nbsp;`vpnkitCIDR` |Overrides the network range used for vpnkit DHCP/DNS for `*.docker.internal`  |
| `windowsContainers` | Parameters and settings related to `windowsContainers` options - grouped together here for convenience.                  |
| &nbsp; &nbsp; &nbsp; &nbsp;`dockerDaemonOptions` | Overrides the options in the linux daemon config file. See the [Docker Engine reference](../../../engine/reference/commandline/dockerd/#daemon-configuration-file).|                                |
|`disableUpdate`|If `value` is set to true, checking for and notifications about Docker Desktop updates is disabled.|
|`analyticsEnabled`|If `value` is set to false, Docker Desktop does not send usage statistics to Docker. |

## Step three: Re-launch and re-authenticate
>Note
>
>Administrators should test the changes made through the `admin-settings.json` file locally to see if the settings work as expected.

Once you have created and configured `admin-settings.json`, developers need to quit Docker Desktop through the Whale menu, and then relaunch Docker Desktop and sign in to receive the changed settings. Selecting **Restart** from the Whale menu isn't sufficient as it only restarts some components of Docker Desktop.

Docker doesn't automatically mandate that developers re-launch and re-authenticate once a change has been made so as not to disrupt your developers workflow. 

In Docker Desktop, developers see the relevant settings grayed out and the message **Locked by your administrator**.

![Proxy settings grayed out](/assets/images/grayed-setting.png){:width="750px"}
