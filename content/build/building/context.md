---
title: Build context
description: Learn how to use the build context to access files from your Dockerfile
keywords: build, buildx, buildkit, context, git, tarball, stdin
---

The [`docker build`](../../engine/reference/commandline/build.md) or
[`docker buildx build`](../../engine/reference/commandline/buildx_build.md)
commands build Docker images from a [Dockerfile](../../engine/reference/builder.md)
and a "context".

The build context is the argument that you pass to the build command:

```console
$ docker build [OPTIONS] PATH | URL | -
                         ^^^^^^^^^^^^^^
```

## What is a build context?

You can pass any of the following inputs as the context for a build:

- The relative or absolute path to a local directory
- The address of a remote Git repository, tarball, or plain-text file
- A piped plain-text file or a tarball using standard input

### Filesystem contexts

When your build context is a local directory, a remote Git repository, or a tar file,
then that becomes the set of files that the builder can access during the build.
Build instructions can refer to any of the files and directories in the context.
For example, when you use a [`COPY` instruction](../../engine/reference/builder.md#copy),
the builder copies the file or directory from the build context, into the build container.
A filesystem build context is processed recursively:

- When you specify a local directory or a tarball, all subdirectories are included
- When you specify a remote Git repository, the repository and all submodules are included

### Text file contexts

When your build context is a plain-text file, the builder interprets the file
as a Dockerfile. With this approach, the builder doesn't receive a filesystem context.
For more information about building with a text file context, see [Text files](#text-files).

## Local directories and tarballs

The following example shows a build command that uses the current directory
(`.`) as a build context:

```console
$ docker build .
...
#16 [internal] load build context
#16 sha256:23ca2f94460dcbaf5b3c3edbaaa933281a4e0ea3d92fe295193e4df44dc68f85
#16 transferring context: 13.16MB 2.2s done
...
```

This makes files and directories in the current working directory available
to the builder. The builder loads the files that it needs from the build context,
when it needs them.

For example, given the following directory structure:

```
.
├── index.ts
├── src/
├── Dockerfile
├── package.json
└── package-lock.json
```

Dockerfile instructions can reference and include these files in the build
if you pass the directory as a context.

```dockerfile
# syntax=docker/dockerfile:1
FROM node:latest
WORKDIR /src
COPY package.json package-lock.json .
RUN npm ci
COPY index.ts src .
```

### `.dockerignore`

You can use a [`.dockerignore`](../../engine/reference/builder.md#dockerignore-file)
file to exclude some files or directories from being sent:

```gitignore
# .dockerignore
node_modules
bar
```

A `.dockerignore` file located at the root of build context is automatically
detected and used.

If you use multiple Dockerfiles, you can use different ignore-files for each
Dockerfile. You do so using a special naming convention for the ignore-files.
Place your ignore-file in the same directory as the Dockerfile, and prefix the
ignore-file with the name of the Dockerfile.

For example:

```console
.
├── index.ts
├── src/
├── docker
│   ├── build.Dockerfile
│   ├── build.Dockerfile.dockerignore
│   ├── lint.Dockerfile
│   ├── lint.Dockerfile.dockerignore
│   ├── test.Dockerfile
│   └── test.Dockerfile.dockerignore
├── package.json
└── package-lock.json
```

A Dockerfile-specific ignore-file takes precedence over the `.dockerignore`
file at the root of the build context if both exist.

## Git repositories

When you pass a URL pointing to the location of a Git repository as an argument
to `docker build`, the builder uses the repository as the build context.

The builder performs a shallow clone of the repository, downloading only
the HEAD commit, not the entire history.

The builder recursively clones the repository and any submodules it contains.

```console
$ docker build https://github.com/user/myrepo.git
```

By default, the builder clones the latest commit on the default branch of the
repository that you specify.

### URL fragments

You can append URL fragments to the Git repository address to make the builder
clone a specific branch, tag, and subdirectory of a repository.

The format of the URL fragment is `#ref:dir`, where:

- `ref` is the name of the branch, tag, or remote reference
- `dir` is a subdirectory inside the repository

For example, the following command uses the `container` branch,
and the `docker` subdirectory in that branch, as the build context:

```console
$ docker build https://github.com/user/myrepo.git#container:docker
```

The following table represents all the valid suffixes with their build
contexts:

| Build Syntax Suffix            | Commit Used                   | Build Context Used |
| ------------------------------ | ----------------------------- | ------------------ |
| `myrepo.git`                   | `refs/heads/<default branch>` | `/`                |
| `myrepo.git#mytag`             | `refs/tags/mytag`             | `/`                |
| `myrepo.git#mybranch`          | `refs/heads/mybranch`         | `/`                |
| `myrepo.git#pull/42/head`      | `refs/pull/42/head`           | `/`                |
| `myrepo.git#:myfolder`         | `refs/heads/<default branch>` | `/myfolder`        |
| `myrepo.git#master:myfolder`   | `refs/heads/master`           | `/myfolder`        |
| `myrepo.git#mytag:myfolder`    | `refs/tags/mytag`             | `/myfolder`        |
| `myrepo.git#mybranch:myfolder` | `refs/heads/mybranch`         | `/myfolder`        |

### Keep `.git` directory

By default, BuildKit doesn't keep the `.git` directory when using Git contexts.
You can configure BuildKit to keep the directory by setting the
[`BUILDKIT_CONTEXT_KEEP_GIT_DIR` build argument](../../engine/reference/builder.md#buildkit-built-in-build-args).
This can be useful to if you want to retrieve Git information during your build:

```dockerfile
# syntax=docker/dockerfile:1
FROM alpine
WORKDIR /src
RUN --mount=target=. \
  make REVISION=$(git rev-parse HEAD) build
```

```console
$ docker build \
  --build-arg BUILDKIT_CONTEXT_KEEP_GIT_DIR=1
  https://github.com/user/myrepo.git#main
```

### Private repositories

When you specify a Git context that's also a private repository, the builder
needs you to provide the necessary authentication credentials. You can use
either SSH or token-based authentication.

Buildx automatically detects and uses SSH credentials if the Git context you
specify is an SSH or Git address. By default, this uses `$SSH_AUTH_SOCK`.
You can configure the SSH credentials to use with the
[`--ssh` flag](../../engine/reference/commandline/buildx_build.md#ssh).

```console
$ docker buildx build --ssh default git@github.com:user/private.git
```

If you want to use token-based authentication instead, you can pass the token
using the
[`--secret` flag](../../engine/reference/commandline/buildx_build.md#secret).

```console
$ GIT_AUTH_TOKEN=<token> docker buildx build \
  --secret id=GIT_AUTH_TOKEN \
  https://github.com/user/private.git
```

> **Note**
>
> Don't use `--build-arg` for secrets, except for
> [HTTP proxies](../../network/proxy.md#set-proxy-using-the-cli)

### Remote tarballs

If you pass the URL to a remote tarball, then the URL itself is sent to the builder.

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
Tarball contexts must be tar archives conforming to the standard `tar` Unix
format and can be compressed with any one of the `xz`, `bzip2`, `gzip` or
`identity` (no compression) formats.

## Pipes

When you pass a single dash `-` as the argument to the build command, you can
pipe a plain-text file or a tarball as the context:

```console
$ docker build - PIPE
```

For example:

```console
$ docker build - < Dockerfile
$ docker build - < archive.tar
$ docker build - <<EOF
FROM node:alpine
COPY . .
RUN npm ci
EOF
```

### Tarballs

When you pipe a tarball to the build command, the build uses the contents of
the tarball as a filesystem context.

For example, given the following project directory:

```
.
├── Dockerfile
├── Makefile
├── README.md
├── main.c
├── scripts
├── src
└── test.Dockerfile
```

You can create a tarball of the directory and pipe it to the build for use as
a context:

```console
$ tar czf foo.tar.gz *
$ docker build - < foo.tar.gz
```

The build resolves the Dockerfile from the tarball context. You can use the
`--file` flag to specify the name and location of the Dockerfile relative to
the root of the tarball. The following command builds using `test.Dockerfile`
in the tarball:

```console
$ docker build --file test.Dockerfile - < foo.tar.gz
```

### Text files

When you use a text file as the build context, the builder interprets the file
as a Dockerfile. Using a text file as context means that the build has no
filesystem context. This can be useful when your build doesn't require any
local files. This means there's no filesystem context when building.

You can pass the text file using a standard input stream, or by pointing at the
URL of a remote text file.

```console
$ docker build - < Dockerfile
```

With PowerShell on Windows, you can run:

```powershell
Get-Content Dockerfile | docker build -
```

To use a remote text file, pass the URL of the text file as the argument to the
build command:

```console
$ docker build https://raw.githubusercontent.com/dvdksn/clockbox/main/Dockerfile
```

Again, this means that the build has no filesystem context,
so Dockerfile commands such as `COPY` can't refer to local files:

```console
$ ls
main.c
$ docker build -<<< $'FROM scratch\nCOPY main.c .'
[+] Building 0.0s (4/4) FINISHED
 => [internal] load build definition from Dockerfile       0.0s
 => => transferring dockerfile: 64B                        0.0s
 => [internal] load .dockerignore                          0.0s
 => => transferring context: 2B                            0.0s
 => [internal] load build context                          0.0s
 => => transferring context: 2B                            0.0s
 => ERROR [1/1] COPY main.c .                              0.0s
------
 > [1/1] COPY main.c .:
------
Dockerfile:2
--------------------
   1 |     FROM scratch
   2 | >>> COPY main.c .
   3 |
--------------------
ERROR: failed to solve: failed to compute cache key: failed to calculate checksum of ref 7ab2bb61-0c28-432e-abf5-a4c3440bc6b6::4lgfpdf54n5uqxnv9v6ymg7ih: "/main.c": not found
```

#### Build using heredocs

The following example builds an image using a `Dockerfile` that is passed
through standard input using
[shell heredocs](https://en.wikipedia.org/wiki/Here_document){: target="_blank" rel="noopener"}:

```bash
docker build -t myimage:latest - <<EOF
FROM busybox
RUN echo "hello world"
EOF
```

This approach is useful when you want to quickly run a build command with a
Dockerfile that's short and concise.
