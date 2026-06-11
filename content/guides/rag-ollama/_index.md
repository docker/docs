---
description: Containerize RAG application using Ollama and Docker
keywords: python, generative ai, genai, llm, ollama, rag, qdrant
title: Build a RAG application using Ollama and Docker
linkTitle: RAG Ollama application
summary: |
  This guide demonstrates how to use Docker to deploy Retrieval-Augmented
  Generation (RAG) models with Ollama.
aliases:
  - /guides/use-case/rag-ollama/
  - /guides/rag-ollama/containerize/
  - /guides/rag-ollama/develop/
params:
  tags: [ai]
  time: 20 minutes
---


The Retrieval Augmented Generation (RAG) guide teaches you how to containerize an existing RAG application using Docker. The example application is a RAG that acts like a sommelier, giving you the best pairings between wines and food. In this guide, you’ll learn how to:

- Containerize and run a RAG application
- Set up a local environment to run the complete RAG stack locally for development

Start by containerizing an existing RAG application.

## Containerize a RAG application

### Overview

This section walks you through containerizing a RAG application using Docker.

> [!NOTE]
> You can see more samples of containerized GenAI applications in the [GenAI Stack](https://github.com/docker/genai-stack) demo applications.

### Get the sample application

The sample application used in this guide is an example of RAG application, made by three main components, which are the building blocks for every RAG application. A Large Language Model hosted somewhere, in this case it is hosted in a container and served via [Ollama](https://ollama.ai/). A vector database, [Qdrant](https://qdrant.tech/), to store the embeddings of local data, and a web application, using [Streamlit](https://streamlit.io/) to offer the best user experience to the user.

Clone the sample application. Open a terminal, change directory to a directory that you want to work in, and run the following command to clone the repository:

```console
$ git clone https://github.com/mfranzon/winy.git
```

You should now have the following files in your `winy` directory.

```text
├── winy/
│ ├── .gitignore
│ ├── app/
│ │ ├── main.py
│ │ ├── Dockerfile
| | └── requirements.txt
│ ├── tools/
│ │ ├── create_db.py
│ │ ├── create_embeddings.py
│ │ ├── requirements.txt
│ │ ├── test.py
| | └── download_model.sh
│ ├── docker-compose.yaml
│ ├── wine_database.db
│ ├── LICENSE
│ └── README.md
```

### Containerizing your application: Essentials

Containerizing an application involves packaging it along with its dependencies into a container, which ensures consistency across different environments. Here’s what you need to containerize an app like Winy :

1. Dockerfile: A Dockerfile that contains instructions on how to build a Docker image for your application. It specifies the base image, dependencies, configuration files, and the command to run your application.

2. Docker Compose File: Docker Compose is a tool for defining and running multi-container Docker applications. A Compose file allows you to configure your application's services, networks, and volumes in a single file.

### Run the application

Inside the `winy` directory, run the following command in a
terminal.

```console
$ docker compose up --build
```

Docker builds and runs your application. Depending on your network connection, it may take several minutes to download all the dependencies. You'll see a message like the following in the terminal when the application is running.

```console
server-1  |   You can now view your Streamlit app in your browser.
server-1  |
server-1  |   URL: http://0.0.0.0:8501
server-1  |
```

Open a browser and view the application at [http://localhost:8501](http://localhost:8501). You should see a simple Streamlit application.

The application requires a Qdrant database service and an LLM service to work properly. If you have access to services that you ran outside of Docker, specify the connection information in the `docker-compose.yaml`.

```yaml
winy:
  build:
    context: ./app
    dockerfile: Dockerfile
  environment:
    - QDRANT_CLIENT=http://qdrant:6333 # Specifies the url for the qdrant database
    - OLLAMA=http://ollama:11434 # Specifies the url for the ollama service
  container_name: winy
  ports:
    - "8501:8501"
  depends_on:
    - qdrant
    - ollama
```

If you don't have the services running, continue with this guide to learn how you can run some or all of these services with Docker.
Remember that the `ollama` service is empty; it doesn't have any model. For this reason you need to pull a model before starting to use the RAG application. All the instructions are in the following page.

In the terminal, press `ctrl`+`c` to stop the application.

### Summary

In this section, you learned how you can containerize and run your RAG
application using Docker.

### Next steps

In the next section, you'll learn how to properly configure the application with your preferred LLM model, completely locally, using Docker.

## Use containers for RAG development

### Prerequisites

Complete [Containerize a RAG application](containerize.md).

### Overview

In this section, you'll learn how to set up a development environment to access all the services that your generative RAG application needs. This includes:

- Adding a local database
- Adding a local or remote LLM service

> [!NOTE]
> You can see more samples of containerized GenAI applications in the [GenAI Stack](https://github.com/docker/genai-stack) demo applications.

### Add a local database

You can use containers to set up local services, like a database. In this section, you'll explore the database service in the `docker-compose.yaml` file.

To run the database service:

1. In the cloned repository's directory, open the `docker-compose.yaml` file in an IDE or text editor.

2. In the `docker-compose.yaml` file, you'll see the following:

   ```yaml
   services:
     qdrant:
       image: qdrant/qdrant
       container_name: qdrant
       ports:
         - "6333:6333"
       volumes:
         - qdrant_data:/qdrant/storage
   ```

   > [!NOTE]
   > To learn more about Qdrant, see the [Qdrant Official Docker Image](https://hub.docker.com/r/qdrant/qdrant).

3. Start the application. Inside the `winy` directory, run the following command in a terminal.

   ```console
   $ docker compose up --build
   ```

4. Access the application. Open a browser and view the application at [http://localhost:8501](http://localhost:8501). You should see a simple Streamlit application.

5. Stop the application. In the terminal, press `ctrl`+`c` to stop the application.

### Add a local or remote LLM service

The sample application supports both [Ollama](https://ollama.ai/). This guide provides instructions for the following scenarios:

- Run Ollama in a container
- Run Ollama outside of a container

While all platforms can use any of the previous scenarios, the performance and
GPU support may vary. You can use the following guidelines to help you choose the appropriate option:

- Run Ollama in a container if you're on Linux, and using a native installation of the Docker Engine, or Windows 10/11, and using Docker Desktop, you
  have a CUDA-supported GPU, and your system has at least 8 GB of RAM.
- Run Ollama outside of a container if running Docker Desktop on a Linux Machine.

Choose one of the following options for your LLM service.

{{< tabs >}}
{{< tab name="Run Ollama in a container" >}}

When running Ollama in a container, you should have a CUDA-supported GPU. While you can run Ollama in a container without a supported GPU, the performance may not be acceptable. Only Linux and Windows 11 support GPU access to containers.

To run Ollama in a container and provide GPU access:

1. Install the prerequisites.
   - For Docker Engine on Linux, install the [NVIDIA Container Toolkilt](https://github.com/NVIDIA/nvidia-container-toolkit).
   - For Docker Desktop on Windows 10/11, install the latest [NVIDIA driver](https://www.nvidia.com/Download/index.aspx) and make sure you are using the [WSL2 backend](/manuals/desktop/features/wsl/_index.md#turn-on-docker-desktop-wsl-2)
2. The `docker-compose.yaml` file already contains the necessary instructions. In your own apps, you'll need to add the Ollama service in your `docker-compose.yaml`. The following is
   the updated `docker-compose.yaml`:

   ```yaml
   ollama:
     image: ollama/ollama
     container_name: ollama
     ports:
       - "8000:8000"
     deploy:
       resources:
         reservations:
           devices:
             - driver: nvidia
               count: 1
               capabilities: [gpu]
   ```

   > [!NOTE]
   > For more details about the Compose instructions, see [Turn on GPU access with Docker Compose](/manuals/compose/how-tos/gpu-support.md).

3. Once the Ollama container is up and running it is possible to use the `download_model.sh` inside the `tools` folder with this command:

   ```console
   . ./download_model.sh <model-name>
   ```

Pulling an Ollama model can take several minutes.

{{< /tab >}}
{{< tab name="Run Ollama outside of a container" >}}

To run Ollama outside of a container:

1. [Install](https://github.com/jmorganca/ollama) and run Ollama on your host
   machine.
2. Pull the model to Ollama using the following command.

   ```console
   $ ollama pull llama2
   ```

3. Remove the `ollama` service from the `docker-compose.yaml` and update properly the connection variables in `winy` service:

   ```diff
   - OLLAMA=http://ollama:11434
   + OLLAMA=<your-url>
   ```

{{< /tab >}}
{{< /tabs >}}

### Run your RAG application

At this point, you have the following services in your Compose file:

- Server service for your main RAG application
- Database service to store vectors in a Qdrant database
- (optional) Ollama service to run the LLM
  service

Once the application is running, open a browser and access the application at [http://localhost:8501](http://localhost:8501).

Depending on your system and the LLM service that you chose, it may take several
minutes to answer.

### Summary

In this section, you learned how to set up a development environment to provide
access all the services that your GenAI application needs.

Related information:

- [Dockerfile reference](/reference/dockerfile.md)
- [Compose file reference](/reference/compose-file/_index.md)
- [Ollama Docker image](https://hub.docker.com/r/ollama/ollama)
- [GenAI Stack demo applications](https://github.com/docker/genai-stack)

### Next steps

See samples of more GenAI applications in the [GenAI Stack demo applications](https://github.com/docker/genai-stack).
