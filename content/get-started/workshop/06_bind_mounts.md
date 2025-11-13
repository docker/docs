---
title: Use bind mounts
weight: 60
linkTitle: "Part 5: Use bind mounts"
keywords: 'get started, setup, orientation, quickstart, intro, concepts, containers, docker desktop'
description: Using bind mounts in our application
aliases:
 - /guides/walkthroughs/access-local-folder/
 - /get-started/06_bind_mounts/
 - /guides/workshop/06_bind_mounts/
---

In [part 4](./05_persisting_data.md), you used a volume mount to persist the
data in your database. A volume mount is a great choice when you need somewhere
persistent to store your application data.

A bind mount is another type of mount, which lets you share a directory from the
host's filesystem into the container. When working on an application, you can
use a bind mount to mount source code into the container. The container sees the
changes you make to the code immediately, as soon as you save a file. This means
that you can run processes in the container that watch for filesystem changes
and respond to them.

In this chapter, you'll see how you can use bind mounts and a tool called
[nodemon](https://npmjs.com/package/nodemon) to watch for file changes, and then restart the application
automatically. There are equivalent tools in most other languages and
frameworks.

## Quick volume type comparisons

The following are examples of a named volume and a bind mount using `--mount`:

- Named volume: `type=volume,src=my-volume,target=/usr/local/data`
- Bind mount: `type=bind,src=/path/to/data,target=/usr/local/data`

The following table outlines the main differences between volume mounts and bind
mounts.

|                                              | Named volumes                                      | Bind mounts                                          |
| -------------------------------------------- | -------------------------------------------------- | ---------------------------------------------------- |
| Host location                                | Docker chooses                                     | You decide                                           |
| Populates new volume with container contents | Yes                                                | No                                                   |
| Supports Volume Drivers                      | Yes                                                | No                                                   |

## Trying out bind mounts

Before looking at how you can use bind mounts for developing your application,
you can run a quick experiment to get a practical understanding of how bind mounts
work.

1. Verify that your `getting-started-app` directory is in a directory defined in
Docker Desktop's file sharing setting. This setting defines which parts of your
filesystem you can share with containers. For details about accessing the setting, see [File sharing](/manuals/desktop/settings-and-maintenance/settings.md#file-sharing).

    > [!NOTE]
    > The **File sharing** tab is only available in Hyper-V mode, because the files are automatically shared in WSL 2 mode and Windows container mode.

2. Open a terminal and change directory to the `getting-started-app`
   directory.

3. Run the following command to start `bash` in an `ubuntu` container with a
   bind mount.

   {{< tabs >}}
   {{< tab name="Mac / Linux" >}}

   ```console
   $ docker run -it --mount type=bind,src="$(pwd)",target=/src ubuntu bash
   ```
   
   {{< /tab >}}
   {{< tab name="Command Prompt" >}}

   ```console
   $ docker run -it --mount "type=bind,src=%cd%,target=/src" ubuntu bash
   ```
   
   {{< /tab >}}
   {{< tab name="Git Bash" >}}

   ```console
   $ docker run -it --mount type=bind,src="/$(pwd)",target=/src ubuntu bash
   ```
   
   {{< /tab >}}
   {{< tab name="PowerShell" >}}

   ```console
   $ docker run -it --mount "type=bind,src=$($pwd),target=/src" ubuntu bash
   ```
   
   {{< /tab >}}
   {{< /tabs >}}
   
   The `--mount type=bind` option tells Docker to create a bind mount, where `src` is the
   current working directory on your host machine (`getting-started-app`), and
   `target` is where that directory should appear inside the container (`/src`).

4. After running the command, Docker starts an interactive `bash` session in the
   root directory of the container's filesystem.

   ```console
   root@ac1237fad8db:/# pwd
   /
   root@ac1237fad8db:/# ls
   bin   dev  home  media  opt   root  sbin  srv  tmp  var
   boot  etc  lib   mnt    proc  run   src   sys  usr
   ```

5. Change directory to the `src` directory.

   This is the directory that you mounted when starting the container. Listing
   the contents of this directory displays the same files as in the
   `getting-started-app` directory on your host machine.

   ```console
   root@ac1237fad8db:/# cd src
   root@ac1237fad8db:/src# ls
   Dockerfile  node_modules  package.json  spec  src  yarn.lock
   ```

6. Create a new file named `myfile.txt`.

   ```console
   root@ac1237fad8db:/src# touch myfile.txt
   root@ac1237fad8db:/src# ls
   Dockerfile  myfile.txt  node_modules  package.json  spec  src  yarn.lock
   ```

7. Open the `getting-started-app` directory on the host and observe that the
   `myfile.txt` file is in the directory.

   ```text
   ├── getting-started-app/
   │ ├── Dockerfile
   │ ├── myfile.txt
   │ ├── node_modules/
   │ ├── package.json
   │ ├── spec/
   │ ├── src/
   │ └── yarn.lock
   ```

8. From the host, delete the `myfile.txt` file.
9. In the container, list the contents of the `app` directory once more. Observe that the file is now gone.

   ```console
   root@ac1237fad8db:/src# ls
   Dockerfile  node_modules  package.json  spec  src  yarn.lock
   ```

10. Stop the interactive container session with `Ctrl` + `D`.

That's all for a brief introduction to bind mounts. This procedure
demonstrated how files are shared between the host and the container, and how
changes are immediately reflected on both sides. Now you can use
bind mounts to develop software.

## Development containers

Using bind mounts is common for local development setups. The advantage is that the development machine doesn’t need to have all of the build tools and environments installed. With a single docker run command, Docker pulls dependencies and tools.

### Run your app in a development container

The following steps describe how to run a development container with a bind
mount that does the following:

- Mount your source code into the container
- Install all dependencies
- Start `nodemon` to watch for filesystem changes

You can use the CLI or Docker Desktop to run your container with a bind mount.

{{< tabs >}}
{{< tab name="Mac / Linux CLI" >}}

1. Make sure you don't have any `getting-started` containers currently running.

2. Run the following command from the `getting-started-app` directory.

   ```console
   $ docker run -dp 127.0.0.1:3000:3000 \
       -w /app --mount type=bind,src="$(pwd)",target=/app \
       node:lts-alpine \
       sh -c "yarn install && yarn run dev"
   ```

   The following is a breakdown of the command:
   - `-dp 127.0.0.1:3000:3000` - same as before. Run in detached (background) mode and
     create a port mapping
   - `-w /app` - sets the "working directory" or the current directory that the
     command will run from
   - `--mount type=bind,src="$(pwd)",target=/app` - bind mount the current
     directory from the host into the `/app` directory in the container
   - `node:lts-alpine` - the image to use. Note that this is the base image for
     your app from the Dockerfile
   - `sh -c "yarn install && yarn run dev"` - the command. You're starting a
     shell using `sh` (alpine doesn't have `bash`) and running `yarn install` to
     install packages and then running `yarn run dev` to start the development
     server. If you look in the `package.json`, you'll see that the `dev` script
     starts `nodemon`.

3. You can watch the logs using `docker logs <container-id>`. You'll know you're
   ready to go when you see this:

   ```console
   $ docker logs -f <container-id>
   nodemon -L src/index.js
   [nodemon] 2.0.20
   [nodemon] to restart at any time, enter `rs`
   [nodemon] watching path(s): *.*
   [nodemon] watching extensions: js,mjs,json
   [nodemon] starting `node src/index.js`
   Using sqlite database at /etc/todos/todo.db
   Listening on port 3000
   ```

   When you're done watching the logs, exit out by hitting `Ctrl`+`C`.

{{< /tab >}}
{{< tab name="PowerShell CLI" >}}

1. Make sure you don't have any `getting-started` containers currently running.

2. Run the following command from the `getting-started-app` directory.

   ```powershell
   $ docker run -dp 127.0.0.1:3000:3000 `
       -w /app --mount "type=bind,src=$pwd,target=/app" `
       node:lts-alpine `
       sh -c "yarn install && yarn run dev"
   ```

   The following is a breakdown of the command:
   - `-dp 127.0.0.1:3000:3000` - same as before. Run in detached (background) mode and
     create a port mapping
   - `-w /app` - sets the "working directory" or the current directory that the
     command will run from
   - `--mount "type=bind,src=$pwd,target=/app"` - bind mount the current
     directory from the host into the `/app` directory in the container
   - `node:lts-alpine` - the image to use. Note that this is the base image for
     your app from the Dockerfile
   - `sh -c "yarn install && yarn run dev"` - the command. You're starting a
     shell using `sh` (alpine doesn't have `bash`) and running `yarn install` to
     install packages and then running `yarn run dev` to start the development
     server. If you look in the `package.json`, you'll see that the `dev` script
     starts `nodemon`.

3. You can watch the logs using `docker logs <container-id>`. You'll know you're
   ready to go when you see this:

   ```console
   $ docker logs -f <container-id>
   nodemon -L src/index.js
   [nodemon] 2.0.20
   [nodemon] to restart at any time, enter `rs`
   [nodemon] watching path(s): *.*
   [nodemon] watching extensions: js,mjs,json
   [nodemon] starting `node src/index.js`
   Using sqlite database at /etc/todos/todo.db
   Listening on port 3000
   ```

   When you're done watching the logs, exit out by hitting `Ctrl`+`C`.

{{< /tab >}}
{{< tab name="Command Prompt CLI" >}}

1. Make sure you don't have any `getting-started` containers currently running.

2. Run the following command from the `getting-started-app` directory.

   ```console
   $ docker run -dp 127.0.0.1:3000:3000 ^
       -w /app --mount "type=bind,src=%cd%,target=/app" ^
       node:lts-alpine ^
       sh -c "yarn install && yarn run dev"
   ```

   The following is a breakdown of the command:
   - `-dp 127.0.0.1:3000:3000` - same as before. Run in detached (background) mode and
     create a port mapping
   - `-w /app` - sets the "working directory" or the current directory that the
     command will run from
   - `--mount "type=bind,src=%cd%,target=/app"` - bind mount the current
     directory from the host into the `/app` directory in the container
   - `node:lts-alpine` - the image to use. Note that this is the base image for
     your app from the Dockerfile
   - `sh -c "yarn install && yarn run dev"` - the command. You're starting a
     shell using `sh` (alpine doesn't have `bash`) and running `yarn install` to
     install packages and then running `yarn run dev` to start the development
     server. If you look in the `package.json`, you'll see that the `dev` script
     starts `nodemon`.

3. You can watch the logs using `docker logs <container-id>`. You'll know you're
   ready to go when you see this:

   ```console
   $ docker logs -f <container-id>
   nodemon -L src/index.js
   [nodemon] 2.0.20
   [nodemon] to restart at any time, enter `rs`
   [nodemon] watching path(s): *.*
   [nodemon] watching extensions: js,mjs,json
   [nodemon] starting `node src/index.js`
   Using sqlite database at /etc/todos/todo.db
   Listening on port 3000
   ```

   When you're done watching the logs, exit out by hitting `Ctrl`+`C`.

{{< /tab >}}
{{< tab name="Git Bash CLI" >}}

1. Make sure you don't have any `getting-started` containers currently running.

2. Run the following command from the `getting-started-app` directory.

   ```console
   $ docker run -dp 127.0.0.1:3000:3000 \
       -w //app --mount type=bind,src="/$(pwd)",target=/app \
       node:lts-alpine \
       sh -c "yarn install && yarn run dev"
   ```

   The following is a breakdown of the command:
   - `-dp 127.0.0.1:3000:3000` - same as before. Run in detached (background) mode and
     create a port mapping
   - `-w //app` - sets the "working directory" or the current directory that the
     command will run from
   - `--mount type=bind,src="/$(pwd)",target=/app` - bind mount the current
     directory from the host into the `/app` directory in the container
   - `node:lts-alpine` - the image to use. Note that this is the base image for
     your app from the Dockerfile
   - `sh -c "yarn install && yarn run dev"` - the command. You're starting a
     shell using `sh` (alpine doesn't have `bash`) and running `yarn install` to
     install packages and then running `yarn run dev` to start the development
     server. If you look in the `package.json`, you'll see that the `dev` script
     starts `nodemon`.

3. You can watch the logs using `docker logs <container-id>`. You'll know you're
   ready to go when you see this:

   ```console
   $ docker logs -f <container-id>
   nodemon -L src/index.js
   [nodemon] 2.0.20
   [nodemon] to restart at any time, enter `rs`
   [nodemon] watching path(s): *.*
   [nodemon] watching extensions: js,mjs,json
   [nodemon] starting `node src/index.js`
   Using sqlite database at /etc/todos/todo.db
   Listening on port 3000
   ```

   When you're done watching the logs, exit out by hitting `Ctrl`+`C`.

{{< /tab >}}
{{< tab name="Docker Desktop" >}}

Make sure you don't have any `getting-started` containers currently running.

Run the image with a bind mount.

1. Select the search box at the top of Docker Desktop.
2. In the search window, select the **Images** tab.
3. In the search box, specify the container name, `getting-started`.

   > [!TIP]
   >
   >  Use the search filter to filter images and only show **Local images**.

4. Select your image and then select **Run**.
5. Select **Optional settings**.
6. In **Host path**, specify the path to the `getting-started-app` directory on your host machine.
7. In **Container path**, specify `/app`.
8. Select **Run**.

You can watch the container logs using Docker Desktop.

1. Select **Containers** in Docker Desktop.
2. Select your container name.

You'll know you're ready to go when you see this:

```console
nodemon -L src/index.js
[nodemon] 2.0.20
[nodemon] to restart at any time, enter `rs`
[nodemon] watching path(s): *.*
[nodemon] watching extensions: js,mjs,json
[nodemon] starting `node src/index.js`
Using sqlite database at /etc/todos/todo.db
Listening on port 3000
```

{{< /tab >}}
{{< /tabs >}}

### Develop your app with the development container

Update your app on your host machine and see the changes reflected in the container.

1. In the `src/static/js/app.js` file, on line
   109, change the "Add Item" button to simply say "Add":

   ```diff
   - {submitting ? 'Adding...' : 'Add Item'}
   + {submitting ? 'Adding...' : 'Add'}
   ```

   Save the file.

2. Refresh the page in your web browser, and you should see the change reflected
   almost immediately because of the bind mount. Nodemon detects the change and
   restarts the server. It might take a few seconds for the Node server to
   restart. If you get an error, try refreshing after a few seconds.

   ![Screenshot of updated label for Add button](images/updated-add-button.webp)

3. Feel free to make any other changes you'd like to make. Each time you make a
   change and save a file, the change is reflected in the container because of
   the bind mount. When Nodemon detects a change, it restarts the app inside the
   container automatically. When you're done, stop the container and build your
   new image using:

   ```console
   $ docker build -t getting-started .
   ```

## Summary

At this point, you can persist your database and see changes in your app as you develop without rebuilding the image.

In addition to volume mounts and bind mounts, Docker also supports other mount
types and storage drivers for handling more complex and specialized use cases.

Related information:

 - [docker CLI reference](/reference/cli/docker/)
 - [Manage data in Docker](https://docs.docker.com/storage/)

## Next steps

In order to prepare your app for production, you need to migrate your database
from working in SQLite to something that can scale a little better. For
simplicity, you'll keep using a relational database and switch your application
to use MySQL. But, how should you run MySQL? How do you allow the containers to
talk to each other? You'll learn about that in the next section.

{{< button text="Multi container apps" url="07_multi_container.md" >}}
