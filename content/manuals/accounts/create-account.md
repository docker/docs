---
title: Create a Docker account
linkTitle: Create an account
weight: 10
description: Learn how to register for a Docker ID and sign in to your account
keywords: accounts, docker ID, billing, paid plans, support, Hub, Store, Forums, knowledge
  base, beta access, email, activation, verification
aliases:
- /docker-hub/accounts/
- /docker-id/
---

You can create a free Docker account with your email address or by signing up
with your Google or GitHub account. After creating a unique Docker ID, you can
access all Docker products, including Docker Hub, Docker Desktop, and Docker Scout.

Your Docker ID becomes your username for hosted Docker services, and
[Docker forums](https://forums.docker.com/).

> [!TIP]
>
> Explore [Docker's subscriptions](https://www.docker.com/pricing/) to see what
else Docker can offer you.

## Create an account

You can sign up with an email address or use your Google or GitHub account.

### Sign up with your email

1. Go to the [Docker sign-up page](https://app.docker.com/signup/).
1. Enter a unique, valid email address.
1. Enter a username to use as your Docker ID. Once you create your Docker ID
you can't reuse it in the future if you deactivate this account.

    Your username:
    - Must be between 4 and 30 characters long
    - Can only contain numbers and lowercase letters

1. Enter a password that's at least 9 characters long.
1. Select **Sign Up**.
1. Open your email client. Docker sends a verification email to the
address you provided.
1. Verify your email address to complete the registration process.

> [!NOTE]
>
> You must verify your email address before you have full access to Docker's
features.

### Sign up with Google or GitHub

> [!IMPORTANT]
>
> To sign up with your social provider, you must verify your email address with
your provider before you begin.

1. Go to the [Docker sign-up page](https://app.docker.com/signup/).
1. Select your social provider, Google or GitHub.
1. Select the social account you want to link to your Docker account.
1. Select **Authorize Docker** to let Docker access your social account
information. You will be re-routed to the sign-up page.
1. Enter a username to use as your Docker ID.

    Your username:
    - Must be between 4 and 30 characters long
    - Can only contain numbers and lowercase letters
1. Select **Sign up**.

## Sign in to your account

You can sign in with your email, Google or GitHub account, or from
the Docker CLI.

### Sign in with email or Docker ID

1. Go to the [Docker sign in page](https://login.docker.com).
1. Enter your email address or Docker ID and select **Continue**.
1. Enter your password and select **Continue**.

To reset your password, see [Reset your password](#reset-your-password).

### Sign in with Google or GitHub

> [!IMPORTANT]
>
> Your Google or GitHub account must have a verified email address.

You can sign in using your Google or GitHub credentials. If your social
account uses the same email address as an existing Docker ID, the
accounts are automatically linked.

If no Docker ID exists, Docker creates a new account for you.

Docker doesn't currently support linking multiple sign-in methods
to the same Docker ID.

### Sign in using the CLI

Use the `docker login` command to authenticate from the command line. For
details, see [`docker login`](/reference/cli/docker/login/).

> [!WARNING]
>
> The `docker login` command stores credentials in your home directory under
> `.docker/config.json`. The password is base64-encoded.
>
> To improve security, use
> [Docker credential helpers](https://github.com/docker/docker-credential-helpers).
> For even stronger protection, use a [personal access token](../security/access-tokens.md)
> instead of a password. This is especially useful in CI/CD environments
> or when credential helpers aren't available.

## Reset your password

To reset your password:

1. Go to the [Docker sign in page](https://login.docker.com/).
1. Enter your email address.
1. When prompted for your password, select **Forgot password?**.

## Troubleshooting

If you have a paid Docker subscription,
[contact the Support team](https://hub.docker.com/support/contact/) for assistance.

All Docker users can seek troubleshooting information and support through the
following resources, where Docker or the community respond on a best effort
basis:
   - [Docker Community Forums](https://forums.docker.com/)
   - [Docker Community Slack](http://dockr.ly/comm-slack)
