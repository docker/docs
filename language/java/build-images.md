---
title: "Build your Java image"
keywords: Java, build, images, dockerfile
description: Learn how to build your first Docker image by writing a Dockerfile
---

{% include_relative nav.html selected="1" %}

## Prerequisites

Work through the orientation and setup in Get started [Part 1](../../get-started/index.md){:target="_blank" rel="noopener" class="_"} to understand Docker concepts. Refer to the following section for Java prerequisites.

## Overview

Now that we have a good overview of containers and the Docker platform, let’s take a look at building our first image. An image includes everything needed to run an application - the code or binary, runtime, dependencies, and any other file system objects required.

To complete this tutorial, you need the following:

- Java OpenJDK version 15 or later. [Download and install Java](https://jdk.java.net/){: target="_blank" rel="noopener" class="_"}
- Docker running locally. Follow the instructions to [download and install Docker](../../get-docker.md)
- A Git client
- An IDE or a text editor to edit files. We recommend using [IntelliJ Community Edition](https://www.jetbrains.com/idea/download/){: target="_blank" rel="noopener" class="_"}.

## Sample application

Let’s clone the sample application that we'll be using in this module to our local development machine. Run the following commands in a terminal to clone the repo.

```console
$ cd /path/to/working/directory
$ git clone https://github.com/spring-projects/spring-petclinic.git
$ cd spring-petclinic
```

## Test the application

Let’s start our application and make sure it is running properly. Maven will manage all the project processes (compiling, tests, packaging, etc). The **Spring Pets Clinic** project you cloned earlier contains an embedded version of Maven. Therefore, you don't need to install Maven separately on your local machine.

Open your terminal and navigate to the working directory you created and run the following command:

```console
$ ./mvnw spring-boot:run
```

This downloads the dependencies, builds the project, and starts it.

To test that the application is working properly, open a new browser and navigate to `http://localhost:8080`.

Switch back to the terminal where our server is running and you should see the following requests in the server logs. The data will be different on your machine.

```console
o.s.s.petclinic.PetClinicApplication     : Started
PetClinicApplication in 11.743 seconds (JVM running for 12.364)
```

## Create a Dockerfile for Java

Now that our application is running properly, let’s take a look at creating a Dockerfile.

A Dockerfile is a text document that contains all the commands a user could call on the command line to assemble an image. When we tell Docker to build our image by executing the `docker build` command, Docker reads these instructions and execute them sequentially and creates a Docker image.

Let’s walk through the steps on creating a Dockerfile for our application. In the root of your working directory, create a file named `Dockerfile` and open this file in your text editor.

> **Note**
>
> The name of the Dockerfile is not important but the default filename for many commands is simply `Dockerfile`. Therefore, we’ll use that as our filename throughout this series.

The first line to add to the Dockerfile is a [`# syntax` parser directive](/engine/reference/builder/#syntax).
While _optional_, this directive instructs the Docker builder what syntax to use
when parsing the Dockerfile, and allows older Docker versions with BuildKit enabled
to upgrade the parser before starting the build. Parser directives
must appear before any other comment, whitespace, or Dockerfile instruction in
your Dockerfile, and should be the first line in Dockerfiles.

```dockerfile
# syntax=docker/dockerfile:1
```

We recommend using `docker/dockerfile:1`, which always points to the latest release
of the version 1 syntax. BuildKit automatically checks for updates of the syntax
before building, making sure you are using the most current version.

Next, we need to add a line in our Dockerfile that tells Docker what base image
we would like to use for our application.

```dockerfile
# syntax=docker/dockerfile:1

FROM openjdk:16-alpine3.13
```

Docker images can be inherited from other images. Therefore, instead of creating our own base image, we’ll use the official Maven image with Java JDK that already has all the tools and packages that we need to run a Java application.

> **Note**
>
> To learn more about creating your own base images, see [Creating base images](https://docs.docker.com/develop/develop-images/baseimages/).

To make things easier when running the rest of our commands, let’s create a working directory. This instructs Docker to use this path as the default location for all subsequent commands. By doing this, we do not have to type out full file paths but can use relative paths based on the working directory.

```dockerfile
WORKDIR /app
```

Usually, the very first thing you do once you’ve downloaded a project written in Java which is using Maven for project management is to install dependencies.

Before we can run `mvnw dependency`, we need to get the Maven wrapper and our `pom.xml` file into our image. We’ll use the `COPY` command to do this. The `COPY` command takes two parameters. The first parameter tells Docker what file(s) you would like to copy into the image. The second parameter tells Docker where you want that file(s) to be copied to. We’ll copy all those files and directories into our working directory - `/app`.

```dockerfile
COPY .mvn/ .mvn
COPY mvnw pom.xml ./
```

Once we have our `pom.xml` file inside the image, we can use the `RUN` command to execute the command `mvnw dependency:go-offline`. This works exactly the same way as if we were running `mvnw` (or mvn) dependency locally on our machine, but this time the dependencies will be installed into the image.

```dockerfile
RUN ./mvnw dependency:go-offline
```

At this point, we have an Alpine version 3.13 image that is based on OpenJDK version 16, and we have also installed our dependencies. The next thing we need to do is to add our source code into the image. We’ll use the `COPY` command just like we did with our `pom.xml` file above.

```dockerfile
COPY src ./src
```

This `COPY` command takes all the files located in the current directory and copies them into the image. Now, all we have to do is to tell Docker what command we want to run when our image is executed inside a container. We do this using the `CMD` command.

```dockerfile
CMD ["./mvnw", "spring-boot:run"]
```

Here's the complete Dockerfile.

```dockerfile
# syntax=docker/dockerfile:1

FROM openjdk:16-alpine3.13

WORKDIR /app

COPY .mvn/ .mvn
COPY mvnw pom.xml ./
RUN ./mvnw dependency:go-offline

COPY src ./src

CMD ["./mvnw", "spring-boot:run"]
```

### Create a `.dockerignore` file

To use a file in the build context, the Dockerfile refers to the file specified in an instruction, for example, a `COPY` instruction. To increase the performance of the build, and to exclude files and directories, we recommend that you create a `.dockerignore` file to the context directory. To improve the context load time, add a `target` directory within the `.dockerignore` file.

## Build an image

Now that we’ve created our Dockerfile, let’s build our image. To do this, we use the `docker build` command. The `docker build` command builds Docker images from a Dockerfile and a “context”. A build’s context is the set of files located in the specified PATH or URL. The Docker build process can access any of the files located in this context.

The build command optionally takes a `--tag` flag. The tag is used to set the name of the image and an optional tag in the format `name:tag`. We’ll leave off the optional `tag` for now to help simplify things. If you do not pass a tag, Docker uses “latest” as its default tag. You can see this in the last line of the build output.

Let’s build our first Docker image.

```console
$ docker build --tag java-docker .
```

```console
Sending build context to Docker daemon  5.632kB
Step 1/7 : FROM java:3.7-alpine
Step 2/7 : WORKDIR /app
...
Successfully built a0bb458aabd0
Successfully tagged java-docker:latest
```

## View local images

To see a list of images we have on our local machine, we have two options. One is to use the CLI and the other is to use [Docker Desktop](../../desktop/dashboard.md#explore-your-images). As we are currently working in the terminal let’s take a look at listing images using the CLI.

To list images, simply run the `docker images` command.

```console
$ docker images
REPOSITORY          TAG                 IMAGE ID            CREATED          SIZE
java-docker         latest              b1b5f29f74f0        47 minutes ago   567MB
```

You should see at least the we just built `java-docker:latest`.

## Tag images

An image name is made up of slash-separated name components. Name components may contain lowercase letters, digits, and separators. A separator is defined as a period, one or two underscores, or one or more dashes. A name component may not start or end with a separator.

An image is made up of a manifest and a list of layers. Do not worry too much about manifests and layers at this point other than a “tag” points to a combination of these artifacts. You can have multiple tags for an image. Let’s create a second tag for the image we built and take a look at its layers.

To create a new tag for the image we’ve built above, run the following command:

```console
$ docker tag java-docker:latest java-docker:v1.0.0
```

The `docker tag` command creates a new tag for an image. It does not create a new image. The tag points to the same image and is just another way to reference the image.

Now, run the `docker images` command to see a list of our local images.

```console
$ docker images
REPOSITORY    TAG      IMAGE ID		  CREATED		  SIZE
java-docker   latest   b1b5f29f74f0	  59 minutes ago	567MB
java-docker   v1.0.0   b1b5f29f74f0	  59 minutes ago	567MB
```

You can see that we have two images that start with `java-docker`. We know they are the same image because if you take a look at the `IMAGE ID` column, you can see that the values are the same for the two images.

Let’s remove the tag that we just created. To do this, we’ll use the `rmi` command. The `rmi` command stands for “remove image”.

```console
$ docker rmi java-docker:v1.0.0
Untagged: java-docker:v1.0.0
```

Note that the response from Docker tells us that the image has not been removed but only “untagged”. You can check this by running the `docker images` command.

```console
$ docker images
REPOSITORY      TAG     IMAGE ID        CREATED              SIZE
java-docker    	latest	b1b5f29f74f0	59 minutes ago	     567MB
```

Our image that was tagged with `:v1.0.0` has been removed, but we still have the `java-docker:latest` tag available on our machine.

## Next steps

In this module, we took a look at setting up our example Java application that we'll use for the rest of the tutorial. We also created a Dockerfile that we used to build our Docker image. Then, we took a look at tagging our images and removing images. In the next module, we’ll take a look at how to:

[Run your image as a container](run-containers.md){: .button .primary-btn}

## Feedback

Help us improve this topic by providing your feedback. Let us know what you think by creating an issue in the [Docker Docs](https://github.com/docker/docker.github.io/issues/new?title=[Java%20docs%20feedback]){:target="_blank" rel="noopener" class="_"} GitHub repository. Alternatively, [create a PR](https://github.com/docker/docker.github.io/pulls){:target="_blank" rel="noopener" class="_"} to suggest updates.

<br />
