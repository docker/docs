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

## Run the image on a new instance

Now that your image has been built and pushed to a registry, any device with Docker can run your image.

To run the image on another device with Docker, use the `docker run` command and replace `<your-docker-id>` with your Docker ID.

   ```console
   $ docker run -dp 3000:3000 <your-docker-id>/getting-started
   ```

If you don't have another device to try, you can delete the image from your device, which will allow you to simulate running the image on a new device.

To remove the image, you must first stop and remove the container.

1. Get the ID of the container by using the `docker ps` command. Use the `--all` flag to list all running and stopped containers.

   ```console
   $ docker ps --all
   ```

2. If the container is running, use the `docker stop` command to stop the container. Replace `<the-container-id>` with the ID from `docker ps`.

   ```console
   $ docker stop <the-container-id>
   ```

3. Once the container has stopped, you can remove it by using the `docker rm` command. Replace `<the-container-id>` with the ID from `docker ps`.

   ```console
   $ docker rm <the-container-id>
   ```

4. Now that the container has been removed, you can remove the image. Get the ID of the image by using the `docker image ls` command.

   ```console
   $ docker image ls
   ```

5. Remove the image using the `docker image rm` command. Replace `<the-image-id>` with the ID from `docker image ls`.

   ```console
   $ docker image rm <the-image-id>
   ```

   Your device no longer has the image. In order to run the image, you can rebuild it, or pull it from Docker Hub.

6. Pull and run the image from Docker Hub using the `docker run` command. Replace `<your-docker-id>` with your Docker ID.

   ```console
   $ docker run -dp 3000:3000 <your-docker-id>/getting-started
   ```

   Docker pulls the image from Docker Hub and runs it.

7. After a few seconds, open your web browser to [http://localhost:3000](http://localhost:3000){:target="_blank" rel="noopener" class="_"}.
   You should see your app.

## Next steps

In this section, you learned how to share an image by pushing it to a registry. You then pulled the image from the registry and ran it. Using images from a registry is quite common in CI pipelines, where the pipeline will create the image and push it to a registry and then the production environment can use the latest version of the image.

In the next part, you'll learn how you can persist data across container restarts.

[Persist the DB](05_persisting_data.md){: .button  .primary-btn}