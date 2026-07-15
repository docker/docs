---
title: "Todo Tool"
description: "Task list management for complex multi-step workflows."
keywords: docker agent, ai agents, tools, toolsets, todo tool
linkTitle: "Todo"
weight: 170
canonical: https://docs.docker.com/ai/docker-agent/tools/todo/
---

_Task list management for complex multi-step workflows._

## Overview

The todo tool provides task list management. Agents can create, update, list, and track progress on tasks with status tracking (pending, in-progress, completed). Useful for complex multi-step workflows where the agent needs to stay organized and ensure all steps are completed.

## Available Tools

| Tool           | Description                              |
| -------------- | ---------------------------------------- |
| `create_todo`  | Create a new task                        |
| `create_todos` | Create multiple tasks at once            |
| `update_todos` | Update status of one or more tasks       |
| `list_todos`   | List all current tasks with their status |

### Task Statuses

| Status        | Description                  |
| ------------- | ---------------------------- |
| `pending`     | Task has not been started    |
| `in-progress` | Task is currently being done |
| `completed`   | Task is finished             |

## Configuration

```yaml
toolsets:
  - type: todo
```

### Options

| Property | Type    | Default | Description                                                             |
| -------- | ------- | ------- | ----------------------------------------------------------------------- |
| `shared` | boolean | `false` | When `true`, todos are shared across all agents in a multi-agent config |

### Shared Todos

In multi-agent setups, enable shared todos so all agents can see and update the same task list:

```yaml
toolsets:
  - type: todo
    shared: true
```
