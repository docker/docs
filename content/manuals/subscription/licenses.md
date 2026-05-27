---
title: Licenses
linkTitle: Licenses
description: Manage product licenses for your organization, including assignment, revocation, and automatic assignment.
keywords: licenses, organization, members, Docker Offload, AI governance, license assignment, admin console
weight: 50
grid:
  - title: Manage licenses
    description: Assign, revoke, and bulk-assign product licenses from the Members page.
    icon: user-plus
    link: /admin/organization/manage/manage-licenses/
  - title: Automatic license assignment
    description: Configure automatic assignment and understand what happens when licenses run out.
    icon: arrow-path
    link: /admin/organization/manage/manage-licenses/#automated-license-assignment
  - title: Docker Offload
    description: Overview of Offload licensing and access requirements.
    icon: cloud
    link: /offload/about/
  - title: AI governance for sandboxes
    description: Organization policies for sandbox network, filesystem, and governance behavior.
    icon: shield-check
    link: /ai/sandboxes/security/governance/
  - title: Docker Sandboxes (sbx)
    description: Install and run sandboxes; understand when organization policies apply.
    icon: cube
    link: /ai/sandboxes/get-started/
  - title: Scale your subscription
    description: Add seats or capacity for Docker subscription products.
    icon: chart-bar
    link: /subscription/scale/
---

You can purchase licenses that give  
members of your organization access Docker products in addition to your Docker subscription. Licenses add an additional layer to your subscription, letting admins assign certain products to select members without buying a Docker product for the entire  
organization. With licenses, members can belong to your organization  
without holding a license for every product your organization owns.

## Key features

Organization membership and product licenses are separate. When someone  
joins your organization, they won’t automatically receive licenses for  
license-based products.

- You can purchase licenses for Docker Offload or AI Governance.
- Once purchased, admins manage licenses from the Members page in Docker Home.

To add licenses to your organization, contact your Docker account team.

## AI governance licenses

Docker Sandboxes (`sbx`) run AI coding agents in isolated virtual environments so agents can build containers, install packages, and use tools without accessing your host system.

- AI governance licenses allow admins to apply network and filesystem policies across your organization for users who hold the license.
- When a member who holds an AI governance license runs `sbx`, the organization policy overrides the member’s local policy rules.

  <Explanation about how licenses work when you aren’t part of a Docker Core subscription\>  
  To learn more about AI governance, see [Organization governance](/manuals/ai/sandboxes/security/governance.md) and [Sandbox policies](/manuals/ai/sandboxes/security/policy.md).

  ## Offload licenses

  Docker Offload lets developers offload building and running containers to the cloud. You can assign Offload licenses to the members who need cloud-backed builds and runs without giving every organization member access to Offload.

  To purchase Docker Offload licenses, you must have a Docker Team or Docker Business subscription.

  ## What’s next

{{< grid >}}
