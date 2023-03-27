---
description: Image Access Management
keywords: image, access, management
title: Image Access Management
---

>Note
>
>Image Access Management is available to [Docker Business](../subscription/details.md) customers only. 

Image Access Management gives administrators control over which types of images, such as Docker Official Images, Docker Verified Publisher Images, or community images, their developers can pull from Docker Hub.

For example, a developer, who is part of an organization, building a new containerized application could accidentally use an untrusted, community image as a component of their application. This image could be malicious and pose a security risk to the company. Using Image Access Management, the organization owner can ensure that the developer can only access trusted content like Docker Official Images, Docker Verified Publisher Images, or the organization’s own images, preventing such a risk.

## Configure Image Access Management permissions

1. Sign into your [Docker Hub](https://hub.docker.com){: target="_blank" rel="noopener" class="_"} account as an organization administrator.
2. Select an organization, and navigate to the **Settings** tab
3. From the **Organizations** page select **Org Permissions**.
4. Enable Image Access Management to set the permissions for the following categories of images you can manage:
- **Organization Images**: When Image Access Management is enabled, images from your organization are always allowed. These images can be public or private created by members within your organization.
- **Docker Official Images**: A curated set of Docker repositories hosted on Hub. They provide OS repositories, best practices for Dockerfiles, drop-in solutions, and applies security updates on time.
- **Docker Verified Publisher Images**: published by Docker partners that are part of the Verified Publisher program and are qualified to be included in the developer secure supply chain. You can set permissions to **Allowed** or **Restricted**.
- **Community Images**: Images are always disabled when Image Access Management is enabled. These images are not trusted because various Docker Hub users contribute them and pose security risks.

    > **Note**
    >
    > Image Access Management is turned off by default. However, members of the `owners` team in your organization have access to all images regardless of the settings.

5. Select the category restrictions for your images by selecting **Allowed**.
     Once the restrictions are applied, your members can view the organization permissions page in a read-only format.
6. Optional: To ensure that each organization member uses images in a safe and secure environment, [enfore sign-in](../docker-hub/configure-sign-in.md).

## Verify the restrictions

   To confirm that the restrictions are successful, have each organization member pull an image onto their local computer after signing in to Docker Desktop. If they don't sign in, they receive an error message.

   For example, if you enable Image Access Management, your members can only pull an Organization Image, Docker Official Image, or Verified Publisher Image onto their local machine. If you disable the restrictions, your members can pull any image, including community images.
