---
description: Getting Started
keywords: linux, tutorial, run, docker, local, machine
title: Docker Desktop for Linux user manual
---

The Docker Desktop for Linux user manual provides information on how to manage your credentials.

For information about Docker Desktop download, system requirements, and installation instructions, see [Install Docker Desktop](../install/linux-install.md).

## Credentials management

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

Once `pass` is initialized, we can sign in on the Docker Dashboard and pull our private images.
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

## Docker Hub

Select **Sign in / Create Docker ID** from the Docker Desktop menu to access your [Docker](https://hub.docker.com/){: target="_blank" rel="noopener" class="_" } account. Once logged in, you can access your Docker Hub repositories and organizations directly from the Docker Desktop menu.

For more information, refer to the following [Docker Hub topics](../../docker-hub/index.md){:target="_blank"
class="_"}:

* [Organizations and Teams in Docker Hub](../../docker-hub/orgs.md){:target="_blank" rel="noopener" class="_"}
* [Builds](../../docker-hub/builds/index.md){:target="_blank" rel="noopener" class="_"}

## Pause/Resume

You can pause your Docker Desktop session when you are not actively using it and save CPU resources on your machine. When you pause Docker Desktop, the Linux VM running Docker Engine is paused, the current state of all your containers are saved in memory, and all processes are frozen. This reduces the CPU usage and helps you retain a longer battery life on your laptop. You can resume Docker Desktop when you want by clicking the Resume option.

From the Docker menu, select ![whale menu](images/whale-x.png){: .inline} > **Pause** to pause Docker Desktop.

![Docker context menu](images/menu/prefs.png){:width="250px"}

Docker Desktop now displays the paused status on the Docker menu and on the  **Containers**, **Images**, **Volumes**, and **Dev Environment** screens on the Docker Dashboard. You can still access the **Settings** and the **Troubleshoot** menu from the Dashboard when you've paused Docker Desktop.

Select ![whale menu](images/whale-x.png){: .inline} > **Resume** to resume Docker Desktop.

> **Note**
>
> When Docker Desktop is paused, running any commands in the Docker CLI will automatically resume Docker Desktop.

## Where to go next

* Try out the walkthrough at [Get Started](/get-started/){: target="_blank"
  class="_"}.

* Dig in deeper with [Docker Labs](https://github.com/docker/labs/) example
  walkthroughs and source code.

* For a summary of Docker command line interface (CLI) commands, see
  [Docker CLI Reference Guide](../../engine/api/index.md){: target="_blank" rel="noopener" class="_"}.
