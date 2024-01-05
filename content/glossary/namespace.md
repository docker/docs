---
title: Namespace
id: namespace
short_description: >
  A Linux kernel feature that isolates and virtualizes system resources.
---

A [Linux namespace](https://man7.org/linux/man-pages/man7/namespaces.7.html) is a Linux kernel feature that isolates and virtualizes system resources. Processes which are restricted to a namespace can only interact with resources or processes that are part of the same namespace. Namespaces
are an important part of Docker's isolation model. Namespaces exist for each type of resource, including `net` (networking), `mnt` (storage), `pid` (processes), `uts` (hostname control), and `user` (UID mapping). For more information about namespaces, see [Docker run reference](/engine/reference/run/) and [Isolate containers with a user namespace](/engine/security/userns-remap/).
