---
title: "Telemetry"
description: "Docker Agent collects usage data to help improve the tool. Telemetry can be disabled at any time."
keywords: docker agent, ai agents, configuration, telemetry
weight: 120
canonical: https://docs.docker.com/ai/docker-agent/configuration/telemetry/
aliases:
  - /ai/docker-agent/community/telemetry/
---

_Docker Agent collects usage data to help improve the tool. Telemetry can be disabled at any time._

On first startup, Docker Agent displays a notice about telemetry collection so you're always informed. Ordinary telemetry events are sent asynchronously; selected error events are sent synchronously before exit.

## Disabling Telemetry

```bash
# Disable via environment variable
$ TELEMETRY_ENABLED=false docker agent run agent.yaml

# Or export it in your shell profile
$ export TELEMETRY_ENABLED=false
```

> [!NOTE]
> **Default**
>
> Telemetry is **enabled by default**. Set `TELEMETRY_ENABLED=false` to opt out.

## What's Collected ✅

- Command names, positional arguments, and success/failure status
- Agent names and model types
- Tool names and whether calls succeed or fail
- Token counts (input/output totals) and estimated costs
- Session metadata (durations, error counts)

## Sensitive Data

Command events include positional arguments. For `run` and `exec`, these can include prompts, file paths, and registry references. Error events also include error text. Do not assume telemetry is free of secrets or personally identifying information: either can appear in arguments or errors.

Conversation transcripts and file contents are not collected as separate telemetry fields, but text included in command arguments or errors can still be transmitted. Set `TELEMETRY_ENABLED=false` before running commands that may contain sensitive data.

> [!TIP]
> **See events locally**
>
> Use `--debug` to see telemetry events in the debug log. Enabled telemetry is still transmitted, using Docker's staging telemetry endpoint in debug mode. Combine it with `TELEMETRY_ENABLED=false` to inspect events without transmitting them.

```bash
TELEMETRY_ENABLED=false docker agent run agent.yaml --debug
```

## Event Types

The telemetry system uses structured, type-safe events:

| Event Type  | What It Tracks                                                      |
| ----------- | ------------------------------------------------------------------- |
| **Command** | CLI command execution with success status                           |
| **Tool**    | Agent tool calls with timing and error information                  |
| **Token**   | LLM token usage by model, session, and cost                         |
| **Session** | Agent session lifecycle with start/end events and aggregate metrics |

## For Developers

Telemetry is automatically wrapped around all commands. To record additional events, use the context-based API:

```go
// Recommended: context-based telemetry (clean, testable)
if telemetryClient := telemetry.FromContext(ctx); telemetryClient != nil {
    telemetryClient.RecordToolCall(ctx, "filesystem", "session-id", "agentName", time.Millisecond*500, nil)
    telemetryClient.RecordTokenUsage(ctx, "gpt-4", 100, 50, 0.01)
}

// Or use direct calls
telemetry.TrackCommand(ctx, "run", args)
```

`Track()` prepares the event and sends it asynchronously. `TrackSynchronous()` waits for the HTTP request; command-error tracking uses it before the process exits.
