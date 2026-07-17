---
title: Use Docker Hub
linktitle: Docker Hub
description: Browse the DHI catalog on Docker Hub to search repositories, inspect image metadata, and view SBOMs, CVEs, and attestations.
weight: 10
keywords: docker hub dhi catalog, hardened images hub, dhi repository details, image variants hub
---

The [Docker Hardened Images catalog](https://hub.docker.com/hardened-images/catalog)
on Docker Hub is the primary web interface for browsing, searching, and inspecting
DHI repositories and their metadata.

## Catalog page

The catalog lists all available DHI repositories. You can filter by name,
image type, or compliance requirements (FIPS, STIG) to find the image you need.

## Repository details page

When you select a repository from the catalog, the repository details page
provides the following:

- Overview: A brief explanation of the image.
- Guides: Several guides on how to use the image and migrate your existing application.
- Images: Select this option to [view image variants](#images-page).
- Security summary: Select a tag name to view a quick security summary,
  including package count and total known vulnerabilities.
- Recently pushed tags: A list of recently updated image variants and when they
  were last updated.
- Use this image: After selecting an image variant, you can select this option to
  view instructions on how to pull and use the image variant, or select **Mirror
  repository** to mirror it to your organization.

## Images page

From the repository details page, select **Images** to see all available image
variants for that repository. The table includes:

- Image version: The image name with its base distribution (for example, `debian
  13`) and associated tags.
- Type: The support lifecycle status of the variant.
- Compliance: Relevant compliance designations, for example `CIS`, `FIPS`, or
  `STIG (100%)`.
- Package manager: Whether a package manager is available. A checkmark indicates
  a package manager is present (for example, `apt` or `apk`), a dash indicates
  none.
- Shell: Whether a shell is available. A checkmark indicates a shell is present
  (for example, `bash` or `busybox`), a dash indicates none.
- User: The user that the container runs as, for example `root` or `nonroot
  (65532)`.
- Last pushed: When the image variant was last updated.
- Vulnerabilities: Vulnerability counts by severity level.

## Image variant details page

Select an image version from the Images table to view detailed information about
that specific variant:

- Packages: A list of all packages included in the image variant, with each
  package's name, version, distribution, and licensing information.
- Specifications:
  - Source and build information: The Dockerfile and Git commit used to build the image.
  - Build parameters, entrypoint, CMD, user, working directory, environment
    variables, labels, and platform.
- Vulnerabilities: A list of known CVEs for the image variant, including CVE ID,
  severity, affected package, fix version, last detected date, status, and
  suppressed CVEs.
- Attestations: Signed security attestations covering the image's build process,
  contents, and security posture. For the full list, see
  [Attestations](/dhi/explore/security-concepts/attestations/).

## Manage page

The Manage page (**My Hub** > **Hardened Images** > **Manage**) is the central
place for administering your organization's mirrored DHI repositories. It has
two tabs:

- Mirrored Images: Lists all image repositories currently mirrored to your
  organization, with their source DHI repository, destination repository name,
  and mirroring status. From here you can stop mirroring or open a repository's
  settings.
- Mirrored Helm charts: The same view for Helm chart repositories.

Selecting a mirrored repository opens its settings, where you can enable or
disable Extended Lifecycle Support (ELS) and access customizations.

For step-by-step instructions, see [Mirror a Docker Hardened Image
repository](/dhi/how-to/mirror/).

## Customizations

Customizations are accessible from **My Hub** > **Hardened Images** > **Manage** > **Mirrored Images**.
Select the menu icon next to a mirrored repository and
then **Customize**. Each customization defines
additional packages, OCI artifacts, environment variables, or labels to layer
onto the base DHI during a rebuild.

The customizations view shows each customization's name, status, and last build
time. Selecting a customization opens its configuration, where you can edit the
definition, trigger a rebuild, or delete it.

For step-by-step instructions, see [Customize a Docker Hardened
Image](/dhi/how-to/customize/).
