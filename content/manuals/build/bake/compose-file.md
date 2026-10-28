---
title: Building with Bake from a Compose file
description: Build your compose services with Bake
keywords: build, buildx, bake, buildkit, compose, yaml
aliases:
  - /build/customize/bake/compose-file/
  - /guides/compose-bake/
---

Bake supports the [Compose file format](/reference/compose-file/_index.md)
to parse a Compose file and translate each service to a [target](reference.md#target).

```yaml
# compose.yaml
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
      "targets": ["db", "webapp-dev", "webapp-release"]
    }
  },
  "target": {
    "db": {
      "context": ".",
      "dockerfile": "Dockerfile.db",
      "tags": ["docker.io/username/db"]
    },
    "webapp-dev": {
      "context": ".",
      "dockerfile": "Dockerfile.webapp",
      "tags": ["docker.io/username/webapp:latest"],
      "cache-from": [
        {
          "ref": "docker.io/username/webapp:cache",
          "type": "registry"
        }
      ],
      "cache-to": [
        {
          "ref": "docker.io/username/webapp:cache",
          "type": "registry"
        }
      ]
    },
    "webapp-release": {
      "context": ".",
      "dockerfile": "Dockerfile.webapp",
      "tags": ["docker.io/username/webapp:latest"],
      "cache-from": [
        {
          "ref": "docker.io/username/webapp:cache",
          "type": "registry"
        }
      ],
      "cache-to": [
        {
          "ref": "docker.io/username/webapp:cache",
          "type": "registry"
        }
      ],
      "platforms": ["linux/amd64", "linux/arm64"]
    }
  }
}
```

The compose format has some limitations compared to the HCL format:

- Specifying variables or global scope attributes is not supported
- `inherits` service field is not supported, but you can use [YAML anchors](/reference/compose-file/fragments.md)
  to reference other services, as demonstrated in the previous example with `&build-dev`.

## Extend Compose builds with HCL

You can keep development settings in Compose and layer release build settings in
`docker-bake.hcl`. This is useful when release builds need different stages or
platforms from local development builds. Compose can also build release images;
separating the settings is a workflow choice.

For example, this Compose file builds an application and a utility service that
seeds development data:

```yaml {title=compose.yaml}
services:
  app:
    image: docker.io/username/app:latest
    build:
      context: ./app
      target: dev
  seed:
    build:
      context: ./seed
```

Create `docker-bake.hcl` alongside `compose.yaml`:

```hcl {title=docker-bake.hcl}
group "default" {
  targets = ["app"]
}

target "_release" {
  platforms = ["linux/amd64", "linux/arm64"]
  attest = ["type=provenance,mode=max", "type=sbom"]
  annotations = ["org.opencontainers.image.authors=username"]
}

target "app" {
  inherits = ["_release"]
  target = "release"
}
```

Bake automatically loads and merges both files. The HCL file replaces the
`default` group's target list, so `docker buildx bake` builds only `app`.
The `seed` target remains available with `docker buildx bake seed`.
See [Default file lookup](overrides.md#default-file-lookup) for file precedence.

The `app` target keeps its context and image tag from Compose, overrides the
Dockerfile stage from `dev` to `release`, and inherits the shared platforms,
attestations, and annotations. The application's Dockerfile must define both
stages. Other release targets can inherit the same settings; see
[Common reusable targets](inheritance.md#common-reusable-targets).

Inspect the merged configuration or select a target to inspect:

```console
$ docker buildx bake --print
$ docker buildx bake --print app
```

Check the `group.default.targets` list to see what runs by default. Omitting
`seed` from that group doesn't delete its definition from the merged
configuration. Selecting `app` limits the printed targets to that selection
and its dependencies.

## `.env` file

You can declare default environment variables in an environment file named
`.env`. This file will be loaded from the current working directory,
where the command is executed and applied to compose definitions passed
with `-f`.

```yaml
# compose.yaml
services:
  webapp:
    image: docker.io/username/webapp:${TAG:-v1.0.0}
    build:
      dockerfile: Dockerfile
```

```sh
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
      "targets": ["webapp"]
    }
  },
  "target": {
    "webapp": {
      "context": ".",
      "dockerfile": "Dockerfile",
      "tags": ["docker.io/username/webapp:v1.1.0"]
    }
  }
}
```

> [!NOTE]
>
> System environment variables take precedence over environment variables
> in `.env` file.

## Extension field with `x-bake`

Where some fields are not available in the compose specification, you can use
the [special extension](/reference/compose-file/extension.md) field
`x-bake` in your compose file to evaluate extra fields:

```yaml
# compose.yaml
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
      "targets": ["addon", "aws"]
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
      "tags": ["ct-addon:foo", "ct-addon:alp"],
      "cache-from": [
        {
          "ref": "user/app:cache",
          "type": "registry"
        },
        {
          "src": "path/to/cache",
          "type": "local"
        }
      ],
      "cache-to": [
        {
          "dest": "path/to/cache",
          "type": "local"
        }
      ],
      "platforms": ["linux/amd64", "linux/arm64"],
      "pull": true
    },
    "aws": {
      "context": ".",
      "dockerfile": "./aws.Dockerfile",
      "args": {
        "CT_ECR": "foo",
        "CT_TAG": "bar"
      },
      "tags": ["ct-fake-aws:bar"],
      "secret": [
        {
          "id": "mysecret",
          "src": "./secret"
        },
        {
          "id": "mysecret2",
          "src": "./secret2"
        }
      ],
      "platforms": ["linux/arm64"],
      "output": [
        {
          "type": "docker"
        }
      ],
      "no-cache": true
    }
  }
}
```

Complete list of valid fields for `x-bake`:

- `cache-from`
- `cache-to`
- `contexts`
- `no-cache`
- `no-cache-filter`
- `output`
- `platforms`
- `pull`
- `secret`
- `ssh`
- `tags`
