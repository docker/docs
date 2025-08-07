---
title: Models
description: Learn about the models top-level element
keywords: compose, compose specification, models, compose file reference
weight: 120
---

{{< summary-bar feature_name="Compose models" >}}

The top-level `models` section declares AI models that are used by your Compose application. These models are typically pulled as OCI artifacts, run by a model runner, and exposed as an API that your service containers can consume.

Services can only access models when explicitly granted by a [`models` attribute](services.md#models) within the `services` top-level element.

## Examples

### Example 1

```yaml
services:
  app:
    image: app
    models:
      - ai_model


models:
  ai_model:
    model: ai/model
```

In this basic example:

 - The app service uses the `ai_model`.
 - The `ai_model` is defined as an OCI artifact (`ai/model`) that is pulled and served by the model runner.
 - Docker Compose injects connection info, for example `AI_MODEL_URL`, into the container. 

### Example 2

```yaml
services:
  app:
    image: app
    models:
      my_model:
        endpoint_var: MODEL_URL

models:
  my_model:
    model: ai/model
    context_size: 1024
    runtime_flags: 
      - "--a-flag"
      - "--another-flag=42"
```

In this advanced setup:

 - The service app references `my_model` using the long syntax.
 - Compose injects the model runner's URL as the environment variable `MODEL_URL`.

## Attributes

- `model` (required): The OCI artifact identifier for the model. This is what Compose pulls and runs via the model runner. 
- `context_size`: Defines the maximum token context size for the model.
- `runtime_flags`: A list of raw command-line flags passed to the inference engine when the model is started.

## Additional resources

For more examples and information on using `model`, see [Use AI models in Compose](/manuals/ai/compose/models-and-compose.md)