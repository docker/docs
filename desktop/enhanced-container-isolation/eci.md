---
description: Enhanced Container Isolation - benefits, why use it, how it differs to Docker rootless, who it is for
keywords: containers, rootless, security, sysbox, runtime
title: What is Enhanced Container Isolation?
---

What it is 

what the benefits of it are

Who is it for

how does it work 

- Who gets this ? (e.g. currently developers in Docker Business customers, requires authentication, etc)
- How the feature works under the hood
    - enables Sysbox under the hood, ensuring containers run using the Linux user namespace and are not root in the VM, etc
- Dive in on details
    - e.g. when enabled, Desktop uses Sysbox runtime by default, for all containers. requires an Apply & restart, etc
- Why this approach is advantageous as compared to traditional ‘rootless Docker' or ‘rootless mode’ in “other products”
    - workload compatibility, ease of use, etc. dive in on why Sysbox is awesome for both security and workloads
- Admins can lock in the use of the ‘Enhanced container isolation’ mode within their org via the ‘Admin Controls’ feature <link to Admin Controls docs>