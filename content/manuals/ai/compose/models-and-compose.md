---
title: Define AI Models in Docker Compose applications
linkTitle: Use AI models in Compose
description: Learn how to define and use AI models in Docker Compose applications using the models top-level element
keywords: compose, docker compose, models, ai, machine learning, cloud providers, specification
aliases:
  - /compose/how-tos/model-runner/
  - /ai/compose/model-runner/
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
- A platform that supports Compose models such as Docker Model Runner (DMR) or compatible cloud providers.
  If you are using DMR, see the [requirements](/manuals/ai/model-runner/_index.md#requirements).

## What are Compose models?

Compose `models` are a standardized way to define AI model dependencies in your application. By using the [`models` top-level element](/reference/compose-file/models.md) in your Compose file, you can:

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
    model: ai/smollm2
```

This example defines:
- A service called `chat-app` that uses a model named `llm`
- A model definition for `llm` that references the `ai/smollm2` model image

## Model configuration options

Models support various configuration options:

```yaml
models:
  llm:
    model: ai/smollm2
    context_size: 1024
    runtime_flags:
      - "--a-flag"
      - "--another-flag=42"
```

Common configuration options include:
- `model` (required): The OCI artifact identifier for the model. This is what Compose pulls and runs via the model runner.
- `context_size`: Defines the maximum token context size for the model.

   > [!NOTE]
   > Each model has its own maximum context size. When increasing the context length,
   > consider your hardware constraints. In general, try to keep context size
   > as small as feasible for your specific needs.

- `runtime_flags`: A list of raw command-line flags passed to the inference engine when the model is started.
   See [Configuration options](/manuals/ai/model-runner/configuration.md) for commonly used parameters and examples.
- Platform-specific options may also be available via extension attributes `x-*`

> [!TIP]
> See more example in the [Common runtime configurations](#common-runtime-configurations) section.

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
    model: ai/smollm2
  embedding-model:
    model: ai/all-minilm
```

With short syntax, the platform automatically generates environment variables based on the model name:
- `LLM_URL` - URL to access the LLM model
- `LLM_MODEL` - Model identifier for the LLM model
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
    model: ai/smollm2
  embedding-model:
    model: ai/all-minilm
```

With this configuration, your service receives:
- `AI_MODEL_URL` and `AI_MODEL_NAME` for the LLM model
- `EMBEDDING_URL` and `EMBEDDING_NAME` for the embedding model

## Platform portability

One of the key benefits of using Compose models is portability across different platforms that support the Compose specification.

### Docker Model Runner

When [Docker Model Runner is enabled](/manuals/ai/model-runner/_index.md):

```yaml
services:
  chat-app:
    image: my-chat-app
    models:
      llm:
        endpoint_var: AI_MODEL_URL
        model_var: AI_MODEL_NAME

models:
  llm:
    model: ai/smollm2
    context_size: 4096
    runtime_flags:
      - "--no-prefill-assistant"
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
    model: ai/smollm2
    # Cloud-specific configurations
    x-cloud-options:
      - "cloud.instance-type=gpu-small"
      - "cloud.region=us-west-2"
```

Cloud providers might:
- Use managed AI services instead of running models locally
- Apply cloud-specific optimizations and scaling
- Provide additional monitoring and logging capabilities
- Handle model versioning and updates automatically

## Common runtime configurations

Below are some example configurations for various use cases.

### Development

```yaml
services:
  app:
    image: app
    models:
      dev_model:
        endpoint_var: DEV_URL
        model_var: DEV_MODEL

models:
  dev_model:
    model: ai/model
    context_size: 4096
    runtime_flags:
      - "--verbose"                       # Set verbosity level to infinity
      - "--verbose-prompt"                # Print a verbose prompt before generation
      - "--log-prefix"                    # Enable prefix in log messages
      - "--log-timestamps"                # Enable timestamps in log messages
      - "--log-colors"                    # Enable colored logging
```

### Conservative with disabled reasoning

```yaml
services:
  app:
    image: app
    models:
      conservative_model:
        endpoint_var: CONSERVATIVE_URL
        model_var: CONSERVATIVE_MODEL

models:
  conservative_model:
    model: ai/model
    context_size: 4096
    runtime_flags:
      - "--temp"                # Temperature
      - "0.1"
      - "--top-k"               # Top-k sampling
      - "1"
      - "--reasoning-budget"    # Disable reasoning
      - "0"
```

### Creative with high randomness

```yaml
services:
  app:
    image: app
    models:
      creative_model:
        endpoint_var: CREATIVE_URL
        model_var: CREATIVE_MODEL

models:
  creative_model:
    model: ai/model
    context_size: 4096
    runtime_flags:
      - "--temp"                # Temperature
      - "1"
      - "--top-p"               # Top-p sampling
      - "0.9"
```

### Highly deterministic

```yaml
services:
  app:
    image: app
    models:
      deterministic_model:
        endpoint_var: DET_URL
        model_var: DET_MODEL

models:
  deterministic_model:
    model: ai/model
    context_size: 4096
    runtime_flags:
      - "--temp"                # Temperature
      - "0"
      - "--top-k"               # Top-k sampling
      - "1"
```

### Concurrent processing

```yaml
services:
  app:
    image: app
    models:
      concurrent_model:
        endpoint_var: CONCURRENT_URL
        model_var: CONCURRENT_MODEL

models:
  concurrent_model:
    model: ai/model
    context_size: 2048
    runtime_flags:
      - "--threads"             # Number of threads to use during generation
      - "8"
      - "--mlock"               # Lock memory to prevent swapping
```

### Rich vocabulary model

```yaml
services:
  app:
    image: app
    models:
      rich_vocab_model:
        endpoint_var: RICH_VOCAB_URL
        model_var: RICH_VOCAB_MODEL

models:
  rich_vocab_model:
    model: ai/model
    context_size: 4096
    runtime_flags:
      - "--temp"                # Temperature
      - "0.1"
      - "--top-p"               # Top-p sampling
      - "0.9"
```

### Embeddings

When using embedding models with the `/v1/embeddings` endpoint, you must include the `--embeddings` runtime flag for the model to be properly configured.

```yaml
services:
  app:
    image: app
    models:
      embedding_model:
        endpoint_var: EMBEDDING_URL
        model_var: EMBEDDING_MODEL

models:
  embedding_model:
    model: ai/all-minilm
    context_size: 2048
    runtime_flags:
      - "--embeddings"          # Required for embedding models
```

## Alternative configuration with provider services

> [!IMPORTANT]
>
> This approach is deprecated. Use the [`models` top-level element](#basic-model-definition) instead.

You can also use the `provider` service type, which allows you to declare platform capabilities required by your application.
For AI models, you can use the `model` type to declare model dependencies.

To define a model provider:

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
        runtime-flags: "--no-prefill-assistant"
```

## Reference

- [`models` top-level element](/reference/compose-file/models.md)
- [`models` attribute](/reference/compose-file/services.md#models)
- [Docker Model Runner documentation](/manuals/ai/model-runner/_index.md)
- [Configuration options](/manuals/ai/model-runner/configuration.md) - Context size and runtime parameters
- [Inference engines](/manuals/ai/model-runner/inference-engines.md) - llama.cpp and vLLM details
- [API reference](/manuals/ai/model-runner/api-reference.md) - OpenAI and Ollama-compatible APIs
