---
title: "Background Jobs Tool"
description: "Run and manage long-running shell commands."
keywords: docker agent, ai agents, tools, toolsets, background jobs, shell
linkTitle: "Background Jobs"
weight: 21
canonical: https://docs.docker.com/ai/docker-agent/tools/background-jobs/
---

_Run and manage long-running shell commands._

## Overview

The `background_jobs` toolset starts shell commands that should keep running while the agent continues with other work, such as local servers, file watchers, long builds, or test suites. It returns a job ID immediately, captures combined stdout/stderr up to 10 MB per job, and terminates all running jobs when the agent session ends.

Use the [`shell`](../shell/index.md) toolset for short synchronous commands. Add both toolsets when an agent needs both synchronous commands and long-running processes.

## Configuration

```yaml
toolsets:
  - type: shell
  - type: background_jobs
```

### Options

| Property       | Type    | Description                                                                                                                                          |
| -------------- | ------- | ---------------------------------------------------------------------------------------------------------------------------------------------------- |
| `env`    | object  | Environment variables to set for all background job commands.                                                                                        |
| `recall` | boolean | Let `run_background_job` expose a `recall` parameter so jobs can steer the agent when they finish (see [Background job recall](#background-job-recall)). Default `false`. |

### Custom Environment Variables

```yaml
toolsets:
  - type: background_jobs
    env:
      MY_VAR: "value"
      PATH: "${env.PATH}:/custom/bin"
```

### Background job recall

Set `recall: true` to let the `run_background_job` tool expose a `recall` boolean parameter:

```yaml
toolsets:
  - type: background_jobs
    recall: true
```

When the agent starts a background job with `recall: true`, docker-agent sends a steering message back into the running agent loop after the job finishes. The message contains a short completion sentence and the job output, so the agent can react without polling `view_background_job`.

Use recall for finite background work where completion matters (for example, a long build or test suite). Avoid it for servers and watchers that are expected to run until stopped. See [`examples/shell_recall.yaml`](https://github.com/docker/docker-agent/blob/main/examples/shell_recall.yaml) for a complete configuration.

## Available Tools

The background jobs toolset exposes five tools:

| Tool Name              | Description                                                                                    |
| ---------------------- | ---------------------------------------------------------------------------------------------- |
| `run_background_job`   | Start a command asynchronously and return a job ID immediately. Use for servers/watchers/etc. |
| `list_background_jobs` | List all background jobs with their status, runtime, and metadata.                             |
| `view_background_job`  | View the buffered output and status of a specific background job by ID.                        |
| `stop_background_job`  | Stop a running background job. Child processes are terminated too.                             |
| `wait_background_job`  | Block until a job finishes and return its exit code and output. Safe on already-finished jobs. |

### `run_background_job` parameters

| Parameter | Type    | Required | Description                                                                                                                                 |
| --------- | ------- | -------- | ------------------------------------------------------------------------------------------------------------------------------------------- |
| `cmd`     | string  | ✓        | The shell command to execute in the background.                                                                                             |
| `cwd`     | string  | ✗        | Working directory to run the command in (default: `.`).                                                                                     |
| `recall`  | boolean | ✗        | Only available when the `background_jobs` toolset has `recall: true`. When true, send a steering message with the job output when it finishes. |

`view_background_job` and `stop_background_job` each take a single required `job_id` string returned by `run_background_job` or `list_background_jobs`.

### `wait_background_job` parameters

| Parameter | Type    | Required | Description                                                                                                    |
| --------- | ------- | -------- | -------------------------------------------------------------------------------------------------------------- |
| `job_id`  | string  | ✓        | Job ID returned by `run_background_job` or `list_background_jobs`.                                             |
| `timeout` | integer | ✗        | Maximum seconds to wait (default: `60`). If the job is still running when the limit fires, the tool returns the current output with a notice and the job continues in the background. |

> [!WARNING]
> **Safety**
>
> Background jobs run shell commands with the same access as the agent process. Stop servers and watchers when they are no longer needed, and use [Sandbox Mode](../../configuration/sandbox/index.md) for additional isolation.
