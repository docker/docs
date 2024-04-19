---
title: Understanding the image layers
keywords: concepts, build, images, container, docker desktop
description: This concept page will teach you about the layers of container image.
---

{{< youtube-embed wJwqtAkmtQA >}}

## Explanation

As you learned in [What is an image?](../the-basics/what-is-an-image/), container images are composed of layers. And each of these layers, once created, are immutable. But, what does that actually mean? And how are those layers then merged back together to create the filesystem a container will use?

### Image layers

Each layer in an image contains a set of filesystem changes - additions, deletions, or modifications. Let’s look at a theoretical image:

1. The first layer adds basic commands and a package manager, such as apt.
2. The second layer installs a Python runtime and pip for dependency management.
3. The third layer copies in an application’s specific requirements.txt file.
4. The fourth layer installs that application’s specific dependencies.
5. The fifth layer copies in the actual source code of the application.

This example might look like:

![screenshot of the flowchart showing the concept of the image layers](images/container_image_layers.webp?border=true)

Why is this beneficial? With layering, you can reuse the layers from one image for another image. For example, imagine you wanted to create another Python application. Due to layering, you can leverage the same Python base. This will make builds faster and reduce the amount of storage and bandwidth required to distribute the images. The image layering might look similar to the following:

![screenshot of the flowchart showing the benefits of the image layering](images/container_image_layer_reuse.webp?border=true)

With this layering, you can easily extend the images of others! As a developer, it allows you to focus on your unique application needs.

### Merging the layers

The magic for merging the layers lies in the use of a union filesystem. While this will get technical, here’s how it works:

1. After each layer is downloaded, it is extracted into its own directory on the host filesystem. 
2. A union filesystem is created. This stacks each of the layers' directories on top of each other into a new directory, giving a new and unified view. 
3. When the container is started, its root directory is set to the location of this merged directory (using `chroot`).

When the union filesystem is created, in addition to the image layers, a directory is created specifically for the running container. This allows the container to make filesystem changes while allowing the original image layers to remain untouched. This enables you to run multiple containers from the same underlying image.

## Try it out

In this hands-on, you will create new image layers manually using the [docker commit](https://docs.docker.com/reference/cli/docker/container/commit/) command. Note that you’ll rarely create images this way, as you’ll normally [use a Dockerfile](./writing-a-dockerfile.md). But, it makes it easier to understand how it’s all working.

### Create a base image

In this first step, you will create your own base image that you will then use for the following steps.

1. [Download and install](https://www.docker.com/products/docker-desktop/) Docker Desktop.


2. In a terminal, run the following command to start a new container:

    ```console
    $ docker run --name=base-container -ti ubuntu
    ```

    Once the image has been downloaded and the container has started, you should see a new shell prompt. This is running inside your container. It will look similar to the following (the container ID will vary):

    ```console
    root@d8c5ca119fcd:/#
    ```

3. Inside the container, run the following command to install Node.js:

    ```console
    $ apt update && apt install -y nodejs
    ```

    When this command runs, it downloads and installs Node inside the container. In the context of the union filesystem, these filesystem changes occur within the directory unique to this container. 

4. Validate if Node is installed by running the following command:

    ```console
    $ node -e 'console.log("Hello world!")'
    ```

    You should then see a “Hello world!” appear in the console.

5. Now that you have Node installed, you’re ready to save the changes you’ve made as a new image layer, from which you can start new containers or build new images. To do so, you will use the [docker commit](https://docs.docker.com/reference/cli/docker/container/commit/) command. Run the following command in a new terminal:

    ```console
    $ docker commit -m "Add node" base-container node-base
    ```

    Congratulations! You have now created a brand new image! 🎉

6. View the layers of your image using the `docker image history` command:

    ```console
    $ docker image history node-base
    ```

    You will see output similar to the following:

    ```console
    IMAGE          CREATED          CREATED BY                                      SIZE      COMMENT
    d5c1fca2cdc4   10 seconds ago   /bin/bash                                       126MB     Add node
    2b7cc08dcdbb   5 weeks ago      /bin/sh -c #(nop)  CMD ["/bin/bash"]            0B
    <missing>      5 weeks ago      /bin/sh -c #(nop) ADD file:07cdbabf782942af0…   69.2MB
    <missing>      5 weeks ago      /bin/sh -c #(nop)  LABEL org.opencontainers.…   0B
    <missing>      5 weeks ago      /bin/sh -c #(nop)  LABEL org.opencontainers.…   0B
    <missing>      5 weeks ago      /bin/sh -c #(nop)  ARG LAUNCHPAD_BUILD_ARCH     0B
    <missing>      5 weeks ago      /bin/sh -c #(nop)  ARG RELEASE                  0B
    ```

    Note the “Add node” comment on the top line. This layer contains the Node.js install you just made.

7. To prove your image has Node installed, you can start a new container using this new image:

    ```console
    $ docker run node-base node -e "console.log('Hello again')"
    ```

    With that, you should get a “Hello again” output in the terminal, showing Node was installed and working!

8. Now that you’re done creating your base image, you can remove that container:

    ```console
    $ docker rm -f base-container
    ```

> **Base image definition**
>
> A base image is a foundation for building other images. It's possible to use any images as a base image. However, some images are intentionally created as building blocks, providing a foundation or starting point for an application.
>
> In this example, you probably won’t deploy this `node-base` image, as it doesn’t actually do anything yet. But, it’s a base you can use for other builds.
{ .tip }


### Build an app image

Now that you have a base image, you can extend that image to build additional images.

1. Start a new container using the newly created node-base image:

    ```console
    $ docker run --name=app-container -ti node-base
    ```

2. Inside of this container, run the following command to create a Node program:

    ```console
    $ echo 'console.log("Hello from an app")' > app.js
    ```

    To run this Node program, you can use the following command and see the message printed on the screen:

    ```console
    $ node app.js
    ```

3. In another terminal, run the following command to save this container’s changes as a new image:

    ```console
    $ docker commit -c "CMD node app.js" -m "Add app" app-container sample-app
    ```

    This command not only creates a new image named `sample-app`, but also adds additional configuration to the image to set the default command when starting a container. In this case, you are setting it to automatically run `node app.js`.

4. In a terminal outside of the container, run the following command to view the updated layers:

    ```console
    $ docker image history sample-app
    ```

    You’ll then see output that looks like the following. Note the top layer comment has “Add app” and the next layer has “Add node”:

    ```console
    IMAGE          CREATED              CREATED BY                                      SIZE      COMMENT
    c1502e2ec875   About a minute ago   /bin/bash                                       33B       Add app
    5310da79c50a   4 minutes ago        /bin/bash                                       126MB     Add node
    2b7cc08dcdbb   5 weeks ago          /bin/sh -c #(nop)  CMD ["/bin/bash"]            0B
    <missing>      5 weeks ago          /bin/sh -c #(nop) ADD file:07cdbabf782942af0…   69.2MB
    <missing>      5 weeks ago          /bin/sh -c #(nop)  LABEL org.opencontainers.…   0B
    <missing>      5 weeks ago          /bin/sh -c #(nop)  LABEL org.opencontainers.…   0B
    <missing>      5 weeks ago          /bin/sh -c #(nop)  ARG LAUNCHPAD_BUILD_ARCH     0B
    <missing>      5 weeks ago          /bin/sh -c #(nop)  ARG RELEASE                  0B
    ```

5. Finally, start a new container using the brand new image. Since you specified the default command, you can use the following command:

    ```console
    $ docker run sample-app
    ```

    You should see your greeting appear in the terminal, coming from your Node program.

6. Now that you’re done with your containers, you can remove them using the following command:

    ```console
    $ docker rm -f app-container
    ```

## Additional resources

If you’d like to dive deeper into the things you learned, check out the following resources:

* [docker image history CLI reference](/reference/cli/docker/image/history/)
* [docker image layers](/build/guide/layers/)
* [docker commit](/reference/cli/docker/container/commit/)


## Next steps

As hinted earlier, most image builds don’t use `docker commit`. Instead, you’ll use a Dockerfile which automates these steps for you.

{{< button text="Writing a Dockerfile" url="writing-a-dockerfile" >}}
