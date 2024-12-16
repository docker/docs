---
description: Learn how to add content to a repository on Docker Hub.
keywords: Docker Hub, Hub, repository content, push
title: Push images to a repository
linkTitle: Push images
weight: 30
---

To add content to a repository on Docker Hub, you need to tag your Docker image
and then push it to your repository. This process lets you share your
images with others or use them in different environments.

1. Tag your Docker image.

   The `docker tag` command assigns a tag to your Docker image, which includes
   your Docker Hub namespace and the repository name. The general syntax is:

   ```console
   $ docker tag [SOURCE_IMAGE[:TAG]] [NAMESPACE/REPOSITORY[:TAG]]
   ```

   Example:

   If your local image is called `my-app` and you want to tag it for the
   repository `my-namespace/my-repo` with the tag `v1.0`, run:

   ```console
   $ docker tag my-app my-namespace/my-repo:v1.0
   ```

2. Push the image to Docker Hub.

   Use the `docker push` command to upload your tagged image to the specified
   repository on Docker Hub.

   Example:

   ```console
   $ docker push my-namespace/my-repo:v1.0
   ```

   This command pushes the image tagged `v1.0` to the `my-namespace/my-repo` repository.

3. Verify the image on Docker Hub.
