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
    link: "/desktop/hardened-desktop/settings-management/"
  - title: "Enhanced Container Isolation"
    description: Understand how Enhanced Container Isolation can prevent container attacks.
    icon: "security"
    link: "/desktop/hardened-desktop/enhanced-container-isolation/"
  - title: "Registry Access Management"
    description: Control the registries developers can access while using Docker Desktop.
    icon: "home_storage"
    link: "/security/for-admins/registry-access-management/"
  - title: "Image Access Management"
    description: Control the images developers can pull from Docker Hub.
    icon: "photo_library"
    link: "/security/for-admins/image-access-management/"
---

>Note
>
>Hardened Docker Desktop is available to Docker Business customers only.

Hardened Docker Desktop is a group of security features for Docker Desktop, designed to improve the security of developer environments without impacting developer experience or productivity.

It is for security conscious organizations who don’t give their users root or admin access on their machines, and who would like Docker Desktop to be within their organization’s centralized control.

Hardened Docker Desktop moves the ownership boundary for Docker Desktop configuration to the organization, meaning that any security controls admins set cannot be altered by the user of Docker Desktop.

Hardened Docker Desktop includes:
- Settings Management, which helps admins to confidently manage and control the usage of Docker Desktop within their organization.
- Enhanced Container Isolation, a setting that instantly enhances security by preventing containers from running as root in Docker Desktop’s Linux VM and ensures that any configurations set using Settings Management cannot be bypassed or modified by containers.
- Registry Access Management, which allows admins to control the registries developers can access.
- Image Access Management, which gives admins control over which images developers can pull from Docker Hub.

The features of Hardened Docker Desktop operate independently of each other. When used together, these mechanisms defend against attacks at different functional layers of the developer workflow, providing a defense-in-depth approach to securing developer environments.

Docker plans to continue adding more security enhancements to the Hardened Docker Desktop security model.

{{< grid >}}
