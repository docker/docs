---
description: Learn how to get started using Docker Hub
keywords: Docker Hub, push image, pull image, repositories
title: Docker Hub quickstart
linkTitle: Quickstart
weight: 10
---

Docker Hub provides a vast library of pre-built images and resources,
accelerating development workflows and reducing setup time. You can build upon
pre-built images from Docker Hub and then use repositories to share and
distribute your own images with your team or millions of other developers.

This guide shows you how to find and run a pre-built image. It then walks you
through creating a custom image and sharing it through Docker Hub.

## Prerequisites

- [Download and install Docker](../../get-started/get-docker.md)
- [Create a Docker account](https://app.docker.com/signup)

## Step 1: Find an image in Docker Hub's library

You can search for content in Docker Hub itself, in the Docker Desktop
Dashboard, or by using the CLI.

To search or browse for content on Docker Hub:

{{< tabs >}}
{{< tab name="Docker Hub" >}}

1. Navigate to the [Docker Hub Explore page](https://hub.docker.com/explore).

   On the **Explore** page, you can browse by catalog or category, or use the search
   to quickly find content.

2. Under **Categories**, select **Web servers**.

   After the results are displayed, you can further filter the results using the
   filters on the left side of the page.

3. In the filters, select **Docker Official Image**.

   Filtering by Trusted Content ensures that you see only high-quality, secure
   images curated by Docker and verified publishing partners.

4. In the results, select the **nginx** image.

   Selecting the image opens the image's page where you can learn more about how
   to use the image. On the page, you'll also find the `docker pull` command to
   pull the image.

{{< /tab >}}
{{< tab name="Docker Desktop" >}}

1. Open the Docker Desktop Dashboard.
2. Select the **Docker Hub** view.

   In the **Docker Hub** view, you can browse by catalog or category, or use the search
   to quickly find content.

3. Leave the search box empty and then select **Search**.

   The search results are shown with additional filters now next to the search box.

4. Select the search filter icon, and then select **Docker Official Image** and **Web Servers**.
5. In the results, select the **nginx** image.

{{< /tab >}}
{{< tab name="CLI" >}}

1. Open a terminal window.

   > [!TIP]
   >
   > The Docker Desktop Dashboard contains a built-in terminal. At the bottom of
   > the Dashboard, select **>_ Terminal** to open it.

2. In the terminal, run the following command.

   ```console
   $ docker search --filter is-official=true nginx
   ```

   Unlike the Docker Hub and Docker Desktop interfaces, you can't browse by
   category using the `docker search` command. For more details about the
   command, see [docker search](/reference/cli/docker/search/).

{{< /tab >}}
{{< /tabs >}}

Now that you've found an image, it's time to pull and run it on your device.

## Step 2: Pull and run an image from Docker Hub

You can run images from Docker Hub using the CLI or Docker Desktop Dashboard.

{{< tabs >}}
{{< tab name="Docker Desktop" >}}

1. In the Docker Desktop Dashboard, select the **nginx** image in the **Docker
   Hub** view. For more details, see [Step 1: Find an image in Docker Hub's
   library](#step-1-find-an-image-in-docker-hubs-library).

2. On the **nginx** screen, select **Run**.

   If the image doesn't exist on your device, it is automatically pulled from
   Docker Hub. Pulling the image may take a few seconds or minutes depending on
   your connection. After the image has been pulled, a window appears in Docker
   Desktop and you can specify run options.

3. In the **Host port** option, specify `8080`.
4. Select **Run**.

   The container logs appear after the container starts.

5. Select the **8080:80** link to open the server, or visit
   [https://localhost:8080](https://localhost:8080) in your web browser.

6. In the Docker Desktop Dashboard, select the **Stop** button to stop the
   container.


{{< /tab >}}
{{< tab name="CLI" >}}

1. Open a terminal window.

   > [!TIP]
   >
   > The Docker Desktop Dashboard contains a built-in terminal. At the bottom of
   > the Dashboard, select **>_ Terminal** to open it.

2. In your terminal, run the following command to pull and run the Nginx image.

   ```console
   $ docker run -p 8080:80 --rm nginx
   ```

   The `docker run` command automatically pulls and runs the image without the
   need to run `docker pull` first. To learn more about the command and its
   options, see the [`docker run` CLI
   reference](../../reference/cli/docker/container/run.md). After running the
   command, you should see output similar to the following.

   ```console {collapse=true}
   Unable to find image 'nginx:latest' locally
   latest: Pulling from library/nginx
   a480a496ba95: Pull complete
   f3ace1b8ce45: Pull complete
   11d6fdd0e8a7: Pull complete
   f1091da6fd5c: Pull complete
   40eea07b53d8: Pull complete
   6476794e50f4: Pull complete
   70850b3ec6b2: Pull complete
   Digest: sha256:28402db69fec7c17e179ea87882667f1e054391138f77ffaf0c3eb388efc3ffb
   Status: Downloaded newer image for nginx:latest
   /docker-entrypoint.sh: /docker-entrypoint.d/ is not empty, will attempt to perform configuration
   /docker-entrypoint.sh: Looking for shell scripts in /docker-entrypoint.d/
   /docker-entrypoint.sh: Launching /docker-entrypoint.d/10-listen-on-ipv6-by-default.sh
   10-listen-on-ipv6-by-default.sh: info: Getting the checksum of /etc/nginx/conf.d/default.conf
   10-listen-on-ipv6-by-default.sh: info: Enabled listen on IPv6 in /etc/nginx/conf.d/default.conf
   /docker-entrypoint.sh: Sourcing /docker-entrypoint.d/15-local-resolvers.envsh
   /docker-entrypoint.sh: Launching /docker-entrypoint.d/20-envsubst-on-templates.sh
   /docker-entrypoint.sh: Launching /docker-entrypoint.d/30-tune-worker-processes.sh
   /docker-entrypoint.sh: Configuration complete; ready for start up
   2024/11/07 21:43:41 [notice] 1#1: using the "epoll" event method
   2024/11/07 21:43:41 [notice] 1#1: nginx/1.27.2
   2024/11/07 21:43:41 [notice] 1#1: built by gcc 12.2.0 (Debian 12.2.0-14)
   2024/11/07 21:43:41 [notice] 1#1: OS: Linux 6.10.11-linuxkit
   2024/11/07 21:43:41 [notice] 1#1: getrlimit(RLIMIT_NOFILE): 1048576:1048576
   2024/11/07 21:43:41 [notice] 1#1: start worker processes
   2024/11/07 21:43:41 [notice] 1#1: start worker process 29
   ...
   ```

3. Visit [https://localhost:8080](https://localhost:8080) to view the default
   Nginx page and verify that the container is running.

4. In the terminal, press <kdb>Ctrl+C</kbd> to stop the container.

{{< /tab >}}
{{< /tabs >}}

You've now run a web server without any set up or configuration. Docker Hub
provides instant access to pre-built, ready-to-use container images, letting you
quickly pull and run applications without needing to install or configure
software manually. With Docker Hub's vast library of images, you can experiment
with and deploy applications effortlessly, boosting productivity and making it
easy to try out new tools, set up development environments, or build on top of
existing software.

You can also extend images from Docker Hub, letting you quickly build and
customize your own images to suit specific needs.


## Step 3: Build and push an image to Docker Hub

1. Create a [Dockerfile](/reference/dockerfile.md) to specify your application:

   ```dockerfile
   FROM nginx
   RUN echo "<h1>Hello world from Docker!</h1>" > /usr/share/nginx/html/index.html
   ```

   This Dockerfile extends the Nginx image from Docker Hub to create a
   simple website. With just a few lines, you can easily set up, customize, and
   share a static website using Docker.

2. Run the following command to build your image. Replace `<YOUR-USERNAME>` with your Docker ID.

   ```console
   $ docker build -t <YOUR-USERNAME>/nginx-custom .
   ```

   This command builds your image and tags it so that Docker understands which
   repository to push it to in Docker Hub. To learn more about the command and
   its options, see the [`docker build` CLI
   reference](../../reference/cli/docker/buildx/build.md). After running the
   command, you should see output similar to the following.

   ```console {collapse=true}
   [+] Building 0.6s (6/6) FINISHED                      docker:desktop-linux
    => [internal] load build definition from Dockerfile                  0.0s
    => => transferring dockerfile: 128B                                  0.0s
    => [internal] load metadata for docker.io/library/nginx:latest       0.0s
    => [internal] load .dockerignore                                     0.0s
    => => transferring context: 2B                                       0.0s
    => [1/2] FROM docker.io/library/nginx:latest                         0.1s
    => [2/2] RUN echo "<h1>Hello world from Docker!</h1>" > /usr/share/  0.2s
    => exporting to image                                                0.1s
    => => exporting layers                                               0.0s
    => => writing image sha256:f85ab68f4987847713e87a95c39009a5c9f4ad78  0.0s
    => => naming to docker.io/mobyismyname/nginx-custom                  0.0s
   ```

3. Run the following command to test your image. Replace `<YOUR-USERNAME>` with
   your Docker ID.

   ```console
   $ docker run -p 8080:80 --rm <YOUR-USERNAME>/nginx-custom
   ```

4. Visit [https://localhost:8080](https://localhost:8080) to view the page. You
   should see `Hello world from Docker!`.

5. In the terminal, press CTRL+C to stop the container.

6. Sign in to Docker Desktop. You must be signed in before pushing an image to
   Docker Hub.

7. Run the following command to push your image to Docker Hub. Replace `<YOUR-USERNAME>` with your Docker ID.

   ```console
   $ docker push <YOUR-USERNAME>/nginx-custom
   ```

    > [!NOTE]
    >
    > You must be signed in to Docker Hub through Docker Desktop or the command line, and you must also name your images correctly, as per the above steps.

   The command pushes the image to Docker Hub and automatically
   creates the repository if it doesn't exist. To learn more about the command,
   see the [`docker push` CLI
   reference](../../reference/cli/docker/image/push.md). After running the
   command, you should see output similar to the following.

   ```console {collapse=true}
   Using default tag: latest
   The push refers to repository [docker.io/mobyismyname/nginx-custom]
   d0e011850342: Pushed
   e4e9e9ad93c2: Mounted from library/nginx
   6ac729401225: Mounted from library/nginx
   8ce189049cb5: Mounted from library/nginx
   296af1bd2844: Mounted from library/nginx
   63d7ce983cd5: Mounted from library/nginx
   b33db0c3c3a8: Mounted from library/nginx
   98b5f35ea9d3: Mounted from library/nginx
   latest: digest: sha256:7f5223ae866e725a7f86b856c30edd3b86f60d76694df81d90b08918d8de1e3f size: 1985
   ```

  Now that you've created a repository and pushed your image, it's time to view
  your repository and explore its options.

## Step 4: View your repository on Docker Hub and explore options

You can view your Docker Hub repositories in the Docker Hub or Docker Desktop interface.

{{< tabs >}}
{{< tab name="Docker Hub" >}}

1. Go to [Docker Hub](https://hub.docker.com) and sign in.

   After signing in, you should be on the **Repositories** page. If not, then go
   to the [**Repositories**](https://hub.docker.com/repositories/) page.

2. Find the **nginx-custom** repository and select that row.

   After selecting the repository, you should see more details and options for
   your repository.

{{< /tab >}}
{{< tab name="Docker Desktop" >}}

1. Sign in to Docker Desktop.
2. Select the **Images** view.
3. Select the **Hub repositories** tab.

   A list of your Docker Hub repositories appears.

4. Find the **nginx-custom** repository, hover over the row, and then select **View in Hub**.

   Docker Hub opens and you are able to view more details about the image.

{{< /tab >}}
{{< /tabs >}}

You've now verified that your repository exists on Docker Hub, and you've
discovered more options for it. View the next steps to learn more about some of
these options.

## Next steps

Add [repository information](./repos/manage/information.md) to help users find and use
your image.

