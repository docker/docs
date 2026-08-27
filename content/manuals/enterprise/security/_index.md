---
linkTitle: Security
title: Security for enterprises
description: Learn about enterprise level security features Docker has to offer and explore best practices
keywords: docker, docker hub, docker desktop, security, enterprises, scale
weight: 10
params:
  sidebar:
    group: Enterprise
grid_admins:
  - title: Settings Management
    description: Learn how Settings Management can secure your developers' workflows.
    icon: shield-check
    link: /enterprise/security/hardened-desktop/settings-management/
  - title: Enhanced Container Isolation
    description: Understand how Enhanced Container Isolation can prevent container attacks.
    icon: shield-check
    link: /enterprise/security/hardened-desktop/enhanced-container-isolation/
  - title: Registry Access Management
    description: Control the registries developers can access while using Docker Desktop.
    icon: server
    link: /enterprise/security/hardened-desktop/registry-access-management/
  - title: Image Access Management
    description: Control the images developers can pull from Docker Hub.
    icon: photo
    link: /enterprise/security/hardened-desktop/image-access-management/
  - title: "Air-Gapped Containers"
    description: Restrict containers from accessing unwanted network resources.
    icon: lock-closed
    link: /enterprise/security/hardened-desktop/air-gapped-containers/
  - title: Namespace access
    description: Control which Kubernetes namespaces developers can access in Docker Desktop.
    icon: lock-closed
    link: /enterprise/security/hardened-desktop/namespace-access/
---

Docker provides security guardrails for both administrators and developers.

If you're an administrator, you can enforce sign-in across Docker products for your developers, and
scale, manage, and secure your instances of Docker Desktop with DevOps security controls like Enhanced Container Isolation and Registry Access Management.

## For administrators

Explore the security features Docker offers to satisfy your company's security policies.

{{< grid items="grid_admins" >}}
