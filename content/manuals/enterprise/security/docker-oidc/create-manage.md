---
title: Create and manage Docker OIDC connections
linkTitle: Create and manage connections
description: Create, update, and delete Docker OIDC connections for your organization
keywords: docker oidc, create oidc connection, manage oidc connection, openid connect, identity provider, enterprise security, admin
tags: [admin]
weight: 10
---

{{< summary-bar feature_name="Docker OIDC" >}}

You can create new OIDC connections or manage existing ones from the Admin Console in Docker Home. Establishing an OIDC connection occurs in two phases. First, you create the OIDC connection in the Admin Console, then you configure your GitHub Actions workflow YAML file.

## Connect Docker OIDC to GitHub Actions

### Step 1: Create the Docker OIDC connection

1. Sign in to [Docker Home](https://app.docker.com/), select your organization, then go to the **Admin Console**.
1. In **Security**, select **OIDC connections**.
1. Select **Create OIDC connection** to go to the creation page. Fill in the OIDC connection form.
   - You must provide rulesets and subject claims. Other values are optional.
   - To learn about rulesets, subject claims, and resources, see [Docker OIDC rulesets and subject claims](/manuals/enterprise/security/docker-oidc/rulesets-claims.md).
1. Select **Create connection**.
1. Copy your OIDC connection ID.

### Step 2: Define GitHub Actions workflow

1. Add a top-level `permissions` key that requests a GitHub OIDC ID token:

   ```yaml
   permissions:
     id-token: write
   ```

1. Define a job that triggers the OIDC exchange. Update `connection_id` with the connection ID you copied from Docker:

   ```yaml
   jobs:
     login:
       runs-on: ubuntu-latest
       steps:
         - name: Docker OIDC
           id: docker_oidc
           uses: docker/oidc-action@v0
           with:
             connection_id: <YOUR_CONNECTION_ID>
   ```

1. Add a step that signs into Docker with an access token once ID token passes authentication:

   ```yaml
   - name: Login to Docker Hub
     uses: docker/login-action@{{% param "login_action_version" %}}
     with:
       username: <YOUR_ORGANIZATION_NAME>
       password: ${{ steps.docker_oidc.outputs.token }}
   ```

   The `username` value must be an organization name. Personal accounts are not supported.

Your updated workflow YAML should look like this:

```yaml
permissions:
  id-token: write

jobs:
  login:
    runs-on: ubuntu-latest
    steps:
      - name: Docker OIDC
        id: docker_oidc
        uses: docker/oidc-action@v0
        with:
          connection_id: <YOUR_CONNECTION_ID>

      - name: Login to Docker Hub
        uses: docker/login-action@{{% param "login_action_version" %}}
        with:
          username: <YOUR_ORGANIZATION_NAME>
          password: ${{ steps.docker_oidc.outputs.token }}
```

### Step 3 (optional): Test

After both phases, open the workflow run in GitHub Actions and select **Stage** to test the job.

## Manage OIDC connections

You can view, edit, deactivate, or delete your connections from the **OIDC connection** page.

1. From the **Admin Console**, go to **OIDC connections**.
1. From the **OIDC connections** page, find the row with your target connection ID.
1. Select the action menu icon for your options.
   - **Edit** opens the **Edit OIDC connection** page where you can copy your connection ID, update rulesets, or view the **Failures** table.
   - **Deactivate** temporarily disables access to your GitHub workflow.
   - **Activate** restores access to your GitHub workflow.
   - **Delete** permanently deletes a connection.

## Deactivation and deletion

You can deactivate an OIDC connection to pause GitHub workflow access to your Docker resources without deleting the connection. While a connection is deactivated:

  - It cannot issue Docker access tokens.
  - Without Docker access tokens, the `docker/oidc-action` step references will fail at the token-exchange step until you activate the connection.

Unlike deactivation, deleting an OIDC connection is permanent. Any workflow whose `docker/oidc-action` step still references the deleted `connection_id` will fail at the token-exchange step, so update that input with a replacement connection's id in every affected workflow before it runs again.

## What’s next

- To update your Docker OIDC connection, see [Manage OIDC connections](#manage-oidc-connections)
- For reference documentation about Docker OIDC rulesets and behaviors, see [Docker OIDC rulesets and subject claims](/manuals/enterprise/security/docker-oidc/rulesets-claims.md)
