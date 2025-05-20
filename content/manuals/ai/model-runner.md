---
title: Docker Model Runner
params:
  sidebar:
    badge:
      color: blue
      text: Beta
    group: AI
weight: 20
description: Learn how to use Docker Model Runner to manage and run AI models.
keywords: Docker, ai, model runner, docker desktop, docker engine, llm
aliases:
  - /desktop/features/model-runner/
  - /ai/model-runner/
---

{{< summary-bar feature_name="Docker Model Runner" >}}

The Docker Model Runner plugin lets you:

- [Pull models from Docker Hub](https://hub.docker.com/u/ai)
- Run AI models directly from the command line
- Manage local models (add, list, remove)
- Interact with models using a submitted prompt or in chat mode in the CLI or Docker Desktop Dashboard
- Push models to Docker Hub

Models are pulled from Docker Hub the first time they're used and stored locally. They're loaded into memory only at runtime when a request is made, and unloaded when not in use to optimize resources. Since models can be large, the initial pull may take some time — but after that, they're cached locally for faster access. You can interact with the model using [OpenAI-compatible APIs](#what-api-endpoints-are-available).

> [!TIP]
>
> Using Testcontainers or Docker Compose? [Testcontainers for Java](https://java.testcontainers.org/modules/docker_model_runner/) and [Go](https://golang.testcontainers.org/modules/dockermodelrunner/), and [Docker Compose](/manuals/compose/how-tos/model-runner.md) now support Docker Model Runner.

## Enable Docker Model Runner

### Enable DMR in Docker Desktop

1. Navigate to the **Features in development** tab in settings.
2. Under the **Experimental features** tab, select **Access experimental features**.
3. Select **Apply and restart**.
4. Quit and reopen Docker Desktop to ensure the changes take effect.
5. Open the **Settings** view in Docker Desktop.
6. Navigate to **Features in development**.
7. From the **Beta** tab, tick the **Enable Docker Model Runner** setting.
8. If you are running on Windows with a supported NVIDIA GPU, you should also see and be able to tick the **Enable GPU-backed inference** setting.

You can now use the `docker model` command in the CLI and view and interact with your local models in the **Models** tab in the Docker Desktop Dashboard.

### Enable DMR in Docker Engine

1. Ensure you have installed [Docker Engine](/engine/install/).
2. DMR is available as a package. To install it, run:

   {{< tabs >}}
   {{< tab name="Ubuntu/Debian">}}

   ```console
   $ sudo apt-get update
   $ sudo apt-get install docker-model-plugin
   ```

   {{< /tab >}}
   {{< tab name="RPM-base distributions">}}

   ```console
   $ sudo dnf update
   $ sudo dnf install docker-model-plugin
   ```

   {{< /tab >}}
   {{< /tabs >}}

3. Test the installation:

   ```console
   $ docker model version
   ```

## Available commands

### Model runner status

Check whether the Docker Model Runner is active and displays the current inference engine:

```console
$ docker model status
```

### View all commands

Displays help information and a list of available subcommands.

```console
$ docker model help
```

Output:

```text
Usage:  docker model COMMAND

Commands:
  list        List models available locally
  pull        Download a model from Docker Hub
  rm          Remove a downloaded model
  run         Run a model interactively or with a prompt
  status      Check if the model runner is running
  version     Show the current version
```

### Pull a model

Pulls a model from Docker Hub to your local environment.

```console
$ docker model pull <model>
```

Example:

```console
$ docker model pull ai/smollm2
```

Output:

```text
Downloaded: 257.71 MB
Model ai/smollm2 pulled successfully
```

The models also display in the Docker Desktop Dashboard.

#### Pull from Hugging Face

You can also pull GGUF models directly from [Hugging Face](https://huggingface.co/models?library=gguf).

```console
$ docker model pull hf.co/<model-you-want-to-pull>
```

For example:

```console
$ docker model pull hf.co/bartowski/Llama-3.2-1B-Instruct-GGUF
```

Pulls the [bartowski/Llama-3.2-1B-Instruct-GGUF](https://huggingface.co/bartowski/Llama-3.2-1B-Instruct-GGUF).

### List available models

Lists all models currently pulled to your local environment.

```console
$ docker model list
```

You will see something similar to:

```text
+MODEL       PARAMETERS  QUANTIZATION    ARCHITECTURE  MODEL ID      CREATED     SIZE
+ai/smollm2  361.82 M    IQ2_XXS/Q4_K_M  llama         354bf30d0aa3  3 days ago  256.35 MiB
```

### Run a model

Run a model and interact with it using a submitted prompt or in chat mode. When you run a model, Docker
calls an Inference Server API endpoint hosted by the Model Runner through Docker Desktop. The model
stays in memory until another model is requested, or until a pre-defined inactivity timeout is reached (currently 5 minutes).

You do not have to use `Docker model run` before interacting with a specific model from a
host process or from within a container. Model Runner transparently loads the requested model on-demand, assuming it has been
pulled beforehand and is locally available.

#### One-time prompt

```console
$ docker model run ai/smollm2 "Hi"
```

Output:

```text
Hello! How can I assist you today?
```

#### Interactive chat

```console
$ docker model run ai/smollm2
```

Output:

```text
Interactive chat mode started. Type '/bye' to exit.
> Hi
Hi there! It's SmolLM, AI assistant. How can I help you today?
> /bye
Chat session ended.
```

> [!TIP]
>
> You can also use chat mode in the Docker Desktop Dashboard when you select the model in the **Models** tab.

### Push a model to Docker Hub

To push your model to Docker Hub:

```console
$ docker model push <namespace>/<model>
```

### Tag a model

To specify a particular version or variant of the model:

```console
$ docker model tag
```

If no tag is provided, Docker defaults to `latest`.

### View the logs

Fetch logs from Docker Model Runner to monitor activity or debug issues.

```console
$ docker model logs
```

The following flags are accepted:

- `-f`/`--follow`: View logs with real-time streaming
- `--no-engines`: Exclude inference engine logs from the output

### Remove a model

Removes a downloaded model from your system.

```console
$ docker model rm <model>
```

Output:

```text
Model <model> removed successfully
```

### Package a model

Packages a GGUF file into a Docker model OCI artifact, with optional licenses, and pushes it to the specified registry.

```console
$ docker model package \
    --gguf ./model.gguf \
    --licenses license1.txt \
    --licenses license2.txt \
    --push registry.example.com/ai/custom-model
```

## Integrate the Docker Model Runner into your software development lifecycle

You can now start building your Generative AI application powered by the Docker Model Runner.

If you want to try an existing GenAI application, follow these instructions.

1. Set up the sample app. Clone and run the following repository:

   ```console
   $ git clone https://github.com/docker/hello-genai.git
   ```

2. In your terminal, navigate to the `hello-genai` directory.

3. Run `run.sh` for pulling the chosen model and run the app(s):

4. Open you app in the browser at the addresses specified in the repository [README](https://github.com/docker/hello-genai).

You'll see the GenAI app's interface where you can start typing your prompts.

You can now interact with your own GenAI app, powered by a local model. Try a few prompts and notice how fast the responses are — all running on your machine with Docker.

## FAQs

### What models are available?

All the available models are hosted in the [public Docker Hub namespace of `ai`](https://hub.docker.com/u/ai).

### What CLI commands are available?

See [the reference docs](/reference/cli/docker/model/).

### What API endpoints are available?

Once the feature is enabled, new API endpoints are available under the following base URLs:

- From containers: `http://model-runner.docker.internal/`
- From host processes: `http://localhost:12434/`, assuming you have enabled TCP host access on default port 12434.

Docker Model management endpoints:

```text
POST /models/create
GET /models
GET /models/{namespace}/{name}
DELETE /models/{namespace}/{name}
```

OpenAI endpoints:

```text
GET /engines/llama.cpp/v1/models
GET /engines/llama.cpp/v1/models/{namespace}/{name}
POST /engines/llama.cpp/v1/chat/completions
POST /engines/llama.cpp/v1/completions
POST /engines/llama.cpp/v1/embeddings
```

To call these endpoints via a Unix socket (`/var/run/docker.sock`), prefix their path with
with `/exp/vDD4.40`.

> [!NOTE]
> You can omit `llama.cpp` from the path. For example: `POST /engines/v1/chat/completions`.


### How do I interact through the OpenAI API?

#### From within a container

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

#### From the host using TCP

To call the `chat/completions` OpenAI endpoint from the host via TCP:

1. Enable the host-side TCP support from the Docker Desktop GUI, or via the [Docker Desktop CLI](/manuals/desktop/features/desktop-cli.md).
   For example: `docker desktop enable model-runner --tcp <port>`.
2. Interact with it as documented in the previous section using `localhost` and the correct port.

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

#### From the host using a Unix socket

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

## Known issues

### `docker model` is not recognised

If you run a Docker Model Runner command and see:

```text
docker: 'model' is not a docker command
```

It means Docker can't find the plugin because it's not in the expected CLI plugins directory.

To fix this, create a symlink so Docker can detect it:

```console
$ ln -s /Applications/Docker.app/Contents/Resources/cli-plugins/docker-model ~/.docker/cli-plugins/docker-model
```

Once linked, re-run the command.

### No safeguard for running oversized models

Currently, Docker Model Runner doesn't include safeguards to prevent you from launching models that exceed their system’s available resources. Attempting to run a model that is too large for the host machine may result in severe slowdowns or render the system temporarily unusable. This issue is particularly common when running LLMs models without sufficient GPU memory or system RAM.

### No consistent digest support in Model CLI

The Docker Model CLI currently lacks consistent support for specifying models by image digest. As a temporary workaround, you should refer to models by name instead of digest.

## Share feedback

Thanks for trying out Docker Model Runner. Give feedback or report any bugs you may find through the **Give feedback** link next to the **Enable Docker Model Runner** setting.

## Disable the feature

To disable Docker Model Runner:

1. Open the **Settings** view in Docker Desktop.
2. Navigate to the **Beta** tab in **Features in development**.
3. Clear the **Enable Docker Model Runner** checkbox.
4. Select **Apply & restart**.
