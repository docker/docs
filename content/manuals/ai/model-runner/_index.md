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
  - /model-runner/
---

{{< summary-bar feature_name="Docker Model Runner" >}}

Docker Model Runner makes it easy to manage, run, and
deploy AI models using Docker. Designed for developers,
Docker Model Runner streamlines the process of pulling, running, and serving
large language models (LLMs) and other AI models directly from Docker Hub or any
OCI-compliant registry.

With seamless integration into Docker Desktop and Docker
Engine, you can serve models via OpenAI-compatible APIs, package GGUF files as
OCI Artifacts, and interact with models from both the command line and graphical
interface.

Whether you're building generative AI applications, experimenting with machine
learning workflows, or integrating AI into your software development lifecycle,
Docker Model Runner provides a consistent, secure, and efficient way to work
with AI models locally.

## Key features

- [Pull and push models to and from Docker Hub](https://hub.docker.com/u/ai)
- Serve models on OpenAI-compatible APIs for easy integration with existing apps
- Package GGUF files as OCI Artifacts and publish them to any Container Registry
- Run and interact with AI models directly from the command line or from the Docker Desktop GUI
- Manage local models and display logs

## How it works

Models are pulled from Docker Hub the first time they're used and stored locally. They're loaded into memory only at runtime when a request is made, and unloaded when not in use to optimize resources. Since models can be large, the initial pull may take some time — but after that, they're cached locally for faster access. You can interact with the model using [OpenAI-compatible APIs](#what-api-endpoints-are-available).

> [!TIP]
>
> Using Testcontainers or Docker Compose?
> [Testcontainers for Java](https://java.testcontainers.org/modules/docker_model_runner/)
> and [Go](https://golang.testcontainers.org/modules/dockermodelrunner/), and
> [Docker Compose](/manuals/compose/how-tos/model-runner.md) now support Docker Model Runner.

## Enable Docker Model Runner

### Enable DMR in Docker Desktop

1. Navigate to the **Beta features** tab in settings.
2. Tick the **Enable Docker Model Runner** setting.
3. If you are running on Windows with a supported NVIDIA GPU, you should also see and be able to tick the **Enable GPU-backed inference** setting.

You can now use the `docker model` command in the CLI and view and interact with your local models in the **Models** tab in the Docker Desktop Dashboard.

> [!IMPORTANT]
>
> For Docker Desktop versions 4.41 and earlier, this settings lived under the **Experimental features** tab on the **Features in development** page.

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
   $ docker model run ai/smollm2
   ```

## Pull a model

Models are cached locally.

{{< tabs group="release" >}}
{{< tab name="From Docker Desktop">}}

1. Select **Models** and select the **Docker Hub** tab.
2. Find the model of your choice and select **Pull**.

{{< /tab >}}
{{< tab name="From the Docker CLI">}}

Use the [`docker model pull` command](/reference/cli/docker/model/pull/).

{{< /tab >}}
{{< /tabs >}}

## Run a model

{{< tabs group="release" >}}
{{< tab name="From Docker Desktop">}}

Select **Models** and select the **Local** tab and click the play button.
The interactive chat screen opens.

{{< /tab >}}
{{< tab name="From the Docker CLI" >}}

Use the [`docker model run` command](/reference/cli/docker/model/run/).

{{< /tab >}}
{{< /tabs >}}

## Troubleshooting

To troubleshoot potential issues, display the logs:

{{< tabs group="release" >}}
{{< tab name="From Docker Desktop">}}

Select **Models** and select the **Logs** tab.

{{< /tab >}}
{{< tab name="From the Docker CLI">}}

Use the [`docker model logs` command](/reference/cli/docker/model/logs/).

{{< /tab >}}
{{< /tabs >}}

## Publish a model

> [!NOTE]
>
> This works for any Container Registry supporting OCI Artifacts, not only Docker Hub.

You can tag existing models with a new name and publish them under a different namespace and repository:

```console
# Tag a pulled model under a new name
$ docker model tag ai/smollm2 myorg/smollm2

# Push it to Docker Hub
$ docker model push myorg/smollm2
```

For more details, see the [`docker model tag`](/reference/cli/docker/model/tag) and [`docker model push`](/reference/cli/docker/model/push) command documentation.

You can also directly package a model file in GGUF format as an OCI Artifact and publish it to Docker Hub.

```console
# Download a model file in GGUF format, e.g. from HuggingFace
$ curl -L -o model.gguf https://huggingface.co/TheBloke/Mistral-7B-v0.1-GGUF/resolve/main/mistral-7b-v0.1.Q4_K_M.gguf

# Package it as OCI Artifact and push it to Docker Hub
$ docker model package --gguf "$(pwd)/model.gguf" --push myorg/mistral-7b-v0.1:Q4_K_M
```

For more details, see the [`docker model package`](/reference/cli/docker/model/package/) command documentation.

## Example: Integrate Docker Model Runner into your software development lifecycle

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

   If you are running on Windows, also enable GPU-backed inference.
   See [Enable Docker Model Runner](#enable-dmr-in-docker-desktop).

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

Once linked, rerun the command.

### No safeguard for running oversized models

Currently, Docker Model Runner doesn't include safeguards to prevent you from
launching models that exceed your system's available resources. Attempting to
run a model that is too large for the host machine may result in severe
slowdowns or may render the system temporarily unusable. This issue is
particularly common when running LLMs without sufficient GPU memory or system
RAM.

### No consistent digest support in Model CLI

The Docker Model CLI currently lacks consistent support for specifying models by image digest. As a temporary workaround, you should refer to models by name instead of digest.

## Share feedback

Thanks for trying out Docker Model Runner. Give feedback or report any bugs you may find through the **Give feedback** link next to the **Enable Docker Model Runner** setting.
