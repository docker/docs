---
title: "OpenTelemetry Tracing"
description: "Export docker-agent traces to any OTLP backend, including Langfuse and LangSmith, for debugging agentic workflows."
keywords: docker agent, ai agents, community, opentelemetry tracing
weight: 40
canonical: https://docs.docker.com/ai/docker-agent/community/opentelemetry/
---

_docker-agent can export OpenTelemetry traces of an agent run to any OTLP/HTTP backend. This is separate from [product-analytics telemetry](../telemetry/index.md) and is opt-in via the `--otel` flag._

When enabled, docker-agent emits OpenTelemetry GenAI (`gen_ai.*`) and MCP (`mcp.*`) spans following the [OpenTelemetry semantic conventions](https://opentelemetry.io/docs/specs/semconv/gen-ai/). Spans cover the agent turn, model calls (with token usage and cost attributes), tool calls, MCP client/server activity, sub-agent hand-offs, and provider fallbacks. W3C `traceparent` context is propagated so the whole run renders as a single connected trace tree.

## Enabling

```bash
docker agent run agent.yaml --otel
```

Without an exporter endpoint configured, spans are recorded locally as no-ops. To ship them somewhere, set the standard OTLP environment variables described below.

## Configuration

docker-agent reads the standard OTLP environment variables:

| Variable | Purpose |
| --- | --- |
| `OTEL_EXPORTER_OTLP_ENDPOINT` | Base OTLP/HTTP endpoint. The signal subpath (`/v1/traces`, `/v1/metrics`, `/v1/logs`) is appended automatically. |
| `OTEL_EXPORTER_OTLP_HEADERS` | Comma-separated `key=value` headers sent with every export request (for example, an `Authorization` header). |
| `OTEL_RESOURCE_ATTRIBUTES` | Extra resource attributes merged into every span. |
| `OTEL_INSTRUMENTATION_GENAI_CAPTURE_MESSAGE_CONTENT` | Set to `true` to capture prompt and response message content as span attributes. Off by default. |

> [!NOTE]
> **Base endpoint, not the full signal URL**
>
> Set `OTEL_EXPORTER_OTLP_ENDPOINT` to the **base** endpoint (for example `https://cloud.langfuse.com/api/public/otel`). docker-agent appends `/v1/traces` for you, matching the value documented by Langfuse and LangSmith. A bare `host:port` is also accepted and gets `https://` (or `http://` for localhost).

> [!WARNING]
> **Message content can contain sensitive data**
>
> `OTEL_INSTRUMENTATION_GENAI_CAPTURE_MESSAGE_CONTENT` is off by default because chat history routinely contains PII, secrets, and internal documents. Enable it only for backends and environments where exporting that content is acceptable.

## Backends

Protocol support is OTLP over HTTP (`http/protobuf`). gRPC endpoints are not currently supported.

### Langfuse

[Langfuse](https://langfuse.com) exposes an OTLP endpoint and authenticates with HTTP Basic auth built from a project's public and secret keys.

```bash
# Base64 of "public_key:secret_key"
LANGFUSE_AUTH=$(echo -n "pk-lf-...:sk-lf-..." | base64)

export OTEL_EXPORTER_OTLP_ENDPOINT="https://cloud.langfuse.com/api/public/otel"
export OTEL_EXPORTER_OTLP_HEADERS="Authorization=Basic ${LANGFUSE_AUTH}"

docker agent run agent.yaml --otel
```

Regional and self-hosted hosts use the same `/api/public/otel` base path:

| Region | Endpoint |
| --- | --- |
| EU | `https://cloud.langfuse.com/api/public/otel` |
| US | `https://us.cloud.langfuse.com/api/public/otel` |
| Self-hosted (>= v3.22.0) | `http://localhost:3000/api/public/otel` |

### LangSmith

[LangSmith](https://docs.langchain.com/langsmith/trace-with-opentelemetry) authenticates with an `x-api-key` header (the raw API key, with no `Basic`/`Bearer` prefix). An optional `Langsmith-Project` header routes traces to a named project.

```bash
export OTEL_EXPORTER_OTLP_ENDPOINT="https://api.smith.langchain.com/otel"
export OTEL_EXPORTER_OTLP_HEADERS="x-api-key=<your-api-key>,Langsmith-Project=<project>"

docker agent run agent.yaml --otel
```

### OpenTelemetry Collector

Any OTLP/HTTP collector (the OpenTelemetry Collector, Grafana Alloy, Jaeger, and so on) works by pointing at its base endpoint:

```bash
export OTEL_EXPORTER_OTLP_ENDPOINT="http://localhost:4318"
docker agent run agent.yaml --otel
```

> [!NOTE]
> **Langfuse and LangSmith ingest traces only**
>
> Both backends accept the traces signal only. docker-agent also wires metric and log exporters at the same endpoint, so their periodic exports return `404` against trace-only backends. This is harmless to traces but appears in the debug log. Point a full OTLP collector at the endpoint if you also want metrics and logs.

## Inspecting traces locally

Use `--debug` to print telemetry activity to the debug log (`~/.cagent/cagent.debug.log` by default) without standing up a backend:

```bash
docker agent run agent.yaml --otel --debug
```
