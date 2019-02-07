---
description: Multi-CPU Architecture Support
keywords: mac, Multi-CPU architecture support
redirect_from:
- /mackit/multi-arch/
title: Leverage multi-CPU architecture support
notoc: true
---

Docker Desktop for Mac provides `binfmt_misc` multi architecture support, so you can run
containers for different Linux architectures, such as `arm`, `mips`, `ppc64le`,
and even `s390x`.

This does not require any special configuration in the container itself as it uses
<a href="http://wiki.qemu.org/" target="_blank">qemu-static</a> from the Docker for
Mac VM.

You can run an ARM container, like the <a href="https://www.balena.io/what-is-balena/" target="_blank">
balena</a> arm builds:

```
$ docker run balenalib/armv7hf-debian uname -a

Linux 3d3ffca44f6e 4.9.125-linuxkit #1 SMP Fri Sep 7 08:20:28 UTC 2018 armv7l GNU/Linux

$ docker run justincormack/ppc64le-debian uname -a

Linux edd13885f316 4.1.12 #1 SMP Tue Jan 12 10:51:00 UTC 2016 ppc64le GNU/Linux

```

Multi architecture support makes it easy to build <a href="https://blog.docker.com/2017/11/multi-arch-all-the-things/" target="_blank">
multi architecture Docker images</a> or experiment with ARM images and binaries
from your Mac.
