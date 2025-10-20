---
title: Use Docker Model Runner with Compose Bridge 
linkTitle: Use Model Runner
weight: 30
description: How to use Docker Model Runner with Compose Bridge for consistent deployments
keywords: docker compose bridge, customize compose bridge, compose bridge templates, compose to kubernetes, compose bridge transformation, go templates docker, model runner, ai, llms
---

Compose Bridge supports model-aware deployments. It can deploy and configure Docker Model Runner, a lightweight service that hosts and serves machine LLMs.

This reduces manual setup for LLM-enabled services and keeps deployments consistent across Docker Desktop and Kubernetes environments.

If you have a `models` top-level element in your `compose.yaml` file, Compose Bridge:

- Automatically injects environment variables for each model’s endpoint and name.
- Configures model endpoints differently for Docker Desktop vs Kubernetes.
- Optionally deploys Docker Model Runner in Kubernetes when enabled in Helm values

## Configure model runner settings

When deploying using generated Helm Charts, you can control the model runner configuration through Helm values.

```yaml
# Model Runner settings
modelRunner:
    # Set to false for Docker Desktop (uses host instance)
    # Set to true for standalone Kubernetes clusters
    enabled: false
    # Endpoint used when enabled=false (Docker Desktop)
    hostEndpoint: "http://host.docker.internal:12434/engines/v1/"
    # Deployment settings when enabled=true
    image: "docker/model-runner:latest"
    imagePullPolicy: "IfNotPresent"
    # GPU support
    gpu:
        enabled: false
        vendor: "nvidia" # nvidia or amd
        count: 1
    # Node scheduling (uncomment and customize as needed)
    # nodeSelector:
    #   accelerator: nvidia-tesla-t4
    # tolerations: []
    # affinity: {}

    # Security context
    securityContext:
        allowPrivilegeEscalation: false
    # Environment variables (uncomment and add as needed)
    # env:
    #   DMR_ORIGINS: "http://localhost:31246"
    resources:
        limits:
            cpu: "1000m"
            memory: "2Gi"
        requests:
            cpu: "100m"
            memory: "256Mi"
    # Storage for models
    storage:
        size: "100Gi"
        storageClass: "" # Empty uses default storage class
    # Models to pre-pull
    models:
        - ai/qwen2.5:latest
        - ai/mxbai-embed-large
```

## Deploying a model runner

### Docker Desktop

When `modelRunner.enabled` is `false`, Compose Bridge configures your workloads to connect to Docker Model Runner running on the host:

```text
http://host.docker.internal:12434/engines/v1/
```

The endpoint is automatically injected into your service containers.

### Kubernetes

When `modelRunner.enabled` is `true`, Compose Bridge generates additional manifests that deploy Docker Model Runner in your cluster, including:

- Deployment: Runs the `docker-model-runner` container
- Service: Rxposes port `80` (maps to container port `12434`)
- `PersistentVolumeClaim`: Stores model files

The `modelRunner.enabled` setting also determines the number of replicas for the `model-runner-deploymen`t:

- When `true`, the deployment replica count is set to 1, and Docker Model Runner is deployed in the Kubernetes cluster.
- When `false`, the replica count is 0, and no Docker Model Runner resources are deployed.