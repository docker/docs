---
description: Containerize a Golang application and monitor it with Prometheus and Grafana.
keywords: golang, prometheus, grafana, monitoring, containerize
title: Monitor a Golang application with Prometheus and Grafana
summary: |
  Learn how to containerize a Golang application and monitor it with Prometheus and Grafana.
linkTitle: Monitor with Prometheus and Grafana
languages: [go]
params:
  time: 45 minutes
---

The guide teaches you how to containerize a Golang application and monitor it with Prometheus and Grafana. 

> **Acknowledgment**
>
> Docker would like to thank [Pradumna Saraf](https://twitter.com/pradumna_saraf) for his contribution to this guide.

## Overview

To make sure our application is working as intended, monitoring is really important. One of the most popular monitoring tools is Prometheus. Prometheus is an open-source monitoring and alerting toolkit that is designed for reliability and scalability. It collects metrics from monitored targets by scraping metrics HTTP endpoints on these targets. To visualize the metrics, we can use Grafana. Grafana is an open-source platform for monitoring and observability that allows you to query, visualize, alert on, and understand your metrics no matter where they are stored.

In this guide, we will be creating a Golang server with some endpoints to simulate a real-world application. Then we will expose metrics from the server using Prometheus. Finally, we will visualize the metrics using Grafana. We will containerize the Golang application and using the Docker Compose file, we will connect all the services- Golang, Prometheus, and Grafana. 

## What will you learn?

* Create a Golang application with custom Prometheus metrics.
* Containerize a Golang application.
* Use Docker Compose to run multiple services and connect them together to monitor a Golang application with Prometheus and Grafana.
* Visualize the metrics using Grafana dashboards.

## Prerequisites

- A good understanding of Golang is assumed.
- You must me familiar with Prometheus and creating dashboards in Grafana.
- You must have familiarity with Docker concepts like containers, images, and Dockerfiles. If you are new to Docker, you can start with the [Docker basics](/get-started/docker-concepts/the-basics/what-is-a-container.md) guide.

## Next steps

We will be creating a Golang server and exposing metrics using Prometheus.
