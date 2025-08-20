---
title: Enhanced Container Isolation FAQs
linkTitle: FAQs
description: Frequently asked questions about Enhanced Container Isolation
keywords: enhanced container isolation, faq, troubleshooting, docker desktop
toc_max: 2
aliases:
 - /desktop/hardened-desktop/enhanced-container-isolation/faq/
 - /security/for-admins/hardened-desktop/enhanced-container-isolation/faq/
weight: 40
---

{{< summary-bar feature_name="Hardened Docker Desktop" >}}

This page answers common questions about Enhanced Container Isolation (ECI) that aren't covered in the main documentation.

## Do I need to change the way I use Docker when ECI is switched on?

No. ECI works automatically in the background by creating more secure containers. You can continue using all your existing Docker commands, workflows, and development tools without any changes.

## Do all container workloads work well with ECI?

Most container workloads run without issues when ECI is turned on. However, some advanced workloads that require specific kernel-level access may not work. For details about which workloads are affected, see [ECI limitations](/manuals/enterprise/security/hardened-desktop/enhanced-container-isolation/limitations.md).

## Why not just restrict usage of the `--privileged` flag?

Privileged containers serve legitimate purposes like Docker-in-Docker, Kubernetes-in-Docker, and accessing hardware devices. ECI provides a better solution by allowing these advanced workloads to run securely while preventing them from compromising the Docker Desktop VM.

## Does ECI affect container performance?

ECI has minimal impact on container performance. The only exception is containers that perform many `mount` and `umount` system calls, as these are inspected by the Sysbox runtime for security. Most development workloads see no noticeable performance difference.

## Can I override the container runtime with ECI turned on?


No. When ECI is turned on, all containers use the Sysbox runtime regardless of any `--runtime` flags:

```console
$ docker run --runtime=runc alpine echo "test"
# This still uses sysbox-runc, not runc
```

The `--runtime` flag is ignored to prevent users from bypassing ECI security by running containers as true root in the Docker Desktop VM.

## Does ECI protect containers created before turning it on?

No. ECI only protects containers created after it's turned on. Remove existing containers before turning on ECI:

```console
$ docker stop $(docker ps -q)
$ docker rm $(docker ps -aq)
```

For more details, see [Enable Enhanced Container Isolation](/manuals/enterprise/security/hardened-desktop/enhanced-container-isolation/enable-eci.md).

## Which containers does ECI protect?

ECI protection varies by container type and Docker Desktop version:

### Always protected

- Containers created with `docker run` and `docker create`
- Containers using the `docker-container` build driver

### Version dependent

- Docker Build: Protected in Docker Desktop 4.30+ (except WSL 2)
- Kubernetes: Protected in Docker Desktop 4.38+ when using the kind provisioner

### Not protected

- Docker Extensions
- Docker Debug containers
- Kubernetes with Kubeadm provisioner

For complete details, see [ECI limitations](/manuals/enterprise/security/hardened-desktop/enhanced-container-isolation/limitations.md).

## Can I mount the Docker socket with ECI turned on?

By default, no. ECI blocks Docker socket bind mounts for security. However, you can configure exceptions for trusted images like Testcontainers.

For configuration details, see [Configure Docker socket exceptions](/manuals/enterprise/security/hardened-desktop/enhanced-container-isolation/config.md).

## What bind mounts does ECI restrict?

ECI restricts bind mounts of Docker Desktop VM directories but allows host directory mounts configured in Docker Desktop Settings.
