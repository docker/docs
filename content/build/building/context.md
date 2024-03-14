---
title: Build context
description: Learn how to use the build context to access files from your Dockerfile
keywords: build, buildx, buildkit, context, git, tarball, stdin
---

The `docker build` and `docker buildx build` commands build Docker images from
a [Dockerfile](../../reference/dockerfile.md) and a context.

## What is a build context?

The build context is the set of files that your build can access.
The positional argument that you pass to the build command specifies the
context that you want to use for the build:

```console
$ docker build [OPTIONS] PATH | URL | -
                         ^^^^^^^^^^^^^^
```

You can pass any of the following inputs as the context for a build:

- The relative or absolute path to a local directory
- A remote URL of a Git repository, tarball, or plain-text file
- A plain-text file or tarball piped to the `docker build` command through standard input

### Filesystem contexts

When your build context is a local directory, a remote Git repository, or a tar
file, then that becomes the set of files that the builder can access during the
build. Build instructions such as `COPY` and `ADD` can refer to any of the
files and directories in the context.

A filesystem build context is processed recursively:

- When you specify a local directory or a tarball, all subdirectories are included
- When you specify a remote Git repository, the repository and all submodules are included

For more information about the different types of filesystem contexts that you
can use with your builds, see:

- [Local files](#local-context)
- [Git repositories](#git-repositories)
- [Remote tarballs](#remote-tarballs)

### Text file contexts

When your build context is a plain-text file, the builder interprets the file
as a Dockerfile. With this approach, the build doesn't use a filesystem context.

For more information, see [empty build context](#empty-context).

## Local context

To use a local build context, you can specify a relative or absolute filepath
to the `docker build` command. The following example shows a build command that
uses the current directory (`.`) as a build context:

```console
$ docker build .
...
#16 [internal] load build context
#16 sha256:23ca2f94460dcbaf5b3c3edbaaa933281a4e0ea3d92fe295193e4df44dc68f85
#16 transferring context: 13.16MB 2.2s done
...
```

This makes files and directories in the current working directory available to
the builder. The builder loads the files it needs from the build context when
needed.

You can also use local tarballs as build context, by piping the tarball
contents to the `docker build` command. See [Tarballs](#local-tarballs).

### Local directories

Consider the following directory structure:

```text
.
├── index.ts
├── src/
├── Dockerfile
├── package.json
└── package-lock.json
```

Dockerfile instructions can reference and include these files in the build if
you pass this directory as a context.

```dockerfile
# syntax=docker/dockerfile:1
FROM node:latest
WORKDIR /src
COPY package.json package-lock.json .
RUN npm ci
COPY index.ts src .
```

```console
$ docker build .
```

### Local context with Dockerfile from stdin

Use the following syntax to build an image using files on your local
filesystem, while using a Dockerfile from stdin.

```console
$ docker build -f- <PATH>
```

The syntax uses the -f (or --file) option to specify the Dockerfile to use, and
it uses a hyphen (-) as filename to instruct Docker to read the Dockerfile from
stdin.

The following example uses the current directory (.) as the build context, and
builds an image using a Dockerfile passed through stdin using a here-document.

```bash
# create a directory to work in
mkdir example
cd example

# create an example file
touch somefile.txt

# build an image using the current directory as context
# and a Dockerfile passed through stdin
docker build -t myimage:latest -f- . <<EOF
FROM busybox
COPY somefile.txt ./
RUN cat /somefile.txt
EOF
```

### Local tarballs

When you pipe a tarball to the build command, the build uses the contents of
the tarball as a filesystem context.

For example, given the following project directory:

```text
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

## Remote context

You can specify the address of a remote Git repository, tarball, or plain-text
file as your build context.

- For Git repositories, the builder automatically clones the repository. See
  [Git repositories](#git-repositories).
- For tarballs, the builder downloads and extracts the contents of the tarball.
  See [Tarballs](#remote-tarballs).

If the remote tarball is a text file, the builder receives no [filesystem
context](#filesystem-contexts), and instead assumes that the remote
file is a Dockerfile. See [Empty build context](#empty-context).

### Git repositories

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

#### URL fragments

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

#### Keep `.git` directory

By default, BuildKit doesn't keep the `.git` directory when using Git contexts.
You can configure BuildKit to keep the directory by setting the
[`BUILDKIT_CONTEXT_KEEP_GIT_DIR` build argument](../../reference/dockerfile.md#buildkit-built-in-build-args).
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

#### Private repositories

When you specify a Git context that's also a private repository, the builder
needs you to provide the necessary authentication credentials. You can use
either SSH or token-based authentication.

Buildx automatically detects and uses SSH credentials if the Git context you
specify is an SSH or Git address. By default, this uses `$SSH_AUTH_SOCK`.
You can configure the SSH credentials to use with the
[`--ssh` flag](../../reference/cli/docker/buildx/build.md#ssh).

```console
$ docker buildx build --ssh default git@github.com:user/private.git
```

If you want to use token-based authentication instead, you can pass the token
using the
[`--secret` flag](../../reference/cli/docker/buildx/build.md#secret).

```console
$ GIT_AUTH_TOKEN=<token> docker buildx build \
  --secret id=GIT_AUTH_TOKEN \
  https://github.com/user/private.git
```

> **Note**
>
> Don't use `--build-arg` for secrets.

### Remote context with Dockerfile from stdin

Use the following syntax to build an image using files on your local
filesystem, while using a Dockerfile from stdin.

```console
$ docker build -f- <URL>
```

The syntax uses the -f (or --file) option to specify the Dockerfile to use, and
it uses a hyphen (-) as filename to instruct Docker to read the Dockerfile from
stdin.

This can be useful in situations where you want to build an image from a
repository that doesn't contain a Dockerfile. Or if you want to build with a
custom Dockerfile, without maintaining your own fork of the repository.

The following example builds an image using a Dockerfile from stdin, and adds
the `hello.c` file from the [hello-world](https://github.com/docker-library/hello-world)
repository on GitHub.

```bash
docker build -t myimage:latest -f- https://github.com/docker-library/hello-world.git <<EOF
FROM busybox
COPY hello.c ./
EOF
```

### Remote tarballs

If you pass the URL to a remote tarball, the URL itself is sent to the builder.

```console
$ docker build http://server/context.tar.gz
#1 [internal] load remote build context
#1 DONE 0.2s

#2 copy /context /
#2 DONE 0.1s
...
```

The download operation will be performed on the host where the BuildKit daemon
is running. Note that if you're using a remote Docker context or a remote
builder, that's not necessarily the same machine as where you issue the build
command. BuildKit fetches the `context.tar.gz` and uses it as the build
context. Tarball contexts must be tar archives conforming to the standard `tar`
Unix format and can be compressed with any one of the `xz`, `bzip2`, `gzip` or
`identity` (no compression) formats.

## Empty context

When you use a text file as the build context, the builder interprets the file
as a Dockerfile. Using a text file as context means that the build has no
filesystem context.

You can build with an empty build context when your Dockerfile doesn't depend
on any local files.

### How to build without a context

You can pass the text file using a standard input stream, or by pointing at the
URL of a remote text file.

{{< tabs >}}
{{< tab name="Unix pipe" >}}

```console
$ docker build - < Dockerfile
```

{{< /tab >}}
{{< tab name="PowerShell" >}}

```powershell
Get-Content Dockerfile | docker build -
```

{{< /tab >}}
{{< tab name="Heredocs" >}}

```bash
docker build -t myimage:latest - <<EOF
FROM busybox
RUN echo "hello world"
EOF
```

{{< /tab >}}
{{< tab name="Remote file" >}}

```console
$ docker build https://raw.githubusercontent.com/dvdksn/clockbox/main/Dockerfile
```

{{< /tab >}}
{{< /tabs >}}

When you build without a filesystem context, Dockerfile instructions such as
`COPY` can't refer to local files:

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

## .dockerignore files

You can use a `.dockerignore` file to exclude files or directories from the
build context.

```text
# .dockerignore
node_modules
bar
```

This helps avoid sending unwanted files and directories to the builder,
improving build speed, especially when using a remote builder.

### Filename and location

When you run a build command, the build client looks for a file named
`.dockerignore` in the root directory of the context. If this file exists, the
files and directories that match patterns in the files are removed from the
build context before it's sent to the builder.

If you use multiple Dockerfiles, you can use different ignore-files for each
Dockerfile. You do so using a special naming convention for the ignore-files.
Place your ignore-file in the same directory as the Dockerfile, and prefix the
ignore-file with the name of the Dockerfile, as shown in the following example.

```text
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

### Syntax

The `.dockerignore` file is a newline-separated list of patterns similar to the
file globs of Unix shells. For the purposes of matching, the root of the
context is considered to be both the working and the root directory. For
example, the patterns `/foo/bar` and `foo/bar` both exclude a file or directory
named `bar` in the `foo` subdirectory of `PATH` or in the root of the Git
repository located at `URL`. Neither excludes anything else.

If a line in `.dockerignore` file starts with `#` in column 1, then this line
is considered as a comment and is ignored before interpreted by the CLI.

If you're interested in learning the precise details of the `.dockerignore`
pattern matching logic, check out the
[moby/patternmatcher repository](https://github.com/moby/patternmatcher/tree/main/ignorefile)
on GitHub, which contains the source code.

#### Matching

The following code snippet shows an example `.dockerignore` file.

```text
# comment
*/temp*
*/*/temp*
temp?
```

This file causes the following build behavior:

| Rule        | Behavior                                                                                                                                                                                                      |
| :---------- | :------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| `# comment` | Ignored.                                                                                                                                                                                                      |
| `*/temp*`   | Exclude files and directories whose names start with `temp` in any immediate subdirectory of the root. For example, the plain file `/somedir/temporary.txt` is excluded, as is the directory `/somedir/temp`. |
| `*/*/temp*` | Exclude files and directories starting with `temp` from any subdirectory that is two levels below the root. For example, `/somedir/subdir/temporary.txt` is excluded.                                         |
| `temp?`     | Exclude files and directories in the root directory whose names are a one-character extension of `temp`. For example, `/tempa` and `/tempb` are excluded.                                                     |

Matching is done using Go's
[`filepath.Match` function](https://golang.org/pkg/path/filepath#Match) rules.
A preprocessing step uses Go's
[`filepath.Clean` function](https://golang.org/pkg/path/filepath/#Clean)
to trim whitespace and remove `.` and `..`.
Lines that are blank after preprocessing are ignored.

> **Note**
>
> For historical reasons, the pattern `.` is ignored.

Beyond Go's `filepath.Match` rules, Docker also supports a special wildcard
string `**` that matches any number of directories (including zero). For
example, `**/*.go` excludes all files that end with `.go` found anywhere in the
build context.

You can use the `.dockerignore` file to exclude the `Dockerfile` and
`.dockerignore` files. These files are still sent to the builder as they're
needed for running the build. But you can't copy the files into the image using
`ADD`, `COPY`, or bind mounts.

#### Negating matches

You can prepend lines with a `!` (exclamation mark) to make exceptions to
exclusions. The following is an example `.dockerignore` file that uses this
mechanism:

```text
*.md
!README.md
```

All markdown files right under the context directory _except_ `README.md` are
excluded from the context. Note that markdown files under subdirectories are
still included.

The placement of `!` exception rules influences the behavior: the last line of
the `.dockerignore` that matches a particular file determines whether it's
included or excluded. Consider the following example:

```text
*.md
!README*.md
README-secret.md
```

No markdown files are included in the context except README files other than
`README-secret.md`.

Now consider this example:

```text
*.md
README-secret.md
!README*.md
```

All of the README files are included. The middle line has no effect because
`!README*.md` matches `README-secret.md` and comes last.
