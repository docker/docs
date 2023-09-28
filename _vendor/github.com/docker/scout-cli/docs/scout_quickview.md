# docker scout quickview

<!---MARKER_GEN_START-->
Quick overview of an image

### Aliases

`docker scout quickview`, `docker scout qv`

### Options

| Name             | Type     | Default | Description                                                                                              |
|:-----------------|:---------|:--------|:---------------------------------------------------------------------------------------------------------|
| `--env`          | `string` |         | Name of the environment                                                                                  |
| `--org`          | `string` |         | Namespace of the Docker organization                                                                     |
| `-o`, `--output` | `string` |         | Write the report to a file.                                                                              |
| `--platform`     | `string` |         | Platform of image to analyze                                                                             |
| `--ref`          | `string` |         | Reference to use if the provided tarball contains multiple references.<br>Can only be used with archive. |


<!---MARKER_GEN_END-->

## Description

The `docker scout quickview` command displays a quick overview of an image.
It displays a summary of the vulnerabilities in the image and the vulnerabilities from the base image.
If available it also displays base image refresh and update recommendations.

If no image is specified, the most recently built image will be used.

The following artifact types are supported:

- Images
- OCI layout directories
- Tarball archives, as created by `docker save`
- Local directory or file

The tool analyzes the provided software artifact, and generates a vulnerability report.

By default, the tool expects an image reference, such as:

- `redis`
- `curlimages/curl:7.87.0`
- `mcr.microsoft.com/dotnet/runtime:7.0`

If the artifact you want to analyze is an OCI directory, a tarball archive, a local file or directory,
or if you want to control from where the image will be resolved, you must prefix the reference with one of the following:

- `image://` (default) use a local image, or fall back to a registry lookup
- `local://` use an image from the local image store (don't do a registry lookup)
- `registry://` use an image from a registry (don't use a local image)
- `oci-dir://` use an OCI layout directory
- `archive://` use a tarball archive, as created by docker save
- `fs://` use a local directory or file

## Examples

### Quick overview of an image

```console
$ docker scout quickview golang:1.19.4
    ...Pulling
    ✓ Pulled
    ✓ SBOM of image already cached, 278 packages indexed

  Your image  golang:1.19.4                          │    5C     3H     6M    63L
  Base image  buildpack-deps:bullseye-scm            │    5C     1H     3M    48L     6?
  Refreshed base image  buildpack-deps:bullseye-scm  │    0C     0H     0M    42L
                                                     │    -5     -1     -3     -6     -6
  Updated base image  buildpack-deps:sid-scm         │    0C     0H     1M    29L
                                                     │    -5     -1     -2    -19     -6
```

### Quick overview of the most recently built image

```console
$ docker scout qv
```
