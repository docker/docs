---
title: OpenTelemetry for the Docker CLI
description: Learn about how to capture OpenTelemetry metrics for the Docker command line
keywords: otel, opentelemetry, telemetry, traces, tracing, metrics, logs
aliases:
  - /config/otel/
---

{{< summary-bar feature_name="Docker CLI OpenTelemetry" >}}

The Docker CLI supports [OpenTelemetry](https://opentelemetry.io/docs/) instrumentation
for emitting metrics about command invocations. This is disabled by default.
You can configure the CLI to start emitting metrics to the endpoint that you
specify. This allows you to capture information about your `docker` command
invocations for more insight into your Docker usage.

Exporting metrics is opt-in, and you control where data is being sent by
specifying the destination address of the metrics collector.

## What is OpenTelemetry?

OpenTelemetry, or OTel for short, is an open observability framework for
creating and managing telemetry data, such as traces, metrics, and logs.
OpenTelemetry is vendor- and tool-agnostic, meaning that it can be used with a
broad variety of Observability backends.

Support for OpenTelemetry instrumentation in the Docker CLI means that the CLI can emit
information about events that take place, using the protocols and conventions
defined in the Open Telemetry specification.

## How it works

The Docker CLI doesn't emit telemetry data by default. Only if you've set an
environment variable on your system will Docker CLI attempt to emit OpenTelemetry
metrics, to the endpoint that you specify.

```bash
DOCKER_CLI_OTEL_EXPORTER_OTLP_ENDPOINT=<endpoint>
```

The variable specifies the endpoint of an OpenTelemetry collector, where telemetry data
about `docker` CLI invocation should be sent. To capture the data, you'll need
an OpenTelemetry collector listening on that endpoint.

The purpose of a collector is to receive the telemetry data, process it, and
exports it to a backend. The backend is where the telemetry data gets stored.
You can choose from a number of different backends, such as Prometheus or
InfluxDB.

Some backends provide tools for visualizing the metrics directly.
Alternatively, you can also run a dedicated frontend with support for
generating more useful graphs, such as Grafana.

## Setup

To get started capturing telemetry data for the Docker CLI, you'll need to:

- Set the `DOCKER_CLI_OTEL_EXPORTER_OTLP_ENDPOINT` environment variable to point to an OpenTelemetry collector endpoint
- Run an OpenTelemetry collector that receives the signals from CLI command invocations
- Run a backend for storing the data received from the collector

The following Docker Compose file bootstraps a set of services to get started with OpenTelemetry.
It includes an OpenTelemetry collector that the CLI can send metrics to,
and a Prometheus backend that scrapes the metrics off the collector.

```yaml {collapse=true,title=compose.yml}
name: cli-otel
services:
  prometheus:
    image: prom/prometheus
    command:
      - "--config.file=/etc/prometheus/prom.yml"
    ports:
      # Publish the Prometheus frontend on localhost:9091
      - 9091:9090
    restart: always
    volumes:
      # Store Prometheus data in a volume:
      - prom_data:/prometheus
      # Mount the prom.yml config file
      - ./prom.yml:/etc/prometheus/prom.yml
  otelcol:
    image: otel/opentelemetry-collector
    restart: always
    depends_on:
      - prometheus
    ports:
      - 4317:4317
    volumes:
      # Mount the otelcol.yml config file
      - ./otelcol.yml:/etc/otelcol/config.yaml

volumes:
  prom_data:
```

This service assumes that the following two configuration files exist alongside
`compose.yml`:

- ```yaml {collapse=true,title=otelcol.yml}
  # Receive signals over gRPC and HTTP
  receivers:
    otlp:
      protocols:
        grpc:
        http:

  # Establish an endpoint for Prometheus to scrape from
  exporters:
    prometheus:
      endpoint: "0.0.0.0:8889"

  service:
    pipelines:
      metrics:
        receivers: [otlp]
        exporters: [prometheus]
  ```

- ```yaml {collapse=true,title=prom.yml}
  # Configure Prometheus to scrape the OpenTelemetry collector endpoint
  scrape_configs:
    - job_name: "otel-collector"
      scrape_interval: 1s
      static_configs:
        - targets: ["otelcol:8889"]
  ```

With these files in place:

1. Start the Docker Compose services:

   ```console
   $ docker compose up
   ```

2. Configure Docker CLI to export telemetry to the OpenTelemetry collector.

   ```console
   $ export DOCKER_CLI_OTEL_EXPORTER_OTLP_ENDPOINT=http://localhost:4317
   ```

3. Run a `docker` command to trigger the CLI into sending a metric signal to
   the OpenTelemetry collector.

   ```console
   $ docker version
   ```

4. To view telemetry metrics created by the CLI, open the Prometheus expression
   browser by going to <http://localhost:9091/graph>.

5. In the **Query** field, enter `command_time_milliseconds_total`, and execute
   the query to see the telemetry data.

## Available metrics

Docker CLI currently exports a single metric, `command.time`, which measures
the execution duration of a command in milliseconds. This metric has the
following attributes:

- `command.name`: the name of the command
- `command.status.code`: the exit code of the command
- `command.stderr.isatty`: true if stderr is attached to a TTY
- `command.stdin.isatty`: true if stdin is attached to a TTY
- `command.stdout.isatty`: true if stdout is attached to a TTY
