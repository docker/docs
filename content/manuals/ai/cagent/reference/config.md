---
title: Configuration file reference
linkTitle: Configuration file
description: Complete reference for the cagent YAML configuration file format
keywords: [ai, agent, cagent, configuration, yaml]
weight: 10
---

This reference documents the YAML configuration file format for cagent agents.
It covers file structure, agent parameters, model configuration, toolset setup,
and RAG sources.

For detailed documentation of each toolset's capabilities and specific options,
see the [Toolsets reference](./toolsets.md).

## File structure

A configuration file has four top-level sections:

```yaml
agents: # Required - agent definitions
  root:
    model: anthropic/claude-sonnet-4-5
    description: What this agent does
    instruction: How it should behave

models: # Optional - model configurations
  custom_model:
    provider: openai
    model: gpt-5

rag: # Optional - RAG sources
  docs:
    docs: [./documents]
    strategies: [...]

metadata: # Optional - author, license, readme
  author: Your Name
```

## Agents

| Property               | Type    | Description                                    | Required |
| ---------------------- | ------- | ---------------------------------------------- | -------- |
| `model`                | string  | Model reference or name                        | Yes      |
| `description`          | string  | Brief description of agent's purpose           | No       |
| `instruction`          | string  | Detailed behavior instructions                 | Yes      |
| `sub_agents`           | array   | Agent names for task delegation                | No       |
| `handoffs`             | array   | Agent names for conversation handoff           | No       |
| `toolsets`             | array   | Available tools                                | No       |
| `welcome_message`      | string  | Message displayed on start                     | No       |
| `add_date`             | boolean | Include current date in context                | No       |
| `add_environment_info` | boolean | Include working directory, OS, Git info        | No       |
| `add_prompt_files`     | array   | Prompt file paths to include                   | No       |
| `max_iterations`       | integer | Maximum tool call loops (unlimited if not set) | No       |
| `num_history_items`    | integer | Conversation history limit                     | No       |
| `code_mode_tools`      | boolean | Enable Code Mode for tools                     | No       |
| `commands`             | object  | Named prompts accessible via `/command_name`   | No       |
| `structured_output`    | object  | JSON schema for structured responses           | No       |
| `rag`                  | array   | RAG source names                               | No       |

### Task delegation versus conversation handoff

Use `sub_agents` to break work into tasks. The root agent assigns work to a
sub-agent and gets results back while staying in control.

Use `handoffs` to transfer the entire conversation to a different agent. The new
agent takes over completely.

### Commands

Named prompts users invoke with `/command_name`. Supports JavaScript template
literals with `${env.VARIABLE}` for environment variables:

```yaml
commands:
  greet: "Say hello to ${env.USER}"
  analyze: "Analyze ${env.PROJECT_NAME || 'demo'}"
```

Run with: `cagent run config.yaml /greet`

### Structured output

Constrain responses to a JSON schema (OpenAI and Gemini only):

```yaml
structured_output:
  name: code_analysis
  strict: true
  schema:
    type: object
    properties:
      issues:
        type: array
        items: { ... }
    required: [issues]
```

## Models

| Property              | Type    | Description                                    | Required |
| --------------------- | ------- | ---------------------------------------------- | -------- |
| `provider`            | string  | `openai`, `anthropic`, `google`, `dmr`         | Yes      |
| `model`               | string  | Model name                                     | Yes      |
| `temperature`         | float   | Randomness (0.0-2.0)                           | No       |
| `max_tokens`          | integer | Maximum response length                        | No       |
| `top_p`               | float   | Nucleus sampling (0.0-1.0)                     | No       |
| `frequency_penalty`   | float   | Repetition penalty (-2.0 to 2.0, OpenAI only)  | No       |
| `presence_penalty`    | float   | Topic penalty (-2.0 to 2.0, OpenAI only)       | No       |
| `base_url`            | string  | Custom API endpoint                            | No       |
| `parallel_tool_calls` | boolean | Enable parallel tool execution (default: true) | No       |
| `token_key`           | string  | Authentication token key                       | No       |
| `track_usage`         | boolean | Track token usage                              | No       |
| `thinking_budget`     | mixed   | Reasoning effort (provider-specific)           | No       |
| `provider_opts`       | object  | Provider-specific options                      | No       |

### Alloy models

Use multiple models in rotation by separating names with commas:

```yaml
model: anthropic/claude-sonnet-4-5,openai/gpt-5
```

### Thinking budget

Controls reasoning depth. Configuration varies by provider:

- **OpenAI**: String values - `minimal`, `low`, `medium`, `high`
- **Anthropic**: Integer token budget (1024-32768, must be less than
  `max_tokens`)
  - Set `provider_opts.interleaved_thinking: true` for tool use during reasoning
- **Gemini**: Integer token budget (0 to disable, -1 for dynamic, max 24576)
  - Gemini 2.5 Pro: 128-32768, cannot disable (minimum 128)

```yaml
# OpenAI
thinking_budget: low

# Anthropic
thinking_budget: 8192
provider_opts:
  interleaved_thinking: true

# Gemini
thinking_budget: 8192    # Fixed
thinking_budget: -1      # Dynamic
thinking_budget: 0       # Disabled
```

### Docker Model Runner (DMR)

Run local models. If `base_url` is omitted, cagent auto-discovers via Docker
Model plugin.

```yaml
provider: dmr
model: ai/qwen3
max_tokens: 8192
base_url: http://localhost:12434/engines/llama.cpp/v1 # Optional
```

Pass llama.cpp options via `provider_opts.runtime_flags` (array, string, or
multiline):

```yaml
provider_opts:
  runtime_flags: ["--ngl=33", "--threads=8"]
  # or: runtime_flags: "--ngl=33 --threads=8"
```

Model config fields auto-map to runtime flags:

- `temperature` → `--temp`
- `top_p` → `--top-p`
- `max_tokens` → `--context-size`

Explicit `runtime_flags` override auto-mapped flags.

Speculative decoding for faster inference:

```yaml
provider_opts:
  speculative_draft_model: ai/qwen3:0.6B-F16
  speculative_num_tokens: 16
  speculative_acceptance_rate: 0.8
```

## Tools

Configure tools in the `toolsets` array. Three types: built-in, MCP
(local/remote), and Docker Gateway.

> [!NOTE] This section covers toolset configuration syntax. For detailed
> documentation of each toolset's capabilities, available tools, and specific
> configuration options, see the [Toolsets reference](./toolsets.md).

All toolsets support common properties like `tools` (whitelist), `defer`
(deferred loading), `toon` (output compression), `env` (environment variables),
and `instruction` (usage guidance). See the [Toolsets reference](./toolsets.md)
for details on these properties and what each toolset does.

### Built-in tools

```yaml
toolsets:
  - type: filesystem
  - type: shell
  - type: think
  - type: todo
    shared: true
  - type: memory
    path: ./memory.db
```

### MCP tools

Local process:

```yaml
- type: mcp
  command: npx
  args:
    ["-y", "@modelcontextprotocol/server-filesystem", "/path/to/allowed/files"]
  tools: ["read_file", "write_file"] # Optional: limit to specific tools
  env:
    NODE_OPTIONS: "--max-old-space-size=8192"
```

Remote server:

```yaml
- type: mcp
  remote:
    url: https://mcp-server.example.com
    transport_type: sse
    headers:
      Authorization: Bearer token
```

### Docker MCP Gateway

Containerized tools from [Docker MCP
Catalog](/manuals/ai/mcp-catalog-and-toolkit/mcp-gateway.md):

```yaml
- type: mcp
  ref: docker:duckduckgo
```

## RAG

Retrieval-augmented generation for document knowledge bases. Define sources at
the top level, reference in agents.

```yaml
rag:
  docs:
    docs: [./documents, ./README.md]
    strategies:
      - type: chunked-embeddings
        embedding_model: openai/text-embedding-3-small
        vector_dimensions: 1536
        database: ./embeddings.db

agents:
  root:
    rag: [docs]
```

### Retrieval strategies

All strategies support chunking configuration. Chunk size and overlap are
measured in characters (Unicode code points), not tokens.

#### Chunked-embeddings

Direct semantic search using vector embeddings. Best for understanding intent,
synonyms, and paraphrasing.

| Field                              | Type    | Default |
| ---------------------------------- | ------- | ------- |
| `embedding_model`                  | string  | -       |
| `database`                         | string  | -       |
| `vector_dimensions`                | integer | -       |
| `similarity_metric`                | string  | cosine  |
| `threshold`                        | float   | 0.5     |
| `limit`                            | integer | 5       |
| `chunking.size`                    | integer | 1000    |
| `chunking.overlap`                 | integer | 75      |
| `chunking.respect_word_boundaries` | boolean | true    |
| `chunking.code_aware`              | boolean | false   |

```yaml
- type: chunked-embeddings
  embedding_model: openai/text-embedding-3-small
  vector_dimensions: 1536
  database: ./vector.db
  similarity_metric: cosine_similarity
  threshold: 0.5
  limit: 10
  chunking:
    size: 1000
    overlap: 100
```

#### Semantic-embeddings

LLM-enhanced semantic search. Uses a language model to generate rich semantic
summaries of each chunk before embedding, capturing deeper meaning.

| Field                              | Type    | Default |
| ---------------------------------- | ------- | ------- |
| `embedding_model`                  | string  | -       |
| `chat_model`                       | string  | -       |
| `database`                         | string  | -       |
| `vector_dimensions`                | integer | -       |
| `similarity_metric`                | string  | cosine  |
| `threshold`                        | float   | 0.5     |
| `limit`                            | integer | 5       |
| `ast_context`                      | boolean | false   |
| `semantic_prompt`                  | string  | -       |
| `chunking.size`                    | integer | 1000    |
| `chunking.overlap`                 | integer | 75      |
| `chunking.respect_word_boundaries` | boolean | true    |
| `chunking.code_aware`              | boolean | false   |

```yaml
- type: semantic-embeddings
  embedding_model: openai/text-embedding-3-small
  vector_dimensions: 1536
  chat_model: openai/gpt-5-mini
  database: ./semantic.db
  threshold: 0.3
  limit: 10
  chunking:
    size: 1000
    overlap: 100
```

#### BM25

Keyword-based search using BM25 algorithm. Best for exact terms, technical
jargon, and code identifiers.

| Field                              | Type    | Default |
| ---------------------------------- | ------- | ------- |
| `database`                         | string  | -       |
| `k1`                               | float   | 1.5     |
| `b`                                | float   | 0.75    |
| `threshold`                        | float   | 0.0     |
| `limit`                            | integer | 5       |
| `chunking.size`                    | integer | 1000    |
| `chunking.overlap`                 | integer | 75      |
| `chunking.respect_word_boundaries` | boolean | true    |
| `chunking.code_aware`              | boolean | false   |

```yaml
- type: bm25
  database: ./bm25.db
  k1: 1.5
  b: 0.75
  threshold: 0.3
  limit: 10
  chunking:
    size: 1000
    overlap: 100
```

### Hybrid retrieval

Combine multiple strategies with fusion:

```yaml
strategies:
  - type: chunked-embeddings
    embedding_model: openai/text-embedding-3-small
    vector_dimensions: 1536
    database: ./vector.db
    limit: 20
  - type: bm25
    database: ./bm25.db
    limit: 15

results:
  fusion:
    strategy: rrf # Options: rrf, weighted, max
    k: 60 # RRF smoothing parameter
  deduplicate: true
  limit: 5
```

Fusion strategies:

- `rrf`: Reciprocal Rank Fusion (recommended, rank-based, no normalization
  needed)
- `weighted`: Weighted combination (`fusion.weights: {chunked-embeddings: 0.7,
bm25: 0.3}`)
- `max`: Maximum score across strategies

### Reranking

Re-score results with a specialized model for improved relevance:

```yaml
results:
  reranking:
    model: openai/gpt-5-mini
    top_k: 10 # Only rerank top K (0 = all)
    threshold: 0.3 # Minimum score after reranking
    criteria: | # Optional domain-specific guidance
      Prioritize official docs over blog posts
  limit: 5
```

DMR native reranking:

```yaml
models:
  reranker:
    provider: dmr
    model: hf.co/ggml-org/qwen3-reranker-0.6b-q8_0-gguf

results:
  reranking:
    model: reranker
```

### Code-aware chunking

For source code, use AST-based chunking. With semantic-embeddings, you can
include AST metadata in the LLM prompts:

```yaml
- type: semantic-embeddings
  embedding_model: openai/text-embedding-3-small
  vector_dimensions: 1536
  chat_model: openai/gpt-5-mini
  database: ./code.db
  ast_context: true # Include AST metadata in semantic prompts
  chunking:
    size: 2000
    code_aware: true # Enable AST-based chunking
```

### RAG properties

Top-level RAG source:

| Field        | Type     | Description                                                     |
| ------------ | -------- | --------------------------------------------------------------- |
| `docs`       | []string | Document paths (suppports glob patterns, respects `.gitignore`) |
| `tool`       | object   | Customize RAG tool name/description/instruction                 |
| `strategies` | []object | Retrieval strategies (see above for strategy-specific fields)   |
| `results`    | object   | Post-processing (fusion, reranking, limits)                     |

Results:

| Field                 | Type    | Default |
| --------------------- | ------- | ------- |
| `limit`               | integer | 15      |
| `deduplicate`         | boolean | true    |
| `include_score`       | boolean | false   |
| `fusion.strategy`     | string  | -       |
| `fusion.k`            | integer | 60      |
| `fusion.weights`      | object  | -       |
| `reranking.model`     | string  | -       |
| `reranking.top_k`     | integer | 0       |
| `reranking.threshold` | float   | 0.5     |
| `reranking.criteria`  | string  | ""      |
| `return_full_content` | boolean | false   |

## Metadata

Documentation and sharing information:

| Property  | Type   | Description                     |
| --------- | ------ | ------------------------------- |
| `author`  | string | Author name                     |
| `license` | string | License (e.g., MIT, Apache-2.0) |
| `readme`  | string | Usage documentation             |

```yaml
metadata:
  author: Your Name
  license: MIT
  readme: |
    Description and usage instructions
```

## Example configuration

Complete configuration demonstrating key features:

```yaml
agents:
  root:
    model: claude
    description: Technical lead
    instruction: Coordinate development tasks and delegate to specialists
    sub_agents: [developer, reviewer]
    toolsets:
      - type: filesystem
      - type: mcp
        ref: docker:duckduckgo
    rag: [readmes]
    commands:
      status: "Check project status"

  developer:
    model: gpt
    description: Software developer
    instruction: Write clean, maintainable code
    toolsets:
      - type: filesystem
      - type: shell

  reviewer:
    model: claude
    description: Code reviewer
    instruction: Review for quality and security
    toolsets:
      - type: filesystem

models:
  gpt:
    provider: openai
    model: gpt-5

  claude:
    provider: anthropic
    model: claude-sonnet-4-5
    max_tokens: 64000

rag:
  readmes:
    docs: ["**/README.md"]
    strategies:
      - type: chunked-embeddings
        embedding_model: openai/text-embedding-3-small
        vector_dimensions: 1536
        database: ./embeddings.db
        limit: 10
      - type: bm25
        database: ./bm25.db
        limit: 10
    results:
      fusion:
        strategy: rrf
        k: 60
      limit: 5
```

## What's next

- Read the [Toolsets reference](./toolsets.md) for detailed toolset
  documentation
- Review the [CLI reference](./cli.md) for command-line options
- Browse [example
  configurations](https://github.com/docker/cagent/tree/main/examples)
- Learn about [sharing agents](../sharing-agents.md)
