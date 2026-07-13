---
title: "HCL Configuration"
description: "Write docker-agent configs in HCL instead of YAML, using labeled blocks, heredocs, and the same underlying schema."
keywords: docker agent, ai agents, configuration, yaml, hcl configuration
weight: 20
---

_Write docker-agent configs in HCL instead of YAML. It maps to the same docker-agent schema and validation rules._

`docker-agent` supports `.hcl` config files anywhere it supports `.yaml` or `.yml` files. HCL is useful if you prefer labeled blocks, less punctuation, and heredocs for long prompts.

> [!TIP]
> **Same config model, different syntax**
>
> YAML and HCL are just two syntaxes for the same docker-agent configuration model. docker-agent converts HCL to the equivalent YAML structure internally, then runs the normal schema validation and loading pipeline.

## Minimal Example

```hcl
#!/usr/bin/env docker agent run

agent "root" {
  model       = "openai/gpt-5"
  description = "A helpful assistant"
  instruction = <<-EOT
  You are a helpful assistant.
  EOT

  toolset "think" {}
}
```

Run it exactly like a YAML config:

```bash
$ docker agent run agent.hcl
$ docker agent run --exec agent.hcl "Summarize this repository"
$ docker agent serve api ./agents/   # directories may mix .yaml, .yml, and .hcl files
```

> [!TIP]
> **See also**
>
> HCL changes the syntax, not the meaning of fields. For what each field does, see [Agent Config](../agents/index.md), [Model Config](../models/index.md), and [Tool Config](../tools/index.md).

## YAML vs HCL

These two configs are equivalent:

```yaml
models:
  claude:
    provider: anthropic
    model: claude-sonnet-4-5

agents:
  root:
    model: claude
    description: Coding assistant
    instruction: You help with software development.
    toolsets:
      - type: filesystem
      - type: shell
```

```hcl
model "claude" {
  provider = "anthropic"
  model    = "claude-sonnet-4-5"
}

agent "root" {
  model       = "claude"
  description = "Coding assistant"
  instruction = "You help with software development."

  toolset "filesystem" {}
  toolset "shell" {}
}
```

## Core Conventions

HCL follows a few simple mapping rules:

| HCL syntax | YAML shape |
| --- | --- |
| `agent "root" { ... }` | `agents.root` |
| `model "claude" { ... }` | `models.claude` |
| `provider "team" { ... }` | `providers.team` |
| `mcp "github" { ... }` | `mcps.github` |
| `rag "docs" { ... }` | `rag.docs` |
| `command "fix" { ... }` inside an agent | `commands.fix` |
| `toolset "shell" {}` | list item in `toolsets` with `type: shell` |
| `metadata { ... }`, `permissions { ... }` | singleton blocks with the same top-level name |

### Top-level keyed maps become labeled blocks

In YAML, several sections are maps keyed by name. In HCL, those become labeled blocks:

```hcl
model "claude" {
  provider = "anthropic"
  model    = "claude-sonnet-4-5"
}

agent "root" {
  model       = "claude"
  description = "Primary assistant"
  instruction = "You are helpful."
}
```

The supported top-level labeled blocks are:

- `agent`
- `model`
- `provider`
- `mcp`
- `rag`

The supported top-level singleton blocks are:

- `metadata`
- `permissions`

### Toolsets use the block label as `type`

Instead of writing list entries with `type: ...`, HCL uses a `toolset` block whose label becomes the tool type:

```hcl
agent "root" {
  model       = "openai/gpt-5"
  description = "Dev assistant"
  instruction = "You can inspect and modify code."

  toolset "filesystem" {}

  toolset "mcp" {
    ref = "docker:github-official"
  }
}
```

### Commands use labeled blocks too

Agent commands are often nicer to write in HCL because each command gets its own block:

```hcl
agent "root" {
  model       = "openai/gpt-5"
  description = "Build helper"
  instruction = "You help with builds."

  command "fix-lint" {
    description = "Fix lint issues"
    instruction = "Run the linter, then fix any problems."
  }
}
```

## Strings and Heredocs

Use quoted strings for short values and heredocs for long prompts, welcome messages, or embedded JSON.

```hcl
agent "root" {
  model       = "openai/gpt-5"
  description = "Friendly assistant"

  instruction = <<-EOT
  You are a helpful assistant.

  Keep answers concise and practical.
  EOT
}
```

### Escaping literal `${...}`

HCL treats `${...}` inside strings and heredocs as template interpolation. If you need the literal text `${...}` in your prompt, escape it as `$${...}`.

This matters for command prompts that intentionally show docker-agent template snippets:

```hcl
command "fix-lint" {
  instruction = <<-EOT
  Run the linter and inspect the result:

  $${shell({cmd: "task lint"})}
  EOT
}
```

The model will receive the literal `${shell({cmd: "task lint"})}` text.

## Repeated Blocks Become Lists

Some YAML sections are lists. In HCL, those are written as repeated blocks.

For example, model routing rules become repeated `routing { ... }` blocks:

```hcl
model "smart_router" {
  provider = "openai"
  model    = "gpt-5"

  routing {
    model    = "anthropic/claude-sonnet-4-5"
    examples = [
      "Write a detailed technical document",
      "Review this code for security issues",
    ]
  }

  routing {
    model    = "openai/gpt-5"
    examples = [
      "Generate some creative ideas",
      "Help me brainstorm",
    ]
  }
}
```

The same idea applies to other list-shaped sections such as RAG `strategy` blocks and hook event entries.

## Important Differences from Terraform

docker-agent uses HCL as a configuration syntax, not as Terraform:

- There are no modules, `locals`, or `variable` blocks.
- There are no docker-agent-specific HCL functions to call from expressions.
- Prefer normal literal values: strings, numbers, booleans, lists, objects, and nested blocks.
- After conversion, the result is validated exactly like the equivalent YAML config.

If you already know Terraform, think of docker-agent HCL as a thin block-based syntax over the existing config schema.

## Examples

See these real configs in the repository:

- [`examples/pirate.hcl`](https://github.com/docker/docker-agent/blob/main/examples/pirate.hcl)
- [`examples/gopher.hcl`](https://github.com/docker/docker-agent/blob/main/examples/gopher.hcl)
