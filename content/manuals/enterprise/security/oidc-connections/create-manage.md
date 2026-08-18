---
title: Create and manage OIDC connections
linkTitle: Create and manage connections
description: Create and manage OIDC connections in Docker Home, then authenticate GitHub Actions to Docker Hub with docker/login-action and short-lived tokens
keywords: oidc connections, create oidc connection, github actions, docker/login-action, DOCKERHUB_OIDC_CONNECTIONID, openid connect, docker hub, enterprise security, admin
tags: [admin]
weight: 10
---

{{< summary-bar feature_name="OIDC connections" >}}

Organization owners and editors create and manage OIDC connections from
**OIDC connections** in Docker Home. After you create a connection,
configure your GitHub Actions workflow so it can sign in to Docker Hub
with a short-lived token.

> [!NOTE]
> OIDC connections support only GitHub as a trusted third party.

## Create an OIDC connection

1. Sign in to [Docker Home](https://app.docker.com/), select your
   organization, then go to **Identity & auth**.
1. Select **OIDC connections**.
1. Select **Create OIDC connection** and fill in the OIDC connection form.
   - Provide rulesets and subject claims. Other values are optional.
   - For rulesets, subject claims, and resources, see
     [OIDC connections rulesets and subject claims](/manuals/enterprise/security/oidc-connections/rulesets-claims.md).
1. Select **Create connection**.
1. Copy your OIDC connection ID.

## Configure a GitHub Actions workflow

Use `docker/login-action` version 4.5.0 or later. The action exchanges
the GitHub OIDC token and signs in to Docker Hub in a single step.

1. Add the following to your workflow YAML. Replace `<YOUR_ORG_NAME>`
   with your Docker organization name and `<YOUR_CONNECTION_ID>` with
   the ID you copied from Docker Home:

   ```yaml
   permissions:
     contents: read
     id-token: write

   jobs:
     build:
       runs-on: ubuntu-latest
       steps:
         - name: Docker login
           uses: docker/login-action@v4 # v4.5.0+
           with:
             username: <YOUR_ORG_NAME>
           env:
             DOCKERHUB_OIDC_CONNECTIONID: <YOUR_CONNECTION_ID>
   ```

   The `username` value must be a Docker organization name. Only
   organization accounts can sign in using OIDC.

1. Run the workflow and confirm it can sign in to Docker.

## Manage OIDC connections

You can view, edit, deactivate, or delete connections from the **OIDC
connections** page.

1. Select **Identity & auth**, then **OIDC connections**.
1. Find the row with your target connection ID.
1. Select the action menu icon for your options.
   - **Edit** opens the **Edit OIDC connection** page where you can copy
     your connection ID, update rulesets, or view the **Failures** table.
   - **Deactivate** temporarily disables access to your GitHub workflow.
   - **Activate** restores access to your GitHub workflow.
   - **Delete** permanently deletes a connection.

## Deactivate or delete a connection

Deactivate an OIDC connection to pause GitHub workflow access to your
Docker resources without deleting the connection. A deactivated
connection does not issue Docker access tokens, so `docker/login-action`
fails at the token-exchange step until you activate the connection.

> [!WARNING]
> Deleting an OIDC connection is permanent. Any workflow that still sets
> `DOCKERHUB_OIDC_CONNECTIONID` to the deleted ID fails at the
> token-exchange step. Update that environment variable with a
> replacement connection ID in every affected workflow before it runs
> again.

## Next steps

- [OIDC connections rulesets and subject claims](/manuals/enterprise/security/oidc-connections/rulesets-claims.md)
