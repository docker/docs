---
title: "Background Agents Tool"
description: "Dispatch work to sub-agents concurrently and collect results asynchronously."
keywords: docker agent, ai agents, tools, toolsets, background agents tool
linkTitle: "Background Agents"
weight: 90
canonical: https://docs.docker.com/ai/docker-agent/tools/background-agents/
---

_Dispatch work to sub-agents concurrently and collect results asynchronously._

## Overview

The background agents tool lets an orchestrator dispatch work to sub-agents concurrently and collect results asynchronously. Unlike [transfer_task](../transfer-task/index.md) (which blocks until the sub-agent finishes), background agent tasks run in parallel — the orchestrator can start several tasks, do other work, and check on them later.

## Available Tools

| Tool                     | Description                                                     |
| ------------------------ | --------------------------------------------------------------- |
| `run_background_agent`   | Start a sub-agent task in the background; returns a task ID     |
| `list_background_agents` | List all background tasks with their status and runtime         |
| `view_background_agent`  | View live output or final result of a task by ID                |
| `stop_background_agent`  | Cancel a running task by ID                                     |
| `wait_background_agents` | Wait for specified tasks to finish and collect their results    |

### `run_background_agent` parameters

| Parameter         | Type   | Required | Description                                                                 |
| ----------------- | ------ | -------- | --------------------------------------------------------------------------- |
| `agent`           | string | ✓        | Name of the sub-agent to run. Must be listed under the caller's `sub_agents`. |
| `task`            | string | ✓        | Clear, concise description of the task the sub-agent should achieve.        |
| `expected_output` | string | ✗        | Optional description of the result format the caller expects.               |
| `worktree`        | bool   | ✗        | Create an isolated Git checkout and branch for this task. Default `false`. |

`run_background_agent` returns a **task ID** string. Tools run by the sub-agent inherit the parent session's permissions. Because background tasks run non-interactively, any tool call that would normally prompt the user for approval will be automatically denied. To allow background agents to run mutating tools, you must explicitly approve them in the parent session (e.g. via YOLO mode or explicit allow rules).

Background delegation shares the same runtime guards as `transfer_task`: delegation cycles are rejected and chains are capped at 10 nested delegations. See [Delegation Limits](../transfer-task/index.md#delegation-limits).

### `view_background_agent` and `stop_background_agent` parameters

| Parameter | Type   | Required | Description                                                    |
| --------- | ------ | -------- | -------------------------------------------------------------- |
| `task_id` | string | ✓        | Task ID returned by `run_background_agent` or `list_background_agents`. |

`list_background_agents` takes no parameters.

### `wait_background_agents` parameters

| Parameter  | Type     | Required | Description                                                              |
| ---------- | -------- | -------- | ------------------------------------------------------------------------ |
| `task_ids` | string[] | ✓        | 1–100 distinct task IDs returned by `run_background_agent`.               |
| `timeout`  | integer  | ✗        | Maximum seconds to wait for the whole group. Default `300`, maximum `3600`; `0` uses the default. |

Start all independent tasks first, then call `wait_background_agents` with their IDs:

```json
{"task_ids": ["agent_task_a", "agent_task_b"], "timeout": 300}
```

The tool waits for **all specified tasks**, including failed or stopped tasks, without polling. It returns JSON with:

- `all_done`: every requested task's execution has exited. This does **not** mean every task succeeded: check each task's `status` before continuing dependent work.
- `timed_out`: present and `true` when the shared timeout expires.
- `tasks`: results in the same order as `task_ids`, each containing `task_id`, `status`, and `done`, plus `output` or `error` when available. Output/error text is capped at 4 KiB per task; `truncated: true` means to use `view_background_agent` for the full result.

A timeout returns available results and live output without stopping unfinished tasks. Call the tool again to continue waiting; already-finished tasks return their cached results immediately. Cancelling the wait also leaves the tasks running and returns a tool error with available task results. Runtime shutdown still stops background tasks as usual.

A stopped task can have `done: false` while its execution is still shutting down. The join waits for execution to exit, not just for the stop request to be acknowledged. Unknown or pruned IDs are reported individually as `status: "not_found"`, leave `all_done` false, and make the result a tool error; other requested tasks are still joined. Completed tasks may be pruned when the task-history limit is reached, so collect their results before dispatching many more tasks.

## Isolated coding tasks

Set `worktree: true` to run concurrent coding tasks without sharing a checkout:

```json
{"agent": "coder", "task": "Implement the backend change and commit it", "worktree": true}
```

This works when the calling session is in a Git repository **or an existing Git worktree**. Each task gets a fresh branch from that checkout's current `HEAD`, not the repository's default branch. Uncommitted and untracked files are **not copied**; commit shared starting changes before dispatching tasks. A missing repository or a failed setup is reported as a failed task.

The child gets fresh toolsets and hooks rooted in its checkout; the parent's tools and process working directory are unchanged. Blocking delegations (`transfer_task` and forked skills) inherit the isolated workspace. Nested background dispatch is rejected: start parallel worktree tasks from the top-level orchestrator. Model clients are reused, and tasks share the parent's run budget and inherit its safety policy and permissions. Configured `worktree_create` hooks run before tools and skills are loaded, so setup-created ignore files and directories take effect.

`list_background_agents` and `view_background_agent` report the checkout path, branch, and base commit once creation finishes. `wait_background_agents` includes them as `worktree`, `branch`, and `base_commit`, independently of truncated output and whether the task succeeded.

### Review and merge at the top level

1. Start all independent tasks with `worktree: true`. Describe files relative to the repository. The child is instructed to translate source-checkout paths to its own checkout; the child is instructed to treat attached files as read-only inputs at their original paths.
2. Have each child commit its changes on its own branch, following the project's signing policy. No commits are made automatically.
3. Wait for all tasks, check each status, and review each branch's diff.
4. Merge branches **one at a time** from the parent checkout, resolve conflicts, and run the combined tests. For example:

   ```sh
   git diff HEAD...worktree-agent-<id>
   git merge --no-edit worktree-agent-<id>
   ```

5. After integrating the work, remove the now-unneeded checkout and branch:

   ```sh
   git worktree remove /path/reported/by/the/task
   git branch -d worktree-agent-<id>
   ```

Worktrees and branches are **retained**, including after failure, cancellation, and runtime shutdown. Uncommitted changes cannot be merged with `git merge`; inspect and commit them first. Never remove a task's checkout while it or its descendants are still running. Save the reported paths before task history is pruned.

> [!NOTE]
> Worktrees isolate checkouts, not permissions or arbitrary external side effects. Explicit absolute tool/hook directories, external services, and shared Git configuration remain shared. External harness agents do not yet support `worktree: true`; their drivers currently inherit the process directory. Programmatic agents need a working-directory factory; configuration-loaded native agents supply one automatically. Adding tools after loading invalidates that factory unless the host supplies a replacement. Shared plans, shared todos, and default persistent memory remain shared intentionally. RAG indexes are scoped to the workspace unless an explicit database path selects shared storage.

## Configuration

```yaml
toolsets:
  - type: background_agents
```

No configuration options. Requires the agent to have `sub_agents` configured so the background tasks have agents to dispatch to.

## Example

```yaml
agents:
  coordinator:
    model: openai/gpt-4o
    description: Orchestrates parallel research
    instruction: Fan out research tasks and synthesize results.
    sub_agents: [researcher]
    toolsets:
      - type: background_agents
      - type: think

  researcher:
    model: openai/gpt-4o
    description: Web researcher
    instruction: Research topics thoroughly.
    toolsets:
      - type: mcp
        ref: docker:duckduckgo
```

> [!TIP]
> **When to Use**
>
> Use `background_agents` when your orchestrator needs to fan out work to multiple specialists in parallel — for example, researching several topics simultaneously or running independent code analyses side by side.

The TUI sidebar's Agents section shows **Background: N running** while background tasks are executing. The count also appears in the collapsed sidebar, counts separate tasks even when they use the same agent, and remains visible after the foreground agent finishes. Tasks leave the count once execution and cleanup finish, including failed or cancelled tasks. Hiding the Agents section hides the count.

For remote sessions, the count tracks task lifecycle events received while connected; tasks already running before the connection or events lost during recovery may not be reflected. Use `list_background_agents` for the authoritative task list.

In the TUI, each background task's token usage is accounted for live: the sidebar's Agents panel shows the sub-agent's context usage percentage on its roster row, the Agent Inspector shows its exact token counts, and the task's cost joins the session total.

## Using Harness Sub-Agents

Background agents work equally well with [harness-backed sub-agents](../../features/harnesses/index.md) — sub-agents driven by external coding CLIs such as Claude Code or Codex. This lets you dispatch multiple independent coding tasks in parallel:

```yaml
agents:
  root:
    model: anthropic/claude-sonnet-4-5
    description: Orchestrator that fans out coding tasks
    instruction: |
      Dispatch the frontend and backend tasks in parallel,
      then collect results and produce a summary.
    sub_agents:
      - claude-coder
      - codex-coder
    toolsets:
      - type: background_agents

  claude-coder:
    description: Frontend specialist (Claude Code)
    harness:
      type: claude-code
      effort: medium

  codex-coder:
    description: Backend specialist (Codex)
    harness:
      type: codex
```

The orchestrator calls `run_background_agent` for each coding task, then uses `wait_background_agents` with all task IDs to join and collect results. It can inspect live progress separately with `list_background_agents` and `view_background_agent`.

> [!NOTE]
> **Harness toolsets are ignored**
>
> Harness agents use the external CLI's own tools — any `toolsets:` configured on the harness agent are silently ignored. See [Coding Harnesses](../../features/harnesses/index.md) for details and caveats.

See [`examples/coding_harness_background_agents.yaml`](https://github.com/docker/docker-agent/blob/main/examples/coding_harness_background_agents.yaml) for a complete configuration.
