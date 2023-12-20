---
title: Provenance attestations
keywords: build, attestations, provenance, slsa, git, metadata
description: >
  Provenance build attestations describe how and where your image was built.
---

The provenance attestations include facts about the build process, including
details such as:

- Build timestamps
- Build parameters and environment
- Version control metadata
- Source code details
- Materials (files, scripts) consumed during the build

Provenance attestations follow the
[SLSA provenance schema, version 0.2](https://slsa.dev/provenance/v0.2#schema).

For more information about how BuildKit populates these provenance properties, refer to
[SLSA definitions](slsa-definitions.md).

## Create provenance attestations

To create a provenance attestation, pass the `--attest type=provenance` option
to the `docker buildx build` command:

```console
$ docker buildx build --tag <namespace>/<image>:<version> \
    --attest type=provenance,mode=[min,max] .
```

Alternatively, you can use the shorthand `--provenance=true` option instead of `--attest type=provenance`.
To specify the `mode` parameter using the shorthand option, use: `--provenance=mode=max`.

For an example on how to add provenance attestations with GitHub Actions, see
[Add attestations with GitHub Actions](../ci/github-actions/attestations.md).

## Mode

You can use the `mode` parameter to define the level of detail to be included in
the provenance attestation. Supported values are `mode=min`, and `mode=max`
(default).

### Min

In `min` mode, the provenance attestations include a minimal set of information,
such as:

- Build timestamps
- The frontend used
- Build materials
- Source repository and revision
- Build platform
- Reproducibility

Values of build arguments, the identities of secrets, and rich layer metadata is
not included `mode=min`. The `min`-level provenance is safe to use for all
builds, as it doesn't leak information from any part of the build environment.

The following JSON example shows the information included in a provenance
attestations created using the `min` mode:

```json
{
  "_type": "https://in-toto.io/Statement/v0.1",
  "predicateType": "https://slsa.dev/provenance/v0.2",
  "subject": [
    {
      "name": "pkg:docker/<registry>/<image>@<tag/digest>?platform=<platform>",
      "digest": {
        "sha256": "e8275b2b76280af67e26f068e5d585eb905f8dfd2f1918b3229db98133cb4862"
      }
    }
  ],
  "predicate": {
    "builder": { "id": "" },
    "buildType": "https://mobyproject.org/buildkit@v1",
    "materials": [
      {
        "uri": "pkg:docker/docker/dockerfile@1",
        "digest": {
          "sha256": "9ba7531bd80fb0a858632727cf7a112fbfd19b17e94c4e84ced81e24ef1a0dbc"
        }
      },
      {
        "uri": "pkg:docker/golang@1.19.4-alpine?platform=linux%2Farm64",
        "digest": {
          "sha256": "a9b24b67dc83b3383d22a14941c2b2b2ca6a103d805cac6820fd1355943beaf1"
        }
      }
    ],
    "invocation": {
      "configSource": { "entryPoint": "Dockerfile" },
      "parameters": {
        "frontend": "gateway.v0",
        "args": {
          "cmdline": "docker/dockerfile:1",
          "source": "docker/dockerfile:1",
          "target": "binaries"
        },
        "locals": [{ "name": "context" }, { "name": "dockerfile" }]
      },
      "environment": { "platform": "linux/arm64" }
    },
    "metadata": {
      "buildInvocationID": "c4a87v0sxhliuewig10gnsb6v",
      "buildStartedOn": "2022-12-16T08:26:28.651359794Z",
      "buildFinishedOn": "2022-12-16T08:26:29.625483253Z",
      "reproducible": false,
      "completeness": {
        "parameters": true,
        "environment": true,
        "materials": false
      },
      "https://mobyproject.org/buildkit@v1#metadata": {
        "vcs": {
          "revision": "a9ba846486420e07d30db1107411ac3697ecab68",
          "source": "git@github.com:<org>/<repo>.git"
        }
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
detailed information for analysis. However, on some builds it may not be
appropriate, as it includes the values of
[build arguments](../../engine/reference/commandline/buildx_build.md#build-arg)
and metadata about secrets and SSH mounts. If you pass sensitive information
using build arguments, consider refactoring builds to pass secret values using
[build secrets](../../engine/reference/commandline/buildx_build.md#secret), to
prevent leaking of sensitive information.

## Inspecting Provenance

To explore created Provenance exported through the `image` exporter, you can
use [`imagetools inspect`](../../engine/reference/commandline/buildx_imagetools_inspect.md).

Using the `--format` option, you can specify a template for the output. All
provenance-related data is available under the `.Provenance` attribute. For
example, to get the raw contents of the Provenance in the SLSA format:

```console
$ docker buildx imagetools inspect <namespace>/<image>:<version> \
    --format "{{ json .Provenance.SLSA }}"
{
  "buildType": "https://mobyproject.org/buildkit@v1",
  ...
}
```

You can also construct more complex expressions using the full functionality of
Go templates. For example, for provenance generated with `mode=max`, you can
extract the full source code of the Dockerfile used to build the image:

```console
$ docker buildx imagetools inspect <namespace>/<image>:<version> \
    --format '{{ range (index .Provenance.SLSA.metadata "https://mobyproject.org/buildkit@v1#metadata").source.infos }}{{ if eq .filename "Dockerfile" }}{{ .data }}{{ end }}{{ end }}' | base64 -d
FROM ubuntu:20.04
RUN apt-get update
...
```

## Provenance attestation example

<!-- TODO: add a link to the definitions page, imported from moby/buildkit -->

The following example shows what a JSON representation of a provenance
attestation with `mode=max` looks like:

```json
{
  "_type": "https://in-toto.io/Statement/v0.1",
  "predicateType": "https://slsa.dev/provenance/v0.2",
  "subject": [
    {
      "name": "pkg:docker/<registry>/<image>@<tag/digest>?platform=<platform>",
      "digest": {
        "sha256": "e8275b2b76280af67e26f068e5d585eb905f8dfd2f1918b3229db98133cb4862"
      }
    }
  ],
  "predicate": {
    "builder": { "id": "" },
    "buildType": "https://mobyproject.org/buildkit@v1",
    "materials": [
      {
        "uri": "pkg:docker/docker/dockerfile@1",
        "digest": {
          "sha256": "9ba7531bd80fb0a858632727cf7a112fbfd19b17e94c4e84ced81e24ef1a0dbc"
        }
      },
      {
        "uri": "pkg:docker/golang@1.19.4-alpine?platform=linux%2Farm64",
        "digest": {
          "sha256": "a9b24b67dc83b3383d22a14941c2b2b2ca6a103d805cac6820fd1355943beaf1"
        }
      }
    ],
    "buildConfig": {
      "llbDefinition": [
        {
          "id": "step4",
          "op": {
            "Op": {
              "exec": {
                "meta": {
                  "args": ["/bin/sh", "-c", "go mod download -x"],
                  "env": [
                    "PATH=/go/bin:/usr/local/go/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin",
                    "GOLANG_VERSION=1.19.4",
                    "GOPATH=/go",
                    "CGO_ENABLED=0"
                  ],
                  "cwd": "/src"
                },
                "mounts": [
                  { "input": 0, "dest": "/", "output": 0 },
                  {
                    "input": -1,
                    "dest": "/go/pkg/mod",
                    "output": -1,
                    "mountType": 3,
                    "cacheOpt": { "ID": "//go/pkg/mod" }
                  },
                  {
                    "input": 1,
                    "selector": "/go.mod",
                    "dest": "/src/go.mod",
                    "output": -1,
                    "readonly": true
                  },
                  {
                    "input": 1,
                    "selector": "/go.sum",
                    "dest": "/src/go.sum",
                    "output": -1,
                    "readonly": true
                  }
                ]
              }
            },
            "platform": { "Architecture": "arm64", "OS": "linux" },
            "constraints": {}
          },
          "inputs": ["step3:0", "step1:0"]
        }
      ]
    },
    "metadata": {
      "buildInvocationID": "edf52vxjyf9b6o5qd7vgx0gru",
      "buildStartedOn": "2022-12-15T15:38:13.391980297Z",
      "buildFinishedOn": "2022-12-15T15:38:14.274565297Z",
      "reproducible": false,
      "completeness": {
        "parameters": true,
        "environment": true,
        "materials": false
      },
      "https://mobyproject.org/buildkit@v1#metadata": {
        "vcs": {
          "revision": "a9ba846486420e07d30db1107411ac3697ecab68-dirty",
          "source": "git@github.com:<org>/<repo>.git"
        },
        "source": {
          "locations": {
            "step4": {
              "locations": [
                {
                  "ranges": [
                    { "start": { "line": 5 }, "end": { "line": 5 } },
                    { "start": { "line": 6 }, "end": { "line": 6 } },
                    { "start": { "line": 7 }, "end": { "line": 7 } },
                    { "start": { "line": 8 }, "end": { "line": 8 } }
                  ]
                }
              ]
            }
          },
          "infos": [
            {
              "filename": "Dockerfile",
              "data": "RlJPTSBhbHBpbmU6bGF0ZXN0Cg==",
              "llbDefinition": [
                {
                  "id": "step0",
                  "op": {
                    "Op": {
                      "source": {
                        "identifier": "local://dockerfile",
                        "attrs": {
                          "local.differ": "none",
                          "local.followpaths": "[\"Dockerfile\",\"Dockerfile.dockerignore\",\"dockerfile\"]",
                          "local.session": "s4j58ngehdal1b5hn7msiqaqe",
                          "local.sharedkeyhint": "dockerfile"
                        }
                      }
                    },
                    "constraints": {}
                  }
                },
                { "id": "step1", "op": { "Op": null }, "inputs": ["step0:0"] }
              ]
            }
          ]
        },
        "layers": {
          "step2:0": [
            [
              {
                "mediaType": "application/vnd.docker.image.rootfs.diff.tar.gzip",
                "digest": "sha256:261da4162673b93e5c0e7700a3718d40bcc086dbf24b1ec9b54bca0b82300626",
                "size": 3259190
              },
              {
                "mediaType": "application/vnd.docker.image.rootfs.diff.tar.gzip",
                "digest": "sha256:bc729abf26b5aade3c4426d388b5ea6907fe357dec915ac323bb2fa592d6288f",
                "size": 286218
              },
              {
                "mediaType": "application/vnd.docker.image.rootfs.diff.tar.gzip",
                "digest": "sha256:7f1d6579712341e8062db43195deb2d84f63b0f2d1ed7c3d2074891085ea1b56",
                "size": 116878653
              },
              {
                "mediaType": "application/vnd.docker.image.rootfs.diff.tar.gzip",
                "digest": "sha256:652874aefa1343799c619d092ab9280b25f96d97939d5d796437e7288f5599c9",
                "size": 156
              }
            ]
          ]
        }
      }
    }
  }
}
```
