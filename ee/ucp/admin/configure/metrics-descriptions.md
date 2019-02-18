---
description: Using UCP cluster metrics with Prometheus
keywords: prometheus, metrics, ucp
title: Using UCP cluster metrics with Prometheus
redirect_from:
- /engine/admin/prometheus/
---

# UCP metrics

The following table lists the metrics that UCP exposes in Prometheus, along with descriptions. Note that only the metrics 
labeled with `ucp_` are documented. Other metrics are exposed in Prometheus but are not documented.

| Name                                                    | Units                | Description                                                                                                                                                                                                                                                                     | Labels                                         | Metric source |
|---------------------------------------------------------|----------------------|---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|------------------------------------------------|---------------|
| `ucp_controller_services`                               | number of services   | The total number of Swarm services                                                                                                                                                                                                                                              |                                                | Controller    |
| `ucp_engine_container_cpu_percent`                      | percentage           | The percentage of CPU time this container is using.                                                                                                                                                                                                                             | container labels                               | Node          |
| `ucp_engine_container_cpu_total_time_nanoseconds`       | nanoseconds          | Total CPU time used by this container in nanoseconds                                                                                                                                                                                                                            | container labels                               | Node          |
| `ucp_engine_container_health`                           | 0.0 or 1.0           | Whether or not this container is healthy, according to its healthcheck. Note that if this value is 0, it just means that the container is not reporting healthy; it might not have a healthcheck defined at all, or its healthcheck might not have returned any results yet     | container labels                               | Node          |
| `ucp_engine_container_memory_max_usage_bytes`           | bytes                | Maximum memory used by this container in bytes                                                                                                                                                                                                                                  | container labels                               | Node          |
| `ucp_engine_container_memory_usage_bytes`               | bytes                | Current memory used by this container in bytes                                                                                                                                                                                                                                  | container labels                               | Node          |
| `ucp_engine_container_memory_usage_percent`             | percentage           | Percentage of total node memory currently being used by this container                                                                                                                                                                                                          | container labels                               | Node          |
| `ucp_engine_container_network_rx_bytes_total`           | bytes                | Number of bytes received by this container on this network in the last sample                                                                                                                                                                                                   | container networking labels                    | Node          |
| `ucp_engine_container_network_rx_dropped_packets_total` | number of packets    | Number of packets bound for this container on this network that were dropped in the last sample                                                                                                                                                                                 | container networking labels                    | Node          |
| `ucp_engine_container_network_rx_errors_total`          | number of errors     | Number of received network errors for this container on this network in the last sample                                                                                                                                                                                         | container networking labels                    | Node          |
| `ucp_engine_container_network_rx_packets_total`         | number of packets    | Number of received packets for this container on this network in the last sample                                                                                                                                                                                                | container networking labels                    | Node          |
| `ucp_engine_container_network_tx_bytes_total`           | bytes                | Number of bytes sent by this container on this network in the last sample                                                                                                                                                                                                       | container networking labels                    | Node          |
| `ucp_engine_container_network_tx_dropped_packets_total` | number of packets    | Number of packets sent from this container on this network that were dropped in the last sample                                                                                                                                                                                 | container networking labels                    | Node          |
| `ucp_engine_container_network_tx_errors_total`          | number of errors     | Number of sent network errors for this container on this network in the last sample                                                                                                                                                                                             | container networking labels                    | Node          |
| `ucp_engine_container_network_tx_packets_total`         | number of packets    | Number of sent packets for this container on this network in the last sample                                                                                                                                                                                                    | container networking labels                    | Node          |
| `ucp_engine_container_unhealth`                         | 0.0 or 1.0           | Whether or not this container is unhealthy, according to its healthcheck. Note that if this value is 0, it just means that the container is not reporting unhealthy; it might not have a healthcheck defined at all, or its healthcheck might not have returned any results yet | container labels                               | Node          |
| `ucp_engine_containers`                                 | number of containers | Total number of containers on this node                                                                                                                                                                                                                                         | node labels                                    | Node          |
| `ucp_engine_cpu_total_time_nanoseconds`                 | nanoseconds          | System CPU time used by this container in nanoseconds                                                                                                                                                                                                                           | container labels                               | Node          |
| `ucp_engine_disk_free_bytes`                            | bytes                | Free disk space on the Docker root directory on this node in bytes. Note that this metric is not available for Windows nodes                                                                                                                                                    | node labels                                    | Node          |
| `ucp_engine_disk_total_bytes`                           | bytes                | Total disk space on the Docker root directory on this node in bytes. Note that this metric is not available for Windows nodes                                                                                                                                                   | node labels                                    | Node          |
| `ucp_engine_images`                                     | number of images     | Total number of images on this node                                                                                                                                                                                                                                             | node labels                                    | Node          |
| `ucp_engine_memory_total_bytes`                         | bytes                | Total amount of memory on this node in bytes                                                                                                                                                                                                                                    | node labels                                    | Node          |
| `ucp_engine_networks`                                   | number of networks   | Total number of networks on this node                                                                                                                                                                                                                                           | node labels                                    | Node          |
| `ucp_engine_node_health`                                | 0.0 or 1.0           | Whether or not this node is healthy, as determined by UCP                                                                                                                                                                                                                       | nodeName: node name, nodeAddr: node IP address | Controller    |
| `ucp_engine_num_cpu_cores`                              | number of cores      | Number of CPU cores on this node                                                                                                                                                                                                                                                | node labels                                    | Node          |
| `ucp_engine_pod_container_ready`                        | 0.0 or 1.0           | Whether or not this container in a Kubernetes pod is ready, as determined by its readiness probe.                                                                                                                                                                               | pod labels                                     | Controller    |
| `ucp_engine_pod_ready`                                  | 0.0 or 1.0           | Whether or not this Kubernetes pod is ready, as determined by its readiness probe.                                                                                                                                                                                              | pod labels                                     | Controller    |
| `ucp_engine_volumes`                                    | number of volumes    | Total number of volumes on this node                                                                                                                                                                                                                                            | node labels                                    | Node          |

## Metrics labels

Metrics exposed by UCP in Prometheus have standardized labels, depending on the resource that they are measuring. 
The following table lists some of the labels that are used, along with their values:

### Container labels

| Label name         | Value                                                                                       |
|--------------------|---------------------------------------------------------------------------------------------|
| `collection`       | The collection ID of the collection this container is in, if any                            |
| `container`        | The ID of this container                                                                    |
| `image`            | The name of this container's image                                                          |
| `manager`          | "true" if the container's node is a UCP manager, "false" otherwise                          |
| `name`             | The name of the container                                                                   |
| `podName`          | If this container is part of a Kubernetes pod, this is the pod's name                       |
| `podNamespace`     | If this container is part of a Kubernetes pod, this is the pod's namespace                  |
| `podContainerName` | If this container is part of a Kubernetes pod, this is the container's name in the pod spec |
| `service`          | If this container is part of a Swarm service, this is the service ID                        |
| `stack`            | If this container is part of a Docker compose stack, this is the name of the stack          |

### Container networking labels

The following metrics measure network activity for a given network attached to a given
container. They have the same labels as Container labels, with one addition:

| Label name | Value                 |
|------------|-----------------------|
| `network`  | The ID of the network |

### Node labels

| Label name | Value                                                  |
|------------|--------------------------------------------------------|
| `manager`    | "true" if the node is a UCP manager, "false" otherwise |

## Metric source

UCP exports metrics on every node and also exports additional metrics from
every controller. The metrics that are exported from controllers are
cluster-scoped, for example, the total number of Swarm services. Metrics that
are exported from nodes are specific to those nodes, for example, the total memory
on that node.
