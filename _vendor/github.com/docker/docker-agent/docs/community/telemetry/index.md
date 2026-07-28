---
title: "Telemetry"
description: "docker-agent collects anonymous usage data to help improve the tool. Telemetry can be disabled at any time."
keywords: docker agent, ai agents, community, telemetry
weight: 30
canonical: https://docs.docker.com/ai/docker-agent/community/telemetry/
---

_docker-agent collects anonymous usage data to help improve the tool. Telemetry can be disabled at any time._

On first startup, docker-agent displays a notice about telemetry collection so you're always informed. All events are processed synchronously when recorded.

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

- Command names and success/failure status
- Agent names and model types
- Tool names and whether calls succeed or fail
- Token counts (input/output totals) and estimated costs
- Session metadata (durations, error counts)

## What's NOT Collected ❌

- User input or prompts
- Agent responses or generated content
- File contents
- API keys or credentials
- Personally identifying information (PII)

> [!TIP]
> **See events locally**
>
> Use `--debug` to see telemetry events printed to the debug log without sending them anywhere additional.

```bash
docker agent run agent.yaml --debug
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

```bash
// Recommended: context-based telemetry (clean, testable)
if telemetryClient := telemetry.FromContext(ctx); telemetryClient != nil {
    telemetryClient.RecordToolCall(ctx, "filesystem", "session-id", "agentName", time.Millisecond*500, nil)
    telemetryClient.RecordTokenUsage(ctx, "gpt-4", 100, 50, 0.01)
}

// Or use direct calls
telemetry.TrackCommand("run", args)
```

Events are processed synchronously when `Track()` is called, sending HTTP requests immediately.
