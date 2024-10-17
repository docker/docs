---
title: Build secrets
linkTitle: Secrets
weight: 30
description: Manage credentials and other secrets securely
keywords: build, secrets, credentials, passwords, tokens, ssh, git, auth, http
tags: [Secrets]
---

A build secret is any piece of sensitive information, such as a password or API
token, consumed as part of your application's build process.

Build arguments and environment variables are inappropriate for passing secrets
to your build, because they persist in the final image. Instead, you should use
secret mounts or SSH mounts, which expose secrets to your builds securely.

## Types of build secrets

- [Secret mounts](#secret-mounts) are general-purpose mounts for passing
  secrets into your build. A secret mount takes a secret from the build client
  and makes it temporarily available inside the build container, for the
  duration of the build instruction. This is useful if, for example, your build
  needs to communicate with a private artifact server or API.
- [SSH mounts](#ssh-mounts) are special-purpose mounts for making SSH sockets
  or keys available inside builds. They're commonly used when you need to fetch
  private Git repositories in your builds.
- [Git authentication for remote contexts](#git-authentication-for-remote-contexts)
  is a set of pre-defined secrets for when you build with a remote Git context
  that's also a private repository. These secrets are "pre-flight" secrets:
  they are not consumed within your build instruction, but they're used to
  provide the builder with the necessary credentials to fetch the context.

## Using build secrets

For secret mounts and SSH mounts, using build secrets is a two-step process.
First you need to pass the secret into the `docker build` command, and then you
need to consume the secret in your Dockerfile.

To pass a secret to a build, use the [`docker build --secret`
flag](/reference/cli/docker/buildx/build.md#secret), or the
equivalent options for [Bake](../bake/reference.md#targetsecret).

{{< tabs >}}
{{< tab name="CLI" >}}

```console
$ docker build --secret id=aws,src=$HOME/.aws/credentials .
```

{{< /tab >}}
{{< tab name="Bake" >}}

```hcl
variable "HOME" {
  default = null
}

target "default" {
  secret = [
    "id=aws,src=${HOME}/.aws/credentials"
  ]
}
```

{{< /tab >}}
{{< /tabs >}}

To consume a secret in a build and make it accessible to the `RUN` instruction,
use the [`--mount=type=secret`](/reference/dockerfile.md#run---mounttypesecret)
flag in the Dockerfile.

```dockerfile
RUN --mount=type=secret,id=aws \
    AWS_SHARED_CREDENTIALS_FILE=/run/secrets/aws \
    aws s3 cp ...
```

## Secret mounts

Secret mounts expose secrets to the build containers, as files or environment
variables. You can use secret mounts to pass sensitive information to your
builds, such as API tokens, passwords, or SSH keys.

### Sources

The source of a secret can be either a
[file](/reference/cli/docker/buildx/build.md#file) or an
[environment variable](/reference/cli/docker/buildx/build.md#env).
When you use the CLI or Bake, the type can be detected automatically. You can
also specify it explicitly with `type=file` or `type=env`.

The following example mounts the environment variable `KUBECONFIG` to secret ID `kube`,
as a file in the build container at `/run/secrets/kube`.

```console
$ docker build --secret id=kube,env=KUBECONFIG .
```

When you use secrets from environment variables, you can omit the `env` parameter
to bind the secret to a file with the same name as the variable.
In the following example, the value of the `API_TOKEN` variable
is mounted to `/run/secrets/API_TOKEN` in the build container.

```console
$ docker build --secret id=API_TOKEN .
```

### Target

When consuming a secret in a Dockerfile, the secret is mounted to a file by
default. The default file path of the secret, inside the build container, is
`/run/secrets/<id>`. You can customize how the secrets get mounted in the build
container using the `target` and `env` options for the `RUN --mount` flag in
the Dockerfile.

The following example takes secret id `aws` and mounts it to a file at
`/run/secrets/aws` in the build container.

```dockerfile
RUN --mount=type=secret,id=aws \
    AWS_SHARED_CREDENTIALS_FILE=/run/secrets/aws \
    aws s3 cp ...
```

To mount a secret as a file with a different name, use the `target` option in
the `--mount` flag.

```dockerfile
RUN --mount=type=secret,id=aws,target=/root/.aws/credentials \
    aws s3 cp ...
```

To mount a secret as an environment variable instead of a file, use the
`env` option in the `--mount` flag.

```dockerfile
RUN --mount=type=secret,id=aws-key-id,env=AWS_ACCESS_KEY_ID \
    --mount=type=secret,id=aws-secret-key,env=AWS_SECRET_ACCESS_KEY \
    --mount=type=secret,id=aws-session-token,env=AWS_SESSION_TOKEN \
    aws s3 cp ...
```

It's possible to use the `target` and `env` options together to mount a secret
as both a file and an environment variable.

## SSH mounts

If the credential you want to use in your build is an SSH agent socket or key,
you can use the SSH mount instead of a secret mount. Cloning private Git
repositories is a common use case for SSH mounts.

The following example clones a private GitHub repository using a [Dockerfile
SSH mount](/reference/dockerfile.md#run---mounttypessh).

```dockerfile
# syntax=docker/dockerfile:1
FROM alpine
ADD git@github.com:me/myprivaterepo.git /src/
```

To pass an SSH socket the build, you use the [`docker build --ssh`
flag](/reference/cli/docker/buildx/build.md#ssh), or equivalent
options for [Bake](../bake/reference.md#targetssh).

```console
$ docker buildx build --ssh default .
```

## Git authentication for remote contexts

BuildKit supports two pre-defined build secrets, `GIT_AUTH_TOKEN` and
`GIT_AUTH_HEADER`. Use them to specify HTTP authentication parameters when
building with remote, private Git repositories, including:

- Building with a private Git repository as build context
- Fetching private Git repositories in a build with `ADD`

For example, say you have a private GitLab project at
`https://gitlab.com/example/todo-app.git`, and you want to run a build using
that repository as the build context. An unauthenticated `docker build` command
fails because the builder isn't authorized to pull the repository:

```console
$ docker build https://gitlab.com/example/todo-app.git
[+] Building 0.4s (1/1) FINISHED
 => ERROR [internal] load git source https://gitlab.com/example/todo-app.git
------
 > [internal] load git source https://gitlab.com/example/todo-app.git:
0.313 fatal: could not read Username for 'https://gitlab.com': terminal prompts disabled
------
```

To authenticate the builder to the Git server, set the `GIT_AUTH_TOKEN`
environment variable to contain a valid GitLab access token, and pass it as a
secret to the build:

```console
$ GIT_AUTH_TOKEN=$(cat gitlab-token.txt) docker build \
  --secret id=GIT_AUTH_TOKEN \
  https://gitlab.com/example/todo-app.git
```

The `GIT_AUTH_TOKEN` also works with `ADD` to fetch private Git repositories as
part of your build:

```dockerfile
FROM alpine
ADD https://gitlab.com/example/todo-app.git /src
```

### HTTP authentication scheme

By default, Git authentication over HTTP uses the Bearer authentication scheme:

```http
Authorization: Bearer <GIT_AUTH_TOKEN>
```

If you need to use a Basic scheme, with a username and password, you can set
the `GIT_AUTH_HEADER` build secret:

```console
$ export GIT_AUTH_TOKEN=$(cat gitlab-token.txt)
$ export GIT_AUTH_HEADER=basic
$ docker build \
  --secret id=GIT_AUTH_TOKEN \
  --secret id=GIT_AUTH_HEADER \
  https://gitlab.com/example/todo-app.git
```

BuildKit currently only supports the Bearer and Basic schemes.

### Multiple hosts

You can set the `GIT_AUTH_TOKEN` and `GIT_AUTH_HEADER` secrets on a per-host
basis, which lets you use different authentication parameters for different
hostnames. To specify a hostname, append the hostname as a suffix to the secret
ID:

```console
$ export GITLAB_TOKEN=$(cat gitlab-token.txt)
$ export GERRIT_TOKEN=$(cat gerrit-username-password.txt)
$ export GERRIT_SCHEME=basic
$ docker build \
  --secret id=GIT_AUTH_TOKEN.gitlab.com,env=GITLAB_TOKEN \
  --secret id=GIT_AUTH_TOKEN.gerrit.internal.example,env=GERRIT_TOKEN \
  --secret id=GIT_AUTH_HEADER.gerrit.internal.example,env=GERRIT_SCHEME \
  https://gitlab.com/example/todo-app.git
```
