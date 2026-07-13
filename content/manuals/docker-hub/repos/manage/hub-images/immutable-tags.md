---
description: Learn about immutable tags and how they help maintain image version consistency on Docker Hub.
keywords: Docker Hub, Hub, repository content, tags, immutable tags, version control
title: Immutable tags on Docker Hub
linkTitle: Immutable tags
weight: 11
---
{{< summary-bar feature_name="Immutable tags" >}}

Immutable tags provide a way to ensure that specific image versions remain unchanged once they are published to Docker Hub. This feature helps maintain consistency and reliability in your container deployments by preventing accidental overwrites of important image versions.

## What are immutable tags?

Immutable tags are image tags that, once pushed to Docker Hub, cannot be overwritten or deleted. This ensures that a specific version of an image remains exactly the same throughout its lifecycle, providing:

- Version consistency
- Reproducible builds
- Protection against accidental overwrites
- Better security and compliance

## Enable immutable tags

To enable immutable tags for your repository:

1. Sign in to [Docker Hub](https://hub.docker.com).
2. Select **My Hub** > **Repositories**.
3. Select the repository where you want to enable immutable tags.
4. Go to **Settings** > **General**.
5. Under **Tag mutability settings**, select one of the following options:
   - **All tags are mutable (Default)**:  
     Tags can be changed to reference a different image. This lets you retarget a tag without creating a new one.
   - **All tags are immutable**:  
     Tags cannot be updated to point to a different image after creation. This ensures consistency and prevents accidental changes. This includes the `latest` tag.
   - **Specific tags are immutable**:  
     Define specific tags that cannot be updated after creation using regex values.
6. Select **Save**.

Once enabled, all tags are locked to their specific images, ensuring that each tag always points to the same image version and cannot be modified.

> [!NOTE]
> This implementation of regular expressions follows the [Go regexp package](https://pkg.go.dev/regexp), which is based on the RE2 engine. For more information, visit [RE2 Regular Expression Syntax](https://github.com/google/re2/wiki/Syntax).

## Working with immutable tags

When immutable tags are enabled:

- You cannot push a new image with the same tag name
- You must use a new tag name for each new image version

To push an image, create a new tag for your updated image and push it to the repository.










