---
title: DMR REST API
description: Reference documentation for the Docker Model Runner REST API endpoints and usage examples.
weight: 30
keywords: Docker, ai, model runner, rest api, openai, endpoints, documentation
---

Once Model Runner is enabled, new API endpoints are available. You can use
these endpoints to interact with a model programmatically.

### Determine the base URL

The base URL to interact with the endpoints depends
on how you run Docker:

{{< tabs >}}
{{< tab name="Docker Desktop">}}

- From containers: `http://model-runner.docker.internal/`
- From host processes: `http://localhost:12434/`, assuming TCP host access is
  enabled on the default port (12434).

{{< /tab >}}
{{< tab name="Docker Engine">}}

- From containers: `http://172.17.0.1:12434/` (with `172.17.0.1` representing the host gateway address)
- From host processes: `http://localhost:12434/`

> [!NOTE]
> The `172.17.0.1` interface may not be available by default to containers
  within a Compose project.
> In this case, add an `extra_hosts` directive to your Compose service YAML:
>
> ```yaml
> extra_hosts:
>   - "model-runner.docker.internal:host-gateway"
> ```
> Then you can access the Docker Model Runner APIs at http://model-runner.docker.internal:12434/

{{< /tab >}}
{{</tabs >}}

### Available DMR endpoints

- Create a model:

  ```text
  POST /models/create
  ```

- List models:

  ```text
  GET /models
  ```

- Get a model:

  ```text
  GET /models/{namespace}/{name}
  ```

- Delete a local model:

  ```text
  DELETE /models/{namespace}/{name}
  ```

### Available OpenAI endpoints

DMR supports the following OpenAI endpoints:

- [List models](https://platform.openai.com/docs/api-reference/models/list):

  ```text
  GET /engines/llama.cpp/v1/models
  ```

- [Retrieve model](https://platform.openai.com/docs/api-reference/models/retrieve):

  ```text
  GET /engines/llama.cpp/v1/models/{namespace}/{name}
  ```

- [List chat completions](https://platform.openai.com/docs/api-reference/chat/list):

  ```text
  POST /engines/llama.cpp/v1/chat/completions
  ```

- [Create completions](https://platform.openai.com/docs/api-reference/completions/create):

  ```text
  POST /engines/llama.cpp/v1/completions
  ```


- [Create embeddings](https://platform.openai.com/docs/api-reference/embeddings/create):

  ```text
  POST /engines/llama.cpp/v1/embeddings
  ```

To call these endpoints via a Unix socket (`/var/run/docker.sock`), prefix their path
with `/exp/vDD4.40`.

> [!NOTE]
> You can omit `llama.cpp` from the path. For example: `POST /engines/v1/chat/completions`.

## REST API examples

### Request from within a container

To call the `chat/completions` OpenAI endpoint from within another container using `curl`:

```bash
#!/bin/sh

curl http://model-runner.docker.internal/engines/llama.cpp/v1/chat/completions \
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

  curl http://localhost:12434/engines/llama.cpp/v1/chat/completions \
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
    localhost/exp/vDD4.40/engines/llama.cpp/v1/chat/completions \
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
