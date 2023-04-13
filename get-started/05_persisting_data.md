---
title: "Persist the DB"
keywords: get started, setup, orientation, quickstart, intro, concepts, containers, docker desktop
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

1. Start an `ubuntu` container that will create a file named `/data.txt` with a random number
   between 1 and 10000.

    ```console
    $ docker run -d ubuntu bash -c "shuf -i 1-10000 -n 1 -o /data.txt && tail -f /dev/null"
    ```

    In case you're curious about the command, you're starting a bash shell and invoking two
    commands (why you have the `&&`). The first portion picks a single random number and writes
    it to `/data.txt`. The second command is simply watching a file to keep the container running.

2. Validate that you can see the output by accessing the terminal in the container. To do so, you can use the CLI or Docker Desktop's graphical interface.

   <ul class="nav nav-tabs">
     <li class="active"><a data-toggle="tab" data-target="#cli">CLI</a></li>
     <li><a data-toggle="tab" data-target="#gui">Docker Desktop</a></li>
   </ul>
   <div class="tab-content">
   <div id="cli" class="tab-pane fade in active" markdown="1">

    On the command line, use the `docker exec` command to access the container. You need to get the
   container's ID (use `docker ps` to get it). In your Mac or Linux terminal, or in Windows Command Prompt or PowerShell, get the content with the following command.

    ```console
    $ docker exec <container-id> cat /data.txt
    ```

   <hr>
   </div>
   <div id="gui" class="tab-pane fade" markdown="1">

    In Docker Desktop, go to **Containers**, hover over the container running the **ubuntu** image, and select the **Show container actions** menu. From the dropdown menu, select **Open in terminal**.

    You will see a terminal that is running a shell in the Ubuntu container. Run the following command to see the content of the `/data.txt` file. Close this terminal afterwards again.

    ```console
    $ cat /data.txt
    ```

   <hr>
   </div>
   </div>

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

1. Create a volume by using the `docker volume create` command.

    ```console
    $ docker volume create todo-db
    ```

2. Stop and remove the todo app container once again in the Dashboard (or with `docker rm -f <id>`), as it is still running without using the persistent volume.

3. Start the todo app container, but add the `--mount` option to specify a volume mount. Give the volume a name, and mount
   it to `/etc/todos` in the container, which captures all files created at the path.

    ```console
    $ docker run -dp 3000:3000 --mount type=volume,src=todo-db,target=/etc/todos getting-started
    ```

4. Once the container starts up, open the app and add a few items to your todo list.

    ![Items added to todo list](images/items-added.png){: style="width: 55%; " }
    {: .text-center }

5. Stop and remove the container for the todo app. Use the Dashboard or `docker ps` to get the ID and then `docker rm -f <id>` to remove it.

6. Start a new container using the same command from above.

7. Open the app. You should see your items still in your list.

8. Go ahead and remove the container when you're done checking out your list.

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
need to have root access to access this directory from the host. But, that's where it is.

> **Accessing volume data directly on Docker Desktop**
> 
> While running in Docker Desktop, the Docker commands are actually running inside a small VM on your machine.
> If you wanted to look at the actual contents of the mount point directory, you would need to look inside of
> that VM.

## Next steps

At this point, you have a functioning application that can survive restarts.

However, you saw earlier that rebuilding images for every change takes quite a bit of time. There's got to be a better
way to make changes, right? With bind mounts, there is a better way.

[Use bind mounts](06_bind_mounts.md){: .button  .primary-btn}
