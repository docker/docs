---
description: How to uninstall Docker Compose
keywords: compose, orchestration, uninstall, uninstallation, docker, documentation

title: Uninstall Docker Compose
---


### Uninstalling Docker Desktop

If you want to uninstall Compose and you have installed Docker Desktop, follow the corresponding link bellow to get instructions on how to remove Docker Desktop. 
> Please note that, unless you have other Docker instances installed in that specific environment, you would be removing Docker altogether by uninstalling the Desktop.

See Uninstall Docker Desktop for:
* [Mac](../../desktop/mac/install.md/#uninstall-docker-desktop){:target="_blank" rel="noopener" class="_"}
* [Windows](../../desktop/windows/install.md/#uninstall-docker-desktop){:target="_blank" rel="noopener" class="_"}
* [Linux](../../desktop/linux/install.md/#uninstall-docker-desktop){:target="_blank" rel="noopener" class="_"}


### Uninstalling Compose CLI plugin

To remove the Compose CLI plugin, run:

```console
$ sudo apt-get remove docker-compose-plugin
```
Or, if using a different distro, use the equivalent package manager to remove `docker-compose-plugin`. 

__Manually installed__

If you used `curl` to install Compose CLI plugin, to uninstall it run:

```console
$ rm $DOCKER_CONFIG/cli-plugins/docker-compose
```
    
or, if you have installed Compose for all users, run:  

```console
$ rm /usr/local/lib/docker/cli-plugins/docker-compose
```

You can also use:

{% raw %}	
```console
docker info --format '{{range .ClientInfo.Plugins}}{{if eq .Name "compose"}}{{.Path}}{{end}}{{end}}'
```
{% endraw %}

to inspect the location of the Compose CLI plugin.


> Got a "Permission denied" error?
>
> If you get a "Permission denied" error using either of the above
> methods, you do not have the permissions allowing you to remove
> `docker-compose`. To force the removal, prepend `sudo` to either of the above instructions and run it again.
