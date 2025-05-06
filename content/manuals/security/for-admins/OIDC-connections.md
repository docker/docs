---
description: Learn how to create and manage OIDC connections in Docker.
keywords: OIDC, GitHub Actions, tokens
title: OIDC connections
---

{{< summary-bar feature_name="OIDC connections" >}}

Docker's OIDC connections enable secure authentication between GitHub Actions workflows and Docker without storing long-lived credentials. By using short-lived tokens, it reduces the risk of credential exposure and allows for fine-grained access control to your Docker resources.

When a GitHub Actions workflow needs to authenticate with Docker services, it can obtain a temporary OIDC token and exchange it for a Docker-issued access token.

## Authentication flow

1. A GitHub Actions workflow requests an OIDC token.
2. GitHub issues a signed OIDC token.
3. The workflow sends this token to Docker’s token service.
4. Docker validates the token’s signature.
5. If the token is valid and matches the configured rules, Docker returns a short-lived access token.
6. The workflow uses this access token to access Docker services.

GitHub OIDC tokens identify the calling workflow or organization and are treated like passwords for authentication. They are valid for a short period, minimizing the risk of compromise.

## Core concepts

Docker Hub manages GitHub OIDC authentication through configurations and rulesets.

- Configurations define the type of OIDC connection. Currently, Docker only supports GitHub Actions connection types.
- Rulesets determine which GitHub workflows can authenticate and what resources they can access. Each ruleset includes:
  - Label: The name of the ruleset.
  - Subject claims: Conditions based on OIDC claims (e.g., repository name, branch, workflow).
  - Resources: Docker Hub repositories or Docker Cloud resources accessible when the ruleset matches.
  - Scopes: The permissions granted.

You can define 1 to 5 rulesets per connection. These rulesets let you tailor access to the minimum required for each workflow.

> [!IMPORTANT]
>
> If multiple rulesets match an incoming token, Docker merges the resources across those rulesets and grants access to all of them.

### Subject claims

A subject claim is a string in the OIDC token that identifies the source workflow. It typically follows this format:

```PHP
repo:<org>/<repo>:ref:refs/heads/<branch>
```

For example, a subject claim could be: `repo:octo-org/octo-repo:ref:refs/heads/main`.

You can use wildcards (`*` and `?`) to match multiple repositories or branches.
For example:

- `repo:my-org/my-repo:ref:refs/heads/main`: only the `main` branch
- `repo:my-org/*`: all repos in the org
- `repo:*:ref:refs/heads/release-*`: all branches starting with `release-`

See [GitHub’s example subject claims](https://docs.github.com/en/actions/security-for-github-actions/security-hardening-your-deployments/about-security-hardening-with-openid-connect#example-subject-claims) for more formats.

## Create an OIDC connection

To let GitHub Actions workflows authenticate with Docker Hub, create an OIDC connection and define one or more rulesets.

1. Sign in to the [Admin Console](https://app.docker.com/admin) and select your organization.
2. In the left sidebar, go to **Security and access** > **OIDC connections**.
3. Select **Create new OIDC connection**.
4. Enter a **Connection name**.
5. Optional. Add a **description**.
6. Under **Connection type**, select **GitHub Actions**.
7. Select **Add ruleset**, then:
   - Enter a **label**.
   - Enter a **subject claim**.
   - Select the Docker Hub repositories or Docker Cloud resources to allow access.
8. Select **Next** and review the configuration.
9. Select **Confirm**, then **Create OIDC connection**.

## Configure your GitHub Actions workflow

After setting up the connection and ruleset, update your workflow to use GitHub OIDC.

### Step 1: Request an OIDC token

```yaml
permissions:
  id-token: write
```

### Step 2: Exchange the OIDC token for a Docker access token

```yaml
jobs:
  login:
    runs-on: ubuntu-latest
    steps:
      - name: Docker OIDC
        id: docker_oidc
        uses: docker/configure-credentials-action
        with:
          connection_id: 77537984-f27a-45c5-9b4c-f20b3f88817e
```

### Use the access token

With Docker CLI:

```yaml
- name: Login to Docker Hub
  uses: docker/login-action@v3
  with:
    username: docker # must match your Docker organization name
    password: ${{ steps.docker_oidc.outputs.token }}
```

With Docker Hub API:

```yaml
- name: List members of the org
  run: |
    curl \
      -H "Authorization: Bearer ${{ steps.docker_oidc.outputs.token }}" \
      https://hub.docker.com/v2/orgs/${ORG_ID}/members
```

## Manage OIDC connections

After creating a connection, you can copy its ID, edit its details, or manage access based on your organization’s needs.

### Copy connection ID

You’ll need the connection ID to configure your workflows.

1. Sign in to the [Admin Console](https://app.docker.com/admin) and select your organization.
2. Go to **Security and access** > **OIDC connections**.
3. Select the **Actions menu**, then **Edit**.
4. Copy the value in the Connection ID field.

### Edit a connection

You can edit a connection to rename the connection, update rulesets, or modify resources.

1. Sign in to the [Admin Console](https://app.docker.com/admin) and select your organization.
2. In **Security and access**, select **OIDC connections**.
3. Select the **Actions menu icon**, then **Edit**.
4. You can:
    - Update connection name or description
    - Edit or remove rulesets
    - Add new rulesets
    - Update resources inside rulesets

### Deactivate a connection

You can deactivate a connection if you want to temporarily disable access or are rotating credentials and want to pause workflow authorization.

1. Sign in to the [Admin Console](https://app.docker.com/admin) and select your organization.
2. In **Security and access**, select **OIDC connections**.
3. Select the **Actions menu icon**, then **Deactivate**.

### Delete a connection

Delete a connection if it’s no longer needed or to clean up unused or outdated rulesets.

1. Sign in to the [Admin Console](https://app.docker.com/admin) and select your organization.
2. In **Security and access**, select **OIDC connections**.
3. Select the **Actions menu icon**, then **Delete**.