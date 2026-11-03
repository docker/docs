---
title: "Think Tool"
description: "Step-by-step reasoning scratchpad for planning and decision-making."
keywords: docker agent, ai agents, tools, toolsets, think tool
linkTitle: "Think"
weight: 140
canonical: https://docs.docker.com/ai/docker-agent/tools/think/
---

_Step-by-step reasoning scratchpad for planning and decision-making._

## Overview

The think tool is a reasoning scratchpad that lets agents think step-by-step before acting. The agent can write its thoughts without producing visible output to the user — ideal for planning complex tasks, breaking down problems, and reasoning through multi-step solutions.

This lightweight tool changes only its internal scratchpad, not external data. It is most useful for models that lack built-in reasoning or thinking capabilities (e.g., smaller or older models). For models that already support native thinking — such as Claude with extended thinking, OpenAI o-series, or Gemini with a thinking budget — this tool is unnecessary since the model can reason internally.

## Configuration

```yaml
toolsets:
  - type: think
```

No configuration options.

> [!TIP]
> **When to use**
>
> Use the think tool with models that don't have native reasoning capabilities. If your model already supports a [thinking budget](../../configuration/models/index.md#thinking-budget), you likely don't need this tool.

## Scope and lifetime

Thoughts are isolated by conversation ID, calling agent and think toolset instance.
They accumulate across turns, runtime closure and compaction. The built-in SQLite
session store retains raw thoughts across reloads and restarts; the in-memory
store retains them only for its lifetime. This internal state is stored before
response-output redaction, so session-database and backup privacy policies apply.

A branch or fork with a new conversation ID starts empty. A same-ID session clone
addresses the same durable scratchpad; cloning messages does not provide a
transactional copy of hidden thoughts. JSON export/import does not transport the
scratchpad, and existing conversations start with empty dedicated state on upgrade.

Toolset namespaces follow configuration order (configured and extra toolsets are
separate). Resuming with reordered toolsets or another configuration with the same
agent name can reassign these namespaces. Stable identity requires unchanged topology.

Root deletion removes all root-owned thoughts, including those from aborted
sub-sessions. Deleting a persisted child removes known subtree thoughts;
never-persisted descendants remain root-owned until root deletion. Stop writers
before deleting sessions. Deletion is logical, not secure erasure or backup removal.

There is no TTL, eviction, quota or output omission: every call returns the full
scratchpad, which can grow beyond model context limits. Storage and output memory
scale with retained thoughts. Persistence, ownership or unsupported-store failures
are tool errors, never silent resets. Only genuinely unscoped hosts, such as
pre-run prompt expansion, return an ephemeral single thought.
