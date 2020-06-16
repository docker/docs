---
title: Troubleshoot Docker Desktop Enterprise issues on Windows
description: Learn about Docker Desktop Enterprise
keywords: Docker EE, Windows, Docker Desktop, Enterprise, troubleshoot
redirect_from:
- /ee/desktop/troubleshoot/windows-issues/
---

This page contains information on how to diagnose Docker Desktop Enterprise (DDE) issues on Windows.

## Creating a diagnostics file in Docker Desktop Enterprise

Right-click the Docker icon in the system tray and select **Diagnose and Feedback** from the menu. When the **Diagnose & Feedback** window opens, it starts collecting diagnostics.

![A diagnostics file is created.](../images/diagnose-windows.png)

When the log capture is complete, select **You can find diagnostics here**. The file explorer window displays the path to the diagnostics .zip file and allows you to view the contents. Diagnostics are provided in .zip files identified by date and time.

Send your diagnostics file to your administrator for assistance.

### Creating a diagnostics file from a terminal

In some cases, it is useful to run diagnostics yourself, for instance if
Docker Desktop Enterprise cannot start.

To run diagnostics from a terminal, enter the following command from a powershell window:

```powershell
'C:\Program Files\Docker\Docker\resources\com.docker.diagnose.exe' gather
```

This command displays the information that it is gathering, and when it finishes, it displays information resembling the following example:

```powershell
Diagnostics Bundle: C:\Users\djs\AppData\Local\Temp\6CE654F6-7B17-4FC7-AAE0-CC53B73B76A2\20190115163621.zip
Diagnostics ID: 6CE654F6-7B17-4FC7-AAE0-CC53B73B76A2/20190115163621
```

The name of the diagnostics file is displayed next to “Diagnostics Bundle”  (`\Temp\6CE654F6-7B17-4FC7-AAE0-CC53B73B76A2\20190115163621.zip` in this example). This is the file that you attach to the support ticket.

### Additional Docker Desktop Enterprise troubleshooting topics

You can also find additional information about various troubleshooting topics in the [Docker Desktop for Windows community](https://docs.docker.com/docker-for-windows/troubleshoot/) documentation.
