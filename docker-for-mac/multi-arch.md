---
description: Multi-CPU Architecture Support
keywords: mac, Multi-CPU architecture support
redirect_from:
- /mackit/multi-arch/
title: Leverage multi-CPU architecture support
notoc: true
---
Docker images can support multiple architectures, which means that a single
image may contain variants for different architectures, and sometimes for different
operating systems, such as Windows.

When running an image with multi-architecture support, `docker` will
automatically select an image variant which matches your OS and architecture.

Most of the official images on Docker Hub provide a [variety of architectures](https://github.com/docker-library/official-images#architectures-other-than-amd64).
For example, the `busybox` image supports `amd64`, `arm32v5`, `arm32v6`,
`arm32v7`, `arm64v8`, `i386`, `ppc64le`, and `s390x`. When running this image
on an `x86_64` / `amd64` machine, the `x86_64` variant will be pulled and run,
which can be seen from the output of the `uname -a` command that's run inside
the container:

```bash
$ docker run busybox uname -a

Linux 82ef1a0c07a2 4.9.125-linuxkit #1 SMP Fri Sep 7 08:20:28 UTC 2018 x86_64 GNU/Linux
```

**Docker Desktop for Mac** provides `binfmt_misc` multi-architecture support,
which means you can run containers for different Linux architectures
such as `arm`, `mips`, `ppc64le`, and even `s390x`.

This does not require any special configuration in the container itself as it uses
<a href="http://wiki.qemu.org/" target="_blank">qemu-static</a> from the **Docker for
Mac VM**. Because of this, you can run an ARM container, like the `arm32v7` or `ppc64le`
variants of the busybox image:

### arm32v7 variant
```bash
$ docker run arm32v7/busybox uname -a

Linux 9e3873123d09 4.9.125-linuxkit #1 SMP Fri Sep 7 08:20:28 UTC 2018 armv7l GNU/Linux
```

### ppc64le variant
```bash
$ docker run ppc64le/busybox uname -a

Linux 57a073cc4f10 4.9.125-linuxkit #1 SMP Fri Sep 7 08:20:28 UTC 2018 ppc64le GNU/Linux
```

Notice that this time, the `uname -a` output shows `armv7l` and
`ppc64le` respectively.

Multi-architecture support makes it easy to build <a href="https://blog.docker.com/2017/11/multi-arch-all-the-things/" target="_blank">multi-architecture Docker images</a> or experiment with ARM images and binaries from your Mac.

