---
title: Build secrets
description: Manage credentials and other secrets securely
keywords: build, secrets, credentials, passwords, tokens
---

A build secret is any piece of sensitive information, such as a password or API
token, consumed as part of your application's build process.

Build arguments and environment variables are inappropriate for passing secrets
to your build, because they persist in the final image. Instead, should use
secret mounts or SSH mounts, which expose secrets to your builds securely.

## Secret mounts

Secret mounts expose secrets to the build containers as files. You [mount the
secrets to the `RUN`
instructions](../../engine/reference/builder.md#run---mounttypesecret) that
need to access them, similar to how you would define a bind mount or cache
mount.

```dockerfile
RUN --mount=type=secret,id=mytoken \
    TOKEN=$(cat /run/secrets/mytoken) ...
```

To pass a secret to a build, use the [`docker build --secret`
flag](../../engine/reference/commandline/buildx_build.md#secret), or the
equivalent options for [Bake](../bake/reference.md#targetsecret).

{{< tabs >}}
{{< tab name="CLI" >}}

```console
$ docker build --secret id=mytoken,src=$HOME/.aws/credentials .
```

{{< /tab >}}
{{< tab name="Bake" >}}

```hcl
variable "HOME" {
  default = null
}

target "default" {
  secret = [
    "id=mytoken,src=${HOME}/.aws/credentials"
  ]
}
```

{{< /tab >}}
{{< /tabs >}}

### Sources

The source of a secret can be either a
[file](../../engine/reference/commandline/buildx_build.md#file) or an
[environment variable](../../engine/reference/commandline/buildx_build.md#env).
When you use the CLI or Bake, the type can be detected automatically. You can
also specify it explicitly with `type=file` or `type=env`.

The following example mounts the environment variable `KUBECONFIG` to secret ID
`kube`.

```console
$ docker build --secret id=kube,env=KUBECONFIG .
```

The following example maps an environment variable directly to a secret ID.

```console
$ docker build --secret env=KUBECONFIG .
```

### Target

By default, secrets are mounted to `/run/secrets/<id>`. You can customize the
mount point in the build container using the `target` option in the Dockerfile.

The following example mounts the secret to a `/root/.aws/credentials` file in
the build container.

```console
$ docker build --secret id=aws,src=/root/.aws/credentials .
```

```dockerfile
RUN --mount=type=secret,id=aws,target=/root/.aws/credentials \
    aws s3 cp ...
```

## SSH mounts

If the credential you want to use in your build is an SSH agent socket or key,
you can use the SSH mount instead of a secret mount. Cloning private Git
repositories is a common use case for SSH mounts.

The following example clones a private GitHub repository using a [Dockerfile
SSH mount](../../engine/reference/builder.md#run---mounttypessh).

```dockerfile
# syntax=docker/dockerfile:1
FROM alpine
ADD git@github.com:me/myprivaterepo.git /src/
```

To pass an SSH socket the build, you use the [`docker build --ssh`
flag](../../engine/reference/commandline/buildx_build.md#ssh), or equivalent
options for [Bake](../bake/reference.md#targetssh).

```console
$ docker buildx build --ssh default .
```
