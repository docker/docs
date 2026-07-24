---
title: "Budget"
description: "Cap what a single run may spend in money, tokens, or working time."
keywords: docker agent, ai agents, configuration, yaml, budget, cost, limits
weight: 75
canonical: https://docs.docker.com/ai/docker-agent/configuration/budget/
---

_Cap what a single run may spend in money, tokens, or working time._

## Overview

A budget sets ceilings on a run. When a ceiling is crossed the run stops with a message naming the exact limit that tripped, and the TUI tracks consumption live while the run is still going.

Every limit is optional and an unset limit is unlimited, so you can cap money, tokens, time, or any combination. Declare no budget at all and runs are unbudgeted, which is the default.

There are two ways to declare one, and they compose:

| Block | Scope |
| --- | --- |
| `budget` | One run-wide ceiling, charged for every agent. |
| `budgets` | Named budgets an agent opts into by name. |

### A run-wide budget

```yaml
agents:
  root:
    model: openai/gpt-4o-mini
    description: An agent on a short leash.
    instruction: You are a helpful assistant.

budget:
  max_cost: 0.50
  max_tokens: 100000
  max_time: 10m
```

| Field | Type | Description |
| --- | --- | --- |
| `max_cost` | number | Maximum spend, in USD. |
| `max_tokens` | integer | Maximum cumulative input+output tokens. |
| `max_time` | string | Maximum time the agents spend working, in Go duration format (`10m`, `30s`, `1h30m`). |

### Named budgets

Define budgets by name under the top-level `budgets` key, then have each agent opt in by listing names in its own `budgets` field. The fields are the same three.

```yaml
budgets:
  shell-work:
    max_cost: 0.03
    max_tokens: 8000
  research:
    max_time: 1m

agents:
  root:
    model: openai/gpt-4o-mini
    description: Does shell work.
    instruction: You are a helpful assistant.
    budgets: [shell-work]
  researcher:
    model: openai/gpt-4o-mini
    description: Answers questions.
    instruction: You answer questions concisely.
    budgets: [shell-work, research]
```

An agent may list several budgets; all of them apply, and the first to be exhausted stops the run. A run-wide `budget` applies on top of any named budget, so the ceiling that binds is whichever runs out first.

Referencing a budget name that isn't defined is a config error, caught at parse time rather than silently leaving the agent uncapped.

## A name is one shared pot

When several agents reference the same budget name they draw from the **same** ceiling — not a copy each. Above, `root` and `researcher` share `shell-work`: together they cannot spend more than `$0.03`.

This is deliberate, and it is the whole reason budgets are worth having. If each agent received its own allowance, a run could spend `max_cost` × N simply by fanning out to N sub-agents, and the ceiling would mean nothing for exactly the workloads that most need one.

Give agents **distinct budget names** when you want independent pots.

The same applies to the run-wide `budget`: every sub-session inside a run (transferred tasks, sub-agents, skills) spends from that one wallet.

## Scope: a budget spans the session

Spend accumulates for the life of the **session**, across every message you send — it does not reset each time you hit enter. A `max_cost: 0.50` you could re-spend on every message would not be a ceiling at all.

Starting a new session starts a fresh budget. It is a per-session ceiling, not a lifetime quota across sessions.

> [!NOTE]
> `max_time` measures the time the agents actually spend **working** — the sum of their turn durations — not wall-clock since the session opened. Because a budget spans a session, and a session sits idle while you read and type, wall-clock would let a budget expire during a coffee break: leave the TUI open for ten minutes and your next message would instantly trip a `max_time` of `2m` without the agent having done anything.

## What stops a run

Crossing any limit — run-wide or named — stops the run and produces:

- an assistant message in the transcript naming the limit and the amounts,
- a `budget_exceeded` event carrying `budget`, `limit`, `used`, `max` and `config_path`,
- a `notification` hook at `warning` level,
- a stream end reason of `budget_exceeded`, so a stopped run is distinguishable from a completed one in telemetry.

The message names the exact YAML path to raise, so there is no ambiguity about which of several budgets tripped:

```text
Execution stopped after reaching the configured budgets.shell-work.max_cost
limit (used $0.0312 of $0.0300).
```

Unlike [`max_iterations`](../agents/index.md), a budget stop is **terminal** — there is no prompt offering to continue. A budget is a ceiling you set deliberately, so raising it means editing the config rather than answering a dialog.

## Tracking spend in the TUI

The sidebar's Token Usage panel lists every active budget by name, with consumption against each ceiling it declares:

```text
run        $0.12/$0.50 · 12.3K/100.0K · 2m14s/10m
shell-work $0.09/$0.10 · 4.3K/20.0K
```

Only the ceilings you configured appear. Each reading is colored by the sidebar's shared gauge bands — the same ones a context gauge uses as it nears compaction — so a budget turns amber well before its ceiling and red just short of it, and a run about to be stopped is visible before it stops.

For a per-agent view of who spent what, set **Sidebar info mode** to `Detailed` in `/settings` → Appearance: the Agents section then reports each agent's cost alongside its effort and context. The budget line deliberately does not repeat that breakdown — the same numbers twice would crowd the sidebar's narrowest column. The per-agent split is still carried on the `budget_usage` event for programmatic consumers, and `/cost` has a **By Agent** section.

## Limits and caveats

### Unpriced models do not count towards `max_cost`

Only responses the runtime can price count towards `max_cost`. A model with no pricing data — an unknown model ID, or a custom endpoint such as a local or private deployment — contributes nothing, because there is no honest number to add.

Such a run emits a warning and the TUI marks the reading `(unpriced spend)`, rather than silently reading low because the spend is invisible. To make a custom endpoint count, price it explicitly with a model-level [`cost`](../models/index.md#custom-token-pricing) block:

```yaml
models:
  local:
    provider: openai
    model: my-model
    base_url: http://localhost:8000/v1
    cost:
      input: 0.15
      output: 0.60

agents:
  root:
    model: local
    description: A locally-served agent with real cost accounting.
    instruction: You are a helpful assistant.

budget:
  max_cost: 0.50
```

### Limits are checked at turn boundaries

A run is checked between turns, so it can overshoot by at most the turn already in flight. `max_time` in particular will not interrupt a model call or tool that has already started; the run stops at the first boundary after the limit is reached.

This is the same granularity [`max_iterations`](../agents/index.md) has, and it keeps the ceiling out of the streaming hot path. Set limits with a little headroom rather than at the exact number you cannot exceed.

### `max_tokens` here is not the model's `max_tokens`

The `max_tokens` in a budget is a **cumulative** count of input+output tokens across the whole run. It is unrelated to the provider- or model-level [`max_tokens`](../models/index.md), which caps the output of a single response.

It is also not the session's context length: compaction resets that, while the budget keeps counting.

## Examples

Cap money only, and let the run take as long as it needs:

```yaml
agents:
  root:
    model: openai/gpt-4o
    description: A cost-capped agent.
    instruction: You are a helpful assistant.

budget:
  max_cost: 5.00
```

Cap working time for an unattended job:

```yaml
agents:
  root:
    model: openai/gpt-4o-mini
    description: A time-boxed agent.
    instruction: You are a helpful assistant.
    toolsets:
      - type: shell

budget:
  max_time: 15m
```

See [`examples/budget.yaml`](https://github.com/docker/docker-agent/blob/main/examples/budget.yaml) for a runnable configuration.
