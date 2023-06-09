---
title: Best practices
description: Best practices for using Docker Desktop with WSL 2 
keywords: wsl, docker desktop, best practices
---


To get the best out of the file system performance when bind-mounting files, we recommend storing source code and other data that is bind-mounted into Linux containers, for instance with docker run -v <host-path>:<container-path>, in the Linux file system, rather than the Windows file system. You can also refer to the recommendation from Microsoft.

Linux containers only receive file change events, “inotify events”, if the original files are stored in the Linux filesystem. For example, some web development workflows rely on inotify events for automatic reloading when files have changed.
Performance is much higher when files are bind-mounted from the Linux filesystem, rather than remoted from the Windows host. Therefore avoid docker run -v /mnt/c/users:/users, where /mnt/c is mounted from Windows.
Instead, from a Linux shell use a command like docker run -v ~/my-project:/sources <my-image> where ~ is expanded by the Linux shell to $HOME.
If you have concerns about the size of the docker-desktop-data VHDX, or need to change it, take a look at the WSL tooling built into Windows.
If you have concerns about CPU or memory usage, you can configure limits on the memory, CPU, and swap size allocated to the WSL 2 utility VM.
To avoid any potential conflicts with using WSL 2 on Docker Desktop, you must uninstall any previous versions of Docker Engine and CLI installed directly through Linux distributions before installing Docker Desktop.
