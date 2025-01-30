---
title: Remote Bake file definition
description: Build with Bake using a remote file definition using Git or HTTP
keywords: build, buildx, bake, file, remote, git, http
---

You can build Bake files directly from a remote Git repository or HTTPS URL:

```console
$ docker buildx bake "https://github.com/docker/cli.git#v20.10.11" --print
#1 [internal] load git source https://github.com/docker/cli.git#v20.10.11
#1 0.745 e8f1871b077b64bcb4a13334b7146492773769f7       refs/tags/v20.10.11
#1 2.022 From https://github.com/docker/cli
#1 2.022  * [new tag]         v20.10.11  -> v20.10.11
#1 DONE 2.9s
```

This fetches the Bake definition from the specified remote location and
executes the groups or targets defined in that file. If the remote Bake
definition doesn't specify a build context, the context is automatically set to
the Git remote. For example, [this case](https://github.com/docker/cli/blob/2776a6d694f988c0c1df61cad4bfac0f54e481c8/docker-bake.hcl#L17-L26)
uses `https://github.com/docker/cli.git`:

```json
{
  "group": {
    "default": {
      "targets": ["binary"]
    }
  },
  "target": {
    "binary": {
      "context": "https://github.com/docker/cli.git#v20.10.11",
      "dockerfile": "Dockerfile",
      "args": {
        "BASE_VARIANT": "alpine",
        "GO_STRIP": "",
        "VERSION": ""
      },
      "target": "binary",
      "platforms": ["local"],
      "output": ["build"]
    }
  }
}
```

## Use the local context with a remote definition

When building with a remote Bake definition, you may want to consume local
files relative to the directory where the Bake command is executed. You can
define contexts as relative to the command context using a `cwd://` prefix.

```hcl {title="https://github.com/dvdksn/buildx/blob/bake-remote-example/docker-bake.hcl"}
target "default" {
  context = "cwd://"
  dockerfile-inline = <<EOT
FROM alpine
WORKDIR /src
COPY . .
RUN ls -l && stop
EOT
}
```

```console
$ touch foo bar
$ docker buildx bake "https://github.com/dvdksn/buildx.git#bake-remote-example"
```

```text
...
 > [4/4] RUN ls -l && stop:
#8 0.101 total 0
#8 0.102 -rw-r--r--    1 root     root             0 Jul 27 18:47 bar
#8 0.102 -rw-r--r--    1 root     root             0 Jul 27 18:47 foo
#8 0.102 /bin/sh: stop: not found
```

You can append a path to the `cwd://` prefix if you want to use a specific
local directory as a context. Note that if you do specify a path, it must be
within the working directory where the command gets executed. If you use an
absolute path, or a relative path leading outside of the working directory,
Bake will throw an error.

### Local named contexts

You can also use the `cwd://` prefix to define local directories in the Bake
execution context as named contexts.

The following example defines the `docs` context as `./src/docs/content`,
relative to the current working directory where Bake is run as a named context.

```hcl {title=docker-bake.hcl}
target "default" {
  contexts = {
    docs = "cwd://src/docs/content"
  }
  dockerfile = "Dockerfile"
}
```

By contrast, if you omit the `cwd://` prefix, the path would be resolved
relative to the build context.

## Specify the Bake definition to use

When loading a Bake file from a remote Git repository, if the repository
contains more than one Bake file, you can specify which Bake definition to use
with the `--file` or `-f` flag:

```console
docker buildx bake -f bake.hcl "https://github.com/crazy-max/buildx.git#remote-with-local"
```

```text
...
#4 [2/2] RUN echo "hello world"
#4 0.270 hello world
#4 DONE 0.3s
```

## Combine local and remote Bake definitions

You can also combine remote definitions with local ones using the `cwd://`
prefix with `-f`.

Given the following local Bake definition in the current working directory:

```hcl
# local.hcl
target "default" {
  args = {
    HELLO = "foo"
  }
}
```

The following example uses `-f` to specify two Bake definitions:

- `-f bake.hcl`: this definition is loaded relative to the Git URL.
- `-f cwd://local.hcl`: this definition is loaded relative to the current
  working directory where the Bake command is executed.

```console
docker buildx bake -f bake.hcl -f cwd://local.hcl "https://github.com/crazy-max/buildx.git#remote-with-local" --print
```

```json
{
  "target": {
    "default": {
      "context": "https://github.com/crazy-max/buildx.git#remote-with-local",
      "dockerfile": "Dockerfile",
      "args": {
        "HELLO": "foo"
      },
      "target": "build",
      "output": [
        {
          "type": "cacheonly"
        }
      ]
    }
  }
}
```

One case where combining local and remote Bake definitions becomes necessary is
when you're building with a remote Bake definition in GitHub Actions and want
to use the [metadata-action](https://github.com/docker/metadata-action) to
generate tags, annotations, or labels. The metadata action generates a Bake
file available in the runner's local Bake execution context. To use both the
remote definition and the local "metadata-only" Bake file, specify both files
and use the `cwd://` prefix for the metadata Bake file:

```yml
      - name: Build
        uses: docker/bake-action@v6
        with:
          files: |
            ./docker-bake.hcl
            cwd://${{ steps.meta.outputs.bake-file }}
          targets: build
```

## Remote definition in a private repository

If you want to use a remote definition that lives in a private repository,
you may need to specify credentials for Bake to use when fetching the definition.

If you can authenticate to the private repository using the default `SSH_AUTH_SOCK`,
then you don't need to specify any additional authentication parameters for Bake.
Bake automatically uses your default agent socket.

For authentication using an HTTP token, or custom SSH agents,
use the following environment variables to configure Bake's authentication strategy:

- [`BUILDX_BAKE_GIT_AUTH_TOKEN`](../building/variables.md#buildx_bake_git_auth_token)
- [`BUILDX_BAKE_GIT_AUTH_HEADER`](../building/variables.md#buildx_bake_git_auth_header)
- [`BUILDX_BAKE_GIT_SSH`](../building/variables.md#buildx_bake_git_ssh)
