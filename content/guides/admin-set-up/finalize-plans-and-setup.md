---
title: Finalize plans and begin setup
description: Collaborate with your MDM team to distribute configurations and set up SSO and Docker product trials.
weight: 20
---

## Step one: Send finalized settings files to MDM team 

Once you have come to an agreement between with the relevant teams regarding your baseline and security configurations outlined in module one, follow the instructions in the [Settings Management]() documentation to create the `admin-settings.json` file which contains these configurations. 

Once this has been done, work with your MDM team to deploy the `admin-settings.json` file and your chosen method for [Enforcing sign-in](/manuals/security/for-admins/enforce-sign-in/_index.md).

It’s highly recommended that the next step in week 3 is a test distribution to a small number of Docker Desktop users to verify the functionality works as expected.

## Step two: Manage your organizations

If you have more than one organization, it’s recommended that you either consolidate them into one organization or use the account hierarchy feature to manage multiple organizations.  Please work with the CS and implementation teams to make this happen.

## Step three: Begin setup

### SSO domain verification

The SSO process has multiple steps involving different teams, so it's recommended that the process is started right away.  The first step is domain verification.  This step ensures that the person setting up SSO actually controls the domain they are requesting.  The detailed steps to verify a domain are located here.  Your DNS team will need to be involved in this step.

### Create SSO Connection

Once the domain is verified, the next step is to create the SSO connection.  This will involve your identity provider team to configure the identity groups and help set up the SSO connection.  Note that this step of creating the SSO connection will not affect the Docker Desktop user experience, and you will be able to test before enforcing SSO for all users.  The steps in the process are located here.  

### Set up free tier Docker product entitlements included in the subscription

Set up the cloud builder for free monthly minutes in Docker Build Cloud, and up to three repositories to monitor via Docker Scout.  Please note that your free entitlements stop when your limits are exceeded so there is no fear of a surprise cost overage.  The instructions on setting up the cloud builder are located on build.docker.com and there is a video walkthrough here, and the instructions on adding a repository for scout monitoring is here for Docker Hub repositories, and here for integration to other image registries.

### Ensure supported version of Docker Desktop

CAUTION: This step could affect the experience for users on older versions of Docker Desktop.  Existing users may have older versions of Docker Desktop that are no longer supported or are out of date.  It is highly recommended that everyone update to a supported version.  We recommend using a MDM solution to manage the version of Docker Desktop for users.  Users may also get Docker Desktop directly from Docker or through a company software portal.  In any of these cases it's important that the users are upgraded to a supported Docker Desktop version.
