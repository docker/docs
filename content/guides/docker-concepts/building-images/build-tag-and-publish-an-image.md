---
title: Build, tag, and publish an image
keywords: concepts, build, images, container, docker desktop
description: This concept page will teach you how to build, tag, and publish an image to Docker Hub or any other registry 
---

{{< youtube-embed chiiGLlYRlY >}}

## Explanation

In this concept, you will learn the following:

- Building images - the process of building an image based on a `Dockerfile`
- Tagging images - the process of giving an image a name, which also determines where the image can be distributed
- Publishing images - the process to distribute or share the newly created image using a container registry

### Building images

Most often, images are built using a `Dockerfile`. The most basic `docker build` command might look like the following:

```bash
docker build .
```

The final `.` in the command provides the path or URL to the [build context](https://docs.docker.com/build/building/context/#what-is-a-build-context). At this location, the builder will find the `Dockerfile` and other referenced files.

When a build occurs, the builder will pull the base image if needed and then run all of the instructions as outlined in the `Dockerfile`.

With the previous command, the image will have no name, but the output will provide the ID of the image. As an example, the previous command might produce the following output:

```console
$ docker build .
[+] Building 3.5s (11/11) FINISHED                                              docker:desktop-linux
 => [internal] load build definition from Dockerfile                                            0.0s
 => => transferring dockerfile: 308B                                                            0.0s
 => [internal] load metadata for docker.io/library/python:3.12                                  0.0s
 => [internal] load .dockerignore                                                               0.0s
 => => transferring context: 2B                                                                 0.0s
 => [1/6] FROM docker.io/library/python:3.12                                                    0.0s
 => [internal] load build context                                                               0.0s
 => => transferring context: 123B                                                               0.0s
 => [2/6] WORKDIR /usr/local/app                                                                0.0s
 => [3/6] RUN useradd app                                                                       0.1s
 => [4/6] COPY ./requirements.txt ./requirements.txt                                            0.0s
 => [5/6] RUN pip install --no-cache-dir --upgrade -r requirements.txt                          3.2s
 => [6/6] COPY ./app ./app                                                                      0.0s
 => exporting to image                                                                          0.1s
 => => exporting layers                                                                         0.1s
 => => writing image sha256:9924dfd9350407b3df01d1a0e1033b1e543523ce7d5d5e2c83a724480ebe8f00    0.0s
```

With the previous output, you could start a container by using the referenced image:

```console
docker run sha256:9924dfd9350407b3df01d1a0e1033b1e543523ce7d5d5e2c83a724480ebe8f00
```

That name certainly isn't memorable, which is where tagging becomes useful.


### Tagging images

Tagging images is the method to provide an image with a memorable name. However, there is a structure to the name of an image. A full image name has the following structure:

```
[HOST[:PORT_NUMBER]/]PATH[:TAG]
```

- `HOST`: The optional registry hostname where the image is located. If no host is specified, Docker's public registry at docker.io is used by default.
- `PORT_NUMBER`: The registry port number if a hostname is provided
- `PATH`: The path of the image, consisting of slash-separated components. For Docker Hub, the format follows [NAMESPACE/]REPOSITORY, where namespace is either a user's or organization's name. If no namespace is specified, library is used.
- `TAG`: A custom, human-readable identifier that's typically used to identify different versions or variants of an image. If no tag is specified, latest is used by default.

Some examples of image names include:

- `nginx` - an image pull would come from the docker.io registry, the library namespace, the nginx image repository, and the latest tag
- `docker/welcome-to-docker` - an image pull would come from the docker.io registry, the docker namespace, the welcome-to-docker image repository, and the latest tag
- `ghcr.io/dockersamples/example-voting-app-vote:pr-311` - will pull from the GitHub Container Registry, the dockersamples namespace, the example-voting-app-vote image repository, and the pr-311 tag

To tag an image during a build, add the `-t` or `--tag` flag:

```console
docker build -t my-username/my-image .
```

If you've already built an image, you can add another tag to the image by using the [docker image tag](https://docs.docker.com/engine/reference/commandline/image_tag/) command:

```console
docker image tag my-username/my-image another-username/another-image:v1
```

### Publishing images

Once you have an image built and tagged, you are ready to push it to a registry. To do so, use the [docker push](https://docs.docker.com/engine/reference/commandline/image_push/) command:

```console
docker push my-username/my-image
```

Within a few seconds, all of the layers for your image will be pushed to the registry.

> **Requiring authentication**
>
> Before you're able to push an image to a repository, you will need to be authenticated.
> To do so, simply use the [docker login](https://docs.docker.com/engine/reference/commandline/login/) command.
{ .information }

## Try it out

In this hands-on, you will build a simple image using a provided `Dockerfile` and push it to Docker Hub.

### Set up

1. Get the sample application.

   If you have Git, you can clone the repository for the sample application. Otherwise, you can download the sample application. Choose one of the following options.

   {{< tabs >}}
   {{< tab name="Clone with git" >}}

   Use the following command in a terminal to clone the sample application repository.

   ```console
   $ git clone https://github.com/docker/getting-started-todo-app
   ```
   {{< /tab >}}
   {{< tab name="Download" >}}

   Download the source and extract it.

   {{< button url="https://github.com/docker/getting-started-todo-app/raw/cd61f824da7a614a8298db503eed6630eeee33a3/app.zip" text="Download the source" >}}

   {{< /tab >}}
   {{< /tabs >}}


2. If you don't have a Docker account yet, create one now. Once you've done that, sign in to Docker Desktop using that account.

3. In [Docker Hub](https://hub.docker.com), create a repository for your new image. Give the new repository a name of `concepts-build-image-demo` and use the defaults for all of the other settings.

### Build an image

Now that you have a repository on Docker Hub, it's time for you to build an image and push it to the repository.

1. Using a terminal in the root of the sample app repository, run the following command. Replace `YOUR_DOCKER_USERNAME` with your Docker Hub username:

    ```console
    $ docker build -t <YOUR_DOCKER_USERNAME>/concepts-build-image-demo .
    ```

    As an example, if your username is `mobywhale`, you would run the command:

    ```console
    $ docker build -t mobywhale/concepts-build-image-demo .
    ```

2. Once the build has completed, you can view the image by using the following command:

    ```console
    $ docker image ls
    ```

    The command will produce output similar to the following:

    ```plaintext
    REPOSITORY                             TAG       IMAGE ID       CREATED          SIZE
    mobywhale/concepts-build-image-demo    latest    746c7e06537f   24 seconds ago   354MB
    ```

3. You can actually view the history (or how the image was created) by using the [docker image history](/reference/cli/docker/image/history/) command:

    ```console
    $ docker image history mobywhale/concepts-build-image-demo
    ```

    You'll then see output similar to the following:

    ```plaintext
    IMAGE          CREATED         CREATED BY                                      SIZE      COMMENT
    f279389d5f01   8 seconds ago   CMD ["node" "./src/index.js"]                   0B        buildkit.dockerfile.v0
    <missing>      8 seconds ago   EXPOSE map[3000/tcp:{}]                         0B        buildkit.dockerfile.v0 
    <missing>      8 seconds ago   WORKDIR /app                                    8.19kB    buildkit.dockerfile.v0
    <missing>      4 days ago      /bin/sh -c #(nop)  CMD ["node"]                 0B
    <missing>      4 days ago      /bin/sh -c #(nop)  ENTRYPOINT ["docker-entry…   0B
    <missing>      4 days ago      /bin/sh -c #(nop) COPY file:4d192565a7220e13…   20.5kB
    <missing>      4 days ago      /bin/sh -c apk add --no-cache --virtual .bui…   7.92MB
    <missing>      4 days ago      /bin/sh -c #(nop)  ENV YARN_VERSION=1.22.19     0B
    <missing>      4 days ago      /bin/sh -c addgroup -g 1000 node     && addu…   126MB
    <missing>      4 days ago      /bin/sh -c #(nop)  ENV NODE_VERSION=20.12.0     0B
    <missing>      2 months ago    /bin/sh -c #(nop)  CMD ["/bin/sh"]              0B
    <missing>      2 months ago    /bin/sh -c #(nop) ADD file:d0764a717d1e9d0af…   8.42MB
    ```

    This output shows the layers of the image, highlighting the layers you added and those that were inherited from your base image.

### Push the image

Now that you have an image built, it's time to push the image to a registry.

1. Push the image using the [docker push](/reference/cli/docker/image/push/) command:

    ```console
    $ docker push <YOUR_DOCKER_USERNAME>/concepts-build-image-demo
    ```

    If you receive a `requested access to the resource is denied`, make sure you are both logged in and that your Docker username is correct in the image tag.

    After a moment, your image should be pushed to Docker Hub! 🎉

## Additional resources

To learn more about building, tagging, and publishing images, visit the following resources:

* [What is a build context?](/build/building/context/#what-is-a-build-context)
* [docker build reference](/engine/reference/commandline/image_build/)
* [docker image tag reference](/engine/reference/commandline/image_tag/)
* [docker push reference](/engine/reference/commandline/image_push/)
* [What is a registry?](/guides/docker-concepts/the-basics/what-is-a-registry/)

Now that you have learned about building and publishing images, it's time to learn how to speed up build process using the Docker build cache.

{{< button text="Using the build cache" url="using-the-build-cache" >}}


