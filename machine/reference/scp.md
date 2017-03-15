---
description: Copy files among machines
keywords: machine, scp, subcommand
title: docker-machine scp
---

Copy files from your local host to a machine, from machine to machine, or from a
machine to your local host using `scp`.

The notation is `machinename:/path/to/files` for the arguments; in the host
machine's case, you don't have to specify the name, just the path.

Consider the following example:

```none
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

In the case of transferring files from machine to machine, they go through the
local host's filesystem first (using `scp`'s `-3` flag).

When transferring large files or updating directories with lots of files,
you can use the `-d` flag, which uses `rsync` to transfer deltas instead of
transferring all of the files.

When transferring directories and not just files, avoid rsync surprises
by using trailing slashes on both the source and destination. For example:

```none
$ mkdir -p bar
$ touch bar/baz
$ docker-machine scp -r -d bar/ dev:/home/docker/bar/
$ docker-machine ssh dev ls bar
baz
```
