---
description: Advanced Image Management dashboard
keywords: dashboard, images, image management, inactive
title: Advanced Image Management dashboard
---

> **Deprecated**
>
> The Advanced Image Management dashboard and API are deprecated, and scheduled
> for removal on November 15th, 2023. You can continue to use the **Tags** in
> Docker Hub to manage tags for your repository.
>
> For more information, see [Deprecation of Advanced Image Management](https://github.com/docker/roadmap/issues/534)
{ .warning }

Advanced Image Management allows you to manage Docker images across all repositories and streamline storage in Docker Hub.

It provides:

- A snapshot of your existing images
- Allows you to view, sort, and filter images by tags, activity status, and date
- Contains options to clean up your workspace by deleting images that are no longer required

## Access the Advanced Image Management dashboard

1. Sign in to [Docker Hub](https://hub.docker.com).
2. Select **Repositories**.
3. Choose a repository.
4. Select **General** or **Tags**, and then select **Go to Advanced Image Management**.

## Understand image activity status and tags

An image retains its 'active' status if it's pulled or pushed in the last month. If there isn’t any activity on the image in the last month, it's considered 'inactive'.

The dashboard also displays the old versions of images you have pushed. When you push an image to Docker Hub, you push a manifest, which is a list of all the layers of your image and the layers themselves. When you update an existing tag, only the new layers are pushed along with the new manifest which references the new layers. This new manifest gets the tag you specify when you push the image, such as `myNamespace/mytag:latest`. This doesn't remove the old manifests or the unique layers referenced by them from Hub. You can still use and reference these using the digest of the manifest if you know the SHA.

## Deleting a tagged image

A Docker image can contain multiple tags. A tag refers to a combination of artifacts or labels associated with the image. When you attempt to delete a tagged image, it deletes the tag associated with the image. This means, if there are any untagged images in your repository that previously held the same tag, those images will also be deleted even if they are active. Therefore, you must be careful when deleting tagged images.

For example, let's assume that Image A is tagged as ‘latest’. You push another image, Image B, and tag it as the new 'latest'. If you now delete Image-B, the 'latest' tag will be deleted, along with all images it previously referred to. However, if those images are tagged by another tag - in this case, if Image-B is also tagged with '1.5.0', for example, that tag would prevent its deletion.

## Advanced Image Management API

The Advanced Image Management API endpoints allow you to manage Docker images across all repositories. For more information, see [Advanced Image management API](./api/latest.md).
