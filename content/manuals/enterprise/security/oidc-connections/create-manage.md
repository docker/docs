---
title: Create and manage OIDC connections
linkTitle: Create and manage connections
description: Create, update, and delete OIDC connections for your organization
keywords: oidc connections, create oidc connection, github actions, docker/oidc-action, openid connect, enterprise security, admin
tags: [admin]
weight: 10
---

{{< summary-bar feature_name="OIDC connections" >}}

Organization owners and editors can create OIDC connections or manage
existing ones from **OIDC connections** in Docker Home. Establishing an
OIDC connection occurs in two phases. First, you create the OIDC
connection in Docker Home, then you configure your GitHub Actions workflow
YAML file.

> [!NOTE]
> OIDC connections support only GitHub as a trusted third party.

## Set up GitHub Actions authentication

### Step 1: Create the OIDC connection

1. Sign in to [Docker Home](https://app.docker.com/), select your
   organization, then go to **Identity & auth**.
1. Select **OIDC connections**.
1. Select **Create OIDC connection** and fill in the OIDC connection form.
   - You must provide rulesets and subject claims. Other values are
     optional.
   - For rulesets, subject claims, and resources, see
     [OIDC connections rulesets and subject claims](/manuals/enterprise/security/oidc-connections/rulesets-claims.md).
1. Select **Create connection**.
1. Copy your OIDC connection ID.

### Step 2: Define the GitHub Actions workflow

1. Add a top-level `permissions` key that requests a GitHub OIDC ID token:

   ```yaml
   permissions:
     id-token: write
   ```

1. Define a job that triggers the OIDC exchange. Update `connection_id`
   with the connection ID you copied from Docker:

   ```yaml
   jobs:
     login:
       runs-on: ubuntu-latest
       steps:
         - name: OIDC connections
           id: docker_oidc
           uses: docker/oidc-action@v1
           with:
             connection_id: <YOUR_CONNECTION_ID>
   ```

1. Add a step that signs in to Docker with an access token once the ID
   token passes authentication:

   ```yaml
   - name: Sign in to Docker Hub
     uses: docker/login-action@{{% param "login_action_version" %}}
     with:
       username: <DOCKER_ORGANIZATION_NAME>
       password: ${{ steps.docker_oidc.outputs.token }}
   ```

   The `username` value must be an organization name. Personal accounts
   aren't supported.

   Your updated workflow YAML should look like this:

   ```yaml
   permissions:
     id-token: write

   jobs:
     login:
       runs-on: ubuntu-latest
       steps:
         - name: OIDC connections
           id: docker_oidc
           uses: docker/oidc-action@v1
           with:
             connection_id: <YOUR_CONNECTION_ID>

         - name: Sign in to Docker Hub
           uses: docker/login-action@{{% param "login_action_version" %}}
           with:
             username: <YOUR_ORGANIZATION_NAME>
             password: ${{ steps.docker_oidc.outputs.token }}
   ```

1. Run your GitHub Action and verify the workflow can sign in to Docker.

## Manage OIDC connections

You can view, edit, deactivate, or delete connections from the **OIDC
connections** page.

1. From **Identity & auth**, go to **OIDC connections**.
1. From the **OIDC connections** page, find the row with your target
   connection ID.
1. Select the action menu icon for your options.
   - **Edit** opens the **Edit OIDC connection** page where you can copy
     your connection ID, update rulesets, or view the **Failures** table.
   - **Deactivate** temporarily disables access to your GitHub workflow.
   - **Activate** restores access to your GitHub workflow.
   - **Delete** permanently deletes a connection.

## Deactivation and deletion

You can deactivate an OIDC connection to pause GitHub workflow access to
your Docker resources without deleting the connection. While a connection
is deactivated:

- It can't issue Docker access tokens.
- Without Docker access tokens, `docker/oidc-action` fails at the
  token-exchange step until you activate the connection.

Unlike deactivation, deleting an OIDC connection is permanent. Any workflow
whose `docker/oidc-action` step still references the deleted
`connection_id` fails at the token-exchange step. Update that input with a
replacement connection's ID in every affected workflow before it runs
again.

## What's next

- [OIDC connections rulesets and subject claims](/manuals/enterprise/security/oidc-connections/rulesets-claims.md)
