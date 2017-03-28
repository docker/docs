---
description: Getting started with Docker
keywords: beginner, getting started, Docker, install
redirect_from:
- /mac/started/
title: Install Docker using a shell script
---

This installation procedure for users who don't want to use a package manager to
install Docker. The script works on a "best effort" basis to determine your
operating system and environment, and attempts to provide reasonable defaults.
The script may allow you to install Docker in environments that are not actually
supported configurations.

If you can use a package manager, you should use the
[recommended installation procedure for your operating system](/engine/installation/)
instead. Using a package manager ensures that you get upgrades when they are
available, and allows you to install a specific version of Docker, rather than
the very latest version.

## Prerequisites

- You need `sudo` access on Linux, or administrator access on Windows.
- You need `curl` or `wget` installed. These instructions use `curl`, but you
  can adapt them to use `wget`.

## Install Docker using the `install.sh` script

> **Warning**: Always examine shell scripts you download from the internet before
> running them.

1.  Open [https://get.docker.com](https://get.docker.com/) in your web browser
    so that you can examine the script before running it. This is important
    because the script will run with elevated privileges.

2.  Run the script, using `curl` to download it and piping it through `sh`:

    ```bash
    $ curl -fsSL https://get.docker.com/ | sh
    ```

    You are prompted for your `sudo` password. The script determines your
    operating system, downloads and installs Docker and its dependencies, starts
    Docker, and attempts to configure your operating system to start Docker
    automatically.

    > **Note**: Ubuntu or Debian users whose host is behind a filtering proxy
    > may experience failure of the `apt-key` step during Docker installation.
    > To work around this, use the following command to manually add the Docker
    > key:
    >
    >       $ curl -fsSL https://get.docker.com/gpg | sudo apt-key add -

  3.  On Ubuntu or Debian systems, the script attempts to start Docker and to
      configure the system to start Docker automatically. On RPM-based platforms,
      use the following command to start Docker:

      ```bash
      $ sudo systemctl start docker
      ```

      If you have an older system that does not have `systemctl`, use the
      `service` command instead:

      ```bash
      $ sudo service docker start
      ```

      To configure Docker to start automatically on RPM-based systems, see
      [Configure Docker to start on boot](/engine/installation/linux/linux-postinstall.md#configure-docker-to-start-on-boot).

  4.  If you installed using this mechanism, Docker will not be upgraded
      automatically when new versions are available. Instead, repeat this
      procedure to upgrade Docker.
