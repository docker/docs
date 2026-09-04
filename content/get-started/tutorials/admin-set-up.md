---
title: Roll out Docker to a pilot group
linkTitle: Pilot an organization rollout
description: Invite a pilot user, apply a Docker Desktop settings policy, and verify the managed experience before a wider deployment.
keywords: Docker, get started, administration, organization setup, pilot deployment, Settings Management
weight: 3
---

Rolling out Docker across an organization changes how developers sign in and
how Docker Desktop behaves. Testing the managed experience with a small group
helps uncover access or workflow problems before they affect everyone.

In this 15-minute tutorial, you'll invite one pilot user, apply one managed
setting to that user, and confirm that Docker Desktop receives the policy.

## Before you start

You need:

- A Docker Business subscription
- Owner access to a Docker organization
- A [verified company domain](/manuals/security/authentication/single-sign-on/connect.md#step-1-add-a-domain)
- [Enforced Docker Desktop sign-in](/manuals/enterprise/security/enforce-sign-in/_index.md)
- One pilot user with Docker Desktop installed

## Invite the pilot user

1. Sign in to [Docker Home](https://app.docker.com) and select your
   organization.
1. Select **Members**, then **Invite**.
1. Select **Emails or usernames** and enter the pilot user's email address.
1. Assign the **Member** role and a Docker Desktop license, then send the
   invitation.
1. Have the pilot user accept the invitation from the email or the
   **Notifications Center**.

The pilot user is now part of the organization and can receive its Docker
Desktop policies.

## Create a pilot policy

1. In Docker Home, select **Docker Desktop**, then **Settings Management**.
1. Select **Create a settings policy**.
1. Name the policy `Docker pilot`.
1. Apply the policy to **Specific users**, then select the pilot user.
1. Set **Automatically update components** to **Always enabled**. Leave the
   other settings as **User-defined**.
1. Select **Create**.

This policy locks one low-impact setting, making the result visible without
changing the pilot user's development workflow.

## Apply and verify the policy

Ask the pilot user to:

1. Quit Docker Desktop completely.
1. Start Docker Desktop and sign in with the account that joined the
   organization.
1. Open Docker Desktop settings and confirm that **Automatically update
   components** is enabled and can't be changed.
1. Run a familiar development project to confirm their normal workflow still
   works.

Docker Desktop checks for policy updates when it starts and every 60 minutes
while running. Fully quitting and reopening it applies the policy without
waiting for the next check.

## What you learned

You added a user, targeted a policy, and verified the managed Docker
Desktop experience before applying controls across the organization. Use the
pilot's feedback to refine the policy and its deployment communication.

## What's next

Continue with these guides when you're ready to expand the deployment:

- [Set up your company for success with Docker](/guides/admin-set-up/) for the
  full deployment handbook
- [Manage organization members](/manuals/accounts/organization/manage/members.md)
  for roles, teams, and bulk invitations
- [Configure Settings Management](/manuals/enterprise/security/hardened-desktop/settings-management/configure-admin-console.md)
  for additional policies and deployment controls
- [Set up SSO and provisioning](/manuals/security/authentication/single-sign-on/_index.md)
  for centralized identity management
- [Govern Docker Sandboxes](/manuals/ai/sandboxes/governance/_index.md) for
  organization-wide agent policies
