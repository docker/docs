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

To push an image, we first need to create a repository on Docker Hub.

1. [Sign up](https://www.docker.com/pricing?utm_source=docker&utm_medium=webreferral&utm_campaign=docs_driven_upgrade){:target="_blank" rel="noopener" class="_"} or Sign in to [Docker Hub](https://hub.docker.com){:target="_blank" rel="noopener" class="_"}.

2. Select the **Create Repository** button.

3. For the repo name, use `getting-started`. Make sure the **Visibility** is **Public**.

    > **Private repositories**
    >
    > Did you know that Docker offers private repositories which allows you to restrict content to specific users or teams? Check out the details on the [Docker pricing](https://www.docker.com/pricing?utm_source=docker&utm_medium=webreferral&utm_campaign=docs_driven_upgrade){:target="_blank" rel="noopener" class="_"} page.

4. Select the **Create** button.

After creating the repository, the repository page appears. On the page, you can view the **Docker commands** section. In this section, you'll see the command to push an image to this repository.

## Push the image

1. In the command line, try running the `docker push` command. Replace `<your-docker-id>` with your Docker ID.

   ```console
    $ docker push <your-docker-id>/getting-started
   ```
   You will see output similar to the following:
   ```console
   The push refers to repository [docker.io/<your-docker-id>/getting-started]
   An image does not exist locally with the tag: <your-docker-id>/getting-started
   ```

   Why did it fail? The push command was looking for an image named `<your-docker-id>/getting-started`, but didn't find one. If you run `docker image ls`, you won't see one either. You will see your image name is `getting-started`.

    To push the image to the repository, you need to "tag" your existing image you've built to give it another name that is prepended with your Docker ID.

2. In a terminal, log in to Docker Hub using the `docker login` command. Replace `<your-docker-id>` with your Docker ID.

   ```console
   $ docker login -u <your-docker-id>
   ```

3. Use the `docker tag` command to rename your `getting-started` image to `<your-docker-id>/getting-started`. Replace `<your-docker-id>` with your Docker ID.

    ```console
    $ docker tag getting-started <your-docker-id>/getting-started
    ```

4. Now try the `docker push` command again. Replace `<your-docker-id>` with your Docker ID.

   ```console
   $ docker push <your-docker-id>/getting-started
   ```

## Run the image on a new instance

Now that your image has been built and pushed into a registry, try running your app on a brand new instance that has never seen this container image. To do this, you will use Play with Docker.


If you are using an Apple silicon device, you must first build a new image because your device's platform differs from the Play with Docker platform. Select your device's platform below.

<ul class="nav nav-tabs">
  <li class="active"><a data-toggle="tab" data-target="#amd">Mac / Linux / Windows with AMD64</a></li>
  <li><a data-toggle="tab" data-target="#arm">Mac with Apple silicon</a></li>
</ul>
<div class="tab-content">
<div id="amd" class="tab-pane fade in active" markdown="1">

### Mac / Linux / Windows with AMD64

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

<hr>
</div>
<div id="arm" class="tab-pane fade" markdown="1">

### Mac with Apple silicon

In the steps below, you will build an additional image that's compatible with the Play with Docker platform. You can also build a single [multi-platform image](./build/building/multi-platform), but that's outside of the scope of this tutorial.

1. In a terminal, change directory to the directory containing your Dockerfile and then run the following command to build a new image that's compatible with the Play with Docker platform. Replace <`your-docker-id>` with your Docker ID.

   ```console
    $ docker build --platform linux/amd64 -t <your-docker-id>/getting-started:amd64 .
   ```
   In the command above, you use `--platform` to specify the platform for the image and you use `-t <your-docker-id>/getting-started:amd64` to name the new image.

2. In a terminal, log in to Docker Hub using  the `docker login` command. Replace `<your-docker-id>` with your Docker ID.

   ```console
   $ docker login -u <your-docker-id>
   ```

3. Now use `docker push` to push the image to Docker Hub. Replace `<your-docker-id>` with your Docker ID.

   ```console
   $ docker push <your-docker-id>/getting-started:amd64
   ```

4. Open your browser to [Play with Docker](https://labs.play-with-docker.com/){:target="_blank" rel="noopener" class="_"}.

5. Select **Login** and then select **docker** from the drop-down list.

6. Connect with your Docker Hub account.

7. Once you're logged in, select the **ADD NEW INSTANCE** option on the left side bar. If you don't see it, make your browser a little wider. After a few seconds, a terminal window opens in your browser.

8. In the terminal, start your freshly pushed app. Replace `<your-docker-id>` with your Docker ID.

    ```console
    $ docker run -dp 3000:3000 <your-docker-id>/getting-started:amd64
    ```

    You should see the image get pulled down and eventually start up.

9. Select the 3000 badge when it comes up and you should see the app with your modifications.
    If the 3000 badge doesn't show up, you can select the **Open Port** button and type in 3000.


<hr>
</div>
</div>

## Next steps

In this section, you learned how to share your images by pushing them to a registry. You then went to a
brand new instance and were able to run the freshly pushed image. This is quite common in CI pipelines,
where the pipeline will create the image and push it to a registry and then the production environment
can use the latest version of the image.

In the next part,  you'll learn how you can persist data across container restarts.

[Persist the DB](05_persisting_data.md){: .button  .primary-btn}