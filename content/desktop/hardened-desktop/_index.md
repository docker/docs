---
title: Overview of Hardened Docker Desktop
description: Overview of what Hardened Docker Desktop is and its key features
keywords: security, hardened desktop, enhanced container isolation, registry access
  management, settings management root access, admins, docker desktop, image access
  management
grid:
  - title: "Settings Management"
    description: Learn how Settings Management can secure your developers' workflows.
    icon: shield_locked
    link: /desktop/hardened-desktop/settings-management/
  - title: "Enhanced Container Isolation"
    description: Understand how Enhanced Container Isolation can prevent container attacks.
    icon: "security"
    link: /desktop/hardened-desktop/enhanced-container-isolation/
  - title: "Registry Access Management"
    description: Control the registries developers can access while using Docker Desktop.
    icon: "home_storage"
    link: /security/for-admins/registry-access-management/
  - title: "Image Access Management"
    description: Control the images developers can pull from Docker Hub.
    icon: "photo_library"
    link: /security/for-admins/image-access-management/
  - title: "Air-Gapped Containers"
    description: Restrict containers from accessing unwanted network resources.
    icon: "shield_locked"
    link: /desktop/hardened-desktop/air-gapped-containers/
---

> **Note**
>
> Hardened Docker Desktop is available to Docker Business customers only.

### What is Hardened Docker Desktop?

Hardened Docker Desktop is a group of security features for Docker Desktop, designed to improve the security of developer environments without impacting developer experience or productivity.

It lets admins define and enforce robust security settings. It guarantees that developers and the containers they deploy are unable to intentionally or unintentionally circumvent these settings. Additionally, you can enhance container isolation, which helps mitigate potential security threats such as malicious payloads breaching the Docker Desktop Linux VM and the underlying host.

Hardened Docker Desktop moves the ownership boundary for Docker Desktop configuration to the organization, meaning that any security controls admins set cannot be altered by the user of Docker Desktop.

### Who is it for?

It is for security conscious organizations who:
- Don’t give their users root or admin access on their machines
- Would like Docker Desktop to be within their organization’s centralized control
- Have certain compliance obligations

### What does Hardened Docker Desktop include?

It includes:

{{< grid >}}

- Settings Management, which helps admins to confidently manage and control the usage of Docker Desktop within their organization.
- Enhanced Container Isolation (ECI), a setting that instantly enhances security by preventing containers from running as root in Docker Desktop’s Linux VM and ensures that any configurations set using Settings Management cannot be bypassed or modified by containers.
- Registry Access Management (RAM), which allows admins to control the registries developers can access.
- Image Access Management (IAM), which gives admins control over which images developers can pull from Docker Hub.
- Air-Gapped Containers, which restrict containers from accessing unwanted network resources.

### How does it help my organisation?

Hardened Desktop features work independently but collectively to create a defense-in-depth strategy, safeguarding developer workstations against potential attacks across various functional layers, such as configuring Docker Desktop, pulling container images, and running container images. This multi-layered defense approach ensures comprehensive security.

It helps mitigate against threats such as:
 - **Malware and supply chain attacks:** RAM and IAM prevent developers from accessing certain container registries and image types, significantly lowering the risk of malicious payloads. Additionally, ECI restricts the impact of containers with malicious payloads by running them without root privileges inside a Linux user namespace.
 - **Lateral movement:** Air gapped containers allows admins to configure network access restrictions for containers, thereby preventing malicious containers from lateral movement within the organization's network.
 - **Insider threats:** Settings Management configures and locks various Docker Desktop settings, such as proxy settings, ECI, and prevents exposure of the Docker API. This helps admins enforce company policies and prevents developers from introducing insecure configurations, intentionally or unintentionally.
