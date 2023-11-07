---
title: Persist data between containers
keywords: get started, quick start, intro, concepts
description: Learn how to persist data between containers
---

This walkthrough shows you how to persist data between containers. To better understand some concepts in this walkthrough, complete the [Run multi-container applications](./multi-container-apps.md) walkthrough first.

Docker isolates all content, code, and data in a container from your local filesystem. When you delete a container, Docker deletes all the content within that container.

![Data isolation diagram](images/getting-started-isolation.png?w=400)

Sometimes, you may want to persist the data that a container generates. To do this, you can use volumes.

Before you start, [get Docker Desktop](../../get-docker.md).

## Step 1: Get the sample application

If you have git, you can clone the repository for the sample application. Otherwise, you can download the sample application. Choose one of the following options.

{{< tabs >}}
{{< tab name="Clone with git" >}}

Use the following command in a terminal to clone the sample application repository.

```console
$ git clone https://github.com/docker/multi-container-app
```

{{< /tab >}}
{{< tab name="Download" >}}

Download the source and extract it.

{{< button url="https://github.com/docker/multi-container-app/archive/refs/heads/main.zip" text="Download the source" >}}

{{< /tab >}}
{{< /tabs >}}

## Step 2: Add a volume to persist data

To persist data after you delete a container, use a volume. A volume is a location in your local filesystem, automatically managed by Docker Desktop.

![Volume diagram](images/getting-started-volume.png?w=400)

To add a volume to this project, open the `compose.yaml` file in a code or text editor, and then uncomment the following lines.

```yaml
todo-database:
    # ...
    volumes:
      - database:/data/db

# ...
volumes:
  database:
```

The `volumes` element that is nested under `todo-database` tells Compose to mount the volume named `database` to `/data/db` in the container for the todo-database service.

The top-level `volumes` element defines and configures a volume named `database` that can be used by any of the services in the Compose file.

## Step 3: Run the application

To run the multi-container application, open a terminal and run the following commands. Replace `/path/to/multi-container-app/` with the path to your applications directory

{{< include "open-terminal.md" >}}

```console
$ cd /path/to/multi-container-app/
```
```console
$ docker compose up -d
```

## Step 4: View the frontend and add todos

In the **Containers** tab of Docker Desktop, you should now have an application stack with two containers running (the todo-app, and todo-database).

To view the frontend and add todos, do the following:

1. In Docker Desktop, expand the application stack in **Containers**.
2. Select the link to port **3000** in the **Port(s)** column or open [https://localhost:3000](https://localhost:3000)⁠.
3. Add some todo tasks in the frontend.

## Step 5: Delete the application stack and run new containers

Now, no matter how often you delete and recreate the containers, Docker Desktop persists your data and it's accessible to any container on your system by mounting the `database` volume. Docker Desktop looks for the `database` volume and creates it if it doesn't exist.

To delete the application stack, do the following:

1. Open the **Containers** tab of Docker Desktop
2. Select the Delete icon next to your application stack.

![Deleting the application stack](images/getting-started-delete-stack.png?w=300&border=true)

After you delete the application stack, follow the steps from [Step 3: Run the
application](#step-3-run-the-application) to run the application again. Note
that when you delete the containers and run them again, Docker Desktop persists any todos that you created.

## Summary

In this walkthrough, you persisted data between containers using a volume. You can use this to persist and share data among isolated and ephemeral containers.

Related information:

- Deep dive into [volumes](../../storage/volumes.md)
- Learn about using volumes in Compose in the [Compose file reference](../../compose/compose-file/_index.md)
- Explore using volumes via the CLI in the [docker volume CLI reference](../../engine/reference/commandline/volume_create.md) and [Docker run reference](/reference/run/)

## Next steps

Continue to the next walkthrough to learn how you can access a local directory from a container.

{{< button url="./access-local-folder.md" text="Access a local folder" >}}