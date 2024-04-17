---
title: Install and configure
description: Understand how to use the MSI installer. Also explore additional configuration options.
keywords: msi, windows, docker desktop, install, deploy, configure
---

## Install interactively

1. Download the installer.

2. Double-click `Docker Desktop Installer.msi` to run the installer. 

3. Once you've accepted the license agreement, you can choose the install location. By default, Docker Desktop is installed at `C:\Program Files\Docker\Docker`.

3. Configure the Docker Desktop installation. You can:

- Create a desktop shortcut

- Set the Docker Desktop service startup type to automatic

- Disable Windows Container usage

- Select the engine for Docker Desktop. Either WSL or Hyper-V. If your system only supports one of the two options, you will not be able to select which backend to use.

4. Follow the instructions on the installation wizard to authorize the installer and proceed with the install.

5. When the installation is successful, select **Finish** to complete the installation process.

If your admin account is different to your user account, you must add the user to the **docker-users** group:
1. Run **Computer Management** as an **administrator**.
2. Navigate to **Local Users and Groups** > **Groups** > **docker-users**. 
3. Right-click to add the user to the group.
4. Sign out and sign back in for the changes to take effect.

## Install from the command line

This section covers command line installations of Docker Desktop using PowerShell. It provides common installation commands that you can run. You can also add additional arguments which are outlined in [configuration options](#configuration-options).

When installing Docker Desktop, you can choose between interactive or non-interactive installations. 

Interactive installations, without specifying `/quiet` or `/qn`, display the user interface and let you select your own properties. 

When installing via the user interface it is possible to:

- Choose the destination folder
- Create a desktop shortcut
- Configure the Docker Desktop service startup type
- Disable Windows Containers
- Choose between the WSL or Hyper-V engine

Non-interactive installations are silent and any additional configuration must be passed as arguments. They must also be executed from an elevated terminal.

### Common installation commands

#### Installing interactively with verbose logging

```powershell
msiexec /i "DockerDesktop.msi" /L*V ".\msi.log"
```

#### Installing non-interactively with verbose logging

```powershell
msiexec /i "DockerDesktop.msi" /L*V ".\msi.log" /quiet
```

#### Installing non-interactively and suppressing reboots

```powershell
msiexec /i "DockerDesktop.msi" /L*V ".\msi.log" /quiet /norestart
```

#### Installing non-interactively with Admin settings

```powershell
msiexec /i "DockerDesktop.msi" /L*V ".\msi.log" /quiet /norestart ADMINSETTINGS="{"configurationFileVersion":2,"enhancedContainerIsolation":{"value":true,"locked":false}}" ALLOWEDORG="docker.com"
```

> **Tip**
>
> Some useful tips to remember when creating a value that expects a JSON string as it’s value:
> 
> - The property expects a JSON formatted string
> - The string should be wrapped in double quotes
> - The string should not contain any whitespace
> - Property names are expected to be in double quotes
{ .tip }

#### Uninstalling interactively

```powershell
msiexec /x "DockerDesktop.msi"
```

#### Uninstalling non-interactively 

```powershell
msiexec /x "DockerDesktop.msi" /quiet
```

### Configuration options 

> **Important**
>
> In addition to the following custom properties, the Docker Desktop MSI installer also supports the standard [Windows Installer command line options](https://learn.microsoft.com/en-us/windows/win32/msi/standard-installer-command-line-options).
{ .important }

| Property | Description | Default |
| :--- | :--- | :--- |
| `ENABLEDESKTOPSHORTCUT` | Creates a desktop shortcut. | 1 |
| `INSTALLFOLDER` | Specifies a custom location where Docker Desktop will be installed. | C:\Program Files\Docker |
| `ADMINSETTINGS` | Automatically creates an `admin-settings.json` file which is used to [control certain Docker Desktop settings](../../hardened-desktop/settings-management/_index.md) on client machines within organizations. It must be used together with the `ALLOWEDORG` property. | None |
| `ALLOWEDORG` | Requires the user to sign in and be part of the specified Docker Hub organization when running the application. This creates the `regsitry.json` file containing the specified organisations. | None |
| `ALWAYSRUNSERVICE` | Lets users switch to Windows containers without needing admin rights | 0 |
| `DISABLEWINDOWSCONTAINERS` | Disables the Windows containers integration | 0 |
| `ENGINE` | Sets the Docker Engine that is used to run containers. This can be either `wsl` , `hyperv`, or `windows` | `wsl` |
| `PROXYHTTPMODE` | Sets the HTTP Proxy mode. This can be either `system` or `manual` | `system` |
| `OVERRIDEPROXYHTTP` | Sets the URL of the HTTP proxy that must be used for outgoing HTTP requests. | None |
| `OVERRIDEPROXYHTTPS` | Sets the URL of the HTTP proxy that must be used for outgoing HTTPS requests. | None |
| `OVERRIDEPROXYEXCLUDE` | Bypasses proxy settings for the hosts and domains. Uses a comma-separated list. | None |
| `HYPERVDEFAULTDATAROOT` | Specifies the default location for the Hyper-V VM disk. | None |
| `WINDOWSCONTAINERSDEFAULTDATAROOT` | Specifies the default location for Windows containers. | None |
| `WSLDEFAULTDATAROOT` | Specifies the default location for the WSL distribution disk. | None |


Additionally, you can also use `/norestart` or `/forcerestart` to control reboot behaviour.

By default, the installer reboots the machine after a successful installation. When ran silently, the reboot is automatic and the user is not prompted.


