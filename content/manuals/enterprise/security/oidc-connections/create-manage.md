---
title: Create and manage OIDC connections
linkTitle: Create and manage connections
description: Create, update, and delete OIDC connections for your organization
keywords: oidc connections, create oidc connection, github actions, docker/login-action, openid connect, enterprise security, admin
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

### Step 2: Update your workflow

Add the OIDC connection ID and your Docker organization name as
[GitHub Actions variables](https://docs.github.com/en/actions/writing-workflows/choosing-what-your-workflow-does/store-information-in-variables)
in your repository or organization settings. Then update your workflow:

```yaml
name: ci

on:
  push:
    branches: main

permissions:
  contents: read
  id-token: write

jobs:
  login:
    runs-on: ubuntu-latest
    steps:
      -
        name: Login to Docker Hub
        uses: docker/login-action@{{% param "login_action_version" %}}
        env:
          DOCKERHUB_OIDC_CONNECTIONID: ${{ vars.DOCKERHUB_OIDC_CONNECTIONID }}
        with:
          username: ${{ vars.DOCKERHUB_ORGANIZATION }}
```

The `id-token: write` permission lets the workflow request a GitHub OIDC
token. The `docker/login-action` handles the OIDC token exchange and
Docker login in a single step when `DOCKERHUB_OIDC_CONNECTIONID` is set
and `password` is omitted.

The `username` value must be an organization name. Personal accounts
aren't supported.

Run your GitHub Action and verify the workflow can sign in to Docker.

> [!TIP]
> If your workflow needs the Docker access token as a separate output
> (for example, for API calls or custom authentication flows), use
> [`docker/oidc-action`](https://github.com/docker/oidc-action) to
> perform the token exchange explicitly.

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
- The login step in your workflow will fail at the token-exchange step
  until you activate the connection.

Unlike deactivation, deleting an OIDC connection is permanent. Any workflow
that still references the deleted connection ID will fail at the
token-exchange step. Update the connection ID in every affected workflow
before it runs again.

## What's next

- [OIDC connections rulesets and subject claims](/manuals/enterprise/security/oidc-connections/rulesets-claims.md)
