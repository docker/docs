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

## Requirements

Docker Model Runner is supported on the following platforms:

{{< tabs >}}
{{< tab name="Windows">}}

Windows(amd64):
-  NVIDIA GPUs
-  NVIDIA drivers 576.57+

Windows(arm64):
- OpenCL for Adreno
- Qualcomm Adreno GPU (6xx series and later)

  > [!NOTE]
  > Some llama.cpp features might not be fully supported on the 6xx series.

{{< /tab >}}
{{< tab name="MacOS">}}

- Apple Silicon

{{< /tab >}}
{{< tab name="Linux">}}

Docker Engine only:

- Linux CPU & Linux NVIDIA
- NVIDIA drivers 575.57.08+

{{< /tab >}}
{{</tabs >}}

## How it works

Models are pulled from Docker Hub the first time you use them and are stored
locally. They load into memory only at runtime when a request is made, and
unload when not in use to optimize resources. Because models can be large, the
initial pull may take some time. After that, they're cached locally for faster
access. You can interact with the model using
[OpenAI-compatible APIs](#what-api-endpoints-are-available).

> [!TIP]
>
> Using Testcontainers or Docker Compose?
> [Testcontainers for Java](https://java.testcontainers.org/modules/docker_model_runner/)
> and [Go](https://golang.testcontainers.org/modules/dockermodelrunner/), and
> [Docker Compose](/manuals/ai/compose/models-and-compose.md) now support Docker
> Model Runner.

## Enable Docker Model Runner

### Enable DMR in Docker Desktop

1. In the settings view, go to the **Beta features** tab.
1. Select the **Enable Docker Model Runner** setting.
1. If you use Windows with a supported NVIDIA GPU, you also see and can select
   **Enable GPU-backed inference**.
1. Optional: To enable TCP support, select **Enable host-side TCP support**.
   1. In the **Port** field, type the port you want to use.
   1. If you interact with Model Runner from a local frontend web app, in
      **CORS Allows Origins**, select the origins that Model Runner should
      accept requests from. An origin is the URL where your web app runs, for
      example `http://localhost:3131`.

You can now use the `docker model` command in the CLI and view and interact
with your local models in the **Models** tab in the Docker Desktop Dashboard.

> [!IMPORTANT]
>
> For Docker Desktop versions 4.41 and earlier, this setting was under the
> **Experimental features** tab on the **Features in development** page.

### Enable DMR in Docker Engine

1. Ensure you have installed [Docker Engine](/engine/install/).
1. DMR is available as a package. To install it, run:

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

1. Test the installation:

   ```console
   $ docker model version
   $ docker model run ai/smollm2
   ```

> [!NOTE]
> TCP support is enabled by default for Docker Engine on port `12434`.

### Update DMR in Docker Engine

To update Docker Model Runner in Docker Engine, uninstall it with
[`docker model uninstall-runner`](/reference/cli/docker/model/uninstall-runner/)
then reinstall it:

```console
docker model uninstall-runner --images && docker model install-runner
```

> [!NOTE]
> With the above command, local models are preserved.
> To delete the models during the upgrade, add the `--models` option to the
> `uninstall-runner` command.

## Pull a model

Models are cached locally.

> [!NOTE]
>
> When you use the Docker CLI, you can also pull models directly from
> [HuggingFace](https://huggingface.co/).

{{< tabs group="release" >}}
{{< tab name="From Docker Desktop">}}

1. Select **Models** and select the **Docker Hub** tab.
1. Find the model you want and select **Pull**.

![Screenshot showing the Docker Hub view.](./images/dmr-catalog.png)

{{< /tab >}}
{{< tab name="From the Docker CLI">}}

Use the [`docker model pull` command](/reference/cli/docker/model/pull/).
For example:

```bash {title="Pulling from Docker Hub"}
docker model pull ai/smollm2:360M-Q4_K_M
```

```bash {title="Pulling from HuggingFace"}
docker model pull hf.co/bartowski/Llama-3.2-1B-Instruct-GGUF
```

{{< /tab >}}
{{< /tabs >}}

## Run a model

{{< tabs group="release" >}}
{{< tab name="From Docker Desktop">}}

1. Select **Models** and select the **Local** tab.
1. Select the play button. The interactive chat screen opens.

![Screenshot showing the Local view.](./images/dmr-run.png)

{{< /tab >}}
{{< tab name="From the Docker CLI" >}}

Use the [`docker model run` command](/reference/cli/docker/model/run/).

{{< /tab >}}
{{< /tabs >}}

## Troubleshooting

To troubleshoot issues, display the logs:

{{< tabs group="release" >}}
{{< tab name="From Docker Desktop">}}

Select **Models** and select the **Logs** tab.

![Screenshot showing the Models view.](./images/dmr-logs.png)

{{< /tab >}}
{{< tab name="From the Docker CLI">}}

Use the [`docker model logs` command](/reference/cli/docker/model/logs/).

{{< /tab >}}
{{< /tabs >}}

## Publish a model

> [!NOTE]
>
> This works for any Container Registry supporting OCI Artifacts, not only
> Docker Hub.

You can tag existing models with a new name and publish them under a different
namespace and repository:

```console
# Tag a pulled model under a new name
$ docker model tag ai/smollm2 myorg/smollm2

# Push it to Docker Hub
$ docker model push myorg/smollm2
```

For more details, see the [`docker model tag`](/reference/cli/docker/model/tag)
and [`docker model push`](/reference/cli/docker/model/push) command
documentation.

You can also package a model file in GGUF format as an OCI Artifact and publish
it to Docker Hub.

```console
# Download a model file in GGUF format, for example from HuggingFace
$ curl -L -o model.gguf https://huggingface.co/TheBloke/Mistral-7B-v0.1-GGUF/resolve/main/mistral-7b-v0.1.Q4_K_M.gguf

# Package it as OCI Artifact and push it to Docker Hub
$ docker model package --gguf "$(pwd)/model.gguf" --push myorg/mistral-7b-v0.1:Q4_K_M
```

For more details, see the
[`docker model package`](/reference/cli/docker/model/package/) command
documentation.

## Example: Integrate Docker Model Runner into your software development lifecycle

### Sample project

You can now start building your generative AI application powered by Docker
Model Runner.

If you want to try an existing GenAI application, follow these steps:

1. Set up the sample app. Clone and run the following repository:

   ```console
   $ git clone https://github.com/docker/hello-genai.git
   ```

1. In your terminal, go to the `hello-genai` directory.

1. Run `run.sh` to pull the chosen model and run the app.

1. Open your app in the browser at the addresses specified in the repository
   [README](https://github.com/docker/hello-genai).

You see the GenAI app's interface where you can start typing your prompts.

You can now interact with your own GenAI app, powered by a local model. Try a
few prompts and notice how fast the responses are — all running on your machine
with Docker.

### Use Model Runner in GitHub Actions

Here is an example of how to use Model Runner as part of a GitHub workflow.
The example installs Model Runner, tests the installation, pulls and runs a
model, interacts with the model via the API, and deletes the model.

```yaml {title="dmr-run.yml", collapse=true}
name: Docker Model Runner Example Workflow

permissions:
  contents: read

on:
  workflow_dispatch:
    inputs:
      test_model:
        description: 'Model to test with (default: ai/smollm2:360M-Q4_K_M)'
        required: false
        type: string
        default: 'ai/smollm2:360M-Q4_K_M'

jobs:
  dmr-test:
    runs-on: ubuntu-latest
    timeout-minutes: 30

    steps:
      - name: Set up Docker
        uses: docker/setup-docker-action@v4

      - name: Install docker-model-plugin
        run: |
          echo "Installing docker-model-plugin..."
          # Add Docker's official GPG key:
          sudo apt-get update
          sudo apt-get install ca-certificates curl
          sudo install -m 0755 -d /etc/apt/keyrings
          sudo curl -fsSL https://download.docker.com/linux/ubuntu/gpg -o /etc/apt/keyrings/docker.asc
          sudo chmod a+r /etc/apt/keyrings/docker.asc
          
          # Add the repository to Apt sources:
          echo \
          "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.asc] https://download.docker.com/linux/ubuntu \
          $(. /etc/os-release && echo "${UBUNTU_CODENAME:-$VERSION_CODENAME}") stable" | \
          sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
          sudo apt-get update
          sudo apt-get install -y docker-model-plugin
          
          echo "Installation completed successfully"

      - name: Test docker model version
        run: |
          echo "Testing docker model version command..."
          sudo docker model version
          
          # Verify the command returns successfully
          if [ $? -eq 0 ]; then
            echo "✅ docker model version command works correctly"
          else
            echo "❌ docker model version command failed"
            exit 1
          fi

      - name: Pull the provided model and run it
        run: |
          MODEL="${{ github.event.inputs.test_model || 'ai/smollm2:360M-Q4_K_M' }}"
          echo "Testing with model: $MODEL"
          
          # Test model pull
          echo "Pulling model..."
          sudo docker model pull "$MODEL"
          
          if [ $? -eq 0 ]; then
            echo "✅ Model pull successful"
          else
            echo "❌ Model pull failed"
            exit 1
          fi
                  
          # Test basic model run (with timeout to avoid hanging)
          echo "Testing docker model run..."
          timeout 60s sudo docker model run "$MODEL" "Give me a fact about whales." || {
            exit_code=$?
            if [ $exit_code -eq 124 ]; then
              echo "✅ Model run test completed (timed out as expected for non-interactive test)"
            else
              echo "❌ Model run failed with exit code: $exit_code"
              exit 1
            fi
          }
               - name: Test model pull and run
        run: |
          MODEL="${{ github.event.inputs.test_model || 'ai/smollm2:360M-Q4_K_M' }}"
          echo "Testing with model: $MODEL"
          
          # Test model pull
          echo "Pulling model..."
          sudo docker model pull "$MODEL"
          
          if [ $? -eq 0 ]; then
            echo "✅ Model pull successful"
          else
            echo "❌ Model pull failed"
            exit 1
          fi
                  
          # Test basic model run (with timeout to avoid hanging)
          echo "Testing docker model run..."
          timeout 60s sudo docker model run "$MODEL" "Give me a fact about whales." || {
            exit_code=$?
            if [ $exit_code -eq 124 ]; then
              echo "✅ Model run test completed (timed out as expected for non-interactive test)"
            else
              echo "❌ Model run failed with exit code: $exit_code"
              exit 1
            fi
          }

      - name: Test API endpoint
        run: |
          MODEL="${{ github.event.inputs.test_model || 'ai/smollm2:360M-Q4_K_M' }}"
          echo "Testing API endpoint with model: $MODEL"
                  
          # Test API call with curl
          echo "Testing API call..."
          RESPONSE=$(curl -s http://localhost:12434/engines/llama.cpp/v1/chat/completions \
            -H "Content-Type: application/json" \
            -d "{
                \"model\": \"$MODEL\",
                \"messages\": [
                    {
                        \"role\": \"user\",
                        \"content\": \"Say hello\"
                    }
                ],
                \"top_k\": 1,
                \"temperature\": 0
            }")
          
          if [ $? -eq 0 ]; then
            echo "✅ API call successful"
            echo "Response received: $RESPONSE"
            
            # Check if response contains "hello" (case-insensitive)
            if echo "$RESPONSE" | grep -qi "hello"; then
              echo "✅ Response contains 'hello' (case-insensitive)"
            else
              echo "❌ Response does not contain 'hello'"
              echo "Full response: $RESPONSE"
              exit 1
            fi
          else
            echo "❌ API call failed"
            exit 1
          fi

      - name: Test model cleanup
        run: |
          MODEL="${{ github.event.inputs.test_model || 'ai/smollm2:360M-Q4_K_M' }}"
          
          echo "Cleaning up test model..."
          sudo docker model rm "$MODEL" || echo "Model removal failed or model not found"
          
          # Verify model was removed
          echo "Verifying model cleanup..."
          sudo docker model ls
          
          echo "✅ Model cleanup completed"

      - name: Report success
        if: success()
        run: |
          echo "🎉 Docker Model Runner daily health check completed successfully!"
          echo "All tests passed:"
          echo "  ✅ docker-model-plugin installation successful"
          echo "  ✅ docker model version command working"
          echo "  ✅ Model pull and run operations successful"
          echo "  ✅ API endpoint operations successful"
          echo "  ✅ Cleanup operations successful"
```

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

### No consistent digest support in Model CLI

The Docker Model CLI currently lacks consistent support for specifying models by image digest. As a temporary workaround, you should refer to models by name instead of digest.

## Share feedback

Thanks for trying out Docker Model Runner. Give feedback or report any bugs you may find through the **Give feedback** link next to the **Enable Docker Model Runner** setting.
