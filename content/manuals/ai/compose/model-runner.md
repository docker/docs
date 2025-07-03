---
title: Use Docker Model Runner
description: Learn how to integrate Docker Model Runner with Docker Compose to build AI-powered applications
keywords: compose, docker compose, model runner, ai, llm, artificial intelligence, machine learning
weight: 20
aliases:
 - /compose/how-tos/model-runner/
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

- Docker Compose v2.38 or later
- Docker Desktop 4.43 or later 
- Docker Desktop for Mac with Apple Silicon or Docker Desktop for Windows with NVIDIA GPU
- [Docker Model Runner enabled in Docker Desktop](/manuals/ai/model-runner.md#enable-docker-model-runner)

## Use `models` definition

The [`models` top-level element](/manuals/ai/compose/models-and-compose.md) in the Compose file lets you define AI models to be used by your application.
Compose can then use Docker Model Runner as the model runtime. 

The following example shows how to provide the minimal configuration to use a model within your Compose application:

```yaml
services:
  my-chat-app:
    image: my-chat-app
    models:
      - smollm2

models:
  smollm2:
    model: ai/smollm2
```

### How it works

During the `docker compose up` process, Docker Model Runner automatically pulls and runs the specified model.  
It also sends Compose the model tag name and the URL to access the model runner.

This information is then passed to services which declare a dependency on the model provider.  
In the example above, the `my-chat-app` service receives 2 environment variables prefixed by the service name:
- `SMOLLM2_ENDPOINT` with the URL to access the model
- `SMOLLM2_MODEL` with the model name

This lets the `my-chat-app` service to interact with the model and use it for its own purposes.

### Customizing environment variables

You can customize the environment variable names which will be passed to your service container using the long syntax:

```yaml
services:
  my-chat-app:
    image: my-chat-app
    models:
      smollm2:
        endpoint_var: AI_MODEL_URL
        model_var: AI_MODEL_NAME

models:
  smollm2:
    model: ai/smollm2
```

With this configuration, your `my-chat-app` service will receive:
- `AI_MODEL_URL` with the URL to access the model
- `AI_MODEL_NAME` with the model name

This allows you to use more descriptive variable names that match your application's expectations.


## Alternative configuration with Provider services

> [!TIP]
>
> Use the []`models` top-level element](#use-models-definition) instead.

Compose introduced a new service type called `provider` that allows you to declare platform capabilities required by your application. For AI models, you can use the `model` type to declare model dependencies.

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
```

Notice the dedicated `provider` attribute in the `ai_runner` service.   
This attribute specifies that the service is a model provider and lets you define options such as the name of the model to be used.

There is also a `depends_on` attribute in the `my-chat-app` service.  
This attribute specifies that the `my-chat-app` service depends on the `ai_runner` service.  
This means that the `ai_runner` service will be started before the `my-chat-app` service to allow injection of model information to the `my-chat-app` service.

## Reference

- [Docker Model Runner documentation](/manuals/ai/model-runner.md)
