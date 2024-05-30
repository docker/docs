---
description: Learn how to register for a Docker ID and sign in to your account
keywords: accounts, docker ID, billing, paid plans, support, Hub, Store, Forums, knowledge
  base, beta access, email, activation, verification
title: Create an account
aliases:
- /docker-cloud/dockerid/
- /docker-hub/accounts/
---

You can create a free Docker account with your email address or by signing up with your Google or GitHub account. Once you've created your account with a unique Docker ID, you can access all Docker products, including Docker Hub. With Docker Hub, you can access repositories and explore images that are available from the community and verified publishers.

Your Docker ID becomes your username for hosted Docker services, and [Docker forums](https://forums.docker.com/).

> **Tip**
>
> Explore [Docker's core subscriptions](https://www.docker.com/pricing/) to see what else Docker can offer you. 
{ .tip }

## Create a Docker ID

### Sign up with your email address

1. Go to the [Docker sign-up page](https://app.docker.com/signup/).

2. Enter a unique, valid email address.

3. Enter a username.

    Your Docker ID must be between 4 and 30 characters long, and can only contain numbers and lowercase letters. Once you create your Docker ID you can't reuse it in the future if you deactivate this account.

4. Enter a password that's at least 9 characters long.

5. Select **Sign Up**.

   Docker sends a verification email to the address you provided.

6. Verify your email address to complete the registration process.

> **Note**
>
> You have limited actions available until you verify your email address.

### Sign up with Google or GitHub

> **Important**
>
> To sign up with your social provider, make sure you verify your email address with your provider before you begin.
{ .important }

1. Go to the [Docker sign-up page](https://hub.docker.com/signup/).

2. Select your social provider, Google or GitHub.

3. Select the social account you want to link to your Docker account.

4. Select **Authorize Docker** to allow Docker to access your social account information and be re-routed to the sign-up page.

5. Enter a username.

    Your Docker ID must be between 4 and 30 characters long, and can only contain numbers and lowercase letters. Once you create your Docker ID you can't reuse it in the future if you deactivate this account.

6. Select **Sign up**.

## Sign in

Once you register and verify your Docker ID email address, you can sign in to [your Docker account](https://login.docker.com/u/login/). You can sign in with your email address (or username) and password. Or, you can sign in with your social provider. See [Sign in with your social provider](#sign-in-with-your-social-provider).

You can also sign in through the CLI using the `docker login` command. For more information, see [`docker login`](../reference/cli/docker/login.md).

> **Warning**
>
> When you use the `docker login` command, your credentials are
stored in your home directory in `.docker/config.json`. The password is base64-encoded in this file.
>
> We recommend using one of the [Docker credential helpers](https://github.com/docker/docker-credential-helpers) for secure storage of passwords. For extra security, you can also use a [personal access token](../security/for-developers/access-tokens.md) to log in instead, which is still encoded in this file (without a Docker credential helper) but doesn't allow admin actions (such as changing the password).
{ .warning }

### Sign in with your social provider

> **Important**
>
> To sign in with your social provider, make sure you verify your email address with your provider before you begin.
{ .important }

Optionally, you can sign in to an existing Docker account with your Google or GitHub account. If a Docker account exists with the same email address as the primary email for your social provider, your Docker account will automatically be linked to the social profile. This lets you sign in with your social provider.

If you try to sign in with your social provider and don't have a Docker account yet, a new account will be created for you. Follow the on-screen instructions to create a Docker ID using your social provider.

## Reset your password at sign in

To reset your password, enter your email address on the [Sign in](https://login.docker.com/u/login) page and continue to sign in. When prompted for your password, select **Forgot password?**.

## Troubleshooting

For support and troubleshooting information, see [Get support](../support.md).