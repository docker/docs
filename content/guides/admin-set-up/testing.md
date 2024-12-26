---
title: Testing
description: Test your Docker setup.
weight: 30
---

## SSO and SCIM testing

You can test SSO and SCIM by signing in to Docker Desktop or Docker Hub with the email address linked to a Docker account that is part of the verified domain. Developers who sign in using their Docker usernames will remain unaffected by the SSO and/or SCIM setup. 

> [!IMPORTANT] 
>
> Some users may need CLI based logins to Docker Hub, and for this they will need a [personal access token (PAT)](/manuals/security/for-developers/access-tokens.md).

## Test RAM and IAM

> [!WARNING]
> Be sure to communicate with your users before proceeding, as this step will impact all existing users signing into your Docker organization

If you plan to use [Registry Access Management (RAM)](/manuals/security/for-admins/hardened-desktop/registry-access-management.md) and/or [Image Access Management (IAM)](/manuals/security/for-admins/hardened-desktop/image-access-management.md), ensure your test developer signs in to Docker Desktop using their organization credentials. Once authenticated, have them attempt to pull an unauthorized image or one from a disallowed registry via the Docker CLI. They should receive an error message indicating that the registry is restricted by the organization.

## Deploy settings and enforce sign in to test group

Deploy the Docker settings and enforce sign-in for a small group of test users via MDM. Have this group test their development workflows with containers on Docker Desktop and Docker Hub to ensure all settings and the sign-in enforcement function as expected.

## Test Docker Build Cloud capabilities

Have one of your Docker Desktop testers [connect to the cloud builder you created and use it to build](/manuals/build-cloud/usage.md). 

## Verify Docker Scout monitoring of repositories

Check the [Docker Scout dashboard](https://scout.docker.com/) to confirm that data is being properly received for the repositories where Docker Scout has been enabled.
