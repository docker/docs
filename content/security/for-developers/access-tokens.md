---
title: Create and manage access tokens
description: Learn how to create and manage your personal Docker Hub access tokens
  to securely push and pull images programmatically.
keywords: docker hub, hub, security, PAT, personal access token
aliases: 
- /docker-hub/access-tokens/
---

You can create personal access tokens (PAT) to use as alternatives to your password for Docker CLI authentication.

Compared to passwords, personal access tokens provide the following advantages:

- You can investigate when the PAT was last used and then disable or delete it if you find any suspicious activity.
- When using an access token, you can't perform any admin activity on the account, including changing the password. It protects your account if your computer is compromised.
  
Access tokens are also valuable for building integrations, as you can issue multiple tokens, one for each integration, and revoke them at
any time.
   > **Note**
   >
   > If you have [two-factor authentication (2FA)](2fa/index.md) enabled on
   > your account, you must create at least one personal access token. Otherwise,
   > you won't be able to sign in to your account from the Docker CLI.

## Create an access token

> **Important**
>
> Treat access tokens like your password and keep them secret. Store your tokens securely in a credential manager for example.
{ .important }

1. Sign in to [Docker Hub](https://hub.docker.com).

2. Select your username in the top-right corner and from the drop-down menu select **Account Settings**.

3. Select the **Security** tab and then **New Access Token**.

4. Add a description for your token. Use something that indicates the use case or purpose of the token.
   
5. Set the access permissions.
   The access permissions are scopes that set restrictions in your
   repositories. For example, for Read & Write permissions, an automation
   pipeline can build an image and then push it to a repository. However, it
   can not delete the repository.

6. Select **Generate** and then copy the token that appears on the screen and save it. You won't be able
   to retrieve the token once you close this prompt.

## Use an access token

You can use an access token in place of your password when you log in using Docker CLI.

Log in from your Docker CLI client with the following command, replacing `YOUR_USERNAME` with your Docker ID:

```console
$ docker login --username <YOUR_USERNAME>
```

When prompted for a password, enter your personal access token instead of a password.

> **Note**
>
> If you have [two-factor authentication (2FA)](2fa/index.md) enabled, you must
> use a personal access token when logging in from the Docker CLI. 2FA is an
> optional, but more secure method of authentication.

## Modify existing tokens

You can rename, activate, deactivate, or delete a token as needed.

1. Access your tokens under **Account Settings > Security**.
   This page shows an overview of all your tokens, and lists if the token was generated manually or if it was [auto-generated](#auto-generated-tokens). You can also view the number
   of tokens that are activated and deactivated in the toolbar.

2. Select a token and then select **Delete** or **Edit**, or use the actions menu on the far right of a token row to bring up the edit screen.
   You can also select multiple tokens to delete at once.

3. After modifying the token, select **Save**.

## Auto-generated tokens

When you sign in to your Docker account with Docker Desktop, Docker Desktop generates an authentication token on your behalf. When you interact with Docker Hub using the Docker CLI, the CLI uses this token for authentication. The token scope has Read, Write, and Delete access. If your Docker Desktop session expires, the token is automatically removed locally.

You can have up to 5 auto-generated tokens associated with your account. These are deleted and created automatically based on usage and creation dates. You can also delete your auto-generated tokens as needed. See [Modify existing tokens](#modify-existing-tokens).