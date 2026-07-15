---
title: "CLI Reference"
description: "Complete reference for all docker-agent command-line commands and flags."
keywords: docker agent, ai agents, features, cli reference
weight: 30
canonical: https://docs.docker.com/ai/docker-agent/features/cli/
aliases:
  - /ai/docker-agent/reference/cli/
---

_Complete reference for all docker-agent command-line commands and flags._

> [!TIP]
> **No config needed**
>
> Running `docker agent run` without a config file uses a built-in default agent. Perfect for quick experimentation.

## Commands

### `docker agent run`

Launch the interactive TUI with an agent configuration (`.yaml`, `.yml`, or `.hcl`).

```bash
$ docker agent run [config] [message...] [flags]
```

| Flag                                    | Description                                                                                                                               |
| --------------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------- |
| `-a, --agent <name>`                    | Run a specific agent from the config                                                                                                      |
| `--yolo`                                | Auto-approve all tool calls                                                                                                               |
| `--model <ref>`                         | Override model(s). Use `provider/model` for all agents, or `agent=provider/model` for specific agents. Comma-separate multiple overrides. |
| `--session <id>`                        | Resume a previous session. Supports relative refs (`-1` = last, `-2` = second to last). An explicit ID that does not exist yet is created with that ID, so a supervisor can own the session ID upfront and reuse it across runs. |
| `-s, --session-db <path>`               | Path to the SQLite session database (default: `<data-dir>/session.db`, so `~/.cagent/session.db` unless `--data-dir` is set)              |
| `--session-read-only`                   | Open the TUI in read-only mode: conversation history is displayed but no new messages can be sent to the LLM. Cannot be used with `--exec`. |
| `--prompt-file <path>`                  | Include file contents as additional system context (repeatable)                                                                           |
| `--attach <path>`                       | Attach an image file to the initial message                                                                                               |
| `--dry-run`                             | Initialize the agent without executing anything (useful for validating a config)                                                          |
| `--remote <addr>`                       | Use a remote runtime at the given address instead of running the agent locally                                                            |
| `--listen <addr>`                       | Expose this run's control plane over HTTP so an external process can drive the running TUI (send follow-ups, stream events, read the title). Accepts `host:port` or `unix://`, `npipe://`, `fd://`. See the [API Server](../api-server/index.md#listen). |
| `--lean`                                | Use a simplified, non-alternate-screen TUI. Unlike the default full-screen TUI, this renders inline in the normal terminal buffer — useful in environments where an alternate screen is unwanted (e.g. inside tmux panes, CI with a tty, or log-friendly pipelines). Displays an ASCII art banner on startup. |
| `--app-name <name>`                     | Override the application name label shown in the TUI (status bar, window title, "/exit" notifications).                                   |
| `--sidebar`                             | Control sidebar visibility. Set to `--sidebar=false` to hide the sidebar and disable the Ctrl+B toggle (default: `true`).                 |
| `--disable-commands <list>`             | Hide and disable specific slash commands in the TUI. Accepts a comma-separated list of command names (leading slash optional, case-insensitive). E.g. `--disable-commands="/cost,/eval,/model"`. |
| `--theme <name>`                        | Preselect a TUI theme by name, or `auto` to follow the terminal's light/dark background (overrides the theme from user config; ignored outside the interactive TUI) |
| `--on-event <type>=<cmd>`               | Run a shell command when an event of the given type fires (`*=<cmd>` matches any event). Repeatable.                                      |
| `--json`                                | Output results as newline-delimited JSON (use with `--exec`)                                                                              |
| `--hide-tool-calls`                     | Hide tool calls in the output                                                                                                             |
| `--hide-tool-results`                   | Hide tool call results in the output                                                                                                      |
| `--sandbox`                             | Run the agent inside a Docker sandbox (see [Sandbox](../../configuration/sandbox/index.md))                                     |
| `--template <image>`                    | Template image for the sandbox (default: `docker/sandbox-templates:docker-agent`)                                                         |
| `--sbx`                                 | Prefer the `sbx` CLI backend when available (default `true`; set `--sbx=false` to force `docker sandbox`)                                 |
| `--no-kit`                              | Disable the [auto-kit](../../configuration/sandbox/index.md#auto-kit): do not stage skills or prompt files into the sandbox    |
| `--agent-picker [refs]`                 | Show a full-screen interactive picker before launching, letting you browse and select an agent. Accepts an optional comma-separated list of agent references to show (defaults to the built-in `default` and `coder` agents plus any agent configs found in `~/.agents`). Arrow keys navigate; `?` toggles the YAML preview panel; `l` (or mouse-click) toggles the **Lean Mode** checkbox to launch in the lean TUI; `b` (or clicking **[ Open Board ]**) opens the Kanban board (`docker agent board`) instead of running an agent; Enter confirms. Not available in `--exec` or non-TTY modes. |
| `-w, --worktree [name]`                 | Run the agent in a fresh git worktree of the working directory, isolating its changes from your checkout. Optionally name it (`--worktree=my-feature`); otherwise a name is generated. Requires the working directory to be inside a git repository. Every tool (the shell included) runs inside the worktree. Combine with `--working-dir` to branch from another repository, and with `--session` to resume into the same worktree later. Cannot be combined with `--remote` or `--sandbox`. When the session ends, a clean worktree is removed automatically; one with work prompts to keep or remove (never in `--exec`). |
| `--worktree-base <ref>`                  | Branch the `--worktree` from `<ref>` (a branch, tag, commit, or remote-tracking ref like `origin/main`) instead of the current `HEAD`. A remote-tracking ref is fetched first so the worktree starts from the latest remote state. Requires `--worktree`; cannot be combined with `--worktree-pr`, `--remote`, or `--sandbox`. |
| `--worktree-pr <number\|url>`            | Run the agent in a git worktree checked out on an existing GitHub pull request (PR number, `#123`, or PR URL). Continues the PR's branch so commits push back to it. Requires the [GitHub CLI](https://cli.github.com/) (`gh`). Cannot be combined with `--worktree`, `--remote`, or `--sandbox`. |
| `--working-dir <path>`                  | Set the working directory for the session (applies to tools and relative paths)                                                           |
| `--env-from-file <path>`                | Load environment variables from file (repeatable)                                                                                         |
| `--code-mode-tools`                     | Provide a single tool to call other tools via JavaScript (forces code-mode tools globally)                                                |
| `--models-gateway <addr>`               | Route model traffic through a gateway. Also reads `DOCKER_AGENT_MODELS_GATEWAY` (legacy `CAGENT_MODELS_GATEWAY`) env var.                  |
| `--hook-pre-tool-use <cmd>`             | Add a pre-tool-use hook command (repeatable). See [Hooks](../../configuration/hooks/index.md).                                  |
| `--hook-post-tool-use <cmd>`            | Add a post-tool-use hook command (repeatable)                                                                                             |
| `--hook-session-start <cmd>`            | Add a session-start hook command (repeatable)                                                                                             |
| `--hook-session-end <cmd>`              | Add a session-end hook command (repeatable)                                                                                               |
| `--hook-on-user-input <cmd>`            | Add an on-user-input hook command (repeatable)                                                                                            |
| `--hook-stop <cmd>`                     | Add a stop hook command, fired when the model finishes responding (repeatable)                                                            |
| `--fake <path>`                         | Replay AI responses from a cassette file (for testing). Mutually exclusive with `--record`.                                               |
| `--fake-stream [ms]`                    | When replaying with `--fake`, simulate streaming with a delay between chunks (defaults to 15ms when given without a value).               |
| `--record [path]`                       | Record AI API interactions to a cassette file and generate a TUI e2e test from the session (auto-generates filename if no path given). Routes through `--models-gateway` when one is configured. |
| `-d, --debug`                           | Enable debug logging                                                                                                                      |
| `--log-file <path>`                     | Custom debug log location                                                                                                                 |
| `-o, --otel`                            | Enable OpenTelemetry observability: traces, metrics, and logs. Requires `OTEL_EXPORTER_OTLP_ENDPOINT` to export to a collector. |

```bash
# Examples
$ docker agent run agent.yaml
$ docker agent run agent.yaml "Fix the bug in auth.go"
$ docker agent run agent.yaml -a developer --yolo
$ docker agent run agent.yaml --model anthropic/claude-sonnet-4-5
$ docker agent run agent.yaml --model "dev=openai/gpt-4o,reviewer=anthropic/claude-sonnet-4-5"
$ docker agent run agent.yaml --session -1  # resume last session
$ docker agent run agent.yaml --session -1 --session-read-only  # review last session without sending messages
$ docker agent run agent.yaml --prompt-file ./context.md  # include file as context

# Add hooks from the command line
$ docker agent run agent.yaml --hook-session-start "./scripts/setup-env.sh"
$ docker agent run agent.yaml --hook-pre-tool-use "./scripts/validate.sh" --hook-post-tool-use "./scripts/log.sh"

# Queue multiple messages (processed in sequence)
$ docker agent run agent.yaml "question 1" "question 2" "question 3"

# Customize TUI display
$ docker agent run agent.yaml --app-name "My Project"
$ docker agent run agent.yaml --sidebar=false
$ docker agent run agent.yaml --disable-commands="/cost,/eval,/model"

# Browse and pick an agent interactively
$ docker agent run --agent-picker
$ docker agent run --agent-picker=agentcatalog/coder,agentcatalog/researcher
```

> [!TIP]
> **Lean, inline TUI**
>
> Pass `--lean` to get a lightweight TUI that renders inline in your terminal (no alternate screen). It displays an ASCII art banner on startup and supports the same slash commands and streaming output as the full TUI, making it handy inside tmux, scripts, or any context where a full-screen takeover is unwanted.

> [!TIP]
> **Isolate a run in a git worktree**
>
> When the working directory is inside a git repository, `--worktree` creates a fresh [git worktree](https://git-scm.com/docs/git-worktree) and points the session at it, so the agent's edits land on a separate branch and never touch your checkout. Every tool — the shell included — runs inside the worktree. The worktree is stored under `<data-dir>/worktrees/<name>` on a branch named `worktree-<name>`.

```bash
# Run in an isolated worktree with a generated name (e.g. "focused_turing")
$ docker agent run agent.yaml --worktree
$ docker agent run agent.yaml -w "Refactor the auth package"

# Give the worktree (and its branch) an explicit name
$ docker agent run agent.yaml --worktree=auth-refactor

# Branch the worktree from another ref instead of the current HEAD.
# A remote-tracking ref is fetched first, so the worktree starts from the
# latest remote state.
$ docker agent run agent.yaml --worktree=auth-refactor --worktree-base origin/main

# Resume a worktree run later: the session remembers its worktree, so you
# don't pass --worktree again (which would fail — the worktree already exists).
$ docker agent run agent.yaml --worktree=auth-refactor   # first run, creates it
$ docker agent run agent.yaml --session -1               # resumes into the same worktree

# Check out an existing GitHub pull request to continue it (requires gh)
$ docker agent run agent.yaml --worktree-pr 123
$ docker agent run agent.yaml --worktree-pr https://github.com/owner/repo/pull/123
```

With `--worktree-pr`, the PR's head branch is checked out tracking its remote (via the [GitHub CLI](https://cli.github.com/)), so commits made during the run push straight back to the pull request. The worktree is stored under `<data-dir>/worktrees/pr-<number>`.

When the interactive session ends, the worktree is cleaned up based on its state:

- **Clean** (no uncommitted changes, untracked files, or new commits): the worktree and its branch are removed automatically.
- **Has work** (uncommitted changes, untracked files, or new commits): you're prompted to keep or remove it. Keeping preserves the directory and branch so you can return later; removing discards the worktree, its branch, and all that work.
- **Non-interactive runs** (`--exec`): the worktree is never cleaned up — it's left in place for inspection.

A worktree is only ever removed if `--worktree` created it for this run; a pre-existing worktree is never touched.

A kept worktree can be resumed: pass `--session` (a relative ref like `-1`, or the session id) and the run reattaches to the same worktree directory and branch automatically. Don't re-pass `--worktree` on resume — that would try to create a new worktree and fail because it already exists.

### `docker agent run --exec`

Run an agent in non-interactive (headless) mode. No TUI — output goes to stdout.

```bash
$ docker agent run --exec [config] [message...] [flags]
```

```bash
# One-shot task
$ docker agent run --exec agent.yaml "Create a Dockerfile for a Python Flask app"

# With auto-approve
$ docker agent run --exec agent.yaml --yolo "Set up CI/CD pipeline"

# Multi-turn conversation
$ docker agent run --exec agent.yaml "question 1" "question 2" "question 3"
```

### `docker agent new`

Interactively generate a new agent configuration file.

```bash
$ docker agent new [flags]

# Examples
$ docker agent new
$ docker agent new --model openai/gpt-5
$ docker agent new --model dmr/ai/gemma3-qat:12B --max-iterations 15
```

### `docker agent models`

List models available for use with `--model`. By default only shows models for providers you have credentials for. Aliases: `docker agent models list`, `docker agent models ls`.

```bash
$ docker agent models [flags]
```

| Flag                   | Default | Description                                                                        |
| ---------------------- | ------- | ---------------------------------------------------------------------------------- |
| `-p, --provider <id>`  | (none)  | Filter models by provider name (e.g. `openai`, `anthropic`, `dmr`, `ollama`, …).   |
| `--format <fmt>`       | `table` | Output format: `table` or `json`.                                                  |
| `-a, --all`            | `false` | Include models from all providers, not just those you have credentials for.        |

```bash
# Examples
$ docker agent models                                 # only providers you can use
$ docker agent models --all                           # every provider the catalog knows about
$ docker agent models --provider openai
$ docker agent models --format json | jq
```

### `docker agent setup`

Set up a model interactively. Two paths: pick a cloud provider, paste its API key, and choose where to store it (macOS Keychain, `pass`, or the docker agent env file `~/.config/cagent/.env`), or check Docker Model Runner and pull a local model (no API key needed). Ends with the exact command to start chatting. Secret values are never printed.

The wizard is also offered automatically when an interactive run finds no usable model (decline-able; set `DOCKER_AGENT_NO_SETUP=1` to suppress the offer).

```bash
$ docker agent setup
```

### `docker agent doctor`

Diagnose the model and credential setup. Reports which model providers have credentials and where each credential comes from (shell environment, env file, pass, keychain, …), whether Docker Model Runner is reachable and which models are pulled, and which model the `auto` selection would pick. Secret values are never printed. Exits with a non-zero status when an issue would prevent an agent from running, which makes it usable in scripts and CI.

```bash
$ docker agent doctor [agent-file|registry-ref] [flags]
```

With an agent file, also lists the environment variables that file requires (model credentials and tool secrets such as `GITHUB_PERSONAL_ACCESS_TOKEN`), whether each one is set, and from which source.

| Flag                     | Default | Description                                                    |
| ------------------------ | ------- | -------------------------------------------------------------- |
| `--json`                 | `false` | Output the full report in JSON format (for scripting).         |
| `--env-from-file <file>` | (none)  | Also check variables supplied by an env file.                  |
| `--models-gateway <url>` | (none)  | Diagnose against a models gateway (credentials come from it).  |

```bash
# Examples
$ docker agent doctor                        # credential, DMR, and auto-selection state
$ docker agent doctor ./agent.yaml           # also check that file's requirements
$ docker agent doctor --json | jq .issues
```

### `docker agent serve api`

Start the HTTP API server for programmatic access. The argument can be a single agent file, a registry reference, or a directory — when given a directory, every `.yaml`/`.yml`/`.hcl` file in it is exposed as a separate entry under `/api/agents`.

```bash
$ docker agent serve api <agent-file>|<agents-dir>|<registry-ref> [flags]
```

| Flag                       | Default            | Description                                                                                                |
| -------------------------- | ------------------ | ---------------------------------------------------------------------------------------------------------- |
| `-l, --listen <addr>`      | `127.0.0.1:8080`   | Address to listen on.                                                                                      |
| `--auth-token <token>`     | (none)             | Bearer token required for all API requests. When set, every request must include `Authorization: Bearer <token>`. Leave empty to disable authentication (safe when listening on loopback interfaces only). |
| `-s, --session-db <path>`  | `session.db`       | Path to the SQLite session database (relative paths resolve against the working directory).                |
| `--pull-interval <minutes>`| `0`                | Periodically re-pull OCI/URL references and refresh the agent definition. `0` disables auto-pull.          |
| `--fake <path>`             | (none)             | Replay AI responses from a cassette file (for testing). Mutually exclusive with `--record`.               |
| `--record [path]`           | (none)             | Record AI API interactions to a cassette file. Routes through `--models-gateway` when one is configured. |
| `--mcp-oauth-redirect-uri <url>` | (none)        | OAuth redirect URI for the unmanaged MCP OAuth flow in server mode. When set, the runtime drives PKCE and code exchange in-process and sends the full authorize URL to the client via elicitation. See [Remote MCP](../remote-mcp/index.md) for details. |

> **Diagnostics:** Set `CAGENT_PPROF_ADDR=127.0.0.1:6060` (or `--pprof-addr`, a hidden flag) to start a live Go pprof server at `/debug/pprof/`. Use a loopback address; a non-loopback binding logs a security warning.

All [runtime configuration flags](#runtime-configuration-flags) (`--working-dir`, `--env-from-file`, `--models-gateway`, `--hook-*`, …) are also accepted.

```bash
# Examples
$ docker agent serve api agent.yaml
$ docker agent serve api agent.yaml --listen :8080
$ docker agent serve api ./agents/                          # directory of agent YAML/HCL configs
$ docker agent serve api ociReference --pull-interval 10    # auto-refresh
```

See [API Server](../api-server/index.md) for the full HTTP API reference.

### `docker agent serve mcp`

Expose agents as MCP tools for use in Claude Desktop, Claude Code, or other MCP clients. Defaults to stdio transport; use `--http` to start a streaming HTTP server instead.

```bash
$ docker agent serve mcp <config> [flags]
```

| Flag                   | Default            | Description                                                                                       |
| ---------------------- | ------------------ | ------------------------------------------------------------------------------------------------- |
| `-a, --agent <name>`   | (all agents)       | Name of the agent to expose. If omitted, every agent in the config is exposed as a separate tool. |
| `--tool-name <name>`   | (agent name)       | Override the MCP tool identifier clients call; only valid when exposing a single agent.           |
| `--http`               | `false`            | Use streaming HTTP transport instead of stdio.                                                    |
| `-l, --listen <addr>`  | `127.0.0.1:8081`   | Address to listen on (only used with `--http`).                                                   |
| `--mcp-keepalive <dur>`| `0` (disabled)     | Interval between MCP keep-alive pings (e.g. `30s`).                                               |
| `--attach [target]`    | (none)             | Attach to a running TUI run by pid, address, or session id; given without a value, selects the most recent run.   |

All [runtime configuration flags](#runtime-configuration-flags) are also accepted.

```bash
# Examples
$ docker agent serve mcp agent.yaml                                # stdio transport
$ docker agent serve mcp agent.yaml --http --listen 127.0.0.1:9090 # streaming HTTP
$ docker agent serve mcp agent.yaml --working-dir /path/to/project
$ docker agent serve mcp agentcatalog/coder
```

See [MCP Mode](../mcp-mode/index.md) for detailed setup.

### `docker agent serve a2a`

Start an A2A (Agent-to-Agent) protocol server.

```bash
$ docker agent serve a2a <config> [flags]
```

| Flag                   | Default            | Description                                                                                |
| ---------------------- | ------------------ | ------------------------------------------------------------------------------------------ |
| `-a, --agent <name>`   | (team default)     | Name of the agent to run. Defaults to the team's first agent if not specified.             |
| `-l, --listen <addr>`  | `127.0.0.1:8082`   | Address to listen on.                                                                       |
| `-s, --session-db <path>` | `<data-dir>/session.db` | Path to the SQLite session database.                                                 |

All [runtime configuration flags](#runtime-configuration-flags) are also accepted.

```bash
# Examples
$ docker agent serve a2a agent.yaml
$ docker agent serve a2a agent.yaml --listen 127.0.0.1:9000
$ docker agent serve a2a agentcatalog/pirate
```

### `docker agent serve acp`

Start an ACP (Agent Client Protocol) server over stdio. This allows external clients to interact with your agents using the ACP protocol.

```bash
$ docker agent serve acp <config> [flags]
```

| Flag                      | Default                     | Description                                       |
| ------------------------- | --------------------------- | ------------------------------------------------- |
| `-s, --session-db <path>` | `<data-dir>/session.db`     | Path to the SQLite session database.              |

All [runtime configuration flags](#runtime-configuration-flags) are also accepted.

```bash
# Examples
$ docker agent serve acp agent.yaml
$ docker agent serve acp ./team.yaml
$ docker agent serve acp agentcatalog/pirate
```

See [ACP](../acp/index.md) for details on the Agent Client Protocol.

### `docker agent serve chat`

Start an HTTP server that exposes one or more agents through an **OpenAI-compatible Chat Completions API** at `/v1/chat/completions` and `/v1/models`. This lets any tool that already speaks the OpenAI protocol — for example [Open WebUI](https://github.com/open-webui/open-webui), `curl`, the OpenAI Python SDK, or LangChain — drive a docker-agent agent without any custom integration.

```bash
$ docker agent serve chat <config> [flags]
```

| Flag                          | Default            | Description                                                                                                       |
| ----------------------------- | ------------------ | ----------------------------------------------------------------------------------------------------------------- |
| `-a, --agent <name>`          | (all agents)       | Name of the agent to expose. If omitted, every agent in the config is exposed as a separate model.                |
| `-l, --listen <addr>`         | `127.0.0.1:8083`   | Address to listen on.                                                                                             |
| `--cors-origin <origin>`      | (none)             | Allowed CORS origin (e.g. `https://example.com`). Empty disables CORS.                                            |
| `--api-key <token>`           | (none)             | Required Bearer token clients must present (`Authorization: Bearer <token>`). Empty disables auth.                |
| `--api-key-env <name>`        | (none)             | Read the API key from this environment variable instead of the command line.                                      |
| `--max-request-size <bytes>`  | `1048576` (1 MiB)  | Maximum request body size.                                                                                        |
| `--request-timeout <dur>`     | `5m`               | Per-request timeout (covers model + tool calls + streaming).                                                      |
| `--conversations-max <n>`     | `0`                | Cache up to N conversations server-side, keyed by `X-Conversation-Id`. `0` disables — clients must resend history. |
| `--conversation-ttl <dur>`    | `30m`              | Idle TTL after which a cached conversation is evicted.                                                            |
| `--max-idle-runtimes <n>`     | `4`                | Maximum number of idle runtimes pooled per agent. `0` disables pooling.                                           |

```bash
# Examples
$ docker agent serve chat agent.yaml
$ docker agent serve chat ./team.yaml --agent reviewer
$ docker agent serve chat agentcatalog/pirate --listen 127.0.0.1:9090
$ docker agent serve chat agent.yaml --api-key-env CHAT_BEARER_TOKEN

# Drive it from any OpenAI-compatible client
$ curl http://127.0.0.1:8083/v1/chat/completions \
    -H 'Content-Type: application/json' \
    -d '{"model": "root", "messages": [{"role": "user", "content": "hello"}]}'
```

See [Chat Server](../chat-server/index.md) for the full feature reference.

### `docker agent board`

Launch a full-screen Kanban TUI for orchestrating multiple agents at once. Each card runs an agent in its own tmux session on an isolated git worktree; moving a card forward through the pipeline (Dev → Review → Push → Done) delivers the destination column's prompt to that card's agent. Projects and column prompts are stored in the global config file (`~/.config/cagent/config.yaml`) and can be managed from the TUI.

```bash
$ docker agent board
```

Takes no arguments or flags. Requires `tmux` and `git` to be installed.

See [Kanban Board](../board/index.md) for key bindings, configuration, and workflow details.

### `docker agent share push` / `docker agent share pull`

Share agents via OCI registries.

```bash
# Push an agent
$ docker agent share push ./agent.yaml docker.io/username/my-agent:latest

# Pull an agent
$ docker agent share pull docker.io/username/my-agent:latest

# Force pull, overwriting the local copy
$ docker agent share pull docker.io/username/my-agent:latest --force
```

| Flag       | Applies to | Description                                                |
| ---------- | ---------- | ---------------------------------------------------------- |
| `--force`  | `pull`     | Force pull even if the configuration already exists locally |

See [Agent Distribution](../../concepts/distribution/index.md) for full registry workflow details.

### `docker agent eval`

Run agent evaluations against a directory of recorded sessions.

```bash
$ docker agent eval <agent-file>|<registry-ref> [<eval-dir>|./evals] [flags]
```

| Flag                | Default                              | Description                                                                |
| ------------------- | ------------------------------------ | -------------------------------------------------------------------------- |
| `-c, --concurrency` | num CPUs                             | Number of concurrent evaluation runs                                       |
| `--judge-model`     | `anthropic/claude-opus-4-5-20251101` | Model for LLM-as-a-judge relevance scoring (format: `provider/model`)      |
| `--output <dir>`    | `<eval-dir>/results`                 | Directory for results, logs, and session databases                         |
| `--only <pattern>`  | (all)                                | Only run evals with file names matching these patterns (repeatable)        |
| `--base-image`      | (default)                            | Custom base Docker image for eval containers                               |
| `--keep-containers` | `false`                              | Keep containers after evaluation (don't remove with `--rm`)                |
| `-e, --env`         | (none)                               | Environment variables to pass to container (`KEY` or `KEY=VALUE`, repeatable) |
| `--repeat <n>`      | `1`                                  | Number of times to repeat each evaluation (useful for computing baselines) |

All [runtime configuration flags](#runtime-configuration-flags) are also accepted.

```bash
# Examples
$ docker agent eval agent.yaml                            # use ./evals
$ docker agent eval agent.yaml ./my-evals                 # custom directory
$ docker agent eval agent.yaml -c 8                       # 8 concurrent evaluations
$ docker agent eval agent.yaml --keep-containers          # keep containers for debugging
$ docker agent eval agent.yaml --only "auth*"             # only run matching evals
$ docker agent eval agent.yaml --repeat 5                 # repeat each eval 5 times
```

See [Evaluation](../evaluation/index.md) for details on creating eval sessions and interpreting results.

### `docker agent version`

Print the version and commit hash for your `docker-agent` install.

```bash
$ docker agent version
docker agent version v1.54.0
Commit: 1737035c
```

### `docker agent alias`

Manage agent aliases for quick access.

```bash
# List aliases
$ docker agent alias ls

# List aliases as JSON
$ docker agent alias list --json

# Add an alias
$ docker agent alias add pirate /path/to/pirate.yaml
$ docker agent alias add other ociReference

# Add an alias with runtime options
$ docker agent alias add yolo-coder agentcatalog/coder --yolo
$ docker agent alias add fast-coder agentcatalog/coder --model openai/gpt-4o-mini
$ docker agent alias add safe-coder agentcatalog/coder --sandbox
$ docker agent alias add turbo agentcatalog/coder --yolo --model anthropic/claude-sonnet-4-5

# Use an alias
$ docker agent run pirate
$ docker agent run yolo-coder
```

**Alias Options:** Aliases can include runtime options that apply automatically when used:

- `--yolo` — Auto-approve all tool calls when running the alias
- `--model <ref>` — Override the model for the alias
- `--hide-tool-results` — Hide tool call results in the TUI when running the alias
- `--sandbox` — Always run the alias inside a [Docker sandbox](../../configuration/sandbox/index.md)

When listing aliases, options are shown in brackets:

```bash
$ docker agent alias ls
Registered aliases (3):

  fast-coder  → agentcatalog/coder [model=openai/gpt-4o-mini]
  turbo       → agentcatalog/coder [yolo, model=anthropic/claude-sonnet-4-5]
  yolo-coder  → agentcatalog/coder [yolo]

Run an alias with: docker agent run <alias>
```

Pass `--json` to output aliases as a JSON array instead of the formatted table. Each entry includes the alias `name` and its options:

```bash
$ docker agent alias list --json
[
  {
    "name": "fast-coder",
    "path": "agentcatalog/coder",
    "model": "openai/gpt-4o-mini"
  },
  {
    "name": "turbo",
    "path": "agentcatalog/coder",
    "yolo": true,
    "model": "anthropic/claude-sonnet-4-5"
  },
  {
    "name": "yolo-coder",
    "path": "agentcatalog/coder",
    "yolo": true
  }
]
```

JSON output is sorted by name and omits false/zero-valued fields. This is useful for scripting and automation.

> [!TIP]
> **Override alias options**
>
> Command-line flags override alias options. For example, `docker agent run yolo-coder --yolo=false` disables yolo mode even though the alias has it enabled.

> [!TIP]
> **Set a default agent**
>
> Create a `default` alias to customize what `docker agent` starts with no arguments:
>
> ```console
> $ docker agent alias add default /my/default/agent.yaml
> ```
>
> Then simply run `docker agent` — it will launch that agent automatically.

### `docker agent sandbox`

Manage settings shared by every [`--sandbox`](../../configuration/sandbox/index.md) run — today, the persistent network allowlist that turns a `Blocked by network policy` 403 into a one-line, durable fix:

```bash
# Allow a host on every subsequent --sandbox run.
$ docker agent sandbox allow api.example.com

# Or several at once.
$ docker agent sandbox allow api.example.com registry.npmjs.org:443

# See what's persisted in ~/.config/cagent/config.yaml.
$ docker agent sandbox list

# Drop a host you no longer need.
$ docker agent sandbox deny api.example.com
```

Entries are unioned with the gateway, the kit-resolved tool install hosts, and any `runtime.network_allowlist` declared by the agent. The launch summary lists every source separately so you can see which holes were punched by which layer.

## Global Flags

These flags are available on every `docker agent` command:

| Flag                      | Description                                                                            |
| ------------------------- | -------------------------------------------------------------------------------------- |
| `-d, --debug`             | Enable debug logging (default location: `~/.cagent/cagent.debug.log`)                  |
| `--log-file <path>`       | Custom debug log location (only used with `--debug`)                                   |
| `-o, --otel`              | Enable OpenTelemetry observability: traces, metrics, and logs. Requires `OTEL_EXPORTER_OTLP_ENDPOINT` to export to a collector. |
| `--cache-dir <path>`      | Override the cache directory (default: `~/Library/Caches/cagent` on macOS)             |
| `--config-dir <path>`     | Override the config directory (default: `~/.config/cagent`). Also reads `DOCKER_AGENT_CONFIG_DIR` (legacy `CAGENT_CONFIG_DIR`) env var. |
| `--data-dir <path>`       | Override the data directory (default: `~/.cagent`; holds `session.db`, worktrees, plans, …)            |
| `--help`                  | Show help for any command                                                              |

### OpenTelemetry environment variables

When `--otel` is enabled, the standard [OTel SDK env vars](https://opentelemetry.io/docs/specs/otel/configuration/sdk-environment-variables/) are honored (`OTEL_EXPORTER_OTLP_ENDPOINT`, `OTEL_RESOURCE_ATTRIBUTES`, etc.). Two additional docker-agent-specific variables control GenAI instrumentation:

| Variable | Default | Description |
| -------- | ------- | ----------- |
| `OTEL_INSTRUMENTATION_GENAI_CAPTURE_MESSAGE_CONTENT` | `false` | Set to `true` to capture prompt text, model responses, tool arguments, and tool results as span attributes. Off by default because these fields may contain PII. |
| `OTEL_SEMCONV_STABILITY_OPT_IN` | (dual-emit) | Set to `gen_ai_latest_experimental` to emit only the spec-defined `gen_ai.*` keys from the [GenAI semantic conventions](https://opentelemetry.io/docs/specs/semconv/gen-ai/). The default dual-emit mode emits both `gen_ai.*` and legacy keys so existing dashboards continue working. |

## Runtime Configuration Flags

These flags are accepted by every command that loads an agent (`run`, `run --exec`, `new`, `eval`, `serve api`, `serve mcp`, `serve a2a`, `serve acp`, `serve chat`). They are listed once here to avoid repetition in the per-command tables above.

| Flag                            | Description                                                                                                              |
| ------------------------------- | ------------------------------------------------------------------------------------------------------------------------ |
| `--working-dir <path>`          | Set the working directory for the session (applies to tools and relative paths).                                         |
| `--env-from-file <path>`        | Load environment variables from file (repeatable).                                                                       |
| `--code-mode-tools`             | Provide a single tool to call other tools via JavaScript (forces code-mode tools globally).                              |
| `--models-gateway <addr>`       | Route model traffic through a gateway. Reads `DOCKER_AGENT_MODELS_GATEWAY` (legacy `CAGENT_MODELS_GATEWAY`) env var.      |
| `--hook-pre-tool-use <cmd>`     | Add a pre-tool-use hook command (repeatable). See [Hooks](../../configuration/hooks/index.md).                 |
| `--hook-post-tool-use <cmd>`    | Add a post-tool-use hook command (repeatable).                                                                           |
| `--hook-session-start <cmd>`    | Add a session-start hook command (repeatable).                                                                           |
| `--hook-session-end <cmd>`      | Add a session-end hook command (repeatable).                                                                             |
| `--hook-on-user-input <cmd>`    | Add an on-user-input hook command (repeatable).                                                                          |
| `--hook-stop <cmd>`             | Add a stop hook command, fired when the model finishes responding (repeatable).                                          |

## Agent References

Commands that accept a config support multiple reference types:

| Type          | Example                                     |
| ------------- | ------------------------------------------- |
| Local file    | `./agent.yaml`                              |
| OCI registry  | `docker.io/username/agent:latest`           |
| Agent catalog | `agentcatalog/pirate`                       |
| Alias         | `pirate` (after `docker agent alias add`)   |
| Default       | (no argument) — uses built-in default agent |

> [!NOTE]
> **Debugging**
>
> Having issues? See [Troubleshooting](../../community/troubleshooting/index.md) for debug mode, log analysis, and common solutions.
