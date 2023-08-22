# docker scout quickview

<!---MARKER_GEN_START-->
Quick overview of an image

### Aliases

`docker scout quickview`, `docker scout qv`

### Options

| Name             | Type     | Default | Description                                                                                                     |
|:-----------------|:---------|:--------|:----------------------------------------------------------------------------------------------------------------|
| `-o`, `--output` | `string` |         | Write the report to a file.                                                                                     |
| `--platform`     | `string` |         | Platform of image to analyze                                                                                    |
| `--ref`          | `string` |         | Reference to use if the provided tarball contains multiple references.<br>Can only be used with --type archive. |
| `--type`         | `string` | `image` | Type of the image to analyze. Can be one of:<br>- image<br>- oci-dir<br>- archive (docker save tarball)<br>     |


<!---MARKER_GEN_END-->

## Description

The `docker scout quickview` command displays a quick overview of an image.
It displays a summary of the vulnerabilities in the image and the vulnerabilities from the base image.
If available it also displays base image refresh and update recommendations.

If no image is specified, the most recently built image will be used.

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

  │ Know more about vulnerabilities:
  │    docker scout cves golang:1.19.4
  │ Know more about base image update recommendations:
  │    docker scout recommendations golang:1.19.4
```

### Quick overview of the most recently built image

```console
$ docker scout qv
```
