---
title: About Docker EE
description: Information about Docker Enterprise Edition Platform 2.0
keywords: enterprise, enterprise edition, ee, docker ee, docker enterprise edition, lts, commercial, cs engine
redirect_from:
  - /enterprise/supported-platforms/
---

Docker Enterprise Edition (*Docker EE*) is designed for enterprise
development and IT teams who build, ship, and run business-critical
applications in production and at scale. Docker EE is integrated, certified,
and supported to provide enterprises with the most secure container platform
in the industry. For more info about Docker EE, including purchasing
options, see [Docker Enterprise Edition](https://www.docker.com/enterprise-edition/).

<!-- This is populated by logic in js/archive.js -->
<p id="ee-version-div"></p>

The free Docker products continue to be available as the Docker Community
Edition (*Docker CE*).

## Supported platforms

The following table shows all of the platforms that are available for Docker EE.
Each link in the first column takes you to the installation
instructions for the corresponding platform. Docker EE is an integrated,
supported, and certified container platform for the listed cloud providers and
operating systems.

{% include docker_platform_matrix.md %}

## Docker EE feature tiers

Docker EE is available in three tiers:

-  **Basic:** The Docker platform for certified infrastructure, with support
   from Docker Inc. and certified containers and plugins from Docker Store.
-  **Standard:** Adds advanced image and container management, LDAP/AD user
   integration, and role-based access control. Together, these features
   comprise Docker Enterprise Edition.
-  **Advanced:** Adds
   [Docker Security Scanning](https://blog.docker.com/2016/05/docker-security-scanning/)
   and continuous vulnerability monitoring.

## Docker Enterprise Edition release cycles

Docker EE is released quarterly. Releases use a time-based versioning
scheme, so for example, Docker EE version 17.03 was released
in March 2017. For schedule details, see
[Time-based release schedule](/engine/installation/#time-based-release-schedule).

Each Docker EE release is supported and maintained for one year and
receives security and critical bug fixes during this period.

The Docker API version is independent of the Docker platform version. The API
version doesn't change from Docker 1.13.1 to Docker 17.03. We maintain
careful API backward compatibility and deprecate APIs and features slowly and
conservatively. We remove features after deprecating them for a period of
three stable releases. Docker 1.13 introduced improved interoperability
between clients and servers using different API versions, including dynamic
feature negotiation.

## Upgrades and support

If you're a Docker DDC or CS Engine customer, you don't need to upgrade to
Docker EE to continue to get support. We will continue to support customers
with valid subscriptions whether the subscription covers Docker EE or
Commercially Supported Docker. You can choose to stay with your current
deployed version, or you can upgrade to the latest Docker EE version. For
more info, see [Scope of Coverage and Maintenance
Lifecycle](https://success.docker.com/Policies/Scope_of_Support).

## Where to go next

- [Install Docker](/engine/installation/index.md)
- [Get Started with Docker](/get-started/index.md)
