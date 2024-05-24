---
title: Analytics
description: Understand how the MSI installer collects analytics and how to opt out
keywords: msi, deploy, docker desktop, analyitcs
---

The MSI installer collects data to better understand user behaviour and to improve the user experience by identifying and addressing issues or optimizing popular features.

## What data is collected?

The table below outlines the properties that are sent in the analytics payload.

| Name                   |   Type          |
| :-------------- | :-------------- |
| action | string |
| adminSettings | boolean |
| allowedOrg | boolean |
| alwaysRunService | boolean |
| appBuildNumber | string |
| appMajorVersion | string |
| appMinorVersion | string |
| appPatchVersion | string |
| appVersionName | string |
| disableWindowsContainers | boolean |
| enableDesktopShortcut | boolean |
| engine | string |
| existingExeInstallPresent | boolean |
| hypervDefaultDataRoot | boolean |
| customInstallFolder | boolean |
| os | string |
| osBuildVersion | string |
| osEdition | string |
| osLanguage | string |
| removeExistingInstall | boolean |
| sessionId | string |
| windowsContainersDefaultDataRoot | boolean |
| wslDefaultDataRoot | boolean |

## How to opt-out

### From the command line

When you install Docker Desktop from the command line, use the `DISABLEANALYTICS` property.

```powershell
msiexec /i "win\msi\bin\en-US\DockerDesktop.msi" /L*V ".\msi.log" DISABLEANALYTICS=1
```

### From the UI 

When you install Docker Desktop from the default installer GUI, select the **Disable analytics** checkbox located on the bottom-left corner of the **Welcome** dialog.

### Persistence

If you decide to disable analytics for an installation, your choice is persisted in the registry and honoured across future upgrades and uninstalls.

However, the key is removed when Docker Desktop is uninstalled and must be configured again via one of the preivous methods.

The registry key is as follows:

```powershell
SOFTWARE\Docker Inc.\Docker Desktop\DisableMsiAnalytics
```

When analytics have been disabled, this key has a value of `1`. 