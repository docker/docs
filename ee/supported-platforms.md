---
title: About Docker Enterprise
description: Information about Docker Enterprise 2.1
keywords: Docker Enterprise, enterprise, enterprise edition, ee, docker ee, docker enterprise edition, lts, commercial, cs engine, commercially supported
redirect_from:
  - /enterprise/supported-platforms/
  - /cs-engine/
  - /cs-engine/1.12/
  - /cs-engine/1.12/upgrade/
  - /cs-engine/1.13/
  - /cs-engine/1.13/upgrade/
green-check: '![yes](/install/images/green-check.svg){: style="height: 14px; margin:auto;"}'
install-prefix-ee: '/install/linux/docker-ee'
---

Docker Enterprise is designed for enterprise development as well as IT teams who build, ship, and run business-critical
applications in production and at scale. Docker Enterprise is integrated, certified,
and supported to provide enterprises with the most secure container platform
in the industry. For more info about Docker Enterprise, including purchasing
options, see [Docker Enterprise](https://www.docker.com/enterprise-edition/).

> Compatibility Matrix
>
> Refer to the [Compatibility Matrix](https://success.docker.com/article/compatibility-matrix) for the latest list of supported platforms.
{: .important}

## Docker EE tiers

{% include docker_ce_ee.md %}

> Note
>
> Starting with Docker Enterprise 2.1, Docker Enterprise --- Basic, Docker Enterprise --- Standard,
> and Docker Enterprise --- Advanced are all now called Docker Enterprise.

### Docker Enterprise

With Docker Enterprise, you can deploy Docker Engine --- Enterprise
to manage your container workloads in a flexible way. You can manage workloads
on Windows, Linux, on site, or on the cloud.

Docker Enterprise has private image management, integrated image signing policies, and cluster
management with support for Kubernetes and Swarm orchestrators. It allows you to implement
node-based RBAC policies, image promotion policies, image mirroring, and
scan your images for vulnerabilities. It also has support with defined SLAs and extended
maintenance cycles for patches for up to 24 months.

### New Licensing for Docker Enterprise 

In version 18.09, the Docker Enterprise --- Engine is aware of the license applied on the system. The 
license summary is available in the `docker info` output on standalone or manager nodes.

For EE platform customers, when you license UCP, this same license is applied to the underlying 
engines in the cluster. Docker recommends platform customers use UCP to manage their license.

Standalone EE engines can be licensed using `docker engine activate`.

Offline activation of standalone EE engines can be performed by downloading the license and 
using the command `docker engine activate --license filename.lic`. 

Additionally, Docker is now distributing the CLI as a separate installation package. 
This gives Enterprise users the ability to install as many CLI packages as needed 
without using the Engine node licenses for client-only systems.

[Learn more about Docker Enterprise](/ee/index.md).


> When using Docker Enterprise
>
> IBM Power is not supported as managers or workers.
> Microsoft Windows Server is not supported as a manager. Microsoft Windows
> Server 1803 is not supported as a worker.

### Docker Certified Infrastructure

Docker Certified Infrastructure is Docker’s prescriptive approach to deploying
Docker Enterprise on a range of infrastructure choices. Each Docker
Certified Infrastructure includes a reference architecture, automation templates,
and third-party ecosystem solution briefs.

| Platform                                                                                | Docker Enterprise Edition |
|:----------------------------------------------------------------------------------------|:-------------------------:|
| [VMware](https://success.docker.com/article/certified-infrastructures-vmware-vsphere)   |  {{ page.green-check }}   |
| [Amazon Web Services](https://success.docker.com/article/certified-infrastructures-aws) |  {{ page.green-check }}   |
| [Microsoft Azure](https://success.docker.com/article/certified-infrastructures-azure)   |  {{ page.green-check }}   |
| IBM Cloud                                                                               |        Coming soon        |


## Docker Enterprise release cycles

Each Docker Enterprise release is supported and maintained for 24 months, and
receives security and critical bug fixes during this period.

The Docker API version is independent of the Docker platform version. We maintain
careful API backward compatibility and deprecate APIs and features slowly and
conservatively. We remove features after deprecating them for a period of
three stable releases. Docker 1.13 introduced improved interoperability
between clients and servers using different API versions, including dynamic
feature negotiation.

## Upgrades and support

If you're a Docker DDC or CS Engine customer, you don't need to upgrade to
Docker Enterprise to continue to get support. We will continue to support customers
with valid subscriptions whether the subscription covers Docker EE or
Commercially Supported Docker. You can choose to stay with your current
deployed version, or you can upgrade to the latest Docker EE version. For
more info, see [Scope of Coverage and Maintenance
Lifecycle](https://success.docker.com/Policies/Scope_of_Support).

## Where to go next

- [Install Docker](/engine/installation/index.md)
- [Get Started with Docker](/get-started/index.md)
