---
description: Multi-CPU Architecture Support
keywords: mac, Multi-CPU architecture support
redirect_from:
- /mackit/multi-arch/
title: Leverage multi-CPU architecture support
notoc: true
---

Docker for Mac provides `binfmt_misc` multi architecture support, so you can run
containers for different Linux architectures, such as `arm`, `mips`, `ppc64le`,
and even `s390x`.

This should just work without any configuration, but the containers you run need
to have the appropriate `qemu` binary inside before you can do
this. (See <a href="http://wiki.qemu.org/" target="_blank">QEMU</a> for more
information.)

So, you can run a container that already has this set up, like the <a
href="https://resin.io/how-it-works/" target="_blank">resin</a> arm builds:

```
$ docker run resin/armv7hf-debian uname -a

Linux 7ed2fca7a3f0 4.1.12 #1 SMP Tue Jan 12 10:51:00 UTC 2016 armv7l GNU/Linux

$ docker run justincormack/ppc64le-debian uname -a

Linux edd13885f316 4.1.12 #1 SMP Tue Jan 12 10:51:00 UTC 2016 ppc64le GNU/Linux

```

Running containers pre-configured with `qemu` has the advantage that you can use
these to do builds `FROM`, so you can build new Multi-CPU architecture packages.

Alternatively, you can bind mount in the `qemu` static binaries to any
cross-architecture package, such as the semi-official ones using a script like
this one [https://github.com/justincormack/cross-docker](https://github.com/justincormack/cross-docker). (See the README at the
given link for details on how to use the script.)
