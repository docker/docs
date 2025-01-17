---
title: Create an account
weight: 10
description: Learn how to register for a Docker ID and sign in to your account
keywords: accounts, docker ID, billing, paid plans, support, Hub, Store, Forums, knowledge
  base, beta access, email, activation, verification
aliases:
- /docker-hub/accounts/
- /docker-id/
---

You can create a free Docker account with your email address or by signing up with your Google or GitHub account. Once you've created your account with a unique Docker ID, you can access all Docker products, including Docker Hub. With Docker Hub, you can access repositories and explore images that are available from the community and verified publishers.

Your Docker ID becomes your username for hosted Docker services, and [Docker forums](https://forums.docker.com/).

> [!TIP]
>
> Explore [Docker's subscriptions](https://www.docker.com/pricing/) to see what else Docker can offer you.

## Create a Docker ID

### Sign up with your email address

1. Go to the [Docker sign-up page](https://app.docker.com/signup/).

2. Enter a unique, valid email address.

3. Enter a username to use as your Docker ID. Once you create your Docker ID you can't reuse it in the future if you deactivate this account.

    Your username:
    - Must be between 4 and 30 characters long
    - Can only contain numbers and lowercase letters

4. Enter a password that's at least 9 characters long.

5. Select **Sign Up**.

6. Open your email client. Docker sends a verification email to the address you provided.

7. Verify your email address to complete the registration process.

> [!NOTE]
>
> You must verify your email address before you have full access to Docker's features.

### Sign up with Google or GitHub

> [!IMPORTANT]
>
> To sign up with your social provider, you must verify your email address with your provider before you begin.

1. Go to the [Docker sign-up page](https://app.docker.com/signup/).

2. Select your social provider, Google or GitHub.

3. Select the social account you want to link to your Docker account.

4. Select **Authorize Docker** to let Docker to access your social account information. You will be re-routed to the sign-up page.

5. Enter a username to use as your Docker ID.

    Your username:
    - Must be between 4 and 30 characters long
    - Can only contain numbers and lowercase letters

6. Select **Sign up**.

## Sign in

Once you register your Docker ID and verify your email address, you can sign in to [your Docker account](https://login.docker.com/u/login/). You can either:
- Sign in with your email address (or username) and password.
- Sign in with your social provider. For more information, see [Sign in with your social provider](#sign-in-with-your-social-provider).
- Sign in through the CLI using the `docker login` command. For more information, see [`docker login`](/reference/cli/docker/login.md).

> [!WARNING]
>
> When you use the `docker login` command, your credentials are
stored in your home directory in `.docker/config.json`. The password is base64-encoded in this file.
>
> We recommend using one of the [Docker credential helpers](https://github.com/docker/docker-credential-helpers) for secure storage of passwords. For extra security, you can also use a [personal access token](../security/for-developers/access-tokens.md) to sign in instead, which is still encoded in this file (without a Docker credential helper) but doesn't permit administrator actions (such as changing the password).

### Sign in with your social provider

> [!IMPORTANT]
>
> To sign in with your social provider, you must verify your email address with your provider before you begin.

You can also sign in to your Docker account with your Google or GitHub account. If a Docker account exists with the same email address as the primary email for your social provider, your Docker account will automatically be linked to the social profile. This lets you sign in with your social provider.

If you try to sign in with your social provider and don't have a Docker account yet, a new account will be created for you. Follow the on-screen instructions to create a Docker ID using your social provider.

## Reset your password at sign in

To reset your password, enter your email address on the [Sign in](https://login.docker.com/u/login) page and continue to sign in. When prompted for your password, select **Forgot password?**.

## Troubleshooting

If you have a paid Docker subscription, you can [contact the Support team](https://hub.docker.com/support/contact/) for assistance.

All Docker users can seek troubleshooting information and support through the following resources, where Docker or the community respond on a best effort basis:
   - [Docker Community Forums](https://forums.docker.com/)
   - [Docker Community Slack](http://dockr.ly/comm-slack)
