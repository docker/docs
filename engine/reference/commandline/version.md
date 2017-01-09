---
datafolder: engine-cli
datafile: docker_version
title: docker version
---

<!--
Sorry, but the contents of this page are automatically generated from
Docker's source code. If you want to suggest a change to the text that appears
here, you'll need to find the string by searching this repo:

https://www.github.com/docker/docker
-->

{% include cli.md %}

## Examples

### Display Docker version information

The default output:

```bash
$ docker version
	Client:
	 Version:      1.8.0
	 API version:  1.20
	 Go version:   go1.4.2
	 Git commit:   f5bae0a
	 Built:        Tue Jun 23 17:56:00 UTC 2015
	 OS/Arch:      linux/amd64

	Server:
	 Version:      1.8.0
	 API version:  1.20
	 Go version:   go1.4.2
	 Git commit:   f5bae0a
	 Built:        Tue Jun 23 17:56:00 UTC 2015
	 OS/Arch:      linux/amd64
```

Get server version:

```bash
{% raw %}
$ docker version --format '{{.Server.Version}}'

	1.8.0
{% endraw %}
```

Dump raw data:

To view all available fields, you can use the format `{% raw %}{{json .}}{% endraw %}`.

```bash
{% raw %}
$ docker version --format '{{json .}}'

{"Client":{"Version":"1.8.0","ApiVersion":"1.20","GitCommit":"f5bae0a","GoVersion":"go1.4.2","Os":"linux","Arch":"amd64","BuildTime":"Tue Jun 23 17:56:00 UTC 2015"},"ServerOK":true,"Server":{"Version":"1.8.0","ApiVersion":"1.20","GitCommit":"f5bae0a","GoVersion":"go1.4.2","Os":"linux","Arch":"amd64","KernelVersion":"3.13.2-gentoo","BuildTime":"Tue Jun 23 17:56:00 UTC 2015"}}
{% endraw %}
```
