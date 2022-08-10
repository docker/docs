---
description: Pause Docker Dashboard
keywords: Docker Dashboard, manage, containers, gui, dashboard, images, user manual
title: Pause Docker Desktop
---

You can pause your Docker Desktop session when you are not actively using it and save CPU resources on your machine. When you pause Docker Desktop, the Linux VM running Docker Engine is paused, the current state of all your containers are saved in memory, and all processes are frozen. This reduces the CPU usage and helps you retain a longer battery life on your laptop. You can resume Docker Desktop when you want by clicking the **Resume** option.

From the Docker menu, select![whale menu](../images/whale-x.png){: .inline} and then **Pause** to pause Docker Desktop.

Docker Desktop displays the paused status on the Docker menu and on the  **Containers**, **Images**, **Volumes**, and **Dev Environment** screens in Docker Dashboard. You can still access the **Preferences** (or **Settings** if you are a Windows user) and the **Troubleshoot** menu from the Dashboard when you've paused Docker Desktop.

Select ![whale menu](../images/whale-x.png){: .inline} then **Resume** to resume Docker Desktop.

> **Note**
>
> When Docker Desktop is paused, running any commands in the Docker CLI will automatically resume Docker Desktop.
