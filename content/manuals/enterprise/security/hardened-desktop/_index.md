---
title: Hardened Docker Desktop
linkTitle: Hardened Docker Desktop
description: Security features that help organizations secure developer environments without impacting productivity
keywords: security, hardened desktop, enhanced container isolation, registry access management, settings management, admins, docker desktop, image access management, air-gapped containers
tags: [admin]
aliases:
 - /desktop/hardened-desktop/
 - /security/for-admins/hardened-desktop/
grid:
  - title: "Settings Management"
    description: Learn how Settings Management can secure your developers' workflows.
    icon: shield_locked
    link: /enterprise/security/hardened-desktop/settings-management/
  - title: "Enhanced Container Isolation"
    description: Understand how Enhanced Container Isolation can prevent container attacks.
    icon: "security"
    link: /enterprise/security/hardened-desktop/enhanced-container-isolation/
  - title: "Registry Access Management"
    description: Control the registries developers can access while using Docker Desktop.
    icon: "home_storage"
    link: /enterprise/security/hardened-desktop/registry-access-management/
  - title: "Image Access Management"
    description: Control the images developers can pull from Docker Hub.
    icon: "photo_library"
    link: /enterprise/security/hardened-desktop/image-access-management/
  - title: "Air-Gapped Containers"
    description: Restrict containers from accessing unwanted network resources.
    icon: "vpn_lock"
    link: /enterprise/security/hardened-desktop/air-gapped-containers/
weight: 60
---

{{< summary-bar feature_name="Hardened Docker Desktop" >}}

Hardened Docker Desktop provides a collection of security features designed to strengthen developer environments without compromising productivity or developer experience.

With Hardened Docker Desktop, you can enforce strict security policies that prevent developers and containers from bypassing organizational controls. You can also enhance container isolation to protect against security threats like malicious payloads that might breach the Docker Desktop Linux VM or underlying host system.

## Who should use Hardened Docker Desktop?

Hardened Docker Desktop is ideal for security-focused organizations that:

- Don't provide root or administrator access to developers' machines
- Want centralized control over Docker Desktop configurations
- Must meet specific compliance requirements

## How Hardened Docker Desktop works

Hardened Docker Desktop features work independently and together to create a defense-in-depth security strategy. They protect developer workstations against attacks across multiple layers, including Docker Desktop configuration, container image management, and container runtime security:

- Registry Access Management and Image Access Management prevent access to unauthorized container registries and image types, reducing exposure to malicious payloads
- Enhanced Container Isolation runs containers without root privileges inside a Linux user namespace, limiting the impact of malicious containers
- Air-gapped containers let you configure network restrictions for containers, preventing malicious containers from accessing your organization's internal network resources
- Settings Management locks down Docker Desktop configurations to enforce company policies and prevent developers from introducing insecure settings, whether intentionally or accidentally

## Next steps

Explore Hardened Docker Desktop features to understand how they can strengthen your organization's security posture:

{{< grid >}}
