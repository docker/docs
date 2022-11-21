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

## Create a repo

To push an image, you first need to create a repository on Docker Hub.

1. [Sign up](https://www.docker.com/pricing?utm_source=docker&utm_medium=webreferral&utm_campaign=docs_driven_upgrade){:target="_blank" rel="noopener" class="_"} or Sign in to [Docker Hub](https://hub.docker.com){:target="_blank" rel="noopener" class="_"}.

2. Select the **Create Repository** button.

3. For the repo name, use `getting-started`. Make sure the **Visibility** is **Public**.

    > **Private repositories**
    >
    > Did you know that Docker offers private repositories which allows you to restrict content to specific users or teams? Check out the details on the [Docker pricing](https://www.docker.com/pricing?utm_source=docker&utm_medium=webreferral&utm_campaign=docs_driven_upgrade){:target="_blank" rel="noopener" class="_"} page.

4. Select the **Create** button.

After creating the repository, the repository page appears. On the page, you can view the **Docker commands** section. In this section, you'll see the command to push an image to this repository.


## Build a multi-platform image and push it to Docker Hub

The image that you have built will only run on platforms using the same architecture as your development machine. If you want to share and run the image on machines with different architectures, you can use buildx to build a [multi-platform image](../build/building/multi-platform.md).

In the following steps, you will build a multi-platform image that can run on AMD64 and ARM64/v8.

1. In a terminal, run the following command to create and use a new builder with the `docker-container` driver which gives you access to more complex features like multi-platform builds.

   ```console
   $ docker buildx create --name mybuilder --driver docker-container --bootstrap --use
   ```

2. In a terminal, log in to Docker Hub using  the `docker login` command. Replace `<your-docker-id>` with your Docker ID.

   ```console
   $ docker login -u <your-docker-id>
   ```

3. In a terminal, change directory to the directory containing your Dockerfile and then run the following command to build a multi-platform image and push it to Docker Hub. Replace <`your-docker-id>` with your Docker ID.

   ```console
    $ docker buildx build --platform linux/amd64,linux/arm/v8 -t <your-docker-id>/getting-started --push .
   ```
   In the command above, you use `--platform` to specify the platforms for the image and you use `--push` to push the image to Docker Hub.

## Run the image on a new instance

Now that your image has been built and pushed into a registry, try running your app on a brand new instance that has never seen this container image! To do this, you will use Play with Docker.

1. Open your browser to [Play with Docker](https://labs.play-with-docker.com/){:target="_blank" rel="noopener" class="_"}.

2. Select **Login** and then select **docker** from the drop-down list.

3. Connect with your Docker Hub account.

4. Once you're logged in, select the **ADD NEW INSTANCE** option on the left side bar. If you don't see it, make your browser a little wider. After a few seconds, a terminal window opens in your browser.

5. In the terminal, start your freshly pushed app. Replace `<your-docker-id>` with your Docker ID.

   ```console
   $ docker run -dp 3000:3000 <your-docker-id>/getting-started
   ```

   You should see the image get pulled down and eventually start up.

6. Select the 3000 badge when it comes up and you should see the app with your modifications.
    If the 3000 badge doesn't show up, you can select the **Open Port** button and type in 3000.


## Next steps

In this section, you learned how to share your images by pushing them to a registry. You then went to a
brand new instance and were able to run the freshly pushed image. This is quite common in CI pipelines,
where the pipeline will create the image and push it to a registry and then the production environment
can use the latest version of the image.

In the next part,  you'll learn how you can persist data across container restarts.

[Persist the DB](05_persisting_data.md){: .button  .primary-btn}