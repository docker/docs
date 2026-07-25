---
title: "Multi-Agent Systems"
description: "Build teams of specialized agents that collaborate and delegate tasks to each other."
keywords: docker agent, ai agents, concepts, multi-agent systems
linkTitle: "Multi-Agent"
weight: 40
canonical: https://docs.docker.com/ai/docker-agent/concepts/multi-agent/
---

_Build teams of specialized agents that collaborate and delegate tasks to each other._

## Why Multi-Agent?

Complex tasks benefit from specialization. Instead of one monolithic agent trying to do everything, you can create a **team** of focused agents:

- A **coordinator** that understands the overall goal and delegates
- A **developer** that writes code with filesystem and shell access
- A **reviewer** that checks code quality
- A **researcher** that searches the web for information

Each agent has its own model, tools, and instructions — optimized for its specific role.

## Two Patterns: Delegation vs. Handoffs

Docker Agent supports two multi-agent patterns:

| | **Delegation** (`sub_agents`) | **Handoffs** (`handoffs`) |
|---|---|---|
| **Topology** | Hierarchical (parent → child → parent) | Peer-to-peer graph (A → B → C → A) |
| **Session** | Child runs in a **sub-session** | Conversation stays in the **same session** |
| **Context** | Child gets a clean task description | Next agent sees the **full conversation history** |
| **Control flow** | Parent blocks until child finishes, then continues | Active agent switches — previous agent is no longer in the loop |
| **Tool** | `transfer_task` | `handoff` |
| **Best for** | Task delegation to specialists | Pipeline workflows, conversational routing |

You can combine both patterns in the same configuration — an agent can have both `sub_agents` and `handoffs`.

> [!TIP]
> **When to use which**
>
> **`sub_agents`** — Use when a coordinator needs to send tasks to specialists and synthesize their results.
>
> **`handoffs`** — Use when agents should take turns processing the same conversation (pipelines, routing).
>
> **`background_agents`** — Use when multiple independent tasks can run simultaneously.

## Delegation with `sub_agents`

Agents delegate tasks using the built-in `transfer_task` tool, which is automatically available to any agent with `sub_agents`. The parent agent sends a task to a child agent, waits for the result, and then continues.

1. **User** sends a message to the root agent
2. **Root agent** analyzes the request and decides which sub-agent should handle it
3. **Root agent** calls `transfer_task` with the target agent, task description, and expected output
4. **Sub-agent** processes the task in its own agentic loop using its tools
5. **Results** flow back to the root agent, which responds to the user

```bash
# The transfer_task tool call looks like:
transfer_task(
  agent="developer",
  task="Create a REST API endpoint for user authentication",
  expected_output="Working Go code with tests"
)
```

> [!NOTE]
> **Auto-Approved**
>
> Unlike other tools, `transfer_task` is always auto-approved — no user confirmation needed. This allows seamless delegation between agents.

## Handoffs Routing

Handoffs are a peer-to-peer routing pattern where agents **hand off the entire conversation** to another agent. Unlike delegation, there is no sub-session — the conversation stays in a single session and the active agent simply switches.

This pattern is ideal for:

- **Pipeline workflows** — data flows through a chain of specialized agents
- **Conversational routing** — a coordinator routes the user to the right specialist, who can route back when done
- **Graph topologies** — agents can form cycles (A → B → C → A), enabling iterative workflows

### How It Works

1. **User** sends a message to the starting agent
2. **Agent A** processes the message, then calls `handoff` to route to **Agent B**
3. **Agent B** becomes the active agent and sees the **full conversation history**
4. **Agent B** can respond, use its own tools, or hand off to another agent
5. This continues until an agent responds directly without handing off

```bash
# The handoff tool call looks like:
handoff(
  agent="summarizer"
)
```

> [!NOTE]
> **Scoped Handoff Targets**
>
> Each agent can only hand off to agents listed in its own `handoffs` array. The `handoff` tool is automatically injected — you don't need to add it manually.

### Example

A coordinator routes to a researcher, who hands off to a summarizer, who returns to the coordinator:

```text
Root ──→ Researcher ──→ Summarizer ──→ Root
```

```yaml
agents:
  root:
    model: anthropic/claude-sonnet-4-5
    description: Coordinator that routes queries
    instruction: |
      Route research queries to the researcher.
    handoffs:
      - researcher

  researcher:
    model: openai/gpt-5
    description: Web researcher
    instruction: |
      Search the web, then hand off to the summarizer.
    toolsets:
      - type: mcp
        ref: docker:duckduckgo
    handoffs:
      - summarizer

  summarizer:
    model: openai/gpt-5
    description: Summarizes findings
    instruction: |
      Summarize the research results, then hand off
      back to root.
    handoffs:
      - root
```

> [!TIP]
> **Full pipeline example**
>
> For a more complex handoff graph with branching and multiple processing stages, see [`examples/handoff.yaml`](https://github.com/docker/docker-agent/blob/main/examples/handoff.yaml).

### Forced Handoffs

With `handoffs`, the **model decides** whether to call the `handoff` tool — which means it can forget to, breaking pipelines that depend on a strict order. `force_handoff` removes that uncertainty: whenever the agent produces a final response, the **runtime itself** routes the conversation to the named agent, bypassing the LLM's tool-calling entirely. The full conversation context carries over.

```yaml
agents:
  root:
    model: anthropic/claude-sonnet-4-5
    description: Extracts key facts from the input
    instruction: |
      Extract the key facts from the user's input as a bullet list.
    force_handoff: summarizer

  summarizer:
    model: anthropic/claude-sonnet-4-5
    description: Produces the final summary
    instruction: |
      Summarize the extracted facts for the user.
```

Rules enforced at config load time:

- The target must be an agent defined in the config (or an external reference)
- An agent cannot `force_handoff` to itself
- Chains of `force_handoff` edges must not form a cycle (A → B → A is rejected)

See [`examples/force_handoff.yaml`](https://github.com/docker/docker-agent/blob/main/examples/force_handoff.yaml) for a runnable example.

## Parallel Delegation with Background Agents

`transfer_task` is **sequential** — the coordinator waits for the sub-agent to finish before continuing. When you need to fan out work to multiple agents at the same time, use the `background_agents` toolset instead.

Add it to your coordinator's toolsets:

```yaml
agents:
  root:
    model: anthropic/claude-sonnet-4-5
    description: Research coordinator
    sub_agents: [researcher, analyst, writer]
    toolsets:
      - type: think
      - type: background_agents
```

The coordinator can then:

1. **Dispatch** several tasks at once with `run_background_agent` — each returns a task ID immediately
2. **Monitor** progress with `list_background_agents` or `view_background_agent`
3. **Collect** results once tasks complete
4. **Cancel** tasks that are no longer needed with `stop_background_agent`

```bash
# Start two tasks in parallel
run_background_agent(agent="researcher", task="Find recent papers on LLM agents")
run_background_agent(agent="analyst", task="Analyze our current architecture")

# Check on all tasks
list_background_agents()

# Read results when ready
view_background_agent(task_id="agent_task_abc123")
```

## External Sub-Agents from Registries

Sub-agents don't have to be defined locally — you can reference agents from OCI registries (such as the [Docker Agent Catalog](https://hub.docker.com/u/agentcatalog)) directly in your `sub_agents` list. This lets you compose teams using pre-built, shared agents without duplicating their configuration.

```yaml
agents:
  root:
    model: openai/gpt-5
    description: Coordinator that delegates to local and catalog sub-agents
    instruction: |
      Delegate tasks to the most appropriate sub-agent.
    sub_agents:
      - local_helper
      - agentcatalog/pirate # pulled from registry automatically

  local_helper:
    model: openai/gpt-5
    description: A local helper agent for simple tasks
    instruction: You are a helpful assistant.
```

External sub-agents are automatically named after their last path segment — for example, `agentcatalog/pirate` becomes `pirate`. You can also give them an explicit name using the `name:reference` syntax:

```yaml
    sub_agents:
      - my_pirate:agentcatalog/pirate  # available as "my_pirate"
      - reviewer:docker.io/myorg/review-agent:latest
```

### Pin external sub-agents to a digest

External references use a tag by default: `agentcatalog/pirate` is shorthand for `agentcatalog/pirate:latest`. Tag references are re-resolved against the registry on **every** `docker agent run`: each unpinned external sub-agent triggers a digest lookup at startup, even when that sub-agent is never invoked in the session. On a healthy connection this typically adds a second or two per reference (it depends on your network and registry), and it is one of the paths that can stall if the registry or credential helper misbehaves.

Pinning a reference to an immutable digest (`@sha256:…`) makes the runtime serve it straight from the local cache with no network round-trip, so startup stays fast and your team is fully reproducible:

```yaml
    sub_agents:
      - reviewer:docker.io/myorg/review-agent@sha256:44117e73263afa5c861bdf3730dae7925918ffdd146827eee5bcff20bc55e8fa
```

Copy the digest from your registry (Docker Hub shows it next to the tag) or read it with a standard OCI tool such as `docker buildx imagetools inspect <reference>`. Docker Agent logs a startup warning for any external OCI reference that still uses a tag instead of a digest.

External references in `handoffs` and `force_handoff` carry the same per-run cost, so pin those to a digest too.

> [!TIP]
> External sub-agents work with any OCI-compatible registry, not just the Docker Agent Catalog. See [Agent Distribution](../distribution/index.md) for more on registry references.
>
> See [`examples/sub-agents-from-catalog.yaml`](https://github.com/docker/docker-agent/blob/main/examples/sub-agents-from-catalog.yaml) for a complete example mixing local and catalog sub-agents.

## Harness-Backed Sub-Agents

Sub-agents can be backed by external coding CLIs — Claude Code, Codex, opencode, or pi — instead of a model API. Add a `harness:` block in place of a `model:` field to create a harness sub-agent:

```yaml
agents:
  root:
    model: anthropic/claude-sonnet-4-5
    description: Orchestrator that plans and delegates
    instruction: |
      Break down coding tasks and delegate to the coding agents.
    sub_agents:
      - claude-coder
      - codex-coder

  claude-coder:
    description: Claude Code specialist
    harness:
      type: claude-code
      effort: high

  codex-coder:
    description: Codex specialist
    harness:
      type: codex
```

The orchestrator uses `transfer_task` to send work to a harness sub-agent just like any other sub-agent. Docker Agent handles the orchestration and hooks; the external CLI drives the coding loop.

> [!TIP]
> **Learn more**
>
> See [Coding Harnesses](../../features/harnesses/index.md) for the full field reference, parallel dispatch patterns, and what does not work inside harness agents.

## Example: Development Team

```yaml
agents:
  root:
    model: anthropic/claude-sonnet-4-5
    description: Technical lead coordinating development
    instruction: |
      You are a technical lead managing a development team.
      Analyze requests and delegate to the right specialist.
      Ensure quality by reviewing results before responding.
    sub_agents: [developer, reviewer, tester]
    toolsets:
      - type: think

  developer:
    model: anthropic/claude-sonnet-4-5
    description: Expert software developer
    instruction: |
      You are an expert developer. Write clean, efficient code
      and follow best practices.
    toolsets:
      - type: filesystem
      - type: shell
      - type: think

  reviewer:
    model: openai/gpt-5
    description: Code review specialist
    instruction: |
      You review code for quality, security, and maintainability.
      Provide actionable feedback.
    toolsets:
      - type: filesystem

  tester:
    model: openai/gpt-5
    description: Quality assurance engineer
    instruction: |
      You write tests and ensure software quality. Run tests
      and report results.
    toolsets:
      - type: shell
      - type: todo
```

## Example: Research Team

```yaml
agents:
  root:
    model: anthropic/claude-sonnet-4-5
    description: Research coordinator
    instruction: |
      Coordinate research tasks. Delegate web searches to
      the researcher and writing to the writer.
    sub_agents: [researcher, writer]
    toolsets:
      - type: think

  researcher:
    model: openai/gpt-5
    description: Web researcher
    instruction: Search the web and gather information.
    toolsets:
      - type: mcp
        ref: docker:duckduckgo
      - type: memory
        path: ./research.db

  writer:
    model: anthropic/claude-sonnet-4-5
    description: Content writer
    instruction: Write clear, well-structured content.
    toolsets:
      - type: filesystem
```

## Multi-Model Teams

A key advantage of multi-agent systems is using different models for different roles — picking the best model for each job:

```yaml
models:
  fast:
    provider: openai
    model: gpt-5-mini
    temperature: 0.2 # precise

  creative:
    provider: openai
    model: gpt-5
    temperature: 0.8 # creative

  local:
    provider: dmr
    model: ai/qwen3 # runs locally, no API cost

agents:
  analyst:
    model: fast # cheap and fast for analysis
  writer:
    model: creative # creative for content
  helper:
    model: local # free for simple tasks
```

## Shared Tools

Tools like `todo` can be shared between agents for collaborative task tracking:

```yaml
toolsets:
  - type: todo
    shared: true # all agents see the same todo list
```

## Best Practices

- **Keep agents focused** — Each agent should have a clear, narrow role
- **Write clear descriptions** — The coordinator uses descriptions to decide who to delegate to
- **Give minimal tools** — Only give each agent the tools it needs for its specific role
- **Use the think tool when needed** — For models without native reasoning, give coordinators the think tool so they reason about delegation. Models with built-in thinking (e.g., via `thinking_budget`) don't need it
- **Use the right model** — Use capable models for complex reasoning, cheap models for simple tasks
- **Choose the right pattern** — Use `sub_agents` for hierarchical task delegation, `handoffs` for pipeline workflows and conversational routing

> [!NOTE]
> **Beyond Docker Agent**
>
> For interoperability with other agent frameworks, Docker Agent supports the [A2A protocol](../../features/a2a/index.md) and can expose agents via [MCP Mode](../../features/mcp-mode/index.md).
