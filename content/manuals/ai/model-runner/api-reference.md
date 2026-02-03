---
title: DMR REST API
description: Reference documentation for the Docker Model Runner REST API endpoints, including OpenAI, Anthropic, and Ollama compatibility.
weight: 30
keywords: Docker, ai, model runner, rest api, openai, anthropic, ollama, endpoints, documentation, cline, continue, cursor
---

Once Model Runner is enabled, new API endpoints are available. You can use
these endpoints to interact with a model programmatically. Docker Model Runner
provides compatibility with OpenAI, Anthropic, and Ollama API formats.

## Determine the base URL

The base URL to interact with the endpoints depends on how you run Docker and
which API format you're using.

{{< tabs >}}
{{< tab name="Docker Desktop">}}

| Access from | Base URL |
|-------------|----------|
| Containers | `http://model-runner.docker.internal` |
| Host processes (TCP) | `http://localhost:12434` |

> [!NOTE]
> TCP host access must be enabled. See [Enable Docker Model Runner](get-started.md#enable-docker-model-runner-in-docker-desktop).

{{< /tab >}}
{{< tab name="Docker Engine">}}

| Access from | Base URL |
|-------------|----------|
| Containers | `http://172.17.0.1:12434` |
| Host processes | `http://localhost:12434` |

> [!NOTE]
> The `172.17.0.1` interface may not be available by default to containers
  within a Compose project.
> In this case, add an `extra_hosts` directive to your Compose service YAML:
>
> ```yaml
> extra_hosts:
>   - "model-runner.docker.internal:host-gateway"
> ```
> Then you can access the Docker Model Runner APIs at `http://model-runner.docker.internal:12434/`

{{< /tab >}}
{{</tabs >}}

### Base URLs for third-party tools

When configuring third-party tools that expect OpenAI-compatible APIs, use these base URLs:

| Tool type | Base URL format |
|-----------|-----------------|
| OpenAI SDK / clients | `http://localhost:12434/engines/v1` |
| Anthropic SDK / clients | `http://localhost:12434` |
| Ollama-compatible clients | `http://localhost:12434` |

See [IDE and tool integrations](ide-integrations.md) for specific configuration examples.

## Supported APIs

Docker Model Runner supports multiple API formats:

| API | Description | Use case |
|-----|-------------|----------|
| [OpenAI API](#openai-compatible-api) | OpenAI-compatible chat completions, embeddings | Most AI frameworks and tools |
| [Anthropic API](#anthropic-compatible-api) | Anthropic-compatible messages endpoint | Tools built for Claude |
| [Ollama API](#ollama-compatible-api) | Ollama-compatible endpoints | Tools built for Ollama |
| [Image Generation API](#image-generation-api-diffusers) | Diffusers-based image generation | Generating images from text prompts |
| [DMR API](#dmr-native-endpoints) | Native Docker Model Runner endpoints | Model management |

## OpenAI-compatible API

DMR implements the OpenAI API specification for maximum compatibility with existing tools and frameworks.

### Endpoints

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/engines/v1/models` | GET | [List models](https://platform.openai.com/docs/api-reference/models/list) |
| `/engines/v1/models/{namespace}/{name}` | GET | [Retrieve model](https://platform.openai.com/docs/api-reference/models/retrieve) |
| `/engines/v1/chat/completions` | POST | [Create chat completion](https://platform.openai.com/docs/api-reference/chat/create) |
| `/engines/v1/completions` | POST | [Create completion](https://platform.openai.com/docs/api-reference/completions/create) |
| `/engines/v1/embeddings` | POST | [Create embeddings](https://platform.openai.com/docs/api-reference/embeddings/create) |

> [!NOTE]
> You can optionally include the engine name in the path: `/engines/llama.cpp/v1/chat/completions`.
> This is useful when running multiple inference engines.

### Model name format

When specifying a model in API requests, use the full model identifier including the namespace:

```json
{
  "model": "ai/smollm2",
  "messages": [...]
}
```

Common model name formats:
- Docker Hub models: `ai/smollm2`, `ai/llama3.2`, `ai/qwen2.5-coder`
- Tagged versions: `ai/smollm2:360M-Q4_K_M`
- Custom models: `myorg/mymodel`

### Supported parameters

The following OpenAI API parameters are supported:

| Parameter | Type | Description |
|-----------|------|-------------|
| `model` | string | Required. The model identifier. |
| `messages` | array | Required for chat completions. The conversation history. |
| `prompt` | string | Required for completions. The prompt text. |
| `max_tokens` | integer | Maximum tokens to generate. |
| `temperature` | float | Sampling temperature (0.0-2.0). |
| `top_p` | float | Nucleus sampling parameter (0.0-1.0). |
| `stream` | Boolean | Enable streaming responses. |
| `stop` | string/array | Stop sequences. |
| `presence_penalty` | float | Presence penalty (-2.0 to 2.0). |
| `frequency_penalty` | float | Frequency penalty (-2.0 to 2.0). |

### Limitations and differences from OpenAI

Be aware of these differences when using DMR's OpenAI-compatible API:

| Feature | DMR behavior |
|---------|--------------|
| API key | Not required. DMR ignores the `Authorization` header. |
| Function calling | Supported with llama.cpp for compatible models. |
| Vision | Supported for multi-modal models (e.g., LLaVA). |
| JSON mode | Supported via `response_format: {"type": "json_object"}`. |
| Logprobs | Supported. |
| Token counting | Uses the model's native token encoder, which may differ from OpenAI's. |

## Anthropic-compatible API

DMR provides [Anthropic Messages API](https://platform.claude.com/docs/en/api/messages) compatibility for tools and frameworks built for Claude.

### Endpoints

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/anthropic/v1/messages` | POST | [Create a message](https://platform.claude.com/docs/en/api/messages/create) |
| `/anthropic/v1/messages/count_tokens` | POST | [Count tokens](https://docs.anthropic.com/en/api/messages-count-tokens) |

### Supported parameters

The following Anthropic API parameters are supported:

| Parameter | Type | Description |
|-----------|------|-------------|
| `model` | string | Required. The model identifier. |
| `messages` | array | Required. The conversation messages. |
| `max_tokens` | integer | Maximum tokens to generate. |
| `temperature` | float | Sampling temperature (0.0-1.0). |
| `top_p` | float | Nucleus sampling parameter. |
| `top_k` | integer | Top-k sampling parameter. |
| `stream` | Boolean | Enable streaming responses. |
| `stop_sequences` | array | Custom stop sequences. |
| `system` | string | System prompt. |

### Example: Chat with Anthropic API

```bash
curl http://localhost:12434/v1/messages \
  -H "Content-Type: application/json" \
  -d '{
    "model": "ai/smollm2",
    "max_tokens": 1024,
    "messages": [
      {"role": "user", "content": "Hello!"}
    ]
  }'
```

### Example: Streaming response

```bash
curl http://localhost:12434/v1/messages \
  -H "Content-Type: application/json" \
  -d '{
    "model": "ai/smollm2",
    "max_tokens": 1024,
    "stream": true,
    "messages": [
      {"role": "user", "content": "Count from 1 to 10"}
    ]
  }'
```

## Ollama-compatible API

DMR also provides Ollama-compatible endpoints for tools and frameworks built for Ollama.

### Endpoints

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/api/tags` | GET | List available models |
| `/api/show` | POST | Show model information |
| `/api/chat` | POST | Generate chat completion |
| `/api/generate` | POST | Generate completion |
| `/api/embeddings` | POST | Generate embeddings |

### Example: Chat with Ollama API

```bash
curl http://localhost:12434/api/chat \
  -H "Content-Type: application/json" \
  -d '{
    "model": "ai/smollm2",
    "messages": [
      {"role": "user", "content": "Hello!"}
    ]
  }'
```

### Example: List models

```bash
curl http://localhost:12434/api/tags
```

## Image generation API (Diffusers)

DMR supports image generation through the Diffusers backend, enabling you to generate
images from text prompts using models like Stable Diffusion.

> [!NOTE]
> The Diffusers backend requires an NVIDIA GPU with CUDA support and is only
> available on Linux (x86_64 and ARM64). See [Inference engines](inference-engines.md#diffusers)
> for setup instructions.

### Endpoint

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/engines/diffusers/v1/images/generations` | POST | Generate an image from a text prompt |

### Supported parameters

| Parameter | Type | Description |
|-----------|------|-------------|
| `model` | string | Required. The model identifier (e.g., `stable-diffusion:Q4`). |
| `prompt` | string | Required. The text description of the image to generate. |
| `size` | string | Image dimensions in `WIDTHxHEIGHT` format (e.g., `512x512`). |

### Response format

The API returns a JSON response with the generated image encoded in base64:

```json
{
  "data": [
    {
      "b64_json": "<base64-encoded-image-data>"
    }
  ]
}
```

### Example: Generate an image

```bash
curl -s -X POST http://localhost:12434/engines/diffusers/v1/images/generations \
  -H "Content-Type: application/json" \
  -d '{
    "model": "stable-diffusion:Q4",
    "prompt": "A picture of a nice cat",
    "size": "512x512"
  }' | jq -r '.data[0].b64_json' | base64 -d > image.png
```

This command:
1. Sends a POST request to the Diffusers image generation endpoint
2. Specifies the model, prompt, and output image size
3. Extracts the base64-encoded image from the response using `jq`
4. Decodes the base64 data and saves it as `image.png`


## DMR native endpoints

These endpoints are specific to Docker Model Runner for model management:

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/models/create` | POST | Pull/create a model |
| `/models` | GET | List local models |
| `/models/{namespace}/{name}` | GET | Get model details |
| `/models/{namespace}/{name}` | DELETE | Delete a local model |

## REST API examples

### Request from within a container

To call the `chat/completions` OpenAI endpoint from within another container using `curl`:

```bash
#!/bin/sh

curl http://model-runner.docker.internal/engines/v1/chat/completions \
    -H "Content-Type: application/json" \
    -d '{
        "model": "ai/smollm2",
        "messages": [
            {
                "role": "system",
                "content": "You are a helpful assistant."
            },
            {
                "role": "user",
                "content": "Please write 500 words about the fall of Rome."
            }
        ]
    }'

```

### Request from the host using TCP

To call the `chat/completions` OpenAI endpoint from the host via TCP:

1. Enable the host-side TCP support from the Docker Desktop GUI, or via the [Docker Desktop CLI](/manuals/desktop/features/desktop-cli.md).
   For example: `docker desktop enable model-runner --tcp <port>`.

   If you are running on Windows, also enable GPU-backed inference.
   See [Enable Docker Model Runner](get-started.md#enable-docker-model-runner-in-docker-desktop).

1. Interact with it as documented in the previous section using `localhost` and the correct port.

```bash
#!/bin/sh

curl http://localhost:12434/engines/v1/chat/completions \
  -H "Content-Type: application/json" \
  -d '{
      "model": "ai/smollm2",
      "messages": [
          {
              "role": "system",
              "content": "You are a helpful assistant."
          },
          {
              "role": "user",
              "content": "Please write 500 words about the fall of Rome."
          }
      ]
  }'
```

### Request from the host using a Unix socket

To call the `chat/completions` OpenAI endpoint through the Docker socket from the host using `curl`:

```bash
#!/bin/sh

curl --unix-socket $HOME/.docker/run/docker.sock \
    localhost/exp/vDD4.40/engines/v1/chat/completions \
    -H "Content-Type: application/json" \
    -d '{
        "model": "ai/smollm2",
        "messages": [
            {
                "role": "system",
                "content": "You are a helpful assistant."
            },
            {
                "role": "user",
                "content": "Please write 500 words about the fall of Rome."
            }
        ]
    }'
```

### Streaming responses

To receive streaming responses, set `stream: true`:

```bash
curl http://localhost:12434/engines/v1/chat/completions \
  -H "Content-Type: application/json" \
  -d '{
      "model": "ai/smollm2",
      "stream": true,
      "messages": [
          {"role": "user", "content": "Count from 1 to 10"}
      ]
  }'
```

## Using with OpenAI SDKs

### Python

```python
from openai import OpenAI

client = OpenAI(
    base_url="http://localhost:12434/engines/v1",
    api_key="not-needed"  # DMR doesn't require an API key
)

response = client.chat.completions.create(
    model="ai/smollm2",
    messages=[
        {"role": "user", "content": "Hello!"}
    ]
)

print(response.choices[0].message.content)
```

### Node.js

```javascript
import OpenAI from 'openai';

const client = new OpenAI({
  baseURL: 'http://localhost:12434/engines/v1',
  apiKey: 'not-needed',
});

const response = await client.chat.completions.create({
  model: 'ai/smollm2',
  messages: [{ role: 'user', content: 'Hello!' }],
});

console.log(response.choices[0].message.content);
```

## What's next

- [IDE and tool integrations](ide-integrations.md) - Configure Cline, Continue, Cursor, and other tools
- [Configuration options](configuration.md) - Adjust context size and runtime parameters
- [Inference engines](inference-engines.md) - Learn about llama.cpp, vLLM, and Diffusers options
