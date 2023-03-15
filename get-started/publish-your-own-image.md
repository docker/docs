---
title: Publish your own image
keywords: get started, quick start, intro, concepts
description: Learn how to publish your own images
---

In this guide, you'll learn how you can share your packaged application in an image using Docker Hub.

## Step 1: Get an image

Before you publish your image, you'll need an image to publish. For this guide, you'll use the `welcome-to-docker` image. To get the image, run the following `docker pull` command in a terminal to pull the image from Docker Hub.

```console
$ docker pull docker/welcome-to-docker
```

## Step 2: Sign in to Docker

To publish images publicly on Docker Hub, you first need an account. Select **Sign in** on the top-right of Docker Desktop to either sign in or create a new account on Docker Hub.

![Signing in to Docker](images/getting-started-sign-in.png)

## Step 3: Rename your image

Before you can publish your image to Docker Hub, you need to rename it so that Docker Hub knows that the image is yours. Run the following `docker tag` command in your terminal to rename your image. Replace `YOUR-USERNAME` with your Docker ID.

```console
$ docker tag docker/welcome-to-docker YOUR-USERNAME/welcome-to-docker
```

## Step 4: Push your image to Docker Hub

In Docker Desktop, go to the **Images** tab and find your image. In the **Actions** column, select the **More image actions** icon and then select **Push to Hub**. Your image will upload to Docker Hub and be publicly available for anyone to use.

## Step 5: Verify the image is on Docker Hub

THat's it! Your image is now being shared on Docker Hub. In browser, go to [Docker Hub](https://hub.docker.com) and verify that you see the `welcome-to-docker` repository.

![Signing in to Docker](images/getting-started-push.gif)