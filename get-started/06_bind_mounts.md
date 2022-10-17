---
title: "Use bind mounts"
keywords: get started, setup, orientation, quickstart, intro, concepts, containers, docker desktop
description: Using bind mounts in our application
---

In the previous chapter, we talked about and used a **named volume** to persist the data in our database.
Named volumes are great if we simply want to store data, as we don't have to worry about _where_ the data
is stored.

With **bind mounts**, we control the exact mountpoint on the host. We can use this to persist data, but it's often
used to provide additional data into containers. When working on an application, we can use a bind mount to
mount our source code into the container to let it see code changes, respond, and let us see the changes right
away.

For Node-based applications, [nodemon](https://npmjs.com/package/nodemon){:target="_blank" rel="noopener" class="_"} is a great tool to watch for file
changes and then restart the application. There are equivalent tools in most other languages and frameworks.

## Quick volume type comparisons

Bind mounts and named volumes are the two main types of volumes that come with the Docker engine. However, additional
volume drivers are available to support other use cases ([SFTP](https://github.com/vieux/docker-volume-sshfs){:target="_blank" rel="noopener" class="_"}, [Ceph](https://ceph.com/geen-categorie/getting-started-with-the-docker-rbd-volume-plugin/){:target="_blank" rel="noopener" class="_"}, [NetApp](https://netappdvp.readthedocs.io/en/stable/){:target="_blank" rel="noopener" class="_"}, [S3](https://github.com/elementar/docker-s3-volume){:target="_blank" rel="noopener" class="_"}, and more).

|   | Named Volumes | Bind Mounts |
| - | ------------- | ----------- |
| Host Location | Docker chooses | You control |
| Mount Example (using `-v`) | my-volume:/usr/local/data | /path/to/data:/usr/local/data |
| Populates new volume with container contents | Yes | No |
| Supports Volume Drivers | Yes | No |

## Start a dev-mode container

To run our container to support a development workflow, we will do the following:

- Mount our source code into the container
- Install all dependencies, including the "dev" dependencies
- Start nodemon to watch for filesystem changes

So, let's do it!

1. Make sure you don't have any previous `getting-started` containers running.

2. Run the following command from the app directory. We'll explain what's going on afterwards.

    If you are using an x86-64 Mac or Linux device, then use the following command.

    ```console
    $ docker run -dp 3000:3000 \
        -w /app -v "$(pwd):/app" \
        node:12-alpine \
        sh -c "yarn install && yarn run dev"
    ```

    If you are using Windows, then use the following command in PowerShell.

    ```powershell
    PS> docker run -dp 3000:3000 `
        -w /app -v "$(pwd):/app" `
        node:12-alpine `
        sh -c "yarn install && yarn run dev"
    ```

    If you are using an Apple silicon Mac or another ARM64 device, then use the following command.

    ```console
    $ docker run -dp 3000:3000 \
        -w /app -v "$(pwd):/app" \
        node:12-alpine \
        sh -c "apk add --no-cache python2 g++ make && yarn install && yarn run dev"
    ```

    - `-dp 3000:3000` - same as before. Run in detached (background) mode and create a port mapping
    - `-w /app` - sets the "working directory" or the current directory that the command will run from
    - `-v "$(pwd):/app"` - bind mount the current directory from the host in the container into the `/app` directory
    - `node:12-alpine` - the image to use. Note that this is the base image for our app from the Dockerfile
    - `sh -c "yarn install && yarn run dev"` - the command. We're starting a shell using `sh` (alpine doesn't have `bash`) and
      running `yarn install` to install _all_ dependencies and then running `yarn run dev`. If we look in the `package.json`,
      we'll see that the `dev` script is starting `nodemon`.

3. You can watch the logs using `docker logs`. You'll know you're ready to go when you see this:

    ```console
    $ docker logs -f <container-id>
    nodemon src/index.js
    [nodemon] 1.19.2
    [nodemon] to restart at any time, enter `rs`
    [nodemon] watching dir(s): *.*
    [nodemon] starting `node src/index.js`
    Using sqlite database at /etc/todos/todo.db
    Listening on port 3000
    ```

    When you're done watching the logs, exit out by hitting `Ctrl`+`C`.

4. Now, let's make a change to the app. In the `src/static/js/app.js` file, let's change the "Add Item" button to simply say
   "Add". This change will be on line 109:

    ```diff
    -                         {submitting ? 'Adding...' : 'Add Item'}
    +                         {submitting ? 'Adding...' : 'Add'}
    ```

5. Simply refresh the page (or open it) and you should see the change reflected in the browser almost immediately. It might
   take a few seconds for the Node server to restart, so if you get an error, just try refreshing after a few seconds.

    ![Screenshot of updated label for Add button](images/updated-add-button.png){: style="width:75%;"}
    {: .text-center }

6. Feel free to make any other changes you'd like to make. When you're done, stop the container and build your new image
   using:

    ```console
    $ docker build -t getting-started .
    ```

Using bind mounts is _very_ common for local development setups. The advantage is that the dev machine doesn't need to have
all of the build tools and environments installed. With a single `docker run` command, the dev environment is pulled and ready
to go. We'll talk about Docker Compose in a future step, as this will help simplify our commands (we're already getting a lot
of flags).

## Next steps

At this point, you can persist your database and respond rapidly to the needs and demands of your investors and founders. Hooray!
But, guess what? You received great news! Your project has been selected for future development!

In order to prepare for production, you need to migrate your database from working in SQLite to something that can scale a
little better. For simplicity, you'll keep with a relational database and switch your application to use MySQL. But, how 
should you run MySQL? How do you allow the containers to talk to each other? You'll learn about that next!

[Multi container apps](07_multi_container.md){: .button .primary-btn}