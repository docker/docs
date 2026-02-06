---
title: Docker Model Runner
linkTitle: Model Runner
params:
  sidebar:
    group: AI
weight: 30
description: Learn how to use Docker Model Runner to manage and run AI models.
keywords: Docker, ai, model runner, docker desktop, docker engine, llm, openai, ollama, llama.cpp, vllm, diffusers, cpu, nvidia, cuda, amd, rocm, vulkan, cline, continue, cursor, image generation, stable diffusion
aliases:
  - /desktop/features/model-runner/
  - /model-runner/
---

{{< summary-bar feature_name="Docker Model Runner" >}}

Docker Model Runner (DMR) makes it easy to manage, run, and
deploy AI models using Docker. Designed for developers,
Docker Model Runner streamlines the process of pulling, running, and serving
large language models (LLMs) and other AI models directly from Docker Hub,
any OCI-compliant registry, or [Hugging Face](https://huggingface.co/).

With seamless integration into Docker Desktop and Docker
Engine, you can serve models via OpenAI and Ollama-compatible APIs, package GGUF files as
OCI Artifacts, and interact with models from both the command line and graphical
interface.

Whether you're building generative AI applications, experimenting with machine
learning workflows, or integrating AI into your software development lifecycle,
Docker Model Runner provides a consistent, secure, and efficient way to work
with AI models locally.

## Key features

- [Pull and push models to and from Docker Hub or any OCI-compliant registry](https://hub.docker.com/u/ai)
- [Pull models from Hugging Face](https://huggingface.co/)
- Serve models on [OpenAI and Ollama-compatible APIs](api-reference.md) for easy integration with existing apps
- Support for [llama.cpp, vLLM, and Diffusers inference engines](inference-engines.md) (vLLM and Diffusers on Linux with NVIDIA GPUs)
- [Generate images from text prompts](inference-engines.md#diffusers) using Stable Diffusion models with the Diffusers backend
- Package GGUF and Safetensors files as OCI Artifacts and publish them to any Container Registry
- Run and interact with AI models directly from the command line or from the Docker Desktop GUI
- [Connect to AI coding tools](ide-integrations.md) like Cline, Continue, Cursor, and Aider
- [Configure context size and model parameters](configuration.md) to tune performance
- [Set up Open WebUI](openwebui-integration.md) for a ChatGPT-like web interface
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

Models are pulled from Docker Hub, an OCI-compliant registry, or
[Hugging Face](https://huggingface.co/) the first time you use them and are
stored locally. They load into memory only at runtime when a request is made,
and unload when not in use to optimize resources. Because models can be large,
the initial pull may take some time. After that, they're cached locally for
faster access. You can interact with the model using
[OpenAI and Ollama-compatible APIs](api-reference.md).

### Inference engines

Docker Model Runner supports three inference engines:

| Engine | Best for | Model format |
|--------|----------|--------------|
| [llama.cpp](inference-engines.md#llamacpp) | Local development, resource efficiency | GGUF (quantized) |
| [vLLM](inference-engines.md#vllm) | Production, high throughput | Safetensors |
| [Diffusers](inference-engines.md#diffusers) | Image generation (Stable Diffusion) | Safetensors |

llama.cpp is the default engine and works on all platforms. vLLM requires NVIDIA GPUs and is supported on Linux x86_64 and Windows with WSL2. Diffusers enables image generation and requires NVIDIA GPUs on Linux (x86_64 or ARM64). See [Inference engines](inference-engines.md) for detailed comparison and setup.

### Context size

Models have a configurable context size (context length) that determines how many tokens they can process. The default varies by model but is typically 2,048-8,192 tokens. You can adjust this per-model:

```console
$ docker model configure --context-size 8192 ai/qwen2.5-coder
```

See [Configuration options](configuration.md) for details on context size and other parameters.

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

## Privacy and data collection

Docker Model Runner respects your privacy settings in Docker Desktop. Data collection is controlled by the **Send usage statistics** setting:

- **Disabled**: No usage data is collected
- **Enabled**: Only minimal, non-personal data is collected:
  - [Model names](https://github.com/docker/model-runner/blob/eb76b5defb1a598396f99001a500a30bbbb48f01/pkg/metrics/metrics.go#L96) (via HEAD requests to Docker Hub)
  - User agent information
  - Whether requests originate from the host or containers

When using Docker Model Runner with Docker Engine, HEAD requests to Docker Hub are made to track model names, regardless of any settings.

No prompt content, responses, or personally identifiable information is ever collected.

## Share feedback

Thanks for trying out Docker Model Runner. To report bugs or request features, [open an issue on GitHub](https://github.com/docker/model-runner/issues). You can also give feedback through the **Give feedback** link next to the **Enable Docker Model Runner** setting.

## Next steps

- [Get started with DMR](get-started.md) - Enable DMR and run your first model
- [API reference](api-reference.md) - OpenAI and Ollama-compatible API documentation
- [Configuration options](configuration.md) - Context size and runtime parameters
- [Inference engines](inference-engines.md) - llama.cpp, vLLM, and Diffusers details
- [IDE integrations](ide-integrations.md) - Connect Cline, Continue, Cursor, and more
- [Open WebUI integration](openwebui-integration.md) - Set up a web chat interface
