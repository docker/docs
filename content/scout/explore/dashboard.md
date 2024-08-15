---
description: The Docker Scout Dashboard helps review and share the analysis of images.
keywords: scout, scanning, analysis, vulnerabilities, Hub, supply chain, security, report,
  reports, dashboard
title: Dashboard
aliases:
- /scout/reports/
- /scout/web-app/
- /scout/dashboard/
---

The [Docker Scout Dashboard](https://scout.docker.com/) helps you share the
analysis of images in an organization with your team. Developers can now see an
overview of their security status across all their images from both Docker Hub
and Artifactory, and get remediation advice at their fingertips. It helps team
members in roles such as security, compliance, and operations to know what
vulnerabilities and issues they need to focus on.

## Overview

![A screenshot of the Docker Scout Dashboard overview](../images/dashboard-overview.webp?border=true)

The **Overview** tab provides a summary for the repositories in the selected
organization.

At the top of this page, you can select which **Environment** to view.
By default, the most recently pushed images are shown. To learn more about
environments, see [Environment monitoring](/scout/integrations/environment/_index.md).

The **Policy** boxes show your current compliance rating for each policy, and a
trend indication for the selected environment. The trend describes the policy
delta for the most recent images compared to the previous version.
For more information about policies, see [Policy Evaluation](/scout/policy/_index.md).

The vulnerability chart shows the total number of vulnerabilities for images in
the selected environment over time. You can configure the timescale for the
chart using the drop-down menu.

Use the header menu at the top of the website to access the different main
sections of the Docker Scout Dashboard:

- **Policies**: shows the policy compliance for the organization, see [Policies](#policies)
- **Images**: lists all Docker Scout-enabled repositories in the organization, see [Images](#images)
- **Base images**: lists all base images used by repositories in an organization
- **Packages**: lists all packages across repositories in the organization
- **Vulnerabilities**: lists all CVEs in the organization's images, see [Vulnerabilities](#vulnerabilities)
- **Integrations**: create and manage third-party integrations, see [Integrations](#integrations)
- **Settings**: manage repository and billing settings, see [Settings](#settings)

## Policies

The **Policies** view shows a breakdown of policy compliance for all of the
images in the selected organization and environment. You can use the **Image**
drop-down menu to view a policy breakdown for a specific environment.

For more information about policies, see [Policy Evaluation](/scout/policy/_index.md).

## Images

The **Images** view shows all images in Scout-enabled repositories for the selected environment.
You can filter the list by selecting a different environment, or by repository name using the text filter.

![Screenshot of the images view](../images/dashboard-images.webp)

For each repository, the list displays the following details:

- The repository name (image reference without the tag or digest)
- The most recent tag of the image in the selected environment
- Operating systems and architectures for the most recent tag
- Vulnerabilities status for the most recent tag
- Policy status for the most recent tag

Selecting a repository link takes you to a list of all images in that repository that have been analyzed.
From here you can view the full analysis results for a specific image,
and compare tags to view the differences in packages and vulnerabilities

Selecting an image link takes you to a details view for the selected tag or digest.
This view contains two tabs that detail the composition and policy compliance for the image:

- **Policy status** shows the policy evaluation results for the selected image.
  Here you also have links for details about the policy violations.

  For more information about policy, see [Policy Evaluation](/scout/policy/_index.md).

- **Image layers** shows a breakdown of the image analysis results.
  You can get a complete view of the vulnerabilities your image contains
  and understand how they got in.

## Vulnerabilities

The **Vulnerabilities** view shows a list of all vulnerabilities for images in the organization.
This list includes details about CVE such as the severity and Common Vulnerability Scoring System (CVSS) score,
as well as whether there's a fix version available.
The CVSS score displayed here is the highest score out of all available [sources](/scout/deep-dive/advisory-db-sources.md).

Selecting the links on this page opens the vulnerability details page,
This page is a publicly visible page, and shows detailed information about a CVE.
You can share the link to a particular CVE description with other people
even if they're not a member of your Docker organization or signed in to Docker Scout.

If you are signed in, the **My images** tab on this page lists all of your images
affected by the CVE.

## Integrations

The **Integrations** page lets you create and manage your Docker Scout
integrations, such as environment integrations and registry integrations. For
more information on how to get started with integrations, see
[Integrating Docker Scout with other systems](/scout/integrations/_index.md).

## Settings

The settings menu in the Docker Scout Dashboard contains:

- [**Billing**](#billing-settings) for managing your Docker Scout subscription and payments
- [**Repository settings**](#repository-settings) for enabling and disabling repositories
- [**Notifications**](#notification-settings) for managing your notification preferences for Docker Scout.

### Billing settings

The [Billing settings](https://scout.docker.com/settings/billing) page shows
you the Docker Scout plan for the current organization. Here you can see what's
included in your plan, compare it with other available plans, and change the
plan if you're an organization owner.

For more information about subscription plans, see
[Docker Scout subscriptions and features](/subscription/scout-details.md)

### Repository settings

When you enable Docker Scout for a repository,
Docker Scout analyzes new tags automatically when you push to that repository.
To enable repositories in Amazon ECR, Azure ACR, or other third-party registries,
you first need to integrate them.
See [Container registry integrations](/scout/integrations/_index.md#container-registries)

### Notification settings

The [Notification settings](https://scout.docker.com/settings/notifications)
page is where you can change the preferences for receiving notifications from
Docker Scout. Notification settings are personal, and changing notification
settings only affects your personal account, not the entire organization.

The purpose of notifications in Docker Scout is to raise awareness about
upstream changes that affect you. Docker Scout will notify you about when a new
vulnerability is disclosed in a security advisory, and it affects one or more
of your images. You will not receive notifications about changes to
vulnerability exposure or policy compliance as a result of pushing a new image.

> [!NOTE]
>
> Notifications are only triggered for the *last pushed* image tags for each
> repository. "Last pushed" refers to the image tag that was most recently
> pushed to the registry and analyzed by Docker Scout. If the last pushed image
> is not affected by a newly disclosed CVE, then no notification will be
> triggered.

The available notification settings are:

- **Repository scope**

  Here you can select whether you want to enable notifications for all
  repositories, or only for specific repositories. These settings apply to the
  currently selected organization, and can be changed for each organization you
  are a member of.

  - **All repositories**: select this option to receive notifications for all
    repositories that you have access to.
  - **Specific repositories**: select this option to receive notifications for
    specific repositories. You can then enter the names of repositories you
    want to receive notifications for.

- **Delivery preferences**

  These settings control how you receive notifications from Docker Scout. They
  apply to all organizations that you're a member of.

  - **Notification pop-ups**: select this check-box to receive notification
    pop-up messages in the Docker Scout Dashboard.
  - **OS notifications**: select this check-box to receive OS-level notifications
    from your browser if you have the Docker Scout Dashboard open in a browser
    tab.
  
  To enable OS notifications, Docker Scout needs permissions to send
  notifications using the browser API.

From this page, you can also go to the settings for Team collaboration
integrations, such as the [Slack](/scout/integrations/team-collaboration/slack.md)
integration.

You can also configure your notification settings in Docker Desktop by going
to **Settings** > **Notifications**.
