---
linkTitle: Best practices
title: WSL 2 best practices for Docker Desktop on Windows 
description: Best practices for using Docker Desktop with WSL 2
keywords: wsl 2, docker desktop, best practices, Windows Subsystem for Linux, Docker Desktop Windows performance
tags: [Best practices]
aliases:
- /desktop/wsl/best-practices/
---

This page covers recommendations when running Docker Desktop on Windows using WSL 2, including version requirements and file system performance.

## Keep WSL up to date

Always use the latest version of WSL. 

At a minimum you must use WSL version 2.1.5, otherwise Docker Desktop may not work as expected. Testing, development, and documentation is based on the newest kernel versions. Older versions of WSL can cause:
- Docker Desktop to hang periodically or when upgrading
- Deployment via SCCM to fail
- The `vmmem.exe` to consume all memory 
- Network filter policies to be applied globally, not to specific objects
- GPU failures with containers

## Optimise file system performance with bind mounts

To get the best out of the file system performance when bind-mounting files, store source code and other data that is bind-mounted into Linux containers. For instance, use `docker run -v <host-path>:<container-path>` in the Linux file system, rather than the Windows file system. You can also refer to [Microsoft's recommendation](https://learn.microsoft.com/en-us/windows/wsl/compare-versions).

Linux containers only receive file change events, “inotify events”, if the original files are stored in the Linux filesystem. For example, some web development workflows rely on inotify events for automatic reloading when files have changed.

Performance is much higher when files are bind-mounted from the Linux filesystem, rather than accessed from the Windows host filesystem. Therefore avoid `docker run -v /mnt/c/users:/users` where `/mnt/c` is mounted from Windows.

Instead, from a Linux shell use a command like `docker run -v ~/my-project:/sources <my-image>` where `~` is expanded by the Linux shell to `$HOME`.

## Limit CPU and memory usage

If you have concerns about CPU or memory usage, configure limits on the memory, CPU, and swap size allocated to the [WSL 2 utility VM](https://learn.microsoft.com/en-us/windows/wsl/wsl-config#global-configuration-options-with-wslconfig).
