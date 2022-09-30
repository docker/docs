---
title: "Containerize an application"
keywords: get started, setup, orientation, quickstart, intro, concepts, containers, docker desktop
redirect_from:
- /get-started/part2/
description: Containerize and run a simple application to learn Docker
---


For the rest of this tutorial, you will be working with a simple todo
list manager that's running in Node.js. If you're not familiar with Node.js,
don't worry. This tutorial doesn't require real JavaScript experience.

At this point, your development team is quite small and you're simply
building an app to prove out your MVP (minimum viable product). You want
to show how it works and what it's capable of doing without needing to
think about how it will work for a large team, multiple developers, etc.

![Todo List Manager Screenshot](images/todo-list-sample.png){: style="width:50%;" }

## Get the app

Before you can run the application, you need to get the application source code onto 
your machine. The repository contains the entire getting started tutorial, but you'll be working with the sample application in the `app` folder.

1. Download the contents from the [getting-started repository](https://github.com/docker/getting-started/tree/master){:target="_blank" rel="noopener" class="_"}. You can either pull the entire project or [download it as a zip](https://github.com/docker/getting-started/archive/refs/heads/master.zip) and extract the `app` folder out to get started with.

    > **Note**
    >
    > While the repository contains many other files and folders, the sample application is inside the `app` folder.


2. Once extracted, use your favorite code editor to open the project. If you're in need of an editor, you can use [Visual Studio Code](https://code.visualstudio.com/){:target="_blank" rel="noopener" class="_"}. You should see the following contents inside the `app` folder.

    ![Screenshot of Visual Studio Code opened with the app loaded](images/ide-screenshot.png){: style="width:650px;margin-top:20px;"}
    {: .text-center }


## Build the app's container image

In order to build the application, you need to use a `Dockerfile`. A
Dockerfile is simply a text-based script of instructions that is used to
create a container image. If you've created Dockerfiles before, you might
see a few flaws in the Dockerfile below. But, don't worry. This tutorial will address them later.

1. Create a file named `Dockerfile` in the same folder as the file `package.json`.

   <ul class="nav nav-tabs">
     <li class="active"><a data-toggle="tab" data-target="#tab3">Mac/Linux</a></li>
     <li><a data-toggle="tab" data-target="#tab4">Windows</a></li>
    </ul>
    <div class="tab-content">
    <div id="tab3" class="tab-pane fade in active" markdown="1">
   
    ```console
    $ cd path/to/app
    $ touch Dockerfile
    ```
    <hr>
    </div>
    <div id="tab4" class="tab-pane fade" markdown="1">
   
     ```console
     C:\> cd path\to\app
     C:\path\to\app> type nul >> Dockerfile
     ```
   <hr>
   </div>
   </div>


     Please check that the file `Dockerfile` has no file extension like `.txt`. Some editors may append this file extension automatically and this would result in an error in the next step.

     You should now see the following contents inside the `app` folder.
     ```
     |— app/
       |— spec/
       |— src/
       |— Dockerfile
       |— package.json
       |— yarn.lock
     ```

2. Add the following contents to the Dockerfile.

   ```dockerfile
   # syntax=docker/dockerfile:1
   FROM node:12-alpine
   RUN apk add --no-cache python2 g++ make
   WORKDIR /app
   COPY . .
   RUN yarn install --production
   CMD ["node", "src/index.js"]
   EXPOSE 3000
   ```


3. If you haven't already done so, open a terminal and go to the `app` directory with the `Dockerfile`. Now build the container image using the `docker build` command.

   ```console
   $ cd path/to/app
   $ docker build -t getting-started .
   ```

   This command used the Dockerfile to build a new container image. You might
   have noticed that a lot of "layers" were downloaded. This is because you instructed
   the builder that you wanted to start from the `node:12-alpine` image. But, since you
   didn't have that on your machine, it downloaded the image.

   After it downloaded the image, the instructions in the Dockerfile copied in your application and used `yarn` to 
   install your application's dependencies. The `CMD` directive specifies the default 
   command to run when starting a container from this image.

   Finally, the `-t` flag tags your image. Think of this simply as a human-readable name
   for the final image. Since you named the image `getting-started`, you can refer to that
   image when you run a container.

   The `.` at the end of the `docker build` command tells Docker that it should look for the `Dockerfile` in the current directory.

## Start an app container

Now that you have an image, you can run the application. To do so, you will use the `docker run`
command.

1. Start your container using the `docker run` command and specify the name of the image you 
   just created:

   ```console
   $ docker run -dp 3000:3000 getting-started
   ```

   Remember the `-d` and `-p` flags? You're running the new container in "detached" mode (in the 
   background) and creating a mapping between the host's port 3000 to the container's port 3000.
   Without the port mapping, you wouldn't be able to access the application.

2. After a few seconds, open your web browser to [http://localhost:3000](http://localhost:3000).
   You should see your app.

   ![Empty Todo List](images/todo-list-empty.png){: style="width:450px;margin-top:20px;"}
   {: .text-center }

3. Go ahead and add an item or two and see that it works as you expect. You can mark items as
   complete and remove items. Your frontend is successfully storing items in the backend.

At this point, you should have a running todo list manager with a few items, all built by you.

If you take a quick look at the Docker Dashboard, you should see at least one  container running now.

![Docker Dashboard with tutorial and app containers running](images/dashboard-two-containers.png)

## Next steps

In this short section, you learned the very basics about building a container image and created a
Dockerfile to do so. Once you built an image, you started the container and saw the running app.

Next, you're going to make a modification to your app and learn how to update your running application
with a new image. Along the way, you'll learn a few other useful commands.

[Update the application](03_updating_app.md){: .button .primary-btn}