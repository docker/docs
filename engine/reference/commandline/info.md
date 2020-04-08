---
datafolder: engine-cli
datafile: docker_info
title: docker info
redirect_from:
  - /edge/engine/reference/commandline/info/
skip_read_time: true
---
<!--
This page is automatically generated from Docker's source code. If you want to
suggest a change to the text that appears here, open a ticket or pull request
in the source repository on GitHub:

https://github.com/docker/cli
-->
{% include cli.md datafolder=page.datafolder datafile=page.datafile %}

## Warnings about kernel support

If your operating system does not enable certain capabilities, you may see
warnings such as one of the following, when you run `docker info`:

```none
WARNING: Your kernel does not support swap limit capabilities. Limitation discarded.
```

```none
WARNING: No swap limit support
```

You can ignore these warnings unless you actually need the ability to
[limit these resources](../../../config/containers/resource_constraints.md), in which case you
should consult your operating system's documentation for enabling them.
[Learn more](../../install/linux-postinstall.md#your-kernel-does-not-support-cgroup-swap-limit-capabilities).
