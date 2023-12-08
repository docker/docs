---
title: Use containers for generative AI development
keywords: python, local, development, generative ai, genai, llm, neo4j, ollama, langchain
description: Learn how to develop your generative AI (GenAI) application locally.
---

## Prerequisites

- Complete [Containerize a generative AI application](containerize.md).
- Your host machine has at least 8 GB of RAM to run the LLM.

## Overview

In this section, you'll learn how to set up a development environment to locally run all the services that your generative AI (GenAI) application needs. This includes:

- Adding a local database
- Adding a local LLM service

> **Note**
>
> You can see more samples of containerized GenAI applications in the [GenAI Stack](https://github.com/docker/genai-stack) demo applications.

## Add a local database

You can use containers to set up local services, like a database. In this section, you'll update the `compose.yaml` file to define a database service. In addition, you'll specify an environment variables file to load the database connection information rather than manually entering the information every time.

In the cloned repository's directory, rename `env.example` file to `.env`. This file contains the environment variables that the containers will use.

In the cloned repository's directory, open the `compose.yaml` file in an IDE or text editor.

In the `compose.yaml` file, you need to do the following:
- Add instructions to run a Neo4j database
- Specify the environment file under the server service in order to pass in the environment variables for the connection

The following is the updated `compose.yaml` file. All comments have been removed.

```yaml
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

> **Note**
>
> To learn more about Neo4j, see the [Neo4j Official Docker Image](https://hub.docker.com/_/neo4j).

Test that your application runs. Inside the `docker-genai-sample` directory, run
the following command in a terminal.

```console
$ docker compose up --build
```

Open a browser and view the application at [http://localhost:8000](http://localhost:8000). You should see a simple Streamlit application. Note that asking questions to a PDF will cause the application to fail because the Ollama service specified in the `.env` file isn't running yet.

In the terminal, press `ctrl`+`c` to stop the application.

## Add a local LLM service

The sample application uses [Ollama](https://hub.docker.com/r/ollama/ollama) as the LLM service. To run Ollama in a container, update your `compose.yaml` file with the following:
- Add the Ollama service
- Add another service to automatically pull the model
- Add a volume to persist the model data

### Add the Ollama service to the Compose file

The following is the updated `compose.yaml` file with the Ollama service and a volume to persist the model.

```yaml {hl_lines=["24-31"]}
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
volumes:
  ollama_volume:
```

### Provide GPU access to Ollama (optional)

The previous updates to the `compose.yaml` file don't provide GPU access to the
Ollama service. Providing GPU access to the Ollama service will improve
performance. To provide GPU access to Ollama, you can do one of the following:
- Run Ollama outside of a container
- Provide GPU access to the container

{{< tabs >}}
{{< tab name="Run Ollama outside of a container" >}}

To run Ollama outside of a container:
1. Remove the Ollama service from your `compose.yaml` file.
2. [Install](https://github.com/jmorganca/ollama) and run Ollama.
3. Update the `OLLAMA_BASE_URL` value in your `.env` file to
   `http://host.docker.internal:1134`.

{{< /tab >}}
{{< tab name="Provide GPU access to the container" >}}

The following steps describe how to provide access to the GPU on Windows 11 or Linux.

> **Note**
>
> You can't provide GPU access to containers on Docker Desktop for Mac. For Mac,
> you can provide GPU access to Ollama by running Ollama outside of a container.

To provide GPU access to the container:
1. Install the prerequisites.
   - For Linux, install the [NVIDIA Container Toolkilt](https://github.com/NVIDIA/nvidia-container-toolkit).
   - For Windows 11, install the latest [NVIDIA driver](https://www.nvidia.com/Download/index.aspx).
2. Update your Ollama service in your `compose.yaml` with the following:
   ```yaml {hl_lines="9-15"}
   services:
     #...
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
     #...
   ```

> **Note**
>
> For more details about the Compose instructions, see [Turn on GPU access with Docker Compose](../../../compose/gpu-support.md).

{{< /tab >}}
{{< /tabs >}}

### Add the model puller service

Ollama requires a model to run. You can automate pulling the model with a script and a service to run the script. To do this, you can use the model puller script from the [GenAI Stack](https://github.com/docker/genai-stack) demo applications.

Create a new file in the `docker-genai-sample` directory called `pull_model.Dockerfile`. Open `pull_model.Dockerfile` in a code or text editor and add the following content:

```dockerfile
#syntax = docker/dockerfile:1.4

FROM ollama/ollama:latest AS ollama
FROM babashka/babashka:latest

# just using as a client - never as a server
COPY --from=ollama /bin/ollama ./bin/ollama

COPY <<EOF pull_model.clj
(ns pull-model
  (:require [babashka.process :as process]
            [clojure.core.async :as async]))

(try
  (let [llm (get (System/getenv) "LLM")
        url (get (System/getenv) "OLLAMA_BASE_URL")]
    (println (format "pulling ollama model %s using %s" llm url))
    (if (and llm url)
      (let [done (async/chan)]
        (async/go-loop [n 0]
          (let [[v _] (async/alts! [done (async/timeout 5000)])]
            (if (= :stop v) :stopped (do (println (format "... pulling model (%ss) - will take several minutes" (* n 10))) (recur (inc n))))))
        (process/shell {:env {"OLLAMA_HOST" url} :out :inherit :err :inherit} (format "./bin/ollama pull %s" llm))
        (async/>!! done :stop))

      (println "OLLAMA model only pulled if both LLM and OLLAMA_BASE_URL are set")))
  (catch Throwable _ (System/exit 1)))
EOF

ENTRYPOINT ["bb", "-f", "pull_model.clj"]
```

The `pull_model.Dockerfile` copies the Ollama client to the Babashka Clojure runtime image, and then defines and runs a Clojure script to pull the model.

To run the pull-model service, update your `compose.yaml` file with the following:
- Add a pull-model service
- Add a condition to the server service so that the server container doesn't start until the model-puller service completes

The following is the new service added to the `compose.yaml` file:

```yaml {hl_lines=["12-18"]}
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
      pull-model:
        condition: service_completed_successfully
  pull-model:
    build:
      dockerfile: pull_model.Dockerfile
    env_file:
      - .env
  # ...
```

## Run your GenAI application

At this point, you have the following services in your Compose file:
- Server service for your main GenAI application
- Database service to store vectors in a Neo4j database
- Ollama service to run the LLM
- Pull-model service to automatically pull the model for the Ollama service

To run all the services, run the following command in your `docker-genai-sample`
directory:

```console
$ docker compose up --build
```

The first time you run the Compose stack, it may take several minutes for the model puller service to pull the model. The puller service will continuously update the console with its status. After pulling the model, the model puller service container will stop and you can access the application.

Once the model puller service finishes, open a browser and access the application at [http://localhost:8000](http://localhost:8000).

Upload a PDF file, for example the [Docker CLI Cheat Sheet](https://docs.docker.com/get-started/docker_cheatsheet.pdf), and ask a question about the PDF.

Keep in mind that it may take several minutes to provide an answer if you aren't providing GPU access to Ollama.

## Summary

In this section, you learned how to set up a development environment to locally run all the services that your GenAI application needs.

Related information:
 - [Dockerfile reference](../../../engine/reference/builder.md)
 - [Compose file reference](../../../compose/compose-file/_index.md)
 - [Ollama Docker image](https://hub.docker.com/r/ollama/ollama)
 - [Neo4j Official Docker Image](https://hub.docker.com/_/neo4j)
 - [GenAI Stack demo applications](https://github.com/docker/genai-stack)

## Next steps

See samples of more GenAI applications in the [GenAI Stack demo applications](https://github.com/docker/genai-stack).