---
title: "Evaluation"
description: "Measure agent quality with automated evaluations — tool call accuracy, response relevance, output size, and more."
keywords: docker agent, ai agents, features, evaluation
weight: 100
canonical: https://docs.docker.com/ai/docker-agent/features/evaluation/
aliases:
  - /ai/docker-agent/evals/
---

_Measure agent quality with automated evaluations — tool call accuracy, response relevance, output size, and more._

## Overview

The `docker agent eval` command runs your agent against a set of recorded sessions and scores the results. Each eval session captures a user question, the expected tool calls, and criteria the response must satisfy. Docker Agent replays the question, compares the agent's behavior to expectations, and produces a report.

> [!NOTE]
> **Container runtime required**
>
> Evaluations run inside containers for isolation. Each eval gets a clean environment with optional setup scripts. A running Docker-compatible container CLI/runtime is required: Docker Desktop or Docker Engine by default, or another Docker-compatible runtime such as Podman selected with `--container-runtime`.

## Quick Start

```bash
# Run evaluations for an agent
$ docker agent eval agent.yaml

# Specify a custom evals directory
$ docker agent eval agent.yaml ./my-evals

# Run with 8 concurrent evaluations
$ docker agent eval agent.yaml -c 8

# Only run evals matching a pattern
$ docker agent eval agent.yaml --only "auth*"

# Repeat each eval 5 times to compute a baseline
$ docker agent eval agent.yaml --repeat 5

# Repeat a specific eval 5 times
$ docker agent eval agent.yaml --only "auth*" --repeat 5

# Use a Docker-compatible runtime such as Podman
$ docker agent eval agent.yaml --container-runtime podman
```

## Eval Directory Structure

By default, Docker Agent looks for eval sessions in an `evals/` directory next to your agent config:

```bash
my-agent/
├── agent.yaml
└── evals/
    ├── 41b179a2-....json          # Eval session 1
    ├── 5d83e247-....json          # Eval session 2
    └── results/                   # Output (auto-created)
        ├── adjective-noun-1234.json
        ├── adjective-noun-1234.log
        ├── adjective-noun-1234.db
        └── adjective-noun-1234-sessions.json
```

## Eval Session Format

Each eval file is a JSON session that captures a complete conversation. The key fields for evaluation are the user message, the expected tool calls (recorded from a real session), and optional eval criteria:

```json
{
  "id": "41b179a2-ed19-4ae2-a45d-95775aaa90f7",
  "title": "Counting Files in Local Folder",
  "messages": [
    {
      "message": {
        "message": {
          "role": "user",
          "content": "How many files in the local folder?"
        }
      }
    },
    {
      "message": {
        "agent_name": "root",
        "message": {
          "role": "assistant",
          "tool_calls": [
            {
              "id": "call_abc123",
              "type": "function",
              "function": {
                "name": "list_directory",
                "arguments": "{\"path\":\"./\"}"
              }
            }
          ]
        }
      }
    },
    {
      "message": {
        "agent_name": "root",
        "message": {
          "role": "assistant",
          "content": "There are 2 files in the local folder..."
        }
      }
    }
  ],
  "evals": {
    "relevance": [
      "The response mentions exactly 2 files",
      "The response lists README.md and agent.yaml"
    ],
    "size": "S",
    "working_dir": "my-project",
    "setup": "echo 'hello' > test.txt"
  }
}
```

## Eval Criteria

The `evals` object inside each session controls what gets scored:

| Field         | Type     | Description                                                                               |
| ------------- | -------- | ----------------------------------------------------------------------------------------- |
| `relevance`   | string[] | Statements that must be true about the agent's response. Scored by an LLM judge.          |
| `size`        | string   | Expected response size: `S`, `M`, `L`, or `XL`. Compared against actual output length.    |
| `working_dir` | string   | Subdirectory under `evals/working_dirs/` to mount as the container's working directory.   |
| `setup`       | string   | Shell script to run in the container before the agent executes (e.g., create test files). |

## Scoring Metrics

Docker Agent evaluates agents across three dimensions:

| Metric              | How It's Measured                                                                                                         |
| ------------------- | ------------------------------------------------------------------------------------------------------------------------- |
| **Tool Calls (F1)** | F1 score between the expected tool call sequence (from the recorded session) and the actual tool calls made by the agent. |
| **Relevance**       | An LLM judge (configurable via `--judge-model`) evaluates whether each relevance statement is satisfied by the response.  |
| **Size**            | Whether the response length matches the expected size category (S/M/L/XL).                                                |

## Serve-safety verification and rollback

When changing an agent served over MCP HTTP, chat, or A2A, add an evaluation that attempts an approval-requiring tool call and verifies the resolved safety policy and authentication behavior. Run the evaluation with the same explicit `--safety` setting used in deployment. If a rollout must be reversed, stop the affected listener, restore the prior agent configuration and explicit safety flag, then restart only after confirming non-loopback listeners still require authentication. Do not restore an unauthenticated network listener as a rollback shortcut.

## Creating Eval Sessions

The easiest way to create eval sessions is from real conversations:

1. Run your agent interactively: `docker agent run agent.yaml`
2. Have a conversation that tests the behavior you care about
3. Use the `/eval` slash command in the TUI to save the session as an eval file
4. Edit the generated JSON to add `evals` criteria (relevance, size, etc.)

> [!TIP]
> Start with tool call scoring (automatic from recorded sessions), then add relevance criteria for the responses you care most about.

## CLI Flags

```bash
$ docker agent eval <agent-file>|<registry-ref> [<eval-dir>|./evals]
```

| Flag                | Default                     | Description                                                       |
| ------------------- | --------------------------- | ----------------------------------------------------------------- |
| `-c, --concurrency` | num CPUs                    | Number of concurrent evaluation runs                              |
| `--judge-model`     | `anthropic/claude-opus-5` | Model for LLM-as-a-judge relevance scoring                        |
| `--output`          | `<eval-dir>/results`  | Directory for results, logs, and session databases                |
| `--only`            | (all)                       | Only run evals with file names matching these patterns            |
| `--base-image`      | (default)                   | Custom base image for eval containers (see [Custom Base Images](#custom-base-images)) |
| `--container-runtime` | `docker`                  | Container runtime executable for building and running evaluations (e.g. `podman`) |
| `--keep-containers` | `false`                     | Keep containers after evaluation (don't remove with `--rm`)       |
| `-e, --env`         | (none)                      | Environment variables to pass to container (`KEY` or `KEY=VALUE`) |
| `--repeat`          | `1`                         | Number of times to repeat each evaluation (useful for computing baselines) |
| `--baseline`        | (none)                      | Compare against a previously saved run JSON and exit non-zero on regression (see [Regression gate](#regression-gate)) |
| `--regression-tolerance` | `0`                    | How far an aggregate quality rate may fall before `--baseline` reports a regression (0–1) |

### Regression gate

`--baseline` compares the run against a previous one and exits non-zero when
quality regressed, so an eval suite can gate CI:

```console
$ docker agent eval ./agent.yaml --baseline results/2026-08-01-run.json
```

The baseline is the run JSON written by a previous invocation —
`<output>/<run-name>.json` — so there is no separate artifact to produce.

Four rules decide the verdict, and they are worth knowing before wiring this
into CI:

- **The tolerance governs aggregate rates only.** An LLM judge does not return
  the same score twice, so without a tolerance the gate flaps. `--regression-tolerance 0.05`
  lets an aggregate rate fall five points before it counts.
- **An evaluation that passed and now fails always gates**, regardless of the
  tolerance. That transition is the signal the gate exists to catch, so it is
  never absorbed.
- **Cost is reported but never gates.** A provider price change is not a quality
  regression.
- **An added *failing* evaluation gates** via the aggregate rate, even though no
  existing evaluation regressed. A suite that got worse should say so — but it
  means committing a known-failing eval needs a tolerance bump or a fix.

A baseline that carries no evaluations, or a run that produced none (an
`--only` pattern that matched nothing), is rejected rather than reported as
passing: a gate that cannot fail is worse than no gate.

### Provider Credentials

Eval containers are isolated from your host environment. Dedicated model
provider API keys (for example `ANTHROPIC_API_KEY` or `OPENAI_API_KEY`) are
forwarded into eval containers automatically, so most provider setups work
without extra flags.

> [!WARNING]
> **`GITHUB_TOKEN` and `GH_TOKEN` are not forwarded automatically**
>
> GitHub tokens are broad credentials (git, `gh`, CI, packages), not dedicated
> model API keys, so for security reasons Docker Agent does not forward them
> into eval containers — even when they are set in your shell or in
> `~/.config/cagent/.env`. If your agent uses the `github-copilot` provider,
> pass the token explicitly by name:
>
> ```bash
> docker agent eval agent.yaml ./evals -e GITHUB_TOKEN
> ```
>
> When using a custom env file, both flags are required:
>
> ```bash
> docker agent eval agent.yaml ./evals \
>   --env-from-file /path/to/secrets.env \
>   -e GITHUB_TOKEN
> ```

Note that the LLM judge runs on the host, not inside the eval container. If
the token is not forwarded, judge validation can succeed while every
evaluated agent run fails to authenticate.

### Custom Base Images

When `--base-image` is set, the eval harness builds a derived image on top of your base image at evaluation time. Two things happen automatically:

1. **The docker-agent binary is injected** — it is copied from `docker/docker-agent:edge` into the derived image at build time, so you don't need to include it in your base image.
2. **The entrypoint is overridden** — Docker Agent replaces your base image's entrypoint with its own `/run.sh` wrapper.

Your base image therefore only needs to provide the runtime environment: language runtimes, installed dependencies, test fixtures, the appropriate working directory, and so on. Any `ENTRYPOINT` or `CMD` defined in your base image is ignored.

## Output

After a run completes, Docker Agent produces:

- **Console summary** — Pass/fail status per eval with metric breakdowns
- **JSON results** — Full structured results for programmatic analysis
- **SQLite database** — Complete sessions for detailed investigation and debugging
- **Sessions JSON** — Exported session data for analysis
- **Log file** — Debug-level log of the entire evaluation run

> [!TIP]
> **Debugging Failed Evals**
>
> Use `--keep-containers` to preserve containers after evaluation. You can then inspect them with your selected runtime's `exec` command (`docker exec` by default, `podman exec` with `--container-runtime podman`) to understand why an eval failed. The session database (`.db` file) contains the full conversation history for each eval.

```bash
$ docker agent eval demo.yaml ./evals

  ✓ Counting Files in Local Folder
    ✓ tool calls  ✓ relevance 2/2
  ✓ Checking the Content of README.md File
    ✓ tool calls  ✓ relevance 1/1

✅     Tool Calls: 100.0% avg F1 (2 evals)
✅      Relevance: 3/3 passed (100.0%)

Total Cost: $0.012345
Total Time: 12s

Sessions DB: ./evals/results/happy-panda-1234.db
Sessions JSON: ./evals/results/happy-panda-1234-sessions.json
Log: ./evals/results/happy-panda-1234.log
```

## Example

Here's a minimal evaluation setup:

```yaml
# agent.yaml
agents:
  root:
    model: openai/gpt-4o
    description: Test agent
    instruction: You know how to read/write and list files.
    toolsets:
      - type: filesystem
```

```bash
# Create evals from interactive sessions
$ docker agent run agent.yaml
# ... have conversations, then use /eval to save them

# Run the evaluations
$ docker agent eval agent.yaml ./evals
```

> [!NOTE]
> **See also**
>
> Use `/eval` in the [TUI](../tui/index.md) to create eval sessions from conversations. See the [CLI Reference](../cli/index.md) for all `docker agent eval` flags. Example eval configs are in [examples/eval](https://github.com/docker/docker-agent/tree/main/examples/eval) on GitHub.
