---
title: Accounts and admin release notes
linkTitle: Release notes
description: >-
  Learn about new features, bug fixes, and breaking changes for Docker accounts
  and admin features, including Docker Home, billing, security, and
  subscriptions.
keywords: accounts, admin, Docker Home, billing, subscription, security,
  release notes, what's new
weight: 60
params:
  sidebar:
    group: Accounts and admin
tags: [Release notes, admin]
---

This page lists new features, enhancements, known issues, and bug fixes for
Docker accounts and admin features, including Docker Home, billing, security,
and subscriptions.

## 2026-09-29

### New

- Organization owners can now
  [assign a license to a team](/manuals/accounts/organization/manage/manage-licenses.md#teams).
  Every member of the team receives the license, including people who join the
  team later. Each member uses one license per product.
- The
  [Licenses page](/manuals/accounts/organization/manage/manage-licenses.md#view-licenses)
  in Docker Home shows how many licenses are available and whether they are
  assigned to teams or to individual members.

## 2026-09-24

### New

- You can now subscribe to the
  [Docker Agentic Platform](/manuals/subscription-billing/plans/docker-agentic-platform.md)
  pay-as-you-go plan with a personal account to run agents in cloud
  sandboxes. Compute is metered by the second while a sandbox runs.

## 2026-09-14

### New

- Organization owners can now
  [export a member list](/manuals/accounts/organization/manage/members.md#export-a-member-list-csv)
  from Docker Home and receive the CSV by email. Docker generates the file
  asynchronously and emails a download link to the owner.

## 2026-08-20

### New

- [Docker Verified Publisher](/manuals/subscription-billing/plans/docker-verified-publisher.md)
  Starter and Growth plans are now available via self-serve. Organizations
  can apply and subscribe without contacting sales.

## 2026-08-14

### New

- Administrators can now
  [select a product license when inviting a member](/manuals/accounts/organization/manage/manage-licenses.md#licenses-and-invites).
  Docker assigns the license when the invitee accepts.

## 2026-07-31

### New

- Administrators can now create
  [OIDC connections](/manuals/security/authentication/oidc-connections/_index.md)
  so GitHub Actions workflows authenticate to Docker with short-lived tokens
  instead of stored personal or organization access tokens. Available for
  Docker Team, Docker Business, Docker Hardened Images, and Docker Sponsored
  Open Source organizations.

## 2026-06-18

### New

- Custom roles now include
  [AI Governance permissions](/manuals/security/roles-and-permissions/custom-roles/permissions-reference.md#ai-governance)
  so owners can delegate policy management to other users and teams.

## 2026-06-02

### New

- Administrators can now provision products to organization members with
  [licenses](/manuals/accounts/organization/manage/manage-licenses.md).
  Licenses were introduced with AI Governance. Owners can assign or revoke
  them from the Members page, or turn on automatic assignment when a member
  first uses a supported product.

## 2026-05-19

### New

- You can now purchase
  [Gordon Plus, Max, and Ultra plans](/manuals/subscription-billing/plans/gordon.md)
  for personal accounts from the billing portal in Docker Home.
- Organizations can now purchase
  [DHI Select](/manuals/subscription-billing/plans/dhi.md) repositories via
  self-serve from the billing portal in Docker Home.

## 2026-05-12

### New

- [AI Governance](/manuals/subscription-billing/plans/ai-governance.md) is
  now available. Administrators can purchase licenses through sales,
  [assign them to members](/manuals/accounts/organization/manage/manage-licenses.md),
  and enforce
  [organization policies](/manuals/ai/sandboxes/governance/_index.md) for
  Docker AI products from Docker Home.

## 2026-03-03

### New

- [DHI Select](/manuals/subscription-billing/plans/dhi.md) is now available
  as a Docker Hardened Images plan for organizations that need SLA-backed
  patching and compliance-ready images.

## 2026-02-18

### New

- Administrators can now
  [configure DVP analytics settings](/manuals/docker-hub/repos/manage/trusted-content/insights-analytics.md#configure-dvp-analytics-settings)
  for consuming domain and benchmark report allocations in the Admin
  Console.

## 2026-02-13

### New

- Administrators can now control whether organization members can push content
  to their personal namespaces on Docker Hub with
  [namespace access control](/manuals/desktop/enterprise/hardened-desktop/namespace-access.md).
- Administrators can now prevent creating public repositories within
  organization namespaces using the
  [Disable public repositories](/manuals/docker-hub/settings.md#disable-creation-of-public-repos)
  setting.

## 2026-01-27

### New

- Administrators can now use an allow list with
  [Image Access Management](/manuals/desktop/enterprise/hardened-desktop/image-access-management.md)
  to approve specific repositories that bypass image access controls.

## 2025-11-04

### New

- Owners can now create
  [custom roles](/manuals/security/roles-and-permissions/custom-roles/_index.md)
  and assign them to members and teams.

## 2025-10-28

### New

- Docker Business subscribers can now add
  [Premium Support](/manuals/support/_index.md#paid-subscription-support),
  with faster response times and 24/7 availability.

## 2025-10-22

### New

- Organizations can now
  [pay by invoice](/manuals/subscription-billing/manage/payment-method.md#pay-by-invoice).

## 2025-10-14

### Bug fixes and enhancements

- Docker Home now keeps
  [activity logs](/manuals/accounts/organization/activity-logs.md#access-activity-logs)
  for 30 days. Use the Docker Hub API to retrieve older events.

## 2025-10-07

### Bug fixes and enhancements

- Organization management has moved out of Docker Hub. Manage organizations in
  Docker Home.

## 2025-06-30

### New

- Organization owners can now
  [export Docker Desktop user data](/manuals/accounts/organization/insights.md#export-docker-desktop-user-data)
  from Insights as a CSV file.

## 2025-06-23

### Bug fixes and enhancements

- [Organization access tokens](/manuals/security/access-tokens/organization-access-tokens.md)
  now work with Docker Scout.

## 2025-06-18

### New

- Organization owners can now
  [resend invitations in bulk](/manuals/accounts/organization/manage/members.md#manage-invitations)
  from the Members page.

## 2025-06-10

### Bug fixes and enhancements

- [Activity logs](/manuals/accounts/organization/activity-logs.md) now record
  single sign-on connection changes, and changes to SCIM and just-in-time
  provisioning.

## 2025-05-12

### New

- You can now connect
  [more than one identity provider](/manuals/security/authentication/single-sign-on/connect.md#configure-multiple-idps)
  to a single sign-on domain. Users choose a provider when they sign in with
  SSO.

## 2025-04-30

### New

- You can now pay for a subscription with a
  [verified US bank account](/manuals/subscription-billing/manage/payment-method.md#verify-a-bank-account).

## 2025-04-22

### Bug fixes and enhancements

- [Personal access tokens](/manuals/security/access-tokens/personal-access-tokens.md)
  now transfer to the organization owners when you
  [convert a user account into an organization](/manuals/accounts/organization/setup/convert-account.md).

## 2025-04-08

### New

- Organization owners can now onboard an organization with
  [guided setup](/manuals/accounts/organization/setup/onboard.md#onboard-with-guided-setup)
  in Docker Home.

## 2025-04-01

### New

- [Organization access tokens](/manuals/security/access-tokens/organization-access-tokens.md)
  are now generally available.
- Single sign-on now supports the
  [`dockerSessionMinutes` attribute](/manuals/security/provisioning/_index.md#sso-attributes),
  so a session can follow the identity provider timeout.
- Administrators can now
  [track whether users comply with Docker Desktop settings policies](/manuals/desktop/enterprise/hardened-desktop/settings-management/compliance-reporting.md)
  from Docker Home (Early Access). Compliance status is reported by Docker
  Desktop version 4.40 and later.

## 2025-03-12

### New

- You can now
  [disconnect a linked Google or GitHub account](/manuals/accounts/individual/manage-account.md#manage-connected-accounts)
  from Account settings.

## 2025-03-11

### New

- [Organization access tokens](/manuals/security/access-tokens/reference.md#available-scopes)
  now include repository scopes and organization management scopes for members,
  invites, and groups.

## 2025-02-21

### Bug fixes and enhancements

- Organization access tokens now work with Docker Build Cloud and the Docker
  Hub APIs. Company owners can manage them.

## 2025-02-11

### New

- Docker Home and the Docker Admin Console are now generally available.

## 2025-01-31

### Bug fixes and enhancements

- Docker began collecting VAT for all European countries on March 1, 2025. See
  [Sales tax exemption and VAT](/manuals/subscription-billing/manage/tax-certificate.md).

## 2025-01-30

### New

- Installing Docker Desktop via the PKG installer is now generally available.
- Enforcing sign-in via configuration profiles is now generally available.

## 2025-01-10

### Bug fixes and enhancements

- [Activity logs](/manuals/accounts/organization/activity-logs.md) now record
  when a settings policy is created, updated, deleted, or transferred.

## 2024-12-10

### New

- New Docker subscriptions are now available. For more information, see
  [Docker subscriptions and features](https://www.docker.com/pricing?ref=Docs&refAction=DocsPlatformReleaseNotes)
  and
  [Announcing Upgraded Docker Plans: Simpler, More Value, Better Development and Productivity](https://www.docker.com/blog/november-2024-updated-plans-announcement/).

## 2024-11-18

### New

- Administrators can now:
  - Enforce sign-in with
    [configuration profiles](/manuals/desktop/enterprise/enforce-sign-in/methods.md#configuration-profiles-method-mac-only)
    (Early Access).
  - Enforce sign-in for more than one organization at a time (Early Access).
  - Deploy Docker Desktop for Mac in bulk with the
    [PKG installer](/manuals/desktop/enterprise/enterprise-deployment/pkg-install-and-configure.md)
    (Early Access).
  - [Use Desktop Settings Management via the Docker Admin Console](/manuals/desktop/enterprise/hardened-desktop/settings-management/configure-admin-console.md)
    (Early Access).

### Bug fixes and enhancements

- Enhanced Container Isolation (ECI) has been improved to:
  - Permit administrators to
    [turn off Docker socket mount restrictions](/manuals/desktop/enterprise/hardened-desktop/enhanced-container-isolation/config.md#allowing-all-containers-to-mount-the-docker-socket).
  - Support wildcard tags when using the
    [`allowedDerivedImages` setting](/manuals/desktop/enterprise/hardened-desktop/enhanced-container-isolation/config.md#docker-socket-mount-permissions-for-derived-images).

## 2024-11-11

### New

- [Personal access tokens](/manuals/security/access-tokens/personal-access-tokens.md)
  (PATs) now support expiration dates.

## 2024-10-15

### New

- Beta: You can now create
  [organization access tokens](/manuals/security/access-tokens/organization-access-tokens.md)
  (OATs) to enhance security for organizations and streamline access
  management for organizations in the Docker Admin Console.

## 2024-08-29

### New

- Deploying Docker Desktop via the
  [MSI installer](/manuals/desktop/enterprise/enterprise-deployment/msi-install-and-configure.md)
  is now generally available.
- Two new methods to
  [enforce sign-in](/manuals/desktop/enterprise/enforce-sign-in/_index.md)
  (Windows registry key and `.plist` file) are now generally available.

## 2024-08-24

### New

- Administrators can now view
  [organization Insights](/manuals/accounts/organization/insights.md).

## 2024-07-17

### New

- You can now centrally access and manage Docker products in
  [Docker Home](https://app.docker.com).
