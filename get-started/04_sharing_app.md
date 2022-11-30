---
title: "Share the application"
keywords: get started, setup, orientation, quickstart, intro, concepts, containers, docker desktop, docker hub, sharing 
redirect_from:
- /get-started/part3/
description: Sharing the image you built for your example application so you can run it else where and other developers can use it
---

Now that you've built an image, you can share it. To share Docker images, you have to use a Docker registry. The default registry is Docker Hub and is where all of the images you've used have come from.

> **Docker ID**
>
> A Docker ID allows you to access Docker Hub which is the world's largest library and community for container images. Create a [Docker ID](https://hub.docker.com/signup){:target="_blank" rel="noopener" class="_"} for free if you don't have one.


## Create a repository

To push an image, you first need to create a repository on Docker Hub.

1. [Sign up](https://www.docker.com/pricing?utm_source=docker&utm_medium=webreferral&utm_campaign=docs_driven_upgrade){:target="_blank" rel="noopener" class="_"} or sign in to [Docker Hub](https://hub.docker.com){:target="_blank" rel="noopener" class="_"}.

2. Select the **Create Repository** button.

3. For the repository name, use `getting-started`. Make sure the **Visibility** is **Public**.

    > **Private repositories**
    >
    > Did you know that Docker offers private repositories which allows you to restrict content to specific users or teams? Check out the details on the [Docker pricing](https://www.docker.com/pricing?utm_source=docker&utm_medium=webreferral&utm_campaign=docs_driven_upgrade){:target="_blank" rel="noopener" class="_"} page.

4. Select the **Create** button.

   After creating the repository, the repository page appears. On the page, you can view the **Docker commands** section. In this section, you'll see the command to push an image to this repository.

## Push the image

1. In a terminal, log in to Docker Hub using the `docker login` command. Replace `<your-docker-id>` with your Docker ID.

   ```console
   $ docker login -u <your-docker-id>
   ```

2. Use the `docker tag` command to rename the `getting-started` image to `<your-docker-id>/getting-started`. To push the image to the repository, you must prepended the image name with your Docker ID. Replace `<your-docker-id>` with your Docker ID.

   ```console
   $ docker tag getting-started <your-docker-id>/getting-started
   ```

3. Push the image to the Docker Hub registry using the `docker push` command. Replace `<your-docker-id>` with your Docker ID.

   ```console
   $ docker push <your-docker-id>/getting-started
   ```

## Pull and run an image

Now that your image has been built and pushed to a registry, any device with Docker can run your image.

To run the image on another device with Docker, use the `docker run` command and replace `<your-docker-id>` with your Docker ID.

   ```console
   $ docker run -dp 3000:3000 <your-docker-id>/getting-started
   ```

You can also pull and run other images from a registry on your device. In the following steps, you will pull and run the [hello-word](https://hub.docker.com/_/hello-world){:target="_blank" rel="noopener" class="_"} image from Docker Hub.

Use the `docker run` command to pull and run the image.

   ```console
   $ docker run hello-world
   ```

You will see output similar to the following, indicating that Docker pulled the image from the registry and ran it on your device.

   ```plaintext
   Unable to find image 'hello-world:latest' locally
   latest: Pulling from library/hello-world
   2db29710123e: Pull complete
   Digest: sha256:faa03e786c97f07ef34423fccceeec2398ec8a5759259f94d99078f264e9d7af
   Status: Downloaded newer image for hello-world:latest

   Hello from Docker!
   This message shows that your installation appears to be working correctly.

   To generate this message, Docker took the following steps:
    1. The Docker client contacted the Docker daemon.
    2. The Docker daemon pulled the "hello-world" image from the Docker Hub.
       (amd64)
    3. The Docker daemon created a new container from that image which runs the
       executable that produces the output you are currently reading.
    4. The Docker daemon streamed that output to the Docker client, which sent it to your terminal.

   To try something more ambitious, you can run an Ubuntu container with:
    $ docker run -it ubuntu bash

   Share images, automate workflows, and more with a free Docker ID:
    https://hub.docker.com/

   For more examples and ideas, visit:
    https://docs.docker.com/get-started/
   ```

Take some time to explore other images on [Docker Hub](https://hub.docker.com/search){:target="_blank" rel="noopener" class="_"}. Docker Hub is the world's largest repository of container images with an array of content sources including container community developers, open source projects and independent software vendors (ISV) building and distributing their code in containers. Store and share your personal projects and see what the container community is building.

## Next steps

In this section, you learned how to share an image by pushing it to a registry. You then pulled the image from the registry and ran it. Using images from a registry is quite common in CI pipelines, where the pipeline will create the image and push it to a registry and then the production environment can use the latest version of the image.

In the next part, you'll learn how you can persist data across container restarts.

[Persist the DB](05_persisting_data.md){: .button  .primary-btn}