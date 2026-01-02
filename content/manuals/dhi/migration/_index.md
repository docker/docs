---
title: Migration
description: Learn how to migrate your existing applications to Docker Hardened Images
weight: 18
keywords: migrate, docker hardened images, dhi, migration guide
aliases:
  - /dhi/how-to/migrate/
params:
  grid_migration_paths:
    - title: Migrate with Docker's AI assistant
      description: Use Docker's AI assistant to automatically migrate your Dockerfile to Docker Hardened Images with guidance and recommendations.
      icon: smart_toy
      link: /dhi/migration/migrate-with-ai/
    - title: Migrate from Alpine or Debian images
      description: Manual migration guide for moving from Docker Official Images (Alpine or Debian-based) to Docker Hardened Images.
      icon: code
      link: /dhi/migration/migrate-from-doi/
    - title: Migrate from Wolfi
      description: Manual migration guide for transitioning from Wolfi-based images to Docker Hardened Images.
      icon: transform
      link: /dhi/migration/migrate-from-wolfi/
  
  grid_migration_resources:
    - title: Migration checklist
      description: A comprehensive checklist of migration considerations to ensure successful transition to Docker Hardened Images.
      icon: checklist
      link: /dhi/migration/checklist/
    - title: Examples
      description: Example Dockerfile migrations for different programming languages and frameworks to guide your migration process.
      icon: preview
      link: /dhi/migration/examples/
---

This section provides guidance for migrating your applications to Docker
Hardened Images (DHI). Migrating to DHI enhances the security posture of your
containerized applications by leveraging hardened base images with built-in
security features.

## Migration paths

Choose the migration approach that best fits your needs:

{{< grid items="grid_migration_paths" >}}

## Resources

{{< grid items="grid_migration_resources" >}}


