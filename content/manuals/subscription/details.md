---
title: Docker subscriptions and features
linkTitle: Subscriptions and features
description: Learn about Docker subscription tiers and their key features
keywords: subscription, personal, pro, team, business, features, docker subscription
aliases:
- /subscription/core-subscription/details/
weight: 10
---

Docker subscriptions provide licensing for commercial use of Docker products and include access to Docker's complete development platform:

- [Docker Desktop](../desktop/_index.md): The industry-leading container-first
  development solution that includes, Docker Engine, Docker CLI, Docker Compose,
  Docker Build/BuildKit, and Kubernetes.
- [Docker Hub](../docker-hub/_index.md): The world's largest cloud-based
  container registry.
- [Docker Build Cloud](../build-cloud/_index.md): Powerful cloud-based builders that accelerate build times by up to 39x.
- [Docker Scout](../scout/_index.md): Tooling for software supply chain security
  that lets you quickly assess image health and accelerate security improvements.
- [Testcontainers Cloud](https://testcontainers.com/cloud/docs):  Container-based
  testing automation that provides faster tests, a unified developer experience,
  and more.
- [Docker Hardened Images](/manuals/dhi/_index.md): Minimal, secure, and
  production-ready container base and application images maintained by Docker.

Choose the subscription that fits your needs, from individual developers to large enterprises.

> [!NOTE]
>
> Legacy Docker plans apply to subscribers who last purchased or renewed before December 10, 2024. These subscribers keep their current subscription and pricing until their next renewal on or after December 10, 2024.

## Subscriptions

{{< tabs >}}
{{< tab name="Docker subscription" >}}

## Docker Personal

**Docker Personal** is ideal for open source communities, individual developers, education, and small businesses.

Docker Personal includes:

- Essential Docker tools at no cost
- 1 Docker Scout repository with vulnerability analysis
- Unlimited public Docker Hub repositories
- 200 pulls per 6 hours for authenticated users
- 7-day trials for Docker Build Cloud and Testcontainers Cloud

Docker Personal users who want to continue using Docker Build Cloud or Docker
Testcontainers Cloud after their trial can upgrade to a Docker Pro subscription at any
time.

All unauthenticated users, including unauthenticated Docker Personal users, get
100 pulls per 6 hours per IPv4 address or IPv6 /64 subnet.

For a list of features available in each tier, see [Docker Pricing](https://www.docker.com/pricing/).

## Docker Pro

**Docker Pro** is ideal for individual developers who need full access to Docker's development platform.

Docker Pro includes:

- Full access to all Docker tools
- 200 Docker Build Cloud minutes per month, Docker Build Cloud minutes do not
rollover month to month
- 2 Docker Scout repositories with vulnerability analysis
- 100 Testcontainers Cloud runtime minutes per month, Testcontainers Cloud runtime minutes do not rollover month to month
- No Docker Hub pull rate limits

For a list of features available in each tier, see [Docker
Pricing](https://www.docker.com/pricing/).

## Docker Team

**Docker Team** is ideal for development teams that need collaboration and security features.

Docker Team includes:

- 500 Docker Build Cloud minutes per month, Docker Build Cloud minutes do not
rollover month to month
- Unlimited Docker Scout repositories with vulnerability analysis
- 500 Testcontainers Cloud runtime minutes per month, Testcontainers Cloud runtime minutes do not rollover month to month
- No Docker Hub pull rate limits
- Advanced collaboration tools including organization management, [Role Based Access Control
(RBAC)](/security/for-admins/roles-and-permissions/), [activity logs](/admin/organization/activity-logs/), and more

For a list of features available in each tier, see [Docker
Pricing](https://www.docker.com/pricing/).

## Docker Business

**Docker Business** is ideal for enterprises that need centralized management and advanced security.

Docker Business includes:

- 1500 Docker Build Cloud minutes per month, Docker Build Cloud minutes do not
rollover month to month
- Unlimited Docker Scout repositories with vulnerability analysis
- 1500 Testcontainers Cloud runtime minutes per month, Testcontainers Cloud runtime minutes do not rollover month to month
- No Docker Hub pull rate limits
- Enterprise security features:
  - [Hardened Docker Desktop](/manuals/enterprise/security/hardened-desktop/_index.md)
  - [Image Access
  Management](/manuals/enterprise/security/hardened-desktop/image-access-management.md)
  which lets admins control what content developers can access
  - [Registry Access
  Management](/manuals/enterprise/security/hardened-desktop/registry-access-management.md)
  which lets admins control what registries developers can access
  - [Company layer](/admin/company/) to manage multiple organizations and settings
  - [Single sign-on](/security/for-admins/single-sign-on/)
  - [System for Cross-domain Identity
  Management](/security/for-admins/provisioning/scim/)

For a list of features available in each tier, see [Docker
Pricing](https://www.docker.com/pricing/).

{{< /tab >}}
{{< tab name="Legacy Docker plans" >}}

> [!IMPORTANT]
>
> Legacy Docker plans apply to subscribers who last purchased or renewed before December 10, 2024. These subscribers keep their current subscription and pricing until their next renewal on or after December 10, 2024.

If you have a legacy subscription, you'll automatically upgrade to the new Docker subscription model when you renew. The new plans provide access to all Docker tools with increased limits and additional features.

## Legacy Docker Pro

**Legacy Docker Pro** enables individual developers to get more control of their
development environment and provides an integrated and reliable developer
experience. It reduces the amount of time developers spend on mundane and
repetitive tasks and empowers developers to spend more time creating value for
their customers.

Legacy Docker Pro includes:
- Unlimited public repositories
- Unlimited [Scoped Access Tokens](/security/access-tokens/)
- Unlimited [collaborators](/docker-hub/repos/manage/access/#collaborators) for public repositories at no cost per month.
- Access to [Legacy Docker Scout Free](#legacy-docker-scout-free) to get started with software supply chain security.
- Unlimited private repositories
- 5000 image [pulls per day](/manuals/docker-hub/usage/pulls.md)
- [Auto Builds](/docker-hub/builds/) with 5 concurrent builds
- 300 [Vulnerability Scans](/docker-hub/vulnerability-scanning/)

For a list of features available in each legacy tier, see [Legacy Docker Pricing](https://www.docker.com/legacy-pricing/).

### Upgrade your Legacy Docker Pro subscription

When you upgrade your Legacy Docker Pro subscription to a Docker Pro subscription, your subscription includes the following changes:

- Docker Build Cloud build minutes increased from 100/month to 200/month and no monthly fee. Docker Build Cloud minutes do not rollover month to month.
- 2 included repositories with continuous vulnerability analysis in Docker Scout.
- 100 Testcontainers Cloud runtime minutes are now included for use either in Docker Desktop or for CI. Testcontainers Cloud runtime minutes do not rollover month to month.
- Docker Hub image pull rate limits are removed.

For a list of features available in each tier, see [Docker Pricing](https://www.docker.com/pricing/).

## Legacy Docker Team

**Legacy Docker Team** offers capabilities for collaboration, productivity, and
security across organizations. It enables groups of developers to unlock the
full power of collaboration and sharing combined with essential security
features and team management capabilities. A Docker Team subscription includes
licensing for commercial use of Docker components including Docker Desktop and
Docker Hub.

Legacy Docker Team includes:
- Everything included in legacy Docker Pro
- Unlimited teams
- [Auto Builds](/docker-hub/builds/) with 15 concurrent builds
- Unlimited [Vulnerability Scanning](/docker-hub/vulnerability-scanning/)
- 5000 image [pulls per day](/manuals/docker-hub/usage/pulls.md) for each team member

There are also advanced collaboration and management tools, including organization and team management with [Role Based Access Control (RBAC)](/security/for-admins/roles-and-permissions/), [activity logs](/admin/organization/activity-logs/), and more.

For a list of features available in each legacy tier, see [Legacy Docker Pricing](https://www.docker.com/legacy-pricing/).

### Upgrade your Legacy Docker Team subscription

When you upgrade your Legacy Docker Team subscription to a Docker Team subscription, your subscription includes the following changes:

- Instead of paying an additional per-seat fee, Docker Build Cloud is now available to all users in your Docker subscription.
- Docker Build Cloud build minutes increase from 400/mo to 500/mo. Docker Build Cloud minutes do not rollover month to month.
- Docker Scout now includes unlimited repositories with continuous vulnerability analysis, an increase from 3.
- 500 Testcontainers Cloud runtime minutes are now included for use either in Docker Desktop or for CI. Testcontainers Cloud runtime minutes do not rollover month to month.
- Docker Hub image pull rate limits are removed.
- The minimum number of users is 1 (lowered from 5).

For a list of features available in each tier, see [Docker Pricing](https://www.docker.com/pricing/).

## Legacy Docker Business

**Legacy Docker Business** offers centralized management and advanced security features
for enterprises that use Docker at scale. It empowers leaders to manage their
Docker development environments and speed up their secure software supply chain
initiatives. A Docker Business subscription includes licensing for commercial
use of Docker components including Docker Desktop and Docker Hub.

Legacy Docker Business includes:
- Everything included in legacy Docker Team
- [Hardened Docker Desktop](/manuals/enterprise/security/hardened-desktop/_index.md)
- [Image Access Management](/manuals/enterprise/security/hardened-desktop/image-access-management.md) which lets admins control what content developers can access
- [Registry Access Management](/manuals/enterprise/security/hardened-desktop/registry-access-management.md) which lets admins control what registries developers can access
- [Company layer](/admin/company/) to manage multiple organizations and settings
- [Single Sign-On](/security/for-admins/single-sign-on/)
- [System for Cross-domain Identity Management](/security/for-admins/provisioning/scim/) and more.

For a list of features available in each tier, see [Legacy Docker Pricing](https://www.docker.com/legacy-pricing/).

### Upgrade your Legacy Docker Business subscription

When you upgrade your Legacy Docker Business subscription to a Docker Business subscription, your subscription includes the following changes:

- Instead of paying an additional per-seat fee, Docker Build Cloud is now available to all users in your Docker subscription.
- Docker Build Cloud included minutes increase from 800/mo to 1500/mo. Docker Build Cloud minutes do not rollover month to month.
- Docker Scout now includes unlimited repositories with continuous vulnerability analysis, an increase from 3.
- 1500 Testcontainers Cloud runtime minutes are now included for use either in Docker Desktop or for CI. Testcontainers Cloud runtime minutes do not rollover month to month.
- Docker Hub image pull rate limits are removed.

For a list of features available in each tier, see [Docker Pricing](https://www.docker.com/pricing/).

## Legacy Docker Scout subscriptions

This section provides an overview of the legacy subscriptions for Docker
Scout.

> [!IMPORTANT]
>
> As of December 10, 2024, Docker Scout subscriptions are no longer available
> and have been replaced by Docker subscriptions that provide access to
> all tools. If you subscribed or renewed your subscriptions before December 10, 2024, your legacy Docker subscriptions still apply to your account until
> you renew. For more details, see [Announcing Upgraded Docker
> Plans](https://www.docker.com/blog/november-2024-updated-plans-announcement/).

### Legacy Docker Scout Free

Legacy Docker Scout Free is available for organizations. If you have a Legacy Docker subscription, you automatically have access to legacy Docker Scout Free.

Legacy Docker Scout Free includes:

- Unlimited local image analysis
- Up to 3 Docker Scout-enabled repositories
- SDLC integration, including policy evaluation and workload integration
- On-prem and cloud container registry integrations
- Security posture reporting

### Legacy Docker Scout Team

Legacy Docker Scout Team includes:

- All the features available in legacy Docker Scout Free
- In addition to 3 Docker Scout-enabled repositories, add up to 100 repositories when you buy your subscription

### Legacy Docker Scout Business

Legacy Docker Scout Business includes:

- All the features available in legacy Docker Scout Team
- Unlimited Docker Scout-enabled repositories

### Upgrade your Legacy Docker Scout subscription

When you upgrade your Legacy Docker Scout subscription to a Docker subscription, your
subscription includes the following changes:

- Docker Business: Unlimited repositories with continuous vulnerability analysis, an increase from 3.
- Docker Team: Unlimited repositories with continuous vulnerability analysis, an increase from 3
- Docker Pro: 2 included repositories with continuous vulnerability analysis.
- Docker Personal: 1 included repository with continuous vulnerability analysis.

For a list of features available in each tier, see [Docker Pricing](https://www.docker.com/pricing/).

## Legacy Docker Build Cloud subscriptions

 This section describes the features available for the different legacy Docker
 Build Cloud subscription tiers.

> [!IMPORTANT]
>
> As of December 10, 2024, Docker Build Cloud is only available with the
> new Docker Pro, Team, and Business plans. When your subscription renews on or after
> December 10, 2024, you will see an increase in your included Build Cloud
> minutes each month. For more details, see [Announcing Upgraded Docker
> Plans](https://www.docker.com/blog/november-2024-updated-plans-announcement/).

### Legacy Docker Build Cloud Starter

If you have a Legacy Docker subscription, a base level of Build Cloud
minutes and cache are included. The features available vary depending on your
Legacy Docker subscription tier.

#### Legacy Docker Pro

- 100 build minutes every month
- Available for one user
- 4 parallel builds

#### Legacy Docker Team

- 400 build minutes every month shared across your organization
- Option to onboard up to 100 members
- Can buy additional seats to add more minutes

#### Legacy Docker Business

- All the features listed for Docker Team
- 800 build minutes every month shared across your organization

### Legacy Docker Build Cloud Team

Legacy Docker Build Cloud Team offers the following features:

- 200 additional build minutes per seat
- Option to buy reserve minutes
- Increased shared cache

The legacy Docker Build Cloud Team subscription is tied to a Docker
[organization](/admin/organization/). To use the build minutes or
shared cache of a legacy Docker Build Cloud Team subscription, users must be a part of
the organization associated with the subscription. See Manage seats and invites.

### Legacy Docker Build Cloud Enterprise

For more details about your enterprise subscription, [contact sales](https://www.docker.com/products/build-cloud/#contact_sales).

### Upgrade your Legacy Docker Build Cloud subscription

You no longer need to subscribe to a separate Docker Build Cloud subscription to access
Docker Build Cloud or to scale your minutes. When you upgrade your Legacy Docker subscription to
a Docker subscription, your subscription includes the following changes:

- Docker Business: Included minutes are increased from 800/mo to 1500/mo with the option to scale more minutes.
- Docker Team: Included minutes are increased from 400/mo to 500/mo with the option to scale more minutes.
- Docker Pro: Included minutes are increased from 100/mo to 200/mo with the option to scale more minutes.
- Docker Personal: You receive a 7-day trial.

{{< /tab >}}
{{< /tabs >}}

## Subscription management options

### Self-serve

You manage everything directly including invoices, seats, billing information, and subscription changes.

### Sales-assisted

A dedicated Docker account manager handles setup and management for Docker Business and Team subscriptions.

## Support

All Docker Pro, Team, and Business subscribers receive email support for their subscriptions.