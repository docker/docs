---
description: Collecting Docker metrics with Prometheus
keywords: prometheus, metrics
title: Collect Docker metrics with Prometheus
aliases:
  - /engine/admin/prometheus/
  - /config/thirdparty/monitoring/
  - /config/thirdparty/prometheus/
  - /config/daemon/prometheus/
---

[Prometheus](https://prometheus.io/) is an open-source systems monitoring and
alerting toolkit. You can configure Docker as a Prometheus target.

> [!WARNING]
>
> The available metrics and the names of those metrics are in active
> development and may change at any time.

Currently, you can only monitor Docker itself. You can't currently monitor your
application using the Docker target.

## Example

The following example shows you how to configure your Docker daemon, set up
Prometheus to run as a container on your local machine, and monitor your Docker
instance using Prometheus.

### Configure the daemon

To configure the Docker daemon as a Prometheus target, you need to specify the
`metrics-address` in the `daemon.json` configuration file. This daemon expects
the file to be located at one of the following locations by default. If the
file doesn't exist, create it.

- **Linux**: `/etc/docker/daemon.json`
- **Windows Server**: `C:\ProgramData\docker\config\daemon.json`
- **Docker Desktop**: Open the Docker Desktop settings and select **Docker Engine** to edit the file.

Add the following configuration:

```json
{
  "metrics-addr": "127.0.0.1:9323"
}
```

Save the file, or in the case of Docker Desktop for Mac or Docker Desktop for
Windows, save the configuration. Restart Docker.

Docker now exposes Prometheus-compatible metrics on port 9323 via the loopback
interface. You can configure it to use the wildcard address `0.0.0.0` instead,
but this will expose the Prometheus port to the wider network. Consider your
threat model carefully when deciding which option best suits your environment.

### Create a Prometheus configuration

Copy the following configuration file and save it to a location of your choice,
for example `/tmp/prometheus.yml`. This is a stock Prometheus configuration file,
except for the addition of the Docker job definition at the bottom of the file.

```yml
# my global config
global:
  scrape_interval: 15s # Set the scrape interval to every 15 seconds. Default is every 1 minute.
  evaluation_interval: 15s # Evaluate rules every 15 seconds. The default is every 1 minute.
  # scrape_timeout is set to the global default (10s).

  # Attach these labels to any time series or alerts when communicating with
  # external systems (federation, remote storage, Alertmanager).
  external_labels:
    monitor: "codelab-monitor"

# Load rules once and periodically evaluate them according to the global 'evaluation_interval'.
rule_files:
  # - "first.rules"
  # - "second.rules"

# A scrape configuration containing exactly one endpoint to scrape:
# Here it's Prometheus itself.
scrape_configs:
  # The job name is added as a label `job=<job_name>` to any timeseries scraped from this config.
  - job_name: prometheus

    # metrics_path defaults to '/metrics'
    # scheme defaults to 'http'.

    static_configs:
      - targets: ["localhost:9090"]

  - job_name: docker
      # metrics_path defaults to '/metrics'
      # scheme defaults to 'http'.

    static_configs:
      - targets: ["host.docker.internal:9323"]
```

### Run Prometheus in a container

Next, start a Prometheus container using this configuration.

```console
$ docker run --name my-prometheus \
    --mount type=bind,source=/tmp/prometheus.yml,destination=/etc/prometheus/prometheus.yml \
    -p 9090:9090 \
    --add-host host.docker.internal=host-gateway \
    prom/prometheus
```

If you're using Docker Desktop, the `--add-host` flag is optional. This flag
makes sure that the host's internal IP gets exposed to the Prometheus
container. Docker Desktop does this by default. The host IP is exposed as the
`host.docker.internal` hostname. This matches the configuration defined in
`prometheus.yml` in the previous step.

### Open the Prometheus Dashboard

Verify that the Docker target is listed at `http://localhost:9090/targets/`.

![Prometheus targets page](images/prometheus-targets.webp)

> [!NOTE]
>
> You can't access the endpoint URLs on this page directly if you use Docker
> Desktop.

### Use Prometheus

Create a graph. Select the **Graphs** link in the Prometheus UI. Choose a metric
from the combo box to the right of the **Execute** button, and click
**Execute**. The screenshot below shows the graph for
`engine_daemon_network_actions_seconds_count`.

![Idle Prometheus report](images/prometheus-graph_idle.webp)

The graph shows a pretty idle Docker instance, unless you're already running
active workloads on your system.

To make the graph more interesting, run a container that uses some network
actions by starting downloading some packages using a package manager:

```console
$ docker run --rm alpine apk add git make musl-dev go
```

Wait a few seconds (the default scrape interval is 15 seconds) and reload your
graph. You should see an uptick in the graph, showing the increased network
traffic caused by the container you just ran.

![Prometheus report showing traffic](images/prometheus-graph_load.webp)

## Available metrics

Docker exposes metrics in Prometheus format. This section describes the available metrics and their meaning.

> [!WARNING]
>
> The available metrics and the names of those metrics are in active
> development and may change at any time.

### Metric types

Docker metrics use the following Prometheus metric types:

- **Counter**: A cumulative metric that only increases (or resets to zero on restart). Use counters for values like total number of events or requests.
- **Gauge**: A metric that can go up or down. Use gauges for values like current memory usage or number of running containers.
- **Histogram**: A metric that samples observations and counts them in configurable buckets. Histograms expose multiple time series:
  - `<basename>_bucket{le="<upper_bound>"}`: Cumulative counters for observation buckets
  - `<basename>_sum`: Total sum of all observed values
  - `<basename>_count`: Count of events that have been observed

For histogram metrics, you can calculate averages, percentiles, and rates. For example, to calculate the average duration: `rate(<basename>_sum[5m]) / rate(<basename>_count[5m])`.

### Engine metrics

These metrics provide information about the Docker Engine's operation and resource usage.

| Metric                                      | Type      | Description                                                                                                                  |
| ------------------------------------------- | --------- | ---------------------------------------------------------------------------------------------------------------------------- |
| `engine_daemon_container_actions_seconds`   | Histogram | Time taken to process container operations (start, stop, create, etc.). Labels indicate the action type.                     |
| `engine_daemon_container_states_containers` | Gauge     | Number of containers currently in each state (running, paused, stopped). Labels indicate the state.                          |
| `engine_daemon_engine_cpus_cpus`            | Gauge     | Number of CPUs available on the host system.                                                                                 |
| `engine_daemon_engine_info`                 | Gauge     | Static information about the Docker Engine. Always set to 1. Labels provide version, architecture, and other engine details. |
| `engine_daemon_engine_memory_bytes`         | Gauge     | Total memory available on the host system in bytes.                                                                          |
| `engine_daemon_events_subscribers_total`    | Gauge     | Number of current subscribers to Docker events.                                                                              |
| `engine_daemon_events_total`                | Counter   | Total number of events processed by the daemon. Labels indicate the event action and type.                                   |
| `engine_daemon_health_checks_failed_total`  | Counter   | Total number of health checks that have failed.                                                                              |
| `engine_daemon_health_checks_total`         | Counter   | Total number of health checks performed.                                                                                     |
| `engine_daemon_host_info_functions_seconds` | Histogram | Time taken to gather host information.                                                                                       |
| `engine_daemon_network_actions_seconds`     | Histogram | Time taken to process network operations (create, connect, disconnect, etc.). Labels indicate the action type.               |

### Swarm metrics

These metrics are only available when the Docker Engine is running in Swarm mode.

| Metric                                           | Type      | Description                                                                                     |
| ------------------------------------------------ | --------- | ----------------------------------------------------------------------------------------------- |
| `swarm_dispatcher_scheduling_delay_seconds`      | Histogram | Time from task creation to scheduling decision. Measures scheduler performance.                 |
| `swarm_manager_configs_total`                    | Gauge     | Total number of configs in the swarm cluster.                                                   |
| `swarm_manager_leader`                           | Gauge     | Indicates if this node is the swarm manager leader (1) or not (0).                              |
| `swarm_manager_networks_total`                   | Gauge     | Total number of networks in the swarm cluster.                                                  |
| `swarm_manager_nodes`                            | Gauge     | Number of nodes in the swarm cluster. Labels indicate node state (ready, down, etc.).           |
| `swarm_manager_secrets_total`                    | Gauge     | Total number of secrets in the swarm cluster.                                                   |
| `swarm_manager_services_total`                   | Gauge     | Total number of services in the swarm cluster.                                                  |
| `swarm_manager_tasks_total`                      | Gauge     | Total number of tasks in the swarm cluster. Labels indicate task state (running, failed, etc.). |
| `swarm_node_manager`                             | Gauge     | Indicates if this node is a swarm manager (1) or worker (0).                                    |
| `swarm_raft_snapshot_latency_seconds`            | Histogram | Time taken to create and restore Raft snapshots.                                                |
| `swarm_raft_transaction_latency_seconds`         | Histogram | Time taken to commit Raft transactions. Measures consensus performance.                         |
| `swarm_store_batch_latency_seconds`              | Histogram | Time taken for batch operations in the swarm store.                                             |
| `swarm_store_lookup_latency_seconds`             | Histogram | Time taken for lookup operations in the swarm store.                                            |
| `swarm_store_memory_store_lock_duration_seconds` | Histogram | Duration of lock acquisitions in the memory store.                                              |
| `swarm_store_read_tx_latency_seconds`            | Histogram | Time taken for read transactions in the swarm store.                                            |
| `swarm_store_write_tx_latency_seconds`           | Histogram | Time taken for write transactions in the swarm store.                                           |

### Using histogram metrics

For histogram metrics (those with `_seconds` in the name), Prometheus creates three time series:

- `<metric_name>_bucket`: Cumulative counters for each configured bucket
- `<metric_name>_sum`: Total sum of all observed values
- `<metric_name>_count`: Total count of observations

For example, `engine_daemon_container_actions_seconds` produces:

- `engine_daemon_container_actions_seconds_bucket{action="start",le="0.005"}`: Count of start actions taking ≤5ms
- `engine_daemon_container_actions_seconds_bucket{action="start",le="0.01"}`: Count of start actions taking ≤10ms
- `engine_daemon_container_actions_seconds_sum{action="start"}`: Total time spent on start actions
- `engine_daemon_container_actions_seconds_count{action="start"}`: Total number of start actions

Use these to calculate percentiles, averages, and rates in your Prometheus queries.

## Next steps

The example provided here shows how to run Prometheus as a container on your
local system. In practice, you'll probably be running Prometheus on another
system or as a cloud service somewhere. You can set up the Docker daemon as a
Prometheus target in such contexts too. Configure the `metrics-addr` of the
daemon and add the address of the daemon as a scrape endpoint in your
Prometheus configuration.

```yaml
- job_name: docker
  static_configs:
    - targets: ["docker.daemon.example:<PORT>"]
```

For more information about Prometheus, refer to the
[Prometheus documentation](https://prometheus.io/docs/introduction/overview/)
