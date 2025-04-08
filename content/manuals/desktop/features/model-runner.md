---
title: Docker Model Runner
params:
  sidebar:
    badge:
      color: blue
      text: Beta
weight: 20
description: Learn how to use Docker Model Runner to manage and run AI models. 
keywords: Docker, ai, model runner, docker deskotp, llm
---

{{< summary-bar feature_name="Docker Model Runner" >}}

The Docker Model Runner plugin lets you:

- [Pull models from Docker Hub](https://hub.docker.com/u/ai)
- Run AI models directly from the command line
- Manage local models (add, list, remove)
- Interact with models using a submitted prompt or in chat mode

Models are pulled from Docker Hub the first time they're used and stored locally. They're loaded into memory only at runtime when a request is made, and unloaded when not in use to optimize resources. Since models can be large, the initial pull may take some time — but after that, they're cached locally for faster access. You can interact with the model using [OpenAI-compatible APIs](#what-api-endpoints-are-available).

## Enable Docker Model Runner

1. Navigate to the **Features in development** tab in settings.
2. Under the **Experimental features** tab, select **Access experimental features**.
3. Select **Apply and restart**. 
4. Quit and reopen Docker Desktop to ensure the changes take effect. 
5. Open the **Settings** view in Docker Desktop.
6. Navigate to **Features in development**.
7. From the **Beta** tab, check the **Enable Docker Model Runner** setting.

## Available commands

### Model runner status

Check whether the Docker Model Runner is active:

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

Run a model and interact with it using a submitted prompt or in chat mode.

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
docker model run ai/smollm2
```

Output:

```text
Interactive chat mode started. Type '/bye' to exit.
> Hi
Hi there! It's SmolLM, AI assistant. How can I help you today?
> /bye
Chat session ended.
```

### Remove a model

Removes a downloaded model from your system.

```console
$ docker model rm <model>
```

Output:

```text
Model <model> removed successfully
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

### What API endpoints are available?

Once the feature is enabled, the following new APIs are available:

```text
#### Inside containers ####

http://model-runner.docker.internal/

    # Docker Model management
    POST /models/create
    GET /models
    GET /models/{namespace}/{name}
    DELETE /models/{namespace}/{name}

    # OpenAI endpoints
    GET /engines/llama.cpp/v1/models
    GET /engines/llama.cpp/v1/models/{namespace}/{name}
    POST /engines/llama.cpp/v1/chat/completions
    POST /engines/llama.cpp/v1/completions
    POST /engines/llama.cpp/v1/embeddings
    Note: You can also omit llama.cpp.
    E.g., POST /engines/v1/chat/completions.

#### Inside or outside containers (host) ####

Same endpoints on /var/run/docker.sock

    # While still in Beta
    Prefixed with /exp/vDD4.40
```

### How do I interact through the OpenAI API?

#### From within a container

Examples of calling an OpenAI endpoint (`chat/completions`) from within another container using `curl`:

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

#### From the host using a Unix socket

Examples of calling an OpenAI endpoint (`chat/completions`) through the Docker socket from the host using `curl`:

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

#### From the host using TCP

In case you want to interact with the API from the host, but use TCP instead of a Docker socket, you can enable the host-side TCP support from the Docker Desktop GUI, or via the [Docker Desktop CLI](/manuals/desktop/features/desktop-cli.md). For example, using `docker desktop enable model-runner --tcp <port>`.

Afterwards, interact with it as previously documented using `localhost` and the chosen, or the default port.

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

### `model run` drops into chat even if pull fails

If a model image fails to pull successfully, for example due to network issues or lack of disk space, the `docker model run` command will still drop you into the chat interface, even though the model isn’t actually available. This can lead to confusion, as the chat will not function correctly without a running model. 

You can manually retry the `docker model pull` command to ensure the image is available before running it again.

### No consistent digest support in Model CLI

The Docker Model CLI currently lacks consistent support for specifying models by image digest. As a temporary workaround, you should refer to models by name instead of digest.

### Misleading pull progress after failed initial attempt

In some cases, if an initial `docker model pull` fails partway through, a subsequent successful pull may misleadingly report “0 bytes” downloaded even though data is being fetched in the background. This can give the impression that nothing is happening, when in fact the model is being retrieved. Despite the incorrect progress output, the pull typically completes as expected.

## Share feedback

Thanks for trying out Docker Model Runner. Give feedback or report any bugs you may find through the **Give feedback** link next to the **Enable Docker Model Runner** setting. 

## Disable the feature

To disable Docker Model Runner:

1. Open the **Settings** view in Docker Desktop.
2. Navigate to the **Beta** tab in **Features in development**.
3. Clear the **Enable Docker Model Runner** checkbox.
4. Select **Apply & restart**.