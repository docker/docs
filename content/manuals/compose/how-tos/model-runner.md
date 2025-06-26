---
title: Use Docker Model Runner
description: Learn how to integrate Docker Model Runner with Docker Compose to build AI-powered applications
keywords: compose, docker compose, model runner, ai, llm, artificial intelligence, machine learning
weight: 111
params:
  sidebar:
    badge:
      color: green
      text: New
---

{{< summary-bar feature_name="Compose model runner" >}}

Docker Model Runner can be integrated with Docker Compose to run AI models as part of your multi-container applications.  
This lets you define and run AI-powered applications alongside your other services.

## Prerequisites

- Docker Compose v2.35 or later
- Docker Desktop 4.41 or later 
- Docker Desktop for Mac with Apple Silicon or Docker Desktop for Windows with NVIDIA GPU
- [Docker Model Runner enabled in Docker Desktop](/manuals/ai/model-runner.md#enable-docker-model-runner)

## Provider services

Compose introduces a new service type called `provider` that allows you to declare platform capabilities required by your application. For AI models, you can use the `model` type to declare model dependencies.

Here's an example of how to define a model provider:

```yaml
services:
  chat:
    image: my-chat-app
    depends_on:
      - ai_runner

  ai_runner:
    provider:
      type: model
      options:
        model: ai/smollm2
        context-size: 1024
```

Notice the following:

- In the `ai_runner` service:

  - `provider.type`: Specifies that the service is a `model` provider.
  - `provider.options`: Specifies the options of the model. In our case, we want to use
    `ai/smollm2`, and we set the context size to 1024 tokens.

     > [!NOTE]
     > Each model has its own maximum context size. When increasing the context length,
     > consider your hardware constraints. In general, try to use the smallest context size
     > possible for your use case.
   
- In the `chat` service:
   
  -  `depends_on` specifies that the `chat` service depends on the `ai_runner` service. The
     `ai_runner` service will be started before the `chat` service, to allow injection of model information to the `chat` service.
   
## How it works

During the `docker compose up` process, Docker Model Runner automatically pulls and runs the specified model.  
It also sends Compose the model tag name and the URL to access the model runner.

This information is then passed to services which declare a dependency on the model provider.  
In the example above, the `chat` service receives 2 environment variables prefixed by the service name:
 - `AI_RUNNER_URL` with the URL to access the model runner
 - `AI_RUNNER_MODEL` with the model name which could be passed with the URL to request the model.

This lets the `chat` service to interact with the model and use it for its own purposes.

## Related pages

- [Docker Model Runner documentation](/manuals/ai/model-runner.md)
