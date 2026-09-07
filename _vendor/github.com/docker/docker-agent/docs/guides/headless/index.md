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

## Choosing a Large-Input Strategy

Docker Agent has several unrelated mechanisms for getting a large document, file, or dataset in front of an agent. Which one fits depends on whether the content is local to the machine running Docker Agent, how often it needs to be revisited, and whether you're driving the agent through the CLI/TUI or over one of the HTTP servers:

| You have… | Use… | Notes |
| --- | --- | --- |
| A file on disk, for one interactive turn | `@path` or `/attach` in the TUI (see [File Attachments](../../features/tui/index.md#file-attachments)), `--attach` at startup (seeds an interactive TUI run just as readily as `--exec`), or `/attach` under `--exec` | Read from the local filesystem and inlined into the message once it clears a local admission check. That never crosses Docker Agent's own inbound HTTP boundary for a local (non-`--remote`) run or for the *first* message of a `--remote` run — but a locally resolved attachment added to a *later* message on a remote run does cross it; see the remote-runtime row below. The admission check itself, and how a rejection is reported, differ by entry point; see the read-time limits below. |
| Instructions/context every agent turn should see | `add_prompt_files` in the agent config, or `--prompt-file` on the CLI (see [Prompt Files](../../configuration/agents/index.md#prompt-files)) | Re-read from disk **in full** on every turn and injected as instruction context; local, not an HTTP upload, and not subject to the attachment inlining checks above — memory and the model's context window are the practical ceilings. |
| A one-off piece of text for a **local** headless run | stdin, without `--remote` (see [`--exec` Mode Basics](#--exec-mode-basics) above) | The piped text becomes the message directly; it never crosses Docker Agent's own inbound HTTP boundary and has no size cap of its own, though Docker Agent may still send it onward to a model/provider over HTTP. |
| A one-off piece of text driving a **remote** runtime | stdin with `docker agent run --remote ... -` | The CLI serializes that stdin text into a native API run request sent to whichever Docker Agent server the `--remote` address points at — a `serve api` process or another run's [`--listen`](../../features/api-server/index.md#listen) control plane, never `serve chat` — so it's measured against that server's own limit: `serve api`'s configurable `--max-request-size`, or a `--listen` control plane's fixed, non-configurable 1 MiB cap. See the HTTP body limit below. That *initial* request carries only message text: conversion currently drops any attachment resolved locally (`@path`/`/attach`, `--attach`) for the first message before it's sent. A *later* message to the same remote run — sent while the agent is still busy, via the default steer behavior or an explicit follow-up (Alt+Enter) — does forward a locally resolved attachment as part of that native API request, so it counts toward the same limit and can 413 like any other oversized request. Client-side `--prompt-file` values are **not** forwarded to a remote runtime either way; prompt files configured on the agent itself (`add_prompt_files`) still resolve, but on the server's filesystem. |
| A document collection you'll query repeatedly, or one too large to inline at all | The [`rag` toolset](../../tools/rag/index.md) | Indexed once in the background; each query retrieves only the relevant chunks into context, instead of the whole collection. |
| An OpenAI-compatible client driving Docker Agent over HTTP | [Chat Server](../../features/chat-server/index.md) | Accepts OpenAI-style `text` and [`image_url`](../../features/chat-server/index.md#image-inputs) content parts. The whole request travels as a single HTTP body capped by that server's `--max-request-size`; a data URL's bytes count toward it, while a remote `http(s)://` image URL is passed to the provider rather than fetched by the chat server, and works only if that provider/model supports it. |
| A native integration driving Docker Agent's own session/control protocol | [API Server](../../features/api-server/index.md) | Supports the documented session, run, and event-streaming flows over a single capped HTTP body; see the note below on the upload contract this does *not* provide. |
| A supervisor process driving an already-running interactive session | An attached run's [`--listen`](../../features/api-server/index.md#listen) control plane | Exposes the same session/follow-up/event-streaming API as the API server, but with a fixed, non-configurable 1 MiB request-body cap and no `--auth-token` — there's no `--max-request-size` or `--auth-token` flag for this surface. Keep it on loopback, a unix socket, or behind an authenticating reverse proxy if it must be reachable from elsewhere. |

Three independent ceilings apply, and hitting one says nothing about the others:

1. **The HTTP body limit** — applies to any request reaching one of Docker Agent's three inbound HTTP surfaces: `docker agent serve api`, `docker agent serve chat`, or an interactive run's [`--listen`](../../features/api-server/index.md#listen) control plane. `docker agent run --remote ... -` builds a native run request from stdin and sends it to whichever of `serve api` or a `--listen` control plane the remote address points at — never `serve chat`, which speaks a different protocol. `serve api` and `serve chat` each enforce their own `--max-request-size` (1 MiB default, configurable) and return `413 Request Entity Too Large` above it; a `--listen` control plane enforces the same 1 MiB default but exposes no `--max-request-size` flag (or `--auth-token`) to change it — it's fixed, and reducing the request is the only remedy. See [Troubleshooting: HTTP 413](../../community/troubleshooting/index.md#http-413-request-body-too-large) for how to diagnose and resolve it. Local (non-`--remote`) stdin and prompt files never cross this boundary; a locally resolved `@path`/`/attach`/`--attach` attachment doesn't either for a remote run's initial message, but one added to a later steer or follow-up message on that same remote run does — see the remote-runtime row above.
2. **Local read-time limits, which differ by path** — the TUI and non-TUI (CLI/`--exec`) message assembly don't share one text/binary rule. The TUI (both `/attach` and a typed `@path` reference) runs every file reference through one flat size ceiling *before* it looks at the file's type; non-TUI assembly — the CLI's `--attach` flag and `--exec`'s own `/attach` parsing — inlines text up to its own ceiling and hands supported binary files (images, PDFs) to a separate path with a higher cap, so a binary file that's fine for `--attach`/`--exec` can still be too big for the TUI's flatter, lower ceiling. Either way an oversized file is rejected rather than silently truncated or allowed to exhaust memory, but whether you're actually told about it depends on the surface: only the interactive TUI's `/attach` command shows a visible error notification (with a file-picker fallback), while a bare `@path` reference is only checked speculatively as you type — a rejection there is logged at debug level only, and the reference is left as plain text in your message with no on-screen notice. Outside the interactive TUI, neither the CLI's `--attach` flag nor `--exec`'s `/attach` (which reuses the same non-TUI assembly and rejection path as `--attach`, not the TUI's notification flow) prints anything about a rejected attachment in normal run output — the message just falls back to text-only, and the reason is visible only in the debug log when `--debug` is enabled. Prompt files (`add_prompt_files`/`--prompt-file`) have no equivalent guard — they're read in full on every turn regardless of size, so memory and the model's context window are what actually bound them.
3. **Model/provider context and media limits** — even content that clears the first two layers still has to fit the model's token/context window, and binary attachments (images, PDF, audio, video) only work if the model declares support for that media type (see [Attachment Capability Overrides](../../configuration/models/index.md#attachment-capability-overrides)). Exceeding the context window is a separate failure from either limit above — see [Context Window Exceeded](../../community/troubleshooting/index.md#context-window-exceeded).

> [!NOTE]
> **No file-upload endpoint**
>
> None of Docker Agent's three inbound HTTP surfaces — the API server, the chat server, or an attached run's `--listen` control plane — exposes a multipart file-upload endpoint, and none of them fetches a remote URL on your behalf — the chat server's `image_url` support only ever passes a remote URL through to the provider (see the table above). Content reaches an agent either locally (`@path`/`/attach`, prompt files, `rag`) or embedded directly in the HTTP message body those servers already document; there is no separate attachment channel over any of them.

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

- **What is allowed to run without asking?** — the safety mode (`--safety strict|balanced|restricted|autonomous`, with `--yolo` as the legacy spelling of `autonomous`) and permission allow-lists answer this.
- **What happens if the model runs something it shouldn't have?** — only `--sandbox` answers that one. The rest of this section explains why, and treats that distinction as the whole point.

### `--sandbox`: the isolation boundary

For an untrusted or autonomous agent — anything acting without a human watching approvals — **`--sandbox` is the isolation boundary to reach for**, not a cleverer allow-list. It runs the entire agent, shell calls included, inside a VM managed by [`sbx`](https://docs.docker.com/ai/sandboxes/): a misbehaving or successfully-prompt-injected agent can't touch anything outside the mounted working directory or reach other host/CI state, regardless of which command it runs. That VM isn't disposable or ephemeral — a sandbox matching the current workspace and mount set is retained and reused across subsequent runs rather than torn down when the session ends (see [How It Works](../../configuration/sandbox/index.md#how-it-works)). See [Sandbox Mode](../../configuration/sandbox/index.md) for the full flag reference, `sbx` requirement, network allowlist, and kit staging behavior.

```bash
$ docker agent run --sandbox --exec agent.yaml --json "Fix the failing test"
```

Because the blast radius is contained by the VM boundary, `--sandbox` also makes unattended operation reasonable in CI — and it defaults to exactly that: unless you already passed a `--yolo` or `--safety` flag of your own, `--sandbox` injects `--yolo` for the agent process it runs inside the VM, so the command above already runs unattended with no confirmation prompts. Passing `--yolo` explicitly (`--sandbox --yolo --exec ...`) is equivalent and can make the intent clearer in a script, but it's optional. To keep confirmation prompts even inside the sandbox, select a stricter mode (`--sandbox --safety strict`) or opt out of the legacy default with `--yolo=false` — `--sandbox` only fills in `--yolo` when neither safety flag was set.

If your CI provider already runs each job in its own disposable VM or container — many hosted runners do — and nothing on the runner matters once the job ends, that may already give you an isolation boundary on its own. `--sandbox` still gives you the same guarantee independent of the CI provider, and starts to matter as soon as the agent runs on a persistent self-hosted runner, a long-lived container, or your own workstation.

### Defense in depth, not a boundary: permissions and shell command matching

Permission allow-lists (`permissions.allow` on the agent, or `settings.permissions.allow` globally — see [Permissions](../../configuration/permissions/index.md)) and the balanced safety mode's shell classifier (see [Safety modes](../../configuration/permissions/index.md#safety-modes)) narrow what runs without asking. Used well, they cut down how often you're prompted and catch obviously destructive calls before they run. They are **not** a security boundary:

- Both work by matching the shell command **string** (or, for `permissions`, the tool's arguments). The classifier's safe-list refuses to vouch for any command carrying shell metacharacters (`;`, `&`, `|`, `<`, `>`, backticks, `$(`, newlines — spaced or not), so `ls && rm -rf ~`, `grep foo|rm -rf /`, and `grep x > /etc/passwd` all fall through to a confirmation instead of inheriting a safe verdict. But string matching still can't reason about what a command actually *does* — see the next point.
- Command-string and argument matching in general can't reason about what a command actually does; a dynamically built string, an unusual quoting form, or a wrapper script can slip past any fixed set of patterns.

For unattended runs, the `restricted` safety mode packages this stance as a fail-closed default: classifier-safe calls run, every other unmatched call is **denied outright** instead of falling through to a confirmation prompt nobody will answer. Pair it with an allow-list scoped to what the job actually needs — explicit `allow` rules still win over the mode, so the job's known-good commands run even when the classifier can't vouch for them:

```yaml
# agent.yaml — the allow-list overrides restricted's deny for these calls
permissions:
  allow:
    - "shell:cmd=go test*"
    - "shell:cmd=go build*"
```

```bash
# Safe and allow-listed calls run; every other call is denied without prompting
$ docker agent run --exec --safety restricted agent.yaml --json "Fix the failing test"
```

Like the allow-list itself, `restricted` is defense in depth, not a security boundary: it narrows what runs unattended, but only `--sandbox` contains what a misbehaving agent can do with the calls that are allowed.

Treat permissions and the balanced/restricted modes as a way to reduce prompt fatigue and catch the obvious cases, paired with least-privilege CI credentials — never as the reason a CI job is safe to run unattended. For that, use `--sandbox`.

> [!WARNING]
> **`--yolo` without `--sandbox` runs untrusted, unattended code with no boundary**
>
> A CI job is exactly the environment where a runaway or misled agent does the most damage before anyone notices — no one is at the keyboard to catch a bad `shell` call before it runs, and, per above, a permission allow-list or the shell classifier can't be trusted to catch everything either. If you can't add `--sandbox`, prefer `--safety restricted` with a permission allow-list scoped to what the job actually needs over blanket `--yolo`, and budget for the credentials and blast radius of the agent's toolsets as if the job itself were compromised — see [`examples/permissions.yaml`](https://github.com/docker/docker-agent/blob/main/examples/permissions.yaml) for a worked allow/deny list.

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

A bare OCI registry reference (`myorg/coder`) has no local config you control, so a security-sensitive CI job should check in a small agent config instead. This example runs a checked-in review agent non-interactively against the repository being built:

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

This job auto-approves every shell call the review agent makes (`--yolo`) rather than trying to allow-list every `git`/`grep`/`cat` invocation a code review might need — the read surface for "review this diff" is open-ended, and a fixed pattern list is exactly the kind of shell-matching boundary the [previous section](#defense-in-depth-not-a-boundary-permissions-and-shell-command-matching) says not to rely on. If your CI environment has `sbx` installed and configured (GitHub-hosted `ubuntu-latest` does not ship it out of the box), add `--sandbox` and get a real isolation boundary around that `--yolo`:

```bash
$ docker-agent run --sandbox --exec --yolo .github/agents/review-agent.yaml --json "..."
```

Without `--sandbox`, this workflow's safety instead rests on least-privilege secrets (only `ANTHROPIC_API_KEY` is injected — no repo-write token), the top-level `permissions: contents: read` block and `persist-credentials: false` on the checkout step (which together mean the job never holds a write-capable `GITHUB_TOKEN` and never persists one to disk for `git` to pick up), and the job running on a GitHub-hosted, ephemeral runner that's discarded after the job.

This example omits the GitHub MCP toolset (`docker:github-official`) shown in earlier revisions of this guide: that server requires a `GITHUB_PERSONAL_ACCESS_TOKEN` this workflow doesn't provide, and — because the toolset above has no `name:` field — its tools would be exposed under their raw MCP names (`get_file_contents`, `search_code`, …) rather than a `github_*`-style qualified name, so permission patterns written against that prefix wouldn't match anything anyway. If your review agent needs GitHub API access, add the toolset back with an explicit `name: github`, wire `GITHUB_PERSONAL_ACCESS_TOKEN` through `env:` from a repository secret, and write any `permissions` patterns against the tool names it actually exposes (`github_get_*` only works once the toolset carries that `name:`).

Swap the model, toolsets, and provider secret for your own — the shape (checkout, install the binary, run `--exec` with `--json` against a checked-in config, upload the transcript) generalizes to any CI provider that can run a shell step.
