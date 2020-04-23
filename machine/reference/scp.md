---
description: Copy files among machines
keywords: machine, scp, subcommand
title: docker-machine scp
hide_from_sitemap: true
---

Copy files from your local host to a machine, from machine to machine, or from a
machine to your local host using `scp`.

The notation is `machinename:/path/to/files` for the arguments; in the host
machine's case, you don't need to specify the name, just the path.

## Example

Consider the following example:

```bash
$ cat foo.txt
cat: foo.txt: No such file or directory

$ docker-machine ssh dev pwd
/home/docker

$ docker-machine ssh dev 'echo A file created remotely! >foo.txt'
$ docker-machine scp dev:/home/docker/foo.txt .
foo.txt                                                           100%   28     0.0KB/s   00:00
$ cat foo.txt
A file created remotely!
```

Just like how `scp` has a `-r` flag for copying files recursively,
`docker-machine` has a `-r` flag for this feature.

In the case of transferring files from machine to machine,
they go through the local host's filesystem first (using `scp`'s `-3` flag).

When transferring large files or updating directories with lots of files,
you can use the `-d` flag, which uses `rsync` to transfer deltas instead of
transferring all of the files.

When transferring directories and not just files, avoid rsync surprises
by using trailing slashes on both the source and destination. For example:

```bash
$ mkdir -p bar
$ touch bar/baz
$ docker-machine scp -r -d bar/ dev:/home/docker/bar/
$ docker-machine ssh dev ls bar
baz
```

## Specifying file paths for remote deployments

When you copy files to a remote server with `docker-machine scp` for app
deployment, make sure `docker-compose` and the Docker daemon know how to find
them. Avoid using relative paths, but specify absolute paths in
[Compose files](../../compose/compose-file/index.md). It's best to specify absolute
paths both for the location on the Docker daemon and within the container.

For example, imagine you want to transfer your local directory
`/Users/<username>/webapp` to a remote machine and bind mount it into a
container on the remote host. If the remote user is `ubuntu`, use a command like
this:

```bash
$ docker-machine scp -r /Users/<username>/webapp MACHINE-NAME:/home/ubuntu/webapp
```

Then write a docker-compose file that bind mounts it in:

```none
version: "3.1"
services:
  webapp:
    image: alpine
    command: cat /app/root.php
    volumes:
    - "/home/ubuntu/webapp:/app"
```

And we can try it out like so:

```bash
$ eval $(docker-machine env MACHINE-NAME)
$ docker-compose run webapp
```
