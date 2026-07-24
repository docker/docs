---
title: "Managing Context & Compaction"
description: "How to keep long Docker Agent sessions from filling the model's context window: automatic and on-demand compaction, trimming tool results, and reading the context gauge."
keywords: docker agent, ai agents, guides, context window, compaction, session compaction
weight: 50
canonical: https://docs.docker.com/ai/docker-agent/guides/compaction/
---

_How to keep long-running sessions from filling the model's context window._

## Why long sessions fill the context window

Every model has a fixed context window — a maximum number of tokens it can read per request. As a session grows, the system prompt, tool definitions, prompt files, and the full message history (including every tool call and its result) all count against that budget. A long-running agent — one that reads many files, runs many commands, or just keeps chatting for a while — eventually approaches the limit. Once a request no longer fits, the model provider rejects it and the session stalls.

Docker Agent addresses this with **compaction**: replacing older parts of the conversation with a compact AI-generated summary, freeing up room for the session to keep going. This guide covers the levers you have — automatic compaction, on-demand compaction, and trimming individual tool results — and how to read the context gauge so you know where a session stands.

## Let Docker Agent compact automatically

By default, every agent proactively compacts its own session once estimated token usage crosses **90%** of the model's context window:

```yaml
agents:
  root:
    model: anthropic/claude-sonnet-4-5
    description: A long-running research assistant
    instruction: You are a helpful assistant.
```

No configuration is required to get this behavior — it's on by default. Three fields let you tune it:

| Field | Where | Description |
| --- | --- | --- |
| `session_compaction` | agent | Set to `false` to disable automatic compaction entirely for this agent (both the proactive threshold trigger and the post-overflow auto-recovery). The manual `/compact` command still works. Default: `true`. |
| `compaction_threshold` | agent or model | Fraction of the context window (greater than `0`, at most `1`) at which proactive compaction fires. A value set on the model takes precedence over the agent-level value. Default: `0.9`. |
| `compaction_model` | agent, model, or provider | Delegate the compaction (summary-generation) call to a different, usually cheaper and faster, model. The agent-level value wins, then the model-level value, then the provider-level default. |

Lower the threshold to compact earlier and keep individual requests smaller and cheaper; raise it to keep more verbatim history in context before the first summary happens:

```yaml
models:
  primary:
    provider: anthropic
    model: claude-sonnet-4-5
    # Compact at 80% of the window instead of the default 90%.
    compaction_threshold: 0.8
```

Compaction itself is a model call — feeding the full conversation to a model and asking for a summary — and it's the most expensive call in a session simply because it's the one that runs when the context is largest. There's rarely a reason to spend your primary reasoning model on it. Point `compaction_model` at something smaller instead; every other call still runs on the primary model:

```yaml
models:
  primary:
    provider: anthropic
    model: claude-sonnet-4-5
    compaction_model: fast
  fast:
    provider: anthropic
    model: claude-haiku-4-5
```

> [!IMPORTANT]
> **Context window mismatch**
>
> If `compaction_model` has a **smaller** context window than the primary model, Docker Agent triggers compaction against the smaller window so the summary call can always ingest the full conversation. Pair the primary model with a compaction model whose window is at least as large to keep the proactive trigger aligned with the primary model's window.

Disable compaction only when you specifically want a session to keep full, unabridged history and are willing to risk hitting the context limit:

```yaml
agents:
  archivist:
    model: anthropic/claude-sonnet-4-5
    description: An assistant that never auto-compacts its sessions.
    instruction: You keep full conversation history and never lose context.
    session_compaction: false
```

See [`examples/compaction_model.yaml`](https://github.com/docker/docker-agent/blob/main/examples/compaction_model.yaml) and [`examples/compaction_threshold.yaml`](https://github.com/docker/docker-agent/blob/main/examples/compaction_threshold.yaml) for complete configurations, and [Delegating Session Compaction](../../configuration/models/index.md#delegating-session-compaction) in the Model Config reference for the full field-level details.

## Compact on demand

You don't have to wait for the automatic threshold. Two TUI commands give you direct control:

- **`/compact`** — summarize and compact the current session's history right now, regardless of how full the context window is. Useful before starting a new phase of work that doesn't need the earlier detail.
- **`/context`** — open a context-window breakdown: estimated tokens per category (system prompt, tool definitions, prompt files, messages, tool results, compaction summary), a **Live sessions** view listing the current session plus every running sub-agent session with its own context budget, and a per-file inventory of attachments and prompt files.

From `/context`, select any live session with the arrow keys and press <kbd>Enter</kbd> to explicitly compact it — including a sub-agent's session, not just the main one. This is the only way cross-agent compaction happens: there's no idle-triggered automatic compaction of sub-agent sessions, so a long-running background agent stays under your control. The request is queued onto the target session's own run loop and applied at the next safe point between model turns, so it never corrupts an in-flight turn.

```bash
$ docker agent run agent.yaml
# ... work for a while ...
# Type /context to see the current breakdown, or /compact to summarize now
```

## Trim tool results to save room

Compaction deals with the whole conversation at once. For sessions dominated by a few oversized tool results — a full build log, a large file dump — three agent-level fields let you cap the damage before it ever reaches compaction:

| Field | What it bounds | Behavior |
| --- | --- | --- |
| `max_tool_result_tokens` | Each tool result, as it's added to the session | Oversized results are truncated **middle-out**: the head and tail are kept (usually the most informative parts) and the removed middle is replaced with a truncation marker. |
| `max_old_tool_call_tokens` | The total budget for **old** tool call arguments and results | Once older tool calls exceed the budget, their content is replaced wholesale with a placeholder — freeing context space without touching recent, still-relevant calls. |
| `num_history_items` | The number of non-system conversation messages kept in history | A message-**count** limit, not a token budget. The oldest non-protected messages are dropped first once the count is exceeded; **system and user messages are always protected** and are never counted against or removed by this limit, so the assembled history can exceed `num_history_items` and every user message survives even a long single-turn agentic loop. |

`max_tool_result_tokens` and `max_old_tool_call_tokens` are approximated as `len/4` tokens (the industry rule-of-thumb of ~4 characters per token); `num_history_items` counts messages, not tokens. All three are disabled by default (`0`). Set a positive value to enable them:

```yaml
agents:
  root:
    model: anthropic/claude-sonnet-4-5
    description: An assistant whose tool results are capped at ~2000 tokens each.
    instruction: |
      You are a helpful assistant with shell access. Very large command
      outputs are truncated in the middle — the beginning and end are
      always preserved, and a marker shows where content was removed.
    max_tool_result_tokens: 2000
    toolsets:
      - type: shell
```

> [!TIP]
> **Use both together**
>
> `max_tool_result_tokens` bounds each result the moment it's recorded; `max_old_tool_call_tokens` reclaims space from calls that are no longer fresh. Combine them on tool-heavy agents (shell, filesystem) to keep the session lean well before it approaches the compaction threshold.

See [`examples/max_tool_result_tokens.yaml`](https://github.com/docker/docker-agent/blob/main/examples/max_tool_result_tokens.yaml) for a complete example, and [Agent Config](../../configuration/agents/index.md#properties-reference) for the full field reference.

## Read the context gauge

The TUI's sidebar token-usage section (and the fill bar in the [lean TUI](../../features/tui/index.md#lean-tui) status line) color-escalates as a session approaches its compaction threshold, so you can see trouble coming before a request fails:

| State | Color | Trigger |
| --- | --- | --- |
| Normal | (default) | Usage below 75% of the compaction threshold |
| Warning | Orange | Usage at or above 75% of the compaction threshold |
| Critical | Red | Usage at or above 95% of the compaction threshold |

While a compaction is running, the percentage is replaced by a **"compacting…"** indicator; token counts remain visible in the lean TUI status line. The thresholds scale with the agent's configured `compaction_threshold` (default `0.9`), so a custom value keeps a predictable visual runway — for example, a `compaction_threshold: 0.8` session turns orange at 60% usage (75% of 0.8) instead of 67.5%.

Open `/context` at any time for the full per-category breakdown behind that percentage.

## What happens to cost across compaction

Compaction summarizes history, but it never resets what a session has actually cost you. Session cost tracking is **monotonic across compaction**: the running total only ever goes up, even though the summarized conversation itself is now smaller. Check `/cost` to see the current breakdown at any point, before or after a compaction has run.

## Related: deferred tool loading

Compaction and result trimming manage context that's already in the session. If large toolsets are inflating your **starting** context instead — many MCP servers, hundreds of tools — look at [deferred tool loading](../../configuration/tools/index.md#deferred-tool-loading), which registers a toolset's tools lazily instead of eagerly at startup.
