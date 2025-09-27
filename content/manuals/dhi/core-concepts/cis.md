---
title: CIS Benchmark
description: Learn how Docker Hardened Images comply with the CIS Docker Benchmark to help organizations harden container images for secure deployments.
keywords: docker cis benchmark, cis docker compliance, cis docker images, docker hardened images, secure container images
---

## What is the CIS Docker Benchmark?

The [CIS Docker Benchmark](https://www.cisecurity.org/benchmark/docker) is part
of the globally recognized CIS Benchmarks, developed by the [Center for
Internet Security (CIS)](https://www.cisecurity.org/). It defines recommended secure
configurations for all aspects of the Docker container ecosystem, including the
container host, Docker daemon, container images, and the container runtime.

## Why CIS Benchmark compliance matters

Following the CIS Docker Benchmark helps organizations:

- Reduce security risk with widely recognized hardening guidance.
- Meet regulatory or contractual requirements that reference CIS controls.
- Standardize image and Dockerfile practices across teams.
- Demonstrate audit readiness with configuration decisions grounded in a public standard.

## How Docker Hardened Images comply with the CIS Benchmark

Docker Hardened Images (DHIs) are designed with security in mind and are
verified to be compliant with the relevant controls from the latest CIS
Docker Benchmark (v1.8.0) for the scope that applies to container images and
Dockerfile configuration.

CIS-compliant DHIs are compliant with all controls in Section 4, with the sole
exception of the control requiring Docker Content Trust (DCT), which [Docker
officially retired](https://www.docker.com/blog/retiring-docker-content-trust/).
Instead, DHIs are [signed](/manuals/dhi/core-concepts/signatures.md) using
Cosign, providing an even higher level of authenticity and integrity. By
starting from a CIS-compliant DHI, teams can adopt image-level best practices
from the benchmark more quickly and confidently.

> [!NOTE]
>
> The CIS Docker Benchmark also includes controls for the host, daemon, and
> runtime. CIS-compliant DHIs address only the image and Dockerfile scope (Section
> 4). Overall compliance still depends on how you configure and operate the
> broader environment.

## Identify CIS-compliant images

CIS-compliant images are labeled as **CIS** in the Docker Hardened Images catalog.
To find them, [explore images](../how-to/explore.md) and look for the **CIS**
designation on individual listings.

## Get the benchmark

Download the latest CIS Docker Benchmark directly from CIS:
https://www.cisecurity.org/benchmark/docker
