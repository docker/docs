---
description: How to uninstall Docker Compose
keywords: compose, orchestration, uninstall, uninstallation, docker, documentation

title: Uninstall Docker Compose
---

Uninstalling Docker Compose depends on the method you have used to install Docker Compose. On this page you can find specific instructions to uninstall Docker Compose.


### Uninstalling Docker Desktop

If you want to uninstall Compose and you have installed Docker Desktop, see, [Uninstall Docker Desktop](../../desktop/uninstall.md) follow the corresponding link bellow to get instructions on how to remove Docker Desktop.

> **Note**
>
> Unless you have other Docker instances installed on that specific environment, you would be removing Docker altogether by uninstalling the Desktop.

### Uninstalling the Docker Compose CLI plugin

To remove the Compose CLI plugin, run:

Ubuntu, Debian:

```console
$ sudo apt-get remove docker-compose-plugin
```
RPM-based distros:

```console
$ sudo yum remove docker-compose-plugin
```

#### Manually installed

If you used `curl` to install Compose CLI plugin, to uninstall it, run:

```console
$ rm $DOCKER_CONFIG/cli-plugins/docker-compose
```

#### Remove for all users

Or, if you have installed Compose for all users, run:

```console
$ rm /usr/local/lib/docker/cli-plugins/docker-compose
```

> Got a **Permission denied** error?
>
> If you get a **Permission denied** error using either of the above
> methods, you do not have the permissions allowing you to remove
> `docker-compose`. To force the removal, prepend `sudo` to either of the above instructions and run it again.

#### Inspect the location of the Compose CLI plugin

To check where Compose is installed, use:

```console
$ docker info --format '{{range .ClientInfo.Plugins}}{{if eq .Name "compose"}}{{.Path}}{{end}}{{end}}'
```