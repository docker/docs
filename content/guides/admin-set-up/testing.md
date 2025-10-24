---
title: Testing
description: Test your Docker setup.
weight: 30
---

## SSO and SCIM testing

Test SSO and SCIM by signing in to Docker Desktop or Docker Hub with the email
address linked to a Docker account that is part of the verified domain.
Developers who sign in using their Docker usernames remain unaffected by the
SSO and SCIM setup.

> [!IMPORTANT]
>
> Some users may need CLI based logins to Docker Hub, and for this they will
need a [personal access token (PAT)](/manuals/security/access-tokens.md).

## Test Registry Access Management and Image Access Management

> [!WARNING]
>
> Communicate with your users before proceeding, as this step will impact all
existing users signing into your Docker organization.

If you plan to use [Registry Access Management (RAM)](/manuals/enterprise/security/hardened-desktop/registry-access-management.md) and/or [Image Access Management (IAM)](/manuals/enterprise/security/hardened-desktop/image-access-management.md):

1. Ensure your test developer signs in to Docker Desktop using their
    organization credentials
2. Have them attempt to pull an unauthorized image or one from a disallowed
    registry via the Docker CLI
3. Verify they receive an error message indicating that the registry is
    restricted by the organization

## Deploy settings and enforce sign in to test group

Deploy the Docker settings and enforce sign-in for a small group of test users
via MDM. Have this group test their development workflows with containers on
Docker Desktop and Docker Hub to ensure all settings and the sign-in enforcement
function as expected.

## Test Docker Build Cloud capabilities

Have one of your Docker Desktop testers [connect to the cloud builder you created and use it to build](/manuals/build-cloud/usage.md).

## Test Testcontainers Cloud

Have a test developer [connect to Testcontainers Cloud](https://testcontainers.com/cloud/docs/#getting-started) and run a container in
the cloud to verify the setup is working correctly.

## Verify Docker Scout monitoring of repositories

Check the [Docker Scout dashboard](https://scout.docker.com/) to confirm that
data is being properly received for the repositories where Docker Scout has
been enabled.

## Verify access to Docker Hardened Images

Have a test developer attempt to [pull a Docker Hardened Image](/manuals/dhi/get-started.md) to confirm that
the team has proper access and can integrate these images into their workflows.
