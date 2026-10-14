---
title: "Model Picker Tool"
description: "Let the agent pick between several models per turn."
keywords: docker agent, ai agents, tools, toolsets, model picker tool
linkTitle: "Model Picker"
weight: 200
canonical: https://docs.docker.com/ai/docker-agent/tools/model-picker/
---

_Let the agent pick between several models per turn._

## Overview

The model picker tool gives an agent the ability to dynamically choose which model to use for each turn of the conversation. This is useful when you want the agent to route different types of requests to different models — for example, using a fast, inexpensive model for simple queries and a more capable model for complex reasoning tasks.

## Configuration

```yaml
toolsets:
  - type: model_picker
    models:
      - openai/gpt-5-mini
      - anthropic/claude-sonnet-4-5
      - openai/gpt-5
```

### Options

| Property | Type           | Required | Description                                                  |
| -------- | -------------- | -------- | ------------------------------------------------------------ |
| `models` | array[string]  | ✓        | List of model references the agent can choose from. Use `provider/model` format. |

## How It Works

When the model picker toolset is enabled, the agent gets two tools: `change_model` to switch to one of the configured models, and `revert_model` to return to its default model. The agent decides which model to use based on the complexity of the task, cost considerations, or other factors you describe in its instruction.

## Example

```yaml
agents:
  root:
    model: openai/gpt-5-mini  # Default model
    instruction: |
      You are a helpful assistant. For simple questions, use gpt-5-mini.
      For complex reasoning or coding tasks, switch to claude-sonnet-4-5 or gpt-5.
    toolsets:
      - type: model_picker
        models:
          - openai/gpt-5-mini
          - anthropic/claude-sonnet-4-5
          - openai/gpt-5
```

> [!TIP]
> **Cost optimization**
>
> The model picker tool is particularly useful for cost optimization: let the agent use a cheap model by default and only escalate to expensive models when necessary.

## Tool Interface

The toolset exposes two tools:

### `change_model`

| Parameter | Type   | Required | Description                                                                 |
| --------- | ------ | -------- | --------------------------------------------------------------------------- |
| `model`   | string | ✓        | The model to switch to. Must be one of the configured models.               |

### `revert_model`

Takes no parameters. Reverts the agent to its original/default model.

The switch takes effect immediately: the next inference call — including the remainder of the current agentic loop — uses the new model.
