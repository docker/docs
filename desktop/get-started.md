---
description: Docker Dashboard
keywords: Docker Dashboard, manage, containers, gui, dashboard, images, user manual
title: Sign in and get started
redirect_from:
- /docker-for-windows/
- /docker-for-windows/index/
- /docker-for-windows/started/
- /engine/installation/windows/
- /installation/windows/
- /win/
- /windows/
- /windows/started/
- /winkit/
- /winkit/getting-started/
- /desktop/linux/index/
- /desktop/windows/index/
- /docker-for-mac/
- /docker-for-mac/index/
- /docker-for-mac/mutagen/
- /docker-for-mac/mutagen-caching/
- /docker-for-mac/osx/
- /docker-for-mac/started/
- /engine/installation/mac/
- /installation/mac/
- /mac/
- /mac/started/
- /mackit/
- /mackit/getting-started/
- /docker-for-mac/osxfs/
- /docker-for-mac/osxfs-caching/
- /desktop/mac/index/
---

## Quick Start Guide

  Once Docker Desktop is installed, the Quick Start Guide launches. It includes a simple exercise to build an example Docker image, run it as a container, push and save the image to Docker Hub.

   ![Docker Quick Start tutorial](images/docker-tutorial-linux.png)

To run the Quick Start Guide on demand, select ![whale menu](images/whale-x.png){: .inline} and then choose **Quick Start Guide**.

For a more detailed guide, see [Get started](../get-started/index.md)

## Sign in to Docker Desktop

We recommend that you authenticate using the **Sign in/Create ID** option in the top-right corner of Docker Desktop.

Once logged in, you can access your Docker Hub repositories directly from Docker Desktop.

Authenticated users get a higher pull rate limit compared to anonymous users. For example, if you are authenticated, you get 200 pulls per 6 hour period, compared to 100 pulls per 6 hour period per IP address for anonymous users. For more information, see [Download rate limit](https://docs.docker.com/docker-hub/download-rate-limit/).

In large enterprises where admin access is restricted, administrators can create a registry.json file and deploy it to the developers’ machines using a device management software as part of the Docker Desktop installation process. Enforcing developers to authenticate through Docker Desktop also allows administrators to set up guardrails using features such as [Image Access Management](https://docs.docker.com/docker-hub/image-access-management/) which allows team members to only have access to Trusted Content on Docker Hub, and pull only from the specified categories of images. For more information, see Configure registry.json to enforce sign in https://docs.docker.com/docker-hub/configure-sign-in/.

### Two-factor authentication

Docker Desktop enables you to sign in to Docker Hub using two-factor authentication. Two-factor authentication provides an extra layer of security when accessing your Docker Hub account.

You must enable two-factor authentication in Docker Hub before signing into your Docker Hub account through Docker Desktop. For instructions, see [Enable two-factor authentication for Docker Hub](/docker-hub/2fa/).

After two-factor authentication is enabled:

1. Go to the Docker Desktop menu and then select **Sign in / Create Docker ID**.

2. Enter your Docker ID and password and click **Sign in**.

3. After you have successfully signed in, Docker Desktop prompts you to enter the authentication code. Enter the six-digit code from your phone and then click **Verify**.

### Credentials management for Linux users

Docker Desktop relies on [`pass`](https://www.passwordstore.org/){: target="_blank" rel="noopener" class="_"} to store credentials in gpg2-encrypted files.
Before signing in to Docker Hub from the Docker Dashboard or the Docker menu, you must initialize `pass`.
Docker Desktop displays a warning if you've not initialized `pass`.

You can intialize pass by using a gpg key. To generate a gpg key, run:

``` console
$ gpg --generate-key
...
GnuPG needs to construct a user ID to identify your key.

Real name: Molly
Email address: molly@example.com
You selected this USER-ID:
    "Molly <molly@example.com>"

Change (N)ame, (E)mail, or (O)kay/(Q)uit? O
...
pub   rsa3072 2022-03-31 [SC] [expires: 2024-03-30]
      7865BA9185AFA2C26C5B505669FC4F36530097C2
uid                      Molly <molly@example.com>
sub   rsa3072 2022-03-31 [E] [expires: 2024-03-30]
```

To initialize `pass`, run:

```console
molly@ubuntu:~$ pass init 7865BA9185AFA2C26C5B505669FC4F36530097C2
mkdir: created directory '/home/molly/.password-store/'
Password store initialized for 7865BA9185AFA2C26C5B505669FC4F36530097C2
```

Once `pass` is initialized, you can sign in on the Docker Dashboard and pull your private images.
When credentials are used by the Docker CLI or Docker Desktop, a user prompt may pop up for the password you set during the gpg key generation.

```console
$ docker pull molly/privateimage
Using default tag: latest
latest: Pulling from molly/privateimage
3b9cc81c3203: Pull complete 
Digest: sha256:3c6b73ce467f04d4897d7a7439782721fd28ec9bf62ea2ad9e81a5fb7fb3ff96
Status: Downloaded newer image for molly/privateimage:latest
docker.io/molly/privateimage:latest
```

## Pause or resume Docker Desktop

You can pause your Docker Desktop session when you are not actively using it and save CPU resources on your machine. When you pause Docker Desktop, the Linux VM running Docker Engine is paused, the current state of all your containers are saved in memory, and all processes are frozen. This reduces the CPU usage and helps you retain a longer battery life on your laptop. You can resume Docker Desktop when you want by clicking the **Resume** option.

From the Docker menu, select![whale menu](images/whale-x.png){: .inline} and then **Pause** to pause Docker Desktop.

Docker Desktop displays the paused status on the Docker menu and on the  **Containers**, **Images**, **Volumes**, and **Dev Environment** screens in Docker Dashboard. You can still access the **Preferences** (or **Settings** if you are a Windows user) and the **Troubleshoot** menu from the Dashboard when you've paused Docker Desktop.

Select ![whale menu](images/whale-x.png){: .inline} then **Resume** to resume Docker Desktop.

> **Note**
>
> When Docker Desktop is paused, running any commands in the Docker CLI will automatically resume Docker Desktop.