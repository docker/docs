---
title: Custom kernels on WSL 
description: Using custom kernels with Docker Desktop on WSL 2
keywords: wsl, docker desktop, custom kernel
tags: [Best practices, troubleshooting]
aliases:
- /desktop/wsl/custom-kernels/
---

Docker Desktop depends on several kernel features built into the default
WSL 2 Linux kernel distributed by Microsoft. Because of that, **using a
custom kernel with Docker Desktop on WSL 2 is not officially supported**
and may cause issues with Docker Desktop startup or operation.

At the same time, we understand that in some cases it may be necessary
to run custom kernels; Docker Desktop does not block their use, and
some users have reported success using them.

If you choose to use a custom kernel, we suggest starting
from the kernel tree distributed by Microsoft at their [official
repository](https://github.com/microsoft/WSL2-Linux-Kernel) and adding
the features you need on top of that.

We suggest to:
- Use the same kernel version as the one distributed by the latest WSL2
release (you can find the version by running `wsl.exe --system uname -r`
in a terminal).
- Start from the default kernel configuration as provided by Microsoft
in their [repository](https://github.com/microsoft/WSL2-Linux-Kernel)
and add the features you need on top of that.
- Make sure that your kernel build environment includes `pahole` and
its version is properly reflected in the corresponding kernel config
(`CONFIG_PAHOLE_VERSION`).

