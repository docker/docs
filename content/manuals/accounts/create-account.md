---
title: Create a Docker account
linkTitle: Create
weight: 10
description: Create a Docker ID with email, Google, or GitHub, then sign in
keywords:
  create docker account, docker ID, sign up, sign in, email, Google, GitHub,
  verification, OTP, password reset, docker login
aliases:
  - /docker-hub/accounts/
  - /docker-id/
---

You can create a free Docker account with your email address or by signing up
with your Google or GitHub account. After creating a unique Docker ID, you can
access Docker products, including Docker Hub, Docker Desktop, and Docker Scout.

Your Docker ID becomes your username for hosted Docker services and
[Docker forums](https://forums.docker.com/).

> [!TIP]
>
> Explore
> [Docker's subscriptions](https://www.docker.com/pricing?ref=Docs&refAction=DocsCreateAccount)
> to see what else Docker can offer you.

## Create and verify your account

Signing up with an email address, Google, or GitHub requires verification
before you can sign in:

- If you sign up with Google or GitHub, verify your email address with that
  provider first.
- If you sign up with an email address, Docker sends a one-time password
  (OTP). Enter the code from the Docker email on the OTP page.

Docker blocks sign-in until you've verified your account.

{{< tabs >}}
{{< tab name="Email" >}}

1. Go to the [Docker sign-up page](https://app.docker.com/signup/) and enter a
   unique, valid email address.
1. Enter a username to use as your Docker ID. Once you create your Docker ID,
   you can't reuse it if you deactivate this account. Your username:
   - Must be between 4 and 30 characters long
   - Can only contain numbers and lowercase letters
1. Choose a password that's at least 9 characters long, then select
   **Sign up**.
1. Verify your email address when you receive the Docker OTP verification
   email. This completes the registration process.

{{< /tab >}}
{{< tab name="Google or GitHub" >}}

1. Go to the [Docker sign-up page](https://app.docker.com/signup/).
1. Select your social provider, Google or GitHub.
1. Select the social account you want to link to your Docker account.
1. Select **Authorize Docker** to let Docker access your social account
   information. Docker redirects you to the sign-up page.
1. Enter a username to use as your Docker ID. Your username:
   - Must be between 4 and 30 characters long
   - Can only contain numbers and lowercase letters
1. Select **Sign up**.

{{< /tab >}}
{{< /tabs >}}

## Sign in to your account

You can sign in with your email, Google or GitHub account, or from the Docker
CLI.

{{< tabs >}}
{{< tab name="Email or Docker ID" >}}

1. Go to the [Docker sign in page](https://login.docker.com/).
1. Enter your email address or Docker ID and select **Continue**.
1. Enter your password and select **Continue**.

To reset your password, see [Reset your password](#reset-your-password).

{{< /tab >}}
{{< tab name="Google or GitHub" >}}

You can sign in using your Google or GitHub credentials. If your social
account uses the same email address as an existing Docker ID, Docker links
that provider to the account.

If no Docker ID exists, Docker creates a new account for you.

You can't link both Google and GitHub to the same Docker ID.

{{< /tab >}}
{{< tab name="CLI" >}}

Use the `docker login` command to authenticate from the command line. For
details, see [`docker login`](/reference/cli/docker/login/).

> [!WARNING]
>
> The `docker login` command stores credentials in your home directory under
> `.docker/config.json`. The password is base64-encoded.
>
> To improve security, use
> [Docker credential helpers](https://github.com/docker/docker-credential-helpers).
> For even stronger protection, use a
> [personal access token](../security/access-tokens.md) instead of a password.
> This is especially useful in CI/CD environments or when credential helpers
> aren't available.

{{< /tab >}}
{{< /tabs >}}

## Reset your password

To reset your password:

1. Go to the [Docker sign in page](https://login.docker.com/).
1. Enter your email address.
1. When prompted for your password, select **Forgot password?**.

## Troubleshooting

If you have a paid Docker subscription,
[contact the Support team](https://hub.docker.com/support/contact/) for
assistance.

All Docker users can seek troubleshooting information and support through the
following resources, where Docker or the community respond on a best-effort
basis:

- [Docker Community Forums](https://forums.docker.com/)
- [Docker Community Slack](https://dockr.ly/comm-slack)
