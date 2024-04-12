---
title: Installation scenarios
description: Find example code samples for different installation scenarios 
keywords: install, msi, deploy, windows, docker desktop, examples
---

## Installation scenarios 

This section covers command line installations of Docker Desktop using PowerShell.

Interactive installations, without specifying `/quiet` or `/qn`, display the user interface and let the user select their own properties. TODO. 

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
> Some useful tips to remember when creating a value that expects a JSON string as it’s value:
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
