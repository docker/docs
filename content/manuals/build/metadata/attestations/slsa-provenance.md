---
title: Provenance attestations
keywords: build, attestations, provenance, slsa, git, metadata
description: >
  Provenance build attestations describe how and where your image was built.
aliases:
  - /build/attestations/slsa-provenance/
---

The provenance attestations include facts about the build process, including
details such as:

- Build timestamps
- Build parameters and environment
- Version control metadata
- Source code details
- Materials (files, scripts) consumed during the build

Provenance attestations default to the
[SLSA provenance schema, version 1](https://slsa.dev/spec/v1.1/provenance#schema).
Use [the `version` parameter](#version) to request the older v0.2 schema for
tools that require it.

For more information about how BuildKit populates these provenance properties, refer to
[SLSA definitions](slsa-definitions.md).

## Create provenance attestations

To create a provenance attestation, pass the `--attest type=provenance` option
to the `docker buildx build` command:

```console
$ docker buildx build --tag <namespace>/<image>:<version> \
    --attest type=provenance,mode=[min,max],version=[v0.2,v1] .
```

Alternatively, you can use the shorthand `--provenance=true` option instead of `--attest type=provenance`.
To specify the `mode` or `version` parameters using the shorthand option, use:
`--provenance=mode=max,version=v1`.

For an example on how to add provenance attestations with GitHub Actions, see
[Add attestations with GitHub Actions](/manuals/build/ci/github-actions/attestations.md).

## Mode

You can use the `mode` parameter to define the level of detail to be included in
the provenance attestation. Supported values are `mode=min` (default) and
`mode=max`.

### Min

In `min` mode, the provenance attestations include a minimal set of information,
such as:

- Build timestamps
- The frontend used
- Build materials
- Source repository and revision
- Build platform
- Reproducibility

Values of build arguments, the identities of secrets, and rich layer metadata are
not included in `mode=min`. The `min`-level provenance is safe to use for all
builds, as it doesn't leak information from any part of the build environment.

The following JSON example shows the information included in a provenance
attestation created using `mode=min` and SLSA v1 for `FROM alpine:latest`:

```json
{
  "_type": "https://in-toto.io/Statement/v1",
  "predicateType": "https://slsa.dev/provenance/v1",
  "subject": [
    {
      "name": "pkg:docker/<registry>/<image>@<tag/digest>?platform=<platform>",
      "digest": {
        "sha256": "e8275b2b76280af67e26f068e5d585eb905f8dfd2f1918b3229db98133cb4862"
      }
    }
  ],
  "predicate": {
    "buildDefinition": {
      "buildType": "https://github.com/moby/buildkit/blob/master/docs/attestations/slsa-definitions.md",
      "externalParameters": {
        "configSource": {
          "path": "Dockerfile"
        },
        "request": {
          "frontend": "dockerfile.v0",
          "args": {},
          "locals": [
            {
              "name": "context"
            },
            {
              "name": "dockerfile"
            }
          ]
        }
      },
      "internalParameters": {
        "builderPlatform": "linux/amd64"
      },
      "resolvedDependencies": [
        {
          "uri": "pkg:docker/alpine@latest?platform=linux%2Famd64",
          "digest": {
            "sha256": "8914eb54f968791faf6a8638949e480fef81e697984fba772b3976835194c6d4"
          }
        }
      ]
    },
    "runDetails": {
      "builder": {
        "id": ""
      },
      "metadata": {
        "invocationId": "yirbp1aosi1vqjmi3z6bc75nb",
        "startedOn": "2022-12-08T11:48:59.466513707Z",
        "finishedOn": "2022-12-08T11:49:01.256820297Z",
        "buildkit_reproducible": false,
        "buildkit_completeness": {
          "request": false,
          "resolvedDependencies": false
        },
        "buildkit_metadata": {}
      }
    }
  }
}
```

### Max

The `max` mode includes all of the information included in the `min` mode, as
well as:

- The LLB definition of the build. These show the exact steps taken to produce
  the image.
- Information about the Dockerfile, including a full base64-encoded version of
  the file.
- Source maps describing the relationship between build steps and image layers.

When possible, you should prefer `mode=max` as it contains significantly more
detailed information for analysis.

> [!WARNING]
>
> Note that `mode=max` exposes the values of
> [build arguments](/reference/cli/docker/buildx/build/#build-arg).
>
> If you're misusing build arguments to pass credentials, authentication
> tokens, or other secrets, you should refactor your build to pass the secrets using
> [secret mounts](/reference/cli/docker/buildx/build/#secret) instead.
> Secret mounts don't leak outside of the build and are never included in provenance attestations.

## Version

The `version` parameter lets you specify which SLSA provenance schema version
to use. Supported values are `version=v1` (default) and `version=v0.2`.

To generate the older SLSA v0.2 format:

```console
$ docker buildx build --tag <namespace>/<image>:<version> \
    --attest type=provenance,mode=max,version=v0.2 .
```

For more information about SLSA Provenance v1, see the
[SLSA specification](https://slsa.dev/spec/v1.1/provenance). To see the
difference between SLSA v0.2 and v1 provenance attestations, refer to
[SLSA definitions](./slsa-definitions.md)

## Inspecting Provenance

To explore created Provenance exported through the `image` exporter, you can
use [`imagetools inspect`](/reference/cli/docker/buildx/imagetools/inspect/).

Using the `--format` option, you can specify a template for the output. All
provenance-related data is available under the `.Provenance` attribute. For
example, to get the raw contents of the Provenance in the SLSA format:

```console
$ docker buildx imagetools inspect <namespace>/<image>:<version> \
    --format "{{ json .Provenance.SLSA }}"
{
  "buildDefinition": {
    "buildType": "https://github.com/moby/buildkit/blob/master/docs/attestations/slsa-definitions.md",
    ...
  },
  ...
}
```

You can also construct more complex expressions using the full functionality of
Go templates. For example, for SLSA v1 provenance generated with `mode=max`, you
can extract the full source code of the Dockerfile used to build the image:

```console
$ docker buildx imagetools inspect <namespace>/<image>:<version> \
    --format '{{ range .Provenance.SLSA.runDetails.metadata.buildkit_metadata.source.infos }}{{ if eq .filename "Dockerfile" }}{{ .data }}{{ end }}{{ end }}' | base64 -d
FROM alpine:latest
```

## Provenance attestation example

The following example shows what a JSON representation of a provenance
attestation with `mode=max` and SLSA v1 looks like for `FROM alpine:latest`:

```json
{
  "_type": "https://in-toto.io/Statement/v1",
  "predicateType": "https://slsa.dev/provenance/v1",
  "subject": [
    {
      "name": "pkg:docker/<registry>/<image>@<tag/digest>?platform=<platform>",
      "digest": {
        "sha256": "e8275b2b76280af67e26f068e5d585eb905f8dfd2f1918b3229db98133cb4862"
      }
    }
  ],
  "predicate": {
    "buildDefinition": {
      "buildType": "https://github.com/moby/buildkit/blob/master/docs/attestations/slsa-definitions.md",
      "externalParameters": {
        "configSource": {
          "path": "Dockerfile"
        },
        "request": {
          "frontend": "dockerfile.v0",
          "args": {},
          "locals": [
            {
              "name": "context"
            },
            {
              "name": "dockerfile"
            }
          ]
        }
      },
      "internalParameters": {
        "builderPlatform": "linux/amd64",
        "buildConfig": {
          "llbDefinition": [
            {
              "id": "step0"
            },
            {
              "id": "step1",
              "inputs": [
                "step0:0"
              ]
            }
          ]
        }
      },
      "resolvedDependencies": [
        {
          "uri": "pkg:docker/alpine@latest?platform=linux%2Famd64",
          "digest": {
            "sha256": "8914eb54f968791faf6a8638949e480fef81e697984fba772b3976835194c6d4"
          }
        }
      ]
    },
    "runDetails": {
      "builder": {
        "id": ""
      },
      "metadata": {
        "invocationId": "46ue2x93k3xj5l463dektwldw",
        "startedOn": "2022-12-08T11:50:54.953375437Z",
        "finishedOn": "2022-12-08T11:50:55.447841328Z",
        "buildkit_reproducible": false,
        "buildkit_completeness": {
          "request": true,
          "resolvedDependencies": false
        },
        "buildkit_metadata": {
          "source": {
            "infos": [
              {
                "filename": "Dockerfile",
                "data": "RlJPTSBhbHBpbmU6bGF0ZXN0Cg=="
              }
            ]
          },
          "layers": {
            "step0:0": [
              [
                {
                  "mediaType": "application/vnd.oci.image.layer.v1.tar+gzip",
                  "digest": "sha256:c158987b05517b6f2c5913f3acef1f2182a32345a304fe357e3ace5fadcad715",
                  "size": 3370706
                }
              ]
            ]
          }
        }
      }
    }
  }
}
```
