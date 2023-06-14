---
description: Registry Access Management
keywords: registry, access, managment
title: Registry Access Management
---

{% include admin-early-access.md %}

Registry Access Management (RAM) is a feature available to organizations with a Docker Business subscription. When RAM is enabled, organization owners can ensure that their developers using Docker Desktop can only access registries that have been allow-listed via the Registry Access Management dashboard on Docker Hub to reflect support for other registries: AWS ECR, GitHub Container Registry, Google Container Registry, Quay, a local private registry, and others.

For example, you can use RAM if you manage engineering teams that use Docker Desktop for local development and want to ensure that the images they are pulling are licensed and reputable before using them.

## Requirements:

Your Docker users must use Docker Desktop v4.8 or a later.

## Configure Registry Access Management permissions

To configure Registry Access Management permissions, perform the following steps:

1. Sign in to [Docker Admin](https://admin.docker.com){: target="_blank" rel="noopener" class="_"}.
2. In the left navigation, select your organization in the drop-down menu.
3. Select **Registry Access**.
4. Enable Registry Access Management to set the permissions for your registry.

   > **Note**
   >
   > When enabled, the Docker Hub registry is set by default, however you can also restrict this registry for your developers.

5. Select **Add** and enter your registry details in the applicable fields, and select **Create** to add the registry to your list.
6. Verify that the registry appears in your list and select **Save & Apply**. 

   > **Note**
   >
   > Once you add a registry, it can take up to 24 hours for the changes to be enforced on your developers’ machines. If you want to apply the changes sooner, you must force a Docker logout on your developers’ machine and have the developers re-authenticate for Docker Desktop. Also, there is no limit on the number of registries you can add. See the [Caveats](#caveats) section to learn more about limitations when using this feature.

## Enforce authentication

To ensure that each organization member uses Registry Access Management on their local machine, you can perform the following steps to enforce sign-in under your organization. To do this:

1. Install the latest version of Docker Desktop on your member's machine.
2. Create a `registry.json` file by following the instructions for [Configure registry.json](../../docker-hub/configure-sign-in.md).

## Verify the restrictions

The new Registry Access Management policy should be in place after the developer successfully authenticates to Docker Desktop using their organization credentials. The developer can attempt to pull an image from a disallowed registry via the Docker CLI. They will then receive an error message that your organization has disallowed this registry.

### Caveats

There are certain limitations when using Registry Access Management; they are as follows:

- Windows image pulls, and image builds are not restricted
- Builds such as `docker buildx` using a Kubernetes driver are not restricted
- Builds such as `docker buildx` using a custom docker-container driver are not restricted
- Blocking is DNS-based; you must use a registry's access control mechanisms to distinguish between “push” and “pull”
- WSL 2 requires at least a 5.4 series Linux kernel (this does not apply to earlier Linux kernel series)
- Under the WSL 2 network, traffic from all Linux distributions is restricted (this will be resolved in the updated 5.15 series Linux kernel)

Also, Registry Access Management operates on the level of hosts, not IP addresses. Developers can bypass this restriction within their domain resolution, for example by running Docker against a local proxy or modifying their operating system's `sts` file. Blocking these forms of manipulation is outside the remit of Docker Desktop.
