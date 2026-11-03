---
title: "RAG Tool"
description: "Give your agents access to document knowledge bases with background indexing, multiple retrieval strategies, and hybrid search."
keywords: docker agent, ai agents, tools, toolsets, rag tool
linkTitle: "RAG"
weight: 110
canonical: https://docs.docker.com/ai/docker-agent/tools/rag/
aliases:
  - /ai/docker-agent/rag/
---

_Give your agents access to document knowledge bases with background indexing, multiple retrieval strategies, and hybrid search._

## Overview

The `rag` toolset lets agents search through your documents to find relevant information before responding. Knowledge bases are declared once at the top of the config under `rag:` and then referenced from any agent via `type: rag, ref: <name>`. Docker Agent supports:

- **Background indexing** — Files are indexed automatically and re-indexed on change
- **Multiple strategies** — Semantic embeddings, BM25 keyword search, and LLM-enhanced search
- **Hybrid search** — Combine strategies with result fusion for best results
- **Reranking** — Re-score results with specialized models for improved relevance

RAG is the strategy to reach for when a document collection is too large to inline directly, or gets queried repeatedly across turns/sessions — see [Choosing a Large-Input Strategy](../../guides/headless/index.md#choosing-a-large-input-strategy) for how it compares to `@`/`/attach` attachments and prompt files.

## Quick Start

```yaml
rag:
  my_docs:
    tool:
      description: "Technical documentation"
    docs: [./documents, ./some-doc.md]
    strategies:
      - type: chunked-embeddings
        embedding_model: openai/text-embedding-3-small
        database: ./docs.db
        vector_dimensions: 1536

agents:
  root:
    model: openai/gpt-4o
    instruction: |
      You have access to a knowledge base. Use it to answer questions.
    toolsets:
      - type: rag
        ref: my_docs
```

### Sharing a knowledge base between agents

When several agents reference the same knowledge base, they share one index.
Docker Agent indexes each document once, runs a single file watcher, and
re-embeds a changed file once, no matter how many agents use it. Embedding
calls are bounded by the knowledge base's own `max_embedding_concurrency`, not
multiplied by the number of agents.

References share an index only when they resolve to the same configuration
and the same credentials. A reference that overrides `name:`, or a session
loaded with different credentials, gets its own index. Background indexing
for a shared knowledge base is not attributed to any session.

```yaml
agents:
  researcher:
    model: anthropic/claude-sonnet-4-5
    toolsets:
      - type: rag
        ref: my_docs   # shares one index...
  writer:
    model: anthropic/claude-sonnet-4-5
    toolsets:
      - type: rag
        ref: my_docs   # ...with this reference
```

## Retrieval Strategies

### Chunked Embeddings (Semantic Search)

Uses embedding models to find semantically similar content. Best for understanding intent, synonyms, and paraphrasing.

```yaml
strategies:
  - type: chunked-embeddings
    embedding_model: openai/text-embedding-3-small
    database: ./vector.db
    vector_dimensions: 1536
    similarity_metric: cosine_similarity
    threshold: 0.5
    limit: 10
    embedding_batch_size: 50
    chunking:
      size: 1000
      overlap: 100
```

Supported embedding providers: **OpenAI** (and OpenAI-compatible custom providers), **DMR**, **Google Gemini** (see [Embeddings](../../providers/google/index.md#embeddings)). `vector_dimensions` must match the model's output size, e.g. 1536 for `openai/text-embedding-3-small` or 3072 for `google/gemini-embedding-2`.

### Semantic Embeddings (LLM-Enhanced)

Uses an LLM to generate semantic summaries of each chunk before embedding, capturing meaning and intent. Best for code search and understanding implementations.

```yaml
strategies:
  - type: semantic-embeddings
    embedding_model: openai/text-embedding-3-small
    vector_dimensions: 1536
    chat_model: openai/gpt-4o-mini
    database: ./semantic.db
    ast_context: true # include AST metadata
    chunking:
      size: 1000
      code_aware: true # AST-aware chunking
```

> [!NOTE]
> **Trade-offs**
>
> Semantic embeddings provide higher quality retrieval but slower indexing (LLM call per chunk) and additional API costs.

### BM25 (Keyword Search)

Traditional keyword matching using the BM25 algorithm. Best for exact terms, technical jargon, and code identifiers.

```yaml
strategies:
  - type: bm25
    database: ./bm25.db
    k1: 1.5 # term frequency saturation
    b: 0.75 # length normalization
    threshold: 0.3
    limit: 10
    chunking:
      size: 1000
      overlap: 100
```

## Hybrid Search

Combine multiple strategies for best results. Strategies run in parallel and results are fused together:

```yaml
rag:
  hybrid:
    docs: [./docs]
    strategies:
      - type: chunked-embeddings
        embedding_model: openai/text-embedding-3-small
        database: ./vector.db
        vector_dimensions: 1536
        limit: 20
        chunking: { size: 1000, overlap: 100 }
      - type: bm25
        database: ./bm25.db
        limit: 15
        chunking: { size: 1000, overlap: 100 }
    results:
      fusion:
        strategy: rrf # Reciprocal Rank Fusion
        k: 60
      deduplicate: true
      limit: 5
```

## Fusion Strategies

| Strategy   | Best For                          | Description                                                        |
| ---------- | --------------------------------- | ------------------------------------------------------------------ |
| `rrf`      | General use (recommended)         | Reciprocal Rank Fusion — rank-based, no score normalization needed |
| `weighted` | Known performance characteristics | Weight strategies differently (e.g., embeddings: 0.7, BM25: 0.3)   |
| `max`      | Same scoring scale                | Takes the maximum score from any strategy                          |

## Reranking

Re-score retrieved documents with a specialized model to improve relevance:

```yaml
results:
  reranking:
    model: openai/gpt-4o-mini
    top_k: 10 # only rerank top 10
    threshold: 0.3 # minimum score after reranking
    criteria: |
      Prioritize official documentation over blog posts.
      Prefer recent information and practical examples.
  limit: 5
```

Supported reranking providers: **DMR** (native `/rerank` endpoint), **OpenAI**, **Anthropic**, **Gemini**.

## Code-Aware Chunking

For source code, enable AST-based chunking to keep functions and methods intact:

```yaml
chunking:
  size: 2000
  code_aware: true # Uses tree-sitter for AST-based chunking
```

> [!NOTE]
> **Language Support**
>
> Currently supports Go (`.go`) files. More languages will be added. Falls back to plain text chunking for unsupported file types.

## Indexing failures, retries and backoff

When a knowledge-base fails to start — because the embedding provider is rate-limiting
your requests or returning a transient server error — Docker Agent spaces out retry
attempts with bounded exponential backoff instead of hammering the provider on
every agent turn.

### What triggers backoff

For **RAG indexing**, the toolset gate arms on two different signals depending
on the failure:

- **HTTP 429 (rate limit)** aborts the whole indexing run at the *first*
  failure — continuing would just keep hammering a provider that asked for
  backoff — and that abort error reaches the gate immediately.
- **HTTP 408 (request timeout) and a fixed 5xx set** (`500, 502, 503, 504,
  529`) are otherwise handled per-file: a single file's transient
  failure is skipped so the run can keep indexing the rest. But if **no
  file in the run is successfully indexed** — every attempted file hit one
  of these retryable statuses — the run treats that as a sustained backend
  failure rather than a one-off hiccup, and surfaces the error so the gate
  arms on the next turn (fixed in
  [#4097](https://github.com/docker/docker-agent/issues/4097); previously
  only 429 reached the gate).

| Failure kind | Behaviour |
|---|---|
| HTTP 429 (rate limit) | Aborts the run on the first failure; backoff: next attempt delayed |
| HTTP 408 or 5xx, isolated to some files | Per-file skip; run succeeds, no backoff (indexed files persist, failures retried next run) |
| HTTP 408 or 5xx, affecting every file | Run fails; backoff: next attempt delayed |
| Other failures (config errors, auth, unrecognized 4xx) | Fail fast: retried every turn with no added delay |
| Context cancellation or agent shutdown | Immediate: no delay |

> [!NOTE]
> This trigger set (429 always, 408/5xx when sustained across every file) is
> specific to the RAG/embedding path. Other toolset types have their own
> trigger sets against the same gate — for example, remote MCP toolsets pace
> every connection attempt (not just a sustained run) on 408 and the same
> fixed 5xx set (see
> [MCP startup failure behaviour](../mcp/index.md#lifecycle-auto-restart-profiles)),
> and the A2A toolset paces its agent-card fetch the same way (see
> [A2A startup failure behaviour](../a2a/index.md#startup-failure-behaviour)).

### Retry policy and parameters

The backoff is **bounded exponential with additive jitter**:

- **Base delay**: 15 seconds
- **Maximum delay**: up to ~6 minutes (5-minute cap plus up to 20% additive jitter)
- **Growth**: doubles after each consecutive retryable failure (15s → 30s → 1m → 2m → 4m → 5m)
- **Jitter**: each wait is a random value in `[nominal, 1.2×nominal]` (additive 0–20%)
  so concurrent knowledge-base sources spread their retries and avoid
  hammering the provider together
- **Retry-After override**: if the embedding provider responds with a `Retry-After`
  header, that hint overrides the computed delay (capped at the 5-minute maximum,
  with the same additive jitter applied to spread concurrent retries)

The gate is a lightweight wall-clock check — it creates no background threads or
timers. A Stop command or agent shutdown takes effect immediately regardless of
how much of the backoff window remains.

### Long indexing runs

Starting any toolset — RAG included — is bounded by a 30-second *wait* budget
(`tools.DefaultStartTimeout`): if a toolset's `Start` has not returned within
30s, the caller stops waiting and the turn proceeds without that toolset's
tools, exactly as it would for a wedged MCP server. This budget exists to
detect toolsets that never come up; it is not a deadline on indexing itself.

A large knowledge base can legitimately take much longer than 30s to index.
Rather than abort in-flight indexing at the 30s mark — which used to discard
any embeddings not yet committed for the file being processed — the RAG
toolset detaches indexing from the caller's wait budget the same way it
already detaches its file watcher. When the 30s budget expires:

- The current turn proceeds without the RAG tool (the existing "taking too
  long to start" warning and TUI progress events still apply).
- Indexing keeps running in the background.
- Files that finish indexing are persisted atomically per file (a file is
  always either fully indexed or untouched — see #4073/#4158), so progress
  is never lost to this budget.
- A later turn's `Start` call picks the toolset up: if indexing has finished,
  the tool is available immediately; if it is still running, that turn also
  proceeds without the tool.

This is a **visible behavior change**: knowledge bases that used to finish
indexing (and thus offer the tool) on the very first turn, within the old
30-second window, now do so on whichever turn happens to land after indexing
completes. If indexing takes under 30s, nothing changes.

`indexing_timeout` bounds indexing itself, independent of the 30s wait
budget — it exists only so a hung provider connection cannot pin a knowledge
base's indexing lock forever:

```yaml
rag:
  codebase:
    indexing_timeout: 2h # Go duration; "0s" = unbounded; default 30m
    docs: [./knowledge-base]
    strategies:
      - type: chunked-embeddings
        embedding_model: openai/text-embedding-3-small
```

- Default: `30m` when omitted.
- `0s` means no bound. A negative value is rejected at config validation time.
- Expiry is a plain timeout, not a provider error: it does **not** arm the
  backoff gate described above, and the next turn's `Start` resumes indexing
  immediately rather than waiting out a cooldown. Files already persisted are
  kept; unfinished files are re-indexed on resume.
- A 429/408/5xx from the provider *during* a detached indexing run is
  unaffected by any of this — it still aborts the run and arms the backoff
  gate exactly as described above.

### Operational impact

**Before**: a rate-limited knowledge base was re-indexed on every agent turn —
`max_indexing_concurrency × max_embedding_concurrency` concurrent provider calls
could relaunch within milliseconds, easily tripping rate limits for both the
knowledge base and the agent's own model calls.

**After**: retries are spaced out and jittered so the provider has room to recover
before the next attempt. The agent continues working with any other toolsets that
are not affected.

### What you will see

- Docker Agent logs a single warning when a knowledge base first fails to start.
  Repeated failures in between are logged at debug level only, so you are not
  flooded with alerts on every turn. Recovery is intentionally silent — the
  tool appearing in the agent's tool list is the signal that indexing succeeded.
- The knowledge-base tool does not appear in the agent's tool list until indexing
  succeeds. A successful start is silent — the tool is listed and the agent uses it.

### Troubleshooting repeated 429/5xx/408 errors

If you see persistent `429`, `5xx`, or `408` errors in the logs:

1. **Check provider rate limits.** Your embedding API key may have a low requests-per-minute
   quota. Upgrading the plan or using a different API key can help.
2. **Reduce concurrency.** The chunked-embeddings and semantic-embeddings strategies
   accept `max_indexing_concurrency` (default `3`) and `max_embedding_concurrency`
   (default `3`) parameters. Lowering these reduces simultaneous requests:

   ```yaml
   rag:
     docs:
       docs: [./knowledge-base]
       strategies:
         - type: chunked-embeddings
           max_indexing_concurrency: 1
           max_embedding_concurrency: 1
   ```

3. **Use a model with a higher quota.** Some providers offer higher rate limits on
   specific embedding model tiers.

## Debugging RAG

Enable debug logging to see retrieval details:

```bash
$ docker agent run config.yaml --debug --log-file debug.log
```

Look for log tags: `[RAG Manager]`, `[Chunked-Embeddings Strategy]`, `[BM25 Strategy]`, `[RRF Fusion]`, `[Reranker]`.

**Permanent model errors abort early.** If the embedding model, semantic-LLM model, or reranking model returns a permanent error (HTTP 400, 401, 404, or 429 — invalid config, bad auth, unknown model, or rate limit), Docker Agent treats the model configuration as invalid and stops immediately rather than retrying doomed requests:

- **Indexing** — the entire indexing run is aborted after the first permanent failure (including 429). The error is surfaced in the logs so you know immediately if a model name or API key is wrong, rather than silently producing incomplete results.
- **Reranking** — a permanent error (including 429) permanently disables the reranker for the lifetime of the manager. Subsequent queries fall back to un-reranked results. Only transient errors (5xx, timeouts) fall back and retry on the next query.

> [!TIP]
> **Examples**
>
> See the [RAG examples](https://github.com/docker/docker-agent/tree/main/examples/rag) in the GitHub repo for complete, runnable configurations.

## Configuration Reference

### Top-Level RAG Fields

| Field              | Type     | Default | Description                                                                                                                                    |
| ------------------ | -------- | ------- | ---------------------------------------------------------------------------------------------------------------------------------------------- |
| `docs`             | []string | —       | Document paths/directories (shared across strategies)                                                                                          |
| `description`      | string   | —       | Human-readable description of this RAG source                                                                                                  |
| `respect_vcs`      | boolean  | `true`  | Respect `.gitignore` files when indexing documents                                                                                             |
| `indexing_timeout` | string   | `30m`   | Cap on a single indexing run, detached from the 30s toolset-start wait budget; `0s` = unbounded. See [Long indexing runs](#long-indexing-runs) |
| `strategies`       | []object | —       | Array of retrieval strategy configurations                                                                                                     |
| `results`          | object   | —       | Post-processing: fusion, reranking, deduplication, final limit                                                                                 |

### Chunked-Embeddings Strategy

| Field                       | Type   | Default             | Description                                                  |
| --------------------------- | ------ | ------------------- | ------------------------------------------------------------ |
| `embedding_model`           | string | —                   | **Required.** Embedding model reference                      |
| `database`                  | string | —                   | Path to local SQLite database                                |
| `vector_dimensions`         | int    | —                   | Embedding dimensions (e.g., 1536 for text-embedding-3-small) |
| `similarity_metric`         | string | `cosine_similarity` | Similarity metric                                            |
| `threshold`                 | float  | `0.5`               | Minimum similarity score (0–1)                               |
| `limit`                     | int    | `5`                 | Max results from this strategy                               |
| `embedding_batch_size`      | int    | `50`                | Chunks per embedding request                                 |
| `max_embedding_concurrency` | int    | `3`                 | Max concurrent embedding requests                            |
| `max_indexing_concurrency`  | int    | `3`                 | Max concurrent file-indexing tasks                           |
| `chunking.size`             | int    | `1500`              | Chunk size in characters (`4000` when `code_aware` is set)   |
| `chunking.overlap`          | int    | `75`                | Overlap between chunks in characters                         |
| `chunking.code_aware`       | bool   | `false`             | AST-based chunking (Go files only)                           |

### Semantic-Embeddings Strategy

| Field                      | Type   | Default    | Description                                                        |
| -------------------------- | ------ | ---------- | ------------------------------------------------------------------ |
| `embedding_model`          | string | —          | **Required.** Embedding model reference                            |
| `chat_model`               | string | —          | **Required.** LLM for generating semantic summaries                |
| `vector_dimensions`        | int    | —          | **Required.** Embedding dimensions                                 |
| `database`                 | string | —          | Path to local SQLite database                                      |
| `semantic_prompt`          | string | (built-in) | Custom prompt template (`${path}`, `${content}`, `${ast_context}`) |
| `ast_context`              | bool   | `false`    | Include tree-sitter AST metadata in prompts                        |
| `threshold`                | float  | `0.5`      | Minimum similarity score (0–1)                                     |
| `limit`                    | int    | `5`        | Max results                                                        |
| `max_indexing_concurrency` | int    | `3`        | Max concurrent file indexing                                       |
| `chunking.size`            | int    | `1500`     | Chunk size in characters (`4000` when `code_aware` is set)         |
| `chunking.overlap`         | int    | `75`       | Overlap between chunks                                             |
| `chunking.code_aware`      | bool   | `false`    | AST-based chunking                                                 |

### BM25 Strategy

| Field              | Type   | Default | Description                                     |
| ------------------ | ------ | ------- | ----------------------------------------------- |
| `database`         | string | —       | Path to local SQLite database                   |
| `k1`               | float  | `1.5`   | Term frequency saturation (1.2–2.0 recommended) |
| `b`                | float  | `0.75`  | Length normalization (0–1)                      |
| `threshold`        | float  | `0.0`   | Minimum BM25 score                              |
| `limit`            | int    | `5`     | Max results                                     |
| `chunking.size`    | int    | `1500`  | Chunk size in characters                        |
| `chunking.overlap` | int    | `75`    | Overlap between chunks                          |

### Results (Post-Processing)

| Field                 | Type   | Default | Description                                                 |
| --------------------- | ------ | ------- | ----------------------------------------------------------- |
| `fusion.strategy`     | string | `rrf`   | Fusion method: `rrf`, `weighted`, or `max`                  |
| `fusion.k`            | int    | `60`    | RRF rank constant                                           |
| `deduplicate`         | bool   | `true`  | Remove duplicate results                                    |
| `limit`               | int    | `15`    | Final number of results                                     |
| `include_score`       | bool   | `false` | Include relevance scores in results                         |
| `return_full_content` | bool   | `false` | Return full document content instead of just matched chunks |
| `reranking.model`     | string | —       | Reranking model reference                                   |
| `reranking.top_k`     | int    | (`limit`) | Only rerank top K results. Defaults to the results `limit` when set.  |
| `reranking.threshold` | float  | `0.5`   | Minimum relevance score after reranking                     |
| `reranking.criteria`  | string | —       | Custom relevance guidance for the reranking model           |
