---
description: How to write to and view a container's logs
keywords: docker, logging
title: View a container's logs
---

The `docker logs` command shows information logged by a running container. The
information that is logged and the format of the log depends almost entirely on
the container's endpoint command.

By default, `docker logs` shows the command's output just as it would appear if
you ran the command interactively in a terminal. UNIX and Linux commands
typically open three I/O streams when they run, called `STDIN`, `STDOUT`, and
`STDERR`. `STDIN` is the commmand's input stream, which may include input from
the keyboard or input from another command. `STDOUT` is usually a command's
normal output, and `STDERR` is typically used to output error messages. By
default, `docker logs` shows the command's `STDOUT` and `STDERR`. To read more
about I/O and Linux, see the
[Linux Documentation Project article on I/O redirection](http://www.tldp.org/LDP/abs/html/io-redirection.html)

In some cases, `docker logs` may not show useful information unless you take
additional steps.

- If you use a [logging driver](overview.md) which sends logs to a file, an
  external host, a database, or another logging back-end, `docker logs` may not
  show useful information.

- If your image runs a non-interactive process such as a web server or a
  database, that application may send its output to log files instead of `STDOUT`
  and `STDERR`.

In the first case, your logs are processed in other ways and you may choose not
to use `docker logs`. In the second case, the official `nginx` image shows one
workaround, and the official Apache `httpd` image shows another.

The official `nginx` image creates a symbolic link from
`/var/log/nginx/access.log` to `/dev/stdout`, and creates another symbolic link
from `/var/log/nginx/error.log` to `/dev/stderr`, overwriting the previous
devices in the process. See the
[Dockerfile](https://github.com/nginxinc/docker-nginx/blob/8921999083def7ba43a06fabd5f80e4406651353/mainline/jessie/Dockerfile#L21-L23).

The official `httpd` driver changes the `httpd` application's configuration to
write its normal output directly to `/proc/self/fd/1` (which is `STDOUT`) and
its errors to `/proc/self/fd/2` (which is `STDERR`). See the
[Dockerfile](https://github.com/docker-library/httpd/blob/b13054c7de5c74bbaa6d595dbe38969e6d4f860c/2.2/Dockerfile#L72-L75).

## Next steps

- Learn about using custom [logging drivers](overview.md).
- Learn about writing a [Dockerfile](/engine/reference/builder.md).
