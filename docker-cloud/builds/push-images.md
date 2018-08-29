---
description: Push images to Docker Cloud
keywords: images, private, registry
redirect_from:
- /docker-cloud/getting-started/intermediate/pushing-images-to-dockercloud/
- /docker-cloud/tutorials/pushing-images-to-dockercloud/
title: Push images to Docker Cloud
notoc: true
---

Docker Cloud uses Docker Hub as its native registry for storing both public and
private repositories. Once you push your images to Docker Hub, they are
available in Docker Cloud.

If you don't have Swarm Mode enable, images pushed to Docker Hub automatically appear for you on the **Services/Wizard** page on Docker Cloud.

> **Note**: You must use Docker Engine 1.6 or later to push to Docker Hub.
Follow the [official installation instructions](/install/index.md){: target="_blank" class="_" } depending on your system.

1. In a terminal window, set the environment variable **DOCKER_ID_USER** as *your username* in Docker Cloud.

    This allows you to copy and paste the commands directly from this tutorial.

    ```
    $ export DOCKER_ID_USER="username"
    ```

    If you don't want to set this environment variable, change the examples in
    this tutorial to replace `DOCKER_ID_USER` with your Docker Cloud username.

2. Log in to Docker Cloud using the `docker login` command.

    ```
    $ docker login
    ```
    This logs you in using your Docker ID, which is shared between both Docker Hub and Docker Cloud.

    If you have never logged in to Docker Hub or Docker Cloud and do not have a Docker ID, running this command prompts you to create a Docker ID.

3. Tag your image using `docker tag`.

    In the example below replace `my_image` with your image's name, and `DOCKER_ID_USER` with your Docker Cloud username if needed.

    ```
    $ docker tag my_image $DOCKER_ID_USER/my_image
    ```

4. Push your image to Docker Hub using `docker push` (making the same replacements as in the previous step).

    ```
    $ docker push $DOCKER_ID_USER/my_image
    ```

5. Check that the image you just pushed appears in Docker Cloud.

    Go to Docker Cloud and navigate to the **Repositories** tab and confirm that your image appears in this list.

>**Note**: If you're a member of any organizations that are using Docker
> Cloud, you might need to switch to the organization account namespace using the
> account menu at the upper right to see other repositories.
