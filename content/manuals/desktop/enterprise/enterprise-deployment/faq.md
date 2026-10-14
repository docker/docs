---
title: Enterprise deployment FAQs
linkTitle: FAQs
description: Frequently asked questions for deploying Docker Desktop at scale
keywords: msi, deploy, docker desktop, faqs, pkg, mdm, jamf, intune, windows, mac, enterprise, admin
tags: [FAQ, admin]
weight: 70
aliases:
 - /enterprise/enterprise-deployment/faq/
---

## MSI

Common questions about installing Docker Desktop using the MSI installer.

### What happens to user data if they have an older Docker Desktop installation (i.e. `.exe`)?

Users must [uninstall](/manuals/desktop/uninstall.md) older `.exe` installations before using the new MSI version. The `.exe` installer includes a `--keep-data` flag that removes Docker Desktop while preserving underlying resources such as the container VMs:

```powershell
# For all-user installations
& 'C:\Program Files\Docker\Docker\Docker Desktop Installer.exe' uninstall --keep-data

# For per-user installations
& '%LOCALAPPDATA%\Programs\DockerDesktop\Docker Desktop Installer.exe' uninstall --keep-data

```

For all-users installations, you can have the MSI do this for you with the `REMOVEEXISTINGINSTALL` property, described in the next answer.

### What happens if the user's machine has an older `.exe` installation?

### What happens if the user's machine has an older `.exe` installation?

The MSI installer detects existing `.exe` installations and, by default, blocks the installation. How you resolve it depends on whether the `.exe` was installed for all users or for a single user.

#### All-users `.exe` installation

The installation stops with:

```text
You need to uninstall the previous Docker Desktop version in order to use the MSI installer.
```

Either uninstall it first with `--keep-data` as described in the previous answer, or let the MSI do it by setting `REMOVEEXISTINGINSTALL=1`:

```powershell
msiexec /i "DockerDesktop.msi" /L*V ".\msi.log" /quiet /norestart REMOVEEXISTINGINSTALL=1
```

This runs the existing uninstaller with `--keep-data`, so settings and container data are preserved.

#### Per-user `.exe` installation

Available with Docker Desktop version 4.84 and later, the installation stops with:

```text
Docker Desktop is installed per-user for one or more accounts on this machine: <usernames>.
Please have each affected user uninstall Docker Desktop first before running the MSI installer.
```

`REMOVEEXISTINGINSTALL` doesn't help here. The MSI runs with machine-wide privileges and can't reliably uninstall software installed under another user's profile, so each listed user must uninstall Docker Desktop themselves before the MSI can proceed.

With Docker Desktop version 4.83 and earlier, the MSI doesn't detect per-user installations.

> [!NOTE]
>
> Per-user installations became more common with Docker Desktop version 4.83, when the EXE installer started selecting a per-user installation by default. Expect to encounter them on machines where developers installed Docker Desktop themselves.

### Can I install the MSI per-user?

No. The MSI installer only supports all-users installations.

Available with Docker Desktop version 4.92 and later, passing `MSIINSTALLPERUSER` fails with:

```text
Docker Desktop does not support per-user installation with the MSI installer.
```

With Docker Desktop version 4.91 and earlier, the MSI accepts `MSIINSTALLPERUSER` but produces an installation that later updates can't upgrade.

Use the EXE installer with the `--user` flag if you need a per-user installation.

### My installation failed, how do I find out what happened?

MSI installations may fail silently, offering little diagnostic feedback.

To debug a failed installation, run the install again with verbose logging enabled:

```powershell
msiexec /i "DockerDesktop.msi" /L*V ".\msi.log"
```

After the installation has failed, open the log file and search for occurrences of `value 3`. This is the exit code Windows Installer outputs when it has failed. Just above the line, you will find the reason for the failure.

### Why does the installer prompt for a reboot at the end of every fresh installation?

The installer prompts for a reboot because it assumes that changes have been made to the system that require a reboot to finish their configuration.

For example, if you select the WSL engine, the installer adds the required Windows features. After these features are installed, the system reboots to complete configurations so the WSL engine is functional.

You can suppress reboots by using the `/norestart` option when launching the installer from the command line:

```powershell
msiexec /i "DockerDesktop.msi" /L*V ".\msi.log" /norestart
```

### Why isn't the `docker-users` group populated when the MSI is installed with Intune or another MDM solution?

It's common for MDM solutions to install applications in the context of the system account. This means that the `docker-users` group isn't populated with the user's account, as the system account doesn't have access to the user's context.

As an example, you can reproduce this by running the installer with `psexec` in an elevated command prompt:

```powershell
psexec -i -s msiexec /i "DockerDesktop.msi"
```
The installation should complete successfully, but the `docker-users` group won't be populated.

As a workaround, you can create a script that runs in the context of the user account.

The script would be responsible for ensuring the `docker-users` group exists and populating it with the correct user.

> [!WARNING]
>
> Membership in `docker-users` grants access to the Docker daemon socket, which is equivalent to granting administrative privileges on the host. Only add users who require access to Windows containers or Hyper-V VM management. For Linux containers using the WSL 2 backend, this group membership is not required. See [Protect the Docker daemon socket](/manuals/engine/security/protect-access.md) for further information.

Here's an example script that creates the `docker-users` group if needed and adds the current user to it (requirements may vary depending on environment):

```powershell
$Group = "docker-users"
$CurrentUser = [System.Security.Principal.WindowsIdentity]::GetCurrent().Name

# Create the group if it doesn't exist
if (-not (Get-LocalGroup -Name $Group -ErrorAction SilentlyContinue)) {
    New-LocalGroup -Name $Group
}

# Add the user to the group
Add-LocalGroupMember -Group $Group -Member $CurrentUser
```

> [!NOTE]
>
> After adding a new user to the `docker-users` group, the user must sign out and then sign back in for the changes to take effect.

## MDM

Common questions about deploying Docker Desktop using mobile device management
(MDM) tools such as Jamf, Intune, or Workspace ONE.

### Why doesn't my MDM tool apply all Docker Desktop configuration settings at once?

Some MDM tools, such as Workspace ONE, may not support applying multiple
configuration settings in a single XML file. In these cases, you may need to
deploy each setting in a separate XML file.

Refer to your MDM provider's documentation for specific deployment
requirements or limitations.
