---
title: Build context
description: Learn how to use the build context to access files from your Dockerfile
keywords: build, buildx, buildkit, context, git, tarball, stdin
---

The [`docker build`](../../engine/reference/commandline/build.md) or
[`docker buildx build`](../../engine/reference/commandline/buildx_build.md)
commands build Docker images from a [Dockerfile](../../engine/reference/builder.md)
and a "context".

A build's context is the set of files located at the `PATH` or `URL` specified
as the positional argument to the build command:

```console
$ docker build [OPTIONS] PATH | URL | -
                         ^^^^^^^^^^^^^^
```

The build process can refer to any of the files in the context. For example,
your build can use a [`COPY` instruction](../../engine/reference/builder.md#copy)
to reference a file in the context or a [`RUN --mount=type=bind` instruction](../../engine/reference/builder.md#run---mounttypebind)
for better performance with [BuildKit](../buildkit/index.md). The build context
is processed recursively. So, a `PATH` includes any subdirectories and the
`URL` includes the repository and its submodules.

## `PATH` context

This example shows a build command that uses the current directory (`.`) as a
build context:

```console
$ docker build .
...
#16 [internal] load build context
#16 sha256:23ca2f94460dcbaf5b3c3edbaaa933281a4e0ea3d92fe295193e4df44dc68f85
#16 transferring context: 13.16MB 2.2s done
...
```

With the following Dockerfile:

```dockerfile
# syntax=docker/dockerfile:1
FROM busybox
WORKDIR /src
COPY foo .
```

And this directory structure:

```
.
├── Dockerfile
├── bar
├── foo
└── node_modules
```

The legacy builder sends the entire directory to the daemon, including `bar`
and `node_modules` directories, even though the `Dockerfile` does not use
them. When using [BuildKit](../buildkit/index.md), the client only sends the
files required by the `COPY` instructions, in this case `foo`.

In some cases you may want to send the entire context:

```dockerfile
# syntax=docker/dockerfile:1
FROM busybox
WORKDIR /src
COPY . .
```

You can use a [`.dockerignore`](../../engine/reference/builder.md#dockerignore-file)
file to exclude some files or directories from being sent:

```gitignore
# .dockerignore
node_modules
bar
```

> **Warning**
>
> Avoid using your root directory, `/`, as the `PATH` for your build context,
> as it causes the build to transfer the entire contents of your hard drive to
> the daemon.
{:.warning}

## `URL` context

The `URL` parameter can refer to three kinds of resources:
* [Git repositories](#git-repositories)
* Pre-packaged [tarball contexts](#tarball-contexts)
* Plain [text files](#text-files)

### Git repositories

When the `URL` parameter points to the location of a Git repository, the
repository acts as the build context. The builder recursively pulls the
repository and its submodules. A shallow clone is performed and therefore pulls
down just the latest commits, not the entire history. A repository is first
pulled into a temporary directory on your host. After that succeeds, the
directory is sent to the daemon as the context. Local copy gives you the ability
to access private repositories using local user credentials, VPN's, and so forth.

> **Note**
>
> If the `URL` parameter contains a fragment the system will recursively clone
> the repository and its submodules using a `git clone --recursive` command.

Git URLs accept a context configuration parameter in the form of a URL fragment,
separated by a colon (`:`). The first part represents the reference that Git
will check out, and can be either a branch, a tag, or a remote reference. The
second part represents a subdirectory inside the repository that will be used
as a build context.

For example, run this command to use a directory called `docker` in the branch
`container`:

```console
$ docker build https://github.com/user/myrepo.git#container:docker
```

The following table represents all the valid suffixes with their build
contexts:

| Build Syntax Suffix            | Commit Used           | Build Context Used |
|--------------------------------|-----------------------|--------------------|
| `myrepo.git`                   | `refs/heads/master`   | `/`                |
| `myrepo.git#mytag`             | `refs/tags/mytag`     | `/`                |
| `myrepo.git#mybranch`          | `refs/heads/mybranch` | `/`                |
| `myrepo.git#pull/42/head`      | `refs/pull/42/head`   | `/`                |
| `myrepo.git#:myfolder`         | `refs/heads/master`   | `/myfolder`        |
| `myrepo.git#master:myfolder`   | `refs/heads/master`   | `/myfolder`        |
| `myrepo.git#mytag:myfolder`    | `refs/tags/mytag`     | `/myfolder`        |
| `myrepo.git#mybranch:myfolder` | `refs/heads/mybranch` | `/myfolder`        |

By default `.git` directory is not kept on Git checkouts. You can set the
[BuildKit built-in arg `BUILDKIT_CONTEXT_KEEP_GIT_DIR=1`](../../engine/reference/builder.md#buildkit-built-in-build-args)
to keep it. It can be useful to keep it around if you want to retrieve Git
information during your build:

```dockerfile
# syntax=docker/dockerfile:1
FROM alpine
WORKDIR /src
RUN --mount=target=. \
  make REVISION=$(git rev-parse HEAD) build
```

```console
$ docker build --build-arg BUILDKIT_CONTEXT_KEEP_GIT_DIR=1 https://github.com/user/myrepo.git#main
```

### Tarball contexts

If you pass a `URL` to a remote tarball, the `URL` itself is sent to the daemon:

```console
$ docker build http://server/context.tar.gz
#1 [internal] load remote build context
#1 DONE 0.2s

#2 copy /context /
#2 DONE 0.1s
...
```

The download operation will be performed on the host the daemon is running on,
which is not necessarily the same host from which the build command is being
issued. The daemon will fetch `context.tar.gz` and use it as the build context.
Tarball contexts must be tar archives conforming to the standard `tar` UNIX
format and can be compressed with any one of the `xz`, `bzip2`, `gzip` or
`identity` (no compression) formats.

### Text files

Instead of specifying a context, you can pass a single `Dockerfile` in the
`URL` or pipe the file in via `STDIN`. To pipe a `Dockerfile` from `STDIN`:

```console
$ docker build - < Dockerfile
```

With Powershell on Windows, you can run:

```powershell
Get-Content Dockerfile | docker build -
```

If you use `STDIN` or specify a `URL` pointing to a plain text file, the system
places the contents into a file called `Dockerfile`, and any `-f`, `--file`
option is ignored. In this scenario, there is no context.

The following example builds an image using a `Dockerfile` that is passed
through stdin. No files are sent as build context to the daemon.

```bash
docker build -t myimage:latest -<<EOF
FROM busybox
RUN echo "hello world"
EOF
```

Omitting the build context can be useful in situations where your `Dockerfile`
does not require files to be copied into the image, and improves the build-speed,
as no files are sent to the daemon.
