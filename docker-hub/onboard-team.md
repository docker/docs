---
description: Docker Team onboarding
keywords: team, organizations, get started, onboarding
title: Docker Team onboarding
toc_min: 1
toc_max: 2
---

The following section contains step-by-step instructions on how to get started onboarding your organization after you obtain a Docker Team subscription.

## Prerequisites

Before you start to on board your organization, ensure that you've completed the following:
- You have a Docker Team subscription. [Buy now](https://www.docker.com/pricing/) if you haven't subscribed to Docker Team yet.
- You are familiar with Docker terminology. If you discover any unfamiliar terms, see the [glossary](/glossary/#docker) or [FAQs](../docker-hub/onboarding-faqs.md).


## Step 1: Identify your Docker users and their Docker accounts

To begin, you should identify which users you will need to add to your Docker Team organization. Identifying your users will help you efficiently allocate your subscription's seats and manage access.

1. Identify the Docker users in your organization.
   - If your organization uses device management software, like MDM or JAMF, you may use the device management software to help identify Docker users. See your device management software's documentation for details. You can identify Docker users by checking if Docker Desktop is installed at the following location on each user's machine:
      - Mac: `/Applications/Docker.app`
      - Windows: `C:\Program Files\Docker\Docker`
   - If your organization doesn't use device management software, you may survey your users.
2. Instruct all your Docker users in your organization to update their existing Docker account's email address to an address that's in your organization's domain, or to create a new account using an email address in your organization's domain.
   - To update an account's email address, instruct your users to sign in to [Docker Hub](https://hub.docker.com){: target="_blank" rel="noopener" class="_"}, go to [Account Settings](https://hub.docker.com/settings/general){: target="_blank" rel="noopener" class="_"}, and update the email address to their email address in your organization's domain.
   - To create a new account, instruct your users to go [sign up](https://hub.docker.com/signup){: target="_blank" rel="noopener" class="_"} using their email address in your organization's domain.
3. Ask your Docker sales representative to provide a list of Docker accounts that use an email address in your organization's domain.

## Step 2: Add members

Now that you have a Docker Team organization, it's time to start adding members.

Before you add members, you must create at least one team. For details, see [Create a team](../docker-hub/orgs.md/#create-a-team){: target="_blank" rel="noopener" class="_"}.

Add members to your organization by inviting them to a team using their Docker ID or email address. For details, see [Invite members](../docker-hub/members.md/#invite-members){: target="_blank" rel="noopener" class="_"}.

## Step 3: Enforce sign-in for Docker Desktop

By default, members of your organization can use Docker Desktop on their machines without signing in to any Docker account. To ensure that a user signs in to a Docker account that is a member of your organization and that the
organization’s settings apply to the user’s session, you can use a `registry.json` file. For details, see [Configure registry.json to enforce sign-in](../docker-hub/configure-sign-in.md){: target="_blank" rel="noopener" class="_"}.

## What's next

Create and manage your repositories:

- Create [repositories](../docker-hub/repos/index.md) to share container images.
- [Consolidate a repository](../docker-hub/repos/index.md/#consolidating-a-repository) from your personal account to your organization.
- Create [teams](../docker-hub/orgs.md/#create-a-team) and configure [repository permissions](../docker-hub/orgs.md/#configure-repository-permissions).

Your Docker Team subscription provides many more [additional features](../subscription/index.md).


