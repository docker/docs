---
title: Install and configure
description:
keywords:
---

## Install 

TBD

## Configuration options 

> **Important**
>
> In addition to the following custom properties, the Docker Desktop MSI installer also supports the standard [Windows Installer command line options](https://learn.microsoft.com/en-us/windows/win32/msi/standard-installer-command-line-options).
{ .important }

| Property | Description | Default |
| --- | --- | --- |
| `ENABLEDESKTOPSHORTCUT` | Creates a shortcut on the current users desktop | 1 |
| `ADMINSETTINGS` | Automatically creates an admin-settings.json file which is used by admins to control certain Docker Desktop settings on client machines within their organization. It must be used together with the `ALLOWEDORG` property. | None |
| `ALLOWEDORG` | Requires the user to sign in and be part of the specified Docker Hub organization when running the application. This creates the regsitry.json file containing the specified organisations. | None |
| `ALWAYSRUNSERVICE` | Lets users switch to Windows containers without needing admin rights | 0 |
| `DISABLEWINDOWSCONTAINERS` | Disables the Windows containers integration | 0 |
| `ENGINE` | The docker engine that will be used to run containers. This can be one of:
`wsl` , `hyperv` or `windows` | `wsl` |
| `PROXYHTTPMODE` | Sets the HTTP Proxy mode. This can be one of:
`system` or `manual` | `system` |
| `OVERRIDEPROXYHTTP` | Sets the URL of the HTTP proxy that must be used for outgoing HTTP requests. | None |
| `OVERRIDEPROXYHTTPS` | Sets the URL of the HTTP proxy that must be used for outgoing HTTPS requests. | None |
| `OVERRIDEPROXYEXCLUD`E | Bypasses proxy settings for the hosts and domains. Uses a comma-separated list. | None |
| `HYPERVDEFAULTDATAROOT` | Specifies the default location for the Hyper-V VM disk. | None |
| `WINDOWSCONTAINERSDEFAULTDATAROOT` | Specifies the default location for the Windows containers. | None |
| `WSLDEFAULTDATAROOT` | Specifies the default location for the WSL distribution disk. | None |
| `DISABLEENGINEINSTALL` | TBA | TBA |
| `INSTALLFOLDER` | Specifies a custom location where Docker Desktop will be installed TODO: We will try to include for a later stage. | C:\Program Files\Docker |

### Silent installations

For silent installations you can use `/quiet` (or `qn`). This runs the installer without displaying a user interface.

Additionally,you can also use `/norestart` or `/forcerestart` to control reboot behaviour.

By default, the installer reboots the machine after a successful installation. When ran silently, the reboot is automatic and the user is not prompted.

## Installation scenarios 

This section covers command line installations of Docker Desktop using PowerShell.

Interactive installations, without specifying `/quiet` or `/qn`, display the user interface and let the user select their own properties.

Non-interactive installations are silent and any additional configuration must be passed as arguments.

> **Note**
> 
> Non-interactive installations must be executed from an elevated terminal session.

### Installing interactively with verbose logging

```powershell
msiexec /i "DockerDesktop.msi" /L*V ".\msi.log"
```

### Installing non-interactively with verbose logging

```powershell
msiexec /i "DockerDesktop.msi" /L*V ".\msi.log" /quiet
```

### Installing non-interactively and suppressing reboots

```powershell
msiexec /i "DockerDesktop.msi" /L*V ".\msi.log" /quiet /norestart
```

### Installing non-interactively with Admin settings

```powershell
msiexec /i "DockerDesktop.msi" /L*V ".\msi.log" /quiet /norestart ADMINSETTINGS="{"configurationFileVersion":2,"enhancedContainerIsolation":{"value":true,"locked":false}}" ALLOWEDORG="docker.com"
```

> **Tip**
>
> Here are some useful tips to remember when creating a value that expects a JSON string as it’s value:
> 
> - The property expects a JSON formatted string
> - The string should be wrapped in double quotes
> - The string should not contain any whitespace
> - Property names are expected to be in double quotes
{ .tip }

### Uninstalling interactively

```powershell
msiexec /x "DockerDesktop.msi"
```

### Uninstalling non-interactively 

```powershell
msiexec /x "DockerDesktop.msi" /quiet
```

