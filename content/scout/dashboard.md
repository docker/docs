---
description: The Docker Scout Dashboard helps review and share the analysis of images.
keywords: scanning, analysis, vulnerabilities, Hub, supply chain, security, report,
  reports, dashboard
title: Dashboard
aliases:
- /scout/reports/
- /scout/web-app/
---

The [Docker Scout Dashboard](https://scout.docker.com/) helps you share the
analysis of images in an organization with your team. Developers can now see an
overview of their security status across all their images from both Docker Hub
and Artifactory, and get remediation advice at their fingertips. It helps team
members in roles such as security, compliance, and operations to know what
vulnerabilities and issues they need to focus on.

## Overview

![A screenshot of the Docker Scout Dashboard overview](./images/dashboard-overview.webp?border=true)

The **Overview** tab provides a summary for the repositories in the selected
organization.

At the top of this page, you can select which **Environment** to view.
By default, the most recently pushed images are shown. To learn more about
environments, see [Environment monitoring](./integrations/environment/_index.md).

The **Policy** boxes show your current compliance rating for each policy, and a
trend indication for the selected environment. The trend describes the policy
delta for the most recent images compared to the previous version.
For more information about policies, see [Policy Evaluation](./policy/_index.md).

The vulnerability chart shows the total number of vulnerabilities for images in
the selected environment over time. You can configure the timescale for the
chart using the drop-down menu.

Use the header menu at the top of the website to access the different main
sections of the Docker Scout Dashboard:

- [Images](#images)
- [Policies](#policies)
- [Vulnerabilities](#vulnerabilities)
- [Base images](#base-images)
- [Packages](#packages)

## Images

The **Images** view shows a list of images in an organization. You can search
for specific repositories using the search box.

Each entry in the list shows the following details:

- The repository name for the image. Selecting the link for the repository opens [the list of tags for the repository](#repository-tag-list).
- The most recent tag of the image in the selected environment. Selecting the link for the base image opens [the image details view](#image-details-view).
- The operating system and architecture of the image.
- The date of the last push for the image.
- The vulnerabilities for the most recent image version.
- Policy status, including the change for the most recent version, and a link to more details for non-compliant images.

### Repository tag list

![Screenshot of tags for a repository](./images/dashboard-repo-tags.webp?border=true)

There are two tabs on this page:

- The **Policy** tab displays the policy delta for the latest version of the
  image.
- The **Tags** tab contains the repository tag list, and shows all tags for the
  repository. 

In the **Tags** tab, you can filter the list by environment, or by tag or
digest using the search box.

Each entry in the list shows the following details:

- A checkbox to mark the tag for comparison to one other.

  > **Tip**
  >
  > Compare two image tags by selecting the checkboxes next to them and selecting the **Compare images** button at the top of the list.
  { .tip }

- The tag version or image digest. Clicking the link for version opens [the image layer view](#image-details-view).
- The [environments](./integrations/environment/_index.md) that the image is assigned to.
- The operating system and architecture of the image.
- The vulnerabilities for the tag version.
- The last push for the tag version.
- The base image and version used by the repository and the vulnerabilities for that version.

#### Compare images

You can compare two or more images in the list. Mark the image versions that
you want to compare, and select **Compare images**.

The top section of the comparison view shows an overview of the two selected
image tags. The tabs section of the view shows the following:

- Select the **Packages** tab to see packages added, removed, or changed in each image. Each entry in the table shows the differences between the versions and vulnerabilities in each image. Select the disclosure triangle next to a package to see more detail on the vulnerabilities changed.
- Select the **Vulnerabilities** tab to see changes to the vulnerabilities present in each image.

### Image details view

Selecting an image tag takes you to the image details view. This view contains
two tabs that let you drill down into the details of the composition and
policy compliance for the image: **Policy status** and **Image layers**.

{{< tabs >}}
{{< tab name="Policy status" >}}

![Screenshot of the policy tab in the image details view](./images/dashboard-image-policies.webp?border=true)

The policy tab shows you the policy evaluation results for the image. Use the
**View details** and **View fixes** links to the right to view the full
evaluation results, and learn how to improve compliance score for non-compliant
images.

For more information about policy, see [Policy Evaluation](./policy/_index.md).

{{< /tab >}}
{{< tab name="Image layers" >}}

![Screenshot showing Docker Scout image layers](./images/dashboard-image-layers.webp?border=true)

The layer view shows a breakdown of the Docker Scout analysis, including
an overview of the digest Secure Hash Algorithms (SHA), version, the image hierarchy (base images), image
layers, packages, and vulnerabilities.

> **Note**
>
> You can find more details on the elements in the image layer view in [the image details view docs](./image-details-view.md).

{{< /tab >}}
{{< /tabs >}}

## Policies

![A screenshot of the Docker Scout policies view](./images/dashboard-policies-view.webp?border=true)

The **Policies** view shows a breakdown of policy compliance for all of the
images in the selected organization and environment. You can use the **Image**
drop-down menu to view a policy breakdown for a specific environment.

For more information about policies, see [Policy Evaluation](./policy/_index.md).

## Base images

![A screenshot of the Docker Scout view showing base images used](./images/dashboard-base-images.webp?border=true)

The **Base images** view shows all base images used by repositories in an organization.

Each entry in the list shows the following details:

- The base image name.
- The versions of the base image used by images in the organization.
- The number of images that use the base image. Selecting the link opens [the list of images that use the base image view](#images-using-base-image).
- The number of packages in the base image.

### Images using base image

The **Images** tab shows all images in an organization that use a specific base image.

Each entry in the list shows the following details:

- The repository name. Selecting the link opens [the list of tags for the repository](#repository-tag-list).
- The most recent tag of the image and its vulnerabilities. Selecting the link for the tag opens [the Image layer detail view](#image-details-view) for the repository.
- The operating system and architecture of the image.
- The base image tag used by the repository. Selecting the link opens [the image layer detail view](#image-details-view) for that version.
- The current base image digest for the repository.
- The date of the last push for the repository.

## Packages

The **Packages** view shows all packages across repositories in an organization.

Each entry in the list shows the following details:

- The package name.
- The package type.
- The versions of the package used by images in the organization.
- The number of images that use the package.

## Vulnerabilities

The **Vulnerabilities** view shows a list of all vulnerabilities from images in
the organization. You can sort and filter the list by severity and search for
Common Vulnerabilities and Exposures (CVE) ID using the search box.

Each entry in the list shows the following details:

- Severity of the vulnerability.

  > **Note**
  >
  > Docker Scout bases the calculation behind this severity level on a variety
  > of sources.

- The severity of the vulnerability.
- The vulnerability CVE ID. Selecting the link for the CVE ID opens [the vulnerability details page](#vulnerability-details-page).
- The package name and version affected by this CVE.
- The Common Vulnerability Scoring System (CVSS) score for the vulnerability. Docker Scout shows the highest CVSS score from multiple sources.
- The number of images in the organization that use the package affected by this CVE. Selecting this link opens the [vulnerability details page](#vulnerability-details-page).
- If Docker Scout knows of a fix for the vulnerability, and if so, the package version of the fix.

### Vulnerability details page

The vulnerability details page shows detailed information about a particular
CVE. This page is a publicly open page. You can share the link to a particular
CVE description with other people even if they're not a member of your Docker
organization.

The page shows the following information:

- The CVE ID and severity.
- A description of the vulnerability.
- The number of packages affected by the vulnerability.
- The vulnerability publish date.

Following this information is a list of all repositories affected by the
vulnerability, searchable by image name. Each entry in the list shows the
following details:

- The repository name. Selecting the link for the repository name opens [the repository tag list view](#repository-tag-list).
- The current tag version of the image. Selecting the link for the tag name opens [the repository tag list layer view](#image-details-view).
- The date the image was last pushed.
- The registry where the image is stored.
- The affected package name and version in the image.

## Settings

The settings menu under the drop-down in the website header contains link to go
to the [Integrations](#integrations) page and [Repository
settings](#repository-settings).

### Integrations

The **Integrations** page lets you create and manage your Docker Scout
integrations, such as environment integrations and registry integrations. For
more information on how to get started with integrations, see [Integrating
Docker Scout with other systems](./integrations/_index.md).

### Repository settings

The **Repository settings** is where you enable and disable Docker Scout for
repositories in your organization.

To enable repositories, select the checkboxes for the repositories on which you
want to enable Docker Scout analysis and select **Enable image analysis**.

When you enable image analysis for a repository, Docker Scout analyzes new tags
automatically when you push to that repository.

Disable Docker Scout analysis on selected repositories by selecting **Disable
image analysis**.
