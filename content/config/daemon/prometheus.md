---
description: Collecting Docker metrics with Prometheus
keywords: prometheus, metrics
title: Collect Docker metrics with Prometheus
aliases:
  - /engine/admin/prometheus/
  - /config/thirdparty/monitoring/
  - /config/thirdparty/prometheus/
---

[Prometheus](https://prometheus.io/) is an open-source systems monitoring and
alerting toolkit. You can configure Docker as a Prometheus target. This topic
shows you how to configure Docker, set up Prometheus to run as a Docker
container, and monitor your Docker instance using Prometheus.

> **Warning**
>
> The available metrics and the names of those metrics are in active
> development and may change at any time.
{ .warning }

Currently, you can only monitor Docker itself. You can't currently monitor your
application using the Docker target.

## Prerequisites

1.  One or more Docker engines are joined into a Docker Swarm, using `docker
swarm init` on one manager and `docker swarm join` on other managers and
    worker nodes.
2.  You need an internet connection to pull the Prometheus image.

## Configure Docker

To configure the Docker daemon as a Prometheus target, you need to specify the
`metrics-address`. The best way to do this is via the `daemon.json`, which is
located at one of the following locations by default. If the file doesn't
exist, create it.

- **Linux**: `/etc/docker/daemon.json`
- **Windows Server**: `C:\ProgramData\docker\config\daemon.json`
- **Docker Desktop**: Open the Docker Desktop settings and select **Docker Engine**.

If the file is currently empty, paste the following:

```json
{
  "metrics-addr": "127.0.0.1:9323"
}
```

If the file isn't empty, add the new key, making sure that the resulting
file is valid JSON. Be careful that every line ends with a comma (`,`) except
for the last line.

Save the file, or in the case of Docker Desktop for Mac or Docker Desktop for Windows, save the
configuration. Restart Docker.

Docker now exposes Prometheus-compatible metrics on port 9323.

## Configure and run Prometheus

Prometheus runs as a Docker service on a Swarm.

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
  - job_name: "prometheus"

    # metrics_path defaults to '/metrics'
    # scheme defaults to 'http'.

    static_configs:
      - targets: ["host.docker.internal:9090"]

  - job_name:
      "docker"
      # metrics_path defaults to '/metrics'
      # scheme defaults to 'http'.

    static_configs:
      - targets: ["localhost:9323"]
```

Next, start a single-replica Prometheus service using this configuration.

- If you're using Docker Desktop, run:

  ```console
  $ docker service create --replicas 1 --name my-prometheus \
      --mount type=bind,source=/tmp/prometheus.yml,destination=/etc/prometheus/prometheus.yml \
      --publish published=9090,target=9090,protocol=tcp \
      prom/prometheus
  ```

- If you're using Docker Engine without Docker Desktop, run:

  ```console
  $ docker service create --replicas 1 --name my-prometheus \
      --mount type=bind,source=/tmp/prometheus.yml,destination=/etc/prometheus/prometheus.yml \
      --publish published=9090,target=9090,protocol=tcp \
      --add-host host.docker.internal:host-gateway \
      prom/prometheus
  ```

Verify that the Docker target is listed at `http://localhost:9090/targets/`.

![Prometheus targets page](images/prometheus-targets.png)

You can't access the endpoint URLs directly if you use Docker Desktop
for Mac or Docker Desktop for Windows.

## Use Prometheus

Create a graph. Select the **Graphs** link in the Prometheus UI. Choose a metric
from the combo box to the right of the **Execute** button, and click
**Execute**. The screenshot below shows the graph for
`engine_daemon_network_actions_seconds_count`.

![Prometheus engine_daemon_network_actions_seconds_count report](images/prometheus-graph_idle.png)

The above graph shows a pretty idle Docker instance. Your graph might look
different if you are running active workloads.

To make the graph more interesting, create some network actions by starting
a service with 10 tasks that just ping Docker non-stop (you can change the
ping target to anything you like):

```console
$ docker service create \
  --replicas 10 \
  --name ping_service \
  alpine ping docker.com
```

Wait a few minutes (the default scrape interval is 15 seconds) and reload
your graph.

![Prometheus engine_daemon_network_actions_seconds_count report](images/prometheus-graph_load.png)

When you are ready, stop and remove the `ping_service` service, so that you
aren't flooding a host with pings for no reason.

```console
$ docker service remove ping_service
```

Wait a few minutes and you should see that the graph falls back to the idle
level.

## Next steps

- Read the [Prometheus documentation](https://prometheus.io/docs/introduction/overview/)
- Set up some [alerts](https://prometheus.io/docs/alerting/overview/)
