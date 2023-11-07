---
title: Publish your image
keywords: get started, quick start, intro, concepts
description: Learn how to publish your image to Docker Hub
aliases:
- /get-started/publish-your-own-image/
---

Follow this walkthrough to learn how to publish and share your images on Docker Hub.

Before you start, [get Docker Desktop](../../get-docker.md).

## Step 1: Get the example image

To get the example image, do the following:

1. In Docker Desktop, select the search bar.
2. In the search bar, specify `docker/welcome-to-docker`.
3. Select **Pull** to pull the image from Docker Hub to you computer.

![Search Docker Desktop for the welcome-to-docker image](images/getting-started-search.png?w=400)

## Step 2: Sign in to Docker

Select **Sign in** on the top-right of Docker Desktop to either sign in or create a new Docker account.

![Signing in to Docker Desktop](images/getting-started-signin.png?w=300)

## Step 3: Rename your image

Before you can publish your image, you need to rename it so that Docker Hub knows that the image is yours. In a terminal, run the following command to rename your image. Replace `YOUR-USERNAME` with your Docker ID.

{{< include "open-terminal.md" >}}

```console
$ docker tag docker/welcome-to-docker YOUR-USERNAME/welcome-to-docker
```

## Step 4: Push your image to Docker Hub

To push your image to Docker Hub, do the following:

1. In Docker Desktop, go to the **Images** tab
2. In the **Actions** column for your image, select the **Show image actions** icon.
3. Select **Push to Hub**.

![Pushing an image to Docker Hub](images/getting-started-push.png?w=200&border=true)

Go to [Docker Hub](https://hub.docker.com)⁠ and verify that the list of your repositories now contains `YOU-USERNAME/welcome-to-docker`.

## Summary

In this walkthrough, you pushed and shared an image on Docker Hub.

Related information:

- Deep dive into the [Docker Hub manual](../../docker-hub/_index.md)
- Learn more about the [docker tag](../../engine/reference/commandline/tag.md)
  command

## Next steps

Continue to the language-specific guides to learn how you can use Docker to containerize and develop applications in your favorite language. Choose one of the following guides.

<div class="grid grid-cols-3 auto-rows-fr sm:flex-col sm:h-auto gap-4">
    <div class="flex items-center flex-1 shadow p-4">
        <a href="/language/nodejs/"><img class="m-auto rounded" src="/language/images/nodejs.png" alt="Develop with Node"></a>
    </div>
    <div class="flex items-center flex-1 shadow p-4">
        <a href="/language/python/"><img class="m-auto rounded" src="/language/images/python.png" alt="Develop with Python"></a>
    </div>
    <div class="flex items-center flex-1 shadow p-4">
        <a href="/language/java/"><img class="m-auto rounded" src="/language/images/java.png" alt="Develop with Java"></a>
    </div>
    <div class="flex items-center flex-1 shadow p-4">
        <a href="/language/golang/"><img class="m-auto rounded" src="/language/images/golang.png" alt="Develop with Go"></a>
    </div>
    <div class="flex items-center flex-1 shadow p-4">
        <a href="/language/dotnet/"><img class="m-auto rounded" src="/language/images/c-sharp.png" alt="Develop with C#"></a>
    </div>
    <div class="flex items-center flex-1 shadow p-4">
        <a href="/language/rust/"><img class="m-auto rounded" src="/language/images/rust-logo.png" alt="Develop with Rust"></a>
    </div>
</div>