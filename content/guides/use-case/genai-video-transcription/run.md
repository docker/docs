---
title: Explore and run a GenAI video transcription app
keywords: containers, images, python, genai, pinecone, llm, openai, whisper, dockerfiles, build
description: Learn how to build and run containerized generative AI applications that transcribes videos
---

## Overview

In this section you're going to look at a containerized generative AI application that does video transcription. You'll learn how to build and run the app using Docker.

## Prerequisites

* You have installed the latest version of [Docker Desktop](../../../get-docker.md). Docker adds new features regularly and some parts of this guide may work only with the latest version of Docker Desktop.
* You have a [git client](https://git-scm.com/downloads). The examples in this section use a command-line based git client, but you can use any client.
* You have a [Pinecone](https://www.pinecone.io/) and an [OpenAI](https://openai.com/) account.
   > **Important**
   >
   > OpenAI and Pinecone are third-party hosted services and charges may apply
   > when using the services.
   { .important }

## Step 1: Clone and meet the example app
Run the following command in a terminal to clone the application repository to your local machine.

```console
$ git clone https://github.com/Davidnet/docker-genai.git
```

This sample repository contains the following two applications:

* docker-bot: A question-answering service that leverages both a vector database and an AI model to provide responses.
* yt-whisper: A YouTube video processing service that uses the OpenAI Whisper model to generate transcriptions of videos and stores them in a Pinecone vector database.

The applications use the following third-party services, tools, and
packages:
* [OpenAI](https://openai.com/) for chat completion and speech recognition using
  Whisper.
* [Langchain](https://www.langchain.com/) for managing and orchestrating the
  models.
* [Streamlit](https://streamlit.io/) for the UI.
* [Pinecone](https://www.pinecone.io/) for the embeddings in a vector database.
* [pytube](https://github.com/pytube/pytube) for downloading the YouTube videos.
* [Poetry](https://python-poetry.org/) for dependency management.

This guide explores the yt-whisper application, but you can also apply the instructions to the docker-bot application.

## Step 2: Explore the generative AI application

The example applications in the repository have already been containerized. Before you run the applications, inspect the source code and Docker assets.

### Explore the app.py file

In the `docker-genai/yt-whisiper/yt_whisper` directory is the `app.py` file. This is the entrypoint, or main script for the application. Open the `app.py` file in a code or text editor.

The main thing to notice about this app is that it requires the following
environment variables to run: `OPENAI_TOKEN`, `PINECONE_TOKEN`, and
`PINECONE_ENVIRONMENT`. In one of the following steps, you'll learn more about
these environment variables and how you can specify these environment variables
when running the container.

### Explore the Dockerfile

In the `docker-genai/yt-whisiper/yt_whisper` directory is the `Dockerfile`. The `Dockerfile` is a text document that contains all the instructions used to create an image. Open the `Dockerfile` in a code or text editor.

The following describes the instructions used in this `Dockerfile`:
* `FROM` specifies the base image to use. The base image is the initial image that your Dockerfile is based on. In this case, it's the [Docker Official Image for Python](https://hub.docker.com/_/python). You can explore [Docker Hub](https://hub.docker.com/) to find more pre-made images.
* `WORKDIR` sets the working directory for any instructions that follow it.
* `COPY` copies files into the image from your host machine.
* `RUN` runs commands when building the image, such as installing packages.
* `EXPOSE` informs Docker that the container listens on the specified network
  ports at runtime.
* `HEALTHCHECK` tells Docker how to test a container to check that it's still
  working.
* `ENTRYPOINT` specifies which process to run when the container is started.

For more details about the Dockerfile instructions, see the [Dockerfile reference](/engine/reference/builder/).

### Explore the Compose file

In the `docker-genai/` directory is the `docker-compose.yaml` file. The Compose file is a YAML file that you can use to configure your application's services. Open the `docker-compose.yaml` file in a code or text editor.

This Compose file specifies two services, `bot` and `yt-whisper`. Under those services it defines where the Dockerfiles are located, which ports to expose, and the environment file that contains environment variables that the application needs.

One thing to notice is that there is no `.env` file contained in the repository.
Create the `.env` file now.

## Step 3: Create the .env file

You can use a `.env` file to [set environment variables with Compose](../../../compose/environment-variables/set-environment-variables.md).
In the `docker-genai/` directory, create a new file named `.env`. Open the
`.env` file in a code or text editor and specify the following environment
variables.

```text
#----------------------------------------------------------------------------
# OpenAI
#----------------------------------------------------------------------------
OPENAI_TOKEN=your-api-key # Replace your-api-key with your personal API key

#----------------------------------------------------------------------------
# Pinecone
#----------------------------------------------------------------------------
PINECONE_TOKEN=your-api-key # Replace your-api-key with your personal API key
PINECONE_ENVIRONMENT=us-west1-gcp-free
```

To learn more about the values of the environment variables, see the following:
* `OPENAI_TOKEN` is your [OpenAI API key](https://help.openai.com/en/articles/4936850-where-do-i-find-my-api-key).
* `PINECONE_TOKEN` is your [Pinecone API key](https://docs.pinecone.io/docs/authentication).
* `PINECONE_ENVIRONMENT` is the [Pinecone cloud environment](https://docs.pinecone.io/docs/projects#project-environment).

## Step 4: Run the generative AI application

To build and run the application, in a terminal, change directory to the `docker-genai/` directory and run the following command.

```console
docker compose up --build
```

Docker Compose builds the images and runs them as containers. Depending on your network connection, it may take several minutes to download the dependencies.

You should see output similar to the following in the terminal after Docker starts the containers.

```console
bot-1         |   You can now view your Streamlit app in your browser.
bot-1         |
bot-1         |   URL: http://0.0.0.0:8504
bot-1         |
yt-whisper-1  |
yt-whisper-1  | Collecting usage statistics. To deactivate, set browser.gatherUsageStats to False.
yt-whisper-1  |
yt-whisper-1  |
yt-whisper-1  |   You can now view your Streamlit app in your browser.
yt-whisper-1  |
yt-whisper-1  |   URL: http://0.0.0.0:8503
yt-whisper-1  |
```

Once the containers start, you can access the yt-whisper application by opening a web browser and navigating to [localhost:8503](http://localhost:8503). Specify the URL to a short YouTube video, for example the Docker in 100 seconds video at [https://www.youtube.com/watch?v=IXifQ8mX8DE](https://www.youtube.com/watch?v=IXifQ8mX8DE), and then select **Submit**.

Once the video has been processed, open the docker-bot application at [localhost:8504](http://localhost:8504). Ask a question about your video and the bot answers.

Stop the containers by pressing `CTRL`+`C` in the terminal.

## Summary

At this point, you have explored the Docker assets required to build and run a containerized application. You can create the assets from scratch, as the author of this application did, or use the `docker init`command to help get the process started.

Related information:
* [Dockerfile reference](/engine/reference/builder/)
* [Compose overview](../../../compose/_index.md)

## Next steps

Continue to the next section to learn how you can containerize generative AI applications using Docker.

{{< button text="Containerize a GenAI app" url="containerize.md" >}}