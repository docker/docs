---
title: Docker Engine release notes
description: Learn about the new features, bug fixes, and breaking changes for Docker Engine CE and EE
keywords: docker, docker engine, ee, ce, whats new, release notes
toc_min: 1
toc_max: 2
---

This document describes the latest changes, additions, known issues, and fixes
for Docker Engine Enterprise Edition (Docker EE) and Community Edition (CE)

Docker EE is a superset of all the features in Docker CE. It incorporates defect fixes 
that you can use in environments where new features cannot be adopted as quickly for 
consistency and compatibility reasons.


## 18.09 (2018-11-XX)

### New features for Docker Engine EE and CE

* Docker CE-EE node activate, which enables a user to apply a license to a CE binary and have it seamlessly upgrade 
to the EE binary
* Docker Build architecture overhaul, which integrates buildKit to provide enhancements to `docker build`
* Integrated containerd runtime to serve as a foundation for Docker Engine

### New features for Docker Engine EE 
* Improved FIPS implementation to include Windows support.
* Docker Content Trust Enforcement for the Enterprise Engine. This prevents Docker Engine from running containers not signed by a specific organization.

### Bug fixes


### Known issues

**Depreciation Notice**

Docker EE 2.1 Platform release, will serve to deprecate support of Device Mapper in a future release. It will continue to be supported at this time, but 
support will be removed in a future release. Docker will continue to support Device Mapper for for existing EE 2.0 and 2.1 customers.
Please contact Sales for more information.

Docker reocmmends that existing customers migrate to using Overlay2 for the storage driver. The Overlay2 storage driver is now the default
for Docker engine implementations.

## Earlier versions

- [Docker Enterprise Engine 18.03 and earlier release notes](/ee/engine/release-notes.md) 
- [Docker CE release notes for 18.06 and earlier](/release-notes/docker-ce/index.md)
