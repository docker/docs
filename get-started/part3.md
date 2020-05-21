---
title: "Share images on Docker Hub"
keywords: docker hub, push, images
description: Learn how to share images on Docker Hub.
redirect_from:
- /get-started/part5/
---

{% include_relative nav.html selected="3" %}

## Prerequisites

Work through the steps to build an image and run it as a containerized application in [Part 2](part2.md).

## Introduction

At this point, you've built a containerized application described in [Part 2](part2.md) on your local development machine.

The final step in developing a containerized application is to share your images on a registry like [Docker Hub](https://hub.docker.com/){: target="_blank" class="_”}, so they can be easily downloaded and run on any destination machine.

## Set up your Docker Hub account

If you don't have a Docker ID, follow these steps to create one. A Docker ID allows you to share images on Docker Hub.

1.  Visit the [Docker Hub sign up](https://hub.docker.com/signup){: target="_blank" class="_”} page.

2.  Fill out the form and submit to create your Docker ID.

3.  Verify your email address to complete the registration process.

4.  Click on the Docker icon in your toolbar or system tray, and click **Sign in / Create Docker ID**.

5.  Fill in your new Docker ID and password. After you have successfully authenticated, your Docker ID appears in the Docker Desktop menu in place of the 'Sign in' option you just used.

    You can also sign into Docker Hub from the command line by typing `docker login`.

## Create a Docker Hub repository and push your image

>
> Before creating a repository, ensure you’ve set up your Docker Hub account and have connected it to your Docker Desktop.

Now let's create your first repository, and push your bulletin board image to Docker Hub.

1.  Click on the Docker icon in your menu bar, and navigate to **Repositories > Create**. You'll be redirected to the **Create Repository** page on Docker Hub.

2.  Type the repository name as `bulletinboard` and click **Create** at the bottom of the page. Do not fill any other details for now.

    ![make a repo](images/newrepo.png){:width="100%"}

3.  You are now ready to share your image on Docker Hub, however, there's one thing you must do first: images must be *namespaced correctly* to share on Docker Hub. Specifically, you must name images like `<Your Docker ID>/<Repository Name>:<tag>`.

    Make sure you’re in the `node-bulletin-board/bulletin-board-app` directory in a terminal or PowerShell then and run:

    ```shell
    docker tag bulletinboard:1.0 <Your Docker ID>/bulletinboard:1.0
    ```

4.  Finally, push your image to Docker Hub:

    ```shell
    docker push <Your Docker ID>/bulletinboard:1.0
    ```

    Visit your repository in [Docker Hub](https://hub.docker.com/repositories){: target="_blank" class="_”}, and you'll see your new image there. Remember, Docker Hub repositories are public by default.

    > **Having trouble pushing?** Remember, you must be signed into Docker Hub through Docker Desktop or the command line, and you must also name your images correctly, as per the above steps. If the push seemed to work, but you don't see it in Docker Hub, refresh your browser after a couple of minutes and check again.

## Conclusion

Now that your image is available on Docker Hub, you'll be able to run it anywhere. If you try to use it on a new machine that doesn't have it yet, Docker will automatically try and download it from Docker Hub. By moving images around in this way, you no longer need to install any dependencies except Docker on the machines you want to run your software on. The dependencies of containerized applications are completely encapsulated and isolated within your images, which you can share using Docker Hub as described above.

Another thing to keep in mind: at the moment, you've only pushed your image to Docker Hub; what about your Dockerfile? A crucial best practice is to keep these in version control, perhaps alongside your source code for your application. You can add a link or note in your Docker Hub repository description indicating where these files can be found, preserving the record not only of how your image was built, but how it's meant to be run as a full application.

## Where to go next

We recommend that you take a look at the topics in [Develop with Docker](../develop/index.md) to learn how to develop your own applications using Docker.
