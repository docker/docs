---
title: "Sessions"
description: "How Docker Agent stores conversations, resumes them across runs, and tracks their token cost."
keywords: docker agent, ai agents, features, sessions
weight: 25
canonical: https://docs.docker.com/ai/docker-agent/features/sessions/
---

_How Docker Agent stores conversations, resumes them across runs, and tracks their token cost._

## What a Session Is

Every `docker agent run` creates or resumes a **session**: the record of a conversation, including every message, tool call, sub-agent run, and cost. Sessions are what makes `/undo`, `/sessions`, `--session <id>`, and cost tracking possible — the agent itself is stateless between runs, but the session is not.

A session is only persisted to disk once it has content (the first message added), so a run you cancel before sending anything never leaves an empty session behind.

## Where Sessions Are Stored

Sessions live in a SQLite database, `session.db`, under the [data directory](../cli/index.md#global-flags) (`~/.cagent` by default):

```bash
$ ls ~/.cagent/session.db
```

Override the location with `-s`/`--session-db`, or by overriding the data directory itself with `--data-dir`:

```bash
# Use a project-local session database instead of the global one
$ docker agent run agent.yaml --session-db ./sessions.db
```

## Resuming a Session

Pass `--session <id>` to continue a previous conversation instead of starting a new one:

```bash
# Resume by explicit session ID
$ docker agent run agent.yaml --session 3f9c1e2a-...

# Resume the most recently created session
$ docker agent run agent.yaml --session -1

# Resume the session created before that one
$ docker agent run agent.yaml --session -2
```

`--session` accepts two kinds of reference:

- **A relative offset** (`-1`, `-2`, …): sessions are ordered by **creation time**, newest first — `-1` is the most recently *created* session, `-2` the one created before it, and so on. This is creation order, not last-used order: resuming an older session with `--session <id>` does not make it the new `-1`; the next `-1` still resolves to whichever session was created most recently. A relative offset that has no matching session (for example `-1` with an empty database) is an error.
- **An explicit ID**: if a session with that ID already exists, it is resumed. If it doesn't exist yet, docker agent **creates it with that ID** instead of failing. This lets a supervisor (for example a board or a script) choose a session ID up front and reuse it across runs — the first run creates the session, later runs resume it.

## Read-Only Sessions

Add `--session-read-only` to open a session for viewing without sending new messages — useful for reviewing a past conversation without accidentally continuing it:

```bash
$ docker agent run agent.yaml --session -1 --session-read-only
```

`--session-read-only` requires the TUI: it cannot be combined with `--exec`, since there would be nothing to display without one.

## Browsing Sessions in the TUI

Press `/sessions` to open the session browser: search and filter past conversations, see which ones were started in the current working directory ("This workspace") versus elsewhere ("Other locations"), and restore one with <kbd>Enter</kbd>. Restoring a session reopens it in its original working directory. See [Session Management](../tui/index.md#session-management) in the Terminal UI docs for the full set of session-browser and session-title features (starring, branching by editing a past message, and so on).

### Tabs

<kbd>Ctrl</kbd>+<kbd>T</kbd> opens a new tab running an additional agent session alongside the current one; <kbd>Ctrl</kbd>+<kbd>N</kbd>/<kbd>Ctrl</kbd>+<kbd>P</kbd> cycle between tabs and <kbd>Ctrl</kbd>+<kbd>W</kbd> closes the current one. By default, tabs are not restored the next time you launch the TUI. Set `restore_tabs: true` in your user config to reopen the same tabs (and their sessions) on the next launch:

```yaml
# ~/.config/cagent/config.yaml
settings:
  restore_tabs: true
```

## Session Titles

Docker Agent auto-generates a short title for each session from your first message, using a one-shot call to the agent's own model. Point that call at a smaller, cheaper model instead with `title_model` on a model definition:

```yaml
# examples/title_model.yaml
models:
  primary:
    provider: anthropic
    model: claude-sonnet-4-5
    # Generate session titles with the cheaper Haiku model instead of Sonnet.
    title_model: fast
  fast:
    provider: anthropic
    model: claude-haiku-4-5

agents:
  root:
    model: primary
    description: An assistant that generates session titles with a cheaper model.
    instruction: You are a helpful assistant.
```

When `title_model` is omitted, title generation reuses the agent's own model. Set or regenerate a title from inside the TUI with `/title` (see [Session Title Editing](../tui/index.md#session-title-editing)) — regenerating sends every user message in the session so far, not just the first — or generate one from the command line without starting a session at all:

```bash
$ docker agent debug title agent.yaml "How do I configure a fallback model?"
```

See [`docker agent debug title`](../cli/index.md#docker-agent-debug) for details.

## Usage & Cost Tracking

Every tracked model call — the main conversation turns and compaction calls — updates the session's cumulative input/output token counts and cost. Check them at any time with `/cost` in the TUI, or disable tracking per-model with `track_usage: false` if you don't want a model's calls counted (for example, a free local model). Auxiliary one-shot calls the runtime makes on your behalf, such as automatic session-title generation, call the model directly and are not folded into this total.

Cost is computed from the [models.dev](https://models.dev/) pricing catalogue by default. For a custom endpoint, a private deployment, or a negotiated enterprise rate the catalogue doesn't know about, declare pricing explicitly with a model's `cost:` block:

```yaml
# examples/custom-pricing.yaml
models:
  internal-gpt:
    provider: internal-llm
    model: gpt-4o
    cost:
      input: 1.25 # USD per 1M input tokens
      output: 5.00 # USD per 1M output tokens
      cache_read: 0.125 # USD per 1M cached input tokens
      cache_write: 1.5625 # USD per 1M cache-write tokens
```

An all-zero `cost:` table means "priced, free" — distinct from a model with no `cost:` at all, which falls back to the catalogue (and bills $0 for models the catalogue doesn't know).

> [!NOTE]
> **Cost never decreases**
>
> A session's cumulative cost is a running total updated after every *tracked* model call — the same main conversation turns and compaction calls described above. Compacting the conversation (manually with `/compact`, or automatically — see the agent's `session_compaction`/`compaction_threshold` fields in the [Agent Config reference](../../configuration/agents/index.md)) reshapes the message history sent back to the model, but never touches this running total. Auxiliary one-shot calls such as automatic session-title generation are not tracked and are not reflected in this total: `/cost` covers everything the session's tracked calls have spent, not literally every model call the runtime makes on your behalf.

## Resuming Into a Worktree

A session created during a `--worktree` run remembers which worktree it used. Resuming it with `--session` reattaches to the same worktree directory and branch automatically — you don't need to pass `--worktree` again. See [`--worktree`](../cli/index.md#docker-agent-run) in the CLI reference for the full worktree lifecycle.
