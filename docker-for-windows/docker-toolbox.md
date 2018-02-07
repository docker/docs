---
description: Docker for Windows and Docker Toolbox
keywords: windows, alpha, beta, toolbox, docker-machine, tutorial
title: Docker for Windows vs. Docker Toolbox
---

If you already have an installation of Docker Toolbox, read these topics
to learn how to migrate to Docker for Windows.

## Migrating from Docker Toolbox to Docker for Mac

Docker for Windows does not propose Toolbox image migration as part of the
Docker for Windows installer since version 18.01.0.  You can migrate existing
Docker Toolbox images with the scripts described below. (Note that this
migration cannot merge images from both Docker and Toolbox: any existing Docker
image are *replaced* by the Toolbox images.)

To run these instructions you need to now how to run shell commands in a
terminal. You also need a working `qemu-img.exe`; download it from https://cloudbase.it/downloads/qemu-img-win-x64-2_3_0.zip. 

Then:
1. After installing Docker for Windows, stop Docker for Windows if you have started it.
2. Move or delete your current Docker VM disk. (Keep in mind the Toolbox migration does not merge any image, it just replaces the entire VM with your previous Toolbox data)
 `mv 'C:\Users\Public\Documents\Hyper-V\Virtual Hard Disks\MobyLinuxVM.vhdx' C:/...`
3. Convert your Toolbox image: 
`qemu-img.exe convert 'C:\Users\<username>\.docker\machine\machines\default\disk.vmdk' -O vhdx -o subformat=dynamic -p 'C:\Users\Public\Documents\Hyper-V\Virtual Hard Disks\MobyLinuxVM.vhdx'`
4. Restart Docker for Windows. It will start using the converted disk at `C:\Users\Public\Documents\Hyper-V\Virtual Hard Disks\MobyLinuxVM.vhdx`

Optionally, if you are done with Docker Toolbox, you may fully uninstall
it, see below.

## How do I uninstall Docker Toolbox?

You might decide that you do not need Toolbox now that you have Docker for
Windowsd, and want to uninstall it. For details on how to perform a clean
uninstall of Toolbox, see [How to uninstall Toolbox](/toolbox/toolbox_install_windows/#how-to-uninstall-toolbox).
