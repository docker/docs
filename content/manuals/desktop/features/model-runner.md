---
title: Docker Model Runner
params:
  sidebar:
    badge:
      color: blue
      text: Beta
weight: 20
description: Learn how to use Docker Model Runner to manage and run AI models. 
keywords: Docker, ai, model runner, docker deskotp,
---

{{< summary-bar feature_name="Docker Model Runner" >}}

The Docker Model Runner plugin lets you:

- Pull models from Docker Hub
- Run AI models directly from the command line
- Manage local models (add, list, remove)
- Interact with models using a submitted prompt or in chat mode

Models are pulled from Docker Hub the first time they're used and stored locally. They're loaded into memory only at runtime when a request is made, and unloaded when not in use to optimize resources. Since models can be large, the initial pull may take some time — but after that, they're cached locally for faster access. You can interact with the model using [OpenAI-compatible APIs](#what-api-endpoints-are-available).

## Enable the feature

To enable Docker Model Runner:

1. Open the **Settings** view in Docker Desktop.
2. Navigate to the **Beta** tab in **Features in development**.
3. Check the **Enable Docker Model Runner** checkbox.
4. Select **Apply & restart**.

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
$ docker model pull ai/llama3.2:1b
```

Output:

```text
Downloaded: 626.05 MB
Model ai/llama3.2:1b pulled successfully
```

### List available models

Lists all models currently pulled to your local environment.

```console
$ docker model list
```

You will something similar to:

```text
MODEL                                     PARAMETERS  QUANTIZATION    ARCHITECTURE  MODEL ID      CREATED       SIZE
ignaciolopezluna020/gemma-3-it:4B-Q4_K_M  3.88 B      IQ2_XXS/Q4_K_M  gemma3        adea14bef2fe  55 years ago  2.31 GiB
```

### Run a model

Run a model and interact with it using a submitted prompt or in chat mode.

#### One-time prompt

```console
$ docker model run ai/llama3.2:1b "Hi"
```

Output:

```text
Hi! How can I assist you today
```

#### Interactive chat

```console
docker model run ai/llama3.2:1b
```

Output:

```text
Interactive chat mode started. Type '/bye' to exit.
> Hi
Hi! How are you doing today?
> /bye
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

1. Pull the required model from Docker Hub so it's ready for use in your app.

   ```console
   $ docker model pull ai/llama3.2:1b
   ```

2. Set up the sample app. Download and unzip the following folder:
   
   [myapp.zip](attachment:abc104c4-e0c9-4163-b90b-e1f06caab687:myapp.zip)

3. In your terminal, navigate to the `myapp` folder.
4. Start the app with Docker Compose:

   ```console
   $ docker compose up 
   ```

5. Open you app in the browser at `http://localhost:3000`. 

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

    # OpenAI endpoints (per-backend)
    GET /engines/{backend}/v1/models
    GET /engines/{backend}/v1/models/{namespace}/{name}
    POST /engines/{backend}/v1/chat/completions
    POST /engines/{backend}/v1/completions
    POST /engines/{backend}/v1/embeddings
    Note: You can also omit {backend} and it will default to llama.cpp
    E.g., POST /engines/v1/chat/completions.

#### Inside or outside containers (host) ####

Same endpoints on /var/run/docker.sock

    # Until stable...
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
        "model": "ai/llama3.2:1b",
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
        "model": "ai/llama3.2:1b",
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

In case you want to interact with the API from the host, but use TCP instead of a Docker socket, it is recommended you use a helper container as a reverse-proxy. For example, in order to forward the API to `8080`:

```bash
docker run -d --name model-runner-proxy -p 8080:80 alpine/socat tcp-listen:80,fork,reuseaddr tcp:model-runner.docker.internal:80
```

Afterwards, interact with it as previously documented using `localhost` and the forward port, in this case `8080`:

```bash
#!/bin/sh

	curl http://localhost:8080/engines/llama.cpp/v1/chat/completions \
    -H "Content-Type: application/json" \
    -d '{
        "model": "ai/llama3.2:1b",
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

## Share feedback

Thanks for trying out Docker Model Runner. Give feedback or report any bugs you may find through the **Give feedback** link next to the **Enable Docker Model Runner** setting. 
