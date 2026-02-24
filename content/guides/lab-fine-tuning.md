---
title: "Lab: Fine-Tuning Local Models"
linkTitle: "Lab: Fine-Tuning Models"
description: |
  Fine-tune AI models using Docker Offload, Docker Model Runner, and Unsloth
  in this hands-on interactive lab.
summary: |
  Hands-on lab: Fine-tune, validate, and share custom AI models using Docker
  Offload, Unsloth, and Docker Model Runner.
keywords: AI, Docker, fine-tuning, Docker Offload, Unsloth, Model Runner, lab, labspace
aliases:
  - /labs/docker-for-ai/fine-tuning/
params:
  tags: [ai, labs]
  time: 60 minutes
  resource_links:
    - title: Docker Model Runner docs
      url: /ai/model-runner/
    - title: Labspace repository
      url: https://github.com/dockersamples/labspace-fine-tuning
---

This lab provides a hands-on walkthrough of fine-tuning AI models using Docker
Offload, Docker Model Runner, and Unsloth. Learn how to customize models for
your specific use case, validate the results, and share them via Docker Hub.

## What you'll learn

- Use Docker Offload to fine-tune a model with GPU acceleration
- Package and share the fine-tuned model on Docker Hub
- Run the custom model with Docker Model Runner
- Understand the end-to-end workflow from training to deployment

## Modules

| # | Module | Description |
|---|--------|-------------|
| 1 | Introduction | Overview of fine-tuning concepts and the Docker AI stack |
| 2 | Fine-Tuning with Docker Offload | Run fine-tuning using Unsloth and Docker Offload |
| 3 | Validate and Publish | Test the fine-tuned model and publish to Docker Hub |
| 4 | Conclusion | Summary, key takeaways, and next steps |

## Prerequisites

- Docker Desktop with Docker Offload enabled
- GPU access (via Docker Offload cloud resources)

## Launch the lab

Ensure you have Docker Offload running, then start the labspace:

```console
$ docker compose -f oci://dockersamples/labspace-fine-tuning up -d
```

Then open your browser to [http://localhost:3030](http://localhost:3030).
