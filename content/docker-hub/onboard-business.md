---
description: Docker Business onboarding
keywords: business, organizations, get started, onboarding
title: Docker Business onboarding
toc_min: 1
toc_max: 2
---

The following section contains step-by-step instructions on how to get started onboarding your organization after you obtain a Docker Business subscription.

## Prerequisites

Before you start to on board your organization, ensure that you've completed the following:
- You have a Docker Business subscription. [Get in touch with us](https://www.docker.com/pricing/contact-sales/) if you haven't subscribed to Docker Business yet.
- Your Docker Business subscription is new. If you upgraded your Docker Team subscription or renewed your Docker Business subscription, see [what's next](#whats-next).
- Your Docker Business subscription has started. You can't complete all the steps until after your subscription start date.
-  You are familiar with Docker terminology. If you discover any unfamiliar terms, see the [glossary](/glossary/#docker) or [FAQs](../docker-hub/onboarding-faqs.md).

## Step 1: Identify your Docker users and their Docker accounts

To begin, you should identify which users you will need to add to your Docker Business organization. Identifying your users will help you efficiently allocate your subscription's seats and manage access.

1. Identify the Docker users in your organization.
   - If your organization uses device management software, like MDM or JAMF, you may use the device management software to help identify Docker users. See your device management software's documentation for details. You can identify Docker users by checking if Docker Desktop is installed at the following location on each user's machine:
      - Mac: `/Applications/Docker.app`
      - Windows: `C:\Program Files\Docker\Docker`
   - If your organization doesn't use device management software, you may survey your users.
2. Instruct all your Docker users in your organization to update their existing Docker account's email address to an address that's in your organization's domain, or to create a new account using an email address in your organization's domain.
   - To update an account's email address, instruct your users to sign in to [Docker Hub](https://hub.docker.com){: target="_blank" rel="noopener" class="_"}, go to [Account Settings](https://hub.docker.com/settings/general){: target="_blank" rel="noopener" class="_"}, and update the email address to their email address in your organization's domain.
   - To create a new account, instruct your users to go [sign up](https://hub.docker.com/signup){: target="_blank" rel="noopener" class="_"} using their email address in your organization's domain.
3. Ask your Docker sales representative to provide a list of Docker accounts that use an email address in your organization's domain.

## Step 2: Add your Docker Business subscription to an organization

On the day that your Docker Business subscription starts, your organization's primary contact will receive a welcome email from Docker to guide you through creating a new organization or to let you choose an existing organization for your Docker Business subscription.

> **Note**
>
> If your organization's primary contact doesn't receive a welcome email from Docker on the day that your subscription starts:
>   - Check your email spam folder.
>   - Use the steps below to verify that your Docker Business organization doesn't already exist.
>   - Contact your Docker sales representative to verify your primary contact's email address.

After completing the steps from the welcome email, verify that your organization exists and your organization has a Docker Business subscription:

1. Go to [Billing Details](https://hub.docker.com/billing){: target="_blank" rel="noopener" class="_"} and then select on your organization's name.
2. Under **Plan**, view your subscription. If you organization has a Docker Business subscription, you will see **Docker Business**.

## Step 3: Invite owners

Now that you have a Docker Business organization, it's time to start adding owners to help you set up and manage your organization. Owners can add or remove members, and configure Single Sign-on as well as other security settings.

To add an owner, invite a user to the **owners** team. For more details, see [Invite members](../docker-hub/members.md/#invite-members){: target="_blank" rel="noopener" class="_"}.

## Step 4: Invite members

Add members to your organization using a CSV file, or by entering their email addresses. For more details, see [Invite members](../docker-hub/members.md/#invite-members){: target="_blank" rel="noopener" class="_"}

## Step 5: Manage members with SSO and SCIM

Automate adding members to your organization using the following:
   - Single Sign-on: Automatically provision and add members when they sign in to Docker Hub through your identity provider. For details, see [Single Sign-on overview](../single-sign-on/index.md).
   - System for Cross-domain Identity Management: Automatically provision, add, and de-provision members from your identity provider. For details, see [SCIM](../docker-hub/scim.md).

## Step 6: Enforce sign-in for Docker Desktop

By default, members of your organization can use Docker Desktop on their machines without signing in to any Docker account. To ensure that a user signs in to a Docker account that is a member of your organization and that the
organization’s settings apply to the user’s session, you can use a `registry.json` file. For details, see [Configure registry.json to enforce sign-in](../docker-hub/configure-sign-in.md){: target="_blank" rel="noopener" class="_"}.

## What's next

Configure security settings and manage your repositories:

- Create [repositories](../docker-hub/repos/index.md) to share container images.
- [Consolidate a repository](../docker-hub/repos/index.md/#consolidating-a-repository) from your personal account to your organization.
- Create [teams](../docker-hub/orgs.md/#create-a-team) and configure [repository permissions](../docker-hub/orgs.md/#configure-repository-permissions).
- Configure [Hardened Docker Desktop](../desktop/hardened-desktop/index.md) to improve your organization’s security posture for containerized development. Hardened Docker Desktop includes:
   - [Settings Management](../desktop/hardened-desktop/settings-management/index.md), which helps you to confidently manage and control the usage of Docker Desktop within your organization.
   - [Enhanced Container Isolation](../desktop/hardened-desktop/enhanced-container-isolation/index.md), a setting that instantly enhances security by preventing containers from running as root in Docker Desktop’s Linux VM.
   - [Image Access Management](../docker-hub/image-access-management.md/), lets you control which images developers can pull from Docker Hub.
   - [Registry Access Management](../docker-hub/registry-access-management.md/), lets you control the registries developers can access.

Your Docker Business subscription provides many more additional features. [Learn more](../subscription/index.md).
