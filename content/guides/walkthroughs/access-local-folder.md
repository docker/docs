---
title: Access a local folder from a container
keywords: get started, quick start, intro, concepts
description: Learn how to access a local folder from a container
---

This walkthrough shows you how to access a local folder from a container. To better understand some concepts in this walkthrough, complete the [Run multi-container applications](./multi-container-apps.md) walkthrough first.

Docker isolates all content, code, and data in a container from your local filesystem. By default, containers can't access directories in your local filesystem.

![Data isolation diagram](images/getting-started-isolation.webp?w=400)

Sometimes, you may want to access a directory from your local filesystem. To do this, you can use bind mounts.

{{< include "guides-get-docker.md" >}}


## Step 1: Get the sample application

If you have git, you can clone the repository for the sample application. Otherwise, you can download the sample application. Choose one of the following options.

{{< tabs >}}
{{< tab name="Clone with git" >}}

Use the following command in a terminal to clone the sample application repository.

```console
$ git clone https://github.com/docker/bindmount-apps
```

{{< /tab >}}
{{< tab name="Download" >}}

Download the source and extract it.

{{< button url="https://github.com/docker/bindmount-apps/archive/refs/heads/main.zip" text="Download the source" >}}

{{< /tab >}}
{{< /tabs >}}

## Step 2: Add a bind mount using Compose

Add a bind mount to access data on your system from a container. A bind mount lets you share a directory from your host's filesystem into the container.

![Bind mount diagram](images/getting-started-bindmount.webp?w=400)

To add a bind mount to this project, open the `compose.yaml` file in a code or text editor, and then uncomment the following lines.

```yaml
todo-app:
    # ...
    volumes:
      - ./app:/usr/src/app
      - /usr/src/app/node_modules

```

The `volumes` element tells Compose to mount the local folder `./app` to `/usr/src/app` in the container for the `todo-app` service. This particular bind mount overwrites the static contents of the `/usr/src/app` directory in the container and creates what is known as a development container. The second instruction, `/usr/src/app/node_modules`, prevents the bind mount from overwriting the container's `node_modules` directory to preserve the packages installed in the container.

## Step 3: Run the application

In a terminal, run the follow commands to bring up your application. Replace `/path/to/bindmount-apps/` with the path to your application's directory.

{{< include "open-terminal.md" >}}

```console
$ cd /path/to/bindmount-apps/
```
```console
$ docker compose up -d
```

## Step 4: Develop the application

Now, you can take advantage of the container’s environment while you develop the application on your local system. Any changes you make to the application on your local system are reflected in the container. In your local directory, open `app/views/todos.ejs` in an code or text editor, update the `Enter your task` string, and save the file. Visit or refresh [localhost:3001](https://localhost:3001)⁠ to view the changes.

## Summary

In this walkthrough, you added a bind mount to access a local folder from a container. You can use this to develop faster without having to rebuild your container when updating your code.

Related information:

- Deep dive into [bind mounts](../../storage/bind-mounts.md)
- Learn about using bind mounts in Compose in the [Compose file reference](../../compose/compose-file/_index.md)
- Explore using bind mounts via the CLI in the [Docker run reference](/engine/reference/run/#volume-shared-filesystems)

## Next steps

Continue to the next walkthrough to learn how you can containerize your own application.

{{< button url="./containerize-your-app.md" text="Containerize your app" >}}