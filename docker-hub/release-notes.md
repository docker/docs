---
title: Docker Hub release notes
description: Learn about the new features, bug fixes, and breaking changes for Docker Hub
keywords: docker hub, whats new, release notes
toc_min: 1
toc_max: 2
---

Here you can learn about the latest changes, new features, bug fixes, and
known issues for each Docker Hub release.

Take a look at the [Docker Public Roadmap](https://github.com/docker/roadmap/projects/1){: target="_blank" rel="noopener" class="_"} to see what's coming next.

## 2022-09-26
### New feature

The new [autobuild feature](../docker-hub/builds/index.md#check-your-active-builds) lets you view your in-progress logs every 30 seconds instead of when the build is complete.

## 2022-09-21
### Enhancement

In Docker Hub, you can now download a [registry.json](../docker-hub/configure-sign-in.md) file or copy the commands to create a registry.json file to enforce sign-in for your organization.

## 2022-09-19
### Enhancement

You can now [export a CSV file of members](../docker-hub/members.md/#export-members) from organizations that you own.


## 2022-07-22
### Enhancement

You can now invite members to your organization with a CSV file containing their email addresses. The CSV file can either be a file you create for this specific purpose or one that’s extracted from another in-house system.

## 2022-07-19

### Enhancement

SCIM is now available for organizations with a Docker Business subscription using an Azure AD identity provider.

## 2022-06-23

### New feature

With [SCIM](scim.md), you can manage users within your Okta identity provider (IdP). In addition, you can enable SCIM on organizations that are part of the Docker Business subscription.

## 2022-05-24

### New feature

[Registry Access Management](registry-access-management.md) is now available for all Docker Business subscriptions. When enabled, your users can access specific registries in Docker Hub.

## 2022-05-03

### New feature

Organization owners can [invite new members](members.md#invite-members) to an organization via Docker ID or email address.

## 2021-11-15

### New

You can now purchase or upgrade to a Docker Business subscription using a credit card. To learn more, see [Upgrade your subscription](../subscription/upgrade.md).

## 2021-08-31

### New

Docker has [announced](https://www.docker.com/blog/updating-product-subscriptions/){: target="*blank" rel="noopener" class="*" id="dkr_docs_relnotes_btl"} updates and extensions to the product subscriptions to increase productivity, collaboration, and added security for our developers and businesses. Docker subscription tiers now include Personal, Pro, Team, and Business.

The updated [Docker Subscription Service Agreement](https://www.docker.com/legal/docker-subscription-service-agreement) includes a change to the terms for **Docker Desktop**.

- Docker Desktop **remains free** for small businesses (fewer than 250 employees AND less than $10 million in annual revenue), personal use, education, and non-commercial open source projects.
- It requires a paid subscription (**Pro, Team, or Business**), for as little as $5 a month, for professional use in larger enterprises.
- The effective date of these terms is August 31, 2021. There is a grace period until January 31, 2022 for those that will require a paid subscription to use Docker Desktop.
- The Docker Pro and Docker Team subscriptions now **include commercial use** of Docker Desktop.
- The existing Docker Free subscription has been renamed **Docker Personal**.
- **No changes** to Docker Engine or any other upstream **open source** Docker or Moby project.

    To understand how these changes affect you, read the [FAQs](https://www.docker.com/pricing/faq){: target="*blank" rel="noopener" class="*" id="dkr_docs_relnotes_btl"}. For more information, see [Docker subscription overview](../subscription/index.md).

## 2021-05-05

### Enhancement

When managing the content of your repositories, you can now filter the results based on the currentness of the tags and more easily identify your untagged images.
For more information, see  [Advanced Management Dashboard](image-management.md).

For Docker Hub API documentation, see [Docker Hub API Reference](api/latest.md#operation/GetNamespacesRepositoriesImages).

## 2021-04-13

### Enhancement

The **Billing Details** page now shows any organizations you own, in addition to your personal account. This allows you to clearly identify the billing details of your chosen namespace, and enables you to switch between your personal and your organization accounts to view or update the details.

## 2021-04-09

### Enhancement

You can now specify any email address to receive billing-related emails for your organization. The email address doesn't have to be associated with an organization owner account. You must be an owner of the organization to update any billing details.

To change the email address receiving billing-related emails, log into Docker Hub and navigate to the **Billing** tab of your organization. Select **Payment Methods** > **Contact Information**. Enter the new email address that you'd like to use in the **Email** field. Click **Update** for the changes to take effect.

For details on how to update your billing information, see [Update billing information](../subscription/upgrade.md#update-billing-information).

## 2021-03-22

### New feature

**Advanced Image Management dashboard**

Docker introduces the Advanced Image Management dashboard that enables you to view and manage Docker images in your repositories. For more information, see [Advanced Image Management dashboard](image-management.md).

## 2021-01-25

### New feature

Docker introduces Audit logs, a new feature that allows team owners to view a list of activities that occur at organization and repository levels. This feature begins tracking the activities from the release date, that is, **from 25 January 2021**.

For more information about this feature and for instructions on how to use it, see [Audit logs](audit-log.md).

## 2020-11-10

### New feature

The **Repositories** view now shows which images have gone stale because they haven't been pulled or pushed recently. For more information, see [repository tags](repos/access/index.md#viewing-repository-tags).

# 2020-10-07

### New feature

Docker introduces Hub Vulnerability Scanning which enables you to automatically scan Docker images for vulnerabilities using Snyk. For more information, see [Hub Vulnerability Scanning](vulnerability-scanning.md).

## 2020-05-14

### New features

* Docker has announced a new, per-seat pricing model to accelerate developer workflows for cloud-native development. The previous private repository/concurrent autobuild-based plans have been replaced with new **Pro** and **Team** plans that include unlimited private repositories. For more information, see [Docker subscription](../subscription/index.md).

* Docker has enabled download rate limits for downloads and pull requests on Docker Hub. This caps the number of objects that users can download within a specified timeframe. For more information, see [Download rate limit](download-rate-limit.md).

## 2019-11-04

### Enhancements

* The [repositories page](repos/index.md) and all
related settings and tabs have been updated and moved from `cloud.docker.com`
to `hub.docker.com`. You can access the page at its new URL: [https://hub.docker.com/repositories](https://hub.docker.com/repositories).

### Known Issues

* Scan results don't appear for some official images.

## 2019-10-21

### New features
* **Beta:** Docker Hub now supports [two-factor authentication (2FA)](2fa/index.md). Enable it in your account settings, under the **[Security](https://hub.docker.com/settings/security)** section.

    > If you lose both your 2FA authentication device and recovery code, you may
    > not be able to recover your account.
    {: .important }

### Enhancements
* As a security measure, when two-factor authentication is enabled, the Docker CLI requires a personal access token instead of a password to log in.

### Known Issues

* Scan results don't appear for some official images.


## 2019-10-02

### Enhancements
* You can now manage teams and members straight from your [organization page](https://hub.docker.com/orgs).
Each organization page now breaks down into these tabs:
  * **New:** Members - manage your members directly from this page (delete,
  add, or open their teams)
  * **New:** Teams - search by team or username, and open up any team page to
  manage the team
  * **New:** Invitees (conditional tab, only if an invite exists) - resend or
  remove invitiations from this tab
  * Repositories
  * Settings
  * Billing

### Bug fixes

* Fixed an issue where Kinematic could not connect and log in to Docker Hub.

### Known Issues

* Scan results don't appear for some official images.


## 2019-09-19

### New features

* You can now [create personal access tokens](access-tokens.md) in Docker Hub and use them to authenticate from the Docker CLI. Find them in your account settings, under the new **[Security](https://hub.docker.com/settings/security)** section.

### Known Issues

* Scan results don't appear for some official images.


## 2019-09-16

### Enhancements

* The [billing page](../subscription/upgrade.md) for personal accounts has been updated. You can access the page at its new URL: [https://hub.docker.com/billing/plan](https://hub.docker.com/billing/plan).

### Known Issues

* Scan results don't appear for some official images.


## 2019-09-05

### Enhancements

* The `Tags` tab on an image page now provides additional information for each tag:
  * A list of digests associated with the tag
  * The architecture it was built on
  * The OS
  * The user who most recently updated an image for a specific tag
* The security scan summary for [official images](official_images.md)
has been updated.

### Known Issues

* Scan results don't appear for some official images.
