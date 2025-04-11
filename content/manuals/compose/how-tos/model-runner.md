---
title: Use Docker Model
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
This allows you to define and run AI-powered applications alongside your other services.

## Prerequisites

- Docker Compose v2.35 or later
- Docker Desktop 4.41 or later 
- Docker Model Runner enabled in Docker Desktop
- Apple Silicon Mac (currently Model Runner is only available for Mac with Apple Silicon)

## Enabling Docker Model Runner

Before you can use Docker Model Runner with Compose, you need to enable it in Docker Desktop, as described in the [Docker Model Runner documentation](/desktop/features/model-runner/).

## Provider services

Compose introduces a new service type called `provider` that allows you to declare platform capabilities required by your application. For AI models, you can use the `model` type to declare model dependencies.

Here's an example of how to define a model provider:

```yaml
services:
  chat:
    image: my-chat-app
    depends_on:
      - ai-runner

  ai-runner:
    provider:
      type: model
      options:
        model: ai/smollm2
```

You should notice the dedicated `provider` attribute in the `ai-runner` service.   
This attribute specifies that the service is a model provider and let you define options such as the name of the model to be used.

There is also a `depends_on` attribute in the `chat` service.  
This attribute specifies that the `chat` service depends on the `ai-runner` service.  
This means that the `ai-runner` service will be started before the `chat` service to allow injection of model information to the `chat` service.

## How it works

During the `docker compose up` process, Docker Model Runner will automatically pull and run the specified model.  
It will also send to Compose the model tag name and the URL to access the model runner.

Those information will be then pass to services which declare a dependency on the model provider.  
In the example above, the `chat` service will receive 2 env variables prefixed by the service name:
 - `AI-RUNNER_URL` with the URL to access the model runner
 - `AI-RUNNER_MODEL` with the model name which could be passed with the URL to request the model.

This allows the `chat` service to interact with the model and use it for its own purposes.


## Reference

- [Docker Model Runner documentation](/desktop/features/model-runner/)

