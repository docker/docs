---
title: Docker Engine release notes
description: Learn about the new features, bug fixes, and breaking changes for Docker Engine CE and EE
keywords: docker, docker engine, ee, ce, whats new, release notes
toc_min: 1
toc_max: 2
---

This document describes the latest changes, additions, known issues, and fixes
for Docker Engine Enterprise Edition (Docker EE) and Community Edition (CE)

Docker EE is functionally equivalent to the corresponding Docker CE that
it references. However, Docker EE also includes back-ported fixes
(security-related and priority defects) from the open source. It incorporates
defect fixes that you can use in environments where new features cannot be
adopted as quickly for consistency and compatibility reasons.


## 18.09 (2018-11-XX)

### New features for Docker Engine EE and CE

* Docker CE-EE node activate, which enables a user to apply a license to a CE binary and have it seamlessly upgrade to the EE binary
* Docker Build architecture overhaul, which provides enhancements to `docker build`
* Integrate containerd for better Windows support and future features

### New features for Docker Engine EE 
* Improved FIPS implementation to include Windows support.
* Docker Content Trust Enforcement for the Enterprise Engine. This prevents Docker Engine from running containers not signed by a specific organization.

### Bug fixes


### Known issues

**Depreciation Notice**

As of EE 2.2, Docker Engine will no longer support Devicemapper as a storage driver.


## Earlier versions

- [Docker Enterprise Engine 18.03 and earlier release notes](/ee/engine/release-notes.md) 
- [Docker CE release notes for 18.06 and earlier](/release-notes/docker-ce/index.md)
