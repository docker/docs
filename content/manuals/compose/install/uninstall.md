---
description: How to uninstall Docker Compose
keywords: compose, orchestration, uninstall, uninstallation, docker, documentation
title: Uninstall Docker Compose
linkTitle: Uninstall 
---

How you uninstall Docker Compose depends on how it was installed. This guide covers uninstallation instructions for:

- Docker Compose installed via Docker Desktop
- Docker Compose installed as a CLI plugin

## Uninstalling Docker Compose with Docker Desktop

If you want to uninstall Docker Compose and you have installed Docker Desktop, see [Uninstall Docker Desktop](/manuals/desktop/uninstall.md).

> [!WARNING]
>
> Unless you have other Docker instances installed on that specific environment, uninstalling Docker Desktop removes all Docker components, including Docker Engine, Docker CLI, and Docker Compose.

## Uninstalling the Docker Compose CLI plugin

If you installed Docker Compose via a package manager, run:

On Ubuntu or Debian:

   ```console
   $ sudo apt-get remove docker-compose-plugin
   ```
On RPM-based distributions:

   ```console
   $ sudo yum remove docker-compose-plugin
   ```

### Manually installed

If you installed Docker Compose manually (using curl), remove it by deleting the binary:

   ```console
   $ rm $DOCKER_CONFIG/cli-plugins/docker-compose
   ```

### Remove for all users

If installed for all users, remove it from the system directory:

   ```console
   $ rm /usr/local/lib/docker/cli-plugins/docker-compose
   ```

> [!NOTE]
>
> If you get a **Permission denied** error using either of the previous
> methods, you do not have the permissions needed to remove
> Docker Compose. To force the removal, prepend `sudo` to either of the previous instructions and run it again.

### Inspect the location of the Compose CLI plugin

To check where Compose is installed, use:

```console
$ docker info --format '{{range .ClientInfo.Plugins}}{{if eq .Name "compose"}}{{.Path}}{{end}}{{end}}'
```
