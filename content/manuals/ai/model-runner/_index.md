---
title: Docker Model Runner
linkTitle: Model Runner
params:
  sidebar:
    group: AI
weight: 30
description: Learn how to use Docker Model Runner to manage and run AI models.
keywords: Docker, ai, model runner, docker desktop, docker engine, llm, openai, llama.cpp, vllm, cpu, nvidia, cuda, amd, rocm, vulkan
aliases:
  - /desktop/features/model-runner/
  - /model-runner/
---

{{< summary-bar feature_name="Docker Model Runner" >}}

Docker Model Runner (DMR) makes it easy to manage, run, and
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
- Support for both llama.cpp and vLLM inference engines (vLLM currently supported on Linux x86_64/amd64 with NVIDIA GPUs only)
- Package GGUF and Safetensors files as OCI Artifacts and publish them to any Container Registry
- Run and interact with AI models directly from the command line or from the Docker Desktop GUI
- Manage local models and display logs
- Display prompt and response details
- Conversational context support for multi-turn interactions

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

- Supports CPU, NVIDIA (CUDA), AMD (ROCm), and Vulkan backends
- Requires NVIDIA driver 575.57.08+ when using NVIDIA GPUs

{{< /tab >}}
{{</tabs >}}

## How Docker Model Runner works

Models are pulled from Docker Hub the first time you use them and are stored
locally. They load into memory only at runtime when a request is made, and
unload when not in use to optimize resources. Because models can be large, the
initial pull may take some time. After that, they're cached locally for faster
access. You can interact with the model using
[OpenAI-compatible APIs](api-reference.md).

Docker Model Runner supports both [llama.cpp](https://github.com/ggerganov/llama.cpp) and [vLLM](https://github.com/vllm-project/vllm) as inference engines, providing flexibility for different model formats and performance requirements. For more details, see the [Docker Model Runner repository](https://github.com/docker/model-runner).

> [!TIP]
>
> Using Testcontainers or Docker Compose?
> [Testcontainers for Java](https://java.testcontainers.org/modules/docker_model_runner/)
> and [Go](https://golang.testcontainers.org/modules/dockermodelrunner/), and
> [Docker Compose](/manuals/ai/compose/models-and-compose.md) now support Docker
> Model Runner.

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

## Share feedback

Thanks for trying out Docker Model Runner. To report bugs or request features, [open an issue on GitHub](https://github.com/docker/model-runner/issues). You can also give feedback through the **Give feedback** link next to the **Enable Docker Model Runner** setting.

## Next steps

[Get started with DMR](get-started.md)
