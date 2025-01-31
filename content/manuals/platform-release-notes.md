---
title: Release notes for Docker Home, the Admin Console, billing, security, and subscription features
linkTitle: Release notes
description: Learn about the new features, bug fixes, and breaking changes for Docker Home, the Admin Console, and billing and subscription features
keywords: Docker Home, Docker Admin Console, billing, subscription, security, admin, releases, what's new
weight: 60
params:
  sidebar:
    group: Platform
tags: [Release notes, admin]
---

This page provides details on new features, enhancements, known issues, and bug fixes across Docker Home, the Admin Console, billing, security, and subscription functionalities.

Take a look at the [Docker Public Roadmap](https://github.com/orgs/docker/projects/51/views/1?filterQuery=) to see what's coming next.

## 2025-01-30

### New

- Installing Docker Desktop via the PKG installer is now generally available.
- Enforcing sign-in via configuration profiles is now generally available.

## 2024-12-10

### New

- New Docker subscriptions are now available. For more information, see [Docker
  subscriptions and features](/manuals/subscription/details.md) and [Announcing
  Upgraded Docker Plans: Simpler, More Value, Better Development and
  Productivity](https://www.docker.com/blog/november-2024-updated-plans-announcement/).

## 2024-11-18

### New

- Administrators can now:
  - Enforce sign-in with [configuration profiles](/manuals/security/for-admins/enforce-sign-in/methods.md#configuration-profiles-method-mac-only) (Early Access).
  - Enforce sign-in for more than one organization at a time (Early Access).
  - Deploy Docker Desktop for Mac in bulk with the [PKG installer](/manuals/desktop/setup/install/enterprise-deployment/pkg-install-and-configure.md) (Early Access).
  - [Use Desktop Settings Management via the Docker Admin Console](/manuals/security/for-admins/hardened-desktop/settings-management/configure-admin-console.md) (Early Access).

### Bug fixes and enhancements

- Enhance Container Isolation (ECI) has been improved to:
  - Permit admins to [turn off Docker socket mount restrictions](/manuals/security/for-admins/hardened-desktop/enhanced-container-isolation/config.md#allowing-all-containers-to-mount-the-docker-socket).
  - Support wildcard tags when using the [`allowedDerivedImages` setting](/manuals/security/for-admins/hardened-desktop/enhanced-container-isolation/config.md#docker-socket-mount-permissions-for-derived-images).

## 2024-11-11

### New

- [Personal access tokens](/security/for-developers/access-tokens/) (PATs) now support expiration dates.

## 2024-10-15

### New

- Beta: You can now create [organization access tokens](/security/for-admins/access-tokens/) (OATs) to enhance security for organizations and streamline access management for organizations in the Docker Admin Console.

## 2024-08-29

### New

- Deploying Docker Desktop via the [MSI installer](/manuals/desktop/setup/install/enterprise-deployment/msi-install-and-configure.md) is now generally available.
- Two new methods to [enforce sign-in](/manuals/security/for-admins/enforce-sign-in/_index.md) (Windows registry key and `.plist` file) are now generally available.

## 2024-08-24

### New

- Administrators can now view [organization insights](/manuals/admin/organization/insights.md) (Early Access).

## 2024-07-17

### New

- You can now centrally access and manage Docker products in [Docker Home](https://app.docker.com) (Early Access).