---
title: "Skills"
description: "Skills provide specialized instructions that agents can load on demand when a task matches a skill's description."
keywords: docker agent, ai agents, features, skills
weight: 110
canonical: https://docs.docker.com/ai/docker-agent/features/skills/
---

_Skills provide specialized instructions that agents can load on demand when a task matches a skill's description._

## How Skills Work

1. docker-agent scans standard directories for `SKILL.md` files
2. Skill metadata (name, description) is injected into the agent's system prompt
3. When a user request matches a skill, the agent reads the full instructions
4. The agent follows the skill's detailed instructions to complete the task

## Enabling Skills

```yaml
agents:
  root:
    model: openai/gpt-4o
    instruction: You are a helpful assistant.
    skills: true
    toolsets:
      - type: filesystem # required for reading skill files
```

> [!TIP]
> Skills are perfect for encoding team-specific workflows (PR review, deployment, coding standards) that apply across projects.

## Filtering Skills

The `skills` field also accepts a list, letting you restrict the agent to a specific subset of skills instead of exposing every discovered one. List items are classified automatically:

- `"local"` or any `http://` / `https://` URL → a **source** to load skills from
- any other string → the **name** of a skill to include

When only names are given, local sources are used by default.

```yaml
agents:
  # Load every discovered local skill (same as `skills: true`).
  full:
    skills: true

  # Load local skills, but only expose "commit" and "poem".
  scoped:
    skills:
      - commit
      - poem

  # Combine an explicit source with a name filter.
  remote_filtered:
    skills:
      - https://skills.example.com
      - commit

  # Disable skills entirely.
  none:
    skills: false
```

A name that doesn't match any discovered skill is logged as a warning at startup but is otherwise ignored.

## Inline Skills

Instead of (or alongside) loading skills from files and URLs, you can define skills directly in the agent config. An inline skill is a mapping item in the `skills` list, freely mixed with the string items above:

```yaml
agents:
  root:
    model: openai/gpt-4o
    instruction: You are a helpful assistant.
    skills:
      - name: changelog
        description: Write a concise changelog entry from a diff or description.
        instructions: |
          Produce a single changelog entry in Keep a Changelog style.
          Pick the right category (Added, Changed, Fixed, Removed) and write
          one imperative sentence summarising the user-visible change.

      # A fork-mode inline skill runs in an isolated sub-agent.
      - name: triage
        description: Triage a bug report in an isolated context.
        context: fork
        instructions: |
          Restate the problem, list likely root causes most-probable-first,
          and propose the smallest reproduction and next concrete action.

      # Inline skills mix freely with sources and name filters.
      - local
    toolsets:
      - type: filesystem
```

Inline skills carry their body in the config itself, so they need no `SKILL.md` file and require no filesystem source. They are **always exposed** — the name filter only applies to file- and URL-based skills. Because inline skills travel inside the agent YAML, they also work in `--sandbox` mode without any kit staging, and they can be shared with the agent via `share push`.

### Inline Skill Fields

| Field           | Required | Description                                                                |
| --------------- | -------- | -------------------------------------------------------------------------- |
| `name`          | Yes      | Skill identifier used by `read_skill` / `run_skill` and the `/<name>` command |
| `description`   | Yes      | Short description shown to the agent for skill matching                    |
| `instructions`  | Yes      | The skill body (what a `SKILL.md` would contain below its frontmatter)     |
| `context`       | No       | Set to `fork` to run the skill as an isolated sub-agent                    |
| `model`         | No       | Override the model used while running a fork-mode skill                    |
| `allowed_tools` | No       | For a fork-mode skill, restricts the sub-session to the parent tools whose names match an entry (glob or exact). See [Scoping a fork skill's tools](#scoping-a-fork-skills-tools). |
| `toolsets`      | No       | For a fork-mode skill, names of top-level [`toolsets`](../../configuration/overview/index.md#reusable-toolsets-toolsets) to expose in the sub-session on top of the inherited tools. |

> [!NOTE]
> **Inline vs. file-based skills**
>
> Inline skills support the subset of the SKILL.md format that fits in YAML. They cannot bundle supporting files (no `read_skill_file`) or use `` !`command` `` expansion. For skills that need bundled resources or executable helpers, use a `SKILL.md` directory instead.

## SKILL.md Format

<!-- yaml-lint:skip -->
```yaml
---
name: create-dockerfile
description: Create optimized Dockerfiles for applications
license: Apache-2.0
metadata:
  author: my-org
  version: "1.0"
---

# Creating Dockerfiles

When asked to create a Dockerfile:

1. Analyze the application type and language
2. Use multi-stage builds for compiled languages
3. Minimize image size by using slim base images
4. Follow security best practices (non-root user, etc.)
```

### Frontmatter Fields

| Field            | Required | Description                                                                 |
| ---------------- | -------- | --------------------------------------------------------------------------- |
| `name`           | Yes      | Unique skill identifier                                                     |
| `description`    | Yes      | Short description shown to the agent for skill matching                     |
| `context`        | No       | Set to `fork` to run the skill as an isolated sub-agent (see below)         |
| `model`          | No       | Override the model used while running the skill as a sub-agent (fork only)  |
| `allowed-tools`  | No       | For a fork-mode skill, restricts the sub-session to the parent tools whose names match an entry (YAML list or comma-separated string). See [Scoping a fork skill's tools](#scoping-a-fork-skills-tools). |
| `toolsets`       | No       | For a fork-mode skill, names of top-level [`toolsets`](../../configuration/overview/index.md#reusable-toolsets-toolsets) to expose in the sub-session (YAML list or comma-separated string). |
| `license`        | No       | License identifier (e.g. `Apache-2.0`)                                      |
| `compatibility`  | No       | Free-text compatibility notes                                               |
| `metadata`       | No       | Arbitrary key-value pairs (e.g. `author`, `version`)                        |

## Running a Skill as a Sub-Agent

By default, when an agent invokes a skill it reads the instructions inline into its own conversation. For complex, multi-step skills this can consume a large portion of the agent's context window and pollute the parent conversation with intermediate tool calls.

Adding `context: fork` to the SKILL.md frontmatter tells the agent to run the skill in an **isolated sub-agent** instead:

<!-- yaml-lint:skip -->
```yaml
---
name: bump-go-dependencies
description: Update Go module dependencies one by one
context: fork
---

# Bump Dependencies

1. List outdated deps
2. Update each one, run tests, commit or revert
3. Produce a summary table
```

When the agent encounters a task that matches a `context: fork` skill, it uses the `run_skill` tool instead of `read_skill`. This:

- **Spawns a child session** with the skill content as the system prompt and the caller's task as the user message
- **Isolates the context window** — the sub-agent has its own conversation history, so lengthy tool-call chains don't eat into the parent's token budget
- **Folds the result** — only the sub-agent's final answer is returned to the parent as the tool result
- **Inherits the parent's model and tools** — the sub-agent can use all tools available to the parent agent (scope this with `allowed_tools` / `toolsets`, see [Scoping a fork skill's tools](#scoping-a-fork-skills-tools))

> [!TIP]
> **When to use context: fork**
>
> Use `context: fork` for skills that involve many steps, heavy tool usage, or that should not clutter the main conversation — for example dependency bumping, large refactors, or code generation pipelines.

### Overriding the model for a fork skill

Fork skills can declare a `model` field in their frontmatter to use a
different model than the parent agent for the duration of the sub-session.
This is useful when a skill is best handled by a faster, cheaper, or more
specialised model — for example a powerful reasoning model for refactors,
or a fast model for routine bookkeeping work. The override only applies
while the skill is running; the parent agent keeps its own model.

The `model` value accepts either a named model from the agent config or
an inline `provider/model` reference (and the same comma-separated alloy
syntax as the rest of the agent config):

<!-- yaml-lint:skip -->
```yaml
---
name: bump-go-dependencies
description: Update Go module dependencies one by one
context: fork
model: openai/gpt-4o-mini
---

# Bump Dependencies

1. ...
```

If the model reference cannot be resolved (unknown name, missing
credentials, runtime not configured for model switching, …) the skill
falls back to the agent's currently-active model (its configured
default, or any override the user previously set via the model picker)
and a warning is logged.

When the skill completes, the agent's previous model is restored — but
only if no one else changed the model in the meantime. If the user
switches the model via the TUI model picker while the fork skill is
running, their choice is preserved (the deferred restore becomes a
no-op).

### Scoping a fork skill's tools

By default a fork skill inherits the parent agent's entire tool set. Two
optional fields let you scope what the sub-session can use. Both apply
**only to fork-mode skills** and work the same whether the skill is
inline or loaded from a `SKILL.md` file.

`allowed_tools` (frontmatter: `allowed-tools`) is an **allow-list** over
the inherited tools: only tools whose names match an entry are kept,
everything else is hidden from the sub-session. Entries support glob
patterns (e.g. `read_*`) and otherwise match exactly. This is the
Claude-Code-compatible `allowed-tools` field, now enforced for fork
skills rather than merely recorded.

`toolsets` references reusable [top-level toolsets](../../configuration/overview/index.md#reusable-toolsets-toolsets)
by name. The referenced toolsets are exposed in the sub-session **in
addition to** the inherited tools, and they bypass the `allowed_tools`
filter (the skill explicitly asked for them).

```yaml
toolsets:
  web:
    type: fetch

agents:
  root:
    model: openai/gpt-4o
    instruction: You are a helpful assistant.
    toolsets:
      - type: filesystem
      - type: shell
    skills:
      # Inherits the parent tools but is restricted to read-only filesystem
      # access while it runs — shell and write tools are hidden.
      - name: audit
        description: Review the repository layout without modifying anything.
        context: fork
        allowed_tools:
          - read_file
          - list_directory
          - directory_tree
        instructions: Inspect the repository structure and summarise it.

      # Brings in the top-level `web` toolset on top of the parent's tools.
      - name: research
        description: Research a topic using web fetches in an isolated context.
        context: fork
        toolsets:
          - web
        instructions: Research the requested topic and summarise with links.
```

The equivalent in a `SKILL.md` file uses frontmatter lists:

<!-- yaml-lint:skip -->
```yaml
---
name: research
description: Research a topic using web fetches
context: fork
allowed-tools:
  - fetch
toolsets:
  - web
---
```

> [!NOTE]
> **Fork only**
>
> Both fields are rejected by config validation when set on a non-fork skill, and a `toolsets` entry that doesn't resolve to a top-level toolset is a load-time error.

## Search Paths

Skills are discovered from these locations (later overrides earlier):

### Global

| Path                | Search Type                             |
| ------------------- | --------------------------------------- |
| `~/.codex/skills/`  | Recursive (searches all subdirectories) |
| `~/.claude/skills/` | Flat (immediate children only)          |
| `~/.agents/skills/` | Recursive (searches all subdirectories) |

### Project (from git root to current directory)

| Path              | Search Type                                |
| ----------------- | ------------------------------------------ |
| `.claude/skills/` | Flat (cwd only)                            |
| `.agents/skills/` | Flat (each directory from git root to cwd) |

## Invoking Skills

Skills can be invoked in multiple ways:

- **Automatic:** The agent detects when your request matches a skill's description and loads it automatically
- **Explicit:** Reference the skill name in your prompt: "Use the create-dockerfile skill to..."
- **Slash command:** Use `/{skill-name}` to invoke a skill directly

```bash
# In the TUI, invoke skill directly:
/create-dockerfile

# Or mention it in your message:
"Create a dockerfile for my Python app (use the create-dockerfile skill)"
```

## Precedence

When multiple skills share the same name:

1. Global skills load first
2. Project skills load next, from git root toward current directory
3. Skills closer to the current directory override those further away

## Skills in Sandbox Mode

When you run an agent with [`--sandbox`](../../configuration/sandbox/index.md), the sandbox VM has its own filesystem with no access to your host's skill directories. docker-agent handles this transparently via the [auto-kit](../../configuration/sandbox/index.md#auto-kit): every discovered local skill is staged into a per-agent kit on the host, run through best-effort secret redaction (see the [auto-kit](../../configuration/sandbox/index.md#secret-redaction) docs), and bind-mounted read-only into the sandbox so the agent sees the same skills inside the VM as on the host. No configuration is required — use `--no-kit` only if you explicitly want to run the sandbox without any host skills.

## Creating a Skill

```bash
# Create the skill directory
$ mkdir -p ~/.agents/skills/create-dockerfile

# Write the SKILL.md file
$ cat > ~/.agents/skills/create-dockerfile/SKILL.md << 'EOF'
---
name: create-dockerfile
description: Create optimized Dockerfiles for applications
---

# Creating Dockerfiles

When asked to create a Dockerfile:

1. Analyze the application type and language
2. Use multi-stage builds for compiled languages
3. Use slim base images to minimize size
4. Run as non-root user for security
EOF
```

The skill will automatically be available to any agent with skills enabled (`skills: true`, or a list that targets its name — see [Filtering Skills](#filtering-skills)).

> [!NOTE]
> **See also**
>
> Skills are enabled in the [Agent Config](../../configuration/agents/index.md) with the `skills` property (boolean or list). For tool-based capabilities, see [Tools](../../concepts/tools/index.md).
>
> Example configs: [`examples/skills_inline.yaml`](https://github.com/docker/docker-agent/blob/main/examples/skills_inline.yaml) (inline skill definition), [`examples/skills_fork_toolsets.yaml`](https://github.com/docker/docker-agent/blob/main/examples/skills_fork_toolsets.yaml) (scoping a fork skill's tools), [`examples/skills_filter.yaml`](https://github.com/docker/docker-agent/blob/main/examples/skills_filter.yaml) (filtering which skills load).
