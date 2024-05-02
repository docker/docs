---
description: Explore the Learning center and understand the benefits of signing in
  to Docker Desktop
keywords: Docker Dashboard, manage, containers, gui, dashboard, images, user manual,
  learning center, guide, sign in
title: Sign in to Docker Desktop
aliases:
- /desktop/linux/
- /desktop/linux/index/
- /desktop/mac/
- /desktop/mac/index/
- /desktop/windows/
- /desktop/windows/index/
- /docker-for-mac/
- /docker-for-mac/index/
- /docker-for-mac/osx/
- /docker-for-mac/started/
- /docker-for-windows/
- /docker-for-windows/index/
- /docker-for-windows/started/
- /mac/
- /mac/started/
- /mackit/
- /mackit/getting-started/
- /win/
- /windows/
- /windows/started/
- /winkit/
- /winkit/getting-started/
---

In large enterprises where admin access is restricted, administrators can [Configure registry.json to enforce sign-in](../security/for-admins/configure-sign-in.md). 

> **Tip**
>
> Explore [Docker's core subscriptions](https://www.docker.com/pricing/) to see what else Docker can offer you. 

## Benefits of signing in

- You can access your Docker Hub repositories directly from Docker Desktop.

- Authenticated users also get a higher pull rate limit compared to anonymous users. For example, if you are authenticated, you get 200 pulls per 6 hour period, compared to 100 pulls per 6 hour period per IP address for anonymous users. For more information, see [Download rate limit](../docker-hub/download-rate-limit.md).

- Improve your organization’s security posture for containerized development by taking advantage of [Hardened Desktop](hardened-desktop/index.md).

## Signing in with Docker Desktop

Docker Desktop supports multiple ways of signing in users:

- Using the **Sign in** option in the top-right corner of the Docker Dashboard. This is the <u>recommended</u> option. When clicking on the button, your browser will open and ask for your credentials. Once authentication is successful, the browser will automatically sign Docker Desktop in.

> **Note**
>
> Docker Desktop automatically signs you out after 90 days, or after 30 days of inactivity.

- Using the Docker CLI: `docker login`. Refer to [Docker CLI "login" documentation](https://docs.docker.com/reference/cli/docker/login/) for additional details. If you sign in against https://hub.docker.com, then Docker Desktop will automatically be signed in.

> **Note**
>
> Docker Hub provides 2 ways of authenticating with the CLI:
> - via a password
> - via a generated Access Token (available here: https://hub.docker.com/settings/security)

## FAQ

Q: When signing in with my Docker Hub credentials and Docker CLI (`docker login`), Docker Desktop is not automatically signed in.

A: Check the following file: `$HOME/.docker/config.json`. When using Docker Desktop, the Credential Store (`credsStore`) needs to be `desktop`.

Docker Desktop will automatically use your native credential manager (Apple macOS keychain, Microsoft Windows Credential Manager, D-Bus Secret Service or [pass](https://www.passwordstore.org))


## Signing in with Docker Desktop for Linux

Docker Desktop for Linux relies on [`pass`](https://www.passwordstore.org/) to store credentials in gpg2-encrypted files.
Before signing in to Docker Desktop with your [Docker ID](../docker-id/_index.md), you must initialize `pass`.
Docker Desktop displays a warning if you've not initialized `pass`.

You can initialize pass by using a gpg key. To generate a gpg key, run:

``` console
$ gpg --generate-key
``` 

The following is an example similar to what you see once you run the previous command:

```console {hl_lines=12}
...
GnuPG needs to construct a user ID to identify your key.

Real name: Molly
Email address: molly@example.com
You selected this USER-ID:
   "Molly <molly@example.com>"

Change (N)ame, (E)mail, or (O)kay/(Q)uit? O
...
pubrsa3072 2022-03-31 [SC] [expires: 2024-03-30]
 <generated gpg-id public key>
uid          Molly <molly@example.com>
subrsa3072  2022-03-31 [E] [expires: 2024-03-30]
```

To initialize `pass`, run the following command using the public key generated from the previous command:

```console
$ pass init <your_generated_gpg-id_public_key>
``` 
The following is an example similar to what you see once you run the previous command:

```console
mkdir: created directory '/home/molly/.password-store/'
Password store initialized for <generated_gpg-id_public_key>
```

Once you initialize `pass`, you can sign in and pull your private images.
When Docker CLI or Docker Desktop use credentials, a user prompt may pop up for the password you set during the gpg key generation.

```console
$ docker pull molly/privateimage
Using default tag: latest
latest: Pulling from molly/privateimage
3b9cc81c3203: Pull complete 
Digest: sha256:3c6b73ce467f04d4897d7a7439782721fd28ec9bf62ea2ad9e81a5fb7fb3ff96
Status: Downloaded newer image for molly/privateimage:latest
docker.io/molly/privateimage:latest
```

## What's next?

- [Explore Docker Desktop](use-desktop/index.md) and its features. 
- Change your Docker Desktop settings
- [Browse common FAQs](faqs/general.md)
