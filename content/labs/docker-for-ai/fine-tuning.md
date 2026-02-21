---
title: Fine-Tuning Local Models
description: |
  Fine-tune AI models using Docker Offload, Docker Model Runner, and Unsloth.
weight: 30
---

This lab provides a hands-on walkthrough of fine-tuning AI models using Docker
Offload, Docker Model Runner, and Unsloth. Learn how to customize models for
your specific use case, validate the results, and share them via Docker Hub.

## What you'll learn

- Use Docker Offload to fine-tune a model with GPU acceleration
- Package and share the fine-tuned model on Docker Hub
- Run the custom model with Docker Model Runner
- Understand the end-to-end workflow from training to deployment

## Prerequisites

- Docker Desktop with Docker Offload enabled
- GPU access (via Docker Offload cloud resources)

## Modules

| # | Module | Description |
|---|--------|-------------|
| 1 | Introduction | Overview of fine-tuning concepts and the Docker AI stack |
| 2 | Fine-Tuning with Docker Offload | Run fine-tuning using Unsloth and Docker Offload |
| 3 | Validate and Publish | Test the fine-tuned model and publish to Docker Hub |
| 4 | Conclusion | Summary, key takeaways, and next steps |

## Launch the lab

1. Clone the repository:

```console
$ git clone https://github.com/dockersamples/labspace-fine-tuning
$ cd labspace-fine-tuning
```

2. Ensure you have Docker Offload running, then start the labspace:

```console
$ docker compose -f oci://dockersamples/labspace-fine-tuning up -d
```

Then open your browser to [http://localhost:3030](http://localhost:3030).

## Resources

- [Labspace repository](https://github.com/dockersamples/labspace-fine-tuning)
- [Docker Model Runner docs](/ai/model-runner/)

<div id="docker-ai-labs-survey-anchor"></div>
