---
title: "Model Routing"
description: "Route requests to different models based on the content of user messages."
keywords: docker agent, ai agents, configuration, yaml, model routing
weight: 100
canonical: https://docs.docker.com/ai/docker-agent/configuration/routing/
---

_Route requests to different models based on the content of user messages._

## Overview

Model routing lets you define a "router" model that automatically selects the best underlying model based on the user's message. This is useful for cost optimization, specialized handling, or load balancing across models.

> [!NOTE]
> **How It Works**
>
> docker-agent uses NLP-based text similarity (via Bleve full-text search) to match user messages against example phrases you define. The route with the best-matching examples wins, and that model handles the request.

## Configuration

Add `routing` rules to any model definition. The model's `provider`/`model` fields become the fallback when no route matches:

```yaml
models:
  smart_router:
    # Fallback model when no routing rule matches
    provider: openai
    model: gpt-5-mini

    # Routing rules
    routing:
      - model: anthropic/claude-sonnet-4-5
        examples:
          - "Write a detailed technical document"
          - "Help me architect this system"
          - "Review this code for security issues"
          - "Explain this complex algorithm"

      - model: openai/gpt-5
        examples:
          - "Generate some creative ideas"
          - "Write a story about"
          - "Help me brainstorm"
          - "Come up with names for"

      - model: openai/gpt-5-mini
        examples:
          - "What time is it"
          - "Convert this to JSON"
          - "Simple math calculation"
          - "Translate this word"

agents:
  root:
    model: smart_router
    description: Assistant with intelligent model routing
    instruction: You are a helpful assistant.
```

## Routing Rules

Each routing rule has:

| Field      | Type   | Required | Description                                                   |
| ---------- | ------ | -------- | ------------------------------------------------------------- |
| `model`    | string | ✓        | Target model (inline format or reference to `models` section) |
| `examples` | array  | ✓        | Example phrases that should route to this model               |

## Matching Behavior

The router:

1. Extracts the last user message from the conversation
2. Searches all examples using full-text search
3. Aggregates match scores by route (best score per route wins)
4. Selects the route with the highest overall score
5. Falls back to the base model if no good match is found

> [!TIP]
> **Writing Good Examples**
>
> - Use diverse phrasing that captures the intent
> - Include keywords users actually use
> - Add 5-10 examples per route for best results
> - Examples don't need to be exact matches — the router uses semantic similarity

## Use Cases

### Cost Optimization

Route simple queries to cheaper models:

```yaml
models:
  cost_optimizer:
    provider: openai
    model: gpt-5-mini # Cheap fallback
    routing:
      - model: anthropic/claude-sonnet-4-5
        examples:
          - "Complex analysis"
          - "Detailed research"
          - "Multi-step reasoning"
```

### Specialized Models

Route coding tasks to code-specialized models:

```yaml
models:
  task_router:
    provider: openai
    model: gpt-5-mini # General fallback
    routing:
      - model: anthropic/claude-sonnet-4-5
        examples:
          - "Write code"
          - "Debug this function"
          - "Review my implementation"
          - "Fix this bug"
      - model: openai/gpt-5
        examples:
          - "Write a blog post"
          - "Help me with writing"
          - "Summarize this document"
```

### Load Balancing

Distribute load across equivalent models from different providers:

```yaml
models:
  load_balancer:
    provider: openai
    model: gpt-5-mini
    routing:
      - model: anthropic/claude-sonnet-4-5
        examples:
          - "First request pattern"
          - "Another request type"
      - model: google/gemini-2.5-flash
        examples:
          - "Different request pattern"
          - "Alternative query style"
```

## Debugging

Enable debug logging to see routing decisions:

```bash
$ docker agent run config.yaml --debug
```

Look for log entries like:

```text
"Rule-based router selected model" router=smart_router selected_model=anthropic/claude-sonnet-4-5
"Route matched" model=anthropic/claude-sonnet-4-5 score=2.45
```

> [!WARNING]
> **Limitations**
>
> - Routing only considers the last user message, not full conversation context
> - Very short messages may not match well — consider your fallback carefully
> - Each routed model creates a separate provider connection
