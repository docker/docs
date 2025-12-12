---
title: Best practices
description: Patterns and techniques for building effective cagent agents
keywords: [cagent, best practices, patterns, agent design, optimization]
weight: 40
---

Patterns you learn from building and running cagent agents. These aren't
features or configuration options - they're approaches that work well in
practice.

## Handling large command outputs

Shell commands that produce large output can overflow your agent's context
window. Validation tools, test suites, and build logs often generate thousands
of lines. If you capture this output directly, it consumes all available context
and the agent fails.

The solution: redirect output to a file, then read the file. The Read tool
automatically truncates large files to 2000 lines, and your agent can navigate
through it if needed.

**Don't do this:**

```yaml
reviewer:
  instruction: |
    Run validation: `docker buildx bake validate`
    Check the output for errors.
  toolsets:
    - type: shell
```

The validation output goes directly into context. If it's large, the agent fails
with a context overflow error.

**Do this:**

```yaml
reviewer:
  instruction: |
    Run validation and save output:
    `docker buildx bake validate > validation.log 2>&1`

    Read validation.log to check for errors.
    The file can be large - read the first 2000 lines.
    Errors usually appear at the beginning.
  toolsets:
    - type: filesystem
    - type: shell
```

The output goes to a file, not context. The agent reads what it needs using the
filesystem toolset.

## Structuring agent teams

A single agent handling multiple responsibilities makes instructions complex and
behavior unpredictable. Breaking work across specialized agents produces better
results.

The coordinator pattern works well: a root agent understands the overall task
and delegates to specialists. Each specialist focuses on one thing.

**Example: Documentation writing team**

```yaml
agents:
  root:
    description: Technical writing coordinator
    instruction: |
      Coordinate documentation work:
      1. Delegate to writer for content creation
      2. Delegate to editor for formatting polish
      3. Delegate to reviewer for validation
      4. Loop back through editor if reviewer finds issues
    sub_agents: [writer, editor, reviewer]
    toolsets: [filesystem, todo]

  writer:
    description: Creates and edits documentation content
    instruction: |
      Write clear, practical documentation.
      Focus on content quality - the editor handles formatting.
    toolsets: [filesystem, think]

  editor:
    description: Polishes formatting and style
    instruction: |
      Fix formatting issues, wrap lines, run prettier.
      Remove AI-isms and polish style.
      Don't change meaning or add content.
    toolsets: [filesystem, shell]

  reviewer:
    description: Runs validation tools
    instruction: |
      Run validation suite, report failures.
    toolsets: [filesystem, shell]
```

Each agent has clear responsibilities. The writer doesn't worry about line
wrapping. The editor doesn't generate content. The reviewer just runs tools.

**When to use teams:**

- Multiple distinct steps in your workflow
- Different skills required (writing ↔ editing ↔ testing)
- One step might need to retry based on later feedback

**When to use a single agent:**

- Simple, focused tasks
- All work happens in one step
- Adding coordination overhead doesn't help

## Optimizing RAG performance

RAG indexing takes time when you have many files. A configuration that indexes
your entire codebase might take minutes to start. Optimize for what your agent
actually needs.

**Narrow the scope:**

Don't index everything. Index what's relevant for the agent's work.

```yaml
# Too broad - indexes entire codebase
rag:
  codebase:
    docs: [./]

# Better - indexes only relevant directories
rag:
  codebase:
    docs: [./src/api, ./docs, ./examples]
```

If your agent only works with API code, don't index tests, vendor directories,
or generated files.

**Increase batching and concurrency:**

Process more chunks per API call and make parallel requests.

```yaml
strategies:
  - type: chunked-embeddings
    embedding_model: openai/text-embedding-3-small
    batch_size: 50 # More chunks per API call
    max_embedding_concurrency: 10 # Parallel requests
    chunking:
      size: 2000 # Larger chunks = fewer total chunks
      overlap: 150
```

This reduces both API calls and indexing time.

**Consider BM25 for fast local search:**

If you need exact term matching (function names, error messages, identifiers),
BM25 is fast and runs locally without API calls.

```yaml
strategies:
  - type: bm25
    database: ./bm25.db
    chunking:
      size: 1500
```

Combine with embeddings using hybrid retrieval when you need both semantic
understanding and exact matching.

## Preserving document scope

When building agents that update documentation, a common problem: the agent
transforms minimal guides into tutorials. It adds prerequisites,
troubleshooting, best practices, examples, and detailed explanations to
everything.

These additions might individually be good, but they change the document's
character. A focused 90-line how-to becomes a 200-line reference.

**Build this into instructions:**

```yaml
writer:
  instruction: |
    When updating documentation:

    1. Understand the current document's scope and length
    2. Match that character - don't transform minimal guides into tutorials
    3. Add only what's genuinely missing
    4. Value brevity - not every topic needs comprehensive coverage

    Good additions fill gaps. Bad additions change the document's character.
    When in doubt, add less rather than more.
```

Tell your agents explicitly to preserve the existing document's scope. Without
this guidance, they default to being comprehensive.

## Model selection

Choose models based on the agent's role and complexity.

**Use larger models (Sonnet, GPT-5) for:**

- Complex reasoning and planning
- Writing and editing content
- Coordinating multiple agents
- Tasks requiring judgment and creativity

**Use smaller models (Haiku, GPT-5 Mini) for:**

- Running validation tools
- Simple structured tasks
- Reading logs and reporting errors
- High-volume, low-complexity work

Example from the documentation writing team:

```yaml
agents:
  root:
    model: anthropic/claude-sonnet-4-5 # Complex coordination
  writer:
    model: anthropic/claude-sonnet-4-5 # Creative content work
  editor:
    model: anthropic/claude-sonnet-4-5 # Judgment about style
  reviewer:
    model: anthropic/claude-haiku-4-5 # Just runs validation
```

The reviewer uses Haiku because it runs commands and checks for errors. No
complex reasoning needed, and Haiku is faster and cheaper.

## What's next

- Review [configuration reference](./reference/config.md) for all available
  options
- Check [toolsets reference](./reference/toolsets.md) to understand what tools
  agents can use
- See [example
  configurations](https://github.com/docker/cagent/tree/main/examples) for
  complete working agents
- Read the [RAG guide](./rag.md) for detailed retrieval optimization
