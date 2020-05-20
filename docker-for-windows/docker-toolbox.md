---
description: Docker Desktop for Windows and Docker Toolbox
keywords: windows, alpha, beta, toolbox, docker-machine, tutorial
title: Migrate Docker Toolbox
---

This page explains how to migrate your Docker Toolbox disk image, or images if
you have them, to Docker Desktop for Windows.

## How to migrate Docker Toolbox disk images to Docker Desktop

> **Warning**
>
> Migrating disk images from Docker Toolbox _clobbers_ Docker images if they
> exist. The migration process replaces the entire VM with your previous Docker
> Toolbox data.
{: .warning }

1.  Install [qemu](https://www.qemu.org/){: target="_blank" class="_"} (a machine emulator): [https://cloudbase.it/downloads/qemu-img-win-x64-2_3_0.zip](https://cloudbase.it/downloads/qemu-img-win-x64-2_3_0.zip).
2.  Install [Docker Desktop for Windows](install.md){: target="_blank" class="_"}.
3.  Stop Docker Desktop, if running.
4.  Move your current Docker VM disk to a safe location:

    ```shell
    mv 'C:\Users\Public\Documents\Hyper-V\Virtual Hard Disks\MobyLinuxVM.vhdx' C:/<any directory>
    ```

5.  Convert your Toolbox disk image:

    ```shell
    qemu-img.exe convert 'C:\Users\<username>\.docker\machine\machines\default\disk.vmdk' -O vhdx -o subformat=dynamic -p 'C:\Users\Public\Documents\Hyper-V\Virtual Hard Disks\MobyLinuxVM.vhdx'
    ```

6.  Restart Docker Desktop (with your converted disk).

## How to uninstall Docker Toolbox

Whether or not you migrate your Docker Toolbox images, you may decide to
uninstall it. For details on how to perform a clean uninstall of Toolbox,
see [How to uninstall Toolbox](../toolbox/toolbox_install_windows.md#how-to-uninstall-toolbox){: target="_blank" class="_"}.
