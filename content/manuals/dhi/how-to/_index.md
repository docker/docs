---
title: How-tos
description: Step-by-step guidance for working with Docker Hardened Images, from discovery to debugging.
weight: 20
params:
  grid_discover:
    - title: Explore Docker Hardened Images
      description: Learn how to find and evaluate image repositories, variants, metadata, and attestations in the DHI catalog on Docker Hub.
      icon: travel_explore
      link: /dhi/how-to/explore/
  grid_adopt:
    - title: Mirror a Docker Hardened Image repository
      description: Learn how to mirror an image into your organization's namespace and optionally push it to another private registry.
      icon: compare_arrows
      link: /dhi/how-to/mirror/
    - title: Customize a Docker Hardened Image or chart
      description: Learn how to customize Docker Hardened Images and charts.
      icon: settings
      link: /dhi/how-to/customize/
    - title: Use a Docker Hardened Image
      description: Learn how to pull, run, and reference Docker Hardened Images in Dockerfiles, CI pipelines, and standard development workflows.
      icon: play_arrow
      link: /dhi/how-to/use/
    - title: Use a Docker Hardened Image in Kubernetes
      description: Learn how to use Docker Hardened Images in Kubernetes deployments.
      icon: play_arrow
      link: /dhi/how-to/k8s/
    - title: Use a Docker Hardened Image chart
      description: Learn how to use a Docker Hardened Image chart.
      icon: leaderboard
      link: /dhi/how-to/helm/
    - title: Use Extended Lifecycle Support for Docker Hardened Images
      description: Learn how to use Extended Lifecycle Support with Docker Hardened Images.
      icon: update
      link: /dhi/how-to/els/
    - title: Manage Docker Hardened Images and charts
      description: Learn how to manage your mirrored and customized Docker Hardened Images in your organization.
      icon: reorder
      link: /dhi/how-to/manage/
  grid_evaluate:
    - title: Compare Docker Hardened Images
      description: Learn how to compare Docker Hardened Images with other container images to evaluate security improvements and differences.
      icon: compare
      link: /dhi/how-to/compare/
  grid_verify:
    - title: Verify a Docker Hardened Image or chart
      description: Use Docker Scout or cosign to verify signed attestations like SBOMs, provenance, and vulnerability data for Docker Hardened Images and charts.
      icon: check_circle
      link: /dhi/how-to/verify/
    - title: Scan Docker Hardened Images
      description: Learn how to scan Docker Hardened Images for known vulnerabilities using Docker Scout, Grype, or Trivy.
      icon: bug_report
      link: /dhi/how-to/scan/
  grid_govern:
    - title: Enforce Docker Hardened Image usage with policies
      description: Learn how to use image policies with Docker Scout for Docker Hardened Images.
      icon: policy
      link: /dhi/how-to/policies/
  grid_troubleshoot:
    - title: Debug a Docker Hardened Image
      description: Use Docker Debug to inspect a running container based on a hardened image without modifying it.
      icon: terminal
      link: /dhi/how-to/debug/
---

This section provides practical, task-based guidance for working with Docker
Hardened Images (DHIs). Whether you're evaluating DHIs for the first time or
integrating them into a production CI/CD pipeline, these topics cover the key
tasks across the adoption journey, from discovery to debugging.

The topics are organized around the typical lifecycle of working with DHIs, but
you can use them as needed based on your specific workflow.

Explore the topics below that match your current needs.

## Discover

Explore available images and metadata in the DHI catalog.

{{< grid
  items="grid_discover"
>}}

## Adopt

Mirror trusted images, customize as needed, and integrate into your workflows.

{{< grid
  items="grid_adopt"
>}}

## Evaluate

Compare with other images to understand security improvements.

{{< grid
  items="grid_evaluate"
>}}

## Verify

Check signatures, SBOMs, and provenance, and scan for vulnerabilities.

{{< grid
  items="grid_verify"
>}}

## Govern

Enforce policies to maintain security and compliance.

{{< grid
  items="grid_govern"
>}}

## Troubleshoot

Debug containers based on DHIs without modifying the image.

{{< grid
  items="grid_troubleshoot"
>}}