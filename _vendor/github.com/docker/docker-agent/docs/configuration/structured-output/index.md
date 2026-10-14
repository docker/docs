---
title: "Structured Output"
description: "Force the agent to respond with JSON matching a specific schema."
keywords: docker agent, ai agents, configuration, yaml, structured output
weight: 90
canonical: https://docs.docker.com/ai/docker-agent/configuration/structured-output/
---

_Force the agent to respond with JSON matching a specific schema._

## Overview

Structured output constrains the agent's responses to match a predefined JSON schema. This is useful for building agents that need to produce machine-readable output for downstream processing, API responses, or integration with other systems.

> [!NOTE]
> **When to Use**
>
> - Building API endpoints that need consistent JSON responses
> - Data extraction and transformation pipelines
> - Agents that feed into other automated systems
> - Ensuring predictable output format for parsing

## Configuration

```yaml
agents:
  analyzer:
    model: openai/gpt-4o
    description: Code analyzer that outputs structured results
    instruction: |
      Analyze the provided code and identify issues.
      Return your findings in the structured format.
    structured_output:
      name: analysis_result
      description: Code analysis findings
      strict: true
      schema:
        type: object
        properties:
          issues:
            type: array
            items:
              type: object
              properties:
                severity:
                  type: string
                  enum: ["error", "warning", "info"]
                line:
                  type: integer
                message:
                  type: string
              required: ["severity", "line", "message"]
          summary:
            type: string
        required: ["issues", "summary"]
```

## Properties

| Property      | Type    | Required | Description                                                       |
| ------------- | ------- | -------- | ----------------------------------------------------------------- |
| `name`        | string  | ✓        | Name identifier for the output schema                             |
| `description` | string  | ✗        | Description of what the output represents                         |
| `strict`      | boolean | ✗        | Enforce strict schema adherence — `native` mode only (default: `false`) |
| `schema`      | object  | ✓        | JSON Schema defining the output structure                         |
| `mode`        | string  | ✗        | Enforcement mode: `native` (default) or `tool` (see [Modes](#modes)) |

## Modes

### `native` (default)

The schema is passed to the provider's native structured-output support (OpenAI JSON mode, Gemini JSON mode, ...). Omitting `mode` keeps this behavior.

### `tool`

```yaml
structured_output:
  mode: tool
  name: analysis_result
  schema:
    type: object
    properties:
      summary:
        type: string
    required: ["summary"]
```

In tool mode nothing is sent to the provider's native structured-output API. Instead, the runtime exposes an internal tool named `__structured_output__` whose parameters are exactly the configured schema. The model works normally — including calling other tools — and delivers its final answer by calling that tool, alone, as the only tool call of its response. The runtime validates the arguments against the schema:

- A valid call ends the turn; the validated (compacted) JSON becomes the final assistant message.
- Invalid JSON gets a detailed tool error so the model can correct itself and retry.
- If the model answers in plain text instead, the runtime injects a transient system reminder and retries (at most 2 reminders), then fails with a `structured_output_failed` error.

Tool-mode validation applies the full JSON Schema, including `additionalProperties` — unexpected fields are rejected when the schema forbids them. External `$ref` references (`http(s)://`, `file://`, cross-document) are rejected when the schema is compiled; only same-document references starting with `#` (e.g. `#/definitions/item`) are allowed. The `strict` flag has no effect in tool mode.

Fork-mode skills (`context: fork`) run as exempt sub-sessions: the skill produces its own plain-text answer for the calling agent and is not required to call the output tool. The parent agent still delivers its final answer through the tool.

Use tool mode when the model must combine tool use with a schema-constrained final answer, or when the provider has no native structured-output support. See [`examples/structured-output-tool-mode.yaml`](https://github.com/docker/docker-agent/blob/main/examples/structured-output-tool-mode.yaml).

## Schema Format

The schema follows [JSON Schema](https://json-schema.org/) specification. Common schema types:

### Simple Object

```yaml
schema:
  type: object
  properties:
    name:
      type: string
    count:
      type: integer
    active:
      type: boolean
  required: ["name", "count"]
```

### Array of Objects

```yaml
schema:
  type: object
  properties:
    items:
      type: array
      items:
        type: object
        properties:
          id:
            type: string
          value:
            type: number
        required: ["id", "value"]
  required: ["items"]
```

### Enum Values

```yaml
schema:
  type: object
  properties:
    status:
      type: string
      enum: ["pending", "approved", "rejected"]
    priority:
      type: string
      enum: ["low", "medium", "high", "critical"]
  required: ["status"]
```

## Strict Mode

`strict` only applies to `native` mode: it is passed to the provider's structured-output API. Tool mode ignores it and always validates against the full schema instead. When `strict: true`, the model is constrained to only produce output that exactly matches the schema. This provides stronger guarantees but may limit the model's flexibility.

- **`strict: false` (default)** — model aims to match the schema but may include additional fields or slight variations.
- **`strict: true`** — model output is constrained to exactly match the schema. Stronger guarantees.

## Provider Support

Structured output support varies by provider:

| Provider      | Support    | Notes                                   |
| ------------- | ---------- | --------------------------------------- |
| OpenAI        | ✓ Full     | Native JSON mode with schema validation |
| Anthropic     | ✓ Full     | Tool-based structured output            |
| Google Gemini | ✓ Full     | Native JSON mode                        |
| AWS Bedrock   | ✓ Partial  | Depends on underlying model             |
| DMR           | ⚠️ Limited | Depends on model capabilities           |

## Example: Data Extraction Agent

```yaml
agents:
  extractor:
    model: openai/gpt-4o
    description: Extract structured data from text
    instruction: |
      Extract contact information from the provided text.
      Return all found contacts in the structured format.
    structured_output:
      name: contacts
      description: Extracted contact information
      strict: true
      schema:
        type: object
        properties:
          contacts:
            type: array
            items:
              type: object
              properties:
                name:
                  type: string
                  description: Full name of the contact
                email:
                  type: string
                  description: Email address
                phone:
                  type: string
                  description: Phone number
                company:
                  type: string
                  description: Company or organization
              required: ["name"]
          total_found:
            type: integer
            description: Total number of contacts found
        required: ["contacts", "total_found"]
```

## Example: Classification Agent

```yaml
agents:
  classifier:
    model: anthropic/claude-sonnet-4-5
    description: Classify support tickets
    instruction: |
      Classify the support ticket into the appropriate category
      and priority level based on its content.
    structured_output:
      name: ticket_classification
      strict: true
      schema:
        type: object
        properties:
          category:
            type: string
            enum:
              ["billing", "technical", "account", "feature_request", "other"]
          priority:
            type: string
            enum: ["low", "medium", "high", "urgent"]
          confidence:
            type: number
            minimum: 0
            maximum: 1
            description: Confidence score between 0 and 1
          reasoning:
            type: string
            description: Brief explanation for the classification
        required: ["category", "priority", "confidence"]
```

> [!WARNING]
> **Tool Limitations**
>
> When using native structured output, the agent typically cannot use tools since its response format is constrained to the schema. Design your agent workflow accordingly — native structured output agents work best for single-turn analysis or extraction tasks. Use `mode: tool` when the agent needs to call tools before producing its schema-constrained final answer.
