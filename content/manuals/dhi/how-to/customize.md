---
title: Customize a Docker Hardened Image
linkTitle: Customize an image
weight: 25
keywords: debug, hardened images, DHI, customize, certificate, artifact
description: Learn how to customize a Docker Hardened Images (DHI).
---

You can customize a Docker Hardened Image (DHI) to suit your specific needs
using the Docker Hub UI. This allows you to select a base image, add packages,
add OCI artifacts (such as custom certificates or additional tools), and
configure settings. In addition, the build pipeline ensures that your customized
image is built securely and includes attestations.

Your customized images stay secure automatically. When the base Docker Hardened
Image receives a security patch or your OCI artifacts are updated, Docker
automatically rebuilds your customized images in the background. This ensures
continuous compliance and protection by default, with no manual work required.
The rebuilt images are signed and attested to the same SLSA Build Level 3
standard as the base images, ensuring a secure and verifiable supply chain.

## Customize a Docker Hardened Image

To add a customized Docker Hardened Image to your organization, an organization
owner must first [mirror](./mirror.md) the DHI repository to your organization.
Once the repository is mirrored, any user with access to the mirrored DHI
repository can create a customized image.

To customize a Docker Hardened Image, follow these steps:

1. Sign in to [Docker Hub](https://hub.docker.com).
1. Select **My Hub**.
1. In the namespace drop-down, select your organization that has a mirrored DHI
   repository.
1. Select **Hardened Images** > **Management**.
1. For the mirrored DHI repository you want to customize, select the menu icon in the far right column.
1. Select **Customize**.

   At this point, the on-screen instructions will guide you through the
   customization process. You can continue with the following steps for more
   details.

1. Select the image version you want to customize.
1. Optional. Add packages.

   1. In the **Packages** drop-down, select the packages you want to add to the
      image.

      The packages available in the drop-down are OS system packages for the
      selected image variant. For example, if you are customizing the Alpine
      variant of the Python DHI, the list will include all Alpine system
      packages.

   1. In the **OCI artifacts** drop-down, first, select the repository that
      contains the OCI artifact image. Then, select the tag you want to use from
      that repository. Finally, specify the specific paths you want to include
      from the OCI artifact image.

      The OCI artifacts are images that you have previously
      built and pushed to a repository in the same namespace as the mirrored
      DHI. For example, you can add a custom root CA certificate or a another
      image that contains a tool you need, like adding Python to a Node.js
      image. For more details on how to create an OCI artifact image, see
      [Create an OCI artifact image](#create-an-oci-artifact-image).

      When combining images that contain directories and files with the same
      path, images later in the list will overwrite files from earlier images.
      To manage this, you must select paths to include and optionally exclude
      from each OCI artifact image. This allows you to control which files are
      included in the final customized image.

      By default, no files are included from the OCI artifact image. You must
      explicitly include the paths you want. After including a path, you can
      then explicitly exclude files or directories underneath it.

      > [!NOTE]
      >
      > When files necessary for runtime are overwritten by OCI artifacts, the
      > image build still succeeds, but you may have issues when running the
      > image.

   1. In the **Scripts** section, you can add, edit, or remove scripts.

      Scripts let you add files to the container image that you can access at runtime. They are not executed during
      the build process. This is useful for services that require pre-start initialization, such as setup scripts or
      file writes to directories like `/var/lock` or `/out`.

      You must specify the following:

      - The path where the script will be placed
      - The script content
      - The UID and GID ownership of the script
      - The octal file permissions of the script

1. Select **Next: Configure** and then configure the following options.
1. Specify a suffix that is appended to the customized image's tag. For
   example, if you specify `custom` when customizing the `dhi-python:3.13`
   image, the customized image will be tagged as `dhi-python:3.13_custom`.
1. Select the platforms you want to build the image for.
1. Add [`ENTRYPOINT`](/reference/dockerfile/#entrypoint) and
   [`CMD`](/reference/dockerfile/#cmd) arguments to the image. These
   arguments are appended to the base image's entrypoint and command.
1. Specify the users to add to the image.
1. Specify the user groups to add to the image.
1. Select which [user](/reference/dockerfile/#user) to run the images as.
1. Specify the [environment variables](/reference/dockerfile/#env) and their
   values that the image will contain.
1. Add [annotations](/build/metadata/annotations/) to the image.
1. Add [labels](/reference/dockerfile/#label) to the image.
1. Select **Create Customization**.

   A summary of the customization appears. It may take some time for the image
   to build. Once built, it will appear in the **Tags** tab of the repository,
   and your team members can pull it like any other image.

## Edit or delete a Docker Hardened Image customization

To edit or delete a Docker Hardened Image customization, follow these steps:

1. Sign in to [Docker Hub](https://hub.docker.com).
2. Select **My Hub**.
3. In the namespace drop-down, select your organization that has a mirrored DHI.
4. Select **Hardened Images** > **Management**.
5. Select **Customizations**.

6. For the customized DHI repository you want to manage, select the menu icon in the far right column.
   From here, you can:

   - **Edit**: Edit the customized image.
   - **Create new**: Create a new customized image based on the source repository.
   - **Delete**: Delete the customized image.

7. Follow the on-screen instructions to complete the edit or deletion.

## Create an OCI artifact image

An OCI artifact image is a Docker image that contains files or directories that
you want to include in your customized Docker Hardened Image (DHI). This can
include additional tools, libraries, or configuration files.

When creating an image to use as an OCI artifact, it should ideally be as
minimal as possible and contain only the necessary files.

For example, to distribute a custom root CA certificate as part of a trusted CA
bundle, you can use a multi-stage build. This approach registers your
certificate with the system and outputs an updated CA bundle, which can be
extracted into a minimal final image:

```dockerfile
# syntax=docker/dockerfile:1

FROM <your-namespace>/dhi-bash:5-dev AS certs

ENV DEBIAN_FRONTEND=noninteractive

RUN mkdir -p /usr/local/share/ca-certificates/my-rootca
COPY certs/rootCA.crt /usr/local/share/ca-certificates/my-rootca

RUN update-ca-certificates

FROM scratch
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
```

You can follow this pattern to create other OCI artifacts, such as images
containing tools or libraries that you want to include in your customized DHI.
Install the necessary tools or libraries in the first stage, and then copy the
relevant files to the final stage that uses `FROM scratch`. This ensures that
your OCI artifact is minimal and contains only the necessary files.

Build and push the OCI artifact image to a repository in your organization's
namespace and it automatically appears in the customization workflow when you
select the OCI artifacts to add to your customized Docker Hardened Image.
