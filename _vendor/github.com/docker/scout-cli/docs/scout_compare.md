# docker scout compare

<!---MARKER_GEN_START-->
Compare two images and display differences (experimental)

### Aliases

`docker scout compare`, `docker scout diff`

### Options

| Name                  | Type          | Default             | Description                                                                                                                                                                    |
|:----------------------|:--------------|:--------------------|:-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `-x`, `--exit-on`     | `stringSlice` |                     |  Comma separated list of conditions to fail the action step if worse or changed, options are: vulnerability, policy, package                                                   |
| `--format`            | `string`      | `text`              | Output format of the generated vulnerability report:<br>- text: default output, plain text with or without colors depending on the terminal<br>- markdown: Markdown output<br> |
| `--hide-policies`     |               |                     | Hide policy status from the output                                                                                                                                             |
| `--ignore-base`       |               |                     | Filter out CVEs introduced from base image                                                                                                                                     |
| `--ignore-suppressed` |               |                     | Filter CVEs found in Scout exceptions based on the specified exception scope                                                                                                   |
| `--ignore-unchanged`  |               |                     | Filter out unchanged packages                                                                                                                                                  |
| `--multi-stage`       |               |                     | Show packages from multi-stage Docker builds                                                                                                                                   |
| `--only-fixed`        |               |                     | Filter to fixable CVEs                                                                                                                                                         |
| `--only-package-type` | `stringSlice` |                     | Comma separated list of package types (like apk, deb, rpm, npm, pypi, golang, etc)                                                                                             |
| `--only-policy`       | `stringSlice` |                     | Comma separated list of policies to evaluate                                                                                                                                   |
| `--only-severity`     | `stringSlice` |                     | Comma separated list of severities (critical, high, medium, low, unspecified) to filter CVEs by                                                                                |
| `--only-stage`        | `stringSlice` |                     | Comma separated list of multi-stage Docker build stage names                                                                                                                   |
| `--only-unfixed`      |               |                     | Filter to unfixed CVEs                                                                                                                                                         |
| `--only-vex-affected` |               |                     | Filter CVEs by VEX statements with status not affected                                                                                                                         |
| `--org`               | `string`      |                     | Namespace of the Docker organization                                                                                                                                           |
| `-o`, `--output`      | `string`      |                     | Write the report to a file                                                                                                                                                     |
| `--platform`          | `string`      |                     | Platform of image to analyze                                                                                                                                                   |
| `--ref`               | `string`      |                     | Reference to use if the provided tarball contains multiple references.<br>Can only be used with archive                                                                        |
| `--to`                | `string`      |                     | Image, directory, or archive to compare to                                                                                                                                     |
| `--to-env`            | `string`      |                     | Name of environment to compare to                                                                                                                                              |
| `--to-latest`         |               |                     | Latest image processed to compare to                                                                                                                                           |
| `--to-ref`            | `string`      |                     | Reference to use if the provided tarball contains multiple references.<br>Can only be used with archive.                                                                       |
| `--vex-author`        | `stringSlice` | `[<.*@docker.com>]` | List of VEX statement authors to accept                                                                                                                                        |
| `--vex-location`      | `stringSlice` |                     | File location of directory or file containing VEX statements                                                                                                                   |


<!---MARKER_GEN_END-->

## Description

The `docker scout compare` command analyzes two images and displays a comparison.

> This command is **experimental** and its behaviour might change in the future

The intended use of this command is to compare two versions of the same image.
For instance, when a new image is built and compared to the version running in production.

If no image is specified, the most recently built image is used
as a comparison target.

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

## Examples

### Compare the most recently built image to the latest tag

```console
$ docker scout compare --to namespace/repo:latest
```

### Compare local build to the same tag from the registry

```console
$ docker scout compare local://namespace/repo:latest --to registry://namespace/repo:latest
```

### Ignore base images

```console
$ docker scout compare --ignore-base --to namespace/repo:latest namespace/repo:v1.2.3-pre
```

### Generate a markdown output

```console
$ docker scout compare --format markdown --to namespace/repo:latest namespace/repo:v1.2.3-pre
```

### Only compare maven packages and only display critical vulnerabilities for maven packages

```console
$ docker scout compare --only-package-type maven --only-severity critical --to namespace/repo:latest namespace/repo:v1.2.3-pre
```

### Show all policy results for both images

```console
docker scout compare --to namespace/repo:latest namespace/repo:v1.2.3-pre
```
