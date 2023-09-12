---
description: Learn how to register for a Docker ID and log in to your account
keywords: accounts, docker ID, billing, paid plans, support, Hub, Store, Forums, knowledge
  base, beta access, email, activation, verification
title: Create an account
aliases:
- /docker-cloud/dockerid/
- /docker-hub/accounts/
---

All you need is an email address to create a Docker account. Once you've created your account with a unique Docker ID, you can access Docker Hub repositories and explore images that are available from the community and verified publishers. 

Your Docker ID becomes your username for hosted Docker services, and [Docker forums](https://forums.docker.com/).

## Create a Docker ID

1. Go to the [Docker Hub signup page](https://hub.docker.com/signup/).

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

## Sign in

Once you register and verify your Docker ID email address, you can sign in to [Docker Hub](https://hub.docker.com).

You can also sign in through the CLI using the `docker login` command. For more information, see [`docker login`](../engine/reference/commandline/login.md).

> **Warning**
>
> When you use the `docker login` command, your credentials are
stored in your home directory in `.docker/config.json`. The password is base64-encoded in this file.
>
> We recommend using one of the [Docker credential helpers](https://github.com/docker/docker-credential-helpers) for secure storage of passwords. For extra security, you can also use a [personal access token](../docker-hub/access-tokens.md) to log in instead, which is still encoded in this file (without a Docker credential helper) but doesn't allow admin actions (such as changing the password).
{ .warning }

## Troubleshooting

If you run into trouble with your Docker ID account, know that we're here to help! If you want to retrieve or reset your password, [enter your email address](https://login.docker.com/u/login) for additional instructions.

You can use the [Docker forums](https://forums.docker.com/) to ask questions amongst other Docker community members, while our [hub-feedback GitHub repository](https://github.com/docker/hub-feedback) allows you to provide feedback on how we can better improve the experience with Docker Hub.

If you still need any help, [create a support ticket](https://hub.docker.com/support/contact/) and let us know how we can help you.