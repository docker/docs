---
title: Overview of Hardened Docker Desktop
linkTitle: Hardened Docker Desktop
description: Overview of what Hardened Docker Desktop is and its key features
keywords: security, hardened desktop, enhanced container isolation, registry access
  management, settings management root access, admins, docker desktop, image access
  management
tags: [admin]
aliases:
 - /desktop/hardened-desktop/
grid:
  - title: "Settings Management"
    description: Learn how Settings Management can secure your developers' workflows.
    icon: shield_locked
    link: /security/for-admins/hardened-desktop/settings-management/
  - title: "Enhanced Container Isolation"
    description: Understand how Enhanced Container Isolation can prevent container attacks.
    icon: "security"
    link: /security/for-admins/hardened-desktop/enhanced-container-isolation/
  - title: "Registry Access Management"
    description: Control the registries developers can access while using Docker Desktop.
    icon: "home_storage"
    link: /security/for-admins/hardened-desktop/registry-access-management/
  - title: "Image Access Management"
    description: Control the images developers can pull from Docker Hub.
    icon: "photo_library"
    link: /security/for-admins/hardened-desktop/image-access-management/
  - title: "Air-Gapped Containers"
    description: Restrict containers from accessing unwanted network resources.
    icon: "vpn_lock"
    link: /security/for-admins/hardened-desktop/air-gapped-containers/
weight: 60
---

{{< summary-bar feature_name="Hardened Docker Desktop" >}}

Hardened Docker Desktop is a group of security features, designed to improve the security of developer environments with minimal impact on developer experience or productivity.

It lets you enforce strict security settings, preventing developers and their containers from bypassing these controls, either intentionally or unintentionally. Additionally, you can enhance container isolation, to mitigate potential security threats such as malicious payloads breaching the Docker Desktop Linux VM and the underlying host.

Hardened Docker Desktop moves the ownership boundary for Docker Desktop configuration to the organization, meaning that any security controls you set cannot be altered by the user of Docker Desktop. 

It is for security conscious organizations who:
- Don’t give their users root or administrator access on their machines
- Would like Docker Desktop to be within their organization’s centralized control
- Have certain compliance obligations

### How does it help my organization?

Hardened Desktop features work independently but collectively to create a defense-in-depth strategy, safeguarding developer workstations against potential attacks across various functional layers, such as configuring Docker Desktop, pulling container images, and running container images. This multi-layered defense approach ensures comprehensive security. It helps mitigate against threats such as:

 - Malware and supply chain attacks: Registry Access Management and Image Access Management prevent developers from accessing certain container registries and image types, significantly lowering the risk of malicious payloads. Additionally, Enhanced Container Isolation (ECI) restricts the impact of containers with malicious payloads by running them without root privileges inside a Linux user namespace.
 - Lateral movement: Air-gapped containers lets you configure network access restrictions for containers, thereby preventing malicious containers from performing lateral movement within the organization's network.
 - Insider threats: Settings Management configures and locks various Docker Desktop settings so you can enforce company policies and prevent developers from introducing insecure configurations, intentionally or unintentionally.

{{< grid >}}
