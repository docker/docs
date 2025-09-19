---
title: MSI installer
description: Understand how to use the MSI installer. Also explore additional configuration options.
keywords: msi, windows, docker desktop, install, deploy, configure, admin, mdm
tags: [admin]
weight: 10
aliases:
 - /desktop/install/msi/install-and-configure/
 - /desktop/setup/install/msi/install-and-configure/
 - /desktop/install/msi/
 - /desktop/setup/install/msi/
 - /desktop/setup/install/enterprise-deployment/msi-install-and-configure/
---

{{< summary-bar feature_name="MSI installer" >}}

The MSI package supports various MDM (Mobile Device Management) solutions, making it ideal for bulk installations and eliminating the need for manual setups by individual users. With this package, IT administrators can ensure standardized, policy-driven installations of Docker Desktop, enhancing efficiency and software management across their organizations.

## Install interactively

1. In [Docker Home](http://app.docker.com), choose your organization.
2. Select **Admin Console**, then **Enterprise deployment**.
3. From the **Windows OS** tab, select the **Download MSI installer** button.
4. Once downloaded, double-click `Docker Desktop Installer.msi` to run the installer.
5. After accepting the license agreement, choose the install location. By default, Docker Desktop is installed at `C:\Program Files\Docker\Docker`.
6. Configure the Docker Desktop installation. You can:

    - Create a desktop shortcut

    - Set the Docker Desktop service startup type to automatic

    - Disable Windows Container usage

    - Select the Docker Desktop backend: WSL or Hyper-V. If only one is supported by your system, you won't be able to choose.
7. Follow the instructions on the installation wizard to authorize the installer and proceed with the install.
8. When the installation is successful, select **Finish** to complete the installation process.

If your administrator account is different from your user account, you must add the user to the **docker-users** group to access features that require higher privileges, such as creating and managing the Hyper-V VM, or using Windows containers:

1. Run **Computer Management** as an **administrator**.
2. Navigate to **Local Users and Groups** > **Groups** > **docker-users**.
3. Right-click to add the user to the group.
4. Sign out and sign back in for the changes to take effect.

> [!NOTE]
>
> When installing Docker Desktop with the MSI, in-app updates are automatically disabled. This ensures organizations can maintain version consistency and prevent unapproved updates. For Docker Desktop installed with the .exe installer, in-app updates remain supported.
>
> Docker Desktop notifies you when an update is available. To update Docker Desktop, download the latest installer from the Docker Admin Console. Navigate to the **Enterprise deployment** page.
>
> To keep up to date with new releases, check the [release notes](/manuals/desktop/release-notes.md) page.

## Install from the command line

This section covers command line installations of Docker Desktop using PowerShell. It provides common installation commands that you can run. You can also add additional arguments which are outlined in [configuration options](#configuration-options).

When installing Docker Desktop, you can choose between interactive or non-interactive installations.

Interactive installations, without specifying `/quiet` or `/qn`, display the user interface and let you select your own properties.

When installing via the user interface it's possible to:

- Choose the destination folder
- Create a desktop shortcut
- Configure the Docker Desktop service startup type
- Disable Windows Containers
- Choose between the WSL or Hyper-V engine

Non-interactive installations are silent and any additional configuration must be passed as arguments.

### Common installation commands

> [!IMPORTANT]
>
> Admin rights are required to run any of the following commands.

#### Install interactively with verbose logging

```powershell
msiexec /i "DockerDesktop.msi" /L*V ".\msi.log"
```

#### Install interactively without verbose logging

```powershell
msiexec /i "DockerDesktop.msi"
```

#### Install non-interactively with verbose logging

```powershell
msiexec /i "DockerDesktop.msi" /L*V ".\msi.log" /quiet
```

#### Install non-interactively and suppressing reboots

```powershell
msiexec /i "DockerDesktop.msi" /L*V ".\msi.log" /quiet /norestart
```

#### Install non-interactively with admin settings

```powershell
msiexec /i "DockerDesktop.msi" /L*V ".\msi.log" /quiet /norestart ADMINSETTINGS="{""configurationFileVersion"":2,""enhancedContainerIsolation"":{""value"":true,""locked"":false}}" ALLOWEDORG="your-organization"
```

#### Install interactively and allow users to switch to Windows containers without admin rights

```powershell
msiexec /i "DockerDesktop.msi" /L*V ".\msi.log" /quiet /norestart ALLOWEDORG="your-organization" ALWAYSRUNSERVICE=1
```

#### Install with the passive display option

You can use the `/passive` display option instead of `/quiet` when you want to perform a non-interactive installation but show a progress dialog.

In passive mode the installer doesn't display any prompts or error messages to the user and the installation cannot be cancelled.

For example:

```powershell
msiexec /i "DockerDesktop.msi" /L*V ".\msi.log" /passive /norestart
```

> [!TIP]
>
> When creating a value that expects a JSON string:
>
> - The property expects a JSON formatted string
> - The string should be wrapped in double quotes
> - The string shouldn't contain any whitespace
> - Property names are expected to be in double quotes

### Common uninstall commands

When uninstalling Docker Desktop, you need to use the same `.msi` file that was originally used to install the application.

If you no longer have the original `.msi` file, you need to use the product code associated with the installation. To find the product code, run:

```powershell
Get-WmiObject Win32_Product | Select-Object IdentifyingNumber, Name | Where-Object {$_.Name -eq "Docker Desktop"}
```

It should return output similar to the following:

```text
IdentifyingNumber                      Name
-----------------                      ----
{10FC87E2-9145-4D7D-B493-2E99E8D8E103} Docker Desktop
```
> [!NOTE]
>
> This command may take some time, depending on the number of installed applications.

`IdentifyingNumber` is the applications product code and can be used to uninstall Docker Desktop. For example:

```powershell
msiexec /x {10FC87E2-9145-4D7D-B493-2E99E8D8E103} /L*V ".\msi.log" /quiet
```

#### Uninstall interactively with verbose logging

```powershell
msiexec /x "DockerDesktop.msi" /L*V ".\msi.log"
```

#### Uninstall interactively without verbose logging

```powershell
msiexec /x "DockerDesktop.msi"
```

#### Uninstall non-interactively with verbose logging

```powershell
msiexec /x "DockerDesktop.msi" /L*V ".\msi.log" /quiet
```

#### Uninstall non-interactively without verbose logging

```powershell
msiexec /x "DockerDesktop.msi" /quiet
```

### Configuration options

> [!IMPORTANT]
>
> In addition to the following custom properties, the Docker Desktop MSI installer also supports the standard [Windows Installer command line options](https://learn.microsoft.com/en-us/windows/win32/msi/standard-installer-command-line-options).

| Property | Description | Default |
| :--- | :--- | :--- |
| `ENABLEDESKTOPSHORTCUT` | Creates a desktop shortcut. | 1 |
| `INSTALLFOLDER` | Specifies a custom location where Docker Desktop will be installed. | C:\Program Files\Docker |
| `ADMINSETTINGS` | Automatically creates an `admin-settings.json` file which is used to [control certain Docker Desktop settings](/manuals/enterprise/security/hardened-desktop/settings-management/_index.md) on client machines within organizations. It must be used together with the `ALLOWEDORG` property. | None |
| `ALLOWEDORG` | Requires the user to sign in and be part of the specified Docker Hub organization when running the application. This creates a registry key called `allowedOrgs` in `HKLM\Software\Policies\Docker\Docker Desktop`. | None |
| `ALWAYSRUNSERVICE` | Lets users switch to Windows containers without needing admin rights | 0 |
| `DISABLEWINDOWSCONTAINERS` | Disables the Windows containers integration | 0 |
| `ENGINE` | Sets the Docker Engine that's used to run containers. This can be either `wsl` , `hyperv`, or `windows` | `wsl` |
| `PROXYENABLEKERBEROSNTLM` | When set to 1, enables support for Kerberos and NTLM proxy authentication. Available with Docker Desktop 4.33 and later| 0 |
| `PROXYHTTPMODE` | Sets the HTTP Proxy mode. This can be either `system` or `manual` | `system` |
| `OVERRIDEPROXYHTTP` | Sets the URL of the HTTP proxy that must be used for outgoing HTTP requests. | None |
| `OVERRIDEPROXYHTTPS` | Sets the URL of the HTTP proxy that must be used for outgoing HTTPS requests. | None |
| `OVERRIDEPROXYEXCLUDE` | Bypasses proxy settings for the hosts and domains. Uses a comma-separated list. | None |
| `HYPERVDEFAULTDATAROOT` | Specifies the default location for the Hyper-V VM disk. | None |
| `WINDOWSCONTAINERSDEFAULTDATAROOT` | Specifies the default location for Windows containers. | None |
| `WSLDEFAULTDATAROOT` | Specifies the default location for the WSL distribution disk. | None |
| `DISABLEANALYTICS` | When set to 1, analytics collection will be disabled for the MSI. For more information, see [Analytics](#analytics). | 0 |


Additionally, you can also use `/norestart` or `/forcerestart` to control reboot behaviour.

By default, the installer reboots the machine after a successful installation. When run silently, the reboot is automatic and the user is not prompted.

## Analytics

The MSI installer collects anonymous usage statistics relating to installation only. This is to better understand user behaviour and to improve the user experience by identifying and addressing issues or optimizing popular features.

### How to opt-out

{{< tabs >}}
{{< tab name="From the GUI" >}}

When you install Docker Desktop from the default installer GUI, select the **Disable analytics** checkbox located on the bottom-left corner of the **Welcome** dialog.

{{< /tab >}}
{{< tab name="From the command line" >}}

When you install Docker Desktop from the command line, use the `DISABLEANALYTICS` property.

```powershell
msiexec /i "win\msi\bin\en-US\DockerDesktop.msi" /L*V ".\msi.log" DISABLEANALYTICS=1
```

{{< /tab >}}
{{< /tabs >}}

### Persistence

If you decide to disable analytics for an installation, your choice is persisted in the registry and honoured across future upgrades and uninstalls.

However, the key is removed when Docker Desktop is uninstalled and must be configured again via one of the previous methods.

The registry key is as follows:

```powershell
SOFTWARE\Docker Inc.\Docker Desktop\DisableMsiAnalytics
```

When analytics is disabled, this key is set to `1`.

## Additional resources

- [Explore the FAQs](/manuals/enterprise/enterprise-deployment/faq.md)
