---
title: Run your own container
keywords: get started, quick start, intro, concepts
description: Learn how to run a container from scratch
---

In this guide you'll learn the basic steps to run a container from scratch. For this guide, you'll use a sample Node.js application, but it's not necessary to know Node.js.

## Step 1: Get the sample application

If you have git, use the following command in a terminal to clone the sample application repository.

```console
$ git clone https://github.com/docker/welcome-to-docker
```

If you don't have git, [download the source](https://github.com/docker/
welcome-to-docker/archive/refs/heads/main.zip) and extract it.

## Step 2: Create a Dockerfile in your project folder

To run your code in a container, the most fundamental thing you need is a Dockerfile. A Dockerfile describes what goes into a container. To add a Dockerfile, create a text file called `Dockerfile` with no file extension in the root directory of your project and add the following contents.

```Dockerfile
# Start your image with a node base image
FROM node:18-alpine

# Create an application directory
RUN mkdir -p /app

# Set the /app directory as the working directory for any command that follows
WORKDIR /app

# Copy the local app package and package-lock.json file to the container
COPY package*.json ./

# Copy local directories to the working directory of our docker image (/app)
COPY ./src ./src
COPY ./public ./public

# Install node packages, install serve, build the app, and remove dependencies at the end
RUN npm install \
    && npm install -g serve \
    && npm run build \
    && rm -fr node_modules

# Specify that the application in the container will listen on port 80
EXPOSE 80

# Start the app using serve command
CMD [ "serve", "-s", "build" ]
```

## Step 3: Build your first image
An image is like a static version of a container. You always need an image to run a container. Once you have a Dockerfile in your repository, run the following `docker build` command in the project folder to create an image.

```console
$ docker build -t welcome-to-docker .
```

> **Breaking down the `docker build` command**
>
> Here are what the different parts of the `docker build` command do:
> - `docker build`: This command builds the image. It needs one argument, the source folder for the Dockerfile that needs to be built. In this case, it’s the Dockerfile in the current folder, `.`.
> - `-t welcome-to-docker`: The `-t` flag tags the image with a unique name. In this case, `welcome-to-docker`.

## Step 4: Run your container

Now that you have your image, use the following `docker run` command to see your container in action.

```console
$ docker run -p 8089:80 welcome-to-docker
```

> **Breaking down the `docker run` command**
>
> Here are what the different parts of the `docker run` command do:
> - `docker run`: This is used to run containers. It needs at least one argument, and that argument is the image you want to run. In this case, it's `welcome-to-docker`.
> - `-p 8089:80`: This lets Docker know that port 80 in the container needs to be accessible from port 8089 on your local host.

## Step 5: Monitor your container

You can use Docker Desktop to monitor your running containers. Go to the **Containers** tab to view the container you just ran. Select the container name to view the logs, stats, and more.

![Viewing the Containers tab in Docker Desktop](images/getting-started-monitor.gif){:width="500px"}

## What's next

In this guide, you built your own image. When running containers on Docker Desktop, you don’t need to build your own image from scratch. You can also run images created by others on Docker Hub.

> **Note**
>
> If you want to learn more about creating images for applications in other languages, check out the following language-specific guides:
> - [Node.js](../language/nodejs/index.md)
> - [Python](../language/python/index.md)
> - [Go](../language/golang/index.md)
> - [Java](../language/java/index.md)
> - [C# (.NET)](../language/dotnet/index.md)


[Run Docker Hub images](run-docker-hub-images.md){: .button .primary-btn}

