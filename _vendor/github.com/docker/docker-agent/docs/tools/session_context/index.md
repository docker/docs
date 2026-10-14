---
title: "Session Context Tool"
description: "Reference a previous session as context in the current one."
keywords: docker agent, ai agents, tools, toolsets, session context tool
linkTitle: "Session Context"
weight: 210
canonical: https://docs.docker.com/ai/docker-agent/tools/session_context/
---

_Reference a previous session as context, without manual export/import._

## Overview

The `session_context` toolset lets an agent discover earlier sessions and pull one in as context for the current session. It removes the manual workaround of exporting a conversation to HTML and re-attaching it with an `@` mention.

The tool surface is two read-only tools:

| Tool            | Description                                                                                                  |
| --------------- | ------------------------------------------------------------------------------------------------------------ |
| `list_sessions` | List previous sessions (most recent first) with id, title, creation time and message count.                  |
| `read_session`  | Return the transcript of a previous session, by id or by a relative reference like `-1`.                      |

The session the agent is currently running in is never listed by `list_sessions` and cannot be read by `read_session` (a circular reference returns an error).

## Configuration

```yaml
toolsets:
  - type: session_context
```

No configuration options. Both tools are read-only and operate against the same session store the runtime already uses for persistence.

Restrict the toolset to a subset of tools the standard way:

```yaml
# An agent that may browse but never pull a full transcript into context.
toolsets:
  - type: session_context
    tools:
      - list_sessions
```

## Selecting a session

`read_session` accepts either form:

- A concrete id returned by `list_sessions`, e.g. `read_session("a1b2c3...")`.
- A relative reference: `-1` is the most recent session, `-2` the second most recent, and so on. Relative references resolve against the same ordering `list_sessions` uses (most recent first), excluding sub-sessions.

## Transcript size

A long session could overflow the current context window, so `read_session` caps the rendered transcript. When a transcript is larger than the budget, the oldest messages are dropped (the most recent are usually the most useful for continuing work) and a note records how many were omitted:

```text
[12 earlier message(s) omitted to fit the context budget; showing the most recent 8]
```

## Notes

- `list_sessions` defaults to 20 sessions and is capped at 100; pass `limit` to request fewer.
- `read_session` returns an error when the session is not found, when the reference cannot be resolved, or when it points at the current session.
- Both tools are read-only: they never modify, branch, or delete sessions.

## Example

See [`examples/session_context.yaml`](https://github.com/docker/docker-agent/blob/main/examples/session_context.yaml) for a complete working example.
