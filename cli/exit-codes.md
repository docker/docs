---
title: Exit codes
description: Description of how Docker CLI uses exit codes
keywords: docker, cli, exit codes, status codes, run
---

Exit codes for the Docker CLI can represent either the success or failure of
invoking the `docker` command itself, or the status of a command run inside a
container.

- Exit code `0` represents a successful invocation.
- Exit code `1` represents a failed invocation.
- For the `docker run` command, exit codes `125-127` indicate an error invoking
  a command inside of a container

## `docker run`

The exit code from `docker run` gives information about why the container
failed to run or why it exited. When `docker run` exits with a non-zero code,
the exit codes follow the `chroot` standard:

- `125` if the error is with Docker daemon **_itself_**

  ```console
  $ docker run --foo busybox; echo $?

  flag provided but not defined: --foo
  See 'docker run --help'.
  125
  ```

- `126` if the **_contained command_** cannot be invoked

  ```console
  $ docker run busybox /etc; echo $?

  docker: Error response from daemon: Container command '/etc' could not be invoked.
  126
  ```

- `127` if the **_contained command_** cannot be found

  ```console
  $ docker run busybox foo; echo $?

  docker: Error response from daemon: Container command 'foo' not found or does not exist.
  127
  ```

If neither of the above scenarios apply, the exit code of the contained command
is used:

```console
$ docker run busybox /bin/sh -c 'exit 3'
$ echo $?
3
```

## `docker attach`

When you use the `docker attach` command to attach your I/O streams to a
running container, the exit code for the `docker attach` command is the same
as the process in the container.

```console
$ docker run --name test -dit alpine
275c44472aebd77c926d4527885bb09f2f6db21d878c75f0a1c212c03d3bcfab
$ docker attach test
/# exit 13
$ echo $?
13
$ docker ps -a --filter name=test
CONTAINER ID   IMAGE     COMMAND     CREATED              STATUS                       PORTS     NAMES
a2fe3fd886db   alpine    "/bin/sh"   About a minute ago   Exited (13) 40 seconds ago             test
```
