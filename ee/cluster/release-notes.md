---
title: Docker Cluster release notes
description: Learn about the new features, bug fixes, and breaking changes for Docker Cluster.
keywords: cluster, whats new, release notes
---

>{% include enterprise_label_shortform.md %}

This page provides information about Docker Cluster versions. 

# Version 1


## 1.2.0
(2019-10-8)

### Features

- Added new env variable type which allows users to supply cluster variable values as an environment variable (DCIS-509)

### Fixes

- Fixed an issue where errors in the cluster commands would return exit code of 0 (DCIS-508)

- New error message displayed when a docker login is required: 
```Checking for licenses on Docker Hub
   Error: no Hub login info found; please run 'docker login' first
```

## Version 1.1.0 
(2019-09-03)

### What's new

* Support for Azure cloud provider
* Support for RHEL 8 operating system

### Bug Fixes

* Container file permissions on linux (DCIS-346)
* License check failure with non-Docker subscriptions (DCIS-465)

## Version 1.0.1 
(2019-07-19)

### What's new

* Support for SLES 12.4
* Support for Windows Server 2016

## Version 1.0 
(2019-06-25)

First major release.
