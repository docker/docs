---
title: "Random Tool"
description: "Generate random integers with bounds chosen by the agent."
keywords: docker agent, ai agents, tools, toolsets, random integer
linkTitle: "Random"
weight: 40
canonical: https://docs.docker.com/ai/docker-agent/tools/random/
---

The `random` toolset exposes `random_int`, which generates a uniformly random
integer between an agent-chosen minimum and maximum, **both inclusive**. It
requires no external dependencies and does not modify files or external state.

## Configuration

```yaml
agents:
  root:
    model: openai/gpt-5-mini
    instruction: Use random_int to roll dice and make random choices.
    toolsets:
      - type: random
```

Bounds are tool-call arguments, not toolset configuration: the agent chooses
them on every call.

## Parameters

| Parameter | Type | Required | Description |
| --- | --- | --- | --- |
| `min` | integer | ✓ | Inclusive minimum value. |
| `max` | integer | ✓ | Inclusive maximum value; must be greater than or equal to `min`. |

The tool returns the generated integer. Negative bounds are supported, and
setting `min` equal to `max` returns that value. Bounds must fit in a signed
64-bit integer. Missing bounds or `min > max` produce an error.

For a six-sided die, the agent calls:

```json
{"min": 1, "max": 6}
```

For a coin flip, use `min: 0` and `max: 1` and map the result to heads or tails.

Successful draws may repeat identical arguments without triggering duplicate-call
loop detection. Failed calls still count toward that limit. To bound the number
of model turns, set the agent's `max_iterations`; it defaults to unlimited.

See [`examples/random.yaml`](https://github.com/docker/docker-agent/blob/main/examples/random.yaml)
for a runnable example.
