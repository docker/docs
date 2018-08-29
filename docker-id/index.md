---
title: Sign up for a Docker ID
description: Register for a Docker ID for a Docker account
keywords: docker-id, account, hub, forums, success-center, support-center
---

Your free Docker ID grants you access to Docker Hub repositories and other
Docker services. Your Docker ID becomes repository namespace. All you need is
an email address.

Your Docker account also allows you to log in to services such as the Docker
Support Center, the Docker Forums, and the Docker Success portal.

## Register for a Docker ID

{% include register-for-docker-id.md %}

## Log in to Docker

Once you register and verify your Docker ID email address, you can log in
to Docker services through the web interface or the commandline.

At the commandline, use `docker login`. See the [CLI reference](/engine/reference/commandline/login.md).

```
$ docker login
Login with your Docker ID to push and pull images from Docker Hub. If you don't have a Docker ID, head over to https://hub.docker.com to create one.
Username: <docker-id>
Password:
WARNING! Your password will be stored unencrypted in /home/<user>/.docker/config.json.
Configure a credential helper to remove this warning. See
https://docs.docker.com/engine/reference/commandline/login/#credentials-store

Login Succeeded
```

> Docker login creds not secure
>
> When you use the `docker login` command, your credentials are
stored in your home directory in `.docker/config.json`. The password is base64
encoded in this file. If you require secure storage for this password, use the
[Docker credential helpers](https://github.com/docker/docker-credential-helpers).
{:.warning}
