---
layout: default
title: "Docker Agent"
description: "Run AI agents like containers. Define them in YAML, share them through any OCI registry, and run them anywhere."
permalink: /
---

<div class="hero">
  <h1>Docker Agent</h1>
  <p><strong>Docker Agent is to AI agents what <code>docker run</code> is to containers.</strong><br>Define an agent in a YAML file, run it with one command, share it through any OCI registry — the same workflow you already use for images.</p>
  <div class="hero-buttons">
    <a href="{{ '/getting-started/quickstart/' | relative_url }}" class="btn btn-primary">Quick Start →</a>
    <a href="https://github.com/docker/docker-agent" target="_blank" rel="noopener noreferrer" class="btn btn-secondary">View on GitHub</a>
  </div>
</div>

<div class="demo-container">
  <img src="{{ '/demo.gif' | relative_url }}" alt="Docker Agent TUI demo showing an interactive agent session" loading="lazy">
</div>

<div class="elevator">
  <div class="elevator-card">
    <div class="elevator-label">What it is</div>
    <p>A CLI that runs AI agents defined declaratively in YAML or HCL.</p>
  </div>
  <div class="elevator-card">
    <div class="elevator-label">What it isn’t</div>
    <p>Not a framework you write code in. Not a hosted SaaS. Not a new model.</p>
  </div>
  <div class="elevator-card">
    <div class="elevator-label">Who it’s for</div>
    <p>Developers who want agents in their workflow without glue code.</p>
  </div>
  <div class="elevator-card">
    <div class="elevator-label">What you get</div>
    <p>TUI · CLI · HTTP API · MCP server · A2A · OCI distribution.</p>
  </div>
</div>

## What Is Docker Agent?

Docker Agent is an open-source tool from **Docker** — makers of [Docker Engine](https://www.docker.com/products/container-runtime/), [Docker Desktop](https://www.docker.com/products/docker-desktop/), [Docker Hub](https://hub.docker.com), and [Docker Scout](https://www.docker.com/products/docker-scout/) — that lets you **build, run, and share AI agents using simple configuration files** instead of writing application code.

You describe what your agent does — its model, personality, tools, and teammates — in a YAML file. Docker Agent handles the LLM orchestration loop, tool execution, multi-agent delegation, and streaming output. You focus on *what* the agent should do, not *how* to wire it up.

```yaml
# agent.yaml — this is all you need
agents:
  root:
    model: anthropic/claude-sonnet-4-5
    description: A coding assistant
    instruction: |
      You are an expert developer. Help users write clean,
      efficient code. Explain your reasoning step by step.
    toolsets:
      - type: filesystem
      - type: shell
      - type: think
```

```bash
$ docker agent run agent.yaml
```

That's it. Your agent can now read and write files, run shell commands, and reason through problems — all through an interactive terminal UI.

## Without vs. with Docker Agent

The same coding assistant, written two different ways:

<div class="compare" markdown="1">
  <div class="compare-side compare-without" markdown="1">
    <div class="compare-label">Without Docker Agent</div>

```python
# ~30 lines of glue code, every project
import anthropic, json, subprocess
from pathlib import Path

client = anthropic.Anthropic()
MODEL = "claude-sonnet-4-5"
TOOLS = [
  {"name": "read_file", "input_schema": {...}},
  {"name": "write_file", "input_schema": {...}},
  {"name": "run_shell", "input_schema": {...}},
]

def dispatch(name, args):
  if name == "read_file":
      return Path(args["path"]).read_text()
  if name == "write_file":
      Path(args["path"]).write_text(args["content"])
      return "ok"
  if name == "run_shell":
      return subprocess.check_output(args["cmd"], shell=True).decode()

messages = [{"role": "user", "content": input("> ")}]
while True:
  resp = client.messages.create(
    model=MODEL, max_tokens=4096, tools=TOOLS,
    system="You are an expert developer…",
    messages=messages,
  )
  # …parse tool_use blocks, dispatch, append, loop…
  if resp.stop_reason == "end_turn": break
```

  </div>
  <div class="compare-side compare-with" markdown="1">
    <div class="compare-label">With Docker Agent</div>

```yaml
# agent.yaml — 8 lines, no glue code
agents:
  root:
    model: anthropic/claude-sonnet-4-5
    description: A coding assistant
    instruction: You are an expert developer.
    toolsets:
      - type: filesystem
      - type: shell
```

```bash
$ docker agent run agent.yaml
```

  </div>
</div>

## Why Docker Agent?

Most AI agent frameworks ask you to write Python or TypeScript to glue together models, tools, and workflows. Docker Agent takes a different approach: **declare everything in config, run it with a single command.**

<div class="pain-grid">
  <div class="pain-row">
    <div class="pain-pain"><span class="pain-x">×</span> “I rebuilt the same agent loop in three projects.”</div>
    <div class="pain-fix"><span class="pain-check">✓</span> Reusable YAML — declare once, run everywhere.</div>
  </div>
  <div class="pain-row">
    <div class="pain-pain"><span class="pain-x">×</span> “Sharing my agent means a repo plus a setup README.”</div>
    <div class="pain-fix"><span class="pain-check">✓</span> <code>docker agent run user/agent</code> — OCI distribution, like images.</div>
  </div>
  <div class="pain-row">
    <div class="pain-pain"><span class="pain-x">×</span> “I’m locked into one model SDK.”</div>
    <div class="pain-fix"><span class="pain-check">✓</span> Swap the <code>model:</code> line — OpenAI, Anthropic, Gemini, Bedrock, local.</div>
  </div>
  <div class="pain-row">
    <div class="pain-pain"><span class="pain-x">×</span> “Tools are one-offs glued to one agent.”</div>
    <div class="pain-fix"><span class="pain-check">✓</span> Built-in toolsets plus any MCP server from <a href="https://hub.docker.com/u/mcp" target="_blank" rel="noopener">Docker's MCP catalog</a> — reuse them across agents.</div>
  </div>
</div>

<div class="features-grid">
  <div class="feature">
    <div class="feature-icon">
      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.6" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><path d="M14.5 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V7.5L14.5 2z"/><polyline points="14 2 14 8 20 8"/><line x1="9" y1="13" x2="15" y2="13"/><line x1="9" y1="17" x2="15" y2="17"/></svg>
    </div>
    <h3>Config, Not Code</h3>
    <p>Define agents in YAML or HCL. Swap models, add tools, or change behavior without touching application code.</p>
  </div>
  <div class="feature">
    <div class="feature-icon">
      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.6" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><path d="M14.7 6.3a1 1 0 0 0 0 1.4l1.6 1.6a1 1 0 0 0 1.4 0l3.77-3.77a6 6 0 0 1-7.94 7.94l-6.91 6.91a2.12 2.12 0 0 1-3-3l6.91-6.91a6 6 0 0 1 7.94-7.94l-3.76 3.76z"/></svg>
    </div>
    <h3>Built-in Tools + MCP</h3>
    <p>Comes with tools for filesystem, shell, memory, web fetch, and more. Extend with any of the hundreds of MCP servers in <a href="https://hub.docker.com/u/mcp" target="_blank" rel="noopener">Docker's MCP catalog</a>.</p>
  </div>
  <div class="feature">
    <div class="feature-icon">
      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.6" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/><circle cx="9" cy="7" r="4"/><path d="M23 21v-2a4 4 0 0 0-3-3.87"/><path d="M16 3.13a4 4 0 0 1 0 7.75"/></svg>
    </div>
    <h3>Multi-Agent Teams</h3>
    <p>Build teams of specialized agents that delegate work to each other. A coordinator routes tasks to the right specialist.</p>
  </div>
  <div class="feature">
    <div class="feature-icon">
      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.6" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><path d="M9.5 2A2.5 2.5 0 0 1 12 4.5v15a2.5 2.5 0 0 1-4.96.44 2.5 2.5 0 0 1-2.96-3.08 3 3 0 0 1-.34-5.58 2.5 2.5 0 0 1 1.32-4.24 2.5 2.5 0 0 1 1.98-3A2.5 2.5 0 0 1 9.5 2z"/><path d="M14.5 2A2.5 2.5 0 0 0 12 4.5v15a2.5 2.5 0 0 0 4.96.44 2.5 2.5 0 0 0 2.96-3.08 3 3 0 0 0 .34-5.58 2.5 2.5 0 0 0-1.32-4.24 2.5 2.5 0 0 0-1.98-3A2.5 2.5 0 0 0 14.5 2z"/></svg>
    </div>
    <h3>Any Model</h3>
    <p>OpenAI, Anthropic, Google Gemini, AWS Bedrock, local models via Docker Model Runner or Ollama — bring your own provider.</p>
  </div>
  <div class="feature">
    <div class="feature-icon">
      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.6" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><path d="M21 16V8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73l7 4a2 2 0 0 0 2 0l7-4A2 2 0 0 0 21 16z"/><polyline points="3.27 6.96 12 12.01 20.73 6.96"/><line x1="12" y1="22.08" x2="12" y2="12"/></svg>
    </div>
    <h3>Package &amp; Share Like Images</h3>
    <p>Push agents to any OCI registry. Pull and run them anywhere with one command — the same workflow you use for containers.</p>
  </div>
  <div class="feature">
    <div class="feature-icon">
      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.6" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><polyline points="4 17 10 11 4 5"/><line x1="12" y1="19" x2="20" y2="19"/></svg>
    </div>
    <h3>Run Anywhere</h3>
    <p>Interactive TUI, headless CLI, HTTP API server, OpenAI-compatible chat endpoint, MCP server, or A2A protocol.</p>
  </div>
</div>

## Use Cases

What people build with Docker Agent today:

<div class="usecase-grid">
  <div class="usecase">
    <div class="usecase-icon">
      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.6" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><polyline points="16 18 22 12 16 6"/><polyline points="8 6 2 12 8 18"/></svg>
    </div>
    <h3>Coding agents</h3>
    <p>Pair-programmer agents with file system, shell, and LSP tools. Read code, edit it, run tests, iterate.</p>
    <a href="https://hub.docker.com/u/agentcatalog" target="_blank" rel="noopener">Browse the catalog →</a>
  </div>
  <div class="usecase">
    <div class="usecase-icon">
      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.6" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><polyline points="4 17 10 11 4 5"/><line x1="12" y1="19" x2="20" y2="19"/></svg>
    </div>
    <h3>Ops &amp; SRE</h3>
    <p>Triage incidents, search logs, run kubectl, build Dockerfiles. Pipe alerts in via <code>--exec</code> for headless runs.</p>
    <a href="{{ '/features/cli/' | relative_url }}">CLI reference →</a>
  </div>
  <div class="usecase">
    <div class="usecase-icon">
      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.6" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><line x1="18" y1="20" x2="18" y2="10"/><line x1="12" y1="20" x2="12" y2="4"/><line x1="6" y1="20" x2="6" y2="14"/><line x1="4" y1="20" x2="20" y2="20"/></svg>
    </div>
    <h3>Data &amp; research</h3>
    <p>Persistent memory, web fetch, RAG over local docs, structured output for downstream pipelines.</p>
    <a href="{{ '/tools/rag/' | relative_url }}">RAG guide →</a>
  </div>
  <div class="usecase">
    <div class="usecase-icon">
      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.6" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><circle cx="12" cy="12" r="10"/><polygon points="16.24 7.76 14.12 14.12 7.76 16.24 9.88 9.88 16.24 7.76"/></svg>
    </div>
    <h3>Custom workflows</h3>
    <p>Multi-agent teams, hooks, model routing, A2A and MCP servers — wire agents into your existing stack.</p>
    <a href="{{ '/concepts/multi-agent/' | relative_url }}">Multi-agent →</a>
  </div>
</div>

## How It Works

Docker Agent follows a simple loop:

<figure class="flow-diagram">
  <img src="{{ '/assets/how-it-works.svg' | relative_url }}" alt="agent.yaml is run by 'docker agent run', which loops through Model, Tools and Sub-agents, then streams results to the TUI or API." loading="lazy">
  <figcaption>Your YAML config is the input; the runtime drives a Model ↔ Tools ↔ Sub-agents loop until the task is done; results stream back to the TUI or any API client.</figcaption>
</figure>

1. **You define an agent** in YAML — its model, instructions, tools, and sub-agents
2. **You run it** with `docker agent run` via TUI, CLI, or API
3. **The agent processes your request** — calling tools, delegating to sub-agents, reasoning step by step
4. **Results stream back** in real time

## A few terms you'll see

<dl class="glossary">
  <dt>Agent</dt>
  <dd>An LLM with instructions, tools, and (optionally) sub-agents — the unit you define in YAML.</dd>

  <dt>Tool</dt>
  <dd>A function the agent can call, like <code>read_file</code> or <code>shell</code>. Tools come from built-in toolsets or external MCP servers.</dd>

  <dt>MCP</dt>
  <dd><em>Model Context Protocol</em> — an open standard for tool servers. Docker Agent can use any MCP server as a toolset.</dd>

  <dt>A2A</dt>
  <dd><em>Agent-to-Agent</em> — an HTTP protocol agents use to talk to each other across machines.</dd>

  <dt>TUI</dt>
  <dd><em>Terminal User Interface</em> — the default interactive front end, launched by <code>docker agent run</code>.</dd>

  <dt>OCI</dt>
  <dd>The same registry format used for Docker images. Docker Agent reuses it to push and pull agents.</dd>
</dl>

### Zero Config

The fastest way to try it — no config file needed:

```bash
# Run the built-in default agent
$ docker agent run
```

### From the Registry

Run pre-built agents from the [agent catalog on **Docker Hub**](https://hub.docker.com/u/agentcatalog) — just like pulling a Docker image:

```bash
# A pirate-themed assistant
$ docker agent run agentcatalog/pirate

# A coding agent
$ docker agent run agentcatalog/coder
```

### Multi-Agent Teams

Build a team where a coordinator delegates tasks to specialists:

```yaml
agents:
  root:
    model: openai/gpt-5
    description: Team coordinator
    instruction: Route tasks to the best specialist.
    sub_agents: [coder, reviewer]

  coder:
    model: anthropic/claude-sonnet-4-5
    description: Writes and modifies code
    instruction: Write clean, tested code.
    toolsets:
      - type: filesystem
      - type: shell

  reviewer:
    model: anthropic/claude-sonnet-4-5
    description: Reviews code for quality
    instruction: Review code for bugs, style, and best practices.
    toolsets:
      - type: filesystem
```

### Non-Interactive Mode

Use `--exec` for scripting and automation:

```bash
# One-shot task
$ docker agent run --exec agent.yaml "Create a Dockerfile for a Node.js app"

# Pipe input
$ cat error.log | docker agent run --exec agent.yaml "What's wrong in this log?"

# Serve as an API
$ docker agent serve api agent.yaml --listen :8080
```

<div class="callout callout-tip" markdown="1">
<div class="callout-title">Prefer HCL?
</div>
  <p>You can also write agent configs in HCL using labeled blocks and heredocs. See <a href="{{ '/configuration/hcl/' | relative_url }}">HCL Configuration</a>.</p>
</div>

## Part of the Docker ecosystem

Docker Agent reuses the tooling and conventions you already know:

<div class="ecosystem">
  <a class="ecosystem-tile" href="https://hub.docker.com/u/agentcatalog" target="_blank" rel="noopener">
    <strong>Docker Hub</strong>
    <span>Pull pre-built agents from the agent catalog — same registry, same auth.</span>
  </a>
  <a class="ecosystem-tile" href="https://www.docker.com/products/docker-desktop/" target="_blank" rel="noopener">
    <strong>Docker Desktop</strong>
    <span>Run MCP toolsets in containers via <code>ref: docker:…</code> with one click.</span>
  </a>
  <a class="ecosystem-tile" href="https://docs.docker.com/desktop/features/model-runner/" target="_blank" rel="noopener">
    <strong>Docker Model Runner</strong>
    <span>Run local OSS models on your machine — just point your agent at <code>dmr/…</code>.</span>
  </a>
  <a class="ecosystem-tile" href="https://www.docker.com/products/docker-scout/" target="_blank" rel="noopener">
    <strong>Docker Scout</strong>
    <span>Same supply-chain visibility you have for images extends to agent images.</span>
  </a>
</div>

## Explore the Docs

<div class="cards">
  <a class="card" href="{{ '/getting-started/introduction/' | relative_url }}">
    <div class="card-icon">
      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.6" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><path d="M4.5 16.5c-1.5 1.26-2 5-2 5s3.74-.5 5-2c.71-.84.7-2.13-.09-2.91a2.18 2.18 0 0 0-2.91-.09z"/><path d="M12 15l-3-3a22 22 0 0 1 2-3.95A12.88 12.88 0 0 1 22 2c0 2.72-.78 7.5-6 11a22.35 22.35 0 0 1-4 2z"/><path d="M9 12H4s.55-3.03 2-4c1.62-1.08 5 0 5 0"/><path d="M12 15v5s3.03-.55 4-2c1.08-1.62 0-5 0-5"/></svg>
    </div>
    <h3>Introduction</h3>
    <p>The full story: what Docker Agent is, why it exists, and how it works.</p>
  </a>
  <a class="card" href="{{ '/getting-started/quickstart/' | relative_url }}">
    <div class="card-icon">
      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.6" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><polygon points="13 2 3 14 12 14 11 22 21 10 12 10 13 2"/></svg>
    </div>
    <h3>Quick Start</h3>
    <p>Get your first agent running in under 5 minutes.</p>
  </a>
  <a class="card" href="{{ '/concepts/agents/' | relative_url }}">
    <div class="card-icon">
      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.6" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><path d="M9 18h6"/><path d="M10 22h4"/><path d="M12 2a7 7 0 0 0-4 12.7c.6.4 1 1 1 1.7V18h6v-1.6c0-.7.4-1.3 1-1.7A7 7 0 0 0 12 2z"/></svg>
    </div>
    <h3>Core Concepts</h3>
    <p>Agents, models, tools, and multi-agent orchestration explained.</p>
  </a>
  <a class="card" href="{{ '/configuration/overview/' | relative_url }}">
    <div class="card-icon">
      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.6" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><circle cx="12" cy="12" r="3"/><path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 1 1-2.83 2.83l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 1 1-4 0v-.09a1.65 1.65 0 0 0-1-1.51 1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 1 1-2.83-2.83l.06-.06a1.65 1.65 0 0 0 .33-1.82 1.65 1.65 0 0 0-1.51-1H3a2 2 0 1 1 0-4h.09a1.65 1.65 0 0 0 1.51-1 1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 1 1 2.83-2.83l.06.06a1.65 1.65 0 0 0 1.82.33h0a1.65 1.65 0 0 0 1-1.51V3a2 2 0 1 1 4 0v.09a1.65 1.65 0 0 0 1 1.51h0a1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 1 1 2.83 2.83l-.06.06a1.65 1.65 0 0 0-.33 1.82v0a1.65 1.65 0 0 0 1.51 1H21a2 2 0 1 1 0 4h-.09a1.65 1.65 0 0 0-1.51 1z"/></svg>
    </div>
    <h3>Configuration</h3>
    <p>Full reference for every YAML and HCL option.</p>
  </a>
  <a class="card" href="{{ '/providers/overview/' | relative_url }}">
    <div class="card-icon">
      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.6" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><path d="M9.5 2A2.5 2.5 0 0 1 12 4.5v15a2.5 2.5 0 0 1-4.96.44 2.5 2.5 0 0 1-2.96-3.08 3 3 0 0 1-.34-5.58 2.5 2.5 0 0 1 1.32-4.24 2.5 2.5 0 0 1 1.98-3A2.5 2.5 0 0 1 9.5 2z"/><path d="M14.5 2A2.5 2.5 0 0 0 12 4.5v15a2.5 2.5 0 0 0 4.96.44 2.5 2.5 0 0 0 2.96-3.08 3 3 0 0 0 .34-5.58 2.5 2.5 0 0 0-1.32-4.24 2.5 2.5 0 0 0-1.98-3A2.5 2.5 0 0 0 14.5 2z"/></svg>
    </div>
    <h3>Model Providers</h3>
    <p>OpenAI, Anthropic, Gemini, Bedrock, Docker Model Runner, and more.</p>
  </a>
  <a class="card" href="{{ '/features/tui/' | relative_url }}">
    <div class="card-icon">
      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.6" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2"/></svg>
    </div>
    <h3>Features</h3>
    <p>TUI, CLI, API server, MCP mode, A2A, RAG, Skills, and distribution.</p>
  </a>
</div>
