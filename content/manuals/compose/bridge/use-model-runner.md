---
title: Use Docker Model Runner with Compose Bridge 
linkTitle: Use Model Runner
weight: 30
description: How to use Docker Model Runner with Compose Bridge for consistent deployments
keywords: docker compose bridge, customize compose bridge, compose bridge templates, compose to kubernetes, compose bridge transformation, go templates docker, model runner, ai, llms
---

Compose Bridge now supports model-aware deployments. It can deploy and configure Docker Model Runner, a lightweight service that hosts and serves machine LLMs.

This reduces manual setup for LLM-enabled services and keeps deployments consistent across Docker Desktop and Kubernetes environments.

If you have a `models` top-level element in your `compose.yaml` file, Compose Bridge:

- Automatically injects environment variables for each model’s endpoint and name.
- Configures model endpoints differently for Docker Desktop vs Kubernetes.
- Optionally deploys Docker Model Runner in Kubernetes when enabled in Helm values

## Configure model runner settings

You can control the model runner through Helm values.

```yaml
modelRunner:
  enabled: true               # true = deploy in-cluster model runner
                              # false = use Docker Desktop host model runner
  endpoint: "http://model-runner/engines/v1/"
  hostEndpoint: "http://host.docker.internal:12434/engines/v1/"
  image: "docker/model-runner:latest"
  resources:
    limits:
      cpu: "1000m"
      memory: "2Gi"
    requests:
      cpu: "100m"
      memory: "256Mi"
  storage:
    size: "100Gi"
    storageClass: ""          # Use default storage class if empty
  models:
    - sentiment
    - toxicity
```

| Setting        | Description                                                                               |
| -------------- | ----------------------------------------------------------------------------------------- |
| `enabled`      | Deploy Docker Model Runner inside your cluster (`true`) or use an external runner (`false`). |
| `endpoint`     | URL for Docker Model Runner used by injected environment variables.                          |
| `hostEndpoint` | Address of Docker Model Runner for Docker Desktop.                                            |
| `models`       | List of models to pre-pull during startup.                                                |
| `storage`      | Persistent storage configuration for model files.                                         |
| `resources`    | Resource requests and limits for the model runner pod.                                    |


## Deploying a model runner

### Docker Desktop

When `modelRunner.enabled` is `false`, Compose Bridge configures your workloads to connect to Docker Model Runner running on the host:

```text
http://host.docker.internal:12434/engines/v1/
```

The endpoint is automatically injected into your service containers.

### Kubernetes

When `modelRunner.enabled` is `true`, Compose Bridge generates additional manifests that deploy Docker Model Runner in your cluster, including:

- Deployment — runs the `docker-model-runner` container
- Service — exposes port `80` (maps to container port `12434`)
- `PersistentVolumeClaim` — stores model files

