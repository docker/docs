---
description: understand what pausing Docker Dashboard means
keywords: Docker Dashboard, manage, containers, gui, dashboard, pause, user manual
title: Pause Docker Desktop
---

To save CPU resources on your machine, Docker Desktop automatically pauses when it is not in use. 

When Docker Desktop is paused, the Linux VM running Docker Engine is paused, the current state of all your containers are saved in memory, and all processes are frozen. This reduces the CPU usage and helps you retain a longer battery life on your laptop. 

You can also manually pause Docker Desktop. From either the menu bar on Mac or the tray icon on Windows, select the Docker menu ![whale menu](../images/whale-x.svg){: .inline} and then **Pause**. You can manually resume Docker Desktop by clicking the **Resume** option in the Docker menu, or by running any Docker CLI command.

When you manually pause Docker Desktop, a paused status displays on the Docker menu and on the Docker Dashboard. You can still access the **Settings** and the **Troubleshoot** menu.
