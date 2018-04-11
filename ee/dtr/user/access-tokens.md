---
title: Manage access tokens
description: Learn how to create and manage your personal DTR access tokens to securely
  integrate DTR with other products.
keywords: dtr, security, access tokens
redirect_from:
  - /datacenter/dtr/2.5/guides/user/access-tokens/
---

Docker Trusted Registry allows you to issue access tokens so that you can
integrate with other services without having to give those services your
credentials. An access token is issued for a user and has the same DTR
permissions the user has.

It's better to use access tokens to build integrations since you can issue
multiple tokens, one for each integration, and revoke them at any time.

## Create an access token

In the **DTR web UI**, navigate to your user profile, and choose **Access tokens**.

![Token list](../images/access-tokens-1.png){: .with-border}

Click **New access token**, and assign a meaningful name to your token.
Choose a name that indicates where the token is going to be used, or what’s the
purpose for the token. Administrators can also create tokens for other users.

![Create token](../images/access-tokens-2.png){: .with-border}

Once the token is created you won’t be able to see it again, but you can
rename it if needed.

## Use the access token

You can use an access token in any place that requires your DTR password.
As an example you can use access tokens to login in from your Docker CLI client:

```bash
docker login dtr.example.org --username <username> --password <token>
```

To use the DTR API to list the repositories your user has access to:

```bash
curl --silent --insecure --user <username>:<token> dtr.example.org/api/v0/repositories
```

