---
title: Configuration options
description: Configure context size, runtime parameters, and model behavior in Docker Model Runner.
weight: 35
keywords: Docker, ai, model runner, configuration, context size, context length, tokens, llama.cpp, parameters
---

Docker Model Runner provides several configuration options to tune model behavior,
memory usage, and inference performance. This guide covers the key settings and
how to apply them.

## Context size (context length)

The context size determines the maximum number of tokens a model can process in
a single request, including both the input prompt and generated output. This is
one of the most important settings affecting memory usage and model capabilities.

### Default context size

By default, Docker Model Runner uses a context size that balances capability with
resource efficiency:

| Engine | Default behavior |
|--------|------------------|
| llama.cpp | 4096 tokens |
| vLLM | Uses the model's maximum trained context size |

> [!NOTE]
> The actual default varies by model. Most models support between 2,048 and 8,192
> tokens by default. Some newer models support 32K, 128K, or even larger contexts.

### Configure context size

You can adjust context size per model using the `docker model configure` command:

```console
$ docker model configure --context-size 8192 ai/qwen2.5-coder
```

Or in a Compose file:

```yaml
models:
  llm:
    model: ai/qwen2.5-coder
    context_size: 8192
```

### Context size guidelines

| Context size | Typical use case | Memory impact |
|--------------|------------------|---------------|
| 2,048 | Simple queries, short code snippets | Low |
| 4,096 | Standard conversations, medium code files | Moderate |
| 8,192 | Long conversations, larger code files | Higher |
| 16,384+ | Extended documents, multi-file context | High |

> [!IMPORTANT]
> Larger context sizes require more memory (RAM/VRAM). If you experience out-of-memory
> errors, reduce the context size. As a rough guide, each additional 1,000 tokens
> requires approximately 100-500 MB of additional memory, depending on the model size.

### Check a model's maximum context

To see a model's configuration including context size:

```console
$ docker model inspect ai/qwen2.5-coder
```

> [!NOTE]
> The `docker model inspect` command shows the model's maximum supported context length
> (e.g., `gemma3.context_length`), not the configured context size. The configured context
> size is what you set with `docker model configure --context-size` and represents the
> actual limit used during inference, which should be less than or equal to the model's
> maximum supported context length.

## Runtime flags

Runtime flags let you pass parameters directly to the underlying inference engine.
This provides fine-grained control over model behavior.

### Using runtime flags

Runtime flags can be provided through multiple mechanisms:

#### Using Docker Compose

In a Compose file:

```yaml
models:
  llm:
    model: ai/qwen2.5-coder
    context_size: 4096
    runtime_flags:
      - "--temp"
      - "0.7"
      - "--top-p"
      - "0.9"
```

#### Using Command Line

With the `docker model configure` command:

```console
$ docker model configure --runtime-flag "--temp" --runtime-flag "0.7" --runtime-flag "--top-p" --runtime-flag "0.9" ai/qwen2.5-coder
```

### Common llama.cpp parameters

These are the most commonly used llama.cpp parameters. You don't need to look up
the llama.cpp documentation for typical use cases.

#### Sampling parameters

| Flag | Description | Default | Range |
|------|-------------|---------|-------|
| `--temp` | Temperature for sampling. Lower = more deterministic, higher = more creative | 0.8 | 0.0-2.0 |
| `--top-k` | Limit sampling to top K tokens. Lower = more focused | 40 | 1-100 |
| `--top-p` | Nucleus sampling threshold. Lower = more focused | 0.9 | 0.0-1.0 |
| `--min-p` | Minimum probability threshold | 0.05 | 0.0-1.0 |
| `--repeat-penalty` | Penalty for repeating tokens | 1.1 | 1.0-2.0 |

**Example: Deterministic output (for code generation)**

```yaml
runtime_flags:
  - "--temp"
  - "0"
  - "--top-k"
  - "1"
```

**Example: Creative output (for storytelling)**

```yaml
runtime_flags:
  - "--temp"
  - "1.2"
  - "--top-p"
  - "0.95"
```

#### Performance parameters

| Flag | Description | Default | Notes |
|------|-------------|---------|-------|
| `--threads` | CPU threads for generation | Auto | Set to number of performance cores |
| `--threads-batch` | CPU threads for batch processing | Auto | Usually same as `--threads` |
| `--batch-size` | Batch size for prompt processing | 512 | Higher = faster prompt processing |
| `--mlock` | Lock model in memory | Off | Prevents swapping, requires sufficient RAM |
| `--no-mmap` | Disable memory mapping | Off | May improve performance on some systems |

**Example: Optimized for multi-core CPU**

```yaml
runtime_flags:
  - "--threads"
  - "8"
  - "--batch-size"
  - "1024"
```

#### GPU parameters

| Flag | Description | Default | Notes |
|------|-------------|---------|-------|
| `--n-gpu-layers` | Layers to offload to GPU | All (if GPU available) | Reduce if running out of VRAM |
| `--main-gpu` | GPU to use for computation | 0 | For multi-GPU systems |
| `--split-mode` | How to split across GPUs | layer | Options: `none`, `layer`, `row` |

**Example: Partial GPU offload (limited VRAM)**

```yaml
runtime_flags:
  - "--n-gpu-layers"
  - "20"
```

#### Advanced parameters

| Flag | Description | Default |
|------|-------------|---------|
| `--rope-scaling` | RoPE scaling method | Auto |
| `--rope-freq-base` | RoPE base frequency | Model default |
| `--rope-freq-scale` | RoPE frequency scale | Model default |
| `--no-prefill-assistant` | Disable assistant pre-fill | Off |
| `--reasoning-budget` | Token budget for reasoning models | 0 (disabled) |

### vLLM parameters

When using the vLLM backend, different parameters are available.

Use `--hf_overrides` to pass HuggingFace model config overrides as JSON:

```console
$ docker model configure --hf_overrides '{"rope_scaling": {"type": "dynamic", "factor": 2.0}}' ai/model-vllm
```

## Configuration presets

Here are complete configuration examples for common use cases.

### Code completion (fast, deterministic)

```yaml
models:
  coder:
    model: ai/qwen2.5-coder
    context_size: 4096
    runtime_flags:
      - "--temp"
      - "0.1"
      - "--top-k"
      - "1"
      - "--batch-size"
      - "1024"
```

### Chat assistant (balanced)

```yaml
models:
  assistant:
    model: ai/llama3.2
    context_size: 8192
    runtime_flags:
      - "--temp"
      - "0.7"
      - "--top-p"
      - "0.9"
      - "--repeat-penalty"
      - "1.1"
```

### Creative writing (high temperature)

```yaml
models:
  writer:
    model: ai/llama3.2
    context_size: 8192
    runtime_flags:
      - "--temp"
      - "1.2"
      - "--top-p"
      - "0.95"
      - "--repeat-penalty"
      - "1.0"
```

### Long document analysis (large context)

```yaml
models:
  analyzer:
    model: ai/qwen2.5-coder:14B
    context_size: 32768
    runtime_flags:
      - "--mlock"
      - "--batch-size"
      - "2048"
```

### Low memory system

```yaml
models:
  efficient:
    model: ai/smollm2:360M-Q4_K_M
    context_size: 2048
    runtime_flags:
      - "--threads"
      - "4"
```

## Environment-based configuration

You can also configure models via environment variables in containers:

| Variable | Description |
|----------|-------------|
| `LLM_URL` | Auto-injected URL of the model endpoint |
| `LLM_MODEL` | Auto-injected model identifier |

See [Models and Compose](/manuals/ai/compose/models-and-compose.md) for details on how these are populated.

## Reset configuration

Configuration set via `docker model configure` persists until the model is removed.
To reset configuration:

```console
$ docker model configure --context-size -1 ai/qwen2.5-coder
```

Using `-1` resets to the default value.

## What's next

- [Inference engines](inference-engines.md) - Learn about llama.cpp and vLLM
- [API reference](api-reference.md) - API parameters for per-request configuration
- [Models and Compose](/manuals/ai/compose/models-and-compose.md) - Configure models in Compose applications
