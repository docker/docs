---
title: Containerize a GenAI video transcription app
keywords: python, generative ai, genai, llm, pinecone, openai, whisper, langchain
description: Learn how to containerize generative AI video transcription applications.
---
## Overview

In this section you're going to walk through the process of containerizing a video transcription app that uses generative AI.

## Prerequisites

Work through the steps of the [Explore and run a genAI video transcription app](run.md) section to get and configure the example application used in this section.

## Step 1: Prepare the example application

Before you walk through the process of creating Docker assets, delete the
current assets in the example application. This section only explores the
yt-whisper application, so delete the existing `Dockerfile` inside the
`docker-genai/yt-whisper/` directory.

You should now have the following contents inside the `docker-genai/yt-whisper/` directory.

```text
├── docker-genai/yt-whisper/
│ ├── scripts/
│ ├── tests/
│ ├── yt_whisper/
│ ├── README.md
│ ├── poetry.lock
│ └── pyproject.toml
```

## Step 2: Initialize Docker assets

Use `docker init` to create the necessary Docker assets to containerize your application. Inside the `docker-genai/yt-whisper/` directory, run the `docker init` command in a terminal. `docker init` provides some default configuration, but you'll need to answer a few questions about your application. For example, this application uses Streamlit to run. Refer to the following docker init example and use the same answers for your prompts.

```console
docker init
Welcome to the Docker Init CLI!

This utility will walk you through creating the following files with sensible defaults for your project:
  - .dockerignore
  - Dockerfile
  - compose.yaml
  - README.Docker.md

Let's get started!

? What application platform does your project use? Python
? What version of Python do you want to use? 3.11
? What port do you want your app to listen on? 8503
? What is the command to run your app? streamlit run yt_whisper/app.py --server.port=8503 --server.address=0.0.0.0
```

You should now have the following contents in your `docker-genai/yt-whisper/` directory.

```text
├── docker-genai/yt-whisper/
│ ├── scripts/
│ ├── tests/
│ ├── yt_whisper/
│ ├── README.md
│ ├── poetry.lock
│ ├── pyproject.toml
│ ├── README.Docker.md
│ ├── .dockerignore
│ ├── compose.yaml
│ └── Dockerfile
```

## Step 3: Explore and update the Docker assets

`docker init` creates Docker assets to help you get started. Depending on your application, you may need to modify the assets. In the following sections, you'll explore and update the Docker assets.

### Explore and update the Dockerfile

1. Open the `docker-genai/yt-whisper/Dockerfile` in a code or text editor.
2. Inspect the contents of the `Dockerfile`. First notice that the `Dockerfile`
   is more extensive than the Dockerfile from
   [Explore and run a genAI video transcription app](run.md). It's more
   extensive because `docker init` implements several [best practices](../../../develop/develop-images/dockerfile_best-practices.md)
   for creating production images. Next, notice the instructions specify a
   `requirements.txt` file for the packages. This particular application uses
   Poetry, so you'll need to update the relevant instructions.
3. Update the `Dockerfile` for Poetry. The following is the updated
   `Dockerfile`.
   ```dockerfile{hl_lines=["36-40"]}
   # syntax=docker/dockerfile:1

   # Comments are provided throughout this file to help you get started.
   # If you need more help, visit the Dockerfile reference guide at
   # https://docs.docker.com/go/dockerfile-reference/
   
   ARG PYTHON_VERSION=3.11
   FROM python:${PYTHON_VERSION}-slim as base
   
   # Prevents Python from writing pyc files.
   ENV PYTHONDONTWRITEBYTECODE=1
   
   # Keeps Python from buffering stdout and stderr to avoid situations where
   # the application crashes without emitting any logs due to buffering.
   ENV PYTHONUNBUFFERED=1
   
   WORKDIR /app
   
   # Create a non-privileged user that the app will run under.
   # See https://docs.docker.com/go/dockerfile-user-best-practices/
   ARG UID=10001
   RUN adduser \
       --disabled-password \
       --gecos "" \
       --home "/nonexistent" \
       --shell "/sbin/nologin" \
       --no-create-home \
       --uid "${UID}" \
       appuser
   
   # Download dependencies as a separate step to take advantage of Docker's caching.
   # Leverage a cache mount to /root/.cache/pip to speed up subsequent builds.
   # Leverage a bind mount to requirements.txt to avoid having to copy them into
   # into this layer.
   RUN --mount=type=cache,target=/root/.cache/pip \
       --mount=type=bind,source=pyproject.toml,target=pyproject.toml \
       --mount=type=bind,source=poetry.lock,target=poetry.lock \
       --mount=type=bind,source=yt_whisper,target=yt_whisper \
       --mount=type=bind,source=README.md,target=README.md \
       python -m pip install .
   
   # Switch to the non-privileged user to run the application.
   USER appuser
   
   # Copy the source code into the container.
   COPY . .
   
   # Expose the port that the application listens on.
   EXPOSE 8503
   
   # Run the application.
   CMD streamlit run yt_whisper/app.py --server.port=8503 --server.address=0.0.0.0
   ```
4. Save and close the file.

### Explore and update the .dockerignore file

One new file that you haven't explored yet is the `.dockerignore` file. The `.dockerignore` file excludes files and directories from the image. For more details, see [.dockerignore files](../../../build/building/context.md#dockerignore-files).

1. Open the `docker-genai/yt-whisper/.dockerignore` file in a code or text
   editor.
2. Inspect the contents of the `.dockerignore` file. Notice that it has
   `README.md`, but you just specified copying this file in the previous
   `Dockerfile`.
3. Update the `.dockerignore` file and remove `README.md`. The following is the
   updated `.dockerignore` file.
   ```text
   # Include any files or directories that you don't want to be copied to your
   # container here (e.g., local build artifacts, temporary files, etc.).
   #
   # For more help, visit the .dockerignore file reference guide at
   # https://docs.docker.com/go/build-context-dockerignore/
   
   **/.DS_Store
   **/__pycache__
   **/.venv
   **/.classpath
   **/.env
   **/.git
   **/.gitignore
   **/.project
   **/.settings
   **/.toolstarget
   **/.vs
   **/.vscode
   **/*.*proj.user
   **/*.dbmdl
   **/*.jfm
   **/bin
   **/charts
   **/docker-compose*
   **/compose*
   **/Dockerfile*
   **/node_modules
   **/npm-debug.log
   **/obj
   **/secrets.dev.yaml
   **/values.dev.yaml
   LICENSE
   ```
4. Save and close the file.

### Explore and update the Compose file

1. Open the `docker-genai/yt-whisper/compose.yaml` file in a code or text
   editor.
2. Inspect the contents of the `compose.yaml` file. Notice that the environment
   variables file isn't specified. Also, notice that only one service named `server` is specified.
3. Update the `compose.yaml` file and add the environment variables file. The
   following is the updated `compose.yaml` file.
   ```yaml{hl_lines=["16-17"]}
   # Comments are provided throughout this file to help you get started.
   # If you need more help, visit the Docker Compose reference guide at
   # https://docs.docker.com/go/compose-spec-reference/
   
   # Here the instructions define your application as a service called "server".
   # This service is built from the Dockerfile in the current directory.
   # You can add other services your application may depend on here, such as a
   # database or a cache. For examples, see the Awesome Compose repository:
   # https://github.com/docker/awesome-compose
   services:
     server:
       build:
         context: .
       ports:
         - 8503:8503
       env_file:
         - ../.env
    # ...
   ```

## Step 4: Run the containerized generative AI application

To run the application, in a terminal, change directory to the `docker-genai/yt-whisper/` directory and run the following command.

```console
docker compose up --build
```

Docker builds the image and runs it as a container. Once the container starts, you can access the yt-whisper application by opening a web browser and navigating to [localhost:8503](http://localhost:8503).

Only the service for the yt-whisper app is started. You should be unable to access the bot at [localhost:8504](http://localhost:8504).

Stop the container by pressing `CTRL`+`C` in the terminal.

To start both the bot and yt-whisper services, you can use the exisiting Compose file in the root of `docker-genai/` directory. Change directory to `docker-genai/` and run the following command to bring up both services.

```console
docker compose up --build
```

Stop the containers by pressing `CTRL`+`C` in the terminal.

## Summary

In this section you learned how to containerize and run a generative AI application.

Related information:
* [docker init CLI reference](../../../engine/reference/commandline/init.md)
* [Dockerfile reference](/engine/reference/builder/)
* [Compose overview](../../../compose/_index.md)
* [.dockerignore files](../../../build/building/context.md#dockerignore-files)

## Next steps

* Try to containerize the docker-bot in the example application using `docker
  init`.
* See the [Python language-specific guide](../../../language/python/_index.md) to learn how to configure CI/CD and deploy Python apps locally to Kubernetes.
* See the [GenAI Stack](https://github.com/docker/genai-stack) demo applications for more examples of containerized generative AI applications.