---
title: "Tips & Best Practices"
description: "Expert guidance for building effective, efficient, and secure agents."
keywords: docker agent, ai agents, guides, tips & best practices
weight: 10
aliases:
  - /ai/docker-agent/best-practices/
---

_Expert guidance for building effective, efficient, and secure agents._

## Configuration Tips

### Auto Mode for Quick Start

Don't have a config file? docker-agent can automatically detect your available API keys and use an appropriate model:

```bash
# Automatically uses the best available provider
$ docker agent run

# Provider priority: Anthropic → OpenAI → Google → Mistral → Amazon Bedrock → DMR
```

The special `auto` model value also works in configs:

```yaml
agents:
  root:
    model: auto # Uses best available provider
    description: Adaptive assistant
    instruction: You are a helpful assistant.
```

### Environment Variable Interpolation

Commands support JavaScript template literal syntax for environment variables:

```yaml
agents:
  root:
    model: openai/gpt-4o
    description: Deployment assistant
    instruction: You help with deployments.
    commands:
      # Simple variable
      greet: "Hello ${env.USER}!"

      # With default value
      deploy: "Deploy to ${env.ENV || 'staging'}"

      # Multiple variables
      release: "Release ${env.PROJECT} v${env.VERSION || '1.0.0'}"
```

### Model Aliases Are Auto-Pinned

docker-agent automatically resolves model aliases to their latest pinned versions. This ensures reproducible behavior:

```yaml
# You write:
model: anthropic/claude-sonnet-4-5

# docker-agent resolves to:
# anthropic/claude-sonnet-4-5-20250929 (or latest available)
```

To use a specific version, specify it explicitly in your config.

## Performance Tips

### Defer Tools for Faster Startup

Large MCP toolsets can slow down agent startup. Use `defer` to load tools on-demand:

```yaml
agents:
  root:
    model: openai/gpt-4o
    description: Multi-tool assistant
    instruction: You have many tools available.
    toolsets:
      - type: mcp
        ref: docker:github-official
        defer: true
      - type: mcp
        ref: docker:slack
        defer: true
      - type: mcp
        ref: docker:linear
        defer: true
```

Or defer specific tools within a toolset:

```yaml
toolsets:
  - type: mcp
    ref: docker:github-official
    defer:
      - "list_issues"
      - "search_repos"
  - type: mcp
    ref: docker:slack
    defer:
      - "list_channels"
```

### Filter MCP Tools

Many MCP servers expose dozens of tools. Filter to only what you need:

```yaml
toolsets:
  - type: mcp
    ref: docker:github-official
    # Only expose these specific tools
    tools:
      - list_issues
      - create_issue
      - get_pull_request
      - create_pull_request
```

Fewer tools means faster tool selection and less confusion for the model.

### Set max_iterations

Always set `max_iterations` for agents with powerful tools to prevent infinite loops:

```yaml
agents:
  developer:
    model: anthropic/claude-sonnet-4-5
    description: Development assistant
    instruction: You are a developer.
    max_iterations: 30 # Reasonable limit for development tasks
    toolsets:
      - type: filesystem
      - type: shell
```

Typical values: 20-30 for development agents, 10-15 for simple tasks.

## Reliability Tips

### Use Fallback Models

Configure fallback models for resilience against provider outages or rate limits:

```yaml
agents:
  root:
    model: anthropic/claude-sonnet-4-5
    description: Reliable assistant
    instruction: You are a helpful assistant.
    fallback:
      models:
        # Different provider for resilience
        - openai/gpt-4o
        # Cheaper model as last resort
        - openai/gpt-4o-mini
      retries: 2 # Retry 5xx errors twice
      cooldown: 1m # Stick with fallback for 1 min after rate limit
```

**Best practices for fallback chains:**

- Use different providers for true redundancy
- Order by preference (best first)
- Include a cheaper/faster model as last resort

### Use Think Tool for Non-Reasoning Models

The `think` tool provides a reasoning scratchpad for models that lack built-in thinking capabilities:

```yaml
toolsets:
  - type: think # Useful for models without native reasoning
```

The agent uses it as a scratchpad for planning and decision-making. If your model already supports a [thinking budget](../../configuration/models/index.md#thinking-budget) (e.g., Claude with extended thinking, OpenAI o-series, Gemini with thinking enabled), you don't need this tool — the model can reason internally.

## Security Tips

### Use --yolo Mode Carefully

The `--yolo` flag auto-approves all tool calls without confirmation:

```bash
# Auto-approve everything (use with caution!)
$ docker agent run agent.yaml --yolo
```

**When it's appropriate:**

- CI/CD pipelines with controlled inputs
- Automated testing
- Agents with only safe, read-only tools

**When to avoid:**

- Interactive sessions with untested prompts
- Agents with shell or filesystem write access
- Any situation where unreviewed actions could cause harm

### Combine Permissions with Sandbox

For defense in depth, use both permissions and [sandbox mode](../../configuration/sandbox/index.md):

```yaml
agents:
  secure_dev:
    model: anthropic/claude-sonnet-4-5
    description: Secure development assistant
    instruction: You are a secure coding assistant.
    toolsets:
      - type: filesystem
      - type: shell

permissions:
  allow:
    - "read_*"
    - "shell:cmd=go*"
    - "shell:cmd=npm*"
  deny:
    - "shell:cmd=sudo*"
    - "shell:cmd=rm*-rf*"
```

```bash
# Run with sandbox enabled
docker-agent run --sandbox agent.yaml
```

### Set Global Permission Guardrails

Use [global permissions](../../configuration/permissions/index.md#global-permissions) in your user config to enforce safety rules across every agent:

```yaml
# ~/.config/cagent/config.yaml
settings:
  permissions:
    deny:
      - "shell:cmd=sudo*"
      - "shell:cmd=rm*-rf*"
      - "shell:cmd=git push --force*"
    allow:
      - "read_*"
      - "shell:cmd=ls*"
      - "shell:cmd=cat*"
```

These rules merge with any agent-level permissions. Deny patterns from your global config cannot be overridden by agent configs, so you can trust that dangerous commands stay blocked regardless of which agent you run.

### Use Hooks for Audit Logging

Log all tool calls for compliance or debugging:

```yaml
agents:
  audited:
    model: openai/gpt-4o
    description: Audited assistant
    instruction: You are a helpful assistant.
    hooks:
      post_tool_use:
        - matcher: "*"
          hooks:
            - type: command
              command: "./scripts/audit-log.sh"
```

## Multi-Agent Tips

### Handoffs vs Sub-Agents

Understand the difference between `sub_agents` and `handoffs`:

- **`sub_agents` (transfer_task)** — delegates a task to a child in a sub-session, waits for the result, then continues. Hierarchical: the parent remains in control.

  ```yaml
  sub_agents: [researcher, writer]
  ```

- **`handoffs` (peer-to-peer)** — hands off the entire conversation to another agent in the same session. The active agent switches and sees the full history. Agents can form cycles.

  ```yaml
  handoffs:
    - specialist
    - summarizer
  ```

See [Multi-Agent Systems](../../concepts/multi-agent/index.md) for a detailed comparison.

### Give Sub-Agents Clear Descriptions

The root agent uses descriptions to decide which sub-agent to delegate to:

```yaml
agents:
  root:
    model: anthropic/claude-sonnet-4-5
    description: Technical lead
    instruction: Delegate to specialists based on the task.
    sub_agents: [frontend, backend, devops]

  frontend:
    model: openai/gpt-4o
    # Good: specific and actionable
    description: |
      Frontend specialist. Handles React, TypeScript, CSS, 
      UI components, and browser-related issues.

  backend:
    model: openai/gpt-4o
    # Good: clear domain boundaries
    description: |
      Backend specialist. Handles APIs, databases, 
      server logic, and Go/Python code.

  devops:
    model: openai/gpt-4o
    description: |
      DevOps specialist. Handles CI/CD, Docker, Kubernetes,
      infrastructure, and deployment pipelines.
```

## Debugging Tips

### Enable Debug Logging

Use the `--debug` flag to see detailed execution logs:

```bash
# Default log location: ~/.cagent/cagent.debug.log
$ docker agent run agent.yaml --debug

# Custom log location
$ docker agent run agent.yaml --debug --log-file ./debug.log
```

### Check Token Usage

Use the `/cost` command during a session to see token consumption:

```text
/cost

Token Usage:
  Input:  12,456 tokens
  Output:  3,789 tokens
  Total:  16,245 tokens
```

### Compact Long Sessions

If a session gets too long, use `/compact` to summarize and reduce context:

```text
/compact

Session compacted. Summary generated and history trimmed.
```

## More Tips

### User-Defined Default Model

Set your preferred default model in `~/.config/cagent/config.yaml`:

```yaml
settings:
  default_model: anthropic/claude-sonnet-4-5
```

This model is used when you run `docker agent run` without a config file.

### GitHub PR Reviewer Example

Use docker-agent as a GitHub Actions PR reviewer:

```yaml
# .github/workflows/pr-review.yml
name: PR Review
on:
  pull_request:
    types: [opened, synchronize]

jobs:
  review:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Run docker-agent review
        env:
          ANTHROPIC_API_KEY: ${{ secrets.ANTHROPIC_API_KEY }}
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
        run: |
          # Install docker-agent
          curl -fsSL https://get.docker-agent.dev | sh

          # Run the review
          docker agent run --exec reviewer.yaml --yolo \
            "Review PR #${{ github.event.pull_request.number }}"
```

With a simple reviewer agent:

```yaml
# reviewer.yaml
agents:
  root:
    model: anthropic/claude-sonnet-4-5
    description: PR reviewer
    instruction: |
      Review pull requests for code quality, bugs, and security issues.
      Be constructive and specific in your feedback.
    toolsets:
      - type: mcp
        ref: docker:github-official
      - type: think
```
