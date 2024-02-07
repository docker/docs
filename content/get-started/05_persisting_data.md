---
title: Persist the DB
keywords: get started, setup, orientation, quickstart, intro, concepts, containers,
  docker desktop
description: Making your DB persistent in your application
---

In case you didn't notice, your todo list is empty every single time
you launch the container. Why is this? In this part, you'll dive into how the container is working.

## The container's filesystem

When a container runs, it uses the various layers from an image for its filesystem.
Each container also gets its own "scratch space" to create/update/remove files. Any
changes won't be seen in another container, even if they're using the same image.

### See this in practice

To see this in action, you're going to start two containers and create a file in each.
What you'll see is that the files created in one container aren't available in another.

> **Note**
>
> If you use Windows and want to use Git Bash to run Docker commands, see [Working with Git Bash](../desktop/troubleshoot/topics.md#working-with-git-bash) for syntax differences.

1. Start an `ubuntu` container that will create a file named `/data.txt` with a random number
   between 1 and 10000.

    ```console
    $ docker run -d ubuntu bash -c "shuf -i 1-10000 -n 1 -o /data.txt && tail -f /dev/null"
    ```

    In case you're curious about the command, you're starting a bash shell and invoking two
    commands (why you have the `&&`). The first portion picks a single random number and writes
    it to `/data.txt`. The second command is simply watching a file to keep the container running.

2. Validate that you can see the output by accessing the terminal in the container. To do so, you can use the CLI or Docker Desktop's graphical interface.

   {{< tabs group="ui" >}}
   {{< tab name="CLI" >}}

   On the command line, use the `docker exec` command to access the container. You need to get the
   container's ID (use `docker ps` to get it). In your Mac or Linux terminal, or in Windows Command Prompt or PowerShell, get the content with the following command.

   ```console
   $ docker exec <container-id> cat /data.txt
   ```
   
   {{< /tab >}}
   {{< tab name="Docker Desktop" >}}
   
   In Docker Desktop, go to **Containers**, hover over the container running the **ubuntu** image, and select the **Show container actions** menu. From the drop-down menu, select **Open in terminal**.

   You will see a terminal that is running a shell in the Ubuntu container. Run the following command to see the content of the `/data.txt` file. Close this terminal afterwards again.

   ```console
   $ cat /data.txt
   ```

   {{< /tab >}}
   {{< /tabs >}}

   You should see a random number.

3. Now, start another `ubuntu` container (the same image) and you'll see you don't have the same file. In your Mac or Linux terminal, or in Windows Command Prompt or PowerShell, get the content with the following command.

    ```console
    $ docker run -it ubuntu ls /
    ```

    In this case the command lists the files in the root directory of the container.
    Look, there's no `data.txt` file there! That's because it was written to the scratch space for
    only the first container.

4. Go ahead and remove the first container using the `docker rm -f <container-id>` command.

## Container volumes

With the previous experiment, you saw that each container starts from the image definition each time it starts. 
While containers can create, update, and delete files, those changes are lost when you remove the container 
and Docker isolates all changes to that container. With volumes, you can change all of this.

[Volumes](../storage/volumes.md) provide the ability to connect specific filesystem paths of 
the container back to the host machine. If you mount a directory in the container, changes in that
directory are also seen on the host machine. If you mount that same directory across container restarts, you'd see
the same files.

There are two main types of volumes. You'll eventually use both, but you'll start with volume mounts.

## Persist the todo data

By default, the todo app stores its data in a SQLite database at
`/etc/todos/todo.db` in the container's filesystem. If you're not familiar with SQLite, no worries! It's simply a relational database that stores all the data in a single file. While this isn't the best for large-scale applications,
it works for small demos. You'll learn how to switch this to a different database engine later.

With the database being a single file, if you can persist that file on the host and make it available to the
next container, it should be able to pick up where the last one left off. By creating a volume and attaching
(often called "mounting") it to the directory where you stored the data, you can persist the data. As your container 
writes to the `todo.db` file, it will persist the data to the host in the volume.

As mentioned, you're going to use a volume mount. Think of a volume mount as an opaque bucket of data. 
Docker fully manages the volume, including the storage location on disk. You only need to remember the
name of the volume.

### Create a volume and start the container

You can create the volume and start the container using the CLI or Docker Desktop's graphical interface.

{{< tabs group="ui" >}}
{{< tab name="CLI" >}}

1. Create a volume by using the `docker volume create` command.

   ```console
   $ docker volume create todo-db
   ```

2. Stop and remove the todo app container once again with `docker rm -f <id>`, as it is still running without using the persistent volume.

3. Start the todo app container, but add the `--mount` option to specify a volume mount. Give the volume a name, and mount
   it to `/etc/todos` in the container, which captures all files created at the path. In your Mac or Linux terminal, or in Windows Command Prompt or PowerShell, run the following command:

   ```console
   $ docker run -dp 127.0.0.1:3000:3000 --mount type=volume,src=todo-db,target=/etc/todos getting-started
   ```

{{< /tab >}}
{{< tab name="Docker Desktop" >}}

To create a volume:

1. Select **Volumes** in Docker Desktop.
2. In **Volumes**, select **Create**.
3. Specify `todo-db` as the volume name, and then select **Create**.

To stop and remove the app container:

1. Select **Containers** in Docker Desktop.
2. Select **Delete** in the **Actions** column for the container.

To start the todo app container with the volume mounted:

1. Select the search box at the top of Docker Desktop.
2. In the search window, select the **Images** tab.
3. In the search box, specify the container name, `getting-started`.

   > **Tip**
   >
   >  Use the search filter to filter images and only show **Local images**.
   { .tip }

4. Select your image and then select **Run**.
5. Select **Optional settings**.
6. In **Host port**, specify the port, for example, `3000`.
7. In **Host path**, specify the name of the volume, `todo-db`.
8. In **Container path**, specify `/etc/todos`.
9. Select **Run**.

{{< /tab >}}
{{< /tabs >}}

### Verify that the data persists

1. Once the container starts up, open the app and add a few items to your todo list.

    ![Items added to todo list](images/get-started/items-added.webp)
    

2. Stop and remove the container for the todo app. Use Docker Desktop or `docker ps` to get the ID and then `docker rm -f <id>` to remove it.

3. Start a new container using the previous steps.

4. Open the app. You should see your items still in your list.

5. Go ahead and remove the container when you're done checking out your list.

You've now learned how to persist data.

## Dive into the volume

A lot of people frequently ask "Where is Docker storing my data when I use a volume?" If you want to know, 
you can use the `docker volume inspect` command.

```console
$ docker volume inspect todo-db
[
    {
        "CreatedAt": "2019-09-26T02:18:36Z",
        "Driver": "local",
        "Labels": {},
        "Mountpoint": "/var/lib/docker/volumes/todo-db/_data",
        "Name": "todo-db",
        "Options": {},
        "Scope": "local"
    }
]
```

The `Mountpoint` is the actual location of the data on the disk. Note that on most machines, you will
need to have root access to access this directory from the host.

## Summary

In this section, you learned how to persist container data.

Related information:

 - [docker CLI reference](/engine/reference/commandline/cli/)
 - [Volumes](../storage/volumes.md)

## Next steps

Next, you'll learn how you can develop your app more efficiently using bind mounts.

{{< button text="Use bind mounts" url="06_bind_mounts.md" >}}
