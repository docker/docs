---
title: Remote Bake file definition
description: Build with Bake using a remote file definition using Git or HTTP
keywords: build, buildx, bake, file, remote, git, http
---

You can also build Bake files directly from a remote Git repository or HTTPS URL:

```console
$ docker buildx bake "https://github.com/docker/cli.git#v20.10.11" --print
#1 [internal] load git source https://github.com/docker/cli.git#v20.10.11
#1 0.745 e8f1871b077b64bcb4a13334b7146492773769f7       refs/tags/v20.10.11
#1 2.022 From https://github.com/docker/cli
#1 2.022  * [new tag]         v20.10.11  -> v20.10.11
#1 DONE 2.9s
```

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

As you can see the context is fixed to `https://github.com/docker/cli.git` even if
[no context is actually defined](https://github.com/docker/cli/blob/2776a6d694f988c0c1df61cad4bfac0f54e481c8/docker-bake.hcl#L17-L26)
in the definition.

If you want to access the main context for bake command from a bake file
that has been imported remotely, you can use the [`BAKE_CMD_CONTEXT` built-in var](reference.md#built-in-variables).

```console
$ curl -s https://raw.githubusercontent.com/tonistiigi/buildx/remote-test/docker-bake.hcl
```

```hcl
target "default" {
  context = BAKE_CMD_CONTEXT
  dockerfile-inline = <<EOT
FROM alpine
WORKDIR /src
COPY . .
RUN ls -l && stop
EOT
}
```

```console
$ docker buildx bake "https://github.com/tonistiigi/buildx.git#remote-test" --print
```

```json
{
  "target": {
    "default": {
      "context": ".",
      "dockerfile": "Dockerfile",
      "dockerfile-inline": "FROM alpine\nWORKDIR /src\nCOPY . .\nRUN ls -l \u0026\u0026 stop\n"
    }
  }
}
```

```console
$ touch foo bar
$ docker buildx bake "https://github.com/tonistiigi/buildx.git#remote-test"
```

```text
...
 > [4/4] RUN ls -l && stop:
#8 0.101 total 0
#8 0.102 -rw-r--r--    1 root     root             0 Jul 27 18:47 bar
#8 0.102 -rw-r--r--    1 root     root             0 Jul 27 18:47 foo
#8 0.102 /bin/sh: stop: not found
```

```console
$ docker buildx bake "https://github.com/tonistiigi/buildx.git#remote-test" "https://github.com/docker/cli.git#v20.10.11" --print
#1 [internal] load git source https://github.com/tonistiigi/buildx.git#remote-test
#1 0.429 577303add004dd7efeb13434d69ea030d35f7888       refs/heads/remote-test
#1 CACHED
```

```json
{
  "target": {
    "default": {
      "context": "https://github.com/docker/cli.git#v20.10.11",
      "dockerfile": "Dockerfile",
      "dockerfile-inline": "FROM alpine\nWORKDIR /src\nCOPY . .\nRUN ls -l \u0026\u0026 stop\n"
    }
  }
}
```

```console
$ docker buildx bake "https://github.com/tonistiigi/buildx.git#remote-test" "https://github.com/docker/cli.git#v20.10.11"
```

```text
...
 > [4/4] RUN ls -l && stop:
#8 0.136 drwxrwxrwx    5 root     root          4096 Jul 27 18:31 kubernetes
#8 0.136 drwxrwxrwx    3 root     root          4096 Jul 27 18:31 man
#8 0.136 drwxrwxrwx    2 root     root          4096 Jul 27 18:31 opts
#8 0.136 -rw-rw-rw-    1 root     root          1893 Jul 27 18:31 poule.yml
#8 0.136 drwxrwxrwx    7 root     root          4096 Jul 27 18:31 scripts
#8 0.136 drwxrwxrwx    3 root     root          4096 Jul 27 18:31 service
#8 0.136 drwxrwxrwx    2 root     root          4096 Jul 27 18:31 templates
#8 0.136 drwxrwxrwx   10 root     root          4096 Jul 27 18:31 vendor
#8 0.136 -rwxrwxrwx    1 root     root          9620 Jul 27 18:31 vendor.conf
#8 0.136 /bin/sh: stop: not found
```

## Remote definition with the --file flag

You can also specify the Bake definition to load from the remote repository,
using the `--file` or `-f` flag:

```console
docker buildx bake -f bake.hcl "https://github.com/crazy-max/buildx.git#remote-with-local"
```

```text
...
#4 [2/2] RUN echo "hello world"
#4 0.270 hello world
#4 DONE 0.3s
```

If you want to use a combination of local and remote definitions, you can
specify a local definition using the `cwd://` prefix with `-f`:

```hcl
# local.hcl
target "default" {
  args = {
    HELLO = "foo"
  }
}
```

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
        "type=cacheonly"
      ]
    }
  }
}
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
