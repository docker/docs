---
description: How to uninstall Docker Compose
keywords: compose, orchestration, uninstall, uninstallation, docker, documentation
title: Uninstall Docker Compose
---

Uninstalling Docker Compose depends on the method you have used to install Docker Compose. On this page you can find specific instructions to uninstall Docker Compose.


## Uninstalling Docker Desktop

If you want to uninstall Docker Compose and you have installed Docker Desktop, see [Uninstall Docker Desktop](/manuals/desktop/uninstall.md).

> [!NOTE]
>
> Unless you have other Docker instances installed on that specific environment, you would be removing Docker altogether by uninstalling Docker Desktop.

## Uninstalling the Docker Compose CLI plugin

To remove the Docker Compose CLI plugin, run:

Ubuntu, Debian:

   ```console
   $ sudo apt-get remove docker-compose-plugin
   ```
RPM-based distributions:

   ```console
   $ sudo yum remove docker-compose-plugin
   ```

### Manually installed

If you used `curl` to install Docker Compose CLI plugin, to uninstall it, run:

   ```console
   $ rm $DOCKER_CONFIG/cli-plugins/docker-compose
   ```

### Remove for all users

Or, if you have installed Docker Compose for all users, run:

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
