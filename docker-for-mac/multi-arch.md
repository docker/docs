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

You can run an ARM container, like the <a href="https://resin.io/how-it-works/" target="_blank">
resin</a> arm builds:

```
$ docker run resin/armv7hf-debian uname -a

Linux 7ed2fca7a3f0 4.1.12 #1 SMP Tue Jan 12 10:51:00 UTC 2016 armv7l GNU/Linux

$ docker run justincormack/ppc64le-debian uname -a

Linux edd13885f316 4.1.12 #1 SMP Tue Jan 12 10:51:00 UTC 2016 ppc64le GNU/Linux

```

Multi architecture support makes it easy to build <a href="https://blog.docker.com/2017/11/multi-arch-all-the-things/" target="_blank">
multi architecture Docker images</a> or experiment with ARM images and binaries
from your Mac.
