---
title: Inference engines
description: Learn about the llama.cpp and vLLM inference engines in Docker Model Runner.
weight: 50
keywords: Docker, ai, model runner, llama.cpp, vllm, inference, gguf, safetensors, cuda, gpu
---

Docker Model Runner supports two inference engines: **llama.cpp** and **vLLM**.
Each engine has different strengths, supported platforms, and model format
requirements. This guide helps you choose the right engine and configure it for
your use case.

## Engine comparison

| Feature | llama.cpp | vLLM |
|---------|-----------|------|
| **Model formats** | GGUF | Safetensors, HuggingFace |
| **Platforms** | All (macOS, Windows, Linux) | Linux x86_64 only |
| **GPU support** | NVIDIA, AMD, Apple Silicon, Vulkan | NVIDIA CUDA only |
| **CPU inference** | Yes | No |
| **Quantization** | Built-in (Q4, Q5, Q8, etc.) | Limited |
| **Memory efficiency** | High (with quantization) | Moderate |
| **Throughput** | Good | High (with batching) |
| **Best for** | Local development, resource-constrained environments | Production, high throughput |

## llama.cpp

[llama.cpp](https://github.com/ggerganov/llama.cpp) is the default inference
engine in Docker Model Runner. It's designed for efficient local inference and
supports a wide range of hardware configurations.

### Platform support

| Platform | GPU support | Notes |
|----------|-------------|-------|
| macOS (Apple Silicon) | Metal | Automatic GPU acceleration |
| Windows (x64) | NVIDIA CUDA | Requires NVIDIA drivers 576.57+ |
| Windows (ARM64) | Adreno OpenCL | Qualcomm 6xx series and later |
| Linux (x64) | NVIDIA, AMD, Vulkan | Multiple backend options |
| Linux | CPU only | Works on any x64/ARM64 system |

### Model format: GGUF

llama.cpp uses the GGUF format, which supports efficient quantization for reduced
memory usage without significant quality loss.

#### Quantization levels

| Quantization | Bits per weight | Memory usage | Quality |
|--------------|-----------------|--------------|---------|
| Q2_K | ~2.5 | Lowest | Reduced |
| Q3_K_M | ~3.5 | Minimal | Acceptable |
| Q4_K_M | ~4.5 | Low | Good |
| Q5_K_M | ~5.5 | Moderate | Excellent |
| Q6_K | ~6.5 | Higher | Excellent |
| Q8_0 | 8 | High | Near-original |
| F16 | 16 | Highest | Original |

**Recommended**: Q4_K_M offers the best balance of quality and memory usage for
most use cases.

#### Pulling quantized models

Models on Docker Hub often include quantization in the tag:

```console
$ docker model pull ai/llama3.2:3B-Q4_K_M
```

### Using llama.cpp

llama.cpp is the default engine. No special configuration is required:

```console
$ docker model run ai/smollm2
```

To explicitly specify llama.cpp when running models:

```console
$ docker model run ai/smollm2 --backend llama.cpp
```

### llama.cpp API endpoints

When using llama.cpp, API calls use the llama.cpp engine path:

```text
POST /engines/llama.cpp/v1/chat/completions
```

Or without the engine prefix:

```text
POST /engines/v1/chat/completions
```

## vLLM

[vLLM](https://github.com/vllm-project/vllm) is a high-performance inference
engine optimized for production workloads with high throughput requirements.

### Platform support

| Platform | GPU | Support status |
|----------|-----|----------------|
| Linux x86_64 | NVIDIA CUDA | Supported |
| Windows with WSL2 | NVIDIA CUDA | Supported (Docker Desktop 4.54+) |
| macOS | - | Not supported |
| Linux ARM64 | - | Not supported |
| AMD GPUs | - | Not supported |

> [!IMPORTANT]
> vLLM requires an NVIDIA GPU with CUDA support. It does not support CPU-only
> inference.

### Model format: Safetensors

vLLM works with models in Safetensors format, which is the standard format for
HuggingFace models. These models typically use more memory than quantized GGUF
models but may offer better quality and faster inference on powerful hardware.

### Setting up vLLM

#### Docker Engine (Linux)

Install the Model Runner with vLLM backend:

```console
$ docker model install-runner --backend vllm --gpu cuda
```

Verify the installation:

```console
$ docker model status
Docker Model Runner is running

Status:
llama.cpp: running llama.cpp version: c22473b
vllm: running vllm version: 0.11.0
```

#### Docker Desktop (Windows with WSL2)

1. Ensure you have:
   - Docker Desktop 4.54 or later
   - NVIDIA GPU with updated drivers
   - WSL2 enabled

2. Install vLLM backend:
   ```console
   $ docker model install-runner --backend vllm --gpu cuda
   ```

### Running models with vLLM

vLLM models are typically tagged with `-vllm` suffix:

```console
$ docker model run ai/smollm2-vllm
```

To specify the vLLM backend explicitly:

```console
$ docker model run ai/model --backend vllm
```

### vLLM API endpoints

When using vLLM, specify the engine in the API path:

```text
POST /engines/vllm/v1/chat/completions
```

### vLLM configuration

#### HuggingFace overrides

Use `--hf_overrides` to pass model configuration overrides:

```console
$ docker model configure --hf_overrides '{"max_model_len": 8192}' ai/model-vllm
```

#### Common vLLM settings

| Setting | Description | Example |
|---------|-------------|---------|
| `max_model_len` | Maximum context length | 8192 |
| `gpu_memory_utilization` | Fraction of GPU memory to use | 0.9 |
| `tensor_parallel_size` | GPUs for tensor parallelism | 2 |

### vLLM and llama.cpp performance comparison

| Scenario | Recommended engine |
|----------|-------------------|
| Single user, local development | llama.cpp |
| Multiple concurrent requests | vLLM |
| Limited GPU memory | llama.cpp (with quantization) |
| Maximum throughput | vLLM |
| CPU-only system | llama.cpp |
| Apple Silicon Mac | llama.cpp |
| Production deployment | vLLM (if hardware supports it) |

## Running both engines

You can run both llama.cpp and vLLM simultaneously. Docker Model Runner routes
requests to the appropriate engine based on the model or explicit engine selection.

Check which engines are running:

```console
$ docker model status
Docker Model Runner is running

Status:
llama.cpp: running llama.cpp version: c22473b
vllm: running vllm version: 0.11.0
```

### Engine-specific API paths

| Engine | API path |
|--------|----------|
| llama.cpp | `/engines/llama.cpp/v1/...` |
| vLLM | `/engines/vllm/v1/...` |
| Auto-select | `/engines/v1/...` |

## Managing inference engines

### Install an engine

```console
$ docker model install-runner --backend <engine> [--gpu <type>]
```

Options:
- `--backend`: `llama.cpp` or `vllm`
- `--gpu`: `cuda`, `rocm`, `vulkan`, or `metal` (depends on platform)

### Reinstall an engine

```console
$ docker model reinstall-runner --backend <engine>
```

### Check engine status

```console
$ docker model status
```

### View engine logs

```console
$ docker model logs
```

## Packaging models for each engine

### Package a GGUF model (llama.cpp)

```console
$ docker model package --gguf ./model.gguf --push myorg/mymodel:Q4_K_M
```

### Package a Safetensors model (vLLM)

```console
$ docker model package --safetensors ./model/ --push myorg/mymodel-vllm
```

## Troubleshooting

### vLLM won't start

1. Verify NVIDIA GPU is available:
   ```console
   $ nvidia-smi
   ```

2. Check Docker has GPU access:
   ```console
   $ docker run --rm --gpus all nvidia/cuda:12.0-base nvidia-smi
   ```

3. Verify you're on a supported platform (Linux x86_64 or Windows WSL2).

### llama.cpp is slow

1. Ensure GPU acceleration is working (check logs for Metal/CUDA messages).

2. Try a more aggressive quantization:
   ```console
   $ docker model pull ai/model:Q4_K_M
   ```

3. Reduce context size:
   ```console
   $ docker model configure --context-size 2048 ai/model
   ```

### Out of memory errors

1. Use a smaller quantization (Q4 instead of Q8).
2. Reduce context size.
3. For vLLM, adjust `gpu_memory_utilization`:
   ```console
   $ docker model configure --hf_overrides '{"gpu_memory_utilization": 0.8}' ai/model
   ```

## What's next

- [Configuration options](configuration.md) - Detailed parameter reference
- [API reference](api-reference.md) - API documentation
- [GPU support](/manuals/desktop/features/gpu.md) - GPU configuration for Docker Desktop
