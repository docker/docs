---
linkTitle: Custom kernels
title: Custom kernels on WSL 
description: Using custom kernels with Docker Desktop on WSL 2
keywords: wsl, docker desktop, custom kernel
tags: [Best practices, troubleshooting]
---

> [!WARNING]
>
> Using a custom kernel with Docker Desktop on WSL 2 is not officially supported
and may cause issues with Docker Desktop startup or operation.

Docker Desktop depends on several kernel features built into the default
WSL 2 Linux kernel distributed by Microsoft.

However, in some cases it may be necessary
to run custom kernels; Docker Desktop does not block their use, and
some users have reported success using them.

## Recommendations if you must use a custom kernel

If you choose to use a custom kernel, start
from the kernel tree distributed by Microsoft from their [official
repository](https://github.com/microsoft/WSL2-Linux-Kernel) and then add
the features you need on top of that.

Also:
- Use the same kernel version as the one distributed by the latest WSL2
release. You can find the version by running `wsl.exe --system uname -r`
in a terminal.
- Make sure that your kernel build environment includes `pahole` and
its version is properly reflected in the corresponding kernel config
(`CONFIG_PAHOLE_VERSION`).

