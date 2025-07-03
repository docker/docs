---
title: Define AI Models in Docker Compose applications
linkTitle: Use AI models in Compose
description: Learn how to define and use AI models in Docker Compose applications using the models top-level element
keywords: compose, docker compose, models, ai, machine learning, cloud providers, specification
weight: 10
params:
  sidebar:
    badge:
      color: green
      text: New
---

{{< summary-bar feature_name="Compose models" >}}

Compose lets you define AI models as core components of your application, so you can declare model dependencies alongside services and run the application on any platform that supports the Compose Specification.

## Prerequisites

- Docker Compose v2.38 or later
- A platform that supports Compose models such as Docker Model Runner or compatible cloud providers

## What are Compose models?

Compose `models` are a standardized way to define AI model dependencies in your application. By using the []`models` top-level element](/reference/compose-file/models.md) in your Compose file, you can:

- Declare which AI models your application needs
- Specify model configurations and requirements
- Make your application portable across different platforms
- Let the platform handle model provisioning and lifecycle management

## Basic model definition

To define models in your Compose application, use the `models` top-level element:

```yaml
services:
  chat-app:
    image: my-chat-app
    models:
      - llm

models:
  llm:
    image: ai/smollm2
```

This example defines:
- A service called `chat-app` that uses a model named `llm`
- A model definition for `llm` that references the `ai/smollm2` model image

## Model configuration options

Models support various configuration options:

```yaml
models:
  llm:
    image: ai/smollm2
    context_size: 1024
    runtime_flags:
      - "--a-flag"
      - "--another-flag=42"
```

Common configuration options include:
- `model` (required): The OCI artifact identifier for the model. This is what Compose pulls and runs via the model runner. 
- `context_size`: Defines the maximum token context size for the model.
- `runtime_flags`: A list of raw command-line flags passed to the inference engine when the model is started.
-  Platform-specific options may also be available via extensions attributes `x-*`

## Service model binding

Services can reference models in two ways: short syntax and long syntax.

### Short syntax

The short syntax is the simplest way to bind a model to a service:

```yaml
services:
  app:
    image: my-app
    models:
      - llm
      - embedding-model

models:
  llm:
    image: ai/smollm2
  embedding-model:
    image: ai/all-minilm
```

With short syntax, the platform automatically generates environment variables based on the model name:
- `LLM_URL` - URL to access the llm model
- `LLM_MODEL` - Model identifier for the llm model
- `EMBEDDING_MODEL_URL` - URL to access the embedding-model
- `EMBEDDING_MODEL_MODEL` - Model identifier for the embedding-model

### Long syntax

The long syntax allows you to customize environment variable names:

```yaml
services:
  app:
    image: my-app
    models:
      llm:
        endpoint_var: AI_MODEL_URL
        model_var: AI_MODEL_NAME
      embedding-model:
        endpoint_var: EMBEDDING_URL
        model_var: EMBEDDING_NAME

models:
  llm:
    image: ai/smollm2
  embedding-model:
    image: ai/all-minilm
```

With this configuration, your service receives:
- `AI_MODEL_URL` and `AI_MODEL_NAME` for the LLM model
- `EMBEDDING_URL` and `EMBEDDING_NAME` for the embedding model

## Platform portability

One of the key benefits of using Compose models is portability across different platforms that support the Compose specification.

### Docker Model Runner

When Docker Model Runner is enabled:

```yaml
services:
  chat-app:
    image: my-chat-app
    models:
      - llm

models:
  llm:
    image: ai/smollm2
```

Docker Model Runner will:
- Pull and run the specified model locally
- Provide endpoint URLs for accessing the model
- Inject environment variables into the service

### Cloud providers

The same Compose file can run on cloud providers that support Compose models:

```yaml
services:
  chat-app:
    image: my-chat-app
    models:
      - llm

models:
  llm:
    image: ai/smollm2
    # Cloud-specific configurations
    labels:
      - "cloud.instance-type=gpu-small"
      - "cloud.region=us-west-2"
```

Cloud providers might:
- Use managed AI services instead of running models locally
- Apply cloud-specific optimizations and scaling
- Provide additional monitoring and logging capabilities
- Handle model versioning and updates automatically

## Reference

- [`models` top-level element](/reference/compose-file/models.md)
- [`models` attribute](/reference/compose-file/services.md#models)
- [Docker Model Runner documentation](/manuals/ai/model-runner.md)
- [Compose Model Runner documentation](/manuals/ai/compose/model-runner.md)