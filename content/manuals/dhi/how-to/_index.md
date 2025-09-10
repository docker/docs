---
title: How-tos
description: Step-by-step guidance for working with Docker Hardened Images, from discovery to debugging.
weight: 20
params:
  grid_howto:
    - title: Explore Docker Hardened Images
      description: Learn how to find and evaluate image repositories, variants, metadata, and attestations in the DHI catalog on Docker Hub.
      icon: travel_explore
      link: /dhi/how-to/explore/
    - title: Mirror a Docker Hardened Image repository
      description: Learn how to mirror an image into your organization's namespace and optionally push it to another private registry.
      icon: compare_arrows
      link: /dhi/how-to/mirror/
    - title: Customize a Docker Hardened Image
      description: Learn how to customize a DHI to suit your organization's needs.
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
    - title: Manage Docker Hardened Images
      description: Learn how to manage your mirrored and customized Docker Hardened Images in your organization.
      icon: reorder
      link: /dhi/how-to/manage/
    - title: Migrate an existing application to use Docker Hardened Images
      description: Follow a step-by-step guide to update your Dockerfiles and adopt Docker Hardened Images for secure, minimal, and production-ready builds.
      icon: directions_run
      link: /dhi/how-to/migrate/
    - title: Verify a Docker Hardened Image
      description: Use Docker Scout or cosign to verify signed attestations like SBOMs, provenance, and vulnerability data for Docker Hardened Images.
      icon: check_circle
      link: /dhi/how-to/verify/
    - title: Scan a Docker Hardened Image
      description: Learn how to scan Docker Hardened Images for known vulnerabilities using Docker Scout, Grype, or Trivy.
      icon: bug_report
      link: /dhi/how-to/scan/
    - title: Enforce Docker Hardened Image usage with policies
      description: Learn how to use image policies with Docker Scout for Docker Hardened Images.
      icon: policy
      link: /dhi/how-to/policies/
    - title: Debug a Docker Hardened Image
      description: Use Docker Debug to inspect a running container based on a hardened image without modifying it.
      icon: terminal
      link: /dhi/how-to/debug/
---

This section provides practical, step-by-step guidance for working with Docker
Hardened Images (DHIs). Whether you're evaluating DHIs for the first time or
integrating them into a production CI/CD pipeline, these topics walk you
through each phase of the adoption journey, from discovery to debugging.

To help you get started and stay secure, the topics are organized around the
typical lifecycle of working with DHIs.

## Lifecycle flow

1. Explore available images and metadata in the DHI catalog.
2. Mirror trusted images into your namespace or registry.
3. Adopt DHIs in your workflows by pulling, using in development and CI, and
   migrating existing applications to use secure, minimal base images.
4. Analyze images by verifying signatures, SBOMs, and provenance, and scanning
   for vulnerabilities.
5. Enforce policies to maintain security and compliance.
6. Debug containers based on DHIs without modifying the image.

Each of the following topics aligns with a step in this lifecycle, so you can progress
confidently through exploration, implementation, and ongoing maintenance.

## Step-by-step topics

{{< grid
  items="grid_howto"
>}}