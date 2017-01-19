---
datafolder: engine-cli
datafile: docker_container_diff
title: docker container diff
---

<!--
Sorry, but the contents of this page are automatically generated from
Docker's source code. If you want to suggest a change to the text that appears
here, you'll need to find the string by searching this repo:

https://www.github.com/docker/docker
-->

{% include cli.md %}

## Examples

Inspect the changes to an `nginx` container:

```bash
$ docker diff 1fdfd1f54c1b

C /dev
C /dev/console
C /dev/core
C /dev/stdout
C /dev/fd
C /dev/ptmx
C /dev/stderr
C /dev/stdin
C /run
A /run/nginx.pid
C /var/lib/nginx/tmp
A /var/lib/nginx/tmp/client_body
A /var/lib/nginx/tmp/fastcgi
A /var/lib/nginx/tmp/proxy
A /var/lib/nginx/tmp/scgi
A /var/lib/nginx/tmp/uwsgi
C /var/log/nginx
A /var/log/nginx/access.log
A /var/log/nginx/error.log
```
