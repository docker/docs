---
title: Overview
description: Docker Enterprise product information
keywords: Docker Enterprise, enterprise, enterprise edition, ee, docker ee, docker enterprise edition, lts, commercial, cs engine, commercially supported
redirect_from:
  - /enterprise/supported-platforms/
  - /cs-engine/
  - /cs-engine/1.12/
  - /cs-engine/1.12/upgrade/
  - /cs-engine/1.13/
  - /cs-engine/1.13/upgrade/
green-check: '![yes](/images/green-check.svg){: style="height: 14px; margin:auto;"}'
install-prefix-ee: '/install/linux/docker-ee'
---

>{% include enterprise_label_shortform.md %}

Docker Enterprise is designed for enterprise development as well as IT teams who build, share, and run business-critical
applications at scale in production. Docker Enterprise is an integrated container platform that includes
Docker Desktop Enterprise, a secure image registry, advanced management control plane, and Docker Engine - Enterprise.
Docker Engine - Enterprise is a certified and supported container runtime that is also available as a standalone
solution  to provide enterprises with the most secure container engine in the industry. For more information
about Docker Enterprise and Docker Engine - Enterprise, including purchasing options,
see [Docker Enterprise](https://www.docker.com/enterprise-edition/).

> Compatibility Matrix
>
> Refer to the [Compatibility Matrix](https://success.docker.com/article/compatibility-matrix)
> for the latest list of supported platforms.
{: .important}

## Docker Enterprise products

{% include docker_ee.md %}

> Note
>
> Starting with Docker Enterprise 2.1, Docker Enterprise - Basic is now Docker Engine - Enterprise,
> and both Docker Enterprise - Standard and Docker Enterprise - Advanced are now called Docker Enterprise.

### Docker Enterprise

With Docker Enterprise, you can manage container workloads on Windows, Linux, on site, or on the cloud
in a flexible way.

Docker Enterprise has private image management, integrated image signing policies, and cluster
management with support for Kubernetes and Swarm orchestrators. It allows you to implement
node-based RBAC policies, image promotion policies, image mirroring, and
scan your images for vulnerabilities. It also has support with defined SLAs and extended
maintenance cycles for patches for up to 24 months.

### New licensing for Docker Enterprise

Starting in version 18.09, Docker Enterprise is aware of the license applied on
the system. The license summary is available in the `docker info` output on
standalone or manager nodes.

For Docker Enterprise customers, when you license Universal Control Plane
(UCP), this same license is applied to the underlying engines in the cluster.
Docker recommends that Enterprise customers use UCP to manage their license.

Docker distributing the CLI as a separate installation package. This gives Docker
Enterprise users the ability to install as many CLI packages as needed without
using the Engine node licenses for client-only systems.

[Learn more about Docker Enterprise](index.md).


> When using Docker Enterprise
> Microsoft Windows Server is not supported as a manager. Microsoft Windows
> Server 1803 is not supported as a worker.

### Docker Certified Infrastructure

Docker Certified Infrastructure is Docker’s prescriptive approach to deploying Docker Enterprise
on a variety of infrastructures. Each Docker Certified Infrastructure option includes a reference architecture,
a CLI plugin for automated deployment and configuration, and third-party ecosystem solution briefs.

| Platform  | Docker Enterprise support |
:----------------------------------------------------------------------------------------|:-------------------------:|
| [Amazon Web Services](cluster/aws.md) |  {{ page.green-check }}   |
| [Azure](cluster/azure.md) |  {{ page.green-check }}   |
| VMware  |  coming soon  |

## Docker Enterprise release cycles

Each Docker Enterprise release is supported and maintained for 24 months, and
receives security and critical bug fixes during this period.

The Docker API version is independent of the Docker version. We maintain
careful API backward compatibility and deprecate APIs and features slowly and
conservatively. We remove features after deprecating them for a period of
three stable releases. Docker 1.13 introduced improved interoperability
between clients and servers using different API versions, including dynamic
feature negotiation.

## Upgrades and support
Docker supports Docker Enterprise minor releases for 24 months. Upgrades to the
latest minor release of Docker Enterprise are not required, however we
recommend staying on the latest maintenance release of the supported minor
release you are on. Please see [Maintenance
Lifecycle](https://success.docker.com/article/maintenance-lifecycle) for more
details on EOL of minor and major versions of Docker Enterprise.

## Where to go next

- [Install Docker Engine - Enterprise for RHEL](docker-ee/rhel.md)
- [Install Docker Engine - Enterprise for Ubuntu](docker-ee/ubuntu.md)
- [Install Docker Engine - Enterprise for Windows Server](docker-ee/windows/docker-ee.md)
