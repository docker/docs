---
title: Explore Docker Hardened Images
linktitle: Explore images
description: Learn how to find and evaluate image repositories, variants, metadata, and attestations in the DHI catalog on Docker Hub.
keywords: explore docker images, image variants, docker hub catalog, container image metadata, signed attestations
weight: 10
---

{{< summary-bar feature_name="Docker Hardened Images" >}}

Docker Hardened Images (DHI) are a curated set of secure, production-ready
container images designed for enterprise use. This page explains how to explore
available DHI repositories, review image metadata, examine variant details, and
understand the security attestations provided. Use this information to evaluate
and select the right image variants for your applications before mirroring them
to your organization.

## Explore Docker Hardened Images

To explore Docker Hardened Images (DHI):

1. Go to [Docker Hub](https://hub.docker.com) and sign in.
2. Select **My Hub**.
3. In the namespace drop-down, select your organization that has access to DHI.
4. Select **Hardened Images** > **Catalog**.

On the DHI page, you can browse images, search images, or filter images by
category.

## View repository details

To view repository details:

1. Go to [Docker Hub](https://hub.docker.com) and sign in.
2. Select **My Hub**.
3. In the namespace drop-down, select your organization that has access to DHI.
4. Select **Hardened Images** > **Catalog**.
5. Select a repository in the DHI catalog list.

The repository details page provides the following:

 - Overview: A brief explanation of the image.
 - Guides: Several guides on how to use the image and migrate your existing application.
 - Tags: Select this option to [view image variants](#view-image-variants).
 - Security summary: Select a tag name to view a quick security summary,
   including package count, total known vulnerabilities, and Scout health score.
 - Recently pushed tags: A list of recently updated image variants and when they
   were last updated.
 - Mirror to repository: Select this option to mirror the image to your
   organization's repository in order to use it. Only organization owners can mirror a repository.
 - View in repository: After a repository has been mirrored, you can select this
   option to view where the repository has been mirrored, or mirror it to another repository.

## View image variants

Tags are used to identify image variants. Image variants are different builds of
the same application or framework tailored for different use-cases.

To explore image variants:

1. Go to [Docker Hub](https://hub.docker.com) and sign in.
2. Select **My Hub**.
3. In the namespace drop-down, select your organization that has access to DHI.
4. Select **Hardened Images** > **Catalog**.
5. Select a repository in the DHI catalog list.
6. Select **Tags**.

The **Tags** page provides the following information:

- Tags: A list of all available tags, also known as image variants.
- Compliance: Lists relevant compliance designations. For example, `FIPS` or `STIG`.
- Distribution: The distribution that the variant is based on. For example, `debian 12` or `alpine 3.21`.
- Package manager: The package manager that is available in the variant. For example, `apt`, `apk`, or `-` (no package manager).
- Shell: The shell that is available in the variant. For example, `bash`, `busybox`, or `-` (no shell).
- User: The user that the container runs as. For example, `root`, `nonroot (65532)`, or `node (1000)`.
- Last pushed: The amount of days ago that the image variant was last pushed.
- Vulnerabilities: The amount of vulnerabilities in the variant based on the severity.
- Health: The Scout health score for the variant. Select the score icon to get more details.

> [!NOTE]
>
> Unlike most images on Docker Hub, Docker Hardened Images do not use the
> `latest` tag. Each image variant is published with a full semantic version tag
> (for example, `3.13`, `3.13-dev`) and is kept up to date. If you need to pin
> to a specific image release for reproducibility, you can reference the image
> by its [digest](../core-concepts/digests.md).

## View image variant details

To explore the details of an image variant:

1. Go to [Docker Hub](https://hub.docker.com) and sign in.
2. Select **My Hub**.
3. In the namespace drop-down, select your organization that has access to DHI.
4. Select **Hardened Images** > **Catalog**.
5. Select a repository in the DHI catalog list.
6. Select **Tags**.
7. Select the image variant's tag in the table.

The image variant details page provides the following information:

- Packages: A list of all packages included in the image variant. This section
  includes details about each package, including its name, version,
  distribution, and licensing information.
- Specifications: The specifications for the image variant include the following
  key details:
   - Source & Build Information: The image is built from the Dockerfile found
     here and the Git commit.
   - Build parameters
   - Entrypoint
   - CMD
   - User
   - Working directory
   - Environment Variables
   - Labels
   - Platform
- Vulnerabilities: The vulnerabilities section provides a list of known CVEs for
  the image variant, including:
   - CVE
   - Severity
   - Package
   - Fix version
   - Last detected
   - Status
   - Suppressed CVEs
- Attestations: Variants include comprehensive security attestations to verify
  the image's build process, contents, and security posture. These attestations
  are signed and can be verified using cosign. For a list of available
  attestations, see [Attestations](../core-concepts/attestations.md).

## What's next

After finding an image you need, you can [mirror the image to your
organization](./mirror.md). If the image is already mirrored, then you can start
[using the image](./use.md).