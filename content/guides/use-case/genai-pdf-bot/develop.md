---
title: Use containers for generative AI development
keywords: python, local, development, generative ai, genai, llm, neo4j, ollama, langchain, openai
description: Learn how to develop your generative AI (GenAI) application locally.
---

## Prerequisites

Complete [Containerize a generative AI application](containerize.md).

## Overview

In this section, you'll learn how to set up a development environment to access all the services that your generative AI (GenAI) application needs. This includes:

- Adding a local database
- Adding a local or remote LLM service

> [!NOTE]
>
> You can see more samples of containerized GenAI applications in the [GenAI Stack](https://github.com/docker/genai-stack) demo applications.

## Add a local database

You can use containers to set up local services, like a database. In this section, you'll update the `compose.yaml` file to define a database service. In addition, you'll specify an environment variables file to load the database connection information rather than manually entering the information every time.

To run the database service:

1. In the cloned repository's directory, rename `env.example` file to `.env`.
   This file contains the environment variables that the containers will use.
2. In the cloned repository's directory, open the `compose.yaml` file in an IDE or text editor.
3. In the `compose.yaml` file, add the following:
   - Add instructions to run a Neo4j database
   - Specify the environment file under the server service in order to pass in the environment variables for the connection

   The following is the updated `compose.yaml` file. All comments have been removed.

   ```yaml{hl_lines=["7-23"]}
   services:
     server:
       build:
         context: .
       ports:
         - 8000:8000
       env_file:
         - .env
       depends_on:
         database:
           condition: service_healthy
     database:
       image: neo4j:5.11
       ports:
         - "7474:7474"
         - "7687:7687"
       environment:
         - NEO4J_AUTH=${NEO4J_USERNAME}/${NEO4J_PASSWORD}
       healthcheck:
         test: ["CMD-SHELL", "wget --no-verbose --tries=1 --spider localhost:7474 || exit 1"]
         interval: 5s
         timeout: 3s
         retries: 5
   ```

   > [!NOTE]
   >
   > To learn more about Neo4j, see the [Neo4j Official Docker Image](https://hub.docker.com/_/neo4j).

4. Run the application. Inside the `docker-genai-sample` directory,
run the following command in a terminal.

   ```console
   $ docker compose up --build
   ```

5. Access the application. Open a browser and view the application at [http://localhost:8000](http://localhost:8000). You should see a simple Streamlit application. Note that asking questions to a PDF will cause the application to fail because the LLM service specified in the `.env` file isn't running yet.

6. Stop the application. In the terminal, press `ctrl`+`c` to stop the application.

## Add a local or remote LLM service

The sample application supports both [Ollama](https://ollama.ai/) and [OpenAI](https://openai.com/). This guide provides instructions for the following scenarios:
- Run Ollama in a container
- Run Ollama outside of a container
- Use OpenAI

While all platforms can use any of the previous scenarios, the performance and
GPU support may vary. You can use the following guidelines to help you choose the appropriate option:
- Run Ollama in a container if you're on Linux, and using a native installation of the Docker Engine, or Windows 10/11, and using Docker Desktop, you
  have a CUDA-supported GPU, and your system has at least 8 GB of RAM.
- Run Ollama outside of a container if you're on an Apple silicon Mac.
- Use OpenAI if the previous two scenarios don't apply to you.

Choose one of the following options for your LLM service.

{{< tabs >}}
{{< tab name="Run Ollama in a container" >}}

When running Ollama in a container, you should have a CUDA-supported GPU. While you can run Ollama in a container without a supported GPU, the performance may not be acceptable. Only Linux and Windows 11 support GPU access to containers.

To run Ollama in a container and provide GPU access:
1. Install the prerequisites.
   - For Docker Engine on Linux, install the [NVIDIA Container Toolkilt](https://github.com/NVIDIA/nvidia-container-toolkit).
   - For Docker Desktop on Windows 10/11, install the latest [NVIDIA driver](https://www.nvidia.com/Download/index.aspx) and make sure you are using the [WSL2 backend](../../../desktop/wsl/index.md/#turn-on-docker-desktop-wsl-2)
2. Add the Ollama service and a volume in your `compose.yaml`. The following is
   the updated `compose.yaml`:

   ```yaml {hl_lines=["24-38"]}
   services:
     server:
       build:
         context: .
       ports:
         - 8000:8000
       env_file:
         - .env
       depends_on:
         database:
           condition: service_healthy
     database:
       image: neo4j:5.11
       ports:
         - "7474:7474"
         - "7687:7687"
       environment:
         - NEO4J_AUTH=${NEO4J_USERNAME}/${NEO4J_PASSWORD}
       healthcheck:
         test: ["CMD-SHELL", "wget --no-verbose --tries=1 --spider localhost:7474 || exit 1"]
         interval: 5s
         timeout: 3s
         retries: 5
     ollama:
       image: ollama/ollama:latest
       ports:
         - "11434:11434"
       volumes:
         - ollama_volume:/root/.ollama
       deploy:
         resources:
           reservations:
             devices:
               - driver: nvidia
                 count: all
                 capabilities: [gpu]
   volumes:
     ollama_volume:
   ```

   > [!NOTE]
   >
   > For more details about the Compose instructions, see [Turn on GPU access with Docker Compose](../../../compose/gpu-support.md).

3. Add the ollama-pull service to your `compose.yaml` file. This service uses
   the `docker/genai:ollama-pull` image, based on the GenAI Stack's
   [pull_model.Dockerfile](https://github.com/docker/genai-stack/blob/main/pull_model.Dockerfile).
   The service will automatically pull the model for your Ollama
   container. The following is the updated section of the `compose.yaml` file:

   ```yaml {hl_lines=["12-17"]}
   services:
     server:
       build:
         context: .
       ports:
         - 8000:8000
       env_file:
         - .env
       depends_on:
         database:
           condition: service_healthy
         ollama-pull:
           condition: service_completed_successfully
     ollama-pull:
       image: docker/genai:ollama-pull
       env_file:
         - .env
     # ...
   ```

{{< /tab >}}
{{< tab name="Run Ollama outside of a container" >}}

To run Ollama outside of a container:
1. [Install](https://github.com/jmorganca/ollama) and run Ollama on your host
   machine.
2. Update the `OLLAMA_BASE_URL` value in your `.env` file to
   `http://host.docker.internal:11434`.
3. Pull the model to Ollama using the following command.
   ```console
   $ ollama pull llama2
   ```

{{< /tab >}}
{{< tab name="Use OpenAI" >}}

> [!IMPORTANT]
>
> Using OpenAI requires an [OpenAI account](https://platform.openai.com/login). OpenAI is a third-party hosted service and charges may apply.

1. Update the `LLM` value in your `.env` file to
   `gpt-3.5`.
2. Uncomment and update the `OPENAI_API_KEY` value in your `.env` file to
   your [OpenAI API key](https://help.openai.com/en/articles/4936850-where-do-i-find-my-api-key).

{{< /tab >}}
{{< /tabs >}}

## Run your GenAI application

At this point, you have the following services in your Compose file:
- Server service for your main GenAI application
- Database service to store vectors in a Neo4j database
- (optional) Ollama service to run the LLM
- (optional) Ollama-pull service to automatically pull the model for the Ollama
  service

To run all the services, run the following command in your `docker-genai-sample`
directory:

```console
$ docker compose up --build
```

If your Compose file has the ollama-pull service, it may take several minutes for the ollama-pull service to pull the model. The ollama-pull service will continuously update the console with its status. After pulling the model, the ollama-pull service container will stop and you can access the application.

Once the application is running, open a browser and access the application at [http://localhost:8000](http://localhost:8000).

Upload a PDF file, for example the [Docker CLI Cheat Sheet](https://docs.docker.com/get-started/docker_cheatsheet.pdf), and ask a question about the PDF.

Depending on your system and the LLM service that you chose, it may take several
minutes to answer. If you are using Ollama and the performance isn't
acceptable, try using OpenAI.

## Summary

In this section, you learned how to set up a development environment to provide
access all the services that your GenAI application needs.

Related information:
 - [Dockerfile reference](../../../reference/dockerfile.md)
 - [Compose file reference](/reference/compose-file/_index.md)
 - [Ollama Docker image](https://hub.docker.com/r/ollama/ollama)
 - [Neo4j Official Docker Image](https://hub.docker.com/_/neo4j)
 - [GenAI Stack demo applications](https://github.com/docker/genai-stack)

## Next steps

See samples of more GenAI applications in the [GenAI Stack demo applications](https://github.com/docker/genai-stack).
