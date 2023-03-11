---
title: "Building from Compose file"
keywords: build, buildx, bake, buildkit, compose
redirect_from:
- /build/customize/bake/compose-file/
---

## Specification

Bake uses the [compose-spec](../../compose/compose-file/index.md) to
parse a compose file and translate each service to a [target](file-definition.md#target).

```yaml
# docker-compose.yml
services:
  webapp-dev: 
    build: &build-dev
      dockerfile: Dockerfile.webapp
      tags:
        - docker.io/username/webapp:latest
      cache_from:
        - docker.io/username/webapp:cache
      cache_to:
        - docker.io/username/webapp:cache

  webapp-release:
    build:
      <<: *build-dev
      x-bake:
        platforms:
          - linux/amd64
          - linux/arm64

  db:
    image: docker.io/username/db
    build:
      dockerfile: Dockerfile.db
```

```console
$ docker buildx bake --print
```

```json
{
  "group": {
    "default": {
      "targets": [
        "db",
        "webapp-dev",
        "webapp-release"
      ]
    }
  },
  "target": {
    "db": {
      "context": ".",
      "dockerfile": "Dockerfile.db",
      "tags": [
        "docker.io/username/db"
      ]
    },
    "webapp-dev": {
      "context": ".",
      "dockerfile": "Dockerfile.webapp",
      "tags": [
        "docker.io/username/webapp:latest"
      ],
      "cache-from": [
        "docker.io/username/webapp:cache"
      ],
      "cache-to": [
        "docker.io/username/webapp:cache"
      ]
    },
    "webapp-release": {
      "context": ".",
      "dockerfile": "Dockerfile.webapp",
      "tags": [
        "docker.io/username/webapp:latest"
      ],
      "cache-from": [
        "docker.io/username/webapp:cache"
      ],
      "cache-to": [
        "docker.io/username/webapp:cache"
      ],
      "platforms": [
        "linux/amd64",
        "linux/arm64"
      ]
    }
  }
}
```

Unlike the [HCL format](file-definition.md#hcl-definition), there are some
limitations with the compose format:

* Specifying variables or global scope attributes is not yet supported
* `inherits` service field is not supported, but you can use [YAML anchors](https://docs.docker.com/compose/compose-file/#fragments){:target="blank" rel="noopener" class=""}
  to reference other services like the example above

## `.env` file

You can declare default environment variables in an environment file named
`.env`. This file will be loaded from the current working directory,
where the command is executed and applied to compose definitions passed
with `-f`.

```yaml
# docker-compose.yml
services:
  webapp:
    image: docker.io/username/webapp:${TAG:-v1.0.0}
    build:
      dockerfile: Dockerfile
```

```
# .env
TAG=v1.1.0
```

```console
$ docker buildx bake --print
```

```json
{
  "group": {
    "default": {
      "targets": [
        "webapp"
      ]
    }
  },
  "target": {
    "webapp": {
      "context": ".",
      "dockerfile": "Dockerfile",
      "tags": [
        "docker.io/username/webapp:v1.1.0"
      ]
    }
  }
}
```

> **Note**
>
> System environment variables take precedence over environment variables
> in `.env` file.

## Extension field with `x-bake`

Even if some fields are not (yet) available in the compose specification, you
can use the [special extension](../../compose/compose-file/index.md#extension)
field `x-bake` in your compose file to evaluate extra fields:

```yaml
# docker-compose.yml
services:
  addon:
    image: ct-addon:bar
    build:
      context: .
      dockerfile: ./Dockerfile
      args:
        CT_ECR: foo
        CT_TAG: bar
      x-bake:
        tags:
          - ct-addon:foo
          - ct-addon:alp
        platforms:
          - linux/amd64
          - linux/arm64
        cache-from:
          - user/app:cache
          - type=local,src=path/to/cache
        cache-to:
          - type=local,dest=path/to/cache
        pull: true

  aws:
    image: ct-fake-aws:bar
    build:
      dockerfile: ./aws.Dockerfile
      args:
        CT_ECR: foo
        CT_TAG: bar
      x-bake:
        secret:
          - id=mysecret,src=./secret
          - id=mysecret2,src=./secret2
        platforms: linux/arm64
        output: type=docker
        no-cache: true
```

```console
$ docker buildx bake --print
```

```json
{
  "group": {
    "default": {
      "targets": [
        "aws",
        "addon"
      ]
    }
  },
  "target": {
    "addon": {
      "context": ".",
      "dockerfile": "./Dockerfile",
      "args": {
        "CT_ECR": "foo",
        "CT_TAG": "bar"
      },
      "tags": [
        "ct-addon:foo",
        "ct-addon:alp"
      ],
      "cache-from": [
        "user/app:cache",
        "type=local,src=path/to/cache"
      ],
      "cache-to": [
        "type=local,dest=path/to/cache"
      ],
      "platforms": [
        "linux/amd64",
        "linux/arm64"
      ],
      "pull": true
    },
    "aws": {
      "context": ".",
      "dockerfile": "./aws.Dockerfile",
      "args": {
        "CT_ECR": "foo",
        "CT_TAG": "bar"
      },
      "tags": [
        "ct-fake-aws:bar"
      ],
      "secret": [
        "id=mysecret,src=./secret",
        "id=mysecret2,src=./secret2"
      ],
      "platforms": [
        "linux/arm64"
      ],
      "output": [
        "type=docker"
      ],
      "no-cache": true
    }
  }
}
```

Complete list of valid fields for `x-bake`:

* `cache-from`
* `cache-to`
* `contexts`
* `no-cache`
* `no-cache-filter`
* `output`
* `platforms`
* `pull`
* `secret`
* `ssh`
* `tags`
