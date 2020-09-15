---
description: Sign up for a Docker ID and log in
keywords: accounts, docker ID, billing, paid plans, support, Hub, Store, Forums, knowledge base, beta access
title: Docker ID accounts
redirect_from:
- /docker-cloud/dockerid/
- /docker-hub/accounts/
---

Your free Docker ID grants you access to Docker Hub repositories and some beta programs. All you need is an email address.

## Register for a Docker ID

Your Docker ID becomes your user namespace for hosted Docker services, and becomes your username on the Docker Forums. To create a new Docker ID:

1. Go to the [Docker Hub signup page](https://hub.docker.com/signup/).

2. Enter a username that is also your Docker ID.

    Your Docker ID must be between 4 and 30 characters long, and can only contain numbers and lowercase letters.

3. Enter a unique, valid email address.

4. Enter a password. Note that the password must be at least 9 characters.

5. Complete the Captcha verification and then then click **Sign up**.

   Docker sends a verification email to the address you provided.

6. Verify your email address to complete the registration process.

> **Note**: You cannot log in with your Docker ID until you verify your email address.

## Log in

Once you register and verify your Docker ID email address, you can log in
to [Docker Hub](https://hub.docker.com).

You can also log in through the CLI using the `docker login` command. For more information, see [`docker login`](../engine/reference/commandline/login.md).

> **Warning**:
> When you use the `docker login` command, your credentials are
stored in your home directory in `.docker/config.json`. The password is base64-encoded in this file.
>
> For extra security, you can use a [personal access token](../docker-hub/access-tokens.md) to log in instead, which is still encoded in this file but doesn't allow admin actions (such as changing the password). If you require secure storage for this password or personal access token, use the [Docker credential helpers](https://github.com/docker/docker-credential-helpers).
{:.warning}
