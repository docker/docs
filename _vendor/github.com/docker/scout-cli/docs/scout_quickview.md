# docker scout quickview

<!---MARKER_GEN_START-->
Quick overview of an image

### Aliases

`docker scout quickview`, `docker scout qv`

### Options

| Name                  | Type          | Default             | Description                                                                                             |
|:----------------------|:--------------|:--------------------|:--------------------------------------------------------------------------------------------------------|
| `--env`               | `string`      |                     | Name of the environment                                                                                 |
| `--ignore-suppressed` |               |                     | Filter CVEs found in Scout exceptions based on the specified exception scope                            |
| `--latest`            |               |                     | Latest indexed image                                                                                    |
| `--only-policy`       | `stringSlice` |                     | Comma separated list of policies to evaluate                                                            |
| `--only-vex-affected` |               |                     | Filter CVEs by VEX statements with status not affected                                                  |
| `--org`               | `string`      |                     | Namespace of the Docker organization                                                                    |
| `-o`, `--output`      | `string`      |                     | Write the report to a file                                                                              |
| `--platform`          | `string`      |                     | Platform of image to analyze                                                                            |
| `--ref`               | `string`      |                     | Reference to use if the provided tarball contains multiple references.<br>Can only be used with archive |
| `--vex-author`        | `stringSlice` | `[<.*@docker.com>]` | List of VEX statement authors to accept                                                                 |
| `--vex-location`      | `stringSlice` |                     | File location of directory or file containing VEX statements                                            |


<!---MARKER_GEN_END-->

## Description

The `docker scout quickview` command displays a quick overview of an image.
It displays a summary of the vulnerabilities in the specified image
and vulnerabilities from the base image.
If available, it also displays base image refresh and update recommendations.

If no image is specified, the most recently built image is used.

The following artifact types are supported:

- Images
- OCI layout directories
- Tarball archives, as created by `docker save`
- Local directory or file

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
- `archive://` use a tarball archive, as created by `docker save`
- `fs://` use a local directory or file
- `sbom://` SPDX file or in-toto attestation file with SPDX predicate or `syft` json SBOM file
    In case of `sbom://` prefix, if the file is not defined then it will try to read it from the standard input.

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

### Quick overview from an SPDX file

```console
$  syft -o spdx-json alpine:3.16.1 | docker scout quickview sbom://
 ✔ Loaded image                                                                                                                              alpine:3.16.1
 ✔ Parsed image                                                                    sha256:3d81c46cd8756ddb6db9ec36fa06a6fb71c287fb265232ba516739dc67a5f07d
 ✔ Cataloged contents                                                                     274a317d88b54f9e67799244a1250cad3fe7080f45249fa9167d1f871218d35f
   ├── ✔ Packages                        [14 packages]
   ├── ✔ File digests                    [75 files]
   ├── ✔ File metadata                   [75 locations]
   └── ✔ Executables                     [16 executables]

  Target   │ <stdin>        │    1C     2H     8M     0L
    digest │  274a317d88b5  │
```
